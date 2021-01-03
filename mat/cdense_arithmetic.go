// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/cblas128"
)

// Add adds a and b element-wise, placing the result in the receiver. Add
// will panic if the two matrices do not have the same shape.
func (m *CDense) Add(a, b CMatrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || ac != bc {
		panic(ErrShape)
	}

	aU, aTrans, aConj := untransposeCmplx(a)
	bU, bTrans, bConj := untransposeCmplx(b)

	m.reuseAsNonZeroed(ar, ac)

	if arm, ok := a.(*CDense); ok {
		if brm, ok := b.(*CDense); ok {
			amat, bmat := arm.mat, brm.mat
			if m != aU {
				m.checkOverlap(amat)
			}
			if m != bU {
				m.checkOverlap(bmat)
			}
			for ja, jb, jm := 0, 0, 0; ja < ar*amat.Stride; ja, jb, jm = ja+amat.Stride, jb+bmat.Stride, jm+m.mat.Stride {
				for i, v := range amat.Data[ja : ja+ac] {
					m.mat.Data[i+jm] = v + bmat.Data[i+jb]
				}
			}
			return
		}
	}

	m.checkOverlapMatrix(aU)
	m.checkOverlapMatrix(bU)
	var restore func()
	if aTrans != aConj && m == aU {
		m, restore = m.isolatedWorkspace(aU)
		defer restore()
	} else if bTrans != bConj && m == bU {
		m, restore = m.isolatedWorkspace(bU)
		defer restore()
	}

	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.set(r, c, a.At(r, c)+b.At(r, c))
		}
	}
}

// Sub subtracts the matrix b from a, placing the result in the receiver. Sub
// will panic if the two matrices do not have the same shape.
func (m *CDense) Sub(a, b CMatrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || ac != bc {
		panic(ErrShape)
	}

	aU, aTrans, aConj := untransposeCmplx(a)
	bU, bTrans, bConj := untransposeCmplx(b)
	m.reuseAsNonZeroed(ar, ac)

	if arm, ok := a.(*CDense); ok {
		if brm, ok := b.(*CDense); ok {
			amat, bmat := arm.mat, brm.mat
			if m != aU {
				m.checkOverlap(amat)
			}
			if m != bU {
				m.checkOverlap(bmat)
			}
			for ja, jb, jm := 0, 0, 0; ja < ar*amat.Stride; ja, jb, jm = ja+amat.Stride, jb+bmat.Stride, jm+m.mat.Stride {
				for i, v := range amat.Data[ja : ja+ac] {
					m.mat.Data[i+jm] = v - bmat.Data[i+jb]
				}
			}
			return
		}
	}

	m.checkOverlapMatrix(aU)
	m.checkOverlapMatrix(bU)
	var restore func()
	if aTrans != aConj && m == aU {
		m, restore = m.isolatedWorkspace(aU)
		defer restore()
	} else if bTrans != bConj && m == bU {
		m, restore = m.isolatedWorkspace(bU)
		defer restore()
	}

	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.set(r, c, a.At(r, c)-b.At(r, c))
		}
	}

}

// MulElem performs element-wise multiplication of a and b, placing the result
// in the receiver. MulElem will panic if the two matrices do not have the same
// shape.
func (m *CDense) MulElem(a, b CMatrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || ac != bc {
		panic(ErrShape)
	}

	aU, aTrans, aConj := untransposeCmplx(a)
	bU, bTrans, bConj := untransposeCmplx(b)
	m.reuseAsNonZeroed(ar, ac)

	if arm, ok := a.(*CDense); ok {
		if brm, ok := b.(*CDense); ok {
			amat, bmat := arm.mat, brm.mat
			if m != aU {
				m.checkOverlap(amat)
			}
			if m != bU {
				m.checkOverlap(bmat)
			}
			for ja, jb, jm := 0, 0, 0; ja < ar*amat.Stride; ja, jb, jm = ja+amat.Stride, jb+bmat.Stride, jm+m.mat.Stride {
				for i, v := range amat.Data[ja : ja+ac] {
					m.mat.Data[i+jm] = v * bmat.Data[i+jb]
				}
			}
			return
		}
	}

	m.checkOverlapMatrix(aU)
	m.checkOverlapMatrix(bU)
	var restore func()
	if aTrans != aConj && m == aU {
		m, restore = m.isolatedWorkspace(aU)
		defer restore()
	} else if bTrans != bConj && m == bU {
		m, restore = m.isolatedWorkspace(bU)
		defer restore()
	}

	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.set(r, c, a.At(r, c)*b.At(r, c))
		}
	}
}

