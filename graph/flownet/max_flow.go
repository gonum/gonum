package flownet

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/traverse"
)

// ResidualEdge represents a residual edge for Edmonds-Karp (Ford-Fulkerson)
// alogrithm and contains information about current flow in the edge and its maximum capacity.
//
//
// Implementations should be immutable.
type ResidualEdge interface {
	graph.Edge

	// MaxCapacity returns integer value of maximum capacity of the given edge
	MaxCapacity() int32
	// CurrentFlow returns integer value of current flow of the given edge
	CurrentFlow() int32
	// IsReverseEdge returns true if the edge is a reverse residual edge
	// (please refer to the Edmonds-Karp algorithm for more in-depth explanation)
	IsReverseEdge() bool
	// WithFlow returns a new instance of ResidualEdge with given flow value
	WithFlow(int32) ResidualEdge
	// ReversedResidualEdge returns an instance of ResidualEdge that it reversed to the current edge.
	//	e.g. if CurrentFlow() == 8, MaxCapacity() == 12 and IsReverseEdge() == false for a given edge,
	//  its ReversedResidualEdge() should return an edge with CurrentFlow() == 4 (which is 12 - 8), MaxCapacity() == 12, IsReverseEdge == true.
	//
	// If for a given edge IsReverseEdge() == false then its ReversedResidualEdge() should return an edge with  IsReverseEdge() == true.
	// In other words, e.ReversedResidualEdge().ReversedResidualEdge() must return an edge with the same return values of all its methods for any given edge e.
	ReversedResidualEdge() ResidualEdge
}

// FordFulkersonGraph is a generalized graph interface for Ford-Fulkerson maximum flow algorithms
type FordFulkersonGraph interface {
	traverse.Graph
	// ResidualGraph must return a new ResidualGraph instance
	ResidualGraph() ResidualGraph
}

// ResidualGraph is the residual graph Edmonds-Karp (Ford-Fulkerson) algorithm operates on
//
// Implementations must be mutable since SetFlow() method mutates a flow value.
type ResidualGraph interface {
	traverse.Graph
	SetFlow(uid, vid int64, flow int32)
	ResidualEdge(uid, vid int64) ResidualEdge
}

// AugmentingPath calculates a new augmenting path for
//  Edmonds-Karp maximum flow algorithm using breadth-first search.
//  The path found must be a shortest path that has available capacity.
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

// AugmentingPath calculates a new augmenting path for
//  Edmonds-Karp maximum flow algorithm using breadth-first search.
//  The path found must be a shortest path that has available capacity.
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

// AvailableCapacity calculates available capacity for a given path
//  which is a minimum capacity among the edges.
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

// EdmondsKarp is an implementation of Edmonds-Karp algorithm
// for computing the maximum flow of a network
type EdmondsKarp struct {
	g FordFulkersonGraph
}

// ResidualGraph returns the residual graph in its final state of the Edmonds-Karp algorithm
func (ek EdmondsKarp) ResidualGraph(source, sink graph.Node) ResidualGraph {
	rg := ek.g.ResidualGraph()
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
	return rg
}

// MaxFlow returns the maximym flow value
func (ek EdmondsKarp) MaxFlow(source, sink graph.Node) int32 {
	rg := ek.ResidualGraph(source, sink)
	var flow int32
	to := rg.From(source.ID())
	for to.Next() {
		flow += rg.ResidualEdge(source.ID(), to.Node().ID()).CurrentFlow()
	}
	return flow
}
