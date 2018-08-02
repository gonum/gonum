// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simple

import (
	"math"

	"gonum.org/v1/gonum/graph"
)

// Node is a simple graph node.
type Node int64

// ID returns the ID number of the node.
func (n Node) ID() int64 {
	return int64(n)
}

// Edge is a simple graph edge.
type Edge struct {
	F, T graph.Node
}

// From returns the from-node of the edge.
func (e Edge) From() graph.Node { return e.F }

// To returns the to-node of the edge.
func (e Edge) To() graph.Node { return e.T }

// WeightedEdge is a simple weighted graph edge.
type WeightedEdge struct {
	F, T graph.Node
	W    float64
}

// From returns the from-node of the edge.
func (e WeightedEdge) From() graph.Node { return e.F }

// To returns the to-node of the edge.
func (e WeightedEdge) To() graph.Node { return e.T }

// Weight returns the weight of the edge.
func (e WeightedEdge) Weight() float64 { return e.W }

// isSame returns whether two float64 values are the same where NaN values
// are equalable.
func isSame(a, b float64) bool {
	return a == b || (math.IsNaN(a) && math.IsNaN(b))
}

// NodeIterator implements the graph.Nodes and graph.NodeSlicer interfaces.
// The iteration order of NodeIterator is the order of nodes passed to
// NewNodeIterator.
type NodeIterator struct {
	idx   int
	nodes []graph.Node
}

// NewNodeIterator returns a NodeIterator initialized with the provided nodes.
func NewNodeIterator(nodes []graph.Node) *NodeIterator {
	return &NodeIterator{idx: -1, nodes: nodes}
}

func (n *NodeIterator) Len() int {
	if n.idx >= len(n.nodes) {
		return 0
	}
	if n.idx <= 0 {
		return len(n.nodes)
	}
	return len(n.nodes[n.idx:])
}

func (n *NodeIterator) Next() bool {
	if uint(n.idx)+1 < uint(len(n.nodes)) {
		n.idx++
		return true
	}
	n.idx = len(n.nodes)
	return false
}

func (n *NodeIterator) Node() graph.Node {
	if n.idx >= len(n.nodes) || n.idx < 0 {
		return nil
	}
	return n.nodes[n.idx]
}

func (n *NodeIterator) NodeSlice() []graph.Node {
	if n.idx >= len(n.nodes) {
		return nil
	}
	idx := n.idx
	if idx == -1 {
		idx = 0
	}
	n.idx = len(n.nodes)
	return n.nodes[idx:]
}

func (n *NodeIterator) Reset() {
	n.idx = -1
}

// implicitIterator implements the graph.Nodes interface.
type implicitIterator struct {
	beg, end int
	curr     int
}

func newImplicitIterator(beg, end int) *implicitIterator {
	return &implicitIterator{beg: beg, end: end, curr: beg - 1}
}

func (n *implicitIterator) Len() int {
	return n.end - n.curr
}

func (n *implicitIterator) Next() bool {
	if n.curr == n.end {
		return false
	}
	n.curr++
	return n.curr < n.end
}

func (n *implicitIterator) Node() graph.Node {
	if n.Len() == 0 || n.curr < n.beg {
		return nil
	}
	return Node(n.curr)
}

func (n *implicitIterator) Reset() {
	n.curr = n.beg - 1
}

// EdgeIterator implements the graph.Edges and graph.EdgeSlicer interfaces.
// The iteration order of EdgeIterator is the order of edges passed to
// NewEdgeIterator.
type EdgeIterator struct {
	idx   int
	edges []graph.Edge
}

// NewEdgeIterator returns an EdgeIterator initialized with the provided edges.
func NewEdgeIterator(edges []graph.Edge) *EdgeIterator {
	return &EdgeIterator{idx: -1, edges: edges}
}

func (e *EdgeIterator) Len() int {
	if e.idx >= len(e.edges) {
		return 0
	}
	if e.idx <= 0 {
		return len(e.edges)
	}
	return len(e.edges[e.idx:])
}

func (e *EdgeIterator) Next() bool {
	if uint(e.idx)+1 < uint(len(e.edges)) {
		e.idx++
		return true
	}
	e.idx = len(e.edges)
	return false
}

func (e *EdgeIterator) Edge() graph.Edge {
	if e.idx >= len(e.edges) || e.idx < 0 {
		return nil
	}
	return e.edges[e.idx]
}

func (e *EdgeIterator) EdgeSlice() []graph.Edge {
	if e.idx >= len(e.edges) {
		return nil
	}
	idx := e.idx
	if idx == -1 {
		idx = 0
	}
	e.idx = len(e.edges)
	return e.edges[idx:]
}

func (e *EdgeIterator) Reset() {
	e.idx = -1
}

// WeightedEdgeIterator implements the graph.Edges and graph.EdgeSlicer interfaces.
// The iteration order of WeightedEdgeIterator is the order of edges passed to
// NewEdgeIterator.
type WeightedEdgeIterator struct {
	idx   int
	edges []graph.WeightedEdge
}

// NewWeightedEdgeIterator returns an WeightedEdgeIterator initialized with the provided edges.
func NewWeightedEdgeIterator(edges []graph.WeightedEdge) *WeightedEdgeIterator {
	return &WeightedEdgeIterator{idx: -1, edges: edges}
}

func (e *WeightedEdgeIterator) Len() int {
	if e.idx >= len(e.edges) {
		return 0
	}
	if e.idx <= 0 {
		return len(e.edges)
	}
	return len(e.edges[e.idx:])
}

func (e *WeightedEdgeIterator) Next() bool {
	if uint(e.idx)+1 < uint(len(e.edges)) {
		e.idx++
		return true
	}
	e.idx = len(e.edges)
	return false
}

func (e *WeightedEdgeIterator) WeightedEdge() graph.WeightedEdge {
	if e.idx >= len(e.edges) || e.idx < 0 {
		return nil
	}
	return e.edges[e.idx]
}

func (e *WeightedEdgeIterator) WeightedEdges() []graph.WeightedEdge {
	if e.idx >= len(e.edges) {
		return nil
	}
	idx := e.idx
	if idx == -1 {
		idx = 0
	}
	e.idx = len(e.edges)
	return e.edges[idx:]
}

func (e *WeightedEdgeIterator) Reset() {
	e.idx = -1
}
