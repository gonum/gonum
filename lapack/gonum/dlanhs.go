// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/lapack"
)

// Dlanhs  returns the value of the one norm,  or the Frobenius norm, or
// the  infinity norm,  or the  element of  largest absolute value  of a
// Hessenberg matrix A.
func (impl Implementation) Dlanhs(norm lapack.MatrixNorm, n int, a []float64, lda int, work []float64) (resultNorm float64) {
	if n == 0 {
		return 0 // Early return.
	}
	var value, sum float64
	switch norm {
	default:
		panic(badNorm)
	case lapack.MaxAbs:
		// Find max(abs(A(i,j))).
		for j := 0; j < n; j++ {
			imax := min(n-1, j+1)
			for i := 0; i <= imax; i++ {
				value = math.Max(value, math.Abs(a[i*lda+j]))
			}
		}
	case lapack.MaxColumnSum:
		// Find norm1(A).
		for j := 0; j < n; j++ {
			sum = 0
			imax := min(n-1, j+1)
			for i := 0; i <= imax; i++ {
				sum += math.Abs(a[i*lda+j])
			}
			value = math.Max(value, sum)
		}
	case lapack.MaxRowSum:
		// Find normI(A).
		for i := 0; i < n; i++ {
			work[i] = 0
		}
		for j := 0; j < n; j++ {
			imax := min(n-1, j+1)
			for i := 0; i <= imax; i++ {
				work[i] += math.Abs(a[i*lda+j])
			}
		}
		for i := 0; i < n; i++ {
			value = math.Max(value, work[i])
		}
	case lapack.Frobenius:
		// Find normF(A).
		scale := 0.0
		sum = 1.0
		for j := 0; j < n; j++ {
			scale, sum = impl.Dlassq(min(n, j+2), a[j:], lda, scale, sum)
		}
		value = scale * math.Sqrt(sum)
	}
	return value
}
