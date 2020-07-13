package flownet

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/simple"
)

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

func (g *SimpleFlowNetwork) ResidualGraph() ResidualGraph {
	nodes := copyNodes(g)
	from := copyFromEdges(g)
	return &SimpleResidualGraph{
		nodes: nodes,
		from:  from,
	}
}

type SimpleResidualGraph struct {
	nodes map[int64]graph.Node
	from  map[int64]map[int64]*ResidualEdge
}

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

func (g *SimpleResidualGraph) Edge(uid, vid int64) graph.Edge {
	return g.ResidualEdge(uid, vid)
}

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

func (g *SimpleResidualGraph) ResidualEdge(uid, vid int64) ResidualEdge {
	if _, ok := g.from[uid]; ok {
		keys := make([]int64, 0, len(g.from[uid]))
		for k := range g.from[uid] {
			keys = append(keys, k)
		}

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

type SimpleResidualEdge struct {
	simple.Edge
	capacity, flow int32
	reverse        bool
}

func (re SimpleResidualEdge) MaxCapacity() int32 {
	return re.capacity
}

func (re SimpleResidualEdge) CurrentFlow() int32 {
	return re.flow
}

func (re SimpleResidualEdge) IsReverseEdge() bool {
	return re.reverse
}

func (re SimpleResidualEdge) ReversedResidualEdge() ResidualEdge {
	return SimpleResidualEdge{
		simple.Edge{F: re.T, T: re.F},
		re.capacity, re.capacity - re.flow,
		!re.reverse,
	}
}

func (re SimpleResidualEdge) WithFlow(flow int32) ResidualEdge {
	return SimpleResidualEdge{
		re.Edge,
		re.capacity, flow,
		re.reverse,
	}
}
