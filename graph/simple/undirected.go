// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simple

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/set/uid"
)

var (
	ug *UndirectedGraph

	_ graph.Graph       = ug
	_ graph.Undirected  = ug
	_ graph.NodeAdder   = ug
	_ graph.NodeRemover = ug
	_ graph.EdgeAdder   = ug
	_ graph.EdgeRemover = ug
)

// UndirectedGraph implements a generalized undirected graph.
type UndirectedGraph struct {
	nodes map[int64]graph.Node
	edges map[int64]map[int64]graph.Edge

	nodeIDs *uid.Set
}

// NewUndirectedGraph returns an UndirectedGraph.
func NewUndirectedGraph() *UndirectedGraph {
	return &UndirectedGraph{
		nodes: make(map[int64]graph.Node),
		edges: make(map[int64]map[int64]graph.Edge),

		nodeIDs: uid.NewSet(),
	}
}

// AddNode adds n to the graph. It panics if the added node ID matches an existing node ID.
func (g *UndirectedGraph) AddNode(n graph.Node) {
	if _, exists := g.nodes[n.ID()]; exists {
		panic(fmt.Sprintf("simple: node ID collision: %d", n.ID()))
	}
	g.nodes[n.ID()] = n
	g.nodeIDs.Use(n.ID())
}

// Edge returns the edge from u to v if such an edge exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
func (g *UndirectedGraph) Edge(uid, vid int64) graph.Edge {
	return g.EdgeBetween(uid, vid)
}

// EdgeBetween returns the edge between nodes x and y.
func (g *UndirectedGraph) EdgeBetween(xid, yid int64) graph.Edge {
	edge, ok := g.edges[xid][yid]
	if !ok {
		return nil
	}
	if edge.From().ID() == xid {
		return edge
	}
	return edge.ReversedEdge()
}

// Edges returns all the edges in the graph.
func (g *UndirectedGraph) Edges() graph.Edges {
	if len(g.edges) == 0 {
		return graph.Empty
	}
	var edges []graph.Edge
	for xid, u := range g.edges {
		for yid, e := range u {
			if yid < xid {
				// Do not consider edges when the To node ID is
				// before the From node ID. Both orientations
				// are stored.
				continue
			}
			edges = append(edges, e)
		}
	}
	if len(edges) == 0 {
		return graph.Empty
	}
	return iterator.NewOrderedEdges(edges)
}

// From returns all nodes in g that can be reached directly from n.
//
// The returned graph.Nodes is only valid until the next mutation of
// the receiver.
func (g *UndirectedGraph) From(id int64) graph.Nodes {
	if len(g.edges[id]) == 0 {
		return graph.Empty
	}
	return iterator.NewNodesByEdge(g.nodes, g.edges[id])
}

// HasEdgeBetween returns whether an edge exists between nodes x and y.
func (g *UndirectedGraph) HasEdgeBetween(xid, yid int64) bool {
	_, ok := g.edges[xid][yid]
	return ok
}

// NewEdge returns a new Edge from the source to the destination node.
func (g *UndirectedGraph) NewEdge(from, to graph.Node) graph.Edge {
	return Edge{F: from, T: to}
}

// NewNode returns a new unique Node to be added to g. The Node's ID does
// not become valid in g until the Node is added to g.
func (g *UndirectedGraph) NewNode() graph.Node {
	if len(g.nodes) == 0 {
		return Node(0)
	}
	if int64(len(g.nodes)) == uid.Max {
		panic("simple: cannot allocate node: no slot")
	}
	return Node(g.nodeIDs.NewID())
}

// Node returns the node with the given ID if it exists in the graph,
// and nil otherwise.
func (g *UndirectedGraph) Node(id int64) graph.Node {
	return g.nodes[id]
}

// Nodes returns all the nodes in the graph.
//
// The returned graph.Nodes is only valid until the next mutation of
// the receiver.
func (g *UndirectedGraph) Nodes() graph.Nodes {
	if len(g.nodes) == 0 {
		return graph.Empty
	}
	return iterator.NewNodes(g.nodes)
}

// NodeWithID returns a Node with the given ID if possible. If a graph.Node
// is returned that is not already in the graph NodeWithID will return true
// for new and the graph.Node must be added to the graph before use.
func (g *UndirectedGraph) NodeWithID(id int64) (n graph.Node, new bool) {
	n, ok := g.nodes[id]
	if ok {
		return n, false
	}
	return Node(id), true
}

// RemoveEdge removes the edge with the given end IDs from the graph, leaving the terminal nodes.
// If the edge does not exist it is a no-op.
func (g *UndirectedGraph) RemoveEdge(fid, tid int64) {
	if _, ok := g.nodes[fid]; !ok {
		return
	}
	if _, ok := g.nodes[tid]; !ok {
		return
	}

	delete(g.edges[fid], tid)
	delete(g.edges[tid], fid)
}

// RemoveNode removes the node with the given ID from the graph, as well as any edges attached
// to it. If the node is not in the graph it is a no-op.
func (g *UndirectedGraph) RemoveNode(id int64) {
	if _, ok := g.nodes[id]; !ok {
		return
	}
	delete(g.nodes, id)

	for from := range g.edges[id] {
		delete(g.edges[from], id)
	}
	delete(g.edges, id)

	g.nodeIDs.Release(id)
}

// SetEdge adds e, an edge from one node to another. If the nodes do not exist, they are added
// and are set to the nodes of the edge otherwise.
// It will panic if the IDs of the e.From and e.To are equal.
func (g *UndirectedGraph) SetEdge(e graph.Edge) {
	var (
		from = e.From()
		fid  = from.ID()
		to   = e.To()
		tid  = to.ID()
	)

	if fid == tid {
		panic("simple: adding self edge")
	}

	if _, ok := g.nodes[fid]; !ok {
		g.AddNode(from)
	} else {
		g.nodes[fid] = from
	}
	if _, ok := g.nodes[tid]; !ok {
		g.AddNode(to)
	} else {
		g.nodes[tid] = to
	}

	if fm, ok := g.edges[fid]; ok {
		fm[tid] = e
	} else {
		g.edges[fid] = map[int64]graph.Edge{tid: e}
	}
	if tm, ok := g.edges[tid]; ok {
		tm[fid] = e
	} else {
		g.edges[tid] = map[int64]graph.Edge{fid: e}
	}
}
