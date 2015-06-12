// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import "math"

// Node is a graph node. It returns a graph-unique integer ID.
type Node interface {
	ID() int
}

// Edge is a graph edge. In directed graphs, the direction of the
// edge is given from -> to, otherwise the edge is semantically
// unordered.
type Edge interface {
	From() Node
	To() Node
}

// Graph is a generalized graph.
type Graph interface {
	// Has returns whether the node exists within the graph.
	Has(Node) bool

	// Nodes returns all the nodes in the graph.
	Nodes() []Node

	// From returns all nodes that can be reached from
	// the given node.
	From(Node) []Node

	// HasEdge returns whether an edge exists between
	// nodes x and y without considering direction.
	HasEdge(x, y Node) bool

	// Edge returns the edge between nodes u and v when
	// the nodes returned by From(u) include v.
	Edge(u, v Node) Edge
}

// Undirected is an undirected graph.
type Undirected interface {
	Graph

	// EdgeBetween returns the edge between nodes u and v.
	EdgeBetween(u, v Node) Edge
}

// Directed is a directed graph.
type Directed interface {
	Graph

	// EdgeFromTo returns the edge leading from u to v.
	EdgeFromTo(u, v Node) Edge

	// To returns all nodes that can be lead to the
	// given node.
	To(Node) []Node
}

// EdgeLister wraps the Edges method.
type EdgeLister interface {
	Edges() []Edge
}

// Weighter wraps the Weight method.
type Weighter interface {
	// Weight returns the edge weight for the parameter,
	Weight(Edge) float64
}

// Mutable wraps generalized graph alteration methods.
type Mutable interface {
	// NewNode returns a node with a unique arbitrary ID.
	NewNode() Node

	// Adds a node to the graph. AddNode panics if
	// the added node ID matches an existing node ID.
	AddNode(Node)

	// RemoveNode removes a node from the graph, as
	// well as any edges attached to it. If the node
	// is not in the graph it is a no-op.
	RemoveNode(Node)

	// SetEdge adds an edge from one node to another.
	// If the nodes do not exist, they are added.
	// SetEdge will panic if the IDs of the e.From
	// and e.To are equal.
	SetEdge(e Edge, cost float64)

	// RemoveEdge removes the given edge, leaving the
	// terminal nodes. If the edge does not exist it
	// is a no-op.
	RemoveEdge(Edge)
}

// MutableUndirected is an undirected graph that can be arbitrarily altered.
type MutableUndirected interface {
	Undirected
	Mutable
}

// MutableDirected is a directed graph that can be arbitrarily altered.
type MutableDirected interface {
	Directed
	Mutable
}

// WeightFunc is a mapping between an edge and an edge weight.
type WeightFunc func(Edge) float64

// UniformCost is a WeightFunc that returns an edge cost of 1 for a non-nil Edge
// and Inf for a nil Edge.
func UniformCost(e Edge) float64 {
	if e == nil {
		return math.Inf(1)
	}
	return 1
}

// CopyUndirected copies nodes and edges as undirected edges from the source to the
// destination without first clearing the destination. If the source does not
// provide edge weights, UniformCost is used.
//
// Note that if the source is a directed graph and a fundamental cycle exists with
// two node where the edge weights differ, the resulting destination graph's edge
// weight between those nodes is undefined.
func CopyUndirected(dst MutableUndirected, src Graph) {
	var weight WeightFunc
	if g, ok := src.(Weighter); ok {
		weight = g.Weight
	} else {
		weight = UniformCost
	}

	for _, node := range src.Nodes() {
		succs := src.From(node)
		dst.AddNode(node)
		for _, succ := range succs {
			edge := src.Edge(node, succ)
			dst.SetEdge(edge, weight(edge))
		}
	}
}

// CopyDirected copies nodes and edges as directed edges from the source to the
// destination without first clearing the destination. If src is undirected both
// directions will be present in the destination after the copy is complete. If
// the source does not provide edge weights, UniformCost is used.
func CopyDirected(dst MutableDirected, src Graph) {
	var weight WeightFunc
	if g, ok := src.(Weighter); ok {
		weight = g.Weight
	} else {
		weight = UniformCost
	}

	for _, node := range src.Nodes() {
		succs := src.From(node)
		dst.AddNode(node)
		for _, succ := range succs {
			edge := src.Edge(node, succ)
			dst.SetEdge(edge, weight(edge))
		}
	}
}
