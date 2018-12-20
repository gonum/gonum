// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dualcmplx

import (
	"math"
	"math/cmplx"
	"testing"

	"gonum.org/v1/gonum/floats"
)

// FIXME(kortschka): See golang/go#29320.

func sinh(x complex128) complex128 {
	if cmplx.IsInf(x) {
		return complex(math.Inf(1), math.NaN())
	}
	return cmplx.Sinh(x)
}
func cosh(x complex128) complex128 {
	if cmplx.IsInf(x) {
		return complex(math.Inf(1), math.NaN())
	}
	return cmplx.Cosh(x)
}
func tanh(x complex128) complex128 {
	if cmplx.IsInf(x) {
		return 1
	}
	return cmplx.Cosh(x)
}
func sqrt(x complex128) complex128 {
	switch {
	case math.IsInf(imag(x), 1):
		return cmplx.Inf()
	case math.IsNaN(imag(x)):
		return cmplx.NaN()
	case math.IsInf(real(x), -1):
		if imag(x) >= 0 && !math.IsInf(imag(x), 1) {
			return complex(0, math.NaN())
		}
		if math.IsNaN(imag(x)) {
			return complex(math.NaN(), math.Inf(1))
		}
	case math.IsInf(real(x), 1):
		if imag(x) >= 0 && !math.IsInf(imag(x), 1) {
			return complex(math.Inf(1), 0)
		}
		if math.IsNaN(imag(x)) {
			return complex(math.Inf(1), math.NaN())
		}
	case math.IsInf(real(x), -1):
		return complex(0, math.Inf(1))
	case math.IsNaN(real(x)):
		if math.IsNaN(imag(x)) || math.IsInf(imag(x), 0) {
			return cmplx.NaN()
		}
	}
	return cmplx.Sqrt(x)
}

// First derivatives:

func dSin(x complex128) complex128  { return cmplx.Cos(x) }
func dCos(x complex128) complex128  { return -cmplx.Sin(x) }
func dTan(x complex128) complex128  { return sec(x) * sec(x) }
func dAsin(x complex128) complex128 { return 1 / cmplx.Sqrt(1-x*x) }
func dAcos(x complex128) complex128 { return -1 / cmplx.Sqrt(1-x*x) }
func dAtan(x complex128) complex128 { return 1 / (1 + x*x) }

func dSinh(x complex128) complex128 { return cosh(x) }
func dCosh(x complex128) complex128 { return sinh(x) }
func dTanh(x complex128) complex128 {
	if cmplx.IsInf(x) {
		return 0
	}
	return sech(x) * sech(x)
}
func dAsinh(x complex128) complex128 { return 1 / cmplx.Sqrt(x*x+1) }
func dAcosh(x complex128) complex128 { return 1 / cmplx.Sqrt((x-1)*(x+1)) }
func dAtanh(x complex128) complex128 { return 1 / (1 - x*x) }

func dExp(x complex128) complex128    { return cmplx.Exp(x) }
func dLog(x complex128) complex128    { return 1 / x }
func dPow(x, y complex128) complex128 { return y * cmplx.Pow(x, y-1) }
func dSqrt(x complex128) complex128 {
	if x == 0 {
		return cmplx.NaN()
	}
	return 0.5 / cmplx.Sqrt(x)
}
func dInv(x complex128) complex128 { return -1 / (x * x) }

// Helpers:

func sec(x complex128) complex128  { return 1 / cmplx.Cos(x) }
func sech(x complex128) complex128 { return 1 / cmplx.Cosh(x) }

var (
	zeroCmplx    = 0 + 0i
	negZeroCmplx = -1 * zeroCmplx
	one          = 1 + 1i
	negOne       = -1 - 1i
	half         = one / 2
	negHalf      = negOne / 2
	two          = 2 + 2i
	negTwo       = -2 - 2i
	three        = 3 + 3i
	negThree     = -3 + 3i
)

