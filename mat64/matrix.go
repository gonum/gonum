// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"fmt"
	"math"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
	"github.com/gonum/lapack"
)

// Matrix is the basic matrix interface type.
type Matrix interface {
	// Dims returns the dimensions of a Matrix.
	Dims() (r, c int)

	// At returns the value of a matrix element at (r, c). It will panic if r or c are
	// out of bounds for the matrix.
	At(r, c int) float64

	// T returns the transpose of the Matrix. Whether T returns a copy of the
	// underlying data is implementation dependent.
	// This method may be implemented using the Transpose type, which
	// provides an implicit matrix transpose.
	T() Matrix
}

var (
	_ Matrix       = Transpose{}
	_ Untransposer = Transpose{}
)

// Transpose is a type for performing an implicit matrix transpose. It implements
// the Matrix interface, returning values from the transpose of the matrix within.
type Transpose struct {
	Matrix Matrix
}

// At returns the value of the element at row i and column j of the transposed
// matrix, that is, row j and column i of the Matrix field.
func (t Transpose) At(i, j int) float64 {
	return t.Matrix.At(j, i)
}

// Dims returns the dimensions of the transposed matrix. The number of rows returned
// is the number of columns in the Matrix field, and the number of columns is
// the number of rows in the Matrix field.
func (t Transpose) Dims() (r, c int) {
	c, r = t.Matrix.Dims()
	return r, c
}

// T performs an implicit transpose by returning the Matrix field.
func (t Transpose) T() Matrix {
	return t.Matrix
}

// Untranspose returns the Matrix field.
func (t Transpose) Untranspose() Matrix {
	return t.Matrix
}

// Untransposer is a type that can undo an implicit transpose.
type Untransposer interface {
	// Note: This interface is needed to unify all of the Transpose types. In
	// the mat64 methods, we need to test if the Matrix has been implicitly
	// transposed. If this is checked by testing for the specific Transpose type
	// then the behavior will be different if the user uses T() or TTri() for a
	// triangular matrix.

	// Untranspose returns the underlying Matrix stored for the implicit transpose.
	Untranspose() Matrix
}

// Mutable is a matrix interface type that allows elements to be altered.
type Mutable interface {
	// Set alters the matrix element at (r, c) to v. It will panic if r or c are out of
	// bounds for the matrix.
	Set(r, c int, v float64)

	Matrix
}

// A Vectorer can return rows and columns of the represented matrix.
type Vectorer interface {
	// Row returns a []float64 for the row specified by the index i. It will
	// panic if the index is out of bounds. len(dst) must equal the number of
	// columns unless dst is nil in which case a new slice must be allocated.
	Row(dst []float64, i int) []float64

	// Col returns a []float64 for the row specified by the index i. It will
	// panic if the index is out of bounds. len(dst) must equal the number of
	// rows unless dst is nil in which case a new slice must be allocated.
	Col(dst []float64, j int) []float64
}

// A VectorSetter can set rows and columns in the represented matrix.
type VectorSetter interface {
	// SetRow sets the values in the specified rows of the matrix to the values
	// in src. len(src) must equal the number of columns in the receiver.
	SetRow(i int, src []float64)

	// SetCol sets the values in the specified column of the matrix to the values
	// in src. len(src) must equal the number of rows in the receiver.
	SetCol(i int, src []float64)
}

// A RowViewer can return a Vector reflecting a row that is backed by the matrix
// data. The Vector returned will have Len() == nCols.
type RowViewer interface {
	RowView(r int) *Vector
}

// A RawRowViewer can return a slice of float64 reflecting a row that is backed by the matrix
// data.
type RawRowViewer interface {
	RawRowView(r int) []float64
}

// A ColViewer can return a Vector reflecting a row that is backed by the matrix
// data. The Vector returned will have Len() == nRows.
type ColViewer interface {
	ColView(c int) *Vector
}