// DivElem performs element-wise division of a by b, placing the result
// in the receiver. DivElem will panic if the two matrices do not have the same
// shape.
func (m *CDense) DivElem(a, b CMatrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || ac != bc {
		panic(ErrShape)
	}

	aU, aTrans, aConj := untransposeCmplx(a)
	bU, bTrans, bConj := untransposeCmplx(b)
	m.reuseAsNonZeroed(ar, ac)

	if arm, ok := a.(*CDense); ok {
		if brm, ok := b.(*CDense); ok {
			amat, bmat := arm.mat, brm.mat
			if m != aU {
				m.checkOverlap(amat)
			}
			if m != bU {
				m.checkOverlap(bmat)
			}
			for ja, jb, jm := 0, 0, 0; ja < ar*amat.Stride; ja, jb, jm = ja+amat.Stride, jb+bmat.Stride, jm+m.mat.Stride {
				for i, v := range amat.Data[ja : ja+ac] {
					m.mat.Data[i+jm] = v / bmat.Data[i+jb]
				}
			}
			return
		}
	}

	m.checkOverlapMatrix(aU)
	m.checkOverlapMatrix(bU)
	var restore func()
	if aTrans != aConj && m == aU {
		m, restore = m.isolatedWorkspace(aU)
		defer restore()
	} else if bTrans != bConj && m == bU {
		m, restore = m.isolatedWorkspace(bU)
		defer restore()
	}

	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.set(r, c, a.At(r, c)/b.At(r, c))
		}
	}
}

// Mul takes the matrix product of a and b, placing the result in the receiver.
// If the number of columns in a does not equal the number of rows in b, Mul will panic.
func (m *CDense) Mul(a, b CMatrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()

	if ac != br {
		panic(ErrShape)
	}

	aU, aTrans, aConj := untransposeCmplx(a)
	bU, bTrans, bConj := untransposeCmplx(b)
	m.reuseAsNonZeroed(ar, bc)
	var restore func()
	if aTrans != aConj && m == aU {
		m, restore = m.isolatedWorkspace(aU)
		defer restore()
	} else if bTrans != bConj && m == bU {
		m, restore = m.isolatedWorkspace(bU)
		defer restore()
	}
	aT := blas.NoTrans
	if aTrans {
		aT = blas.Trans
	} else if aConj {
		aT = blas.ConjTrans
	}
	bT := blas.NoTrans
	if bTrans != bConj {
		bT = blas.Trans
	} else if bConj {
		bT = blas.ConjTrans
	}

	if aU, ok := aU.(*CDense); ok {
		switch bU := bU.(type) {
		case *CDense:
			if restore == nil {
				m.checkOverlap(bU.mat)
			}
			cblas128.Gemm(aT, bT, 1, aU.mat, bU.mat, 0, m.mat)
			return
		}
	}

	m.checkOverlapMatrix(aU)
	m.checkOverlapMatrix(bU)
	row := getZs(ac, false)
	defer putZs(row)
	for r := 0; r < ar; r++ {
		for i := range row {
			row[i] = a.At(r, i)
		}
		for c := 0; c < bc; c++ {
			var v complex128
			for i, e := range row {
				v += e * b.At(i, c)
			}
			m.mat.Data[r*m.mat.Stride+c] = v
		}
	}
}

