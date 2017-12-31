// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multi

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/uid"
)

var (
	dg *DirectedGraph

	_ graph.Graph              = dg
	_ graph.Directed           = dg
	_ graph.Multigraph         = dg
	_ graph.DirectedMultigraph = dg
)

// DirectedGraph implements a generalized directed graph.
type DirectedGraph struct {
	nodes map[int64]graph.Node
	from  map[int64]map[int64]map[int64]graph.Line
	to    map[int64]map[int64]map[int64]graph.Line

	nodeIDs uid.Set
	lineIDs uid.Set
}

// NewDirectedGraph returns a DirectedGraph.
func NewDirectedGraph() *DirectedGraph {
	return &DirectedGraph{
		nodes: make(map[int64]graph.Node),
		from:  make(map[int64]map[int64]map[int64]graph.Line),
		to:    make(map[int64]map[int64]map[int64]graph.Line),

		nodeIDs: uid.NewSet(),
		lineIDs: uid.NewSet(),
	}
}

// NewNode returns a new unique Node to be added to g. The Node's ID does
// not become valid in g until the Node is added to g.
func (g *DirectedGraph) NewNode() graph.Node {
	if len(g.nodes) == 0 {
		return Node(0)
	}
	if int64(len(g.nodes)) == uid.Max {
		panic("simple: cannot allocate node: no slot")
	}
	return Node(g.nodeIDs.NewID())
}

// AddNode adds n to the graph. It panics if the added node ID matches an existing node ID.
func (g *DirectedGraph) AddNode(n graph.Node) {
	if _, exists := g.nodes[n.ID()]; exists {
		panic(fmt.Sprintf("simple: node ID collision: %d", n.ID()))
	}
	g.nodes[n.ID()] = n
	g.from[n.ID()] = make(map[int64]map[int64]graph.Line)
	g.to[n.ID()] = make(map[int64]map[int64]graph.Line)
	g.nodeIDs.Use(n.ID())
}

// RemoveNode removes n from the graph, as well as any edges attached to it. If the node
// is not in the graph it is a no-op.
func (g *DirectedGraph) RemoveNode(n graph.Node) {
	if _, ok := g.nodes[n.ID()]; !ok {
		return
	}
	delete(g.nodes, n.ID())

	for from := range g.from[n.ID()] {
		delete(g.to[from], n.ID())
	}
	delete(g.from, n.ID())

	for to := range g.to[n.ID()] {
		delete(g.from[to], n.ID())
	}
	delete(g.to, n.ID())

	g.nodeIDs.Release(n.ID())
}

// NewLine returns a new Line from the source to the destination node.
// The returned Line will have a graph-unique ID.
// The Line's ID does not become valid in g until the Line is added to g.
func (g *DirectedGraph) NewLine(from, to graph.Node) graph.Line {
	return &Line{F: from, T: to, UID: g.lineIDs.NewID()}
}

// SetLine adds l, a line from one node to another. If the nodes do not exist, they are added.
func (g *DirectedGraph) SetLine(l graph.Line) {
	var (
		from = l.From()
		fid  = from.ID()
		to   = l.To()
		tid  = to.ID()
		lid  = l.ID()
	)

	if !g.Has(from) {
		g.AddNode(from)
	}
	if g.from[fid][tid] == nil {
		g.from[fid][tid] = make(map[int64]graph.Line)
	}
	if !g.Has(to) {
		g.AddNode(to)
	}
	if g.to[tid][fid] == nil {
		g.to[tid][fid] = make(map[int64]graph.Line)
	}

	g.from[fid][tid][lid] = l
	g.to[tid][fid][lid] = l
	g.lineIDs.Use(lid)
}

