// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/cmplx"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/lapack"
)

type Zlarfber interface {
	Zlarfter
	Zlarfer
	Zlarfb(side blas.Side, trans blas.Transpose, direct lapack.Direct, store lapack.StoreV,
		m, n, k int, v []complex128, ldv int, t []complex128, ldt int,
		c []complex128, ldc int, work []complex128, ldwork int)
}

// ZlarfbTest checks that the block-reflector application by Zlarfb matches
// the result of applying k unblocked reflectors with Zlarf.
func ZlarfbTest(t *testing.T, impl Zlarfber) {
	const tol = 1e-12
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, store := range []lapack.StoreV{lapack.ColumnWise, lapack.RowWise} {
		for _, direct := range []lapack.Direct{lapack.Forward, lapack.Backward} {
			for _, side := range []blas.Side{blas.Left, blas.Right} {
				for _, trans := range []blas.Transpose{blas.NoTrans, blas.ConjTrans} {
					for _, sz := range []struct{ m, n, k int }{
						{6, 5, 3}, {5, 6, 3}, {7, 7, 4}, {4, 4, 4},
						{8, 4, 2}, {4, 8, 2}, {10, 10, 5},
					} {
						name := fmt.Sprintf("store=%c dir=%c side=%c trans=%c m=%d n=%d k=%d",
							storeChar(store), directChar(direct), sideChar(side), transChar(trans),
							sz.m, sz.n, sz.k)
						zlarfbCheck(t, impl, rnd, side, trans, direct, store, sz.m, sz.n, sz.k, name, tol)
					}
				}
			}
		}
	}
}

