// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import "gonum.org/v1/gonum/blas/cblas128"

// Dense is a dense matrix representation with complex data.
type CDense struct {
	mat cblas128.General

	capRows, capCols int
}

// Dims returns the number of rows and columns in the matrix.
func (m *CDense) Dims() (r, c int) {
	return m.mat.Rows, m.mat.Cols
}

// H performs an implicit conjugate transpose by returning the receiver inside a
// Conjugate.
func (m *CDense) H() CMatrix {
	return Conjugate{m}
}

// NewCDense creates a new complex Dense matrix with r rows and c columns.
// If data == nil, a new slice is allocated for the backing slice.
// If len(data) == r*c, data is used as the backing slice, and changes to the
// elements of the returned CDense will be reflected in data.
// If neither of these is true, NewCDense will panic.
// NewCDense will panic if either r or c is zero.
//
// The data must be arranged in row-major order, i.e. the (i*c + j)-th
// element in the data slice is the {i, j}-th element in the matrix.
func NewCDense(r, c int, data []complex128) *CDense {
	if r <= 0 || c <= 0 {
		if r == 0 || c == 0 {
			panic(ErrZeroLength)
		}
		panic("mat: negative dimension")
	}
	if data != nil && r*c != len(data) {
		panic(ErrShape)
	}
	if data == nil {
		data = make([]complex128, r*c)
	}
	return &CDense{
		mat: cblas128.General{
			Rows:   r,
			Cols:   c,
			Stride: c,
			Data:   data,
		},
		capRows: r,
		capCols: c,
	}
}
