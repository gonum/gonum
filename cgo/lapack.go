// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cgo provides an interface to bindings for a C LAPACK library.
package cgo

import (
	"github.com/gonum/blas"
	"github.com/gonum/lapack"
	"github.com/gonum/lapack/cgo/clapack"
)

// Copied from lapack/native. Keep in sync.
const (
	badDirect     = "lapack: bad direct"
	badIpiv       = "lapack: insufficient permutation length"
	badLdA        = "lapack: index of a out of range"
	badSide       = "lapack: bad side"
	badStore      = "lapack: bad store"
	badTau        = "lapack: tau has insufficient length"
	badTrans      = "lapack: bad trans"
	badUplo       = "lapack: illegal triangle"
	badWork       = "lapack: insufficient working memory"
	badWorkStride = "lapack: insufficient working array stride"
	negDimension  = "lapack: negative matrix dimension"
	nLT0          = "lapack: n < 0"
	shortWork     = "lapack: working array shorter than declared"
)

func min(m, n int) int {
	if m < n {
		return m
	}
	return n
}

// checkMatrix verifies the parameters of a matrix input.
// Copied from lapack/native. Keep in sync.
func checkMatrix(m, n int, a []float64, lda int) {
	if m < 0 {
		panic("lapack: has negative number of rows")
	}
	if m < 0 {
		panic("lapack: has negative number of columns")
	}
	if lda < n {
		panic("lapack: stride less than number of columns")
	}
	if len(a) < (m-1)*lda+n {
		panic("lapack: insufficient matrix slice length")
	}
}

// Implementation is the cgo-based C implementation of LAPACK routines.
type Implementation struct{}

var _ lapack.Float64 = Implementation{}

// Dpotrf computes the cholesky decomposition of the symmetric positive definite
// matrix a. If ul == blas.Upper, then a is stored as an upper-triangular matrix,
// and a = U U^T is stored in place into a. If ul == blas.Lower, then a = L L^T
// is computed and stored in-place into a. If a is not positive definite, false
// is returned. This is the blocked version of the algorithm.
func (impl Implementation) Dpotrf(ul blas.Uplo, n int, a []float64, lda int) (ok bool) {
	// ul is checked in clapack.Dpotrf.
	if n < 0 {
		panic(nLT0)
	}
	if lda < n {
		panic(badLdA)
	}
	if n == 0 {
		return true
	}
	return clapack.Dpotrf(ul, n, a, lda)
}

// Dgetf2 computes the LU decomposition of the m×n matrix a.
// The LU decomposition is a factorization of a into
//  A = P * L * U
// where P is a permutation matrix, L is a unit lower triangular matrix, and
// U is a (usually) non-unit upper triangular matrix. On exit, L and U are stored
// in place into a.
//
// ipiv is a permutation vector. It indicates that row i of the matrix was
// changed with ipiv[i]. ipiv must have length at least min(m,n), and will panic
// otherwise.
//
// Dgetf2 returns whether the matrix a is singular. The LU decomposition will
// be computed regardless of the singularity of A, but division by zero
// will occur if the false is returned and the result is used to solve a
// system of equations.
func (Implementation) Dgetf2(m, n int, a []float64, lda int, ipiv []int) (ok bool) {
	mn := min(m, n)
	checkMatrix(m, n, a, lda)
	if len(ipiv) < mn {
		panic(badIpiv)
	}
	ipiv32 := make([]int32, len(ipiv))
	ok = clapack.Dgetf2(m, n, a, lda, ipiv32)
	for i, v := range ipiv32 {
		ipiv[i] = int(v) - 1 // OpenBLAS returns one indexed.
	}
	return ok
}
