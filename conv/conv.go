// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package conv provides matrix type interconversion utilities.
package conv

import (
	"github.com/gonum/matrix/cmat128"
	"github.com/gonum/matrix/mat64"
)

// Complex is a complex matrix constructed from two real matrices.
type Complex struct {
	// r and i are not exposed to ensure that
	// their dimensions can not be altered by
	// clients behind our back.
	r, i mat64.Matrix
}

var (
	_ Realer = Complex{}
	_ Imager = Complex{}
)

// NewComplex returns a complex matrix constructed from r and i. At least one of
// r or i must be non-nil otherwise NewComples will panic. If one of the inputs
// is nil, that part of the complex number will be zero when returned by At.
// If both are non-nil but differ in their sizes, NewComplex will panic.
func NewComplex(r, i mat64.Matrix) Complex {
	if r == nil && i == nil {
		panic("conv: no matrix")
	} else if r != nil && i != nil {
		rr, rc := r.Dims()
		ir, ic := i.Dims()
		if rr != ir || rc != ic {
			panic(mat64.ErrShape)
		}
	}
	return Complex{r: r, i: i}
}

// Dims returns the number of rows and columns in the matrix.
func (m Complex) Dims() (r, c int) {
	if m.r == nil {
		return m.i.Dims()
	}
	return m.r.Dims()
}

// At returns the element at row r, column c.
func (m Complex) At(r, c int) complex128 {
	if m.i == nil {
		return complex(m.r.At(r, c), 0)
	}
	if m.r == nil {
		return complex(0, m.i.At(r, c))
	}
	return complex(m.r.At(r, c), m.i.At(r, c))
}

// T performs an implicit transpose.
func (m Complex) T() cmat128.Matrix {
	if m.i == nil {
		return Complex{r: m.r.T()}
	}
	if m.r == nil {
		return Complex{i: m.i.T()}
	}
	return Complex{r: m.r.T(), i: m.i.T()}
}

// Real returns the real part of the receiver.
func (m Complex) Real() mat64.Matrix { return m.r }

// Imag returns the imaginary part of the receiver.
func (m Complex) Imag() mat64.Matrix { return m.i }

// Realer is a complex matrix that can return its real part.
type Realer interface {
	Real() mat64.Matrix
}

// Imager is a complex matrix that can return its imaginary part.
type Imager interface {
	Imag() mat64.Matrix
}

// Real is the real part of a complex matrix.
type Real struct{ Matrix cmat128.Matrix }

// NewReal returns a mat64.Matrix representing the real part of m. If m is a Realer,
// the real part is returned.
func NewReal(m cmat128.Matrix) mat64.Matrix {
	if m, ok := m.(Realer); ok {
		return m.Real()
	}
	return Real{m}
}

// Dims returns the number of rows and columns in the matrix.
func (m Real) Dims() (r, c int) { return m.Matrix.Dims() }

// At returns the element at row r, column c.
func (m Real) At(r, c int) float64 { return real(m.Matrix.At(r, c)) }

// T performs an implicit transpose.
func (m Real) T() mat64.Matrix { return Real{m.Matrix.T()} }

// Imag is the imaginary part of a complex matrix.
type Imag struct{ Matrix cmat128.Matrix }

// NewImage returns a mat64.Matrix representing the imaginary part of m. If m is an Imager,
// the imaginary part is returned.
func NewImag(m cmat128.Matrix) mat64.Matrix {
	if m, ok := m.(Imager); ok {
		return m.Imag()
	}
	return Imag{m}
}

// Dims returns the number of rows and columns in the matrix.
func (m Imag) Dims() (r, c int) { return m.Matrix.Dims() }

// At returns the element at row r, column c.
func (m Imag) At(r, c int) float64 { return imag(m.Matrix.At(r, c)) }

// T performs an implicit transpose.
func (m Imag) T() mat64.Matrix { return Imag{m.Matrix.T()} }
