// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package blas64 provides a simple interface to the float64 BLAS API.
package blas64

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/native"
)

var blas64 blas.Float64 = native.Implementation{}

// Use sets the BLAS float64 implementation to be used by subsequent BLAS calls.
// The default implementation is native.Implementation.
func Use(b blas.Float64) {
	blas64 = b
}

// Implementation returns the current BLAS float64 implementation.
//
// Implementation allows direct calls to the current the BLAS float64 implementation
// giving finer control of parameters.
func Implementation() blas.Float64 {
	return blas64
}

// Vector represents a vector with an associated element increment.
type Vector struct {
	Inc  int
	Data []float64
}

// General represents a matrix using the conventional storage scheme.
type General struct {
	Rows, Cols int
	Stride     int
	Data       []float64
}

// Band represents a band matrix using the band storage scheme.
type Band struct {
	Rows, Cols int
	KL, KU     int
	Stride     int
	Data       []float64
}

// Triangular represents a triangular matrix using the conventional storage scheme.
type Triangular struct {
	N      int
	Stride int
	Data   []float64
	Uplo   blas.Uplo
	Diag   blas.Diag
}

// TriangularBand represents a triangular matrix using the band storage scheme.
type TriangularBand struct {
	N, K   int
	Stride int
	Data   []float64
	Uplo   blas.Uplo
	Diag   blas.Diag
}

// TriangularPacked represents a triangular matrix using the packed storage scheme.
type TriangularPacked struct {
	N    int
	Data []float64
	Uplo blas.Uplo
	Diag blas.Diag
}

// Symmetric represents a symmetric matrix using the conventional storage scheme.
type Symmetric struct {
	N      int
	Stride int
	Data   []float64
	Uplo   blas.Uplo
}

// SymmetricBand represents a symmetric matrix using the band storage scheme.
type SymmetricBand struct {
	N, K   int
	Stride int
	Data   []float64
	Uplo   blas.Uplo
}

// SymmetricPacked represents a symmetric matrix using the packed storage scheme.
type SymmetricPacked struct {
	N    int
	Data []float64
	Uplo blas.Uplo
}

// Level 1

const negInc = "blas64: negative vector increment"

func Dot(n int, x, y Vector) float64 {
	return blas64.Ddot(n, x.Data, x.Inc, y.Data, y.Inc)
}

// Nrm2 will panic if the vector increment is negative.
func Nrm2(n int, x Vector) float64 {
	if x.Inc < 0 {
		panic(negInc)
	}
	return blas64.Dnrm2(n, x.Data, x.Inc)
}

// Asum will panic if the vector increment is negative.
func Asum(n int, x Vector) float64 {
	if x.Inc < 0 {
		panic(negInc)
	}
	return blas64.Dasum(n, x.Data, x.Inc)
}

// Iamax will panic if the vector increment is negative.
func Iamax(n int, x Vector) int {
	if x.Inc < 0 {
		panic(negInc)
	}
	return blas64.Idamax(n, x.Data, x.Inc)
}

func Swap(n int, x, y Vector) {
	blas64.Dswap(n, x.Data, x.Inc, y.Data, y.Inc)
}

func Copy(n int, x, y Vector) {
	blas64.Dcopy(n, x.Data, x.Inc, y.Data, y.Inc)
}

func Axpy(n int, alpha float64, x, y Vector) {
	blas64.Daxpy(n, alpha, x.Data, x.Inc, y.Data, y.Inc)
}

func Rotg(a, b float64) (c, s, r, z float64) {
	return blas64.Drotg(a, b)
}

func Rotmg(d1, d2, b1, b2 float64) (p blas.DrotmParams, rd1, rd2, rb1 float64) {
	return blas64.Drotmg(d1, d2, b1, b2)
}

func Rot(n int, x, y Vector, c, s float64) {
	blas64.Drot(n, x.Data, x.Inc, y.Data, y.Inc, c, s)
}

func Rotm(n int, x, y Vector, p blas.DrotmParams) {
	blas64.Drotm(n, x.Data, x.Inc, y.Data, y.Inc, p)
}

// Scal will panic if the vector increment is negative
func Scal(n int, alpha float64, x Vector) {
	if x.Inc < 0 {
		panic(negInc)
	}
	blas64.Dscal(n, alpha, x.Data, x.Inc)
}

// Level 2

