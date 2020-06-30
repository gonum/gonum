package flownet

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/traverse"
)

type ResidualEdge interface {
	graph.Edge
	MaxCapacity() int32
	CurrentFlow() int32
	IsReverseEdge() bool
	WithFlow(int32) ResidualEdge
	ReversedResidualEdge() ResidualEdge
}

type Graph interface {
	traverse.Graph
	ResidualGraph() ResidualGraph
}

type ResidualGraph interface {
	traverse.Graph
	SetFlow(uid, vid int64, flow int32)
	ResidualEdge(uid, vid int64) ResidualEdge
}

func AugmentingPath(g ResidualGraph, source, sink graph.Node) []ResidualEdge {
	parentOf := make(map[int64]int64)
	bfs := traverse.BreadthFirst{
		Traverse: func(e graph.Edge) bool {
			edge := g.ResidualEdge(e.From().ID(), e.To().ID())
			if edge.MaxCapacity()-edge.CurrentFlow() > 0 {
				if _, ok := parentOf[e.To().ID()]; !ok {
					parentOf[e.To().ID()] = e.From().ID()
				}
				return true
			}
			return false
		},
	}
	bfs.Walk(g, source, func(n graph.Node, d int) bool {
		return n.ID() == sink.ID()
	},
	)

	return makePath(g, parentOf, sink)
}

func makePath(g ResidualGraph, parentOf map[int64]int64, sink graph.Node) []ResidualEdge {
	child := sink.ID()
	if _, ok := parentOf[child]; !ok {
		return nil
	}
	path := []ResidualEdge{}

	for {
		if _, ok := parentOf[child]; ok {
			parent := parentOf[child]
			path = append(path, g.ResidualEdge(parent, child))
			child = parent
		} else {
			return path
		}
	}
}

func AvailableCapacity(edges []ResidualEdge) int32 {
	min := edges[0].MaxCapacity() - edges[0].CurrentFlow()
	for _, edge := range edges {
		capacity := edge.MaxCapacity() - edge.CurrentFlow()
		if capacity < min {
			min = capacity
		}
	}
	return min
}

func MaxFlow(g Graph, source, sink graph.Node) int32 {
	rg := g.ResidualGraph()
	for {
		path := AugmentingPath(rg, source, sink)
		if path == nil {
			break
		}
		flow := AvailableCapacity(path)
		for _, edge := range path {
			newFlow := edge.CurrentFlow() + flow
			rg.SetFlow(
				edge.From().ID(),
				edge.To().ID(),
				newFlow,
			)
		}
	}

	var flow int32
	to := rg.From(source.ID())
	for to.Next() {
		flow += rg.ResidualEdge(source.ID(), to.Node().ID()).CurrentFlow()
	}

	return flow
}
