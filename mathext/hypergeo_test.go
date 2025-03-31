// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import (
	"fmt"
	"math"
	"testing"

	"gonum.org/v1/gonum/floats/scalar"
)

func TestHypergeo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		a    float64
		b    float64
		c    float64
		z    float64
		want float64
		tol  float64
	}{
		// Constants based on https://github.com/scipy/scipy/blob/main/scipy/special/tests/test_hyp2f1.py
		{a: 1.3, b: -0.2, c: 0.3, z: -2.1, want: 1.8202169687521206, tol: 5e-15},

		// Constants based on https://github.com/scipy/scipy/blob/main/scipy/special/tests/test_mpmath.py
		{a: 1. / 3, b: 2. / 3, c: 5. / 6, z: 27. / 32, want: 1.6, tol: 1e-10},
		{a: 1. / 4, b: 1. / 2, c: 3. / 4, z: 80. / 81, want: 1.8, tol: 1e-10},
		{a: 0.7235, b: -1, c: -5, z: 0.3, want: 1.04341, tol: 1e-10},
		{a: 0.25, b: 1. / 3, c: 2, z: 0.999, want: 1.0682644949603062, tol: 1e-10},
		{a: 0.25, b: 1. / 3, c: 2, z: -1, want: 0.9665658449252437, tol: 1e-10},
		{a: 2, b: 3, c: 5, z: 0.99, want: 27.699347904322664, tol: 1e-10},
		{a: 3. / 2, b: -0.5, c: 3, z: 0.99, want: 0.6840303684391167, tol: 1e-10},
		{a: 2, b: 2.5, c: -3.25, z: 0.999, want: 2.183739328012162e+26, tol: 1e-10},
		{a: -8, b: 18.016500331508873, c: 10.805295997850628, z: 0.90875647507000001, want: -3.566216341442061e-09, tol: 1e-10},
		{a: -10, b: 900, c: -10.5, z: 0.99, want: 2.5101757354622962e+22, tol: 1e-10},
		{a: -10, b: 900, c: 10.5, z: 0.99, want: 5.5748237303615776e+17, tol: 1e-10},
		{a: -1, b: 2, c: 1, z: -1, want: 3, tol: 1e-10},
		{a: 0.5, b: 1 - 270.5, c: 1.5, z: 0.999 * 0.999, want: 0.053963052503373715, tol: 1e-10},

		// Constants based on https://github.com/scipy/scipy/issues/1561
		{a: 10, b: 5, c: -300.5, z: 0.5, want: -3.85202708152391e32, tol: 5e-15},

		// Constants based on Table 26,
		// John Pearson, Computation of Hypergeometric Functions,
		// Master thesis for Worcester College, Oxford University.
		// https://api.semanticscholar.org/CorpusID:124333574
		{a: 0.1, b: 0.2, c: 0.3, z: 0.5, want: 1.046432811217352, tol: 1e-15},
		{a: -0.1, b: 0.2, c: 0.3, z: 0.5, want: 0.956434210968214, tol: 1e-15},
		{a: 1e-8, b: 1e-8, c: 1e-8, z: 1e-6, want: 1, tol: 1e-14},
		{a: 2 + 1e-9, b: 3, c: 5, z: -0.75, want: 0.492238858852651, tol: 1e-15},
		{a: -2, b: -3, c: -5 + 1e-9, z: 0.5, want: 0.474999999913750, tol: 1e-15},
		{a: -1, b: -1.5, c: -2 - 1e-15, z: 0.5, want: 0.625, tol: 1e-15},
		{a: 500, b: -500, c: 500, z: 0.75, want: 9.332636185032189e-302, tol: 1e-15},
		{a: 500, b: 500, c: 500, z: -0.6, want: 8.709809816217217e-103, tol: 1e-15},
		{a: -1000, b: -2000, c: -4000.1, z: -0.5, want: 5.233580403196932e94, tol: 5e-7},
		{a: -100, b: -200, c: -300 + 1e-9, z: 0.5 * math.Sqrt(2), want: 2.653635302903707e-31, tol: 1e-15},
		{a: 300, b: 10, c: 5, z: 0.5, want: 3.912238919961547e98, tol: 1e-15},
		{a: 5, b: -300, c: 10, z: 0.5, want: 1.661006238211309e-7, tol: 1e-15},
		{a: 2.25, b: 3.75, c: -0.5, z: -1, want: -0.631220676949703, tol: 1e-15},

		// Additional cases for hardening against large |a|, |b|, or |c|.
		// These wanted values are based on the agreed values of both Mathematica 14.2.0 and mpmath 1.3.0.
		// mpmath is chosen because it is often treated as
		// ground truth in scipy's issues:
		// * https://github.com/scipy/scipy/issues/1561
		// * https://github.com/scipy/scipy/issues/5349
		{a: -290, b: 5, c: -300.5, z: 0.5, want: 26.9076853843542, tol: 5e-15},
		{a: 10, b: 5, c: -300.5, z: -0.5, want: 1.08800279612753, tol: 1e-14},
		{a: 10, b: 5, c: 300.5, z: 0.5, want: 1.08796774775660, tol: 5e-15},
		{a: -5.28, b: -3, c: -12.28, z: 0.95, want: 0.17167484351795846, tol: 1e-15},
		{a: -5, b: 0.5, c: 0.3, z: -2.1, want: 607.226576917773834, tol: 1e-15},

		// Additional cases for hardening against large |a|, |b|, or |c|.
		// These wanted values are based on the agreed values
		// of both Mathematica 14.2.0 and Miller's algorithm described in https://github.com/scipy/scipy/issues/1561#issuecomment-130488352
		// We trust Mathematica and Miller's algorithm because
		// * Mathematica is the gold standard held in discussions in R's hypergeo, scipy, and mpmath:
		//   * https://github.com/RobinHankin/hypergeo/issues/7
		//   * https://github.com/scipy/scipy/issues/5349
		//   * https://github.com/mpmath/mpmath/issues/296
		// * Miller's algorithm is discussed in the often cited paper John Pearson, Computation of Hypergeometric Functions
		// Note that these are particular hard cases,
		// as all four softwares R's hypergeo, scipy, mpmath, and Mathematica give wildly different results.
		{a: 100, b: 5, c: -300.5, z: 0.5, want: -1.07414870340863581e139, tol: 1e-15},
		{a: -100, b: -200, c: 300 + 1e-9, z: -0.5 * math.Sqrt(2), want: -1.073671875630690e-30, tol: 1e-15},
		{a: -100, b: -200, c: -300 + 0.1, z: 0.5 * math.Sqrt(2), want: 2.568211952590272e-31, tol: 1e-15},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			y := Hypergeo(test.a, test.b, test.c, test.z)
			eps := math.Nextafter(1.0, 2.0) - 1.0
			if !scalar.EqualWithinAbsOrRel(y, test.want, eps, test.tol) {
				t.Errorf("expected result from Hypergeo(%f, %f, %f, %f): got %f want %f", test.a, test.b, test.c, test.z, y, test.want)
			}
		})
	}
}

