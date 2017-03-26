// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !amd64 noasm appengine

package c128

func AxpyUnitary(alpha complex128, x, y []complex128) {
	for i, v := range x {
		y[i] += alpha * v
	}
}

func AxpyUnitaryTo(dst []complex128, alpha complex128, x, y []complex128) {
	for i, v := range x {
		dst[i] = alpha*v + y[i]
	}
}

func AxpyInc(alpha complex128, x, y []complex128, n, incX, incY, ix, iy uintptr) {
	for i := 0; i < int(n); i++ {
		y[iy] += alpha * x[ix]
		ix += incX
		iy += incY
	}
}

func AxpyIncTo(dst []complex128, incDst, idst uintptr, alpha complex128, x, y []complex128, n, incX, incY, ix, iy uintptr) {
	for i := 0; i < int(n); i++ {
		dst[idst] = alpha*x[ix] + y[iy]
		ix += incX
		iy += incY
		idst += incDst
	}
}
