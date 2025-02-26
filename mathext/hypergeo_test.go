// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import (
	"fmt"
	"math"
	"testing"
)

func TestHyp2f1(t *testing.T) {
	t.Parallel()

	// Constants taken from https://github.com/RobinHankin/hypergeo/blob/master/tests/testthat/test_aaa.R
	var tests = []struct {
		x float64
		y float64
	}{
		{x: 0.28, y: 1.3531156987873853569937},
		{x: -0.79, y: 0.5773356740314405932679},
		{x: 0.56, y: 2.1085704049533617876477},
		{x: -2.13, y: 0.3352446571148822718200},
		{x: -0.43, y: 0.7150355048137748692483},
		{x: -1.23, y: 0.4670987707934830535095},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			y := Hypergeo(1.21, 1.443, 1.88, test.x)
			if d := y - test.y; math.Abs(d) > 1e-12 {
				t.Errorf("%f %f %f %f", test.x, y, test.y, d)
			}
		})
	}
}

func TestHyp2f1_15_1_15(t *testing.T) {
	t.Parallel()

	// eqn15_1_15a is the left hand side of equation 15.1.15 of
	// M. Abramowitz and I. A. Stegun 1965. Handbook of Mathematical Functions, New York: Dover.
	eqn15_1_15a := func(a, z float64) float64 {
		return Hypergeo(a, 1-a, 3./2, math.Pow(math.Sin(z), 2))
	}
	// eqn15_1_15b is the right hand side of equation 15.1.15 of Abramowitz.
	eqn15_1_15b := func(a, z float64) float64 {
		return math.Sin((2*a-1)*z) / ((2*a - 1) * math.Sin(z))
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
			t.Parallel()

			// Ignore z=-2.13, since both R's hypergeo and Maple don't handle this case either.
			if test.x == -2.13 {
				return
			}

			a := eqn15_1_15a(0.2, test.x)
			b := eqn15_1_15b(0.2, test.x)
			if d := a - b; math.Abs(d) > 1e-6 {
				t.Errorf("%f %f %f %f", test.x, a, b, d)
			}
		})
	}
}

func TestHyp2f1_15_2_10(t *testing.T) {
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
			t.Parallel()
			v := eqn15_2_10(0.1, 0.44, 0.611, test.x)
			if math.Abs(v) > 1e-6 {
				t.Errorf("%f %f", test.x, v)
			}
		})
	}
}

func TestHyp2f1_15_1(t *testing.T) {
	t.Parallel()

	type equation struct {
		name string
		f    func(float64) float64
	}

	equations := []equation{
		{
			name: "15_1_3",
			f: func(z float64) float64 {
				lhs := Hypergeo(1, 1, 2, z)
				rhs := -math.Log(1-z) / z
				return math.Abs(rhs - lhs)
			},
		},
		{
			name: "15_1_5",
			f: func(z float64) float64 {
				lhs := Hypergeo(1./2, 1, 3./2, -z*z)
				rhs := math.Atan(z) / z
				return math.Abs(rhs - lhs)
			},
		},
		{
			name: "15_1_7a",
			f: func(z float64) float64 {
				lhs := Hypergeo(1./2, 1./2, 3./2, -z*z)
				rhs := math.Sqrt(1+z*z) * Hypergeo(1, 1, 3./2, -z*z)
				return math.Abs(rhs - lhs)
			},
		},
		{
			name: "15_1_7b",
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
			t.Parallel()
			for _, eqn := range equations {
				if d := eqn.f(test.x); d > 1e-10 {
					t.Errorf("%s %f %f", eqn.name, test.x, d)
				}
			}
		})
	}
}

func TestHyp2f1_15_1_zz(t *testing.T) {
	t.Parallel()

	type equation struct {
		name string
		f    func(float64) float64
	}

	equations := []equation{
		{
			name: "15_1_4",
			f: func(z float64) float64 {
				lhs := Hypergeo(1./2, 1, 3./2, z*z)
				rhs := 0.5 * math.Log((1+z)/(1-z)) / z
				return math.Abs(rhs - lhs)
			},
		},
		{
			name: "15_1_6a",
			f: func(z float64) float64 {
				lhs := Hypergeo(1./2, 1./2, 3./2, z*z)
				rhs := math.Sqrt(1-z*z) * Hypergeo(1, 1, 3./2, z*z)
				return math.Abs(rhs - lhs)
			},
		},
		{
			name: "15_1_6b",
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
		{x: -2.13},
		{x: -0.43},
		{x: -1.23},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			if test.x*test.x > 1 {
				return
			}

			for _, eqn := range equations {
				if d := eqn.f(test.x); d > 1e-10 {
					t.Errorf("%s %f %f", eqn.name, test.x, d)
				}
			}
		})
	}
}

func TestHyp2f1_Igor_Kojanov(t *testing.T) {
	t.Parallel()

	var y float64
	y = Hypergeo(1, 2, 3, 0)
	if d := y - 1; d != 0 {
		t.Errorf("%f %f", y, d)
	}

	y = Hypergeo(1, 1.64, 2.64, -0.1111)
	if d := y - 0.9361003540660249866434; math.Abs(d) > 1e-15 {
		t.Errorf("%f %f", y, d)
	}
}

func TestHyp2f1_John_Ormerod(t *testing.T) {
	t.Parallel()

	y := Hypergeo(5.25, 1, 6.5, 0.501)
	if d := y - 1.70239432012007391092082702795; math.Abs(d) > 1e-10 {
		t.Errorf("%f %f", y, d)
	}
}

func TestHyp2f1Scipy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		a    float64
		b    float64
		c    float64
		z    float64
		want float64
	}{
		// Constants taken from https://github.com/scipy/scipy/blob/main/scipy/special/tests/test_hyp2f1.py
		{a: 1.3, b: -0.2, c: 0.3, z: -2.1, want: 1.8202169687521206},

		// Constants taken from https://github.com/scipy/scipy/blob/main/scipy/special/tests/test_mpmath.py
		{a: 1. / 3, b: 2. / 3, c: 5. / 6, z: 27. / 32, want: 1.6},
		{a: 1. / 4, b: 1. / 2, c: 3. / 4, z: 80. / 81, want: 1.8},
		{a: 0.7235, b: -1, c: -5, z: 0.3, want: 1.04341},
		{a: 0.25, b: 1. / 3, c: 2, z: 0.999, want: 1.0682644949603062},
		{a: 0.25, b: 1. / 3, c: 2, z: -1, want: 0.9665658449252437},
		{a: 2, b: 3, c: 5, z: 0.99, want: 27.699347904322664},
		{a: 3. / 2, b: -0.5, c: 3, z: 0.99, want: 0.6840303684391167},
		{a: 2, b: 2.5, c: -3.25, z: 0.999, want: 2.183739328012162e+26},
		{a: -8, b: 18.016500331508873, c: 10.805295997850628, z: 0.90875647507000001, want: -3.566216341442061e-09},
		{a: -10, b: 900, c: -10.5, z: 0.99, want: 2.5101757354622962e+22},
		{a: -10, b: 900, c: 10.5, z: 0.99, want: 5.5748237303615776e+17},
		{a: -1, b: 2, c: 1, z: -1, want: 3},
		{a: 0.5, b: 1 - 270.5, c: 1.5, z: 0.999 * 0.999, want: 0.053963052503373715},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			y := Hypergeo(test.a, test.b, test.c, test.z)
			if d := (y - test.want) / test.want; math.Abs(d) > 1e-10 {
				t.Errorf("%f %f %f", y, test.want, d)
			}
		})
	}
}
