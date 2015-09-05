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
	absIncNotOne  = "lapack: increment not one or negative one"
	badDiag       = "lapack: bad diag"
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

func max(m, n int) int {
	if m < n {
		return n
	}
	return m
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

// Dgelq2 computes the LQ factorization of the m×n matrix A.
//
// In an LQ factorization, L is a lower triangular m×n matrix, and Q is an n×n
// orthornormal matrix.
//
// a is modified to contain the information to construct L and Q.
// The lower triangle of a contains the matrix L. The upper triangular elements
// (not including the diagonal) contain the elementary reflectors. Tau is modified
// to contain the reflector scales. tau must have length of at least k = min(m,n)
// and this function will panic otherwise.
//
// See Dgeqr2 for a description of the elementary reflectors and orthonormal
// matrix Q. Q is constructed as a product of these elementary reflectors,
// Q = H_k ... H_2*H_1.
//
// Work is temporary storage of length at least m and this function will panic otherwise.
func (impl Implementation) Dgelq2(m, n int, a []float64, lda int, tau, work []float64) {
	checkMatrix(m, n, a, lda)
	if len(tau) < min(m, n) {
		panic(badTau)
	}
	if len(work) < m {
		panic(badWork)
	}
	clapack.Dgelq2(m, n, a, lda, tau)
}

// Dgelqf computes the LQ factorization of the m×n matrix A using a blocked
// algorithm. See the documentation for Dgelq2 for a description of the
// parameters at entry and exit.
//
// The C interface does not support providing temporary storage. To provide compatibility
// with native, lwork == -1 will not run Dgeqrf but will instead write the minimum
// work necessary to work[0]. If len(work) < lwork, Dgeqrf will panic.
//
// tau must have length at least min(m,n), and this function will panic otherwise.
func (impl Implementation) Dgelqf(m, n int, a []float64, lda int, tau, work []float64, lwork int) {
	if lwork == -1 {
		work[0] = float64(m)
		return
	}
	checkMatrix(m, n, a, lda)
	if len(work) < lwork {
		panic(shortWork)
	}
	if lwork < m {
		panic(badWork)
	}
	if len(tau) < min(m, n) {
		panic(badTau)
	}
	clapack.Dgelqf(m, n, a, lda, tau)
}

// Dgeqr2 computes a QR factorization of the m×n matrix A.
//
// In a QR factorization, Q is an m×m orthonormal matrix, and R is an
// upper triangular m×n matrix.
//
// A is modified to contain the information to construct Q and R.
// The upper triangle of a contains the matrix R. The lower triangular elements
// (not including the diagonal) contain the elementary reflectors. Tau is modified
// to contain the reflector scales. tau must have length at least min(m,n), and
// this function will panic otherwise.
//
// The ith elementary reflector can be explicitly constructed by first extracting
// the
//  v[j] = 0           j < i
//  v[j] = i           j == i
//  v[j] = a[i*lda+j]  j > i
// and computing h_i = I - tau[i] * v * v^T.
//
// The orthonormal matrix Q can be constucted from a product of these elementary
// reflectors, Q = H_1*H_2 ... H_k, where k = min(m,n).
//
// Work is temporary storage of length at least n and this function will panic otherwise.
func (impl Implementation) Dgeqr2(m, n int, a []float64, lda int, tau, work []float64) {
	checkMatrix(m, n, a, lda)
	if len(work) < n {
		panic(badWork)
	}
	k := min(m, n)
	if len(tau) < k {
		panic(badTau)
	}
	clapack.Dgeqr2(m, n, a, lda, tau)
}

// Dgeqrf computes the QR factorization of the m×n matrix A using a blocked
// algorithm. See the documentation for Dgeqr2 for a description of the
// parameters at entry and exit.
//
// The C interface does not support providing temporary storage. To provide compatibility
// with native, lwork == -1 will not run Dgeqrf but will instead write the minimum
// work necessary to work[0]. If len(work) < lwork, Dgeqrf will panic.
//
// tau must have length at least min(m,n), and this function will panic otherwise.
func (impl Implementation) Dgeqrf(m, n int, a []float64, lda int, tau, work []float64, lwork int) {
	if lwork == -1 {
		work[0] = float64(n)
		return
	}
	checkMatrix(m, n, a, lda)
	if len(work) < lwork {
		panic(shortWork)
	}
	if lwork < n {
		panic(badWork)
	}
	k := min(m, n)
	if len(tau) < k {
		panic(badTau)
	}
	clapack.Dgeqrf(m, n, a, lda, tau)
}

// Dgels finds a minimum-norm solution based on the matrices A and B using the
// QR or LQ factorization. Dgels returns false if the matrix
// A is singular, and true if this solution was successfully found.
//
// The minimization problem solved depends on the input parameters.
//
//  1. If m >= n and trans == blas.NoTrans, Dgels finds X such that || A*X - B||_2
//     is minimized.
//  2. If m < n and trans == blas.NoTrans, Dgels finds the minimum norm solution of
//     A * X = B.
//  3. If m >= n and trans == blas.Trans, Dgels finds the minimum norm solution of
//     A^T * X = B.
//  4. If m < n and trans == blas.Trans, Dgels finds X such that || A*X - B||_2
//     is minimized.
// Note that the least-squares solutions (cases 1 and 3) perform the minimization
// per column of B. This is not the same as finding the minimum-norm matrix.
//
// The matrix A is a general matrix of size m×n and is modified during this call.
// The input matrix B is of size max(m,n)×nrhs, and serves two purposes. On entry,
// the elements of b specify the input matrix B. B has size m×nrhs if
// trans == blas.NoTrans, and n×nrhs if trans == blas.Trans. On exit, the
// leading submatrix of b contains the solution vectors X. If trans == blas.NoTrans,
// this submatrix is of size n×nrhs, and of size m×nrhs otherwise.
//
// The C interface does not support providing temporary storage. To provide compatibility
// with native, lwork == -1 will not run Dgeqrf but will instead write the minimum
// work necessary to work[0]. If len(work) < lwork, Dgeqrf will panic.
func (impl Implementation) Dgels(trans blas.Transpose, m, n, nrhs int, a []float64, lda int, b []float64, ldb int, work []float64, lwork int) bool {
	mn := min(m, n)
	if lwork == -1 {
		work[0] = float64(mn + max(mn, nrhs))
		return true
	}
	checkMatrix(m, n, a, lda)
	checkMatrix(max(m, n), nrhs, b, ldb)
	if len(work) < lwork {
		panic(shortWork)
	}
	if lwork < mn+max(mn, nrhs) {
		panic(badWork)
	}
	return clapack.Dgels(trans, m, n, nrhs, a, lda, b, ldb)
}

// Dgetf2 computes the LU decomposition of the m×n matrix A.
// The LU decomposition is a factorization of a into
//  A = P * L * U
// where P is a permutation matrix, L is a unit lower triangular matrix, and
// U is a (usually) non-unit upper triangular matrix. On exit, L and U are stored
// in place into a.
//
// ipiv is a permutation vector. It indicates that row i of the matrix was
// changed with ipiv[i]. ipiv must have length at least min(m,n), and will panic
// otherwise. ipiv is zero-indexed.
//
// Dgetf2 returns whether the matrix A is singular. The LU decomposition will
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
		ipiv[i] = int(v) - 1 // Transform to zero-indexed.
	}
	return ok
}

