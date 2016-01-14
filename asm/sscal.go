// Generated code do not edit. Run `go generate`.

// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

func SscalUnitary(alpha float32, x []float32) {
	for i := range x {
		x[i] *= alpha
	}
}

func SscalUnitaryTo(dst []float32, alpha float32, x []float32) {
	for i, v := range x {
		dst[i] = alpha * v
	}
}

func SscalInc(alpha float32, x []float32, n, incX, ix uintptr) {
	for i := 0; i < int(n); i++ {
		x[ix] *= alpha
		ix += incX
	}
}

func SscalIncTo(dst []float32, incDst, idst uintptr, alpha float32, x []float32, n, incX, ix uintptr) {
	for i := 0; i < int(n); i++ {
		dst[idst] = alpha * x[ix]
		ix += incX
		idst += incDst
	}
}
