// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/cmplx"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/lapack"
)

type Zlarfter interface {
	Zlarfgerhelper
	Zlarft(direct lapack.Direct, store lapack.StoreV, n, k int,
		v []complex128, ldv int, tau []complex128, t []complex128, ldt int)
}

// Zlarfgerhelper lets Zlarft tests generate reflectors via Zlarfg.
type Zlarfgerhelper interface {
	Zlarfg(n int, alpha complex128, x []complex128, incX int) (beta float64, tau complex128)
}

// ZlarftTest verifies that the T produced by Zlarft satisfies
//
//	H = I - V * T * V^H     (columnwise)
//	H = I - V^H * T * V     (rowwise)
//
// where H is the product of the k reflectors stored in V/tau.
func ZlarftTest(t *testing.T, impl Zlarfter) {
	const tol = 1e-12
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, store := range []lapack.StoreV{lapack.ColumnWise, lapack.RowWise} {
		for _, direct := range []lapack.Direct{lapack.Forward, lapack.Backward} {
			for _, sz := range []struct{ n, k int }{
				{5, 3}, {6, 4}, {7, 5}, {4, 4}, {3, 2}, {10, 6},
			} {
				name := fmt.Sprintf("store=%c dir=%c n=%d k=%d",
					storeChar(store), directChar(direct), sz.n, sz.k)
				zlarftCheck(t, impl, rnd, direct, store, sz.n, sz.k, name, tol)
			}
		}
	}
}

