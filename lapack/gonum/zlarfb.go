// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math/cmplx"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/cblas128"
	"gonum.org/v1/gonum/lapack"
)

// Zlarfb applies a complex block reflector H or its conjugate transpose H^H to
// a complex m×n matrix C.
//
//	C = H   * C  if side == blas.Left  and trans == blas.NoTrans
//	C = C   * H  if side == blas.Right and trans == blas.NoTrans
//	C = H^H * C  if side == blas.Left  and trans == blas.ConjTrans
//	C = C * H^H  if side == blas.Right and trans == blas.ConjTrans
//
// H is a product of elementary reflectors. direct sets the direction of
// multiplication:
//
//	H = H_0 * H_1 * ... * H_{k-1}  if direct == lapack.Forward
//	H = H_{k-1} * ... * H_1 * H_0  if direct == lapack.Backward
//
// The combination of direct and store defines the orientation of the elementary
// reflectors. See Dlarfb for the V layout described with the unit diagonal
// elements implicitly represented.
//
// t is a k×k matrix containing the triangular factor T. See Zlarft.
//
// work is a temporary storage matrix with stride ldwork. It must be at least
// of size n×k if side == blas.Left and m×k if side == blas.Right.
//
// Zlarfb is an internal routine. It is exported for testing purposes.
func (impl Implementation) Zlarfb(side blas.Side, trans blas.Transpose, direct lapack.Direct, store lapack.StoreV, m, n, k int, v []complex128, ldv int, t []complex128, ldt int, c []complex128, ldc int, work []complex128, ldwork int) {
	nv := m
	if side == blas.Right {
		nv = n
	}
	switch {
	case side != blas.Left && side != blas.Right:
		panic(badSide)
	case trans != blas.NoTrans && trans != blas.ConjTrans:
		panic(badTrans)
	case direct != lapack.Forward && direct != lapack.Backward:
		panic(badDirect)
	case store != lapack.ColumnWise && store != lapack.RowWise:
		panic(badStoreV)
	case m < 0:
		panic(mLT0)
	case n < 0:
		panic(nLT0)
	case k < 0:
		panic(kLT0)
	case store == lapack.ColumnWise && ldv < max(1, k):
		panic(badLdV)
	case store == lapack.RowWise && ldv < max(1, nv):
		panic(badLdV)
	case ldt < max(1, k):
		panic(badLdT)
	case ldc < max(1, n):
		panic(badLdC)
	case ldwork < max(1, k):
		panic(badLdWork)
	}

	if m == 0 || n == 0 {
		return
	}

	nw := n
	if side == blas.Right {
		nw = m
	}
	switch {
	case store == lapack.ColumnWise && len(v) < (nv-1)*ldv+k:
		panic(shortV)
	case store == lapack.RowWise && len(v) < (k-1)*ldv+nv:
		panic(shortV)
	case len(t) < (k-1)*ldt+k:
		panic(shortT)
	case len(c) < (m-1)*ldc+n:
		panic(shortC)
	case len(work) < (nw-1)*ldwork+k:
		panic(shortWork)
	}

	bi := cblas128.Implementation()

	transt := blas.ConjTrans
	if trans == blas.ConjTrans {
		transt = blas.NoTrans
	}
	one := complex(1, 0)
	negOne := complex(-1, 0)

	if store == lapack.ColumnWise {
		if direct == lapack.Forward {
			// V1 is the first k rows of C. V2 is the remaining.
			if side == blas.Left {
				// W = C^H * V = C1^H * V1 + C2^H * V2 (stored in work).

				// W = C1^H.  W has shape n×k. W[i, j] = conj(C1[j, i]).
				// Copy then conjugate each row of W.
				for j := 0; j < k; j++ {
					bi.Zcopy(n, c[j*ldc:], 1, work[j:], ldwork)
				}
				// Conjugate work in-place (all n×k elements).
				zlacgvMat(n, k, work, ldwork)
				// W = W * V1 (V1 is the upper k×k of V, unit lower triangular).
				bi.Ztrmm(blas.Right, blas.Lower, blas.NoTrans, blas.Unit,
					n, k, one, v, ldv, work, ldwork)
				if m > k {
					// W += C2^H * V2.
					bi.Zgemm(blas.ConjTrans, blas.NoTrans, n, k, m-k,
						one, c[k*ldc:], ldc, v[k*ldv:], ldv,
						one, work, ldwork)
				}
				// W = W * T^op.
				bi.Ztrmm(blas.Right, blas.Upper, transt, blas.NonUnit,
					n, k, one, t, ldt, work, ldwork)
				// C -= V * W^H.
				if m > k {
					// C2 -= V2 * W^H.
					bi.Zgemm(blas.NoTrans, blas.ConjTrans, m-k, n, k,
						negOne, v[k*ldv:], ldv, work, ldwork,
						one, c[k*ldc:], ldc)
				}
				// W = W * V1^H.
				bi.Ztrmm(blas.Right, blas.Lower, blas.ConjTrans, blas.Unit,
					n, k, one, v, ldv, work, ldwork)
				// C1 -= W^H.
				for i := 0; i < n; i++ {
					for j := 0; j < k; j++ {
						c[j*ldc+i] -= cmplx.Conj(work[i*ldwork+j])
					}
				}
				return
			}
			// side == blas.Right. Form C = C * H or C * H^H where C = (C1 C2).
			// W = C * V.

			// W = C1. W has shape m×k. W[i, j] = C1[i, j].
			for j := 0; j < k; j++ {
				bi.Zcopy(m, c[j:], ldc, work[j:], ldwork)
			}
			// W = W * V1.
			bi.Ztrmm(blas.Right, blas.Lower, blas.NoTrans, blas.Unit,
				m, k, one, v, ldv, work, ldwork)
			if n > k {
				bi.Zgemm(blas.NoTrans, blas.NoTrans, m, k, n-k,
					one, c[k:], ldc, v[k*ldv:], ldv,
					one, work, ldwork)
			}
			// W = W * T^op (Right variant uses plain trans, not transt).
			bi.Ztrmm(blas.Right, blas.Upper, trans, blas.NonUnit,
				m, k, one, t, ldt, work, ldwork)
			if n > k {
				bi.Zgemm(blas.NoTrans, blas.ConjTrans, m, n-k, k,
					negOne, work, ldwork, v[k*ldv:], ldv,
					one, c[k:], ldc)
			}
			bi.Ztrmm(blas.Right, blas.Lower, blas.ConjTrans, blas.Unit,
				m, k, one, v, ldv, work, ldwork)
			for i := 0; i < m; i++ {
				for j := 0; j < k; j++ {
					c[i*ldc+j] -= work[i*ldwork+j]
				}
			}
			return
		}
		// direct == lapack.Backward, store = Columnwise.
		// V = (V1; V2) where V2 is unit upper triangular, last k rows.
		if side == blas.Left {
			// W = C2^H.
			for j := 0; j < k; j++ {
				bi.Zcopy(n, c[(m-k+j)*ldc:], 1, work[j:], ldwork)
			}
			zlacgvMat(n, k, work, ldwork)
			// W = W * V2.
			bi.Ztrmm(blas.Right, blas.Upper, blas.NoTrans, blas.Unit,
				n, k, one, v[(m-k)*ldv:], ldv, work, ldwork)
			if m > k {
				bi.Zgemm(blas.ConjTrans, blas.NoTrans, n, k, m-k,
					one, c, ldc, v, ldv,
					one, work, ldwork)
			}
			bi.Ztrmm(blas.Right, blas.Lower, transt, blas.NonUnit,
				n, k, one, t, ldt, work, ldwork)
			if m > k {
				bi.Zgemm(blas.NoTrans, blas.ConjTrans, m-k, n, k,
					negOne, v, ldv, work, ldwork,
					one, c, ldc)
			}
			bi.Ztrmm(blas.Right, blas.Upper, blas.ConjTrans, blas.Unit,
				n, k, one, v[(m-k)*ldv:], ldv, work, ldwork)
			for i := 0; i < n; i++ {
				for j := 0; j < k; j++ {
					c[(m-k+j)*ldc+i] -= cmplx.Conj(work[i*ldwork+j])
				}
			}
			return
		}
		// Right, backward, columnwise.
		// W = C2.
		for j := 0; j < k; j++ {
			bi.Zcopy(m, c[n-k+j:], ldc, work[j:], ldwork)
		}
		bi.Ztrmm(blas.Right, blas.Upper, blas.NoTrans, blas.Unit,
			m, k, one, v[(n-k)*ldv:], ldv, work, ldwork)
		if n > k {
			bi.Zgemm(blas.NoTrans, blas.NoTrans, m, k, n-k,
				one, c, ldc, v, ldv,
				one, work, ldwork)
		}
		bi.Ztrmm(blas.Right, blas.Lower, trans, blas.NonUnit,
			m, k, one, t, ldt, work, ldwork)
		if n > k {
			bi.Zgemm(blas.NoTrans, blas.ConjTrans, m, n-k, k,
				negOne, work, ldwork, v, ldv,
				one, c, ldc)
		}
		bi.Ztrmm(blas.Right, blas.Upper, blas.ConjTrans, blas.Unit,
			m, k, one, v[(n-k)*ldv:], ldv, work, ldwork)
		for i := 0; i < m; i++ {
			for j := 0; j < k; j++ {
				c[i*ldc+n-k+j] -= work[i*ldwork+j]
			}
		}
		return
	}
	// store == lapack.RowWise.
	if direct == lapack.Forward {
		// V = (V1 V2) where V1 is unit upper triangular.
		if side == blas.Left {
			// W = C1^H.
			for j := 0; j < k; j++ {
				bi.Zcopy(n, c[j*ldc:], 1, work[j:], ldwork)
			}
			zlacgvMat(n, k, work, ldwork)
			// W = W * V1^H.
			bi.Ztrmm(blas.Right, blas.Upper, blas.ConjTrans, blas.Unit,
				n, k, one, v, ldv, work, ldwork)
			if m > k {
				// W += C2^H * V2^H.
				bi.Zgemm(blas.ConjTrans, blas.ConjTrans, n, k, m-k,
					one, c[k*ldc:], ldc, v[k:], ldv,
					one, work, ldwork)
			}
			bi.Ztrmm(blas.Right, blas.Upper, transt, blas.NonUnit,
				n, k, one, t, ldt, work, ldwork)
			if m > k {
				// C2 -= V2^H * W^H  -> but netlib has ZGEMM('ConjTrans','ConjTrans') for this.
				bi.Zgemm(blas.ConjTrans, blas.ConjTrans, m-k, n, k,
					negOne, v[k:], ldv, work, ldwork,
					one, c[k*ldc:], ldc)
			}
			// W *= V1.
			bi.Ztrmm(blas.Right, blas.Upper, blas.NoTrans, blas.Unit,
				n, k, one, v, ldv, work, ldwork)
			for i := 0; i < n; i++ {
				for j := 0; j < k; j++ {
					c[j*ldc+i] -= cmplx.Conj(work[i*ldwork+j])
				}
			}
			return
		}
		// Right, forward, rowwise.
		// W = C1.
		for j := 0; j < k; j++ {
			bi.Zcopy(m, c[j:], ldc, work[j:], ldwork)
		}
		bi.Ztrmm(blas.Right, blas.Upper, blas.ConjTrans, blas.Unit,
			m, k, one, v, ldv, work, ldwork)
		if n > k {
			bi.Zgemm(blas.NoTrans, blas.ConjTrans, m, k, n-k,
				one, c[k:], ldc, v[k:], ldv,
				one, work, ldwork)
		}
		bi.Ztrmm(blas.Right, blas.Upper, trans, blas.NonUnit,
			m, k, one, t, ldt, work, ldwork)
		if n > k {
			bi.Zgemm(blas.NoTrans, blas.NoTrans, m, n-k, k,
				negOne, work, ldwork, v[k:], ldv,
				one, c[k:], ldc)
		}
		bi.Ztrmm(blas.Right, blas.Upper, blas.NoTrans, blas.Unit,
			m, k, one, v, ldv, work, ldwork)
		for i := 0; i < m; i++ {
			for j := 0; j < k; j++ {
				c[i*ldc+j] -= work[i*ldwork+j]
			}
		}
		return
	}
	// direct == Backward, store == Rowwise.
	// V = (V1 V2) where V2 is unit lower triangular, last k columns.
	if side == blas.Left {
		// W = C2^H.
		for j := 0; j < k; j++ {
			bi.Zcopy(n, c[(m-k+j)*ldc:], 1, work[j:], ldwork)
		}
		zlacgvMat(n, k, work, ldwork)
		bi.Ztrmm(blas.Right, blas.Lower, blas.ConjTrans, blas.Unit,
			n, k, one, v[m-k:], ldv, work, ldwork)
		if m > k {
			bi.Zgemm(blas.ConjTrans, blas.ConjTrans, n, k, m-k,
				one, c, ldc, v, ldv,
				one, work, ldwork)
		}
		bi.Ztrmm(blas.Right, blas.Lower, transt, blas.NonUnit,
			n, k, one, t, ldt, work, ldwork)
		if m > k {
			bi.Zgemm(blas.ConjTrans, blas.ConjTrans, m-k, n, k,
				negOne, v, ldv, work, ldwork,
				one, c, ldc)
		}
		bi.Ztrmm(blas.Right, blas.Lower, blas.NoTrans, blas.Unit,
			n, k, one, v[m-k:], ldv, work, ldwork)
		for i := 0; i < n; i++ {
			for j := 0; j < k; j++ {
				c[(m-k+j)*ldc+i] -= cmplx.Conj(work[i*ldwork+j])
			}
		}
		return
	}
	// Right, backward, rowwise.
	// W = C2.
	for j := 0; j < k; j++ {
		bi.Zcopy(m, c[n-k+j:], ldc, work[j:], ldwork)
	}
	bi.Ztrmm(blas.Right, blas.Lower, blas.ConjTrans, blas.Unit,
		m, k, one, v[n-k:], ldv, work, ldwork)
	if n > k {
		bi.Zgemm(blas.NoTrans, blas.ConjTrans, m, k, n-k,
			one, c, ldc, v, ldv,
			one, work, ldwork)
	}
	bi.Ztrmm(blas.Right, blas.Lower, trans, blas.NonUnit,
		m, k, one, t, ldt, work, ldwork)
	if n > k {
		bi.Zgemm(blas.NoTrans, blas.NoTrans, m, n-k, k,
			negOne, work, ldwork, v, ldv,
			one, c, ldc)
	}
	bi.Ztrmm(blas.Right, blas.Lower, blas.NoTrans, blas.Unit,
		m, k, one, v[n-k:], ldv, work, ldwork)
	for i := 0; i < m; i++ {
		for j := 0; j < k; j++ {
			c[i*ldc+n-k+j] -= work[i*ldwork+j]
		}
	}
}

// zlacgvMat conjugates every element of an m×n complex matrix stored with
// leading dimension ld, in place.
func zlacgvMat(m, n int, a []complex128, ld int) {
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			a[i*ld+j] = cmplx.Conj(a[i*ld+j])
		}
	}
}
