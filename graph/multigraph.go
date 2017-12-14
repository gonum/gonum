// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

// Line is an edge in a multigraph. A Line returns an ID that must
// distinguish Lines sharing Node end points.
type Line interface {
	Edge
	ID() int64
}

// WeightedLine is a weighted multigraph edge.
type WeightedLine interface {
	Line
	Weight() float64
}

// Multigraph is a generalized multigraph.
type Multigraph interface {
	// Has returns whether the node exists within the multigraph.
	Has(Node) bool

	// Nodes returns all the nodes in the multigraph.
	Nodes() []Node

	// From returns all nodes that can be reached directly
	// from the given node.
	From(Node) []Node

	// HasEdgeBetween returns whether an edge exists between
	// nodes x and y without considering direction.
	HasEdgeBetween(x, y Node) bool

	// Lines returns the lines from u to v if any such lines
	// exist and nil otherwise. The node v must be directly
	// reachable from u as defined by the From method.
	Lines(u, v Node) []Line
}

// WeightedMultigraph is a weighted multigraph.
type WeightedMultigraph interface {
	Multigraph

	// WeightedLines returns the weighted lines from u to v if
	// any such lines exist and nil otherwise. The node v must
	// be directly reachable from u as defined by the
	// From method.
	WeightedLines(u, v Node) []WeightedLine
}

// UndirectedMultigraph is an undirected multigraph.
type UndirectedMultigraph interface {
	Multigraph

	// LinesBetween returns the lines between nodes x and y.
	LinesBetween(x, y Node) []Line
}

// WeightedUndirectedMultigraph is a weighted undirected multigraph.
type WeightedUndirectedMultigraph interface {
	WeightedMultigraph

	// WeightedLinesBetween returns the lines between nodes
	// x and y.
	WeightedLinesBetween(x, y Node) []WeightedLine
}

// DirectedMultigraph is a directed multigraph.
type DirectedMultigraph interface {
	Multigraph

	// HasEdgeFromTo returns whether an edge exists
	// in the multigraph from u to v.
	HasEdgeFromTo(u, v Node) bool

	// To returns all nodes that can reach directly
	// to the given node.
	To(Node) []Node
}

// WeightedDirectedMultigraph is a weighted directed multigraph.
type WeightedDirectedMultigraph interface {
	WeightedMultigraph

	// HasEdgeFromTo returns whether an edge exists
	// in the multigraph from u to v.
	HasEdgeFromTo(u, v Node) bool

	// To returns all nodes that can reach directly
	// to the given node.
	To(Node) []Node
}

// LineAdder is an interface for adding lines to a multigraph.
type LineAdder interface {
	// NewLine returns a new Line from the source to the destination node.
	NewLine(from, to Node) Line

	// SetLine adds a Line from one node to another.
	// If the multigraph supports node addition the nodes
	// will be added if they do not exist, otherwise
	// SetLine will panic.
	SetLine(l Line)
}

// WeightedLineAdder is an interface for adding lines to a multigraph.
type WeightedLineAdder interface {
	// NewWeightedLine returns a new WeightedLine from
	// the source to the destination node.
	NewWeightedLine(from, to Node, weight float64) WeightedLine

	// SetWeightedLine adds a weighted line from one node
	// to another. If the multigraph supports node addition
	// the nodes will be added if they do not exist,
	// otherwise SetWeightedLine will panic.
	SetWeightedLine(e WeightedLine)
}

// LineRemover is an interface for removing lines from a multigraph.
type LineRemover interface {
	// RemoveLine removes the given line, leaving the
	// terminal nodes. If the line does not exist it
	// is a no-op.
	RemoveLine(Line)
}

// MultigraphBuilder is a multigraph that can have nodes and lines added.
type MultigraphBuilder interface {
	NodeAdder
	LineAdder
}

// WeightedMultigraphBuilder is a multigraph that can have nodes and weighted lines added.
type WeightedMultigraphBuilder interface {
	NodeAdder
	WeightedLineAdder
}

// UndirectedMultgraphBuilder is an undirected multigraph builder.
type UndirectedMultigraphBuilder interface {
	UndirectedMultigraph
	MultigraphBuilder
}

// UndirectedWeightedMultigraphBuilder is an undirected weighted multigraph builder.
type UndirectedWeightedMultigraphBuilder interface {
	UndirectedMultigraph
	WeightedMultigraphBuilder
}

// DirectedMultigraphBuilder is a directed multigraph builder.
type DirectedMultigraphBuilder interface {
	DirectedMultigraph
	MultigraphBuilder
}

// DirectedWeightedMultigraphBuilder is a directed weighted multigraph builder.
type DirectedWeightedMultigraphBuilder interface {
	DirectedMultigraph
	WeightedMultigraphBuilder
}
