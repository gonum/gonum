package mat64

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

// Symmetric represents a symmetric matrix.
type Symmetric struct {
	mat blas64.Symmetric
}

// NewSymmetric constructs an n x n symmetric matrix. If len(mat) == n * n,
// mat will be used to hold the underlying data, or if mat == nil, new data will be allocated.
// The underlying data representation is the same as a Dense matrix, except
// the values of the entries in the lower triangular portion are completely ignored.
func NewSymmetric(n int, mat []float64) *Symmetric {
	if n < 0 {
		panic("mat64: negative dimension")
	}
	if mat != nil && n*n != len(mat) {
		panic(ErrShape)
	}
	if mat == nil {
		mat = make([]float64, n*n)
	}
	return &Symmetric{blas64.Symmetric{
		N:      n,
		Stride: n,
		Data:   mat,
		Uplo:   blas.Upper,
	}}
}

func (s *Symmetric) Dims() (r, c int) {
	return s.mat.N, s.mat.N
}

func (s *Symmetric) RawSymmetric() blas64.Symmetric {
	return s.mat
}
