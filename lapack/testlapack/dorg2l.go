// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas/blas64"
)

type Dorg2ler interface {
	Dorg2l(m, n, k int, a []float64, lda int, tau, work []float64)
	Dgeql2er
}

func Dorg2lTest(t *testing.T, impl Dorg2ler) {
	const tol = 1e-14

	rnd := rand.New(rand.NewPCG(1, 1))
	for _, test := range []struct {
		m, n, k, lda int
	}{
		{5, 4, 3, 0},
		{5, 4, 4, 0},
		{3, 3, 2, 0},
		{5, 5, 5, 0},
		{5, 4, 3, 11},
		{5, 4, 4, 11},
		{3, 3, 2, 11},
		{5, 5, 5, 11},
	} {
		m := test.m
		n := test.n
		k := test.k
		lda := test.lda
		if lda == 0 {
			lda = n
		}

		a := make([]float64, m*lda)
		for i := range a {
			a[i] = rnd.NormFloat64()
		}
		tau := nanSlice(max(m, n))
		work := make([]float64, n)
		impl.Dgeql2(m, n, a, lda, tau, work)

		impl.Dorg2l(m, n, k, a, lda, tau[n-k:], work)

		q := blas64.General{Rows: m, Cols: n, Data: a, Stride: lda}
		if resid := residualOrthogonal(q, false); resid > tol {
			t.Errorf("Case m=%v, n=%v, k=%v, lda=%v: columns of Q not orthonormal; resid=%v, want<=%v", m, n, k, lda, resid, tol)
		}
	}
}
