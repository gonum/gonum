// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dpttrser interface {
	Dpttrs(n, nrhs int, d, e []float64, b []float64, ldb int)

	Dpttrfer
}

func DpttrsTest(t *testing.T, impl Dpttrser) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 20, 50, 51, 52, 53, 54, 100} {
		for _, nrhs := range []int{0, 1, 2, 3, 4, 5, 10, 20, 50} {
			for _, ldb := range []int{max(1, nrhs), nrhs + 3} {
				dpttrsTest(t, impl, rnd, n, nrhs, ldb)
			}
		}
	}
}

func dpttrsTest(t *testing.T, impl Dpttrser, rnd *rand.Rand, n, nrhs, ldb int) {
	const tol = 1e-15

	name := fmt.Sprintf("n=%v", n)

	// Generate a random diagonally dominant symmetric tridiagonal matrix A.
	d, e := newRandomSymTridiag(n, rnd)

	// Make a copy of d and e to hold the factorization.
	dFac := make([]float64, len(d))
	copy(dFac, d)
	eFac := make([]float64, len(e))
	copy(eFac, e)

	// Compute the Cholesky factorization of A.
	ok := impl.Dpttrf(n, dFac, eFac)
	if !ok {
		t.Errorf("%v: bad test matrix, Dpttrf failed", name)
		return
	}

	// Generate a random solution matrix X.
	xWant := randomGeneral(n, nrhs, ldb, rnd)

	// Compute the right-hand side.
	b := zeros(n, nrhs, ldb)
	dstmm(n, nrhs, d, e, xWant.Data, xWant.Stride, b.Data, b.Stride)

	// Solve A*X=B.
	impl.Dpttrs(n, nrhs, dFac, eFac, b.Data, b.Stride)

	resid := dpttrsResidual(b, xWant)
	if resid > tol {
		t.Errorf("%v: unexpected solution: |diff| = %v, want <= %v", name, resid, tol)
	}
}

// dstmm computes the matrix-matrix product
//
//	C = A*B
//
// where A is an m×m symmetric tridiagonal matrix represented by the diagonal d
// and subdiagonal e, and B and C are m×n matrices.
func dstmm(m, n int, d, e []float64, b []float64, ldb int, c []float64, ldc int) {
	if m == 0 || n == 0 {
		return
	}
	if m == 1 {
		d0 := d[0]
		for j, b0j := range b[:n] {
			c[j] = d0 * b0j
		}
		return
	}
	for j := 0; j < n; j++ {
		c[j] = d[0]*b[j] + e[0]*b[ldb+j]
	}
	for i := 1; i < m-1; i++ {
		for j := 0; j < n; j++ {
			c[i*ldc+j] = e[i-1]*b[(i-1)*ldb+j] + d[i]*b[i*ldb+j] + e[i]*b[(i+1)*ldb+j]
		}
	}
	for j := 0; j < n; j++ {
		c[(m-1)*ldc+j] = e[m-2]*b[(m-2)*ldb+j] + d[m-1]*b[(m-1)*ldb+j]
	}
}

// dpttrsResidual returns |XGOT - XWANT|_1 / n.
func dpttrsResidual(xGot, xWant blas64.General) float64 {
	n, nrhs := xGot.Rows, xGot.Cols
	d := zeros(n, nrhs, nrhs)
	for i := 0; i < n; i++ {
		for j := 0; j < nrhs; j++ {
			d.Data[i*d.Stride+j] = xGot.Data[i*xGot.Stride+j] - xWant.Data[i*xWant.Stride+j]
		}
	}
	return dlange(lapack.MaxColumnSum, n, nrhs, d.Data, d.Stride) / float64(n)
}
