// Copyright ©2017 The Gonum Authors. All rights reserved.
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

type Dlags2er interface {
	Dlags2(upper bool, a1, a2, a3, b1, b2, b3 float64) (csu, snu, csv, snv, csq, snq float64)
}

func Dlags2Test(t *testing.T, impl Dlags2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, upper := range []bool{true, false} {
		for i := 0; i < 100; i++ {
			// Generate randomly the elements of a 2×2 matrix A
			//  [ a1 a2 ] or [ a1 0  ]
			//  [ 0  a3 ]    [ a2 a3 ]
			a1 := rnd.Float64()
			a2 := rnd.Float64()
			a3 := rnd.Float64()
			// Generate randomly the elements of a 2×2 matrix B.
			//  [ b1 b2 ] or [ b1 0  ]
			//  [ 0  b3 ]    [ b2 b3 ]
			b1 := rnd.Float64()
			b2 := rnd.Float64()
			b3 := rnd.Float64()

			// Compute orthogal matrices U, V, Q
			//  U = [  csu  snu ], V = [  csv snv ], Q = [  csq   snq ]
			//      [ -snu  csu ]      [ -snv csv ]      [ -snq   csq ]
			// that transform A and B.
			csu, snu, csv, snv, csq, snq := impl.Dlags2(upper, a1, a2, a3, b1, b2, b3)

			// Check that U, V, Q are orthogonal matrices (their
			// determinant is equal to 1).
			detU := det2x2(csu, snu, -snu, csu)
			if !floats.EqualWithinAbsOrRel(math.Abs(detU), 1, 1e-14, 1e-14) {
				t.Errorf("U not orthogonal: det(U)=%v", detU)
			}
			detV := det2x2(csv, snv, -snv, csv)
			if !floats.EqualWithinAbsOrRel(math.Abs(detV), 1, 1e-14, 1e-14) {
				t.Errorf("V not orthogonal: det(V)=%v", detV)
			}
			detQ := det2x2(csq, snq, -snq, csq)
			if !floats.EqualWithinAbsOrRel(math.Abs(detQ), 1, 1e-14, 1e-14) {
				t.Errorf("Q not orthogonal: det(Q)=%v", detQ)
			}

			// Create U, V, Q explicitly as dense matrices.
			u := blas64.General{
				Rows:   2,
				Cols:   2,
				Stride: 2,
				Data:   []float64{csu, snu, -snu, csu},
			}
			v := blas64.General{
				Rows:   2,
				Cols:   2,
				Stride: 2,
				Data:   []float64{csv, snv, -snv, csv},
			}
			q := blas64.General{
				Rows:   2,
				Cols:   2,
				Stride: 2,
				Data:   []float64{csq, snq, -snq, csq},
			}

			// Create A and B explicitly as dense matrices.
			a := blas64.General{Rows: 2, Cols: 2, Stride: 2}
			b := blas64.General{Rows: 2, Cols: 2, Stride: 2}
			if upper {
				a.Data = []float64{a1, a2, 0, a3}
				b.Data = []float64{b1, b2, 0, b3}
			} else {
				a.Data = []float64{a1, 0, a2, a3}
				b.Data = []float64{b1, 0, b2, b3}
			}

			tmp := blas64.General{Rows: 2, Cols: 2, Stride: 2, Data: make([]float64, 4)}
			// Transform A as U^T*A*Q.
			blas64.Gemm(blas.Trans, blas.NoTrans, 1, u, a, 0, tmp)
			blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, tmp, q, 0, a)
			// Transform B as V^T*A*Q.
			blas64.Gemm(blas.Trans, blas.NoTrans, 1, v, b, 0, tmp)
			blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, tmp, q, 0, b)

			// Extract elements of transformed A and B that should be equal to zero.
			var gotA, gotB float64
			if upper {
				gotA = a.Data[1]
				gotB = b.Data[1]
			} else {
				gotA = a.Data[2]
				gotB = b.Data[2]
			}
			// Check that they are indeed zero.
			if !floats.EqualWithinAbsOrRel(gotA, 0, 1e-14, 1e-14) {
				t.Errorf("unexpected non-zero value for zero triangle of U^T*A*Q: %v", gotA)
			}
			if !floats.EqualWithinAbsOrRel(gotB, 0, 1e-14, 1e-14) {
				t.Errorf("unexpected non-zero value for zero triangle of V^T*B*Q: %v", gotB)
			}
		}
	}
}

func det2x2(a, b, c, d float64) float64 { return a*d - b*c }
