// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"gonum.org/v1/gonum/blas/blas64"
)

var (
	bandDense *BandDense
	_         Matrix  = bandDense
	_         Banded  = bandDense
	_         RawBand = bandDense
)

// BandDense represents a band matrix in dense storage format.
type BandDense struct {
	mat blas64.Band
}

type Banded interface {
	Matrix
	// Bandwidth returns the lower and upper bandwidth values for
	// the matrix. The total bandwidth of the matrix is kl+ku+1.
	Bandwidth() (kl, ku int)

	// TBand is the equivalent of the T() method in the Matrix
	// interface but guarantees the transpose is of banded type.
	TBand() Banded
}

type RawBand interface {
	RawBand() blas64.Band
}

var (
	_ Matrix            = TransposeBand{}
	_ Banded            = TransposeBand{}
	_ UntransposeBander = TransposeBand{}
)

// TransposeBand is a type for performing an implicit transpose of a band
// matrix. It implements the Banded interface, returning values from the
// transpose of the matrix within.
type TransposeBand struct {
	Banded Banded
}

// At returns the value of the element at row i and column j of the transposed
// matrix, that is, row j and column i of the Banded field.
func (t TransposeBand) At(i, j int) float64 {
	return t.Banded.At(j, i)
}

// Dims returns the dimensions of the transposed matrix.
func (t TransposeBand) Dims() (r, c int) {
	c, r = t.Banded.Dims()
	return r, c
}

// T performs an implicit transpose by returning the Banded field.
func (t TransposeBand) T() Matrix {
	return t.Banded
}

// Bandwidth returns the number of rows/columns in the matrix and its orientation.
func (t TransposeBand) Bandwidth() (kl, ku int) {
	kl, ku = t.Banded.Bandwidth()
	return ku, kl
}

// TBand performs an implicit transpose by returning the Banded field.
func (t TransposeBand) TBand() Banded {
	return t.Banded
}

// Untranspose returns the Banded field.
func (t TransposeBand) Untranspose() Matrix {
	return t.Banded
}

// UntransposeBand returns the Banded field.
func (t TransposeBand) UntransposeBand() Banded {
	return t.Banded
}

// NewBandDense creates a new Band matrix with r rows and c columns. If data == nil,
// a new slice is allocated for the backing slice. If len(data) == min(r, c+kl)*(kl+ku+1),
// data is used as the backing slice, and changes to the elements of the returned
// BandDense will be reflected in data. If neither of these is true, NewBandDense
// will panic. kl must be at least zero and less r, and ku must be at least zero and
// less than c, otherwise NewBandDense will panic.
//
// The data must be arranged in row-major order constructed by removing the zeros
// from the rows outside the band and aligning the diagonals. For example, the matrix
//    1  2  3  0  0  0
//    4  5  6  7  0  0
//    0  8  9 10 11  0
//    0  0 12 13 14 15
//    0  0  0 16 17 18
//    0  0  0  0 19 20
// becomes (* entries are never accessed)
//     *  1  2  3
//     4  5  6  7
//     8  9 10 11
//    12 13 14 15
//    16 17 18  *
//    19 20  *  *
// which is passed to NewBandDense as []float64{*, 1, 2, 3, 4, ...} with kl=1 and ku=2.
// Only the values in the band portion of the matrix are used.
func NewBandDense(r, c, kl, ku int, data []float64) *BandDense {
	if r < 0 || c < 0 || kl < 0 || ku < 0 {
		panic("mat: negative dimension")
	}
	if kl+1 > r || ku+1 > c {
		panic("mat: band out of range")
	}
	bc := kl + ku + 1
	if data != nil && len(data) != min(r, c+kl)*bc {
		panic(ErrShape)
	}
	if data == nil {
		data = make([]float64, min(r, c+kl)*bc)
	}
	return &BandDense{
		mat: blas64.Band{
			Rows:   r,
			Cols:   c,
			KL:     kl,
			KU:     ku,
			Stride: bc,
			Data:   data,
		},
	}
}

// Dims returns the number of rows and columns in the matrix.
func (b *BandDense) Dims() (r, c int) {
	return b.mat.Rows, b.mat.Cols
}

// Bandwidth returns the upper and lower bandwidths of the matrix.
func (b *BandDense) Bandwidth() (kl, ku int) {
	return b.mat.KL, b.mat.KU
}

// T performs an implicit transpose by returning the receiver inside a Transpose.
func (b *BandDense) T() Matrix {
	return Transpose{b}
}

// TBand performs an implicit transpose by returning the receiver inside a TransposeBand.
func (b *BandDense) TBand() Banded {
	return TransposeBand{b}
}

// RawBand returns the underlying blas64.Band used by the receiver.
// Changes to elements in the receiver following the call will be reflected
// in returned blas64.Band.
func (b *BandDense) RawBand() blas64.Band {
	return b.mat
}
