// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import "math"

// All a node needs to do is identify itself. This allows the user to pass in nodes more
// interesting than an int, but also allow us to reap the benefits of having a map-storable,
// comparable type.
type Node interface {
	ID() int
}

// Allows edges to do something more interesting that just be a group of nodes. While the methods
// are called From and To, they are not considered directed unless the given interface specifies
// otherwise.
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

type EdgeListGraph interface {
	Graph
	EdgeLister
}

// A Graph that implements Coster has an actual cost between adjacent nodes, also known as a
// weighted graph. If a graph implements coster and a function needs to read cost (e.g. A*),
// this function will take precedence over the Uniform Cost function (all weights are 1) if "nil"
// is passed in for the function argument.
//
// If the argument is nil, or the edge is invalid for some reason, this should return math.Inf(1)
type Coster interface {
	Cost(Edge) float64
}

type CostGraph interface {
	Coster
	Graph
}

type CostDirectedGraph interface {
	Coster
	Directed
}

// A graph that implements HeuristicCoster implements a heuristic between any two given nodes.
// Like Coster, if a graph implements this and a function needs a heuristic cost (e.g. A*), this
// function will take precedence over the Null Heuristic (always returns 0) if "nil" is passed in
// for the function argument. If HeuristicCost is not intended to be used, it can be implemented as
// the null heuristic (always returns 0).
type HeuristicCoster interface {
	// HeuristicCost returns a heuristic cost between any two nodes.
	HeuristicCost(n1, n2 Node) float64
}

// A Mutable is a graph that can have arbitrary nodes and edges added or removed.
//
// Anything implementing Mutable is required to store the actual argument. So if AddNode(myNode) is
// called and later a user calls on the graph graph.Nodes(), the node added by AddNode must be
// an the exact node, not a new node with the same ID.
//
// In any case where conflict is possible (e.g. adding two nodes with the same ID), the later
// call always supercedes the earlier one.
//
// Functions will generally expect one of MutableGraph or MutableDirectedGraph and not Mutable
// itself. That said, any function that takes Mutable[x], the destination mutable should
// always be a different graph than the source.
type Mutable interface {
	// NewNode returns a node with a unique arbitrary ID.
	NewNode() Node

	// Adds a node to the graph. AddNode panics if
	// the added node ID matches an existing node ID.
	AddNode(Node)

	// RemoveNode removes a node from the graph, as well as any edges
	// attached to it. If no such node exists, this is a no-op, not an error.
	RemoveNode(Node)
}

// MutableGraph is an interface ensuring the implementation of the ability to construct
// an arbitrary undirected graph. It is very important to note that any implementation
// of MutableGraph absolutely cannot safely implement the DirectedGraph interface.
//
// A MutableGraph is required to store any Edge argument in the same way Mutable must
// store a Node argument -- any retrieval call is required to return the exact supplied edge.
// This is what makes it incompatible with DirectedGraph.
//
// The reasoning is this: if you call AddUndirectedEdge(Edge{head,tail}); you are required
// to return the exact edge passed in when a retrieval method (EdgeTo/EdgeBetween) is called.
// If I call EdgeTo(tail,head), this means that since the edge exists, and was added as
// Edge{head,tail} this function MUST return Edge{head,tail}. However, EdgeTo requires this
// be returned as Edge{tail,head}. Thus there's a conflict that cannot be resolved between the
// two interface requirements.
type MutableGraph interface {
	CostGraph
	Mutable

	// AddUndirected adds an undirected edge between
	// distinct nodes. If the nodes do not exist, they
	// are added. AddUndirectedEdge will panic if the
	// IDs of the e.From and e.To are equal.
	AddUndirectedEdge(e Edge, cost float64)

	// RemoveEdge clears the stored edge between two nodes. Calling this will never
	// remove a node. If the edge does not exist this is a no-op, not an error.
	RemoveUndirectedEdge(Edge)
}

// MutableDirectedGraph is an interface that ensures one can construct an arbitrary directed
// graph. Naturally, a MutableDirectedGraph works for both undirected and directed cases,
// but simply using a MutableGraph may be cleaner. As the documentation for MutableGraph
// notes, however, a graph cannot safely implement MutableGraph and MutableDirectedGraph
// at the same time, because of the functionality of a EdgeTo in DirectedGraph.
type MutableDirectedGraph interface {
	CostDirectedGraph
	Mutable

	// AddUndirected adds a directed edge from one
	// node to another. If the nodes do not exist, they
	// are added. AddDirectedEdge will panic if the
	// IDs of the e.From and e.To are equal.
	AddDirectedEdge(e Edge, cost float64)

	// Removes an edge FROM e.From TO e.To. If no such edge exists, this is a no-op,
	// not an error.
	RemoveDirectedEdge(Edge)
}

// A function that returns the cost of following an edge
type CostFunc func(Edge) float64

// UniformCost returns an edge cost of 1 for a non-nil Edge and Inf for a nil Edge.
func UniformCost(e Edge) float64 {
	if e == nil {
		return math.Inf(1)
	}
	return 1
}

// Estimates the cost of travelling between two nodes
type HeuristicCostFunc func(Node, Node) float64

// CopyUndirected copies nodes and edges as undirected edges from the source to the
// destination without first clearing the destination. If the source does not
// provide edge weights, UniformCost is used.
func CopyUndirected(dst MutableGraph, src Graph) {
	var weight CostFunc
	if g, ok := src.(Coster); ok {
		weight = g.Cost
	} else {
		weight = UniformCost
	}

	for _, node := range src.Nodes() {
		succs := src.From(node)
		dst.AddNode(node)
		for _, succ := range succs {
			edge := src.Edge(node, succ)
			dst.AddUndirectedEdge(edge, weight(edge))
		}
	}
}

// CopyDirected copies nodes and edges as directed edges from the source to the
// destination without first clearing the destination. If src is undirected both
// directions will be present in the destination after the copy is complete. If
// the source does not provide edge weights, UniformCost is used.
func CopyDirected(dst MutableDirectedGraph, src Graph) {
	var weight CostFunc
	if g, ok := src.(Coster); ok {
		weight = g.Cost
	} else {
		weight = UniformCost
	}

	for _, node := range src.Nodes() {
		succs := src.From(node)
		dst.AddNode(node)
		for _, succ := range succs {
			edge := src.Edge(node, succ)
			dst.AddDirectedEdge(edge, weight(edge))
		}
	}
}
