// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mat64 provides basic linear algebra operations for float64 matrices.
//
// Note that in all interfaces that assign the result to the receiver, the receiver must
// be either the correct dimensions for the result or the zero value for the concrete type
// of the matrix. In the latter case, matrix data is allocated and stored in the receiver.
// If the matrix dimensions do not match the result, the method must panic.
package mat64

import (
	"github.com/gonum/blas"
)

// Matrix is the basic matrix interface type.
type Matrix interface {
	// Dims returns the dimensions of a Matrix.
	Dims() (r, c int)

	// At returns the value of a matrix element at (r, c). It will panic if r or c are
	// out of bounds for the matrix.
	At(r, c int) float64
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
	// is out of bounds. If the call requires a copy and row is not nil it will be used and
	// returned, if it is not nil the number of elements copied will be the minimum of the
	// length of the slice and the number of columns in the matrix.
	Row(row []float64, r int) []float64

	// Col returns a slice of float64 for the column specified. It will panic if the index
	// is out of bounds. If the call requires a copy and col is not nil it will be used and
	// returned, if it is not nil the number of elements copied will be the minimum of the
	// length of the slice and the number of rows in the matrix.
	Col(col []float64, c int) []float64
}

// A VectorSetter can set rows and columns in the represented matrix.
type VectorSetter interface {
	// SetRow sets the values of the specified row to the values held in a slice of float64.
	// It will panic if the index is out of bounds. The number of elements copied is
	// returned and will be the minimum of the length of the slice and the number of columns
	// in the matrix.
	SetRow(r int, row []float64) int

	// SetCol sets the values of the specified column to the values held in a slice of float64.
	// It will panic if the index is out of bounds. The number of elements copied is
	// returned and will be the minimum of the length of the slice and the number of rows
	// in the matrix.
	SetCol(c int, col []float64) int
}

// A Cloner can make a copy of a into the receiver, overwriting the previous value of the
// receiver. The clone operation does not make any restriction on shape.
type Cloner interface {
	Clone(a Matrix)
}

// A Copier can make a copy of elements of a into the receiver. The copy operation fills the
// submatrix in m with the values from the submatrix of a with the dimensions equal to the
// minumum of two matrices. The number of row and columns copied it returned.
type Copier interface {
	Copy(a Matrix) (r, c int)
}

// A Viewer can extract a submatrix view of of the receiver, starting at row i, column j
// and extending r rows and c columns. If i or j are out of range, or r or c extend beyond
// the bounds of the matrix View will panic with ErrIndexOutOfRange. View must retain the
// receiver's reference to the original matrix such that changes in the elements of the
// submatrix are reflected in the original and vice versa.
type Viewer interface {
	View(i, j, r, c int)
}

