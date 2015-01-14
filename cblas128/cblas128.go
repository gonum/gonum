// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cblas128 provides a simple interface to the complex128 BLAS API.
package cblas128

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/cgo"
)

// TODO(kortschak): Change this and the comment below to native.Implementation
// when blas/native covers the complex BLAS API.
var cblas128 blas.Complex128 = cgo.Implementation{}

// Use sets the BLAS complex128 implementation to be used by subsequent BLAS calls.
// The default implementation is cgo.Implementation.
func Use(b blas.Complex128) {
	cblas128 = b
}

// Implementation returns the current BLAS complex128 implementation.
//
// Implementation allows direct calls to the current the BLAS complex128 implementation
// giving finer control of parameters.
func Implementation() blas.Complex128 {
	return cblas128
}

// Vector represents a vector with an associated element increment.
type Vector struct {
	Inc  int
	Data []complex128
}

// General represents a matrix using the conventional storage scheme.
type General struct {
	Rows, Cols int
	Stride     int
	Data       []complex128
}

// Band represents a band matrix using the band storage scheme.
type Band struct {
	Rows, Cols int
	KL, KU     int
	Stride     int
	Data       []complex128
}

// Triangular represents a triangular matrix using the conventional storage scheme.
type Triangular struct {
	N      int
	Stride int
	Data   []complex128
	Uplo   blas.Uplo
	Diag   blas.Diag
}

// TriangularBand represents a triangular matrix using the band storage scheme.
type TriangularBand struct {
	N, K   int
	Stride int
	Data   []complex128
	Uplo   blas.Uplo
	Diag   blas.Diag
}

// TriangularPacked represents a triangular matrix using the packed storage scheme.
type TriangularPacked struct {
	N    int
	Data []complex128
	Uplo blas.Uplo
	Diag blas.Diag
}

// Symmetric represents a symmetric matrix using the conventional storage scheme.
type Symmetric struct {
	N      int
	Stride int
	Data   []complex128
	Uplo   blas.Uplo
}

// SymmetricBand represents a symmetric matrix using the band storage scheme.
type SymmetricBand struct {
	N, K   int
	Stride int
	Data   []complex128
	Uplo   blas.Uplo
}

// SymmetricPacked represents a symmetric matrix using the packed storage scheme.
type SymmetricPacked struct {
	N    int
	Data []complex128
	Uplo blas.Uplo
}

// Hermitian represents an Hermitian matrix using the conventional storage scheme.
type Hermitian Symmetric

// HermitianBand represents an Hermitian matrix using the band storage scheme.
type HermitianBand SymmetricBand

// HermitianPacked represents an Hermitian matrix using the packed storage scheme.
type HermitianPacked SymmetricPacked

// Level 1

const negInc = "cblas128: negative vector increment"

func Dotu(n int, x, y Vector) complex128 {
	return cblas128.Zdotu(n, x.Data, x.Inc, y.Data, y.Inc)
}

func Dotc(n int, x, y Vector) complex128 {
	return cblas128.Zdotc(n, x.Data, x.Inc, y.Data, y.Inc)
}

// Nrm2 will panic if the vector increment is negative.
func Nrm2(n int, x Vector) float64 {
	if x.Inc < 0 {
		panic(negInc)
	}
	return cblas128.Dznrm2(n, x.Data, x.Inc)
}

// Asum will panic if the vector increment is negative.
func Asum(n int, x Vector) float64 {
	if x.Inc < 0 {
		panic(negInc)
	}
	return cblas128.Dzasum(n, x.Data, x.Inc)
}

// Iamax will panic if the vector increment is negative.
func Iamax(n int, x Vector) int {
	if x.Inc < 0 {
		panic(negInc)
	}
	return cblas128.Izamax(n, x.Data, x.Inc)
}

func Swap(n int, x, y Vector) {
	cblas128.Zswap(n, x.Data, x.Inc, y.Data, y.Inc)
}

func Copy(n int, x, y Vector) {
	cblas128.Zcopy(n, x.Data, x.Inc, y.Data, y.Inc)
}

