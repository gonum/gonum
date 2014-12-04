// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import "github.com/gonum/blas"

var blasEngine blas.Float64

func Register(b blas.Float64) { blasEngine = b }

func Registered() blas.Float64 { return blasEngine }

var (
	matrix *Dense

	_ Matrix       = matrix
	_ Mutable      = matrix
	_ Vectorer     = matrix
	_ VectorSetter = matrix

	_ Cloner    = matrix
	_ Viewer    = matrix
	_ RowViewer = matrix

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

	_ Stacker   = matrix
	_ Augmenter = matrix

	_ Equaler       = matrix
	_ ApproxEqualer = matrix

	_ RawMatrixLoader = matrix
	_ RawMatrixer     = matrix
)

type Dense struct {
	mat RawMatrix
}

func NewDense(r, c int, mat []float64) *Dense {
	if mat != nil && r*c != len(mat) {
		panic(ErrShape)
	}
	if mat == nil {
		mat = make([]float64, r*c)
	}
	return &Dense{RawMatrix{
		Rows:   r,
		Cols:   c,
		Stride: c,
		Data:   mat,
	}}
}

// DenseCopyOf returns a newly allocated copy of the elements of a.
func DenseCopyOf(a Matrix) *Dense {
	d := &Dense{}
	d.Clone(a)
	return d
}

func (m *Dense) LoadRawMatrix(b RawMatrix) { m.mat = b }

func (m *Dense) RawMatrix() RawMatrix { return m.mat }

func (m *Dense) isZero() bool {
	return m.mat.Cols == 0 || m.mat.Rows == 0
}

func (m *Dense) At(r, c int) float64 {
	if r >= m.mat.Rows || r < 0 {
		panic("index error: row access out of bounds")
	}
	if c >= m.mat.Cols || c < 0 {
		panic("index error: column access out of bounds")
	}
	return m.at(r, c)
}

func (m *Dense) at(r, c int) float64 {
	return m.mat.Data[r*m.mat.Stride+c]
}

func (m *Dense) Set(r, c int, v float64) {
	if r >= m.mat.Rows || r < 0 {
		panic("index error: row access out of bounds")
	}
	if c >= m.mat.Cols || c < 0 {
		panic("index error: column access out of bounds")
	}
	m.mat.Data[r*m.mat.Stride+c] = v
}

func (m *Dense) Dims() (r, c int) { return m.mat.Rows, m.mat.Cols }

func (m *Dense) Col(dst []float64, j int) []float64 {
	if j >= m.mat.Cols || j < 0 {
		panic(ErrIndexOutOfRange)
	}

	if dst == nil {
		dst = make([]float64, m.mat.Rows)
	}
	dst = dst[:min(len(dst), m.mat.Rows)]
	if blasEngine == nil {
		panic(ErrNoEngine)
	}
	blasEngine.Dcopy(len(dst), m.mat.Data[j:], m.mat.Stride, dst, 1)

	return dst
}

func (m *Dense) SetCol(j int, src []float64) int {
	if j >= m.mat.Cols || j < 0 {
		panic(ErrIndexOutOfRange)
	}

	if blasEngine == nil {
		panic(ErrNoEngine)
	}
	blasEngine.Dcopy(min(len(src), m.mat.Rows), src, 1, m.mat.Data[j:], m.mat.Stride)

	return min(len(src), m.mat.Rows)
}

func (m *Dense) Row(dst []float64, i int) []float64 {
	if i >= m.mat.Rows || i < 0 {
		panic(ErrIndexOutOfRange)
	}

	if dst == nil {
		dst = make([]float64, m.mat.Cols)
	}
	copy(dst, m.rowView(i))

	return dst
}

func (m *Dense) SetRow(i int, src []float64) int {
	if i >= m.mat.Rows || i < 0 {
		panic(ErrIndexOutOfRange)
	}

	copy(m.rowView(i), src)

	return min(len(src), m.mat.Cols)
}

