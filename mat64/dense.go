// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"github.com/gonum/blas"
	"math"
)

var blasEngine blas.Float64

func Register(b blas.Float64) { blasEngine = b }

func Registered() blas.Float64 { return blasEngine }

const BlasOrder = blas.RowMajor

var (
	matrix *Dense

	_ Matrix       = matrix
	_ Mutable      = matrix
	_ Vectorer     = matrix
	_ VectorSetter = matrix

	_ Cloner      = matrix
	_ Viewer      = matrix
	_ Submatrixer = matrix

	_ Adder     = matrix
	_ Suber     = matrix
	_ Muler     = matrix
	_ Dotter    = matrix
	_ ElemMuler = matrix

	_ Scaler  = matrix
	_ Applyer = matrix

	_ TransposeCopier = matrix
	// _ TransposeViewer = matrix

	_ Tracer = matrix
	_ Normer = matrix
	_ Sumer  = matrix

	_ Uer = matrix
	_ Ler = matrix

	// _ Stacker   = matrix
	// _ Augmenter = matrix

	_ Equaler       = matrix
	_ ApproxEqualer = matrix

	_ BlasLoader = matrix
	_ Blasser    = matrix
)

type Dense struct {
	mat BlasMatrix
}

func NewDense(r, c int, mat []float64) (*Dense, error) {
	if r*c != len(mat) {
		return nil, ErrShape
	}
	return &Dense{BlasMatrix{
		Order:  BlasOrder,
		Rows:   r,
		Cols:   c,
		Stride: c,
		Data:   mat,
	}}, nil
}

// DenseCopyOf returns a newly allocated copy of the elements of a.
func DenseCopyOf(a Matrix) *Dense {
	d := &Dense{}
	d.Clone(a)
	return d
}

func (m *Dense) LoadBlas(b BlasMatrix) {
	if b.Order != BlasOrder {
		panic(ErrIllegalOrder)
	}
	m.mat = b
}

func (m *Dense) BlasMatrix() BlasMatrix { return m.mat }

func (m *Dense) isZero() bool {
	return m.mat.Cols == 0 || m.mat.Rows == 0
}

func (m *Dense) At(r, c int) float64 {
	return m.mat.Data[r*m.mat.Stride+c]
}

func (m *Dense) Set(r, c int, v float64) {
	m.mat.Data[r*m.mat.Stride+c] = v
}

func (m *Dense) Dims() (r, c int) { return m.mat.Rows, m.mat.Cols }

func (m *Dense) Col(col []float64, c int) []float64 {
	if c >= m.mat.Cols || c < 0 {
		panic(ErrIndexOutOfRange)
	}

	if col == nil {
		col = make([]float64, m.mat.Rows)
	}
	col = col[:min(len(col), m.mat.Rows)]
	if blasEngine == nil {
		panic(ErrNoEngine)
	}
	blasEngine.Dcopy(len(col), m.mat.Data[c:], m.mat.Stride, col, 1)

	return col
}

func (m *Dense) SetCol(c int, v []float64) int {
	if c >= m.mat.Cols || c < 0 {
		panic(ErrIndexOutOfRange)
	}

	if blasEngine == nil {
		panic(ErrNoEngine)
	}
	blasEngine.Dcopy(min(len(v), m.mat.Rows), v, 1, m.mat.Data[c:], m.mat.Stride)

	return min(len(v), m.mat.Rows)
}

func (m *Dense) Row(row []float64, r int) []float64 {
	if r >= m.mat.Rows || r < 0 {
		panic(ErrIndexOutOfRange)
	}

	if row == nil {
		row = make([]float64, m.mat.Cols)
	}
	copy(row, m.mat.Data[r*m.mat.Stride:r*m.mat.Stride+m.mat.Cols])

	return row
}

func (m *Dense) SetRow(r int, v []float64) int {
	if r >= m.mat.Rows || r < 0 {
		panic(ErrIndexOutOfRange)
	}

	copy(m.mat.Data[r*m.mat.Stride:r*m.mat.Stride+m.mat.Cols], v)

	return min(len(v), m.mat.Cols)
}

// View returns a view on the receiver.
func (m *Dense) View(i, j, r, c int) {
	m.mat.Data = m.mat.Data[i*m.mat.Stride+j : (i+r-1)*m.mat.Stride+(j+c)]
	m.mat.Rows = r
	m.mat.Cols = c
}

