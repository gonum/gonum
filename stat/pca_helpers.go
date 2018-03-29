// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"gonum.org/v1/gonum/mat"
)

// TransformDisjunctive transforms a Complete Disjunctive Table to a Transformed
// Complete Disjunctive Table in order to be able to perform Multiple
// Correspondence Analysis using PC.PrincipalComponents. TransformDisjunctive
// places the TCDT in dst and returns it.
//
// If dst is nil, a new mat.Dense is allocated. If dst is not a zero matrix,
// the dimensions of dst and cdt must match otherwise TransformDisjunctive will
// panic. If cdt contains values other 0 or 1 TransformDisjunctive will panic.
//
// It is safe to reuse cdt as dst.
func TransformDisjunctive(dst *mat.Dense, cdt mat.Matrix) *mat.Dense {
	r, c := cdt.Dims()

	if dst == nil {
		dst = mat.NewDense(r, c, nil)
	} else if dr, dc := dst.Dims(); !dst.IsZero() && (dr != r || dc != c) {
		panic(mat.ErrShape)
	}

	for j := 0; j < c; j++ {
		var p float64
		for i := 0; i < r; i++ {
			v := cdt.At(i, j)
			if v != 0 && v != 1 {
				panic("stat: input is not a complete disjunctive table")
			}
			p += v
		}
		for i := 0; i < r; i++ {
			dst.Set(i, j, cdt.At(i, j)/p-1)
		}
	}

	return dst
}
