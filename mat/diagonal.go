// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

var (
	diagDense *DiagDense
	_         Matrix          = diagDense
	_         Diagonal        = diagDense
	_         MutableDiagonal = diagDense
	_         Triangular      = diagDense
	_         Symmetric       = diagDense
	_         Banded          = diagDense
	_         RawBander       = diagDense
	_         RawSymBander    = diagDense
)

// Diagonal represents a diagonal matrix, that is a square matrix that only
// has non-zero terms on the diagonal.
type Diagonal interface {
	Matrix
	// Diag and Symmetric return the number of rows/columns in
	// the matrix. Both methods are included to allow diagonal
	// matrices to be used in functions taking symmetric inputs.
	Diag() int
	Symmetric() int
}

// MutableDiagonal is a Diagonal matrix whose elements can be set.
type MutableDiagonal interface {
	Diagonal
	SetDiag(i int, v float64)
}

// DiagDense represents a diagonal matrix in dense storage format.
type DiagDense struct {
	mat blas64.Vector
	n   int
}

// NewDiagonal creates a new Diagonal matrix with n rows and n columns.
// The length of data must be n or data must be nil, otherwise NewDiagonal
// will panic. NewDiagonal will panic if n is zero.
func NewDiagonal(n int, data []float64) *DiagDense {
	if n <= 0 {
		if n == 0 {
			panic(ErrZeroLength)
		}
		panic("mat: negative dimension")
	}
	if data == nil {
		data = make([]float64, n)
	}
	if len(data) != n {
		panic(ErrShape)
	}
	return &DiagDense{
		mat: blas64.Vector{Data: data, Inc: 1},
		n:   n,
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

// TBand performs an implicit transpose by returning the receiver inside a
// TransposeBand.
func (d *DiagDense) TBand() Banded {
	return TransposeBand{d}
}

// Bandwidth returns the upper and lower bandwidths of the matrix.
// These values are always zero for diagonal matrices.
func (d *DiagDense) Bandwidth() (kl, ku int) {
	return 0, 0
}

// Symmetric implements the Symmetric interface.
func (d *DiagDense) Symmetric() int {
	return d.n
}

// Triangle implements the Triangular interface.
func (d *DiagDense) Triangle() (int, TriKind) {
	return d.n, Upper
}

// RawBand returns the underlying data used by the receiver represented
// as a blas64.Band.
// Changes to elements in the receiver following the call will be reflected
// in returned blas64.Band.
func (d *DiagDense) RawBand() blas64.Band {
	return blas64.Band{
		Rows:   d.n,
		Cols:   d.n,
		KL:     0,
		KU:     0,
		Stride: d.mat.Inc,
		Data:   d.mat.Data,
	}
}

// RawSymBand returns the underlying data used by the receiver represented
// as a blas64.SymmetricBand.
// Changes to elements in the receiver following the call will be reflected
// in returned blas64.Band.
func (d *DiagDense) RawSymBand() blas64.SymmetricBand {
	return blas64.SymmetricBand{
		N:      d.n,
		K:      0,
		Stride: d.mat.Inc,
		Uplo:   blas.Upper,
		Data:   d.mat.Data,
	}
}
