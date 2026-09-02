// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/cmplx"
	"math/rand/v2"
	"testing"
)

type Zlarfger interface {
	Zlarfg(n int, alpha complex128, x []complex128, incX int) (beta float64, tau complex128)
}

// ZlarfgTest verifies that Zlarfg produces a reflector H such that
//
//	H * (alpha; x) = (beta; 0)
//	H^H * H = I
//
// where H = I - tau * u * u^H with u = (1; v) and v overwrites x on exit.
func ZlarfgTest(t *testing.T, impl Zlarfger) {
	const tol = 1e-13
	rnd := rand.New(rand.NewPCG(1, 1))
	type tc struct {
		n     int
		alpha complex128
		x     []complex128
	}
	cases := []tc{
		{n: 1, alpha: complex(2, 3)},
		{n: 1, alpha: 0},
		{n: 3, alpha: complex(4, 0)},
		{n: 3, alpha: complex(-2, 1)},
		{n: 3, alpha: 0, x: []complex128{0, 0}},
		{n: 3, alpha: 0, x: []complex128{complex(1, 2), complex(-0.5, 0.25)}},
		{n: 4, alpha: complex(1, 0), x: []complex128{
			complex(4, -1), complex(5, 2), complex(6, 0),
		}},
		{n: 5},
		{n: 8},
	}
	for i, test := range cases {
		n := test.n
		incX := 1
		var x []complex128
		if test.x == nil {
			x = make([]complex128, max(0, n-1))
			for j := range x {
				x[j] = complex(rnd.NormFloat64(), rnd.NormFloat64())
			}
		} else {
			if len(test.x) != n-1 {
				panic("bad test data")
			}
			x = make([]complex128, n-1)
			copy(x, test.x)
		}
		alpha := test.alpha
		if alpha == 0 && test.x == nil && n != 0 {
			alpha = complex(rnd.NormFloat64(), rnd.NormFloat64())
		}
		xCopy := make([]complex128, len(x))
		copy(xCopy, x)
		alphaOrig := alpha

		beta, tau := impl.Zlarfg(n, alpha, x, incX)

		name := fmt.Sprintf("case %d (n=%d)", i, n)

		// Construct v = (1; x_on_exit) of length n.
		if n == 0 {
			continue
		}
		v := make([]complex128, n)
		v[0] = 1
		for j := 0; j < n-1; j++ {
			v[1+j] = x[j]
		}

		// Zlarfg returns tau such that H**H * (alpha; x) = (beta; 0), where
		// H = I - tau * v * v**H. Applying H**H = I - conj(tau)*v*v**H zeros x.
		input := make([]complex128, n)
		input[0] = alphaOrig
		for j := 0; j < n-1; j++ {
			input[1+j] = xCopy[j]
		}
		var vhu complex128
		for j := 0; j < n; j++ {
			vhu += cmplx.Conj(v[j]) * input[j]
		}
		got := make([]complex128, n)
		for j := 0; j < n; j++ {
			got[j] = input[j] - cmplx.Conj(tau)*v[j]*vhu
		}
		if d := cmplx.Abs(got[0] - complex(beta, 0)); d > tol {
			t.Errorf("%s: H^H*(alpha;x)[0] != beta; got=%v, beta=%v", name, got[0], beta)
		}
		for j := 1; j < n; j++ {
			if cmplx.Abs(got[j]) > tol {
				t.Errorf("%s: H^H*(alpha;x)[%d] not zero: %v", name, j, got[j])
			}
		}

		// Verify H^H H = I. Since H = I - tau v v^H,
		//   H^H H = I - (tau + conj(tau))*v*v^H + |tau|^2 * (v^H v) * v*v^H.
		// For H^H H = I we need tau + conj(tau) == |tau|^2 * (v^H v).
		var vhv complex128
		for j := 0; j < n; j++ {
			vhv += cmplx.Conj(v[j]) * v[j]
		}
		lhs := tau + cmplx.Conj(tau)
		rhs := tau * cmplx.Conj(tau) * vhv
		if d := cmplx.Abs(lhs - rhs); d > tol*(1+cmplx.Abs(tau)) {
			t.Errorf("%s: H^H H != I; residual %v (tau=%v, vhv=%v)", name, d, tau, vhv)
		}
	}
}
