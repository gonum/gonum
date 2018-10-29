// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hyperdual

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
func dAtanh(x float64) float64 { return 1 / (1 - x*x) }

func dExp(x float64) float64    { return math.Exp(x) }
func dLog(x float64) float64    { return 1 / x }
func dPow(x, y float64) float64 { return y * math.Pow(x, y-1) }
func dSqrt(x float64) float64   { return 0.5 * math.Pow(x, -0.5) }
func dInv(x float64) float64    { return -1 / (x * x) }

// Second derivatives:

func d2Sin(x float64) float64  { return -math.Sin(x) }
func d2Cos(x float64) float64  { return -math.Cos(x) }
func d2Tan(x float64) float64  { return 2 * math.Tan(x) * sec(x) * sec(x) }
func d2Asin(x float64) float64 { return x / math.Pow(1-x*x, 1.5) }
func d2Acos(x float64) float64 { return -x / math.Pow(1-x*x, 1.5) }
func d2Atan(x float64) float64 { return -2 * x / ((x*x + 1) * (x*x + 1)) }

func d2Sinh(x float64) float64  { return math.Sinh(x) }
func d2Cosh(x float64) float64  { return math.Cosh(x) }
func d2Tanh(x float64) float64  { return -2 * math.Tanh(x) * sech(x) * sech(x) }
func d2Asinh(x float64) float64 { return -x / math.Pow((x*x+1), 1.5) }
func d2Acosh(x float64) float64 { return -x / (math.Pow(x-1, 1.5) * math.Pow(x+1, 1.5)) }
func d2Atanh(x float64) float64 { return 2 * x / ((1 - x*x) * (1 - x*x)) }

func d2Exp(x float64) float64    { return math.Exp(x) }
func d2Log(x float64) float64    { return -1 / (x * x) }
func d2Pow(x, y float64) float64 { return y * (y - 1) * math.Pow(x, y-2) }
func d2Sqrt(x float64) float64   { return -0.25 * math.Pow(x, -1.5) }
func d2Inv(x float64) float64    { return 2 / (x * x * x) }

// Helpers:

func sec(x float64) float64  { return 1 / math.Cos(x) }
func sech(x float64) float64 { return 1 / math.Cosh(x) }