func Gemv(tA blas.Transpose, alpha float64, a General, x Vector, beta float64, y Vector) {
	blas64.Dgemv(tA, a.Rows, a.Cols, alpha, a.Data, a.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Gbmv(tA blas.Transpose, alpha float64, a Band, x Vector, beta float64, y Vector) {
	blas64.Dgbmv(tA, a.Rows, a.Cols, a.KL, a.KU, alpha, a.Data, a.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Trmv(tA blas.Transpose, a Triangular, x Vector) {
	blas64.Dtrmv(a.Uplo, tA, a.Diag, a.N, a.Data, a.Stride, x.Data, x.Inc)
}

func Tbmv(tA blas.Transpose, a TriangularBand, x Vector) {
	blas64.Dtbmv(a.Uplo, tA, a.Diag, a.N, a.K, a.Data, a.Stride, x.Data, x.Inc)
}

func Tpmv(tA blas.Transpose, a TriangularPacked, x Vector) {
	blas64.Dtpmv(a.Uplo, tA, a.Diag, a.N, a.Data, x.Data, x.Inc)
}

func Trsv(tA blas.Transpose, a Triangular, x Vector) {
	blas64.Dtrsv(a.Uplo, tA, a.Diag, a.N, a.Data, a.Stride, x.Data, x.Inc)
}

func Tbsv(tA blas.Transpose, a TriangularBand, x Vector) {
	blas64.Dtbsv(a.Uplo, tA, a.Diag, a.N, a.K, a.Data, a.Stride, x.Data, x.Inc)
}

func Tpsv(tA blas.Transpose, a TriangularPacked, x Vector) {
	blas64.Dtpsv(a.Uplo, tA, a.Diag, a.N, a.Data, x.Data, x.Inc)
}

func Symv(alpha float64, a Symmetric, x Vector, beta float64, y Vector) {
	blas64.Dsymv(a.Uplo, a.N, alpha, a.Data, a.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Sbmv(alpha float64, a SymmetricBand, x Vector, beta float64, y Vector) {
	blas64.Dsbmv(a.Uplo, a.N, a.K, alpha, a.Data, a.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Spmv(alpha float64, a SymmetricPacked, x Vector, beta float64, y Vector) {
	blas64.Dspmv(a.Uplo, a.N, alpha, a.Data, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Ger(alpha float64, x, y Vector, a General) {
	blas64.Dger(a.Rows, a.Cols, alpha, x.Data, x.Inc, y.Data, y.Inc, a.Data, a.Stride)
}

func Syr(alpha float64, x Vector, a Symmetric) {
	blas64.Dsyr(a.Uplo, a.N, alpha, x.Data, x.Inc, a.Data, a.Stride)
}

func Spr(alpha float64, x Vector, a SymmetricPacked) {
	blas64.Dspr(a.Uplo, a.N, alpha, x.Data, x.Inc, a.Data)
}

func Syr2(alpha float64, x, y Vector, a Symmetric) {
	blas64.Dsyr2(a.Uplo, a.N, alpha, x.Data, x.Inc, y.Data, y.Inc, a.Data, a.Stride)
}

func Spr2(alpha float64, x, y Vector, a SymmetricPacked) {
	blas64.Dspr2(a.Uplo, a.N, alpha, x.Data, x.Inc, y.Data, y.Inc, a.Data)
}

// Level 3

func Gemm(tA, tB blas.Transpose, alpha float64, a, b General, beta float64, c General) {
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
	blas64.Dgemm(tA, tB, m, n, k, alpha, a.Data, a.Stride, b.Data, b.Stride, beta, c.Data, c.Stride)
}

func Symm(s blas.Side, alpha float64, a Symmetric, b General, beta float64, c General) {
	var m, n int
	if s == blas.Left {
		m, n = a.N, b.Cols
	} else {
		m, n = b.Rows, a.N
	}
	blas64.Dsymm(s, a.Uplo, m, n, alpha, a.Data, a.Stride, b.Data, b.Stride, beta, c.Data, c.Stride)
}

func Syrk(t blas.Transpose, alpha float64, a General, beta float64, c Symmetric) {
	var n, k int
	if t == blas.NoTrans {
		n, k = a.Rows, a.Cols
	} else {
		n, k = a.Cols, a.Rows
	}
	blas64.Dsyrk(c.Uplo, t, n, k, alpha, a.Data, a.Stride, beta, c.Data, c.Stride)
}

func Syr2k(t blas.Transpose, alpha float64, a, b General, beta float64, c Symmetric) {
	var n, k int
	if t == blas.NoTrans {
		n, k = a.Rows, a.Cols
	} else {
		n, k = a.Cols, a.Rows
	}
	blas64.Dsyr2k(c.Uplo, t, n, k, alpha, a.Data, a.Stride, b.Data, b.Stride, beta, c.Data, c.Stride)
}

func Trmm(s blas.Side, tA blas.Transpose, alpha float64, a Triangular, b General) {
	blas64.Dtrmm(s, a.Uplo, tA, a.Diag, b.Rows, b.Cols, alpha, a.Data, a.Stride, b.Data, b.Stride)
}

func Trsm(s blas.Side, tA blas.Transpose, alpha float64, a Triangular, b General) {
	blas64.Dtrsm(s, a.Uplo, tA, a.Diag, b.Rows, b.Cols, alpha, a.Data, a.Stride, b.Data, b.Stride)
}
