// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"math"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

var inf = math.Inf(1)

const (
	epsilon = 2.2204e-16
	small   = math.SmallestNonzeroFloat64
)

// Norm returns the specified matrix p-norm of the receiver.
//
// See the Normer interface for more information.
func (m *Dense) Norm(ord float64) float64 {
	var n float64
	switch {
	case ord == 1:
		col := make([]float64, m.mat.Rows)
		for i := 0; i < m.mat.Cols; i++ {
			var s float64
			for _, e := range m.Col(col, i) {
				s += math.Abs(e)
			}
			n = math.Max(s, n)
		}
	case math.IsInf(ord, +1):
		row := make([]float64, m.mat.Cols)
		for i := 0; i < m.mat.Rows; i++ {
			var s float64
			for _, e := range m.Row(row, i) {
				s += math.Abs(e)
			}
			n = math.Max(s, n)
		}
	case ord == -1:
		n = math.MaxFloat64
		col := make([]float64, m.mat.Rows)
		for i := 0; i < m.mat.Cols; i++ {
			var s float64
			for _, e := range m.Col(col, i) {
				s += math.Abs(e)
			}
			n = math.Min(s, n)
		}
	case math.IsInf(ord, -1):
		n = math.MaxFloat64
		row := make([]float64, m.mat.Cols)
		for i := 0; i < m.mat.Rows; i++ {
			var s float64
			for _, e := range m.Row(row, i) {
				s += math.Abs(e)
			}
			n = math.Min(s, n)
		}
	case ord == 0:
		for i := 0; i < len(m.mat.Data); i += m.mat.Stride {
			for _, v := range m.mat.Data[i : i+m.mat.Cols] {
				n = math.Hypot(n, v)
			}
		}
		return n
	case ord == 2, ord == -2:
		s := SVD(m, epsilon, small, false, false).Sigma
		if ord == 2 {
			return s[0]
		}
		return s[len(s)-1]
	default:
		panic(ErrNormOrder)
	}

	return n
}

// Add adds a and b element-wise, placing the result in the receiver.
//
// See the Adder interface for more information.
func (m *Dense) Add(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || ac != bc {
		panic(ErrShape)
	}

	aMat, _ := untranspose(a)
	bMat, _ := untranspose(b)
	m.reuseAs(ar, ac)
	var restore func()
	if m == aMat {
		m, restore = m.isolatedWorkspace(aMat)
		defer restore()
	} else if m == bMat {
		m, restore = m.isolatedWorkspace(bMat)
		defer restore()
	}

	if a, ok := a.(RawMatrixer); ok {
		if b, ok := b.(RawMatrixer); ok {
			amat, bmat := a.RawMatrix(), b.RawMatrix()
			for ja, jb, jm := 0, 0, 0; ja < ar*amat.Stride; ja, jb, jm = ja+amat.Stride, jb+bmat.Stride, jm+m.mat.Stride {
				for i, v := range amat.Data[ja : ja+ac] {
					m.mat.Data[i+jm] = v + bmat.Data[i+jb]
				}
			}
			return
		}
	}

	if a, ok := a.(Vectorer); ok {
		if b, ok := b.(Vectorer); ok {
			rowa := make([]float64, ac)
			rowb := make([]float64, bc)
			for r := 0; r < ar; r++ {
				a.Row(rowa, r)
				for i, v := range b.Row(rowb, r) {
					rowa[i] += v
				}
				copy(m.rowView(r), rowa)
			}
			return
		}
	}

	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.set(r, c, a.At(r, c)+b.At(r, c))
		}
	}
}

