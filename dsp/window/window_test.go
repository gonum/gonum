// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package window

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/cmplxs"
	"gonum.org/v1/gonum/floats"
)

const tol = 1e-6
const (
	gaussWin = iota
	tukeyWin
)

// these types describe zero, mono, and multi parameter windowing
// functions for both real and complex data
//
// when looking at the truth tables, it is helpful to know that
// Go treats []complex128{1} as [1+0i].  More specifically, a variable which
// is of type complex128 to which you assign an untyped constant has a zero
// imaginary part.

type transformFunc func([]float64) []float64
type transformFuncComplex func([]complex128) []complex128

type transformer interface {
	Transform([]float64) []float64
}

type transformerComplex interface {
	TransformComplex([]complex128) []complex128
}

type dualTransformer interface {
	transformer
	transformerComplex
}

var windowTestsReal = []struct {
	name string
	fn   transformFunc
	want []float64
}{
	{
		name: "Rectangular", fn: Rectangular,
		want: []float64{
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		},
	},
	{
		name: "Sine", fn: Sine,
		want: []float64{
			0.078459, 0.233445, 0.382683, 0.522499, 0.649448, 0.760406, 0.852640, 0.923880, 0.972370, 0.996917,
			0.996917, 0.972370, 0.923880, 0.852640, 0.760406, 0.649448, 0.522499, 0.382683, 0.233445, 0.078459,
		},
	},
	{
		name: "Lanczos", fn: Lanczos,
		want: []float64{
			0.052415, 0.170011, 0.300105, 0.436333, 0.57162, 0.698647, 0.810332, 0.900316, 0.963398, 0.995893,
			0.995893, 0.963398, 0.900316, 0.810332, 0.698647, 0.57162, 0.436333, 0.300105, 0.170011, 0.052415,
		},
	},
	// This case tests Lanczos for a NaN condition. The Lanczos NaN condition is k=(N-1)/2, that is when N is odd.
	{
		name: "LanczosOdd", fn: Lanczos,
		want: []float64{
			0.049813, 0.161128, 0.284164, 0.413497, 0.543076, 0.666582, 0.777804, 0.871026, 0.941379, 0.985147,
			1,
			0.985147, 0.941379, 0.871026, 0.777804, 0.666582, 0.543076, 0.413497, 0.284164, 0.161128, 0.049813,
		},
	},
	{
		name: "Triangular", fn: Triangular,
		want: []float64{
			0.05, 0.15, 0.25, 0.35, 0.45, 0.55, 0.65, 0.75, 0.85, 0.95,
			0.95, 0.85, 0.75, 0.65, 0.55, 0.45, 0.35, 0.25, 0.15, 0.05,
		},
	},
	{
		name: "Hann", fn: Hann,
		want: []float64{
			0.006155, 0.054496, 0.146447, 0.273005, 0.421783, 0.578217, 0.726995, 0.853553, 0.945503, 0.993844,
			0.993844, 0.945503, 0.853553, 0.726995, 0.578217, 0.421783, 0.273005, 0.146447, 0.054496, 0.006155,
		},
	},
	{
		name: "BartlettHann", fn: BartlettHann,
		want: []float64{
			0.016678, 0.077417, 0.171299, 0.291484, 0.428555, 0.571445, 0.708516, 0.828701, 0.922582, 0.983322,
			0.983322, 0.922582, 0.828701, 0.708516, 0.571445, 0.428555, 0.291484, 0.171299, 0.077417, 0.016678,
		},
	},
	{
		name: "Hamming", fn: Hamming,
		want: []float64{
			0.092577, 0.136714, 0.220669, 0.336222, 0.472063, 0.614894, 0.750735, 0.866288, 0.950242, 0.994379,
			0.994379, 0.950242, 0.866288, 0.750735, 0.614894, 0.472063, 0.336222, 0.220669, 0.136714, 0.092577,
		},
	},
	{
		name: "Blackman", fn: Blackman,
		want: []float64{
			0.002240, 0.021519, 0.066446, 0.145982, 0.265698, 0.422133, 0.599972, 0.773553, 0.912526, 0.989929,
			0.989929, 0.912526, 0.773553, 0.599972, 0.422133, 0.265698, 0.145982, 0.066446, 0.021519, 0.002240,
		},
	},
	{
		name: "BlackmanHarris", fn: BlackmanHarris,
		want: []float64{
			0.000429, 0.004895, 0.021735, 0.065564, 0.153302, 0.295468, 0.485851, 0.695764, 0.878689, 0.985801,
			0.985801, 0.878689, 0.695764, 0.485851, 0.295468, 0.153302, 0.065564, 0.021735, 0.004895, 0.000429,
		},
	},
	{
		name: "Nuttall", fn: Nuttall,
		want: []float64{
			0.000315, 0.004300, 0.020039, 0.062166, 0.148072, 0.289119, 0.479815, 0.691497, 0.876790, 0.985566,
			0.985566, 0.876790, 0.691497, 0.479815, 0.289119, 0.148072, 0.062166, 0.020039, 0.004300, 0.000315,
		},
	},
	{
		name: "BlackmanNuttall", fn: BlackmanNuttall,
		want: []float64{
			0.000859, 0.006348, 0.025205, 0.071718, 0.161975, 0.305361, 0.494863, 0.701958, 0.881398, 0.986132,
			0.986132, 0.881398, 0.701958, 0.494863, 0.305361, 0.161975, 0.071718, 0.025205, 0.006348, 0.000859,
		},
	},
	{
		name: "FlatTop", fn: FlatTop,
		want: []float64{
			-0.001079, -0.007892, -0.026872, -0.056135, -0.069724, -0.015262, 0.157058, 0.444135, 0.760699, 0.970864,
			0.970864, 0.760699, 0.444135, 0.157058, -0.015262, -0.069724, -0.056135, -0.026872, -0.007892, -0.001079,
		},
	},
}