// A RawColViewer can return a slice of float64 reflecting a column that is backed by the matrix
// data.
type RawColViewer interface {
	RawColView(c int) *Vector
}

// A Cloner can make a copy of a into the receiver, overwriting the previous value of the
// receiver. The clone operation does not make any restriction on shape.
type Cloner interface {
	Clone(a Matrix)
}

// A Reseter can reset the matrix so that it can be reused as the receiver of a dimensionally
// restricted operation. This is commonly used when the matrix is being used a a workspace
// or temporary matrix.
//
// If the matrix is a view, using the reset matrix may result in data corruption in elements
// outside the view.
type Reseter interface {
	Reset()
}

// A Copier can make a copy of elements of a into the receiver. The submatrix copied
// starts at row and column 0 and has dimensions equal to the minimum dimensions of
// the two matrices. The number of row and columns copied is returned.
// Note that the behavior of Copy from a Matrix with backing data that aliases the
// receiver is undefined.
type Copier interface {
	Copy(a Matrix) (r, c int)
}

// A Viewer returns a submatrix view of the Matrix parameter, starting at row i, column j
// and extending r rows and c columns. If i or j are out of range, or r or c are zero or
// extend beyond the bounds of the matrix View will panic with ErrIndexOutOfRange. The
// returned matrix must retain the receiver's reference to the original matrix such that
// changes in the elements of the submatrix are reflected in the original and vice versa.
type Viewer interface {
	View(i, j, r, c int) Matrix
}

// A Grower can grow the size of the represented matrix by the given number of rows and columns.
// Growing beyond the size given by the Caps method will result in the allocation of a new
// matrix and copying of the elements. If Grow is called with negative increments it will
// panic with ErrIndexOutOfRange.
type Grower interface {
	Caps() (r, c int)
	Grow(r, c int) Matrix
}

// A Normer can return the specified matrix norm, o of the matrix represented by the receiver.
//
// Valid order values are:
//
//     1 - max of the sum of the absolute values of columns
//    -1 - min of the sum of the absolute values of columns
//   Inf - max of the sum of the absolute values of rows
//  -Inf - min of the sum of the absolute values of rows
//     0 - Frobenius norm
//
// Norm will panic with ErrNormOrder if an illegal norm order is specified.
type Normer interface {
	Norm(o float64) float64
}

// An Inver can calculate the inverse of the matrix represented by a and stored in the receiver.
// ErrSingular is returned if there is no inverse of the matrix.
type Inver interface {
	Inv(a Matrix) error
}

// An Adder can add the matrices represented by a and b, placing the result in the receiver. Add
// will panic if the two matrices do not have the same shape.
type Adder interface {
	Add(a, b Matrix)
}

// A Suber can subtract the matrix b from a, placing the result in the receiver. Sub will panic if
// the two matrices do not have the same shape.
type Suber interface {
	Sub(a, b Matrix)
}

// An ElemMuler can perform element-wise multiplication of the matrices represented by a and b,
// placing the result in the receiver. MulEmen will panic if the two matrices do not have the same
// shape.
type ElemMuler interface {
	MulElem(a, b Matrix)
}

// An ElemDiver can perform element-wise division a / b of the matrices represented by a and b,
// placing the result in the receiver. DivElem will panic if the two matrices do not have the same
// shape.
type ElemDiver interface {
	DivElem(a, b Matrix)
}

// A Scaler can perform scalar multiplication of the matrix represented by a with c, placing
// the result in the receiver.
type Scaler interface {
	Scale(c float64, a Matrix)
}

// A Sumer can return the sum of elements of the matrix represented by the receiver.
type Sumer interface {
	Sum() float64
}

// A Muler can determine the matrix product of a and b, placing the result in the receiver.
// If the number of columns in a does not equal the number of rows in b, Mul will panic.
type Muler interface {
	Mul(a, b Matrix)
}

