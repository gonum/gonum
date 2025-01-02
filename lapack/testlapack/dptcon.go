// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/lapack"
)

type Dptconer interface {
	Dptcon(n int, d, e []float64, anorm float64, work []float64) (rcond float64)

	Dpttrf(n int, d, e []float64) (ok bool)
	Dpttrs(n, nrhs int, d, e []float64, b []float64, ldb int)
}

func DptconTest(t *testing.T, impl Dptconer) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 20, 50, 51, 52, 53, 54, 100} {
		dptconTest(t, impl, rnd, n)
	}
}

func dptconTest(t *testing.T, impl Dptconer, rnd *rand.Rand, n int) {
	const tol = 1e-15

	name := fmt.Sprintf("n=%v", n)

	// Generate a random diagonally dominant symmetric tridiagonal matrix A.
	d, e := newRandomSymTridiag(n, rnd)
	aNorm := dlanst(lapack.MaxColumnSum, n, d, e)

	// Compute the Cholesky factorization of A.
	ok := impl.Dpttrf(n, d, e)
	if !ok {
		t.Errorf("%v: bad test matrix, Dpttrf failed", name)
		return
	}

	// Compute the reciprocal of the condition number of A.
	dCopy := make([]float64, len(d))
	copy(dCopy, d)
	eCopy := make([]float64, len(e))
	copy(eCopy, e)
	work := make([]float64, 3*n)
	rcondGot := impl.Dptcon(n, d, e, aNorm, work)

	// Check that Dptcon didn't modify d and e.
	if !floats.Equal(d, dCopy) {
		t.Errorf("%v: unexpected modification of d", name)
	}
	if !floats.Equal(e, eCopy) {
		t.Errorf("%v: unexpected modification of e", name)
	}

	// Compute the norm of A⁻¹.
	aInv, lda := make([]float64, n*n), max(1, n)
	for i := 0; i < n; i++ {
		aInv[i*lda+i] = 1
	}
	impl.Dpttrs(n, n, d, e, aInv, lda)
	aInvNorm := dlange(lapack.MaxColumnSum, n, n, aInv, lda)

	rcondWant := 1.0
	if aNorm > 0 && aInvNorm > 0 {
		rcondWant = 1 / aNorm / aInvNorm
	}

	diff := math.Abs(rcondGot - rcondWant)
	if diff > tol {
		t.Errorf("%v: unexpected value of rcond. got=%v, want=%v (diff=%v)", name, rcondGot, rcondWant, diff)
	}
}
