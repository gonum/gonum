// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blas64

import "gonum.org/v1/gonum/blas"

// GeneralCols represents a matrix using the conventional column-major storage scheme.
type GeneralCols General

// From fills the receiver with elements from a. The receiver
// must have the same dimensions as a and have adequate backing
// data storage.
func (t GeneralCols) From(a General) {
	if t.Rows != a.Rows || t.Cols != a.Cols {
		panic("blas64: mismatched dimension")
	}
	if len(t.Data) < (t.Cols-1)*t.Stride+t.Rows {
		panic("blas64: short data slice")
	}
	for i := 0; i < a.Rows; i++ {
		for j, v := range a.Data[i*a.Stride : i*a.Stride+a.Cols] {
			t.Data[i+j*t.Stride] = v
		}
	}
}

// From fills the receiver with elements from a. The receiver
// must have the same dimensions as a and have adequate backing
// data storage.
func (t General) From(a GeneralCols) {
	if t.Rows != a.Rows || t.Cols != a.Cols {
		panic("blas64: mismatched dimension")
	}
	if len(t.Data) < (t.Rows-1)*t.Stride+t.Cols {
		panic("blas64: short data slice")
	}
	for i := 0; i < a.Cols; i++ {
		for j, v := range a.Data[i*a.Stride : i*a.Stride+a.Rows] {
			t.Data[i+j*t.Stride] = v
		}
	}
}

// SymmetricCols represents a matrix using the conventional column-major storage scheme.
type SymmetricCols Symmetric

// From fills the receiver with elements from a. The receiver
// must have the same dimensions and uplo as a and have adequate
// backing data storage.
func (t SymmetricCols) From(a Symmetric) {
	if t.N != a.N {
		panic("blas64: mismatched dimension")
	}
	if t.Uplo != a.Uplo {
		panic("blas64: mismatched BLAS uplo")
	}
	switch a.Uplo {
	default:
		panic("blas64: bad BLAS uplo")
	case blas.Upper:
		for i := 0; i < a.N; i++ {
			for j := i; j < a.N; j++ {
				t.Data[i+j*t.Stride] = a.Data[i*a.Stride+j]
			}
		}
	case blas.Lower:
		for i := 0; i < a.N; i++ {
			for j := 0; j <= i; j++ {
				t.Data[i+j*t.Stride] = a.Data[i*a.Stride+j]
			}
		}
	}
}

// From fills the receiver with elements from a. The receiver
// must have the same dimensions and uplo as a and have adequate
// backing data storage.
func (t Symmetric) From(a SymmetricCols) {
	if t.N != a.N {
		panic("blas64: mismatched dimension")
	}
	if t.Uplo != a.Uplo {
		panic("blas64: mismatched BLAS uplo")
	}
	switch a.Uplo {
	default:
		panic("blas64: bad BLAS uplo")
	case blas.Upper:
		for i := 0; i < a.N; i++ {
			for j := i; j < a.N; j++ {
				t.Data[i*t.Stride+j] = a.Data[i+j*a.Stride]
			}
		}
	case blas.Lower:
		for i := 0; i < a.N; i++ {
			for j := 0; j <= i; j++ {
				t.Data[i*t.Stride+j] = a.Data[i+j*a.Stride]
			}
		}
	}
}

// TriangularCols represents a matrix using the conventional column-major storage scheme.
type TriangularCols Triangular

// From fills the receiver with elements from a. The receiver
// must have the same dimensions as a and have adequate backing
// data storage.
func (t TriangularCols) From(a Triangular) {
	if t.N != a.N {
		panic("blas64: mismatched dimension")
	}
	t.Uplo = a.Uplo
	t.Diag = a.Diag
	switch a.Uplo {
	default:
		panic("blas64: bad BLAS uplo")
	case blas.Upper:
		for i := 0; i < a.N; i++ {
			for j := i; j < a.N; j++ {
				t.Data[i+j*t.Stride] = a.Data[i*a.Stride+j]
			}
		}
	case blas.Lower:
		for i := 0; i < a.N; i++ {
			for j := 0; j <= i; j++ {
				t.Data[i+j*t.Stride] = a.Data[i*a.Stride+j]
			}
		}
	case blas.All:
		for i := 0; i < a.N; i++ {
			for j := 0; j < a.N; j++ {
				t.Data[i+j*t.Stride] = a.Data[i*a.Stride+j]
			}
		}
	}
}

// From fills the receiver with elements from a. The receiver
// must have the same dimensions as a and have adequate backing
// data storage.
func (t Triangular) From(a TriangularCols) {
	if t.N != a.N {
		panic("blas64: mismatched dimension")
	}
	t.Uplo = a.Uplo
	t.Diag = a.Diag
	switch a.Uplo {
	default:
		panic("blas64: bad BLAS uplo")
	case blas.Upper:
		for i := 0; i < a.N; i++ {
			for j := i; j < a.N; j++ {
				t.Data[i*t.Stride+j] = a.Data[i+j*a.Stride]
			}
		}
	case blas.Lower:
		for i := 0; i < a.N; i++ {
			for j := 0; j <= i; j++ {
				t.Data[i*t.Stride+j] = a.Data[i+j*a.Stride]
			}
		}
	case blas.All:
		for i := 0; i < a.N; i++ {
			for j := 0; j < a.N; j++ {
				t.Data[i*t.Stride+j] = a.Data[i+j*a.Stride]
			}
		}
	}
}

