// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"bytes"
	"encoding/binary"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

var (
	matrix *Dense

	_ Matrix  = matrix
	_ Mutable = matrix

	_ VectorSetter = matrix

	_ Cloner       = matrix
	_ Viewer       = matrix
	_ RowViewer    = matrix
	_ ColViewer    = matrix
	_ RawRowViewer = matrix
	_ Grower       = matrix

	_ Adder     = matrix
	_ Suber     = matrix
	_ Muler     = matrix
	_ Dotter    = matrix
	_ ElemMuler = matrix
	_ ElemDiver = matrix
	_ Exper     = matrix

	_ Scaler  = matrix
	_ Applyer = matrix

	_ Normer = matrix
	_ Sumer  = matrix

	_ Stacker   = matrix
	_ Augmenter = matrix

	_ RawMatrixSetter = matrix
	_ RawMatrixer     = matrix

	_ Reseter = matrix
)

// Dense is a dense matrix representation.
type Dense struct {
	mat blas64.General

	capRows, capCols int
}

// NewDense creates a new matrix of type Dense with dimensions r and c.
// If the mat argument is nil, a new data slice is allocated.
//
// The data must be arranged in row-major order, i.e. the (i*c + j)-th
// element in mat is the {i, j}-th element in the matrix.
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

// reuseAs resizes an empty matrix to a r×c matrix,
// or checks that a non-empty matrix is r×c.
func (m *Dense) reuseAs(r, c int) {
	if m.mat.Rows > m.capRows || m.mat.Cols > m.capCols {
		// Panic as a string, not a mat64.Error.
		panic("mat64: caps not correctly set")
	}
	if m.isZero() {
		m.mat = blas64.General{
			Rows:   r,
			Cols:   c,
			Stride: c,
			Data:   use(m.mat.Data, r*c),
		}
		m.capRows = r
		m.capCols = c
		return
	}
	if r != m.mat.Rows || c != m.mat.Cols {
		panic(ErrShape)
	}
}

// untranspose untransposes a matrix if applicable. If a is an Untransposer, then
// untranspose returns the underlying matrix and true. If it is not, then it returns
// the input matrix and false.
func untranspose(a Matrix) (Matrix, bool) {
	if ut, ok := a.(Untransposer); ok {
		return ut.Untranspose(), true
	}
	return a, false
}

// isolatedWorkspace returns a new dense matrix w with the size of a and
// returns a callback to defer which performs cleanup at the return of the call.
// This should be used when a method receiver is the same pointer as an input argument.
func (m *Dense) isolatedWorkspace(a Matrix) (w *Dense, restore func()) {
	r, c := a.Dims()
	w = getWorkspace(r, c, false)
	return w, func() {
		m.Copy(w)
		putWorkspace(w)
	}
}

func (m *Dense) isZero() bool {
	// It must be the case that m.Dims() returns
	// zeros in this case. See comment in Reset().
	return m.mat.Stride == 0
}

// asTriDense returns a TriDense with the given size and side. The backing data
// of the TriDense is the same as the receiver.
func (m *Dense) asTriDense(n int, diag blas.Diag, uplo blas.Uplo) *TriDense {
	return &TriDense{
		blas64.Triangular{
			N:      n,
			Stride: m.mat.Stride,
			Data:   m.mat.Data,
			Uplo:   uplo,
			Diag:   diag,
		},
	}
}

// DenseCopyOf returns a newly allocated copy of the elements of a.
func DenseCopyOf(a Matrix) *Dense {
	d := &Dense{}
	d.Clone(a)
	return d
}

// SetRawMatrix sets the underlying blas64.General used by the receiver.
// Changes to elements in the receiver following the call will be reflected
// in b.
func (m *Dense) SetRawMatrix(b blas64.General) {
	m.capRows, m.capCols = b.Rows, b.Cols
	m.mat = b
}

// RawMatrix returns the underlying blas64.General used by the receiver.
// Changes to elements in the receiver following the call will be reflected
// in returned blas64.General.
func (m *Dense) RawMatrix() blas64.General { return m.mat }

// Dims returns the number of rows and columns in the matrix.
func (m *Dense) Dims() (r, c int) { return m.mat.Rows, m.mat.Cols }

