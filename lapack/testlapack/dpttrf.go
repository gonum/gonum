// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/lapack"
)

type Dpttrfer interface {
	Dpttrf(n int, d, e []float64) (ok bool)
}

// DpttrfTest tests a tridiagonal Cholesky factorization on random symmetric
// positive definite tridiagonal matrices by checking that the Cholesky factors
// multiply back to the original matrix.
func DpttrfTest(t *testing.T, impl Dpttrfer) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 20, 50, 51, 52, 53, 54, 100} {
		dpttrfTest(t, impl, rnd, n)
	}
}

func dpttrfTest(t *testing.T, impl Dpttrfer, rnd *rand.Rand, n int) {
	const tol = 1e-15

	name := fmt.Sprintf("n=%v", n)

	// Generate a random diagonally dominant symmetric tridiagonal matrix A.
	d, e := newRandomSymTridiag(n, rnd)

	// Make a copy of d and e to hold the factorization.
	var dFac, eFac []float64
	if n > 0 {
		dFac = make([]float64, len(d))
		copy(dFac, d)
		if n > 1 {
			eFac = make([]float64, len(e))
			copy(eFac, e)
		}
	}

	// Compute the Cholesky factorization of A.
	ok := impl.Dpttrf(n, dFac, eFac)
	if !ok {
		t.Errorf("%v: bad test matrix, Dpttrf failed", name)
		return
	}

	// Check the residual norm(L*D*Lᵀ - A)/(n * norm(A)).
	resid := dpttrfResidual(n, d, e, dFac, eFac)
	if resid > tol {
		t.Errorf("%v: unexpected residual |L*D*Lᵀ - A|/(n * norm(A)); got %v, want <= %v", name, resid, tol)
	}
}

func dpttrfResidual(n int, d, e, dFac, eFac []float64) float64 {
	if n == 0 {
		return 0
	}

	// Construct the difference L*D*Lᵀ - A.
	dDiff := make([]float64, n)
	eDiff := make([]float64, n-1)
	dDiff[0] = dFac[0] - d[0]
	for i, ef := range eFac {
		de := dFac[i] * ef
		dDiff[i+1] = de*ef + dFac[i+1] - d[i+1]
		eDiff[i] = de - e[i]
	}

	// Compute the 1-norm of the difference L*D*Lᵀ - A.
	var resid float64
	if n == 1 {
		resid = math.Abs(dDiff[0])
	} else {
		resid = math.Max(math.Abs(dDiff[0])+math.Abs(eDiff[0]), math.Abs(dDiff[n-1])+math.Abs(eDiff[n-2]))
		for i := 1; i < n-1; i++ {
			resid = math.Max(resid, math.Abs(dDiff[i])+math.Abs(eDiff[i-1])+math.Abs(eDiff[i]))
		}
	}

	anorm := dlanst(lapack.MaxColumnSum, n, d, e)

	// Compute norm(L*D*Lᵀ - A)/(n * norm(A)).
	if anorm == 0 {
		if resid != 0 {
			return math.Inf(1)
		}
		return 0
	}
	return resid / float64(n) / anorm
}

func newRandomSymTridiag(n int, rnd *rand.Rand) (d, e []float64) {
	if n == 0 {
		return nil, nil
	}

	if n == 1 {
		d = make([]float64, 1)
		d[0] = rnd.Float64()
		return d, nil
	}

	// Allocate the diagonal d and fill it with numbers from [0,1).
	d = make([]float64, n)
	dlarnv(d, 1, rnd)
	// Allocate the subdiagonal e and fill it with numbers from [-1,1).
	e = make([]float64, n-1)
	dlarnv(e, 2, rnd)

	// Make A diagonally dominant by adding the absolute value of off-diagonal
	// elements to it.
	d[0] += math.Abs(e[0])
	for i := 1; i < n-1; i++ {
		d[i] += math.Abs(e[i]) + math.Abs(e[i-1])
	}
	d[n-1] += math.Abs(e[n-2])

	return d, e
}
