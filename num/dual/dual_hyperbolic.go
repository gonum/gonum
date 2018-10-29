// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dual

import "math"

// Sinh returns the hyperbolic sine of d.
func Sinh(d Number) Number {
	fn := math.Sinh(d.Real)
	deriv := math.Cosh(d.Real)
	return Number{
		Real: fn,
		Emag: deriv * d.Emag,
	}
}

// Cosh returns the hyperbolic cosine of d.
func Cosh(d Number) Number {
	fn := math.Cosh(d.Real)
	deriv := math.Sinh(d.Real)
	return Number{
		Real: fn,
		Emag: deriv * d.Emag,
	}
}

// Tanh returns the hyperbolic tangent of d.
func Tanh(d Number) Number {
	fn := math.Tanh(d.Real)
	deriv := 1 - fn*fn
	return Number{
		Real: fn,
		Emag: deriv * d.Emag,
	}
}

// Asinh returns the inverse hyperbolic sine of d.
func Asinh(d Number) Number {
	fn := math.Asinh(d.Real)
	deriv := 1 / math.Sqrt(d.Real*d.Real+1)
	return Number{
		Real: fn,
		Emag: deriv * d.Emag,
	}
}

// Acosh returns the inverse hyperbolic cosine of d.
func Acosh(d Number) Number {
	fn := math.Acosh(d.Real)
	deriv := 1 / math.Sqrt(d.Real*d.Real-1)
	return Number{
		Real: fn,
		Emag: deriv * d.Emag,
	}
}

// Atanh returns the inverse hyperbolic tangent of d.
func Atanh(d Number) Number {
	fn := math.Atanh(d.Real)
	deriv := 1 / (1 - d.Real*d.Real)
	return Number{
		Real: fn,
		Emag: deriv * d.Emag,
	}
}