// Sub subtracts the matrix b from a, placing the result in the receiver.
//
// See the Suber interface for more information.
func (m *Dense) Sub(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()

	if ar != br || ac != bc {
		panic(ErrShape)
	}

	m.reuseAs(ar, ac)

	if a, ok := a.(RawMatrixer); ok {
		if b, ok := b.(RawMatrixer); ok {
			amat, bmat := a.RawMatrix(), b.RawMatrix()
			for ja, jb, jm := 0, 0, 0; ja < ar*amat.Stride; ja, jb, jm = ja+amat.Stride, jb+bmat.Stride, jm+m.mat.Stride {
				for i, v := range amat.Data[ja : ja+ac] {
					m.mat.Data[i+jm] = v - bmat.Data[i+jb]
				}
			}
			return
		}
	}

	if a, ok := a.(Vectorer); ok {
		if b, ok := b.(Vectorer); ok {
			rowa := make([]float64, ac)
			rowb := make([]float64, bc)
			for r := 0; r < ar; r++ {
				a.Row(rowa, r)
				for i, v := range b.Row(rowb, r) {
					rowa[i] -= v
				}
				copy(m.rowView(r), rowa)
			}
			return
		}
	}

	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.set(r, c, a.At(r, c)-b.At(r, c))
		}
	}
}

// MulElem performs element-wise multiplication of a and b, placing the result
// in the receiver.
//
// See the ElemMuler interface for more information.
func (m *Dense) MulElem(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()

	if ar != br || ac != bc {
		panic(ErrShape)
	}

	m.reuseAs(ar, ac)

	if a, ok := a.(RawMatrixer); ok {
		if b, ok := b.(RawMatrixer); ok {
			amat, bmat := a.RawMatrix(), b.RawMatrix()
			for ja, jb, jm := 0, 0, 0; ja < ar*amat.Stride; ja, jb, jm = ja+amat.Stride, jb+bmat.Stride, jm+m.mat.Stride {
				for i, v := range amat.Data[ja : ja+ac] {
					m.mat.Data[i+jm] = v * bmat.Data[i+jb]
				}
			}
			return
		}
	}

	if a, ok := a.(Vectorer); ok {
		if b, ok := b.(Vectorer); ok {
			rowa := make([]float64, ac)
			rowb := make([]float64, bc)
			for r := 0; r < ar; r++ {
				a.Row(rowa, r)
				for i, v := range b.Row(rowb, r) {
					rowa[i] *= v
				}
				copy(m.rowView(r), rowa)
			}
			return
		}
	}

	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.set(r, c, a.At(r, c)*b.At(r, c))
		}
	}
}

// DivElem performs element-wise division of a by b, placing the result
// in the receiver.
//
// See the ElemDiver interface for more information.
func (m *Dense) DivElem(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()

	if ar != br || ac != bc {
		panic(ErrShape)
	}

	m.reuseAs(ar, ac)

	if a, ok := a.(RawMatrixer); ok {
		if b, ok := b.(RawMatrixer); ok {
			amat, bmat := a.RawMatrix(), b.RawMatrix()
			for ja, jb, jm := 0, 0, 0; ja < ar*amat.Stride; ja, jb, jm = ja+amat.Stride, jb+bmat.Stride, jm+m.mat.Stride {
				for i, v := range amat.Data[ja : ja+ac] {
					m.mat.Data[i+jm] = v / bmat.Data[i+jb]
				}
			}
			return
		}
	}

	if a, ok := a.(Vectorer); ok {
		if b, ok := b.(Vectorer); ok {
			rowa := make([]float64, ac)
			rowb := make([]float64, bc)
			for r := 0; r < ar; r++ {
				a.Row(rowa, r)
				for i, v := range b.Row(rowb, r) {
					rowa[i] /= v
				}
				copy(m.rowView(r), rowa)
			}
			return
		}
	}

	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.set(r, c, a.At(r, c)/b.At(r, c))
		}
	}
}

// Dot returns the sum of the element-wise products of the elements of the
// receiver and b.
//
// See the Dotter interface for more information.
func (m *Dense) Dot(b Matrix) float64 {
	mr, mc := m.Dims()
	br, bc := b.Dims()

	if mr != br || mc != bc {
		panic(ErrShape)
	}

	var d float64

	if b, ok := b.(RawMatrixer); ok {
		bmat := b.RawMatrix()
		for jm, jb := 0, 0; jm < mr*m.mat.Stride; jm, jb = jm+m.mat.Stride, jb+bmat.Stride {
			for i, v := range m.mat.Data[jm : jm+mc] {
				d += v * bmat.Data[i+jb]
			}
		}
		return d
	}

	if b, ok := b.(Vectorer); ok {
		row := make([]float64, bc)
		for r := 0; r < br; r++ {
			for i, v := range b.Row(row, r) {
				d += m.mat.Data[r*m.mat.Stride+i] * v
			}
		}
		return d
	}

	for r := 0; r < mr; r++ {
		for c := 0; c < mc; c++ {
			d += m.At(r, c) * b.At(r, c)
		}
	}
	return d
}

