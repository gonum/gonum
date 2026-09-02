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

// Zlarft forms the triangular factor T of a complex block reflector H of order
// n, which is defined as a product of k elementary reflectors. T is stored in
// t.
//
//	H = I - V * T * V^H  if store == lapack.ColumnWise
//	H = I - V^H * T * V  if store == lapack.RowWise
//
// H is defined by a product of the elementary reflectors where
//
//	H = H_0 * H_1 * ... * H_{k-1}  if direct == lapack.Forward
//	H = H_{k-1} * ... * H_1 * H_0  if direct == lapack.Backward
//
// t is a k×k triangular matrix. t is upper triangular if direct == lapack.Forward
// and lower triangular otherwise. This function will panic if t is not of
// sufficient size.
//
// store describes the storage of the elementary reflectors in v. See Zlarfb for
// a description of the storage layouts.
//
// tau contains the scalar factors of the elementary reflectors H_i.
//
// Zlarft is an internal routine. It is exported for testing purposes.
func (impl Implementation) Zlarft(direct lapack.Direct, store lapack.StoreV, n, k int, v []complex128, ldv int, tau []complex128, t []complex128, ldt int) {
	mv, nv := n, k
	if store == lapack.RowWise {
		mv, nv = k, n
	}
	switch {
	case direct != lapack.Forward && direct != lapack.Backward:
		panic(badDirect)
	case store != lapack.RowWise && store != lapack.ColumnWise:
		panic(badStoreV)
	case n < 0:
		panic(nLT0)
	case k < 1:
		panic(kLT1)
	case ldv < max(1, nv):
		panic(badLdV)
	case len(tau) < k:
		panic(shortTau)
	case ldt < max(1, k):
		panic(shortT)
	}

	if n == 0 {
		return
	}

	switch {
	case len(v) < (mv-1)*ldv+nv:
		panic(shortV)
	case len(t) < (k-1)*ldt+k:
		panic(shortT)
	}

	bi := cblas128.Implementation()
	one := complex(1, 0)

	if direct == lapack.Forward {
		for i := 0; i < k; i++ {
			if tau[i] == 0 {
				// H(i) = I.
				for j := 0; j <= i; j++ {
					t[j*ldt+i] = 0
				}
				continue
			}
			if store == lapack.ColumnWise {
				// T(0:i, i) = -tau(i) * V(i:n, 0:i)^H * V(i:n, i)
				// Temporarily set V(i,i) = 1 to include the unit diagonal.
				vii := v[i*ldv+i]
				v[i*ldv+i] = one
				if i > 0 {
					bi.Zgemv(blas.ConjTrans, n-i, i,
						-tau[i], v[i*ldv:], ldv, v[i*ldv+i:], ldv,
						0, t[i:], ldt)
				}
				v[i*ldv+i] = vii
			} else {
				// Row-wise.
				// T(0:i, i) = -tau(i) * V(0:i, i:n) * V(i, i:n)^H
				vii := v[i*ldv+i]
				v[i*ldv+i] = one
				// Conjugate V(i, i+1:n), preserving unit at V(i,i).
				if i < n-1 {
					zlacgv(n-i-1, v[i*ldv+i+1:], 1)
				}
				if i > 0 {
					bi.Zgemv(blas.NoTrans, i, n-i,
						-tau[i], v[i:], ldv, v[i*ldv+i:], 1,
						0, t[i:], ldt)
				}
				if i < n-1 {
					zlacgv(n-i-1, v[i*ldv+i+1:], 1)
				}
				v[i*ldv+i] = vii
			}
			// T(0:i, i) = T(0:i, 0:i) * T(0:i, i).
			if i > 0 {
				bi.Ztrmv(blas.Upper, blas.NoTrans, blas.NonUnit, i, t, ldt, t[i:], ldt)
			}
			t[i*ldt+i] = tau[i]
		}
		return
	}
	// direct == lapack.Backward.
	for i := k - 1; i >= 0; i-- {
		if tau[i] == 0 {
			for j := i; j < k; j++ {
				t[j*ldt+i] = 0
			}
			continue
		}
		if i < k-1 {
			if store == lapack.ColumnWise {
				// T(i+1:k, i) = -tau(i) * V(0:n-k+i+1, i+1:k)^H * V(0:n-k+i+1, i)
				vii := v[(n-k+i)*ldv+i]
				v[(n-k+i)*ldv+i] = one
				bi.Zgemv(blas.ConjTrans, n-k+i+1, k-i-1,
					-tau[i], v[i+1:], ldv, v[i:], ldv,
					0, t[(i+1)*ldt+i:], ldt)
				v[(n-k+i)*ldv+i] = vii
			} else {
				// T(i+1:k, i) = -tau(i) * V(i+1:k, 0:n-k+i+1) * V(i, 0:n-k+i+1)^H
				vii := v[i*ldv+n-k+i]
				v[i*ldv+n-k+i] = one
				// Conjugate V(i, 0:n-k+i) (of length n-k+i).
				if n-k+i > 0 {
					zlacgv(n-k+i, v[i*ldv:], 1)
				}
				bi.Zgemv(blas.NoTrans, k-i-1, n-k+i+1,
					-tau[i], v[(i+1)*ldv:], ldv, v[i*ldv:], 1,
					0, t[(i+1)*ldt+i:], ldt)
				if n-k+i > 0 {
					zlacgv(n-k+i, v[i*ldv:], 1)
				}
				v[i*ldv+n-k+i] = vii
			}
			// T(i+1:k, i) = T(i+1:k, i+1:k) * T(i+1:k, i).
			bi.Ztrmv(blas.Lower, blas.NoTrans, blas.NonUnit, k-i-1,
				t[(i+1)*ldt+i+1:], ldt,
				t[(i+1)*ldt+i:], ldt)
		}
		t[i*ldt+i] = tau[i]
	}
}

// zlacgv conjugates a complex vector of length n in place, with stride incx.
// This is the LAPACK reference ZLACGV routine.
func zlacgv(n int, x []complex128, incx int) {
	if n <= 0 {
		return
	}
	ix := 0
	for i := 0; i < n; i++ {
		x[ix] = cmplx.Conj(x[ix])
		ix += incx
	}
}
