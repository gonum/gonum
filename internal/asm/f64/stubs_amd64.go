// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

package f64

func AbsSum(x []float64) (sum float64)

func AbsSumInc(x []float64, n, incX int) (sum float64)

func AddConst(alpha float64, x []float64)

func Add(dst, s []float64)

func AxpyUnitary(alpha float64, x, y []float64)

func AxpyUnitaryTo(dst []float64, alpha float64, x, y []float64)

func AxpyInc(alpha float64, x, y []float64, n, incX, incY, ix, iy uintptr)

func AxpyIncTo(dst []float64, incDst, idst uintptr, alpha float64, x, y []float64, n, incX, incY, ix, iy uintptr)

func CumSum(dst, s []float64) []float64

func CumProd(dst, s []float64) []float64

func Div(dst, s []float64)

func DivTo(dst, x, y []float64) []float64

func DotUnitary(x, y []float64) (sum float64)

func DotInc(x, y []float64, n, incX, incY, ix, iy uintptr) (sum float64)

func L1Norm(s, t []float64) float64

func LinfNorm(s, t []float64) float64

func ScalUnitary(alpha float64, x []float64)

func ScalUnitaryTo(dst []float64, alpha float64, x []float64)

// incX must be positive.
func ScalInc(alpha float64, x []float64, n, incX uintptr)

// incDst and incX must be positive.
func ScalIncTo(dst []float64, incDst uintptr, alpha float64, x []float64, n, incX uintptr)