// Mul takes the matrix product of a and b, placing the result in the receiver.
//
// See the Muler interface for more information.
func (m *Dense) Mul(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()

	if ac != br {
		panic(ErrShape)
	}

	aU, aTrans := untranspose(a)
	bU, bTrans := untranspose(b)
	m.reuseAs(ar, bc)
	var restore func()
	if m == aU {
		m, restore = m.isolatedWorkspace(aU)
		defer restore()
	} else if m == bU {
		m, restore = m.isolatedWorkspace(bU)
		defer restore()
	}
	aT := blas.NoTrans
	if aTrans {
		aT = blas.Trans
	}
	bT := blas.NoTrans
	if bTrans {
		bT = blas.Trans
	}

	// Some of the cases do not have a transpose option, so create
	// temporary memory.
	// C = A^T * B = (B^T * A)^T
	// C^T = B^T * A.
	if aU, ok := aU.(RawMatrixer); ok {
		amat := aU.RawMatrix()
		if bU, ok := bU.(RawMatrixer); ok {
			bmat := bU.RawMatrix()
			blas64.Gemm(aT, bT, 1, amat, bmat, 0, m.mat)
			return
		}
		if bU, ok := bU.(RawSymmetricer); ok {
			bmat := bU.RawSymmetric()
			if aTrans {
				c := getWorkspace(ac, ar, false)
				blas64.Symm(blas.Left, 1, bmat, amat, 0, c.mat)
				strictCopy(m, c.T())
				putWorkspace(c)
				return
			}
			blas64.Symm(blas.Right, 1, bmat, amat, 0, m.mat)
			return
		}
		if bU, ok := bU.(RawTriangular); ok {
			// Trmm updates in place, so copy aU first.
			bmat := bU.RawTriangular()
			if aTrans {
				c := getWorkspace(ac, ar, false)
				var tmp Dense
				tmp.SetRawMatrix(aU.RawMatrix())
				c.Copy(&tmp)
				bT := blas.Trans
				if bTrans {
					bT = blas.NoTrans
				}
				blas64.Trmm(blas.Left, bT, 1, bmat, c.mat)
				strictCopy(m, c.T())
				putWorkspace(c)
				return
			}
			m.Copy(a)
			blas64.Trmm(blas.Right, bT, 1, bmat, m.mat)
			return
		}
		if bU, ok := bU.(*Vector); ok {
			bvec := bU.RawVector()
			if bTrans {
				// {ar,1} x {1,bc}, which is not a vector.
				// Instead, construct B as a General.
				bmat := blas64.General{
					Rows:   bc,
					Cols:   1,
					Stride: bvec.Inc,
					Data:   bvec.Data,
				}
				blas64.Gemm(aT, bT, 1, amat, bmat, 0, m.mat)
				return
			}
			cvec := blas64.Vector{
				Inc:  m.mat.Stride,
				Data: m.mat.Data,
			}
			blas64.Gemv(aT, 1, amat, bvec, 0, cvec)
			return
		}
	}
	if bU, ok := bU.(RawMatrixer); ok {
		bmat := bU.RawMatrix()
		if aU, ok := aU.(RawSymmetricer); ok {
			amat := aU.RawSymmetric()
			if bTrans {
				c := getWorkspace(bc, br, false)
				blas64.Symm(blas.Right, 1, amat, bmat, 0, c.mat)
				strictCopy(m, c.T())
				putWorkspace(c)
				return
			}
			blas64.Symm(blas.Left, 1, amat, bmat, 0, m.mat)
			return
		}
		if aU, ok := aU.(RawTriangular); ok {
			// Trmm updates in place, so copy bU first.
			amat := aU.RawTriangular()
			if bTrans {
				c := getWorkspace(bc, br, false)
				var tmp Dense
				tmp.SetRawMatrix(bU.RawMatrix())
				c.Copy(&tmp)
				aT := blas.Trans
				if aTrans {
					aT = blas.NoTrans
				}
				blas64.Trmm(blas.Right, aT, 1, amat, c.mat)
				strictCopy(m, c.T())
				putWorkspace(c)
				return
			}
			m.Copy(b)
			blas64.Trmm(blas.Left, aT, 1, amat, m.mat)
			return
		}
		if aU, ok := aU.(*Vector); ok {
			avec := aU.RawVector()
			if aTrans {
				// {1,ac} x {ac, bc}
				// Transpose B so that the vector is on the right.
				cvec := blas64.Vector{
					Inc:  1,
					Data: m.mat.Data,
				}
				bT := blas.Trans
				if bTrans {
					bT = blas.NoTrans
				}
				blas64.Gemv(bT, 1, bmat, avec, 0, cvec)
				return
			}
			// {ar,1} x {1,bc} which is not a vector result.
			// Instead, construct A as a General.
			amat := blas64.General{
				Rows:   ar,
				Cols:   1,
				Stride: avec.Inc,
				Data:   avec.Data,
			}
			blas64.Gemm(aT, bT, 1, amat, bmat, 0, m.mat)
			return
		}
	}

	if aU, ok := aU.(Vectorer); ok {
		if bU, ok := bU.(Vectorer); ok {
			row := make([]float64, ac)
			col := make([]float64, br)
			if aTrans {
				if bTrans {
					for r := 0; r < ar; r++ {
						dataTmp := m.mat.Data[r*m.mat.Stride : r*m.mat.Stride+bc]
						for c := 0; c < bc; c++ {
							dataTmp[c] = blas64.Dot(ac,
								blas64.Vector{Inc: 1, Data: aU.Col(row, r)},
								blas64.Vector{Inc: 1, Data: bU.Row(col, c)},
							)
						}
					}
					return
				}
				// TODO(jonlawlor): determine if (b*a)' is more efficient
				for r := 0; r < ar; r++ {
					dataTmp := m.mat.Data[r*m.mat.Stride : r*m.mat.Stride+bc]
					for c := 0; c < bc; c++ {
						dataTmp[c] = blas64.Dot(ac,
							blas64.Vector{Inc: 1, Data: aU.Col(row, r)},
							blas64.Vector{Inc: 1, Data: bU.Col(col, c)},
						)
					}
				}
				return
			}
			if bTrans {
				for r := 0; r < ar; r++ {
					dataTmp := m.mat.Data[r*m.mat.Stride : r*m.mat.Stride+bc]
					for c := 0; c < bc; c++ {
						dataTmp[c] = blas64.Dot(ac,
							blas64.Vector{Inc: 1, Data: aU.Row(row, r)},
							blas64.Vector{Inc: 1, Data: bU.Row(col, c)},
						)
					}
				}
				return
			}
			for r := 0; r < ar; r++ {
				dataTmp := m.mat.Data[r*m.mat.Stride : r*m.mat.Stride+bc]
				for c := 0; c < bc; c++ {
					dataTmp[c] = blas64.Dot(ac,
						blas64.Vector{Inc: 1, Data: aU.Row(row, r)},
						blas64.Vector{Inc: 1, Data: bU.Col(col, c)},
					)
				}
			}
			return
		}
	}

	row := make([]float64, ac)
	for r := 0; r < ar; r++ {
		for i := range row {
			row[i] = a.At(r, i)
		}
		for c := 0; c < bc; c++ {
			var v float64
			for i, e := range row {
				v += e * b.At(i, c)
			}
			m.mat.Data[r*m.mat.Stride+c] = v
		}
	}
}

