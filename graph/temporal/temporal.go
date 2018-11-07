// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package temporal

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
)

// Node is a simple graph node.
type Node int64

// ID returns the ID number of the node.
func (n Node) ID() int64 {
	return int64(n)
}

// TemporalLine is a temporal graph edge.
type TemporalLine struct {
	F, T graph.Node

	UID int64

	S, E uint64
}

// From returns the from-node of the line.
func (l TemporalLine) From() graph.Node { return l.F }

// To returns the to-node of the line.
func (l TemporalLine) To() graph.Node { return l.T }

// ID returns the ID of the line.
func (l TemporalLine) ID() int64 { return l.UID }

// Interval returns the edge starting time and ending time. The
// edge traversal time is the difference between the two.
func (l TemporalLine) Interval() (start, end uint64) {
	return l.S, l.E
}

// Create a LineStreamer with lines ordered by start time i.e.
// a temporal edge stream.
func NewLineStream(lines []graph.TemporalLine) graph.LineStreamer {
	 in := make([]graph.Line, len(lines))
	 for i := range lines {
	 	in[i] = lines[i]
	 }
	 return iterator.NewLineStream(in, func (a, b graph.Line) bool {
		ta, _ := a.(graph.TemporalLine).Interval()
		tb, _ := b.(graph.TemporalLine).Interval()
		return ta < tb
	 })
}