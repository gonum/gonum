// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package window

import (
	"testing"

	"gonum.org/v1/gonum/floats"
)

// testing table (except gaussian window)
var windowTests = []struct {
	name string
	fn   func([]float64) []float64
	want []float64
}{
	{
		name: "Rectangular", fn: Rectangular,
		want: []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	},
	{
		name: "Sine", fn: Sine,
		want: []float64{0.000000, 0.164595, 0.324699, 0.475947, 0.614213, 0.735724, 0.837166, 0.915773, 0.969400, 0.996584,
			0.996584, 0.969400, 0.915773, 0.837166, 0.735724, 0.614213, 0.475947, 0.324699, 0.164595, 0.000000},
	},
	{
		name: "Lanczos", fn: Lanczos,
		want: []float64{0.000000, 0.115514, 0.247646, 0.389468, 0.532984, 0.669692, 0.791213, 0.889915, 0.959492, 0.995450,
			0.995450, 0.959492, 0.889915, 0.791213, 0.669692, 0.532984, 0.389468, 0.247646, 0.115514, 0.000000},
	},
	{
		name: "Triangular", fn: Triangular,
		want: []float64{0.000000, 0.105263, 0.210526, 0.315789, 0.421053, 0.526316, 0.631579, 0.736842, 0.842105, 0.947368,
			0.947368, 0.842105, 0.736842, 0.631579, 0.526316, 0.421053, 0.315789, 0.210526, 0.105263, 0.000000},
	},
	{
		name: "Hann", fn: Hann,
		want: []float64{0.000000, 0.027091, 0.105430, 0.226526, 0.377257, 0.541290, 0.700848, 0.838641, 0.939737, 0.993181,
			0.993181, 0.939737, 0.838641, 0.700848, 0.541290, 0.377257, 0.226526, 0.105430, 0.027091, 0.000000},
	},
	{
		name: "BartlettHann", fn: BartlettHann,
		want: []float64{0.000000, 0.045853, 0.130653, 0.247949, 0.387768, 0.537696, 0.684223, 0.814209, 0.916305, 0.982186,
			0.982186, 0.916305, 0.814209, 0.684223, 0.537696, 0.387768, 0.247949, 0.130653, 0.045853, 0.000000},
	},
	{
		name: "Hamming", fn: Hamming,
		want: []float64{0.086957, 0.111692, 0.183218, 0.293785, 0.431408, 0.581178, 0.726861, 0.852672, 0.944977, 0.993774,
			0.993774, 0.944977, 0.852672, 0.726861, 0.581178, 0.431409, 0.293785, 0.183218, 0.111692, 0.086957},
	},
	{
		name: "Blackman", fn: Blackman,
		want: []float64{0.000000, 0.010223, 0.045069, 0.114390, 0.226899, 0.382381, 0.566665, 0.752034, 0.903493, 0.988846,
			0.988846, 0.903493, 0.752034, 0.566665, 0.382381, 0.226899, 0.114390, 0.045069, 0.010223, 0.000000},
	},
	{
		name: "BlackmanHarris", fn: BlackmanHarris,
		want: []float64{0.000060, 0.002018, 0.012795, 0.046450, 0.122540, 0.256852, 0.448160, 0.668576, 0.866426, 0.984278,
			0.984278, 0.866426, 0.668576, 0.448160, 0.256852, 0.122540, 0.046450, 0.012795, 0.002018, 0.000060},
	},
	{
		name: "Nuttall", fn: Nuttall,
		want: []float64{0.000000, 0.001706, 0.011614, 0.043682, 0.117808, 0.250658, 0.441946, 0.664015, 0.864348, 0.984019,
			0.984019, 0.864348, 0.664015, 0.441946, 0.250658, 0.117808, 0.043682, 0.011614, 0.001706, 0.000000},
	},
	{
		name: "BlackmanNuttall", fn: BlackmanNuttall,
		want: []float64{0.000363, 0.002885, 0.015360, 0.051652, 0.130567, 0.266629, 0.457501, 0.675215, 0.869392, 0.984644,
			0.984644, 0.869392, 0.675215, 0.457501, 0.266629, 0.130567, 0.051652, 0.015360, 0.002885, 0.000363},
	},
	{
		name: "FlatTop", fn: FlatTop,
		want: []float64{-0.000421, -0.003687, -0.017675, -0.045939, -0.070137, -0.037444, 0.115529, 0.402051, 0.737755, 0.967756,
			0.967756, 0.737755, 0.402051, 0.115529, -0.037444, -0.070137, -0.045939, -0.017675, -0.003687, -0.000421},
	},
}

//testing table (gaussian window)
var gausWindowTests = []struct {
	name  string
	sigma float64
	want  []float64
}{
	{
		name: "Gaussian (sigma=0.3)", sigma: 0.3,
		want: []float64{0.003866, 0.011708, 0.031348, 0.074214, 0.155344, 0.287499, 0.470444, 0.680632, 0.870660, 0.984728,
			0.984728, 0.870660, 0.680632, 0.470444, 0.287499, 0.155344, 0.074214, 0.031348, 0.011708, 0.003866},
	},
	{
		name: "Gaussian (sigma=0.5)", sigma: 0.5,
		want: []float64{0.135335, 0.201673, 0.287499, 0.392081, 0.511524, 0.638423, 0.762260, 0.870660, 0.951361, 0.994475,
			0.994475, 0.951361, 0.870660, 0.762260, 0.638423, 0.511524, 0.392081, 0.287499, 0.201673, 0.135335},
	},
	{
		name: "Gaussian (sigma=1.2)", sigma: 1.2,
		want: []float64{0.706648, 0.757319, 0.805403, 0.849974, 0.890135, 0.925049, 0.953963, 0.976241, 0.991381, 0.999039,
			0.999039, 0.991381, 0.976241, 0.953963, 0.925049, 0.890135, 0.849974, 0.805403, 0.757319, 0.706648},
	},
}

func TestWindows(t *testing.T) {
	const tol = 1e-6

	//Input data
	src := make([]float64, 20)
	for i := range src {
		src[i] = 1
	}

	for _, test := range windowTests {
		t.Run(test.name, func(t *testing.T) {
			// Copy the input since we are mutating it.
			srcCpy := make([]float64, len(src))
			copy(srcCpy, src)
			dst := test.fn(srcCpy)

			if !floats.EqualApprox(dst, test.want, tol) {
				t.Errorf("unexpected result for window function %q:\ngot:%v\nwant:%v", test.name, dst, test.want)
			}
		})
	}
}

func TestGausWindows(t *testing.T) {
	const tol = 1e-6

	//Input data
	src := make([]float64, 20)
	for i := range src {
		src[i] = 1
	}

	for _, test := range gausWindowTests {
		t.Run(test.name, func(t *testing.T) {
			// Copy the input since we are mutating it.
			srcCpy := make([]float64, len(src))
			copy(srcCpy, src)
			dst := Gaussian(srcCpy, test.sigma)

			if !floats.EqualApprox(dst, test.want, tol) {
				t.Errorf("unexpected result for window function %q:\ngot:%v\nwant:%v", test.name, dst, test.want)
			}
		})
	}
}
