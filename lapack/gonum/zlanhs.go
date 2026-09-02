// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/lapack"
)

// Zlanhs returns the value of the one norm, the Frobenius norm, the infinity
// norm, or the element of largest absolute value of a complex upper Hessenberg
// matrix A.
//
// If norm is lapack.MaxColumnSum, work must have length at least n.
//
// Zlanhs is an internal routine. It is exported for testing purposes.
func (impl Implementation) Zlanhs(norm lapack.MatrixNorm, n int, a []complex128, lda int, work []float64) float64 {
	switch {
	case norm != lapack.MaxRowSum && norm != lapack.MaxAbs && norm != lapack.MaxColumnSum && norm != lapack.Frobenius:
		panic(badNorm)
	case n < 0:
		panic(nLT0)
	case lda < max(1, n):
		panic(badLdA)
	}

	if n == 0 {
		return 0
	}

	switch {
	case len(a) < (n-1)*lda+n:
		panic(shortA)
	case norm == lapack.MaxColumnSum && len(work) < n:
		panic(shortWork)
	}

	var value float64
	switch norm {
	case lapack.MaxAbs:
		for i := 0; i < n; i++ {
			minj := max(0, i-1)
			for j := minj; j < n; j++ {
				value = math.Max(value, cmplx.Abs(a[i*lda+j]))
			}
		}
	case lapack.MaxColumnSum:
		for i := 0; i < n; i++ {
			work[i] = 0
		}
		for i := 0; i < n; i++ {
			for j := max(0, i-1); j < n; j++ {
				work[j] += cmplx.Abs(a[i*lda+j])
			}
		}
		for _, v := range work[:n] {
			value = math.Max(value, v)
		}
	case lapack.MaxRowSum:
		for i := 0; i < n; i++ {
			minj := max(0, i-1)
			var sum float64
			for j := minj; j < n; j++ {
				sum += cmplx.Abs(a[i*lda+j])
			}
			value = math.Max(value, sum)
		}
	case lapack.Frobenius:
		scale := 0.0
		sum := 1.0
		for i := 0; i < n; i++ {
			minj := max(0, i-1)
			scale, sum = zlassq(impl, n-minj, a[i*lda+minj:], 1, scale, sum)
		}
		value = scale * math.Sqrt(sum)
	}
	return value
}
