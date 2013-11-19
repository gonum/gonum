// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Based on the LUDecomposition class from Jama 1.0.3.

package mat64

import (
	"math"
)

type LUFactors struct {
	LU    *Dense
	Pivot []int
	Sign  int
}

// LUD performs an LU Decomposition for an m-by-n matrix a.
//
// If m >= n, the LU decomposition is an m-by-n unit lower triangular matrix L,
// an n-by-n upper triangular matrix U, and a permutation vector piv of length m
// so that A(piv,:) = L*U.
//
// If m < n, then L is m-by-m and U is m-by-n.
//
// The LU decompostion with pivoting always exists, even if the matrix is
// singular, so the LUD will never fail. The primary use of the LU decomposition
// is in the solution of square systems of simultaneous linear equations.  This
// will fail if IsSingular() returns true.
func LU(a *Dense) LUFactors {
	// Use a "left-looking", dot-product, Crout/Doolittle algorithm.
	m, n := a.Dims()
	lu := a

	piv := make([]int, m)
	for i := range piv {
		piv[i] = i
	}
	sign := 1

	var (
		luRowi = make([]float64, n)
		luColj = make([]float64, m)
	)

	// Outer loop.
	for j := 0; j < n; j++ {

		// Make a copy of the j-th column to localize references.
		for i := 0; i < m; i++ {
			luColj[i] = lu.At(i, j)
		}

		// Apply previous transformations.
		for i := 0; i < m; i++ {
			lu.Row(luRowi, i)

			// Most of the time is spent in the following dot product.
			kmax := min(i, j)
			var s float64
			for k := 0; k < kmax; k++ {
				s += luRowi[k] * luColj[k]
			}

			luColj[i] -= s
			luRowi[j] = luColj[i]

			lu.SetRow(i, luRowi)
		}

		// Find pivot and exchange if necessary.
		p := j
		for i := j + 1; i < m; i++ {
			if math.Abs(luColj[i]) > math.Abs(luColj[p]) {
				p = i
			}
		}
		if p != j {
			for k := 0; k < n; k++ {
				t := lu.At(p, k)
				lu.Set(p, k, lu.At(j, k))
				lu.Set(j, k, t)
			}
			piv[p], piv[j] = piv[j], piv[p]
			sign = -sign
		}

		// Compute multipliers.
		if j < m && lu.At(j, j) != 0 {
			for i := j + 1; i < m; i++ {
				lu.Set(i, j, lu.At(i, j)/lu.At(j, j))
			}
		}
	}

	return LUFactors{lu, piv, sign}
}

// LUGaussian performs an LU Decomposition for an m-by-n matrix a using Gaussian elimination.
// L and U are found using the "daxpy"-based elimination algorithm used in LINPACK and
// MATLAB.
//
// If m >= n, the LU decomposition is an m-by-n unit lower triangular matrix L,
// an n-by-n upper triangular matrix U, and a permutation vector piv of length m
// so that A(piv,:) = L*U.
//
// If m < n, then L is m-by-m and U is m-by-n.
//
// The LU decompostion with pivoting always exists, even if the matrix is
// singular, so the LUD will never fail. The primary use of the LU decomposition
// is in the solution of square systems of simultaneous linear equations.  This
// will fail if IsSingular() returns true.
func LUGaussian(a *Dense) LUFactors {
	// Initialize.
	m, n := a.Dims()
	lu := a

	piv := make([]int, m)
	for i := range piv {
		piv[i] = i
	}
	sign := 1

	// Main loop.
	for k := 0; k < n; k++ {
		// Find pivot.
		p := k
		for i := k + 1; i < m; i++ {
			if math.Abs(lu.At(i, k)) > math.Abs(lu.At(p, k)) {
				p = i
			}
		}

		// Exchange if necessary.
		if p != k {
			for j := 0; j < n; j++ {
				t := lu.At(p, j)
				lu.Set(p, j, lu.At(k, j))
				lu.Set(k, j, t)
			}
			piv[p], piv[k] = piv[k], piv[p]
			sign = -sign
		}

		// Compute multipliers and eliminate k-th column.
		if lu.At(k, k) != 0 {
			for i := k + 1; i < m; i++ {
				lu.Set(i, k, lu.At(i, k)/lu.At(k, k))
				for j := k + 1; j < n; j++ {
					lu.Set(i, j, lu.At(i, j)-lu.At(i, k)*lu.At(k, j))
				}
			}
		}
	}

	return LUFactors{lu, piv, sign}
}

