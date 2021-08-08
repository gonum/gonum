// Copyright ©2021 The Gonum Authors. All rights reserved.
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

type Dorgr2er interface {
	Dorgr2(m, n, k int, a []float64, lda int, tau []float64, work []float64)

	Dgerqfer
}

func Dorgr2Test(t *testing.T, impl Dorgr2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, k := range []int{0, 1, 2, 5} {
		for _, m := range []int{k, k + 1, k + 2, k + 4} {
			for _, n := range []int{m, m + 1, m + 2, m + 4, m + 7} {
				for _, lda := range []int{n, n + 5} {
					dorgr2Test(t, impl, rnd, m, n, k, lda)
				}
			}
		}
	}
}

func dorgr2Test(t *testing.T, impl Dorgr2er, rnd *rand.Rand, m, n, k, lda int) {
	const tol = 1e-12
	name := fmt.Sprintf("m=%v,n=%v,k=%v,lda=%v", m, n, k, lda)

	if lda == 0 {
		lda = n
	}

	a := randomGeneral(m, n, lda, rnd)
	aCopy := cloneGeneral(a)
	// Compute the RQ decomposition of A.
	tau := make([]float64, m)
	work := make([]float64, 1)
	impl.Dgerqf(m, n, a.Data, a.Stride, tau, work, -1)
	work = make([]float64, int(work[0]))
	impl.Dgerqf(m, n, a.Data, a.Stride, tau, work, len(work))

	// Generate the upper triangular matrix R from subarray A[m-k:m, n-m:n]
	r := zeros(k, m, m)
	for i := 0; i < k; i++ {
		for j := 0; j < k; j++ {
			ia := i + a.Rows - k
			ja := j + a.Cols - k
			jr := j + r.Cols - k
			if i <= j {
				r.Data[i*r.Stride+jr] = a.Data[ia*a.Stride+ja]
			}
		}
	}
	// Compute the matrix Q using Dorg2r.
	impl.Dorgr2(m, n, k, a.Data, a.Stride, tau[m-k:m], work)
	if m == 0 {
		return
	}
	q := a
	// Test Q orthogonality.
	res := residualOrthogonal(q, true)
	if res > tol {
		t.Errorf("%v: |I - Q * Qᵀ| residual too large (%g)", name, res)
	}
	if k == 0 {
		return
	}

	// Reconstruct last rows of A.
	aRec := blas64.General{
		Rows:   k,
		Cols:   aCopy.Cols,
		Data:   aCopy.Data[(a.Rows-k)*lda:],
		Stride: lda,
	}
	// Test |A[m-k:m,0:n] - R[m-k:m,0:m] * Q| is small
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, r, q, -1, aRec)
	res = dlange(lapack.MaxColumnSum, aRec.Rows, aRec.Cols, aRec.Data, aRec.Stride)
	if res > tol {
		t.Errorf("%v: |A[m-k:m,0:n] - R[m-k:m,0:m] * Q| residual too large (%g)", name, res)
	}
}
