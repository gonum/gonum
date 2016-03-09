// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Based on the CholeskyDecomposition class from Jama 1.0.3.

package mat64

import (
	"math"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/lapack/lapack64"
	"github.com/gonum/matrix"
)

const badTriangle = "mat64: invalid triangle"

// Cholesky is a type for creating and using the Cholesky factorization of a
// symmetric positive definite matrix.
type Cholesky struct {
	chol *TriDense
	cond float64
}

// updateCond updates the condition number of the Cholesky decomposition. If
// norm > 0, then that norm is used as the norm of the original matrix A, otherwise
// the norm is estimated from the decomposition.
func (c *Cholesky) updateCond(norm float64) {
	n := c.chol.mat.N
	work := make([]float64, 3*n)
	if norm < 0 {
		// This is an approximation. By the definition of a norm, ||AB|| <= ||A|| ||B||.
		// Here, A = U^T * U.
		// The condition number is ||A|| || A^-1||, so this will underestimate
		// the condition number somewhat.
		// The norm of the original factorized matrix cannot be stored because of
		// update possibilities.
		unorm := lapack64.Lantr(matrix.CondNorm, c.chol.mat, work)
		lnorm := lapack64.Lantr(matrix.CondNormTrans, c.chol.mat, work)
		norm = unorm * lnorm
	}
	sym := c.chol.asSymBlas()
	iwork := make([]int, n)
	v := lapack64.Pocon(sym, norm, work, iwork)
	c.cond = 1 / v
}

// Factorize calculates the Cholesky decomposition of the matrix A and returns
// whether the matrix is positive definite.
func (c *Cholesky) Factorize(a Symmetric) (ok bool) {
	n := a.Symmetric()
	if c.chol == nil {
		c.chol = NewTriDense(n, true, nil)
	} else {
		c.chol = NewTriDense(n, true, use(c.chol.mat.Data, n*n))
	}
	copySymIntoTriangle(c.chol, a)

	sym := c.chol.asSymBlas()
	work := make([]float64, c.chol.mat.N)
	norm := lapack64.Lansy(matrix.CondNorm, sym, work)
	_, ok = lapack64.Potrf(sym)
	if ok {
		c.updateCond(norm)
	} else {
		c.cond = math.Inf(1)
	}
	return ok
}

// Det returns the determinant of the matrix that has been factorized.
func (c *Cholesky) Det() float64 {
	return math.Exp(c.LogDet())
}

// Size returns the dimension of the factorized matrix.
func (c *Cholesky) Size() int {
	return c.chol.mat.N
}

// LogDet returns the log of the determinant of the matrix that has been factorized.
func (c *Cholesky) LogDet() float64 {
	var det float64
	for i := 0; i < c.chol.mat.N; i++ {
		det += 2 * math.Log(c.chol.mat.Data[i*c.chol.mat.Stride+i])
	}
	return det
}

// SolveCholesky finds the matrix m that solves A * m = b where A is represented
// by the cholesky decomposition, placing the result in the receiver.
func (m *Dense) SolveCholesky(chol *Cholesky, b Matrix) error {
	n := chol.chol.mat.N
	bm, bn := b.Dims()
	if n != bm {
		panic(matrix.ErrShape)
	}

	m.reuseAs(bm, bn)
	if b != m {
		m.Copy(b)
	}
	blas64.Trsm(blas.Left, blas.Trans, 1, chol.chol.mat, m.mat)
	blas64.Trsm(blas.Left, blas.NoTrans, 1, chol.chol.mat, m.mat)
	if chol.cond > matrix.ConditionTolerance {
		return matrix.Condition(chol.cond)
	}
	return nil
}

// SolveCholeskyVec finds the vector v that solves A * v = b where A is represented
// by the Cholesky decomposition, placing the result in the receiver.
func (v *Vector) SolveCholeskyVec(chol *Cholesky, b *Vector) error {
	n := chol.chol.mat.N
	vn := b.Len()
	if vn != n {
		panic(matrix.ErrShape)
	}
	if v != b {
		v.checkOverlap(b.mat)
	}
	v.reuseAs(n)
	if v != b {
		v.CopyVec(b)
	}
	blas64.Trsv(blas.Trans, chol.chol.mat, v.mat)
	blas64.Trsv(blas.NoTrans, chol.chol.mat, v.mat)
	if chol.cond > matrix.ConditionTolerance {
		return matrix.Condition(chol.cond)
	}
	return nil

}

// UFromCholesky extracts the n×n upper triangular matrix U from a Choleksy
// decomposition
//  A = U^T * U.
func (t *TriDense) UFromCholesky(chol *Cholesky) {
	n := chol.chol.mat.N
	t.reuseAs(n, true)
	t.Copy(chol.chol)
}

// LFromCholesky extracts the n×n lower triangular matrix U from a Choleksy
// decomposition
//  A = L * L^T.
func (t *TriDense) LFromCholesky(chol *Cholesky) {
	n := chol.chol.mat.N
	t.reuseAs(n, false)
	t.Copy(chol.chol.TTri())
}

// FromCholesky reconstructs the original positive definite matrix given its
// Cholesky decomposition.
func (s *SymDense) FromCholesky(chol *Cholesky) {
	n := chol.chol.mat.N
	s.reuseAs(n)
	s.SymOuterK(1, chol.chol.T())
}

// InverseCholesky computes the inverse of the matrix represented by its Cholesky
// factorization and stores the result into the receiver. If the factorized
// matrix is ill-conditioned, a Condition error will be returned.
// Note that matrix inversion is numerically unstable, and should generally be
// avoided where possible, for example by using the Solve routines.
func (s *SymDense) InverseCholesky(chol *Cholesky) error {
	// TODO(btracey): Replace this code with a direct call to Dpotri when it
	// is available.
	s.reuseAs(chol.chol.mat.N)
	// If:
	//  chol(A) = U^T * U
	// Then:
	//  chol(A^-1) = S * S^T
	// where S = U^-1
	var t TriDense
	err := t.InverseTri(chol.chol)
	s.SymOuterK(1, &t)
	return err
}
