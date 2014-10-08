// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

// Inner computes the generalized inner product between x and y with matrix A.
//  x^T A y
// This is only a true inner product if m is symmetric positive definite, though
// the operation works for any matrix A.
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
	case RawMatrixer:
		bmat := b.RawMatrix()
		for i, xi := range x {
			for j, yj := range y {
				sum += xi * bmat.Data[i*bmat.Stride+j] * yj
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
