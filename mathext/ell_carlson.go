// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import (
	"math"
)

// CarlsonRF computes the symmetric elliptic integral R_F(x,y,z):
//
//	R_F(x,y,z) = (1/2)\int_{0}^{\infty}{s^{-1}(t)} dt,
//	s(t) = \sqrt{t+x}\sqrt{t+y}\sqrt{t+z}.
//
// See: http://dlmf.nist.gov/19.16.E1 for the definition.
//
// See: http://doi.org/10.1145/355958.355970 for the original Fortran code.
//
// See: http://dx.doi.org/10.1007/BF02198293 for the modified method of computation.
func CarlsonRF(x, y, z float64) float64 {
	const lower = 1.1125369292536006915451163586662020321096079902312e-307 // 5*2^-1022
	const upper = 1 / lower
	if x < 0 || y < 0 || z < 0 || math.IsNaN(x) || math.IsNaN(y) || math.IsNaN(z) {
		return math.NaN()
	}
	if upper < x || upper < y || upper < z {
		return math.NaN()
	}
	if x+y < lower || y+z < lower || z+x < lower {
		return math.NaN()
	}

	const tol = 1.2674918778210762260320167734407048051023273568443e-02 // (3ε)^(1/8)
	A0 := (x + y + z) / 3
	An := A0
	Q := math.Max(math.Max(math.Abs(A0-x), math.Abs(A0-y)), math.Abs(A0-z)) / tol
	xn, yn, zn := x, y, z
	mul := 1.0

	for Q >= mul*math.Abs(An) {
		xnsqrt, ynsqrt, znsqrt := math.Sqrt(xn), math.Sqrt(yn), math.Sqrt(zn)
		lambda := xnsqrt*ynsqrt + ynsqrt*znsqrt + znsqrt*xnsqrt
		An = (An + lambda) * 0.25
		xn = (xn + lambda) * 0.25
		yn = (yn + lambda) * 0.25
		zn = (zn + lambda) * 0.25
		mul *= 4
	}

	X := (A0 - x) / (mul * An)
	Y := (A0 - y) / (mul * An)
	Z := -(X + Y)
	E2 := X*Y - Z*Z
	E3 := X * Y * Z

	// http://dlmf.nist.gov/19.36.E1
	return (E3*(6930*E3+E2*(15015*E2-16380)+17160) + E2*((10010-5775*E2)*E2-24024) + 240240) / (240240 * math.Sqrt(An))
}

// CarlsonRD computes the symmetric elliptic integral R_D(x,y,z):
//
//	R_D(x,y,z) = (1/2)\int_{0}^{\infty}{s^{-1}(t)}{(t+z)^{-1}} dt,
//	s(t) = \sqrt{t+x}\sqrt{t+y}\sqrt{t+z}.
//
// See: http://dlmf.nist.gov/19.16.E5 for the definition.
//
// See: http://doi.org/10.1145/355958.355970 for the original Fortran code.
//
// See: http://dx.doi.org/10.1007/BF02198293 for the modified method of computation.
func CarlsonRD(x, y, z float64) float64 {
	const lower = 4.8095540743116787026618007863123676393525016818363e-103 // (5*2^-1022)^(1/3)
	const upper = 1 / lower
	if x < 0 || y < 0 || math.IsNaN(x) || math.IsNaN(y) || math.IsNaN(z) {
		return math.NaN()
	}
	if upper < x || upper < y || upper < z {
		return math.NaN()
	}
	if x+y < lower || z < lower {
		return math.NaN()
	}

	const tol = 9.03511693393157704747601225470683249938574888493817e-03 // (ε/5)^(1/8)
	A0 := (x + y + 3*z) / 5
	An := A0
	Q := math.Max(math.Max(math.Abs(A0-x), math.Abs(A0-y)), math.Abs(A0-z)) / tol
	xn, yn, zn := x, y, z
	mul, s := 1.0, 0.0

	for Q >= mul*math.Abs(An) {
		xnsqrt, ynsqrt, znsqrt := math.Sqrt(xn), math.Sqrt(yn), math.Sqrt(zn)
		lambda := xnsqrt*ynsqrt + ynsqrt*znsqrt + znsqrt*xnsqrt
		s += 1 / (mul * znsqrt * (zn + lambda))
		An = (An + lambda) * 0.25
		xn = (xn + lambda) * 0.25
		yn = (yn + lambda) * 0.25
		zn = (zn + lambda) * 0.25
		mul *= 4
	}

	X := (A0 - x) / (mul * An)
	Y := (A0 - y) / (mul * An)
	Z := -(X + Y) / 3
	E2 := X*Y - 6*Z*Z
	E3 := (3*X*Y - 8*Z*Z) * Z
	E4 := 3 * (X*Y - Z*Z) * Z * Z
	E5 := X * Y * Z * Z * Z

	// http://dlmf.nist.gov/19.36.E2
	return ((471240-540540*E2)*E5+(612612*E2-540540*E3-556920)*E4+E3*(306306*E3+E2*(675675*E2-706860)+680680)+E2*((417690-255255*E2)*E2-875160)+4084080)/(4084080*mul*An*math.Sqrt(An)) + 3*s
}

// EllipticF computes the Legendre's elliptic integral of the 1st kind F(\phi|m):
//
//	F(\phi|m) = \int_{0}^{\phi}1 / \sqrt{1-m\sin^2\theta} d\theta
//
// Legendre's elliptic integrals can be expressed as symmetric elliptic integrals, in this case:
//
//	F(\phi|m) = \sin\phi R_F(\cos^2\phi,1-m\sin^2\phi,1)
//
// See http://dlmf.nist.gov/19.2.E4 for the definition.
func EllipticF(phi, m float64) float64 {
	s, c := math.Sincos(phi)
	return s * CarlsonRF(c*c, 1-m*s*s, 1)
}

// EllipticE computes the Legendre's elliptic integral of the 2nd kind E(\phi|m):
//
//	E(\phi|m) = \int_{0}^{\phi} \sqrt{1-m\sin^2\theta} d\theta
//
// Legendre's elliptic integrals can be expressed as symmetric elliptic integrals, in this case:
//
//	E(\phi|m) = \sin\phi R_F(\cos^2\phi,1-m\sin^2\phi,1)-(m/3)\sin^3\phi R_D(\cos^2\phi,1-m\sin^2\phi,1)
//
// See http://dlmf.nist.gov/19.2.E5 for the definition.
func EllipticE(phi, m float64) float64 {
	s, c := math.Sincos(phi)
	x, y := c*c, 1-m*s*s
	return s * (CarlsonRF(x, y, 1) - (m/3)*s*s*CarlsonRD(x, y, 1))
}
