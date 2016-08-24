// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"github.com/gonum/blas"
	"github.com/gonum/lapack"
)

// Dormqr multiplies the matrix C by the orthogonal matrix Q defined by the
// slices a and tau. A and tau are as returned from Dgeqrf.
//  C = Q * C    if side == blas.Left and trans == blas.NoTrans
//  C = Q^T * C  if side == blas.Left and trans == blas.Trans
//  C = C * Q    if side == blas.Right and trans == blas.NoTrans
//  C = C * Q^T  if side == blas.Right and trans == blas.Trans
// If side == blas.Left, A is a matrix of side m×k, and if side == blas.Right
// A is of size n×k. This uses a blocked algorithm.
//
// work is temporary storage, and lwork specifies the usable memory length.
// At minimum, lwork >= m if side == blas.Left and lwork >= n if side == blas.Right,
// and this function will panic otherwise.
// Dormqr uses a block algorithm, but the block size is limited
// by the temporary space available. If lwork == -1, instead of performing Dormqr,
// the optimal work length will be stored into work[0].
//
// tau contains the Householder scales and must have length at least k, and
// this function will panic otherwise.
func (impl Implementation) Dormqr(side blas.Side, trans blas.Transpose, m, n, k int, a []float64, lda int, tau, c []float64, ldc int, work []float64, lwork int) {
	var nq, nw int
	switch side {
	default:
		panic(badSide)
	case blas.Left:
		nq = m
		nw = n
	case blas.Right:
		nq = n
		nw = m
	}
	switch {
	case trans != blas.NoTrans && trans != blas.Trans:
		panic(badTrans)
	case m < 0 || n < 0:
		panic(negDimension)
	case k < 0 || nq < k:
		panic("lapack: invalid value of k")
	case len(work) < lwork:
		panic(shortWork)
	case lwork < max(1, nw) && lwork != -1:
		panic(badWork)
	}
	if lwork != -1 {
		checkMatrix(nq, k, a, lda)
		checkMatrix(m, n, c, ldc)
		if len(tau) != k {
			panic(badTau)
		}
	}

	if m == 0 || n == 0 || k == 0 {
		work[0] = 1
		return
	}

	const (
		nbmax = 64
		ldt   = nbmax
		tsize = nbmax * ldt
	)
	opts := string(side) + string(trans)
	nb := min(nbmax, impl.Ilaenv(1, "DORMQR", opts, m, n, k, -1))
	lworkopt := max(1, nw)*nb + tsize
	if lwork == -1 {
		work[0] = float64(lworkopt)
		return
	}

	nbmin := 2
	if 1 < nb && nb < k {
		if lwork < nw*nb+tsize {
			nb = (lwork - tsize) / nw
			nbmin = max(2, impl.Ilaenv(2, "DORMQR", opts, m, n, k, -1))
		}
	}

	if nb < nbmin || k <= nb {
		// Call unblocked code.
		impl.Dorm2r(side, trans, m, n, k, a, lda, tau, c, ldc, work)
		work[0] = float64(lworkopt)
		return
	}

	var (
		ldwork = nb
		left   = side == blas.Left
		notran = trans == blas.NoTrans
	)
	switch {
	case left && notran:
		for i := ((k - 1) / nb) * nb; i >= 0; i -= nb {
			ib := min(nb, k-i)
			impl.Dlarft(lapack.Forward, lapack.ColumnWise, m-i, ib,
				a[i*lda+i:], lda,
				tau[i:],
				work[:tsize], ldt)
			impl.Dlarfb(side, trans, lapack.Forward, lapack.ColumnWise, m-i, n, ib,
				a[i*lda+i:], lda,
				work[:tsize], ldt,
				c[i*ldc:], ldc,
				work[tsize:], ldwork)
		}

	case left && !notran:
		for i := 0; i < k; i += nb {
			ib := min(nb, k-i)
			impl.Dlarft(lapack.Forward, lapack.ColumnWise, m-i, ib,
				a[i*lda+i:], lda,
				tau[i:],
				work[:tsize], ldt)
			impl.Dlarfb(side, trans, lapack.Forward, lapack.ColumnWise, m-i, n, ib,
				a[i*lda+i:], lda,
				work[:tsize], ldt,
				c[i*ldc:], ldc,
				work[tsize:], ldwork)
		}

	case !left && notran:
		for i := 0; i < k; i += nb {
			ib := min(nb, k-i)
			impl.Dlarft(lapack.Forward, lapack.ColumnWise, n-i, ib,
				a[i*lda+i:], lda,
				tau[i:],
				work[:tsize], ldt)
			impl.Dlarfb(side, trans, lapack.Forward, lapack.ColumnWise, m, n-i, ib,
				a[i*lda+i:], lda,
				work[:tsize], ldt,
				c[i:], ldc,
				work[tsize:], ldwork)
		}

	case !left && !notran:
		for i := ((k - 1) / nb) * nb; i >= 0; i -= nb {
			ib := min(nb, k-i)
			impl.Dlarft(lapack.Forward, lapack.ColumnWise, n-i, ib,
				a[i*lda+i:], lda,
				tau[i:],
				work[:tsize], ldt)
			impl.Dlarfb(side, trans, lapack.Forward, lapack.ColumnWise, m, n-i, ib,
				a[i*lda+i:], lda,
				work[:tsize], ldt,
				c[i:], ldc,
				work[tsize:], ldwork)
		}
	}
	work[0] = float64(lworkopt)
}
