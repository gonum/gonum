// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network

import (
	"math"

	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/linear"
	"gonum.org/v1/gonum/graph/simple"
)

// MaxFlowDinic computes the maximum flow from source to target in a directed,
// weighted graph using [Dinic's algorithm]. It repeatedly builds level graphs
// and augments blocking flows until no more augmenting paths exist.
//
// MaxFlowDinic will panic if s and t are the same node or g has any
// reachable negative edge weight.
//
// The eps parameter specifies an absolute tolerance for treating tiny flow
// updates as zero. If eps is negative a default of 1e-12 is used.
//
// [Dinic's algorithm]: https://en.wikipedia.org/wiki/Dinic%27s_algorithm
func MaxFlowDinic(g graph.WeightedDirected, s, t graph.Node, eps float64) float64 {
	if s.ID() == t.ID() {
		panic("no cut between s and t")
	}
	parents := make([][]int64, g.Nodes().Len())
	r := initializeResidualGraph(g)

	if eps < 0 {
		eps = 1e-12
	}
	var maxFlow float64
	for canReachTargetInLevelGraph(r, s, t, parents) {
		flow := computeBlockingPath(r, s, t, parents)
		if scalar.EqualWithinAbs(flow, 0, eps) {
			break
		}
		maxFlow += flow
	}
	return maxFlow
}

// initializeResidualGraph builds the residual graph for Dinic’s algorithm.
// It copies all nodes and, for each original edge u→v:
//   - adds a forward edge u→v with capacity equal to the original capacity
//   - adds a reverse edge v→u with capacity 0 if one doesn’t already exist
func initializeResidualGraph(g graph.WeightedDirected) *simple.WeightedDirectedGraph {
	r := simple.NewWeightedDirectedGraph(0, 0)

	nodes := g.Nodes()
	for nodes.Next() {
		r.AddNode(nodes.Node())
	}

	nodes.Reset()
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
				panic("negative edge weight")
			}

			// Add forward edge with full capacity.
			forward := r.NewWeightedEdge(u, v, capacity)
			r.SetWeightedEdge(forward)
			// Add reverse edge if it does not exist
			if _, ok := r.Weight(v.ID(), u.ID()); !ok {
				r.SetWeightedEdge(r.NewWeightedEdge(v, u, 0))
			}
		}
	}
	return r
}

// canReachTargetInLevelGraph builds a level graph using BFS on residualGraph.
// It records, for each reachable node, the list of parents at the previous level.
// It returns whether target is reachable from source via positive-capacity edges.
func canReachTargetInLevelGraph(r *simple.WeightedDirectedGraph, s, t graph.Node, parents [][]int64) bool {
	// Reset parents slices in place.
	for i := range parents {
		parents[i] = parents[i][:0]
	}

	// levels[i] holds the BFS level of node i, or -1 if unvisited.
	levels := make([]int32, r.Nodes().Len())
	for i := range levels {
		levels[i] = -1
	}
	levels[s.ID()] = 0

	var queue linear.NodeQueue
	queue.Enqueue(s)

	for queue.Len() > 0 {
		pid := queue.Dequeue().ID()

		// Explore all outgoing edges with capacity > 0.
		for it := r.From(pid); it.Next(); {
			cid := it.Node().ID()
			capacity, ok := r.Weight(pid, cid)
			if !ok || capacity <= 0 {
				continue
			}
			// First time we visit cid: set its level and record p.
			if levels[cid] == -1 {
				levels[cid] = levels[pid] + 1
				parents[cid] = append(parents[cid], pid)
				queue.Enqueue(it.Node())
				// If we reach cid again at the same level, record an additional p.
			} else if levels[cid] == levels[pid]+1 {
				parents[cid] = append(parents[cid], pid)
			}
		}
	}
	// The t node is reachable iff it was assigned a level ≥ 0.
	return levels[t.ID()] >= 0
}

// computeBlockingPath finds and augments all blocking‐flow paths in the current
// level graph of Dinic’s algorithm. It backtracks from target to source using
// the parents slices, computes each path’s bottleneck capacity, updates both
// the residual capacities and the flow map, and returns the total flow added.
func computeBlockingPath(r *simple.WeightedDirectedGraph, s, t graph.Node, parents [][]int64) float64 {
	var totalFlow float64

	// path holds node IDs from t back to (eventually) s.
	path := []int64{t.ID()}
	uid := t.ID()

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
			bottleNeckOnPath := math.Inf(0)
			for i := 0; i+1 < len(path); i++ {
				pid := path[i+1]
				cid := path[i]
				w, ok := r.Weight(pid, cid)
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
				currentCapacity, ok := r.Weight(pid, cid)
				if !ok {
					panic("expected a weight for existing edge")
				}
				parent := r.Node(pid)
				child := r.Node(cid)
				forwardCapacity := r.NewWeightedEdge(parent, child, currentCapacity-bottleNeckOnPath)
				r.SetWeightedEdge(forwardCapacity)
				reverseCapacity, ok := r.Weight(cid, pid)
				if !ok {
					panic("expected reverse residual edge")
				}
				r.SetWeightedEdge(r.NewWeightedEdge(r.Node(cid), r.Node(pid), reverseCapacity+bottleNeckOnPath))
			}
			totalFlow += bottleNeckOnPath
			path = []int64{t.ID()}
		}
		uid = vid
	}
	return totalFlow
}
