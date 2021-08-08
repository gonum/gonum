// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

// Dorgr2 generates an m×n real matrix Q with orthonormal rows,
// which is defined as the last m rows of a product of k elementary
// reflectors of order n
//  Q = H_0 * H_1 * ... * H_{k-1}
// as returned by Dgerqf.
//
// Each entry of tau contains the scalar factor of the elementary reflector H_i
// and the length of tau must be at least k. The length of work must be at least m.
// a is a matrix of dimensions (n,lda). On entry the [m-k+i-1]-th row must contain (counting from zero)
// the vector which defines the elementary reflector H_i, for i = 0,1,2,...,k-1, as
// returned by Dgerqf in the last k rows of its array argument A.
// On exit, the m×n matrix Q.
//  n >= m >= k >= 0
//
// Dorgr2 will panic if the conditions on input values are not met.
//
// Dorgr2 is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dorgr2(m, n, k int, a []float64, lda int, tau, work []float64) {
	switch {
	case m < 0:
		panic(mLT0)
	case n < 0:
		panic(nLT0)
	case m > n:
		panic(mGTN)
	case k < 0:
		panic(kLT0)
	case k > m:
		panic(kGTM)
	case lda < max(1, n):
		panic(badLdA)
	case len(work) < m:
		panic(shortWork)
	}

	// Quick return if possible.
	if m == 0 {
		return
	}

	switch {
	case len(tau) != k:
		panic(badLenTau)
	case len(a) < (m-1)*lda+n:
		panic(shortA)
	}

	if k < m {
		// Initialise rows 0:m-k to rows of the unit matrix.
		for l := 0; l < m-k; l++ {
			for j := 0; j < n; j++ {
				a[l*lda+j] = 0
			}
			a[l*lda+n-m+l] = 1
		}
	}
	bi := blas64.Implementation()
	for i := 0; i < k; i++ {
		ii := m - k + i

		// Apply H_i to A[0:m-k+i+1, 0:n-k+i+1] from the right.
		a[ii*lda+n-m+ii] = 1
		impl.Dlarf(blas.Right, ii, n-m+ii+1, a[ii*lda:], 1, tau[i], a, lda, work)
		bi.Dscal(n-m+ii, -tau[i], a[ii*lda:], 1)
		a[ii*lda+n-m+ii] = 1 - tau[i]

		// Set A[m-k+i, n-k+i:n] to zero.
		for l := n - m + ii + 1; l < n; l++ {
			a[ii*lda+l] = 0
		}
	}
}
