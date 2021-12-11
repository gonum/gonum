// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import (
	"unsafe"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
)

const (
	badDim = "bad matrix dimensions"
	badIdx = "bad matrix index"
)

// Mat represents a 3×3 matrix. Useful for rotation matrices and such.
type Mat struct {
	data *[3][3]float64
}

var _ mat.Matrix = (*Mat)(nil)

// NewMat returns a new 3×3 matrix Mat type and populates its elements
// with values passed as argument in row-major form. If val argument
// is nil then NewMat returns a matrix filled with zeros.
func NewMat(val []float64) *Mat {
	if len(val) != 9 {
		if val == nil {
			return &Mat{data: new([3][3]float64)}
		}
		panic(badDim)
	}
	m := Mat{}
	m.setBackingSlice(val)
	return &m
}

// Dims returns the number of rows and columns of this matrix.
// This method will always return 3×3 for a Mat.
func (m *Mat) Dims() (r, c int) { return 3, 3 }

// At returns the value of a matrix element at row i, column j.
// At expects indices in the range [0,2].
// It will panic if i or j are out of bounds for the matrix.
func (m *Mat) At(i, j int) float64 {
	return m.data[i][j]
}

// Set sets the element at row i, column j to the value v.
func (m *Mat) Set(i, j int, v float64) {
	m.data[i][j] = v
}

// T returns the transpose of Mat. Changes in the receiver will be reflected in the returned matrix.
func (m *Mat) T() mat.Matrix { return mat.Transpose{Matrix: m} }

// RawMatrix returns the blas representation of the matrix with the backing data of this matrix.
// Changes to the returned matrix will be reflected in the receiver.
func (m *Mat) RawMatrix() blas64.General {
	return blas64.General{Rows: 3, Cols: 3, Data: m.backingSlice(), Stride: 3}
}

// Eye returns the 3×3 Identity matrix
func Eye() *Mat {
	return &Mat{data: &[3][3]float64{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}}
}

// Scale multiplies the elements of a by f, placing the result in the receiver.
//
// See the mat.Scaler interface for more information.
func (m *Mat) Scale(f float64, a mat.Matrix) {
	r, c := a.Dims()
	if r != 3 || c != 3 {
		panic(badDim)
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			m.Set(i, j, f*a.At(i, j))
		}
	}
}

// Performs matrix multiplication on v:
//  result = M * v
func (m *Mat) MulVec(v Vec) Vec {
	return Vec{
		X: v.X*m.At(0, 0) + v.Y*m.At(0, 1) + v.Z*m.At(0, 2),
		Y: v.X*m.At(1, 0) + v.Y*m.At(1, 1) + v.Z*m.At(1, 2),
		Z: v.X*m.At(2, 0) + v.Y*m.At(2, 1) + v.Z*m.At(2, 2),
	}
}

// Performs transposed matrix multiplication on v:
//  result = Mᵀ * v
func (m *Mat) MulVecTrans(v Vec) Vec {
	return Vec{
		X: v.X*m.At(0, 0) + v.Y*m.At(1, 0) + v.Z*m.At(2, 0),
		Y: v.X*m.At(0, 1) + v.Y*m.At(1, 1) + v.Z*m.At(2, 1),
		Z: v.X*m.At(0, 2) + v.Y*m.At(1, 2) + v.Z*m.At(2, 2),
	}
}

// Skew returns the 3×3 skew symmetric matrix (right hand system) of v.
//                  ⎡ 0 -z  y⎤
//  Skew({x,y,z}) = ⎢ z  0 -x⎥
//                  ⎣-y  x  0⎦
func Skew(v Vec) (M *Mat) {
	return &Mat{data: &[3][3]float64{
		{0, -v.Z, v.Y},
		{v.Z, 0, -v.X},
		{-v.Y, v.X, 0},
	}}
}

