// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package simple provides a suite of simple graph implementations satisfying
// the gonum/graph interfaces.
package simple // import "gonum.org/v1/gonum/graph/simple"

import (
	"math"

	"golang.org/x/tools/container/intsets"

	"gonum.org/v1/gonum/graph"
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

// maxInt is the maximum value of the machine-dependent int type.
const maxInt int = int(^uint(0) >> 1)

// isSame returns whether two float64 values are the same where NaN values
// are equalable.
func isSame(a, b float64) bool {
	return a == b || (math.IsNaN(a) && math.IsNaN(b))
}

// idSet implements available ID storage.
type idSet struct {
	used, free intsets.Sparse
}

// newID returns a new unique ID.
func (s *idSet) newID() int {
	var id int
	if s.free.TakeMin(&id) {
		return id
	}
	if id = s.used.Max(); id < maxInt {
		return id + 1
	}
	for id = 0; id < maxInt; id++ {
		if !s.used.Has(id) {
			return id
		}
	}
	panic("unreachable")
}

// use adds the id to the used IDs in the idSet.
func (s *idSet) use(id int) {
	s.free.Remove(id)
	s.used.Insert(id)
}

// free frees the id for reuse.
func (s *idSet) release(id int) {
	s.free.Insert(id)
	s.used.Remove(id)
}
