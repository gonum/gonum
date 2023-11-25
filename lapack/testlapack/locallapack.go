// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/internal/asm/f64"
	"gonum.org/v1/gonum/lapack"
)

// dlagtm is a local implementation of Dlagtm to keep code paths independent.
func dlagtm(trans blas.Transpose, m, n int, alpha float64, dl, d, du []float64, b []float64, ldb int, beta float64, c []float64, ldc int) {
	if m == 0 || n == 0 {
		return
	}

	if beta != 1 {
		if beta == 0 {
			for i := 0; i < m; i++ {
				ci := c[i*ldc : i*ldc+n]
				for j := range ci {
					ci[j] = 0
				}
			}
		} else {
			for i := 0; i < m; i++ {
				ci := c[i*ldc : i*ldc+n]
				for j := range ci {
					ci[j] *= beta
				}
			}
		}
	}

	if alpha == 0 {
		return
	}

	if m == 1 {
		if alpha == 1 {
			for j := 0; j < n; j++ {
				c[j] += d[0] * b[j]
			}
		} else {
			for j := 0; j < n; j++ {
				c[j] += alpha * d[0] * b[j]
			}
		}
		return
	}

	if trans != blas.NoTrans {
		dl, du = du, dl
	}

	if alpha == 1 {
		for j := 0; j < n; j++ {
			c[j] += d[0]*b[j] + du[0]*b[ldb+j]
		}
		for i := 1; i < m-1; i++ {
			for j := 0; j < n; j++ {
				c[i*ldc+j] += dl[i-1]*b[(i-1)*ldb+j] + d[i]*b[i*ldb+j] + du[i]*b[(i+1)*ldb+j]
			}
		}
		for j := 0; j < n; j++ {
			c[(m-1)*ldc+j] += dl[m-2]*b[(m-2)*ldb+j] + d[m-1]*b[(m-1)*ldb+j]
		}
	} else {
		for j := 0; j < n; j++ {
			c[j] += alpha * (d[0]*b[j] + du[0]*b[ldb+j])
		}
		for i := 1; i < m-1; i++ {
			for j := 0; j < n; j++ {
				c[i*ldc+j] += alpha * (dl[i-1]*b[(i-1)*ldb+j] + d[i]*b[i*ldb+j] + du[i]*b[(i+1)*ldb+j])
			}
		}
		for j := 0; j < n; j++ {
			c[(m-1)*ldc+j] += alpha * (dl[m-2]*b[(m-2)*ldb+j] + d[m-1]*b[(m-1)*ldb+j])
		}
	}
}

// dlangt is a local implementation of Dlangt to keep code paths independent.
func dlangt(norm lapack.MatrixNorm, n int, dl, d, du []float64) float64 {
	if n == 0 {
		return 0
	}

	dl = dl[:n-1]
	d = d[:n]
	du = du[:n-1]

	var anorm float64
	switch norm {
	case lapack.MaxAbs:
		for _, diag := range [][]float64{dl, d, du} {
			for _, di := range diag {
				if math.IsNaN(di) {
					return di
				}
				di = math.Abs(di)
				if di > anorm {
					anorm = di
				}
			}
		}
	case lapack.MaxColumnSum:
		if n == 1 {
			return math.Abs(d[0])
		}
		anorm = math.Abs(d[0]) + math.Abs(dl[0])
		if math.IsNaN(anorm) {
			return anorm
		}
		tmp := math.Abs(du[n-2]) + math.Abs(d[n-1])
		if math.IsNaN(tmp) {
			return tmp
		}
		if tmp > anorm {
			anorm = tmp
		}
		for i := 1; i < n-1; i++ {
			tmp = math.Abs(du[i-1]) + math.Abs(d[i]) + math.Abs(dl[i])
			if math.IsNaN(tmp) {
				return tmp
			}
			if tmp > anorm {
				anorm = tmp
			}
		}
	case lapack.MaxRowSum:
		if n == 1 {
			return math.Abs(d[0])
		}
		anorm = math.Abs(d[0]) + math.Abs(du[0])
		if math.IsNaN(anorm) {
			return anorm
		}
		tmp := math.Abs(dl[n-2]) + math.Abs(d[n-1])
		if math.IsNaN(tmp) {
			return tmp
		}
		if tmp > anorm {
			anorm = tmp
		}
		for i := 1; i < n-1; i++ {
			tmp = math.Abs(dl[i-1]) + math.Abs(d[i]) + math.Abs(du[i])
			if math.IsNaN(tmp) {
				return tmp
			}
			if tmp > anorm {
				anorm = tmp
			}
		}
	case lapack.Frobenius:
		panic("not implemented")
	default:
		panic("invalid norm")
	}
	return anorm
}