// Pow calculates the integral power of the matrix a to n, placing the result
// in the receiver. Pow will panic if n is negative or if a is not square.
func (m *CDense) Pow(a CMatrix, n int) {
	if n < 0 {
		panic("mat: illegal power")
	}
	r, c := a.Dims()
	if r != c {
		panic(ErrShape)
	}

	m.reuseAsNonZeroed(r, c)

	// Take possible fast paths.
	switch n {
	case 0:
		for i := 0; i < r; i++ {
			zeroC(m.mat.Data[i*m.mat.Stride : i*m.mat.Stride+c])
			m.mat.Data[i*m.mat.Stride+i] = 1
		}
		return
	case 1:
		m.Copy(a)
		return
	case 2:
		m.Mul(a, a)
		return
	}

	// Perform iterative exponentiation by squaring in work space.
	w := getWorkspaceCmplx(r, r, false)
	w.Copy(a)
	s := getWorkspaceCmplx(r, r, false)
	s.Copy(a)
	x := getWorkspaceCmplx(r, r, false)
	for n--; n > 0; n >>= 1 {
		if n&1 != 0 {
			x.Mul(w, s)
			w, x = x, w
		}
		if n != 1 {
			x.Mul(s, s)
			s, x = x, s
		}
	}
	m.Copy(w)
	putWorkspaceCmplx(w)
	putWorkspaceCmplx(s)
	putWorkspaceCmplx(x)
}

// Kronecker calculates the Kronecker product of a and b, placing the result in
// the receiver.
func (m *CDense) Kronecker(a, b CMatrix) {
	ra, ca := a.Dims()
	rb, cb := b.Dims()

	m.reuseAsNonZeroed(ra*rb, ca*cb)
	for i := 0; i < ra; i++ {
		for j := 0; j < ca; j++ {
			m.slice(i*rb, (i+1)*rb, j*cb, (j+1)*cb).Scale(a.At(i, j), b)
		}
	}
}

// Scale multiplies the elements of a by f, placing the result in the receiver.
//
// See the Scaler interface for more information.
func (m *CDense) Scale(f complex128, a CMatrix) {
	ar, ac := a.Dims()

	m.reuseAsNonZeroed(ar, ac)

	aU, aTrans, aConj := untransposeExtractCmplx(a)
	if rm, ok := aU.(*CDense); ok {
		amat := rm.mat
		if (aTrans != aConj && m == aU) || m.checkOverlap(amat) {
			var restore func()
			m, restore = m.isolatedWorkspace(a)
			defer restore()
		}
		if aTrans == aConj {
			for ja, jm := 0, 0; ja < ar*amat.Stride; ja, jm = ja+amat.Stride, jm+m.mat.Stride {
				for i, v := range amat.Data[ja : ja+ac] {
					m.mat.Data[i+jm] = v * f
				}
			}
		} else {
			for ja, jm := 0, 0; ja < ac*amat.Stride; ja, jm = ja+amat.Stride, jm+1 {
				for i, v := range amat.Data[ja : ja+ar] {
					m.mat.Data[i*m.mat.Stride+jm] = v * f
				}
			}
		}
		return
	}

	m.checkOverlapMatrix(a)
	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.set(r, c, f*a.At(r, c))
		}
	}
}

// Apply applies the function fn to each of the elements of a, placing the
// resulting matrix in the receiver. The function fn takes a row/column
// index and element value and returns some function of that tuple.
func (m *CDense) Apply(fn func(i, j int, v complex128) complex128, a CMatrix) {
	ar, ac := a.Dims()

	m.reuseAsNonZeroed(ar, ac)

	aU, aTrans, aConj := untransposeExtractCmplx(a)
	if rm, ok := aU.(*CDense); ok {
		amat := rm.mat
		if (aTrans != aConj) && m == aU || m.checkOverlap(amat) {
			var restore func()
			m, restore = m.isolatedWorkspace(a)
			defer restore()
		}
		if aTrans == aConj {
			for j, ja, jm := 0, 0, 0; ja < ar*amat.Stride; j, ja, jm = j+1, ja+amat.Stride, jm+m.mat.Stride {
				for i, v := range amat.Data[ja : ja+ac] {
					m.mat.Data[i+jm] = fn(j, i, v)
				}
			}
		} else {
			for j, ja, jm := 0, 0, 0; ja < ac*amat.Stride; j, ja, jm = j+1, ja+amat.Stride, jm+1 {
				for i, v := range amat.Data[ja : ja+ar] {
					m.mat.Data[i*m.mat.Stride+jm] = fn(i, j, v)
				}
			}
		}
		return
	}

	m.checkOverlapMatrix(a)
	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.set(r, c, fn(r, c, a.At(r, c)))
		}
	}
}
