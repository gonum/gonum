// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

// All a node needs to do is identify itself. This allows the user to pass in nodes more
// interesting than an int, but also allow us to reap the benefits of having a map-storable,
// ==able type.
type Node interface {
	ID() int
}

// Allows edges to do something more interesting that just be a group of nodes. While the methods
// are called Head and Tail, they are not considered directed unless the given interface specifies
// otherwise.
type Edge interface {
	Head() Node
	Tail() Node
}

// A Graph implements the behavior of an undirected graph.
//
// All methods in Graph are implicitly undirected. Graph algorithms that care about directionality
// will intelligently choose the DirectedGraph behavior if that interface is also implemented,
// even if the function itself only takes in a Graph (or a super-interface of graph).
type Graph interface {
	// NodeExists returns true when node is currently in the graph.
	NodeExists(Node) bool

	// NodeList returns a list of all nodes in no particular order, useful for
	// determining things like if a graph is fully connected. The caller is
	// free to modify this list. Implementations should construct a new list
	// and not return internal representation.
	NodeList() []Node

	// Neighbors returns all nodes connected by any edge to this node.
	Neighbors(Node) []Node

	// EdgeBetween returns an edge between node and neighbor such that
	// Head is one argument and Tail is the other. If no
	// such edge exists, this function returns nil.
	EdgeBetween(node, neighbor Node) Edge
}

// Directed graphs are characterized by having seperable Heads and Tails in their edges.
// That is, if node1 goes to node2, that does not necessarily imply that node2 goes to node1.
//
// While it's possible for a directed graph to have fully reciprocal edges (i.e. the graph is
// symmetric) -- it is not required to be. The graph is also required to implement Graph
// because in many cases it can be useful to know all neighbors regardless of direction.
type DirectedGraph interface {
	Graph
	// Successors gives the nodes connected by OUTBOUND edges.
	// If the graph is an undirected graph, this set is equal to Predecessors.
	Successors(Node) []Node

	// EdgeTo returns an edge between node and successor such that
	// Head returns node and Tail returns successor, if no
	// such edge exists, this function returns nil.
	EdgeTo(node, successor Node) Edge

	// Predecessors gives the nodes connected by INBOUND edges.
	// If the graph is an undirected graph, this set is equal to Successors.
	Predecessors(Node) []Node
}

// Returns all undirected edges in the graph
type EdgeLister interface {
	EdgeList() []Edge
}

type EdgeListGraph interface {
	Graph
	EdgeLister
}

// Returns all directed edges in the graph.
type DirectedEdgeLister interface {
	DirectedEdgeList() []Edge
}

type DirectedEdgeListGraph interface {
	Graph
	DirectedEdgeLister
}

// A crunch graph forces a sparse graph to become a dense graph. That is, if the node IDs are
// [1,4,9,7] it would "crunch" the ids into the contiguous block [0,1,2,3]. Order is not
// required to be preserved between the non-cruched and crunched instances (that means in
// the example above 0 may correspond to 4 or 7 or 9, not necessarily 1).
//
// All dense graphs must have the first ID as 0.
type CrunchGraph interface {
	Graph
	Crunch()
}

// A Graph that implements Coster has an actual cost between adjacent nodes, also known as a
// weighted graph. If a graph implements coster and a function needs to read cost (e.g. A*),
// this function will take precedence over the Uniform Cost function (all weights are 1) if "nil"
// is passed in for the function argument.
//
// If the argument is nil, or the edge is invalid for some reason, this should return math.Inf(1)
type Coster interface {
	Cost(edge Edge) float64
}

// Guarantees that something implementing Coster is also a Graph.
type CostGraph interface {
	Coster
	Graph
}

// A graph that implements HeuristicCoster implements a heuristic between any two given nodes.
// Like Coster, if a graph implements this and a function needs a heuristic cost (e.g. A*), this
// function will take precedence over the Null Heuristic (always returns 0) if "nil" is passed in
// for the function argument. If HeuristicCost is not intended to be used, it can be implemented as
// the null heuristic (always returns 0.)
type HeuristicCoster interface {
	// HeuristicCost returns a heuristic cost between any two nodes.
	HeuristicCost(node1, node2 Node) float64
}

// A MutableGraph is a graph that can be changed in an arbitrary way. It is useful for several
// algorithms; for instance, Johnson's Algorithm requires adding a temporary node and changing
// edge weights. Another case where this is used is computing minimum spanning trees. Since trees
// are graphs, a minimum spanning tree can be created using this interface.
//
// Note that just because a graph does not implement MutableGraph does not mean that this package
// expects it to be invariant (though even a MutableGraph should be treated as invariant while an
// algorithm is operating on it), it simply means that without this interface this package can not
// properly handle the graph in order to, say, fill it with a minimum spanning tree.
//
// In functions that take a MutableGraph as an argument, it should not be the same as the Graph
// argument as concurrent modification will likely cause problems.
//
// MutableGraphs should always record the IDs as they are represented -- which means they are
// sparse by nature.
//
// MutableGraphs are required to keep the exact Nodes and Edges passed in, and return
// the originals when asked.
type MutableGraph interface {
	CostGraph
	// NewNode adds a node with an arbitrary ID and returns the new, unique ID
	// used.
	NewNode() Node

	// Adds a node to the graph
	AddNode(Node)

	// AddEdge connects two nodes in the graph. Neither node is required
	// to have been added before this is called. If directed is false,
	// it also adds the reciprocal edge. If this is called a second time,
	// it overrides any existing edge.
	AddEdge(e Edge, cost float64, directed bool)

	// RemoveNode removes a node from the graph, as well as any edges
	// attached to it
	RemoveNode(Node)

	// RemoveEdge removes a connection between two nodes, but does not
	// remove Head nor Tail under any circumstance. As with AddEdge, if
	// directed is false it also removes the reciprocal edge. This function
	// should be treated as a no-op and not an error if the edge doesn't exist.
	RemoveEdge(e Edge, directed bool)

	// EmptyGraph clears the graph of all nodes and edges.
	EmptyGraph()
}

// A function that returns the cost of following an edge
type CostFunc func(Edge) float64

// Estimates the cost of travelling between two nodes
type HeuristicCostFunc func(Node, Node) float64

// Convenience constants for AddEdge and RemoveEdge
const (
	Directed   bool = true
	Undirected      = false
)
