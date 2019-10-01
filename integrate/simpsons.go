// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integrate

import (
	"sort"
)

// Simpson's method approximates the integral of a function f(x)
// by means of subdividing the interval of integration into segments
// and applying a method that fits a polynomial to each subinterval.
// this implementation makes use of the following composite Simpson's method
// to estimate \int_a^b f(x) dx:
//  \sum_{i=0}^{(N/2)-1} {a_0}_i * f_{2i-1} + {a_1}_i * f_{2i} + {a_2}_i * f_{2i+1}
// where N is the number of subintervals and the number of subintervals
// is even and {a_0}_i, {a_1}_i, and {a_2}_i are constants at index i given by:
//  {a_0}_i = 2 * h^{3}_{2i} - h^{3}_{2i+1} + 3 * h_{2i+1} * h^{2}_{2i} /
//            6 * h_{2i} * (h_{2i+1} + h_{2i})
//  {a_1}_i = h^{3}_{2i+1} + h^{3}_{2i} +
//            3 * h_{2i+1} * (h_{2i+1} + h_{2i}) /
//            6 * h_{2i+1} * h_{2i}
//  {a_2}_i = 2 * h^{3}_{2i+1} - h^{3}_{2i} + 3 * h_{2i} * h^{2}_{2i+1} /
//            6 * h_{2i+1} * (h_{2i+1} + h_{2i})
// where h_{k} is the length of the subinterval at k.
//
// The formula above approximates the integral of function of f
// if N is an even number of subintervals. If the number of subintervals are odd,
// the subintervals up i=0..n-2 are given by the above
// and the approximations over the second to last and last subintervals are given by:
//  {a_0}_i * f_{N-2} + {a_1}_i * f_{N-1} + {a_2}_i * f_{N}
// where the coefficients are:
//  {a_0}_i = -1 * h^{3}_{N-1} /
//             6 * h_{N-2} * (h_{N-2} + h_{N-1})
//  {a_1}_i = h^{2}_{N-1} + 3 * h^{3}_{N-1} * h_{N-2} /
//            6 * h_{N-2}
//  {a_2}_i = 2 * h^{2}_{N-1} + 3 * h^{2}_{N-1} * h_{N-2} /
//            6 * (h_{N-2} + h_{N-1})
// More information is available at:
// https://en.wikipedia.org/wiki/Simpson%27s_rule#Composite_Simpson's_rule_for_irregularly_spaced_data
//
// The (x,f) input data points must be sorted along x. One can use stat.SortWeighted to do that.
// The x and f slices must be of equal length and have length > 1.
func Simpsons(x, f []float64) float64 {
	validateInterval(x, f)

	var integral float64

	// If only three points are supplied, approximate the integral
	// as a midpoint and interpolate with Simpson's method.
	if len(x) == 3 {
		// Panic if either subinterval coordinate points match
		if x[0] == x[1] && f[0] == f[1] || x[1] == x[2] && f[0] == f[1] {
			panic("integrate: at least three unique points are required")
		}
		return (x[2] - x[0]) * (f[0]/6.0 + 4.0*f[1]/6.0 + f[2]/6.0)
	}

	for i := 1; i < len(x)-1; i += 2 {
		h0 := x[i] - x[i-1]
		h0p2 := h0 * h0
		h0p3 := h0 * h0 * h0
		h1 := x[i+1] - x[i]
		h1p2 := h1 * h1
		h1p3 := h1 * h1 * h1
		hph := h0 + h1
		f0 := f[i-1]
		f1 := f[i]
		f2 := f[i+1]
		a0 := (2*h0p3 - h1p3 + 3*h1*h0p2) / (6 * h0 * hph)
		a1 := (h1p3 + h0p3 + 3*h1*h0*hph) / (6 * h1 * h0)
		a2 := (2*h1p3 - h0p3 + 3*h0*h1p2) / (6 * h1 * hph)
		integral += a0 * f0
		integral += a1 * f1
		integral += a2 * f2
	}

	if (len(x)-1)%2 != 0 {
		h0 := x[len(x)-2] - x[len(x)-3]
		h1 := x[len(x)-1] - x[len(x)-2]
		h1p2 := h1 * h1
		h1p3 := h1 * h1 * h1
		hph := h1 + h0
		a0 := -1 * h1p3 / (6 * h0 * hph)
		a1 := (h1p2 + 3*h1*h0) / (6 * h0)
		a2 := (2*h1p2 + 3*h1*h0) / (6 * hph)
		integral += a0 * f[len(x)-3]
		integral += a1 * f[len(x)-2]
		integral += a2 * f[len(x)-1]
	}
	return integral
}

func validateInterval(x, f []float64) {
	switch {
	case len(x) != len(f):
		panic("integrate: slice length mismatch")
	case len(x) <= 2:
		panic("integrate: input data too small")
	case !sort.Float64sAreSorted(x):
		panic("integrate: must be sorted")
	}
}
