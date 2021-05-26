// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/cblas128"
)

// Ztrti2 computes the inverse of a triangular matrix, storing the result in place
// into a. This is the BLAS level 2 version of the algorithm.
//
// Ztrti2 is an internal routine. It is exported for testing purposes.
func (impl Implementation) Ztrti2(uplo blas.Uplo, diag blas.Diag, n int, a []complex128, lda int) {
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
		return
	}

	if len(a) < (n-1)*lda+n {
		panic(shortA)
	}

	bi := cblas128.Implementation()

	nonUnit := diag == blas.NonUnit
	// TODO(btracey): Replace this with a row-major ordering.
	if uplo == blas.Upper {
		for j := 0; j < n; j++ {
			var ajj complex128
			if nonUnit {
				ajj = 1 / a[j*lda+j]
				a[j*lda+j] = ajj
				ajj *= -1
			} else {
				ajj = -1
			}
			bi.Ztrmv(blas.Upper, blas.NoTrans, diag, j, a, lda, a[j:], lda)
			bi.Zscal(j, ajj, a[j:], lda)
		}
		return
	}
	for j := n - 1; j >= 0; j-- {
		var ajj complex128
		if nonUnit {
			ajj = 1 / a[j*lda+j]
			a[j*lda+j] = ajj
			ajj *= -1
		} else {
			ajj = -1
		}
		if j < n-1 {
			bi.Ztrmv(blas.Lower, blas.NoTrans, diag, n-j-1, a[(j+1)*lda+j+1:], lda, a[(j+1)*lda+j:], lda)
			bi.Zscal(n-j-1, ajj, a[(j+1)*lda+j:], lda)
		}
	}
}
