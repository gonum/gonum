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

package dual

import "math"

func PowReal(d Number, p float64) Number {
	const tol = 1e-15

	r := d.Real
	if math.Abs(r) < tol {
		if r >= 0 {
			r = tol
		}
		if r < 0 {
			r = -tol
		}
	}
	deriv := p * math.Pow(r, p-1)
	return Number{
		Real: math.Pow(d.Real, p),
		Emag: d.Emag * deriv,
	}
}

// Pow return d**r, the base-d exponential of r.
func Pow(d, p Number) Number {
	return Exp(Mul(p, Log(d)))
}

// Exp returns e**d, the base-e exponential of d.
func Exp(d Number) Number {
	fnDeriv := math.Exp(d.Real)
	return Number{
		Real: fnDeriv,
		Emag: fnDeriv * d.Emag,
	}
}

// Log returns the natural logarithm of d.
func Log(d Number) Number {
	return Number{
		Real: math.Log(d.Real),
		Emag: d.Emag / d.Real,
	}
}

// Sin returns the sine of d.
func Sin(d Number) Number {
	fn := math.Sin(d.Real)
	deriv := math.Cos(d.Real)
	return Number{
		Real: fn,
		Emag: deriv * d.Emag,
	}
}

// Cos returns the cosine of d.
func Cos(d Number) Number {
	fn := math.Cos(d.Real)
	deriv := -math.Sin(d.Real)
	return Number{
		Real: fn,
		Emag: deriv * d.Emag,
	}
}

// Tan returns the tangent of d.
func Tan(d Number) Number {
	fn := math.Tan(d.Real)
	deriv := 1 + fn*fn
	return Number{
		Real: fn,
		Emag: deriv * d.Emag,
	}
}

// Asin returns the inverse sine of d.
func Asin(d Number) Number {
	fn := math.Asin(d.Real)
	deriv := 1 / math.Sqrt(1-d.Real*d.Real)
	return Number{
		Real: fn,
		Emag: deriv * d.Emag,
	}
}

// Acos returns the inverse cosine of d.
func Acos(d Number) Number {
	fn := math.Acos(d.Real)
	deriv := -1 / math.Sqrt(1-d.Real*d.Real)
	return Number{
		Real: fn,
		Emag: deriv * d.Emag,
	}
}

// Atan returns the inverse tangent of d.
func Atan(d Number) Number {
	fn := math.Atan(d.Real)
	deriv := 1 / (1 + d.Real*d.Real)
	return Number{
		Real: fn,
		Emag: deriv * d.Emag,
	}
}

// Sqrt returns the square root of d.
func Sqrt(d Number) Number {
	return PowReal(d, 0.5)
}
