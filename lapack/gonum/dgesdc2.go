// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas/blas64"
)

// Dgesdc2 solves a system of linear equations
//   A * X = scale* RHS
// with a general N-by-N matrix A using the LU factorization with
// complete pivoting computed by Dgetc2.
func (impl Implementation) Dgesdc2(n int, a []float64, lda int, rhs []float64, ipiv, jpiv []int, scale float64) {
	const eps = dlamchP
	var smlnum, _ = impl.Dlabad(dlamchS/eps, 1./(dlamchS/eps))
	if len(a) < (n-1)*lda+n {
		panic(shortA)
	}
	//Apply permutations IPIV to RHS
	impl.Dlaswp(1, rhs, lda, 1, n-1, ipiv, 1)

	// solve for L part
	for i := 0; i < n; i++ {
		// of course this can be optimized
		// I'm still struggling with this indexing.
		for j := i + 1; j < n; j++ {
			rhs[j] = rhs[j] - a[(i*n)+j]*rhs[i]
		}
	}
	scale = 1.
	//
	i := impl.Idamax(n, rhs, 1)
	if 2*smlnum*math.Abs(rhs[i]) > math.Abs(a[n*n+n]) {
		temp := 0.5 / math.Abs(rhs[i])
		blas64.Implementation().Dscal(n, temp, rhs[1:], 1) // what?
		scale = scale * temp
	}
	for i := n; i > 0; i-- {
		temp := 1 / a[(i*n)+i]
		rhs[i] = rhs[i] * temp
		for j := i + 1; j < n; j++ {
			rhs[i] = rhs[i] - rhs[j]*a[(i*n)+j]*temp
		}
	}
	impl.Dlaswp(1, rhs, lda, 1, n-1, jpiv, -1)
}
