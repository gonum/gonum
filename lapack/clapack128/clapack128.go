// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clapack128

import (
	"gonum.org/v1/gonum/blas/cblas128"
	"gonum.org/v1/gonum/lapack"
	"gonum.org/v1/gonum/lapack/gonum"
)

var clapack128 lapack.Complex128 = gonum.Implementation{}

// Use sets the LAPACK complex128 implementation to be used by subsequent BLAS calls.
// The default implementation is native.Implementation.
func Use(l lapack.Complex128) {
	clapack128 = l
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Getri computes the inverse of the matrix A using the LU factorization computed
// by Getrf. On entry, a contains the PLU decomposition of A as computed by
// Getrf and on exit contains the reciprocal of the original matrix.
//
// Getri will not perform the inversion if the matrix is singular, and returns
// a boolean indicating whether the inversion was successful.
//
// Work is temporary storage, and lwork specifies the usable memory length.
// At minimum, lwork >= n and this function will panic otherwise.
// Getri is a blocked inversion, but the block size is limited
// by the temporary space available. If lwork == -1, instead of performing Getri,
// the optimal work length will be stored into work[0].
func Getri(a cblas128.General, ipiv []int, work []complex128, lwork int) (ok bool) {
	return clapack128.Zgetri(a.Cols, a.Data, max(1, a.Stride), ipiv, work, lwork)
}
