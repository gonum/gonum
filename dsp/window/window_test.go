// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package window

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/floats/scalar"
)

var windowTests = []struct {
	name    string
	fn      func([]float64) []float64
	fnCmplx func([]complex128) []complex128
	want    []float64
}{
	{
		name: "Rectangular", fn: Rectangular, fnCmplx: RectangularComplex,
		want: []float64{
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		},
	},
	{
		name: "Sine", fn: Sine, fnCmplx: SineComplex,
		want: []float64{
			0.078459, 0.233445, 0.382683, 0.522499, 0.649448, 0.760406, 0.852640, 0.923880, 0.972370, 0.996917,
			0.996917, 0.972370, 0.923880, 0.852640, 0.760406, 0.649448, 0.522499, 0.382683, 0.233445, 0.078459,
		},
	},
	{
		name: "Lanczos", fn: Lanczos, fnCmplx: LanczosComplex,
		want: []float64{
			0.052415, 0.170011, 0.300105, 0.436333, 0.57162, 0.698647, 0.810332, 0.900316, 0.963398, 0.995893,
			0.995893, 0.963398, 0.900316, 0.810332, 0.698647, 0.57162, 0.436333, 0.300105, 0.170011, 0.052415,
		},
	},
	// This case tests Lanczos for a NaN condition. The Lanczos NaN condition is k=(N-1)/2, that is when N is odd.
	{
		name: "LanczosOdd", fn: Lanczos, fnCmplx: LanczosComplex,
		want: []float64{
			0.049813, 0.161128, 0.284164, 0.413497, 0.543076, 0.666582, 0.777804, 0.871026, 0.941379, 0.985147,
			1,
			0.985147, 0.941379, 0.871026, 0.777804, 0.666582, 0.543076, 0.413497, 0.284164, 0.161128, 0.049813,
		},
	},
	{
		name: "Triangular", fn: Triangular, fnCmplx: TriangularComplex,
		want: []float64{
			0.05, 0.15, 0.25, 0.35, 0.45, 0.55, 0.65, 0.75, 0.85, 0.95,
			0.95, 0.85, 0.75, 0.65, 0.55, 0.45, 0.35, 0.25, 0.15, 0.05,
		},
	},
	{
		name: "Hann", fn: Hann, fnCmplx: HannComplex,
		want: []float64{
			0.006155, 0.054496, 0.146447, 0.273005, 0.421783, 0.578217, 0.726995, 0.853553, 0.945503, 0.993844,
			0.993844, 0.945503, 0.853553, 0.726995, 0.578217, 0.421783, 0.273005, 0.146447, 0.054496, 0.006155,
		},
	},
	{
		name: "BartlettHann", fn: BartlettHann, fnCmplx: BartlettHannComplex,
		want: []float64{
			0.016678, 0.077417, 0.171299, 0.291484, 0.428555, 0.571445, 0.708516, 0.828701, 0.922582, 0.983322,
			0.983322, 0.922582, 0.828701, 0.708516, 0.571445, 0.428555, 0.291484, 0.171299, 0.077417, 0.016678,
		},
	},
	{
		name: "Hamming", fn: Hamming, fnCmplx: HammingComplex,
		want: []float64{
			0.092577, 0.136714, 0.220669, 0.336222, 0.472063, 0.614894, 0.750735, 0.866288, 0.950242, 0.994379,
			0.994379, 0.950242, 0.866288, 0.750735, 0.614894, 0.472063, 0.336222, 0.220669, 0.136714, 0.092577,
		},
	},
	{
		name: "Blackman", fn: Blackman, fnCmplx: BlackmanComplex,
		want: []float64{
			0.002240, 0.021519, 0.066446, 0.145982, 0.265698, 0.422133, 0.599972, 0.773553, 0.912526, 0.989929,
			0.989929, 0.912526, 0.773553, 0.599972, 0.422133, 0.265698, 0.145982, 0.066446, 0.021519, 0.002240,
		},
	},
	{
		name: "BlackmanHarris", fn: BlackmanHarris, fnCmplx: BlackmanHarrisComplex,
		want: []float64{
			0.000429, 0.004895, 0.021735, 0.065564, 0.153302, 0.295468, 0.485851, 0.695764, 0.878689, 0.985801,
			0.985801, 0.878689, 0.695764, 0.485851, 0.295468, 0.153302, 0.065564, 0.021735, 0.004895, 0.000429,
		},
	},
	{
		name: "Nuttall", fn: Nuttall, fnCmplx: NuttallComplex,
		want: []float64{
			0.000315, 0.004300, 0.020039, 0.062166, 0.148072, 0.289119, 0.479815, 0.691497, 0.876790, 0.985566,
			0.985566, 0.876790, 0.691497, 0.479815, 0.289119, 0.148072, 0.062166, 0.020039, 0.004300, 0.000315,
		},
	},
	{
		name: "BlackmanNuttall", fn: BlackmanNuttall, fnCmplx: BlackmanNuttallComplex,
		want: []float64{
			0.000859, 0.006348, 0.025205, 0.071718, 0.161975, 0.305361, 0.494863, 0.701958, 0.881398, 0.986132,
			0.986132, 0.881398, 0.701958, 0.494863, 0.305361, 0.161975, 0.071718, 0.025205, 0.006348, 0.000859,
		},
	},
	{
		name: "FlatTop", fn: FlatTop, fnCmplx: FlatTopComplex,
		want: []float64{
			-0.001079, -0.007892, -0.026872, -0.056135, -0.069724, -0.015262, 0.157058, 0.444135, 0.760699, 0.970864,
			0.970864, 0.760699, 0.444135, 0.157058, -0.015262, -0.069724, -0.056135, -0.026872, -0.007892, -0.001079,
		},
	},
}

