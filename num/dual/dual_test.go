// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dual

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
)

// First derivatives:

func dSin(x float64) float64  { return math.Cos(x) }
func dCos(x float64) float64  { return -math.Sin(x) }
func dTan(x float64) float64  { return sec(x) * sec(x) }
func dAsin(x float64) float64 { return 1 / math.Sqrt(1-x*x) }
func dAcos(x float64) float64 { return -1 / math.Sqrt(1-x*x) }
func dAtan(x float64) float64 { return 1 / (1 + x*x) }

func dSinh(x float64) float64  { return math.Cosh(x) }
func dCosh(x float64) float64  { return math.Sinh(x) }
func dTanh(x float64) float64  { return sech(x) * sech(x) }
func dAsinh(x float64) float64 { return 1 / math.Sqrt(x*x+1) }
func dAcosh(x float64) float64 { return 1 / (math.Sqrt(x-1) * math.Sqrt(x+1)) }
func dAtanh(x float64) float64 {
	switch {
	case math.Abs(x) == 1:
		return math.NaN()
	case math.IsInf(x, 0):
		return negZero
	}
	return 1 / (1 - x*x)
}

func dExp(x float64) float64    { return math.Exp(x) }
func dLog(x float64) float64    { return 1 / x }
func dPow(x, y float64) float64 { return y * math.Pow(x, y-1) }
func dSqrt(x float64) float64 {
	// For whatever reason, math.Sqrt(-0) returns -0.
	// In this case, that is clearly a wrong approach.
	if x == 0 {
		return math.Inf(1)
	}
	return 0.5 / math.Sqrt(x)
}
func dInv(x float64) float64 { return -1 / (x * x) }

// Helpers:

func sec(x float64) float64  { return 1 / math.Cos(x) }
func sech(x float64) float64 { return 1 / math.Cosh(x) }

var negZero = math.Float64frombits(1 << 63)

var dualTests = []struct {
	name   string
	x      []float64
	fnDual func(x Number) Number
	fn     func(x float64) float64
	dFn    func(x float64) float64
}{
	{
		name:   "sin",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Sin,
		fn:     math.Sin,
		dFn:    dSin,
	},
	{
		name:   "cos",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Cos,
		fn:     math.Cos,
		dFn:    dCos,
	},
	{
		name:   "tan",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Tan,
		fn:     math.Tan,
		dFn:    dTan,
	},
	{
		name:   "sinh",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Sinh,
		fn:     math.Sinh,
		dFn:    dSinh,
	},
	{
		name:   "cosh",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Cosh,
		fn:     math.Cosh,
		dFn:    dCosh,
	},
	{
		name:   "tanh",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Tanh,
		fn:     math.Tanh,
		dFn:    dTanh,
	},

	{
		name:   "asin",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Asin,
		fn:     math.Asin,
		dFn:    dAsin,
	},
	{
		name:   "acos",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Acos,
		fn:     math.Acos,
		dFn:    dAcos,
	},
	{
		name:   "atan",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Atan,
		fn:     math.Atan,
		dFn:    dAtan,
	},
	{
		name:   "asinh",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Asinh,
		fn:     math.Asinh,
		dFn:    dAsinh,
	},
	{
		name:   "acosh",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Acosh,
		fn:     math.Acosh,
		dFn:    dAcosh,
	},
	{
		name:   "atanh",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Atanh,
		fn:     math.Atanh,
		dFn:    dAtanh,
	},

	{
		name:   "exp",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Exp,
		fn:     math.Exp,
		dFn:    dExp,
	},
	{
		name:   "log",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Log,
		fn:     math.Log,
		dFn:    dLog,
	},
	{
		name:   "inv",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Inv,
		fn:     func(x float64) float64 { return 1 / x },
		dFn:    dInv,
	},
	{
		name:   "sqrt",
		x:      []float64{math.NaN(), math.Inf(-1), -3, -2, -1, -0.5, negZero, 0, 0.5, 1, 2, 3, math.Inf(1)},
		fnDual: Sqrt,
		fn:     math.Sqrt,
		dFn:    dSqrt,
	},

	{
		name: "Fike example fn",
		x:    []float64{1, 2, 3, 4, 5},
		fnDual: func(x Number) Number {
			return Mul(
				Exp(x),
				Inv(Sqrt(
					Add(
						PowReal(Sin(x), 3),
						PowReal(Cos(x), 3)))))
		},
		fn: func(x float64) float64 {
			return math.Exp(x) / math.Sqrt(math.Pow(math.Sin(x), 3)+math.Pow(math.Cos(x), 3))
		},
		dFn: func(x float64) float64 {
			return math.Exp(x) * (3*math.Cos(x) + 5*math.Cos(3*x) + 9*math.Sin(x) + math.Sin(3*x)) /
				(8 * math.Pow(math.Pow(math.Sin(x), 3)+math.Pow(math.Cos(x), 3), 1.5))
		},
	},
}

func TestDual(t *testing.T) {
	const tol = 1e-15
	for _, test := range dualTests {
		for _, x := range test.x {
			fxDual := test.fnDual(Number{Real: x, Emag: 1})
			fx := test.fn(x)
			dFx := test.dFn(x)
			if !same(fxDual.Real, fx, tol) {
				t.Errorf("unexpected %s(%v): got:%v want:%v", test.name, x, fxDual.Real, fx)
			}
			if !same(fxDual.Emag, dFx, tol) {
				t.Errorf("unexpected %s'(%v): got:%v want:%v", test.name, x, fxDual.Emag, dFx)
			}
		}
	}
}

func same(a, b, tol float64) bool {
	return (math.IsNaN(a) && math.IsNaN(b)) || floats.EqualWithinAbsOrRel(a, b, tol, tol)
}
