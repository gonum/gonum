// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lapack

import "sort"

// Dlasrt sorts the numbers in the input slice d. If s == SortIncreasing,
// the elements are sorted in increasing order. If s == SortDecreasing,
// the elements are sorted in decreasing order. For other values of s Dlasrt
// will panic.
//
// Dlasrt is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dlasrt(s Sort, n int, d []float64) {
	checkVector(n, d, 1)
	d = d[:n]
	switch s {
	default:
		panic(badSort)
	case SortIncreasing:
		sort.Float64s(d)
	case SortDecreasing:
		sort.Sort(sort.Reverse(sort.Float64Slice(d)))
	}
}
