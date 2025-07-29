package network

import (
	"container/list"
	"fmt"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
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


func canReachTargetInLevelGraph(graph graph.WeightedDirected, source, target graph.Node, parents []int32) bool {
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
		parentID := parent.Value.(int64)
		queue.Remove(parent)
		for it := graph.From(parentID); it.Next(); {
			child := it.Node()
			childID := child.ID()
			if capacity, ok := graph.Weight(parentID, childID); ok && capacity > 0 {
				if levels[childID] == -1 {
					levels[childID] = levels[parentID] + 1
					parents[childID]
				}
			}
		}
	}


}


func Dinic(graph graph.WeightedDirected, source, target graph.Node) (float64, error {
	if source.ID() == target.ID() {
		return 0, fmt.Errorf("source and target must be different")
	}
	parents := make([]int32, graph.Nodes().Len())
	for i := range parents {
		parents[i] = -1
	}
	residualGraph := initializeResidualGraph(graph)

}