var windowTestsComplex = []struct {
	name string
	fn   transformFuncComplex
	want []complex128
}{
	{
		name: "Rectangular", fn: RectangularComplex,
		want: []complex128{
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		},
	},
	{
		name: "Sine", fn: SineComplex,
		want: []complex128{
			0.078459, 0.233445, 0.382683, 0.522499, 0.649448, 0.760406, 0.852640, 0.923880, 0.972370, 0.996917,
			0.996917, 0.972370, 0.923880, 0.852640, 0.760406, 0.649448, 0.522499, 0.382683, 0.233445, 0.078459,
		},
	},
	{
		name: "Lanczos", fn: LanczosComplex,
		want: []complex128{
			0.052415, 0.170011, 0.300105, 0.436333, 0.57162, 0.698647, 0.810332, 0.900316, 0.963398, 0.995893,
			0.995893, 0.963398, 0.900316, 0.810332, 0.698647, 0.57162, 0.436333, 0.300105, 0.170011, 0.052415,
		},
	},
	// This case tests Lanczos for a NaN condition. The Lanczos NaN condition is k=(N-1)/2, that is when N is odd.
	{
		name: "LanczosOdd", fn: LanczosComplex,
		want: []complex128{
			0.049813, 0.161128, 0.284164, 0.413497, 0.543076, 0.666582, 0.777804, 0.871026, 0.941379, 0.985147,
			1,
			0.985147, 0.941379, 0.871026, 0.777804, 0.666582, 0.543076, 0.413497, 0.284164, 0.161128, 0.049813,
		},
	},
	{
		name: "Triangular", fn: TriangularComplex,
		want: []complex128{
			0.05, 0.15, 0.25, 0.35, 0.45, 0.55, 0.65, 0.75, 0.85, 0.95,
			0.95, 0.85, 0.75, 0.65, 0.55, 0.45, 0.35, 0.25, 0.15, 0.05,
		},
	},
	{
		name: "Hann", fn: HannComplex,
		want: []complex128{
			0.006155, 0.054496, 0.146447, 0.273005, 0.421783, 0.578217, 0.726995, 0.853553, 0.945503, 0.993844,
			0.993844, 0.945503, 0.853553, 0.726995, 0.578217, 0.421783, 0.273005, 0.146447, 0.054496, 0.006155,
		},
	},
	{
		name: "BartlettHann", fn: BartlettHannComplex,
		want: []complex128{
			0.016678, 0.077417, 0.171299, 0.291484, 0.428555, 0.571445, 0.708516, 0.828701, 0.922582, 0.983322,
			0.983322, 0.922582, 0.828701, 0.708516, 0.571445, 0.428555, 0.291484, 0.171299, 0.077417, 0.016678,
		},
	},
	{
		name: "Hamming", fn: HammingComplex,
		want: []complex128{
			0.092577, 0.136714, 0.220669, 0.336222, 0.472063, 0.614894, 0.750735, 0.866288, 0.950242, 0.994379,
			0.994379, 0.950242, 0.866288, 0.750735, 0.614894, 0.472063, 0.336222, 0.220669, 0.136714, 0.092577,
		},
	},
	{
		name: "Blackman", fn: BlackmanComplex,
		want: []complex128{
			0.002240, 0.021519, 0.066446, 0.145982, 0.265698, 0.422133, 0.599972, 0.773553, 0.912526, 0.989929,
			0.989929, 0.912526, 0.773553, 0.599972, 0.422133, 0.265698, 0.145982, 0.066446, 0.021519, 0.002240,
		},
	},
	{
		name: "BlackmanHarris", fn: BlackmanHarrisComplex,
		want: []complex128{
			0.000429, 0.004895, 0.021735, 0.065564, 0.153302, 0.295468, 0.485851, 0.695764, 0.878689, 0.985801,
			0.985801, 0.878689, 0.695764, 0.485851, 0.295468, 0.153302, 0.065564, 0.021735, 0.004895, 0.000429,
		},
	},
	{
		name: "Nuttall", fn: NuttallComplex,
		want: []complex128{
			0.000315, 0.004300, 0.020039, 0.062166, 0.148072, 0.289119, 0.479815, 0.691497, 0.876790, 0.985566,
			0.985566, 0.876790, 0.691497, 0.479815, 0.289119, 0.148072, 0.062166, 0.020039, 0.004300, 0.000315,
		},
	},
	{
		name: "BlackmanNuttall", fn: BlackmanNuttallComplex,
		want: []complex128{
			0.000859, 0.006348, 0.025205, 0.071718, 0.161975, 0.305361, 0.494863, 0.701958, 0.881398, 0.986132,
			0.986132, 0.881398, 0.701958, 0.494863, 0.305361, 0.161975, 0.071718, 0.025205, 0.006348, 0.000859,
		},
	},
	{
		name: "FlatTop", fn: FlatTopComplex,
		want: []complex128{
			-0.001079, -0.007892, -0.026872, -0.056135, -0.069724, -0.015262, 0.157058, 0.444135, 0.760699, 0.970864,
			0.970864, 0.760699, 0.444135, 0.157058, -0.015262, -0.069724, -0.056135, -0.026872, -0.007892, -0.001079,
		},
	},
}

