// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"gonum.org/v1/gonum/internal/asm/c128"
)

// Zaxpy adds alpha times x to y:
//  y[i] += alpha * x[i] for all i
func (Implementation) Zaxpy(n int, alpha complex128, x []complex128, incX int, y []complex128, incY int) {
	if incX == 0 {
		panic(zeroIncX)
	}
	if incY == 0 {
		panic(zeroIncY)
	}
	if n < 1 {
		if n == 0 {
			return
		}
		panic(negativeN)
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic(badX)
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic(badY)
	}
	if zabs(alpha) == 0 {
		return
	}
	if incX == 1 && incY == 1 {
		c128.AxpyUnitary(alpha, x[:n], y[:n])
		return
	}
	var ix, iy int
	if incX < 0 {
		ix = (1 - n) * incX
	}
	if incY < 0 {
		iy = (1 - n) * incY
	}
	c128.AxpyInc(alpha, x, y, uintptr(n), uintptr(incX), uintptr(incY), uintptr(ix), uintptr(iy))
}

// Zcopy copies the vector x to vector y.
func (Implementation) Zcopy(n int, x []complex128, incX int, y []complex128, incY int) {
	if incX == 0 {
		panic(zeroIncX)
	}
	if incY == 0 {
		panic(zeroIncY)
	}
	if n < 1 {
		if n == 0 {
			return
		}
		panic(negativeN)
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic(badX)
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic(badY)
	}
	if incX == 1 && incY == 1 {
		copy(y[:n], x[:n])
		return
	}
	var ix, iy int
	if incX < 0 {
		ix = (-n + 1) * incX
	}
	if incY < 0 {
		iy = (-n + 1) * incY
	}
	for i := 0; i < n; i++ {
		y[iy] = x[ix]
		ix += incX
		iy += incY
	}
}

// Zdotc computes the dot product
//  x^H · y
// of two complex vectors x and y.
func (Implementation) Zdotc(n int, x []complex128, incX int, y []complex128, incY int) complex128 {
	if incX == 0 {
		panic(zeroIncX)
	}
	if incY == 0 {
		panic(zeroIncY)
	}
	if n <= 0 {
		if n == 0 {
			return 0
		}
		panic(negativeN)
	}
	if incX == 1 && incY == 1 {
		if len(x) < n {
			panic(badX)
		}
		if len(y) < n {
			panic(badY)
		}
		return c128.DotcUnitary(x[:n], y[:n])
	}
	var ix, iy int
	if incX < 0 {
		ix = (-n + 1) * incX
	}
	if incY < 0 {
		iy = (-n + 1) * incY
	}
	if ix >= len(x) || (n-1)*incX >= len(x) {
		panic(badX)
	}
	if iy >= len(y) || (n-1)*incY >= len(y) {
		panic(badY)
	}
	return c128.DotcInc(x, y, uintptr(n), uintptr(incX), uintptr(incY), uintptr(ix), uintptr(iy))
}

// Zdotu computes the dot product
//  x^T · y
// of two complex vectors x and y.
func (Implementation) Zdotu(n int, x []complex128, incX int, y []complex128, incY int) complex128 {
	if incX == 0 {
		panic(zeroIncX)
	}
	if incY == 0 {
		panic(zeroIncY)
	}
	if n <= 0 {
		if n == 0 {
			return 0
		}
		panic(negativeN)
	}
	if incX == 1 && incY == 1 {
		if len(x) < n {
			panic(badX)
		}
		if len(y) < n {
			panic(badY)
		}
		return c128.DotuUnitary(x[:n], y[:n])
	}
	var ix, iy int
	if incX < 0 {
		ix = (-n + 1) * incX
	}
	if incY < 0 {
		iy = (-n + 1) * incY
	}
	if ix >= len(x) || (n-1)*incX >= len(x) {
		panic(badX)
	}
	if iy >= len(y) || (n-1)*incY >= len(y) {
		panic(badY)
	}
	return c128.DotuInc(x, y, uintptr(n), uintptr(incX), uintptr(incY), uintptr(ix), uintptr(iy))
}
