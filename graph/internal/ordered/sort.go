// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ordered

import (
	"sort"

	"gonum.org/v1/gonum/graph"
)

// ByID sorts a slice of graph.Node by ID.
func ByID(n []graph.Node) {
	sort.Slice(n, func(i, j int) bool { return n[i].ID() < n[j].ID() })
}

// BySliceValues sorts a slice of []int64 lexically by the values of the
// []int64.
func BySliceValues(c [][]int64) {
	sort.Slice(c, func(i, j int) bool {
		a, b := c[i], c[j]
		l := len(a)
		if len(b) < l {
			l = len(b)
		}
		for k, v := range a[:l] {
			if v < b[k] {
				return true
			}
			if v > b[k] {
				return false
			}
		}
		return len(a) < len(b)
	})
}

// BySliceIDs sorts a slice of []graph.Node lexically by the IDs of the
// []graph.Node.
func BySliceIDs(c [][]graph.Node) {
	sort.Slice(c, func(i, j int) bool {
		a, b := c[i], c[j]
		l := len(a)
		if len(b) < l {
			l = len(b)
		}
		for k, v := range a[:l] {
			if v.ID() < b[k].ID() {
				return true
			}
			if v.ID() > b[k].ID() {
				return false
			}
		}
		return len(a) < len(b)
	})
}

// Int64s sorts a slice of int64.
func Int64s(s []int64) {
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
}

// LinesByIDs sort a slice of graph.LinesByIDs lexically by the From IDs,
// then by the To IDs, finally by the Line IDs.
func LinesByIDs(n []graph.Line) {
	sort.Slice(n, func(i, j int) bool {
		a, b := n[i], n[j]
		if a.From().ID() != b.From().ID() {
			return a.From().ID() < b.From().ID()
		}
		if a.To().ID() != b.To().ID() {
			return a.To().ID() < b.To().ID()
		}
		return n[i].ID() < n[j].ID()
	})
}

// Reverse reverses the order of nodes.
func Reverse(nodes []graph.Node) {
	for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}
}
