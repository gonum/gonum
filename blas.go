// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package blas provides interfaces for the BLAS linear algebra standard.

All methods must perform appropriate parameter checking and panic if
provided parameters that do not conform to the requirements specified
by the BLAS standard.

Quick Reference Guide to the BLAS from http://www.netlib.org/lapack/lug/node145.html

Level 1 BLAS

	        dim scalar vector   vector   scalars              5-element prefixes
	                                                          struct

	_rotg (                                      a, b )                S, D
	_rotmg(                              d1, d2, a, b )                S, D
	_rot  ( n,         x, incX, y, incY,               c, s )          S, D
	_rotm ( n,         x, incX, y, incY,                      param )  S, D
	_swap ( n,         x, incX, y, incY )                              S, D, C, Z
	_scal ( n,  alpha, x, incX )                                       S, D, C, Z, Cs, Zd
	_copy ( n,         x, incX, y, incY )                              S, D, C, Z
	_axpy ( n,  alpha, x, incX, y, incY )                              S, D, C, Z
	_dot  ( n,         x, incX, y, incY )                              S, D, Ds
	_dotu ( n,         x, incX, y, incY )                              C, Z
	_dotc ( n,         x, incX, y, incY )                              C, Z
	__dot ( n,  alpha, x, incX, y, incY )                              Sds
	_nrm2 ( n,         x, incX )                                       S, D, Sc, Dz
	_asum ( n,         x, incX )                                       S, D, Sc, Dz
	I_amax( n,         x, incX )                                       s, d, c, z

Level 2 BLAS

	        options                   dim   b-width scalar matrix  vector   scalar vector   prefixes

	_gemv ( order,        trans,      m, n,         alpha, a, lda, x, incX, beta,  y, incY ) S, D, C, Z
	_gbmv ( order,        trans,      m, n, kL, kU, alpha, a, lda, x, incX, beta,  y, incY ) S, D, C, Z
	_hemv ( order, uplo,                 n,         alpha, a, lda, x, incX, beta,  y, incY ) C, Z
	_hbmv ( order, uplo,                 n, k,      alpha, a, lda, x, incX, beta,  y, incY ) C, Z
	_hpmv ( order, uplo,                 n,         alpha, ap,     x, incX, beta,  y, incY ) C, Z
	_symv ( order, uplo,                 n,         alpha, a, lda, x, incX, beta,  y, incY ) S, D
	_sbmv ( order, uplo,                 n, k,      alpha, a, lda, x, incX, beta,  y, incY ) S, D
	_spmv ( order, uplo,                 n,         alpha, ap,     x, incX, beta,  y, incY ) S, D
	_trmv ( order, uplo, trans, diag,    n,                a, lda, x, incX )                 S, D, C, Z
	_tbmv ( order, uplo, trans, diag,    n, k,             a, lda, x, incX )                 S, D, C, Z
	_tpmv ( order, uplo, trans, diag,    n,                ap,     x, incX )                 S, D, C, Z
	_trsv ( order, uplo, trans, diag,    n,                a, lda, x, incX )                 S, D, C, Z
	_tbsv ( order, uplo, trans, diag,    n, k,             a, lda, x, incX )                 S, D, C, Z
	_tpsv ( order, uplo, trans, diag,    n,                ap,     x, incX )                 S, D, C, Z

	        options                   dim   scalar vector   vector   matrix  prefixes

	_ger  ( order,                    m, n, alpha, x, incX, y, incY, a, lda ) S, D
	_geru ( order,                    m, n, alpha, x, incX, y, incY, a, lda ) C, Z
	_gerc ( order,                    m, n, alpha, x, incX, y, incY, a, lda ) C, Z
	_her  ( order, uplo,                 n, alpha, x, incX,          a, lda ) C, Z
	_hpr  ( order, uplo,                 n, alpha, x, incX,          ap )     C, Z
	_her2 ( order, uplo,                 n, alpha, x, incX, y, incY, a, lda ) C, Z
	_hpr2 ( order, uplo,                 n, alpha, x, incX, y, incY, ap )     C, Z
	_syr  ( order, uplo,                 n, alpha, x, incX,          a, lda ) S, D
	_spr  ( order, uplo,                 n, alpha, x, incX,          ap )     S, D
	_syr2 ( order, uplo,                 n, alpha, x, incX, y, incY, a, lda ) S, D
	_spr2 ( order, uplo,                 n, alpha, x, incX, y, incY, ap )     S, D

