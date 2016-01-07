// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
	"github.com/gonum/lapack"
)

type Dgesvder interface {
	Dgesvd(jobU, jobVT lapack.SVDJob, m, n int, a []float64, lda int, s, u []float64, ldu int, vt []float64, ldvt int, work []float64, lwork int) (ok bool)
}

func DgesvdTest(t *testing.T, impl Dgesvder) {
	// TODO(btracey): Add tests for all of the cases when the SVD implementation
	// is finished.
	// TODO(btracey): Add tests for m > mnthr and n > mnthr when other SVD
	// conditions are implemented. Right now mnthr is 5,000,000 which is too
	// large to create a square matrix of that size.
	for _, jobU := range []lapack.SVDJob{lapack.SVDAll} {
		for _, jobVT := range []lapack.SVDJob{lapack.SVDAll} {
			for _, test := range []struct {
				m, n, lda, ldu, ldvt int
			}{
				{5, 5, 0, 0, 0},
				{5, 7, 0, 0, 0},
				{7, 5, 0, 0, 0},

				{5, 5, 10, 11, 12},
				{5, 7, 10, 11, 12},
				{7, 5, 10, 11, 12},
			} {
				m := test.m
				n := test.n
				lda := test.lda
				if lda == 0 {
					lda = n
				}
				ldu := test.ldu
				if ldu == 0 {
					ldu = m
				}
				ldvt := test.ldvt
				if ldvt == 0 {
					ldvt = n
				}

				a := make([]float64, m*lda)
				for i := range a {
					a[i] = rand.NormFloat64()
				}

				u := make([]float64, m*ldu)
				for i := range u {
					u[i] = rand.NormFloat64()
				}

				vt := make([]float64, n*ldvt)
				for i := range vt {
					vt[i] = rand.NormFloat64()
				}

				aCopy := make([]float64, len(a))
				copy(aCopy, a)

				s := make([]float64, min(m, n))

				work := make([]float64, 1)
				impl.Dgesvd(jobU, jobVT, m, n, a, lda, s, u, ldu, vt, ldvt, work, -1)

				work = make([]float64, int(work[0]))
				impl.Dgesvd(jobU, jobVT, m, n, a, lda, s, u, ldu, vt, ldvt, work, len(work))

				// Test the decomposition
				sigma := blas64.General{
					Rows:   m,
					Cols:   n,
					Stride: n,
					Data:   make([]float64, m*n),
				}
				for i := 0; i < min(m, n); i++ {
					sigma.Data[i*sigma.Stride+i] = s[i]
				}

				uMat := blas64.General{
					Rows:   m,
					Cols:   m,
					Stride: ldu,
					Data:   u,
				}
				vTMat := blas64.General{
					Rows:   n,
					Cols:   n,
					Stride: ldvt,
					Data:   vt,
				}

				tmp := blas64.General{
					Rows:   m,
					Cols:   n,
					Stride: n,
					Data:   make([]float64, m*n),
				}
				ans := blas64.General{
					Rows:   m,
					Cols:   n,
					Stride: lda,
					Data:   make([]float64, m*lda),
				}
				copy(ans.Data, a)

				blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, uMat, sigma, 0, tmp)
				blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, tmp, vTMat, 0, ans)

				errStr := fmt.Sprintf("jobU = %v, jobVT = %v, m = %v, n = %v, lda = %v, ldu = %v, ldv = %v", jobU, jobVT, m, n, lda, ldu, ldvt)
				if !floats.EqualApprox(ans.Data, aCopy, 1e-8) {
					t.Errorf("Decomposition mismatch %s", errStr)
				}

				// Check that U and V are orthogonal
				for i := 0; i < uMat.Rows; i++ {
					for j := i + 1; j < uMat.Rows; j++ {
						dot := blas64.Dot(uMat.Cols,
							blas64.Vector{Inc: 1, Data: uMat.Data[i*uMat.Stride:]},
							blas64.Vector{Inc: 1, Data: uMat.Data[j*uMat.Stride:]},
						)
						if dot > 1e-8 {
							t.Errorf("U not orthogonal %s", errStr)
						}
					}
				}
				for i := 0; i < vTMat.Rows; i++ {
					for j := i + 1; j < vTMat.Rows; j++ {
						dot := blas64.Dot(vTMat.Cols,
							blas64.Vector{Inc: 1, Data: vTMat.Data[i*vTMat.Stride:]},
							blas64.Vector{Inc: 1, Data: vTMat.Data[j*vTMat.Stride:]},
						)
						if dot > 1e-8 {
							t.Errorf("V not orthogonal %s", errStr)
						}
					}
				}
			}
		}
	}
}
