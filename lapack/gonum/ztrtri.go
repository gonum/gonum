// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/cblas128"
)

// Ztrtri computes the inverse of a triangular matrix, storing the result in place
// into a. This is the BLAS level 3 version of the algorithm which builds upon
// Dtrti2 to operate on matrix blocks instead of only individual columns.
//
// Ztrtri will not perform the inversion if the matrix is singular, and returns
// a boolean indicating whether the inversion was successful.
func (impl Implementation) Ztrtri(uplo blas.Uplo, diag blas.Diag, n int, a []complex128, lda int) (ok bool) {
	switch {
	case uplo != blas.Upper && uplo != blas.Lower:
		panic(badUplo)
	case diag != blas.NonUnit && diag != blas.Unit:
		panic(badDiag)
	case n < 0:
		panic(nLT0)
	case lda < max(1, n):
		panic(badLdA)
	}

	if n == 0 {
		return true
	}

	if len(a) < (n-1)*lda+n {
		panic(shortA)
	}

	if diag == blas.NonUnit {
		for i := 0; i < n; i++ {
			if a[i*lda+i] == 0 {
				return false
			}
		}
	}

	bi := cblas128.Implementation()

	nb := impl.Ilaenv(1, "ZTRTRI", "UD", n, -1, -1, -1)
	if nb <= 1 || nb > n {
		impl.Ztrti2(uplo, diag, n, a, lda)
		return true
	}
	if uplo == blas.Upper {
		for j := 0; j < n; j += nb {
			jb := min(nb, n-j)
			bi.Ztrmm(blas.Left, blas.Upper, blas.NoTrans, diag, j, jb, 1, a, lda, a[j:], lda)
			bi.Ztrsm(blas.Right, blas.Upper, blas.NoTrans, diag, j, jb, -1, a[j*lda+j:], lda, a[j:], lda)
			impl.Ztrti2(blas.Upper, diag, jb, a[j*lda+j:], lda)
		}
		return true
	}
	nn := ((n - 1) / nb) * nb
	for j := nn; j >= 0; j -= nb {
		jb := min(nb, n-j)
		if j+jb <= n-1 {
			bi.Ztrmm(blas.Left, blas.Lower, blas.NoTrans, diag, n-j-jb, jb, 1, a[(j+jb)*lda+j+jb:], lda, a[(j+jb)*lda+j:], lda)
			bi.Ztrsm(blas.Right, blas.Lower, blas.NoTrans, diag, n-j-jb, jb, -1, a[j*lda+j:], lda, a[(j+jb)*lda+j:], lda)
		}
		impl.Ztrti2(blas.Lower, diag, jb, a[j*lda+j:], lda)
	}
	return true
}
