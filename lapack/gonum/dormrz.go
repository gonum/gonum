// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

// Dormrz multiplies an m×n matrix C by the orthogonal matrix Z from a
// RZ factorization determined by Dtzrzf.
//
//	C = Z * C    if side == blas.Left  and trans == blas.NoTrans,
//	C = Zᵀ * C   if side == blas.Left  and trans == blas.Trans,
//	C = C * Z    if side == blas.Right and trans == blas.NoTrans,
//	C = C * Zᵀ   if side == blas.Right and trans == blas.Trans,
//
// where Z is defined as the product of k elementary reflectors
//
//	Z = Z(0) * Z(1) * ... * Z(k-1)
//
// as returned by Dtzrzf. Z is of order m if side == blas.Left
// and of order n if side == blas.Right.
//
// Each Z(i) has the form
//
//	Z(i) = I - tau[i] * v * vᵀ
//
// where tau[i] is stored in tau[i] and v is a vector with
//
//	v[0:i] = 0, v[i] = 1, v[i+1:nq-l] = 0, v[nq-l:nq] = A[i, nq-l:nq]
//
// where nq = m if side == blas.Left and nq = n if side == blas.Right.
// l is the number of columns of the matrix V containing the meaningful
// part of the reflectors.
//
// tau must have length at least k. work must have length at least max(1,lwork).
//
// work is temporary storage, and lwork specifies the usable memory length. At
// minimum, lwork >= n if side == blas.Left and lwork >= m if side ==
// blas.Right, and this function will panic otherwise.
//
// If lwork is -1, instead of performing Dormrz, the optimal workspace size will
// be stored into work[0].
//
// Dormrz is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dormrz(side blas.Side, trans blas.Transpose, m, n, k, l int, a []float64, lda int, tau, c []float64, ldc int, work []float64, lwork int) {
	left := side == blas.Left
	nq := n
	nw := m
	if left {
		nq = m
		nw = n
	}
	switch {
	case !left && side != blas.Right:
		panic(badSide)
	case trans != blas.NoTrans && trans != blas.Trans:
		panic(badTrans)
	case m < 0:
		panic(mLT0)
	case n < 0:
		panic(nLT0)
	case k < 0:
		panic(kLT0)
	case left && k > m:
		panic(kGTM)
	case !left && k > n:
		panic(kGTN)
	case l < 0 || l > nq-k:
		panic("lapack: l out of range")
	case lda < max(1, nq):
		panic(badLdA)
	case ldc < max(1, n):
		panic(badLdC)
	case lwork < max(1, nw) && lwork != -1:
		panic(badLWork)
	case len(work) < max(1, lwork):
		panic(shortWork)
	}

	if lwork == -1 {
		work[0] = float64(max(1, nw))
		return
	}

	if m == 0 || n == 0 || k == 0 {
		work[0] = 1
		return
	}

	switch {
	case len(a) < (k-1)*lda+nq:
		panic(shortA)
	case len(tau) < k:
		panic(shortTau)
	case len(c) < (m-1)*ldc+n:
		panic(shortC)
	}

	bi := blas64.Implementation()

	// Determine loop direction based on side and trans.
	// Left,  NoTrans (Z*C):   apply Z(k-1)...Z(0) → loop i=k-1 downto 0
	// Left,  Trans   (Zᵀ*C):  apply Z(0)...Z(k-1) → loop i=0 to k-1
	// Right, NoTrans (C*Z):   apply Z(0)...Z(k-1) → loop i=0 to k-1
	// Right, Trans   (C*Zᵀ):  apply Z(k-1)...Z(0) → loop i=k-1 downto 0
	i0, iend, istep := 0, k, 1
	if (left && trans == blas.NoTrans) || (!left && trans == blas.Trans) {
		i0, iend, istep = k-1, -1, -1
	}

	for i := i0; i != iend; i += istep {
		if tau[i] == 0 {
			continue
		}
		if left {
			// Apply Z(i) from the left: C := (I - tau*v*vᵀ) * C
			// v[i]=1, v[m-l:m] = A[i, nq-l:nq], rest zero.
			// w = C[i,:] + A[i,m-l:m]ᵀ * C[m-l:m,:]  (length n)
			bi.Dcopy(n, c[i*ldc:], 1, work, 1)
			bi.Dgemv(blas.Trans, l, n, 1, c[(m-l)*ldc:], ldc, a[i*lda+(nq-l):], 1, 1, work, 1)
			// C[i,:] -= tau * w
			bi.Daxpy(n, -tau[i], work, 1, c[i*ldc:], 1)
			// C[m-l:m,:] -= tau * v * wᵀ
			bi.Dger(l, n, -tau[i], a[i*lda+(nq-l):], 1, work, 1, c[(m-l)*ldc:], ldc)
		} else {
			// Apply Z(i) from the right: C := C * (I - tau*v*vᵀ)
			// v[i]=1, v[n-l:n] = A[i, nq-l:nq], rest zero.
			// w = C[:,i] + C[:,n-l:n] * A[i,n-l:n]ᵀ  (length m)
			bi.Dcopy(m, c[i:], ldc, work, 1)
			bi.Dgemv(blas.NoTrans, m, l, 1, c[n-l:], ldc, a[i*lda+(nq-l):], 1, 1, work, 1)
			// C[:,i] -= tau * w
			bi.Daxpy(m, -tau[i], work, 1, c[i:], ldc)
			// C[:,n-l:n] -= tau * w * vᵀ
			bi.Dger(m, l, -tau[i], work, 1, a[i*lda+(nq-l):], 1, c[n-l:], ldc)
		}
	}
	work[0] = float64(max(1, nw))
}
