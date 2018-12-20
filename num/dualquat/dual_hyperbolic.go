// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dualquat

import "gonum.org/v1/gonum/num/quat"

// Sinh returns the hyperbolic sine of d.
//
// Special cases are:
//	Sinh(±0) = (±0+Nϵ)
//	Sinh(±Inf) = ±Inf
//	Sinh(NaN) = NaN
func Sinh(d Number) Number {
	if d.Real == zeroQuat {
		return d
	}
	fn := quat.Sinh(d.Real)
	deriv := quat.Cosh(d.Real)
	return Number{
		Real: fn,
		Dual: quat.Mul(deriv, d.Dual),
	}
}

// Cosh returns the hyperbolic cosine of d.
//
// Special cases are:
//	Cosh(±0) = 1
//	Cosh(±Inf) = +Inf
//	Cosh(NaN) = NaN
func Cosh(d Number) Number {
	fn := quat.Cosh(d.Real)
	deriv := quat.Sinh(d.Real)
	return Number{
		Real: fn,
		Dual: quat.Mul(deriv, d.Dual),
	}
}

// Tanh returns the hyperbolic tangent of d.
//
// Special cases are:
//	Tanh(±0) = (±0+Nϵ)
//	Tanh(±Inf) = (±1+0ϵ)
//	Tanh(NaN) = NaN
func Tanh(d Number) Number {
	if d.Real == zeroQuat {
		return d
	}
	fn := quat.Tanh(d.Real)
	deriv := subRealQuat(1, quat.Mul(fn, fn))
	return Number{
		Real: fn,
		Dual: quat.Mul(deriv, d.Dual),
	}
}

// Asinh returns the inverse hyperbolic sine of d.
//
// Special cases are:
//	Asinh(±0) = (±0+Nϵ)
//	Asinh(±Inf) = ±Inf
//	Asinh(NaN) = NaN
func Asinh(d Number) Number {
	if d.Real == zeroQuat {
		return d
	}
	fn := quat.Asinh(d.Real)
	deriv := quat.Inv(quat.Sqrt(addQuatReal(quat.Mul(d.Real, d.Real), 1)))
	return Number{
		Real: fn,
		Dual: quat.Mul(deriv, d.Dual),
	}
}

// Acosh returns the inverse hyperbolic cosine of d.
//
// Special cases are:
//	Acosh(+Inf) = +Inf
//	Acosh(1) = (0+Infϵ)
//	Acosh(x) = NaN if x < 1
//	Acosh(NaN) = NaN
func Acosh(d Number) Number {
	fn := quat.Acosh(d.Real)
	deriv := quat.Inv(quat.Sqrt(subQuatReal(quat.Mul(d.Real, d.Real), 1)))
	return Number{
		Real: fn,
		Dual: quat.Mul(deriv, d.Dual),
	}
}

// Atanh returns the inverse hyperbolic tangent of d.
//
// Special cases are:
//	Atanh(1) = +Inf
//	Atanh(±0) = (±0+Nϵ)
//	Atanh(-1) = -Inf
//	Atanh(x) = NaN if x < -1 or x > 1
//	Atanh(NaN) = NaN
func Atanh(d Number) Number {
	if d.Real == zeroQuat {
		return d
	}
	fn := quat.Atanh(d.Real)
	deriv := quat.Inv(subRealQuat(1, quat.Mul(d.Real, d.Real)))
	return Number{
		Real: fn,
		Dual: quat.Mul(deriv, d.Dual),
	}
}
