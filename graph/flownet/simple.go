package flownet

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/simple"
)

// SimpleFlowNetwork is a simple implementation
// of FordFulkersonGraph interface.
//
// NOTE: The underlying WeightedDirectedGraph should
// have integer values for its weights since they will be converted to int32
type SimpleFlowNetwork struct {
	simple.WeightedDirectedGraph
}

func copyNodes(g *SimpleFlowNetwork) map[int64]graph.Node {
	nodes := make(map[int64]graph.Node)

	srcNodes := g.Nodes()
	for srcNodes.Next() {
		node := srcNodes.Node()
		nodes[node.ID()] = node
	}
	return nodes
}

func copyFromEdges(g *SimpleFlowNetwork) map[int64]map[int64]*ResidualEdge {
	from := make(map[int64]map[int64]*ResidualEdge)
	srcNodes := g.Nodes()
	for srcNodes.Next() {
		u := srcNodes.Node()
		uid := u.ID()
		to := g.From(uid)
		for to.Next() {
			v := to.Node()
			vid := v.ID()
			var edge ResidualEdge = SimpleResidualEdge{
				simple.Edge{F: u, T: v},
				int32(g.WeightedEdge(uid, vid).Weight()),
				0,
				false,
			}
			if fm, ok := from[uid]; ok {
				fm[vid] = &edge
			} else {
				from[uid] = map[int64]*ResidualEdge{vid: &edge}
			}
		}
	}
	return from
}

// ResidualGraph returns a new ResidualGraph instance which is constructed
// based on the underlying weighted graph.
func (g *SimpleFlowNetwork) ResidualGraph() ResidualGraph {
	nodes := copyNodes(g)
	from := copyFromEdges(g)
	return &SimpleResidualGraph{
		nodes: nodes,
		from:  from,
	}
}

// SimpleResidualGraph is a simple
// implementation of ResidualGraph interface
type SimpleResidualGraph struct {
	nodes map[int64]graph.Node
	from  map[int64]map[int64]*ResidualEdge
}

// From returns all nodes that can be reached directly
// from the node with the given ID.
func (g *SimpleResidualGraph) From(id int64) graph.Nodes {
	if _, ok := g.from[id]; !ok {
		return graph.Empty
	}

	from := make([]graph.Node, len(g.from[id]))
	i := 0
	for vid := range g.from[id] {
		from[i] = g.nodes[vid]
		i++
	}
	if len(from) == 0 {
		return graph.Empty
	}
	return iterator.NewOrderedNodes(from)
}

// Edge returns the edge from u to v, with IDs uid and vid,
// if such an edge exists and nil otherwise. The node v
// must be directly reachable from u as defined by
// the From method.
func (g *SimpleResidualGraph) Edge(uid, vid int64) graph.Edge {
	return g.ResidualEdge(uid, vid)
}

// SetFlow mutates the value of the given edge from uid to vid.
func (g *SimpleResidualGraph) SetFlow(uid, vid int64, flow int32) {
	if uid == vid {
		panic("max_flow: adding self edge")
	}
	edge := g.ResidualEdge(uid, vid)
	newEdge := edge.WithFlow(flow)
	if edge.IsReverseEdge() {
		newEdge = newEdge.ReversedResidualEdge()
	}
	g.from[newEdge.From().ID()][newEdge.To().ID()] = &newEdge
}

// ResidualEdge returns a ResidualEdge instance.
// If theres is no edge from uid to vid but there's an edge from vid to uid in the original graph,
// then a reverse edge is returned.
// If the two nodes are not connected at all nil is returned.
func (g *SimpleResidualGraph) ResidualEdge(uid, vid int64) ResidualEdge {
	if _, ok := g.from[uid]; ok {
		if _, ok := g.from[uid][vid]; ok {
			return *g.from[uid][vid]
		}
	}
	if _, ok := g.from[vid]; ok {
		if _, ok := g.from[vid][uid]; ok {
			edge := g.from[vid][uid]
			return (*edge).ReversedResidualEdge()
		}
	}
	return nil
}

// SimpleResidualEdge is an implementation of ResidualEdge interface
type SimpleResidualEdge struct {
	simple.Edge
	capacity, flow int32
	reverse        bool
}

// MaxCapacity returns integer value of maximum capacity of the given edge
func (re SimpleResidualEdge) MaxCapacity() int32 {
	return re.capacity
}

// CurrentFlow returns integer value of current flow of the given edge
func (re SimpleResidualEdge) CurrentFlow() int32 {
	return re.flow
}

// IsReverseEdge returns true if the edge is a reverse residual edge
func (re SimpleResidualEdge) IsReverseEdge() bool {
	return re.reverse
}

// ReversedResidualEdge returns an instance of ResidualEdge that it reversed to the current edge
func (re SimpleResidualEdge) ReversedResidualEdge() ResidualEdge {
	return SimpleResidualEdge{
		simple.Edge{F: re.T, T: re.F},
		re.capacity, re.capacity - re.flow,
		!re.reverse,
	}
}

// WithFlow returns a new instance of ResidualEdge with given flow value
func (re SimpleResidualEdge) WithFlow(flow int32) ResidualEdge {
	return SimpleResidualEdge{
		re.Edge,
		re.capacity, flow,
		re.reverse,
	}
}