func (m *Dense) Submatrix(a Matrix, i, j, r, c int) {
	// This is probably a bad idea, but for the moment, we do it.
	v := *m
	v.View(i, j, r, c)
	m.Clone(&Dense{v.BlasMatrix()})
}

func (m *Dense) Clone(a Matrix) {
	r, c := a.Dims()
	m.mat = BlasMatrix{
		Order: BlasOrder,
		Rows:  r,
		Cols:  c,
	}
	data := make([]float64, r*c)
	switch a := a.(type) {
	case Blasser:
		amat := a.BlasMatrix()
		for i := 0; i < r; i++ {
			copy(data[i*c:(i+1)*c], amat.Data[i*amat.Stride:i*amat.Stride+c])
		}
		m.mat.Stride = c
		m.mat.Data = data
	case Vectorer:
		for i := 0; i < r; i++ {
			a.Row(data[i*c:(i+1)*c], i)
		}
		m.mat.Stride = c
		m.mat.Data = data
	default:
		m.mat.Data = data
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				m.Set(i, j, a.At(i, j))
			}
		}
	}
}

func (m *Dense) Copy(a Matrix) (r, c int) {
	r, c = a.Dims()
	r = min(r, m.mat.Rows)
	c = min(c, m.mat.Cols)

	switch a := a.(type) {
	case Blasser:
		amat := a.BlasMatrix()
		for i := 0; i < r; i++ {
			copy(m.mat.Data[i*m.mat.Stride:i*m.mat.Stride+c], amat.Data[i*amat.Stride:i*amat.Stride+c])
		}
	case Vectorer:
		for i := 0; i < r; i++ {
			a.Row(m.mat.Data[i*m.mat.Stride:i*m.mat.Stride+c], i)
		}
	default:
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				m.Set(r, c, a.At(r, c))
			}
		}
	}

	return r, c
}

func (m *Dense) Min() float64 {
	min := m.mat.Data[0]
	for k := 0; k < m.mat.Rows; k++ {
		for _, v := range m.mat.Data[k*m.mat.Stride : k*m.mat.Stride+m.mat.Cols] {
			min = math.Min(min, v)
		}
	}
	return min
}