// Mul takes the matrix product of a and b, placing the result in the receiver.
// If the number of columns in a does not equal 3, Mul will panic.
func (m *Mat) Mul(a, b mat.Matrix) {
	ra, ca := a.Dims()
	rb, cb := b.Dims()
	switch {
	case ra != 3:
		panic(badDim)
	case cb != 3:
		panic(badDim)
	case ca != rb:
		panic(badDim)
	}
	if ca != 3 {
		// General matrix multiplication for the case where the inner dimension is not 3.
		t := mat.NewDense(3, 3, m.backingSlice())
		t.Mul(a, b)
		return
	}

	a00 := a.At(0, 0)
	b00 := b.At(0, 0)
	a01 := a.At(0, 1)
	b01 := b.At(0, 1)
	a02 := a.At(0, 2)
	b02 := b.At(0, 2)
	a10 := a.At(1, 0)
	b10 := b.At(1, 0)
	a11 := a.At(1, 1)
	b11 := b.At(1, 1)
	a12 := a.At(1, 2)
	b12 := b.At(1, 2)
	a20 := a.At(2, 0)
	b20 := b.At(2, 0)
	a21 := a.At(2, 1)
	b21 := b.At(2, 1)
	a22 := a.At(2, 2)
	b22 := b.At(2, 2)
	m.data[0][0] = a00*b00 + a01*b10 + a02*b20
	m.data[0][1] = a00*b01 + a01*b11 + a02*b21
	m.data[0][2] = a00*b02 + a01*b12 + a02*b22
	m.data[1][0] = a10*b00 + a11*b10 + a12*b20
	m.data[1][1] = a10*b01 + a11*b11 + a12*b21
	m.data[1][2] = a10*b02 + a11*b12 + a12*b22
	m.data[2][0] = a20*b00 + a21*b10 + a22*b20
	m.data[2][1] = a20*b01 + a21*b11 + a22*b21
	m.data[2][2] = a20*b02 + a21*b12 + a22*b22
}

// CloneFrom makes a copy of a into the receiver m.
// Mat expects a 3×3 input matrix.
func (m *Mat) CloneFrom(a mat.Matrix) {
	r, c := a.Dims()
	if r != 3 || c != 3 {
		panic(badDim)
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			m.Set(i, j, a.At(i, j))
		}
	}
}

// Sub subtracts the matrix b from a, placing the result in the receiver.
// Sub will panic if the two matrices do not have the same shape.
func (m *Mat) Sub(a, b mat.Matrix) {
	if r, c := a.Dims(); r != 3 || c != 3 {
		panic(badDim)
	}
	if r, c := b.Dims(); r != 3 || c != 3 {
		panic(badDim)
	}

	m.data[0][0] = a.At(0, 0) - b.At(0, 0)
	m.data[0][1] = a.At(0, 1) - b.At(0, 1)
	m.data[0][2] = a.At(0, 2) - b.At(0, 2)
	m.data[1][0] = a.At(1, 0) - b.At(1, 0)
	m.data[1][1] = a.At(1, 1) - b.At(1, 1)
	m.data[1][2] = a.At(1, 2) - b.At(1, 2)
	m.data[2][0] = a.At(2, 0) - b.At(2, 0)
	m.data[2][1] = a.At(2, 1) - b.At(2, 1)
	m.data[2][2] = a.At(2, 2) - b.At(2, 2)
}

// Add adds a and b element-wise, placing the result in the receiver. Add will panic if the two matrices do not have the same shape.
func (m *Mat) Add(a, b mat.Matrix) {
	if r, c := a.Dims(); r != 3 || c != 3 {
		panic(badDim)
	}
	if r, c := b.Dims(); r != 3 || c != 3 {
		panic(badDim)
	}

	m.data[0][0] = a.At(0, 0) + b.At(0, 0)
	m.data[0][1] = a.At(0, 1) + b.At(0, 1)
	m.data[0][2] = a.At(0, 2) + b.At(0, 2)
	m.data[1][0] = a.At(1, 0) + b.At(1, 0)
	m.data[1][1] = a.At(1, 1) + b.At(1, 1)
	m.data[1][2] = a.At(1, 2) + b.At(1, 2)
	m.data[2][0] = a.At(2, 0) + b.At(2, 0)
	m.data[2][1] = a.At(2, 1) + b.At(2, 1)
	m.data[2][2] = a.At(2, 2) + b.At(2, 2)
}

// VecRow returns the elements in the ith row of the receiver.
func (m *Mat) VecRow(i int) Vec {
	if i > 2 {
		panic(badIdx)
	}
	return Vec{X: m.At(i, 0), Y: m.At(i, 1), Z: m.At(i, 2)}
}

// VecCol returns the elements in the jth column of the receiver.
func (m *Mat) VecCol(j int) Vec {
	if j > 2 {
		panic(badIdx)
	}
	return Vec{X: m.At(0, j), Y: m.At(1, j), Z: m.At(2, j)}
}

// setBackingSlice requires unsafe.
func (m *Mat) setBackingSlice(vals []float64) {
	m.data = (*[3][3]float64)(unsafe.Pointer(&vals[0]))
}

// backingSlice requires unsafe.
func (m *Mat) backingSlice() []float64 {
	return (*[9]float64)(unsafe.Pointer(m.data))[:]
}
