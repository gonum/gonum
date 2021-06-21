// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import "gonum.org/v1/gonum/blas"

var _ blas.UsFloat64Level1 = Implementation{}

// Ddot computes the dot product of a sparse x vector and a dense y vector
//  \sum_i x[i]*y[i]
func (Implementation) Dusdot(nz int, x []float64, index []int, y []float64, incY int) (w float64) {
	switch {
	case incY == 0:
		panic(zeroIncY)
	case len(x) < nz:
		panic(shortX)
	case len(index) != len(x):
		panic(badXIndexLen)
	}

	if nz <= 0 {
		if nz == 0 {
			return 0
		}
		panic(nLT0)
	}
	for i := 0; i < nz; i++ {
		w += x[i] * y[index[i]*incY]
	}
	return w
}

// Daxpy adds alpha times sparse x to dense y
//  y[i] += alpha * x[i] for all i
func (Implementation) Dusaxpy(nz int, alpha float64, x []float64, index []int, y []float64, incY int) {
	switch {
	case incY == 0:
		panic(zeroIncY)
	case len(x) < nz:
		panic(shortX)
	case len(index) != len(x):
		panic(badXIndexLen)
	}

	if nz <= 0 {
		if nz == 0 {
			return
		}
		panic(nLT0)
	}
	if alpha == 0 {
		return
	}

	for i := 0; i < nz; i++ {
		y[index[i]*incY] += alpha * x[i]
	}
}

// Dusga Gathers y values at non-zero places of sparse x
// and places them in x.
//  x[i] = y[i]   for all i where x[i] != 0
func (Implementation) Dusga(nz int, y []float64, incY int, x []float64, index []int) {
	switch {
	case incY == 0:
		panic(zeroIncY)
	case len(x) < nz:
		panic(shortX)
	case len(index) != len(x):
		panic(badXIndexLen)
	}

	if nz <= 0 {
		if nz == 0 {
			return
		}
		panic(nLT0)
	}
	for i := 0; i < nz; i++ {
		x[i] = y[index[i]*incY]
	}
}

// Dusgz gathers values of y into sparse x at non-zero entries and
// zeros those y values.
//  x[i] = y[i]  where x[i] != 0
//  y[i] = 0     where x[i] != 0
func (Implementation) Dusgz(nz int, y []float64, incY int, x []float64, index []int) {
	switch {
	case incY == 0:
		panic(zeroIncY)
	case len(x) < nz:
		panic(shortX)
	case len(index) != len(x):
		panic(badXIndexLen)
	}

	if nz <= 0 {
		if nz == 0 {
			return
		}
		panic(nLT0)
	}
	for i := 0; i < nz; i++ {
		x[i] = y[index[i]*incY]
		y[index[i]*incY] = 0
	}
}

// Dussc copies non zero values of x into y.
//  y[i] = x[i]    where x[i] != 0
func (Implementation) Dussc(nz int, x []float64, y []float64, incY int, index []int) {
	switch {
	case incY == 0:
		panic(zeroIncY)
	case len(x) < nz:
		panic(shortX)
	case len(index) != len(x):
		panic(badXIndexLen)
	}

	if nz <= 0 {
		if nz == 0 {
			return
		}
		panic(nLT0)
	}
	for i := 0; i < nz; i++ {
		y[index[i]*incY] = x[i]
	}
}