// Caps returns the number of rows and columns in the backing matrix.
func (m *Dense) Caps() (r, c int) { return m.capRows, m.capCols }

// T performs an implicit transpose by returning the receiver inside a Transpose.
func (m *Dense) T() Matrix {
	return Transpose{m}
}

// ColView returns a Vector reflecting col j, backed by the matrix data.
//
// See ColViewer for more information.
func (m *Dense) ColView(j int) *Vector {
	if j >= m.mat.Cols || j < 0 {
		panic(ErrColAccess)
	}
	return &Vector{
		mat: blas64.Vector{
			Inc:  m.mat.Stride,
			Data: m.mat.Data[j : (m.mat.Rows-1)*m.mat.Stride+j+1],
		},
		n: m.mat.Rows,
	}
}

// SetCol sets the values in the specified column of the matrix to the values
// in src. len(src) must equal the number of rows in the receiver.
func (m *Dense) SetCol(j int, src []float64) {
	if j >= m.mat.Cols || j < 0 {
		panic(ErrColAccess)
	}
	if len(src) != m.mat.Rows {
		panic(ErrColLength)
	}

	blas64.Copy(m.mat.Rows,
		blas64.Vector{Inc: 1, Data: src},
		blas64.Vector{Inc: m.mat.Stride, Data: m.mat.Data[j:]},
	)
}

// SetRow sets the values in the specified rows of the matrix to the values
// in src. len(src) must equal the number of columns in the receiver.
func (m *Dense) SetRow(i int, src []float64) {
	if i >= m.mat.Rows || i < 0 {
		panic(ErrRowAccess)
	}
	if len(src) != m.mat.Cols {
		panic(ErrRowLength)
	}

	copy(m.rowView(i), src)
}

// RowView returns row i of the matrix data represented as a column vector,
// backed by the matrix data.
//
// See RowViewer for more information.
func (m *Dense) RowView(i int) *Vector {
	if i >= m.mat.Rows || i < 0 {
		panic(ErrRowAccess)
	}
	return &Vector{
		mat: blas64.Vector{
			Inc:  1,
			Data: m.mat.Data[i*m.mat.Stride : i*m.mat.Stride+m.mat.Cols],
		},
		n: m.mat.Cols,
	}
}

// RawRowView returns a slice backed by the same array as backing the
// receiver.
func (m *Dense) RawRowView(i int) []float64 {
	if i >= m.mat.Rows || i < 0 {
		panic(ErrRowAccess)
	}
	return m.rowView(i)
}

func (m *Dense) rowView(r int) []float64 {
	return m.mat.Data[r*m.mat.Stride : r*m.mat.Stride+m.mat.Cols]
}

// View returns a new Matrix that shares backing data with the receiver.
// The new matrix is located from row i, column j extending r rows and c
// columns.
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