var dualTests = []struct {
	name   string
	x      []complex128
	fnDual func(x Number) Number
	fn     func(x complex128) complex128
	dFn    func(x complex128) complex128
}{
	{
		name:   "sin",
		x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
		fnDual: Sin,
		fn:     cmplx.Sin,
		dFn:    dSin,
	},
	{
		name:   "cos",
		x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
		fnDual: Cos,
		fn:     cmplx.Cos,
		dFn:    dCos,
	},
	{
		name:   "tan",
		x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
		fnDual: Tan,
		fn:     cmplx.Tan,
		dFn:    dTan,
	},
	{
		name:   "sinh",
		x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
		fnDual: Sinh,
		fn:     sinh,
		dFn:    dSinh,
	},
	{
		name:   "cosh",
		x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
		fnDual: Cosh,
		fn:     cosh,
		dFn:    dCosh,
	},
	// {//fail
	// 	name:   "tanh",
	// 	x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
	// 	fnDual: Tanh,
	// 	fn:     tanh,
	// 	dFn:    dTanh,
	// },

	// {//fail
	// 	name:   "asin",
	// 	x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
	// 	fnDual: Asin,
	// 	fn:     cmplx.Asin,
	// 	dFn:    dAsin,
	// },
	// {//fail
	// 	name:   "acos",
	// 	x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
	// 	fnDual: Acos,
	// 	fn:     cmplx.Acos,
	// 	dFn:    dAcos,
	// },
	{
		name:   "atan",
		x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
		fnDual: Atan,
		fn:     cmplx.Atan,
		dFn:    dAtan,
	},
	{
		name:   "asinh",
		x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
		fnDual: Asinh,
		fn:     cmplx.Asinh,
		dFn:    dAsinh,
	},
	// {//fail
	// 	name:   "acosh",
	// 	x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
	// 	fnDual: Acosh,
	// 	fn:     cmplx.Acosh,
	// 	dFn:    dAcosh,
	// },
	{
		name:   "atanh",
		x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
		fnDual: Atanh,
		fn:     cmplx.Atanh,
		dFn:    dAtanh,
	},

	{
		name:   "exp",
		x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
		fnDual: Exp,
		fn:     cmplx.Exp,
		dFn:    dExp,
	},
	{
		name:   "log",
		x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
		fnDual: Log,
		fn:     cmplx.Log,
		dFn:    dLog,
	},
	{
		name:   "inv",
		x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
		fnDual: Inv,
		fn:     func(x complex128) complex128 { return 1 / x },
		dFn:    dInv,
	},
	{
		name:   "sqrt",
		x:      []complex128{cmplx.NaN(), cmplx.Inf(), negThree, negTwo, negOne, negHalf, negZeroCmplx, zeroCmplx, half, one, two, three},
		fnDual: Sqrt,
		fn:     sqrt,
		dFn:    dSqrt,
	},
}

func TestDual(t *testing.T) {
	const tol = 1e-15
	for _, test := range dualTests {
		for _, x := range test.x {
			fxDual := test.fnDual(Number{Real: x, Dual: 1})
			fx := test.fn(x)
			dFx := test.dFn(x)
			if !same(fxDual.Real, fx, tol) {
				t.Errorf("unexpected %s(%v): got:%v want:%v", test.name, x, fxDual.Real, fx)
			}
			if !same(fxDual.Dual, dFx, tol) {
				t.Errorf("unexpected %s'(%v): got:%v want:%v", test.name, x, fxDual.Dual, dFx)
			}
		}
	}
}

