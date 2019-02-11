// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"math/cmplx"

	"gonum.org/v1/gonum/floats"
)

// CMatrix is the basic matrix interface type for complex matrices.
type CMatrix interface {
	// Dims returns the dimensions of a Matrix.
	Dims() (r, c int)

	// At returns the value of a matrix element at row i, column j.
	// It will panic if i or j are out of bounds for the matrix.
	At(i, j int) complex128

	// H returns the conjugate transpose of the Matrix. Whether H
	// returns a copy of the underlying data is implementation dependent.
	// This method may be implemented using the Conjugate type, which
	// provides an implicit matrix conjugate transpose.
	H() CMatrix
}

var (
	_ CMatrix      = Conjugate{}
	_ Unconjugator = Conjugate{}
)

// Conjugate is a type for performing an implicit matrix conjugate transpose.
// It implements the Matrix interface, returning values from the conjugate
// transpose of the matrix within.
type Conjugate struct {
	CMatrix CMatrix
}

// At returns the value of the element at row i and column j of the conjugate
// transposed matrix, that is, row j and column i of the Matrix field.
func (t Conjugate) At(i, j int) complex128 {
	z := t.CMatrix.At(j, i)
	return cmplx.Conj(z)
}

// Dims returns the dimensions of the transposed matrix. The number of rows returned
// is the number of columns in the Matrix field, and the number of columns is
// the number of rows in the Matrix field.
func (t Conjugate) Dims() (r, c int) {
	c, r = t.CMatrix.Dims()
	return r, c
}

// H performs an implicit conjugate transpose by returning the Matrix field.
func (t Conjugate) H() CMatrix {
	return t.CMatrix
}

// Unconjugate returns the Matrix field.
func (t Conjugate) Unconjugate() CMatrix {
	return t.CMatrix
}

// Unconjugator is a type that can undo an implicit conjugate transpose.
type Unconjugator interface {
	// Note: This interface is needed to unify all of the Conjugate types. In
	// the cmat128 methods, we need to test if the Matrix has been implicitly
	// transposed. If this is checked by testing for the specific Conjugate type
	// then the behavior will be different if the user uses H() or HTri() for a
	// triangular matrix.

	// Unconjugate returns the underlying Matrix stored for the implicit
	// conjugate transpose.
	Unconjugate() CMatrix
}

// useC returns a complex128 slice with l elements, using f if it
// has the necessary capacity, otherwise creating a new slice.
func useC(f []complex128, l int) []complex128 {
	if l <= cap(f) {
		return f[:l]
	}
	return make([]complex128, l)
}

// zeroC zeros the given slice's elements.
func zeroC(f []complex128) {
	for i := range f {
		f[i] = 0
	}
}

// unconjugate unconjugates a matrix if applicable. If a is an Unconjugator, then
// unconjugate returns the underlying matrix and true. If it is not, then it returns
// the input matrix and false.
func unconjugate(a CMatrix) (CMatrix, bool) {
	if ut, ok := a.(Unconjugator); ok {
		return ut.Unconjugate(), true
	}
	return a, false
}

// CEqual returns whether the matrices a and b have the same size
// and are element-wise equal.
func CEqual(a, b CMatrix) bool {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || ac != bc {
		return false
	}
	// TODO(btracey): Add in fast-paths.
	for i := 0; i < ar; i++ {
		for j := 0; j < ac; j++ {
			if a.At(i, j) != b.At(i, j) {
				return false
			}
		}
	}
	return true
}

// CEqualApprox returns whether the matrices a and b have the same size and contain all equal
// elements with tolerance for element-wise equality specified by epsilon. Matrices
// with non-equal shapes are not equal. Element-wise equality is compared on
// the real and imaginary parts separately.
func CEqualApprox(a, b CMatrix, epsilon float64) bool {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || ac != bc {
		return false
	}
	for i := 0; i < ar; i++ {
		for j := 0; j < ac; j++ {
			ca := a.At(i, j)
			cb := b.At(i, j)
			if !floats.EqualWithinAbsOrRel(real(ca), real(cb), epsilon, epsilon) {
				return false
			}
			if !floats.EqualWithinAbsOrRel(imag(ca), imag(cb), epsilon, epsilon) {
				return false
			}
		}
	}
	return true
}
