// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Based on the LUDecomposition class from Jama 1.0.3.

package mat64

import (
	"math"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/lapack/lapack64"
)

// LU is a type for creating and using the LU factorization of a matrix.
type LU struct {
	lu    *Dense
	pivot []int
}

// Factorize computes the LU factorization of the square matrix a and stores the
// result. The LU decomposition will complete regardless of the singularity of a.
//
// The LU factorization is computed with pivoting, and so really the decomposition
// is a PLU decomposition where P is a permutation matrix. The individual matrix
// factors can be extracted from the factorization using the Permutation method
// on Dense, and the LFrom and UFrom methods on TriDense.
func (lu *LU) Factorize(a Matrix) {
	r, c := a.Dims()
	if r != c {
		panic(ErrSquare)
	}
	if lu.lu == nil {
		lu.lu = &Dense{}
	}
	lu.lu.Clone(a)
	if cap(lu.pivot) < r {
		lu.pivot = make([]int, r)
	}
	lu.pivot = lu.pivot[:r]
	lapack64.Getrf(lu.lu.mat, lu.pivot)
}

// Det returns the determinant of the matrix that has been factorized.
func (lu *LU) Det() float64 {
	_, n := lu.lu.Dims()
	det := 1.0
	for i := 0; i < n; i++ {
		det *= lu.lu.at(i, i)
	}
	return det
}

// Pivot returns pivot indices that enable the construction of the permutation
// matrix P (see Dense.Permutation). If pivot == nil, then new memory will be
// allocated, otherwise the length of the input must be equal to the size of the
// factorized matrix.
func (lu *LU) Pivot(swaps []int) []int {
	_, n := lu.lu.Dims()
	if swaps == nil {
		swaps = make([]int, n)
	}
	if len(swaps) != n {
		panic(badSliceLength)
	}
	// Perform the inverse of the row swaps in order to find the final
	// row swap position.
	for i := range swaps {
		swaps[i] = i
	}
	for i := n - 1; i >= 0; i-- {
		v := lu.pivot[i]
		swaps[i], swaps[v] = swaps[v], swaps[i]
	}
	return swaps
}

// RankOne updates an LU factorization as if a rank-one update had been applied to
// the original matrix A, storing the result into the receiver. That is, if in
// the original LU decomposition P * L * U = A, in the updated decomposition
// P * L * U = A + alpha * x^T * y.
func (lu *LU) RankOne(orig *LU, alpha float64, x, y *Vector) {
	// RankOne uses algorithm a1 on page 28 of "Multiple-Rank Updates to Matrix
	// Factorizations for Nonlinear Analysis and Circuit Design" by Linzhong Deng.
	// http://web.stanford.edu/group/SOL/dissertations/Linzhong-Deng-thesis.pdf
	_, n := orig.lu.Dims()
	if x.Len() != n {
		panic(ErrShape)
	}
	if y.Len() != n {
		panic(ErrShape)
	}
	if orig != lu {
		if len(lu.pivot) == 0 {
			// lu is zero
			lu.pivot = make([]int, n)
			lu.lu = NewDense(n, n, nil)
		} else {
			if len(lu.pivot) != n {
				panic(ErrShape)
			}
		}
		copy(lu.pivot, orig.pivot)
		lu.lu.Copy(orig.lu)
	}

	xs := make([]float64, n)
	ys := make([]float64, n)
	for i := 0; i < n; i++ {
		xs[i] = x.at(i)
		ys[i] = y.at(i)
	}

	// Adjust for the pivoting in the LU factorization
	for i, v := range lu.pivot {
		xs[i], xs[v] = xs[v], xs[i]
	}

	lum := lu.lu.mat
	omega := alpha
	for j := 0; j < n; j++ {
		ujj := lum.Data[j*lum.Stride+j]
		ys[j] /= ujj
		theta := 1 + xs[j]*ys[j]*omega
		beta := omega * ys[j] / theta
		gamma := omega * xs[j]
		omega -= beta * gamma
		lum.Data[j*lum.Stride+j] *= theta
		for i := j + 1; i < n; i++ {
			xs[i] -= lum.Data[i*lum.Stride+j] * xs[j]
			tmp := ys[i]
			ys[i] -= lum.Data[j*lum.Stride+i] * ys[j]
			lum.Data[i*lum.Stride+j] += beta * xs[i]
			lum.Data[j*lum.Stride+i] += gamma * tmp
		}
	}
}