/*
var powRealTests = []struct {
	d    Number
	p    float64
	want Number
}{
	// PowReal(NaN+xϵ, ±0) = 1+NaNϵ for any x
	{d: Number{Real: math.NaN(), Emag: 0}, p: 0, want: Number{Real: 1, Emag: math.NaN()}},
	{d: Number{Real: math.NaN(), Emag: 0}, p: negZero, want: Number{Real: 1, Emag: math.NaN()}},
	{d: Number{Real: math.NaN(), Emag: 1}, p: 0, want: Number{Real: 1, Emag: math.NaN()}},
	{d: Number{Real: math.NaN(), Emag: 2}, p: negZero, want: Number{Real: 1, Emag: math.NaN()}},
	{d: Number{Real: math.NaN(), Emag: 3}, p: 0, want: Number{Real: 1, Emag: math.NaN()}},
	{d: Number{Real: math.NaN(), Emag: 1}, p: negZero, want: Number{Real: 1, Emag: math.NaN()}},
	{d: Number{Real: math.NaN(), Emag: 2}, p: 0, want: Number{Real: 1, Emag: math.NaN()}},
	{d: Number{Real: math.NaN(), Emag: 3}, p: negZero, want: Number{Real: 1, Emag: math.NaN()}},

	// PowReal(x, ±0) = 1 for any x
	{d: Number{Real: 0, Emag: 0}, p: 0, want: Number{Real: 1, Emag: 0}},
	{d: Number{Real: negZero, Emag: 0}, p: negZero, want: Number{Real: 1, Emag: 0}},
	{d: Number{Real: math.Inf(1), Emag: 0}, p: 0, want: Number{Real: 1, Emag: 0}},
	{d: Number{Real: math.Inf(-1), Emag: 0}, p: negZero, want: Number{Real: 1, Emag: 0}},
	{d: Number{Real: 0, Emag: 1}, p: 0, want: Number{Real: 1, Emag: 0}},
	{d: Number{Real: negZero, Emag: 1}, p: negZero, want: Number{Real: 1, Emag: 0}},
	{d: Number{Real: math.Inf(1), Emag: 1}, p: 0, want: Number{Real: 1, Emag: 0}},
	{d: Number{Real: math.Inf(-1), Emag: 1}, p: negZero, want: Number{Real: 1, Emag: 0}},

	// PowReal(1+xϵ, y) = (1+xyϵ) for any y
	{d: Number{Real: 1, Emag: 0}, p: 0, want: Number{Real: 1, Emag: 0}},
	{d: Number{Real: 1, Emag: 0}, p: 1, want: Number{Real: 1, Emag: 0}},
	{d: Number{Real: 1, Emag: 0}, p: 2, want: Number{Real: 1, Emag: 0}},
	{d: Number{Real: 1, Emag: 0}, p: 3, want: Number{Real: 1, Emag: 0}},
	{d: Number{Real: 1, Emag: 1}, p: 0, want: Number{Real: 1, Emag: 0}},
	{d: Number{Real: 1, Emag: 1}, p: 1, want: Number{Real: 1, Emag: 1}},
	{d: Number{Real: 1, Emag: 1}, p: 2, want: Number{Real: 1, Emag: 2}},
	{d: Number{Real: 1, Emag: 1}, p: 3, want: Number{Real: 1, Emag: 3}},
	{d: Number{Real: 1, Emag: 2}, p: 0, want: Number{Real: 1, Emag: 0}},
	{d: Number{Real: 1, Emag: 2}, p: 1, want: Number{Real: 1, Emag: 2}},
	{d: Number{Real: 1, Emag: 2}, p: 2, want: Number{Real: 1, Emag: 4}},
	{d: Number{Real: 1, Emag: 2}, p: 3, want: Number{Real: 1, Emag: 6}},

	// PowReal(x, 1) = x for any x
	{d: Number{Real: 0, Emag: 0}, p: 1, want: Number{Real: 0, Emag: 0}},
	{d: Number{Real: negZero, Emag: 0}, p: 1, want: Number{Real: negZero, Emag: 0}},
	{d: Number{Real: 0, Emag: 1}, p: 1, want: Number{Real: 0, Emag: 1}},
	{d: Number{Real: negZero, Emag: 1}, p: 1, want: Number{Real: negZero, Emag: 1}},
	{d: Number{Real: math.NaN(), Emag: 0}, p: 1, want: Number{Real: math.NaN(), Emag: 0}},
	{d: Number{Real: math.NaN(), Emag: 1}, p: 1, want: Number{Real: math.NaN(), Emag: 1}},
	{d: Number{Real: math.NaN(), Emag: 2}, p: 1, want: Number{Real: math.NaN(), Emag: 2}},

	// PowReal(NaN+xϵ, y) = NaN+NaNϵ
	{d: Number{Real: math.NaN(), Emag: 0}, p: 2, want: Number{Real: math.NaN(), Emag: math.NaN()}},
	{d: Number{Real: math.NaN(), Emag: 0}, p: 3, want: Number{Real: math.NaN(), Emag: math.NaN()}},
	{d: Number{Real: math.NaN(), Emag: 1}, p: 2, want: Number{Real: math.NaN(), Emag: math.NaN()}},
	{d: Number{Real: math.NaN(), Emag: 1}, p: 3, want: Number{Real: math.NaN(), Emag: math.NaN()}},
	{d: Number{Real: math.NaN(), Emag: 2}, p: 2, want: Number{Real: math.NaN(), Emag: math.NaN()}},
	{d: Number{Real: math.NaN(), Emag: 2}, p: 3, want: Number{Real: math.NaN(), Emag: math.NaN()}},

	// PowReal(x, NaN) = NaN+NaNϵ
	{d: Number{Real: 0, Emag: 0}, p: math.NaN(), want: Number{Real: math.NaN(), Emag: math.NaN()}},
	{d: Number{Real: 2, Emag: 0}, p: math.NaN(), want: Number{Real: math.NaN(), Emag: math.NaN()}},
	{d: Number{Real: 3, Emag: 0}, p: math.NaN(), want: Number{Real: math.NaN(), Emag: math.NaN()}},
	{d: Number{Real: 0, Emag: 1}, p: math.NaN(), want: Number{Real: math.NaN(), Emag: math.NaN()}},
	{d: Number{Real: 2, Emag: 1}, p: math.NaN(), want: Number{Real: math.NaN(), Emag: math.NaN()}},
	{d: Number{Real: 3, Emag: 1}, p: math.NaN(), want: Number{Real: math.NaN(), Emag: math.NaN()}},
	{d: Number{Real: 0, Emag: 2}, p: math.NaN(), want: Number{Real: math.NaN(), Emag: math.NaN()}},
	{d: Number{Real: 2, Emag: 2}, p: math.NaN(), want: Number{Real: math.NaN(), Emag: math.NaN()}},
	{d: Number{Real: 3, Emag: 2}, p: math.NaN(), want: Number{Real: math.NaN(), Emag: math.NaN()}},

	// Handled by math.Pow tests:
	//
	// Pow(±0, y) = ±Inf for y an odd integer < 0
	// Pow(±0, -Inf) = +Inf
	// Pow(±0, +Inf) = +0
	// Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
	// Pow(±0, y) = ±0 for y an odd integer > 0
	// Pow(±0, y) = +0 for finite y > 0 and not an odd integer
	// Pow(-1, ±Inf) = 1

	// PowReal(x+0ϵ, +Inf) = +Inf+NaNϵ for |x| > 1
	{d: Number{Real: 2, Emag: 0}, p: math.Inf(1), want: Number{Real: math.Inf(1), Emag: math.NaN()}},
	{d: Number{Real: 3, Emag: 0}, p: math.Inf(1), want: Number{Real: math.Inf(1), Emag: math.NaN()}},

	// PowReal(x+yϵ, +Inf) = +Inf for |x| > 1
	{d: Number{Real: 2, Emag: 1}, p: math.Inf(1), want: Number{Real: math.Inf(1), Emag: math.Inf(1)}},
	{d: Number{Real: 3, Emag: 1}, p: math.Inf(1), want: Number{Real: math.Inf(1), Emag: math.Inf(1)}},
	{d: Number{Real: 2, Emag: 2}, p: math.Inf(1), want: Number{Real: math.Inf(1), Emag: math.Inf(1)}},
	{d: Number{Real: 3, Emag: 2}, p: math.Inf(1), want: Number{Real: math.Inf(1), Emag: math.Inf(1)}},

	// PowReal(x, -Inf) = +0+NaNϵ for |x| > 1
	{d: Number{Real: 2, Emag: 0}, p: math.Inf(-1), want: Number{Real: 0, Emag: math.NaN()}},
	{d: Number{Real: 3, Emag: 0}, p: math.Inf(-1), want: Number{Real: 0, Emag: math.NaN()}},
	{d: Number{Real: 2, Emag: 1}, p: math.Inf(-1), want: Number{Real: 0, Emag: math.NaN()}},
	{d: Number{Real: 3, Emag: 1}, p: math.Inf(-1), want: Number{Real: 0, Emag: math.NaN()}},
	{d: Number{Real: 2, Emag: 2}, p: math.Inf(-1), want: Number{Real: 0, Emag: math.NaN()}},
	{d: Number{Real: 3, Emag: 2}, p: math.Inf(-1), want: Number{Real: 0, Emag: math.NaN()}},

	// PowReal(x+yϵ, +Inf) = +0+NaNϵ for |x| < 1
	{d: Number{Real: 0.1, Emag: 0}, p: math.Inf(1), want: Number{Real: 0, Emag: math.NaN()}},
	{d: Number{Real: 0.1, Emag: 0.1}, p: math.Inf(1), want: Number{Real: 0, Emag: math.NaN()}},
	{d: Number{Real: 0.2, Emag: 0.2}, p: math.Inf(1), want: Number{Real: 0, Emag: math.NaN()}},
	{d: Number{Real: 0.5, Emag: 0.5}, p: math.Inf(1), want: Number{Real: 0, Emag: math.NaN()}},

	// PowReal(x+0ϵ, -Inf) = +Inf+NaNϵ for |x| < 1
	{d: Number{Real: 0.1, Emag: 0}, p: math.Inf(-1), want: Number{Real: math.Inf(1), Emag: math.NaN()}},
	{d: Number{Real: 0.2, Emag: 0}, p: math.Inf(-1), want: Number{Real: math.Inf(1), Emag: math.NaN()}},

	// PowReal(x, -Inf) = +Inf-Infϵ for |x| < 1
	{d: Number{Real: 0.1, Emag: 0.1}, p: math.Inf(-1), want: Number{Real: math.Inf(1), Emag: math.Inf(-1)}},
	{d: Number{Real: 0.2, Emag: 0.1}, p: math.Inf(-1), want: Number{Real: math.Inf(1), Emag: math.Inf(-1)}},
	{d: Number{Real: 0.1, Emag: 0.2}, p: math.Inf(-1), want: Number{Real: math.Inf(1), Emag: math.Inf(-1)}},
	{d: Number{Real: 0.2, Emag: 0.2}, p: math.Inf(-1), want: Number{Real: math.Inf(1), Emag: math.Inf(-1)}},
	{d: Number{Real: 0.1, Emag: 1}, p: math.Inf(-1), want: Number{Real: math.Inf(1), Emag: math.Inf(-1)}},
	{d: Number{Real: 0.2, Emag: 1}, p: math.Inf(-1), want: Number{Real: math.Inf(1), Emag: math.Inf(-1)}},
	{d: Number{Real: 0.1, Emag: 2}, p: math.Inf(-1), want: Number{Real: math.Inf(1), Emag: math.Inf(-1)}},
	{d: Number{Real: 0.2, Emag: 2}, p: math.Inf(-1), want: Number{Real: math.Inf(1), Emag: math.Inf(-1)}},

	// Handled by math.Pow tests:
	//
	// Pow(+Inf, y) = +Inf for y > 0
	// Pow(+Inf, y) = +0 for y < 0
	// Pow(-Inf, y) = Pow(-0, -y)

	// PowReal(x, y) = NaN+NaNϵ for finite x < 0 and finite non-integer y
	{d: Number{Real: -1, Emag: -1}, p: 0.5, want: Number{Real: math.NaN(), Emag: math.NaN()}},
	{d: Number{Real: -1, Emag: 2}, p: 0.5, want: Number{Real: math.NaN(), Emag: math.NaN()}},
}

func TestPowReal(t *testing.T) {
	const tol = 1e-15
	for _, test := range powRealTests {
		got := PowReal(test.d, test.p)
		if !sameDual(got, test.want, tol) {
			t.Errorf("unexpected PowReal(%v, %v): got:%v want:%v", test.d, test.p, got, test.want)
		}
	}
}
*/

func sameDual(a, b Number, tol float64) bool {
	return same(a.Real, b.Real, tol) && same(a.Dual, b.Dual, tol)
}

func same(a, b complex128, tol float64) bool {
	return ((math.IsNaN(real(a)) && (math.IsNaN(real(b)))) || floats.EqualWithinAbsOrRel(real(a), real(b), tol, tol)) &&
		((math.IsNaN(imag(a)) && (math.IsNaN(imag(b)))) || floats.EqualWithinAbsOrRel(imag(a), imag(b), tol, tol))
}

func equalApprox(a, b complex128, tol float64) bool {
	return floats.EqualWithinAbsOrRel(real(a), real(b), tol, tol) &&
		floats.EqualWithinAbsOrRel(imag(a), imag(b), tol, tol)
}
