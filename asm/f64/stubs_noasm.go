// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !amd64 noasm appengine

package f64

import "math"

func AddConst(alpha float64, x []float64) {
	for i := range x {
		x[i] += alpha
	}
}

func AbsSum(x []float64) (sum float64) {
	for _, v := range x {
		sum += math.Abs(v)
	}
	return sum
}

func AbsSumInc(x []float64, n, incX int) (sum float64) {
	for i := 0; i < n; i++ {
		sum += math.Abs(x[i*incX])
	}
	return sum
}