var monoParamWindowTests = []struct {
	name       string
	param      float64
	windowType int
	want       []float64
	wantCmplx  []complex128
}{
	{
		name: "Gaussian", param: 0.3, windowType: gaussWin,
		want: []float64{
			0.006645, 0.018063, 0.043936, 0.095634, 0.186270, 0.324652, 0.506336, 0.706648, 0.882497, 0.986207,
			0.986207, 0.882497, 0.706648, 0.506336, 0.324652, 0.186270, 0.095634, 0.043936, 0.018063, 0.006645},
		wantCmplx: []complex128{
			0.006645, 0.018063, 0.043936, 0.095634, 0.186270, 0.324652, 0.506336, 0.706648, 0.882497, 0.986207,
			0.986207, 0.882497, 0.706648, 0.506336, 0.324652, 0.186270, 0.095634, 0.043936, 0.018063, 0.006645},
	},
	{
		name: "Gaussian", param: 0.5, windowType: gaussWin,
		want: []float64{
			0.164474, 0.235746, 0.324652, 0.429557, 0.546074, 0.666977, 0.782705, 0.882497, 0.955997, 0.995012,
			0.995012, 0.955997, 0.882497, 0.782705, 0.666977, 0.546074, 0.429557, 0.324652, 0.235746, 0.164474,
		},
		wantCmplx: []complex128{
			0.164474, 0.235746, 0.324652, 0.429557, 0.546074, 0.666977, 0.782705, 0.882497, 0.955997, 0.995012,
			0.995012, 0.955997, 0.882497, 0.782705, 0.666977, 0.546074, 0.429557, 0.324652, 0.235746, 0.164474,
		},
	},
	{
		name: "Gaussian", param: 1.2, windowType: gaussWin,
		want: []float64{
			0.730981, 0.778125, 0.822578, 0.863552, 0.900293, 0.932102, 0.958357, 0.978532, 0.992218, 0.999132,
			0.999132, 0.992218, 0.978532, 0.958357, 0.932102, 0.900293, 0.863552, 0.822578, 0.778125, 0.730981,
		},
		wantCmplx: []complex128{
			0.730981, 0.778125, 0.822578, 0.863552, 0.900293, 0.932102, 0.958357, 0.978532, 0.992218, 0.999132,
			0.999132, 0.992218, 0.978532, 0.958357, 0.932102, 0.900293, 0.863552, 0.822578, 0.778125, 0.730981,
		},
	},
	{
		name: "Tukey", param: 1, windowType: tukeyWin,
		want: []float64{ // copied from Hann
			0.006155, 0.054496, 0.146447, 0.273005, 0.421783, 0.578217, 0.726995, 0.853553, 0.945503, 0.993844,
			0.993844, 0.945503, 0.853553, 0.726995, 0.578217, 0.421783, 0.273005, 0.146447, 0.054496, 0.006155,
		},
		wantCmplx: []complex128{ // copied from Hann
			0.006155, 0.054496, 0.146447, 0.273005, 0.421783, 0.578217, 0.726995, 0.853553, 0.945503, 0.993844,
			0.993844, 0.945503, 0.853553, 0.726995, 0.578217, 0.421783, 0.273005, 0.146447, 0.054496, 0.006155,
		},
	},
	{
		name: "Tukey", param: 0, windowType: tukeyWin,
		want: []float64{ // copied from rectangular
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		},
		wantCmplx: []complex128{ // copied from rectangular
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		},
	},
	{
		name: "Tukey", param: 0.5, windowType: tukeyWin,
		want: []float64{
			0.000000, 0.105430, 0.377257, 0.700847, 0.939737, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
			1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 0.939737, 0.700847, 0.377257, 0.105429, 0.000000,
		},
		wantCmplx: []complex128{
			0.000000, 0.105430, 0.377257, 0.700847, 0.939737, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
			1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 0.939737, 0.700847, 0.377257, 0.105429, 0.000000,
		},
	},
}

