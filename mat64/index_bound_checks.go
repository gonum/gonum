// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file must be kept in sync with index_no_bound_checks.go.

//+build bounds

package mat64

import "github.com/gonum/matrix"

// At returns the element at row r, column c.
func (m *Dense) At(r, c int) float64 {
	return m.at(r, c)
}

func (m *Dense) at(r, c int) float64 {
	if r >= m.mat.Rows || r < 0 {
		panic(matrix.ErrRowAccess)
	}
	if c >= m.mat.Cols || c < 0 {
		panic(matrix.ErrColAccess)
	}
	return m.mat.Data[r*m.mat.Stride+c]
}

// Set sets the element at row r, column c to the value v.
func (m *Dense) Set(r, c int, v float64) {
	m.set(r, c, v)
}

func (m *Dense) set(r, c int, v float64) {
	if r >= m.mat.Rows || r < 0 {
		panic(matrix.ErrRowAccess)
	}
	if c >= m.mat.Cols || c < 0 {
		panic(matrix.ErrColAccess)
	}
	m.mat.Data[r*m.mat.Stride+c] = v
}

// At returns the element at row i.
// It panics if i is out of bounds or if j is not zero.
func (v *Vector) At(i, j int) float64 {
	if j != 0 {
		panic(matrix.ErrColAccess)
	}
	return v.at(i)
}

func (v *Vector) at(i int) float64 {
	if i < 0 || i >= v.n {
		panic(matrix.ErrRowAccess)
	}
	return v.mat.Data[i*v.mat.Inc]
}

// SetVec sets the element at row i to the value val.
// It panics if i is out of bounds.
func (v *Vector) SetVec(i int, val float64) {
	v.setVec(i, val)
}

func (v *Vector) setVec(i int, val float64) {
	if i < 0 || i >= v.n {
		panic(matrix.ErrVectorAccess)
	}
	v.mat.Data[i*v.mat.Inc] = val
}

// At returns the element at row r and column c.
func (t *SymDense) At(r, c int) float64 {
	return t.at(r, c)
}

func (t *SymDense) at(r, c int) float64 {
	if r >= t.mat.N || r < 0 {
		panic(matrix.ErrRowAccess)
	}
	if c >= t.mat.N || c < 0 {
		panic(matrix.ErrColAccess)
	}
	if r > c {
		r, c = c, r
	}
	return t.mat.Data[r*t.mat.Stride+c]
}

// SetSym sets the elements at (r,c) and (c,r) to the value v.
func (t *SymDense) SetSym(r, c int, v float64) {
	t.set(r, c, v)
}

func (t *SymDense) set(r, c int, v float64) {
	if r >= t.mat.N || r < 0 {
		panic(matrix.ErrRowAccess)
	}
	if c >= t.mat.N || c < 0 {
		panic(matrix.ErrColAccess)
	}
	if r > c {
		r, c = c, r
	}
	t.mat.Data[r*t.mat.Stride+c] = v
}

// At returns the element at row r, column c.
func (t *TriDense) At(r, c int) float64 {
	return t.at(r, c)
}

func (t *TriDense) at(r, c int) float64 {
	if r >= t.mat.N || r < 0 {
		panic(matrix.ErrRowAccess)
	}
	if c >= t.mat.N || c < 0 {
		panic(matrix.ErrColAccess)
	}
	isUpper := t.isUpper()
	if (isUpper && r > c) || (!isUpper && r < c) {
		return 0
	}
	return t.mat.Data[r*t.mat.Stride+c]
}

// SetTri sets the element of the triangular matrix at row r, column c to the value v.
// It panics if the location is outside the appropriate half of the matrix.
func (t *TriDense) SetTri(r, c int, v float64) {
	t.set(r, c, v)
}

func (t *TriDense) set(r, c int, v float64) {
	if r >= t.mat.N || r < 0 {
		panic(matrix.ErrRowAccess)
	}
	if c >= t.mat.N || c < 0 {
		panic(matrix.ErrColAccess)
	}
	isUpper := t.isUpper()
	if (isUpper && r > c) || (!isUpper && r < c) {
		panic(matrix.ErrTriangleSet)
	}
	t.mat.Data[r*t.mat.Stride+c] = v
}
