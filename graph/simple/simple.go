// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package simple provides a suite of simple graph implementations satisfying
// the gonum/graph interfaces.
package simple // import "gonum.org/v1/gonum/graph/simple"

import (
	"math"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/set"
)

// Node is a simple graph node.
type Node int

// ID returns the ID number of the node.
func (n Node) ID() int {
	return int(n)
}

// Edge is a simple graph edge.
type Edge struct {
	F, T graph.Node
	W    float64
}

// From returns the from-node of the edge.
func (e Edge) From() graph.Node { return e.F }

// To returns the to-node of the edge.
func (e Edge) To() graph.Node { return e.T }

// Weight returns the weight of the edge.
func (e Edge) Weight() float64 { return e.W }

// isSame returns whether two float64 values are the same where NaN values
// are equalable.
func isSame(a, b float64) bool {
	return a == b || (math.IsNaN(a) && math.IsNaN(b))
}

// maxInt is the maximum value of the machine-dependent int type.
const maxInt int = int(^uint(0) >> 1)

// idSet implements available ID storage.
type idSet struct {
	maxID      int
	used, free set.Ints
}

// newIDSet returns a new idSet. The returned value should not be passed
// except by pointer.
func newIDSet() idSet {
	return idSet{maxID: -1, used: make(set.Ints), free: make(set.Ints)}
}

// newID returns a new unique ID. The ID returned is not considered used
// until passed in a call to use.
func (s *idSet) newID() int {
	for id := range s.free {
		return id
	}
	if s.maxID != maxInt {
		return s.maxID + 1
	}
	for id := 0; id <= s.maxID+1; id++ {
		if !s.used.Has(id) {
			return id
		}
	}
	panic("unreachable")
}

// use adds the id to the used IDs in the idSet.
func (s *idSet) use(id int) {
	s.used.Add(id)
	s.free.Remove(id)
	if id > s.maxID {
		s.maxID = id
	}
}

// free frees the id for reuse.
func (s *idSet) release(id int) {
	s.free.Add(id)
	s.used.Remove(id)
}
