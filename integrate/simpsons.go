// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integrate

import "sort"

// Simpsons returns an approximate value of the integral
//
//	\int_a^b f(x)dx
//
// computed using the Simpsons's method. The function f is given as a slice of
// samples evaluated at locations in x, that is,
//
//	f[i] = f(x[i]), x[0] = a, x[len(x)-1] = b
//
// The slice x must be sorted in strictly increasing order. x and f must be of
// equal length and the length must be at least 3.
//
// See https://en.wikipedia.org/wiki/Simpson%27s_rule#Composite_Simpson's_rule_for_irregularly_spaced_data
// for more information.
func Simpsons(x, f []float64) float64 {
	n := len(x)
	switch {
	case len(f) != n:
		panic("integrate: slice length mismatch")
	case n < 3:
		panic("integrate: input data too small")
	case !sort.Float64sAreSorted(x):
		panic("integrate: must be sorted")
	}

	var integral float64
	for i := 1; i < n-1; i += 2 {
		h0 := x[i] - x[i-1]
		h1 := x[i+1] - x[i]
		if h0 == 0 || h1 == 0 {
			panic("integrate: repeated abscissa")
		}
		h0p2 := h0 * h0
		h0p3 := h0 * h0 * h0
		h1p2 := h1 * h1
		h1p3 := h1 * h1 * h1
		hph := h0 + h1
		a0 := (2*h0p3 - h1p3 + 3*h1*h0p2) / (6 * h0 * hph)
		a1 := (h0p3 + h1p3 + 3*h0*h1*hph) / (6 * h0 * h1)
		a2 := (-h0p3 + 2*h1p3 + 3*h0*h1p2) / (6 * h1 * hph)
		integral += a0 * f[i-1]
		integral += a1 * f[i]
		integral += a2 * f[i+1]
	}

	if n%2 == 0 {
		h0 := x[n-2] - x[n-3]
		h1 := x[n-1] - x[n-2]
		if h0 == 0 || h1 == 0 {
			panic("integrate: repeated abscissa")
		}
		h1p2 := h1 * h1
		h1p3 := h1 * h1 * h1
		hph := h0 + h1
		a0 := -1 * h1p3 / (6 * h0 * hph)
		a1 := (h1p2 + 3*h0*h1) / (6 * h0)
		a2 := (2*h1p2 + 3*h0*h1) / (6 * hph)
		integral += a0 * f[n-3]
		integral += a1 * f[n-2]
		integral += a2 * f[n-1]
	}

	return integral
}
