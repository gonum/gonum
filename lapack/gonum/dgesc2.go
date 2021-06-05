// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas/blas64"
)

// Dgesc2 solves a system of linear equations
//  A * X = scale * RHS
// with a general N-by-N matrix A using the LU factorization with
// complete pivoting computed by Dgetc2. The result is placed in
// rhs on exit.
//
// Dgesc2 is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dgesc2(n int, a []float64, lda int, rhs []float64, ipiv, jpiv []int) (scale float64) {
	switch {
	case n < 0:
		panic(nLT0)
	case lda < max(1, n):
		panic(badLdA)
	}

	// Quick return if possible.
	if n == 0 {
		return 0
	}

	switch {
	case len(a) < (n-1)*lda+n:
		panic(shortA)
	case len(rhs) < n:
		panic(shortRHS)
	case len(ipiv) != n:
		panic(badLenIpiv)
	case len(jpiv) != n:
		panic(badLenJpiv)
	}

	const smlnum = dlamchS / dlamchP
	if len(a) < (n-1)*lda+n {
		panic(shortA)
	}

	// Apply permutations ipiv to RHS.
	impl.Dlaswp(1, rhs, 1, 0, n-1, ipiv[:n], 1)

	// Solve for L part.
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			rhs[j] -= float64(a[j*lda+i] * rhs[i])
		}
	}

	// Solve for U part.

	scale = 1.0

	// Check for scaling.
	bi := blas64.Implementation()
	i := bi.Idamax(n, rhs, 1)
	if 2*smlnum*math.Abs(rhs[i]) > math.Abs(a[(n-1)*lda+(n-1)]) {
		temp := 0.5 / math.Abs(rhs[i])
		bi.Dscal(n, temp, rhs, 1)
		scale *= temp
	}

	for i := n - 1; i >= 0; i-- {
		temp := 1.0 / a[i*lda+i]
		rhs[i] *= temp
		for j := i + 1; j < n; j++ {
			rhs[i] -= float64(rhs[j] * (a[i*lda+j] * temp))
		}
	}

	// Apply permutations jpiv to the solution (rhs).
	impl.Dlaswp(1, rhs, 1, 0, n-1, jpiv[:n], -1)
	return scale
}
