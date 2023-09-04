// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf

import "gonum.org/v1/gonum/graph"

var (
	g *Graph

	_ graph.Graph              = g
	_ graph.Directed           = g
	_ graph.Multigraph         = g
	_ graph.DirectedMultigraph = g
	_ graph.NodeAdder          = g
	_ graph.NodeRemover        = g
	_ graph.LineAdder          = g
	_ graph.LineRemover        = g
)

// AddNode adds n to the graph. It panics if the added node ID matches an existing node ID.
func (g *Graph) AddNode(n graph.Node) {
	g.addNode(n)
}

// NewLine returns a new Line from the source to the destination node.
// The returned Line will have a graph-unique ID.
// The Line's ID does not become valid in g until the Line is added to g.
func (g *Graph) NewLine(from, to graph.Node) graph.Line {
	return g.newLine(from, to)
}

// NewNode returns a new unique Node to be added to g. The Node's ID does
// not become valid in g until the Node is added to g.
func (g *Graph) NewNode() graph.Node {
	return g.newNode()
}

// RemoveLine removes the line with the given end point and line IDs from the graph, leaving
// the terminal nodes. If the line does not exist it is a no-op.
func (g *Graph) RemoveLine(fid, tid, id int64) {
	g.removeLine(fid, tid, id)
}

// RemoveNode removes the node with the given ID from the graph, as well as any edges attached
// to it. If the node is not in the graph it is a no-op.
func (g *Graph) RemoveNode(id int64) {
	g.removeNode(id)
}

// SetLine adds l, a line from one node to another. If the nodes do not exist, they are added
// and are set to the nodes of the line otherwise.
func (g *Graph) SetLine(l graph.Line) {
	g.setLine(l)
}
