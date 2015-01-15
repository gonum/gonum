// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"bytes"
	"encoding/binary"

	"github.com/gonum/blas/blas64"
)

var (
	matrix *Dense

	_ Matrix       = matrix
	_ Mutable      = matrix
	_ Vectorer     = matrix
	_ VectorSetter = matrix

	_ Cloner    = matrix
	_ Viewer    = matrix
	_ RowViewer = matrix
	_ Grower    = matrix

	_ Adder     = matrix
	_ Suber     = matrix
	_ Muler     = matrix
	_ Dotter    = matrix
	_ ElemMuler = matrix
	_ ElemDiver = matrix
	_ Exper     = matrix

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

	_ RawMatrixSetter = matrix
	_ RawMatrixer     = matrix

	_ Reseter = matrix
)

type Dense struct {
	mat blas64.General

	capRows, capCols int
}

func NewDense(r, c int, mat []float64) *Dense {
	if mat != nil && r*c != len(mat) {
		panic(ErrShape)
	}
	if mat == nil {
		mat = make([]float64, r*c)
	}
	return &Dense{
		mat: blas64.General{
			Rows:   r,
			Cols:   c,
			Stride: c,
			Data:   mat,
		},
		capRows: r,
		capCols: c,
	}
}

// DenseCopyOf returns a newly allocated copy of the elements of a.
func DenseCopyOf(a Matrix) *Dense {
	d := &Dense{}
	d.Clone(a)
	return d
}

func (m *Dense) SetRawMatrix(b blas64.General) {
	m.capRows, m.capCols = b.Rows, b.Cols
	m.mat = b
}

func (m *Dense) RawMatrix() blas64.General { return m.mat }

func (m *Dense) isZero() bool {
	// It must be the case that m.Dims() returns
	// zeros in this case. See comment in Reset().
	return m.mat.Stride == 0
}

func (m *Dense) Dims() (r, c int) { return m.mat.Rows, m.mat.Cols }

func (m *Dense) Caps() (r, c int) { return m.capRows, m.capCols }

func (m *Dense) Col(dst []float64, j int) []float64 {
	if j >= m.mat.Cols || j < 0 {
		panic(ErrIndexOutOfRange)
	}

	if dst == nil {
		dst = make([]float64, m.mat.Rows)
	}
	dst = dst[:min(len(dst), m.mat.Rows)]
	blas64.Copy(len(dst),
		blas64.Vector{Inc: m.mat.Stride, Data: m.mat.Data[j:]},
		blas64.Vector{Inc: 1, Data: dst},
	)

	return dst
}