var hyperdualTests = []struct {
	name        string
	x           []float64
	fnHyperdual func(x Number) Number
	fn          func(x float64) float64
	dFn         func(x float64) float64
	d2Fn        func(x float64) float64
}{
	{
		name:        "sin",
		x:           []float64{0, math.Pi / 4, 1, math.Pi / 2, math.Pi, 2 * math.Pi},
		fnHyperdual: Sin,
		fn:          math.Sin,
		dFn:         dSin,
		d2Fn:        d2Sin,
	},
	{
		name:        "cos",
		x:           []float64{0, math.Pi / 4, 1, math.Pi / 2, math.Pi, 2 * math.Pi},
		fnHyperdual: Cos,
		fn:          math.Cos,
		dFn:         dCos,
		d2Fn:        d2Cos,
	},
	{
		name:        "tan",
		x:           []float64{0, math.Pi / 4, 1, math.Pi / 2, math.Pi, 2 * math.Pi},
		fnHyperdual: Tan,
		fn:          math.Tan,
		dFn:         dTan,
		d2Fn:        d2Tan,
	},
	{
		name:        "sinh",
		x:           []float64{0, math.Pi / 4, 1, math.Pi / 2, math.Pi, 2 * math.Pi},
		fnHyperdual: Sinh,
		fn:          math.Sinh,
		dFn:         dSinh,
		d2Fn:        d2Sinh,
	},
	{
		name:        "cosh",
		x:           []float64{0, math.Pi / 4, 1, math.Pi / 2, math.Pi, 2 * math.Pi},
		fnHyperdual: Cosh,
		fn:          math.Cosh,
		dFn:         dCosh,
		d2Fn:        d2Cosh,
	},
	{
		name:        "tanh",
		x:           []float64{0, math.Pi / 4, 1, math.Pi / 2, math.Pi, 2 * math.Pi},
		fnHyperdual: Tanh,
		fn:          math.Tanh,
		dFn:         dTanh,
		d2Fn:        d2Tanh,
	},

	{
		name:        "asin",
		x:           []float64{0, math.Pi / 4, math.Pi / 2, math.Pi, 2 * math.Pi},
		fnHyperdual: Asin,
		fn:          math.Asin,
		dFn:         dAsin,
		d2Fn:        d2Asin,
	},
	{
		name:        "acos",
		x:           []float64{0, math.Pi / 4, math.Pi / 2, math.Pi, 2 * math.Pi},
		fnHyperdual: Acos,
		fn:          math.Acos,
		dFn:         dAcos,
		d2Fn:        d2Acos,
	},
	{
		name:        "atan",
		x:           []float64{0, math.Pi / 4, 1, math.Pi / 2, math.Pi, 2 * math.Pi},
		fnHyperdual: Atan,
		fn:          math.Atan,
		dFn:         dAtan,
		d2Fn:        d2Atan,
	},
	{
		name:        "asinh",
		x:           []float64{0, math.Pi / 4, 1, math.Pi / 2, math.Pi, 2 * math.Pi},
		fnHyperdual: Asinh,
		fn:          math.Asinh,
		dFn:         dAsinh,
		d2Fn:        d2Asinh,
	},
	{
		name:        "acosh",
		x:           []float64{ /*0,*/ math.Pi / 2, math.Pi, 2 * math.Pi, 5},
		fnHyperdual: Acosh,
		fn:          math.Acosh,
		dFn:         dAcosh,
		d2Fn:        d2Acosh,
	},
	{
		name:        "atanh",
		x:           []float64{0, math.Pi / 4, math.Pi / 2, math.Pi},
		fnHyperdual: Atanh,
		fn:          math.Atanh,
		dFn:         dAtanh,
		d2Fn:        d2Atanh,
	},

	{
		name:        "exp",
		x:           []float64{0, 1, 2, 3},
		fnHyperdual: Exp,
		fn:          math.Exp,
		dFn:         dExp,
		d2Fn:        d2Exp,
	},
	{
		name:        "log",
		x:           []float64{ /*0,*/ 1, 2, 3},
		fnHyperdual: Log,
		fn:          math.Log,
		dFn:         dLog,
		d2Fn:        d2Log,
	},
	{
		name:        "inv",
		x:           []float64{ /*0,*/ 1, 2, 3},
		fnHyperdual: Inv,
		fn:          func(x float64) float64 { return 1 / x },
		dFn:         dInv,
		d2Fn:        d2Inv,
	},
	{
		name:        "sqrt",
		x:           []float64{ /*0,*/ 1, 2, 3},
		fnHyperdual: Sqrt,
		fn:          math.Sqrt,
		dFn:         dSqrt,
		d2Fn:        d2Sqrt,
	},

	{
		name: "Fike example fn",
		x:    []float64{1, 2, 3, 4, 5},
		fnHyperdual: func(x Number) Number {
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
		d2Fn: func(x float64) float64 {
			return math.Exp(x) * (130 - 12*math.Cos(2*x) + 30*math.Cos(4*x) + 12*math.Cos(6*x) - 111*math.Sin(2*x) + 48*math.Sin(4*x) + 5*math.Sin(6*x)) /
				(64 * math.Pow(math.Pow(math.Sin(x), 3)+math.Pow(math.Cos(x), 3), 2.5))
		},
	},
}

func TestHyperdual(t *testing.T) {
	const tol = 1e-14
	for _, test := range hyperdualTests {
		for _, x := range test.x {
			fxHyperdual := test.fnHyperdual(Number{Real: x, E1mag: 1, E2mag: 1})
			fx := test.fn(x)
			dFx := test.dFn(x)
			d2Fx := test.d2Fn(x)
			if !same(fxHyperdual.Real, fx, tol) {
				t.Errorf("unexpected %s(%v): got:%v want:%v", test.name, x, fxHyperdual.Real, fx)
			}
			if !same(fxHyperdual.E1mag, dFx, tol) {
				t.Errorf("unexpected %s'(%v) (ϵ₁): got:%v want:%v", test.name, x, fxHyperdual.E1mag, dFx)
			}
			if !same(fxHyperdual.E1mag, fxHyperdual.E2mag, tol) {
				t.Errorf("mismatched ϵ₁ and ϵ₂ for %s(%v): ϵ₁:%v ϵ₂:%v", test.name, x, fxHyperdual.E1mag, fxHyperdual.E2mag)
			}
			if !same(fxHyperdual.E1E2mag, d2Fx, tol) {
				t.Errorf("unexpected %s''(%v): got:%v want:%v", test.name, x, fxHyperdual.E1E2mag, d2Fx)
			}
		}
	}
}

func same(a, b, tol float64) bool {
	return (math.IsNaN(a) && math.IsNaN(b)) || floats.EqualWithinAbsOrRel(a, b, tol, tol)
}