Level 3 BLAS

	        options                                 dim      scalar matrix  matrix  scalar matrix  prefixes

	_gemm ( order,             transA, transB,      m, n, k, alpha, a, lda, b, ldb, beta,  c, ldc ) S, D, C, Z
	_symm ( order, side, uplo,                      m, n,    alpha, a, lda, b, ldb, beta,  c, ldc ) S, D, C, Z
	_hemm ( order, side, uplo,                      m, n,    alpha, a, lda, b, ldb, beta,  c, ldc ) C, Z
	_syrk ( order,       uplo, trans,                  n, k, alpha, a, lda,         beta,  c, ldc ) S, D, C, Z
	_herk ( order,       uplo, trans,                  n, k, alpha, a, lda,         beta,  c, ldc ) C, Z
	_syr2k( order,       uplo, trans,                  n, k, alpha, a, lda, b, ldb, beta,  c, ldc ) S, D, C, Z
	_her2k( order,       uplo, trans,                  n, k, alpha, a, lda, b, ldb, beta,  c, ldc ) C, Z
	_trmm ( order, side, uplo, transA,        diag, m, n,    alpha, a, lda, b, ldb )                S, D, C, Z
	_trsm ( order, side, uplo, transA,        diag, m, n,    alpha, a, lda, b, ldb )                S, D, C, Z

Meaning of prefixes

	S - float32	C - complex64
	D - float64	Z - complex128

Matrix types

	GE - GEneral 		GB - General Band
	SY - SYmmetric 		SB - Symmetric Band 	SP - Symmetric Packed
	HE - HErmitian 		HB - Hermitian Band 	HP - Hermitian Packed
	TR - TRiangular 	TB - Triangular Band 	TP - Triangular Packed

Options

	trans 	= NoTrans, Trans, ConjTrans
	uplo 	= Upper, Lower
	diag 	= Nonunit, Unit
	side 	= Left, Right (A or op(A) on the left, or A or op(A) on the right)

For real matrices, Trans and ConjTrans have the same meaning.
For Hermitian matrices, trans = Trans is not allowed.
For complex symmetric matrices, trans = ConjTrans is not allowed.
*/
package blas

// Type SrotmParams contains Givens transformation parameters returned
// by the Float32 Srotm method.
type SrotmParams struct {
	Flag float32
	H    [4]float32 // Column-major 2 by 2 matrix.
}

// Type DrotmParams contains Givens transformation parameters returned
// by the Float64 Drotm method.
type DrotmParams struct {
	Flag float64
	H    [4]float64 // Column-major 2 by 2 matrix.
}

// Type Order is used to specify the matrix storage format. An implementation
// may not implement both orders and must panic if a routine is called using
// an unimplemented order.
type Order int

const (
	RowMajor Order = 101 + iota
	ColMajor
)

// Type Transpose is used to specify the transposition operation for a
// routine.
type Transpose int

const (
	NoTrans Transpose = 111 + iota
	Trans
	ConjTrans
)

// Type Uplo is used to specify whether the matrix is an upper or lower
// triangular matrix.
type Uplo int

const (
	Upper Uplo = 121 + iota
	Lower
)

// Type Diag is used to specify whether the matrix is a unit or non-unit
// triangular matrix.
type Diag int

const (
	NonUnit Diag = 131 + iota
	Unit
)

// Type side is used to specify from which side a multiplication operation
// is performed.
type Side int

const (
	Left Side = 141 + iota
	Right
)

