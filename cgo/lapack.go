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
	absIncNotOne    = "lapack: increment not one or negative one"
	badD            = "lapack: d has insufficient length"
	badDecompUpdate = "lapack: bad decomp update"
	badDiag         = "lapack: bad diag"
	badDims         = "lapack: bad input dimensions"
	badDirect       = "lapack: bad direct"
	badE            = "lapack: e has insufficient length"
	badIpiv         = "lapack: insufficient permutation length"
	badLdA          = "lapack: index of a out of range"
	badNorm         = "lapack: bad norm"
	badPivot        = "lapack: bad pivot"
	badS            = "lapack: s has insufficient length"
	badSide         = "lapack: bad side"
	badSlice        = "lapack: bad input slice length"
	badStore        = "lapack: bad store"
	badTau          = "lapack: tau has insufficient length"
	badTauQ         = "lapack: tauQ has insufficient length"
	badTauP         = "lapack: tauP has insufficient length"
	badTrans        = "lapack: bad trans"
	badUplo         = "lapack: illegal triangle"
	badWork         = "lapack: insufficient working memory"
	badWorkStride   = "lapack: insufficient working array stride"
	badZ            = "lapack: insufficient z length"
	kGTM            = "lapack: k > m"
	kGTN            = "lapack: k > n"
	kLT0            = "lapack: k < 0"
	mLTN            = "lapack: m < n"
	negDimension    = "lapack: negative matrix dimension"
	negZ            = "lapack: negative z value"
	nLT0            = "lapack: n < 0"
	nLTM            = "lapack: n < m"
	shortWork       = "lapack: working array shorter than declared"
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

// Dlacpy copies the elements of A specified by uplo into B. Uplo can specify
// a triangular portion with blas.Upper or blas.Lower, or can specify all of the
// elemest with blas.All.
func (impl Implementation) Dlacpy(uplo blas.Uplo, m, n int, a []float64, lda int, b []float64, ldb int) {
	checkMatrix(m, n, a, lda)
	checkMatrix(m, n, b, ldb)
	clapack.Dlacpy(uplo, m, n, a, lda, b, ldb)
}

// Dlange computes the matrix norm of the general m×n matrix a. The input norm
// specifies the norm computed.
//  lapack.MaxAbs: the maximum absolute value of an element.
//  lapack.MaxColumnSum: the maximum column sum of the absolute values of the entries.
//  lapack.MaxRowSum: the maximum row sum of the absolute values of the entries.
//  lapack.Frobenius: the square root of the sum of the squares of the entries.
// If norm == lapack.MaxColumnSum, work must be of length n, and this function will panic otherwise.
// There are no restrictions on work for the other matrix norms.
func (impl Implementation) Dlange(norm lapack.MatrixNorm, m, n int, a []float64, lda int, work []float64) float64 {
	checkMatrix(m, n, a, lda)
	switch norm {
	case lapack.MaxRowSum, lapack.MaxColumnSum, lapack.NormFrob, lapack.MaxAbs:
	default:
		panic(badNorm)
	}
	if norm == lapack.MaxColumnSum && len(work) < n {
		panic(badWork)
	}
	return clapack.Dlange(byte(norm), m, n, a, lda)
}

// Dlansy computes the specified norm of an n×n symmetric matrix. If
// norm == lapack.MaxColumnSum or norm == lapackMaxRowSum work must have length
// at least n, otherwise work is unused.
func (impl Implementation) Dlansy(norm lapack.MatrixNorm, uplo blas.Uplo, n int, a []float64, lda int, work []float64) float64 {
	checkMatrix(n, n, a, lda)
	switch norm {
	case lapack.MaxRowSum, lapack.MaxColumnSum, lapack.NormFrob, lapack.MaxAbs:
	default:
		panic(badNorm)
	}
	if (norm == lapack.MaxColumnSum || norm == lapack.MaxRowSum) && len(work) < n {
		panic(badWork)
	}
	if uplo != blas.Upper && uplo != blas.Lower {
		panic(badUplo)
	}
	return clapack.Dlansy(byte(norm), uplo, n, a, lda)
}

