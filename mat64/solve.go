// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"math"

	"github.com/gonum/blas"
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
	// TODO(btracey): Add a test for the condition number of A.
	// TODO(btracey): Add special cases for TriDense, SymDense, etc.
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
		// Solve using an LU decomposition.
		var lu LU
		lu.Factorize(a)
		if lu.Det() == 0 {
			return Condition(math.Inf(1))
		}
		bMat, bTrans := untranspose(b)
		if m == bMat && bTrans {
			var restore func()
			m, restore = m.isolatedWorkspace(bMat)
			defer restore()
		}
		if m != bMat {
			m.Copy(b)
		}
		lapack64.Getrs(blas.NoTrans, lu.lu.mat, m.mat, lu.pivot)
		return nil
	default:
		// Solve using QR/LQ.

		// Copy a since the corresponding LAPACK argument is modified during
		// the call.
		var aCopy Dense
		aCopy.Clone(a)

		x := getWorkspace(max(ar, ac), bc, false)
		defer putWorkspace(x)
		x.Copy(b)

		work := make([]float64, 1)
		lapack64.Gels(blas.NoTrans, aCopy.mat, x.mat, work, -1)
		work = make([]float64, int(work[0]))
		ok := lapack64.Gels(blas.NoTrans, aCopy.mat, x.mat, work, len(work))
		if !ok {
			return Condition(math.Inf(1))
		}
		m.Copy(x)
		return nil
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
	bm := vecAsDense(b)
	return m.Solve(a, bm)
}