// An Exper can perform a matrix exponentiation of the square matrix a. Exp will panic with ErrShape
// if a is not square.
type Exper interface {
	Exp(a Matrix)
}

// A Power can raise a square matrix, a to a positive integral power, n. Pow will panic if n is negative
// or if a is not square.
type Power interface {
	Pow(a Matrix, n int)
}

// A Dotter can determine the sum of the element-wise products of the elements of the receiver and b.
// If the shapes of the two matrices differ, Dot will panic.
type Dotter interface {
	Dot(b Matrix) float64
}

// A Stacker can create the stacked matrix of a with b, where b is placed in the greater indexed rows.
// The result of stacking is placed in the receiver, overwriting the previous value of the receiver.
// Stack will panic if the two input matrices do not have the same number of columns.
type Stacker interface {
	Stack(a, b Matrix)
}

// An Augmenter can create the augmented matrix of a with b, where b is placed in the greater indexed
// columns. The result of augmentation is placed in the receiver, overwriting the previous value of the
// receiver. Augment will panic if the two input matrices do not have the same number of rows.
type Augmenter interface {
	Augment(a, b Matrix)
}

// An Applyer can apply fn to each of the elements of the matrix represented by a, placing the
// resulting matrix in the receiver. The function fn takes a row/column index and element value
// and returns some function of that tuple.
type Applyer interface {
	Apply(fn func(r, c int, v float64) float64, a Matrix)
}

// A BandWidther represents a banded matrix and can return the left and right half-bandwidths, k1 and
// k2.
type BandWidther interface {
	BandWidth() (k1, k2 int)
}

// A RawMatrixSetter can set the underlying blas64.General used by the receiver. There is no restriction
// on the shape of the receiver. Changes to the receiver's elements will be reflected in the blas64.General.Data.
type RawMatrixSetter interface {
	SetRawMatrix(a blas64.General)
}

// A RawMatrixer can return a blas64.General representation of the receiver. Changes to the blas64.General.Data
// slice will be reflected in the original matrix, changes to the Rows, Cols and Stride fields will not.
type RawMatrixer interface {
	RawMatrix() blas64.General
}

// A RawVectorer can return a blas64.Vector representation of the receiver. Changes to the blas64.Vector.Data
// slice will be reflected in the original matrix, changes to the Inc field will not.
type RawVectorer interface {
	RawVector() blas64.Vector
}

// TODO(btracey): Consider adding CopyCol/CopyRow if the behavior seems useful.
// TODO(btracey): Add in fast paths to Row/Col for the other concrete types
// (TriDense, etc.) as well as relevant interfaces (Vectorer, RawRowViewer, etc.)

// Col copies the elements in the jth column of the matrix into the slice dst.
// The length of the provided slice must equal the number of rows, unless the
// slice is nil in which case a new slice is first allocated.
func Col(dst []float64, j int, a Matrix) []float64 {
	r, c := a.Dims()
	if j < 0 || j >= c {
		panic(ErrColAccess)
	}
	if dst == nil {
		dst = make([]float64, r)
	} else {
		if len(dst) != r {
			panic(ErrRowLength)
		}
	}
	aMat, aTrans := untranspose(a)
	if rm, ok := aMat.(RawMatrixer); ok {
		m := rm.RawMatrix()
		if aTrans {
			copy(dst, m.Data[j*m.Stride:j*m.Stride+m.Cols])
			return dst
		}
		blas64.Copy(r,
			blas64.Vector{Inc: m.Stride, Data: m.Data[j:]},
			blas64.Vector{Inc: 1, Data: dst},
		)
		return dst
	}
	for i := 0; i < r; i++ {
		dst[i] = a.At(i, j)
	}
	return dst
}