// dlansy is a local implementation of Dlansy to keep code paths independent.
func dlansy(norm lapack.MatrixNorm, uplo blas.Uplo, n int, a []float64, lda int) float64 {
	if n == 0 {
		return 0
	}
	work := make([]float64, n)
	switch norm {
	case lapack.MaxAbs:
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
	case lapack.MaxRowSum, lapack.MaxColumnSum:
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
	case lapack.Frobenius:
		panic("not implemented")
	default:
		panic("invalid norm")
	}
}

// dlange is a local implementation of Dlange to keep code paths independent.
func dlange(norm lapack.MatrixNorm, m, n int, a []float64, lda int) float64 {
	if m == 0 || n == 0 {
		return 0
	}
	var value float64
	switch norm {
	case lapack.MaxAbs:
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				value = math.Max(value, math.Abs(a[i*lda+j]))
			}
		}
	case lapack.MaxColumnSum:
		work := make([]float64, n)
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				work[j] += math.Abs(a[i*lda+j])
			}
		}
		for i := 0; i < n; i++ {
			value = math.Max(value, work[i])
		}
	case lapack.MaxRowSum:
		for i := 0; i < m; i++ {
			var sum float64
			for j := 0; j < n; j++ {
				sum += math.Abs(a[i*lda+j])
			}
			value = math.Max(value, sum)
		}
	case lapack.Frobenius:
		for i := 0; i < m; i++ {
			row := f64.L2NormUnitary(a[i*lda : i*lda+n])
			value = math.Hypot(value, row)
		}
	default:
		panic("invalid norm")
	}
	return value
}

// dlansb is a local implementation of Dlansb to keep code paths independent.
func dlansb(norm lapack.MatrixNorm, uplo blas.Uplo, n, kd int, ab []float64, ldab int, work []float64) float64 {
	if n == 0 {
		return 0
	}
	var value float64
	switch norm {
	case lapack.MaxAbs:
		if uplo == blas.Upper {
			for i := 0; i < n; i++ {
				for j := 0; j < min(n-i, kd+1); j++ {
					aij := math.Abs(ab[i*ldab+j])
					if aij > value || math.IsNaN(aij) {
						value = aij
					}
				}
			}
		} else {
			for i := 0; i < n; i++ {
				for j := max(0, kd-i); j < kd+1; j++ {
					aij := math.Abs(ab[i*ldab+j])
					if aij > value || math.IsNaN(aij) {
						value = aij
					}
				}
			}
		}
	case lapack.MaxColumnSum, lapack.MaxRowSum:
		work = work[:n]
		var sum float64
		if uplo == blas.Upper {
			for i := range work {
				work[i] = 0
			}
			for i := 0; i < n; i++ {
				sum := work[i] + math.Abs(ab[i*ldab])
				for j := i + 1; j < min(i+kd+1, n); j++ {
					aij := math.Abs(ab[i*ldab+j-i])
					sum += aij
					work[j] += aij
				}
				if sum > value || math.IsNaN(sum) {
					value = sum
				}
			}
		} else {
			for i := 0; i < n; i++ {
				sum = 0
				for j := max(0, i-kd); j < i; j++ {
					aij := math.Abs(ab[i*ldab+kd+j-i])
					sum += aij
					work[j] += aij
				}
				work[i] = sum + math.Abs(ab[i*ldab+kd])
			}
			for _, sum := range work {
				if sum > value || math.IsNaN(sum) {
					value = sum
				}
			}
		}
	case lapack.Frobenius:
		panic("not implemented")
	default:
		panic("invalid norm")
	}
	return value
}

