// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/floats"
)

type Dlarfger interface {
	Dlarfg(n int, alpha float64, x []float64, incX int) (beta, tau float64)
}

func DlarfgTest(t *testing.T, impl Dlarfger) {
	const tol = 1e-14
	rnd := rand.New(rand.NewSource(1))
	for i, test := range []struct {
		alpha float64
		n     int
		x     []float64
	}{
		{
			alpha: 4,
			n:     3,
		},
		{
			alpha: -2,
			n:     3,
		},
		{
			alpha: 0,
			n:     3,
		},
		{
			alpha: 1,
			n:     1,
		},
		{
			alpha: 1,
			n:     4,
			x:     []float64{4, 5, 6},
		},
		{
			alpha: 1,
			n:     4,
			x:     []float64{0, 0, 0},
		},
		{
			alpha: dlamchS,
			n:     4,
			x:     []float64{dlamchS, dlamchS, dlamchS},
		},
	} {
		n := test.n
		incX := 1
		var x []float64
		if test.x == nil {
			x = make([]float64, n-1)
			for i := range x {
				x[i] = rnd.Float64()
			}
		} else {
			if len(test.x) != n-1 {
				panic("bad test")
			}
			x = make([]float64, n-1)
			copy(x, test.x)
		}
		xcopy := make([]float64, n-1)
		copy(xcopy, x)
		alpha := test.alpha
		beta, tau := impl.Dlarfg(n, alpha, x, incX)

		// Verify the returns and the values in v. Construct h and perform
		// the explicit multiplication.
		h := make([]float64, n*n)
		for i := 0; i < n; i++ {
			h[i*n+i] = 1
		}
		hmat := blas64.General{
			Rows:   n,
			Cols:   n,
			Stride: n,
			Data:   h,
		}
		v := make([]float64, n)
		copy(v[1:], x)
		v[0] = 1
		vVec := blas64.Vector{
			Inc:  1,
			Data: v,
		}
		blas64.Ger(-tau, vVec, vVec, hmat)
		eye := blas64.General{
			Rows:   n,
			Cols:   n,
			Stride: n,
			Data:   make([]float64, n*n),
		}
		blas64.Gemm(blas.Trans, blas.NoTrans, 1, hmat, hmat, 0, eye)
		dist := distFromIdentity(n, eye.Data, n)
		if dist > tol {
			t.Errorf("H^T * H is not close to I, dist=%v", dist)
		}

		xVec := blas64.Vector{
			Inc:  1,
			Data: make([]float64, n),
		}
		xVec.Data[0] = test.alpha
		copy(xVec.Data[1:], xcopy)

		ans := make([]float64, n)
		ansVec := blas64.Vector{
			Inc:  1,
			Data: ans,
		}
		blas64.Gemv(blas.NoTrans, 1, hmat, xVec, 0, ansVec)
		if math.Abs(ans[0]-beta) > tol {
			t.Errorf("Case %v, beta mismatch. Want %v, got %v", i, ans[0], beta)
		}
		if floats.Norm(ans[1:n], math.Inf(1)) > tol {
			t.Errorf("Case %v, nonzero answer %v", i, ans[1:n])
		}
	}
}
