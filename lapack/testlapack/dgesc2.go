// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package testlapack

import (
	"testing"

	"golang.org/x/exp/rand"
)

type Dgesc2er interface {
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
		{10, 0},
		// {5, 0},
		// {10, 0},
		// {300, 0},
		// {3, 0},
		// {200, 0},
		// {300, 0},
		// {204, 0},
		// {1, 0},
		// {3000, 0},
	} {
		n := test.n
		lda := test.lda
		if lda == 0 {
			lda = n
		}
		// Generate a random square matrix A with elements uniformly in [-1,1).
		a := make([]float64, max(0, (n-1)*lda+n))
		for i := range a {
			a[i] = 2*rnd.Float64() - 1
		}

		// Create ipiv and jpiv indices
		ipiv := make([]int, n)
		jpiv := make([]int, n)
		for i := range ipiv {
			ipiv[i] = i
			jpiv[i] = i
		}
		// a := randomGeneral(n, n, n+extra, rnd)
		// Store a copy of A for later comparison.
		aCopy := make([]float64, len(a))
		copy(aCopy, a)

		// Allocate a slice for scalar rhs (b in equation A*x = scale * b)
		b := make([]float64, n)
		for i := 0; i < n; i++ {
			b[i] = rnd.NormFloat64()
		}

		// Compute the expected result
		want := make([]float64, len(a))
		copy(want, a)
		scale := impl.Dgesc2(n, a, lda, b, ipiv, jpiv)
		if scale <= 0. || scale > 100 {
			t.Errorf("resulting scale out of bounds (0,100]")
		}
	}
}
