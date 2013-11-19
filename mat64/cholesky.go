// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Based on the CholeskyDecomposition class from Jama 1.0.3.

package mat64

import (
	"math"
)

type CholeskyFactor struct {
	L   *Dense
	SPD bool
}

// CholeskyL returns the left Cholesky decomposition of the matrix a and whether
// the matrix is symmetric or positive definite, the returned matrix l is a lower
// triangular matrix such that a = l.l'.
func Cholesky(a *Dense) CholeskyFactor {
	// Initialize.
	m, n := a.Dims()
	spd := m == n
	l, _ := NewDense(n, n, make([]float64, n*n))

	// Main loop.
	lRowj := make([]float64, n)
	lRowk := make([]float64, n)
	for j := 0; j < n; j++ {
		var d float64
		l.Row(lRowj, j)
		for k := 0; k < j; k++ {
			var s float64
			l.Row(lRowk, k)
			for i := 0; i < k; i++ {
				s += lRowk[i] * lRowj[i]
			}
			s = (a.At(j, k) - s) / l.At(k, k)
			lRowj[k] = s
			d += s * s
			spd = spd && a.At(k, j) == a.At(j, k)
		}
		l.SetRow(j, lRowj)
		d = a.At(j, j) - d
		spd = spd && d > 0
		l.Set(j, j, math.Sqrt(math.Max(d, 0)))
		for k := j + 1; k < n; k++ {
			l.Set(j, k, 0)
		}
	}

	return CholeskyFactor{L: l, SPD: spd}
}

// CholeskyR returns the right Cholesky decomposition of the matrix a and whether
// the matrix is symmetric or positive definite, the returned matrix r is an upper
// triangular matrix such that a = r'.r.
func CholeskyR(a *Dense) (r *Dense, spd bool) {
	// Initialize.
	m, n := a.Dims()
	spd = m == n
	r, _ = NewDense(n, n, make([]float64, n*n))

	// Main loop.
	for j := 0; j < n; j++ {
		var d float64
		for k := 0; k < j; k++ {
			s := a.At(k, j)
			for i := 0; i < k; i++ {
				s -= r.At(i, k) * r.At(i, j)
			}
			s /= r.At(k, k)
			r.Set(k, j, s)
			d += s * s
			spd = spd && a.At(k, j) == a.At(j, k)
		}
		d = a.At(j, j) - d
		spd = spd && d > 0
		r.Set(j, j, math.Sqrt(math.Max(d, 0)))
		for k := j + 1; k < n; k++ {
			r.Set(k, j, 0)
		}
	}

	return r, spd
}

// CholeskySolve returns a matrix x that solves a.x = b where a = l.l'. The matrix b must
// have the same number of rows as a, and a must be symmetric and positive definite. The
// matrix b is overwritten by the operation.
func (f CholeskyFactor) Solve(b *Dense) (x *Dense) {
	if !f.SPD {
		panic("mat64: matrix not symmetric positive definite")
	}
	l := f.L

	_, n := l.Dims()
	_, bn := b.Dims()
	if n != bn {
		panic(ErrShape)
	}

	nx := bn
	x = b

	// Solve L*Y = B;
	for k := 0; k < n; k++ {
		for j := 0; j < nx; j++ {
			for i := 0; i < k; i++ {
				x.Set(k, j, x.At(k, j)-x.At(i, j)*l.At(k, i))
			}
			x.Set(k, j, x.At(k, j)/l.At(k, k))
		}
	}

	// Solve L'*X = Y;
	for k := n - 1; k >= 0; k-- {
		for j := 0; j < nx; j++ {
			for i := k + 1; i < n; i++ {
				x.Set(k, j, x.At(k, j)-x.At(i, j)*l.At(i, k))
			}
			x.Set(k, j, x.At(k, j)/l.At(k, k))
		}
	}

	return x
}
