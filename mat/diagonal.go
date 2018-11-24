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
	_         SymBanded       = diagDense
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

	// Bandwidth and TBand are included in the Diagonal interface
	// to allow the use of Diagonal types in banded functions.
	// Bandwidth will always return (0, 0).
	Bandwidth() (kl, ku int)
	TBand() Banded
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

// Reset zeros the length of the matrix so that it can be reused as the
// receiver of a dimensionally restricted operation.
//
// See the Reseter interface for more information.
func (d *DiagDense) Reset() {
	// No change of Inc or n to 0 may be
	// made unless both are set to 0.
	d.mat.Inc = 0
	d.n = 0
	d.mat.Data = d.mat.Data[:0]
}

// DiagFrom copies the diagonal of m into the receiver. The receiver must
// be min(r, c) long or zero. Otherwise DiagOf will panic.
func (d *DiagDense) DiagFrom(m Matrix) {
	n := min(m.Dims())
	d.reuseAs(n)

	var vec blas64.Vector
	switch r := m.(type) {
	case *DiagDense:
		vec = r.mat
	case RawBander:
		mat := r.RawBand()
		vec = blas64.Vector{
			Inc:  mat.Stride,
			Data: mat.Data[mat.KL : (n-1)*mat.Stride+mat.KL+1],
		}
	case RawMatrixer:
		mat := r.RawMatrix()
		vec = blas64.Vector{
			Inc:  mat.Stride + 1,
			Data: mat.Data[:(n-1)*mat.Stride+n],
		}
	case RawSymBander:
		mat := r.RawSymBand()
		vec = blas64.Vector{
			Inc:  mat.Stride,
			Data: mat.Data[:(n-1)*mat.Stride+1],
		}
	case RawSymmetricer:
		mat := r.RawSymmetric()
		vec = blas64.Vector{
			Inc:  mat.Stride + 1,
			Data: mat.Data[:(n-1)*mat.Stride+n],
		}
	// TODO(kortschak): Add banded triangular handling when the type exists.
	case RawTriangular:
		mat := r.RawTriangular()
		if mat.Diag == blas.Unit {
			for i := 0; i < n; i += d.mat.Inc {
				d.mat.Data[i] = 1
			}
			return
		}
		vec = blas64.Vector{
			Inc:  mat.Stride + 1,
			Data: mat.Data[:(n-1)*mat.Stride+n],
		}
	case RawVectorer:
		d.mat.Data[0] = r.RawVector().Data[0]
		return
	default:
		for i := 0; i < n; i++ {
			d.setDiag(i, m.At(i, i))
		}
		return
	}
	blas64.Copy(n, vec, d.mat)
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

// reuseAs resizes an empty diagonal to a r×r diagonal,
// or checks that a non-empty matrix is r×r.
func (d *DiagDense) reuseAs(r int) {
	if r == 0 {
		panic(ErrZeroLength)
	}
	if d.IsZero() {
		d.mat = blas64.Vector{
			Inc:  1,
			Data: use(d.mat.Data, r),
		}
		d.n = r
		return
	}
	if r != d.n {
		panic(ErrShape)
	}
}

// IsZero returns whether the receiver is zero-sized. Zero-sized vectors can be the
// receiver for size-restricted operations. DiagDenses can be zeroed using Reset.
func (d *DiagDense) IsZero() bool {
	// It must be the case that d.Dims() returns
	// zeros in this case. See comment in Reset().
	return d.mat.Inc == 0
}