// strictCopy copies a into m panicking if the shape of a and m differ.
func strictCopy(m *Dense, a Matrix) {
	r, c := m.Copy(a)
	if r != m.mat.Rows || c != m.mat.Cols {
		// Panic with a string since this
		// is not a user-facing panic.
		panic(ErrShape.string)
	}
}

// Exp calculates the exponential of the matrix a, e^a, placing the result
// in the receiver.
//
// See the Exper interface for more information.
//
// Exp uses the scaling and squaring method described in section 3 of
// http://www.cs.cornell.edu/cv/researchpdf/19ways+.pdf.
func (m *Dense) Exp(a Matrix) {
	r, c := a.Dims()
	if r != c {
		panic(ErrShape)
	}

	var w *Dense
	switch {
	case m.isZero():
		m.mat = blas64.General{
			Rows:   r,
			Cols:   c,
			Stride: c,
			Data:   useZeroed(m.mat.Data, r*r),
		}
		m.capRows = r
		m.capCols = c
		for i := 0; i < r*r; i += r + 1 {
			m.mat.Data[i] = 1
		}
		w = m
	case r == m.mat.Rows && c == m.mat.Cols:
		w = getWorkspace(r, r, true)
		for i := 0; i < r; i++ {
			w.mat.Data[i*w.mat.Stride+i] = 1
		}
	default:
		panic(ErrShape)
	}

	const (
		terms   = 10
		scaling = 4
	)

	small := getWorkspace(r, r, false)
	small.Scale(math.Pow(2, -scaling), a)
	power := getWorkspace(r, r, false)
	power.Copy(small)

	var (
		tmp   = getWorkspace(r, r, false)
		factI = 1.
	)
	for i := 1.; i < terms; i++ {
		factI *= i

		// This is OK to do because power and tmp are
		// new Dense values so all rows are contiguous.
		// TODO(kortschak) Make this explicit in the NewDense doc comment.
		for j, v := range power.mat.Data {
			tmp.mat.Data[j] = v / factI
		}

		w.Add(w, tmp)
		if i < terms-1 {
			tmp.Mul(power, small)
			tmp, power = power, tmp
		}
	}
	putWorkspace(small)
	putWorkspace(power)
	for i := 0; i < scaling; i++ {
		tmp.Mul(w, w)
		tmp, w = w, tmp
	}
	putWorkspace(tmp)

	if w != m {
		m.Copy(w)
		putWorkspace(w)
	}
}

