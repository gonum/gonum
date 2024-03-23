// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sorted

import (
	"cmp"
	"slices"
)

// BySliceValues sorts a slice of []cmp.Ordered (such as []int64)
// lexically by the values of the []cmp.Ordered.
func BySliceValues[S interface{ ~[]E }, E cmp.Ordered](c []S) {
	slices.SortFunc(c, func(a, b S) int {
		l := len(a)
		if len(b) < l {
			l = len(b)
		}
		for k, v := range a[:l] {
			if n := cmp.Compare(v, b[k]); n != 0 {
				return n
			}
		}
		return cmp.Compare(len(a), len(b))
	})
}
