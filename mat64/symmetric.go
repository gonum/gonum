// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

var (
	symDense *SymDense

	_ Matrix           = symDense
	_ Symmetric        = symDense
	_ RawSymmetricer   = symDense
	_ MutableSymmetric = symDense
)

const badSymTriangle = "mat64: blas64.Symmetric not upper"

// SymDense is a symmetric matrix that uses Dense storage.
type SymDense struct {
	mat blas64.Symmetric
}

// Symmetric represents a symmetric matrix (where the element at {i, j} equals
// the element at {j, i}). Symmetric matrices are always square.
type Symmetric interface {
	Matrix
	// Symmetric returns the number of rows/columns in the matrix.
	Symmetric() int
}

// A RawSymmetricer can return a view of itself as a BLAS Symmetric matrix.
type RawSymmetricer interface {
	RawSymmetric() blas64.Symmetric
}

type MutableSymmetric interface {
	Symmetric
	SetSym(i, j int, v float64)
}

// NewSymDense constructs an n x n symmetric matrix. If len(mat) == n * n,
// mat will be used to hold the underlying data, or if mat == nil, new data will be allocated.
// The underlying data representation is the same as a Dense matrix, except
// the values of the entries in the lower triangular portion are completely ignored.
func NewSymDense(n int, mat []float64) *SymDense {
	if n < 0 {
		panic("mat64: negative dimension")
	}
	if mat != nil && n*n != len(mat) {
		panic(ErrShape)
	}
	if mat == nil {
		mat = make([]float64, n*n)
	}
	return &SymDense{blas64.Symmetric{
		N:      n,
		Stride: n,
		Data:   mat,
		Uplo:   blas.Upper,
	}}
}

func (s *SymDense) Dims() (r, c int) {
	return s.mat.N, s.mat.N
}

// T implements the Matrix interface. Symmetric matrices, by definition, are
// equal to their transpose, and this is a no-op.
func (s *SymDense) T() Matrix {
	return s
}

func (s *SymDense) Symmetric() int {
	return s.mat.N
}

// RawSymmetric returns the matrix as a blas64.Symmetric. The returned
// value must be stored in upper triangular format.
func (s *SymDense) RawSymmetric() blas64.Symmetric {
	return s.mat
}

func (s *SymDense) isZero() bool {
	return s.mat.N == 0
}

// reuseAs resizes an empty matrix to a n×n matrix,
// or checks that a non-empty matrix is n×n.
func (s *SymDense) reuseAs(n int) {
	if s.isZero() {
		s.mat = blas64.Symmetric{
			N:      n,
			Stride: n,
			Data:   use(s.mat.Data, n*n),
			Uplo:   blas.Upper,
		}
		return
	}
	if s.mat.Uplo != blas.Upper {
		panic(badSymTriangle)
	}
	if s.mat.N != n {
		panic(ErrShape)
	}
}

func (s *SymDense) AddSym(a, b Symmetric) {
	n := a.Symmetric()
	if n != b.Symmetric() {
		panic(ErrShape)
	}
	if s.isZero() {
		s.mat = blas64.Symmetric{
			N:      n,
			Stride: n,
			Data:   use(s.mat.Data, n*n),
			Uplo:   blas.Upper,
		}
	} else if s.mat.N != n {
		panic(ErrShape)
	}

	if a, ok := a.(RawSymmetricer); ok {
		if b, ok := b.(RawSymmetricer); ok {
			amat, bmat := a.RawSymmetric(), b.RawSymmetric()
			for i := 0; i < n; i++ {
				btmp := bmat.Data[i*bmat.Stride+i : i*bmat.Stride+n]
				stmp := s.mat.Data[i*s.mat.Stride+i : i*s.mat.Stride+n]
				for j, v := range amat.Data[i*amat.Stride+i : i*amat.Stride+n] {
					stmp[j] = v + btmp[j]
				}
			}
			return
		}
	}

	for i := 0; i < n; i++ {
		stmp := s.mat.Data[i*s.mat.Stride : i*s.mat.Stride+n]
		for j := i; j < n; j++ {
			stmp[j] = a.At(i, j) + b.At(i, j)
		}
	}
}

