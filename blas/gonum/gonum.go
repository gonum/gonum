// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate ./single_precision.bash

package gonum

import "math"

type Implementation struct{}

// [SD]gemm behavior constants. These are kept here to keep them out of the
// way during single precision code genration.
const (
	blockSize   = 64 // b x b matrix
	minParBlock = 4  // minimum number of blocks needed to go parallel
	buffMul     = 4  // how big is the buffer relative to the number of workers
)

// subMul is a common type shared by [SD]gemm.
type subMul struct {
	i, j int // index of block
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func checkZhbMatrix(name byte, n, k int, ab []complex128, ldab int) {
	if n < 0 {
		panic(nLT0)
	}
	if k < 0 {
		panic(kLT0)
	}
	if ldab < k+1 {
		panic("blas: illegal stride of Hermitian band matrix " + string(name))
	}
	if len(ab) < (n-1)*ldab+k+1 {
		panic("blas: insufficient " + string(name) + " Hermitian band matrix slice length")
	}
}

func checkZtbMatrix(name byte, n, k int, ab []complex128, ldab int) {
	if n < 0 {
		panic(nLT0)
	}
	if k < 0 {
		panic(kLT0)
	}
	if ldab < k+1 {
		panic("blas: illegal stride of triangular band matrix " + string(name))
	}
	if len(ab) < (n-1)*ldab+k+1 {
		panic("blas: insufficient " + string(name) + " triangular band matrix slice length")
	}
}

func checkZVector(name byte, n int, x []complex128, incX int) {
	if n < 0 {
		panic(nLT0)
	}
	if incX == 0 {
		panic(zeroIncX)
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: insufficient " + string(name) + " vector slice length")
	}
}

// blocks returns the number of divisions of the dimension length with the given
// block size.
func blocks(dim, bsize int) int {
	return (dim + bsize - 1) / bsize
}

// dcabs1 returns |real(z)|+|imag(z)|.
func dcabs1(z complex128) float64 {
	return math.Abs(real(z)) + math.Abs(imag(z))
}
