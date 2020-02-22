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

type Dgeqp3er interface {
	Dlapmter
	Dgeqp3(m, n int, a []float64, lda int, jpvt []int, tau, work []float64, lwork int)
}

func Dgeqp3Test(t *testing.T, impl Dgeqp3er) {
	rnd := rand.New(rand.NewSource(1))
	for _, m := range []int{0, 1, 2, 3, 4, 5, 12, 23, 129} {
		for _, n := range []int{0, 1, 2, 3, 4, 5, 12, 23, 129} {
			for _, lda := range []int{max(1, n), n + 3} {
				dgeqp3Test(t, impl, rnd, m, n, lda)
			}
		}
	}
}

func dgeqp3Test(t *testing.T, impl Dgeqp3er, rnd *rand.Rand, m, n, lda int) {
	const (
		tol = 1e-14

		all = iota
		some
		none
	)
	for _, free := range []int{all, some, none} {
		name := fmt.Sprintf("m=%d,n=%d,lda=%d,", m, n, lda)

		// Allocate m×n matrix A and fill it with random numbers.
		a := randomGeneral(m, n, lda, rnd)
		// Store a copy of A for later comparison.
		aCopy := cloneGeneral(a)
		// Allocate a slice of column pivots.
		jpvt := make([]int, n)
		for j := range jpvt {
			switch free {
			case all:
				// All columns are free.
				jpvt[j] = -1
				name += "free=all"
			case some:
				// Some columns are free, some are leading columns.
				jpvt[j] = rnd.Intn(2) - 1 // -1 or 0
				name += "free=some"
			case none:
				// All columns are leading.
				jpvt[j] = 0
				name += "free=none"
			default:
				panic("bad freedom")
			}
		}
		// Allocate a slice for scalar factors of elementary
		// reflectors and fill it with random numbers. Dgeqp3
		// will overwrite them with valid data.
		k := min(m, n)
		tau := make([]float64, k)
		for i := range tau {
			tau[i] = rnd.Float64()
		}
		// Get optimal workspace size for Dgeqp3.
		work := make([]float64, 1)
		impl.Dgeqp3(m, n, a.Data, a.Stride, jpvt, tau, work, -1)
		lwork := int(work[0])
		work = make([]float64, lwork)
		for i := range work {
			work[i] = rnd.Float64()
		}

		// Compute a QR factorization of A with column pivoting.
		impl.Dgeqp3(m, n, a.Data, a.Stride, jpvt, tau, work, lwork)

		// Compute Q based on the elementary reflectors stored in A.
		q := constructQ("QR", m, n, a.Data, a.Stride, tau)
		// Check that Q is orthogonal.
		if resid := residualOrthogonal(q, false); resid > tol*float64(max(m, n)) {
			t.Errorf("Case %v: Q not orthogonal; resid=%v, want<=%v", name, resid, tol*float64(max(m, n)))
		}

		// Copy the upper triangle of A into R.
		r := zeros(m, n, lda)
		for i := 0; i < m; i++ {
			for j := i; j < n; j++ {
				r.Data[i*r.Stride+j] = a.Data[i*a.Stride+j]
			}
		}
		// Compute Q*R - A*P:
		// 1. Rearrange the columns of A based on the permutation in jpvt.
		qrap := cloneGeneral(aCopy)
		impl.Dlapmt(true, qrap.Rows, qrap.Cols, qrap.Data, qrap.Stride, jpvt)
		// Compute Q*R - A*P.
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, q, r, -1, qrap)
		// Check that |Q*R - A*P| is small.
		resid := dlange(lapack.MaxColumnSum, qrap.Rows, qrap.Cols, qrap.Data, qrap.Stride)
		if resid > tol*float64(max(m, n)) {
			t.Errorf("Case %v: |Q*R - A*P|=%v, want<=%v", name, resid, tol*float64(max(m, n)))

		}
	}
}
