// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/lapack"
)

type Dorgr2er interface {
	Dorgr2(m, n, k int, a []float64, lda int, tau []float64, work []float64)

	Dgerqfer
}

func Dorgr2Test(t *testing.T, impl Dorgr2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, k := range []int{0, 1, 2, 5} {
		for _, m := range []int{k, k + 1, k + 2, k + 4} {
			for _, n := range []int{m, m + 1, m + 2, m + 4, m + 7} {
				for _, lda := range []int{max(1, n), n + 5} {
					dorgr2Test(t, impl, rnd, m, n, k, lda)
				}
			}
		}
	}
}

func dorgr2Test(t *testing.T, impl Dorgr2er, rnd *rand.Rand, m, n, k, lda int) {
	const tol = 1e-14

	name := fmt.Sprintf("m=%v,n=%v,k=%v,lda=%v", m, n, k, lda)

	// Generate a random m×n matrix A.
	a := randomGeneral(m, n, lda, rnd)

	// Compute the RQ decomposition of A.
	rq := cloneGeneral(a)
	tau := make([]float64, m)
	work := make([]float64, 1)
	impl.Dgerqf(m, n, rq.Data, rq.Stride, tau, work, -1)
	work = make([]float64, int(work[0]))
	impl.Dgerqf(m, n, rq.Data, rq.Stride, tau, work, len(work))

	tauCopy := make([]float64, len(tau))
	copy(tauCopy, tau)

	// Compute the matrix Q using Dorg2r.
	q := cloneGeneral(rq)
	impl.Dorgr2(m, n, k, q.Data, q.Stride, tau[m-k:m], work)

	if m == 0 {
		return
	}

	// Check that tau hasn't been modified.
	if !floats.Equal(tau, tauCopy) {
		t.Errorf("%v: unexpected modification in tau", name)
	}

	// Check that Q has orthonormal rows.
	res := residualOrthogonal(q, true)
	if res > tol || math.IsNaN(res) {
		t.Errorf("%v: residual |I - Q*Qᵀ| too large, got %v, want <= %v", name, res, tol)
	}

	if k == 0 {
		return
	}

	// Extract the k×m upper triangular matrix R from RQ[m-k:m,n-k:n].
	r := zeros(k, m, m)
	for i := 0; i < k; i++ {
		for j := 0; j < k; j++ {
			ii := rq.Rows - k + i
			jj := rq.Cols - k + j
			jr := r.Cols - k + j
			if i <= j {
				r.Data[i*r.Stride+jr] = rq.Data[ii*rq.Stride+jj]
			}
		}
	}

	// Construct a view A[m-k:m,0:n] of the last k rows of A.
	aRec := blas64.General{
		Rows:   k,
		Cols:   n,
		Data:   a.Data[(m-k)*a.Stride:],
		Stride: a.Stride,
	}
	// Compute A - R*Q.
	blas64.Gemm(blas.NoTrans, blas.NoTrans, -1, r, q, 1, aRec)
	// Check that |A - R*Q| is small.
	res = dlange(lapack.MaxColumnSum, aRec.Rows, aRec.Cols, aRec.Data, aRec.Stride)
	if res > tol || math.IsNaN(res) {
		t.Errorf("%v: residual |A - R*Q| too large, got %v, want <= %v", name, res, tol)
	}
}