// Float32 implements the single precision real BLAS routines.
type Float32 interface {
	// Level 1 routines.
	Sdsdot(n int, alpha float32, x []float32, incX int, y []float32, incY int) float32
	Dsdot(n int, x []float32, incX int, y []float32, incY int) float64
	Sdot(n int, x []float32, incX int, y []float32, incY int) float32
	Snrm2(n int, x []float32, incX int) float32
	Sasum(n int, x []float32, incX int) float32
	Isamax(n int, x []float32, incX int) int
	Sswap(n int, x []float32, incX int, y []float32, incY int)
	Scopy(n int, x []float32, incX int, y []float32, incY int)
	Saxpy(n int, alpha float32, x []float32, incX int, y []float32, incY int)
	Srotg(a, b float32) (c, s, r, z float32)
	Srotmg(d1, d2, b1, b2 float32) (p *SrotmParams, rd1, rd2, rb1 float32)
	Srot(n int, x []float32, incX int, y []float32, incY int, c, s float32)
	Srotm(n int, x []float32, incX int, y []float32, incY int, p *SrotmParams)
	Sscal(n int, alpha float32, x []float32, incX int)

	// Level 2 routines.
	Sgemv(o Order, tA Transpose, m, n int, alpha float32, a []float32, lda int, x []float32, incX int, beta float32, y []float32, incY int)
	Sgbmv(o Order, tA Transpose, m, n, kL, kU int, alpha float32, a []float32, lda int, x []float32, incX int, beta float32, y []float32, incY int)
	Strmv(o Order, ul Uplo, tA Transpose, d Diag, n int, a []float32, lda int, x []float32, incX int)
	Stbmv(o Order, ul Uplo, tA Transpose, d Diag, n, k int, a []float32, lda int, x []float32, incX int)
	Stpmv(o Order, ul Uplo, tA Transpose, d Diag, n int, ap []float32, x []float32, incX int)
	Strsv(o Order, ul Uplo, tA Transpose, d Diag, n int, a []float32, lda int, x []float32, incX int)
	Stbsv(o Order, ul Uplo, tA Transpose, d Diag, n, k int, a []float32, lda int, x []float32, incX int)
	Stpsv(o Order, ul Uplo, tA Transpose, d Diag, n int, ap []float32, x []float32, incX int)
	Ssymv(o Order, ul Uplo, n int, alpha float32, a []float32, lda int, x []float32, incX int, beta float32, y []float32, incY int)
	Ssbmv(o Order, ul Uplo, n, k int, alpha float32, a []float32, lda int, x []float32, incX int, beta float32, y []float32, incY int)
	Sspmv(o Order, ul Uplo, n int, alpha float32, ap []float32, x []float32, incX int, beta float32, y []float32, incY int)
	Sger(o Order, m, n int, alpha float32, x []float32, incX int, y []float32, incY int, a []float32, lda int)
	Ssyr(o Order, ul Uplo, n int, alpha float32, x []float32, incX int, a []float32, lda int)
	Sspr(o Order, ul Uplo, n int, alpha float32, x []float32, incX int, ap []float32)
	Ssyr2(o Order, ul Uplo, n int, alpha float32, x []float32, incX int, y []float32, incY int, a []float32, lda int)
	Sspr2(o Order, ul Uplo, n int, alpha float32, x []float32, incX int, y []float32, incY int, a []float32)

	// Level 3 routines.
	Sgemm(o Order, tA, tB Transpose, m, n, k int, alpha float32, a []float32, lda int, b []float32, ldb int, beta float32, c []float32, ldc int)
	Ssymm(o Order, s Side, ul Uplo, m, n int, alpha float32, a []float32, lda int, b []float32, ldb int, beta float32, c []float32, ldc int)
	Ssyrk(o Order, ul Uplo, t Transpose, n, k int, alpha float32, a []float32, lda int, beta float32, c []float32, ldc int)
	Ssyr2k(o Order, ul Uplo, t Transpose, n, k int, alpha float32, a []float32, lda int, b []float32, ldb int, beta float32, c []float32, ldc int)
	Strmm(o Order, s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha float32, a []float32, lda int, b []float32, ldb int)
	Strsm(o Order, s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha float32, a []float32, lda int, b []float32, ldb int)
}

