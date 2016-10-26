// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Cephes Math Library Release 2.1:  December, 1988
 * Copyright 1984, 1987, 1988 by Stephen L. Moshier
 * Direct inquiries to 30 Frost Street, Cambridge, MA 02140
 */

package cephes

// polevl evaluates a polynomial of degree N
//  y = c_0 + c_1 x_1 + c_2 x_2^2 ...
// where the coefficients are stored in reverse order, i.e. coef[0] = c_n and
// coef[n] = c_0.
func polevl(x float64, coef []float64, n int) float64 {
	ans := coef[0]
	for i := 1; i <= n; i++ {
		ans = ans*x + coef[i]
	}
	return ans
}

// p1evl is the same as polevl, except c_n is assumed to be 1 and is not included
// in the slice.
func p1evl(x float64, coef []float64, n int) float64 {
	ans := x + coef[0]
	for i := 1; i <= n-1; i++ {
		ans = ans*x + coef[i]
	}
	return ans
}
