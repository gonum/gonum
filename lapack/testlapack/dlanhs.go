// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/lapack"
)

type Dlanhser interface {
	Dlanhs(norm lapack.MatrixNorm, n int, a []float64, lda int, work []float64) float64
}

func DlanhsTest(t *testing.T, impl Dlanhser) {
	const tol = 1e-15
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, n := range []int{0, 1, 2, 4, 9} {
		for _, lda := range []int{max(1, n), n + 5} {
			a := randomGeneral(n, n, lda, rnd)
			for _, norm := range []lapack.MatrixNorm{lapack.MaxAbs, lapack.MaxRowSum, lapack.MaxColumnSum, lapack.Frobenius} {
				var work []float64
				if norm == lapack.MaxColumnSum {
					work = nanSlice(n)
				}

				got := impl.Dlanhs(norm, a.Rows, a.Data, lda, work)

				// Zero out A below the first subdiagonal.
				for i := 2; i < n; i++ {
					for j := 0; j < max(0, i-1); j++ {
						a.Data[i*a.Stride+j] = 0
					}
				}
				want := dlange(norm, a.Rows, a.Cols, a.Data, a.Stride)

				if math.Abs(want-got) > tol*want {
					t.Errorf("Case n=%v,lda=%v,norm=%v: unexpected result. Want %v, got %v.", n, lda, normToString(norm), want, got)
				}
			}
		}
	}
}