func (m *Dense) Max() float64 {
	max := m.mat.Data[0]
	for k := 0; k < m.mat.Rows; k++ {
		for _, v := range m.mat.Data[k*m.mat.Stride : k*m.mat.Stride+m.mat.Cols] {
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

// Norm(±2) depends on SVD, and so m must be tall or square.
func (m *Dense) Norm(ord float64) float64 {
	var n float64
	switch {
	case ord == 1:
		col := make([]float64, m.mat.Rows)
		for i := 0; i < m.mat.Cols; i++ {
			var s float64
			for _, e := range m.Col(col, i) {
				s += e
			}
			n = math.Max(math.Abs(s), n)
		}
	case math.IsInf(ord, +1):
		row := make([]float64, m.mat.Cols)
		for i := 0; i < m.mat.Rows; i++ {
			var s float64
			for _, e := range m.Row(row, i) {
				s += e
			}
			n = math.Max(math.Abs(s), n)
		}
	case ord == -1:
		n = math.MaxFloat64
		col := make([]float64, m.mat.Rows)
		for i := 0; i < m.mat.Cols; i++ {
			var s float64
			for _, e := range m.Col(col, i) {
				s += e
			}
			n = math.Min(math.Abs(s), n)
		}
	case math.IsInf(ord, -1):
		n = math.MaxFloat64
		row := make([]float64, m.mat.Cols)
		for i := 0; i < m.mat.Rows; i++ {
			var s float64
			for _, e := range m.Row(row, i) {
				s += e
			}
			n = math.Min(math.Abs(s), n)
		}
	case ord == 0:
		for i := 0; i < len(m.mat.Data); i += m.mat.Stride {
			for _, v := range m.mat.Data[i : i+m.mat.Cols] {
				n += v * v
			}
		}
		return math.Sqrt(n)
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
		m.mat = BlasMatrix{
			Order:  BlasOrder,
			Rows:   ar,
			Cols:   ac,
			Stride: ac,
			Data:   realloc(m.mat.Data, ar*ac),
		}
	} else if ar != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
	}

	if a, ok := a.(Blasser); ok {
		if b, ok := b.(Blasser); ok {
			amat, bmat := a.BlasMatrix(), b.BlasMatrix()
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
				copy(m.mat.Data[r*m.mat.Stride:r*m.mat.Stride+m.mat.Cols], rowa)
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
		m.mat = BlasMatrix{
			Order:  BlasOrder,
			Rows:   ar,
			Cols:   ac,
			Stride: ac,
			Data:   realloc(m.mat.Data, ar*ac),
		}
	} else if ar != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
	}

	if a, ok := a.(Blasser); ok {
		if b, ok := b.(Blasser); ok {
			amat, bmat := a.BlasMatrix(), b.BlasMatrix()
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
				copy(m.mat.Data[r*m.mat.Stride:r*m.mat.Stride+m.mat.Cols], rowa)
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
		m.mat = BlasMatrix{
			Order:  BlasOrder,
			Rows:   ar,
			Cols:   ac,
			Stride: ac,
			Data:   realloc(m.mat.Data, ar*ac),
		}
	} else if ar != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
	}

	if a, ok := a.(Blasser); ok {
		if b, ok := b.(Blasser); ok {
			amat, bmat := a.BlasMatrix(), b.BlasMatrix()
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
				copy(m.mat.Data[r*m.mat.Stride:r*m.mat.Stride+m.mat.Cols], rowa)
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

	if b, ok := b.(Blasser); ok {
		bmat := b.BlasMatrix()
		if m.mat.Order != BlasOrder || bmat.Order != BlasOrder {
			panic(ErrIllegalOrder)
		}
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
		w.mat = BlasMatrix{
			Order:  BlasOrder,
			Rows:   ar,
			Cols:   bc,
			Stride: bc,
			Data:   realloc(w.mat.Data, ar*bc),
		}
	} else if ar != w.mat.Rows || bc != w.mat.Cols {
		panic(ErrShape)
	}

	if a, ok := a.(Blasser); ok {
		if b, ok := b.(Blasser); ok {
			amat, bmat := a.BlasMatrix(), b.BlasMatrix()
			if blasEngine == nil {
				panic(ErrNoEngine)
			}
			blasEngine.Dgemm(
				BlasOrder,
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
					w.mat.Data[r*w.mat.Stride+w.mat.Cols] = blasEngine.Ddot(ac, a.Row(row, r), 1, b.Col(col, c), 1)
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
			w.mat.Data[r*w.mat.Stride+w.mat.Cols] = v
		}
	}
	*m = w
}

func (m *Dense) Scale(f float64, a Matrix) {
	ar, ac := a.Dims()

	if m.isZero() {
		m.mat = BlasMatrix{
			Order:  BlasOrder,
			Rows:   ar,
			Cols:   ac,
			Stride: ac,
			Data:   realloc(m.mat.Data, ar*ac),
		}
	} else if ar != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
	}

	if a, ok := a.(Blasser); ok {
		amat := a.BlasMatrix()
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
			copy(m.mat.Data[r*m.mat.Stride:r*m.mat.Stride+m.mat.Cols], row)
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
		m.mat = BlasMatrix{
			Order:  BlasOrder,
			Rows:   ar,
			Cols:   ac,
			Stride: ac,
			Data:   realloc(m.mat.Data, ar*ac),
		}
	} else if ar != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
	}

	if a, ok := a.(Blasser); ok {
		amat := a.BlasMatrix()
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
			copy(m.mat.Data[r*m.mat.Stride:r*m.mat.Stride+m.mat.Cols], row)
		}
		return
	}

	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.Set(r, c, f(r, c, a.At(r, c)))
		}
	}
}

func zero(f []float64) {
	f[0] = 0
	for i := 1; i < len(f); {
		i += copy(f[i:], f[:i])
	}
}

