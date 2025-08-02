// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network

import (
	"container/list"
	"fmt"
	"math"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

// MaxFlowDinic computes the maximum flow from source to target in a directed,
// weighted graph using Dinic’s algorithm. It repeatedly builds level graphs
// and augments blocking flows until no more augmenting paths exist.
func MaxFlowDinic(graph graph.WeightedDirected, source, target graph.Node) (float64, error) {
	if source.ID() == target.ID() {
		return 0, fmt.Errorf("source and target must be different")
	}
	parents := make([][]int64, graph.Nodes().Len())
	residualGraph, err := initializeResidualGraph(graph)
	if err != nil {
		return 0, fmt.Errorf("could not build residual graph: %v", err)
	}
	epsilon := 1.e-12
	var maxFlow = 0.0
	for canReachTargetInLevelGraph(residualGraph, source, target, parents) {
		flow := computeBlockingPath(residualGraph, source, target, parents)
		if flow < epsilon {
			break
		}
		maxFlow += flow
	}
	return maxFlow, nil
}

// initializeResidualGraph builds the residual graph for Dinic’s algorithm.
// It copies all nodes, adds directed edges with their original capacities,
// and initializes a zero flow map for each forward edge.
func initializeResidualGraph(originalGraph graph.WeightedDirected) (*residualGraph, error) {
	graphCopy := simple.NewWeightedDirectedGraph(0, 0)
	// flow tracks the current flow for each forward edge (initially zero).
	flow := make(map[edgeKey]float64)

	for nodes := originalGraph.Nodes(); nodes.Next(); {
		graphCopy.AddNode(nodes.Node())
	}

	nodes := originalGraph.Nodes()
	// For each edge current->neighbor in the original graph:
	// 1) add a forward edge with its capacity
	// 2) initialize its flow to 0.0
	for nodes.Next() {
		current := nodes.Node()
		for it := originalGraph.From(current.ID()); it.Next(); {
			neighbor := it.Node()

			capacity, ok := originalGraph.Weight(current.ID(), neighbor.ID())
			if !ok {
				panic("expected a weight for existing edge")
			}
			if capacity < 0.0 {
				return nil, fmt.Errorf("edge weights (capacities) can not be negative")
			}

			// Add forward edge with full capacity.
			forward := graphCopy.NewWeightedEdge(current, neighbor, capacity)
			graphCopy.SetWeightedEdge(forward)

			// Initialize flow for this edge to zero.
			flow[edgeKey{from: current.ID(), to: neighbor.ID()}] = 0.0
		}
	}
	return &residualGraph{
		Graph: graphCopy,
		Flow:  flow,
	}, nil
}

// canReachTargetInLevelGraph builds a level graph using BFS on residualGraph.
// It records, for each reachable node, the list of parents at the previous level.
// Returns true iff target is reachable from source via positive-capacity edges.
func canReachTargetInLevelGraph(residualGraph *residualGraph, source, target graph.Node, parents [][]int64) bool {
	// Reset parents slices in place.
	for i := range parents {
		parents[i] = parents[i][:0]
	}

	// levels[i] holds the BFS level of node i, or -1 if unvisited.
	levels := make([]int32, residualGraph.Graph.Nodes().Len())
	for i := range levels {
		levels[i] = -1
	}

	sourceID := source.ID()
	levels[sourceID] = 0

	queue := list.New()
	queue.PushBack(sourceID)

	for queue.Len() > 0 {
		parent := queue.Front()
		queue.Remove(parent)
		parentID := parent.Value.(int64)

		// Explore all outgoing edges with capacity > 0.
		for it := residualGraph.Graph.From(parentID); it.Next(); {
			childID := it.Node().ID()
			capacity, ok := residualGraph.Graph.Weight(parentID, childID)
			if !ok || capacity <= 0 {
				continue
			}
			// First time we visit childID: set its level and record parent.
			if levels[childID] == -1 {
				levels[childID] = levels[parentID] + 1
				parents[childID] = append(parents[childID], parentID)
				queue.PushBack(childID)
				// If we reach childID again at the same level, record an additional parent.
			} else if levels[childID] == levels[parentID]+1 {
				parents[childID] = append(parents[childID], parentID)
			}
		}
	}
	// The target node is reachable iff it was assigned a level ≥ 0.
	return levels[target.ID()] > -1
}

// computeBlockingPath finds and augments all blocking‐flow paths in the current
// level graph of Dinic’s algorithm. It backtracks from target to source using
// the parents slices, computes each path’s bottleneck capacity, updates both
// the residual capacities and the flow map, and returns the total flow added.
func computeBlockingPath(residualGraph *residualGraph, source, target graph.Node, parents [][]int64) float64 {
	var totalFlow = 0.0

	// path holds node IDs from target back to (eventually) source.
	path := []int64{target.ID()}
	var currentID = target.ID()

	for {
		var currentParentID int64

		// If there is a recorded parent, step “up” toward the source.
		if len(parents[currentID]) > 0 {
			currentParentID = parents[currentID][0]
			path = append(path, currentParentID)
		} else {
			// No further parent: backtrack.
			path = path[:len(path)-1]
			if len(path) == 0 {
				break
			}
			currentParentID = path[len(path)-1]
		}

		// When we’ve backtracked all the way to source:
		if currentParentID == source.ID() {
			// 1) Find the minimum residual capacity (bottleneck) along this path.
			bottleNeckOnPath := math.MaxFloat64
			for i := 0; i+1 < len(path); i++ {
				parentID := path[i+1]
				childID := path[i]
				weight, ok := residualGraph.Graph.Weight(parentID, childID)
				if !ok {
					panic("expected a weight for existing edge")
				}

				if weight < bottleNeckOnPath {
					bottleNeckOnPath = weight
				}
			}
			// 2) Augment flow: subtract bottleneck from forward capacities, and add to the flow map.
			for i := 0; i+1 < len(path); i++ {
				parentID := path[i+1]
				childID := path[i]
				currentCapacity, ok := residualGraph.Graph.Weight(parentID, childID)
				if !ok {
					panic("expected a weight for existing edge")
				}
				parent := residualGraph.Graph.Node(parentID)
				child := residualGraph.Graph.Node(childID)
				newCapacity := residualGraph.Graph.NewWeightedEdge(parent, child, currentCapacity-bottleNeckOnPath)
				residualGraph.Graph.SetWeightedEdge(newCapacity)
				edgeID := edgeKey{from: parentID, to: childID}
				currentFlow, ok := residualGraph.Flow[edgeID]
				if ok {
					residualGraph.Flow[edgeID] = currentFlow + bottleNeckOnPath
				} else {
					residualGraph.Flow[edgeID] = bottleNeckOnPath
				}
			}
			totalFlow += bottleNeckOnPath
			path = []int64{target.ID()}
		}
		currentID = currentParentID
	}
	return totalFlow
}

type residualGraph struct {
	Graph *simple.WeightedDirectedGraph
	Flow  map[edgeKey]float64
}

type edgeKey struct{ from, to int64 }