// Row copies the elements in the jth column of the matrix into the slice dst.
// The length of the provided slice must equal the number of columns, unless the
// slice is nil in which case a new slice is first allocated.
func Row(dst []float64, i int, a Matrix) []float64 {
	r, c := a.Dims()
	if i < 0 || i >= r {
		panic(ErrColAccess)
	}
	if dst == nil {
		dst = make([]float64, c)
	} else {
		if len(dst) != c {
			panic(ErrColLength)
		}
	}
	aMat, aTrans := untranspose(a)
	if rm, ok := aMat.(RawMatrixer); ok {
		m := rm.RawMatrix()
		if aTrans {
			blas64.Copy(c,
				blas64.Vector{Inc: m.Stride, Data: m.Data[i:]},
				blas64.Vector{Inc: 1, Data: dst},
			)
			return dst
		}
		copy(dst, m.Data[i*m.Stride:i*m.Stride+m.Cols])
		return dst
	}
	for j := 0; j < c; j++ {
		dst[j] = a.At(i, j)
	}
	return dst
}

// Det returns the determinant of the matrix a. In many expressions using LogDet
// will be more numerically stable.
func Det(a Matrix) float64 {
	det, sign := LogDet(a)
	return math.Exp(det) * sign
}

// Equal returns whether element a and b have the same size and contain all equal
// elements.
func Equal(a, b Matrix) bool {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || ac != bc {
		return false
	}
	aMat, aTrans := untranspose(a)
	bMat, bTrans := untranspose(b)
	if rma, ok := aMat.(RawMatrixer); ok {
		if rmb, ok := bMat.(RawMatrixer); ok {
			ra := rma.RawMatrix()
			rb := rmb.RawMatrix()
			if aTrans == bTrans {
				for i := 0; i < ra.Rows; i++ {
					for j := 0; j < ra.Cols; j++ {
						if ra.Data[i*ra.Stride+j] != rb.Data[i*rb.Stride+j] {
							return false
						}
					}
				}
				return true
			}
			for i := 0; i < ra.Rows; i++ {
				for j := 0; j < ra.Cols; j++ {
					if ra.Data[i*ra.Stride+j] != rb.Data[j*rb.Stride+i] {
						return false
					}
				}
			}
			return true
		}
	}
	if rma, ok := aMat.(RawSymmetricer); ok {
		if rmb, ok := bMat.(RawSymmetricer); ok {
			ra := rma.RawSymmetric()
			rb := rmb.RawSymmetric()
			// Symmetric matrices are always upper and equal to their transpose.
			for i := 0; i < ra.N; i++ {
				for j := i; j < ra.N; j++ {
					if ra.Data[i*ra.Stride+j] != rb.Data[i*rb.Stride+j] {
						return false
					}
				}
			}
			return true
		}
	}
	if ra, ok := aMat.(*Vector); ok {
		if rb, ok := bMat.(*Vector); ok {
			// If the raw vectors are the same length they must either both be
			// transposed or both not transposed (or have length 1).
			for i := 0; i < ra.n; i++ {
				if ra.mat.Data[i*ra.mat.Inc] != rb.mat.Data[i*rb.mat.Inc] {
					return false
				}
			}
			return true
		}
	}
	for i := 0; i < ar; i++ {
		for j := 0; j < ac; j++ {
			if a.At(i, j) != b.At(i, j) {
				return false
			}
		}
	}
	return true
}

