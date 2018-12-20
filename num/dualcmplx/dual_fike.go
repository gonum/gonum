// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Derived from code by Jeffrey A. Fike at http://adl.stanford.edu/hyperdual/

// The MIT License (MIT)
//
// Copyright (c) 2006 Jeffrey A. Fike
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package dualcmplx

import (
	"math"
	"math/cmplx"
)

// PowReal returns x**p, the base-x exponential of p.
//
// Special cases are (in order):
//	PowReal(NaN+xϵ, ±0) = 1+NaNϵ for any x
//	PowReal(x, ±0) = 1 for any x
//	PowReal(1+xϵ, y) = 1+xyϵ for any y
//	PowReal(x, 1) = x for any x
//	PowReal(NaN+xϵ, y) = NaN+NaNϵ
//	PowReal(x, NaN) = NaN+NaNϵ
//	PowReal(±0, y) = ±Inf for y an odd integer < 0
//	PowReal(±0, -Inf) = +Inf
//	PowReal(±0, +Inf) = +0
//	PowReal(±0, y) = +Inf for finite y < 0 and not an odd integer
//	PowReal(±0, y) = ±0 for y an odd integer > 0
//	PowReal(±0, y) = +0 for finite y > 0 and not an odd integer
//	PowReal(-1, ±Inf) = 1
//	PowReal(x+0ϵ, +Inf) = +Inf+NaNϵ for |x| > 1
//	PowReal(x+yϵ, +Inf) = +Inf for |x| > 1
//	PowReal(x, -Inf) = +0+NaNϵ for |x| > 1
//	PowReal(x, +Inf) = +0+NaNϵ for |x| < 1
//	PowReal(x+0ϵ, -Inf) = +Inf+NaNϵ for |x| < 1
//	PowReal(x, -Inf) = +Inf-Infϵ for |x| < 1
//	PowReal(+Inf, y) = +Inf for y > 0
//	PowReal(+Inf, y) = +0 for y < 0
//	PowReal(-Inf, y) = Pow(-0, -y)
//	PowReal(x, y) = NaN+NaNϵ for finite x < 0 and finite non-integer y
func PowReal(d Number, p complex128) Number {
	deriv := p * cmplx.Pow(d.Real, p-(1+0i))
	return Number{
		Real: cmplx.Pow(d.Real, p),
		Dual: d.Dual * deriv,
	}
}

// Pow return d**r, the base-d exponential of r.
func Pow(d, p Number) Number {
	return Exp(Mul(p, Log(d)))
}

// Sqrt returns the square root of d
//
// Special cases are:
//	Sqrt(+Inf) = +Inf
//	Sqrt(±0) = (±0+Infϵ)
//	Sqrt(x < 0) = NaN
//	Sqrt(NaN) = NaN
func Sqrt(d Number) Number {
	return PowReal(d, 0.5)
}

// Exp returns e**q, the base-e exponential of d.
//
// Special cases are:
//	Exp(+Inf) = +Inf
//	Exp(NaN) = NaN
// Very large values overflow to 0 or +Inf.
// Very small values underflow to 1.
func Exp(d Number) Number {
	fnDeriv := cmplx.Exp(d.Real)
	return Number{
		Real: fnDeriv,
		Dual: fnDeriv * d.Dual,
	}
}

// Log returns the natural logarithm of d.
//
// Special cases are:
//	Log(+Inf) = (+Inf+0ϵ)
//	Log(0) = (-Inf±Infϵ)
//	Log(x < 0) = NaN
//	Log(NaN) = NaN
func Log(d Number) Number {
	switch {
	case d.Real == 0:
		return Number{
			Real: cmplx.Log(d.Real),
			Dual: complex(math.Copysign(math.Inf(1), real(d.Real)), math.NaN()),
		}
	case cmplx.IsInf(d.Real):
		return Number{
			Real: cmplx.Log(d.Real),
			Dual: 0,
		}
	}
	return Number{
		Real: cmplx.Log(d.Real),
		Dual: d.Dual / d.Real,
	}
}

// Sin returns the sine of d.
//
// Special cases are:
//	Sin(±0) = (±0+Nϵ)
//	Sin(±Inf) = NaN
//	Sin(NaN) = NaN
func Sin(d Number) Number {
	if d.Real == 0 {
		return d
	}
	fn := cmplx.Sin(d.Real)
	deriv := cmplx.Cos(d.Real)
	return Number{
		Real: fn,
		Dual: deriv * d.Dual,
	}
}

// Cos returns the cosine of d.
//
// Special cases are:
//	Cos(±Inf) = NaN
//	Cos(NaN) = NaN
func Cos(d Number) Number {
	fn := cmplx.Cos(d.Real)
	deriv := -cmplx.Sin(d.Real)
	return Number{
		Real: fn,
		Dual: deriv * d.Dual,
	}
}

// Tan returns the tangent of d.
//
// Special cases are:
//	Tan(±0) = (±0+Nϵ)
//	Tan(±Inf) = NaN
//	Tan(NaN) = NaN
func Tan(d Number) Number {
	if d.Real == 0 {
		return d
	}
	fn := cmplx.Tan(d.Real)
	deriv := 1 + fn*fn
	return Number{
		Real: fn,
		Dual: deriv * d.Dual,
	}
}

// Asin returns the inverse sine of d.
//
// Special cases are:
//	Asin(±0) = (±0+Nϵ)
//	Asin(±1) = (±Inf+Infϵ)
//	Asin(x) = NaN if x < -1 or x > 1
func Asin(d Number) Number {
	if d.Real == 0 {
		return d
	} else if m := cmplx.Abs(d.Real); m >= 1 {
		if m == 1 {
			return Number{
				Real: cmplx.Asin(d.Real),
				Dual: cmplx.Inf(),
			}
		}
		return Number{
			Real: cmplx.NaN(),
			Dual: cmplx.NaN(),
		}
	}
	fn := cmplx.Asin(d.Real)
	deriv := 1 / cmplx.Sqrt(1-d.Real*d.Real)
	return Number{
		Real: fn,
		Dual: deriv * d.Dual,
	}
}

// Acos returns the inverse cosine of d.
//
// Special cases are:
//	Acos(-1) = (Pi-Infϵ)
//	Acos(1) = (0-Infϵ)
//	Acos(x) = NaN if x < -1 or x > 1
func Acos(d Number) Number {
	if m := cmplx.Abs(d.Real); m >= 1 {
		if m == 1 {
			return Number{
				Real: cmplx.Acos(d.Real),
				Dual: cmplx.Inf(),
			}
		}
		return Number{
			Real: cmplx.NaN(),
			Dual: cmplx.NaN(),
		}
	}
	fn := cmplx.Acos(d.Real)
	deriv := -1 / cmplx.Sqrt(1-d.Real*d.Real)
	return Number{
		Real: fn,
		Dual: deriv * d.Dual,
	}
}

// Atan returns the inverse tangent of d.
//
// Special cases are:
//	Atan(±0) = (±0+Nϵ)
//	Atan(±Inf) = (±Pi/2+0ϵ)
func Atan(d Number) Number {
	if d.Real == 0 {
		return d
	}
	fn := cmplx.Atan(d.Real)
	deriv := 1 / (1 + d.Real*d.Real)
	return Number{
		Real: fn,
		Dual: deriv * d.Dual,
	}
}
