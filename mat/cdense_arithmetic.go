// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

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
