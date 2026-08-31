// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import "gonum.org/v1/gonum/blas/blas64"

// Dtzrzf reduces the m×n (m ≤ n) upper trapezoidal matrix A to upper
// triangular form by means of orthogonal transformations.
//
// The upper trapezoidal matrix A is factored as
//
//	A = [R 0] * Z
//
// where Z is an n×n orthogonal matrix and R is an m×m upper triangular matrix.
//
// Z is stored as a product of m elementary reflectors
//
//	Z = Z(0) * Z(1) * ... * Z(m-1)
//
// Each Z(i) has the form
//
//	Z(i) = I - tau[i] * v * vᵀ
//
// where tau[i] is stored in tau[i] and v is a vector with
//
//	v[0:i] = 0, v[i] = 1, v[i+1:m] = 0, v[m:n] = A[i, m:n]
//
// on exit.
//
// tau must have length at least m. work must have length at least max(1,lwork).
//
// If lwork is -1, instead of performing Dtzrzf, only the optimal workspace
// size is stored into work[0].
//
// Dtzrzf is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dtzrzf(m, n int, a []float64, lda int, tau, work []float64, lwork int) {
	switch {
	case m < 0:
		panic(mLT0)
	case n < m:
		panic(nLTM)
	case lda < max(1, n):
		panic(badLdA)
	case lwork < max(1, m) && lwork != -1:
		panic(badLWork)
	case len(work) < max(1, lwork):
		panic(shortWork)
	}

	if lwork == -1 {
		work[0] = float64(m)
		return
	}

	if m == 0 {
		work[0] = 1
		return
	}

	switch {
	case len(a) < (m-1)*lda+n:
		panic(shortA)
	case len(tau) < m:
		panic(shortTau)
	}

	if m == n {
		for i := range tau[:m] {
			tau[i] = 0
		}
		work[0] = 1
		return
	}

	bi := blas64.Implementation()
	l := n - m
	for i := m - 1; i >= 0; i-- {
		beta, t := impl.Dlarfg(l+1, a[i*lda+i], a[i*lda+m:i*lda+n], 1)
		tau[i] = t
		a[i*lda+i] = beta

		if t != 0 && i > 0 {
			// Apply Z(i) to A[0:i, i:n] from the right.
			// w = A[0:i, i] + A[0:i, m:n] * v
			bi.Dcopy(i, a[i:], lda, work, 1)
			bi.Dgemv('N', i, l, 1, a[m:], lda, a[i*lda+m:], 1, 1, work, 1)
			// A[0:i, i] -= tau * w
			bi.Daxpy(i, -t, work, 1, a[i:], lda)
			// A[0:i, m:n] -= tau * w * v^T
			bi.Dger(i, l, -t, work, 1, a[i*lda+m:], 1, a[m:], lda)
		}
	}
	work[0] = float64(m)
}
