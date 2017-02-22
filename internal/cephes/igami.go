// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Cephes Math Library Release 2.0:  April, 1987
 * Copyright 1985, 1987 by Stephen L. Moshier
 * Direct inquiries to 30 Frost Street, Cambridge, MA 02140
 */

package cephes

import (
	"math"
)

/*
 * Adapted from scipy's cephes igami.c
 */

/*
 *
 *      Inverse of complemented incomplete Gamma integral
 *
 *
 *
 * SYNOPSIS:
 *
 * a, x, p float64
 *
 * x = IgamI(a, p)
 *
 * DESCRIPTION:
 *
 * Given p, the function finds x such that
 *
 *  IgamC(a, x) = p
 *
 * Starting with the approximate value
 *
 *         3
 *  x = a t
 *
 *  where
 *
 *  t = 1 - d - ndtri(p) sqrt(d)
 *
 * and
 *
 *  d = 1/9a,
 *
 * the routine performs up to 10 Newton iterations to find the
 * root of IgamC(a, x) - p = 0.
 *
 * ACCURACY:
 *
 * Tested at random a, p in the intervals indicated.
 *
 *                a        p                      Relative error:
 * arithmetic   domain   domain     # trials      peak         rms
 *    IEEE     0.5,100   0,0.5       100000       1.0e-14     1.7e-15
 *    IEEE     0.01,0.5  0,0.5       100000       9.0e-14     3.4e-15
 *    IEEE    0.5,10000  0,0.5        20000       2.3e-13     3.8e-14
 */

/*
 * Cephes Math Library Release 2.3:  March, 1995
 * Copyright 1984, 1987, 1995 by Stephen L. Moshier
 */

// IgamI calculates the inverse of complemented incomplete Gamma integral
func IgamI(a, y0 float64) float64 {
	// bound the solution
	x0 := math.MaxFloat64
	yl := 0.0
	x1 := 0.0
	yh := 1.0
	dithresh := 5.0 * machEp

	if y0 < 0 || y0 > 1 || a <= 0 {
		panic("IgamI: Domain error")
	}

	if y0 == 0 {
		return math.MaxFloat64
	}

	if y0 == 1 {
		return 0.0
	}

	// approximation to inverse function
	d := 1.0 / (9.0 * a)
	y := 1.0 - d - Ndtri(y0)*math.Sqrt(d)
	x := a * y * y * y

	lgm := lgam(a)

	for i := 0; i < 10; i++ {
		if x > x0 || x < x1 {
			break
		}

		y = IgamC(a, x)

		if y < yl || y > yh {
			break
		}

		if y < y0 {
			x0 = x
			yl = y
		} else {
			x1 = x
			yh = y
		}

		// compute the derivative of the function at this point
		d = (a-1)*math.Log(x) - x - lgm
		if d < -maxLog {
			break
		}
		d = -math.Exp(d)

		// compute the step to the next approximation of x
		d = (y - y0) / d
		if math.Abs(d/x) < machEp {
			return x
		}
		x = x - d
	}

	d = 0.0625
	if x0 == math.MaxFloat64 {
		if x <= 0 {
			x = 1
		}
		for x0 == math.MaxFloat64 {
			x = (1 + d) * x
			y = IgamC(a, x)
			if y < y0 {
				x0 = x
				yl = y
				break
			}
			d = d + d
		}
	}

	d = 0.5
	dir := 0
	for i := 0; i < 400; i++ {
		x = x1 + d*(x0-x1)
		y = IgamC(a, x)

		lgm = (x0 - x1) / (x1 + x0)
		if math.Abs(lgm) < dithresh {
			break
		}

		lgm = (y - y0) / y0
		if math.Abs(lgm) < dithresh {
			break
		}

		if x <= 0 {
			break
		}

		if y >= y0 {
			x1 = x
			yh = y
			if dir < 0 {
				dir = 0
				d = 0.5
			} else if dir > 1 {
				d = 0.5*d + 0.5
			} else {
				d = (y0 - yl) / (yh - yl)
			}
			dir++
		} else {
			x0 = x
			yl = y
			if dir > 0 {
				dir = 0
				d = 0.5
			} else if dir < -1 {
				d = 0.5 * d
			} else {
				d = (y0 - yl) / (yh - yl)
			}
			dir--
		}
	}

	if x == 0 {
		panic("IgamI: Underflow error")
	}

	return x
}