// Float64 implements the double precision real BLAS routines.
type Float64 interface {
	// Level 1 routines.
	Ddot(n int, x []float64, incX int, y []float64, incY int) float64
	Dnrm2(n int, x []float64, incX int) float64
	Dasum(n int, x []float64, incX int) float64
	Idamax(n int, x []float64, incX int) int
	Dswap(n int, x []float64, incX int, y []float64, incY int)
	Dcopy(n int, x []float64, incX int, y []float64, incY int)
	Daxpy(n int, alpha float64, x []float64, incX int, y []float64, incY int)
	Drotg(a, b float64) (c, s, r, z float64)
	Drotmg(d1, d2, b1, b2 float64) (p *DrotmParams, rd1, rd2, rb1 float64)
	Drot(n int, x []float64, incX int, y []float64, incY int, c float64, s float64)
	Drotm(n int, x []float64, incX int, y []float64, incY int, p *DrotmParams)
	Dscal(n int, alpha float64, x []float64, incX int)

	// Level 2 routines.
	Dgemv(o Order, tA Transpose, m, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
	Dgbmv(o Order, tA Transpose, m, n, kL, kU int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
	Dtrmv(o Order, ul Uplo, tA Transpose, d Diag, n int, a []float64, lda int, x []float64, incX int)
	Dtbmv(o Order, ul Uplo, tA Transpose, d Diag, n, k int, a []float64, lda int, x []float64, incX int)
	Dtpmv(o Order, ul Uplo, tA Transpose, d Diag, n int, ap []float64, x []float64, incX int)
	Dtrsv(o Order, ul Uplo, tA Transpose, d Diag, n int, a []float64, lda int, x []float64, incX int)
	Dtbsv(o Order, ul Uplo, tA Transpose, d Diag, n, k int, a []float64, lda int, x []float64, incX int)
	Dtpsv(o Order, ul Uplo, tA Transpose, d Diag, n int, ap []float64, x []float64, incX int)
	Dsymv(o Order, ul Uplo, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
	Dsbmv(o Order, ul Uplo, n, k int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
	Dspmv(o Order, ul Uplo, n int, alpha float64, ap []float64, x []float64, incX int, beta float64, y []float64, incY int)
	Dger(o Order, m, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int)
	Dsyr(o Order, ul Uplo, n int, alpha float64, x []float64, incX int, a []float64, lda int)
	Dspr(o Order, ul Uplo, n int, alpha float64, x []float64, incX int, ap []float64)
	Dsyr2(o Order, ul Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int)
	Dspr2(o Order, ul Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64)

	// Level 3 routines.
	Dgemm(o Order, tA, tB Transpose, m, n, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int)
	Dsymm(o Order, s Side, ul Uplo, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int)
	Dsyrk(o Order, ul Uplo, t Transpose, n, k int, alpha float64, a []float64, lda int, beta float64, c []float64, ldc int)
	Dsyr2k(o Order, ul Uplo, t Transpose, n, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int)
	Dtrmm(o Order, s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int)
	Dtrsm(o Order, s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int)
}

// Complex64 implements the single precision complex BLAS routines.
type Complex64 interface {
	// Level 1 routines
	Cdotu(n int, x []complex64, incX int, y []complex64, incY int) (dotu complex64)
	Cdotc(n int, x []complex64, incX int, y []complex64, incY int) (dotc complex64)
	Scnrm2(n int, x []complex64, incX int) float32
	Scasum(n int, x []complex64, incX int) float32
	Icamax(n int, x []complex64, incX int) int
	Cswap(n int, x []complex64, incX int, y []complex64, incY int)
	Ccopy(n int, x []complex64, incX int, y []complex64, incY int)
	Caxpy(n int, alpha complex64, x []complex64, incX int, y []complex64, incY int)
	Cscal(n int, alpha complex64, x []complex64, incX int)
	Csscal(n int, alpha float32, x []complex64, incX int)

	// Level 2 routines.
	Cgemv(o Order, tA Transpose, m, n int, alpha complex64, a []complex64, lda int, x []complex64, incX int, beta complex64, y []complex64, incY int)
	Cgbmv(o Order, tA Transpose, m, n, kL, kU int, alpha complex64, a []complex64, lda int, x []complex64, incX int, beta complex64, y []complex64, incY int)
	Ctrmv(o Order, ul Uplo, tA Transpose, d Diag, n int, a []complex64, lda int, x []complex64, incX int)
	Ctbmv(o Order, ul Uplo, tA Transpose, d Diag, n, k int, a []complex64, lda int, x []complex64, incX int)
	Ctpmv(o Order, ul Uplo, tA Transpose, d Diag, n int, ap []complex64, x []complex64, incX int)
	Ctrsv(o Order, ul Uplo, tA Transpose, d Diag, n int, a []complex64, lda int, x []complex64, incX int)
	Ctbsv(o Order, ul Uplo, tA Transpose, d Diag, n, k int, a []complex64, lda int, x []complex64, incX int)
	Ctpsv(o Order, ul Uplo, tA Transpose, d Diag, n int, ap []complex64, x []complex64, incX int)
	Chemv(o Order, ul Uplo, n int, alpha complex64, a []complex64, lda int, x []complex64, incX int, beta complex64, y []complex64, incY int)
	Chbmv(o Order, ul Uplo, n, k int, alpha complex64, a []complex64, lda int, x []complex64, incX int, beta complex64, y []complex64, incY int)
	Chpmv(o Order, ul Uplo, n int, alpha complex64, ap []complex64, x []complex64, incX int, beta complex64, y []complex64, incY int)
	Cgeru(o Order, m, n int, alpha complex64, x []complex64, incX int, y []complex64, incY int, a []complex64, lda int)
	Cgerc(o Order, m, n int, alpha complex64, x []complex64, incX int, y []complex64, incY int, a []complex64, lda int)
	Cher(o Order, ul Uplo, n int, alpha float32, x []complex64, incX int, a []complex64, lda int)
	Chpr(o Order, ul Uplo, n int, alpha float32, x []complex64, incX int, a []complex64)
	Cher2(o Order, ul Uplo, n int, alpha complex64, x []complex64, incX int, y []complex64, incY int, a []complex64, lda int)
	Chpr2(o Order, ul Uplo, n int, alpha complex64, x []complex64, incX int, y []complex64, incY int, ap []complex64)

	// Level 3 routines.
	Cgemm(o Order, tA, tB Transpose, m, n, k int, alpha complex64, a []complex64, lda int, b []complex64, ldb int, beta complex64, c []complex64, ldc int)
	Csymm(o Order, s Side, ul Uplo, m, n int, alpha complex64, a []complex64, lda int, b []complex64, ldb int, beta complex64, c []complex64, ldc int)
	Csyrk(o Order, ul Uplo, t Transpose, n, k int, alpha complex64, a []complex64, lda int, beta complex64, c []complex64, ldc int)
	Csyr2k(o Order, ul Uplo, t Transpose, n, k int, alpha complex64, a []complex64, lda int, b []complex64, ldb int, beta complex64, c []complex64, ldc int)
	Ctrmm(o Order, s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha complex64, a []complex64, lda int, b []complex64, ldb int)
	Ctrsm(o Order, s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha complex64, a []complex64, lda int, b []complex64, ldb int)
	Chemm(o Order, s Side, ul Uplo, m, n int, alpha complex64, a []complex64, lda int, b []complex64, ldb int, beta complex64, c []complex64, ldc int)
	Cherk(o Order, ul Uplo, t Transpose, n, k int, alpha float32, a []complex64, lda int, beta float32, c []complex64, ldc int)
	Cher2k(o Order, ul Uplo, t Transpose, n, k int, alpha complex64, a []complex64, lda int, b []complex64, ldb int, beta float32, c []complex64, ldc int)
}

// Complex128 implements the double precision complex BLAS routines.
type Complex128 interface {
	// Level 1 routines.
	Zdotu(n int, x []complex128, incX int, y []complex128, incY int) (dotu complex128)
	Zdotc(n int, x []complex128, incX int, y []complex128, incY int) (dotc complex128)
	Dznrm2(n int, x []complex128, incX int) float64
	Dzasum(n int, x []complex128, incX int) float64
	Izamax(n int, x []complex128, incX int) int
	Zswap(n int, x []complex128, incX int, y []complex128, incY int)
	Zcopy(n int, x []complex128, incX int, y []complex128, incY int)
	Zaxpy(n int, alpha complex128, x []complex128, incX int, y []complex128, incY int)
	Zscal(n int, alpha complex128, x []complex128, incX int)
	Zdscal(n int, alpha float64, x []complex128, incX int)

	// Level 2 routines.
	Zgemv(o Order, tA Transpose, m, n int, alpha complex128, a []complex128, lda int, x []complex128, incX int, beta complex128, y []complex128, incY int)
	Zgbmv(o Order, tA Transpose, m, n int, kL int, kU int, alpha complex128, a []complex128, lda int, x []complex128, incX int, beta complex128, y []complex128, incY int)
	Ztrmv(o Order, ul Uplo, tA Transpose, d Diag, n int, a []complex128, lda int, x []complex128, incX int)
	Ztbmv(o Order, ul Uplo, tA Transpose, d Diag, n, k int, a []complex128, lda int, x []complex128, incX int)
	Ztpmv(o Order, ul Uplo, tA Transpose, d Diag, n int, ap []complex128, x []complex128, incX int)
	Ztrsv(o Order, ul Uplo, tA Transpose, d Diag, n int, a []complex128, lda int, x []complex128, incX int)
	Ztbsv(o Order, ul Uplo, tA Transpose, d Diag, n, k int, a []complex128, lda int, x []complex128, incX int)
	Ztpsv(o Order, ul Uplo, tA Transpose, d Diag, n int, ap []complex128, x []complex128, incX int)
	Zhemv(o Order, ul Uplo, n int, alpha complex128, a []complex128, lda int, x []complex128, incX int, beta complex128, y []complex128, incY int)
	Zhbmv(o Order, ul Uplo, n, k int, alpha complex128, a []complex128, lda int, x []complex128, incX int, beta complex128, y []complex128, incY int)
	Zhpmv(o Order, ul Uplo, n int, alpha complex128, ap []complex128, x []complex128, incX int, beta complex128, y []complex128, incY int)
	Zgeru(o Order, m, n int, alpha complex128, x []complex128, incX int, y []complex128, incY int, a []complex128, lda int)
	Zgerc(o Order, m, n int, alpha complex128, x []complex128, incX int, y []complex128, incY int, a []complex128, lda int)
	Zher(o Order, ul Uplo, n int, alpha float64, x []complex128, incX int, a []complex128, lda int)
	Zhpr(o Order, ul Uplo, n int, alpha float64, x []complex128, incX int, a []complex128)
	Zher2(o Order, ul Uplo, n int, alpha complex128, x []complex128, incX int, y []complex128, incY int, a []complex128, lda int)
	Zhpr2(o Order, ul Uplo, n int, alpha complex128, x []complex128, incX int, y []complex128, incY int, ap []complex128)

	// Level 3 routines.
	Zgemm(o Order, tA, tB Transpose, m, n, k int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta complex128, c []complex128, ldc int)
	Zsymm(o Order, s Side, ul Uplo, m, n int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta complex128, c []complex128, ldc int)
	Zsyrk(o Order, ul Uplo, t Transpose, n, k int, alpha complex128, a []complex128, lda int, beta complex128, c []complex128, ldc int)
	Zsyr2k(o Order, ul Uplo, t Transpose, n, k int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta complex128, c []complex128, ldc int)
	Ztrmm(o Order, s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha complex128, a []complex128, lda int, b []complex128, ldb int)
	Ztrsm(o Order, s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha complex128, a []complex128, lda int, b []complex128, ldb int)
	Zhemm(o Order, s Side, ul Uplo, m, n int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta complex128, c []complex128, ldc int)
	Zherk(o Order, ul Uplo, t Transpose, n, k int, alpha float64, a []complex128, lda int, beta float64, c []complex128, ldc int)
	Zher2k(o Order, ul Uplo, t Transpose, n, k int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta float64, c []complex128, ldc int)
}