// BandCols represents a matrix using the band column-major storage scheme.
type BandCols Band

// From fills the receiver with elements from a. The receiver
// must have the same dimensions and bandwidth as a and have
// adequate backing data storage.
func (t BandCols) From(a Band) {
	if t.Rows != a.Rows || t.Cols != a.Cols {
		panic("blas64: mismatched dimension")
	}
	if t.KL != a.KL || t.KU != a.KU {
		panic("blas64: mismatched bandwidth")
	}
	if a.Stride < a.KL+a.KU+1 {
		panic("blas64: short stride for source")
	}
	if t.Stride < t.KL+t.KU+1 {
		panic("blas64: short stride for destination")
	}
	for i := 0; i < a.Rows; i++ {
		for j := max(0, i-a.KL); j < min(i+a.KU+1, a.Cols); j++ {
			t.Data[i+t.KU-j+j*t.Stride] = a.Data[j+a.KL-i+i*a.Stride]
		}
	}
}

// From fills the receiver with elements from a. The receiver
// must have the same dimensions and bandwidth as a and have
// adequate backing data storage.
func (t Band) From(a BandCols) {
	if t.Rows != a.Rows || t.Cols != a.Cols {
		panic("blas64: mismatched dimension")
	}
	if t.KL != a.KL || t.KU != a.KU {
		panic("blas64: mismatched bandwidth")
	}
	if a.Stride < a.KL+a.KU+1 {
		panic("blas64: short stride for source")
	}
	if t.Stride < t.KL+t.KU+1 {
		panic("blas64: short stride for destination")
	}
	for j := 0; j < a.Cols; j++ {
		for i := max(0, j-a.KU); i < min(j+a.KL+1, a.Rows); i++ {
			t.Data[j+a.KL-i+i*a.Stride] = a.Data[i+t.KU-j+j*t.Stride]
		}
	}
}

// SymmetricBandCols represents a symmetric matrix using the band column-major storage scheme.
type SymmetricBandCols SymmetricBand

// From fills the receiver with elements from a. The receiver
// must have the same dimensions, bandwidth and uplo as a and
// have adequate backing data storage.
func (t SymmetricBandCols) From(a SymmetricBand) {
	if t.N != a.N {
		panic("blas64: mismatched dimension")
	}
	if t.K != a.K {
		panic("blas64: mismatched bandwidth")
	}
	if a.Stride < a.K+1 {
		panic("blas64: short stride for source")
	}
	if t.Stride < t.K+1 {
		panic("blas64: short stride for destination")
	}
	if t.Uplo != a.Uplo {
		panic("blas64: mismatched BLAS uplo")
	}
	switch a.Uplo {
	default:
		panic("blas64: bad BLAS uplo")
	case blas.Upper:
		for i := 0; i < a.N; i++ {
			for j := i; j < min(i+a.K+1, a.N); j++ {
				t.Data[i+t.K-j+j*t.Stride] = a.Data[j+a.K-i+i*a.Stride]
			}
		}
	case blas.Lower:
		for i := 0; i < a.N; i++ {
			for j := max(0, i-a.K); j <= i; j++ {
				t.Data[i+t.K-j+j*t.Stride] = a.Data[j+a.K-i+i*a.Stride]
			}
		}
	}
}

// From fills the receiver with elements from a. The receiver
// must have the same dimensions, bandwidth and uplo as a and
// have adequate backing data storage.
func (t SymmetricBand) From(a SymmetricBandCols) {
	if t.N != a.N {
		panic("blas64: mismatched dimension")
	}
	if t.K != a.K {
		panic("blas64: mismatched bandwidth")
	}
	if a.Stride < a.K+1 {
		panic("blas64: short stride for source")
	}
	if t.Stride < t.K+1 {
		panic("blas64: short stride for destination")
	}
	if t.Uplo != a.Uplo {
		panic("blas64: mismatched BLAS uplo")
	}
	switch a.Uplo {
	default:
		panic("blas64: bad BLAS uplo")
	case blas.Upper:
		for j := 0; j < a.N; j++ {
			for i := j; i < min(j+a.K+1, a.N); i++ {
				t.Data[j+a.K-i+i*a.Stride] = a.Data[i+t.K-j+j*t.Stride]
			}
		}
	case blas.Lower:
		for j := 0; j < a.N; j++ {
			for i := max(0, j-a.K); i <= i; i++ {
				t.Data[j+a.K-i+i*a.Stride] = a.Data[i+t.K-j+j*t.Stride]
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
