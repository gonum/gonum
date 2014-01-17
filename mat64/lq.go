// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"fmt"
	"github.com/gonum/blas"
	"math"
)

type LQFactor struct {
	LQ    *Dense
	lDiag []float64
}

// LQ computes a LQ Decomposition for an m-by-n matrix a with m <= n by Householder
// reflections, the LQ decomposition is an m-by-n orthogonal matrix q and an n-by-n
// upper triangular matrix r so that a = q.r. LQ will panic with ErrShape if m > n.
//
// The LQ decomposition always exists, even if the matrix does not have full rank,
// so LQ will never fail unless m > n. The primary use of the LQ decomposition is
// in the least squares solution of non-square systems of simultaneous linear equations.
// This will fail if LQIsFullRank() returns false. The matrix a is overwritten by the
// decomposition.
func LQ(a *Dense) LQFactor {
	// Initialize.
	m, n := a.Dims()
	if m > n {
		panic(ErrShape)
	}

	lq := &Dense{}
	*lq = *a

	lDiag := make([]float64, m)
	projs := make(Vec, m)

	// Main loop.
	for k := 0; k < m; k++ {
		hh := Vec(lq.RowView(k))[k:]
		norm := blasEngine.Dnrm2(len(hh), hh, 1)
		lDiag[k] = norm

		hhNorm := (norm * math.Sqrt(1-hh[0]/norm))
		if norm != 0 && hhNorm != 0 {
			// Form k-th Householder vector.
			s := 1 / hhNorm
			hh[0] -= norm
			blasEngine.Dscal(len(hh), s, hh, 1)

			fmt.Println("hh", hh)

			// Apply transformation to remaining columns.
			if k < m-1 {
				*a = *lq
				a.View(k+1, k, m-k-1, n-k)
				projs = projs[0 : m-k-1]
				projs.Mul(a, &hh)

				for j := 0; j < m-k-1; j++ {
					dst := a.RowView(j)
					blasEngine.Daxpy(len(dst), -projs[j], hh, 1, dst, 1)
				}
			}
		}
	}

	return LQFactor{lq, lDiag}
}

// IsFullRank returns whether the L matrix and hence a has full rank.
func (f LQFactor) IsFullRank() bool {
	for _, v := range f.lDiag {
		if v == 0 {
			return false
		}
	}
	return true
}

// L returns the lower triangular factor for the LQ decomposition.
func (f LQFactor) L() *Dense {
	lq, lDiag := f.LQ, f.lDiag
	m, _ := lq.Dims()
	l := NewDense(m, m, nil)
	for i, v := range lDiag {
		for j := 0; j < m; j++ {
			if i < j {
				l.Set(j, i, lq.At(j, i))
			} else if i == j {
				l.Set(j, i, v)
			}
		}
	}
	return l
}

// replaces x with Q.x
func (f LQFactor) ApplyQ(x *Dense, trans bool) {
	nh, nc := f.LQ.Dims()
	m, n := x.Dims()
	if m != nc {
		panic(ErrShape)
	}
	proj := make([]float64, n)

	if trans {
		for k := 0; k < nh; k++ {
			sub := &Dense{}
			*sub = *x
			hh := f.LQ.RowView(k)[k:]

			sub.View(k, 0, m-k, n)

			blasEngine.Dgemv(blas.ColMajor, blas.NoTrans, n, m-k, 1,
				sub.mat.Data, sub.mat.Stride, hh, 1, 0, proj, 1)
			for i := k; i < m; i++ {
				row := x.RowView(i)
				blasEngine.Daxpy(n, -hh[i-k], proj, 1, row, 1)
			}
		}
	} else {
		for k := nh - 1; k >= 0; k-- {
			sub := &Dense{}
			*sub = *x
			hh := f.LQ.RowView(k)[k:]

			sub.View(k, 0, m-k, n)
			ms, ns := sub.Dims()

			blasEngine.Dgemv(blas.ColMajor, blas.NoTrans, ns, ms, 1,
				sub.mat.Data, sub.mat.Stride, hh, 1, 0, proj, 1)
			for i := k; i < m; i++ {
				row := x.RowView(k)
				blasEngine.Daxpy(n, hh[i-k], proj, 1, row, 1)
			}
		}
	}
}

// Solve a computes minimum norm least squares solution of a.x = b where b has as many rows as a.
// A matrix x is returned that minimizes the two norm of Q*R*X-B. Solve will panic
// if a is not full rank. The matrix b is overwritten during the call.
func (f LQFactor) Solve(b *Dense) (x *Dense) {
	lq := f.LQ
	lDiag := f.lDiag
	m, n := lq.Dims()
	bm, bn := b.Dims()
	if bm != m {
		panic(ErrShape)
	}
	if !f.IsFullRank() {
		panic("mat64: matrix is rank deficient")
	}

	x = NewDense(n, bn, nil)
	tau := make([]float64, m)
	for i := range tau {
		tau[i] = lq.At(i, i)
		lq.Set(i, i, lDiag[i])
	}
	blasEngine.Dtrsm(blas.RowMajor, blas.Right, blas.Lower, blas.NoTrans, blas.NonUnit,
		bm, bn, 1, lq.mat.Data, lq.mat.Stride, b.mat.Data, b.mat.Stride)
	for i := range tau {
		lq.Set(i, i, tau[i])
	}
	f.ApplyQ(x, true)

	return x
}