func (s *SymDense) CopySym(a Symmetric) int {
	n := a.Symmetric()
	n = min(n, s.mat.N)
	if n == 0 {
		return 0
	}
	switch a := a.(type) {
	case RawSymmetricer:
		amat := a.RawSymmetric()
		if amat.Uplo != blas.Upper {
			panic(badSymTriangle)
		}
		for i := 0; i < n; i++ {
			copy(s.mat.Data[i*s.mat.Stride+i:i*s.mat.Stride+n], amat.Data[i*amat.Stride+i:i*amat.Stride+n])
		}
	default:
		for i := 0; i < n; i++ {
			stmp := s.mat.Data[i*s.mat.Stride : i*s.mat.Stride+n]
			for j := i; j < n; j++ {
				stmp[j] = a.At(i, j)
			}
		}
	}
	return n
}

// SymRankOne performs a symetric rank-one update to the matrix a and stores
// the result in the receiver
//  s = a + alpha * x * x'
func (s *SymDense) SymRankOne(a Symmetric, alpha float64, x *Vector) {
	n := s.mat.N
	if x.Len() != n {
		panic(ErrShape)
	}
	var w SymDense
	if s == a {
		w = *s
	}
	if w.isZero() {
		w.mat = blas64.Symmetric{
			N:      n,
			Stride: n,
			Uplo:   blas.Upper,
			Data:   use(w.mat.Data, n*n),
		}
	} else if n != w.mat.N {
		panic(ErrShape)
	}
	if s != a {
		w.CopySym(a)
	}
	blas64.Syr(alpha, x.mat, w.mat)
	*s = w
	return
}

// SymRankK performs a symmetric rank-k update to the matrix a and stores the
// result into the receiver. If a is zero, see SymOuterK.
//  s = a + alpha * x * x'
func (s *SymDense) SymRankK(a Symmetric, alpha float64, x Matrix) {
	n := a.Symmetric()
	r, _ := x.Dims()
	if r != n {
		panic(ErrShape)
	}
	xMat, aTrans := untranspose(x)
	var g blas64.General
	if rm, ok := xMat.(RawMatrixer); ok {
		g = rm.RawMatrix()
	} else {
		g = DenseCopyOf(x).mat
		aTrans = false
	}
	if a != s {
		s.reuseAs(n)
		s.CopySym(a)
	}
	t := blas.NoTrans
	if aTrans {
		t = blas.Trans
	}
	blas64.Syrk(t, alpha, g, 1, s.mat)
}

// SymOuterK calculates the outer product of a times its transpose and stores
// the result into the receiver. In order to update an existing matrix, see
// SymRankOne
//  s = x * x'
func (s *SymDense) SymOuterK(x Matrix) {
	r, _ := x.Dims()
	s.reuseAs(r)
	s.SymRankK(s, 1, x)
}

// RankTwo performs a symmmetric rank-two update to the matrix a and stores
// the result in the receiver
//  m = a + alpha * (x * y' + y * x')
func (s *SymDense) RankTwo(a Symmetric, alpha float64, x, y *Vector) {
	n := s.mat.N
	if x.Len() != n {
		panic(ErrShape)
	}
	if y.Len() != n {
		panic(ErrShape)
	}
	var w SymDense
	if s == a {
		w = *s
	}
	if w.isZero() {
		w.mat = blas64.Symmetric{
			N:      n,
			Stride: n,
			Uplo:   blas.Upper,
			Data:   use(w.mat.Data, n*n),
		}
	} else if n != w.mat.N {
		panic(ErrShape)
	}
	if s != a {
		w.CopySym(a)
	}
	blas64.Syr2(alpha, x.mat, y.mat, w.mat)
	*s = w
	return
}

// ScaleSym multiplies the elements of a by f, placing the result in the receiver.
func (s *SymDense) ScaleSym(f float64, a Symmetric) {
	n := a.Symmetric()
	s.reuseAs(n)
	if a, ok := a.(RawSymmetricer); ok {
		amat := a.RawSymmetric()
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				s.mat.Data[i*s.mat.Stride+j] = f * amat.Data[i*amat.Stride+j]
			}
		}
		return
	}
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			s.mat.Data[i*s.mat.Stride+j] = f * a.At(i, j)
		}
	}
}