// dlantr is a local implementation of Dlantr to keep code paths independent.
func dlantr(norm lapack.MatrixNorm, uplo blas.Uplo, diag blas.Diag, m, n int, a []float64, lda int, work []float64) float64 {
	// Quick return if possible.
	minmn := min(m, n)
	if minmn == 0 {
		return 0
	}
	switch norm {
	case lapack.MaxAbs:
		if diag == blas.Unit {
			value := 1.0
			if uplo == blas.Upper {
				for i := 0; i < m; i++ {
					for j := i + 1; j < n; j++ {
						tmp := math.Abs(a[i*lda+j])
						if math.IsNaN(tmp) {
							return tmp
						}
						if tmp > value {
							value = tmp
						}
					}
				}
				return value
			}
			for i := 1; i < m; i++ {
				for j := 0; j < min(i, n); j++ {
					tmp := math.Abs(a[i*lda+j])
					if math.IsNaN(tmp) {
						return tmp
					}
					if tmp > value {
						value = tmp
					}
				}
			}
			return value
		}
		var value float64
		if uplo == blas.Upper {
			for i := 0; i < m; i++ {
				for j := i; j < n; j++ {
					tmp := math.Abs(a[i*lda+j])
					if math.IsNaN(tmp) {
						return tmp
					}
					if tmp > value {
						value = tmp
					}
				}
			}
			return value
		}
		for i := 0; i < m; i++ {
			for j := 0; j <= min(i, n-1); j++ {
				tmp := math.Abs(a[i*lda+j])
				if math.IsNaN(tmp) {
					return tmp
				}
				if tmp > value {
					value = tmp
				}
			}
		}
		return value
	case lapack.MaxColumnSum:
		if diag == blas.Unit {
			for i := 0; i < minmn; i++ {
				work[i] = 1
			}
			for i := minmn; i < n; i++ {
				work[i] = 0
			}
			if uplo == blas.Upper {
				for i := 0; i < m; i++ {
					for j := i + 1; j < n; j++ {
						work[j] += math.Abs(a[i*lda+j])
					}
				}
			} else {
				for i := 1; i < m; i++ {
					for j := 0; j < min(i, n); j++ {
						work[j] += math.Abs(a[i*lda+j])
					}
				}
			}
		} else {
			for i := 0; i < n; i++ {
				work[i] = 0
			}
			if uplo == blas.Upper {
				for i := 0; i < m; i++ {
					for j := i; j < n; j++ {
						work[j] += math.Abs(a[i*lda+j])
					}
				}
			} else {
				for i := 0; i < m; i++ {
					for j := 0; j <= min(i, n-1); j++ {
						work[j] += math.Abs(a[i*lda+j])
					}
				}
			}
		}
		var max float64
		for _, v := range work[:n] {
			if math.IsNaN(v) {
				return math.NaN()
			}
			if v > max {
				max = v
			}
		}
		return max
	case lapack.MaxRowSum:
		var maxsum float64
		if diag == blas.Unit {
			if uplo == blas.Upper {
				for i := 0; i < m; i++ {
					var sum float64
					if i < minmn {
						sum = 1
					}
					for j := i + 1; j < n; j++ {
						sum += math.Abs(a[i*lda+j])
					}
					if math.IsNaN(sum) {
						return math.NaN()
					}
					if sum > maxsum {
						maxsum = sum
					}
				}
				return maxsum
			} else {
				for i := 0; i < m; i++ {
					var sum float64
					if i < minmn {
						sum = 1
					}
					for j := 0; j < min(i, n); j++ {
						sum += math.Abs(a[i*lda+j])
					}
					if math.IsNaN(sum) {
						return math.NaN()
					}
					if sum > maxsum {
						maxsum = sum
					}
				}
				return maxsum
			}
		} else {
			if uplo == blas.Upper {
				for i := 0; i < m; i++ {
					var sum float64
					for j := i; j < n; j++ {
						sum += math.Abs(a[i*lda+j])
					}
					if math.IsNaN(sum) {
						return sum
					}
					if sum > maxsum {
						maxsum = sum
					}
				}
				return maxsum
			} else {
				for i := 0; i < m; i++ {
					var sum float64
					for j := 0; j <= min(i, n-1); j++ {
						sum += math.Abs(a[i*lda+j])
					}
					if math.IsNaN(sum) {
						return sum
					}
					if sum > maxsum {
						maxsum = sum
					}
				}
				return maxsum
			}
		}
	case lapack.Frobenius:
		panic("not implemented")
	default:
		panic("invalid norm")
	}
}