func (m *Dense) SetCol(j int, src []float64) int {
	if j >= m.mat.Cols || j < 0 {
		panic(ErrIndexOutOfRange)
	}

	blas64.Copy(min(len(src), m.mat.Rows),
		blas64.Vector{Inc: 1, Data: src},
		blas64.Vector{Inc: m.mat.Stride, Data: m.mat.Data[j:]},
	)

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

func (m *Dense) View(i, j, r, c int) Matrix {
	mr, mc := m.Dims()
	if i < 0 || i >= mr || j < 0 || j >= mc || r <= 0 || i+r > mr || c <= 0 || j+c > mc {
		panic(ErrIndexOutOfRange)
	}
	t := *m
	t.mat.Data = t.mat.Data[i*t.mat.Stride+j : (i+r-1)*t.mat.Stride+(j+c)]
	t.mat.Rows = r
	t.mat.Cols = c
	t.capRows -= i
	t.capCols -= j
	return &t
}

func (m *Dense) Grow(r, c int) Matrix {
	if r < 0 || c < 0 {
		panic(ErrIndexOutOfRange)
	}
	if r == 0 || c == 0 {
		return m
	}

	r += m.mat.Rows
	c += m.mat.Cols

	var t Dense
	switch {
	case m.mat.Rows == 0 || m.mat.Cols == 0:
		t.mat = blas64.General{
			Rows:   r,
			Cols:   c,
			Stride: c,
			// We zero because we don't know how the matrix will be used.
			// In other places, the mat is immediately filled with a result;
			// this is not the case here.
			Data: useZeroed(m.mat.Data, r*c),
		}
	case r > m.capRows || c > m.capCols:
		cr := max(r, m.capRows)
		cc := max(c, m.capCols)
		t.mat = blas64.General{
			Rows:   r,
			Cols:   c,
			Stride: cc,
			Data:   make([]float64, cr*cc),
		}
		t.capRows = cr
		t.capCols = cc
		// Copy the complete matrix over to the new matrix.
		// Including elements not currently visible.
		r, c, m.mat.Rows, m.mat.Cols = m.mat.Rows, m.mat.Cols, m.capRows, m.capCols
		t.Copy(m)
		m.mat.Rows, m.mat.Cols = r, c
		return &t
	default:
		t.mat = blas64.General{
			Data:   m.mat.Data[:(r-1)*m.mat.Stride+c],
			Rows:   r,
			Cols:   c,
			Stride: m.mat.Stride,
		}
	}
	t.capRows = r
	t.capCols = c
	return &t
}

func (m *Dense) Reset() {
	// No change of Stride, Rows and Cols to 0
	// may be made unless all are set to 0.
	m.mat.Rows, m.mat.Cols, m.mat.Stride = 0, 0, 0
	m.capRows, m.capCols = 0, 0
	m.mat.Data = m.mat.Data[:0]
}

func (m *Dense) Clone(a Matrix) {
	r, c := a.Dims()
	mat := blas64.General{
		Rows:   r,
		Cols:   c,
		Stride: c,
	}
	m.capRows, m.capCols = r, c
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
				m.set(i, j, a.At(i, j))
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
				m.set(r, c, a.At(r, c))
			}
		}
	}

	return r, c
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
		m.mat = blas64.General{
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
			m.set(r, c, a.At(r, c))
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
		m.mat = blas64.General{
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
			m.set(r, c, a.At(r, c))
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
		w.mat = blas64.General{
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
				w.set(i, j, a.at(j, i))
			}
		}
	default:
		for i := 0; i < ac; i++ {
			for j := 0; j < ar; j++ {
				w.set(i, j, a.At(j, i))
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
		m.mat = blas64.General{
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
		m.mat = blas64.General{
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

// MarshalBinary encodes the receiver into a binary form and returns the result.
//
// Dense is little-endian encoded as follows:
//  0 -  8  number of rows    (int64)
//  8 - 16  number of columns (int64)
// 16 - ..  matrix data elements (float64)
//          [0,0] [0,1] ... [0,ncols-1]
//          [1,0] [1,1] ... [1,ncols-1]
//          ...
//          [nrows-1,0] ... [nrows-1,ncols-1]
func (m Dense) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, m.mat.Rows*m.mat.Cols*sizeFloat64+2*sizeInt64))
	err := binary.Write(buf, defaultEndian, int64(m.mat.Rows))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, defaultEndian, int64(m.mat.Cols))
	if err != nil {
		return nil, err
	}

	for i := 0; i < m.mat.Rows; i++ {
		for _, v := range m.rowView(i) {
			err = binary.Write(buf, defaultEndian, v)
			if err != nil {
				return nil, err
			}
		}
	}
	return buf.Bytes(), err
}

// UnmarshalBinary decodes the binary form into the receiver.
// It panics if the receiver is a non-zero Dense matrix.
//
// See MarshalBinary for the on-disk layout.
func (m *Dense) UnmarshalBinary(data []byte) error {
	if !m.isZero() {
		panic("mat64: unmarshal into non-zero matrix")
	}

	buf := bytes.NewReader(data)
	var rows int64
	err := binary.Read(buf, defaultEndian, &rows)
	if err != nil {
		return err
	}
	var cols int64
	err = binary.Read(buf, defaultEndian, &cols)
	if err != nil {
		return err
	}

	m.mat.Rows = int(rows)
	m.mat.Cols = int(cols)
	m.mat.Stride = int(cols)
	m.capRows = int(rows)
	m.capCols = int(cols)
	m.mat.Data = use(m.mat.Data, m.mat.Rows*m.mat.Cols)

	for i := range m.mat.Data {
		err = binary.Read(buf, defaultEndian, &m.mat.Data[i])
		if err != nil {
			return err
		}
	}

	return err
}