func (m *Dense) RowView(r int) []float64 {
	if r >= m.mat.Rows || r < 0 {
		panic(ErrIndexOutOfRange)
	}
	return m.rowView(r)
}

func (m *Dense) rowView(r int) []float64 {
	return m.mat.Data[r*m.mat.Stride : r*m.mat.Stride+m.mat.Cols]
}

func (a *Dense) View(i, j, r, c int) Matrix {
	m := *a
	m.mat.Data = m.mat.Data[i*m.mat.Stride+j : (i+r-1)*m.mat.Stride+(j+c)]
	m.mat.Rows = r
	m.mat.Cols = c
	return &m
}

func (m *Dense) Reset() {
	m.mat.Rows, m.mat.Cols = 0, 0
	m.mat.Data = m.mat.Data[:0]
}

func (m *Dense) Clone(a Matrix) {
	r, c := a.Dims()
	mat := RawMatrix{
		Rows:   r,
		Cols:   c,
		Stride: c,
	}
	switch a := a.(type) {
	case RawMatrixer:
		amat := a.RawMatrix()
		mat.Data = make([]float64, r*c)
		for i := 0; i < r; i++ {
			copy(mat.Data[i*c:(i+1)*c], amat.Data[i*amat.Stride:i*amat.Stride+c])
		}
	case Vectorer:
		mat.Data = use(m.mat.Data, r*c)
		for i := 0; i < r; i++ {
			a.Row(mat.Data[i*c:(i+1)*c], i)
		}
	default:
		mat.Data = use(m.mat.Data, r*c)
		m.mat = mat
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				m.Set(i, j, a.At(i, j))
			}
		}
		return
	}
	m.mat = mat
}

func (m *Dense) Copy(a Matrix) (r, c int) {
	r, c = a.Dims()
	r = min(r, m.mat.Rows)
	c = min(c, m.mat.Cols)

	switch a := a.(type) {
	case RawMatrixer:
		amat := a.RawMatrix()
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
		m.mat = RawMatrix{
			Rows:   ar,
			Cols:   ac,
			Stride: ac,
			Data:   use(m.mat.Data, ar*ac),
		}
	case ar != m.mat.Rows || ac != m.mat.Cols:
		panic(ErrShape)
	}

	if a, ok := a.(RawMatrixer); ok {
		amat := a.RawMatrix()
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
		m.mat = RawMatrix{
			Rows:   ar,
			Cols:   ac,
			Stride: ac,
			Data:   use(m.mat.Data, ar*ac),
		}
	case ar != m.mat.Rows || ac != m.mat.Cols:
		panic(ErrShape)
	}

	if a, ok := a.(RawMatrixer); ok {
		amat := a.RawMatrix()
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
		w.mat = RawMatrix{
			Rows: ac,
			Cols: ar,
			Data: use(w.mat.Data, ar*ac),
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

func (m *Dense) Stack(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ac != bc || m == a || m == b {
		panic(ErrShape)
	}

	if m.isZero() {
		m.mat = RawMatrix{
			Rows:   ar + br,
			Cols:   ac,
			Stride: ac,
			Data:   use(m.mat.Data, (ar+br)*ac),
		}
	} else if ar+br != m.mat.Rows || ac != m.mat.Cols {
		panic(ErrShape)
	}

	m.Copy(a)
	w := m.View(ar, 0, br, bc).(*Dense)
	w.Copy(b)
}

func (m *Dense) Augment(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || m == a || m == b {
		panic(ErrShape)
	}

	if m.isZero() {
		m.mat = RawMatrix{
			Rows:   ar,
			Cols:   ac + bc,
			Stride: ac + bc,
			Data:   use(m.mat.Data, ar*(ac+bc)),
		}
	} else if ar != m.mat.Rows || ac+bc != m.mat.Cols {
		panic(ErrShape)
	}

	m.Copy(a)
	w := m.View(0, ac, br, bc).(*Dense)
	w.Copy(b)
}
