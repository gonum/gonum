// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dualquat

import (
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/num/quat"
)

// First derivatives:

func dSin(x quat.Number) quat.Number  { return quat.Cos(x) }
func dCos(x quat.Number) quat.Number  { return quat.Scale(-1, quat.Sin(x)) }
func dTan(x quat.Number) quat.Number  { return quat.Mul(sec(x), sec(x)) }
func dAsin(x quat.Number) quat.Number { return quat.Inv(quat.Sqrt(subRealQuat(1, quat.Mul(x, x)))) }
func dAcos(x quat.Number) quat.Number {
	return quat.Scale(-1, quat.Inv(quat.Sqrt(subRealQuat(1, quat.Mul(x, x)))))
}
func dAtan(x quat.Number) quat.Number { return quat.Inv(addRealQuat(1, quat.Mul(x, x))) }

func dSinh(x quat.Number) quat.Number  { return quat.Cosh(x) }
func dCosh(x quat.Number) quat.Number  { return quat.Sinh(x) }
func dTanh(x quat.Number) quat.Number  { return quat.Mul(sech(x), sech(x)) }
func dAsinh(x quat.Number) quat.Number { return quat.Inv(quat.Sqrt(addQuatReal(quat.Mul(x, x), 1))) }
func dAcosh(x quat.Number) quat.Number {
	return quat.Inv(quat.Mul(quat.Sqrt(subQuatReal(x, 1)), quat.Sqrt(addQuatReal(x, 1))))
}
func dAtanh(x quat.Number) quat.Number { return quat.Inv(subRealQuat(1, quat.Mul(x, x))) }

func dExp(x quat.Number) quat.Number { return quat.Exp(x) }
func dLog(x quat.Number) quat.Number {
	switch {
	case x == zeroQuat:
		return quat.Inf()
	case quat.IsInf(x):
		return zeroQuat
	}
	return quat.Inv(x)
}
func dPow(x, y quat.Number) quat.Number { return quat.Mul(y, quat.Pow(x, subQuatReal(y, 1))) }
func dSqrt(x quat.Number) quat.Number   { return quat.Scale(0.5, quat.Inv(quat.Sqrt(x))) }
func dInv(x quat.Number) quat.Number    { return quat.Scale(-1, quat.Inv(quat.Mul(x, x))) }

// Helpers:

func sec(x quat.Number) quat.Number  { return quat.Inv(quat.Cos(x)) }
func sech(x quat.Number) quat.Number { return quat.Inv(quat.Cosh(x)) }

var (
	negZeroQuat = quat.Scale(-1, zeroQuat)
	one         = quat.Number{1, 1, 1, 1}
	negOne      = quat.Scale(-1, one)
	half        = quat.Scale(0.5, one)
	negHalf     = quat.Scale(-1, half)
	two         = quat.Scale(2, one)
	negTwo      = quat.Scale(-1, two)
	three       = quat.Scale(3, one)
	negThree    = quat.Scale(-1, three)
)

var dualTests = []struct {
	name   string
	x      []quat.Number
	fnDual func(x Number) Number
	fn     func(x quat.Number) quat.Number
	dFn    func(x quat.Number) quat.Number
}{
	{
		name:   "sin",
		x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
		fnDual: Sin,
		fn:     quat.Sin,
		dFn:    dSin,
	},
	{
		name:   "cos",
		x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
		fnDual: Cos,
		fn:     quat.Cos,
		dFn:    dCos,
	},
	{
		name:   "tan",
		x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
		fnDual: Tan,
		fn:     quat.Tan,
		dFn:    dTan,
	},
	{
		name:   "sinh",
		x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
		fnDual: Sinh,
		fn:     quat.Sinh,
		dFn:    dSinh,
	},
	{
		name:   "cosh",
		x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
		fnDual: Cosh,
		fn:     quat.Cosh,
		dFn:    dCosh,
	},
	// {//fail
	// 	name:   "tanh",
	// 	x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
	// 	fnDual: Tanh,
	// 	fn:     quat.Tanh,
	// 	dFn:    dTanh,
	// },

	// {//fail
	// 	name:   "asin",
	// 	x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
	// 	fnDual: Asin,
	// 	fn:     quat.Asin,
	// 	dFn:    dAsin,
	// },
	// {//fail
	// 	name:   "acos",
	// 	x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
	// 	fnDual: Acos,
	// 	fn:     quat.Acos,
	// 	dFn:    dAcos,
	// },
	{
		name:   "atan",
		x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
		fnDual: Atan,
		fn:     quat.Atan,
		dFn:    dAtan,
	},
	{
		name:   "asinh",
		x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
		fnDual: Asinh,
		fn:     quat.Asinh,
		dFn:    dAsinh,
	},
	// { //fail
	// 	name:   "acosh",
	// 	x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
	// 	fnDual: Acosh,
	// 	fn:     quat.Acosh,
	// 	dFn:    dAcosh,
	// },
	{
		name:   "atanh",
		x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
		fnDual: Atanh,
		fn:     quat.Atanh,
		dFn:    dAtanh,
	},

	{
		name:   "exp",
		x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
		fnDual: Exp,
		fn:     quat.Exp,
		dFn:    dExp,
	},
	{
		name:   "log",
		x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
		fnDual: Log,
		fn:     quat.Log,
		dFn:    dLog,
	},
	{
		name:   "inv",
		x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
		fnDual: Inv,
		fn:     quat.Inv,
		dFn:    dInv,
	},
	{
		name:   "sqrt",
		x:      []quat.Number{quat.NaN(), quat.Inf(), negThree, negTwo, negOne, negHalf, negZeroQuat, zeroQuat, half, one, two, three},
		fnDual: Sqrt,
		fn:     quat.Sqrt,
		dFn:    dSqrt,
	},
}

func TestDual(t *testing.T) {
	const tol = 1e-15
	for _, test := range dualTests {
		for _, x := range test.x {
			fxDual := test.fnDual(Number{Real: x, Dual: quat.Number{Real: 1}})
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

func sameDual(a, b Number, tol float64) bool {
	return same(a.Real, b.Real, tol) && same(a.Emag, b.Emag, tol)
}
*/
func same(a, b quat.Number, tol float64) bool {
	return (quat.IsNaN(a) && quat.IsNaN(b)) || (quat.IsInf(a) && quat.IsInf(b)) || equalApprox(a, b, tol)
}

func equalApprox(a, b quat.Number, tol float64) bool {
	return floats.EqualWithinAbsOrRel(a.Real, b.Real, tol, tol) &&
		floats.EqualWithinAbsOrRel(a.Imag, b.Imag, tol, tol) &&
		floats.EqualWithinAbsOrRel(a.Jmag, b.Jmag, tol, tol) &&
		floats.EqualWithinAbsOrRel(a.Kmag, b.Kmag, tol, tol)
}
