// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/cmplx"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/cblas128"
)

type Zgetrser interface {
	Zgetrfer
	Zgetrs(trans blas.Transpose, n, nrhs int, a []complex128, lda int, ipiv []int, b []complex128, ldb int)
}

func ZgetrsTest(t *testing.T, impl Zgetrser) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, trans := range []blas.Transpose{blas.NoTrans, blas.Trans, blas.ConjTrans} {
		for _, test := range []struct {
			n, nrhs, lda, ldb int
			tol               float64
		}{
			{1, 1, 0, 0, 1e-12},
			{2, 1, 0, 0, 1e-12},
			{3, 3, 0, 0, 1e-12},
			{3, 5, 0, 0, 1e-12},
			{5, 3, 0, 0, 1e-12},
			{16, 4, 0, 0, 1e-10},
			{64, 8, 0, 0, 1e-9},
			{3, 3, 8, 10, 1e-12},
			{3, 5, 8, 10, 1e-12},
			{5, 3, 8, 10, 1e-12},
			{100, 20, 0, 0, 1e-8},
			{100, 20, 130, 40, 1e-8},
		} {
			n := test.n
			nrhs := test.nrhs
			lda := test.lda
			if lda == 0 {
				lda = max(1, n)
			}
			ldb := test.ldb
			if ldb == 0 {
				ldb = max(1, nrhs)
			}

			// Build a well-conditioned A = I + small random perturbation.
			a := make([]complex128, n*lda)
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					a[i*lda+j] = complex(0.1*rnd.NormFloat64(), 0.1*rnd.NormFloat64())
				}
				a[i*lda+i] += complex(float64(n), 0)
			}
			b := make([]complex128, n*ldb)
			for i := 0; i < n; i++ {
				for j := 0; j < nrhs; j++ {
					b[i*ldb+j] = complex(rnd.NormFloat64(), rnd.NormFloat64())
				}
			}
			aCopy := make([]complex128, len(a))
			copy(aCopy, a)
			bCopy := make([]complex128, len(b))
			copy(bCopy, b)

			ipiv := make([]int, n)

			// Factor A.
			impl.Zgetrf(n, n, a, lda, ipiv)
			// Solve using factorization.
			impl.Zgetrs(trans, n, nrhs, a, lda, ipiv, b, ldb)

			// Check that op(A)*X == bCopy where op = trans.
			A := cblas128.General{Rows: n, Cols: n, Stride: lda, Data: aCopy}
			X := cblas128.General{Rows: n, Cols: nrhs, Stride: ldb, Data: b}
			B := cblas128.General{Rows: n, Cols: nrhs, Stride: ldb, Data: make([]complex128, n*ldb)}
			cblas128.Gemm(trans, blas.NoTrans, 1, A, X, 0, B)

			name := fmt.Sprintf("trans=%v,n=%v,nrhs=%v,lda=%v,ldb=%v", trans, n, nrhs, lda, ldb)
			for i := 0; i < n; i++ {
				for j := 0; j < nrhs; j++ {
					diff := B.Data[i*ldb+j] - bCopy[i*ldb+j]
					if cmplx.Abs(diff) > test.tol {
						t.Errorf("%s: solve mismatch at (%d,%d): residual %v > tol %v", name, i, j, cmplx.Abs(diff), test.tol)
						break
					}
				}
			}
		}
	}
}
