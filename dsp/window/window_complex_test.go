// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package window

import (
	"math"
	"testing"
	"unsafe"
)

type testingCaseCmplx struct {
	Name           string
	Function       func([]complex128) []complex128
	ExpectedResult []float64
}

type testingCaseCmplxGauss struct {
	Name           string
	Sigma          float64
	ExpectedResult []float64
}

func equalApprox(result []complex128, expectedResult []float64, tol float64) bool {
	if len(result) != len(expectedResult) {
		return false
	}
	for i := range result {
		if math.Abs(real(result[i])-expectedResult[i]) > tol {
			return false
		}
		if math.Abs(imag(result[i])-expectedResult[i]) > tol {
			return false
		}
	}
	return true
}

func TestWindowComplex(t *testing.T) {
	//test precission
	const tol = 1e-6
	//Input data
	src := make([]complex128, 20, 20)
	for i := range src {
		src[i] = complex(1.0, 1.0)
	}

	//testing table (except gauss window)
	var tt []*testingCaseCmplx = []*testingCaseCmplx{
		&testingCaseCmplx{Name: "RectangleComplex", Function: RectangleComplex,
			ExpectedResult: []float64{1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
				1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000}},
		&testingCaseCmplx{Name: "SinComplex", Function: SinComplex,
			ExpectedResult: []float64{0.000000, 0.164595, 0.324699, 0.475947, 0.614213, 0.735724, 0.837166, 0.915773, 0.969400, 0.996584,
				0.996584, 0.969400, 0.915773, 0.837166, 0.735724, 0.614213, 0.475947, 0.324699, 0.164595, 0.000000}},
		&testingCaseCmplx{Name: "LanczosComplex", Function: LanczosComplex,
			ExpectedResult: []float64{0.000000, 0.115514, 0.247646, 0.389468, 0.532984, 0.669692, 0.791213, 0.889915, 0.959492, 0.995450,
				0.995450, 0.959492, 0.889915, 0.791213, 0.669692, 0.532984, 0.389468, 0.247646, 0.115514, 0.000000}},
		&testingCaseCmplx{Name: "BartlettComplex", Function: BartlettComplex,
			ExpectedResult: []float64{0.000000, 0.105263, 0.210526, 0.315789, 0.421053, 0.526316, 0.631579, 0.736842, 0.842105, 0.947368,
				0.947368, 0.842105, 0.736842, 0.631579, 0.526316, 0.421053, 0.315789, 0.210526, 0.105263, 0.000000}},
		&testingCaseCmplx{Name: "HannComplex", Function: HannComplex,
			ExpectedResult: []float64{0.000000, 0.027091, 0.105430, 0.226526, 0.377257, 0.541290, 0.700848, 0.838641, 0.939737, 0.993181,
				0.993181, 0.939737, 0.838641, 0.700848, 0.541290, 0.377257, 0.226526, 0.105430, 0.027091, 0.000000}},
		&testingCaseCmplx{Name: "BartlettHannComplex", Function: BartlettHannComplex,
			ExpectedResult: []float64{0.000000, 0.045853, 0.130653, 0.247949, 0.387768, 0.537696, 0.684223, 0.814209, 0.916305, 0.982186,
				0.982186, 0.916305, 0.814209, 0.684223, 0.537696, 0.387768, 0.247949, 0.130653, 0.045853, 0.000000}},
		&testingCaseCmplx{Name: "HammingComplex", Function: HammingComplex,
			ExpectedResult: []float64{0.080000, 0.104924, 0.176995, 0.288404, 0.427077, 0.577986, 0.724780, 0.851550, 0.944558, 0.993726,
				0.993726, 0.944558, 0.851550, 0.724780, 0.577986, 0.427077, 0.288404, 0.176995, 0.104924, 0.080000}},
		&testingCaseCmplx{Name: "BlackmanComplex", Function: BlackmanComplex,
			ExpectedResult: []float64{0.000000, 0.010223, 0.045069, 0.114390, 0.226899, 0.382381, 0.566665, 0.752034, 0.903493, 0.988846,
				0.988846, 0.903493, 0.752034, 0.566665, 0.382381, 0.226899, 0.114390, 0.045069, 0.010223, 0.000000}},
		&testingCaseCmplx{Name: "BlackmanHarrisComplex", Function: BlackmanHarrisComplex,
			ExpectedResult: []float64{0.000060, 0.002018, 0.012795, 0.046450, 0.122540, 0.256852, 0.448160, 0.668576, 0.866426, 0.984278,
				0.984278, 0.866426, 0.668576, 0.448160, 0.256852, 0.122540, 0.046450, 0.012795, 0.002018, 0.000060}},
		&testingCaseCmplx{Name: "NuttallComplex", Function: NuttallComplex,
			ExpectedResult: []float64{0.000000, 0.001706, 0.011614, 0.043682, 0.117808, 0.250658, 0.441946, 0.664015, 0.864348, 0.984019,
				0.984019, 0.864348, 0.664015, 0.441946, 0.250658, 0.117808, 0.043682, 0.011614, 0.001706, 0.000000}},
		&testingCaseCmplx{Name: "BlackmanNuttallComplex", Function: BlackmanNuttallComplex,
			ExpectedResult: []float64{0.000363, 0.002885, 0.015360, 0.051652, 0.130567, 0.266629, 0.457501, 0.675215, 0.869392, 0.984644,
				0.984644, 0.869392, 0.675215, 0.457501, 0.266629, 0.130567, 0.051652, 0.015360, 0.002885, 0.000363}},
		&testingCaseCmplx{Name: "FlatTopComplex", Function: FlatTopComplex,
			ExpectedResult: []float64{0.004000, -0.011796, -0.078650, -0.212762, -0.328021, -0.178010, 0.531959, 1.862876, 3.422134, 4.490270,
				4.490270, 3.422134, 1.862876, 0.531959, -0.178010, -0.328021, -0.212762, -0.078650, -0.011796, 0.004000}},
	}

	//testing table (gauss window)
	var ttg []*testingCaseCmplxGauss = []*testingCaseCmplxGauss{
		&testingCaseCmplxGauss{Name: "GaussComplex (sigma=0.3)", Sigma: 0.3,
			ExpectedResult: []float64{0.003866, 0.011708, 0.031348, 0.074214, 0.155344, 0.287499, 0.470444, 0.680632, 0.870660, 0.984728,
				0.984728, 0.870660, 0.680632, 0.470444, 0.287499, 0.155344, 0.074214, 0.031348, 0.011708, 0.003866}},
		&testingCaseCmplxGauss{Name: "GaussComplex (sigma=0.5)", Sigma: 0.5,
			ExpectedResult: []float64{0.135335, 0.201673, 0.287499, 0.392081, 0.511524, 0.638423, 0.762260, 0.870660, 0.951361, 0.994475,
				0.994475, 0.951361, 0.870660, 0.762260, 0.638423, 0.511524, 0.392081, 0.287499, 0.201673, 0.135335}},
		&testingCaseCmplxGauss{Name: "GaussComplex (sigma=1.2)", Sigma: 1.2,
			ExpectedResult: []float64{0.706648, 0.757319, 0.805403, 0.849974, 0.890135, 0.925049, 0.953963, 0.976241, 0.991381, 0.999039,
				0.999039, 0.991381, 0.976241, 0.953963, 0.925049, 0.890135, 0.849974, 0.805403, 0.757319, 0.706648}},
	}

	// run tests on testing tables
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			//copy src, because window functions change data inplace
			srcCpy := append([]complex128(nil), src...)
			dst := tc.Function(srcCpy)

			//check that function returned the slice with the same array, not a newly allocated
			addressSrc := unsafe.Pointer(&srcCpy)
			addArraySrc := unsafe.Pointer(*(*uintptr)(addressSrc))
			addressDst := unsafe.Pointer(&dst)
			addArrayDst := unsafe.Pointer(*(*uintptr)(addressDst))
			if addArraySrc != addArrayDst {
				t.Errorf("unexpected result for window function %q : function returned a slice with new underling array at %v instead of source array at %v modification ",
					tc.Name, addArrayDst, addArraySrc)
			}

			//check the data
			if !equalApprox(dst, tc.ExpectedResult, tol) {
				t.Errorf("unexpected result for window function %q: result \n%f\n is different from expectation \n%f\n more then it is allowed by test tolerance %e",
					tc.Name, dst, tc.ExpectedResult, tol)
			}
		})
	}

	for _, tc := range ttg {
		t.Run(tc.Name, func(t *testing.T) {
			//copy src, because window functions change data inplace
			srcCpy := append([]complex128(nil), src...)
			dst := GaussComplex(srcCpy, tc.Sigma)

			//check that function returned the slice with the same array, not a newly allocated
			addressSrc := unsafe.Pointer(&srcCpy)
			addArraySrc := unsafe.Pointer(*(*uintptr)(addressSrc))
			addressDst := unsafe.Pointer(&dst)
			addArrayDst := unsafe.Pointer(*(*uintptr)(addressDst))
			if addArraySrc != addArrayDst {
				t.Errorf("unexpected result for window function %q : function returned a slice with new underling array at %v instead of source array at %v modification ",
					tc.Name, addArrayDst, addArraySrc)
			}

			//check the data
			if !equalApprox(dst, tc.ExpectedResult, tol) {
				t.Errorf("unexpected result for window function %q: result \n%f\n is different from expectation \n%f\n more then it is allowed by test tolerance %e",
					tc.Name, dst, tc.ExpectedResult, tol)
			}
		})
	}
}
