// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dgerq2er interface {
	Dgerq2(m, n int, a []float64, lda int, tau []float64, work []float64)
}

func Dgerq2Test(t *testing.T, impl Dgerq2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, m := range []int{0, 1, 2, 3, 4, 5, 6, 12, 23} {
		for _, n := range []int{0, 1, 2, 3, 4, 5, 6, 12, 23} {
			for _, lda := range []int{max(1, n), n + 4} {
				dgerq2Test(t, impl, rnd, m, n, lda)
			}
		}
	}
}

func dgerq2Test(t *testing.T, impl Dgerq2er, rnd *rand.Rand, m, n, lda int) {
	const tol = 1e-14

	name := fmt.Sprintf("m=%d,n=%d,lda=%d", m, n, lda)

	a := randomGeneral(m, n, lda, rnd)
	aCopy := cloneGeneral(a)

	k := min(m, n)
	tau := make([]float64, k)
	for i := range tau {
		tau[i] = rnd.Float64()
	}

	work := make([]float64, m)
	for i := range work {
		work[i] = rnd.Float64()
	}

	impl.Dgerq2(m, n, a.Data, a.Stride, tau, work)

	// Test that the RQ factorization has completed successfully. Compute
	// Q based on the vectors.
	q := constructQ("RQ", m, n, a.Data, a.Stride, tau)

	// Check that Q is orthogonal.
	if resid := residualOrthogonal(q, false); resid > tol {
		t.Errorf("Case %v: Q not orthogonal; resid=%v, want<=%v", name, resid, tol)
	}

	// Check that |R*Q - A| is small.
	r := zeros(m, n, n)
	for i := 0; i < m; i++ {
		off := m - n
		for j := max(0, i-off); j < n; j++ {
			r.Data[i*r.Stride+j] = a.Data[i*a.Stride+j]
		}
	}
	qra := cloneGeneral(aCopy)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, r, q, -1, qra)
	resid := dlange(lapack.MaxColumnSum, qra.Rows, qra.Cols, qra.Data, qra.Stride)
	if resid > tol*float64(m) {
		t.Errorf("Case %v: |R*Q - A|=%v, want<=%v", name, resid, tol*float64(m))
	}
}
