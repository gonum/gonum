// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dgghrder interface {
	Dgghrd(compq, compz lapack.OrthoComp, n, ilo, ihi int, a []float64, lda int, b []float64, ldb int, q []float64, ldq int, z []float64, ldz int)
}

func DgghrdTest(t *testing.T, impl Dgghrder) {
	rnd := rand.New(rand.NewPCG(1, 1))
	comps := []lapack.OrthoComp{lapack.OrthoExplicit, lapack.OrthoNone, lapack.OrthoPostmul}
	for _, compq := range comps {
		for _, compz := range comps {
			for _, n := range []int{0, 1, 2, 3, 4, 15} {
				for _, ld := range []int{max(1, n), n + 5} {
					testDgghrd(t, impl, rnd, compq, compz, n, 0, n-1, ld, ld, ld, ld)
				}
			}
		}
	}
}

func testDgghrd(t *testing.T, impl Dgghrder, rnd *rand.Rand, compq, compz lapack.OrthoComp, n, ilo, ihi, lda, ldb, ldq, ldz int) {
	const tol = 1e-13

	a := randomGeneral(n, n, lda, rnd)
	b := randomGeneral(n, n, ldb, rnd)

	var q, q1 blas64.General
	switch compq {
	case lapack.OrthoExplicit:
		// Initialize q to a non-orthogonal matrix, Dgghrd should overwrite it
		// with an orthogonal Q.
		q = randomGeneral(n, n, ldq, rnd)
	case lapack.OrthoPostmul:
		// Initialize q to an orthogonal matrix Q1, so that the result Q1*Q is
		// again orthogonal.
		q = randomOrthogonal(n, rnd)
		q1 = cloneGeneral(q)
	}

	var z, z1 blas64.General
	switch compz {
	case lapack.OrthoExplicit:
		z = randomGeneral(n, n, ldz, rnd)
	case lapack.OrthoPostmul:
		z = randomOrthogonal(n, rnd)
		z1 = cloneGeneral(z)
	}

	hGot := cloneGeneral(a)
	tGot := cloneGeneral(b)
	impl.Dgghrd(compq, compz, n, ilo, ihi, hGot.Data, hGot.Stride, tGot.Data, tGot.Stride, q.Data, max(1, q.Stride), z.Data, max(1, z.Stride))

	if n == 0 {
		return
	}

	name := fmt.Sprintf("Case compq=%v,compz=%v,n=%v,ilo=%v,ihi=%v", compq, compz, n, ilo, ihi)

	if !isUpperHessenberg(hGot) {
		t.Errorf("%v: H is not upper Hessenberg", name)
	}
	if !isUpperTriangular(tGot) {
		t.Errorf("%v: T is not upper triangular", name)
	}
	if compq != lapack.OrthoNone {
		if resid := residualOrthogonal(q, true); resid > tol {
			t.Errorf("%v: Q is not orthogonal, resid=%v", name, resid)
		}
	}
	if compz != lapack.OrthoNone {
		if resid := residualOrthogonal(z, true); resid > tol {
			t.Errorf("%v: Z is not orthogonal, resid=%v", name, resid)
		}
	}

	if compq != compz {
		// Verify reduction only when both Q and Z are computed.
		return
	}

	// Zero out the lower triangle of B.
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			b.Data[i*b.Stride+j] = 0
		}
	}

	aux := zeros(n, n, n)
	switch compq {
	case lapack.OrthoExplicit:
		// Qᵀ*A*Z = H
		hCalc := zeros(n, n, n)
		blas64.Gemm(blas.Trans, blas.NoTrans, 1, q, a, 0, aux)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, aux, z, 1, hCalc)
		if !equalApproxGeneral(hGot, hCalc, tol) {
			t.Errorf("%v: Qᵀ*A*Z != H", name)
		}

		// Qᵀ*B*Z = T
		tCalc := zeros(n, n, n)
		blas64.Gemm(blas.Trans, blas.NoTrans, 1, q, b, 0, aux)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, aux, z, 1, tCalc)
		if !equalApproxGeneral(tGot, tCalc, tol) {
			t.Errorf("%v: Qᵀ*B*Z != T", name)
		}
	case lapack.OrthoPostmul:
		//	Q1 * A * Z1ᵀ = (Q1*Q) * H * (Z1*Z)ᵀ
		lhs := zeros(n, n, n)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, q1, a, 0, aux)
		blas64.Gemm(blas.NoTrans, blas.Trans, 1, aux, z1, 0, lhs)

		rhs := zeros(n, n, n)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, q, hGot, 0, aux)
		blas64.Gemm(blas.NoTrans, blas.Trans, 1, aux, z, 0, rhs)
		if !equalApproxGeneral(lhs, rhs, tol) {
			t.Errorf("%v: Q1 * A * Z1ᵀ != (Q1*Q) * H * (Z1*Z)ᵀ", name)
		}

		//	Q1 * B * Z1ᵀ = (Q1*Q) * T * (Z1*Z)ᵀ
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, q1, b, 0, aux)
		blas64.Gemm(blas.NoTrans, blas.Trans, 1, aux, z1, 0, lhs)

		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, q, tGot, 0, aux)
		blas64.Gemm(blas.NoTrans, blas.Trans, 1, aux, z, 0, rhs)
		if !equalApproxGeneral(lhs, rhs, tol) {
			t.Errorf("%v: Q1 * B * Z1ᵀ != (Q1*Q) * T * (Z1*Z)ᵀ", name)
		}
	}
}