// Dlantr computes the specified norm of an m×n trapezoidal matrix A. If
// norm == lapack.MaxColumnSum work must have length at least n, otherwise work
// is unused.
func (impl Implementation) Dlantr(norm lapack.MatrixNorm, uplo blas.Uplo, diag blas.Diag, m, n int, a []float64, lda int, work []float64) float64 {
	checkMatrix(m, n, a, lda)
	switch norm {
	case lapack.MaxRowSum, lapack.MaxColumnSum, lapack.NormFrob, lapack.MaxAbs:
	default:
		panic(badNorm)
	}
	if uplo != blas.Upper && uplo != blas.Lower {
		panic(badUplo)
	}
	if diag != blas.Unit && diag != blas.NonUnit {
		panic(badDiag)
	}
	if norm == lapack.MaxColumnSum && len(work) < n {
		panic(badWork)
	}
	return clapack.Dlantr(byte(norm), uplo, diag, m, n, a, lda)
}

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

// Dbdsqr performs a singular value decomposition of a real n×n bidiagonal matrix.
//
// The SVD of the bidiagonal matrix B is
//  B = Q * S * P^T
// where S is a diagonal matrix of singular values, Q is an orthogonal matrix of
// left singular vectors, and P is an orthogonal matrix of right singular vectors.
//
// Q and P are only computed if requested. If left singular vectors are requested,
// this routine returns U * Q instead of Q, and if right singular vectors are
// requested P^T * VT is returned instead of P^T.
//
// Frequently Dbdsqr is used in conjuction with Dgebrd which reduces a general
// matrix A into bidiagonal form. In this case, the SVD of A is
//  A = (U * Q) * S * (P^T * VT)
//
// This routine may also compute Q^T * C.
//
// d and e contain the elements of the bidiagonal matrix b. d must have length at
// least n, and e must have length at least n-1. Dbdsqr will panic if there is
// insufficient length. On exit, D contains the singular values of B in decreasing
// order.
//
// VT is a matrix of size n×ncvt whose elements are stored in vt. The elements
// of vt are modified to contain P^T * VT on exit. VT is not used if ncvt == 0.
//
// U is a matrix of size nru×n whose elements are stored in u. The elements
// of u are modified to contain U * Q on exit. U is not used if nru == 0.
//
// C is a matrix of size n×ncc whose elements are stored in c. The elements
// of c are modified to contain Q^T * C on exit. C is not used if ncc == 0.
//
// work contains temporary storage and must have length at least 4*n. Dbdsqr
// will panic if there is insufficient working memory.
//
// Dbdsqr returns whether the decomposition was successful.
func (impl Implementation) Dbdsqr(uplo blas.Uplo, n, ncvt, nru, ncc int, d, e, vt []float64, ldvt int, u []float64, ldu int, c []float64, ldc int, work []float64) (ok bool) {
	if uplo != blas.Upper && uplo != blas.Lower {
		panic(badUplo)
	}
	if ncvt != 0 {
		checkMatrix(n, ncvt, vt, ldvt)
	}
	if nru != 0 {
		checkMatrix(nru, n, u, ldu)
	}
	if ncc != 0 {
		checkMatrix(n, ncc, c, ldc)
	}
	if len(d) < n {
		panic(badD)
	}
	if len(e) < n-1 {
		panic(badE)
	}
	if len(work) < 4*n {
		panic(badWork)
	}
	// An address must be passed to cgo. If lengths are zero, allocate a slice.
	if len(vt) == 0 {
		vt = make([]float64, 1)
	}
	if len(u) == 0 {
		vt = make([]float64, 1)
	}
	if len(c) == 0 {
		c = make([]float64, 1)
	}
	return clapack.Dbdsqr(uplo, n, ncvt, nru, ncc, d, e, vt, ldvt, u, ldu, c, ldc)
}