func Test_15_1_15(t *testing.T) {
	t.Parallel()

	// eqn15_1_15_lhs is the left hand side of equation 15.1.15 of
	// M. Abramowitz and I. A. Stegun 1965. Handbook of Mathematical Functions, New York: Dover.
	eqn15_1_15_lhs := func(a, z float64) float64 {
		return Hypergeo(a, 1-a, 3./2, math.Pow(math.Sin(z), 2))
	}
	// eqn15_1_15_rhs is the right hand side of equation 15.1.15 of Abramowitz.
	eqn15_1_15_rhs := func(a, z float64) float64 {
		return math.Sin((2*a-1)*z) / ((2*a - 1) * math.Sin(z))
	}

	var tests = []struct {
		x float64
	}{
		{x: 0.28},
		{x: -0.79},
		{x: 0.56},
		{x: -0.43},
		{x: -1.23},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			lhs := eqn15_1_15_lhs(0.2, test.x)
			rhs := eqn15_1_15_rhs(0.2, test.x)
			if !scalar.EqualWithinAbs(lhs, rhs, 1e-6) {
				t.Errorf("unexpected result from eqn15.1.15(0.2, %f): lhs %f rhs %f", test.x, lhs, rhs)
			}
		})
	}
}

func Test_15_2_10(t *testing.T) {
	t.Parallel()

	eqn15_2_10 := func(a, b, c, z float64) float64 {
		return (c-a)*Hypergeo(a-1, b, c, z) + (2*a-c-a*z+b*z)*Hypergeo(a, b, c, z) + a*(z-1)*Hypergeo(a+1, b, c, z)
	}

	var tests = []struct {
		x float64
	}{
		{x: 0.28},
		{x: -0.79},
		{x: 0.56},
		{x: -2.13},
		{x: -0.43},
		{x: -1.23},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			v := eqn15_2_10(0.1, 0.44, 0.611, test.x)
			if !scalar.EqualWithinAbs(v, 0, 1e-6) {
				t.Errorf("unexpected result from eqn15.2.10(0.1, 0.44, 0.611, %f): got %f want 0", test.x, v)
			}
		})
	}
}

