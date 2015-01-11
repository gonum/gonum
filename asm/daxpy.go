// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !amd64

package asm

func DaxpyUnitary(alpha float64, x, y []float64) {
	for i, v := range x {
		y[i] += alpha * v
	}
}

func DaxpyInc(alpha float64, x, y []float64, n, incX, incY, ix, iy uintptr) {
	for i := 0; i < int(n); i++ {
		y[iy] += alpha * x[ix]
		ix += incX
		iy += incY
	}
}
