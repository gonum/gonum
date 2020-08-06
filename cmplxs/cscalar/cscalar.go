// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cscalar

import (
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/floats/scalar"
)

// EqualWithinAbs returns true when a and b have an absolute difference
// not greater than tol.
func EqualWithinAbs(a, b complex128, tol float64) bool {
	return a == b || cmplx.Abs(a-b) <= tol
}

// minNormalFloat64 is the smallest normal number. For 64 bit IEEE-754
// floats this is 2^{-1022}.
const minNormalFloat64 = 0x1p-1022

// EqualWithinRel returns true when the difference between a and b
// is not greater than tol times the greater absolute value of a and b,
//  abs(a-b) <= tol * max(abs(a), abs(b)).
func EqualWithinRel(a, b complex128, tol float64) bool {
	if a == b {
		return true
	}

	delta := cmplx.Abs(a - b)
	if delta <= minNormalFloat64 {
		return delta <= tol*minNormalFloat64
	}
	// We depend on the division in this relationship to identify
	// infinities.
	return delta/math.Max(cmplx.Abs(a), cmplx.Abs(b)) <= tol
}

// EqualWithinAbsOrRel returns true when a and b are equal to within
// the absolute or relative tolerances. See EqualWithinAbs and
// EqualWithinRel for details.
func EqualWithinAbsOrRel(a, b complex128, absTol, relTol float64) bool {
	return EqualWithinAbs(a, b, absTol) || EqualWithinRel(a, b, relTol)
}

// ParseWithNA converts the string s to a complex128 in value.
// If s equals missing, weight is returned as 0, otherwise 1.
func ParseWithNA(s, missing string) (value complex128, weight float64, err error) {
	if s == missing {
		return 0, 0, nil
	}
	value, err = parse(s)
	if err == nil {
		weight = 1
	}
	return value, weight, err
}

// Round returns the half away from zero rounded value of x with prec precision.
//
// Special cases are:
// 	Round(±0) = +0
// 	Round(±Inf) = ±Inf
// 	Round(NaN) = NaN
func Round(x complex128, prec int) complex128 {
	if x == 0 {
		// Make sure zero is returned
		// without the negative bit set.
		return 0
	}
	return complex(scalar.Round(real(x), prec), scalar.Round(imag(x), prec))
}

// RoundEven returns the half even rounded value of x with prec precision.
//
// Special cases are:
// 	RoundEven(±0) = +0
// 	RoundEven(±Inf) = ±Inf
// 	RoundEven(NaN) = NaN
func RoundEven(x complex128, prec int) complex128 {
	if x == 0 {
		// Make sure zero is returned
		// without the negative bit set.
		return 0
	}
	return complex(scalar.RoundEven(real(x), prec), scalar.RoundEven(imag(x), prec))
}

// Same returns true when the inputs have the same value, allowing NaN equality.
func Same(a, b complex128) bool {
	return a == b || (cmplx.IsNaN(a) && cmplx.IsNaN(b))
}
