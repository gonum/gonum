// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dualcmplx

import (
	"math"
	"math/cmplx"
)

// Sinh returns the hyperbolic sine of d.
//
// Special cases are:
//	Sinh(±0) = (±0+Nϵ)
//	Sinh(±Inf) = ±Inf
//	Sinh(NaN) = NaN
func Sinh(d Number) Number {
	if d.Real == 0 {
		return d
	}
	if cmplx.IsInf(d.Real) {
		// FIXME(kortschka): See golang/go#29320.
		return Number{
			Real: complex(math.Inf(1), math.NaN()),
			Dual: complex(math.Inf(1), math.NaN()),
		}
	}
	fn := cmplx.Sinh(d.Real)
	deriv := cmplx.Cosh(d.Real)
	return Number{
		Real: fn,
		Dual: deriv * d.Dual,
	}
}

// Cosh returns the hyperbolic cosine of d.
//
// Special cases are:
//	Cosh(±0) = 1
//	Cosh(±Inf) = +Inf
//	Cosh(NaN) = NaN
func Cosh(d Number) Number {
	if cmplx.IsInf(d.Real) {
		// FIXME(kortschka): See golang/go#29320.
		return Number{
			Real: complex(math.Inf(1), math.NaN()),
			Dual: complex(math.Inf(1), math.NaN()),
		}
	}
	fn := cmplx.Cosh(d.Real)
	deriv := cmplx.Sinh(d.Real)
	return Number{
		Real: fn,
		Dual: deriv * d.Dual,
	}
}

// Tanh returns the hyperbolic tangent of d.
//
// Special cases are:
//	Tanh(±0) = (±0+Nϵ)
//	Tanh(±Inf) = (±1+0ϵ)
//	Tanh(NaN) = NaN
func Tanh(d Number) Number {
	switch {
	case d.Real == 0:
		return d
	case cmplx.IsInf(d.Real):
		return Number{
			Real: 1,
			Dual: 0,
		}
	}
	fn := cmplx.Tanh(d.Real)
	deriv := 1 - fn*fn
	return Number{
		Real: fn,
		Dual: deriv * d.Dual,
	}
}

// Asinh returns the inverse hyperbolic sine of d.
//
// Special cases are:
//	Asinh(±0) = (±0+Nϵ)
//	Asinh(±Inf) = ±Inf
//	Asinh(NaN) = NaN
func Asinh(d Number) Number {
	if d.Real == 0 {
		return d
	}
	fn := cmplx.Asinh(d.Real)
	deriv := 1 / cmplx.Sqrt(d.Real*d.Real+1)
	return Number{
		Real: fn,
		Dual: deriv * d.Dual,
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
	fn := cmplx.Acosh(d.Real)
	deriv := 1 / cmplx.Sqrt(d.Real*d.Real-1)
	return Number{
		Real: fn,
		Dual: deriv * d.Dual,
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
	if d.Real == 0 {
		return d
	}
	fn := cmplx.Atanh(d.Real)
	deriv := 1 / (1 - d.Real*d.Real)
	return Number{
		Real: fn,
		Dual: deriv * d.Dual,
	}
}
