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

var blasOrder = blas.RowMajor

func Order(o blas.Order) blas.Order {
	if o == blas.RowMajor || o == blas.ColMajor {
		o, blasOrder = blasOrder, o
		return o
	}
	return blasOrder
}

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

	_ Transposer = matrix
	// _ TransposeViewer = matrix

	// _ Deter  = matrix
	// _ Inver  = matrix
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
	var stride int
	switch blasOrder {
	case blas.RowMajor:
		stride = c
	case blas.ColMajor:
		stride = r
	default:
		panic(ErrIllegalOrder)
	}
	return &Dense{BlasMatrix{
		Order:  blasOrder,
		Rows:   r,
		Cols:   c,
		Stride: stride,
		Data:   mat,
	}}, nil
}

func (m *Dense) LoadBlas(b BlasMatrix) { m.mat = b }

func (m *Dense) BlasMatrix() BlasMatrix { return m.mat }

func (m *Dense) isZero() bool {
	return m.mat.Cols == 0 || m.mat.Rows == 0
}

func (m *Dense) At(r, c int) float64 {
	switch m.mat.Order {
	case blas.RowMajor:
		return m.mat.Data[r*m.mat.Stride+c]
	case blas.ColMajor:
		return m.mat.Data[c*m.mat.Stride+r]
	default:
		panic(ErrIllegalOrder)
	}
}

func (m *Dense) Set(r, c int, v float64) {
	switch m.mat.Order {
	case blas.RowMajor:
		m.mat.Data[r*m.mat.Stride+c] = v
	case blas.ColMajor:
		m.mat.Data[c*m.mat.Stride+r] = v
	default:
		panic(ErrIllegalOrder)
	}
}

func (m *Dense) Dims() (r, c int) { return m.mat.Rows, m.mat.Cols }

func (m *Dense) Col(col []float64, c int) []float64 {
	if c >= m.mat.Cols || c < 0 {
		panic(ErrIndexOutOfRange)
	}

	if col == nil {
		col = make([]float64, m.mat.Rows)
	}
	switch m.mat.Order {
	case blas.RowMajor:
		col = col[:min(len(col), m.mat.Rows)]
		blasEngine.Dcopy(len(col), m.mat.Data[c:], m.mat.Stride, col, 1)
	case blas.ColMajor:
		copy(col, m.mat.Data[c*m.mat.Stride:c*m.mat.Stride+m.mat.Rows])
	default:
		panic(ErrIllegalOrder)
	}

	return col
}

func (m *Dense) SetCol(c int, v []float64) int {
	if c >= m.mat.Cols || c < 0 {
		panic(ErrIndexOutOfRange)
	}

	switch m.mat.Order {
	case blas.RowMajor:
		blasEngine.Dcopy(min(len(v), m.mat.Rows), v, 1, m.mat.Data[c:], m.mat.Stride)
	case blas.ColMajor:
		copy(m.mat.Data[c*m.mat.Stride:c*m.mat.Stride+m.mat.Rows], v)
	default:
		panic(ErrIllegalOrder)
	}

	return min(len(v), m.mat.Rows)
}

func (m *Dense) Row(row []float64, r int) []float64 {
	if r >= m.mat.Rows || r < 0 {
		panic(ErrIndexOutOfRange)
	}

	if row == nil {
		row = make([]float64, m.mat.Cols)
	}
	switch m.mat.Order {
	case blas.RowMajor:
		copy(row, m.mat.Data[r*m.mat.Stride:r*m.mat.Stride+m.mat.Cols])
	case blas.ColMajor:
		row = row[:min(len(row), m.mat.Cols)]
		blasEngine.Dcopy(len(row), m.mat.Data[r:], m.mat.Stride, row, 1)
	default:
		panic(ErrIllegalOrder)
	}

	return row
}

func (m *Dense) SetRow(r int, v []float64) int {
	if r >= m.mat.Rows || r < 0 {
		panic(ErrIndexOutOfRange)
	}

	switch m.mat.Order {
	case blas.RowMajor:
		copy(m.mat.Data[r*m.mat.Stride:r*m.mat.Stride+m.mat.Cols], v)
	case blas.ColMajor:
		blasEngine.Dcopy(min(len(v), m.mat.Cols), v, 1, m.mat.Data[r:], m.mat.Stride)
	default:
		panic(ErrIllegalOrder)
	}

	return min(len(v), m.mat.Cols)
}