// EqualApprox returns whether the matrices a and b have the same size and contain all equal
// elements with tolerance for element-wise equality specified by epsilon. Matrices
// with non-equal shapes are not equal.
func EqualApprox(a, b Matrix, epsilon float64) bool {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || ac != bc {
		return false
	}
	aMat, aTrans := untranspose(a)
	bMat, bTrans := untranspose(b)
	if rma, ok := aMat.(RawMatrixer); ok {
		if rmb, ok := bMat.(RawMatrixer); ok {
			ra := rma.RawMatrix()
			rb := rmb.RawMatrix()
			if aTrans == bTrans {
				for i := 0; i < ra.Rows; i++ {
					for j := 0; j < ra.Cols; j++ {
						if !floats.EqualWithinAbsOrRel(ra.Data[i*ra.Stride+j], rb.Data[i*rb.Stride+j], epsilon, epsilon) {
							return false
						}
					}
				}
				return true
			}
			for i := 0; i < ra.Rows; i++ {
				for j := 0; j < ra.Cols; j++ {
					if !floats.EqualWithinAbsOrRel(ra.Data[i*ra.Stride+j], rb.Data[j*rb.Stride+i], epsilon, epsilon) {
						return false
					}
				}
			}
			return true
		}
	}
	if rma, ok := aMat.(RawSymmetricer); ok {
		if rmb, ok := bMat.(RawSymmetricer); ok {
			ra := rma.RawSymmetric()
			rb := rmb.RawSymmetric()
			// Symmetric matrices are always upper and equal to their transpose.
			for i := 0; i < ra.N; i++ {
				for j := i; j < ra.N; j++ {
					if !floats.EqualWithinAbsOrRel(ra.Data[i*ra.Stride+j], rb.Data[i*rb.Stride+j], epsilon, epsilon) {
						return false
					}
				}
			}
			return true
		}
	}
	if ra, ok := aMat.(*Vector); ok {
		if rb, ok := bMat.(*Vector); ok {
			// If the raw vectors are the same length they must either both be
			// transposed or both not transposed (or have length 1).
			for i := 0; i < ra.n; i++ {
				if !floats.EqualWithinAbsOrRel(ra.mat.Data[i*ra.mat.Inc], rb.mat.Data[i*rb.mat.Inc], epsilon, epsilon) {
					return false
				}
			}
			return true
		}
	}
	for i := 0; i < ar; i++ {
		for j := 0; j < ac; j++ {
			if !floats.EqualWithinAbsOrRel(a.At(i, j), b.At(i, j), epsilon, epsilon) {
				return false
			}
		}
	}
	return true
}

// LogDet returns the log of the determinant and the sign of the determinant
// for the matrix that has been factorized. Numerical stability in product and
// division expressions is generally improved by working in log space.
func LogDet(a Matrix) (det float64, sign float64) {
	// TODO(btracey): Add specialized routines for TriDense, etc.
	var lu LU
	lu.Factorize(a)
	return lu.LogDet()
}

// Inverse returns the inverse or pseudoinverse of the matrix a.
// It returns a nil matrix and ErrSingular if a is singular.
func Inverse(a Matrix) (*Dense, error) {
	m, _ := a.Dims()
	d := make([]float64, m*m)
	for i := 0; i < m*m; i += m + 1 {
		d[i] = 1
	}
	eye := NewDense(m, m, d)
	x := &Dense{}
	err := x.Solve(a, eye)
	return x, err
}

// Max returns the largest element value of the matrix A.
func Max(a Matrix) float64 {
	r, c := a.Dims()
	if r == 0 || c == 0 {
		return 0
	}
	// Max(A) = Max(A^T)
	aMat, _ := untranspose(a)
	switch m := aMat.(type) {
	case RawMatrixer:
		rm := m.RawMatrix()
		max := math.Inf(-1)
		for i := 0; i < rm.Rows; i++ {
			for _, v := range rm.Data[i*rm.Stride : i*rm.Stride+rm.Cols] {
				if v > max {
					max = v
				}
			}
		}
		return max
	case RawTriangular:
		rm := m.RawTriangular()
		// The max of a triangular is at least 0 unless the size is 1.
		if rm.N == 1 {
			return rm.Data[0]
		}
		max := 0.0
		if rm.Uplo == blas.Upper {
			for i := 0; i < rm.N; i++ {
				for _, v := range rm.Data[i*rm.Stride+i : i*rm.Stride+rm.N] {
					if v > max {
						max = v
					}
				}
			}
			return max
		}
		for i := 0; i < rm.N; i++ {
			for _, v := range rm.Data[i*rm.Stride : i*rm.Stride+i+1] {
				if v > max {
					max = v
				}
			}
		}
		return max
	case RawSymmetricer:
		rm := m.RawSymmetric()
		if rm.Uplo == blas.Upper {
			max := math.Inf(-1)
			for i := 0; i < rm.N; i++ {
				for _, v := range rm.Data[i*rm.Stride+i : i*rm.Stride+rm.N] {
					if v > max {
						max = v
					}
				}
			}
			return max
		}
		max := math.Inf(-1)
		for i := 0; i < rm.N; i++ {
			for _, v := range rm.Data[i*rm.Stride : i*rm.Stride+i+1] {
				if v > max {
					max = v
				}
			}
		}
		return max
	default:
		r, c := aMat.Dims()
		max := math.Inf(-1)
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				v := aMat.At(i, j)
				if v > max {
					max = v
				}
			}
		}
		return max
	}
}

