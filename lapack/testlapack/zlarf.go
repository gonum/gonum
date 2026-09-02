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
)

type Zlarfer interface {
	Zlarf(side blas.Side, m, n int, v []complex128, incv int, tau complex128, c []complex128, ldc int, work []complex128)
}

// ZlarfTest verifies that Zlarf applies the Householder reflector
// H = I - tau * v * v^H to C, matching a naive reference implementation.
func ZlarfTest(t *testing.T, impl Zlarfer) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, side := range []blas.Side{blas.Left, blas.Right} {
		name := "Left"
		if side == blas.Right {
			name = "Right"
		}
		t.Run(name, func(t *testing.T) {
			runZlarfTest(t, impl, rnd, side)
		})
	}
}

func runZlarfTest(t *testing.T, impl Zlarfer, rnd *rand.Rand, side blas.Side) {
	const tol = 1e-13
	for _, m := range []int{0, 1, 2, 3, 5, 10} {
		for _, n := range []int{0, 1, 2, 3, 5, 10} {
			for _, incv := range []int{1, 3} {
				for _, ldcExtra := range []int{0, 3} {
					for _, tauVariant := range []int{0, 1, 2} {
						zlarfTestCase(t, impl, rnd, side, m, n, incv, max(1, n)+ldcExtra, tauVariant, tol)
					}
				}
			}
		}
	}
}

func zlarfTestCase(t *testing.T, impl Zlarfer, rnd *rand.Rand, side blas.Side, m, n, incv, ldc, tauVariant int, tol float64) {
	var tau complex128
	switch tauVariant {
	case 0:
		tau = 0
	case 1:
		tau = complex(rnd.NormFloat64(), rnd.NormFloat64())
	case 2:
		tau = complex(1, 0)
	}

	if m == 0 || n == 0 {
		// Still call to exercise quick-return paths.
		c := make([]complex128, max(1, m*ldc))
		v := make([]complex128, 1)
		work := make([]complex128, max(1, m+n))
		impl.Zlarf(side, m, n, v, max(1, incv), tau, c, ldc, work)
		return
	}

	c := make([]complex128, m*ldc)
	for i := range c {
		c[i] = complex(rnd.NormFloat64(), rnd.NormFloat64())
	}
	cCopy := make([]complex128, len(c))
	copy(cCopy, c)

	vlen := n
	if side == blas.Left {
		vlen = m
	}
	v := make([]complex128, 1+(vlen-1)*incv)
	for i := 0; i < vlen; i++ {
		v[i*incv] = complex(rnd.NormFloat64(), rnd.NormFloat64())
	}

	var work []complex128
	if side == blas.Left {
		work = make([]complex128, n)
	} else {
		work = make([]complex128, m)
	}

	impl.Zlarf(side, m, n, v, incv, tau, c, ldc, work)

	// Compute reference result: H = I - tau * v * v^H applied to cCopy.
	want := make([]complex128, len(cCopy))
	copy(want, cCopy)
	if tau != 0 {
		if side == blas.Left {
			// H*C: for each col j, w_j = tau * v^H * C_j = tau * sum_i conj(v_i) * C[i,j].
			// Then C[:, j] -= v * w_j.
			for j := 0; j < n; j++ {
				var w complex128
				for i := 0; i < m; i++ {
					w += cmplx.Conj(v[i*incv]) * cCopy[i*ldc+j]
				}
				w *= tau
				for i := 0; i < m; i++ {
					want[i*ldc+j] -= v[i*incv] * w
				}
			}
		} else {
			// C*H: for each row i, w_i = tau * C_i * v = tau * sum_j C[i,j] * v_j.
			// Then C[i, :] -= w_i * v^H.
			for i := 0; i < m; i++ {
				var w complex128
				for j := 0; j < n; j++ {
					w += cCopy[i*ldc+j] * v[j*incv]
				}
				w *= tau
				for j := 0; j < n; j++ {
					want[i*ldc+j] -= w * cmplx.Conj(v[j*incv])
				}
			}
		}
	}

	name := fmt.Sprintf("m=%d,n=%d,incv=%d,ldc=%d,tau=%v", m, n, incv, ldc, tau)
	maxDiff := 0.0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			d := cmplx.Abs(c[i*ldc+j] - want[i*ldc+j])
			if d > maxDiff {
				maxDiff = d
			}
		}
	}
	if maxDiff > tol*float64(max(m, n)+1) {
		t.Errorf("%s: unexpected result; maxDiff=%v, tol=%v", name, maxDiff, tol*float64(max(m, n)+1))
	}
}
