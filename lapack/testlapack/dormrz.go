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

type Dormrzer interface {
	Dtzrzfer
	Dormrz(side blas.Side, trans blas.Transpose, m, n, k, l int, a []float64, lda int, tau, c []float64, ldc int, work []float64, lwork int)
}

func DormrzTest(t *testing.T, impl Dormrzer) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, side := range []blas.Side{blas.Left, blas.Right} {
		for _, trans := range []blas.Transpose{blas.NoTrans, blas.Trans} {
			for _, test := range []struct {
				k, mc, nc int
			}{
				{1, 4, 5},
				{2, 5, 5},
				{3, 8, 6},
				{2, 6, 6},
				{3, 10, 8},
				{1, 1, 5},
				{1, 5, 1},
				{4, 4, 6},
				{4, 6, 4},
			} {
				k := test.k
				mc := test.mc
				nc := test.nc

				var nq int
				if side == blas.Left {
					nq = mc
				} else {
					nq = nc
				}
				if k > nq {
					continue
				}
				l := nq - k

				// Build upper trapezoidal matrix of size k×nq and factor with Dtzrzf.
				ldaz := nq
				az := make([]float64, k*ldaz)
				for i := 0; i < k; i++ {
					for j := i; j < nq; j++ {
						az[i*ldaz+j] = rnd.NormFloat64()
					}
				}
				tau := make([]float64, k)
				work := make([]float64, 1)
				impl.Dtzrzf(k, nq, nil, ldaz, nil, work, -1)
				lwork := max(int(work[0]), max(mc, nc))
				work = make([]float64, lwork)
				impl.Dtzrzf(k, nq, az, ldaz, tau, work, k)

				// Construct Z explicitly.
				z := blas64.General{Rows: nq, Cols: nq, Stride: nq, Data: make([]float64, nq*nq)}
				for i := 0; i < nq; i++ {
					z.Data[i*nq+i] = 1
				}
				for i := k - 1; i >= 0; i-- {
					if tau[i] == 0 {
						continue
					}
					v := make([]float64, nq)
					v[i] = 1
					for j := k; j < nq; j++ {
						v[j] = az[i*ldaz+j]
					}
					tmp := make([]float64, nq)
					for c := 0; c < nq; c++ {
						for r := 0; r < nq; r++ {
							tmp[c] += v[r] * z.Data[r*nq+c]
						}
					}
					for r := 0; r < nq; r++ {
						for c := 0; c < nq; c++ {
							z.Data[r*nq+c] -= tau[i] * v[r] * tmp[c]
						}
					}
				}

				// Generate random C.
				ldc := nc
				c := make([]float64, mc*ldc)
				for i := range c {
					c[i] = rnd.NormFloat64()
				}
				cOrig := make([]float64, len(c))
				copy(cOrig, c)

				// Compute expected result using Dgemm.
				cMat := blas64.General{Rows: mc, Cols: nc, Stride: ldc, Data: make([]float64, len(c))}
				copy(cMat.Data, cOrig)
				cExp := blas64.General{Rows: mc, Cols: nc, Stride: ldc, Data: make([]float64, len(c))}
				switch {
				case side == blas.Left && trans == blas.NoTrans:
					blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, z, cMat, 0, cExp)
				case side == blas.Left && trans == blas.Trans:
					blas64.Gemm(blas.Trans, blas.NoTrans, 1, z, cMat, 0, cExp)
				case side == blas.Right && trans == blas.NoTrans:
					blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, cMat, z, 0, cExp)
				case side == blas.Right && trans == blas.Trans:
					blas64.Gemm(blas.NoTrans, blas.Trans, 1, cMat, z, 0, cExp)
				}

				// Apply Dormrz.
				copy(c, cOrig)
				impl.Dormrz(side, trans, mc, nc, k, l, az, ldaz, tau, c, ldc, work, lwork)

				if !floats.EqualApprox(c, cExp.Data, 1e-13) {
					t.Errorf("side=%c trans=%c k=%d l=%d mc=%d nc=%d: result mismatch",
						byte(side), byte(trans), k, l, mc, nc)
				}
			}
		}
	}
}
