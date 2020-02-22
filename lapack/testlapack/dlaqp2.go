// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dlaqp2er interface {
	Dlapmter
	Dlaqp2(m, n, offset int, a []float64, lda int, jpvt []int, tau, vn1, vn2, work []float64)
}

func Dlaqp2Test(t *testing.T, impl Dlaqp2er) {
	const tol = 1e-14

	rnd := rand.New(rand.NewSource(1))
	for ti, test := range []struct {
		m, n, offset int
	}{
		{m: 4, n: 3, offset: 0},
		{m: 4, n: 3, offset: 2},
		{m: 4, n: 3, offset: 4},
		{m: 3, n: 4, offset: 0},
		{m: 3, n: 4, offset: 1},
		{m: 3, n: 4, offset: 2},
		{m: 8, n: 3, offset: 0},
		{m: 8, n: 3, offset: 4},
		{m: 8, n: 3, offset: 8},
		{m: 3, n: 8, offset: 0},
		{m: 3, n: 8, offset: 1},
		{m: 3, n: 8, offset: 2},
		{m: 10, n: 10, offset: 0},
		{m: 10, n: 10, offset: 5},
		{m: 10, n: 10, offset: 10},
	} {
		m := test.m
		n := test.n
		jpiv := make([]int, n)

		for _, extra := range []int{0, 11} {
			a := randomGeneral(m, n, n+extra, rnd)
			aCopy := cloneGeneral(a)

			for j := range jpiv {
				jpiv[j] = j
			}

			tau := make([]float64, n)
			vn1 := columnNorms(m, n, a.Data, a.Stride)
			vn2 := columnNorms(m, n, a.Data, a.Stride)
			work := make([]float64, n)

			impl.Dlaqp2(m, n, test.offset, a.Data, a.Stride, jpiv, tau, vn1, vn2, work)

			prefix := fmt.Sprintf("Case %v (offset=%d,m=%v,n=%v,extra=%v)", ti, test.offset, m, n, extra)
			if !generalOutsideAllNaN(a) {
				t.Errorf("%v: out-of-range write to A", prefix)
			}

			mo := m - test.offset
			if mo == 0 {
				continue
			}
			q := constructQ("QR", mo, n, a.Data[test.offset*a.Stride:], a.Stride, tau)

			// Check that Q is orthogonal.
			if resid := residualOrthogonal(q, false); resid > tol {
				t.Errorf("%v: Q is not orthogonal; resid=%v, want<=%v", prefix, resid, tol)
			}

			// Check that |A*P - Q*R| is small.
			impl.Dlapmt(true, aCopy.Rows, aCopy.Cols, aCopy.Data, aCopy.Stride, jpiv)
			qrap := blas64.General{
				Rows:   mo,
				Cols:   aCopy.Cols,
				Stride: aCopy.Stride,
				Data:   aCopy.Data[test.offset*aCopy.Stride:],
			}
			r := zeros(mo, n, n)
			for i := 0; i < mo; i++ {
				for j := i; j < n; j++ {
					r.Data[i*r.Stride+j] = a.Data[(test.offset+i)*a.Stride+j]
				}
			}
			blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, q, r, -1, qrap)
			resid := dlange(lapack.MaxColumnSum, qrap.Rows, qrap.Cols, qrap.Data, qrap.Stride)
			if resid > tol {
				t.Errorf("%v: |Q*R - A*P|=%v, want<=%v", prefix, resid, tol)
			}
		}
	}
}
