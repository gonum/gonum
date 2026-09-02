// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math/cmplx"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/cblas128"
)

// Zgetrf computes the LU decomposition of an m×n complex matrix A using partial
// pivoting with row interchanges.
//
// The LU decomposition is a factorization of A into
//
//	A = P * L * U
//
// where P is a permutation matrix, L is a lower triangular with unit diagonal
// elements (lower trapezoidal if m > n), and U is upper triangular (upper
// trapezoidal if m < n).
//
// On entry, a contains the matrix A. On return, L and U are stored in place
// into a, and P is represented by ipiv.
//
// ipiv contains a sequence of row interchanges. It indicates that row i of the
// matrix was interchanged with ipiv[i]. ipiv must have length min(m,n), and
// Zgetrf will panic otherwise. ipiv is zero-indexed.
//
// Zgetrf returns whether the matrix A is nonsingular. The LU decomposition will
// be computed regardless of the singularity of A, but the result should not be
// used to solve a system of equations.
func (impl Implementation) Zgetrf(m, n int, a []complex128, lda int, ipiv []int) (ok bool) {
	mn := min(m, n)
	switch {
	case m < 0:
		panic(mLT0)
	case n < 0:
		panic(nLT0)
	case lda < max(1, n):
		panic(badLdA)
	}

	// Quick return if possible.
	if mn == 0 {
		return true
	}

	switch {
	case len(a) < (m-1)*lda+n:
		panic(shortA)
	case len(ipiv) != mn:
		panic(badLenIpiv)
	}

	bi := cblas128.Implementation()

	nb := impl.Ilaenv(1, "ZGETRF", " ", m, n, -1, -1)
	if nb <= 1 || mn <= nb {
		// Use the unblocked algorithm.
		return impl.zgetf2(m, n, a, lda, ipiv)
	}
	ok = true
	for j := 0; j < mn; j += nb {
		jb := min(mn-j, nb)
		blockOk := impl.zgetf2(m-j, jb, a[j*lda+j:], lda, ipiv[j:j+jb])
		if !blockOk {
			ok = false
		}
		// Adjust pivot indices to be global row indices.
		for i := j; i <= min(m-1, j+jb-1); i++ {
			ipiv[i] = j + ipiv[i]
		}
		// Apply row swaps to the left of the panel.
		impl.zlaswp(j, a, lda, j, j+jb-1, ipiv[:j+jb], 1)
		if j+jb < n {
			// Apply row swaps to the right of the panel.
			impl.zlaswp(n-j-jb, a[j+jb:], lda, j, j+jb-1, ipiv[:j+jb], 1)
			// Compute U block in the panel row.
			bi.Ztrsm(blas.Left, blas.Lower, blas.NoTrans, blas.Unit,
				jb, n-j-jb, 1,
				a[j*lda+j:], lda,
				a[j*lda+j+jb:], lda)
			if j+jb < m {
				// Update the trailing submatrix.
				bi.Zgemm(blas.NoTrans, blas.NoTrans, m-j-jb, n-j-jb, jb, -1,
					a[(j+jb)*lda+j:], lda,
					a[j*lda+j+jb:], lda,
					1, a[(j+jb)*lda+j+jb:], lda)
			}
		}
	}
	return ok
}

// zgetf2 computes the unblocked LU decomposition of an m×n complex matrix A
// using partial pivoting with row interchanges. It is an internal helper used
// by Zgetrf to factor the column panels.
//
// zgetf2 returns whether the matrix A is nonsingular.
func (Implementation) zgetf2(m, n int, a []complex128, lda int, ipiv []int) (ok bool) {
	mn := min(m, n)
	if mn == 0 {
		return true
	}

	bi := cblas128.Implementation()

	sfmin := dlamchS
	ok = true
	for j := 0; j < mn; j++ {
		// Find pivot and test for singularity.
		jp := j + bi.Izamax(m-j, a[j*lda+j:], lda)
		ipiv[j] = jp
		if a[jp*lda+j] == 0 {
			ok = false
		} else {
			// Swap the rows if necessary.
			if jp != j {
				bi.Zswap(n, a[j*lda:], 1, a[jp*lda:], 1)
			}
			// Compute elements below the diagonal of column j.
			if j < m-1 {
				aj := a[j*lda+j]
				if cmplx.Abs(aj) >= sfmin {
					// 1/aj is safe.
					bi.Zscal(m-j-1, 1/aj, a[(j+1)*lda+j:], lda)
				} else {
					// Use explicit division to avoid overflow. Mirrors the
					// DGETF2 safe-minimum branch but for complex divides.
					for i := 0; i < m-j-1; i++ {
						a[(j+1+i)*lda+j] = a[(j+1+i)*lda+j] / aj
					}
				}
			}
		}
		if j < mn-1 {
			// Rank-1 update of trailing submatrix. Use Zgeru (unconjugated).
			bi.Zgeru(m-j-1, n-j-1, -1, a[(j+1)*lda+j:], lda, a[j*lda+j+1:], 1, a[(j+1)*lda+j+1:], lda)
		}
	}
	return ok
}

// zlaswp performs a sequence of row interchanges on the complex matrix a.
// Row k (for k in [k1, k2]) is swapped with row ipiv[k]. If incX is 1, swaps
// are applied in increasing k; if incX is -1, swaps are applied in reverse.
//
// zlaswp is an internal helper; it mirrors Dlaswp for complex matrices. ipiv
// must have length at least k2+1 and all indices are zero-based.
func (Implementation) zlaswp(n int, a []complex128, lda int, k1, k2 int, ipiv []int, incX int) {
	if n <= 0 {
		return
	}
	bi := cblas128.Implementation()
	if incX == 1 {
		for k := k1; k <= k2; k++ {
			if k == ipiv[k] {
				continue
			}
			bi.Zswap(n, a[k*lda:], 1, a[ipiv[k]*lda:], 1)
		}
		return
	}
	// incX == -1
	for k := k2; k >= k1; k-- {
		if k == ipiv[k] {
			continue
		}
		bi.Zswap(n, a[k*lda:], 1, a[ipiv[k]*lda:], 1)
	}
}
