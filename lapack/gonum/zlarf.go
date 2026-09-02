// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/cblas128"
)

// Zlarf applies a complex elementary reflector H to an m×n complex matrix C:
//
//	C = H * C  if side == blas.Left
//	C = C * H  if side == blas.Right
//
// H is represented in the form
//
//	H = I - tau * v * v^H
//
// where tau is a complex scalar and v is a complex vector. Note that H is not
// Hermitian in general.
//
// To apply H^H (the conjugate transpose of H), supply cmplx.Conj(tau) instead.
//
// work must have length at least n if side == blas.Left and at least m if
// side == blas.Right.
//
// Zlarf is an internal routine. It is exported for testing purposes.
func (impl Implementation) Zlarf(side blas.Side, m, n int, v []complex128, incv int, tau complex128, c []complex128, ldc int, work []complex128) {
	switch {
	case side != blas.Left && side != blas.Right:
		panic(badSide)
	case m < 0:
		panic(mLT0)
	case n < 0:
		panic(nLT0)
	case incv == 0:
		panic(zeroIncV)
	case ldc < max(1, n):
		panic(badLdC)
	}

	if m == 0 || n == 0 {
		return
	}

	applyleft := side == blas.Left
	lenV := n
	if applyleft {
		lenV = m
	}

	switch {
	case len(v) < 1+(lenV-1)*abs(incv):
		panic(shortV)
	case len(c) < (m-1)*ldc+n:
		panic(shortC)
	case (applyleft && len(work) < n) || (!applyleft && len(work) < m):
		panic(shortWork)
	}

	// Find the last non-zero element of v, scanning trailing zeros so we can
	// minimize the amount of work. lastv is the index of the last non-zero
	// entry, or -1 if v is entirely zero.
	lastv := -1
	lastc := -1
	if tau != 0 {
		if applyleft {
			lastv = m - 1
		} else {
			lastv = n - 1
		}
		var i int
		if incv > 0 {
			i = lastv * incv
		}
		for lastv >= 0 && v[i] == 0 {
			lastv--
			i -= incv
		}
		if applyleft {
			lastc = ilazlc(lastv+1, n, c, ldc)
		} else {
			lastc = ilazlr(m, lastv+1, c, ldc)
		}
	}
	if lastv == -1 || lastc == -1 {
		return
	}

	bi := cblas128.Implementation()
	if applyleft {
		// Form H * C.
		// work[0:lastc+1] = C[0:lastv+1, 0:lastc+1]^H * v[0:lastv+1]
		bi.Zgemv(blas.ConjTrans, lastv+1, lastc+1,
			1, c, ldc, v, incv,
			0, work, 1)
		// C[0:lastv+1, 0:lastc+1] -= tau * v * work^H
		bi.Zgerc(lastv+1, lastc+1,
			-tau, v, incv, work, 1,
			c, ldc)
	} else {
		// Form C * H.
		// work[0:lastc+1] = C[0:lastc+1, 0:lastv+1] * v[0:lastv+1]
		bi.Zgemv(blas.NoTrans, lastc+1, lastv+1,
			1, c, ldc, v, incv,
			0, work, 1)
		// C[0:lastc+1, 0:lastv+1] -= tau * work * v^H
		bi.Zgerc(lastc+1, lastv+1,
			-tau, work, 1, v, incv,
			c, ldc)
	}
}

// ilazlc scans the m×n complex matrix A (column-indexed, row-major) for the
// last column containing any non-zero entry. Returns -1 if A is entirely zero.
func ilazlc(m, n int, a []complex128, lda int) int {
	if n == 0 {
		return -1
	}
	// Quick check of first and last columns.
	if a[n-1] != 0 {
		return n - 1
	}
	for i := 0; i < m; i++ {
		if a[i*lda+n-1] != 0 {
			return n - 1
		}
	}
	// Scan columns right-to-left.
	for j := n - 1; j >= 0; j-- {
		for i := 0; i < m; i++ {
			if a[i*lda+j] != 0 {
				return j
			}
		}
	}
	return -1
}

// ilazlr scans the m×n complex matrix A for the last row containing any
// non-zero entry. Returns -1 if A is entirely zero.
func ilazlr(m, n int, a []complex128, lda int) int {
	if m == 0 {
		return -1
	}
	// Check last row first.
	for j := 0; j < n; j++ {
		if a[(m-1)*lda+j] != 0 {
			return m - 1
		}
	}
	for i := m - 1; i >= 0; i-- {
		for j := 0; j < n; j++ {
			if a[i*lda+j] != 0 {
				return i
			}
		}
	}
	return -1
}
