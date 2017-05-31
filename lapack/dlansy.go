// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lapack

import (
	"math"

	"gonum.org/v1/gonum/blas"
)

// Dlansy computes the specified norm of an n×n symmetric matrix. If
// norm == MaxColumnSum or norm == lapackMaxRowSum work must have length
// at least n, otherwise work is unused.
func (impl Implementation) Dlansy(norm MatrixNorm, uplo blas.Uplo, n int, a []float64, lda int, work []float64) float64 {
	checkMatrix(n, n, a, lda)
	switch norm {
	case MaxRowSum, MaxColumnSum, NormFrob, MaxAbs:
	default:
		panic(badNorm)
	}
	if (norm == MaxColumnSum || norm == MaxRowSum) && len(work) < n {
		panic(badWork)
	}
	if uplo != blas.Upper && uplo != blas.Lower {
		panic(badUplo)
	}

	if n == 0 {
		return 0
	}
	switch norm {
	default:
		panic("unreachable")
	case MaxAbs:
		if uplo == blas.Upper {
			var max float64
			for i := 0; i < n; i++ {
				for j := i; j < n; j++ {
					v := math.Abs(a[i*lda+j])
					if math.IsNaN(v) {
						return math.NaN()
					}
					if v > max {
						max = v
					}
				}
			}
			return max
		}
		var max float64
		for i := 0; i < n; i++ {
			for j := 0; j <= i; j++ {
				v := math.Abs(a[i*lda+j])
				if math.IsNaN(v) {
					return math.NaN()
				}
				if v > max {
					max = v
				}
			}
		}
		return max
	case MaxRowSum, MaxColumnSum:
		// A symmetric matrix has the same 1-norm and ∞-norm.
		for i := 0; i < n; i++ {
			work[i] = 0
		}
		if uplo == blas.Upper {
			for i := 0; i < n; i++ {
				work[i] += math.Abs(a[i*lda+i])
				for j := i + 1; j < n; j++ {
					v := math.Abs(a[i*lda+j])
					work[i] += v
					work[j] += v
				}
			}
		} else {
			for i := 0; i < n; i++ {
				for j := 0; j < i; j++ {
					v := math.Abs(a[i*lda+j])
					work[i] += v
					work[j] += v
				}
				work[i] += math.Abs(a[i*lda+i])
			}
		}
		var max float64
		for i := 0; i < n; i++ {
			v := work[i]
			if math.IsNaN(v) {
				return math.NaN()
			}
			if v > max {
				max = v
			}
		}
		return max
	case NormFrob:
		if uplo == blas.Upper {
			var sum float64
			for i := 0; i < n; i++ {
				v := a[i*lda+i]
				sum += v * v
				for j := i + 1; j < n; j++ {
					v := a[i*lda+j]
					sum += 2 * v * v
				}
			}
			return math.Sqrt(sum)
		}
		var sum float64
		for i := 0; i < n; i++ {
			for j := 0; j < i; j++ {
				v := a[i*lda+j]
				sum += 2 * v * v
			}
			v := a[i*lda+i]
			sum += v * v
		}
		return math.Sqrt(sum)
	}
}
