// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

var (
	diagDense *DiagDense
	_         Matrix          = diagDense
	_         Diagonal        = diagDense
	_         MutableDiagonal = diagDense
	_         Triangular      = diagDense
	_         Symmetric       = diagDense
)

// Diagonal represents a diagonal matrix, that is a square matrix that only
// has non-zero terms on the diagonal.
type Diagonal interface {
	Matrix
	Symmetric() int
	TTri() Triangular
	// Triangle implements the Triangular interface. Implementers of Diagonal
	// should return Upper as the default TriKind.
	Triangle() (int, TriKind)
	// Diag returns the number of rows/columns in the matrix
	Diag() int
}

// MutableDiagonal is a Diagonal matrix whose elements can be set.
type MutableDiagonal interface {
	Diagonal
	SetDiag(i int, v float64)
}

// DiagDense represents a diagonal matrix in dense storage format.
type DiagDense struct {
	data []float64
	n    int
	cap  int
}

// NewDiagonal creates a new Diagonal matrix with n rows and n columns.
// The length of data must be n or data must be nil, otherwise NewDiagonal
// will panic.
func NewDiagonal(n int, data []float64) *DiagDense {
	if n < 0 {
		panic("mat: negative dimension")
	}
	if data == nil {
		data = make([]float64, n)
	}
	if len(data) != n {
		panic(ErrShape)
	}
	return &DiagDense{
		data: data,
		n:    n,
		cap:  n,
	}
}

// Diag returns the dimension of the receiver.
func (d *DiagDense) Diag() int {
	return d.n
}

// Dims returns the dimensions of the matrix.
func (d *DiagDense) Dims() (r, c int) {
	return d.n, d.n
}

// T returns the transpose of the matrix.
func (d *DiagDense) T() Matrix {
	return d
}

// TTri returns the transpose of the matrix. Note that Diagonal matrices are
// Upper by default
func (d *DiagDense) TTri() Triangular {
	return TransposeTri{d}
}

// Symmetric implements the Symmetric interface.
func (d *DiagDense) Symmetric() int {
	return d.n
}

// Triangle implements the Triangular interface.
func (d *DiagDense) Triangle() (int, TriKind) {
	return d.n, Upper
}
