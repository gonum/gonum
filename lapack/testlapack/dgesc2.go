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
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/lapack"
)

type Dgesc2er interface {
	Dgesc2(n int, a []float64, lda int, rhs []float64, ipiv, jpiv []int) (scale float64)

	Dgetc2er
}

func Dgesc2Test(t *testing.T, impl Dgesc2er) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 20, 50} {
		for _, lda := range []int{n, n + 3} {
			testDgesc2(t, impl, rnd, n, lda, false)
			testDgesc2(t, impl, rnd, n, lda, true)
		}
	}
}

func testDgesc2(t *testing.T, impl Dgesc2er, rnd *rand.Rand, n, lda int, big bool) {
	const tol = 1e-14

	name := fmt.Sprintf("n=%v,lda=%v,big=%v", n, lda, big)

	// Generate random general matrix.
	a := randomGeneral(n, n, max(1, lda), rnd)

	// Generate a random right hand side vector.
	b := randomGeneral(n, 1, 1, rnd)
	if big {
		for i := 0; i < n; i++ {
			b.Data[i] *= bignum
		}
	}

	// Compute the LU factorization of A with full pivoting.
	lu := cloneGeneral(a)
	ipiv := make([]int, n)
	jpiv := make([]int, n)
	impl.Dgetc2(n, lu.Data, lu.Stride, ipiv, jpiv)

	// Make copies of const input to Dgesc2.
	luCopy := cloneGeneral(lu)
	ipivCopy := make([]int, len(ipiv))
	copy(ipivCopy, ipiv)
	jpivCopy := make([]int, len(jpiv))
	copy(jpivCopy, jpiv)

	// Call Dgesc2 to solve A*x = scale*b.
	x := cloneGeneral(b)
	scale := impl.Dgesc2(n, lu.Data, lu.Stride, x.Data, ipiv, jpiv)

	if n == 0 {
		return
	}

	// Check that const input to Dgesc2 hasn't been modified.
	if !floats.Same(lu.Data, luCopy.Data) {
		t.Errorf("%v: unexpected modification in lu", name)
	}
	if !intsEqual(ipiv, ipivCopy) {
		t.Errorf("%v: unexpected modification in ipiv", name)
	}
	if !intsEqual(jpiv, jpivCopy) {
		t.Errorf("%v: unexpected modification in jpiv", name)
	}

	if scale <= 0 || 1 < scale {
		t.Errorf("%v: scale %v out of bounds (0,1]", name, scale)
	}
	if !big && scale != 1 {
		t.Errorf("%v: unexpected scaling, scale=%v", name, scale)
	}

	// Compute the difference rhs := A*x - scale*b.
	diff := b
	for i := 0; i < n; i++ {
		diff.Data[i] *= scale
	}
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, a, x, -1, diff)

	// Compute the residual |A*x - scale*b| / |x|.
	xnorm := dlange(lapack.MaxColumnSum, n, 1, x.Data, 1)
	resid := dlange(lapack.MaxColumnSum, n, 1, diff.Data, 1) / xnorm
	if resid > tol || math.IsNaN(resid) {
		t.Errorf("%v: unexpected result; resid=%v, want<=%v", name, resid, tol)
	}
}
