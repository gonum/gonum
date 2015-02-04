// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import "github.com/gonum/internal/asm"

// Inner computes the generalized inner product
//   x^T A y
// between vectors x and y with matrix A. This is only a true inner product if
// A is symmetric positive definite, though the operation works for any matrix A.
//
// Inner panics if len(x) != m or len(y) != n when A is an m x n matrix.
func Inner(x []float64, A Matrix, y []float64) float64 {
	m, n := A.Dims()
	if len(x) != m {
		panic(ErrShape)
	}
	if len(y) != n {
		panic(ErrShape)
	}
	if m == 0 || n == 0 {
		return 0
	}

	var sum float64

	switch b := A.(type) {
	case RawSymmetricer:
		bmat := b.RawSymmetric()
		for i, xi := range x {
			if xi != 0 {
				sum += xi * asm.DdotUnitary(bmat.Data[i*bmat.Stride+i:i*bmat.Stride+n], y[i:])
			}
			yi := y[i]
			if i != n-1 && yi != 0 {
				sum += yi * asm.DdotUnitary(bmat.Data[i*bmat.Stride+i+1:i*bmat.Stride+n], x[i+1:])
			}
		}
	case RawMatrixer:
		bmat := b.RawMatrix()
		for i, xi := range x {
			if xi != 0 {
				sum += xi * asm.DdotUnitary(bmat.Data[i*bmat.Stride:i*bmat.Stride+n], y)
			}
		}
	default:
		for i, xi := range x {
			for j, yj := range y {
				sum += xi * A.At(i, j) * yj
			}
		}
	}
	return sum
}
