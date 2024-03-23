// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ordered

import (
	"cmp"
	"slices"

	"gonum.org/v1/gonum/graph"
)

// ByID sorts a slice of graph.Node by ID.
func ByID[S interface{ ~[]E }, E graph.Node](n S) {
	slices.SortFunc(n, func(a, b E) int { return cmp.Compare(a.ID(), b.ID()) })
}

// BySliceIDs sorts a slice of []graph.Node lexically by the IDs of the
// []graph.Node.
func BySliceIDs(c [][]graph.Node) {
	slices.SortFunc(c, func(a, b []graph.Node) int {
		l := len(a)
		if len(b) < l {
			l = len(b)
		}
		for k, v := range a[:l] {
			if n := cmp.Compare(v.ID(), b[k].ID()); n != 0 {
				return n
			}
		}
		return cmp.Compare(len(a), len(b))
	})
}

// LinesByIDs sort a slice of graph.LinesByIDs lexically by the From IDs,
// then by the To IDs, finally by the Line IDs.
func LinesByIDs(n []graph.Line) {
	slices.SortFunc(n, func(a, b graph.Line) int {
		if n := cmp.Compare(a.From().ID(), b.From().ID()); n != 0 {
			return n
		}
		if n := cmp.Compare(a.To().ID(), b.To().ID()); n != 0 {
			return n
		}
		return cmp.Compare(a.ID(), b.ID())
	})
}
