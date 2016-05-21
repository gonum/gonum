// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !amd64 noasm appengine

package f64

func AddConst(alpha float64, x []float64) {
	for i := range x {
		x[i] += alpha
	}
}
