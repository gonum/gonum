// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/floats"
)

type Dgesc2er interface {
	Dgetc2er
	// Dgesc2 solves a system of linear equations
	//  A * X = scale * RHS
	// with a general N-by-N matrix A using the LU factorization with
	// complete pivoting computed by Dgetc2. The result is placed in
	// rhs on exit.
	Dgesc2(n int, a []float64, lda int, rhs []float64, ipiv, jpiv []int) (scale float64)
}

func Dgesc2Test(t *testing.T, impl Dgesc2er) {
	const tol = 1e-12
	rnd := rand.New(rand.NewSource(1))
	for _, test := range []struct {
		n, lda int
	}{
		{3, 0},
		{5, 0},
		{20, 30},
		{200, 0},
	} {
		testSolveDgesc2(t, impl, rnd, test.n, test.lda, tol)
	}
}

func testSolveDgesc2(t *testing.T, impl Dgesc2er, rnd *rand.Rand, n, lda int, tol float64) {
	name := fmt.Sprintf("n=%v,lda=%v", n, lda)
	if lda == 0 {
		lda = n
	}
	// Generate random general matrix.
	a := randomGeneral(n, n, lda, rnd)
	anorm := floats.Norm(a.Data, 1)

	// Generate a random solution.
	xWant := randomGeneral(n, 1, 1, rnd)
	xnorm := floats.Norm(xWant.Data, 1)

	// Compute RHS vector that solves for X such that  A*X = scale * RHS
	rhs := zeros(n, 1, 1)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, a, xWant, 1, rhs)
	rhsCopy := zeros(n, 1, 1) // Will contain A*x result.
	copyGeneral(rhsCopy, rhs)
	// Compute LU factorization with full pivoting.
	lu := zeros(n, n, lda)
	copyGeneral(lu, a)
	ipiv := make([]int, n)
	jpiv := make([]int, n)
	impl.Dgetc2(n, lu.Data, lu.Stride, ipiv, jpiv)

	// Solve using lu factorization.
	scale := impl.Dgesc2(lu.Rows, lu.Data, lu.Stride, rhs.Data, ipiv, jpiv)
	x := rhs
	if scale < 0 || scale > 1 {
		t.Errorf("%v: resulting scale out of bounds [0,1]", name)
	}

	var diff float64
	for i := range x.Data {
		diff = math.Max(diff, math.Abs(xWant.Data[i]-x.Data[i]))
	}
	if diff > tol {
		t.Errorf("%v: unexpected result, diff=%v", name, diff)
	}
	// |A*X - scale*RHS| / |A| / |X| is an indicator that solution is good
	// AxResult := zeros(n, 1, 1)
	// blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, a, x, 1, AxResult)
	// blas64.Scal(scale, blas64.Vector{N: n, Data: rhsCopy.Data, Inc: 1})
	// floats.Sub(AxResult.Data, rhsCopy.Data)

	// residualNorm := floats.Norm(rhsCopy.Data, 1) / anorm / xnorm
	// if residualNorm > tol {
	// 	t.Errorf("%v: |A*X - scale*RHS| / |A| / |X| = %g is greater than permissible tol", name, residualNorm)
	// }
}