func TestWindows(t *testing.T) {
	t.Parallel()

	for _, test := range windowTestsReal {
		t.Run(test.name, func(t *testing.T) {
			src := make([]float64, len(test.want))
			for i := range src {
				src[i] = 1
			}

			dst := test.fn(src)
			if !floats.EqualApprox(dst, test.want, tol) {
				t.Errorf("unexpected result for window function %q:\ngot:%#.6v\nwant:%#v", test.name, dst, test.want)
			}
		})
	}
}

func TestWindowsComplex(t *testing.T) {
	t.Parallel()

	for _, test := range windowTestsComplex {
		t.Run(test.name+"Complex", func(t *testing.T) {
			src := make([]complex128, len(test.want))
			for i := range src {
				src[i] = 1 + 0i
			}

			dst := test.fn(src)
			if !cmplxs.EqualApprox(dst, test.want, tol) {
				t.Errorf("unexpected result for window function %q:\ngot:%#.6v\nwant:%#.6v", test.name, dst, test.want)
			}
		})
	}
}

func TestValuesWindows(t *testing.T) {
	win := Values{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	src := make([]float64, len(win))
	for i := range src {
		src[i] = 1
	}

	dst := win.Transform(src)
	if !floats.EqualApprox(dst, win, tol) {
		t.Errorf("unexpected result for lookup window function: got:%#.6v\nwant:%#v", dst, win)
	}
}

func TestValuesWindowsComplex(t *testing.T) {
	win := ValuesComplex{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	src := make([]complex128, len(win))
	for i := range src {
		src[i] = 1
	}

	dst := win.Transform(src)
	if !cmplxs.EqualApprox(dst, win, tol) {
		t.Errorf("unexpected result for lookup window function: got:%#.6v\nwant:%#v", dst, win)
	}
}

func TestMonoParametricWindows(t *testing.T) {
	t.Parallel()

	for _, test := range monoParamWindowTests {
		t.Run(fmt.Sprintf("%s (param=%.1f)", test.name, test.param), func(t *testing.T) {
			trans := []dualTransformer{
				Gaussian{test.param},
				Tukey{test.param},
			}[test.windowType]

			srcRe := make([]float64, len(test.want))
			for i := range srcRe {
				srcRe[i] = 1
			}
			dst := trans.Transform(srcRe)
			if !floats.EqualApprox(dst, test.want, tol) {
				t.Errorf("unexpected result for window function %q:\ngot:%#.6v\nwant:%#v", test.name, dst, test.want)
			}

			srcCmplx := make([]complex128, len(test.wantCmplx))
			for i := range srcCmplx {
				srcCmplx[i] = 1 + 0i
			}
			dstCmplx := trans.TransformComplex(srcCmplx)
			if !cmplxs.EqualApprox(dstCmplx, test.wantCmplx, tol) {
				t.Errorf("unexpected result for window function %q:\ngot:%#.6v\nwant:%#v", test.name, dstCmplx, test.wantCmplx)
			}
		})
	}
}
