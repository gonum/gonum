package network

import (
	"fmt"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

func initializeResidualGraph(g graph.WeightedDirected) *simple.WeightedDirectedGraph {
	residualGraph := simple.NewWeightedDirectedGraph(0, 0)

	// 1) Add all nodes
	for nodes := g.Nodes(); nodes.Next(); {
		residualGraph.AddNode(nodes.Node())
	}

	// 2) For each u -> v in g:
	for nodes := g.Nodes(); nodes.Next(); {
		u := nodes.Node()
		// Iterate over all parents of u
		for it := g.From(u.ID()); it.Next(); {
			v := it.Node()
			// get the weight/capacity
			capacity, ok := g.Weight(u.ID(), v.ID())
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
