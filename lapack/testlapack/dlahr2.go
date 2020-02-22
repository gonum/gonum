// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

type Dlahr2er interface {
	Dlahr2(n, k, nb int, a []float64, lda int, tau, t []float64, ldt int, y []float64, ldy int)
}

func Dlahr2Test(t *testing.T, impl Dlahr2er) {
	const tol = 1e-14

	rnd := rand.New(rand.NewSource(1))
	for _, test := range []struct {
		n, k, nb int
	}{
		{3, 0, 3},
		{3, 1, 2},
		{3, 1, 1},

		{5, 0, 5},
		{5, 1, 4},
		{5, 1, 3},
		{5, 1, 2},
		{5, 1, 1},
		{5, 2, 3},
		{5, 2, 2},
		{5, 2, 1},
		{5, 3, 2},
		{5, 3, 1},

		{7, 3, 4},
		{7, 3, 3},
		{7, 3, 2},
		{7, 3, 1},

		{10, 0, 10},
		{10, 1, 9},
		{10, 1, 5},
		{10, 1, 1},
		{10, 5, 5},
		{10, 5, 3},
		{10, 5, 1},
	} {
		for cas := 0; cas < 100; cas++ {
			for _, extraStride := range []int{0, 1, 10} {
				n := test.n
				k := test.k
				nb := test.nb

				a := randomGeneral(n, n-k+1, n-k+1+extraStride, rnd)
				aCopy := a
				aCopy.Data = make([]float64, len(a.Data))
				copy(aCopy.Data, a.Data)
				tmat := nanTriangular(blas.Upper, nb, nb+extraStride)
				y := nanGeneral(n, nb, nb+extraStride)
				tau := nanSlice(nb)

				impl.Dlahr2(n, k, nb, a.Data, a.Stride, tau, tmat.Data, tmat.Stride, y.Data, y.Stride)

				prefix := fmt.Sprintf("Case n=%v, k=%v, nb=%v, ldex=%v", n, k, nb, extraStride)

				if !generalOutsideAllNaN(a) {
					t.Errorf("%v: out-of-range write to A\n%v", prefix, a.Data)
				}
				if !triangularOutsideAllNaN(tmat) {
					t.Errorf("%v: out-of-range write to T\n%v", prefix, tmat.Data)
				}
				if !generalOutsideAllNaN(y) {
					t.Errorf("%v: out-of-range write to Y\n%v", prefix, y.Data)
				}

				// Check that A[:k,:] and A[:,nb:] blocks were not modified.
				for i := 0; i < n; i++ {
					for j := 0; j < n-k+1; j++ {
						if i >= k && j < nb {
							continue
						}
						if a.Data[i*a.Stride+j] != aCopy.Data[i*aCopy.Stride+j] {
							t.Errorf("%v: unexpected write to A[%v,%v]", prefix, i, j)
						}
					}
				}

				// Check that all elements of tau were assigned.
				for i, v := range tau {
					if math.IsNaN(v) {
						t.Errorf("%v: tau[%v] not assigned", prefix, i)
					}
				}

				// Extract V from a.
				v := blas64.General{
					Rows:   n - k + 1,
					Cols:   nb,
					Stride: nb,
					Data:   make([]float64, (n-k+1)*nb),
				}
				for j := 0; j < v.Cols; j++ {
					v.Data[(j+1)*v.Stride+j] = 1
					for i := j + 2; i < v.Rows; i++ {
						v.Data[i*v.Stride+j] = a.Data[(i+k-1)*a.Stride+j]
					}
				}

				// VT = V.
				vt := v
				vt.Data = make([]float64, len(v.Data))
				copy(vt.Data, v.Data)
				// VT = V * T.
				blas64.Trmm(blas.Right, blas.NoTrans, 1, tmat, vt)
				// YWant = A * V * T.
				ywant := blas64.General{
					Rows:   n,
					Cols:   nb,
					Stride: nb,
					Data:   make([]float64, n*nb),
				}
				blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, aCopy, vt, 0, ywant)

				// Compare Y and YWant.
				for i := 0; i < n; i++ {
					for j := 0; j < nb; j++ {
						diff := math.Abs(ywant.Data[i*ywant.Stride+j] - y.Data[i*y.Stride+j])
						if diff > tol {
							t.Errorf("%v: unexpected Y[%v,%v], diff=%v", prefix, i, j, diff)
						}
					}
				}

				// Construct Q directly from the first nb columns of a.
				q := constructQ("QR", n-k, nb, a.Data[k*a.Stride:], a.Stride, tau)
				if resid := residualOrthogonal(q, false); resid > tol*float64(n) {
					t.Errorf("Case %v: Q is not orthogonal; resid=%v, want<=%v", prefix, resid, tol*float64(n))
				}
				// Construct Q as the product Q = I - V*T*Vᵀ.
				qwant := blas64.General{
					Rows:   n - k + 1,
					Cols:   n - k + 1,
					Stride: n - k + 1,
					Data:   make([]float64, (n-k+1)*(n-k+1)),
				}
				for i := 0; i < qwant.Rows; i++ {
					qwant.Data[i*qwant.Stride+i] = 1
				}
				blas64.Gemm(blas.NoTrans, blas.Trans, -1, vt, v, 1, qwant)
				if resid := residualOrthogonal(qwant, false); resid > tol*float64(n) {
					t.Errorf("Case %v: Q = I - V*T*Vᵀ is not orthogonal; resid=%v, want<=%v", prefix, resid, tol*float64(n))
				}

				// Compare Q and QWant. Note that since Q is
				// (n-k)×(n-k) and QWant is (n-k+1)×(n-k+1), we
				// ignore the first row and column of QWant.
				for i := 0; i < n-k; i++ {
					for j := 0; j < n-k; j++ {
						diff := math.Abs(q.Data[i*q.Stride+j] - qwant.Data[(i+1)*qwant.Stride+j+1])
						if diff > tol {
							t.Errorf("%v: unexpected Q[%v,%v], diff=%v", prefix, i, j, diff)
						}
					}
				}
			}
		}
	}
}
