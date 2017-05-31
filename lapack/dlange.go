// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lapack

import "math"

// Dlange computes the matrix norm of the general m×n matrix a. The input norm
// specifies the norm computed.
//  MaxAbs: the maximum absolute value of an element.
//  MaxColumnSum: the maximum column sum of the absolute values of the entries.
//  MaxRowSum: the maximum row sum of the absolute values of the entries.
//  NormFrob: the square root of the sum of the squares of the entries.
// If norm == MaxColumnSum, work must be of length n, and this function will panic otherwise.
// There are no restrictions on work for the other matrix norms.
func (impl Implementation) Dlange(norm MatrixNorm, m, n int, a []float64, lda int, work []float64) float64 {
	// TODO(btracey): These should probably be refactored to use BLAS calls.
	checkMatrix(m, n, a, lda)
	switch norm {
	case MaxRowSum, MaxColumnSum, NormFrob, MaxAbs:
	default:
		panic(badNorm)
	}
	if norm == MaxColumnSum && len(work) < n {
		panic(badWork)
	}
	if m == 0 && n == 0 {
		return 0
	}
	if norm == MaxAbs {
		var value float64
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				value = math.Max(value, math.Abs(a[i*lda+j]))
			}
		}
		return value
	}
	if norm == MaxColumnSum {
		if len(work) < n {
			panic(badWork)
		}
		for i := 0; i < n; i++ {
			work[i] = 0
		}
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				work[j] += math.Abs(a[i*lda+j])
			}
		}
		var value float64
		for i := 0; i < n; i++ {
			value = math.Max(value, work[i])
		}
		return value
	}
	if norm == MaxRowSum {
		var value float64
		for i := 0; i < m; i++ {
			var sum float64
			for j := 0; j < n; j++ {
				sum += math.Abs(a[i*lda+j])
			}
			value = math.Max(value, sum)
		}
		return value
	}
	if norm == NormFrob {
		var value float64
		scale := 0.0
		sum := 1.0
		for i := 0; i < m; i++ {
			scale, sum = impl.Dlassq(n, a[i*lda:], 1, scale, sum)
		}
		value = scale * math.Sqrt(sum)
		return value
	}
	panic("lapack: bad matrix norm")
}
