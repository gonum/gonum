// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/lapack/lapack64"
)

// Solve solves a minimum-norm solution to a system of linear equations defined
// by the matrices a and b. If a is singular or near-singular a Condition error
// is returned. Please see the documentation for Condition for more information.
//
// The minimization problem solved depends on the input parameters.
//  1. If m >= n and trans == false, find X such that ||a*X - b||_2 is minimized.
//  2. If m < n and trans == false, find the minimum norm solution of a * X = b.
//  3. If m >= n and trans == true, find the minimum norm solution of a^T * X = b.
//  4. If m < n and trans == true, find X such that ||a*X - b||_2 is minimized.
// The solution matrix, X, is stored in place into the receiver.
func (m *Dense) Solve(a, b Matrix) error {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br {
		panic(ErrShape)
	}
	m.reuseAs(ac, bc)

	// TODO(btracey): Add special cases for SymDense, etc.
	aMat, aTrans := untranspose(a)
	bMat, bTrans := untranspose(b)
	switch rma := aMat.(type) {
	case RawTriangular:
		side := blas.Left
		tA := blas.NoTrans
		if aTrans {
			tA = blas.Trans
		}

		switch rm := bMat.(type) {
		case RawMatrixer:
			if m != bMat || bTrans {
				if m == bMat || m.checkOverlap(rm.RawMatrix()) {
					tmp := getWorkspace(br, bc, false)
					tmp.Copy(b)
					m.Copy(tmp)
					putWorkspace(tmp)
					break
				}
				m.Copy(b)
			}
		default:
			if m != bMat {
				m.Copy(b)
			} else if bTrans {
				// m and b share data so Copy cannot be used directly.
				tmp := getWorkspace(br, bc, false)
				tmp.Copy(b)
				m.Copy(tmp)
				putWorkspace(tmp)
			}
		}

		rm := rma.RawTriangular()
		blas64.Trsm(side, tA, 1, rm, m.mat)
		work := make([]float64, 3*rm.N)
		iwork := make([]int, rm.N)
		cond := lapack64.Trcon(condNorm, rm, work, iwork)
		if cond > condTol {
			return Condition(cond)
		}
		return nil
	}

	switch {
	case ar == ac:
		if a == b {
			// x = I.
			if ar == 1 {
				m.mat.Data[0] = 1
				return nil
			}
			for i := 0; i < ar; i++ {
				v := m.mat.Data[i*m.mat.Stride : i*m.mat.Stride+ac]
				zero(v)
				v[i] = 1
			}
			return nil
		}
		var lu LU
		lu.Factorize(a)
		return m.SolveLU(&lu, false, b)
	case ar > ac:
		var qr QR
		qr.Factorize(a)
		return m.SolveQR(&qr, false, b)
	default:
		var lq LQ
		lq.Factorize(a)
		return m.SolveLQ(&lq, false, b)
	}
}

// SolveVec solves a minimum-norm solution to a system of linear equations defined
// by the matrices A and B. If A is singular or near-singular a Condition error
// is returned. Please see the documentation for more information.
func (v *Vector) SolveVec(a Matrix, b *Vector) error {
	_, c := a.Dims()
	// The Solve implementation is non-trivial, so rather than duplicate the code,
	// instead recast the Vectors as Dense and call the matrix code.
	v.reuseAs(c)
	m := vecAsDense(v)
	// We conditionally create bm as m when b and v are identical
	// to prevent the overlap detection code from identifying m
	// and bm as overlapping but not identical.
	bm := m
	if v != b {
		bm = vecAsDense(b)
	}
	return m.Solve(a, bm)
}