// A Submatrixer can extract a copy of submatrix from a into the receiver, starting at row i,
// column j and extending r rows and c columns. If i or j are out of range, or r or c extend
// beyond the bounds of the matrix Submatrix will panic with ErrIndexOutOfRange. There is no
// restriction on the shape of the receiver but changes in the elements of the submatrix must
// not be reflected in the original.
type Submatrixer interface {
	Submatrix(a Matrix, i, j, r, c int)
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

// A TransposeCopier can make a copy of the transpose the matrix represented by a, placing the elements
// into the receiver.
type TransposeCopier interface {
	TCopy(a Matrix)
}

// A Transposer can create a transposed view matrix from the represented by the receiver.
// Changes made to the returned Matrix may be reflected in the original.
type Transposer interface {
	T() Matrix
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

// An ApplyFunc takes a row/column index and element value and returns some function of that tuple.
type ApplyFunc func(r, c int, v float64) float64

// An Applyer can apply an Applyfunc f to each of the elements of the matrix represented by a,
// placing the resulting matrix in the receiver.
type Applyer interface {
	Apply(f ApplyFunc, a Matrix)
}

// A Tracer can return the trace of the matrix represented by the receiver. Trace will panic if the
// matrix is not square.
type Tracer interface {
	Trace() float64
}

// A Uer can return the upper triangular matrix of the matrix represented by a, placing the result
// in the receiver. If the concrete value of a is the receiver, the lower residue is zeroed in place.
type Uer interface {
	U(a Matrix)
}

// An Ler can return the lower triangular matrix of the matrix represented by a, placing the result
// in the receiver. If the concrete value of a is the receiver, the upper residue is zeroed in place.
type Ler interface {
	L(a Matrix)
}

// BlasMatrix represents a cblas native representation of a matrix.
type BlasMatrix struct {
	Order      blas.Order
	Rows, Cols int
	Stride     int
	Data       []float64
}

// Matrix converts a BlasMatrix to a Matrix, writing the data to the matrix represented by c. If c is a
// BlasLoader, that method will be called, otherwise the matrix must be the correct shape.
func (b BlasMatrix) Matrix(c Mutable) {
	if c, ok := c.(BlasLoader); ok {
		c.LoadBlas(b)
		return
	}
	if rows, cols := c.Dims(); rows != b.Rows || cols != b.Cols {
		panic(ErrShape)
	}
	if b.Order == blas.ColMajor {
		for col := 0; col < b.Cols; col++ {
			for row, v := range b.Data[col*b.Stride : col*b.Stride+b.Rows] {
				c.Set(row, col, v)
			}
		}
	} else if b.Order == blas.RowMajor {
		for row := 0; row < b.Rows; row++ {
			for col, v := range b.Data[row*b.Stride : row*b.Stride+b.Cols] {
				c.Set(row, col, v)
			}
		}
	} else {
		panic("matrix: illegal order")
	}
}

// A BlasLoader can directly load a BlasMatrix representation. There is no restriction on the shape of the
// receiver.
type BlasLoader interface {
	LoadBlas(a BlasMatrix)
}

// A Blasser can return a BlasMatrix representation of the receiver. Changes to the BlasMatrix.Data
// slice will be reflected in the original matrix, changes to the Rows, Cols and Stride fields will not.
type Blasser interface {
	BlasMatrix() BlasMatrix
}

// Det returns the determinant of the matrix a.
func Det(a Matrix) float64 {
	if a, ok := a.(Deter); ok {
		return a.Det()
	}
	return LU(DenseCopyOf(a)).Det()
}

// Inverse returns the inverse or pseudoinverse of the matrix a.
func Inverse(a Matrix) *Dense {
	m, _ := a.Dims()
	d := make([]float64, m*m)
	for i := 0; i < m*m; i += m + 1 {
		d[i] = 1
	}
	eye, _ := NewDense(m, m, d)
	return Solve(a, eye)
}

// Solve returns a matrix x that satisfies ax = b.
func Solve(a, b Matrix) (x *Dense) {
	m, n := a.Dims()
	if m == n {
		return LU(DenseCopyOf(a)).Solve(DenseCopyOf(b))
	}
	return QR(DenseCopyOf(a)).Solve(DenseCopyOf(b))
}

// A Panicker is a function that may panic.
type Panicker func()

// Maybe will recover a panic with a type matrix.Error from fn, and return this error.
// Any other error is re-panicked.
func Maybe(fn Panicker) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(Error); ok {
				return
			}
			panic(r)
		}
	}()
	fn()
	return
}

// A FloatPanicker is a function that returns a float64 and may panic.
type FloatPanicker func() float64

// MaybeFloat will recover a panic with a type matrix.Error from fn, and return this error.
// Any other error is re-panicked.
func MaybeFloat(fn FloatPanicker) (f float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(Error); ok {
				err = e
				return
			}
			panic(r)
		}
	}()
	return fn(), nil
}

// Must can be used to wrap a function returning an error.
// If the returned error is not nil, Must will panic.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Type Error represents matrix package errors. These errors can be recovered by Maybe wrappers.
type Error string

func (err Error) Error() string { return string(err) }

const (
	ErrIndexOutOfRange = Error("mat64: index out of range")
	ErrZeroLength      = Error("mat64: zero length in matrix definition")
	ErrRowLength       = Error("mat64: row length mismatch")
	ErrColLength       = Error("mat64: col length mismatch")
	ErrSquare          = Error("mat64: expect square matrix")
	ErrNormOrder       = Error("mat64: invalid norm order for matrix")
	ErrSingular        = Error("mat64: matrix is singular")
	ErrShape           = Error("mat64: dimension mismatch")
	ErrIllegalStride   = Error("mat64: illegal stride")
	ErrPivot           = Error("mat64: malformed pivot list")
	ErrIllegalOrder    = Error("mat64: illegal order")
	ErrNoEngine        = Error("mat64: no blas engine registered: call Register()")
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

func realloc(f []float64, l int) []float64 {
	if l < cap(f) {
		return f[:l]
	}
	return make([]float64, l)
}