func Axpy(n int, alpha complex128, x, y Vector) {
	cblas128.Zaxpy(n, alpha, x.Data, x.Inc, y.Data, y.Inc)
}

// Scal will panic if the vector increment is negative
func Scal(n int, alpha complex128, x Vector) {
	if x.Inc < 0 {
		panic(negInc)
	}
	cblas128.Zscal(n, alpha, x.Data, x.Inc)
}

// Dscal will panic if the vector increment is negative
func Dscal(n int, alpha float64, x Vector) {
	if x.Inc < 0 {
		panic(negInc)
	}
	cblas128.Zdscal(n, alpha, x.Data, x.Inc)
}

// Level 2

func Gemv(tA blas.Transpose, alpha complex128, a General, x Vector, beta complex128, y Vector) {
	cblas128.Zgemv(tA, a.Rows, a.Cols, alpha, a.Data, a.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Gbmv(tA blas.Transpose, alpha complex128, a Band, x Vector, beta complex128, y Vector) {
	cblas128.Zgbmv(tA, a.Rows, a.Cols, a.KL, a.KU, alpha, a.Data, a.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Trmv(tA blas.Transpose, a Triangular, x Vector) {
	cblas128.Ztrmv(a.Uplo, tA, a.Diag, a.N, a.Data, a.Stride, x.Data, x.Inc)
}

func Tbmv(tA blas.Transpose, a TriangularBand, x Vector) {
	cblas128.Ztbmv(a.Uplo, tA, a.Diag, a.N, a.K, a.Data, a.Stride, x.Data, x.Inc)
}

func Tpmv(tA blas.Transpose, a TriangularPacked, x Vector) {
	cblas128.Ztpmv(a.Uplo, tA, a.Diag, a.N, a.Data, x.Data, x.Inc)
}

func Trsv(tA blas.Transpose, a Triangular, x Vector) {
	cblas128.Ztrsv(a.Uplo, tA, a.Diag, a.N, a.Data, a.Stride, x.Data, x.Inc)
}

func Tbsv(tA blas.Transpose, a TriangularBand, x Vector) {
	cblas128.Ztbsv(a.Uplo, tA, a.Diag, a.N, a.K, a.Data, a.Stride, x.Data, x.Inc)
}

func Tpsv(tA blas.Transpose, a TriangularPacked, x Vector) {
	cblas128.Ztpsv(a.Uplo, tA, a.Diag, a.N, a.Data, x.Data, x.Inc)
}

func Hemv(alpha complex128, a Hermitian, x Vector, beta complex128, y Vector) {
	cblas128.Zhemv(a.Uplo, a.N, alpha, a.Data, a.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Hbmv(alpha complex128, a HermitianBand, x Vector, beta complex128, y Vector) {
	cblas128.Zhbmv(a.Uplo, a.N, a.K, alpha, a.Data, a.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Hpmv(alpha complex128, a HermitianPacked, x Vector, beta complex128, y Vector) {
	cblas128.Zhpmv(a.Uplo, a.N, alpha, a.Data, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Geru(alpha complex128, x, y Vector, a General) {
	cblas128.Zgeru(a.Rows, a.Cols, alpha, x.Data, x.Inc, y.Data, y.Inc, a.Data, a.Stride)
}

func Gerc(alpha complex128, x, y Vector, a General) {
	cblas128.Zgerc(a.Rows, a.Cols, alpha, x.Data, x.Inc, y.Data, y.Inc, a.Data, a.Stride)
}

func Her(alpha float64, x Vector, a Hermitian) {
	cblas128.Zher(a.Uplo, a.N, alpha, x.Data, x.Inc, a.Data, a.Stride)
}

func Hpr(alpha float64, x Vector, a HermitianPacked) {
	cblas128.Zhpr(a.Uplo, a.N, alpha, x.Data, x.Inc, a.Data)
}

func Her2(alpha complex128, x, y Vector, a Hermitian) {
	cblas128.Zher2(a.Uplo, a.N, alpha, x.Data, x.Inc, y.Data, y.Inc, a.Data, a.Stride)
}

func Hpr2(alpha complex128, x, y Vector, a HermitianPacked) {
	cblas128.Zhpr2(a.Uplo, a.N, alpha, x.Data, x.Inc, y.Data, y.Inc, a.Data)
}

// Level 3

func Gemm(tA, tB blas.Transpose, alpha complex128, a, b General, beta complex128, c General) {
	var m, n, k int
	if tA == blas.NoTrans {
		m, k = a.Rows, a.Cols
	} else {
		m, k = a.Cols, a.Rows
	}
	if tB == blas.NoTrans {
		n = b.Cols
	} else {
		n = b.Rows
	}
	cblas128.Zgemm(tA, tB, m, n, k, alpha, a.Data, a.Stride, b.Data, b.Stride, beta, c.Data, c.Stride)
}

func Symm(s blas.Side, alpha complex128, a Symmetric, b General, beta complex128, c General) {
	var m, n int
	if s == blas.Left {
		m, n = a.N, b.Cols
	} else {
		m, n = b.Rows, a.N
	}
	cblas128.Zsymm(s, a.Uplo, m, n, alpha, a.Data, a.Stride, b.Data, b.Stride, beta, c.Data, c.Stride)
}

func Syrk(t blas.Transpose, alpha complex128, a General, beta complex128, c Symmetric) {
	var n, k int
	if t == blas.NoTrans {
		n, k = a.Rows, a.Cols
	} else {
		n, k = a.Cols, a.Rows
	}
	cblas128.Zsyrk(c.Uplo, t, n, k, alpha, a.Data, a.Stride, beta, c.Data, c.Stride)
}

func Syr2k(t blas.Transpose, alpha complex128, a, b General, beta complex128, c Symmetric) {
	var n, k int
	if t == blas.NoTrans {
		n, k = a.Rows, a.Cols
	} else {
		n, k = a.Cols, a.Rows
	}
	cblas128.Zsyr2k(c.Uplo, t, n, k, alpha, a.Data, a.Stride, b.Data, b.Stride, beta, c.Data, c.Stride)
}

func Trmm(s blas.Side, tA blas.Transpose, alpha complex128, a Triangular, b General) {
	cblas128.Ztrmm(s, a.Uplo, tA, a.Diag, b.Rows, b.Cols, alpha, a.Data, a.Stride, b.Data, b.Stride)
}

func Trsm(s blas.Side, tA blas.Transpose, alpha complex128, a Triangular, b General) {
	cblas128.Ztrsm(s, a.Uplo, tA, a.Diag, b.Rows, b.Cols, alpha, a.Data, a.Stride, b.Data, b.Stride)
}

func Hemm(s blas.Side, alpha complex128, a Hermitian, b General, beta complex128, c General) {
	var m, n int
	if s == blas.Left {
		m, n = a.N, b.Cols
	} else {
		m, n = b.Rows, a.N
	}
	cblas128.Zhemm(s, a.Uplo, m, n, alpha, a.Data, a.Stride, b.Data, b.Stride, beta, c.Data, c.Stride)
}

func Herk(t blas.Transpose, alpha float64, a General, beta float64, c Hermitian) {
	var n, k int
	if t == blas.NoTrans {
		n, k = a.Rows, a.Cols
	} else {
		n, k = a.Cols, a.Rows
	}
	cblas128.Zherk(c.Uplo, t, n, k, alpha, a.Data, a.Stride, beta, c.Data, c.Stride)
}

func Her2k(t blas.Transpose, alpha complex128, a, b General, beta float64, c Hermitian) {
	var n, k int
	if t == blas.NoTrans {
		n, k = a.Rows, a.Cols
	} else {
		n, k = a.Cols, a.Rows
	}
	cblas128.Zher2k(c.Uplo, t, n, k, alpha, a.Data, a.Stride, b.Data, b.Stride, beta, c.Data, c.Stride)
}
