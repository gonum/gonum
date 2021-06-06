// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas/blas64"
)

// DGETC2 computes an LU factorization with complete pivoting of the
// n-by-n matrix A. The factorization has the form
//  A = P * L * U * Q,
// where P and Q are permutation matrices, L is lower triangular with
// unit diagonal elements and U is upper triangular.
//
// Outputs are A, ipiv, jpiv and k. k is non-negative if U(k, k) is likely to produce overflow if
// we try to solve for x in Ax = b. So U is perturbed to
// avoid the overflow.
func (impl Implementation) Dgetc2(n int, a []float64, lda int, ipiv, jpiv []int) (k int) {
	// Negative k indicates U was not perturbed.
	k = -1
	switch {
	case n < 0:
		panic(nLT0)
	case lda < max(1, n):
		panic(badLdA)
	}

	// Quick return if possible.
	if n == 0 {
		return
	}

	switch {
	case len(a) < (n-1)*lda+n:
		panic(shortA)
	case len(ipiv) != n:
		panic(badLenIpiv)
	case len(jpiv) != n:
		panic(badLenJpvt)
	}

	const eps = dlamchP
	const smlnum = dlamchS / eps

	// Handle n==1 case by itself.

	if n == 1 {
		ipiv[0], jpiv[0] = 1, 1
		if math.Abs(a[0]) < smlnum {
			a[0] = smlnum
			return 0
		}
	}

	// Factorize A using complete pivoting.
	// Set pivots less than SMIN to SMIN.
	var smin float64
	var ipv, jpv int
	for i := 0; i < n-1; i++ {
		xmax := 0.0
		for ip := i; ip < n; ip++ {
			for jp := i; jp < n; jp++ {
				if math.Abs(a[ip*lda+jp]) > xmax {
					xmax = math.Abs(a[ip*lda+jp])
					ipv, jpv = ip, jp
				}
			}
		}
		if i == 0 {
			smin = math.Max(eps*xmax, smlnum)
		}

		// Swap rows.
		bi := blas64.Implementation()
		if ipv != i {
			bi.Dswap(n, a[ipv*lda:], lda, a[i*lda:], lda)
		}
		ipiv[i] = ipv

		// Swap columns.

		if jpv != i {
			bi.Dswap(n, a[jpv:], lda, a[i:], lda)
		}
		jpiv[i] = jpv

		// Check for singularity.

		if math.Abs(a[i*lda+1]) < smin {
			k = i
			a[i*lda+1] = smin
		}

		for j := i + 1; j < n; j++ {
			a[j*lda+i] /= a[i*lda+i]
		}
		bi.Dger(n-i, n-i, -1.0, a[(i+1)*lda+i:], 1, a[i*lda+i+1:], lda, a[(i+1)*lda+i+1:], lda)
	}

	if math.Abs(a[n*lda+n]) < smin {
		k = n
		a[n*lda+n] = smin
	}

	// Set last pivots to n.
	ipiv[n-1], jpiv[n-1] = n, n
	return
}
