// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network

import (
	"container/list"
	"fmt"
	"math"

	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

// MaxFlowDinic computes the maximum flow from source to target in a directed,
// weighted graph using [Dinic's algorithm]. It repeatedly builds level graphs
// and augments blocking flows until no more augmenting paths exist.
//
// [Dinic's algorithm]: https://en.wikipedia.org/wiki/Dinic%27s_algorithm
func MaxFlowDinic(g graph.WeightedDirected, s, t graph.Node) float64 {
	if s.ID() == t.ID() {
		return 0, fmt.Errorf("s and t must be different")
	}
	parents := make([][]int64, graph.Nodes().Len())
	r, err := initializeResidualGraph(graph)
	if err != nil {
		return 0, fmt.Errorf("could not build residual graph: %v", err)
	}
	const epsilon = 1.e-12
	var maxFlow = 0.0
	for canReachTargetInLevelGraph(r, s, t, parents) {
		flow := computeBlockingPath(r, s, t, parents)
		if scalar.EqualWithinAbs(flow, 0, epsilon) {
			break
		}
		maxFlow += flow
	}
	return maxFlow, nil
}

// initializeResidualGraph builds the residual graph for Dinic’s algorithm.
// It copies all nodes, adds directed edges with their original capacities,
// and initializes a zero flow map for each forward edge.
func initializeResidualGraph(g graph.WeightedDirected) (*residualGraph, error) {
	r := &residualGraph{
		g: simple.NewWeightedDirectedGraph(0, 0),
		// flow tracks the current flow for each forward edge (initially zero).
		flow: make(map[edgeKey]float64),
	}

	nodes := g.Nodes()
	for nodes.Next() {
		r.g.AddNode(nodes.Node())
	}

	nodes = g.Nodes()
	// For each edge current->neighbor in the original graph:
	// 1) add a forward edge with its capacity
	// 2) initialize its flow to 0.0
	for nodes.Next() {
		u := nodes.Node()
		it := g.From(u.ID())
		for it.Next() {
			v := it.Node()

			capacity, ok := g.Weight(u.ID(), v.ID())
			if !ok {
				panic("expected a weight for existing edge")
			}
			if capacity < 0 {
				return nil, fmt.Errorf("edge weights (capacities) can not be negative")
			}

			// Add forward edge with full capacity.
			forward := r.g.NewWeightedEdge(u, v, capacity)
			r.g.SetWeightedEdge(forward)

			// Initialize flow for this edge to zero.
			r.flow[edgeKey{from: u.ID(), to: v.ID()}] = 0.0
		}
	}
	return r, nil
}

// canReachTargetInLevelGraph builds a level graph using BFS on residualGraph.
// It records, for each reachable node, the list of parents at the previous level.
// Returns true iff target is reachable from source via positive-capacity edges.
func canReachTargetInLevelGraph(residualGraph *residualGraph, s, t graph.Node, parents [][]int64) bool {
	// Reset parents slices in place.
	for i := range parents {
		parents[i] = parents[i][:0]
	}

	// levels[i] holds the BFS level of node i, or -1 if unvisited.
	levels := make([]int32, residualGraph.g.Nodes().Len())
	for i := range levels {
		levels[i] = -1
	}

	sID := s.ID()
	levels[sID] = 0

	queue := list.New()
	queue.PushBack(sID)

	for queue.Len() > 0 {
		p := queue.Front()
		queue.Remove(p)
		pid := p.Value.(int64)

		// Explore all outgoing edges with capacity > 0.
		for it := residualGraph.g.From(pid); it.Next(); {
			cid := it.Node().ID()
			capacity, ok := residualGraph.g.Weight(pid, cid)
			if !ok || capacity <= 0 {
				continue
			}
			// First time we visit cid: set its level and record p.
			if levels[cid] == -1 {
				levels[cid] = levels[pid] + 1
				parents[cid] = append(parents[cid], pid)
				queue.PushBack(cid)
				// If we reach cid again at the same level, record an additional p.
			} else if levels[cid] == levels[pid]+1 {
				parents[cid] = append(parents[cid], pid)
			}
		}
	}
	// The t node is reachable iff it was assigned a level ≥ 0.
	return levels[t.ID()] > -1
}

// computeBlockingPath finds and augments all blocking‐flow paths in the current
// level graph of Dinic’s algorithm. It backtracks from target to source using
// the parents slices, computes each path’s bottleneck capacity, updates both
// the residual capacities and the flow map, and returns the total flow added.
func computeBlockingPath(r *residualGraph, s, t graph.Node, parents [][]int64) float64 {
	var totalFlow = 0.0

	// path holds node IDs from t back to (eventually) s.
	path := []int64{t.ID()}
	var uid = t.ID()

	for {
		var vid int64

		// If there is a recorded parent, step “up” toward the s.
		if len(parents[uid]) > 0 {
			vid = parents[uid][0]
			path = append(path, vid)
		} else {
			// No further parent: backtrack.
			path = path[:len(path)-1]
			if len(path) == 0 {
				break
			}
			vid = path[len(path)-1]
		}

		// When we’ve backtracked all the way to s:
		if vid == s.ID() {
			// 1) Find the minimum residual capacity (bottleneck) along this path.
			bottleNeckOnPath := math.MaxFloat64
			for i := 0; i+1 < len(path); i++ {
				pid := path[i+1]
				cid := path[i]
				w, ok := r.g.Weight(pid, cid)
				if !ok {
					panic("expected a w for existing edge")
				}

				if w < bottleNeckOnPath {
					bottleNeckOnPath = w
				}
			}
			// 2) Augment flow: subtract bottleneck from forward capacities, and add to the flow map.
			for i := 0; i+1 < len(path); i++ {
				pid := path[i+1]
				cid := path[i]
				currentCapacity, ok := r.g.Weight(pid, cid)
				if !ok {
					panic("expected a weight for existing edge")
				}
				parent := r.g.Node(pid)
				child := r.g.Node(cid)
				newCapacity := r.g.NewWeightedEdge(parent, child, currentCapacity-bottleNeckOnPath)
				r.g.SetWeightedEdge(newCapacity)
				edgeID := edgeKey{from: pid, to: cid}
				currentFlow, ok := r.flow[edgeID]
				if ok {
					r.flow[edgeID] = currentFlow + bottleNeckOnPath
				} else {
					r.flow[edgeID] = bottleNeckOnPath
				}
			}
			totalFlow += bottleNeckOnPath
			path = []int64{t.ID()}
		}
		uid = vid
	}
	return totalFlow
}

type residualGraph struct {
	g    *simple.WeightedDirectedGraph
	flow map[edgeKey]float64
}

type edgeKey struct{ from, to int64 }
