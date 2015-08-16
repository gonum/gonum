// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"fmt"

	"github.com/gonum/blas/blas64"
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
	// Row returns a slice of float64 for the row specified. It will panic if the index
	// is out of bounds. If the call requires a copy and dst is not nil it will be used and
	// returned, if it is not nil the number of elements copied will be the minimum of the
	// length of the slice and the number of columns in the matrix.
	Row(dst []float64, i int) []float64

	// Col returns a slice of float64 for the column specified. It will panic if the index
	// is out of bounds. If the call requires a copy and dst is not nil it will be used and
	// returned, if it is not nil the number of elements copied will be the minimum of the
	// length of the slice and the number of rows in the matrix.
	Col(dst []float64, j int) []float64
}

// A VectorSetter can set rows and columns in the represented matrix.
type VectorSetter interface {
	// SetRow sets the values of the specified row to the values held in a slice of float64.
	// It will panic if the index is out of bounds. The number of elements copied is
	// returned and will be the minimum of the length of the slice and the number of columns
	// in the matrix.
	SetRow(i int, src []float64) int

	// SetCol sets the values of the specified column to the values held in a slice of float64.
	// It will panic if the index is out of bounds. The number of elements copied is
	// returned and will be the minimum of the length of the slice and the number of rows
	// in the matrix.
	SetCol(i int, src []float64) int
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

// A Deter can return the determinant of the represented matrix.
type Deter interface {
	Det() float64
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

// An Equaler can compare the matrices represented by b and the receiver. Matrices with non-equal shapes
// are not equal.
type Equaler interface {
	Equals(b Matrix) bool
}

// An ApproxEqualer can compare the matrices represented by b and the receiver, with tolerance for
// element-wise equailty specified by epsilon. Matrices with non-equal shapes are not equal.
type ApproxEqualer interface {
	EqualsApprox(b Matrix, epsilon float64) bool
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

// A Tracer can return the trace of the matrix represented by the receiver. Trace will panic if the
// matrix is not square.
type Tracer interface {
	Trace() float64
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

// Det returns the determinant of the matrix a.
func Det(a Matrix) float64 {
	if a, ok := a.(Deter); ok {
		return a.Det()
	}
	var lu LU
	lu.Factorize(a)
	return lu.Det()
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

// Condition is the reciprocal of the condition number of a matrix. The condition
// number is defined as ||A|| * ||A^-1||.
//
// One important use of Condition is during linear solve routines (finding x such
// that A * x = b). The condition number of A indicates the accuracy of
// the computed solution. A Condition error will be returned if the inverse condition
// number of A is sufficiently large. If A is exactly singular to working precision,
// Condition == ∞, and the solve algorithm may have completed early. If Condition
// is large and finite the solve algorithm will be performed, but the computed
// solution may be innacurate. Due to the nature of finite precision arithmetic,
// the value of Condition is only an approximate test of singularity.
type Condition float64

func (c Condition) Error() string {
	return fmt.Sprintf("matrix singular or near-singular with inverse condition number %.4e", c)
}

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