// View returns a view on the receiver.
func (m *Dense) View(i, j, r, c int) Blasser {
	v := Dense{BlasMatrix{
		Order:  m.mat.Order,
		Rows:   r - i,
		Cols:   c - j,
		Stride: m.mat.Stride,
	}}
	switch m.mat.Order {
	case blas.RowMajor:
		v.mat.Data = m.mat.Data[i*m.mat.Stride+j : (i+r-1)*m.mat.Stride+(j+c)]
	case blas.ColMajor:
		v.mat.Data = m.mat.Data[i+j*m.mat.Stride : (i+r)+(j+c-1)*m.mat.Stride]
	default:
		panic(ErrIllegalOrder)
	}
	return &v
}

func (m *Dense) Submatrix(a Matrix, i, j, r, c int) {
	// This is probably a bad idea, but for the moment, we do it.
	m.Clone(&Dense{m.View(i, j, r, c).BlasMatrix()})
}

func (m *Dense) Clone(a Matrix) {
	r, c := a.Dims()
	m.mat = BlasMatrix{
		Order: blasOrder,
		Rows:  r,
		Cols:  c,
	}
	data := make([]float64, r*c)
	switch a := a.(type) {
	case Blasser:
		amat := a.BlasMatrix()
		switch blasOrder {
		case blas.RowMajor:
			for i := 0; i < r; i++ {
				copy(data[i*c:(i+1)*c], amat.Data[i*amat.Stride:i*amat.Stride+c])
			}
			m.mat.Stride = c
		case blas.ColMajor:
			for i := 0; i < c; i++ {
				copy(data[i*r:(i+1)*r], amat.Data[i*amat.Stride:i*amat.Stride+r])
			}
			m.mat.Stride = r
		default:
			panic(ErrIllegalOrder)
		}
		m.mat.Data = data
	case Vectorer:
		switch blasOrder {
		case blas.RowMajor:
			for i := 0; i < r; i++ {
				a.Row(data[i*c:(i+1)*c], i)
			}
			m.mat.Stride = c
		case blas.ColMajor:
			for i := 0; i < c; i++ {
				a.Col(data[i*r:(i+1)*r], i)
			}
			m.mat.Stride = r
		default:
			panic(ErrIllegalOrder)
		}
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
		switch blasOrder {
		case blas.RowMajor:
			for i := 0; i < r; i++ {
				copy(m.mat.Data[i*m.mat.Stride:i*m.mat.Stride+c], amat.Data[i*amat.Stride:i*amat.Stride+c])
			}
		case blas.ColMajor:
			for i := 0; i < c; i++ {
				copy(m.mat.Data[i*m.mat.Stride:i*m.mat.Stride+r], amat.Data[i*amat.Stride:i*amat.Stride+r])
			}
		default:
			panic(ErrIllegalOrder)
		}
	case Vectorer:
		switch blasOrder {
		case blas.RowMajor:
			for i := 0; i < r; i++ {
				a.Row(m.mat.Data[i*m.mat.Stride:i*m.mat.Stride+c], i)
			}
		case blas.ColMajor:
			for i := 0; i < c; i++ {
				a.Col(m.mat.Data[i*m.mat.Stride:i*m.mat.Stride+r], i)
			}
		default:
			panic(ErrIllegalOrder)
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
	var i, j int
	switch m.mat.Order {
	case blas.RowMajor:
		i, j = m.mat.Rows, m.mat.Cols
	case blas.ColMajor:
		i, j = m.mat.Cols, m.mat.Rows
	default:
		panic(ErrIllegalOrder)
	}
	min := m.mat.Data[0]
	for k := 0; k < i; k++ {
		for _, v := range m.mat.Data[k*m.mat.Stride : k*m.mat.Stride+j] {
			min = math.Min(min, v)
		}
	}
	return min
}

func (m *Dense) Max() float64 {
	var i, j int
	switch m.mat.Order {
	case blas.RowMajor:
		i, j = m.mat.Rows, m.mat.Cols
	case blas.ColMajor:
		i, j = m.mat.Cols, m.mat.Rows
	default:
		panic(ErrIllegalOrder)
	}
	max := m.mat.Data[0]
	for k := 0; k < i; k++ {
		for _, v := range m.mat.Data[k*m.mat.Stride : k*m.mat.Stride+j] {
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
		var l int
		switch blasOrder {
		case blas.RowMajor:
			l = m.mat.Cols
		case blas.ColMajor:
			l = m.mat.Rows
		default:
			panic(ErrIllegalOrder)
		}
		for i := 0; i < len(m.mat.Data); i += m.mat.Stride {
			for _, v := range m.mat.Data[i : i+l] {
				n += v * v
			}
		}
		return math.Sqrt(n)
	case ord == 2, ord == -2:
		panic("matrix: 2-norm not implemented (pull requests for svd implementation welcomed)")
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

	var k, l int
	if m.isZero() {
		m.mat = BlasMatrix{
			Order: blasOrder,
			Rows:  ar,
			Cols:  ac,
			Data:  realloc(m.mat.Data, ar*ac),
		}
		switch blasOrder {
		case blas.RowMajor:
			m.mat.Stride, k, l = ac, ar, ac
		case blas.ColMajor:
			m.mat.Stride, k, l = ar, ac, ar
		default:
			panic(ErrIllegalOrder)
		}
	} else if ar != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
	} else {
		switch blasOrder {
		case blas.RowMajor:
			k, l = ar, ac
		case blas.ColMajor:
			k, l = ac, ar
		default:
			panic(ErrIllegalOrder)
		}
	}

	if a, ok := a.(Blasser); ok {
		if b, ok := b.(Blasser); ok {
			amat, bmat := a.BlasMatrix(), b.BlasMatrix()
			if amat.Order != blasOrder || bmat.Order != blasOrder {
				panic(ErrIllegalOrder)
			}
			for ja, jb, jm := 0, 0, 0; ja < k*amat.Stride; ja, jb, jm = ja+amat.Stride, jb+bmat.Stride, jm+m.mat.Stride {
				for i, v := range amat.Data[ja : ja+l] {
					m.mat.Data[i+jm] = v + bmat.Data[i+jb]
				}
			}
			return
		}
	}

	if a, ok := a.(Vectorer); ok {
		if b, ok := b.(Vectorer); ok {
			switch blasOrder {
			case blas.RowMajor:
				rowa := make([]float64, ac)
				rowb := make([]float64, bc)
				for r := 0; r < ar; r++ {
					a.Row(rowa, r)
					for i, v := range b.Row(rowb, r) {
						rowa[i] += v
					}
					copy(m.mat.Data[r*m.mat.Stride:r*m.mat.Stride+m.mat.Cols], rowa)
				}
			case blas.ColMajor:
				cola := make([]float64, ar)
				colb := make([]float64, br)
				for c := 0; c < ac; c++ {
					a.Col(cola, c)
					for i, v := range b.Col(colb, c) {
						cola[i] += v
					}
					copy(m.mat.Data[c*m.mat.Stride:c*m.mat.Stride+m.mat.Rows], cola)
				}
			default:
				panic(ErrIllegalOrder)
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

	var k, l int
	if m.isZero() {
		m.mat = BlasMatrix{
			Order: blasOrder,
			Rows:  ar,
			Cols:  ac,
			Data:  realloc(m.mat.Data, ar*ac),
		}
		switch blasOrder {
		case blas.RowMajor:
			m.mat.Stride, k, l = ac, ar, ac
		case blas.ColMajor:
			m.mat.Stride, k, l = ar, ac, ar
		default:
			panic(ErrIllegalOrder)
		}
	} else if ar != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
	} else {
		switch blasOrder {
		case blas.RowMajor:
			k, l = ar, ac
		case blas.ColMajor:
			k, l = ac, ar
		default:
			panic(ErrIllegalOrder)
		}
	}

	if a, ok := a.(Blasser); ok {
		if b, ok := b.(Blasser); ok {
			amat, bmat := a.BlasMatrix(), b.BlasMatrix()
			if amat.Order != blasOrder || bmat.Order != blasOrder {
				panic(ErrIllegalOrder)
			}
			for ja, jb, jm := 0, 0, 0; ja < k*amat.Stride; ja, jb, jm = ja+amat.Stride, jb+bmat.Stride, jm+m.mat.Stride {
				for i, v := range amat.Data[ja : ja+l] {
					m.mat.Data[i+jm] = v - bmat.Data[i+jb]
				}
			}
			return
		}
	}

	if a, ok := a.(Vectorer); ok {
		if b, ok := b.(Vectorer); ok {
			switch blasOrder {
			case blas.RowMajor:
				rowa := make([]float64, ac)
				rowb := make([]float64, bc)
				for r := 0; r < ar; r++ {
					a.Row(rowa, r)
					for i, v := range b.Row(rowb, r) {
						rowa[i] -= v
					}
					copy(m.mat.Data[r*m.mat.Stride:r*m.mat.Stride+m.mat.Cols], rowa)
				}
			case blas.ColMajor:
				cola := make([]float64, ar)
				colb := make([]float64, br)
				for c := 0; c < ac; c++ {
					a.Col(cola, c)
					for i, v := range b.Col(colb, c) {
						cola[i] -= v
					}
					copy(m.mat.Data[c*m.mat.Stride:c*m.mat.Stride+m.mat.Rows], cola)
				}
			default:
				panic(ErrIllegalOrder)
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

	var k, l int
	if m.isZero() {
		m.mat = BlasMatrix{
			Order: blasOrder,
			Rows:  ar,
			Cols:  ac,
			Data:  realloc(m.mat.Data, ar*ac),
		}
		switch blasOrder {
		case blas.RowMajor:
			m.mat.Stride, k, l = ac, ar, ac
		case blas.ColMajor:
			m.mat.Stride, k, l = ar, ac, ar
		default:
			panic(ErrIllegalOrder)
		}
	} else if ar != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
	} else {
		switch blasOrder {
		case blas.RowMajor:
			k, l = ar, ac
		case blas.ColMajor:
			k, l = ac, ar
		default:
			panic(ErrIllegalOrder)
		}
	}

	if a, ok := a.(Blasser); ok {
		if b, ok := b.(Blasser); ok {
			amat, bmat := a.BlasMatrix(), b.BlasMatrix()
			if amat.Order != blasOrder || bmat.Order != blasOrder {
				panic(ErrIllegalOrder)
			}
			for ja, jb, jm := 0, 0, 0; ja < k*amat.Stride; ja, jb, jm = ja+amat.Stride, jb+bmat.Stride, jm+m.mat.Stride {
				for i, v := range amat.Data[ja : ja+l] {
					m.mat.Data[i+jm] = v * bmat.Data[i+jb]
				}
			}
			return
		}
	}

	if a, ok := a.(Vectorer); ok {
		if b, ok := b.(Vectorer); ok {
			switch blasOrder {
			case blas.RowMajor:
				rowa := make([]float64, ac)
				rowb := make([]float64, bc)
				for r := 0; r < ar; r++ {
					a.Row(rowa, r)
					for i, v := range b.Row(rowb, r) {
						rowa[i] *= v
					}
					copy(m.mat.Data[r*m.mat.Stride:r*m.mat.Stride+m.mat.Cols], rowa)
				}
			case blas.ColMajor:
				cola := make([]float64, ar)
				colb := make([]float64, br)
				for c := 0; c < ac; c++ {
					a.Col(cola, c)
					for i, v := range b.Col(colb, c) {
						cola[i] *= v
					}
					copy(m.mat.Data[c*m.mat.Stride:c*m.mat.Stride+m.mat.Rows], cola)
				}
			default:
				panic(ErrIllegalOrder)
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

	var k, l int
	switch blasOrder {
	case blas.RowMajor:
		k, l = mr, mc
	case blas.ColMajor:
		k, l = mc, mr
	default:
		panic(ErrIllegalOrder)
	}

	var d float64

	if b, ok := b.(Blasser); ok {
		bmat := b.BlasMatrix()
		if m.mat.Order != blasOrder || bmat.Order != blasOrder {
			panic(ErrIllegalOrder)
		}
		for jm, jb := 0, 0; jm < k*m.mat.Stride; jm, jb = jm+m.mat.Stride, jb+bmat.Stride {
			for i, v := range m.mat.Data[jm : jm+l] {
				d += v * bmat.Data[i+jb]
			}
		}
		return d
	}

	if b, ok := b.(Vectorer); ok {
		switch blasOrder {
		case blas.RowMajor:
			row := make([]float64, bc)
			for r := 0; r < br; r++ {
				for i, v := range b.Row(row, r) {
					d += m.mat.Data[r*m.mat.Stride+i] * v
				}
			}
		case blas.ColMajor:
			col := make([]float64, br)
			for c := 0; c < bc; c++ {
				for i, v := range b.Col(col, c) {
					d += m.mat.Data[c*m.mat.Stride+i] * v
				}
			}
		default:
			panic(ErrIllegalOrder)
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
			Order: blasOrder,
			Rows:  ar,
			Cols:  bc,
			Data:  realloc(w.mat.Data, ar*bc),
		}
		switch blasOrder {
		case blas.RowMajor:
			w.mat.Stride = bc
		case blas.ColMajor:
			w.mat.Stride = br
		default:
			panic(ErrIllegalOrder)
		}
	} else if ar != w.mat.Rows || bc != w.mat.Cols {
		panic(ErrShape)
	}

	if a, ok := a.(Blasser); ok {
		if b, ok := b.(Blasser); ok {
			amat, bmat := a.BlasMatrix(), b.BlasMatrix()
			if amat.Order != blasOrder || bmat.Order != blasOrder {
				panic(ErrIllegalOrder)
			}
			blasEngine.Dgemm(
				blasOrder,
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
			for r := 0; r < ar; r++ {
				for c := 0; c < bc; c++ {
					switch blasOrder {
					case blas.RowMajor:
						w.mat.Data[r*w.mat.Stride+w.mat.Cols] = blasEngine.Ddot(ac, a.Row(row, r), 1, b.Col(col, c), 1)
					case blas.ColMajor:
						w.mat.Data[c*w.mat.Stride+w.mat.Rows] = blasEngine.Ddot(ac, a.Row(row, r), 1, b.Col(col, c), 1)
					default:
						panic(ErrIllegalOrder)
					}
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
			switch blasOrder {
			case blas.RowMajor:
				w.mat.Data[r*w.mat.Stride+w.mat.Cols] = v
			case blas.ColMajor:
				w.mat.Data[c*w.mat.Stride+w.mat.Rows] = v
			default:
				panic(ErrIllegalOrder)
			}
		}
	}
	*m = w
}

func (m *Dense) Scale(f float64, a Matrix) {
	ar, ac := a.Dims()

	var k, l int
	if m.isZero() {
		m.mat = BlasMatrix{
			Order: blasOrder,
			Rows:  ar,
			Cols:  ac,
			Data:  realloc(m.mat.Data, ar*ac),
		}
		switch blasOrder {
		case blas.RowMajor:
			m.mat.Stride, k, l = ac, ar, ac
		case blas.ColMajor:
			m.mat.Stride, k, l = ar, ac, ar
		default:
			panic(ErrIllegalOrder)
		}
	} else if ar != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
	} else {
		switch blasOrder {
		case blas.RowMajor:
			k, l = ar, ac
		case blas.ColMajor:
			k, l = ac, ar
		default:
			panic(ErrIllegalOrder)
		}
	}

	if a, ok := a.(Blasser); ok {
		amat := a.BlasMatrix()
		if amat.Order != blasOrder {
			panic(ErrIllegalOrder)
		}
		for ja, jm := 0, 0; ja < k*amat.Stride; ja, jm = ja+amat.Stride, jm+m.mat.Stride {
			for i, v := range amat.Data[ja : ja+l] {
				m.mat.Data[i+jm] = v * f
			}
		}
		return
	}

	if a, ok := a.(Vectorer); ok {
		switch blasOrder {
		case blas.RowMajor:
			row := make([]float64, ac)
			for r := 0; r < ar; r++ {
				for i, v := range a.Row(row, r) {
					row[i] = f * v
				}
				copy(m.mat.Data[r*m.mat.Stride:r*m.mat.Stride+m.mat.Cols], row)
			}
		case blas.ColMajor:
			col := make([]float64, ar)
			for c := 0; c < ac; c++ {
				for i, v := range a.Col(col, c) {
					col[i] = f * v
				}
				copy(m.mat.Data[c*m.mat.Stride:c*m.mat.Stride+m.mat.Rows], col)
			}
		default:
			panic(ErrIllegalOrder)
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

	var k, l int
	if m.isZero() {
		m.mat = BlasMatrix{
			Order: blasOrder,
			Rows:  ar,
			Cols:  ac,
			Data:  realloc(m.mat.Data, ar*ac),
		}
		switch blasOrder {
		case blas.RowMajor:
			m.mat.Stride, k, l = ac, ar, ac
		case blas.ColMajor:
			m.mat.Stride, k, l = ar, ac, ar
		default:
			panic(ErrIllegalOrder)
		}
	} else if ar != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
	} else {
		switch blasOrder {
		case blas.RowMajor:
			k, l = ar, ac
		case blas.ColMajor:
			k, l = ac, ar
		default:
			panic(ErrIllegalOrder)
		}
	}

	if a, ok := a.(Blasser); ok {
		amat := a.BlasMatrix()
		if amat.Order != blasOrder {
			panic(ErrIllegalOrder)
		}
		var r, c int
		for j, ja, jm := 0, 0, 0; ja < k*amat.Stride; j, ja, jm = j+1, ja+amat.Stride, jm+m.mat.Stride {
			for i, v := range amat.Data[ja : ja+l] {
				if blasOrder == blas.RowMajor {
					r, c = j, i
				} else {
					r, c = i, j
				}
				m.mat.Data[i+jm] = f(r, c, v)
			}
		}
		return
	}

	if a, ok := a.(Vectorer); ok {
		switch blasOrder {
		case blas.RowMajor:
			row := make([]float64, ac)
			for r := 0; r < ar; r++ {
				for i, v := range a.Row(row, r) {
					row[i] = f(r, i, v)
				}
				copy(m.mat.Data[r*m.mat.Stride:r*m.mat.Stride+m.mat.Cols], row)
			}
		case blas.ColMajor:
			col := make([]float64, ar)
			for c := 0; c < ac; c++ {
				for i, v := range a.Col(col, c) {
					col[i] = f(i, c, v)
				}
				copy(m.mat.Data[c*m.mat.Stride:c*m.mat.Stride+m.mat.Rows], col)
			}
		default:
			panic(ErrIllegalOrder)
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

	var k, l int
	switch {
	case m == a:
		m.zeroLower()
		return
	case m.isZero():
		m.mat = BlasMatrix{
			Order: blasOrder,
			Rows:  ar,
			Cols:  ac,
			Data:  realloc(m.mat.Data, ar*ac),
		}
		switch blasOrder {
		case blas.RowMajor:
			m.mat.Stride, k, l = ac, ar, ac
		case blas.ColMajor:
			m.mat.Stride, k, l = ar, ac, ar
		default:
			panic(ErrIllegalOrder)
		}
	case ar != m.mat.Rows || ac != m.mat.Cols:
		panic(ErrShape)
	default:
		switch blasOrder {
		case blas.RowMajor:
			k, l = ar, ac
		case blas.ColMajor:
			k, l = ac, ar
		default:
			panic(ErrIllegalOrder)
		}
	}

	if a, ok := a.(Blasser); ok {
		amat := a.BlasMatrix()
		if amat.Order != blasOrder {
			panic(ErrIllegalOrder)
		}
		copy(m.mat.Data[:l], amat.Data[:l])
		for j, ja, jm := 1, amat.Stride, m.mat.Stride; ja < k*amat.Stride; j, ja, jm = j+1, ja+amat.Stride, jm+m.mat.Stride {
			zero(m.mat.Data[jm : jm+j])
			copy(m.mat.Data[jm+j:jm+l], amat.Data[ja+j:ja+l])
		}
		return
	}

	if a, ok := a.(Vectorer); ok {
		switch blasOrder {
		case blas.RowMajor:
			row := make([]float64, ac)
			copy(m.mat.Data[:m.mat.Cols], a.Row(row, 0))
			for r := 1; r < ar; r++ {
				zero(m.mat.Data[r*m.mat.Stride : r*(m.mat.Stride+1)])
				copy(m.mat.Data[r*(m.mat.Stride+1):r*m.mat.Stride+m.mat.Cols], a.Row(row, r))
			}
		case blas.ColMajor:
			col := make([]float64, ar)
			for c := 0; c < ac; c++ {
				a.Col(col[:c+1], c)
				m.SetCol(c, col)
			}
		default:
			panic(ErrIllegalOrder)
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
	switch blasOrder {
	case blas.RowMajor:
		for i := 1; i < m.mat.Rows; i++ {
			zero(m.mat.Data[i*m.mat.Stride : i*m.mat.Stride+i])
		}
	case blas.ColMajor:
		for i := 0; i < m.mat.Cols-1; i++ {
			zero(m.mat.Data[i*m.mat.Stride+i+1 : (i+1)*m.mat.Stride])
		}
	default:
		panic(ErrIllegalOrder)
	}
}

func (m *Dense) L(a Matrix) {
	ar, ac := a.Dims()
	if ar != ac {
		panic(ErrSquare)
	}

	var k, l int
	switch {
	case m == a:
		m.zeroUpper()
		return
	case m.isZero():
		m.mat = BlasMatrix{
			Order: blasOrder,
			Rows:  ar,
			Cols:  ac,
			Data:  realloc(m.mat.Data, ar*ac),
		}
		switch blasOrder {
		case blas.RowMajor:
			m.mat.Stride, k, l = ac, ar, ac
		case blas.ColMajor:
			m.mat.Stride, k, l = ar, ac, ar
		default:
			panic(ErrIllegalOrder)
		}
	case ar != m.mat.Rows || ac != m.mat.Cols:
		panic(ErrShape)
	default:
		switch blasOrder {
		case blas.RowMajor:
			k, l = ac, ar
		case blas.ColMajor:
			k, l = ar, ac
		default:
			panic(ErrIllegalOrder)
		}
	}

	if a, ok := a.(Blasser); ok {
		amat := a.BlasMatrix()
		if amat.Order != blasOrder {
			panic(ErrIllegalOrder)
		}
		copy(m.mat.Data[:l], amat.Data[:l])
		for j, ja, jm := 1, amat.Stride, m.mat.Stride; ja < k*amat.Stride; j, ja, jm = j+1, ja+amat.Stride, jm+m.mat.Stride {
			zero(m.mat.Data[jm : jm+j])
			copy(m.mat.Data[jm+j:jm+l], amat.Data[ja+j:ja+l])
		}
		return
	}

	if a, ok := a.(Vectorer); ok {
		switch blasOrder {
		case blas.RowMajor:
			row := make([]float64, ac)
			for r := 0; r < ar; r++ {
				a.Row(row[:r+1], r)
				m.SetRow(r, row)
			}
		case blas.ColMajor:
			col := make([]float64, ar)
			copy(m.mat.Data[:m.mat.Rows], a.Col(col, 0))
			for c := 1; c < ac; c++ {
				zero(m.mat.Data[c*m.mat.Stride : c*(m.mat.Stride+1)])
				copy(m.mat.Data[c*(m.mat.Stride+1):c*m.mat.Stride+m.mat.Rows], a.Col(col, c))
			}
		default:
			panic(ErrIllegalOrder)
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
	switch blasOrder {
	case blas.RowMajor:
		for i := 1; i < m.mat.Rows-1; i++ {
			zero(m.mat.Data[i*m.mat.Stride+i+1 : (i+1)*m.mat.Stride])
		}
	case blas.ColMajor:
		for i := 0; i < m.mat.Cols-1; i++ {
			zero(m.mat.Data[i*m.mat.Stride : i*m.mat.Stride+i])
		}
	default:
		panic(ErrIllegalOrder)
	}
}

func (m *Dense) T(a Matrix) {
	ar, ac := a.Dims()

	var w Dense
	if m != a {
		w = *m
	}
	if w.isZero() {
		w.mat = BlasMatrix{
			Order: blasOrder,
			Rows:  ac,
			Cols:  ar,
			Data:  realloc(w.mat.Data, ar*ac),
		}
		switch blasOrder {
		case blas.RowMajor:
			w.mat.Stride = ar
		case blas.ColMajor:
			w.mat.Stride = ac
		default:
			panic(ErrIllegalOrder)
		}
	} else if ar != m.mat.Cols || ac != m.mat.Rows {
		panic(ErrShape)
	} else if blasOrder != blas.RowMajor && blasOrder != blas.ColMajor {
		panic(ErrIllegalOrder)
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
	var l int
	switch blasOrder {
	case blas.RowMajor:
		l = m.mat.Cols
	case blas.ColMajor:
		l = m.mat.Rows
	default:
		panic(ErrIllegalOrder)
	}
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
		var k, l int
		switch blasOrder {
		case blas.RowMajor:
			k, l = br, bc
		case blas.ColMajor:
			k, l = bc, br
		default:
			panic(ErrIllegalOrder)
		}
		bmat := b.BlasMatrix()
		for jb, jm := 0, 0; jm < k*m.mat.Stride; jb, jm = jb+bmat.Stride, jm+m.mat.Stride {
			for i, v := range m.mat.Data[jm : jm+l] {
				if v != bmat.Data[i+jb] {
					return false
				}
			}
		}
		return true
	}

	if b, ok := b.(Vectorer); ok {
		switch blasOrder {
		case blas.RowMajor:
			rowb := make([]float64, bc)
			for r := 0; r < br; r++ {
				rowm := m.mat.Data[r*m.mat.Stride : r*m.mat.Stride+m.mat.Cols]
				for i, v := range b.Row(rowb, r) {
					if rowm[i] != v {
						return false
					}
				}
			}
		case blas.ColMajor:
			colb := make([]float64, br)
			for c := 0; c < bc; c++ {
				colm := m.mat.Data[c*m.mat.Stride : c*m.mat.Stride+m.mat.Rows]
				for i, v := range b.Col(colb, c) {
					if colm[i] != v {
						return false
					}
				}
			}
		default:
			panic(ErrIllegalOrder)
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
		var k, l int
		switch blasOrder {
		case blas.RowMajor:
			k, l = br, bc
		case blas.ColMajor:
			k, l = bc, br
		default:
			panic(ErrIllegalOrder)
		}
		bmat := b.BlasMatrix()
		for jb, jm := 0, 0; jm < k*m.mat.Stride; jb, jm = jb+bmat.Stride, jm+m.mat.Stride {
			for i, v := range m.mat.Data[jm : jm+l] {
				if math.Abs(v-bmat.Data[i+jb]) > epsilon {
					return false
				}
			}
		}
		return true
	}

	if b, ok := b.(Vectorer); ok {
		switch blasOrder {
		case blas.RowMajor:
			rowb := make([]float64, bc)
			for r := 0; r < br; r++ {
				rowm := m.mat.Data[r*m.mat.Stride : r*m.mat.Stride+m.mat.Cols]
				for i, v := range b.Row(rowb, r) {
					if math.Abs(rowm[i]-v) > epsilon {
						return false
					}
				}
			}
		case blas.ColMajor:
			colb := make([]float64, br)
			for c := 0; c < bc; c++ {
				colm := m.mat.Data[c*m.mat.Stride : c*m.mat.Stride+m.mat.Rows]
				for i, v := range b.Col(colb, c) {
					if math.Abs(colm[i]-v) > epsilon {
						return false
					}
				}
			}
		default:
			panic(ErrIllegalOrder)
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
