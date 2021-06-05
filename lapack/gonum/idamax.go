// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"
)

// Idamax finds the index of the first element having maximum absolute value.
//
//  n is number of elements in input vector(s)
//  dx is array, dimension ( 1 + ( N - 1 )*abs( incx ) )
//  incx is storage spacing between elements of dx
func (Implementation) Idamax(n int, dx []float64, incx int) (idamax int) {
	if n < 1 || incx < 0 {
		return -1 // how to handle zero length vector?
	}
	if n == 1 {
		return
	}
	dmax := math.Abs(dx[0])
	if incx == 1 {
		for i, v := range dx {
			if math.Abs(v) > dmax {
				idamax = i
			}
		}
		return
	}

	for i := incx; i < n; i += incx {
		if math.Abs(dx[i]) > dmax {
			idamax = i
		}
	}
	return
}