func (m *Dense) U(a Matrix) {
	ar, ac := a.Dims()
	if ar != ac {
		panic(ErrSquare)
	}

	switch {
	case m == a:
		m.zeroLower()
		return
	case m.isZero():
		m.mat = BlasMatrix{
			Order:  BlasOrder,
			Rows:   ar,
			Cols:   ac,
			Stride: ac,
			Data:   realloc(m.mat.Data, ar*ac),
		}
	case ar != m.mat.Rows || ac != m.mat.Cols:
		panic(ErrShape)
	}

	if a, ok := a.(Blasser); ok {
		amat := a.BlasMatrix()
		copy(m.mat.Data[:ac], amat.Data[:ac])
		for j, ja, jm := 1, amat.Stride, m.mat.Stride; ja < ar*amat.Stride; j, ja, jm = j+1, ja+amat.Stride, jm+m.mat.Stride {
			zero(m.mat.Data[jm : jm+j])
			copy(m.mat.Data[jm+j:jm+ac], amat.Data[ja+j:ja+ac])
		}
		return
	}

	if a, ok := a.(Vectorer); ok {
		row := make([]float64, ac)
		copy(m.mat.Data[:m.mat.Cols], a.Row(row, 0))
		for r := 1; r < ar; r++ {
			zero(m.mat.Data[r*m.mat.Stride : r*(m.mat.Stride+1)])
			copy(m.mat.Data[r*(m.mat.Stride+1):r*m.mat.Stride+m.mat.Cols], a.Row(row, r))
		}
		return
	}

	m.zeroLower()
	for r := 0; r < ar; r++ {
		for c := r; c < ac; c++ {
			m.Set(r, c, a.At(r, c))
		}
	}
}

func (m *Dense) zeroLower() {
	for i := 1; i < m.mat.Rows; i++ {
		zero(m.mat.Data[i*m.mat.Stride : i*m.mat.Stride+i])
	}
}

func (m *Dense) L(a Matrix) {
	ar, ac := a.Dims()
	if ar != ac {
		panic(ErrSquare)
	}

	switch {
	case m == a:
		m.zeroUpper()
		return
	case m.isZero():
		m.mat = BlasMatrix{
			Order:  BlasOrder,
			Rows:   ar,
			Cols:   ac,
			Stride: ac,
			Data:   realloc(m.mat.Data, ar*ac),
		}
	case ar != m.mat.Rows || ac != m.mat.Cols:
		panic(ErrShape)
	}

	if a, ok := a.(Blasser); ok {
		amat := a.BlasMatrix()
		copy(m.mat.Data[:ar], amat.Data[:ar])
		for j, ja, jm := 1, amat.Stride, m.mat.Stride; ja < ac*amat.Stride; j, ja, jm = j+1, ja+amat.Stride, jm+m.mat.Stride {
			zero(m.mat.Data[jm : jm+j])
			copy(m.mat.Data[jm+j:jm+ar], amat.Data[ja+j:ja+ar])
		}
		return
	}

	if a, ok := a.(Vectorer); ok {
		row := make([]float64, ac)
		for r := 0; r < ar; r++ {
			a.Row(row[:r+1], r)
			m.SetRow(r, row)
		}
		return
	}

	m.zeroUpper()
	for c := 0; c < ac; c++ {
		for r := c; r < ar; r++ {
			m.Set(r, c, a.At(r, c))
		}
	}
}

func (m *Dense) zeroUpper() {
	for i := 0; i < m.mat.Rows-1; i++ {
		zero(m.mat.Data[i*m.mat.Stride+i+1 : (i+1)*m.mat.Stride])
	}
}

func (m *Dense) TCopy(a Matrix) {
	ar, ac := a.Dims()

	var w Dense
	if m != a {
		w = *m
	}
	if w.isZero() {
		w.mat = BlasMatrix{
			Order: BlasOrder,
			Rows:  ac,
			Cols:  ar,
			Data:  realloc(w.mat.Data, ar*ac),
		}
		w.mat.Stride = ar
	} else if ar != m.mat.Cols || ac != m.mat.Rows {
		panic(ErrShape)
	}
	switch a := a.(type) {
	case *Dense:
		for i := 0; i < ac; i++ {
			for j := 0; j < ar; j++ {
				w.Set(i, j, a.At(j, i))
			}
		}
	default:
		for i := 0; i < ac; i++ {
			for j := 0; j < ar; j++ {
				w.Set(i, j, a.At(j, i))
			}
		}
	}
	*m = w
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

	if b, ok := b.(Blasser); ok {
		bmat := b.BlasMatrix()
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

	if b, ok := b.(Blasser); ok {
		bmat := b.BlasMatrix()
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
