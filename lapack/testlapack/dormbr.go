// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math/rand"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/lapack"
)

type Dormbrer interface {
	Dormbr(vect lapack.DecompUpdate, side blas.Side, trans blas.Transpose, m, n, k int, a []float64, lda int, tau, c []float64, ldc int, work []float64, lwork int)
	Dgebrder
}

func DormbrTest(t *testing.T, impl Dormbrer) {
	rnd := rand.New(rand.NewSource(1))
	bi := blas64.Implementation()
	for _, vect := range []lapack.DecompUpdate{lapack.ApplyQ, lapack.ApplyP} {
		for _, side := range []blas.Side{blas.Left, blas.Right} {
			for _, trans := range []blas.Transpose{blas.NoTrans, blas.Trans} {
				for _, wl := range []worklen{minimumWork, mediumWork, optimumWork} {
					for _, test := range []struct {
						m, n, k, lda, ldc int
					}{
						{3, 4, 5, 0, 0},
						{3, 5, 4, 0, 0},
						{4, 3, 5, 0, 0},
						{4, 5, 3, 0, 0},
						{5, 3, 4, 0, 0},
						{5, 4, 3, 0, 0},

						{3, 4, 5, 10, 12},
						{3, 5, 4, 10, 12},
						{4, 3, 5, 10, 12},
						{4, 5, 3, 10, 12},
						{5, 3, 4, 10, 12},
						{5, 4, 3, 10, 12},

						{150, 140, 130, 0, 0},
					} {
						m := test.m
						n := test.n
						k := test.k
						ldc := test.ldc
						if ldc == 0 {
							ldc = n
						}
						nq := n
						nw := m
						if side == blas.Left {
							nq = m
							nw = n
						}

						// Compute a decomposition.
						var ma, na int
						var a []float64
						if vect == lapack.ApplyQ {
							ma = nq
							na = k
						} else {
							ma = k
							na = nq
						}
						lda := test.lda
						if lda == 0 {
							lda = na
						}
						a = make([]float64, ma*lda)
						for i := range a {
							a[i] = rnd.NormFloat64()
						}
						nTau := min(nq, k)
						tauP := make([]float64, nTau)
						tauQ := make([]float64, nTau)
						d := make([]float64, nTau)
						e := make([]float64, nTau)

						work := make([]float64, 1)
						impl.Dgebrd(ma, na, a, lda, d, e, tauQ, tauP, work, -1)
						work = make([]float64, int(work[0]))
						impl.Dgebrd(ma, na, a, lda, d, e, tauQ, tauP, work, len(work))

						// Apply and compare update.
						c := make([]float64, m*ldc)
						for i := range c {
							c[i] = rnd.NormFloat64()
						}
						cCopy := make([]float64, len(c))
						copy(cCopy, c)

						var lwork int
						switch wl {
						case minimumWork:
							lwork = nw
						case optimumWork:
							impl.Dormbr(vect, side, trans, m, n, k, a, lda, tauQ, c, ldc, work, -1)
							lwork = int(work[0])
						case mediumWork:
							work := make([]float64, 1)
							impl.Dormbr(vect, side, trans, m, n, k, a, lda, tauQ, c, ldc, work, -1)
							lwork = (int(work[0]) + nw) / 2
						}
						lwork = max(1, lwork)
						work = make([]float64, lwork)

						if vect == lapack.ApplyQ {
							impl.Dormbr(vect, side, trans, m, n, k, a, lda, tauQ, c, ldc, work, lwork)
						} else {
							impl.Dormbr(vect, side, trans, m, n, k, a, lda, tauP, c, ldc, work, lwork)
						}

						// Check that the multiplication was correct.
						cOrig := blas64.General{
							Rows:   m,
							Cols:   n,
							Stride: ldc,
							Data:   make([]float64, len(cCopy)),
						}
						copy(cOrig.Data, cCopy)
						cAns := blas64.General{
							Rows:   m,
							Cols:   n,
							Stride: ldc,
							Data:   make([]float64, len(cCopy)),
						}
						copy(cAns.Data, cCopy)
						nb := min(ma, na)
						var mulMat blas64.General
						if vect == lapack.ApplyQ {
							mulMat = constructQPBidiagonal(lapack.ApplyQ, ma, na, nb, a, lda, tauQ)
						} else {
							mulMat = constructQPBidiagonal(lapack.ApplyP, ma, na, nb, a, lda, tauP)
						}

						mulTrans := trans

						if side == blas.Left {
							bi.Dgemm(mulTrans, blas.NoTrans, m, n, m, 1, mulMat.Data, mulMat.Stride, cOrig.Data, cOrig.Stride, 0, cAns.Data, cAns.Stride)
						} else {
							bi.Dgemm(blas.NoTrans, mulTrans, m, n, n, 1, cOrig.Data, cOrig.Stride, mulMat.Data, mulMat.Stride, 0, cAns.Data, cAns.Stride)
						}

						if !floats.EqualApprox(cAns.Data, c, 1e-13) {
							isApplyQ := vect == lapack.ApplyQ
							isLeft := side == blas.Left
							isTrans := trans == blas.Trans

							t.Errorf("C mismatch. isApplyQ: %v, isLeft: %v, isTrans: %v, m = %v, n = %v, k = %v, lda = %v, ldc = %v",
								isApplyQ, isLeft, isTrans, m, n, k, lda, ldc)
						}
					}
				}
			}
		}
	}
}
