// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/lapack"
)

// Zlange returns the value of the specified norm of a general m×n complex
// matrix A:
//
//	lapack.MaxAbs:       the maximum absolute value of any element.
//	lapack.MaxColumnSum: the maximum column sum of the absolute values of the elements (1-norm).
//	lapack.MaxRowSum:    the maximum row sum of the absolute values of the elements (infinity-norm).
//	lapack.Frobenius:    the square root of the sum of the squares of the absolute values of the elements (Frobenius norm).
//
// If norm == lapack.MaxColumnSum, work must be of length n, and this function
// will panic otherwise. There are no restrictions on work for the other matrix
// norms.
//
// Zlange is an internal routine. It is exported for testing purposes.
func (impl Implementation) Zlange(norm lapack.MatrixNorm, m, n int, a []complex128, lda int, work []float64) float64 {
	switch {
	case norm != lapack.MaxRowSum && norm != lapack.MaxColumnSum && norm != lapack.Frobenius && norm != lapack.MaxAbs:
		panic(badNorm)
	case m < 0:
		panic(mLT0)
	case n < 0:
		panic(nLT0)
	case lda < max(1, n):
		panic(badLdA)
	}

	// Quick return if possible.
	if m == 0 || n == 0 {
		return 0
	}

	switch {
	case len(a) < (m-1)*lda+n:
		panic(shortA)
	case norm == lapack.MaxColumnSum && len(work) < n:
		panic(shortWork)
	}

	switch norm {
	case lapack.MaxAbs:
		var value float64
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				value = math.Max(value, cmplx.Abs(a[i*lda+j]))
			}
		}
		return value
	case lapack.MaxColumnSum:
		for i := 0; i < n; i++ {
			work[i] = 0
		}
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				work[j] += cmplx.Abs(a[i*lda+j])
			}
		}
		var value float64
		for i := 0; i < n; i++ {
			value = math.Max(value, work[i])
		}
		return value
	case lapack.MaxRowSum:
		var value float64
		for i := 0; i < m; i++ {
			var sum float64
			for j := 0; j < n; j++ {
				sum += cmplx.Abs(a[i*lda+j])
			}
			value = math.Max(value, sum)
		}
		return value
	default:
		// lapack.Frobenius: accumulate real and imag parts via Dlassq.
		scale := 0.0
		sum := 1.0
		for i := 0; i < m; i++ {
			scale, sum = zlassq(impl, n, a[i*lda:], 1, scale, sum)
		}
		return scale * math.Sqrt(sum)
	}
}

// zlassq updates a scaled sum of squares for a complex vector by passing the
// real and imaginary parts through Dlassq. Mirrors netlib ZLASSQ.
func zlassq(impl Implementation, n int, x []complex128, incx int, scale, sumsq float64) (float64, float64) {
	if n <= 0 {
		return scale, sumsq
	}
	// Build real and imaginary slices with the requested increment.
	// n is usually small (a row), so allocation is cheap; alternatively one
	// could accumulate in-line, but Dlassq already handles overflow/underflow.
	re := make([]float64, n)
	im := make([]float64, n)
	for i, ix := 0, 0; i < n; i++ {
		re[i] = real(x[ix])
		im[i] = imag(x[ix])
		ix += incx
	}
	scale, sumsq = impl.Dlassq(n, re, 1, scale, sumsq)
	scale, sumsq = impl.Dlassq(n, im, 1, scale, sumsq)
	return scale, sumsq
}