// Pow calculates the integral power of the matrix a to n, placing the result
// in the receiver.
//
// See the Power interface for more information.
func (m *Dense) Pow(a Matrix, n int) {
	if n < 0 {
		panic("matrix: illegal power")
	}
	r, c := a.Dims()
	if r != c {
		panic(ErrShape)
	}

	m.reuseAs(r, c)

	// Take possible fast paths.
	switch n {
	case 0:
		for i := 0; i < r; i++ {
			zero(m.mat.Data[i*m.mat.Stride : i*m.mat.Stride+c])
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
	w := getWorkspace(r, r, false)
	w.Copy(a)
	s := getWorkspace(r, r, false)
	s.Copy(a)
	x := getWorkspace(r, r, false)
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
	putWorkspace(w)
	putWorkspace(s)
	putWorkspace(x)
}

// Scale multiplies the elements of a by f, placing the result in the receiver.
//
// See the Scaler interface for more information.
func (m *Dense) Scale(f float64, a Matrix) {
	ar, ac := a.Dims()

	m.reuseAs(ar, ac)

	if a, ok := a.(RawMatrixer); ok {
		amat := a.RawMatrix()
		for ja, jm := 0, 0; ja < ar*amat.Stride; ja, jm = ja+amat.Stride, jm+m.mat.Stride {
			for i, v := range amat.Data[ja : ja+ac] {
				m.mat.Data[i+jm] = v * f
			}
		}
		return
	}

	if a, ok := a.(Vectorer); ok {
		row := make([]float64, ac)
		for r := 0; r < ar; r++ {
			for i, v := range a.Row(row, r) {
				row[i] = f * v
			}
			copy(m.rowView(r), row)
		}
		return
	}

	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.set(r, c, f*a.At(r, c))
		}
	}
}

// Apply applies the function fn to each of the elements of a, placing the
// resulting matrix in the receiver.
//
// See the Applyer interface for more information.
func (m *Dense) Apply(fn func(r, c int, v float64) float64, a Matrix) {
	ar, ac := a.Dims()

	m.reuseAs(ar, ac)

	if a, ok := a.(RawMatrixer); ok {
		amat := a.RawMatrix()
		for j, ja, jm := 0, 0, 0; ja < ar*amat.Stride; j, ja, jm = j+1, ja+amat.Stride, jm+m.mat.Stride {
			for i, v := range amat.Data[ja : ja+ac] {
				m.mat.Data[i+jm] = fn(j, i, v)
			}
		}
		return
	}

	if a, ok := a.(Vectorer); ok {
		row := make([]float64, ac)
		for r := 0; r < ar; r++ {
			for i, v := range a.Row(row, r) {
				row[i] = fn(r, i, v)
			}
			copy(m.rowView(r), row)
		}
		return
	}

	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.set(r, c, fn(r, c, a.At(r, c)))
		}
	}
}