// Dgetrf computes the LU decomposition of the m×n matrix A.
// The LU decomposition is a factorization of A into
//  A = P * L * U
// where P is a permutation matrix, L is a unit lower triangular matrix, and
// U is a (usually) non-unit upper triangular matrix. On exit, L and U are stored
// in place into a.
//
// ipiv is a permutation vector. It indicates that row i of the matrix was
// changed with ipiv[i]. ipiv must have length at least min(m,n), and will panic
// otherwise. ipiv is zero-indexed.
//
// Dgetrf is the blocked version of the algorithm.
//
// Dgetrf returns whether the matrix A is singular. The LU decomposition will
// be computed regardless of the singularity of A, but division by zero
// will occur if the false is returned and the result is used to solve a
// system of equations.
func (impl Implementation) Dgetrf(m, n int, a []float64, lda int, ipiv []int) (ok bool) {
	mn := min(m, n)
	checkMatrix(m, n, a, lda)
	if len(ipiv) < mn {
		panic(badIpiv)
	}
	ipiv32 := make([]int32, len(ipiv))
	ok = clapack.Dgetrf(m, n, a, lda, ipiv32)
	for i, v := range ipiv32 {
		ipiv[i] = int(v) - 1 // Transform to zero-indexed.
	}
	return ok
}

// Dgetrs solves a system of equations using an LU factorization.
// The system of equations solved is
//  A * X = B if trans == blas.Trans
//  A^T * X = B if trans == blas.NoTrans
// A is a general n×n matrix with stride lda. B is a general matrix of size n×nrhs.
//
// On entry b contains the elements of the matrix B. On exit, b contains the
// elements of X, the solution to the system of equations.
//
// a and ipiv contain the LU factorization of A and the permutation indices as
// computed by Dgetrf. ipiv is zero-indexed.
func (impl Implementation) Dgetrs(trans blas.Transpose, n, nrhs int, a []float64, lda int, ipiv []int, b []float64, ldb int) {
	checkMatrix(n, n, a, lda)
	checkMatrix(n, nrhs, b, ldb)
	if len(ipiv) < n {
		panic(badIpiv)
	}
	ipiv32 := make([]int32, len(ipiv))
	for i, v := range ipiv {
		ipiv32[i] = int32(v) + 1 // Transform to one-indexed.
	}
	clapack.Dgetrs(trans, n, nrhs, a, lda, ipiv32, b, ldb)
}