func Test_15_1(t *testing.T) {
	t.Parallel()

	type equation struct {
		name string
		f    func(float64) float64
	}

	equations := []equation{
		{
			name: "eqn15.1.3",
			f: func(z float64) float64 {
				lhs := Hypergeo(1, 1, 2, z)
				rhs := -math.Log(1-z) / z
				return math.Abs(rhs - lhs)
			},
		},
		{
			name: "eqn15.1.5",
			f: func(z float64) float64 {
				lhs := Hypergeo(1./2, 1, 3./2, -z*z)
				rhs := math.Atan(z) / z
				return math.Abs(rhs - lhs)
			},
		},
		{
			name: "eqn15.1.7a",
			f: func(z float64) float64 {
				lhs := Hypergeo(1./2, 1./2, 3./2, -z*z)
				rhs := math.Sqrt(1+z*z) * Hypergeo(1, 1, 3./2, -z*z)
				return math.Abs(rhs - lhs)
			},
		},
		{
			name: "eqn15.1.7b",
			f: func(z float64) float64 {
				lhs := Hypergeo(1./2, 1./2, 3./2, -z*z)
				rhs := math.Log(z+math.Sqrt(1+z*z)) / z
				return math.Abs(rhs - lhs)
			},
		},
	}

	var tests = []struct {
		x float64
	}{
		{x: 0.28},
		{x: -0.79},
		{x: 0.56},
		{x: -2.13},
		{x: -0.43},
		{x: -1.23},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			for _, eqn := range equations {
				y := eqn.f(test.x)
				if !scalar.EqualWithinAbs(y, 0, 1e-10) {
					t.Errorf("unexpected result from %s(%f): got %f want 0", eqn.name, test.x, y)
				}
			}
		})
	}
}

func Test_15_1_zz(t *testing.T) {
	t.Parallel()

	type equation struct {
		name string
		f    func(float64) float64
	}

	equations := []equation{
		{
			name: "eqn15.1.4",
			f: func(z float64) float64 {
				lhs := Hypergeo(1./2, 1, 3./2, z*z)
				rhs := 0.5 * math.Log((1+z)/(1-z)) / z
				return math.Abs(rhs - lhs)
			},
		},
		{
			name: "eqn15.1.6a",
			f: func(z float64) float64 {
				lhs := Hypergeo(1./2, 1./2, 3./2, z*z)
				rhs := math.Sqrt(1-z*z) * Hypergeo(1, 1, 3./2, z*z)
				return math.Abs(rhs - lhs)
			},
		},
		{
			name: "eqn15.1.6b",
			f: func(z float64) float64 {
				lhs := Hypergeo(1./2, 1./2, 3./2, z*z)
				rhs := math.Asin(z) / z
				return math.Abs(rhs - lhs)
			},
		},
	}

	var tests = []struct {
		x float64
	}{
		{x: 0.28},
		{x: -0.79},
		{x: 0.56},
		{x: -0.43},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			for _, eqn := range equations {
				y := eqn.f(test.x)
				if !scalar.EqualWithinAbs(y, 0, 1e-10) {
					t.Errorf("unexpected result from %s(%f): got %f want 0", eqn.name, test.x, y)
				}
			}
		})
	}
}

func Test_Igor_Kojanov(t *testing.T) {
	t.Parallel()

	// Constants taken from https://github.com/RobinHankin/hypergeo/blob/master/tests/testthat/test_aaa.R
	y := Hypergeo(1, 2, 3, 0)
	if y != 1 {
		t.Errorf("unexpected result from Hypergeo(1, 2, 3, 0): got %f want 1", y)
	}

	y = Hypergeo(1, 1.64, 2.64, -0.1111)
	want := 0.9361003540660249866434
	if !scalar.EqualWithinAbs(y, want, 1e-15) {
		t.Errorf("unexpected result from Hypergeo(1, 1.64, 2.64, -0.1111): got %f want %f", y, want)
	}
}

func Test_John_Ormerod(t *testing.T) {
	t.Parallel()

	// Constants taken from https://github.com/RobinHankin/hypergeo/blob/master/tests/testthat/test_aaa.R
	y := Hypergeo(5.25, 1, 6.5, 0.501)
	want := 1.70239432012007391092082702795
	if !scalar.EqualWithinAbs(y, want, 1e-10) {
		t.Errorf("unexpected result from Hypergeo(5.25, 1, 6.5, 0.501): got %f want %f", y, want)
	}
}