// IsSingular returns whether the the upper triangular factor and hence a is
// singular.
func (f LUFactors) IsSingular() bool {
	lu := f.LU
	_, n := lu.Dims()
	for j := 0; j < n; j++ {
		if lu.At(j, j) == 0 {
			return true
		}
	}
	return false
}

// L returns the lower triangular factor of the LU decomposition.
func (f LUFactors) L() *Dense {
	lu := f.LU
	m, n := lu.Dims()
	l, _ := NewDense(m, n, make([]float64, m*n))
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if i > j {
				l.Set(i, j, lu.At(i, j))
			} else if i == j {
				l.Set(i, j, 1)
			}
		}
	}
	return l
}

// U returns the upper triangular factor of the LU decomposition.
func (f LUFactors) U() *Dense {
	lu := f.LU
	m, n := lu.Dims()
	u, _ := NewDense(m, n, make([]float64, m*n))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i <= j {
				u.Set(i, j, lu.At(i, j))
			}
		}
	}
	return u
}

// Det returns the determinant of matrix a decomposed into lu. The matrix
// a must have been square.
func (f LUFactors) Det() float64 {
	lu, sign := f.LU, f.Sign
	m, n := lu.Dims()
	if m != n {
		panic(ErrSquare)
	}
	d := float64(sign)
	for j := 0; j < n; j++ {
		d *= lu.At(j, j)
	}
	return d
}

// Solve computes a solution of a.x = b where b has as many rows as a. A matrix x
// is returned that minimizes the two norm of L*U*X = B(piv,:). QRSolve will panic
// if a is singular. The matrix b is overwritten during the call.
func (f LUFactors) Solve(b *Dense) (x *Dense) {
	lu, piv := f.LU, f.Pivot
	m, n := lu.Dims()
	bm, _ := b.Dims()
	if bm != m {
		panic(ErrShape)
	}
	if f.IsSingular() {
		panic("mat64: matrix is singular")
	}

	// Copy right hand side with pivoting
	nx := bm
	x = pivotRows(b, piv)

	// Solve L*Y = B(piv,:)
	for k := 0; k < n; k++ {
		for i := k + 1; i < n; i++ {
			for j := 0; j < nx; j++ {
				x.Set(i, j, x.At(i, j)-x.At(k, j)*lu.At(i, k))
			}
		}
	}

	// Solve U*X = Y;
	for k := n - 1; k >= 0; k-- {
		for j := 0; j < nx; j++ {
			x.Set(k, j, x.At(k, j)/lu.At(k, k))
		}
		for i := 0; i < k; i++ {
			for j := 0; j < nx; j++ {
				x.Set(i, j, x.At(i, j)-x.At(k, j)*lu.At(i, k))
			}
		}
	}

	return x
}

func pivotRows(a *Dense, piv []int) *Dense {
	visit := make([]bool, len(piv))
	_, n := a.Dims()
	fromRow := make([]float64, n)
	toRow := make([]float64, n)
	for to, from := range piv {
		for to != from && !visit[from] {
			visit[from], visit[to] = true, true
			a.Row(fromRow, from)
			a.Row(toRow, to)
			a.SetRow(from, toRow)
			a.SetRow(to, fromRow)
			to, from = from, piv[from]
		}
	}
	return a
}
