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
// the matrix is symmetric and positive definite. The returned matrix l is a lower
// triangular matrix such that a = l.l'.
func Cholesky(a *Dense) CholeskyFactor {
	// Initialize.
	m, n := a.Dims()
	spd := m == n
	l := NewDense(n, n, nil)

	// Main loop.
	for j := 0; j < n; j++ {
		lRowj := l.RawRowView(j)
		var d float64
		for k := 0; k < j; k++ {
			var s float64
			for i, v := range l.RawRowView(k)[:k] {
				s += v * lRowj[i]
			}
			s = (a.at(j, k) - s) / l.at(k, k)
			lRowj[k] = s
			d += s * s
			spd = spd && a.at(k, j) == a.at(j, k)
		}
		d = a.at(j, j) - d
		spd = spd && d > 0
		l.set(j, j, math.Sqrt(math.Max(d, 0)))
		for k := j + 1; k < n; k++ {
			l.set(j, k, 0)
		}
	}

	return CholeskyFactor{L: l, SPD: spd}
}

// CholeskySolve returns a matrix x that solves a.x = b where a = l.l'. The matrix b must
// have the same number of rows as a, and a must be symmetric and positive definite. The
// matrix b is overwritten by the operation.
func (f CholeskyFactor) Solve(b *Dense) (x *Dense) {
	if !f.SPD {
		panic("mat64: matrix not symmetric positive definite")
	}
	l := f.L

	m, n := l.Dims()
	bm, bn := b.Dims()
	if m != bm {
		panic(ErrShape)
	}

	nx := bn
	x = b

	// Solve L*Y = B;
	for k := 0; k < n; k++ {
		for j := 0; j < nx; j++ {
			for i := 0; i < k; i++ {
				x.set(k, j, x.at(k, j)-x.at(i, j)*l.at(k, i))
			}
			x.set(k, j, x.at(k, j)/l.at(k, k))
		}
	}

	// Solve L'*X = Y;
	for k := n - 1; k >= 0; k-- {
		for j := 0; j < nx; j++ {
			for i := k + 1; i < n; i++ {
				x.set(k, j, x.at(k, j)-x.at(i, j)*l.at(i, k))
			}
			x.set(k, j, x.at(k, j)/l.at(k, k))
		}
	}

	return x
}
