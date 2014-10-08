// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"math"

	"github.com/gonum/blas"
)

func (m *Dense) Min() float64 {
	min := m.mat.Data[0]
	for k := 0; k < m.mat.Rows; k++ {
		for _, v := range m.rowView(k) {
			min = math.Min(min, v)
		}
	}
	return min
}

func (m *Dense) Max() float64 {
	max := m.mat.Data[0]
	for k := 0; k < m.mat.Rows; k++ {
		for _, v := range m.rowView(k) {
			max = math.Max(max, v)
		}
	}
	return max
}

func (m *Dense) Trace() float64 {
	if m.mat.Rows != m.mat.Cols {
		panic(ErrSquare)
	}
	var t float64
	for i := 0; i < len(m.mat.Data); i += m.mat.Stride + 1 {
		t += m.mat.Data[i]
	}
	return t
}

var inf = math.Inf(1)

const (
	epsilon = 2.2204e-16
	small   = math.SmallestNonzeroFloat64
)

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

func (m *Dense) Add(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()

	if ar != br || ac != bc {
		panic(ErrShape)
	}

	if m.isZero() {
		m.mat = RawMatrix{
			Rows:   ar,
			Cols:   ac,
			Stride: ac,
			Data:   use(m.mat.Data, ar*ac),
		}
	} else if ar != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
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
			m.Set(r, c, a.At(r, c)+b.At(r, c))
		}
	}
}

func (m *Dense) Sub(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()

	if ar != br || ac != bc {
		panic(ErrShape)
	}

	if m.isZero() {
		m.mat = RawMatrix{
			Rows:   ar,
			Cols:   ac,
			Stride: ac,
			Data:   use(m.mat.Data, ar*ac),
		}
	} else if ar != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
	}

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
			m.Set(r, c, a.At(r, c)-b.At(r, c))
		}
	}
}

func (m *Dense) MulElem(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()

	if ar != br || ac != bc {
		panic(ErrShape)
	}

	if m.isZero() {
		m.mat = RawMatrix{
			Rows:   ar,
			Cols:   ac,
			Stride: ac,
			Data:   use(m.mat.Data, ar*ac),
		}
	} else if ar != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
	}

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
			m.Set(r, c, a.At(r, c)*b.At(r, c))
		}
	}
}

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

func (m *Dense) Mul(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()

	if ac != br {
		panic(ErrShape)
	}

	var w Dense
	if m != a && m != b {
		w = *m
	}
	if w.isZero() {
		w.mat = RawMatrix{
			Rows:   ar,
			Cols:   bc,
			Stride: bc,
			Data:   use(w.mat.Data, ar*bc),
		}
	} else if ar != w.mat.Rows || bc != w.mat.Cols {
		panic(ErrShape)
	}

	if a, ok := a.(RawMatrixer); ok {
		if b, ok := b.(RawMatrixer); ok {
			amat, bmat := a.RawMatrix(), b.RawMatrix()
			if blasEngine == nil {
				panic(ErrNoEngine)
			}
			blasEngine.Dgemm(
				blas.NoTrans, blas.NoTrans,
				ar, bc, ac,
				1.,
				amat.Data, amat.Stride,
				bmat.Data, bmat.Stride,
				0.,
				w.mat.Data, w.mat.Stride)
			*m = w
			return
		}
	}

	if a, ok := a.(Vectorer); ok {
		if b, ok := b.(Vectorer); ok {
			row := make([]float64, ac)
			col := make([]float64, br)
			if blasEngine == nil {
				panic(ErrNoEngine)
			}
			for r := 0; r < ar; r++ {
				for c := 0; c < bc; c++ {
					w.mat.Data[r*w.mat.Stride+c] = blasEngine.Ddot(ac, a.Row(row, r), 1, b.Col(col, c), 1)
				}
			}
			*m = w
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
			w.mat.Data[r*w.mat.Stride+c] = v
		}
	}
	*m = w
}

func (m *Dense) Scale(f float64, a Matrix) {
	ar, ac := a.Dims()

	if m.isZero() {
		m.mat = RawMatrix{
			Rows:   ar,
			Cols:   ac,
			Stride: ac,
			Data:   use(m.mat.Data, ar*ac),
		}
	} else if ar != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
	}

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
			m.Set(r, c, f*a.At(r, c))
		}
	}
}

func (m *Dense) Apply(f ApplyFunc, a Matrix) {
	ar, ac := a.Dims()

	if m.isZero() {
		m.mat = RawMatrix{
			Rows:   ar,
			Cols:   ac,
			Stride: ac,
			Data:   use(m.mat.Data, ar*ac),
		}
	} else if ar != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
	}

	if a, ok := a.(RawMatrixer); ok {
		amat := a.RawMatrix()
		for j, ja, jm := 0, 0, 0; ja < ar*amat.Stride; j, ja, jm = j+1, ja+amat.Stride, jm+m.mat.Stride {
			for i, v := range amat.Data[ja : ja+ac] {
				m.mat.Data[i+jm] = f(j, i, v)
			}
		}
		return
	}

	if a, ok := a.(Vectorer); ok {
		row := make([]float64, ac)
		for r := 0; r < ar; r++ {
			for i, v := range a.Row(row, r) {
				row[i] = f(r, i, v)
			}
			copy(m.rowView(r), row)
		}
		return
	}

	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.Set(r, c, f(r, c, a.At(r, c)))
		}
	}
}

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

func (m *Dense) Equals(b Matrix) bool {
	br, bc := b.Dims()
	if br != m.mat.Rows || bc != m.mat.Cols {
		return false
	}

	if b, ok := b.(RawMatrixer); ok {
		bmat := b.RawMatrix()
		for jb, jm := 0, 0; jm < br*m.mat.Stride; jb, jm = jb+bmat.Stride, jm+m.mat.Stride {
			for i, v := range m.mat.Data[jm : jm+bc] {
				if v != bmat.Data[i+jb] {
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
				if rowm[i] != v {
					return false
				}
			}
		}
		return true
	}

	for r := 0; r < br; r++ {
		for c := 0; c < bc; c++ {
			if m.At(r, c) != b.At(r, c) {
				return false
			}
		}
	}
	return true
}

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

// RankOne performs a rank-one update to the matrix b and stores the result
// in the receiver
//  m = a + alpha * x * y'
func (m *Dense) RankOne(a Matrix, alpha float64, x, y []float64) {
	ar, ac := a.Dims()

	var w Dense
	if m == a {
		w = *m
	}
	if w.isZero() {
		w.mat = RawMatrix{
			Rows:   ar,
			Cols:   ac,
			Stride: ac,
			Data:   use(w.mat.Data, ar*ac),
		}
	} else if ar != w.mat.Rows || ac != w.mat.Cols {
		panic(ErrShape)
	}
	// Copy over to the new memory if necessary
	if m != a {
		w.Copy(a)
	}
	if len(x) != ar {
		panic(ErrShape)
	}
	if len(y) != ac {
		panic(ErrShape)
	}
	blasEngine.Dger(ar, ac, alpha, x, 1, y, 1, w.mat.Data, w.mat.Stride)
	*m = w
	return
}