// dlantb is a local implementation of Dlantb to keep code paths independent.
func dlantb(norm lapack.MatrixNorm, uplo blas.Uplo, diag blas.Diag, n, k int, a []float64, lda int, work []float64) float64 {
	if n == 0 {
		return 0
	}
	var value float64
	switch norm {
	case lapack.MaxAbs:
		if uplo == blas.Upper {
			var jfirst int
			if diag == blas.Unit {
				value = 1
				jfirst = 1
			}
			for i := 0; i < n; i++ {
				for _, aij := range a[i*lda+jfirst : i*lda+min(n-i, k+1)] {
					if math.IsNaN(aij) {
						return aij
					}
					aij = math.Abs(aij)
					if aij > value {
						value = aij
					}
				}
			}
		} else {
			jlast := k + 1
			if diag == blas.Unit {
				value = 1
				jlast = k
			}
			for i := 0; i < n; i++ {
				for _, aij := range a[i*lda+max(0, k-i) : i*lda+jlast] {
					if math.IsNaN(aij) {
						return math.NaN()
					}
					aij = math.Abs(aij)
					if aij > value {
						value = aij
					}
				}
			}
		}
	case lapack.MaxRowSum:
		var sum float64
		if uplo == blas.Upper {
			var jfirst int
			if diag == blas.Unit {
				jfirst = 1
			}
			for i := 0; i < n; i++ {
				sum = 0
				if diag == blas.Unit {
					sum = 1
				}
				for _, aij := range a[i*lda+jfirst : i*lda+min(n-i, k+1)] {
					sum += math.Abs(aij)
				}
				if math.IsNaN(sum) {
					return math.NaN()
				}
				if sum > value {
					value = sum
				}
			}
		} else {
			jlast := k + 1
			if diag == blas.Unit {
				jlast = k
			}
			for i := 0; i < n; i++ {
				sum = 0
				if diag == blas.Unit {
					sum = 1
				}
				for _, aij := range a[i*lda+max(0, k-i) : i*lda+jlast] {
					sum += math.Abs(aij)
				}
				if math.IsNaN(sum) {
					return math.NaN()
				}
				if sum > value {
					value = sum
				}
			}
		}
	case lapack.MaxColumnSum:
		work = work[:n]
		if diag == blas.Unit {
			for i := range work {
				work[i] = 1
			}
		} else {
			for i := range work {
				work[i] = 0
			}
		}
		if uplo == blas.Upper {
			var jfirst int
			if diag == blas.Unit {
				jfirst = 1
			}
			for i := 0; i < n; i++ {
				for j, aij := range a[i*lda+jfirst : i*lda+min(n-i, k+1)] {
					work[i+jfirst+j] += math.Abs(aij)
				}
			}
		} else {
			jlast := k + 1
			if diag == blas.Unit {
				jlast = k
			}
			for i := 0; i < n; i++ {
				off := max(0, k-i)
				for j, aij := range a[i*lda+off : i*lda+jlast] {
					work[i+j+off-k] += math.Abs(aij)
				}
			}
		}
		for _, wi := range work {
			if math.IsNaN(wi) {
				return math.NaN()
			}
			if wi > value {
				value = wi
			}
		}
	case lapack.Frobenius:
		panic("not implemented")
	default:
		panic("invalid norm")
	}
	return value
}

func dlanst(norm lapack.MatrixNorm, n int, d, e []float64) float64 {
	if n == 0 {
		return 0
	}
	var value float64
	switch norm {
	case lapack.MaxAbs:
		if n == 1 {
			value = math.Abs(d[0])
		} else {
			for _, di := range d[:n] {
				value = math.Max(value, math.Abs(di))
			}
			for _, ei := range e[:n-1] {
				value = math.Max(value, math.Abs(ei))
			}
		}
	case lapack.MaxColumnSum, lapack.MaxRowSum:
		if n == 1 {
			value = math.Abs(d[0])
		} else {
			value = math.Abs(d[0]) + math.Abs(e[0])
			value = math.Max(value, math.Abs(d[n-1])+math.Abs(e[n-2]))
			for i := 1; i < n-1; i++ {
				sum := math.Abs(d[i]) + math.Abs(e[i]) + math.Abs(e[i-1])
				value = math.Max(value, sum)
			}
		}
	case lapack.Frobenius:
		panic("not implemented")
	default:
		panic("invalid norm")
	}
	return value
}
