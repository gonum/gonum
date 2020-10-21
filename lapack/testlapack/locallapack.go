// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/lapack"
)

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
