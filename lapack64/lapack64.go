// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lapack64 provides a set of convenient wrapper functions for LAPACK
// calls, as specified in the netlib standard (www.netlib.org).
//
// The native Go routines are used by default, and the Use function can be used
// to set an alternate implementation.
//
// If the type of matrix (General, Symmetric, etc.) is known and fixed, it is
// used in the wrapper signature. In many cases, however, the type of the matrix
// changes during the call to the routine, for example the matrix is symmetric on
// entry and is triangular on exit. In these cases the correct types should be checked
// in the documentation.
//
// The full set of Lapack functions is very large, and it is not clear that a
// full implementation is desirable, let alone feasible. Please open up an issue
// if there is a specific function you need and/or are willing to implement.
package lapack64

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/lapack"
	"github.com/gonum/lapack/native"
)

var lapack64 lapack.Float64 = native.Implementation{}

// Use sets the LAPACK float64 implementation to be used by subsequent BLAS calls.
// The default implementation is native.Implementation.
func Use(l lapack.Float64) {
	lapack64 = l
}

// Potrf computes the cholesky factorization of a.
//  A = U^T * U if ul == blas.Upper
//  A = L * L^T if ul == blas.Lower
// The underlying data between the input matrix and output matrix is shared.
func Potrf(a blas64.Symmetric) (t blas64.Triangular, ok bool) {
	ok = lapack64.Dpotrf(a.Uplo, a.N, a.Data, a.Stride)
	t.Uplo = a.Uplo
	t.N = a.N
	t.Data = a.Data
	t.Stride = a.Stride
	t.Diag = blas.NonUnit
	return
}

// Gels finds a minimum-norm solution based on the matrices A and B using the
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
// Work is temporary storage, and lwork specifies the usable memory length.
// At minimum, lwork >= max(m,n) + max(m,n,nrhs), and this function will panic
// otherwise. A longer work will enable blocked algorithms to be called.
// In the special case that lwork == -1, work[0] will be set to the optimal working
// length.
func Gels(trans blas.Transpose, a blas64.General, b blas64.General, work []float64, lwork int) bool {
	return lapack64.Dgels(trans, a.Rows, a.Cols, b.Cols, a.Data, a.Stride, b.Data, b.Stride, work, lwork)
}

// Geqrf computes the QR factorization of the m×n matrix A using a blocked
// algorithm. A is modified to contain the information to construct Q and R.
// The upper triangle of a contains the matrix R. The lower triangular elements
// (not including the diagonal) contain the elementary reflectors. Tau is modified
// to contain the reflector scales. Tau must have length at least min(m,n), and
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
// Work is temporary storage, and lwork specifies the usable memory length.
// At minimum, lwork >= m and this function will panic otherwise.
// Dgeqrf is a blocked QR factorization, but the block size is limited
// by the temporary space available. If lwork == -1, instead of performing Geqrf,
// the optimal work length will be stored into work[0].
func Geqrf(a blas64.General, tau, work []float64, lwork int) {
	lapack64.Dgeqrf(a.Rows, a.Cols, a.Data, a.Stride, tau, work, lwork)
}

// Gelqf computes the QR factorization of the m×n matrix A using a blocked
// algorithm. A is modified to contain the information to construct L and Q.
// The lower triangle of a contains the matrix L. The lower triangular elements
// (not including the diagonal) contain the elementary reflectors. Tau is modified
// to contain the reflector scales. Tau must have length at least min(m,n), and
// this function will panic otherwise.
//
// See Geqrf for a description of the elementary reflectors and orthonormal
// matrix Q. Q is constructed as a product of these elementary reflectors,
// Q = H_k ... H_2*H_1.
//
// Work is temporary storage, and lwork specifies the usable memory length.
// At minimum, lwork >= m and this function will panic otherwise.
// Dgeqrf is a blocked LQ factorization, but the block size is limited
// by the temporary space available. If lwork == -1, instead of performing Gelqf,
// the optimal work length will be stored into work[0].
func Gelqf(a blas64.General, tau, work []float64, lwork int) {
	lapack64.Dgelqf(a.Rows, a.Cols, a.Data, a.Stride, tau, work, lwork)
}

// Getrf computes the LU decomposition of the m×n matrix A.
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
func Getrf(a blas64.General, ipiv []int) bool {
	return lapack64.Dgetrf(a.Rows, a.Cols, a.Data, a.Stride, ipiv)
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
// computed by Getrf. ipiv is zero-indexed.
func Getrs(trans blas.Transpose, a blas64.General, b blas64.General, ipiv []int) {
	lapack64.Dgetrs(trans, a.Cols, b.Cols, a.Data, a.Stride, ipiv, b.Data, b.Stride)
}