// Min returns the smallest element value of the matrix A.
func Min(a Matrix) float64 {
	r, c := a.Dims()
	if r == 0 || c == 0 {
		return 0
	}
	// Min(A) = Min(A^T)
	aMat, _ := untranspose(a)
	switch m := aMat.(type) {
	case RawMatrixer:
		rm := m.RawMatrix()
		min := math.Inf(1)
		for i := 0; i < rm.Rows; i++ {
			for _, v := range rm.Data[i*rm.Stride : i*rm.Stride+rm.Cols] {
				if v < min {
					min = v
				}
			}
		}
		return min
	case RawTriangular:
		rm := m.RawTriangular()
		// The min of a triangular is at most 0 unless the size is 1.
		if rm.N == 1 {
			return rm.Data[0]
		}
		min := 0.0
		if rm.Uplo == blas.Upper {
			for i := 0; i < rm.N; i++ {
				for _, v := range rm.Data[i*rm.Stride+i : i*rm.Stride+rm.N] {
					if v < min {
						min = v
					}
				}
			}
			return min
		}
		for i := 0; i < rm.N; i++ {
			for _, v := range rm.Data[i*rm.Stride : i*rm.Stride+i+1] {
				if v < min {
					min = v
				}
			}
		}
		return min
	case RawSymmetricer:
		rm := m.RawSymmetric()
		if rm.Uplo == blas.Upper {
			min := math.Inf(1)
			for i := 0; i < rm.N; i++ {
				for _, v := range rm.Data[i*rm.Stride+i : i*rm.Stride+rm.N] {
					if v < min {
						min = v
					}
				}
			}
			return min
		}
		min := math.Inf(1)
		for i := 0; i < rm.N; i++ {
			for _, v := range rm.Data[i*rm.Stride : i*rm.Stride+i+1] {
				if v < min {
					min = v
				}
			}
		}
		return min
	default:
		r, c := aMat.Dims()
		min := math.Inf(1)
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				v := aMat.At(i, j)
				if v < min {
					min = v
				}
			}
		}
		return min
	}
}

// Maybe will recover a panic with a type mat64.Error from fn, and return this error.
// Any other error is re-panicked.
func Maybe(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(Error); ok {
				if e.string == "" {
					panic("mat64: invalid error")
				}
				err = e
				return
			}
			panic(r)
		}
	}()
	fn()
	return
}

// MaybeFloat will recover a panic with a type mat64.Error from fn, and return this error.
// Any other error is re-panicked.
func MaybeFloat(fn func() float64) (f float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(Error); ok {
				if e.string == "" {
					panic("mat64: invalid error")
				}
				err = e
				return
			}
			panic(r)
		}
	}()
	return fn(), nil
}

