// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"math"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack/lapack64"
	"gonum.org/v1/gonum/matrix"
)

// LQ is a type for creating and using the LQ factorization of a matrix.
type LQ struct {
	lq   *Dense
	tau  []float64
	cond float64
}

func (lq *LQ) updateCond() {
	// A = LQ, where Q is orthonormal. Orthonormal multiplications do not change
	// the condition number. Thus, ||A|| = ||L|| ||Q|| = ||Q||.
	m := lq.lq.mat.Rows
	work := getFloats(3*m, false)
	iwork := getInts(m, false)
	l := lq.lq.asTriDense(m, blas.NonUnit, blas.Lower)
	v := lapack64.Trcon(matrix.CondNorm, l.mat, work, iwork)
	lq.cond = 1 / v
	putFloats(work)
	putInts(iwork)
}

// Factorize computes the LQ factorization of an m×n matrix a where n <= m. The LQ
// factorization always exists even if A is singular.
//
// The LQ decomposition is a factorization of the matrix A such that A = L * Q.
// The matrix Q is an orthonormal n×n matrix, and L is an m×n upper triangular matrix.
// L and Q can be extracted from the LTo and QTo methods.
func (lq *LQ) Factorize(a Matrix) {
	m, n := a.Dims()
	if m > n {
		panic(matrix.ErrShape)
	}
	k := min(m, n)
	if lq.lq == nil {
		lq.lq = &Dense{}
	}
	lq.lq.Clone(a)
	work := []float64{0}
	lq.tau = make([]float64, k)
	lapack64.Gelqf(lq.lq.mat, lq.tau, work, -1)
	work = getFloats(int(work[0]), false)
	lapack64.Gelqf(lq.lq.mat, lq.tau, work, len(work))
	putFloats(work)
	lq.updateCond()
}

// TODO(btracey): Add in the "Reduced" forms for extracting the m×m orthogonal
// and upper triangular matrices.

// LTo extracts the m×n lower trapezoidal matrix from a LQ decomposition.
// If dst is nil, a new matrix is allocated. The resulting L matrix is returned.
func (lq *LQ) LTo(dst *Dense) *Dense {
	r, c := lq.lq.Dims()
	if dst == nil {
		dst = NewDense(r, c, nil)
	} else {
		dst.reuseAs(r, c)
	}

	// Disguise the LQ as a lower triangular.
	t := &TriDense{
		mat: blas64.Triangular{
			N:      r,
			Stride: lq.lq.mat.Stride,
			Data:   lq.lq.mat.Data,
			Uplo:   blas.Lower,
			Diag:   blas.NonUnit,
		},
		cap: lq.lq.capCols,
	}
	dst.Copy(t)

	if r == c {
		return dst
	}
	// Zero right of the triangular.
	for i := 0; i < r; i++ {
		zero(dst.mat.Data[i*dst.mat.Stride+r : i*dst.mat.Stride+c])
	}

	return dst
}

// QTo extracts the n×n orthonormal matrix Q from an LQ decomposition.
// If dst is nil, a new matrix is allocated. The resulting Q matrix is returned.
func (lq *LQ) QTo(dst *Dense) *Dense {
	_, c := lq.lq.Dims()
	if dst == nil {
		dst = NewDense(c, c, nil)
	} else {
		dst.reuseAsZeroed(c, c)
	}
	q := dst.mat

	// Set Q = I.
	ldq := q.Stride
	for i := 0; i < c; i++ {
		q.Data[i*ldq+i] = 1
	}

	// Construct Q from the elementary reflectors.
	work := []float64{0}
	lapack64.Ormlq(blas.Left, blas.NoTrans, lq.lq.mat, lq.tau, q, work, -1)
	work = getFloats(int(work[0]), false)
	lapack64.Ormlq(blas.Left, blas.NoTrans, lq.lq.mat, lq.tau, q, work, len(work))
	putFloats(work)

	return dst
}

// SolveLQ finds a minimum-norm solution to a system of linear equations defined
// by the matrices A and b, where A is an m×n matrix represented in its LQ factorized
// form. If A is singular or near-singular a Condition error is returned. Please
// see the documentation for Condition for more information.
//
// The minimization problem solved depends on the input parameters.
//  If trans == false, find the minimum norm solution of A * X = b.
//  If trans == true, find X such that ||A*X - b||_2 is minimized.
// The solution matrix, X, is stored in place into the receiver.
func (m *Dense) SolveLQ(lq *LQ, trans bool, b Matrix) error {
	r, c := lq.lq.Dims()
	br, bc := b.Dims()

	// The LQ solve algorithm stores the result in-place into the right hand side.
	// The storage for the answer must be large enough to hold both b and x.
	// However, this method's receiver must be the size of x. Copy b, and then
	// copy the result into m at the end.
	if trans {
		if c != br {
			panic(matrix.ErrShape)
		}
		m.reuseAs(r, bc)
	} else {
		if r != br {
			panic(matrix.ErrShape)
		}
		m.reuseAs(c, bc)
	}
	// Do not need to worry about overlap between m and b because x has its own
	// independent storage.
	x := getWorkspace(max(r, c), bc, false)
	x.Copy(b)
	t := lq.lq.asTriDense(lq.lq.mat.Rows, blas.NonUnit, blas.Lower).mat
	if trans {
		work := []float64{0}
		lapack64.Ormlq(blas.Left, blas.NoTrans, lq.lq.mat, lq.tau, x.mat, work, -1)
		work = getFloats(int(work[0]), false)
		lapack64.Ormlq(blas.Left, blas.NoTrans, lq.lq.mat, lq.tau, x.mat, work, len(work))
		putFloats(work)

		ok := lapack64.Trtrs(blas.Trans, t, x.mat)
		if !ok {
			return matrix.Condition(math.Inf(1))
		}
	} else {
		ok := lapack64.Trtrs(blas.NoTrans, t, x.mat)
		if !ok {
			return matrix.Condition(math.Inf(1))
		}
		for i := r; i < c; i++ {
			zero(x.mat.Data[i*x.mat.Stride : i*x.mat.Stride+bc])
		}
		work := []float64{0}
		lapack64.Ormlq(blas.Left, blas.Trans, lq.lq.mat, lq.tau, x.mat, work, -1)
		work = getFloats(int(work[0]), false)
		lapack64.Ormlq(blas.Left, blas.Trans, lq.lq.mat, lq.tau, x.mat, work, len(work))
		putFloats(work)
	}
	// M was set above to be the correct size for the result.
	m.Copy(x)
	putWorkspace(x)
	if lq.cond > matrix.ConditionTolerance {
		return matrix.Condition(lq.cond)
	}
	return nil
}

// SolveLQVec finds a minimum-norm solution to a system of linear equations.
// Please see Dense.SolveLQ for the full documentation.
func (v *Vector) SolveLQVec(lq *LQ, trans bool, b *Vector) error {
	if v != b {
		v.checkOverlap(b.mat)
	}
	r, c := lq.lq.Dims()
	// The Solve implementation is non-trivial, so rather than duplicate the code,
	// instead recast the Vectors as Dense and call the matrix code.
	if trans {
		v.reuseAs(r)
	} else {
		v.reuseAs(c)
	}
	return v.asDense().SolveLQ(lq, trans, b.asDense())
}
