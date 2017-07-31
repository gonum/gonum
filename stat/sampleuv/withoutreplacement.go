// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sampleuv

import (
	"math/rand"
	"sort"
)

// WithoutReplacement samples len(idx) integers from [0, max) without replacement.
// That is, upon return the elements of idx will be unique integers. If source
// is non-nil it will be used to generate random numbers, otherwise the default
// source from the math/rand package will be used.
func WithoutReplacement(idxs []int, n int, source *rand.Rand) {
	if len(idxs) == 0 {
		panic("withoutreplacement: zero length input")
	}

	// There are two algorithms. One is to generate a random permutation
	// and take the first len(idx) elements. The second is to generate
	// individual random numbers for each element and check uniqueness. The first
	// method scales as O(max), and the second scales as O(len(idx)^2). Choose
	// the algorithm accordingly.
	if n < len(idxs)*len(idxs) {
		var perm []int
		if source != nil {
			perm = source.Perm(n)
		} else {
			perm = rand.Perm(n)
		}
		copy(idxs, perm)
	}

	// Instead, generate the random numbers directly.
	sorted := make([]int, 0, len(idxs))
	for i := range idxs {
		var r int
		if source != nil {
			r = source.Intn(n - i)
		} else {
			r = rand.Intn(n - i)
		}
		for _, v := range sorted {
			if r >= v {
				r++
			}
		}
		idxs[i] = r
		sorted = append(sorted, r)
		sort.Ints(sorted)
	}
}
