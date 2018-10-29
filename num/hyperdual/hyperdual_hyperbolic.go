// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hyperdual

import "math"

// Sinh returns the hyperbolic sine of d.
func Sinh(d Number) Number {
	fn := math.Sinh(d.Real)
	deriv := math.Cosh(d.Real)
	return Number{
		Real:    fn,
		E1mag:   deriv * d.E1mag,
		E2mag:   deriv * d.E2mag,
		E1E2mag: deriv*d.E1E2mag + fn*d.E1mag*d.E2mag,
	}
}

// Cosh returns the hyperbolic cosine of d.
func Cosh(d Number) Number {
	fn := math.Cosh(d.Real)
	deriv := math.Sinh(d.Real)
	return Number{
		Real:    fn,
		E1mag:   deriv * d.E1mag,
		E2mag:   deriv * d.E2mag,
		E1E2mag: deriv*d.E1E2mag + fn*d.E1mag*d.E2mag,
	}
}

// Tanh returns the hyperbolic tangent of d.
func Tanh(d Number) Number {
	fn := math.Tanh(d.Real)
	deriv := 1 - fn*fn
	return Number{
		Real:    fn,
		E1mag:   deriv * d.E1mag,
		E2mag:   deriv * d.E2mag,
		E1E2mag: deriv*d.E1E2mag - d.E1mag*d.E2mag*(2*fn*deriv),
	}
}

// Asinh returns the hyperbolic inverse sine of d.
func Asinh(d Number) Number {
	fn := math.Asinh(d.Real)
	deriv1 := d.Real*d.Real + 1
	deriv := 1 / math.Sqrt(deriv1)
	return Number{
		Real:    fn,
		E1mag:   deriv * d.E1mag,
		E2mag:   deriv * d.E2mag,
		E1E2mag: deriv*d.E1E2mag + d.E1mag*d.E2mag*(-d.Real*(deriv/deriv1)),
	}
}

// Acosh returns the hyperbolic inverse cosine of d.
func Acosh(d Number) Number {
	fn := math.Acosh(d.Real)
	deriv1 := d.Real*d.Real - 1
	deriv := 1 / math.Sqrt(deriv1)
	return Number{
		Real:    fn,
		E1mag:   deriv * d.E1mag,
		E2mag:   deriv * d.E2mag,
		E1E2mag: deriv*d.E1E2mag + d.E1mag*d.E2mag*(-d.Real*(deriv/deriv1)),
	}
}

// Atanh returns the hyperbolic inverse tangent of d.
func Atanh(d Number) Number {
	fn := math.Atanh(d.Real)
	deriv1 := 1 - d.Real*d.Real
	deriv := 1 / deriv1
	return Number{
		Real:    fn,
		E1mag:   deriv * d.E1mag,
		E2mag:   deriv * d.E2mag,
		E1E2mag: deriv*d.E1E2mag + d.E1mag*d.E2mag*(2*d.Real/(deriv1*deriv1)),
	}
}