// Trace returns the trace of the matrix. Trace will panic if the
// matrix is not square.
func Trace(a Matrix) float64 {
	r, c := a.Dims()
	if r != c {
		panic(ErrSquare)
	}

	aMat, _ := untranspose(a)
	switch m := aMat.(type) {
	case RawMatrixer:
		rm := m.RawMatrix()
		var t float64
		for i := 0; i < r; i++ {
			t += rm.Data[i*rm.Stride+i]
		}
		return t
	case RawTriangular:
		rm := m.RawTriangular()
		var t float64
		for i := 0; i < r; i++ {
			t += rm.Data[i*rm.Stride+i]
		}
		return t
	case RawSymmetricer:
		rm := m.RawSymmetric()
		var t float64
		for i := 0; i < r; i++ {
			t += rm.Data[i*rm.Stride+i]
		}
		return t
	default:
		var t float64
		for i := 0; i < r; i++ {
			t += a.At(i, i)
		}
		return t
	}
}

// Condition is the condition number of a matrix. The condition
// number is defined as ||A|| * ||A^-1||.
//
// One important use of Condition is during linear solve routines (finding x such
// that A * x = b). The condition number of A indicates the accuracy of
// the computed solution. A Condition error will be returned if the condition
// number of A is sufficiently large. If A is exactly singular to working precision,
// Condition == ∞, and the solve algorithm may have completed early. If Condition
// is large and finite the solve algorithm will be performed, but the computed
// solution may be innacurate. Due to the nature of finite precision arithmetic,
// the value of Condition is only an approximate test of singularity.
type Condition float64

func (c Condition) Error() string {
	return fmt.Sprintf("matrix singular or near-singular with inverse condition number %.4e", c)
}

// condTol describes the limit of the condition number. If the inverse of the
// condition number is above this value, the matrix is considered singular.
var condTol float64 = 1e16

// condNorm describes the matrix norm to use for computing the condition number.
var condNorm = lapack.MaxRowSum

// condNormTrans is the norm to compute on A^T to get the same result as computing
// condNorm on A.
var condNormTrans = lapack.MaxColumnSum

// Type Error represents matrix handling errors. These errors can be recovered by Maybe wrappers.
type Error struct{ string }

func (err Error) Error() string { return err.string }

var (
	ErrIndexOutOfRange = Error{"mat64: index out of range"}
	ErrRowAccess       = Error{"mat64: row index out of range"}
	ErrColAccess       = Error{"mat64: column index out of range"}
	ErrVectorAccess    = Error{"mat64: vector index out of range"}
	ErrZeroLength      = Error{"mat64: zero length in matrix definition"}
	ErrRowLength       = Error{"mat64: row length mismatch"}
	ErrColLength       = Error{"mat64: col length mismatch"}
	ErrSquare          = Error{"mat64: expect square matrix"}
	ErrNormOrder       = Error{"mat64: invalid norm order for matrix"}
	ErrSingular        = Error{"mat64: matrix is singular"}
	ErrShape           = Error{"mat64: dimension mismatch"}
	ErrIllegalStride   = Error{"mat64: illegal stride"}
	ErrPivot           = Error{"mat64: malformed pivot list"}
	ErrTriangle        = Error{"mat64: triangular storage mismatch"}
	ErrTriangleSet     = Error{"mat64: triangular set out of bounds"}
)

var (
	badSliceLength = "mat64: improper slice length"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// use returns a float64 slice with l elements, using f if it
// has the necessary capacity, otherwise creating a new slice.
func use(f []float64, l int) []float64 {
	if l <= cap(f) {
		return f[:l]
	}
	return make([]float64, l)
}

// useZeroed returns a float64 slice with l elements, using f if it
// has the necessary capacity, otherwise creating a new slice. The
// elements of the returned slice are guaranteed to be zero.
func useZeroed(f []float64, l int) []float64 {
	if l <= cap(f) {
		f = f[:l]
		zero(f)
		return f
	}
	return make([]float64, l)
}

// zero zeros the given slice's elements.
func zero(f []float64) {
	for i := range f {
		f[i] = 0
	}
}