var gausWindowTests = []struct {
	name  string
	sigma float64
	want  []float64
}{
	{
		name: "Gaussian", sigma: 0.3,
		want: []float64{
			0.006645, 0.018063, 0.043936, 0.095634, 0.186270, 0.324652, 0.506336, 0.706648, 0.882497, 0.986207,
			0.986207, 0.882497, 0.706648, 0.506336, 0.324652, 0.186270, 0.095634, 0.043936, 0.018063, 0.006645},
	},
	{
		name: "Gaussian", sigma: 0.5,
		want: []float64{
			0.164474, 0.235746, 0.324652, 0.429557, 0.546074, 0.666977, 0.782705, 0.882497, 0.955997, 0.995012,
			0.995012, 0.955997, 0.882497, 0.782705, 0.666977, 0.546074, 0.429557, 0.324652, 0.235746, 0.164474,
		},
	},
	{
		name: "Gaussian", sigma: 1.2,
		want: []float64{
			0.730981, 0.778125, 0.822578, 0.863552, 0.900293, 0.932102, 0.958357, 0.978532, 0.992218, 0.999132,
			0.999132, 0.992218, 0.978532, 0.958357, 0.932102, 0.900293, 0.863552, 0.822578, 0.778125, 0.730981,
		},
	},
}

func TestWindows(t *testing.T) {
	t.Parallel()
	const tol = 1e-6

	for _, test := range windowTests {
		t.Run(test.name, func(t *testing.T) {
			src := make([]float64, len(test.want))
			for i := range src {
				src[i] = 1
			}

			dst := test.fn(src)
			if !floats.EqualApprox(dst, test.want, tol) {
				t.Errorf("unexpected result for window function %q:\ngot:%#.6v\nwant:%#v", test.name, dst, test.want)
			}

			for i := range src {
				src[i] = 1
			}

			dst = NewValues(test.fn, len(src)).Transform(src)
			if !floats.EqualApprox(dst, test.want, tol) {
				t.Errorf("unexpected result for lookup window function %q:\ngot:%#.6v\nwant:%#.6v", test.name, dst, test.want)
			}
		})
	}
}

func TestGausWindows(t *testing.T) {
	t.Parallel()
	const tol = 1e-6

	for _, test := range gausWindowTests {
		t.Run(fmt.Sprintf("%s (sigma=%.1f)", test.name, test.sigma), func(t *testing.T) {
			src := make([]float64, 20)
			for i := range src {
				src[i] = 1
			}

			gaussian := Gaussian{test.sigma}

			dst := gaussian.Transform(src)
			if !floats.EqualApprox(dst, test.want, tol) {
				t.Errorf("unexpected result for window function %q:\ngot:%#.6v\nwant:%#v", test.name, dst, test.want)
			}

			for i := range src {
				src[i] = 1
			}

			dst = NewValues(gaussian.Transform, len(src)).Transform(src)
			if !floats.EqualApprox(dst, test.want, tol) {
				t.Errorf("unexpected result for lookup window function %q:\ngot:%#.6v\nwant:%#.6v", test.name, dst, test.want)
			}
		})
	}
}

func TestWindowsComplex(t *testing.T) {
	t.Parallel()
	const tol = 1e-6

	for _, test := range windowTests {
		t.Run(test.name+"Complex", func(t *testing.T) {
			src := make([]complex128, len(test.want))
			for i := range src {
				src[i] = complex(1, 1)
			}

			dst := test.fnCmplx(src)
			if !equalApprox(dst, test.want, tol) {
				t.Errorf("unexpected result for window function %q:\ngot:%#.6v\nwant:%#.6v", test.name, dst, test.want)
			}

			for i := range src {
				src[i] = complex(1, 1)
			}

			dst = NewValuesComplex(test.fnCmplx, len(src)).Transform(src)
			if !equalApprox(dst, test.want, tol) {
				t.Errorf("unexpected result for lookup window function %q:\ngot:%#.6v\nwant:%#.6v", test.name, dst, test.want)
			}
		})
	}
}

func TestGausWindowComplex(t *testing.T) {
	t.Parallel()
	const tol = 1e-6

	for _, test := range gausWindowTests {
		t.Run(fmt.Sprintf("%sComplex (sigma=%.1f)", test.name, test.sigma), func(t *testing.T) {
			src := make([]complex128, 20)
			for i := range src {
				src[i] = complex(1, 1)
			}

			gaussian := GaussianComplex{test.sigma}

			dst := gaussian.Transform(src)
			if !equalApprox(dst, test.want, tol) {
				t.Errorf("unexpected result for window function %q:\ngot:%#.6v\nwant:%#.6v", test.name, dst, test.want)
			}

			for i := range src {
				src[i] = complex(1, 1)
			}

			dst = NewValuesComplex(gaussian.Transform, len(src)).Transform(src)
			if !equalApprox(dst, test.want, tol) {
				t.Errorf("unexpected result for lookup window function %q:\ngot:%#.6v\nwant:%#.6v", test.name, dst, test.want)
			}
		})
	}
}

func equalApprox(seq1 []complex128, seq2 []float64, tol float64) bool {
	if len(seq1) != len(seq2) {
		return false
	}
	for i := range seq1 {
		if !scalar.EqualWithinAbsOrRel(real(seq1[i]), seq2[i], tol, tol) {
			return false
		}
		if !scalar.EqualWithinAbsOrRel(imag(seq1[i]), seq2[i], tol, tol) {
			return false
		}
	}
	return true
}