func zlarfbCheck(t *testing.T, impl Zlarfber, rnd *rand.Rand, side blas.Side, trans blas.Transpose,
	direct lapack.Direct, store lapack.StoreV, m, n, k int, name string, tol float64) {
	// nv is the size of the "long" dimension of V (equal to m for Left, n for Right).
	nv := m
	if side == blas.Right {
		nv = n
	}
	if k > nv {
		return
	}
	// Build V as the result of random reflectors.
	// For columnwise storage, V is nv×k. For rowwise, V is k×nv.
	var ldv int
	var v []complex128
	if store == lapack.ColumnWise {
		ldv = k
		v = make([]complex128, nv*ldv)
	} else {
		ldv = nv
		v = make([]complex128, k*ldv)
	}
	tau := make([]complex128, k)

	// Generate each reflector: pick a random vector of appropriate length,
	// use Zlarfg to get tau and overwrite x with v-below-diagonal, write the
	// unit diagonal and zeros in the right triangle.
	for i := 0; i < k; i++ {
		var length int
		if direct == lapack.Forward {
			length = nv - i
		} else {
			length = nv - (k - 1 - i)
		}
		if length < 1 {
			length = 1
		}
		xvec := make([]complex128, length)
		for j := range xvec {
			xvec[j] = complex(rnd.NormFloat64(), rnd.NormFloat64())
		}
		alpha := xvec[0]
		_, tauI := impl.Zlarfg(length, alpha, xvec[1:], 1)
		tau[i] = tauI

		// Store v_i in V according to LAPACK convention:
		//   - ColumnWise: V[:, i] stores u_i directly.
		//   - RowWise:   V[i, :] stores conj(u_i^T) so that H = I - V^H T V.
		if direct == lapack.Forward {
			if store == lapack.ColumnWise {
				for r := 0; r < i; r++ {
					v[r*ldv+i] = 0
				}
				v[i*ldv+i] = 1
				for j := 0; j < length-1; j++ {
					v[(i+1+j)*ldv+i] = xvec[1+j]
				}
			} else {
				for c := 0; c < i; c++ {
					v[i*ldv+c] = 0
				}
				v[i*ldv+i] = 1
				for j := 0; j < length-1; j++ {
					v[i*ldv+i+1+j] = cmplx.Conj(xvec[1+j])
				}
			}
		} else {
			if store == lapack.ColumnWise {
				anchor := nv - k + i
				for r := anchor + 1; r < nv; r++ {
					v[r*ldv+i] = 0
				}
				v[anchor*ldv+i] = 1
				for j := 0; j < length-1; j++ {
					v[(anchor-1-j)*ldv+i] = xvec[1+j]
				}
			} else {
				anchor := nv - k + i
				for c := anchor + 1; c < nv; c++ {
					v[i*ldv+c] = 0
				}
				v[i*ldv+anchor] = 1
				for j := 0; j < length-1; j++ {
					v[i*ldv+anchor-1-j] = cmplx.Conj(xvec[1+j])
				}
			}
		}
	}

	// Build T with Zlarft.
	ldt := k
	tmat := make([]complex128, k*ldt)
	impl.Zlarft(direct, store, nv, k, v, ldv, tau, tmat, ldt)

	// Generate random C.
	ldc := n + 2
	c := make([]complex128, m*ldc)
	for i := range c {
		c[i] = complex(rnd.NormFloat64(), rnd.NormFloat64())
	}
	cBlocked := make([]complex128, len(c))
	copy(cBlocked, c)
	cUnblocked := make([]complex128, len(c))
	copy(cUnblocked, c)

	// Apply blocked reflector.
	ldwork := k
	nw := n
	if side == blas.Right {
		nw = m
	}
	work := make([]complex128, nw*ldwork)
	impl.Zlarfb(side, trans, direct, store, m, n, k, v, ldv, tmat, ldt, cBlocked, ldc, work, ldwork)

	// Apply same reflectors one by one via Zlarf.
	// H = H_0 * H_1 * ... * H_{k-1} for direct == Forward
	// H = H_{k-1} * ... * H_0 for direct == Backward.
	// For NoTrans Left: C := H * C.
	//   - Forward: apply H_{k-1}, H_{k-2}, ..., H_0 (since H*C = H_0*(H_1*(...*(H_{k-1}*C))))
	//   - Backward: apply H_0, H_1, ..., H_{k-1}.
	// For ConjTrans Left: C := H^H * C.
	//   - H^H = H_{k-1}^H * ... * H_0^H for Forward (reverse order and conj).
	//   - i.e. apply H_0^H, H_1^H, ..., H_{k-1}^H for Forward.
	// For Right NoTrans: C := C * H.
	//   - Forward: apply H_0, H_1, ..., H_{k-1} (C*H = C*H_0*H_1*...).
	//   - Backward: reverse.
	// For Right ConjTrans: C := C * H^H, again reversing.
	//
	// The sequence to apply unblocked reflectors:
	order := make([]int, k)
	if side == blas.Left {
		if (direct == lapack.Forward && trans == blas.NoTrans) ||
			(direct == lapack.Backward && trans == blas.ConjTrans) {
			// Apply in reverse of index.
			for i := 0; i < k; i++ {
				order[i] = k - 1 - i
			}
		} else {
			for i := 0; i < k; i++ {
				order[i] = i
			}
		}
	} else {
		// Right.
		if (direct == lapack.Forward && trans == blas.NoTrans) ||
			(direct == lapack.Backward && trans == blas.ConjTrans) {
			for i := 0; i < k; i++ {
				order[i] = i
			}
		} else {
			for i := 0; i < k; i++ {
				order[i] = k - 1 - i
			}
		}
	}

	workU := make([]complex128, max(m, n))
	for _, idx := range order {
		// Extract full-length u_idx of size nv (always reconstructing the
		// un-conjugated reflector column form).
		vi := make([]complex128, nv)
		if direct == lapack.Forward {
			if store == lapack.ColumnWise {
				vi[idx] = 1
				for r := idx + 1; r < nv; r++ {
					vi[r] = v[r*ldv+idx]
				}
			} else {
				vi[idx] = 1
				for c := idx + 1; c < nv; c++ {
					vi[c] = cmplx.Conj(v[idx*ldv+c])
				}
			}
		} else {
			anchor := nv - k + idx
			if store == lapack.ColumnWise {
				vi[anchor] = 1
				for r := 0; r < anchor; r++ {
					vi[r] = v[r*ldv+idx]
				}
			} else {
				vi[anchor] = 1
				for c := 0; c < anchor; c++ {
					vi[c] = cmplx.Conj(v[idx*ldv+c])
				}
			}
		}
		effTau := tau[idx]
		if trans == blas.ConjTrans {
			effTau = cmplx.Conj(effTau)
		}
		impl.Zlarf(side, m, n, vi, 1, effTau, cUnblocked, ldc, workU)
	}

	// Compare.
	maxDiff := 0.0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			d := cmplx.Abs(cBlocked[i*ldc+j] - cUnblocked[i*ldc+j])
			if d > maxDiff {
				maxDiff = d
			}
		}
	}
	bound := tol * float64(max(1, max(m, n))) * float64(k)
	if maxDiff > bound {
		t.Errorf("%s: blocked vs unblocked mismatch, maxDiff=%v, tol=%v", name, maxDiff, bound)
	}
}

func storeChar(s lapack.StoreV) byte {
	if s == lapack.ColumnWise {
		return 'C'
	}
	return 'R'
}

func directChar(d lapack.Direct) byte {
	if d == lapack.Forward {
		return 'F'
	}
	return 'B'
}

func sideChar(s blas.Side) byte {
	if s == blas.Left {
		return 'L'
	}
	return 'R'
}

func transChar(t blas.Transpose) byte {
	switch t {
	case blas.NoTrans:
		return 'N'
	case blas.Trans:
		return 'T'
	case blas.ConjTrans:
		return 'C'
	}
	return '?'
}
