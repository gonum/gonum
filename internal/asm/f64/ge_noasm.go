// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !go1.8 !amd64 noasm appengine

package f64

// Ger performs the rank-one operation
//  A += alpha * x * y^T
// where A is an m×n dense matrix, x and y are vectors, and alpha is a scalar.
func Ger(m, n uintptr, alpha float64, x []float64, incX uintptr, y []float64, incY uintptr, a []float64, lda uintptr) {
	if incX == 1 && incY == 1 {
		x = x[:m]
		y = y[:n]
		for i, xv := range x {
			AxpyUnitary(alpha*xv, y, a[uintptr(i)*lda:uintptr(i)*lda+n])
		}
		return
	}

	var ky, kx uintptr
	if int(incY) < 0 {
		ky = uintptr(-int(n-1) * int(incY))
	}
	if int(incX) < 0 {
		kx = uintptr(-int(m-1) * int(incX))
	}

	ix := kx
	for i := 0; i < int(m); i++ {
		AxpyInc(alpha*x[ix], y, a[uintptr(i)*lda:uintptr(i)*lda+n], n, incY, 1, ky, 0)
		ix += incX
	}
}