// Sum returns the sum of the elements of the matrix.
//
// See the Sumer interface for more information.
func (m *Dense) Sum() float64 {
	l := m.mat.Cols
	var s float64
	for i := 0; i < len(m.mat.Data); i += m.mat.Stride {
		for _, v := range m.mat.Data[i : i+l] {
			s += v
		}
	}
	return s
}

// EqualsApprox compares the matrices represented by b and the receiver, with
// tolerance for element-wise equality specified by epsilon.
//
// See the ApproxEqualer interface for more information.
func (m *Dense) EqualsApprox(b Matrix, epsilon float64) bool {
	br, bc := b.Dims()
	if br != m.mat.Rows || bc != m.mat.Cols {
		return false
	}

	if b, ok := b.(RawMatrixer); ok {
		bmat := b.RawMatrix()
		for jb, jm := 0, 0; jm < br*m.mat.Stride; jb, jm = jb+bmat.Stride, jm+m.mat.Stride {
			for i, v := range m.mat.Data[jm : jm+bc] {
				if math.Abs(v-bmat.Data[i+jb]) > epsilon {
					return false
				}
			}
		}
		return true
	}

	if b, ok := b.(Vectorer); ok {
		rowb := make([]float64, bc)
		for r := 0; r < br; r++ {
			rowm := m.mat.Data[r*m.mat.Stride : r*m.mat.Stride+m.mat.Cols]
			for i, v := range b.Row(rowb, r) {
				if math.Abs(rowm[i]-v) > epsilon {
					return false
				}
			}
		}
		return true
	}

	for r := 0; r < br; r++ {
		for c := 0; c < bc; c++ {
			if math.Abs(m.At(r, c)-b.At(r, c)) > epsilon {
				return false
			}
		}
	}
	return true
}

// RankOne performs a rank-one update to the matrix a and stores the result
// in the receiver. If a is zero, see Outer.
//  m = a + alpha * x * y'
func (m *Dense) RankOne(a Matrix, alpha float64, x, y *Vector) {
	ar, ac := a.Dims()
	if x.Len() != ar {
		panic(ErrShape)
	}
	if y.Len() != ac {
		panic(ErrShape)
	}

	var w Dense
	if m == a {
		w = *m
	}
	w.reuseAs(ar, ac)

	// Copy over to the new memory if necessary
	if m != a {
		w.Copy(a)
	}
	blas64.Ger(alpha, x.mat, y.mat, w.mat)
	*m = w
}

// Outer calculates the outer product of x and y, and stores the result
// in the receiver. In order to update to an existing matrix, see RankOne.
//  m = x * y'
func (m *Dense) Outer(x, y *Vector) {
	r := x.Len()
	c := y.Len()

	// Copied from reuseAs with use replaced by useZeroed
	// and a final zero of the matrix elements if we pass
	// the shape checks.
	// TODO(kortschak): Factor out into reuseZeroedAs if
	// we find another case that needs it.
	if m.mat.Rows > m.capRows || m.mat.Cols > m.capCols {
		// Panic as a string, not a mat64.Error.
		panic("mat64: caps not correctly set")
	}
	if m.isZero() {
		m.mat = blas64.General{
			Rows:   r,
			Cols:   c,
			Stride: c,
			Data:   useZeroed(m.mat.Data, r*c),
		}
		m.capRows = r
		m.capCols = c
	} else if r != m.mat.Rows || c != m.mat.Cols {
		panic(ErrShape)
	} else {
		for i := 0; i < r; i++ {
			zero(m.mat.Data[i*m.mat.Stride : i*m.mat.Stride+c])
		}
	}

	blas64.Ger(1, x.mat, y.mat, m.mat)
}
