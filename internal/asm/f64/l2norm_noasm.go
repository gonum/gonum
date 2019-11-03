// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !amd64 noasm appengine safe

package f64

import "math"

// L2NormUnitary returns the L2-norm of x.
func L2NormUnitary(x []float64) (sum float64) {
	var scale float64
	sumSquares := 1.0
	for _, v := range x {
		if v == 0 {
			continue
		}
		absxi := math.Abs(v)
		if math.IsNaN(absxi) {
			return math.NaN()
		}
		if scale < absxi {
			s := scale / absxi
			sumSquares = 1 + sumSquares*s*s
			scale = absxi
		} else {
			s := absxi / scale
			sumSquares += s * s
		}
	}
	if math.IsInf(scale, 1) {
		return math.Inf(1)
	}
	return scale * math.Sqrt(sumSquares)
}