// Dgebrd reduces a general m×n matrix A to upper or lower bidiagonal form B by
// an orthogonal transformation:
//  Q^T * A * P = B.
// The diagonal elements of B are stored in d and the off-diagonal elements are
// stored in e. These are additionally stored along the diagonal of A and the
// off-diagonal of A. If m >= n B is an upper-bidiagonal matrix, and if m < n B
// is a lower-bidiagonal matrix.
//
// The remaining elements of A store the data needed to construct Q and P.
// The matrices Q and P are products of elementary reflectors
//  Q = H_1 * H_2 * ... * H_nb
//  P = G_1 * G_2 * ... * G_nb
// where
//  H_i = I - tauQ[i] * v_i * v_i^T
//  G_i = I - tauP[i] * u_i * u_i^T
//
// As an example, on exit the entries of A when m = 6, and n = 5
//  (  d   e   u1  u1  u1 )
//  (  v1  d   e   u2  u2 )
//  (  v1  v2  d   e   u3 )
//  (  v1  v2  v3  d   e  )
//  (  v1  v2  v3  v4  d  )
//  (  v1  v2  v3  v4  v5 )
// and when m = 5, n = 6
//  (  d   u1  u1  u1  u1  u1 )
//  (  e   d   u2  u2  u2  u2 )
//  (  v1  e   d   u3  u3  u3 )
//  (  v1  v2  e   d   u4  u4 )
//  (  v1  v2  v3  e   d   u5 )
//
// d, tauQ, and tauP must all have length at least min(m,n), and e must have
// length min(m,n) - 1.
//
// Work is temporary storage, and lwork specifies the usable memory length.
// At minimum, lwork >= max(m,n) and this function will panic otherwise.
// The C interface does not support providing temporary storage. To provide compatibility
// with native, lwork == -1 will not run Dgeqrf but will instead write the minimum
// work necessary to work[0]. If len(work) < lwork, Dgbrd will panic.
func (impl Implementation) Dgebrd(m, n int, a []float64, lda int, d, e, tauQ, tauP, work []float64, lwork int) {
	checkMatrix(m, n, a, lda)
	minmn := min(m, n)
	if len(d) < minmn {
		panic(badD)
	}
	if len(e) < minmn-1 {
		panic(badE)
	}
	if len(tauQ) < minmn {
		panic(badTauQ)
	}
	if len(tauP) < minmn {
		panic(badTauP)
	}
	ws := max(m, n)
	if lwork == -1 {
		work[0] = float64(ws)
		return
	}
	if lwork < ws {
		panic(badWork)
	}
	if len(work) < lwork {
		panic(badWork)
	}

	clapack.Dgebrd(m, n, a, lda, d, e, tauQ, tauP)
}

