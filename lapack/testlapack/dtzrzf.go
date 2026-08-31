// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/floats"
)

type Dtzrzfer interface {
	Dlarfger
	Dtzrzf(m, n int, a []float64, lda int, tau, work []float64, lwork int)
}

func DtzrzfTest(t *testing.T, impl Dtzrzfer) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, test := range []struct {
		m, n, lda int
	}{
		{1, 1, 0},
		{1, 5, 0},
		{2, 2, 0},
		{2, 5, 0},
		{3, 5, 0},
		{3, 10, 0},
		{5, 5, 0},
		{5, 10, 0},
		{5, 20, 0},
		{10, 20, 0},
		{3, 10, 15},
		{5, 10, 15},
	} {
		m := test.m
		n := test.n
		lda := test.lda
		if lda == 0 {
			lda = n
		}

		aOrig := make([]float64, m*lda)
		for i := 0; i < m; i++ {
			for j := i; j < n; j++ {
				aOrig[i*lda+j] = rnd.NormFloat64()
			}
		}

		a := make([]float64, len(aOrig))
		copy(a, aOrig)
		tau := make([]float64, m)

		// Workspace query.
		work := make([]float64, 1)
		impl.Dtzrzf(m, n, nil, lda, nil, work, -1)
		lwork := int(work[0])
		work = make([]float64, lwork)

		impl.Dtzrzf(m, n, a, lda, tau, work, lwork)

		// Check R is upper triangular.
		for i := 0; i < m; i++ {
			for j := 0; j < i; j++ {
				if a[i*lda+j] != 0 {
					t.Errorf("m=%d n=%d: R not upper triangular at (%d,%d)", m, n, i, j)
				}
			}
		}

		// Construct Z = Z(0)*Z(1)*...*Z(m-1) by applying Z(i) from the left.
		z := blas64.General{Rows: n, Cols: n, Stride: n, Data: make([]float64, n*n)}
		for i := 0; i < n; i++ {
			z.Data[i*n+i] = 1
		}
		for i := m - 1; i >= 0; i-- {
			if tau[i] == 0 {
				continue
			}
			v := make([]float64, n)
			v[i] = 1
			for j := m; j < n; j++ {
				v[j] = a[i*lda+j]
			}
			// Z = (I - tau*v*v^T) * Z
			// tmp = v^T * Z  (length n)
			tmp := make([]float64, n)
			for c := 0; c < n; c++ {
				for r := 0; r < n; r++ {
					tmp[c] += v[r] * z.Data[r*n+c]
				}
			}
			// Z -= tau * v * tmp^T
			for r := 0; r < n; r++ {
				for c := 0; c < n; c++ {
					z.Data[r*n+c] -= tau[i] * v[r] * tmp[c]
				}
			}
		}

		// Verify Z is orthogonal: Z*Z^T = I.
		var zzT blas64.General
		zzT.Rows, zzT.Cols, zzT.Stride = n, n, n
		zzT.Data = make([]float64, n*n)
		blas64.Gemm(blas.NoTrans, blas.Trans, 1, z, z, 0, zzT)
		eye := make([]float64, n*n)
		for i := 0; i < n; i++ {
			eye[i*n+i] = 1
		}
		if !floats.EqualApprox(zzT.Data, eye, 1e-13) {
			t.Errorf("m=%d n=%d: Z is not orthogonal", m, n)
		}

		// Verify [R 0]*Z = A_orig.
		rz := blas64.General{Rows: m, Cols: n, Stride: n, Data: make([]float64, m*n)}
		rMat := blas64.General{Rows: m, Cols: n, Stride: n, Data: make([]float64, m*n)}
		for i := 0; i < m; i++ {
			for j := i; j < m; j++ {
				rMat.Data[i*n+j] = a[i*lda+j]
			}
		}
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, rMat, z, 0, rz)
		aOrigFlat := make([]float64, m*n)
		for i := 0; i < m; i++ {
			copy(aOrigFlat[i*n:i*n+n], aOrig[i*lda:i*lda+n])
		}
		if !floats.EqualApprox(rz.Data, aOrigFlat, 1e-13) {
			t.Errorf("m=%d n=%d: [R 0]*Z != A_orig", m, n)
		}
	}
}
