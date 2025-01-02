// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dgesver interface {
	Dgesv(n, nrhs int, a []float64, lda int, ipiv []int, b []float64, ldb int) bool

	Dgetri(n int, a []float64, lda int, ipiv []int, work []float64, lwork int) bool
}

func DgesvTest(t *testing.T, impl Dgesver) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 50, 100} {
		for _, nrhs := range []int{0, 1, 2, 5} {
			for _, lda := range []int{max(1, n), n + 5} {
				for _, ldb := range []int{max(1, nrhs), nrhs + 5} {
					dgesvTest(t, impl, rnd, n, nrhs, lda, ldb)
				}
			}
		}
	}
}

func dgesvTest(t *testing.T, impl Dgesver, rnd *rand.Rand, n, nrhs, lda, ldb int) {
	const tol = 1e-15

	name := fmt.Sprintf("n=%v,nrhs=%v,lda=%v,ldb=%v", n, nrhs, lda, ldb)

	// Create a random system matrix A and the solution X.
	a := randomGeneral(n, n, lda, rnd)
	xWant := randomGeneral(n, nrhs, ldb, rnd)

	// Compute the right hand side matrix B = A*X.
	b := zeros(n, nrhs, ldb)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, a, xWant, 0, b)

	// Allocate a slice for row pivots and fill it with invalid indices.
	ipiv := make([]int, n)
	for i := range ipiv {
		ipiv[i] = -1
	}

	// Call Dgesv to solve A*X = B.
	lu := cloneGeneral(a)
	xGot := cloneGeneral(b)
	ok := impl.Dgesv(n, nrhs, lu.Data, lu.Stride, ipiv, xGot.Data, xGot.Stride)

	if !ok {
		t.Errorf("%v: unexpected failure in Dgesv", name)
		return
	}

	if n == 0 || nrhs == 0 {
		return
	}

	// Check that all elements of ipiv have been set.
	ipivSet := true
	for _, ipv := range ipiv {
		if ipv == -1 {
			ipivSet = false
			break
		}
	}
	if !ipivSet {
		t.Fatalf("%v: not all elements of ipiv set", name)
		return
	}

	// Compute the reciprocal of the condition number of A from its LU
	// decomposition before it's overwritten further below.
	aInv := cloneGeneral(lu)
	impl.Dgetri(n, aInv.Data, aInv.Stride, ipiv, make([]float64, n), n)
	ainvnorm := dlange(lapack.MaxColumnSum, n, n, aInv.Data, aInv.Stride)
	anorm := dlange(lapack.MaxColumnSum, n, n, a.Data, a.Stride)
	rcond := 1 / anorm / ainvnorm

	// Reconstruct matrix A from factors and compute residual.
	//
	// Extract L and U from lu.
	l := zeros(n, n, n)
	u := zeros(n, n, n)
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			l.Data[i*l.Stride+j] = lu.Data[i*lu.Stride+j]
		}
		l.Data[i*l.Stride+i] = 1
		for j := i; j < n; j++ {
			u.Data[i*u.Stride+j] = lu.Data[i*lu.Stride+j]
		}
	}
	// Compute L*U.
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, l, u, 0, lu)
	// Apply P to L*U.
	for i := n - 1; i >= 0; i-- {
		ip := ipiv[i]
		if ip == i {
			continue
		}
		row1 := blas64.Vector{N: n, Data: lu.Data[i*lu.Stride:], Inc: 1}
		row2 := blas64.Vector{N: n, Data: lu.Data[ip*lu.Stride:], Inc: 1}
		blas64.Swap(row1, row2)
	}
	// Compute P*L*U - A.
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			lu.Data[i*lu.Stride+j] -= a.Data[i*a.Stride+j]
		}
	}
	// Compute the residual |P*L*U - A|.
	resid := dlange(lapack.MaxColumnSum, n, n, lu.Data, lu.Stride)
	resid /= float64(n) * anorm
	if resid > tol || math.IsNaN(resid) {
		t.Errorf("%v: residual |P*L*U - A| is too large, got %v, want <= %v", name, resid, tol)
	}

	// Compute residual of the computed solution.
	//
	// Compute B - A*X.
	blas64.Gemm(blas.NoTrans, blas.NoTrans, -1, a, xGot, 1, b)
	// Compute the maximum over the number of right hand sides of |B - A*X| / (|A| * |X|).
	resid = 0
	for j := 0; j < nrhs; j++ {
		bnorm := blas64.Asum(blas64.Vector{N: n, Data: b.Data[j:], Inc: b.Stride})
		xnorm := blas64.Asum(blas64.Vector{N: n, Data: xGot.Data[j:], Inc: xGot.Stride})
		resid = math.Max(resid, bnorm/anorm/xnorm)
	}
	if resid > tol || math.IsNaN(resid) {
		t.Errorf("%v: residual |B - A*X| is too large, got %v, want <= %v", name, resid, tol)
	}

	// Compare the computed solution with the generated exact solution.
	//
	// Compute X - XWANT.
	for i := 0; i < n; i++ {
		for j := 0; j < nrhs; j++ {
			xGot.Data[i*xGot.Stride+j] -= xWant.Data[i*xWant.Stride+j]
		}
	}
	// Compute the maximum of |X - XWANT|/|XWANT| over all the vectors X and XWANT.
	resid = 0
	for j := 0; j < nrhs; j++ {
		xnorm := dlange(lapack.MaxAbs, n, 1, xWant.Data[j:], xWant.Stride)
		diff := dlange(lapack.MaxAbs, n, 1, xGot.Data[j:], xGot.Stride)
		resid = math.Max(resid, diff/xnorm*rcond)
	}
	if resid > tol || math.IsNaN(resid) {
		t.Errorf("%v: residual |X-XWANT| is too large, got %v, want <= %v", name, resid, tol)
	}
}