// Dgecon estimates the reciprocal of the condition number of the n×n matrix A
// given the LU decomposition of the matrix. The condition number computed may
// be based on the 1-norm or the ∞-norm.
//
// The slice a contains the result of the LU decomposition of A as computed by Dgetrf.
//
// anorm is the corresponding 1-norm or ∞-norm of the original matrix A.
//
// work is a temporary data slice of length at least 4*n and Dgecon will panic otherwise.
//
// iwork is a temporary data slice of length at least n and Dgecon will panic otherwise.
func (impl Implementation) Dgecon(norm lapack.MatrixNorm, n int, a []float64, lda int, anorm float64, work []float64, iwork []int) float64 {
	checkMatrix(n, n, a, lda)
	if norm != lapack.MaxColumnSum && norm != lapack.MaxRowSum {
		panic("bad norm")
	}
	if len(work) < 4*n {
		panic(badWork)
	}
	if len(iwork) < n {
		panic(badWork)
	}
	rcond := make([]float64, 1)
	clapack.Dgecon(byte(norm), n, a, lda, anorm, rcond)
	return rcond[0]
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

// Dgesvd computes the singular value decomposition of the input matrix A.
//
// The singular value decomposition is
//  A = U * Sigma * V^T
// where Sigma is an m×n diagonal matrix containing the singular values of A,
// U is an m×m orthogonal matrix and V is an n×n orthogonal matrix. The first
// min(m,n) columns of U and V are the left and right singular vectors of A
// respectively.
//
// jobU and jobVT are options for computing the singular vectors. The behavior
// is as follows
//  jobU == lapack.SVDAll		All M columns of U are returned in u
//  jobU == lapack.SVDInPlace	The first min(m,n) columns are returned in u
//  jobU == lapack.SVDOverwrite	The first min(m,n) columns of U are written into a
//	jobU == lapack.SVDNone		The columns of U are not computed.
// The behavior is the same for jobVT and the rows of V^T. At most one of jobU
// and jobVT can equal lapack.SVDOverwrite.
//
// On entry, a contains the data for the m×n matrix A. During the call to Dgesvd
// the data is overwritten. On exit, A contains the appropriate singular vectors
// if either job is lapack.SVDOverwrite.
//
// s is a slice of length at least min(m,n) and on exit contains the singular
// values in decreasing order.
//
// u contains the left singular vectors on exit, stored columnwise. If
// jobU == lapack.SVDAll, u is of size m×m. If jobU == lapack.SVDInPlace u is
// of size m×min(m,n). If jobU == lapack.SVDOverwrite or lapack.SVDNone, u is
// not used.
//
// vt contains the left singular vectors on exit, stored rowwise. If
// jobV == lapack.SVDAll, vt is of size n×m. If jobV == lapack.SVDInPlace vt is
// of size min(m,n)×n. If jobU == lapack.SVDOverwrite or lapack.SVDNone, vt is
// not used.
//
// The C interface does not support providing temporary storage. To provide compatibility
// with native, lwork == -1 will not run Dgeqrf but will instead write the minimum
// work necessary to work[0]. If len(work) < lwork, Dgeqrf will panic.
//
// Dgesvd returns whether the decomposition successfully completed.
func (impl Implementation) Dgesvd(jobU, jobVT lapack.SVDJob, m, n int, a []float64, lda int, s, u []float64, ldu int, vt []float64, ldvt int, work []float64, lwork int) (ok bool) {
	checkMatrix(m, n, a, lda)
	if jobU == lapack.SVDAll {
		checkMatrix(m, m, u, ldu)
	} else if jobU == lapack.SVDInPlace {
		checkMatrix(m, min(m, n), u, ldu)
	}
	if jobVT == lapack.SVDAll {
		checkMatrix(n, n, vt, ldvt)
	} else if jobVT == lapack.SVDInPlace {
		checkMatrix(min(m, n), n, vt, ldvt)
	}
	if jobU == lapack.SVDOverwrite && jobVT == lapack.SVDOverwrite {
		panic("lapack: both jobU and jobVT are lapack.SVDOverwrite")
	}
	if len(s) < min(m, n) {
		panic(badS)
	}
	if jobU != lapack.SVDAll || jobVT != lapack.SVDAll {
		panic("lapack: SVD only coded for SVDAll job inputs")
	}
	minWork := max(5*min(m, n), 3*min(m, n)+max(m, n))
	if lwork != -1 {
		if len(work) < lwork {
			panic(badWork)
		}
		if lwork < minWork {
			panic(badWork)
		}
	}
	if lwork == -1 {
		work[0] = float64(minWork)
		return true
	}
	return clapack.Dgesvd(byte(jobU), byte(jobVT), m, n, a, lda, s, u, ldu, vt, ldvt, work)
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

// Dgetri computes the inverse of the matrix A using the LU factorization computed
// by Dgetrf. On entry, a contains the PLU decomposition of A as computed by
// Dgetrf and on exit contains the reciprocal of the original matrix.
//
// Dtrtri will not perform the inversion if the matrix is singular, and returns
// a boolean indicating whether the inversion was successful.
//
// The C interface does not support providing temporary storage. To provide compatibility
// with native, lwork == -1 will not run Dgetri but will instead write the minimum
// work necessary to work[0]. If len(work) < lwork, Dgetri will panic.
func (impl Implementation) Dgetri(n int, a []float64, lda int, ipiv []int, work []float64, lwork int) (ok bool) {
	checkMatrix(n, n, a, lda)
	if len(ipiv) < n {
		panic(badIpiv)
	}
	if lwork == -1 {
		work[0] = float64(n)
		return true
	}
	if lwork < n {
		panic(badWork)
	}
	if len(work) < lwork {
		panic(badWork)
	}
	ipiv32 := make([]int32, len(ipiv))
	for i, v := range ipiv {
		ipiv32[i] = int32(v) + 1 // Transform to one-indexed.
	}
	return clapack.Dgetri(n, a, lda, ipiv32)
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

// Dorgbr generates one of the matrices Q or P^T computed by Dgebrd
// computed from the decomposition Dgebrd. See Dgebd2 for the description of
// Q and P^T.
//
// If vect == lapack.ApplyQ, then a is assumed to have been an m×k matrix and
// Q is of order m. If m >= k, then Dorgbr returns the first n columns of Q
// where m >= n >= k. If m < k, then Dorgbr returns Q as an m×m matrix.
//
// If vect == lapack.ApplyP, then A is assumed to have been a k×n matrix, and
// P^T is of order n. If k < n, then Dorgbr returns the first m rows of P^T,
// where n >= m >= k. If k >= n, then Dorgbr returns P^T as an n×n matrix.
func (impl Implementation) Dorgbr(vect lapack.DecompUpdate, m, n, k int, a []float64, lda int, tau, work []float64, lwork int) {
	mn := min(m, n)
	wantq := vect == lapack.ApplyQ
	if wantq {
		if m < n || n < min(m, k) || m < min(m, k) {
			panic(badDims)
		}
	} else {
		if n < m || m < min(n, k) || n < min(n, k) {
			panic(badDims)
		}
	}
	if wantq {
		checkMatrix(m, k, a, lda)
	} else {
		checkMatrix(k, n, a, lda)
	}
	if lwork == -1 {
		work[0] = float64(mn)
		return
	}
	if len(work) < lwork {
		panic(badWork)
	}
	if lwork < mn {
		panic(badWork)
	}
	clapack.Dorgbr(byte(vect), m, n, k, a, lda, tau)
}

// Dorglq generates an m×n matrix Q with orthonormal rows defined by the
// product of elementary reflectors as computed by Dgelqf.
//  Q = H(0) * H(2) * ... * H(k-1)
// Dorglq is the blocked version of dorgl2 that makes greater use of level-3 BLAS
// routines.
//
// len(tau) >= k, 0 <= k <= n, and 0 <= m <= n.
//
// Work is temporary storage, and lwork specifies the usable memory length.
// The C interface does not support providing temporary storage. To provide compatibility
// with native, lwork == -1 will not run Dorglq but will instead write the minimum
// work necessary to work[0]. If len(work) < lwork, Dgeqrf will panic, and at minimum
// lwork >= m.
//
// Dorgqr will panic if the conditions on input values are not met.
func (impl Implementation) Dorglq(m, n, k int, a []float64, lda int, tau, work []float64, lwork int) {
	if lwork == -1 {
		work[0] = float64(m)
		return
	}
	checkMatrix(m, n, a, lda)
	if k < 0 {
		panic(kLT0)
	}
	if k > m {
		panic(kGTM)
	}
	if m > n {
		panic(nLTM)
	}
	if len(tau) < k {
		panic(badTau)
	}
	if len(work) < lwork {
		panic(shortWork)
	}
	if lwork < m {
		panic(badWork)
	}
	clapack.Dorglq(m, n, k, a, lda, tau)
}

// Dorgqr generates an m×n matrix Q with orthonormal columns defined by the
// product of elementary reflectors as computed by Dgeqrf.
//  Q = H(0) * H(2) * ... * H(k-1)
// Dorgqr is the blocked version of dorg2r that makes greater use of level-3 BLAS
// routines.
//
// len(tau) >= k, 0 <= k <= n, and 0 <= n <= m.
//
// Work is temporary storage, and lwork specifies the usable memory length.
// The C interface does not support providing temporary storage. To provide compatibility
// with native, lwork == -1 will not run Dorgqr but will instead write the minimum
// work necessary to work[0]. If len(work) < lwork, Dgeqrf will panic, and at minimum
// lwork >= n.
//
// Dorgqr will panic if the conditions on input values are not met.
func (impl Implementation) Dorgqr(m, n, k int, a []float64, lda int, tau, work []float64, lwork int) {
	if lwork == -1 {
		work[0] = float64(n)
		return
	}
	checkMatrix(m, n, a, lda)
	if k < 0 {
		panic(kLT0)
	}
	if k > n {
		panic(kGTN)
	}
	if n > m {
		panic(mLTN)
	}
	if len(tau) < k {
		panic(badTau)
	}
	if len(work) < lwork {
		panic(shortWork)
	}
	if lwork < n {
		panic(badWork)
	}
	clapack.Dorgqr(m, n, k, a, lda, tau)
}

// Dormbr applies a multiplicative update to the matrix C based on a
// decomposition computed by Dgebrd.
//
// Dormbr computes
//  Q * C if vect == lapack.ApplyQ, side == blas.Left, and trans == blas.NoTrans
//  C * Q if vect == lapack.ApplyQ, side == blas.Right, and trans == blas.NoTrans
//  Q^T * C if vect == lapack.ApplyQ, side == blas.Left, and trans == blas.Trans
//  C * Q^T if vect == lapack.ApplyQ, side == blas.Right, and trans == blas.Trans
//
//  P * C if vect == lapack.ApplyP, side == blas.Left, and trans == blas.NoTrans
//  C * P if vect == lapack.ApplyP, side == blas.Left, and trans == blas.NoTrans
//  P^T * C if vect == lapack.ApplyP, side == blas.Right, and trans == blas.Trans
//  C * P^T if vect == lapack.ApplyP, side == blas.Right, and trans == blas.Trans
// where P and Q are the orthogonal matrices determined by Dgebrd, A = Q * B * P^T.
// See Dgebrd for the definitions of Q and P.
//
// If vect == lapack.ApplyQ, A is assumed to have been an nq×k matrix, while if
// vect == lapack.ApplyP, A is assumed to have been a k×nq matrix. nq = m if
// side == blas.Left, while nq = n if side == blas.Right.
//
// C is an m×n matrix. On exit it is updated by the multiplication listed above.
//
// Tau must have length min(nq,k), and Dormbr will panic otherwise. Tau contains
// the elementary reflectors to construct Q or P depending on the value of
// vect.
func (impl Implementation) Dormbr(vect lapack.DecompUpdate, side blas.Side, trans blas.Transpose, m, n, k int, a []float64, lda int, tau, c []float64, ldc int, work []float64, lwork int) {
	if side != blas.Left && side != blas.Right {
		panic(badSide)
	}
	if trans != blas.NoTrans && trans != blas.Trans {
		panic(badTrans)
	}
	if vect != lapack.ApplyP && vect != lapack.ApplyQ {
		panic(badDecompUpdate)
	}
	nq := n
	if side == blas.Left {
		nq = m
	}
	if vect == lapack.ApplyQ {
		checkMatrix(nq, min(nq, k), a, lda)
	} else {
		checkMatrix(min(nq, k), nq, a, lda)
	}
	clapack.Dormbr(byte(vect), side, trans, m, n, k, a, lda, tau, c, ldc)
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

// Dtrcon estimates the reciprocal of the condition number of a positive-definite
// matrix A given the Cholesky decmposition of A. The condition number computed
// is based on the 1-norm and the ∞-norm.
//
// anorm is the 1-norm and the ∞-norm of the original matrix A.
//
// work is a temporary data slice of length at least 3*n and Dpocon will panic otherwise.
//
// iwork is a temporary data slice of length at least n and Dpocon will panic otherwise.
func (impl Implementation) Dpocon(uplo blas.Uplo, n int, a []float64, lda int, anorm float64, work []float64, iwork []int) float64 {
	checkMatrix(n, n, a, lda)
	if uplo != blas.Upper && uplo != blas.Lower {
		panic(badUplo)
	}
	if len(work) < 3*n {
		panic(badWork)
	}
	if len(iwork) < n {
		panic(badWork)
	}
	rcond := make([]float64, 1)
	clapack.Dpocon(uplo, n, a, lda, anorm, rcond)
	return rcond[0]
}

// Dtrcon estimates the reciprocal of the condition number of a triangular matrix A.
// The condition number computed may be based on the 1-norm or the ∞-norm.
//
// work is a temporary data slice of length at least 3*n and Dtrcon will panic otherwise.
//
// iwork is a temporary data slice of length at least n and Dtrcon will panic otherwise.
func (impl Implementation) Dtrcon(norm lapack.MatrixNorm, uplo blas.Uplo, diag blas.Diag, n int, a []float64, lda int, work []float64, iwork []int) float64 {
	if norm != lapack.MaxColumnSum && norm != lapack.MaxRowSum {
		panic(badNorm)
	}
	if uplo != blas.Upper && uplo != blas.Lower {
		panic(badUplo)
	}
	if diag != blas.NonUnit && diag != blas.Unit {
		panic(badDiag)
	}
	if len(work) < 3*n {
		panic(badWork)
	}
	if len(iwork) < n {
		panic(badWork)
	}
	rcond := []float64{0}
	clapack.Dtrcon(byte(norm), uplo, diag, n, a, lda, rcond)
	return rcond[0]
}

// Dtrtri computes the inverse of a triangular matrix, storing the result in place
// into a. This is the BLAS level 3 version of the algorithm which builds upon
// Dtrti2 to operate on matrix blocks instead of only individual columns.
//
// Dtrti returns whether the matrix a is singular or whether it's not singular.
// If the matrix is singular the inversion is not performed.
func (impl Implementation) Dtrtri(uplo blas.Uplo, diag blas.Diag, n int, a []float64, lda int) (ok bool) {
	checkMatrix(n, n, a, lda)
	if uplo != blas.Upper && uplo != blas.Lower {
		panic(badUplo)
	}
	if diag != blas.NonUnit && diag != blas.Unit {
		panic(badDiag)
	}
	return clapack.Dtrtri(uplo, diag, n, a, lda)
}

// Dtrtrs solves a triangular system of the form A * X = B or A^T * X = B. Dtrtrs
// returns whether the solve completed successfully. If A is singular, no solve is performed.
func (impl Implementation) Dtrtrs(uplo blas.Uplo, trans blas.Transpose, diag blas.Diag, n, nrhs int, a []float64, lda int, b []float64, ldb int) (ok bool) {
	return clapack.Dtrtrs(uplo, trans, diag, n, nrhs, a, lda, b, ldb)
}
