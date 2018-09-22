// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

// Iterator is an item iterator.
type Iterator interface {
	// Next advances the iterator and returns whether
	// the next call to the item method will return a
	// non-nil item.
	//
	// Next should be called prior to any call to the
	// iterator's item retrieval method after the
	// iterator has been obtained or reset.
	//
	// The order of iteration is implementation
	// dependent.
	Next() bool

	// Len returns the number of items remaining in the
	// iterator.
	Len() int

	// Reset returns the iterator to its start position.
	Reset()
}

// Nodes is a Node iterator.
type Nodes interface {
	Iterator

	// Node returns the current Node from the iterator.
	Node() Node
}

// NodeSlicer wraps the NodeSlice method.
type NodeSlicer interface {
	// NodeSlice returns the set of nodes remaining
	// to be iterated by a Nodes iterator.
	// The holder of the iterator may arbitrarily
	// change elements in the returned slice, but
	// those changes may be reflected to other
	// iterators.
	NodeSlice() []Node
}

// NodesOf returns it.Len() nodes from it. If it is a NodeSlicer, the NodeSlice method
// is used to obtain the nodes. It is safe to pass a nil Nodes to NodesOf.
func NodesOf(it Nodes) []Node {
	if it == nil {
		return nil
	}
	switch it := it.(type) {
	case NodeSlicer:
		return it.NodeSlice()
	}
	n := make([]Node, 0, it.Len())
	for it.Next() {
		n = append(n, it.Node())
	}
	return n
}

// Edges is an Edge iterator.
type Edges interface {
	Iterator

	// Edge returns the current Edge from the iterator.
	Edge() Edge
}

// EdgeSlicer wraps the EdgeSlice method.
type EdgeSlicer interface {
	// EdgeSlice returns the set of edges remaining
	// to be iterated by an Edges iterator.
	// The holder of the iterator may arbitrarily
	// change elements in the returned slice, but
	// those changes may be reflected to other
	// iterators.
	EdgeSlice() []Edge
}

// EdgesOf returns it.Len() nodes from it. If it is an EdgeSlicer, the EdgeSlice method is used
// to obtain the edges. It is safe to pass a nil Edges to EdgesOf.
func EdgesOf(it Edges) []Edge {
	if it == nil {
		return nil
	}
	switch it := it.(type) {
	case EdgeSlicer:
		return it.EdgeSlice()
	}
	n := make([]Edge, 0, it.Len())
	for it.Next() {
		n = append(n, it.Edge())
	}
	return n
}

// WeightedEdges is a WeightedEdge iterator.
type WeightedEdges interface {
	Iterator

	// Edge returns the current Edge from the iterator.
	WeightedEdge() WeightedEdge
}

// WeightedEdgeSlicer wraps the WeightedEdgeSlice method.
type WeightedEdgeSlicer interface {
	// EdgeSlice returns the set of edges remaining
	// to be iterated by an Edges iterator.
	// The holder of the iterator may arbitrarily
	// change elements in the returned slice, but
	// those changes may be reflected to other
	// iterators.
	WeightedEdgeSlice() []WeightedEdge
}

// WeightedEdgesOf returns it.Len() weighted edge from it. If it is a WeightedEdgeSlicer, the
// WeightedEdgeSlice method is used to obtain the edges. It is safe to pass a nil WeightedEdges
// to WeightedEdgesOf.
func WeightedEdgesOf(it WeightedEdges) []WeightedEdge {
	if it == nil {
		return nil
	}
	switch it := it.(type) {
	case WeightedEdgeSlicer:
		return it.WeightedEdgeSlice()
	}
	n := make([]WeightedEdge, 0, it.Len())
	for it.Next() {
		n = append(n, it.WeightedEdge())
	}
	return n
}

// Lines is a Line iterator.
type Lines interface {
	Iterator

	// Line returns the current Line from the iterator.
	Line() Line
}

// LineSlicer wraps the LineSlice method.
type LineSlicer interface {
	// LineSlice returns the set of lines remaining
	// to be iterated by an Lines iterator.
	// The holder of the iterator may arbitrarily
	// change elements in the returned slice, but
	// those changes may be reflected to other
	// iterators.
	LineSlice() []Line
}

// LinesOf returns it.Len() nodes from it. If it is a LineSlicer, the LineSlice method is used
// to obtain the lines. It is safe to pass a nil Lines to LinesOf.
func LinesOf(it Lines) []Line {
	if it == nil {
		return nil
	}
	switch it := it.(type) {
	case LineSlicer:
		return it.LineSlice()
	}
	n := make([]Line, 0, it.Len())
	for it.Next() {
		n = append(n, it.Line())
	}
	return n
}

// WeightedLines is a WeightedLine iterator.
type WeightedLines interface {
	Iterator

	// Line returns the current Line from the iterator.
	WeightedLine() WeightedLine
}

// WeightedLineSlicer wraps the WeightedLineSlice method.
type WeightedLineSlicer interface {
	// LineSlice returns the set of lines remaining
	// to be iterated by an Lines iterator.
	// The holder of the iterator may arbitrarily
	// change elements in the returned slice, but
	// those changes may be reflected to other
	// iterators.
	WeightedLineSlice() []WeightedLine
}

// WeightedLinesOf returns it.Len() weighted line from it. If it is a WeightedLineSlicer, the
// WeightedLineSlice method is used to obtain the lines. It is safe to pass a nil WeightedLines
// to WeightedLinesOf.
func WeightedLinesOf(it WeightedLines) []WeightedLine {
	if it == nil {
		return nil
	}
	switch it := it.(type) {
	case WeightedLineSlicer:
		return it.WeightedLineSlice()
	}
	n := make([]WeightedLine, 0, it.Len())
	for it.Next() {
		n = append(n, it.WeightedLine())
	}
	return n
}
