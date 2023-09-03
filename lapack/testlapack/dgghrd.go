// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dgghrder interface {
	Dgghrd(compq, compz lapack.OrthoComp, n, ilo, ihi int, a []float64, lda int, b []float64, ldb int, q []float64, ldq int, z []float64, ldz int)
}

func DgghrdTest(t *testing.T, impl Dgghrder) {
	const tol = 1e-13
	const ldAdd = 5
	rnd := rand.New(rand.NewSource(1))
	comps := []lapack.OrthoComp{lapack.OrthoUnit, lapack.OrthoNone, lapack.OrthoEntry}
	for _, compq := range comps {
		for _, compz := range comps {
			for _, n := range []int{0, 1, 2, 4, 15} {
				ldMin := max(1, n)
				for _, lda := range []int{ldMin, ldMin + ldAdd} {
					for _, ldb := range []int{ldMin, ldMin + ldAdd} {
						for _, ldq := range []int{ldMin, ldMin + ldAdd} {
							for _, ldz := range []int{ldMin, ldMin + ldAdd} {
								testDgghrd(t, impl, rnd, tol, compq, compz, n, 0, n-1, lda, ldb, ldq, ldz)
							}
						}
					}
				}
			}
		}
	}
}

func testDgghrd(t *testing.T, impl Dgghrder, rnd *rand.Rand, tol float64, compq, compz lapack.OrthoComp, n, ilo, ihi, lda, ldb, ldq, ldz int) {
	a := randomGeneral(n, n, lda, rnd)
	b := blockedUpperTriGeneral(n, n, 0, n, ldb, false, rnd)

	var q, q1, z, z1 blas64.General
	if compq == lapack.OrthoEntry {
		q = randomOrthogonal(n, rnd)
		q1 = cloneGeneral(q)
	} else {
		q = nanGeneral(n, n, ldq)
	}
	if compz == lapack.OrthoEntry {
		z = randomOrthogonal(n, rnd)
		z1 = cloneGeneral(z)
	} else {
		z = nanGeneral(n, n, ldz)
	}

	hGot := cloneGeneral(a)
	tGot := cloneGeneral(b)
	impl.Dgghrd(compq, compz, n, ilo, ihi, hGot.Data, hGot.Stride, tGot.Data, tGot.Stride, q.Data, q.Stride, z.Data, z.Stride)
	if n == 0 {
		return
	}
	if !isUpperHessenberg(hGot) {
		t.Error("H is not upper Hessenberg")
	}
	if !isUpperTriangular(tGot) {
		t.Error("T is not upper triangular")
	}
	if compq == lapack.OrthoNone {
		if !isAllNaN(q.Data) {
			t.Errorf("Q is not NaN")
		}
		return
	}
	if compz == lapack.OrthoNone {
		if !isAllNaN(z.Data) {
			t.Errorf("Z is not NaN")
		}
		return
	}
	if compq != compz {
		return // Do not handle mixed case
	}
	comp := compq
	aux := zeros(n, n, n)

	switch comp {
	case lapack.OrthoUnit:
		// Qᵀ*A*Z = H
		hCalc := zeros(n, n, n)
		blas64.Gemm(blas.Trans, blas.NoTrans, 1, q, a, 0, aux)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, aux, z, 1, hCalc)
		if !equalApproxGeneral(hGot, hCalc, tol) {
			t.Errorf("Qᵀ*A*Z != H")
		}

		// Qᵀ*B*Z = T
		tCalc := zeros(n, n, n)
		blas64.Gemm(blas.Trans, blas.NoTrans, 1, q, b, 0, aux)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, aux, z, 1, tCalc)
		if !equalApproxGeneral(hGot, hCalc, tol) {
			t.Errorf("Qᵀ*B*Z != T")
		}
	case lapack.OrthoEntry:
		//	Q1 * A * Z1ᵀ = (Q1*Q) * H * (Z1*Z)ᵀ
		lhs := zeros(n, n, n)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, q1, a, 0, aux)
		blas64.Gemm(blas.NoTrans, blas.Trans, 1, aux, z1, 0, lhs) // lhs = Q1 * A * Z1ᵀ

		rhs := zeros(n, n, n)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, q, hGot, 0, aux)
		blas64.Gemm(blas.NoTrans, blas.Trans, 1, aux, z, 0, rhs)
		if !equalApproxGeneral(lhs, rhs, tol) {
			t.Errorf("Q1 * A * Z1ᵀ != (Q1*Q) * H * (Z1*Z)ᵀ")
		}

		//	Q1 * B * Z1ᵀ = (Q1*Q) * T * (Z1*Z)ᵀ
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, q1, b, 0, aux)
		blas64.Gemm(blas.NoTrans, blas.Trans, 1, aux, z1, 0, lhs)

		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, q, tGot, 0, aux)
		blas64.Gemm(blas.NoTrans, blas.Trans, 1, aux, z, 0, rhs)
		if !equalApproxGeneral(lhs, rhs, tol) {
			t.Errorf("Q1 * B * Z1ᵀ != (Q1*Q) * T * (Z1*Z)ᵀ")
		}
	}
}