// Grow returns an expanded copy of the receiver. The copy is expanded
// by r rows and c columns. If the dimensions of the new copy are outside
// the caps of the receiver a new allocation is made, otherwise not.
func (m *Dense) Grow(r, c int) Matrix {
	if r < 0 || c < 0 {
		panic(ErrIndexOutOfRange)
	}
	if r == 0 && c == 0 {
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

// Reset zeros the dimensions of the matrix so that it can be reused as the
// receiver of a dimensionally restricted operation.
//
// See the Reseter interface for more information.
func (m *Dense) Reset() {
	// No change of Stride, Rows and Cols to 0
	// may be made unless all are set to 0.
	m.mat.Rows, m.mat.Cols, m.mat.Stride = 0, 0, 0
	m.capRows, m.capCols = 0, 0
	m.mat.Data = m.mat.Data[:0]
}

// Clone makes a copy of a into the receiver, overwriting the previous value of
// the receiver. The clone operation does not make any restriction on shape.
//
// See the Cloner interface for more information.
func (m *Dense) Clone(a Matrix) {
	r, c := a.Dims()
	mat := blas64.General{
		Rows:   r,
		Cols:   c,
		Stride: c,
	}
	m.capRows, m.capCols = r, c

	aU, trans := untranspose(a)
	switch aU := aU.(type) {
	case RawMatrixer:
		amat := aU.RawMatrix()
		// TODO(kortschak): Consider being more precise with determining whether a and m are aliases.
		// The current approach is that all RawMatrixers are considered potential aliases.
		// Note that below we assume that non-RawMatrixers are not aliases; this is not necessarily
		// true, but cases where it is not are not sensible. We should probably fix or document
		// this though.
		mat.Data = make([]float64, r*c)
		if trans {
			for i := 0; i < r; i++ {
				blas64.Copy(c,
					blas64.Vector{Inc: amat.Stride, Data: amat.Data[i : i+(c-1)*amat.Stride+1]},
					blas64.Vector{Inc: 1, Data: mat.Data[i*c : (i+1)*c]})
			}
		} else {
			for i := 0; i < r; i++ {
				copy(mat.Data[i*c:(i+1)*c], amat.Data[i*amat.Stride:i*amat.Stride+c])
			}
		}
	case Vectorer:
		mat.Data = use(m.mat.Data, r*c)
		if trans {
			for i := 0; i < r; i++ {
				aU.Col(mat.Data[i*c:(i+1)*c], i)
			}
		} else {
			for i := 0; i < r; i++ {
				aU.Row(mat.Data[i*c:(i+1)*c], i)
			}
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

// Copy makes a copy of elements of a into the receiver. It is similar to the
// built-in copy; it copies as much as the overlap between the two matrices and
// returns the number of rows and columns it copied.
//
// See the Copier interface for more information.
func (m *Dense) Copy(a Matrix) (r, c int) {
	r, c = a.Dims()
	if a == m {
		return r, c
	}
	r = min(r, m.mat.Rows)
	c = min(c, m.mat.Cols)
	if r == 0 || c == 0 {
		return 0, 0
	}

	aU, trans := untranspose(a)
	switch aU := aU.(type) {
	case RawMatrixer:
		amat := aU.RawMatrix()
		if trans {
			for i := 0; i < r; i++ {
				blas64.Copy(c,
					blas64.Vector{Inc: amat.Stride, Data: amat.Data[i : i+(c-1)*amat.Stride+1]},
					blas64.Vector{Inc: 1, Data: m.mat.Data[i*m.mat.Stride : i*m.mat.Stride+c]})
			}
		} else {
			for i := 0; i < r; i++ {
				copy(m.mat.Data[i*m.mat.Stride:i*m.mat.Stride+c], amat.Data[i*amat.Stride:i*amat.Stride+c])
			}
		}
	case Vectorer:
		if trans {
			for i := 0; i < r; i++ {
				aU.Col(m.mat.Data[i*m.mat.Stride:i*m.mat.Stride+c], i)
			}
		} else {
			for i := 0; i < r; i++ {
				aU.Row(m.mat.Data[i*m.mat.Stride:i*m.mat.Stride+c], i)
			}
		}
	default:
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				m.set(i, j, a.At(i, j))
			}
		}
	}

	return r, c
}

// Stack appends the rows of b onto the rows of a, placing the result into the
// receiver.
//
// See the Stacker interface for more information.
func (m *Dense) Stack(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ac != bc || m == a || m == b {
		panic(ErrShape)
	}

	m.reuseAs(ar+br, ac)

	m.Copy(a)
	w := m.View(ar, 0, br, bc).(*Dense)
	w.Copy(b)
}

// Augment creates the augmented matrix of a and b, where b is placed in the
// greater indexed columns.
//
// See the Augmenter interface for more information.
func (m *Dense) Augment(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || m == a || m == b {
		panic(ErrShape)
	}

	m.reuseAs(ar, ac+bc)

	m.Copy(a)
	w := m.View(0, ac, br, bc).(*Dense)
	w.Copy(b)
}

// MarshalBinary encodes the receiver into a binary form and returns the result.
//
// Dense is little-endian encoded as follows:
//   0 -  8  number of rows    (int64)
//   8 - 16  number of columns (int64)
//  16 - ..  matrix data elements (float64)
//           [0,0] [0,1] ... [0,ncols-1]
//           [1,0] [1,1] ... [1,ncols-1]
//           ...
//           [nrows-1,0] ... [nrows-1,ncols-1]
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