func zlarftCheck(t *testing.T, impl Zlarfter, rnd *rand.Rand, direct lapack.Direct,
	store lapack.StoreV, n, k int, name string, tol float64) {
	if k > n {
		return
	}
	// Build V / tau.
	var ldv int
	var v []complex128
	if store == lapack.ColumnWise {
		ldv = k
		v = make([]complex128, n*ldv)
	} else {
		ldv = n
		v = make([]complex128, k*ldv)
	}
	tau := make([]complex128, k)
	// Keep full-length v_i for assembling H directly.
	vFull := make([][]complex128, k)

	for i := 0; i < k; i++ {
		var length int
		if direct == lapack.Forward {
			length = n - i
		} else {
			length = n - (k - 1 - i)
		}
		if length < 1 {
			length = 1
		}
		xvec := make([]complex128, length)
		for j := range xvec {
			xvec[j] = complex(rnd.NormFloat64(), rnd.NormFloat64())
		}
		_, tauI := impl.Zlarfg(length, xvec[0], xvec[1:], 1)
		tau[i] = tauI

		// Build the full-length reflector vector u_i (the "column" form). For
		// ColumnWise storage V[:, i] stores u_i directly; for RowWise storage
		// V[i, :] stores conj(u_i^T) -- the LAPACK convention so that
		// H = I - V^H T V holds with u_i as the reflector column.
		vi := make([]complex128, n)
		if direct == lapack.Forward {
			vi[i] = 1
			for j := 0; j < length-1; j++ {
				vi[i+1+j] = xvec[1+j]
			}
			if store == lapack.ColumnWise {
				for r := 0; r < n; r++ {
					v[r*ldv+i] = 0
				}
				v[i*ldv+i] = 1
				for j := 0; j < length-1; j++ {
					v[(i+1+j)*ldv+i] = vi[i+1+j]
				}
			} else {
				for c := 0; c < n; c++ {
					v[i*ldv+c] = 0
				}
				v[i*ldv+i] = 1
				for j := 0; j < length-1; j++ {
					v[i*ldv+i+1+j] = cmplx.Conj(vi[i+1+j])
				}
			}
		} else {
			anchor := n - k + i
			vi[anchor] = 1
			for j := 0; j < length-1; j++ {
				vi[anchor-1-j] = xvec[1+j]
			}
			if store == lapack.ColumnWise {
				for r := 0; r < n; r++ {
					v[r*ldv+i] = 0
				}
				v[anchor*ldv+i] = 1
				for j := 0; j < length-1; j++ {
					v[(anchor-1-j)*ldv+i] = vi[anchor-1-j]
				}
			} else {
				for c := 0; c < n; c++ {
					v[i*ldv+c] = 0
				}
				v[i*ldv+anchor] = 1
				for j := 0; j < length-1; j++ {
					v[i*ldv+anchor-1-j] = cmplx.Conj(vi[anchor-1-j])
				}
			}
		}
		vFull[i] = vi
	}

	// Call Zlarft.
	ldt := k
	tmat := make([]complex128, k*ldt)
	impl.Zlarft(direct, store, n, k, v, ldv, tau, tmat, ldt)

	// Assemble H_direct: product of individual reflectors per direct order.
	// direct == Forward: H = H_0 * H_1 * ... * H_{k-1}
	// direct == Backward: H = H_{k-1} * ... * H_1 * H_0
	// Each H_i = I - tau_i * v_i * v_i^H.
	H := eyeComplex(n)
	applyReflectorRight := func(H []complex128, vi []complex128, ti complex128) {
		// H := H * (I - ti * vi * vi^H)
		// For each row r of H: H[r,:] -= ti * (H[r,:] . vi) * vi^H
		for r := 0; r < n; r++ {
			var dot complex128
			for c := 0; c < n; c++ {
				dot += H[r*n+c] * vi[c]
			}
			dot *= ti
			for c := 0; c < n; c++ {
				H[r*n+c] -= dot * cmplx.Conj(vi[c])
			}
		}
	}
	if direct == lapack.Forward {
		for i := 0; i < k; i++ {
			applyReflectorRight(H, vFull[i], tau[i])
		}
	} else {
		for i := k - 1; i >= 0; i-- {
			applyReflectorRight(H, vFull[i], tau[i])
		}
	}

	// Assemble H_T from (V, T):
	// columnwise: H_T = I - V * T * V^H (V is n×k)
	// rowwise:    H_T = I - V^H * T * V (V is k×n)
	HT := eyeComplex(n)
	if store == lapack.ColumnWise {
		// VT = V * T (n×k).
		VT := make([]complex128, n*k)
		for r := 0; r < n; r++ {
			for c := 0; c < k; c++ {
				var s complex128
				for p := 0; p < k; p++ {
					s += v[r*ldv+p] * tmat[p*ldt+c]
				}
				VT[r*k+c] = s
			}
		}
		// HT -= VT * V^H.
		for r := 0; r < n; r++ {
			for c := 0; c < n; c++ {
				var s complex128
				for p := 0; p < k; p++ {
					s += VT[r*k+p] * cmplx.Conj(v[c*ldv+p])
				}
				HT[r*n+c] -= s
			}
		}
	} else {
		// rowwise: V is k×n.
		// VhT = V^H * T -> n×k.
		VhT := make([]complex128, n*k)
		for r := 0; r < n; r++ {
			for c := 0; c < k; c++ {
				var s complex128
				for p := 0; p < k; p++ {
					s += cmplx.Conj(v[p*ldv+r]) * tmat[p*ldt+c]
				}
				VhT[r*k+c] = s
			}
		}
		// HT -= VhT * V.
		for r := 0; r < n; r++ {
			for c := 0; c < n; c++ {
				var s complex128
				for p := 0; p < k; p++ {
					s += VhT[r*k+p] * v[p*ldv+c]
				}
				HT[r*n+c] -= s
			}
		}
	}

	// Compare H vs HT.
	maxDiff := 0.0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			d := cmplx.Abs(H[i*n+j] - HT[i*n+j])
			if d > maxDiff {
				maxDiff = d
			}
		}
	}
	bound := tol * float64(n) * float64(k)
	if maxDiff > bound {
		t.Errorf("%s: T does not reconstruct H; maxDiff=%v, tol=%v", name, maxDiff, bound)
	}
}

func eyeComplex(n int) []complex128 {
	out := make([]complex128, n*n)
	for i := 0; i < n; i++ {
		out[i*n+i] = 1
	}
	return out
}