// RemoveLine removes l from the graph, leaving the terminal nodes. If the line does not exist
// it is a no-op.
func (g *DirectedGraph) RemoveLine(l graph.Line) {
	from, to := l.From(), l.To()
	if _, ok := g.nodes[from.ID()]; !ok {
		return
	}
	if _, ok := g.nodes[to.ID()]; !ok {
		return
	}

	delete(g.from[from.ID()][to.ID()], l.ID())
	if len(g.from[from.ID()][to.ID()]) == 0 {
		delete(g.from[from.ID()], to.ID())
	}
	delete(g.to[to.ID()][from.ID()], l.ID())
	if len(g.to[to.ID()][from.ID()]) == 0 {
		delete(g.to[to.ID()], from.ID())
	}
	g.lineIDs.Release(l.ID())
}

// Node returns the node in the graph with the given ID.
func (g *DirectedGraph) Node(id int64) graph.Node {
	return g.nodes[id]
}

// Has returns whether the node exists within the graph.
func (g *DirectedGraph) Has(n graph.Node) bool {
	_, ok := g.nodes[n.ID()]
	return ok
}

// Nodes returns all the nodes in the graph.
func (g *DirectedGraph) Nodes() []graph.Node {
	if len(g.nodes) == 0 {
		return nil
	}
	nodes := make([]graph.Node, len(g.nodes))
	i := 0
	for _, n := range g.nodes {
		nodes[i] = n
		i++
	}
	return nodes
}

// Edges returns all the edges in the graph. Each edge in the returned slice
// is a multi.Edge.
func (g *DirectedGraph) Edges() []graph.Edge {
	var edges []graph.Edge
	for _, u := range g.nodes {
		for _, e := range g.from[u.ID()] {
			var lines Edge
			for _, l := range e {
				lines = append(lines, l)
			}
			if len(lines) != 0 {
				edges = append(edges, lines)
			}
		}
	}
	return edges
}

// From returns all nodes in g that can be reached directly from n.
func (g *DirectedGraph) From(n graph.Node) []graph.Node {
	if _, ok := g.from[n.ID()]; !ok {
		return nil
	}

	from := make([]graph.Node, len(g.from[n.ID()]))
	i := 0
	for id := range g.from[n.ID()] {
		from[i] = g.nodes[id]
		i++
	}
	return from
}

// To returns all nodes in g that can reach directly to n.
func (g *DirectedGraph) To(n graph.Node) []graph.Node {
	if _, ok := g.from[n.ID()]; !ok {
		return nil
	}

	to := make([]graph.Node, len(g.to[n.ID()]))
	i := 0
	for id := range g.to[n.ID()] {
		to[i] = g.nodes[id]
		i++
	}
	return to
}

// HasEdgeBetween returns whether an edge exists between nodes x and y without
// considering direction.
func (g *DirectedGraph) HasEdgeBetween(x, y graph.Node) bool {
	xid := x.ID()
	yid := y.ID()
	if _, ok := g.from[xid][yid]; ok {
		return true
	}
	_, ok := g.from[yid][xid]
	return ok
}

// Edge returns the edge from u to v if such an edge exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
// The returned graph.Edge is a multi.Edge if an edge exists.
func (g *DirectedGraph) Edge(u, v graph.Node) graph.Edge {
	lines := g.Lines(u, v)
	if len(lines) == 0 {
		return nil
	}
	return Edge(lines)
}

// Lines returns the lines from u to v if such any such lines exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
func (g *DirectedGraph) Lines(u, v graph.Node) []graph.Line {
	edge := g.from[u.ID()][v.ID()]
	if len(edge) == 0 {
		return nil
	}
	var lines []graph.Line
	for _, l := range edge {
		lines = append(lines, l)
	}
	return lines
}

// HasEdgeFromTo returns whether an edge exists in the graph from u to v.
func (g *DirectedGraph) HasEdgeFromTo(u, v graph.Node) bool {
	_, ok := g.from[u.ID()][v.ID()]
	return ok
}

// Degree returns the in+out degree of n in g.
func (g *DirectedGraph) Degree(n graph.Node) int {
	if _, ok := g.nodes[n.ID()]; !ok {
		return 0
	}
	var deg int
	for _, e := range g.from[n.ID()] {
		deg += len(e)
	}
	for _, e := range g.to[n.ID()] {
		deg += len(e)
	}
	return deg
}