// Dormlq multiplies the matrix C by the othogonal matrix Q defined by the
// slices a and tau. A and tau are as returned from Dgelqf.
//  C = Q * C    if side == blas.Left and trans == blas.NoTrans
//  C = Q^T * C  if side == blas.Left and trans == blas.Trans
//  C = C * Q    if side == blas.Right and trans == blas.NoTrans
//  C = C * Q^T  if side == blas.Right and trans == blas.Trans
// If side == blas.Left, A is a matrix of side k×m, and if side == blas.Right
// A is of size k×n. This uses a blocked algorithm.
//
// Work is temporary storage, and lwork specifies the usable memory length.
// At minimum, lwork >= m if side == blas.Left and lwork >= n if side == blas.Right,
// and this function will panic otherwise.
// Dormlq uses a block algorithm, but the block size is limited
// by the temporary space available. If lwork == -1, instead of performing Dormlq,
// the optimal work length will be stored into work[0].
//
// tau contains the householder scales and must have length at least k, and
// this function will panic otherwise.
func (impl Implementation) Dormlq(side blas.Side, trans blas.Transpose, m, n, k int, a []float64, lda int, tau, c []float64, ldc int, work []float64, lwork int) {
	if side != blas.Left && side != blas.Right {
		panic(badSide)
	}
	if trans != blas.Trans && trans != blas.NoTrans {
		panic(badTrans)
	}
	left := side == blas.Left
	if left {
		checkMatrix(k, m, a, lda)
	} else {
		checkMatrix(k, n, a, lda)
	}
	checkMatrix(m, n, c, ldc)
	if len(tau) < k {
		panic(badTau)
	}
	if lwork == -1 {
		if left {
			work[0] = float64(n)
			return
		}
		work[0] = float64(m)
		return
	}
	if left {
		if lwork < n {
			panic(badWork)
		}
	} else {
		if lwork < m {
			panic(badWork)
		}
	}
	clapack.Dormlq(side, trans, m, n, k, a, lda, tau, c, ldc)
}

// Dormqr multiplies the matrix C by the othogonal matrix Q defined by the
// slices a and tau. a and tau are as returned from Dgeqrf.
//  C = Q * C    if side == blas.Left and trans == blas.NoTrans
//  C = Q^T * C  if side == blas.Left and trans == blas.Trans
//  C = C * Q    if side == blas.Right and trans == blas.NoTrans
//  C = C * Q^T  if side == blas.Right and trans == blas.Trans
// If side == blas.Left, A is a matrix of side k×m, and if side == blas.Right
// A is of size k×n. This uses a blocked algorithm.
//
// tau contains the householder scales and must have length at least k, and
// this function will panic otherwise.
//
// The C interface does not support providing temporary storage. To provide compatibility
// with native, lwork == -1 will not run Dgeqrf but will instead write the minimum
// work necessary to work[0]. If len(work) < lwork, Dgeqrf will panic.
func (impl Implementation) Dormqr(side blas.Side, trans blas.Transpose, m, n, k int, a []float64, lda int, tau, c []float64, ldc int, work []float64, lwork int) {
	left := side == blas.Left
	if left {
		checkMatrix(m, k, a, lda)
	} else {
		checkMatrix(n, k, a, lda)
	}
	checkMatrix(m, n, c, ldc)

	if len(tau) < k {
		panic(badTau)
	}

	if lwork == -1 {
		if left {
			work[0] = float64(m)
			return
		}
		work[0] = float64(n)
		return
	}

	if left {
		if lwork < n {
			panic(badWork)
		}
	} else {
		if lwork < m {
			panic(badWork)
		}
	}

	clapack.Dormqr(side, trans, m, n, k, a, lda, tau, c, ldc)
}

// Dtrtrs solves a triangular system of the form A * X = B or A^T * X = B. Dtrtrs
// returns whether the solve completed successfully. If A is singular, no solve is performed.
func (impl Implementation) Dtrtrs(uplo blas.Uplo, trans blas.Transpose, diag blas.Diag, n, nrhs int, a []float64, lda int, b []float64, ldb int) (ok bool) {
	return clapack.Dtrtrs(uplo, trans, diag, n, nrhs, a, lda, b, ldb)
}