// LFromLU extracts the lower triangular matrix from an LU factorization.
func (t *TriDense) LFromLU(lu *LU) {
	_, n := lu.lu.Dims()
	t.reuseAs(n, blas.Lower)
	// Extract the lower triangular elements.
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			t.mat.Data[i*t.mat.Stride+j] = lu.lu.mat.Data[i*lu.lu.mat.Stride+j]
		}
	}
	// Set ones on the diagonal.
	for i := 0; i < n; i++ {
		t.mat.Data[i*t.mat.Stride+i] = 1
	}
}

// UFromLU extracts the upper triangular matrix from an LU factorization.
func (t *TriDense) UFromLU(lu *LU) {
	_, n := lu.lu.Dims()
	t.reuseAs(n, blas.Upper)
	// Extract the upper triangular elements.
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			t.mat.Data[i*t.mat.Stride+j] = lu.lu.mat.Data[i*lu.lu.mat.Stride+j]
		}
	}
}

// Permutation constructs an r×r permutation matrix with the given row swaps.
// A permutation matrix has exactly one element equal to one in each row and column
// and all other elements equal to zero. swaps[i] specifies the row with which
// i will be swapped, which is equivalent to the non-zero column of row i.
func (m *Dense) Permutation(r int, swaps []int) {
	m.reuseAs(r, r)
	for i := 0; i < r; i++ {
		zero(m.mat.Data[i*m.mat.Stride : i*m.mat.Stride+r])
		v := swaps[i]
		if v < 0 || v >= r {
			panic(ErrRowAccess)
		}
		m.mat.Data[i*m.mat.Stride+v] = 1
	}
}

// SolveLU solves a system of linear equations using the LU decomposition of a matrix.
// It computes
//  A * x = b if trans == false
//  A^T * x = b if trans == true
// In both cases, A is represeneted in LU factorized form, and the matrix x is
// stored into the receiver.
//
// If A is singular or near-singular a Condition error is returned. Please see
// the documentation for Condition for more information.
func (m *Dense) SolveLU(lu *LU, trans bool, b Matrix) error {
	_, n := lu.lu.Dims()
	br, bc := b.Dims()
	if br != n {
		panic(ErrShape)
	}
	// TODO(btracey): Should test the condition number instead of testing that
	// the determinant is exactly zero.
	if lu.Det() == 0 {
		return Condition(math.Inf(1))
	}

	m.reuseAs(n, bc)
	bMat, _ := untranspose(b)
	var restore func()
	if m == bMat {
		m, restore = m.isolatedWorkspace(bMat)
		defer restore()
	}
	m.Copy(b)
	t := blas.NoTrans
	if trans {
		t = blas.Trans
	}
	lapack64.Getrs(t, lu.lu.mat, m.mat, lu.pivot)
	return nil
}

// SolveLUVec solves a system of linear equations using the LU decomposition of a matrix.
// It computes
//  A * x = b if trans == false
//  A^T * x = b if trans == true
// In both cases, A is represeneted in LU factorized form, and the matrix x is
// stored into the receiver.
//
// If A is singular or near-singular a Condition error is returned. Please see
// the documentation for Condition for more information.
func (v *Vector) SolveLUVec(lu *LU, trans bool, b *Vector) error {
	_, n := lu.lu.Dims()
	bn := b.Len()
	if bn != n {
		panic(ErrShape)
	}
	// TODO(btracey): Should test the condition number instead of testing that
	// the determinant is exactly zero.
	if lu.Det() == 0 {
		return Condition(math.Inf(1))
	}

	v.reuseAs(n)
	var restore func()
	if v == b {
		v, restore = v.isolatedWorkspace(b)
		defer restore()
	}
	v.CopyVec(b)
	vMat := blas64.General{
		Rows:   n,
		Cols:   1,
		Stride: v.mat.Inc,
		Data:   v.mat.Data,
	}
	t := blas.NoTrans
	if trans {
		t = blas.Trans
	}
	lapack64.Getrs(t, lu.lu.mat, vMat, lu.pivot)
	return nil
}
