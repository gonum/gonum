package network

import (
	"container/list"
	"fmt"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"math"
)

func initializeResidualGraph(graph graph.WeightedDirected) *simple.WeightedDirectedGraph {
	residualGraph := simple.NewWeightedDirectedGraph(0, 0)

	// Add all nodes
	for nodes := graph.Nodes(); nodes.Next(); {
		residualGraph.AddNode(nodes.Node())
	}

	// For each node u :
	for nodes := graph.Nodes(); nodes.Next(); {
		u := nodes.Node()
		// Iterate over all children of u
		for it := graph.From(u.ID()); it.Next(); {
			v := it.Node()
			// get the weight/capacity
			capacity, ok := graph.Weight(u.ID(), v.ID())
			if !ok {
				panic("expected a weight for existing edge")
			}
			// add forward edge to residualGraph (capacity)
			forward := residualGraph.NewWeightedEdge(u, v, capacity)
			residualGraph.SetWeightedEdge(forward)
			// add reverse edge v->u with zero weight (flow)
			reverse := residualGraph.NewWeightedEdge(v, u, 0)
			residualGraph.SetWeightedEdge(reverse)
		}
	}

	return residualGraph
}

func computeBlockingPath(graph *simple.WeightedDirectedGraph, source, target graph.Node, parents [][]int64) float64 {
	var totalFlow = 0.0
	var path []int64
	targetID := target.ID()
	path = append(path, targetID)
	var uID = targetID
	for {
		var vID int64
		if len(parents[uID]) > 0 {
			vID = parents[uID][0]
			path = append(path, vID)
		} else {
			path = path[:len(path)-1]
			if len(path) == 0 {
				break
			}
			vID = path[len(path)-1]
		}
		// path has been build from target to source, so the parent is on i+1 position of the ith child
		if vID == source.ID() {
			bottleNeckObPath := math.MaxFloat64
			// determine minimal flow on path
			for i := 0; i+1 < len(path); i++ {
				parentID := path[i+1]
				childID := path[i]
				weight, ok := graph.Weight(parentID, childID)
				if !ok {
					panic("expected a weight for existing edge")
				}

				if weight < bottleNeckObPath {
					bottleNeckObPath = weight
				}
			}
			// update the capacities and flows in the other edges
			for i := 0; i+1 < len(path); i++ {
				parentID := path[i+1]
				childID := path[i]
				currentCapacity, ok := graph.Weight(parentID, childID)
				if !ok {
					panic("expected a weight for existing edge")
				}
				parent := graph.Node(parentID)
				child := graph.Node(childID)
				newCapacity := graph.NewWeightedEdge(parent, child, currentCapacity-bottleNeckObPath)
				graph.SetWeightedEdge(newCapacity)
				if graph.HasEdgeFromTo(childID, parentID) {
					reverseCapacity, ok := graph.Weight(childID, parentID)
					if !ok {
						panic("expected a weight for existing edge")
					}
					newReverseCapacity := graph.NewWeightedEdge(child, parent, reverseCapacity+bottleNeckObPath)
					graph.SetWeightedEdge(newReverseCapacity)
				} else {
					newReverseCapacity := graph.NewWeightedEdge(child, parent, bottleNeckObPath)
					graph.SetWeightedEdge(newReverseCapacity)
				}
			}
			totalFlow += bottleNeckObPath
			path = path[:0]
			path = append(path, targetID)
		}
		uID = vID
	}
	return totalFlow
}

func canReachTargetInLevelGraph(graph graph.WeightedDirected, source, target graph.Node, parents [][]int64) bool {
	for i := range parents {
		parents[i] = parents[i][:0]
	}
	levels := make([]int32, graph.Nodes().Len())
	for i := range levels {
		levels[i] = -1
	}
	sourceID := source.ID()
	queue := list.New()
	levels[sourceID] = 0
	queue.PushBack(sourceID)
	for queue.Len() > 0 {
		parent := queue.Front()
		queue.Remove(parent)
		parentID := parent.Value.(int64)
		for it := graph.From(parentID); it.Next(); {
			childID := it.Node().ID()
			if capacity, ok := graph.Weight(parentID, childID); ok && capacity > 0 {
				if levels[childID] == -1 {
					levels[childID] = levels[parentID] + 1
					parents[childID] = append(parents[childID], parentID)
					queue.PushBack(childID)
				} else if levels[childID] == levels[parentID]+1 {
					parents[childID] = append(parents[childID], parentID)
				}
			}
		}
	}
	return levels[target.ID()] > -1
}

func MaxFlowDinic(graph graph.WeightedDirected, source, target graph.Node) (float64, error) {
	if source.ID() == target.ID() {
		return 0, fmt.Errorf("source and target must be different")
	}
	parents := make([][]int64, graph.Nodes().Len())
	residualGraph := initializeResidualGraph(graph)
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
