// Do not manually edit this file. It was created by the genLapack.pl script from lapacke.h.

// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package clapack provides bindings to a C LAPACK library.
package clapack

/*
#cgo CFLAGS: -g -O2
#include "lapacke.h"
*/
import "C"

import (
	"github.com/gonum/blas"
	"github.com/gonum/lapack"
	"unsafe"
)

// Type order is used to specify the matrix storage format. We still interact with
// an API that allows client calls to specify order, so this is here to document that fact.
type order int

const (
	rowMajor order = 101 + iota
	colMajor
)

func isZero(ret C.int) bool { return ret == 0 }

func Sbdsdc(ul blas.Uplo, compq lapack.CompSV, n int, d []float32, e []float32, u []float32, ldu int, vt []float32, ldvt int, q []float32, iq []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sbdsdc((C.int)(rowMajor), (C.char)(ul), (C.char)(compq), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&u[0]), (C.lapack_int)(ldu), (*C.float)(&vt[0]), (C.lapack_int)(ldvt), (*C.float)(&q[0]), (*C.lapack_int)(&iq[0])))
}

func Dbdsdc(ul blas.Uplo, compq lapack.CompSV, n int, d []float64, e []float64, u []float64, ldu int, vt []float64, ldvt int, q []float64, iq []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dbdsdc((C.int)(rowMajor), (C.char)(ul), (C.char)(compq), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&u[0]), (C.lapack_int)(ldu), (*C.double)(&vt[0]), (C.lapack_int)(ldvt), (*C.double)(&q[0]), (*C.lapack_int)(&iq[0])))
}

func Sbdsqr(ul blas.Uplo, n int, ncvt int, nru int, ncc int, d []float32, e []float32, vt []float32, ldvt int, u []float32, ldu int, c []float32, ldc int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sbdsqr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ncvt), (C.lapack_int)(nru), (C.lapack_int)(ncc), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&vt[0]), (C.lapack_int)(ldvt), (*C.float)(&u[0]), (C.lapack_int)(ldu), (*C.float)(&c[0]), (C.lapack_int)(ldc)))
}

func Dbdsqr(ul blas.Uplo, n int, ncvt int, nru int, ncc int, d []float64, e []float64, vt []float64, ldvt int, u []float64, ldu int, c []float64, ldc int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dbdsqr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ncvt), (C.lapack_int)(nru), (C.lapack_int)(ncc), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&vt[0]), (C.lapack_int)(ldvt), (*C.double)(&u[0]), (C.lapack_int)(ldu), (*C.double)(&c[0]), (C.lapack_int)(ldc)))
}

func Cbdsqr(ul blas.Uplo, n int, ncvt int, nru int, ncc int, d []float32, e []float32, vt []complex64, ldvt int, u []complex64, ldu int, c []complex64, ldc int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cbdsqr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ncvt), (C.lapack_int)(nru), (C.lapack_int)(ncc), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.lapack_complex_float)(&vt[0]), (C.lapack_int)(ldvt), (*C.lapack_complex_float)(&u[0]), (C.lapack_int)(ldu), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc)))
}

func Zbdsqr(ul blas.Uplo, n int, ncvt int, nru int, ncc int, d []float64, e []float64, vt []complex128, ldvt int, u []complex128, ldu int, c []complex128, ldc int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zbdsqr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ncvt), (C.lapack_int)(nru), (C.lapack_int)(ncc), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.lapack_complex_double)(&vt[0]), (C.lapack_int)(ldvt), (*C.lapack_complex_double)(&u[0]), (C.lapack_int)(ldu), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc)))
}

func Sdisna(job lapack.Job, m int, n int, d []float32, sep []float32) bool {
	return isZero(C.LAPACKE_sdisna((C.char)(job), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&sep[0])))
}

func Ddisna(job lapack.Job, m int, n int, d []float64, sep []float64) bool {
	return isZero(C.LAPACKE_ddisna((C.char)(job), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&sep[0])))
}

func Sgbbrd(vect byte, m int, n int, ncc int, kl int, ku int, ab []float32, ldab int, d []float32, e []float32, q []float32, ldq int, pt []float32, ldpt int, c []float32, ldc int) bool {
	return isZero(C.LAPACKE_sgbbrd((C.int)(rowMajor), (C.char)(vect), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(ncc), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&q[0]), (C.lapack_int)(ldq), (*C.float)(&pt[0]), (C.lapack_int)(ldpt), (*C.float)(&c[0]), (C.lapack_int)(ldc)))
}

func Dgbbrd(vect byte, m int, n int, ncc int, kl int, ku int, ab []float64, ldab int, d []float64, e []float64, q []float64, ldq int, pt []float64, ldpt int, c []float64, ldc int) bool {
	return isZero(C.LAPACKE_dgbbrd((C.int)(rowMajor), (C.char)(vect), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(ncc), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&q[0]), (C.lapack_int)(ldq), (*C.double)(&pt[0]), (C.lapack_int)(ldpt), (*C.double)(&c[0]), (C.lapack_int)(ldc)))
}

func Cgbbrd(vect byte, m int, n int, ncc int, kl int, ku int, ab []complex64, ldab int, d []float32, e []float32, q []complex64, ldq int, pt []complex64, ldpt int, c []complex64, ldc int) bool {
	return isZero(C.LAPACKE_cgbbrd((C.int)(rowMajor), (C.char)(vect), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(ncc), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.lapack_complex_float)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_complex_float)(&pt[0]), (C.lapack_int)(ldpt), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc)))
}

func Zgbbrd(vect byte, m int, n int, ncc int, kl int, ku int, ab []complex128, ldab int, d []float64, e []float64, q []complex128, ldq int, pt []complex128, ldpt int, c []complex128, ldc int) bool {
	return isZero(C.LAPACKE_zgbbrd((C.int)(rowMajor), (C.char)(vect), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(ncc), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.lapack_complex_double)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_complex_double)(&pt[0]), (C.lapack_int)(ldpt), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc)))
}

func Sgbcon(norm byte, n int, kl int, ku int, ab []float32, ldab int, ipiv []int32, anorm float32, rcond []float32) bool {
	return isZero(C.LAPACKE_sgbcon((C.int)(rowMajor), (C.char)(norm), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0]), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Dgbcon(norm byte, n int, kl int, ku int, ab []float64, ldab int, ipiv []int32, anorm float64, rcond []float64) bool {
	return isZero(C.LAPACKE_dgbcon((C.int)(rowMajor), (C.char)(norm), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0]), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Cgbcon(norm byte, n int, kl int, ku int, ab []complex64, ldab int, ipiv []int32, anorm float32, rcond []float32) bool {
	return isZero(C.LAPACKE_cgbcon((C.int)(rowMajor), (C.char)(norm), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0]), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Zgbcon(norm byte, n int, kl int, ku int, ab []complex128, ldab int, ipiv []int32, anorm float64, rcond []float64) bool {
	return isZero(C.LAPACKE_zgbcon((C.int)(rowMajor), (C.char)(norm), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0]), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Sgbequ(m int, n int, kl int, ku int, ab []float32, ldab int, r []float32, c []float32, rowcnd []float32, colcnd []float32, amax []float32) bool {
	return isZero(C.LAPACKE_sgbequ((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&r[0]), (*C.float)(&c[0]), (*C.float)(&rowcnd[0]), (*C.float)(&colcnd[0]), (*C.float)(&amax[0])))
}

func Dgbequ(m int, n int, kl int, ku int, ab []float64, ldab int, r []float64, c []float64, rowcnd []float64, colcnd []float64, amax []float64) bool {
	return isZero(C.LAPACKE_dgbequ((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&r[0]), (*C.double)(&c[0]), (*C.double)(&rowcnd[0]), (*C.double)(&colcnd[0]), (*C.double)(&amax[0])))
}

func Cgbequ(m int, n int, kl int, ku int, ab []complex64, ldab int, r []float32, c []float32, rowcnd []float32, colcnd []float32, amax []float32) bool {
	return isZero(C.LAPACKE_cgbequ((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&r[0]), (*C.float)(&c[0]), (*C.float)(&rowcnd[0]), (*C.float)(&colcnd[0]), (*C.float)(&amax[0])))
}

func Zgbequ(m int, n int, kl int, ku int, ab []complex128, ldab int, r []float64, c []float64, rowcnd []float64, colcnd []float64, amax []float64) bool {
	return isZero(C.LAPACKE_zgbequ((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&r[0]), (*C.double)(&c[0]), (*C.double)(&rowcnd[0]), (*C.double)(&colcnd[0]), (*C.double)(&amax[0])))
}

func Sgbequb(m int, n int, kl int, ku int, ab []float32, ldab int, r []float32, c []float32, rowcnd []float32, colcnd []float32, amax []float32) bool {
	return isZero(C.LAPACKE_sgbequb((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&r[0]), (*C.float)(&c[0]), (*C.float)(&rowcnd[0]), (*C.float)(&colcnd[0]), (*C.float)(&amax[0])))
}

func Dgbequb(m int, n int, kl int, ku int, ab []float64, ldab int, r []float64, c []float64, rowcnd []float64, colcnd []float64, amax []float64) bool {
	return isZero(C.LAPACKE_dgbequb((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&r[0]), (*C.double)(&c[0]), (*C.double)(&rowcnd[0]), (*C.double)(&colcnd[0]), (*C.double)(&amax[0])))
}

func Cgbequb(m int, n int, kl int, ku int, ab []complex64, ldab int, r []float32, c []float32, rowcnd []float32, colcnd []float32, amax []float32) bool {
	return isZero(C.LAPACKE_cgbequb((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&r[0]), (*C.float)(&c[0]), (*C.float)(&rowcnd[0]), (*C.float)(&colcnd[0]), (*C.float)(&amax[0])))
}

func Zgbequb(m int, n int, kl int, ku int, ab []complex128, ldab int, r []float64, c []float64, rowcnd []float64, colcnd []float64, amax []float64) bool {
	return isZero(C.LAPACKE_zgbequb((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&r[0]), (*C.double)(&c[0]), (*C.double)(&rowcnd[0]), (*C.double)(&colcnd[0]), (*C.double)(&amax[0])))
}

func Sgbrfs(trans blas.Transpose, n int, kl int, ku int, nrhs int, ab []float32, ldab int, afb []float32, ldafb int, ipiv []int32, b []float32, ldb int, x []float32, ldx int, ferr []float32, berr []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sgbrfs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&afb[0]), (C.lapack_int)(ldafb), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dgbrfs(trans blas.Transpose, n int, kl int, ku int, nrhs int, ab []float64, ldab int, afb []float64, ldafb int, ipiv []int32, b []float64, ldb int, x []float64, ldx int, ferr []float64, berr []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dgbrfs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&afb[0]), (C.lapack_int)(ldafb), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Cgbrfs(trans blas.Transpose, n int, kl int, ku int, nrhs int, ab []complex64, ldab int, afb []complex64, ldafb int, ipiv []int32, b []complex64, ldb int, x []complex64, ldx int, ferr []float32, berr []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cgbrfs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_float)(&afb[0]), (C.lapack_int)(ldafb), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zgbrfs(trans blas.Transpose, n int, kl int, ku int, nrhs int, ab []complex128, ldab int, afb []complex128, ldafb int, ipiv []int32, b []complex128, ldb int, x []complex128, ldx int, ferr []float64, berr []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zgbrfs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_double)(&afb[0]), (C.lapack_int)(ldafb), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Sgbsv(n int, kl int, ku int, nrhs int, ab []float32, ldab int, ipiv []int32, b []float32, ldb int) bool {
	return isZero(C.LAPACKE_sgbsv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dgbsv(n int, kl int, ku int, nrhs int, ab []float64, ldab int, ipiv []int32, b []float64, ldb int) bool {
	return isZero(C.LAPACKE_dgbsv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cgbsv(n int, kl int, ku int, nrhs int, ab []complex64, ldab int, ipiv []int32, b []complex64, ldb int) bool {
	return isZero(C.LAPACKE_cgbsv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zgbsv(n int, kl int, ku int, nrhs int, ab []complex128, ldab int, ipiv []int32, b []complex128, ldb int) bool {
	return isZero(C.LAPACKE_zgbsv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Sgbsvx(fact byte, trans blas.Transpose, n int, kl int, ku int, nrhs int, ab []float32, ldab int, afb []float32, ldafb int, ipiv []int32, equed []byte, r []float32, c []float32, b []float32, ldb int, x []float32, ldx int, rcond []float32, ferr []float32, berr []float32, rpivot []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sgbsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&afb[0]), (C.lapack_int)(ldafb), (*C.lapack_int)(&ipiv[0]), (*C.char)(unsafe.Pointer(&equed[0])), (*C.float)(&r[0]), (*C.float)(&c[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0]), (*C.float)(&rpivot[0])))
}

func Dgbsvx(fact byte, trans blas.Transpose, n int, kl int, ku int, nrhs int, ab []float64, ldab int, afb []float64, ldafb int, ipiv []int32, equed []byte, r []float64, c []float64, b []float64, ldb int, x []float64, ldx int, rcond []float64, ferr []float64, berr []float64, rpivot []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dgbsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&afb[0]), (C.lapack_int)(ldafb), (*C.lapack_int)(&ipiv[0]), (*C.char)(unsafe.Pointer(&equed[0])), (*C.double)(&r[0]), (*C.double)(&c[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0]), (*C.double)(&rpivot[0])))
}

func Cgbsvx(fact byte, trans blas.Transpose, n int, kl int, ku int, nrhs int, ab []complex64, ldab int, afb []complex64, ldafb int, ipiv []int32, equed []byte, r []float32, c []float32, b []complex64, ldb int, x []complex64, ldx int, rcond []float32, ferr []float32, berr []float32, rpivot []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cgbsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_float)(&afb[0]), (C.lapack_int)(ldafb), (*C.lapack_int)(&ipiv[0]), (*C.char)(unsafe.Pointer(&equed[0])), (*C.float)(&r[0]), (*C.float)(&c[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0]), (*C.float)(&rpivot[0])))
}

func Zgbsvx(fact byte, trans blas.Transpose, n int, kl int, ku int, nrhs int, ab []complex128, ldab int, afb []complex128, ldafb int, ipiv []int32, equed []byte, r []float64, c []float64, b []complex128, ldb int, x []complex128, ldx int, rcond []float64, ferr []float64, berr []float64, rpivot []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zgbsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_double)(&afb[0]), (C.lapack_int)(ldafb), (*C.lapack_int)(&ipiv[0]), (*C.char)(unsafe.Pointer(&equed[0])), (*C.double)(&r[0]), (*C.double)(&c[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0]), (*C.double)(&rpivot[0])))
}

func Sgbtrf(m int, n int, kl int, ku int, ab []float32, ldab int, ipiv []int32) bool {
	return isZero(C.LAPACKE_sgbtrf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0])))
}

func Dgbtrf(m int, n int, kl int, ku int, ab []float64, ldab int, ipiv []int32) bool {
	return isZero(C.LAPACKE_dgbtrf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0])))
}

func Cgbtrf(m int, n int, kl int, ku int, ab []complex64, ldab int, ipiv []int32) bool {
	return isZero(C.LAPACKE_cgbtrf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0])))
}

func Zgbtrf(m int, n int, kl int, ku int, ab []complex128, ldab int, ipiv []int32) bool {
	return isZero(C.LAPACKE_zgbtrf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0])))
}

func Sgbtrs(trans blas.Transpose, n int, kl int, ku int, nrhs int, ab []float32, ldab int, ipiv []int32, b []float32, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sgbtrs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dgbtrs(trans blas.Transpose, n int, kl int, ku int, nrhs int, ab []float64, ldab int, ipiv []int32, b []float64, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dgbtrs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cgbtrs(trans blas.Transpose, n int, kl int, ku int, nrhs int, ab []complex64, ldab int, ipiv []int32, b []complex64, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cgbtrs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zgbtrs(trans blas.Transpose, n int, kl int, ku int, nrhs int, ab []complex128, ldab int, ipiv []int32, b []complex128, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zgbtrs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(kl), (C.lapack_int)(ku), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Sgebak(job lapack.Job, s blas.Side, n int, ilo int, ihi int, scale []float32, m int, v []float32, ldv int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	return isZero(C.LAPACKE_sgebak((C.int)(rowMajor), (C.char)(job), (C.char)(s), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.float)(&scale[0]), (C.lapack_int)(m), (*C.float)(&v[0]), (C.lapack_int)(ldv)))
}

func Dgebak(job lapack.Job, s blas.Side, n int, ilo int, ihi int, scale []float64, m int, v []float64, ldv int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	return isZero(C.LAPACKE_dgebak((C.int)(rowMajor), (C.char)(job), (C.char)(s), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.double)(&scale[0]), (C.lapack_int)(m), (*C.double)(&v[0]), (C.lapack_int)(ldv)))
}

func Cgebak(job lapack.Job, s blas.Side, n int, ilo int, ihi int, scale []float32, m int, v []complex64, ldv int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	return isZero(C.LAPACKE_cgebak((C.int)(rowMajor), (C.char)(job), (C.char)(s), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.float)(&scale[0]), (C.lapack_int)(m), (*C.lapack_complex_float)(&v[0]), (C.lapack_int)(ldv)))
}

func Zgebak(job lapack.Job, s blas.Side, n int, ilo int, ihi int, scale []float64, m int, v []complex128, ldv int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	return isZero(C.LAPACKE_zgebak((C.int)(rowMajor), (C.char)(job), (C.char)(s), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.double)(&scale[0]), (C.lapack_int)(m), (*C.lapack_complex_double)(&v[0]), (C.lapack_int)(ldv)))
}

func Sgebal(job lapack.Job, n int, a []float32, lda int, ilo []int32, ihi []int32, scale []float32) bool {
	return isZero(C.LAPACKE_sgebal((C.int)(rowMajor), (C.char)(job), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.float)(&scale[0])))
}

func Dgebal(job lapack.Job, n int, a []float64, lda int, ilo []int32, ihi []int32, scale []float64) bool {
	return isZero(C.LAPACKE_dgebal((C.int)(rowMajor), (C.char)(job), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.double)(&scale[0])))
}

func Cgebal(job lapack.Job, n int, a []complex64, lda int, ilo []int32, ihi []int32, scale []float32) bool {
	return isZero(C.LAPACKE_cgebal((C.int)(rowMajor), (C.char)(job), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.float)(&scale[0])))
}

func Zgebal(job lapack.Job, n int, a []complex128, lda int, ilo []int32, ihi []int32, scale []float64) bool {
	return isZero(C.LAPACKE_zgebal((C.int)(rowMajor), (C.char)(job), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.double)(&scale[0])))
}

func Sgebrd(m int, n int, a []float32, lda int, d []float32, e []float32, tauq []float32, taup []float32) bool {
	return isZero(C.LAPACKE_sgebrd((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&tauq[0]), (*C.float)(&taup[0])))
}

func Dgebrd(m int, n int, a []float64, lda int, d []float64, e []float64, tauq []float64, taup []float64) bool {
	return isZero(C.LAPACKE_dgebrd((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&tauq[0]), (*C.double)(&taup[0])))
}

func Cgebrd(m int, n int, a []complex64, lda int, d []float32, e []float32, tauq []complex64, taup []complex64) bool {
	return isZero(C.LAPACKE_cgebrd((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.lapack_complex_float)(&tauq[0]), (*C.lapack_complex_float)(&taup[0])))
}

func Zgebrd(m int, n int, a []complex128, lda int, d []float64, e []float64, tauq []complex128, taup []complex128) bool {
	return isZero(C.LAPACKE_zgebrd((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.lapack_complex_double)(&tauq[0]), (*C.lapack_complex_double)(&taup[0])))
}

func Sgecon(norm byte, n int, a []float32, lda int, anorm float32, rcond []float32) bool {
	return isZero(C.LAPACKE_sgecon((C.int)(rowMajor), (C.char)(norm), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Dgecon(norm byte, n int, a []float64, lda int, anorm float64, rcond []float64) bool {
	return isZero(C.LAPACKE_dgecon((C.int)(rowMajor), (C.char)(norm), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Cgecon(norm byte, n int, a []complex64, lda int, anorm float32, rcond []float32) bool {
	return isZero(C.LAPACKE_cgecon((C.int)(rowMajor), (C.char)(norm), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Zgecon(norm byte, n int, a []complex128, lda int, anorm float64, rcond []float64) bool {
	return isZero(C.LAPACKE_zgecon((C.int)(rowMajor), (C.char)(norm), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Sgeequ(m int, n int, a []float32, lda int, r []float32, c []float32, rowcnd []float32, colcnd []float32, amax []float32) bool {
	return isZero(C.LAPACKE_sgeequ((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&r[0]), (*C.float)(&c[0]), (*C.float)(&rowcnd[0]), (*C.float)(&colcnd[0]), (*C.float)(&amax[0])))
}

func Dgeequ(m int, n int, a []float64, lda int, r []float64, c []float64, rowcnd []float64, colcnd []float64, amax []float64) bool {
	return isZero(C.LAPACKE_dgeequ((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&r[0]), (*C.double)(&c[0]), (*C.double)(&rowcnd[0]), (*C.double)(&colcnd[0]), (*C.double)(&amax[0])))
}

func Cgeequ(m int, n int, a []complex64, lda int, r []float32, c []float32, rowcnd []float32, colcnd []float32, amax []float32) bool {
	return isZero(C.LAPACKE_cgeequ((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&r[0]), (*C.float)(&c[0]), (*C.float)(&rowcnd[0]), (*C.float)(&colcnd[0]), (*C.float)(&amax[0])))
}

func Zgeequ(m int, n int, a []complex128, lda int, r []float64, c []float64, rowcnd []float64, colcnd []float64, amax []float64) bool {
	return isZero(C.LAPACKE_zgeequ((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&r[0]), (*C.double)(&c[0]), (*C.double)(&rowcnd[0]), (*C.double)(&colcnd[0]), (*C.double)(&amax[0])))
}

func Sgeequb(m int, n int, a []float32, lda int, r []float32, c []float32, rowcnd []float32, colcnd []float32, amax []float32) bool {
	return isZero(C.LAPACKE_sgeequb((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&r[0]), (*C.float)(&c[0]), (*C.float)(&rowcnd[0]), (*C.float)(&colcnd[0]), (*C.float)(&amax[0])))
}

func Dgeequb(m int, n int, a []float64, lda int, r []float64, c []float64, rowcnd []float64, colcnd []float64, amax []float64) bool {
	return isZero(C.LAPACKE_dgeequb((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&r[0]), (*C.double)(&c[0]), (*C.double)(&rowcnd[0]), (*C.double)(&colcnd[0]), (*C.double)(&amax[0])))
}

func Cgeequb(m int, n int, a []complex64, lda int, r []float32, c []float32, rowcnd []float32, colcnd []float32, amax []float32) bool {
	return isZero(C.LAPACKE_cgeequb((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&r[0]), (*C.float)(&c[0]), (*C.float)(&rowcnd[0]), (*C.float)(&colcnd[0]), (*C.float)(&amax[0])))
}

func Zgeequb(m int, n int, a []complex128, lda int, r []float64, c []float64, rowcnd []float64, colcnd []float64, amax []float64) bool {
	return isZero(C.LAPACKE_zgeequb((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&r[0]), (*C.double)(&c[0]), (*C.double)(&rowcnd[0]), (*C.double)(&colcnd[0]), (*C.double)(&amax[0])))
}

func Sgeev(jobvl lapack.Job, jobvr lapack.Job, n int, a []float32, lda int, wr []float32, wi []float32, vl []float32, ldvl int, vr []float32, ldvr int) bool {
	return isZero(C.LAPACKE_sgeev((C.int)(rowMajor), (C.char)(jobvl), (C.char)(jobvr), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&wr[0]), (*C.float)(&wi[0]), (*C.float)(&vl[0]), (C.lapack_int)(ldvl), (*C.float)(&vr[0]), (C.lapack_int)(ldvr)))
}

func Dgeev(jobvl lapack.Job, jobvr lapack.Job, n int, a []float64, lda int, wr []float64, wi []float64, vl []float64, ldvl int, vr []float64, ldvr int) bool {
	return isZero(C.LAPACKE_dgeev((C.int)(rowMajor), (C.char)(jobvl), (C.char)(jobvr), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&wr[0]), (*C.double)(&wi[0]), (*C.double)(&vl[0]), (C.lapack_int)(ldvl), (*C.double)(&vr[0]), (C.lapack_int)(ldvr)))
}

func Cgeev(jobvl lapack.Job, jobvr lapack.Job, n int, a []complex64, lda int, w []complex64, vl []complex64, ldvl int, vr []complex64, ldvr int) bool {
	return isZero(C.LAPACKE_cgeev((C.int)(rowMajor), (C.char)(jobvl), (C.char)(jobvr), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&w[0]), (*C.lapack_complex_float)(&vl[0]), (C.lapack_int)(ldvl), (*C.lapack_complex_float)(&vr[0]), (C.lapack_int)(ldvr)))
}

func Zgeev(jobvl lapack.Job, jobvr lapack.Job, n int, a []complex128, lda int, w []complex128, vl []complex128, ldvl int, vr []complex128, ldvr int) bool {
	return isZero(C.LAPACKE_zgeev((C.int)(rowMajor), (C.char)(jobvl), (C.char)(jobvr), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&w[0]), (*C.lapack_complex_double)(&vl[0]), (C.lapack_int)(ldvl), (*C.lapack_complex_double)(&vr[0]), (C.lapack_int)(ldvr)))
}

func Sgeevx(balanc byte, jobvl lapack.Job, jobvr lapack.Job, sense byte, n int, a []float32, lda int, wr []float32, wi []float32, vl []float32, ldvl int, vr []float32, ldvr int, ilo []int32, ihi []int32, scale []float32, abnrm []float32, rconde []float32, rcondv []float32) bool {
	return isZero(C.LAPACKE_sgeevx((C.int)(rowMajor), (C.char)(balanc), (C.char)(jobvl), (C.char)(jobvr), (C.char)(sense), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&wr[0]), (*C.float)(&wi[0]), (*C.float)(&vl[0]), (C.lapack_int)(ldvl), (*C.float)(&vr[0]), (C.lapack_int)(ldvr), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.float)(&scale[0]), (*C.float)(&abnrm[0]), (*C.float)(&rconde[0]), (*C.float)(&rcondv[0])))
}

func Dgeevx(balanc byte, jobvl lapack.Job, jobvr lapack.Job, sense byte, n int, a []float64, lda int, wr []float64, wi []float64, vl []float64, ldvl int, vr []float64, ldvr int, ilo []int32, ihi []int32, scale []float64, abnrm []float64, rconde []float64, rcondv []float64) bool {
	return isZero(C.LAPACKE_dgeevx((C.int)(rowMajor), (C.char)(balanc), (C.char)(jobvl), (C.char)(jobvr), (C.char)(sense), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&wr[0]), (*C.double)(&wi[0]), (*C.double)(&vl[0]), (C.lapack_int)(ldvl), (*C.double)(&vr[0]), (C.lapack_int)(ldvr), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.double)(&scale[0]), (*C.double)(&abnrm[0]), (*C.double)(&rconde[0]), (*C.double)(&rcondv[0])))
}

func Cgeevx(balanc byte, jobvl lapack.Job, jobvr lapack.Job, sense byte, n int, a []complex64, lda int, w []complex64, vl []complex64, ldvl int, vr []complex64, ldvr int, ilo []int32, ihi []int32, scale []float32, abnrm []float32, rconde []float32, rcondv []float32) bool {
	return isZero(C.LAPACKE_cgeevx((C.int)(rowMajor), (C.char)(balanc), (C.char)(jobvl), (C.char)(jobvr), (C.char)(sense), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&w[0]), (*C.lapack_complex_float)(&vl[0]), (C.lapack_int)(ldvl), (*C.lapack_complex_float)(&vr[0]), (C.lapack_int)(ldvr), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.float)(&scale[0]), (*C.float)(&abnrm[0]), (*C.float)(&rconde[0]), (*C.float)(&rcondv[0])))
}

func Zgeevx(balanc byte, jobvl lapack.Job, jobvr lapack.Job, sense byte, n int, a []complex128, lda int, w []complex128, vl []complex128, ldvl int, vr []complex128, ldvr int, ilo []int32, ihi []int32, scale []float64, abnrm []float64, rconde []float64, rcondv []float64) bool {
	return isZero(C.LAPACKE_zgeevx((C.int)(rowMajor), (C.char)(balanc), (C.char)(jobvl), (C.char)(jobvr), (C.char)(sense), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&w[0]), (*C.lapack_complex_double)(&vl[0]), (C.lapack_int)(ldvl), (*C.lapack_complex_double)(&vr[0]), (C.lapack_int)(ldvr), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.double)(&scale[0]), (*C.double)(&abnrm[0]), (*C.double)(&rconde[0]), (*C.double)(&rcondv[0])))
}

func Sgehrd(n int, ilo int, ihi int, a []float32, lda int, tau []float32) bool {
	return isZero(C.LAPACKE_sgehrd((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dgehrd(n int, ilo int, ihi int, a []float64, lda int, tau []float64) bool {
	return isZero(C.LAPACKE_dgehrd((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Cgehrd(n int, ilo int, ihi int, a []complex64, lda int, tau []complex64) bool {
	return isZero(C.LAPACKE_cgehrd((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Zgehrd(n int, ilo int, ihi int, a []complex128, lda int, tau []complex128) bool {
	return isZero(C.LAPACKE_zgehrd((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Sgejsv(joba lapack.Job, jobu lapack.Job, jobv lapack.Job, jobr lapack.Job, jobt lapack.Job, jobp lapack.Job, m int, n int, a []float32, lda int, sva []float32, u []float32, ldu int, v []float32, ldv int, stat []float32, istat []int32) bool {
	return isZero(C.LAPACKE_sgejsv((C.int)(rowMajor), (C.char)(joba), (C.char)(jobu), (C.char)(jobv), (C.char)(jobr), (C.char)(jobt), (C.char)(jobp), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&sva[0]), (*C.float)(&u[0]), (C.lapack_int)(ldu), (*C.float)(&v[0]), (C.lapack_int)(ldv), (*C.float)(&stat[0]), (*C.lapack_int)(&istat[0])))
}

func Dgejsv(joba lapack.Job, jobu lapack.Job, jobv lapack.Job, jobr lapack.Job, jobt lapack.Job, jobp lapack.Job, m int, n int, a []float64, lda int, sva []float64, u []float64, ldu int, v []float64, ldv int, stat []float64, istat []int32) bool {
	return isZero(C.LAPACKE_dgejsv((C.int)(rowMajor), (C.char)(joba), (C.char)(jobu), (C.char)(jobv), (C.char)(jobr), (C.char)(jobt), (C.char)(jobp), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&sva[0]), (*C.double)(&u[0]), (C.lapack_int)(ldu), (*C.double)(&v[0]), (C.lapack_int)(ldv), (*C.double)(&stat[0]), (*C.lapack_int)(&istat[0])))
}

func Sgelq2(m int, n int, a []float32, lda int, tau []float32) bool {
	return isZero(C.LAPACKE_sgelq2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dgelq2(m int, n int, a []float64, lda int, tau []float64) bool {
	return isZero(C.LAPACKE_dgelq2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Cgelq2(m int, n int, a []complex64, lda int, tau []complex64) bool {
	return isZero(C.LAPACKE_cgelq2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Zgelq2(m int, n int, a []complex128, lda int, tau []complex128) bool {
	return isZero(C.LAPACKE_zgelq2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Sgelqf(m int, n int, a []float32, lda int, tau []float32) bool {
	return isZero(C.LAPACKE_sgelqf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dgelqf(m int, n int, a []float64, lda int, tau []float64) bool {
	return isZero(C.LAPACKE_dgelqf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Cgelqf(m int, n int, a []complex64, lda int, tau []complex64) bool {
	return isZero(C.LAPACKE_cgelqf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Zgelqf(m int, n int, a []complex128, lda int, tau []complex128) bool {
	return isZero(C.LAPACKE_zgelqf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Sgels(trans blas.Transpose, m int, n int, nrhs int, a []float32, lda int, b []float32, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sgels((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dgels(trans blas.Transpose, m int, n int, nrhs int, a []float64, lda int, b []float64, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dgels((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cgels(trans blas.Transpose, m int, n int, nrhs int, a []complex64, lda int, b []complex64, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cgels((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zgels(trans blas.Transpose, m int, n int, nrhs int, a []complex128, lda int, b []complex128, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zgels((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Sgelsd(m int, n int, nrhs int, a []float32, lda int, b []float32, ldb int, s []float32, rcond float32, rank []int32) bool {
	return isZero(C.LAPACKE_sgelsd((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&s[0]), (C.float)(rcond), (*C.lapack_int)(&rank[0])))
}

func Dgelsd(m int, n int, nrhs int, a []float64, lda int, b []float64, ldb int, s []float64, rcond float64, rank []int32) bool {
	return isZero(C.LAPACKE_dgelsd((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&s[0]), (C.double)(rcond), (*C.lapack_int)(&rank[0])))
}

func Cgelsd(m int, n int, nrhs int, a []complex64, lda int, b []complex64, ldb int, s []float32, rcond float32, rank []int32) bool {
	return isZero(C.LAPACKE_cgelsd((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&s[0]), (C.float)(rcond), (*C.lapack_int)(&rank[0])))
}

func Zgelsd(m int, n int, nrhs int, a []complex128, lda int, b []complex128, ldb int, s []float64, rcond float64, rank []int32) bool {
	return isZero(C.LAPACKE_zgelsd((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&s[0]), (C.double)(rcond), (*C.lapack_int)(&rank[0])))
}

func Sgelss(m int, n int, nrhs int, a []float32, lda int, b []float32, ldb int, s []float32, rcond float32, rank []int32) bool {
	return isZero(C.LAPACKE_sgelss((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&s[0]), (C.float)(rcond), (*C.lapack_int)(&rank[0])))
}

func Dgelss(m int, n int, nrhs int, a []float64, lda int, b []float64, ldb int, s []float64, rcond float64, rank []int32) bool {
	return isZero(C.LAPACKE_dgelss((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&s[0]), (C.double)(rcond), (*C.lapack_int)(&rank[0])))
}

func Cgelss(m int, n int, nrhs int, a []complex64, lda int, b []complex64, ldb int, s []float32, rcond float32, rank []int32) bool {
	return isZero(C.LAPACKE_cgelss((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&s[0]), (C.float)(rcond), (*C.lapack_int)(&rank[0])))
}

func Zgelss(m int, n int, nrhs int, a []complex128, lda int, b []complex128, ldb int, s []float64, rcond float64, rank []int32) bool {
	return isZero(C.LAPACKE_zgelss((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&s[0]), (C.double)(rcond), (*C.lapack_int)(&rank[0])))
}

func Sgelsy(m int, n int, nrhs int, a []float32, lda int, b []float32, ldb int, jpvt []int32, rcond float32, rank []int32) bool {
	return isZero(C.LAPACKE_sgelsy((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_int)(&jpvt[0]), (C.float)(rcond), (*C.lapack_int)(&rank[0])))
}

func Dgelsy(m int, n int, nrhs int, a []float64, lda int, b []float64, ldb int, jpvt []int32, rcond float64, rank []int32) bool {
	return isZero(C.LAPACKE_dgelsy((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_int)(&jpvt[0]), (C.double)(rcond), (*C.lapack_int)(&rank[0])))
}

func Cgelsy(m int, n int, nrhs int, a []complex64, lda int, b []complex64, ldb int, jpvt []int32, rcond float32, rank []int32) bool {
	return isZero(C.LAPACKE_cgelsy((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_int)(&jpvt[0]), (C.float)(rcond), (*C.lapack_int)(&rank[0])))
}

func Zgelsy(m int, n int, nrhs int, a []complex128, lda int, b []complex128, ldb int, jpvt []int32, rcond float64, rank []int32) bool {
	return isZero(C.LAPACKE_zgelsy((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_int)(&jpvt[0]), (C.double)(rcond), (*C.lapack_int)(&rank[0])))
}

func Sgeqlf(m int, n int, a []float32, lda int, tau []float32) bool {
	return isZero(C.LAPACKE_sgeqlf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dgeqlf(m int, n int, a []float64, lda int, tau []float64) bool {
	return isZero(C.LAPACKE_dgeqlf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Cgeqlf(m int, n int, a []complex64, lda int, tau []complex64) bool {
	return isZero(C.LAPACKE_cgeqlf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Zgeqlf(m int, n int, a []complex128, lda int, tau []complex128) bool {
	return isZero(C.LAPACKE_zgeqlf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Sgeqp3(m int, n int, a []float32, lda int, jpvt []int32, tau []float32) bool {
	return isZero(C.LAPACKE_sgeqp3((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&jpvt[0]), (*C.float)(&tau[0])))
}

func Dgeqp3(m int, n int, a []float64, lda int, jpvt []int32, tau []float64) bool {
	return isZero(C.LAPACKE_dgeqp3((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&jpvt[0]), (*C.double)(&tau[0])))
}

func Cgeqp3(m int, n int, a []complex64, lda int, jpvt []int32, tau []complex64) bool {
	return isZero(C.LAPACKE_cgeqp3((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&jpvt[0]), (*C.lapack_complex_float)(&tau[0])))
}

func Zgeqp3(m int, n int, a []complex128, lda int, jpvt []int32, tau []complex128) bool {
	return isZero(C.LAPACKE_zgeqp3((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&jpvt[0]), (*C.lapack_complex_double)(&tau[0])))
}

func Sgeqpf(m int, n int, a []float32, lda int, jpvt []int32, tau []float32) bool {
	return isZero(C.LAPACKE_sgeqpf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&jpvt[0]), (*C.float)(&tau[0])))
}

func Dgeqpf(m int, n int, a []float64, lda int, jpvt []int32, tau []float64) bool {
	return isZero(C.LAPACKE_dgeqpf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&jpvt[0]), (*C.double)(&tau[0])))
}

func Cgeqpf(m int, n int, a []complex64, lda int, jpvt []int32, tau []complex64) bool {
	return isZero(C.LAPACKE_cgeqpf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&jpvt[0]), (*C.lapack_complex_float)(&tau[0])))
}

func Zgeqpf(m int, n int, a []complex128, lda int, jpvt []int32, tau []complex128) bool {
	return isZero(C.LAPACKE_zgeqpf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&jpvt[0]), (*C.lapack_complex_double)(&tau[0])))
}

func Sgeqr2(m int, n int, a []float32, lda int, tau []float32) bool {
	return isZero(C.LAPACKE_sgeqr2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dgeqr2(m int, n int, a []float64, lda int, tau []float64) bool {
	return isZero(C.LAPACKE_dgeqr2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Cgeqr2(m int, n int, a []complex64, lda int, tau []complex64) bool {
	return isZero(C.LAPACKE_cgeqr2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Zgeqr2(m int, n int, a []complex128, lda int, tau []complex128) bool {
	return isZero(C.LAPACKE_zgeqr2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Sgeqrf(m int, n int, a []float32, lda int, tau []float32) bool {
	return isZero(C.LAPACKE_sgeqrf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dgeqrf(m int, n int, a []float64, lda int, tau []float64) bool {
	return isZero(C.LAPACKE_dgeqrf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Cgeqrf(m int, n int, a []complex64, lda int, tau []complex64) bool {
	return isZero(C.LAPACKE_cgeqrf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Zgeqrf(m int, n int, a []complex128, lda int, tau []complex128) bool {
	return isZero(C.LAPACKE_zgeqrf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Sgeqrfp(m int, n int, a []float32, lda int, tau []float32) bool {
	return isZero(C.LAPACKE_sgeqrfp((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dgeqrfp(m int, n int, a []float64, lda int, tau []float64) bool {
	return isZero(C.LAPACKE_dgeqrfp((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Cgeqrfp(m int, n int, a []complex64, lda int, tau []complex64) bool {
	return isZero(C.LAPACKE_cgeqrfp((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Zgeqrfp(m int, n int, a []complex128, lda int, tau []complex128) bool {
	return isZero(C.LAPACKE_zgeqrfp((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Sgerfs(trans blas.Transpose, n int, nrhs int, a []float32, lda int, af []float32, ldaf int, ipiv []int32, b []float32, ldb int, x []float32, ldx int, ferr []float32, berr []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sgerfs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dgerfs(trans blas.Transpose, n int, nrhs int, a []float64, lda int, af []float64, ldaf int, ipiv []int32, b []float64, ldb int, x []float64, ldx int, ferr []float64, berr []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dgerfs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Cgerfs(trans blas.Transpose, n int, nrhs int, a []complex64, lda int, af []complex64, ldaf int, ipiv []int32, b []complex64, ldb int, x []complex64, ldx int, ferr []float32, berr []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cgerfs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zgerfs(trans blas.Transpose, n int, nrhs int, a []complex128, lda int, af []complex128, ldaf int, ipiv []int32, b []complex128, ldb int, x []complex128, ldx int, ferr []float64, berr []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zgerfs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Sgerqf(m int, n int, a []float32, lda int, tau []float32) bool {
	return isZero(C.LAPACKE_sgerqf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dgerqf(m int, n int, a []float64, lda int, tau []float64) bool {
	return isZero(C.LAPACKE_dgerqf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Cgerqf(m int, n int, a []complex64, lda int, tau []complex64) bool {
	return isZero(C.LAPACKE_cgerqf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Zgerqf(m int, n int, a []complex128, lda int, tau []complex128) bool {
	return isZero(C.LAPACKE_zgerqf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Sgesdd(jobz lapack.Job, m int, n int, a []float32, lda int, s []float32, u []float32, ldu int, vt []float32, ldvt int) bool {
	return isZero(C.LAPACKE_sgesdd((C.int)(rowMajor), (C.char)(jobz), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&s[0]), (*C.float)(&u[0]), (C.lapack_int)(ldu), (*C.float)(&vt[0]), (C.lapack_int)(ldvt)))
}

func Dgesdd(jobz lapack.Job, m int, n int, a []float64, lda int, s []float64, u []float64, ldu int, vt []float64, ldvt int) bool {
	return isZero(C.LAPACKE_dgesdd((C.int)(rowMajor), (C.char)(jobz), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&s[0]), (*C.double)(&u[0]), (C.lapack_int)(ldu), (*C.double)(&vt[0]), (C.lapack_int)(ldvt)))
}

func Cgesdd(jobz lapack.Job, m int, n int, a []complex64, lda int, s []float32, u []complex64, ldu int, vt []complex64, ldvt int) bool {
	return isZero(C.LAPACKE_cgesdd((C.int)(rowMajor), (C.char)(jobz), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&s[0]), (*C.lapack_complex_float)(&u[0]), (C.lapack_int)(ldu), (*C.lapack_complex_float)(&vt[0]), (C.lapack_int)(ldvt)))
}

func Zgesdd(jobz lapack.Job, m int, n int, a []complex128, lda int, s []float64, u []complex128, ldu int, vt []complex128, ldvt int) bool {
	return isZero(C.LAPACKE_zgesdd((C.int)(rowMajor), (C.char)(jobz), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&s[0]), (*C.lapack_complex_double)(&u[0]), (C.lapack_int)(ldu), (*C.lapack_complex_double)(&vt[0]), (C.lapack_int)(ldvt)))
}

func Sgesv(n int, nrhs int, a []float32, lda int, ipiv []int32, b []float32, ldb int) bool {
	return isZero(C.LAPACKE_sgesv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dgesv(n int, nrhs int, a []float64, lda int, ipiv []int32, b []float64, ldb int) bool {
	return isZero(C.LAPACKE_dgesv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cgesv(n int, nrhs int, a []complex64, lda int, ipiv []int32, b []complex64, ldb int) bool {
	return isZero(C.LAPACKE_cgesv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zgesv(n int, nrhs int, a []complex128, lda int, ipiv []int32, b []complex128, ldb int) bool {
	return isZero(C.LAPACKE_zgesv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Dsgesv(n int, nrhs int, a []float64, lda int, ipiv []int32, b []float64, ldb int, x []float64, ldx int, iter []int32) bool {
	return isZero(C.LAPACKE_dsgesv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.lapack_int)(&iter[0])))
}

func Zcgesv(n int, nrhs int, a []complex128, lda int, ipiv []int32, b []complex128, ldb int, x []complex128, ldx int, iter []int32) bool {
	return isZero(C.LAPACKE_zcgesv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.lapack_int)(&iter[0])))
}

func Sgesvd(jobu lapack.Job, jobvt lapack.Job, m int, n int, a []float32, lda int, s []float32, u []float32, ldu int, vt []float32, ldvt int, superb []float32) bool {
	return isZero(C.LAPACKE_sgesvd((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobvt), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&s[0]), (*C.float)(&u[0]), (C.lapack_int)(ldu), (*C.float)(&vt[0]), (C.lapack_int)(ldvt), (*C.float)(&superb[0])))
}

func Dgesvd(jobu lapack.Job, jobvt lapack.Job, m int, n int, a []float64, lda int, s []float64, u []float64, ldu int, vt []float64, ldvt int, superb []float64) bool {
	return isZero(C.LAPACKE_dgesvd((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobvt), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&s[0]), (*C.double)(&u[0]), (C.lapack_int)(ldu), (*C.double)(&vt[0]), (C.lapack_int)(ldvt), (*C.double)(&superb[0])))
}

func Cgesvd(jobu lapack.Job, jobvt lapack.Job, m int, n int, a []complex64, lda int, s []float32, u []complex64, ldu int, vt []complex64, ldvt int, superb []float32) bool {
	return isZero(C.LAPACKE_cgesvd((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobvt), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&s[0]), (*C.lapack_complex_float)(&u[0]), (C.lapack_int)(ldu), (*C.lapack_complex_float)(&vt[0]), (C.lapack_int)(ldvt), (*C.float)(&superb[0])))
}

func Zgesvd(jobu lapack.Job, jobvt lapack.Job, m int, n int, a []complex128, lda int, s []float64, u []complex128, ldu int, vt []complex128, ldvt int, superb []float64) bool {
	return isZero(C.LAPACKE_zgesvd((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobvt), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&s[0]), (*C.lapack_complex_double)(&u[0]), (C.lapack_int)(ldu), (*C.lapack_complex_double)(&vt[0]), (C.lapack_int)(ldvt), (*C.double)(&superb[0])))
}

func Sgesvj(joba lapack.Job, jobu lapack.Job, jobv lapack.Job, m int, n int, a []float32, lda int, sva []float32, mv int, v []float32, ldv int, stat []float32) bool {
	return isZero(C.LAPACKE_sgesvj((C.int)(rowMajor), (C.char)(joba), (C.char)(jobu), (C.char)(jobv), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&sva[0]), (C.lapack_int)(mv), (*C.float)(&v[0]), (C.lapack_int)(ldv), (*C.float)(&stat[0])))
}

func Dgesvj(joba lapack.Job, jobu lapack.Job, jobv lapack.Job, m int, n int, a []float64, lda int, sva []float64, mv int, v []float64, ldv int, stat []float64) bool {
	return isZero(C.LAPACKE_dgesvj((C.int)(rowMajor), (C.char)(joba), (C.char)(jobu), (C.char)(jobv), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&sva[0]), (C.lapack_int)(mv), (*C.double)(&v[0]), (C.lapack_int)(ldv), (*C.double)(&stat[0])))
}

func Sgesvx(fact byte, trans blas.Transpose, n int, nrhs int, a []float32, lda int, af []float32, ldaf int, ipiv []int32, equed []byte, r []float32, c []float32, b []float32, ldb int, x []float32, ldx int, rcond []float32, ferr []float32, berr []float32, rpivot []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sgesvx((C.int)(rowMajor), (C.char)(fact), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.char)(unsafe.Pointer(&equed[0])), (*C.float)(&r[0]), (*C.float)(&c[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0]), (*C.float)(&rpivot[0])))
}

func Dgesvx(fact byte, trans blas.Transpose, n int, nrhs int, a []float64, lda int, af []float64, ldaf int, ipiv []int32, equed []byte, r []float64, c []float64, b []float64, ldb int, x []float64, ldx int, rcond []float64, ferr []float64, berr []float64, rpivot []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dgesvx((C.int)(rowMajor), (C.char)(fact), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.char)(unsafe.Pointer(&equed[0])), (*C.double)(&r[0]), (*C.double)(&c[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0]), (*C.double)(&rpivot[0])))
}

func Cgesvx(fact byte, trans blas.Transpose, n int, nrhs int, a []complex64, lda int, af []complex64, ldaf int, ipiv []int32, equed []byte, r []float32, c []float32, b []complex64, ldb int, x []complex64, ldx int, rcond []float32, ferr []float32, berr []float32, rpivot []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cgesvx((C.int)(rowMajor), (C.char)(fact), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.char)(unsafe.Pointer(&equed[0])), (*C.float)(&r[0]), (*C.float)(&c[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0]), (*C.float)(&rpivot[0])))
}

func Zgesvx(fact byte, trans blas.Transpose, n int, nrhs int, a []complex128, lda int, af []complex128, ldaf int, ipiv []int32, equed []byte, r []float64, c []float64, b []complex128, ldb int, x []complex128, ldx int, rcond []float64, ferr []float64, berr []float64, rpivot []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zgesvx((C.int)(rowMajor), (C.char)(fact), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.char)(unsafe.Pointer(&equed[0])), (*C.double)(&r[0]), (*C.double)(&c[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0]), (*C.double)(&rpivot[0])))
}

func Sgetf2(m int, n int, a []float32, lda int, ipiv []int32) bool {
	return isZero(C.LAPACKE_sgetf2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Dgetf2(m int, n int, a []float64, lda int, ipiv []int32) bool {
	return isZero(C.LAPACKE_dgetf2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Cgetf2(m int, n int, a []complex64, lda int, ipiv []int32) bool {
	return isZero(C.LAPACKE_cgetf2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Zgetf2(m int, n int, a []complex128, lda int, ipiv []int32) bool {
	return isZero(C.LAPACKE_zgetf2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Sgetrf(m int, n int, a []float32, lda int, ipiv []int32) bool {
	return isZero(C.LAPACKE_sgetrf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Dgetrf(m int, n int, a []float64, lda int, ipiv []int32) bool {
	return isZero(C.LAPACKE_dgetrf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Cgetrf(m int, n int, a []complex64, lda int, ipiv []int32) bool {
	return isZero(C.LAPACKE_cgetrf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Zgetrf(m int, n int, a []complex128, lda int, ipiv []int32) bool {
	return isZero(C.LAPACKE_zgetrf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Sgetri(n int, a []float32, lda int, ipiv []int32) bool {
	return isZero(C.LAPACKE_sgetri((C.int)(rowMajor), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Dgetri(n int, a []float64, lda int, ipiv []int32) bool {
	return isZero(C.LAPACKE_dgetri((C.int)(rowMajor), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Cgetri(n int, a []complex64, lda int, ipiv []int32) bool {
	return isZero(C.LAPACKE_cgetri((C.int)(rowMajor), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Zgetri(n int, a []complex128, lda int, ipiv []int32) bool {
	return isZero(C.LAPACKE_zgetri((C.int)(rowMajor), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Sgetrs(trans blas.Transpose, n int, nrhs int, a []float32, lda int, ipiv []int32, b []float32, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sgetrs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dgetrs(trans blas.Transpose, n int, nrhs int, a []float64, lda int, ipiv []int32, b []float64, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dgetrs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cgetrs(trans blas.Transpose, n int, nrhs int, a []complex64, lda int, ipiv []int32, b []complex64, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cgetrs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zgetrs(trans blas.Transpose, n int, nrhs int, a []complex128, lda int, ipiv []int32, b []complex128, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zgetrs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Sggbak(job lapack.Job, s blas.Side, n int, ilo int, ihi int, lscale []float32, rscale []float32, m int, v []float32, ldv int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	return isZero(C.LAPACKE_sggbak((C.int)(rowMajor), (C.char)(job), (C.char)(s), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.float)(&lscale[0]), (*C.float)(&rscale[0]), (C.lapack_int)(m), (*C.float)(&v[0]), (C.lapack_int)(ldv)))
}

func Dggbak(job lapack.Job, s blas.Side, n int, ilo int, ihi int, lscale []float64, rscale []float64, m int, v []float64, ldv int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	return isZero(C.LAPACKE_dggbak((C.int)(rowMajor), (C.char)(job), (C.char)(s), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.double)(&lscale[0]), (*C.double)(&rscale[0]), (C.lapack_int)(m), (*C.double)(&v[0]), (C.lapack_int)(ldv)))
}

func Cggbak(job lapack.Job, s blas.Side, n int, ilo int, ihi int, lscale []float32, rscale []float32, m int, v []complex64, ldv int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	return isZero(C.LAPACKE_cggbak((C.int)(rowMajor), (C.char)(job), (C.char)(s), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.float)(&lscale[0]), (*C.float)(&rscale[0]), (C.lapack_int)(m), (*C.lapack_complex_float)(&v[0]), (C.lapack_int)(ldv)))
}

func Zggbak(job lapack.Job, s blas.Side, n int, ilo int, ihi int, lscale []float64, rscale []float64, m int, v []complex128, ldv int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	return isZero(C.LAPACKE_zggbak((C.int)(rowMajor), (C.char)(job), (C.char)(s), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.double)(&lscale[0]), (*C.double)(&rscale[0]), (C.lapack_int)(m), (*C.lapack_complex_double)(&v[0]), (C.lapack_int)(ldv)))
}

func Sggbal(job lapack.Job, n int, a []float32, lda int, b []float32, ldb int, ilo []int32, ihi []int32, lscale []float32, rscale []float32) bool {
	return isZero(C.LAPACKE_sggbal((C.int)(rowMajor), (C.char)(job), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.float)(&lscale[0]), (*C.float)(&rscale[0])))
}

func Dggbal(job lapack.Job, n int, a []float64, lda int, b []float64, ldb int, ilo []int32, ihi []int32, lscale []float64, rscale []float64) bool {
	return isZero(C.LAPACKE_dggbal((C.int)(rowMajor), (C.char)(job), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.double)(&lscale[0]), (*C.double)(&rscale[0])))
}

func Cggbal(job lapack.Job, n int, a []complex64, lda int, b []complex64, ldb int, ilo []int32, ihi []int32, lscale []float32, rscale []float32) bool {
	return isZero(C.LAPACKE_cggbal((C.int)(rowMajor), (C.char)(job), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.float)(&lscale[0]), (*C.float)(&rscale[0])))
}

func Zggbal(job lapack.Job, n int, a []complex128, lda int, b []complex128, ldb int, ilo []int32, ihi []int32, lscale []float64, rscale []float64) bool {
	return isZero(C.LAPACKE_zggbal((C.int)(rowMajor), (C.char)(job), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.double)(&lscale[0]), (*C.double)(&rscale[0])))
}

func Sggev(jobvl lapack.Job, jobvr lapack.Job, n int, a []float32, lda int, b []float32, ldb int, alphar []float32, alphai []float32, beta []float32, vl []float32, ldvl int, vr []float32, ldvr int) bool {
	return isZero(C.LAPACKE_sggev((C.int)(rowMajor), (C.char)(jobvl), (C.char)(jobvr), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&alphar[0]), (*C.float)(&alphai[0]), (*C.float)(&beta[0]), (*C.float)(&vl[0]), (C.lapack_int)(ldvl), (*C.float)(&vr[0]), (C.lapack_int)(ldvr)))
}

func Dggev(jobvl lapack.Job, jobvr lapack.Job, n int, a []float64, lda int, b []float64, ldb int, alphar []float64, alphai []float64, beta []float64, vl []float64, ldvl int, vr []float64, ldvr int) bool {
	return isZero(C.LAPACKE_dggev((C.int)(rowMajor), (C.char)(jobvl), (C.char)(jobvr), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&alphar[0]), (*C.double)(&alphai[0]), (*C.double)(&beta[0]), (*C.double)(&vl[0]), (C.lapack_int)(ldvl), (*C.double)(&vr[0]), (C.lapack_int)(ldvr)))
}

func Cggev(jobvl lapack.Job, jobvr lapack.Job, n int, a []complex64, lda int, b []complex64, ldb int, alpha []complex64, beta []complex64, vl []complex64, ldvl int, vr []complex64, ldvr int) bool {
	return isZero(C.LAPACKE_cggev((C.int)(rowMajor), (C.char)(jobvl), (C.char)(jobvr), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&alpha[0]), (*C.lapack_complex_float)(&beta[0]), (*C.lapack_complex_float)(&vl[0]), (C.lapack_int)(ldvl), (*C.lapack_complex_float)(&vr[0]), (C.lapack_int)(ldvr)))
}

func Zggev(jobvl lapack.Job, jobvr lapack.Job, n int, a []complex128, lda int, b []complex128, ldb int, alpha []complex128, beta []complex128, vl []complex128, ldvl int, vr []complex128, ldvr int) bool {
	return isZero(C.LAPACKE_zggev((C.int)(rowMajor), (C.char)(jobvl), (C.char)(jobvr), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&alpha[0]), (*C.lapack_complex_double)(&beta[0]), (*C.lapack_complex_double)(&vl[0]), (C.lapack_int)(ldvl), (*C.lapack_complex_double)(&vr[0]), (C.lapack_int)(ldvr)))
}

func Sggevx(balanc byte, jobvl lapack.Job, jobvr lapack.Job, sense byte, n int, a []float32, lda int, b []float32, ldb int, alphar []float32, alphai []float32, beta []float32, vl []float32, ldvl int, vr []float32, ldvr int, ilo []int32, ihi []int32, lscale []float32, rscale []float32, abnrm []float32, bbnrm []float32, rconde []float32, rcondv []float32) bool {
	return isZero(C.LAPACKE_sggevx((C.int)(rowMajor), (C.char)(balanc), (C.char)(jobvl), (C.char)(jobvr), (C.char)(sense), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&alphar[0]), (*C.float)(&alphai[0]), (*C.float)(&beta[0]), (*C.float)(&vl[0]), (C.lapack_int)(ldvl), (*C.float)(&vr[0]), (C.lapack_int)(ldvr), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.float)(&lscale[0]), (*C.float)(&rscale[0]), (*C.float)(&abnrm[0]), (*C.float)(&bbnrm[0]), (*C.float)(&rconde[0]), (*C.float)(&rcondv[0])))
}

func Dggevx(balanc byte, jobvl lapack.Job, jobvr lapack.Job, sense byte, n int, a []float64, lda int, b []float64, ldb int, alphar []float64, alphai []float64, beta []float64, vl []float64, ldvl int, vr []float64, ldvr int, ilo []int32, ihi []int32, lscale []float64, rscale []float64, abnrm []float64, bbnrm []float64, rconde []float64, rcondv []float64) bool {
	return isZero(C.LAPACKE_dggevx((C.int)(rowMajor), (C.char)(balanc), (C.char)(jobvl), (C.char)(jobvr), (C.char)(sense), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&alphar[0]), (*C.double)(&alphai[0]), (*C.double)(&beta[0]), (*C.double)(&vl[0]), (C.lapack_int)(ldvl), (*C.double)(&vr[0]), (C.lapack_int)(ldvr), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.double)(&lscale[0]), (*C.double)(&rscale[0]), (*C.double)(&abnrm[0]), (*C.double)(&bbnrm[0]), (*C.double)(&rconde[0]), (*C.double)(&rcondv[0])))
}

func Cggevx(balanc byte, jobvl lapack.Job, jobvr lapack.Job, sense byte, n int, a []complex64, lda int, b []complex64, ldb int, alpha []complex64, beta []complex64, vl []complex64, ldvl int, vr []complex64, ldvr int, ilo []int32, ihi []int32, lscale []float32, rscale []float32, abnrm []float32, bbnrm []float32, rconde []float32, rcondv []float32) bool {
	return isZero(C.LAPACKE_cggevx((C.int)(rowMajor), (C.char)(balanc), (C.char)(jobvl), (C.char)(jobvr), (C.char)(sense), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&alpha[0]), (*C.lapack_complex_float)(&beta[0]), (*C.lapack_complex_float)(&vl[0]), (C.lapack_int)(ldvl), (*C.lapack_complex_float)(&vr[0]), (C.lapack_int)(ldvr), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.float)(&lscale[0]), (*C.float)(&rscale[0]), (*C.float)(&abnrm[0]), (*C.float)(&bbnrm[0]), (*C.float)(&rconde[0]), (*C.float)(&rcondv[0])))
}

func Zggevx(balanc byte, jobvl lapack.Job, jobvr lapack.Job, sense byte, n int, a []complex128, lda int, b []complex128, ldb int, alpha []complex128, beta []complex128, vl []complex128, ldvl int, vr []complex128, ldvr int, ilo []int32, ihi []int32, lscale []float64, rscale []float64, abnrm []float64, bbnrm []float64, rconde []float64, rcondv []float64) bool {
	return isZero(C.LAPACKE_zggevx((C.int)(rowMajor), (C.char)(balanc), (C.char)(jobvl), (C.char)(jobvr), (C.char)(sense), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&alpha[0]), (*C.lapack_complex_double)(&beta[0]), (*C.lapack_complex_double)(&vl[0]), (C.lapack_int)(ldvl), (*C.lapack_complex_double)(&vr[0]), (C.lapack_int)(ldvr), (*C.lapack_int)(&ilo[0]), (*C.lapack_int)(&ihi[0]), (*C.double)(&lscale[0]), (*C.double)(&rscale[0]), (*C.double)(&abnrm[0]), (*C.double)(&bbnrm[0]), (*C.double)(&rconde[0]), (*C.double)(&rcondv[0])))
}

func Sggglm(n int, m int, p int, a []float32, lda int, b []float32, ldb int, d []float32, x []float32, y []float32) bool {
	return isZero(C.LAPACKE_sggglm((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(m), (C.lapack_int)(p), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&d[0]), (*C.float)(&x[0]), (*C.float)(&y[0])))
}

func Dggglm(n int, m int, p int, a []float64, lda int, b []float64, ldb int, d []float64, x []float64, y []float64) bool {
	return isZero(C.LAPACKE_dggglm((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(m), (C.lapack_int)(p), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&d[0]), (*C.double)(&x[0]), (*C.double)(&y[0])))
}

func Cggglm(n int, m int, p int, a []complex64, lda int, b []complex64, ldb int, d []complex64, x []complex64, y []complex64) bool {
	return isZero(C.LAPACKE_cggglm((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(m), (C.lapack_int)(p), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&d[0]), (*C.lapack_complex_float)(&x[0]), (*C.lapack_complex_float)(&y[0])))
}

func Zggglm(n int, m int, p int, a []complex128, lda int, b []complex128, ldb int, d []complex128, x []complex128, y []complex128) bool {
	return isZero(C.LAPACKE_zggglm((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(m), (C.lapack_int)(p), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&d[0]), (*C.lapack_complex_double)(&x[0]), (*C.lapack_complex_double)(&y[0])))
}

func Sgghrd(compq lapack.CompSV, compz lapack.CompSV, n int, ilo int, ihi int, a []float32, lda int, b []float32, ldb int, q []float32, ldq int, z []float32, ldz int) bool {
	return isZero(C.LAPACKE_sgghrd((C.int)(rowMajor), (C.char)(compq), (C.char)(compz), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&q[0]), (C.lapack_int)(ldq), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dgghrd(compq lapack.CompSV, compz lapack.CompSV, n int, ilo int, ihi int, a []float64, lda int, b []float64, ldb int, q []float64, ldq int, z []float64, ldz int) bool {
	return isZero(C.LAPACKE_dgghrd((C.int)(rowMajor), (C.char)(compq), (C.char)(compz), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&q[0]), (C.lapack_int)(ldq), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Cgghrd(compq lapack.CompSV, compz lapack.CompSV, n int, ilo int, ihi int, a []complex64, lda int, b []complex64, ldb int, q []complex64, ldq int, z []complex64, ldz int) bool {
	return isZero(C.LAPACKE_cgghrd((C.int)(rowMajor), (C.char)(compq), (C.char)(compz), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz)))
}

func Zgghrd(compq lapack.CompSV, compz lapack.CompSV, n int, ilo int, ihi int, a []complex128, lda int, b []complex128, ldb int, q []complex128, ldq int, z []complex128, ldz int) bool {
	return isZero(C.LAPACKE_zgghrd((C.int)(rowMajor), (C.char)(compq), (C.char)(compz), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz)))
}

func Sgglse(m int, n int, p int, a []float32, lda int, b []float32, ldb int, c []float32, d []float32, x []float32) bool {
	return isZero(C.LAPACKE_sgglse((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(p), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&c[0]), (*C.float)(&d[0]), (*C.float)(&x[0])))
}

func Dgglse(m int, n int, p int, a []float64, lda int, b []float64, ldb int, c []float64, d []float64, x []float64) bool {
	return isZero(C.LAPACKE_dgglse((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(p), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&c[0]), (*C.double)(&d[0]), (*C.double)(&x[0])))
}

func Cgglse(m int, n int, p int, a []complex64, lda int, b []complex64, ldb int, c []complex64, d []complex64, x []complex64) bool {
	return isZero(C.LAPACKE_cgglse((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(p), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&c[0]), (*C.lapack_complex_float)(&d[0]), (*C.lapack_complex_float)(&x[0])))
}

func Zgglse(m int, n int, p int, a []complex128, lda int, b []complex128, ldb int, c []complex128, d []complex128, x []complex128) bool {
	return isZero(C.LAPACKE_zgglse((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(p), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&c[0]), (*C.lapack_complex_double)(&d[0]), (*C.lapack_complex_double)(&x[0])))
}

func Sggqrf(n int, m int, p int, a []float32, lda int, taua []float32, b []float32, ldb int, taub []float32) bool {
	return isZero(C.LAPACKE_sggqrf((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(m), (C.lapack_int)(p), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&taua[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&taub[0])))
}

func Dggqrf(n int, m int, p int, a []float64, lda int, taua []float64, b []float64, ldb int, taub []float64) bool {
	return isZero(C.LAPACKE_dggqrf((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(m), (C.lapack_int)(p), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&taua[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&taub[0])))
}

func Cggqrf(n int, m int, p int, a []complex64, lda int, taua []complex64, b []complex64, ldb int, taub []complex64) bool {
	return isZero(C.LAPACKE_cggqrf((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(m), (C.lapack_int)(p), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&taua[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&taub[0])))
}

func Zggqrf(n int, m int, p int, a []complex128, lda int, taua []complex128, b []complex128, ldb int, taub []complex128) bool {
	return isZero(C.LAPACKE_zggqrf((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(m), (C.lapack_int)(p), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&taua[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&taub[0])))
}

func Sggrqf(m int, p int, n int, a []float32, lda int, taua []float32, b []float32, ldb int, taub []float32) bool {
	return isZero(C.LAPACKE_sggrqf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&taua[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&taub[0])))
}

func Dggrqf(m int, p int, n int, a []float64, lda int, taua []float64, b []float64, ldb int, taub []float64) bool {
	return isZero(C.LAPACKE_dggrqf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&taua[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&taub[0])))
}

func Cggrqf(m int, p int, n int, a []complex64, lda int, taua []complex64, b []complex64, ldb int, taub []complex64) bool {
	return isZero(C.LAPACKE_cggrqf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&taua[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&taub[0])))
}

func Zggrqf(m int, p int, n int, a []complex128, lda int, taua []complex128, b []complex128, ldb int, taub []complex128) bool {
	return isZero(C.LAPACKE_zggrqf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&taua[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&taub[0])))
}

func Sggsvd(jobu lapack.Job, jobv lapack.Job, jobq lapack.Job, m int, n int, p int, k []int32, l []int32, a []float32, lda int, b []float32, ldb int, alpha []float32, beta []float32, u []float32, ldu int, v []float32, ldv int, q []float32, ldq int, iwork []int32) bool {
	return isZero(C.LAPACKE_sggsvd((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobv), (C.char)(jobq), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(p), (*C.lapack_int)(&k[0]), (*C.lapack_int)(&l[0]), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&alpha[0]), (*C.float)(&beta[0]), (*C.float)(&u[0]), (C.lapack_int)(ldu), (*C.float)(&v[0]), (C.lapack_int)(ldv), (*C.float)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_int)(&iwork[0])))
}

func Dggsvd(jobu lapack.Job, jobv lapack.Job, jobq lapack.Job, m int, n int, p int, k []int32, l []int32, a []float64, lda int, b []float64, ldb int, alpha []float64, beta []float64, u []float64, ldu int, v []float64, ldv int, q []float64, ldq int, iwork []int32) bool {
	return isZero(C.LAPACKE_dggsvd((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobv), (C.char)(jobq), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(p), (*C.lapack_int)(&k[0]), (*C.lapack_int)(&l[0]), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&alpha[0]), (*C.double)(&beta[0]), (*C.double)(&u[0]), (C.lapack_int)(ldu), (*C.double)(&v[0]), (C.lapack_int)(ldv), (*C.double)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_int)(&iwork[0])))
}

func Cggsvd(jobu lapack.Job, jobv lapack.Job, jobq lapack.Job, m int, n int, p int, k []int32, l []int32, a []complex64, lda int, b []complex64, ldb int, alpha []float32, beta []float32, u []complex64, ldu int, v []complex64, ldv int, q []complex64, ldq int, iwork []int32) bool {
	return isZero(C.LAPACKE_cggsvd((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobv), (C.char)(jobq), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(p), (*C.lapack_int)(&k[0]), (*C.lapack_int)(&l[0]), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&alpha[0]), (*C.float)(&beta[0]), (*C.lapack_complex_float)(&u[0]), (C.lapack_int)(ldu), (*C.lapack_complex_float)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_float)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_int)(&iwork[0])))
}

func Zggsvd(jobu lapack.Job, jobv lapack.Job, jobq lapack.Job, m int, n int, p int, k []int32, l []int32, a []complex128, lda int, b []complex128, ldb int, alpha []float64, beta []float64, u []complex128, ldu int, v []complex128, ldv int, q []complex128, ldq int, iwork []int32) bool {
	return isZero(C.LAPACKE_zggsvd((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobv), (C.char)(jobq), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(p), (*C.lapack_int)(&k[0]), (*C.lapack_int)(&l[0]), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&alpha[0]), (*C.double)(&beta[0]), (*C.lapack_complex_double)(&u[0]), (C.lapack_int)(ldu), (*C.lapack_complex_double)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_double)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_int)(&iwork[0])))
}

func Sggsvp(jobu lapack.Job, jobv lapack.Job, jobq lapack.Job, m int, p int, n int, a []float32, lda int, b []float32, ldb int, tola float32, tolb float32, k []int32, l []int32, u []float32, ldu int, v []float32, ldv int, q []float32, ldq int) bool {
	return isZero(C.LAPACKE_sggsvp((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobv), (C.char)(jobq), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (C.float)(tola), (C.float)(tolb), (*C.lapack_int)(&k[0]), (*C.lapack_int)(&l[0]), (*C.float)(&u[0]), (C.lapack_int)(ldu), (*C.float)(&v[0]), (C.lapack_int)(ldv), (*C.float)(&q[0]), (C.lapack_int)(ldq)))
}

func Dggsvp(jobu lapack.Job, jobv lapack.Job, jobq lapack.Job, m int, p int, n int, a []float64, lda int, b []float64, ldb int, tola float64, tolb float64, k []int32, l []int32, u []float64, ldu int, v []float64, ldv int, q []float64, ldq int) bool {
	return isZero(C.LAPACKE_dggsvp((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobv), (C.char)(jobq), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (C.double)(tola), (C.double)(tolb), (*C.lapack_int)(&k[0]), (*C.lapack_int)(&l[0]), (*C.double)(&u[0]), (C.lapack_int)(ldu), (*C.double)(&v[0]), (C.lapack_int)(ldv), (*C.double)(&q[0]), (C.lapack_int)(ldq)))
}

func Cggsvp(jobu lapack.Job, jobv lapack.Job, jobq lapack.Job, m int, p int, n int, a []complex64, lda int, b []complex64, ldb int, tola float32, tolb float32, k []int32, l []int32, u []complex64, ldu int, v []complex64, ldv int, q []complex64, ldq int) bool {
	return isZero(C.LAPACKE_cggsvp((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobv), (C.char)(jobq), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (C.float)(tola), (C.float)(tolb), (*C.lapack_int)(&k[0]), (*C.lapack_int)(&l[0]), (*C.lapack_complex_float)(&u[0]), (C.lapack_int)(ldu), (*C.lapack_complex_float)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_float)(&q[0]), (C.lapack_int)(ldq)))
}

func Zggsvp(jobu lapack.Job, jobv lapack.Job, jobq lapack.Job, m int, p int, n int, a []complex128, lda int, b []complex128, ldb int, tola float64, tolb float64, k []int32, l []int32, u []complex128, ldu int, v []complex128, ldv int, q []complex128, ldq int) bool {
	return isZero(C.LAPACKE_zggsvp((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobv), (C.char)(jobq), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (C.double)(tola), (C.double)(tolb), (*C.lapack_int)(&k[0]), (*C.lapack_int)(&l[0]), (*C.lapack_complex_double)(&u[0]), (C.lapack_int)(ldu), (*C.lapack_complex_double)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_double)(&q[0]), (C.lapack_int)(ldq)))
}

func Sgtcon(norm byte, n int, dl []float32, d []float32, du []float32, du2 []float32, ipiv []int32, anorm float32, rcond []float32) bool {
	return isZero(C.LAPACKE_sgtcon((C.char)(norm), (C.lapack_int)(n), (*C.float)(&dl[0]), (*C.float)(&d[0]), (*C.float)(&du[0]), (*C.float)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Dgtcon(norm byte, n int, dl []float64, d []float64, du []float64, du2 []float64, ipiv []int32, anorm float64, rcond []float64) bool {
	return isZero(C.LAPACKE_dgtcon((C.char)(norm), (C.lapack_int)(n), (*C.double)(&dl[0]), (*C.double)(&d[0]), (*C.double)(&du[0]), (*C.double)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Cgtcon(norm byte, n int, dl []complex64, d []complex64, du []complex64, du2 []complex64, ipiv []int32, anorm float32, rcond []float32) bool {
	return isZero(C.LAPACKE_cgtcon((C.char)(norm), (C.lapack_int)(n), (*C.lapack_complex_float)(&dl[0]), (*C.lapack_complex_float)(&d[0]), (*C.lapack_complex_float)(&du[0]), (*C.lapack_complex_float)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Zgtcon(norm byte, n int, dl []complex128, d []complex128, du []complex128, du2 []complex128, ipiv []int32, anorm float64, rcond []float64) bool {
	return isZero(C.LAPACKE_zgtcon((C.char)(norm), (C.lapack_int)(n), (*C.lapack_complex_double)(&dl[0]), (*C.lapack_complex_double)(&d[0]), (*C.lapack_complex_double)(&du[0]), (*C.lapack_complex_double)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Sgtrfs(trans blas.Transpose, n int, nrhs int, dl []float32, d []float32, du []float32, dlf []float32, df []float32, duf []float32, du2 []float32, ipiv []int32, b []float32, ldb int, x []float32, ldx int, ferr []float32, berr []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sgtrfs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&dl[0]), (*C.float)(&d[0]), (*C.float)(&du[0]), (*C.float)(&dlf[0]), (*C.float)(&df[0]), (*C.float)(&duf[0]), (*C.float)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dgtrfs(trans blas.Transpose, n int, nrhs int, dl []float64, d []float64, du []float64, dlf []float64, df []float64, duf []float64, du2 []float64, ipiv []int32, b []float64, ldb int, x []float64, ldx int, ferr []float64, berr []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dgtrfs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&dl[0]), (*C.double)(&d[0]), (*C.double)(&du[0]), (*C.double)(&dlf[0]), (*C.double)(&df[0]), (*C.double)(&duf[0]), (*C.double)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Cgtrfs(trans blas.Transpose, n int, nrhs int, dl []complex64, d []complex64, du []complex64, dlf []complex64, df []complex64, duf []complex64, du2 []complex64, ipiv []int32, b []complex64, ldb int, x []complex64, ldx int, ferr []float32, berr []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cgtrfs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&dl[0]), (*C.lapack_complex_float)(&d[0]), (*C.lapack_complex_float)(&du[0]), (*C.lapack_complex_float)(&dlf[0]), (*C.lapack_complex_float)(&df[0]), (*C.lapack_complex_float)(&duf[0]), (*C.lapack_complex_float)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zgtrfs(trans blas.Transpose, n int, nrhs int, dl []complex128, d []complex128, du []complex128, dlf []complex128, df []complex128, duf []complex128, du2 []complex128, ipiv []int32, b []complex128, ldb int, x []complex128, ldx int, ferr []float64, berr []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zgtrfs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&dl[0]), (*C.lapack_complex_double)(&d[0]), (*C.lapack_complex_double)(&du[0]), (*C.lapack_complex_double)(&dlf[0]), (*C.lapack_complex_double)(&df[0]), (*C.lapack_complex_double)(&duf[0]), (*C.lapack_complex_double)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Sgtsv(n int, nrhs int, dl []float32, d []float32, du []float32, b []float32, ldb int) bool {
	return isZero(C.LAPACKE_sgtsv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&dl[0]), (*C.float)(&d[0]), (*C.float)(&du[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dgtsv(n int, nrhs int, dl []float64, d []float64, du []float64, b []float64, ldb int) bool {
	return isZero(C.LAPACKE_dgtsv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&dl[0]), (*C.double)(&d[0]), (*C.double)(&du[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cgtsv(n int, nrhs int, dl []complex64, d []complex64, du []complex64, b []complex64, ldb int) bool {
	return isZero(C.LAPACKE_cgtsv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&dl[0]), (*C.lapack_complex_float)(&d[0]), (*C.lapack_complex_float)(&du[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zgtsv(n int, nrhs int, dl []complex128, d []complex128, du []complex128, b []complex128, ldb int) bool {
	return isZero(C.LAPACKE_zgtsv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&dl[0]), (*C.lapack_complex_double)(&d[0]), (*C.lapack_complex_double)(&du[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Sgtsvx(fact byte, trans blas.Transpose, n int, nrhs int, dl []float32, d []float32, du []float32, dlf []float32, df []float32, duf []float32, du2 []float32, ipiv []int32, b []float32, ldb int, x []float32, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sgtsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&dl[0]), (*C.float)(&d[0]), (*C.float)(&du[0]), (*C.float)(&dlf[0]), (*C.float)(&df[0]), (*C.float)(&duf[0]), (*C.float)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dgtsvx(fact byte, trans blas.Transpose, n int, nrhs int, dl []float64, d []float64, du []float64, dlf []float64, df []float64, duf []float64, du2 []float64, ipiv []int32, b []float64, ldb int, x []float64, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dgtsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&dl[0]), (*C.double)(&d[0]), (*C.double)(&du[0]), (*C.double)(&dlf[0]), (*C.double)(&df[0]), (*C.double)(&duf[0]), (*C.double)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Cgtsvx(fact byte, trans blas.Transpose, n int, nrhs int, dl []complex64, d []complex64, du []complex64, dlf []complex64, df []complex64, duf []complex64, du2 []complex64, ipiv []int32, b []complex64, ldb int, x []complex64, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cgtsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&dl[0]), (*C.lapack_complex_float)(&d[0]), (*C.lapack_complex_float)(&du[0]), (*C.lapack_complex_float)(&dlf[0]), (*C.lapack_complex_float)(&df[0]), (*C.lapack_complex_float)(&duf[0]), (*C.lapack_complex_float)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zgtsvx(fact byte, trans blas.Transpose, n int, nrhs int, dl []complex128, d []complex128, du []complex128, dlf []complex128, df []complex128, duf []complex128, du2 []complex128, ipiv []int32, b []complex128, ldb int, x []complex128, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zgtsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&dl[0]), (*C.lapack_complex_double)(&d[0]), (*C.lapack_complex_double)(&du[0]), (*C.lapack_complex_double)(&dlf[0]), (*C.lapack_complex_double)(&df[0]), (*C.lapack_complex_double)(&duf[0]), (*C.lapack_complex_double)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Sgttrf(n int, dl []float32, d []float32, du []float32, du2 []float32, ipiv []int32) bool {
	return isZero(C.LAPACKE_sgttrf((C.lapack_int)(n), (*C.float)(&dl[0]), (*C.float)(&d[0]), (*C.float)(&du[0]), (*C.float)(&du2[0]), (*C.lapack_int)(&ipiv[0])))
}

func Dgttrf(n int, dl []float64, d []float64, du []float64, du2 []float64, ipiv []int32) bool {
	return isZero(C.LAPACKE_dgttrf((C.lapack_int)(n), (*C.double)(&dl[0]), (*C.double)(&d[0]), (*C.double)(&du[0]), (*C.double)(&du2[0]), (*C.lapack_int)(&ipiv[0])))
}

func Cgttrf(n int, dl []complex64, d []complex64, du []complex64, du2 []complex64, ipiv []int32) bool {
	return isZero(C.LAPACKE_cgttrf((C.lapack_int)(n), (*C.lapack_complex_float)(&dl[0]), (*C.lapack_complex_float)(&d[0]), (*C.lapack_complex_float)(&du[0]), (*C.lapack_complex_float)(&du2[0]), (*C.lapack_int)(&ipiv[0])))
}

func Zgttrf(n int, dl []complex128, d []complex128, du []complex128, du2 []complex128, ipiv []int32) bool {
	return isZero(C.LAPACKE_zgttrf((C.lapack_int)(n), (*C.lapack_complex_double)(&dl[0]), (*C.lapack_complex_double)(&d[0]), (*C.lapack_complex_double)(&du[0]), (*C.lapack_complex_double)(&du2[0]), (*C.lapack_int)(&ipiv[0])))
}

func Sgttrs(trans blas.Transpose, n int, nrhs int, dl []float32, d []float32, du []float32, du2 []float32, ipiv []int32, b []float32, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sgttrs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&dl[0]), (*C.float)(&d[0]), (*C.float)(&du[0]), (*C.float)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dgttrs(trans blas.Transpose, n int, nrhs int, dl []float64, d []float64, du []float64, du2 []float64, ipiv []int32, b []float64, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dgttrs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&dl[0]), (*C.double)(&d[0]), (*C.double)(&du[0]), (*C.double)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cgttrs(trans blas.Transpose, n int, nrhs int, dl []complex64, d []complex64, du []complex64, du2 []complex64, ipiv []int32, b []complex64, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cgttrs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&dl[0]), (*C.lapack_complex_float)(&d[0]), (*C.lapack_complex_float)(&du[0]), (*C.lapack_complex_float)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zgttrs(trans blas.Transpose, n int, nrhs int, dl []complex128, d []complex128, du []complex128, du2 []complex128, ipiv []int32, b []complex128, ldb int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zgttrs((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&dl[0]), (*C.lapack_complex_double)(&d[0]), (*C.lapack_complex_double)(&du[0]), (*C.lapack_complex_double)(&du2[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Chbev(jobz lapack.Job, ul blas.Uplo, n int, kd int, ab []complex64, ldab int, w []float32, z []complex64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chbev((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz)))
}

func Zhbev(jobz lapack.Job, ul blas.Uplo, n int, kd int, ab []complex128, ldab int, w []float64, z []complex128, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhbev((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz)))
}

func Chbevd(jobz lapack.Job, ul blas.Uplo, n int, kd int, ab []complex64, ldab int, w []float32, z []complex64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chbevd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz)))
}

func Zhbevd(jobz lapack.Job, ul blas.Uplo, n int, kd int, ab []complex128, ldab int, w []float64, z []complex128, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhbevd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz)))
}

func Chbevx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, kd int, ab []complex64, ldab int, q []complex64, ldq int, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []complex64, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chbevx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_float)(&q[0]), (C.lapack_int)(ldq), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Zhbevx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, kd int, ab []complex128, ldab int, q []complex128, ldq int, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []complex128, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhbevx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_double)(&q[0]), (C.lapack_int)(ldq), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Chbgst(vect byte, ul blas.Uplo, n int, ka int, kb int, ab []complex64, ldab int, bb []complex64, ldbb int, x []complex64, ldx int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chbgst((C.int)(rowMajor), (C.char)(vect), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_float)(&bb[0]), (C.lapack_int)(ldbb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx)))
}

func Zhbgst(vect byte, ul blas.Uplo, n int, ka int, kb int, ab []complex128, ldab int, bb []complex128, ldbb int, x []complex128, ldx int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhbgst((C.int)(rowMajor), (C.char)(vect), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_double)(&bb[0]), (C.lapack_int)(ldbb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx)))
}

func Chbgv(jobz lapack.Job, ul blas.Uplo, n int, ka int, kb int, ab []complex64, ldab int, bb []complex64, ldbb int, w []float32, z []complex64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chbgv((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_float)(&bb[0]), (C.lapack_int)(ldbb), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz)))
}

func Zhbgv(jobz lapack.Job, ul blas.Uplo, n int, ka int, kb int, ab []complex128, ldab int, bb []complex128, ldbb int, w []float64, z []complex128, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhbgv((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_double)(&bb[0]), (C.lapack_int)(ldbb), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz)))
}

func Chbgvd(jobz lapack.Job, ul blas.Uplo, n int, ka int, kb int, ab []complex64, ldab int, bb []complex64, ldbb int, w []float32, z []complex64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chbgvd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_float)(&bb[0]), (C.lapack_int)(ldbb), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz)))
}

func Zhbgvd(jobz lapack.Job, ul blas.Uplo, n int, ka int, kb int, ab []complex128, ldab int, bb []complex128, ldbb int, w []float64, z []complex128, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhbgvd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_double)(&bb[0]), (C.lapack_int)(ldbb), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz)))
}

func Chbgvx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, ka int, kb int, ab []complex64, ldab int, bb []complex64, ldbb int, q []complex64, ldq int, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []complex64, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chbgvx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_float)(&bb[0]), (C.lapack_int)(ldbb), (*C.lapack_complex_float)(&q[0]), (C.lapack_int)(ldq), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Zhbgvx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, ka int, kb int, ab []complex128, ldab int, bb []complex128, ldbb int, q []complex128, ldq int, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []complex128, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhbgvx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_double)(&bb[0]), (C.lapack_int)(ldbb), (*C.lapack_complex_double)(&q[0]), (C.lapack_int)(ldq), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Chbtrd(vect byte, ul blas.Uplo, n int, kd int, ab []complex64, ldab int, d []float32, e []float32, q []complex64, ldq int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chbtrd((C.int)(rowMajor), (C.char)(vect), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.lapack_complex_float)(&q[0]), (C.lapack_int)(ldq)))
}

func Zhbtrd(vect byte, ul blas.Uplo, n int, kd int, ab []complex128, ldab int, d []float64, e []float64, q []complex128, ldq int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhbtrd((C.int)(rowMajor), (C.char)(vect), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.lapack_complex_double)(&q[0]), (C.lapack_int)(ldq)))
}

func Checon(ul blas.Uplo, n int, a []complex64, lda int, ipiv []int32, anorm float32, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_checon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Zhecon(ul blas.Uplo, n int, a []complex128, lda int, ipiv []int32, anorm float64, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhecon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Cheequb(ul blas.Uplo, n int, a []complex64, lda int, s []float32, scond []float32, amax []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cheequb((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&s[0]), (*C.float)(&scond[0]), (*C.float)(&amax[0])))
}

func Zheequb(ul blas.Uplo, n int, a []complex128, lda int, s []float64, scond []float64, amax []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zheequb((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&s[0]), (*C.double)(&scond[0]), (*C.double)(&amax[0])))
}

func Cheev(jobz lapack.Job, ul blas.Uplo, n int, a []complex64, lda int, w []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cheev((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&w[0])))
}

func Zheev(jobz lapack.Job, ul blas.Uplo, n int, a []complex128, lda int, w []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zheev((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&w[0])))
}

func Cheevd(jobz lapack.Job, ul blas.Uplo, n int, a []complex64, lda int, w []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cheevd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&w[0])))
}

func Zheevd(jobz lapack.Job, ul blas.Uplo, n int, a []complex128, lda int, w []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zheevd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&w[0])))
}

func Cheevr(jobz lapack.Job, rng byte, ul blas.Uplo, n int, a []complex64, lda int, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []complex64, ldz int, isuppz []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cheevr((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&isuppz[0])))
}

func Zheevr(jobz lapack.Job, rng byte, ul blas.Uplo, n int, a []complex128, lda int, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []complex128, ldz int, isuppz []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zheevr((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&isuppz[0])))
}

func Cheevx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, a []complex64, lda int, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []complex64, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cheevx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Zheevx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, a []complex128, lda int, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []complex128, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zheevx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Chegst(itype int, ul blas.Uplo, n int, a []complex64, lda int, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chegst((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zhegst(itype int, ul blas.Uplo, n int, a []complex128, lda int, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhegst((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Chegv(itype int, jobz lapack.Job, ul blas.Uplo, n int, a []complex64, lda int, b []complex64, ldb int, w []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chegv((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&w[0])))
}

func Zhegv(itype int, jobz lapack.Job, ul blas.Uplo, n int, a []complex128, lda int, b []complex128, ldb int, w []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhegv((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&w[0])))
}

func Chegvd(itype int, jobz lapack.Job, ul blas.Uplo, n int, a []complex64, lda int, b []complex64, ldb int, w []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chegvd((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&w[0])))
}

func Zhegvd(itype int, jobz lapack.Job, ul blas.Uplo, n int, a []complex128, lda int, b []complex128, ldb int, w []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhegvd((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&w[0])))
}

func Chegvx(itype int, jobz lapack.Job, rng byte, ul blas.Uplo, n int, a []complex64, lda int, b []complex64, ldb int, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []complex64, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chegvx((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Zhegvx(itype int, jobz lapack.Job, rng byte, ul blas.Uplo, n int, a []complex128, lda int, b []complex128, ldb int, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []complex128, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhegvx((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Cherfs(ul blas.Uplo, n int, nrhs int, a []complex64, lda int, af []complex64, ldaf int, ipiv []int32, b []complex64, ldb int, x []complex64, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cherfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zherfs(ul blas.Uplo, n int, nrhs int, a []complex128, lda int, af []complex128, ldaf int, ipiv []int32, b []complex128, ldb int, x []complex128, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zherfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Chesv(ul blas.Uplo, n int, nrhs int, a []complex64, lda int, ipiv []int32, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chesv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zhesv(ul blas.Uplo, n int, nrhs int, a []complex128, lda int, ipiv []int32, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhesv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Chesvx(fact byte, ul blas.Uplo, n int, nrhs int, a []complex64, lda int, af []complex64, ldaf int, ipiv []int32, b []complex64, ldb int, x []complex64, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chesvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zhesvx(fact byte, ul blas.Uplo, n int, nrhs int, a []complex128, lda int, af []complex128, ldaf int, ipiv []int32, b []complex128, ldb int, x []complex128, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhesvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Chetrd(ul blas.Uplo, n int, a []complex64, lda int, d []float32, e []float32, tau []complex64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chetrd((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.lapack_complex_float)(&tau[0])))
}

func Zhetrd(ul blas.Uplo, n int, a []complex128, lda int, d []float64, e []float64, tau []complex128) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhetrd((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.lapack_complex_double)(&tau[0])))
}

func Chetrf(ul blas.Uplo, n int, a []complex64, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chetrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Zhetrf(ul blas.Uplo, n int, a []complex128, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhetrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Chetri(ul blas.Uplo, n int, a []complex64, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chetri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Zhetri(ul blas.Uplo, n int, a []complex128, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhetri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Chetrs(ul blas.Uplo, n int, nrhs int, a []complex64, lda int, ipiv []int32, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chetrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zhetrs(ul blas.Uplo, n int, nrhs int, a []complex128, lda int, ipiv []int32, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhetrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Chfrk(transr blas.Transpose, ul blas.Uplo, trans blas.Transpose, n int, k int, alpha float32, a []complex64, lda int, beta float32, c []complex64) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_chfrk((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(k), (C.float)(alpha), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (C.float)(beta), (*C.lapack_complex_float)(&c[0])))
}

func Zhfrk(transr blas.Transpose, ul blas.Uplo, trans blas.Transpose, n int, k int, alpha float64, a []complex128, lda int, beta float64, c []complex128) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zhfrk((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(k), (C.double)(alpha), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (C.double)(beta), (*C.lapack_complex_double)(&c[0])))
}

func Shgeqz(job lapack.Job, compq lapack.CompSV, compz lapack.CompSV, n int, ilo int, ihi int, h []float32, ldh int, t []float32, ldt int, alphar []float32, alphai []float32, beta []float32, q []float32, ldq int, z []float32, ldz int) bool {
	return isZero(C.LAPACKE_shgeqz((C.int)(rowMajor), (C.char)(job), (C.char)(compq), (C.char)(compz), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.float)(&h[0]), (C.lapack_int)(ldh), (*C.float)(&t[0]), (C.lapack_int)(ldt), (*C.float)(&alphar[0]), (*C.float)(&alphai[0]), (*C.float)(&beta[0]), (*C.float)(&q[0]), (C.lapack_int)(ldq), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dhgeqz(job lapack.Job, compq lapack.CompSV, compz lapack.CompSV, n int, ilo int, ihi int, h []float64, ldh int, t []float64, ldt int, alphar []float64, alphai []float64, beta []float64, q []float64, ldq int, z []float64, ldz int) bool {
	return isZero(C.LAPACKE_dhgeqz((C.int)(rowMajor), (C.char)(job), (C.char)(compq), (C.char)(compz), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.double)(&h[0]), (C.lapack_int)(ldh), (*C.double)(&t[0]), (C.lapack_int)(ldt), (*C.double)(&alphar[0]), (*C.double)(&alphai[0]), (*C.double)(&beta[0]), (*C.double)(&q[0]), (C.lapack_int)(ldq), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Chgeqz(job lapack.Job, compq lapack.CompSV, compz lapack.CompSV, n int, ilo int, ihi int, h []complex64, ldh int, t []complex64, ldt int, alpha []complex64, beta []complex64, q []complex64, ldq int, z []complex64, ldz int) bool {
	return isZero(C.LAPACKE_chgeqz((C.int)(rowMajor), (C.char)(job), (C.char)(compq), (C.char)(compz), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.lapack_complex_float)(&h[0]), (C.lapack_int)(ldh), (*C.lapack_complex_float)(&t[0]), (C.lapack_int)(ldt), (*C.lapack_complex_float)(&alpha[0]), (*C.lapack_complex_float)(&beta[0]), (*C.lapack_complex_float)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz)))
}

func Zhgeqz(job lapack.Job, compq lapack.CompSV, compz lapack.CompSV, n int, ilo int, ihi int, h []complex128, ldh int, t []complex128, ldt int, alpha []complex128, beta []complex128, q []complex128, ldq int, z []complex128, ldz int) bool {
	return isZero(C.LAPACKE_zhgeqz((C.int)(rowMajor), (C.char)(job), (C.char)(compq), (C.char)(compz), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.lapack_complex_double)(&h[0]), (C.lapack_int)(ldh), (*C.lapack_complex_double)(&t[0]), (C.lapack_int)(ldt), (*C.lapack_complex_double)(&alpha[0]), (*C.lapack_complex_double)(&beta[0]), (*C.lapack_complex_double)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz)))
}

func Chpcon(ul blas.Uplo, n int, ap []complex64, ipiv []int32, anorm float32, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chpcon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Zhpcon(ul blas.Uplo, n int, ap []complex128, ipiv []int32, anorm float64, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhpcon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Chpev(jobz lapack.Job, ul blas.Uplo, n int, ap []complex64, w []float32, z []complex64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chpev((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz)))
}

func Zhpev(jobz lapack.Job, ul blas.Uplo, n int, ap []complex128, w []float64, z []complex128, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhpev((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz)))
}

func Chpevd(jobz lapack.Job, ul blas.Uplo, n int, ap []complex64, w []float32, z []complex64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chpevd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz)))
}

func Zhpevd(jobz lapack.Job, ul blas.Uplo, n int, ap []complex128, w []float64, z []complex128, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhpevd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz)))
}

func Chpevx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, ap []complex64, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []complex64, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chpevx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Zhpevx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, ap []complex128, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []complex128, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhpevx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Chpgst(itype int, ul blas.Uplo, n int, ap []complex64, bp []complex64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chpgst((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&bp[0])))
}

func Zhpgst(itype int, ul blas.Uplo, n int, ap []complex128, bp []complex128) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhpgst((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&bp[0])))
}

func Chpgv(itype int, jobz lapack.Job, ul blas.Uplo, n int, ap []complex64, bp []complex64, w []float32, z []complex64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chpgv((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&bp[0]), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz)))
}

func Zhpgv(itype int, jobz lapack.Job, ul blas.Uplo, n int, ap []complex128, bp []complex128, w []float64, z []complex128, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhpgv((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&bp[0]), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz)))
}

func Chpgvd(itype int, jobz lapack.Job, ul blas.Uplo, n int, ap []complex64, bp []complex64, w []float32, z []complex64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chpgvd((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&bp[0]), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz)))
}

func Zhpgvd(itype int, jobz lapack.Job, ul blas.Uplo, n int, ap []complex128, bp []complex128, w []float64, z []complex128, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhpgvd((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&bp[0]), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz)))
}

func Chpgvx(itype int, jobz lapack.Job, rng byte, ul blas.Uplo, n int, ap []complex64, bp []complex64, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []complex64, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chpgvx((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&bp[0]), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Zhpgvx(itype int, jobz lapack.Job, rng byte, ul blas.Uplo, n int, ap []complex128, bp []complex128, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []complex128, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhpgvx((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&bp[0]), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Chprfs(ul blas.Uplo, n int, nrhs int, ap []complex64, afp []complex64, ipiv []int32, b []complex64, ldb int, x []complex64, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chprfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&afp[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zhprfs(ul blas.Uplo, n int, nrhs int, ap []complex128, afp []complex128, ipiv []int32, b []complex128, ldb int, x []complex128, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhprfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&afp[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Chpsv(ul blas.Uplo, n int, nrhs int, ap []complex64, ipiv []int32, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chpsv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zhpsv(ul blas.Uplo, n int, nrhs int, ap []complex128, ipiv []int32, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhpsv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Chpsvx(fact byte, ul blas.Uplo, n int, nrhs int, ap []complex64, afp []complex64, ipiv []int32, b []complex64, ldb int, x []complex64, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chpsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&afp[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zhpsvx(fact byte, ul blas.Uplo, n int, nrhs int, ap []complex128, afp []complex128, ipiv []int32, b []complex128, ldb int, x []complex128, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhpsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&afp[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Chptrd(ul blas.Uplo, n int, ap []complex64, d []float32, e []float32, tau []complex64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chptrd((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.lapack_complex_float)(&tau[0])))
}

func Zhptrd(ul blas.Uplo, n int, ap []complex128, d []float64, e []float64, tau []complex128) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhptrd((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.lapack_complex_double)(&tau[0])))
}

func Chptrf(ul blas.Uplo, n int, ap []complex64, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chptrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_int)(&ipiv[0])))
}

func Zhptrf(ul blas.Uplo, n int, ap []complex128, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhptrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_int)(&ipiv[0])))
}

func Chptri(ul blas.Uplo, n int, ap []complex64, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chptri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_int)(&ipiv[0])))
}

func Zhptri(ul blas.Uplo, n int, ap []complex128, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhptri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_int)(&ipiv[0])))
}

func Chptrs(ul blas.Uplo, n int, nrhs int, ap []complex64, ipiv []int32, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chptrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zhptrs(ul blas.Uplo, n int, nrhs int, ap []complex128, ipiv []int32, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhptrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Shseqr(job lapack.Job, compz lapack.CompSV, n int, ilo int, ihi int, h []float32, ldh int, wr []float32, wi []float32, z []float32, ldz int) bool {
	return isZero(C.LAPACKE_shseqr((C.int)(rowMajor), (C.char)(job), (C.char)(compz), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.float)(&h[0]), (C.lapack_int)(ldh), (*C.float)(&wr[0]), (*C.float)(&wi[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dhseqr(job lapack.Job, compz lapack.CompSV, n int, ilo int, ihi int, h []float64, ldh int, wr []float64, wi []float64, z []float64, ldz int) bool {
	return isZero(C.LAPACKE_dhseqr((C.int)(rowMajor), (C.char)(job), (C.char)(compz), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.double)(&h[0]), (C.lapack_int)(ldh), (*C.double)(&wr[0]), (*C.double)(&wi[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Chseqr(job lapack.Job, compz lapack.CompSV, n int, ilo int, ihi int, h []complex64, ldh int, w []complex64, z []complex64, ldz int) bool {
	return isZero(C.LAPACKE_chseqr((C.int)(rowMajor), (C.char)(job), (C.char)(compz), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.lapack_complex_float)(&h[0]), (C.lapack_int)(ldh), (*C.lapack_complex_float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz)))
}

func Zhseqr(job lapack.Job, compz lapack.CompSV, n int, ilo int, ihi int, h []complex128, ldh int, w []complex128, z []complex128, ldz int) bool {
	return isZero(C.LAPACKE_zhseqr((C.int)(rowMajor), (C.char)(job), (C.char)(compz), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.lapack_complex_double)(&h[0]), (C.lapack_int)(ldh), (*C.lapack_complex_double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz)))
}

func Clacgv(n int, x []complex64, incx int) bool {
	return isZero(C.LAPACKE_clacgv((C.lapack_int)(n), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(incx)))
}

func Zlacgv(n int, x []complex128, incx int) bool {
	return isZero(C.LAPACKE_zlacgv((C.lapack_int)(n), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(incx)))
}

func Slacpy(ul blas.Uplo, m int, n int, a []float32, lda int, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_slacpy((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dlacpy(ul blas.Uplo, m int, n int, a []float64, lda int, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dlacpy((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Clacpy(ul blas.Uplo, m int, n int, a []complex64, lda int, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_clacpy((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zlacpy(ul blas.Uplo, m int, n int, a []complex128, lda int, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zlacpy((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Slamch(cmach byte) float32 {
	return float32(C.LAPACKE_slamch((C.char)(cmach)))
}

func Dlamch(cmach byte) float64 {
	return float64(C.LAPACKE_dlamch((C.char)(cmach)))
}

func Slange(norm byte, m int, n int, a []float32, lda int) float32 {
	return float32(C.LAPACKE_slange((C.int)(rowMajor), (C.char)(norm), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda)))
}

func Dlange(norm byte, m int, n int, a []float64, lda int) float64 {
	return float64(C.LAPACKE_dlange((C.int)(rowMajor), (C.char)(norm), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda)))
}

func Clange(norm byte, m int, n int, a []complex64, lda int) float32 {
	return float32(C.LAPACKE_clange((C.int)(rowMajor), (C.char)(norm), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda)))
}

func Zlange(norm byte, m int, n int, a []complex128, lda int) float64 {
	return float64(C.LAPACKE_zlange((C.int)(rowMajor), (C.char)(norm), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda)))
}

func Clanhe(norm byte, ul blas.Uplo, n int, a []complex64, lda int) float32 {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return float32(C.LAPACKE_clanhe((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda)))
}

func Zlanhe(norm byte, ul blas.Uplo, n int, a []complex128, lda int) float64 {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return float64(C.LAPACKE_zlanhe((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda)))
}

func Slansy(norm byte, ul blas.Uplo, n int, a []float32, lda int) float32 {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return float32(C.LAPACKE_slansy((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda)))
}

func Dlansy(norm byte, ul blas.Uplo, n int, a []float64, lda int) float64 {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return float64(C.LAPACKE_dlansy((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda)))
}

func Clansy(norm byte, ul blas.Uplo, n int, a []complex64, lda int) float32 {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return float32(C.LAPACKE_clansy((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda)))
}

func Zlansy(norm byte, ul blas.Uplo, n int, a []complex128, lda int) float64 {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return float64(C.LAPACKE_zlansy((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda)))
}

func Slantr(norm byte, ul blas.Uplo, d blas.Diag, m int, n int, a []float32, lda int) float32 {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return float32(C.LAPACKE_slantr((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda)))
}

func Dlantr(norm byte, ul blas.Uplo, d blas.Diag, m int, n int, a []float64, lda int) float64 {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return float64(C.LAPACKE_dlantr((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda)))
}

func Clantr(norm byte, ul blas.Uplo, d blas.Diag, m int, n int, a []complex64, lda int) float32 {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return float32(C.LAPACKE_clantr((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda)))
}

func Zlantr(norm byte, ul blas.Uplo, d blas.Diag, m int, n int, a []complex128, lda int) float64 {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return float64(C.LAPACKE_zlantr((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda)))
}

func Slarfb(s blas.Side, trans blas.Transpose, direct byte, storev byte, m int, n int, k int, v []float32, ldv int, t []float32, ldt int, c []float32, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_slarfb((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.char)(direct), (C.char)(storev), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.float)(&v[0]), (C.lapack_int)(ldv), (*C.float)(&t[0]), (C.lapack_int)(ldt), (*C.float)(&c[0]), (C.lapack_int)(ldc)))
}

func Dlarfb(s blas.Side, trans blas.Transpose, direct byte, storev byte, m int, n int, k int, v []float64, ldv int, t []float64, ldt int, c []float64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dlarfb((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.char)(direct), (C.char)(storev), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.double)(&v[0]), (C.lapack_int)(ldv), (*C.double)(&t[0]), (C.lapack_int)(ldt), (*C.double)(&c[0]), (C.lapack_int)(ldc)))
}

func Clarfb(s blas.Side, trans blas.Transpose, direct byte, storev byte, m int, n int, k int, v []complex64, ldv int, t []complex64, ldt int, c []complex64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_clarfb((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.char)(direct), (C.char)(storev), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_float)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_float)(&t[0]), (C.lapack_int)(ldt), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc)))
}

func Zlarfb(s blas.Side, trans blas.Transpose, direct byte, storev byte, m int, n int, k int, v []complex128, ldv int, t []complex128, ldt int, c []complex128, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zlarfb((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.char)(direct), (C.char)(storev), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_double)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_double)(&t[0]), (C.lapack_int)(ldt), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc)))
}

func Slarfg(n int, alpha []float32, x []float32, incx int, tau []float32) bool {
	return isZero(C.LAPACKE_slarfg((C.lapack_int)(n), (*C.float)(&alpha[0]), (*C.float)(&x[0]), (C.lapack_int)(incx), (*C.float)(&tau[0])))
}

func Dlarfg(n int, alpha []float64, x []float64, incx int, tau []float64) bool {
	return isZero(C.LAPACKE_dlarfg((C.lapack_int)(n), (*C.double)(&alpha[0]), (*C.double)(&x[0]), (C.lapack_int)(incx), (*C.double)(&tau[0])))
}

func Clarfg(n int, alpha []complex64, x []complex64, incx int, tau []complex64) bool {
	return isZero(C.LAPACKE_clarfg((C.lapack_int)(n), (*C.lapack_complex_float)(&alpha[0]), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(incx), (*C.lapack_complex_float)(&tau[0])))
}

func Zlarfg(n int, alpha []complex128, x []complex128, incx int, tau []complex128) bool {
	return isZero(C.LAPACKE_zlarfg((C.lapack_int)(n), (*C.lapack_complex_double)(&alpha[0]), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(incx), (*C.lapack_complex_double)(&tau[0])))
}

func Slarft(direct byte, storev byte, n int, k int, v []float32, ldv int, tau []float32, t []float32, ldt int) bool {
	return isZero(C.LAPACKE_slarft((C.int)(rowMajor), (C.char)(direct), (C.char)(storev), (C.lapack_int)(n), (C.lapack_int)(k), (*C.float)(&v[0]), (C.lapack_int)(ldv), (*C.float)(&tau[0]), (*C.float)(&t[0]), (C.lapack_int)(ldt)))
}

func Dlarft(direct byte, storev byte, n int, k int, v []float64, ldv int, tau []float64, t []float64, ldt int) bool {
	return isZero(C.LAPACKE_dlarft((C.int)(rowMajor), (C.char)(direct), (C.char)(storev), (C.lapack_int)(n), (C.lapack_int)(k), (*C.double)(&v[0]), (C.lapack_int)(ldv), (*C.double)(&tau[0]), (*C.double)(&t[0]), (C.lapack_int)(ldt)))
}

func Clarft(direct byte, storev byte, n int, k int, v []complex64, ldv int, tau []complex64, t []complex64, ldt int) bool {
	return isZero(C.LAPACKE_clarft((C.int)(rowMajor), (C.char)(direct), (C.char)(storev), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_float)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_float)(&tau[0]), (*C.lapack_complex_float)(&t[0]), (C.lapack_int)(ldt)))
}

func Zlarft(direct byte, storev byte, n int, k int, v []complex128, ldv int, tau []complex128, t []complex128, ldt int) bool {
	return isZero(C.LAPACKE_zlarft((C.int)(rowMajor), (C.char)(direct), (C.char)(storev), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_double)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_double)(&tau[0]), (*C.lapack_complex_double)(&t[0]), (C.lapack_int)(ldt)))
}

func Slarfx(s blas.Side, m int, n int, v []float32, tau float32, c []float32, ldc int, work []float32) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	return isZero(C.LAPACKE_slarfx((C.int)(rowMajor), (C.char)(s), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&v[0]), (C.float)(tau), (*C.float)(&c[0]), (C.lapack_int)(ldc), (*C.float)(&work[0])))
}

func Dlarfx(s blas.Side, m int, n int, v []float64, tau float64, c []float64, ldc int, work []float64) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	return isZero(C.LAPACKE_dlarfx((C.int)(rowMajor), (C.char)(s), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&v[0]), (C.double)(tau), (*C.double)(&c[0]), (C.lapack_int)(ldc), (*C.double)(&work[0])))
}

func Clarfx(s blas.Side, m int, n int, v []complex64, tau complex64, c []complex64, ldc int, work []complex64) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	return isZero(C.LAPACKE_clarfx((C.int)(rowMajor), (C.char)(s), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&v[0]), (C.lapack_complex_float)(tau), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc), (*C.lapack_complex_float)(&work[0])))
}

func Zlarfx(s blas.Side, m int, n int, v []complex128, tau complex128, c []complex128, ldc int, work []complex128) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	return isZero(C.LAPACKE_zlarfx((C.int)(rowMajor), (C.char)(s), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&v[0]), (C.lapack_complex_double)(tau), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc), (*C.lapack_complex_double)(&work[0])))
}

func Slarnv(idist int, iseed []int32, n int, x []float32) bool {
	return isZero(C.LAPACKE_slarnv((C.lapack_int)(idist), (*C.lapack_int)(&iseed[0]), (C.lapack_int)(n), (*C.float)(&x[0])))
}

func Dlarnv(idist int, iseed []int32, n int, x []float64) bool {
	return isZero(C.LAPACKE_dlarnv((C.lapack_int)(idist), (*C.lapack_int)(&iseed[0]), (C.lapack_int)(n), (*C.double)(&x[0])))
}

func Clarnv(idist int, iseed []int32, n int, x []complex64) bool {
	return isZero(C.LAPACKE_clarnv((C.lapack_int)(idist), (*C.lapack_int)(&iseed[0]), (C.lapack_int)(n), (*C.lapack_complex_float)(&x[0])))
}

func Zlarnv(idist int, iseed []int32, n int, x []complex128) bool {
	return isZero(C.LAPACKE_zlarnv((C.lapack_int)(idist), (*C.lapack_int)(&iseed[0]), (C.lapack_int)(n), (*C.lapack_complex_double)(&x[0])))
}

func Slaset(ul blas.Uplo, m int, n int, alpha float32, beta float32, a []float32, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_slaset((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(m), (C.lapack_int)(n), (C.float)(alpha), (C.float)(beta), (*C.float)(&a[0]), (C.lapack_int)(lda)))
}

func Dlaset(ul blas.Uplo, m int, n int, alpha float64, beta float64, a []float64, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dlaset((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(m), (C.lapack_int)(n), (C.double)(alpha), (C.double)(beta), (*C.double)(&a[0]), (C.lapack_int)(lda)))
}

func Claset(ul blas.Uplo, m int, n int, alpha complex64, beta complex64, a []complex64, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_claset((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_complex_float)(alpha), (C.lapack_complex_float)(beta), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda)))
}

func Zlaset(ul blas.Uplo, m int, n int, alpha complex128, beta complex128, a []complex128, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zlaset((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_complex_double)(alpha), (C.lapack_complex_double)(beta), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda)))
}

func Slasrt(id byte, n int, d []float32) bool {
	return isZero(C.LAPACKE_slasrt((C.char)(id), (C.lapack_int)(n), (*C.float)(&d[0])))
}

func Dlasrt(id byte, n int, d []float64) bool {
	return isZero(C.LAPACKE_dlasrt((C.char)(id), (C.lapack_int)(n), (*C.double)(&d[0])))
}

func Slaswp(n int, a []float32, lda int, k1 int, k2 int, ipiv []int32, incx int) bool {
	return isZero(C.LAPACKE_slaswp((C.int)(rowMajor), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (C.lapack_int)(k1), (C.lapack_int)(k2), (*C.lapack_int)(&ipiv[0]), (C.lapack_int)(incx)))
}

func Dlaswp(n int, a []float64, lda int, k1 int, k2 int, ipiv []int32, incx int) bool {
	return isZero(C.LAPACKE_dlaswp((C.int)(rowMajor), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (C.lapack_int)(k1), (C.lapack_int)(k2), (*C.lapack_int)(&ipiv[0]), (C.lapack_int)(incx)))
}

func Claswp(n int, a []complex64, lda int, k1 int, k2 int, ipiv []int32, incx int) bool {
	return isZero(C.LAPACKE_claswp((C.int)(rowMajor), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (C.lapack_int)(k1), (C.lapack_int)(k2), (*C.lapack_int)(&ipiv[0]), (C.lapack_int)(incx)))
}

func Zlaswp(n int, a []complex128, lda int, k1 int, k2 int, ipiv []int32, incx int) bool {
	return isZero(C.LAPACKE_zlaswp((C.int)(rowMajor), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (C.lapack_int)(k1), (C.lapack_int)(k2), (*C.lapack_int)(&ipiv[0]), (C.lapack_int)(incx)))
}

func Slauum(ul blas.Uplo, n int, a []float32, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_slauum((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda)))
}

func Dlauum(ul blas.Uplo, n int, a []float64, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dlauum((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda)))
}

func Clauum(ul blas.Uplo, n int, a []complex64, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_clauum((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda)))
}

func Zlauum(ul blas.Uplo, n int, a []complex128, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zlauum((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda)))
}

func Sopgtr(ul blas.Uplo, n int, ap []float32, tau []float32, q []float32, ldq int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sopgtr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.float)(&tau[0]), (*C.float)(&q[0]), (C.lapack_int)(ldq)))
}

func Dopgtr(ul blas.Uplo, n int, ap []float64, tau []float64, q []float64, ldq int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dopgtr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.double)(&tau[0]), (*C.double)(&q[0]), (C.lapack_int)(ldq)))
}

func Sopmtr(s blas.Side, ul blas.Uplo, trans blas.Transpose, m int, n int, ap []float32, tau []float32, c []float32, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sopmtr((C.int)(rowMajor), (C.char)(s), (C.char)(ul), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.float)(&tau[0]), (*C.float)(&c[0]), (C.lapack_int)(ldc)))
}

func Dopmtr(s blas.Side, ul blas.Uplo, trans blas.Transpose, m int, n int, ap []float64, tau []float64, c []float64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dopmtr((C.int)(rowMajor), (C.char)(s), (C.char)(ul), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.double)(&tau[0]), (*C.double)(&c[0]), (C.lapack_int)(ldc)))
}

func Sorgbr(vect byte, m int, n int, k int, a []float32, lda int, tau []float32) bool {
	return isZero(C.LAPACKE_sorgbr((C.int)(rowMajor), (C.char)(vect), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dorgbr(vect byte, m int, n int, k int, a []float64, lda int, tau []float64) bool {
	return isZero(C.LAPACKE_dorgbr((C.int)(rowMajor), (C.char)(vect), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Sorghr(n int, ilo int, ihi int, a []float32, lda int, tau []float32) bool {
	return isZero(C.LAPACKE_sorghr((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dorghr(n int, ilo int, ihi int, a []float64, lda int, tau []float64) bool {
	return isZero(C.LAPACKE_dorghr((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Sorglq(m int, n int, k int, a []float32, lda int, tau []float32) bool {
	return isZero(C.LAPACKE_sorglq((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dorglq(m int, n int, k int, a []float64, lda int, tau []float64) bool {
	return isZero(C.LAPACKE_dorglq((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Sorgql(m int, n int, k int, a []float32, lda int, tau []float32) bool {
	return isZero(C.LAPACKE_sorgql((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dorgql(m int, n int, k int, a []float64, lda int, tau []float64) bool {
	return isZero(C.LAPACKE_dorgql((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Sorgqr(m int, n int, k int, a []float32, lda int, tau []float32) bool {
	return isZero(C.LAPACKE_sorgqr((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dorgqr(m int, n int, k int, a []float64, lda int, tau []float64) bool {
	return isZero(C.LAPACKE_dorgqr((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Sorgrq(m int, n int, k int, a []float32, lda int, tau []float32) bool {
	return isZero(C.LAPACKE_sorgrq((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dorgrq(m int, n int, k int, a []float64, lda int, tau []float64) bool {
	return isZero(C.LAPACKE_dorgrq((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Sorgtr(ul blas.Uplo, n int, a []float32, lda int, tau []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sorgtr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dorgtr(ul blas.Uplo, n int, a []float64, lda int, tau []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dorgtr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Sormbr(vect byte, s blas.Side, trans blas.Transpose, m int, n int, k int, a []float32, lda int, tau []float32, c []float32, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sormbr((C.int)(rowMajor), (C.char)(vect), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0]), (*C.float)(&c[0]), (C.lapack_int)(ldc)))
}

func Dormbr(vect byte, s blas.Side, trans blas.Transpose, m int, n int, k int, a []float64, lda int, tau []float64, c []float64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dormbr((C.int)(rowMajor), (C.char)(vect), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0]), (*C.double)(&c[0]), (C.lapack_int)(ldc)))
}

func Sormhr(s blas.Side, trans blas.Transpose, m int, n int, ilo int, ihi int, a []float32, lda int, tau []float32, c []float32, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sormhr((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0]), (*C.float)(&c[0]), (C.lapack_int)(ldc)))
}

func Dormhr(s blas.Side, trans blas.Transpose, m int, n int, ilo int, ihi int, a []float64, lda int, tau []float64, c []float64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dormhr((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0]), (*C.double)(&c[0]), (C.lapack_int)(ldc)))
}

func Sormlq(s blas.Side, trans blas.Transpose, m int, n int, k int, a []float32, lda int, tau []float32, c []float32, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sormlq((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0]), (*C.float)(&c[0]), (C.lapack_int)(ldc)))
}

func Dormlq(s blas.Side, trans blas.Transpose, m int, n int, k int, a []float64, lda int, tau []float64, c []float64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dormlq((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0]), (*C.double)(&c[0]), (C.lapack_int)(ldc)))
}

func Sormql(s blas.Side, trans blas.Transpose, m int, n int, k int, a []float32, lda int, tau []float32, c []float32, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sormql((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0]), (*C.float)(&c[0]), (C.lapack_int)(ldc)))
}

func Dormql(s blas.Side, trans blas.Transpose, m int, n int, k int, a []float64, lda int, tau []float64, c []float64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dormql((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0]), (*C.double)(&c[0]), (C.lapack_int)(ldc)))
}

func Sormqr(s blas.Side, trans blas.Transpose, m int, n int, k int, a []float32, lda int, tau []float32, c []float32, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sormqr((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0]), (*C.float)(&c[0]), (C.lapack_int)(ldc)))
}

func Dormqr(s blas.Side, trans blas.Transpose, m int, n int, k int, a []float64, lda int, tau []float64, c []float64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dormqr((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0]), (*C.double)(&c[0]), (C.lapack_int)(ldc)))
}

func Sormrq(s blas.Side, trans blas.Transpose, m int, n int, k int, a []float32, lda int, tau []float32, c []float32, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sormrq((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0]), (*C.float)(&c[0]), (C.lapack_int)(ldc)))
}

func Dormrq(s blas.Side, trans blas.Transpose, m int, n int, k int, a []float64, lda int, tau []float64, c []float64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dormrq((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0]), (*C.double)(&c[0]), (C.lapack_int)(ldc)))
}

func Sormrz(s blas.Side, trans blas.Transpose, m int, n int, k int, l int, a []float32, lda int, tau []float32, c []float32, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sormrz((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0]), (*C.float)(&c[0]), (C.lapack_int)(ldc)))
}

func Dormrz(s blas.Side, trans blas.Transpose, m int, n int, k int, l int, a []float64, lda int, tau []float64, c []float64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dormrz((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0]), (*C.double)(&c[0]), (C.lapack_int)(ldc)))
}

func Sormtr(s blas.Side, ul blas.Uplo, trans blas.Transpose, m int, n int, a []float32, lda int, tau []float32, c []float32, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sormtr((C.int)(rowMajor), (C.char)(s), (C.char)(ul), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0]), (*C.float)(&c[0]), (C.lapack_int)(ldc)))
}

func Dormtr(s blas.Side, ul blas.Uplo, trans blas.Transpose, m int, n int, a []float64, lda int, tau []float64, c []float64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dormtr((C.int)(rowMajor), (C.char)(s), (C.char)(ul), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0]), (*C.double)(&c[0]), (C.lapack_int)(ldc)))
}

func Spbcon(ul blas.Uplo, n int, kd int, ab []float32, ldab int, anorm float32, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spbcon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Dpbcon(ul blas.Uplo, n int, kd int, ab []float64, ldab int, anorm float64, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpbcon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Cpbcon(ul blas.Uplo, n int, kd int, ab []complex64, ldab int, anorm float32, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpbcon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Zpbcon(ul blas.Uplo, n int, kd int, ab []complex128, ldab int, anorm float64, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpbcon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Spbequ(ul blas.Uplo, n int, kd int, ab []float32, ldab int, s []float32, scond []float32, amax []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spbequ((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&s[0]), (*C.float)(&scond[0]), (*C.float)(&amax[0])))
}

func Dpbequ(ul blas.Uplo, n int, kd int, ab []float64, ldab int, s []float64, scond []float64, amax []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpbequ((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&s[0]), (*C.double)(&scond[0]), (*C.double)(&amax[0])))
}

func Cpbequ(ul blas.Uplo, n int, kd int, ab []complex64, ldab int, s []float32, scond []float32, amax []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpbequ((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&s[0]), (*C.float)(&scond[0]), (*C.float)(&amax[0])))
}

func Zpbequ(ul blas.Uplo, n int, kd int, ab []complex128, ldab int, s []float64, scond []float64, amax []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpbequ((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&s[0]), (*C.double)(&scond[0]), (*C.double)(&amax[0])))
}

func Spbrfs(ul blas.Uplo, n int, kd int, nrhs int, ab []float32, ldab int, afb []float32, ldafb int, b []float32, ldb int, x []float32, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spbrfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&afb[0]), (C.lapack_int)(ldafb), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dpbrfs(ul blas.Uplo, n int, kd int, nrhs int, ab []float64, ldab int, afb []float64, ldafb int, b []float64, ldb int, x []float64, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpbrfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&afb[0]), (C.lapack_int)(ldafb), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Cpbrfs(ul blas.Uplo, n int, kd int, nrhs int, ab []complex64, ldab int, afb []complex64, ldafb int, b []complex64, ldb int, x []complex64, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpbrfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_float)(&afb[0]), (C.lapack_int)(ldafb), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zpbrfs(ul blas.Uplo, n int, kd int, nrhs int, ab []complex128, ldab int, afb []complex128, ldafb int, b []complex128, ldb int, x []complex128, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpbrfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_double)(&afb[0]), (C.lapack_int)(ldafb), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Spbstf(ul blas.Uplo, n int, kb int, bb []float32, ldbb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spbstf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kb), (*C.float)(&bb[0]), (C.lapack_int)(ldbb)))
}

func Dpbstf(ul blas.Uplo, n int, kb int, bb []float64, ldbb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpbstf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kb), (*C.double)(&bb[0]), (C.lapack_int)(ldbb)))
}

func Cpbstf(ul blas.Uplo, n int, kb int, bb []complex64, ldbb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpbstf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kb), (*C.lapack_complex_float)(&bb[0]), (C.lapack_int)(ldbb)))
}

func Zpbstf(ul blas.Uplo, n int, kb int, bb []complex128, ldbb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpbstf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kb), (*C.lapack_complex_double)(&bb[0]), (C.lapack_int)(ldbb)))
}

func Spbsv(ul blas.Uplo, n int, kd int, nrhs int, ab []float32, ldab int, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spbsv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dpbsv(ul blas.Uplo, n int, kd int, nrhs int, ab []float64, ldab int, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpbsv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cpbsv(ul blas.Uplo, n int, kd int, nrhs int, ab []complex64, ldab int, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpbsv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zpbsv(ul blas.Uplo, n int, kd int, nrhs int, ab []complex128, ldab int, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpbsv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Spbsvx(fact byte, ul blas.Uplo, n int, kd int, nrhs int, ab []float32, ldab int, afb []float32, ldafb int, equed []byte, s []float32, b []float32, ldb int, x []float32, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spbsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&afb[0]), (C.lapack_int)(ldafb), (*C.char)(unsafe.Pointer(&equed[0])), (*C.float)(&s[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dpbsvx(fact byte, ul blas.Uplo, n int, kd int, nrhs int, ab []float64, ldab int, afb []float64, ldafb int, equed []byte, s []float64, b []float64, ldb int, x []float64, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpbsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&afb[0]), (C.lapack_int)(ldafb), (*C.char)(unsafe.Pointer(&equed[0])), (*C.double)(&s[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Cpbsvx(fact byte, ul blas.Uplo, n int, kd int, nrhs int, ab []complex64, ldab int, afb []complex64, ldafb int, equed []byte, s []float32, b []complex64, ldb int, x []complex64, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpbsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_float)(&afb[0]), (C.lapack_int)(ldafb), (*C.char)(unsafe.Pointer(&equed[0])), (*C.float)(&s[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zpbsvx(fact byte, ul blas.Uplo, n int, kd int, nrhs int, ab []complex128, ldab int, afb []complex128, ldafb int, equed []byte, s []float64, b []complex128, ldb int, x []complex128, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpbsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_double)(&afb[0]), (C.lapack_int)(ldafb), (*C.char)(unsafe.Pointer(&equed[0])), (*C.double)(&s[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Spbtrf(ul blas.Uplo, n int, kd int, ab []float32, ldab int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spbtrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.float)(&ab[0]), (C.lapack_int)(ldab)))
}

func Dpbtrf(ul blas.Uplo, n int, kd int, ab []float64, ldab int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpbtrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.double)(&ab[0]), (C.lapack_int)(ldab)))
}

func Cpbtrf(ul blas.Uplo, n int, kd int, ab []complex64, ldab int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpbtrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab)))
}

func Zpbtrf(ul blas.Uplo, n int, kd int, ab []complex128, ldab int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpbtrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab)))
}

func Spbtrs(ul blas.Uplo, n int, kd int, nrhs int, ab []float32, ldab int, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spbtrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dpbtrs(ul blas.Uplo, n int, kd int, nrhs int, ab []float64, ldab int, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpbtrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cpbtrs(ul blas.Uplo, n int, kd int, nrhs int, ab []complex64, ldab int, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpbtrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zpbtrs(ul blas.Uplo, n int, kd int, nrhs int, ab []complex128, ldab int, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpbtrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Spftrf(transr blas.Transpose, ul blas.Uplo, n int, a []float32) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spftrf((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0])))
}

func Dpftrf(transr blas.Transpose, ul blas.Uplo, n int, a []float64) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpftrf((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0])))
}

func Cpftrf(transr blas.Transpose, ul blas.Uplo, n int, a []complex64) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpftrf((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0])))
}

func Zpftrf(transr blas.Transpose, ul blas.Uplo, n int, a []complex128) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpftrf((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0])))
}

func Spftri(transr blas.Transpose, ul blas.Uplo, n int, a []float32) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spftri((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0])))
}

func Dpftri(transr blas.Transpose, ul blas.Uplo, n int, a []float64) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpftri((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0])))
}

func Cpftri(transr blas.Transpose, ul blas.Uplo, n int, a []complex64) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpftri((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0])))
}

func Zpftri(transr blas.Transpose, ul blas.Uplo, n int, a []complex128) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpftri((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0])))
}

func Spftrs(transr blas.Transpose, ul blas.Uplo, n int, nrhs int, a []float32, b []float32, ldb int) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spftrs((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dpftrs(transr blas.Transpose, ul blas.Uplo, n int, nrhs int, a []float64, b []float64, ldb int) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpftrs((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cpftrs(transr blas.Transpose, ul blas.Uplo, n int, nrhs int, a []complex64, b []complex64, ldb int) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpftrs((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zpftrs(transr blas.Transpose, ul blas.Uplo, n int, nrhs int, a []complex128, b []complex128, ldb int) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpftrs((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Spocon(ul blas.Uplo, n int, a []float32, lda int, anorm float32, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spocon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Dpocon(ul blas.Uplo, n int, a []float64, lda int, anorm float64, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpocon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Cpocon(ul blas.Uplo, n int, a []complex64, lda int, anorm float32, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpocon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Zpocon(ul blas.Uplo, n int, a []complex128, lda int, anorm float64, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpocon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Spoequ(n int, a []float32, lda int, s []float32, scond []float32, amax []float32) bool {
	return isZero(C.LAPACKE_spoequ((C.int)(rowMajor), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&s[0]), (*C.float)(&scond[0]), (*C.float)(&amax[0])))
}

func Dpoequ(n int, a []float64, lda int, s []float64, scond []float64, amax []float64) bool {
	return isZero(C.LAPACKE_dpoequ((C.int)(rowMajor), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&s[0]), (*C.double)(&scond[0]), (*C.double)(&amax[0])))
}

func Cpoequ(n int, a []complex64, lda int, s []float32, scond []float32, amax []float32) bool {
	return isZero(C.LAPACKE_cpoequ((C.int)(rowMajor), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&s[0]), (*C.float)(&scond[0]), (*C.float)(&amax[0])))
}

func Zpoequ(n int, a []complex128, lda int, s []float64, scond []float64, amax []float64) bool {
	return isZero(C.LAPACKE_zpoequ((C.int)(rowMajor), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&s[0]), (*C.double)(&scond[0]), (*C.double)(&amax[0])))
}

func Spoequb(n int, a []float32, lda int, s []float32, scond []float32, amax []float32) bool {
	return isZero(C.LAPACKE_spoequb((C.int)(rowMajor), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&s[0]), (*C.float)(&scond[0]), (*C.float)(&amax[0])))
}

func Dpoequb(n int, a []float64, lda int, s []float64, scond []float64, amax []float64) bool {
	return isZero(C.LAPACKE_dpoequb((C.int)(rowMajor), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&s[0]), (*C.double)(&scond[0]), (*C.double)(&amax[0])))
}

func Cpoequb(n int, a []complex64, lda int, s []float32, scond []float32, amax []float32) bool {
	return isZero(C.LAPACKE_cpoequb((C.int)(rowMajor), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&s[0]), (*C.float)(&scond[0]), (*C.float)(&amax[0])))
}

func Zpoequb(n int, a []complex128, lda int, s []float64, scond []float64, amax []float64) bool {
	return isZero(C.LAPACKE_zpoequb((C.int)(rowMajor), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&s[0]), (*C.double)(&scond[0]), (*C.double)(&amax[0])))
}

func Sporfs(ul blas.Uplo, n int, nrhs int, a []float32, lda int, af []float32, ldaf int, b []float32, ldb int, x []float32, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sporfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&af[0]), (C.lapack_int)(ldaf), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dporfs(ul blas.Uplo, n int, nrhs int, a []float64, lda int, af []float64, ldaf int, b []float64, ldb int, x []float64, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dporfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&af[0]), (C.lapack_int)(ldaf), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Cporfs(ul blas.Uplo, n int, nrhs int, a []complex64, lda int, af []complex64, ldaf int, b []complex64, ldb int, x []complex64, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cporfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zporfs(ul blas.Uplo, n int, nrhs int, a []complex128, lda int, af []complex128, ldaf int, b []complex128, ldb int, x []complex128, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zporfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Sposv(ul blas.Uplo, n int, nrhs int, a []float32, lda int, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sposv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dposv(ul blas.Uplo, n int, nrhs int, a []float64, lda int, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dposv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cposv(ul blas.Uplo, n int, nrhs int, a []complex64, lda int, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cposv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zposv(ul blas.Uplo, n int, nrhs int, a []complex128, lda int, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zposv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Dsposv(ul blas.Uplo, n int, nrhs int, a []float64, lda int, b []float64, ldb int, x []float64, ldx int, iter []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsposv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.lapack_int)(&iter[0])))
}

func Zcposv(ul blas.Uplo, n int, nrhs int, a []complex128, lda int, b []complex128, ldb int, x []complex128, ldx int, iter []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zcposv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.lapack_int)(&iter[0])))
}

func Sposvx(fact byte, ul blas.Uplo, n int, nrhs int, a []float32, lda int, af []float32, ldaf int, equed []byte, s []float32, b []float32, ldb int, x []float32, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sposvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&af[0]), (C.lapack_int)(ldaf), (*C.char)(unsafe.Pointer(&equed[0])), (*C.float)(&s[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dposvx(fact byte, ul blas.Uplo, n int, nrhs int, a []float64, lda int, af []float64, ldaf int, equed []byte, s []float64, b []float64, ldb int, x []float64, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dposvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&af[0]), (C.lapack_int)(ldaf), (*C.char)(unsafe.Pointer(&equed[0])), (*C.double)(&s[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Cposvx(fact byte, ul blas.Uplo, n int, nrhs int, a []complex64, lda int, af []complex64, ldaf int, equed []byte, s []float32, b []complex64, ldb int, x []complex64, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cposvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&af[0]), (C.lapack_int)(ldaf), (*C.char)(unsafe.Pointer(&equed[0])), (*C.float)(&s[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zposvx(fact byte, ul blas.Uplo, n int, nrhs int, a []complex128, lda int, af []complex128, ldaf int, equed []byte, s []float64, b []complex128, ldb int, x []complex128, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zposvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&af[0]), (C.lapack_int)(ldaf), (*C.char)(unsafe.Pointer(&equed[0])), (*C.double)(&s[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Spotrf(ul blas.Uplo, n int, a []float32, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spotrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda)))
}

func Dpotrf(ul blas.Uplo, n int, a []float64, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpotrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda)))
}

func Cpotrf(ul blas.Uplo, n int, a []complex64, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpotrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda)))
}

func Zpotrf(ul blas.Uplo, n int, a []complex128, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpotrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda)))
}

func Spotri(ul blas.Uplo, n int, a []float32, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spotri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda)))
}

func Dpotri(ul blas.Uplo, n int, a []float64, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpotri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda)))
}

func Cpotri(ul blas.Uplo, n int, a []complex64, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpotri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda)))
}

func Zpotri(ul blas.Uplo, n int, a []complex128, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpotri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda)))
}

func Spotrs(ul blas.Uplo, n int, nrhs int, a []float32, lda int, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spotrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dpotrs(ul blas.Uplo, n int, nrhs int, a []float64, lda int, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpotrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cpotrs(ul blas.Uplo, n int, nrhs int, a []complex64, lda int, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpotrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zpotrs(ul blas.Uplo, n int, nrhs int, a []complex128, lda int, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpotrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Sppcon(ul blas.Uplo, n int, ap []float32, anorm float32, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sppcon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Dppcon(ul blas.Uplo, n int, ap []float64, anorm float64, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dppcon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Cppcon(ul blas.Uplo, n int, ap []complex64, anorm float32, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cppcon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Zppcon(ul blas.Uplo, n int, ap []complex128, anorm float64, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zppcon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Sppequ(ul blas.Uplo, n int, ap []float32, s []float32, scond []float32, amax []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sppequ((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.float)(&s[0]), (*C.float)(&scond[0]), (*C.float)(&amax[0])))
}

func Dppequ(ul blas.Uplo, n int, ap []float64, s []float64, scond []float64, amax []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dppequ((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.double)(&s[0]), (*C.double)(&scond[0]), (*C.double)(&amax[0])))
}

func Cppequ(ul blas.Uplo, n int, ap []complex64, s []float32, scond []float32, amax []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cppequ((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.float)(&s[0]), (*C.float)(&scond[0]), (*C.float)(&amax[0])))
}

func Zppequ(ul blas.Uplo, n int, ap []complex128, s []float64, scond []float64, amax []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zppequ((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.double)(&s[0]), (*C.double)(&scond[0]), (*C.double)(&amax[0])))
}

func Spprfs(ul blas.Uplo, n int, nrhs int, ap []float32, afp []float32, b []float32, ldb int, x []float32, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spprfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&ap[0]), (*C.float)(&afp[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dpprfs(ul blas.Uplo, n int, nrhs int, ap []float64, afp []float64, b []float64, ldb int, x []float64, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpprfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&ap[0]), (*C.double)(&afp[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Cpprfs(ul blas.Uplo, n int, nrhs int, ap []complex64, afp []complex64, b []complex64, ldb int, x []complex64, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpprfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&afp[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zpprfs(ul blas.Uplo, n int, nrhs int, ap []complex128, afp []complex128, b []complex128, ldb int, x []complex128, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpprfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&afp[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Sppsv(ul blas.Uplo, n int, nrhs int, ap []float32, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sppsv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&ap[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dppsv(ul blas.Uplo, n int, nrhs int, ap []float64, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dppsv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&ap[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cppsv(ul blas.Uplo, n int, nrhs int, ap []complex64, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cppsv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zppsv(ul blas.Uplo, n int, nrhs int, ap []complex128, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zppsv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Sppsvx(fact byte, ul blas.Uplo, n int, nrhs int, ap []float32, afp []float32, equed []byte, s []float32, b []float32, ldb int, x []float32, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sppsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&ap[0]), (*C.float)(&afp[0]), (*C.char)(unsafe.Pointer(&equed[0])), (*C.float)(&s[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dppsvx(fact byte, ul blas.Uplo, n int, nrhs int, ap []float64, afp []float64, equed []byte, s []float64, b []float64, ldb int, x []float64, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dppsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&ap[0]), (*C.double)(&afp[0]), (*C.char)(unsafe.Pointer(&equed[0])), (*C.double)(&s[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Cppsvx(fact byte, ul blas.Uplo, n int, nrhs int, ap []complex64, afp []complex64, equed []byte, s []float32, b []complex64, ldb int, x []complex64, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cppsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&afp[0]), (*C.char)(unsafe.Pointer(&equed[0])), (*C.float)(&s[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zppsvx(fact byte, ul blas.Uplo, n int, nrhs int, ap []complex128, afp []complex128, equed []byte, s []float64, b []complex128, ldb int, x []complex128, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zppsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&afp[0]), (*C.char)(unsafe.Pointer(&equed[0])), (*C.double)(&s[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Spptrf(ul blas.Uplo, n int, ap []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spptrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0])))
}

func Dpptrf(ul blas.Uplo, n int, ap []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpptrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0])))
}

func Cpptrf(ul blas.Uplo, n int, ap []complex64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpptrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0])))
}

func Zpptrf(ul blas.Uplo, n int, ap []complex128) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpptrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0])))
}

func Spptri(ul blas.Uplo, n int, ap []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spptri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0])))
}

func Dpptri(ul blas.Uplo, n int, ap []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpptri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0])))
}

func Cpptri(ul blas.Uplo, n int, ap []complex64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpptri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0])))
}

func Zpptri(ul blas.Uplo, n int, ap []complex128) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpptri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0])))
}

func Spptrs(ul blas.Uplo, n int, nrhs int, ap []float32, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spptrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&ap[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dpptrs(ul blas.Uplo, n int, nrhs int, ap []float64, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpptrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&ap[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cpptrs(ul blas.Uplo, n int, nrhs int, ap []complex64, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpptrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zpptrs(ul blas.Uplo, n int, nrhs int, ap []complex128, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpptrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Spstrf(ul blas.Uplo, n int, a []float32, lda int, piv []int32, rank []int32, tol float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_spstrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&piv[0]), (*C.lapack_int)(&rank[0]), (C.float)(tol)))
}

func Dpstrf(ul blas.Uplo, n int, a []float64, lda int, piv []int32, rank []int32, tol float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dpstrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&piv[0]), (*C.lapack_int)(&rank[0]), (C.double)(tol)))
}

func Cpstrf(ul blas.Uplo, n int, a []complex64, lda int, piv []int32, rank []int32, tol float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpstrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&piv[0]), (*C.lapack_int)(&rank[0]), (C.float)(tol)))
}

func Zpstrf(ul blas.Uplo, n int, a []complex128, lda int, piv []int32, rank []int32, tol float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpstrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&piv[0]), (*C.lapack_int)(&rank[0]), (C.double)(tol)))
}

func Sptcon(n int, d []float32, e []float32, anorm float32, rcond []float32) bool {
	return isZero(C.LAPACKE_sptcon((C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Dptcon(n int, d []float64, e []float64, anorm float64, rcond []float64) bool {
	return isZero(C.LAPACKE_dptcon((C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Cptcon(n int, d []float32, e []complex64, anorm float32, rcond []float32) bool {
	return isZero(C.LAPACKE_cptcon((C.lapack_int)(n), (*C.float)(&d[0]), (*C.lapack_complex_float)(&e[0]), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Zptcon(n int, d []float64, e []complex128, anorm float64, rcond []float64) bool {
	return isZero(C.LAPACKE_zptcon((C.lapack_int)(n), (*C.double)(&d[0]), (*C.lapack_complex_double)(&e[0]), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Spteqr(compz lapack.CompSV, n int, d []float32, e []float32, z []float32, ldz int) bool {
	return isZero(C.LAPACKE_spteqr((C.int)(rowMajor), (C.char)(compz), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dpteqr(compz lapack.CompSV, n int, d []float64, e []float64, z []float64, ldz int) bool {
	return isZero(C.LAPACKE_dpteqr((C.int)(rowMajor), (C.char)(compz), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Cpteqr(compz lapack.CompSV, n int, d []float32, e []float32, z []complex64, ldz int) bool {
	return isZero(C.LAPACKE_cpteqr((C.int)(rowMajor), (C.char)(compz), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz)))
}

func Zpteqr(compz lapack.CompSV, n int, d []float64, e []float64, z []complex128, ldz int) bool {
	return isZero(C.LAPACKE_zpteqr((C.int)(rowMajor), (C.char)(compz), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz)))
}

func Sptrfs(n int, nrhs int, d []float32, e []float32, df []float32, ef []float32, b []float32, ldb int, x []float32, ldx int, ferr []float32, berr []float32) bool {
	return isZero(C.LAPACKE_sptrfs((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&df[0]), (*C.float)(&ef[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dptrfs(n int, nrhs int, d []float64, e []float64, df []float64, ef []float64, b []float64, ldb int, x []float64, ldx int, ferr []float64, berr []float64) bool {
	return isZero(C.LAPACKE_dptrfs((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&df[0]), (*C.double)(&ef[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Cptrfs(ul blas.Uplo, n int, nrhs int, d []float32, e []complex64, df []float32, ef []complex64, b []complex64, ldb int, x []complex64, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cptrfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&d[0]), (*C.lapack_complex_float)(&e[0]), (*C.float)(&df[0]), (*C.lapack_complex_float)(&ef[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zptrfs(ul blas.Uplo, n int, nrhs int, d []float64, e []complex128, df []float64, ef []complex128, b []complex128, ldb int, x []complex128, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zptrfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&d[0]), (*C.lapack_complex_double)(&e[0]), (*C.double)(&df[0]), (*C.lapack_complex_double)(&ef[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Sptsv(n int, nrhs int, d []float32, e []float32, b []float32, ldb int) bool {
	return isZero(C.LAPACKE_sptsv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dptsv(n int, nrhs int, d []float64, e []float64, b []float64, ldb int) bool {
	return isZero(C.LAPACKE_dptsv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cptsv(n int, nrhs int, d []float32, e []complex64, b []complex64, ldb int) bool {
	return isZero(C.LAPACKE_cptsv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&d[0]), (*C.lapack_complex_float)(&e[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zptsv(n int, nrhs int, d []float64, e []complex128, b []complex128, ldb int) bool {
	return isZero(C.LAPACKE_zptsv((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&d[0]), (*C.lapack_complex_double)(&e[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Sptsvx(fact byte, n int, nrhs int, d []float32, e []float32, df []float32, ef []float32, b []float32, ldb int, x []float32, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	return isZero(C.LAPACKE_sptsvx((C.int)(rowMajor), (C.char)(fact), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&df[0]), (*C.float)(&ef[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dptsvx(fact byte, n int, nrhs int, d []float64, e []float64, df []float64, ef []float64, b []float64, ldb int, x []float64, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	return isZero(C.LAPACKE_dptsvx((C.int)(rowMajor), (C.char)(fact), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&df[0]), (*C.double)(&ef[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Cptsvx(fact byte, n int, nrhs int, d []float32, e []complex64, df []float32, ef []complex64, b []complex64, ldb int, x []complex64, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	return isZero(C.LAPACKE_cptsvx((C.int)(rowMajor), (C.char)(fact), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&d[0]), (*C.lapack_complex_float)(&e[0]), (*C.float)(&df[0]), (*C.lapack_complex_float)(&ef[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zptsvx(fact byte, n int, nrhs int, d []float64, e []complex128, df []float64, ef []complex128, b []complex128, ldb int, x []complex128, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	return isZero(C.LAPACKE_zptsvx((C.int)(rowMajor), (C.char)(fact), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&d[0]), (*C.lapack_complex_double)(&e[0]), (*C.double)(&df[0]), (*C.lapack_complex_double)(&ef[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Spttrf(n int, d []float32, e []float32) bool {
	return isZero(C.LAPACKE_spttrf((C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0])))
}

func Dpttrf(n int, d []float64, e []float64) bool {
	return isZero(C.LAPACKE_dpttrf((C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0])))
}

func Cpttrf(n int, d []float32, e []complex64) bool {
	return isZero(C.LAPACKE_cpttrf((C.lapack_int)(n), (*C.float)(&d[0]), (*C.lapack_complex_float)(&e[0])))
}

func Zpttrf(n int, d []float64, e []complex128) bool {
	return isZero(C.LAPACKE_zpttrf((C.lapack_int)(n), (*C.double)(&d[0]), (*C.lapack_complex_double)(&e[0])))
}

func Spttrs(n int, nrhs int, d []float32, e []float32, b []float32, ldb int) bool {
	return isZero(C.LAPACKE_spttrs((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dpttrs(n int, nrhs int, d []float64, e []float64, b []float64, ldb int) bool {
	return isZero(C.LAPACKE_dpttrs((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cpttrs(ul blas.Uplo, n int, nrhs int, d []float32, e []complex64, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cpttrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&d[0]), (*C.lapack_complex_float)(&e[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zpttrs(ul blas.Uplo, n int, nrhs int, d []float64, e []complex128, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zpttrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&d[0]), (*C.lapack_complex_double)(&e[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Ssbev(jobz lapack.Job, ul blas.Uplo, n int, kd int, ab []float32, ldab int, w []float32, z []float32, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssbev((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dsbev(jobz lapack.Job, ul blas.Uplo, n int, kd int, ab []float64, ldab int, w []float64, z []float64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsbev((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Ssbevd(jobz lapack.Job, ul blas.Uplo, n int, kd int, ab []float32, ldab int, w []float32, z []float32, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssbevd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dsbevd(jobz lapack.Job, ul blas.Uplo, n int, kd int, ab []float64, ldab int, w []float64, z []float64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsbevd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Ssbevx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, kd int, ab []float32, ldab int, q []float32, ldq int, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []float32, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssbevx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&q[0]), (C.lapack_int)(ldq), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Dsbevx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, kd int, ab []float64, ldab int, q []float64, ldq int, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []float64, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsbevx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&q[0]), (C.lapack_int)(ldq), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Ssbgst(vect byte, ul blas.Uplo, n int, ka int, kb int, ab []float32, ldab int, bb []float32, ldbb int, x []float32, ldx int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssbgst((C.int)(rowMajor), (C.char)(vect), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&bb[0]), (C.lapack_int)(ldbb), (*C.float)(&x[0]), (C.lapack_int)(ldx)))
}

func Dsbgst(vect byte, ul blas.Uplo, n int, ka int, kb int, ab []float64, ldab int, bb []float64, ldbb int, x []float64, ldx int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsbgst((C.int)(rowMajor), (C.char)(vect), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&bb[0]), (C.lapack_int)(ldbb), (*C.double)(&x[0]), (C.lapack_int)(ldx)))
}

func Ssbgv(jobz lapack.Job, ul blas.Uplo, n int, ka int, kb int, ab []float32, ldab int, bb []float32, ldbb int, w []float32, z []float32, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssbgv((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&bb[0]), (C.lapack_int)(ldbb), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dsbgv(jobz lapack.Job, ul blas.Uplo, n int, ka int, kb int, ab []float64, ldab int, bb []float64, ldbb int, w []float64, z []float64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsbgv((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&bb[0]), (C.lapack_int)(ldbb), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Ssbgvd(jobz lapack.Job, ul blas.Uplo, n int, ka int, kb int, ab []float32, ldab int, bb []float32, ldbb int, w []float32, z []float32, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssbgvd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&bb[0]), (C.lapack_int)(ldbb), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dsbgvd(jobz lapack.Job, ul blas.Uplo, n int, ka int, kb int, ab []float64, ldab int, bb []float64, ldbb int, w []float64, z []float64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsbgvd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&bb[0]), (C.lapack_int)(ldbb), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Ssbgvx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, ka int, kb int, ab []float32, ldab int, bb []float32, ldbb int, q []float32, ldq int, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []float32, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssbgvx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&bb[0]), (C.lapack_int)(ldbb), (*C.float)(&q[0]), (C.lapack_int)(ldq), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Dsbgvx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, ka int, kb int, ab []float64, ldab int, bb []float64, ldbb int, q []float64, ldq int, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []float64, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsbgvx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(ka), (C.lapack_int)(kb), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&bb[0]), (C.lapack_int)(ldbb), (*C.double)(&q[0]), (C.lapack_int)(ldq), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Ssbtrd(vect byte, ul blas.Uplo, n int, kd int, ab []float32, ldab int, d []float32, e []float32, q []float32, ldq int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssbtrd((C.int)(rowMajor), (C.char)(vect), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&q[0]), (C.lapack_int)(ldq)))
}

func Dsbtrd(vect byte, ul blas.Uplo, n int, kd int, ab []float64, ldab int, d []float64, e []float64, q []float64, ldq int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsbtrd((C.int)(rowMajor), (C.char)(vect), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&q[0]), (C.lapack_int)(ldq)))
}

func Ssfrk(transr blas.Transpose, ul blas.Uplo, trans blas.Transpose, n int, k int, alpha float32, a []float32, lda int, beta float32, c []float32) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_ssfrk((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(k), (C.float)(alpha), (*C.float)(&a[0]), (C.lapack_int)(lda), (C.float)(beta), (*C.float)(&c[0])))
}

func Dsfrk(transr blas.Transpose, ul blas.Uplo, trans blas.Transpose, n int, k int, alpha float64, a []float64, lda int, beta float64, c []float64) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dsfrk((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.char)(trans), (C.lapack_int)(n), (C.lapack_int)(k), (C.double)(alpha), (*C.double)(&a[0]), (C.lapack_int)(lda), (C.double)(beta), (*C.double)(&c[0])))
}

func Sspcon(ul blas.Uplo, n int, ap []float32, ipiv []int32, anorm float32, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sspcon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Dspcon(ul blas.Uplo, n int, ap []float64, ipiv []int32, anorm float64, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dspcon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Cspcon(ul blas.Uplo, n int, ap []complex64, ipiv []int32, anorm float32, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cspcon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Zspcon(ul blas.Uplo, n int, ap []complex128, ipiv []int32, anorm float64, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zspcon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Sspev(jobz lapack.Job, ul blas.Uplo, n int, ap []float32, w []float32, z []float32, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sspev((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dspev(jobz lapack.Job, ul blas.Uplo, n int, ap []float64, w []float64, z []float64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dspev((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Sspevd(jobz lapack.Job, ul blas.Uplo, n int, ap []float32, w []float32, z []float32, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sspevd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dspevd(jobz lapack.Job, ul blas.Uplo, n int, ap []float64, w []float64, z []float64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dspevd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Sspevx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, ap []float32, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []float32, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sspevx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Dspevx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, ap []float64, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []float64, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dspevx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Sspgst(itype int, ul blas.Uplo, n int, ap []float32, bp []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sspgst((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.float)(&bp[0])))
}

func Dspgst(itype int, ul blas.Uplo, n int, ap []float64, bp []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dspgst((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.double)(&bp[0])))
}

func Sspgv(itype int, jobz lapack.Job, ul blas.Uplo, n int, ap []float32, bp []float32, w []float32, z []float32, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sspgv((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.float)(&bp[0]), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dspgv(itype int, jobz lapack.Job, ul blas.Uplo, n int, ap []float64, bp []float64, w []float64, z []float64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dspgv((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.double)(&bp[0]), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Sspgvd(itype int, jobz lapack.Job, ul blas.Uplo, n int, ap []float32, bp []float32, w []float32, z []float32, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sspgvd((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.float)(&bp[0]), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dspgvd(itype int, jobz lapack.Job, ul blas.Uplo, n int, ap []float64, bp []float64, w []float64, z []float64, ldz int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dspgvd((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.double)(&bp[0]), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Sspgvx(itype int, jobz lapack.Job, rng byte, ul blas.Uplo, n int, ap []float32, bp []float32, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []float32, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sspgvx((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.float)(&bp[0]), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Dspgvx(itype int, jobz lapack.Job, rng byte, ul blas.Uplo, n int, ap []float64, bp []float64, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []float64, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dspgvx((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.double)(&bp[0]), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Ssprfs(ul blas.Uplo, n int, nrhs int, ap []float32, afp []float32, ipiv []int32, b []float32, ldb int, x []float32, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssprfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&ap[0]), (*C.float)(&afp[0]), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dsprfs(ul blas.Uplo, n int, nrhs int, ap []float64, afp []float64, ipiv []int32, b []float64, ldb int, x []float64, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsprfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&ap[0]), (*C.double)(&afp[0]), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Csprfs(ul blas.Uplo, n int, nrhs int, ap []complex64, afp []complex64, ipiv []int32, b []complex64, ldb int, x []complex64, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csprfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&afp[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zsprfs(ul blas.Uplo, n int, nrhs int, ap []complex128, afp []complex128, ipiv []int32, b []complex128, ldb int, x []complex128, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsprfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&afp[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Sspsv(ul blas.Uplo, n int, nrhs int, ap []float32, ipiv []int32, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sspsv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dspsv(ul blas.Uplo, n int, nrhs int, ap []float64, ipiv []int32, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dspsv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Cspsv(ul blas.Uplo, n int, nrhs int, ap []complex64, ipiv []int32, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cspsv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zspsv(ul blas.Uplo, n int, nrhs int, ap []complex128, ipiv []int32, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zspsv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Sspsvx(fact byte, ul blas.Uplo, n int, nrhs int, ap []float32, afp []float32, ipiv []int32, b []float32, ldb int, x []float32, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_sspsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&ap[0]), (*C.float)(&afp[0]), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dspsvx(fact byte, ul blas.Uplo, n int, nrhs int, ap []float64, afp []float64, ipiv []int32, b []float64, ldb int, x []float64, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dspsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&ap[0]), (*C.double)(&afp[0]), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Cspsvx(fact byte, ul blas.Uplo, n int, nrhs int, ap []complex64, afp []complex64, ipiv []int32, b []complex64, ldb int, x []complex64, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cspsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&afp[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zspsvx(fact byte, ul blas.Uplo, n int, nrhs int, ap []complex128, afp []complex128, ipiv []int32, b []complex128, ldb int, x []complex128, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zspsvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&afp[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Ssptrd(ul blas.Uplo, n int, ap []float32, d []float32, e []float32, tau []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssptrd((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&tau[0])))
}

func Dsptrd(ul blas.Uplo, n int, ap []float64, d []float64, e []float64, tau []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsptrd((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&tau[0])))
}

func Ssptrf(ul blas.Uplo, n int, ap []float32, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssptrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.lapack_int)(&ipiv[0])))
}

func Dsptrf(ul blas.Uplo, n int, ap []float64, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsptrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.lapack_int)(&ipiv[0])))
}

func Csptrf(ul blas.Uplo, n int, ap []complex64, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csptrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_int)(&ipiv[0])))
}

func Zsptrf(ul blas.Uplo, n int, ap []complex128, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsptrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_int)(&ipiv[0])))
}

func Ssptri(ul blas.Uplo, n int, ap []float32, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssptri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.lapack_int)(&ipiv[0])))
}

func Dsptri(ul blas.Uplo, n int, ap []float64, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsptri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.lapack_int)(&ipiv[0])))
}

func Csptri(ul blas.Uplo, n int, ap []complex64, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csptri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_int)(&ipiv[0])))
}

func Zsptri(ul blas.Uplo, n int, ap []complex128, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsptri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_int)(&ipiv[0])))
}

func Ssptrs(ul blas.Uplo, n int, nrhs int, ap []float32, ipiv []int32, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssptrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dsptrs(ul blas.Uplo, n int, nrhs int, ap []float64, ipiv []int32, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsptrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Csptrs(ul blas.Uplo, n int, nrhs int, ap []complex64, ipiv []int32, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csptrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zsptrs(ul blas.Uplo, n int, nrhs int, ap []complex128, ipiv []int32, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsptrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Sstebz(rng byte, order byte, n int, vl float32, vu float32, il int, iu int, abstol float32, d []float32, e []float32, m []int32, nsplit []int32, w []float32, iblock []int32, isplit []int32) bool {
	return isZero(C.LAPACKE_sstebz((C.char)(rng), (C.char)(order), (C.lapack_int)(n), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.lapack_int)(&m[0]), (*C.lapack_int)(&nsplit[0]), (*C.float)(&w[0]), (*C.lapack_int)(&iblock[0]), (*C.lapack_int)(&isplit[0])))
}

func Dstebz(rng byte, order byte, n int, vl float64, vu float64, il int, iu int, abstol float64, d []float64, e []float64, m []int32, nsplit []int32, w []float64, iblock []int32, isplit []int32) bool {
	return isZero(C.LAPACKE_dstebz((C.char)(rng), (C.char)(order), (C.lapack_int)(n), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.lapack_int)(&m[0]), (*C.lapack_int)(&nsplit[0]), (*C.double)(&w[0]), (*C.lapack_int)(&iblock[0]), (*C.lapack_int)(&isplit[0])))
}

func Sstedc(compz lapack.CompSV, n int, d []float32, e []float32, z []float32, ldz int) bool {
	return isZero(C.LAPACKE_sstedc((C.int)(rowMajor), (C.char)(compz), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dstedc(compz lapack.CompSV, n int, d []float64, e []float64, z []float64, ldz int) bool {
	return isZero(C.LAPACKE_dstedc((C.int)(rowMajor), (C.char)(compz), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Cstedc(compz lapack.CompSV, n int, d []float32, e []float32, z []complex64, ldz int) bool {
	return isZero(C.LAPACKE_cstedc((C.int)(rowMajor), (C.char)(compz), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz)))
}

func Zstedc(compz lapack.CompSV, n int, d []float64, e []float64, z []complex128, ldz int) bool {
	return isZero(C.LAPACKE_zstedc((C.int)(rowMajor), (C.char)(compz), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz)))
}

func Sstegr(jobz lapack.Job, rng byte, n int, d []float32, e []float32, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []float32, ldz int, isuppz []int32) bool {
	return isZero(C.LAPACKE_sstegr((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&isuppz[0])))
}

func Dstegr(jobz lapack.Job, rng byte, n int, d []float64, e []float64, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []float64, ldz int, isuppz []int32) bool {
	return isZero(C.LAPACKE_dstegr((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&isuppz[0])))
}

func Cstegr(jobz lapack.Job, rng byte, n int, d []float32, e []float32, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []complex64, ldz int, isuppz []int32) bool {
	return isZero(C.LAPACKE_cstegr((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&isuppz[0])))
}

func Zstegr(jobz lapack.Job, rng byte, n int, d []float64, e []float64, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []complex128, ldz int, isuppz []int32) bool {
	return isZero(C.LAPACKE_zstegr((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&isuppz[0])))
}

func Sstein(n int, d []float32, e []float32, m int, w []float32, iblock []int32, isplit []int32, z []float32, ldz int, ifailv []int32) bool {
	return isZero(C.LAPACKE_sstein((C.int)(rowMajor), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (C.lapack_int)(m), (*C.float)(&w[0]), (*C.lapack_int)(&iblock[0]), (*C.lapack_int)(&isplit[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifailv[0])))
}

func Dstein(n int, d []float64, e []float64, m int, w []float64, iblock []int32, isplit []int32, z []float64, ldz int, ifailv []int32) bool {
	return isZero(C.LAPACKE_dstein((C.int)(rowMajor), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (C.lapack_int)(m), (*C.double)(&w[0]), (*C.lapack_int)(&iblock[0]), (*C.lapack_int)(&isplit[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifailv[0])))
}

func Cstein(n int, d []float32, e []float32, m int, w []float32, iblock []int32, isplit []int32, z []complex64, ldz int, ifailv []int32) bool {
	return isZero(C.LAPACKE_cstein((C.int)(rowMajor), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (C.lapack_int)(m), (*C.float)(&w[0]), (*C.lapack_int)(&iblock[0]), (*C.lapack_int)(&isplit[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifailv[0])))
}

func Zstein(n int, d []float64, e []float64, m int, w []float64, iblock []int32, isplit []int32, z []complex128, ldz int, ifailv []int32) bool {
	return isZero(C.LAPACKE_zstein((C.int)(rowMajor), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (C.lapack_int)(m), (*C.double)(&w[0]), (*C.lapack_int)(&iblock[0]), (*C.lapack_int)(&isplit[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifailv[0])))
}

func Sstemr(jobz lapack.Job, rng byte, n int, d []float32, e []float32, vl float32, vu float32, il int, iu int, m []int32, w []float32, z []float32, ldz int, nzc int, isuppz []int32, tryrac []int32) bool {
	return isZero(C.LAPACKE_sstemr((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz), (C.lapack_int)(nzc), (*C.lapack_int)(&isuppz[0]), (*C.lapack_logical)(&tryrac[0])))
}

func Dstemr(jobz lapack.Job, rng byte, n int, d []float64, e []float64, vl float64, vu float64, il int, iu int, m []int32, w []float64, z []float64, ldz int, nzc int, isuppz []int32, tryrac []int32) bool {
	return isZero(C.LAPACKE_dstemr((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz), (C.lapack_int)(nzc), (*C.lapack_int)(&isuppz[0]), (*C.lapack_logical)(&tryrac[0])))
}

func Cstemr(jobz lapack.Job, rng byte, n int, d []float32, e []float32, vl float32, vu float32, il int, iu int, m []int32, w []float32, z []complex64, ldz int, nzc int, isuppz []int32, tryrac []int32) bool {
	return isZero(C.LAPACKE_cstemr((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz), (C.lapack_int)(nzc), (*C.lapack_int)(&isuppz[0]), (*C.lapack_logical)(&tryrac[0])))
}

func Zstemr(jobz lapack.Job, rng byte, n int, d []float64, e []float64, vl float64, vu float64, il int, iu int, m []int32, w []float64, z []complex128, ldz int, nzc int, isuppz []int32, tryrac []int32) bool {
	return isZero(C.LAPACKE_zstemr((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz), (C.lapack_int)(nzc), (*C.lapack_int)(&isuppz[0]), (*C.lapack_logical)(&tryrac[0])))
}

func Ssteqr(compz lapack.CompSV, n int, d []float32, e []float32, z []float32, ldz int) bool {
	return isZero(C.LAPACKE_ssteqr((C.int)(rowMajor), (C.char)(compz), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dsteqr(compz lapack.CompSV, n int, d []float64, e []float64, z []float64, ldz int) bool {
	return isZero(C.LAPACKE_dsteqr((C.int)(rowMajor), (C.char)(compz), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Csteqr(compz lapack.CompSV, n int, d []float32, e []float32, z []complex64, ldz int) bool {
	return isZero(C.LAPACKE_csteqr((C.int)(rowMajor), (C.char)(compz), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz)))
}

func Zsteqr(compz lapack.CompSV, n int, d []float64, e []float64, z []complex128, ldz int) bool {
	return isZero(C.LAPACKE_zsteqr((C.int)(rowMajor), (C.char)(compz), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz)))
}

func Ssterf(n int, d []float32, e []float32) bool {
	return isZero(C.LAPACKE_ssterf((C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0])))
}

func Dsterf(n int, d []float64, e []float64) bool {
	return isZero(C.LAPACKE_dsterf((C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0])))
}

func Sstev(jobz lapack.Job, n int, d []float32, e []float32, z []float32, ldz int) bool {
	return isZero(C.LAPACKE_sstev((C.int)(rowMajor), (C.char)(jobz), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dstev(jobz lapack.Job, n int, d []float64, e []float64, z []float64, ldz int) bool {
	return isZero(C.LAPACKE_dstev((C.int)(rowMajor), (C.char)(jobz), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Sstevd(jobz lapack.Job, n int, d []float32, e []float32, z []float32, ldz int) bool {
	return isZero(C.LAPACKE_sstevd((C.int)(rowMajor), (C.char)(jobz), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz)))
}

func Dstevd(jobz lapack.Job, n int, d []float64, e []float64, z []float64, ldz int) bool {
	return isZero(C.LAPACKE_dstevd((C.int)(rowMajor), (C.char)(jobz), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz)))
}

func Sstevr(jobz lapack.Job, rng byte, n int, d []float32, e []float32, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []float32, ldz int, isuppz []int32) bool {
	return isZero(C.LAPACKE_sstevr((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&isuppz[0])))
}

func Dstevr(jobz lapack.Job, rng byte, n int, d []float64, e []float64, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []float64, ldz int, isuppz []int32) bool {
	return isZero(C.LAPACKE_dstevr((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&isuppz[0])))
}

func Sstevx(jobz lapack.Job, rng byte, n int, d []float32, e []float32, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []float32, ldz int, ifail []int32) bool {
	return isZero(C.LAPACKE_sstevx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.lapack_int)(n), (*C.float)(&d[0]), (*C.float)(&e[0]), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Dstevx(jobz lapack.Job, rng byte, n int, d []float64, e []float64, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []float64, ldz int, ifail []int32) bool {
	return isZero(C.LAPACKE_dstevx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.lapack_int)(n), (*C.double)(&d[0]), (*C.double)(&e[0]), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Ssycon(ul blas.Uplo, n int, a []float32, lda int, ipiv []int32, anorm float32, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssycon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Dsycon(ul blas.Uplo, n int, a []float64, lda int, ipiv []int32, anorm float64, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsycon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Csycon(ul blas.Uplo, n int, a []complex64, lda int, ipiv []int32, anorm float32, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csycon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (C.float)(anorm), (*C.float)(&rcond[0])))
}

func Zsycon(ul blas.Uplo, n int, a []complex128, lda int, ipiv []int32, anorm float64, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsycon((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (C.double)(anorm), (*C.double)(&rcond[0])))
}

func Ssyequb(ul blas.Uplo, n int, a []float32, lda int, s []float32, scond []float32, amax []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssyequb((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&s[0]), (*C.float)(&scond[0]), (*C.float)(&amax[0])))
}

func Dsyequb(ul blas.Uplo, n int, a []float64, lda int, s []float64, scond []float64, amax []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsyequb((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&s[0]), (*C.double)(&scond[0]), (*C.double)(&amax[0])))
}

func Csyequb(ul blas.Uplo, n int, a []complex64, lda int, s []float32, scond []float32, amax []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csyequb((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&s[0]), (*C.float)(&scond[0]), (*C.float)(&amax[0])))
}

func Zsyequb(ul blas.Uplo, n int, a []complex128, lda int, s []float64, scond []float64, amax []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsyequb((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&s[0]), (*C.double)(&scond[0]), (*C.double)(&amax[0])))
}

func Ssyev(jobz lapack.Job, ul blas.Uplo, n int, a []float32, lda int, w []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssyev((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&w[0])))
}

func Dsyev(jobz lapack.Job, ul blas.Uplo, n int, a []float64, lda int, w []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsyev((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&w[0])))
}

func Ssyevd(jobz lapack.Job, ul blas.Uplo, n int, a []float32, lda int, w []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssyevd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&w[0])))
}

func Dsyevd(jobz lapack.Job, ul blas.Uplo, n int, a []float64, lda int, w []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsyevd((C.int)(rowMajor), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&w[0])))
}

func Ssyevr(jobz lapack.Job, rng byte, ul blas.Uplo, n int, a []float32, lda int, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []float32, ldz int, isuppz []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssyevr((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&isuppz[0])))
}

func Dsyevr(jobz lapack.Job, rng byte, ul blas.Uplo, n int, a []float64, lda int, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []float64, ldz int, isuppz []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsyevr((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&isuppz[0])))
}

func Ssyevx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, a []float32, lda int, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []float32, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssyevx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Dsyevx(jobz lapack.Job, rng byte, ul blas.Uplo, n int, a []float64, lda int, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []float64, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsyevx((C.int)(rowMajor), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Ssygst(itype int, ul blas.Uplo, n int, a []float32, lda int, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssygst((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dsygst(itype int, ul blas.Uplo, n int, a []float64, lda int, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsygst((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Ssygv(itype int, jobz lapack.Job, ul blas.Uplo, n int, a []float32, lda int, b []float32, ldb int, w []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssygv((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&w[0])))
}

func Dsygv(itype int, jobz lapack.Job, ul blas.Uplo, n int, a []float64, lda int, b []float64, ldb int, w []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsygv((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&w[0])))
}

func Ssygvd(itype int, jobz lapack.Job, ul blas.Uplo, n int, a []float32, lda int, b []float32, ldb int, w []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssygvd((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&w[0])))
}

func Dsygvd(itype int, jobz lapack.Job, ul blas.Uplo, n int, a []float64, lda int, b []float64, ldb int, w []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsygvd((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&w[0])))
}

func Ssygvx(itype int, jobz lapack.Job, rng byte, ul blas.Uplo, n int, a []float32, lda int, b []float32, ldb int, vl float32, vu float32, il int, iu int, abstol float32, m []int32, w []float32, z []float32, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssygvx((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (C.float)(vl), (C.float)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.float)(abstol), (*C.lapack_int)(&m[0]), (*C.float)(&w[0]), (*C.float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Dsygvx(itype int, jobz lapack.Job, rng byte, ul blas.Uplo, n int, a []float64, lda int, b []float64, ldb int, vl float64, vu float64, il int, iu int, abstol float64, m []int32, w []float64, z []float64, ldz int, ifail []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsygvx((C.int)(rowMajor), (C.lapack_int)(itype), (C.char)(jobz), (C.char)(rng), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (C.double)(vl), (C.double)(vu), (C.lapack_int)(il), (C.lapack_int)(iu), (C.double)(abstol), (*C.lapack_int)(&m[0]), (*C.double)(&w[0]), (*C.double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifail[0])))
}

func Ssyrfs(ul blas.Uplo, n int, nrhs int, a []float32, lda int, af []float32, ldaf int, ipiv []int32, b []float32, ldb int, x []float32, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssyrfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dsyrfs(ul blas.Uplo, n int, nrhs int, a []float64, lda int, af []float64, ldaf int, ipiv []int32, b []float64, ldb int, x []float64, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsyrfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Csyrfs(ul blas.Uplo, n int, nrhs int, a []complex64, lda int, af []complex64, ldaf int, ipiv []int32, b []complex64, ldb int, x []complex64, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csyrfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zsyrfs(ul blas.Uplo, n int, nrhs int, a []complex128, lda int, af []complex128, ldaf int, ipiv []int32, b []complex128, ldb int, x []complex128, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsyrfs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Ssysv(ul blas.Uplo, n int, nrhs int, a []float32, lda int, ipiv []int32, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssysv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dsysv(ul blas.Uplo, n int, nrhs int, a []float64, lda int, ipiv []int32, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsysv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Csysv(ul blas.Uplo, n int, nrhs int, a []complex64, lda int, ipiv []int32, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csysv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zsysv(ul blas.Uplo, n int, nrhs int, a []complex128, lda int, ipiv []int32, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsysv((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Ssysvx(fact byte, ul blas.Uplo, n int, nrhs int, a []float32, lda int, af []float32, ldaf int, ipiv []int32, b []float32, ldb int, x []float32, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssysvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dsysvx(fact byte, ul blas.Uplo, n int, nrhs int, a []float64, lda int, af []float64, ldaf int, ipiv []int32, b []float64, ldb int, x []float64, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsysvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Csysvx(fact byte, ul blas.Uplo, n int, nrhs int, a []complex64, lda int, af []complex64, ldaf int, ipiv []int32, b []complex64, ldb int, x []complex64, ldx int, rcond []float32, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csysvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&rcond[0]), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Zsysvx(fact byte, ul blas.Uplo, n int, nrhs int, a []complex128, lda int, af []complex128, ldaf int, ipiv []int32, b []complex128, ldb int, x []complex128, ldx int, rcond []float64, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsysvx((C.int)(rowMajor), (C.char)(fact), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&af[0]), (C.lapack_int)(ldaf), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&rcond[0]), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Ssytrd(ul blas.Uplo, n int, a []float32, lda int, d []float32, e []float32, tau []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssytrd((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&d[0]), (*C.float)(&e[0]), (*C.float)(&tau[0])))
}

func Dsytrd(ul blas.Uplo, n int, a []float64, lda int, d []float64, e []float64, tau []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsytrd((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&d[0]), (*C.double)(&e[0]), (*C.double)(&tau[0])))
}

func Ssytrf(ul blas.Uplo, n int, a []float32, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssytrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Dsytrf(ul blas.Uplo, n int, a []float64, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsytrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Csytrf(ul blas.Uplo, n int, a []complex64, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csytrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Zsytrf(ul blas.Uplo, n int, a []complex128, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsytrf((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Ssytri(ul blas.Uplo, n int, a []float32, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssytri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Dsytri(ul blas.Uplo, n int, a []float64, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsytri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Csytri(ul blas.Uplo, n int, a []complex64, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csytri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Zsytri(ul blas.Uplo, n int, a []complex128, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsytri((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Ssytrs(ul blas.Uplo, n int, nrhs int, a []float32, lda int, ipiv []int32, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssytrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dsytrs(ul blas.Uplo, n int, nrhs int, a []float64, lda int, ipiv []int32, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsytrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Csytrs(ul blas.Uplo, n int, nrhs int, a []complex64, lda int, ipiv []int32, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csytrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zsytrs(ul blas.Uplo, n int, nrhs int, a []complex128, lda int, ipiv []int32, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsytrs((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Stbcon(norm byte, ul blas.Uplo, d blas.Diag, n int, kd int, ab []float32, ldab int, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_stbcon((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&rcond[0])))
}

func Dtbcon(norm byte, ul blas.Uplo, d blas.Diag, n int, kd int, ab []float64, ldab int, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_dtbcon((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&rcond[0])))
}

func Ctbcon(norm byte, ul blas.Uplo, d blas.Diag, n int, kd int, ab []complex64, ldab int, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ctbcon((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&rcond[0])))
}

func Ztbcon(norm byte, ul blas.Uplo, d blas.Diag, n int, kd int, ab []complex128, ldab int, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ztbcon((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(kd), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&rcond[0])))
}

func Stbrfs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, kd int, nrhs int, ab []float32, ldab int, b []float32, ldb int, x []float32, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_stbrfs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dtbrfs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, kd int, nrhs int, ab []float64, ldab int, b []float64, ldb int, x []float64, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_dtbrfs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Ctbrfs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, kd int, nrhs int, ab []complex64, ldab int, b []complex64, ldb int, x []complex64, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ctbrfs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Ztbrfs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, kd int, nrhs int, ab []complex128, ldab int, b []complex128, ldb int, x []complex128, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ztbrfs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Stbtrs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, kd int, nrhs int, ab []float32, ldab int, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_stbtrs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.float)(&ab[0]), (C.lapack_int)(ldab), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dtbtrs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, kd int, nrhs int, ab []float64, ldab int, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_dtbtrs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.double)(&ab[0]), (C.lapack_int)(ldab), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Ctbtrs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, kd int, nrhs int, ab []complex64, ldab int, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ctbtrs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Ztbtrs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, kd int, nrhs int, ab []complex128, ldab int, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ztbtrs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(kd), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ab[0]), (C.lapack_int)(ldab), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Stfsm(transr blas.Transpose, s blas.Side, ul blas.Uplo, trans blas.Transpose, d blas.Diag, m int, n int, alpha float32, a []float32, b []float32, ldb int) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_stfsm((C.int)(rowMajor), (C.char)(transr), (C.char)(s), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(m), (C.lapack_int)(n), (C.float)(alpha), (*C.float)(&a[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dtfsm(transr blas.Transpose, s blas.Side, ul blas.Uplo, trans blas.Transpose, d blas.Diag, m int, n int, alpha float64, a []float64, b []float64, ldb int) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_dtfsm((C.int)(rowMajor), (C.char)(transr), (C.char)(s), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(m), (C.lapack_int)(n), (C.double)(alpha), (*C.double)(&a[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Ctfsm(transr blas.Transpose, s blas.Side, ul blas.Uplo, trans blas.Transpose, d blas.Diag, m int, n int, alpha complex64, a []complex64, b []complex64, ldb int) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ctfsm((C.int)(rowMajor), (C.char)(transr), (C.char)(s), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_complex_float)(alpha), (*C.lapack_complex_float)(&a[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Ztfsm(transr blas.Transpose, s blas.Side, ul blas.Uplo, trans blas.Transpose, d blas.Diag, m int, n int, alpha complex128, a []complex128, b []complex128, ldb int) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ztfsm((C.int)(rowMajor), (C.char)(transr), (C.char)(s), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_complex_double)(alpha), (*C.lapack_complex_double)(&a[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Stftri(transr blas.Transpose, ul blas.Uplo, d blas.Diag, n int, a []float32) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_stftri((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.float)(&a[0])))
}

func Dtftri(transr blas.Transpose, ul blas.Uplo, d blas.Diag, n int, a []float64) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_dtftri((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.double)(&a[0])))
}

func Ctftri(transr blas.Transpose, ul blas.Uplo, d blas.Diag, n int, a []complex64) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ctftri((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0])))
}

func Ztftri(transr blas.Transpose, ul blas.Uplo, d blas.Diag, n int, a []complex128) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ztftri((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0])))
}

func Stfttp(transr blas.Transpose, ul blas.Uplo, n int, arf []float32, ap []float32) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_stfttp((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&arf[0]), (*C.float)(&ap[0])))
}

func Dtfttp(transr blas.Transpose, ul blas.Uplo, n int, arf []float64, ap []float64) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dtfttp((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&arf[0]), (*C.double)(&ap[0])))
}

func Ctfttp(transr blas.Transpose, ul blas.Uplo, n int, arf []complex64, ap []complex64) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ctfttp((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&arf[0]), (*C.lapack_complex_float)(&ap[0])))
}

func Ztfttp(transr blas.Transpose, ul blas.Uplo, n int, arf []complex128, ap []complex128) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ztfttp((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&arf[0]), (*C.lapack_complex_double)(&ap[0])))
}

func Stfttr(transr blas.Transpose, ul blas.Uplo, n int, arf []float32, a []float32, lda int) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_stfttr((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&arf[0]), (*C.float)(&a[0]), (C.lapack_int)(lda)))
}

func Dtfttr(transr blas.Transpose, ul blas.Uplo, n int, arf []float64, a []float64, lda int) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dtfttr((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&arf[0]), (*C.double)(&a[0]), (C.lapack_int)(lda)))
}

func Ctfttr(transr blas.Transpose, ul blas.Uplo, n int, arf []complex64, a []complex64, lda int) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ctfttr((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&arf[0]), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda)))
}

func Ztfttr(transr blas.Transpose, ul blas.Uplo, n int, arf []complex128, a []complex128, lda int) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ztfttr((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&arf[0]), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda)))
}

func Stgexc(wantq int32, wantz int32, n int, a []float32, lda int, b []float32, ldb int, q []float32, ldq int, z []float32, ldz int, ifst []int32, ilst []int32) bool {
	return isZero(C.LAPACKE_stgexc((C.int)(rowMajor), (C.lapack_logical)(wantq), (C.lapack_logical)(wantz), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&q[0]), (C.lapack_int)(ldq), (*C.float)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifst[0]), (*C.lapack_int)(&ilst[0])))
}

func Dtgexc(wantq int32, wantz int32, n int, a []float64, lda int, b []float64, ldb int, q []float64, ldq int, z []float64, ldz int, ifst []int32, ilst []int32) bool {
	return isZero(C.LAPACKE_dtgexc((C.int)(rowMajor), (C.lapack_logical)(wantq), (C.lapack_logical)(wantz), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&q[0]), (C.lapack_int)(ldq), (*C.double)(&z[0]), (C.lapack_int)(ldz), (*C.lapack_int)(&ifst[0]), (*C.lapack_int)(&ilst[0])))
}

func Ctgexc(wantq int32, wantz int32, n int, a []complex64, lda int, b []complex64, ldb int, q []complex64, ldq int, z []complex64, ldz int, ifst int, ilst int) bool {
	return isZero(C.LAPACKE_ctgexc((C.int)(rowMajor), (C.lapack_logical)(wantq), (C.lapack_logical)(wantz), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_complex_float)(&z[0]), (C.lapack_int)(ldz), (C.lapack_int)(ifst), (C.lapack_int)(ilst)))
}

func Ztgexc(wantq int32, wantz int32, n int, a []complex128, lda int, b []complex128, ldb int, q []complex128, ldq int, z []complex128, ldz int, ifst int, ilst int) bool {
	return isZero(C.LAPACKE_ztgexc((C.int)(rowMajor), (C.lapack_logical)(wantq), (C.lapack_logical)(wantz), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_complex_double)(&z[0]), (C.lapack_int)(ldz), (C.lapack_int)(ifst), (C.lapack_int)(ilst)))
}

func Stgsja(jobu lapack.Job, jobv lapack.Job, jobq lapack.Job, m int, p int, n int, k int, l int, a []float32, lda int, b []float32, ldb int, tola float32, tolb float32, alpha []float32, beta []float32, u []float32, ldu int, v []float32, ldv int, q []float32, ldq int, ncycle []int32) bool {
	return isZero(C.LAPACKE_stgsja((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobv), (C.char)(jobq), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (C.float)(tola), (C.float)(tolb), (*C.float)(&alpha[0]), (*C.float)(&beta[0]), (*C.float)(&u[0]), (C.lapack_int)(ldu), (*C.float)(&v[0]), (C.lapack_int)(ldv), (*C.float)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_int)(&ncycle[0])))
}

func Dtgsja(jobu lapack.Job, jobv lapack.Job, jobq lapack.Job, m int, p int, n int, k int, l int, a []float64, lda int, b []float64, ldb int, tola float64, tolb float64, alpha []float64, beta []float64, u []float64, ldu int, v []float64, ldv int, q []float64, ldq int, ncycle []int32) bool {
	return isZero(C.LAPACKE_dtgsja((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobv), (C.char)(jobq), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (C.double)(tola), (C.double)(tolb), (*C.double)(&alpha[0]), (*C.double)(&beta[0]), (*C.double)(&u[0]), (C.lapack_int)(ldu), (*C.double)(&v[0]), (C.lapack_int)(ldv), (*C.double)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_int)(&ncycle[0])))
}

func Ctgsja(jobu lapack.Job, jobv lapack.Job, jobq lapack.Job, m int, p int, n int, k int, l int, a []complex64, lda int, b []complex64, ldb int, tola float32, tolb float32, alpha []float32, beta []float32, u []complex64, ldu int, v []complex64, ldv int, q []complex64, ldq int, ncycle []int32) bool {
	return isZero(C.LAPACKE_ctgsja((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobv), (C.char)(jobq), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (C.float)(tola), (C.float)(tolb), (*C.float)(&alpha[0]), (*C.float)(&beta[0]), (*C.lapack_complex_float)(&u[0]), (C.lapack_int)(ldu), (*C.lapack_complex_float)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_float)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_int)(&ncycle[0])))
}

func Ztgsja(jobu lapack.Job, jobv lapack.Job, jobq lapack.Job, m int, p int, n int, k int, l int, a []complex128, lda int, b []complex128, ldb int, tola float64, tolb float64, alpha []float64, beta []float64, u []complex128, ldu int, v []complex128, ldv int, q []complex128, ldq int, ncycle []int32) bool {
	return isZero(C.LAPACKE_ztgsja((C.int)(rowMajor), (C.char)(jobu), (C.char)(jobv), (C.char)(jobq), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (C.double)(tola), (C.double)(tolb), (*C.double)(&alpha[0]), (*C.double)(&beta[0]), (*C.lapack_complex_double)(&u[0]), (C.lapack_int)(ldu), (*C.lapack_complex_double)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_double)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_int)(&ncycle[0])))
}

func Stgsyl(trans blas.Transpose, ijob lapack.Job, m int, n int, a []float32, lda int, b []float32, ldb int, c []float32, ldc int, d []float32, ldd int, e []float32, lde int, f []float32, ldf int, scale []float32, dif []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_stgsyl((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(ijob), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&c[0]), (C.lapack_int)(ldc), (*C.float)(&d[0]), (C.lapack_int)(ldd), (*C.float)(&e[0]), (C.lapack_int)(lde), (*C.float)(&f[0]), (C.lapack_int)(ldf), (*C.float)(&scale[0]), (*C.float)(&dif[0])))
}

func Dtgsyl(trans blas.Transpose, ijob lapack.Job, m int, n int, a []float64, lda int, b []float64, ldb int, c []float64, ldc int, d []float64, ldd int, e []float64, lde int, f []float64, ldf int, scale []float64, dif []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dtgsyl((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(ijob), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&c[0]), (C.lapack_int)(ldc), (*C.double)(&d[0]), (C.lapack_int)(ldd), (*C.double)(&e[0]), (C.lapack_int)(lde), (*C.double)(&f[0]), (C.lapack_int)(ldf), (*C.double)(&scale[0]), (*C.double)(&dif[0])))
}

func Ctgsyl(trans blas.Transpose, ijob lapack.Job, m int, n int, a []complex64, lda int, b []complex64, ldb int, c []complex64, ldc int, d []complex64, ldd int, e []complex64, lde int, f []complex64, ldf int, scale []float32, dif []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_ctgsyl((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(ijob), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc), (*C.lapack_complex_float)(&d[0]), (C.lapack_int)(ldd), (*C.lapack_complex_float)(&e[0]), (C.lapack_int)(lde), (*C.lapack_complex_float)(&f[0]), (C.lapack_int)(ldf), (*C.float)(&scale[0]), (*C.float)(&dif[0])))
}

func Ztgsyl(trans blas.Transpose, ijob lapack.Job, m int, n int, a []complex128, lda int, b []complex128, ldb int, c []complex128, ldc int, d []complex128, ldd int, e []complex128, lde int, f []complex128, ldf int, scale []float64, dif []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_ztgsyl((C.int)(rowMajor), (C.char)(trans), (C.lapack_int)(ijob), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc), (*C.lapack_complex_double)(&d[0]), (C.lapack_int)(ldd), (*C.lapack_complex_double)(&e[0]), (C.lapack_int)(lde), (*C.lapack_complex_double)(&f[0]), (C.lapack_int)(ldf), (*C.double)(&scale[0]), (*C.double)(&dif[0])))
}

func Stpcon(norm byte, ul blas.Uplo, d blas.Diag, n int, ap []float32, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_stpcon((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.float)(&rcond[0])))
}

func Dtpcon(norm byte, ul blas.Uplo, d blas.Diag, n int, ap []float64, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_dtpcon((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.double)(&rcond[0])))
}

func Ctpcon(norm byte, ul blas.Uplo, d blas.Diag, n int, ap []complex64, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ctpcon((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.float)(&rcond[0])))
}

func Ztpcon(norm byte, ul blas.Uplo, d blas.Diag, n int, ap []complex128, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ztpcon((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.double)(&rcond[0])))
}

func Stprfs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, ap []float32, b []float32, ldb int, x []float32, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_stprfs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&ap[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dtprfs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, ap []float64, b []float64, ldb int, x []float64, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_dtprfs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&ap[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Ctprfs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, ap []complex64, b []complex64, ldb int, x []complex64, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ctprfs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Ztprfs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, ap []complex128, b []complex128, ldb int, x []complex128, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ztprfs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Stptri(ul blas.Uplo, d blas.Diag, n int, ap []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_stptri((C.int)(rowMajor), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.float)(&ap[0])))
}

func Dtptri(ul blas.Uplo, d blas.Diag, n int, ap []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_dtptri((C.int)(rowMajor), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.double)(&ap[0])))
}

func Ctptri(ul blas.Uplo, d blas.Diag, n int, ap []complex64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ctptri((C.int)(rowMajor), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0])))
}

func Ztptri(ul blas.Uplo, d blas.Diag, n int, ap []complex128) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ztptri((C.int)(rowMajor), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0])))
}

func Stptrs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, ap []float32, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_stptrs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&ap[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dtptrs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, ap []float64, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_dtptrs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&ap[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Ctptrs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, ap []complex64, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ctptrs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Ztptrs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, ap []complex128, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ztptrs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Stpttf(transr blas.Transpose, ul blas.Uplo, n int, ap []float32, arf []float32) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_stpttf((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.float)(&arf[0])))
}

func Dtpttf(transr blas.Transpose, ul blas.Uplo, n int, ap []float64, arf []float64) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dtpttf((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.double)(&arf[0])))
}

func Ctpttf(transr blas.Transpose, ul blas.Uplo, n int, ap []complex64, arf []complex64) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ctpttf((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&arf[0])))
}

func Ztpttf(transr blas.Transpose, ul blas.Uplo, n int, ap []complex128, arf []complex128) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ztpttf((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&arf[0])))
}

func Stpttr(ul blas.Uplo, n int, ap []float32, a []float32, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_stpttr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&ap[0]), (*C.float)(&a[0]), (C.lapack_int)(lda)))
}

func Dtpttr(ul blas.Uplo, n int, ap []float64, a []float64, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dtpttr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&ap[0]), (*C.double)(&a[0]), (C.lapack_int)(lda)))
}

func Ctpttr(ul blas.Uplo, n int, ap []complex64, a []complex64, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ctpttr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda)))
}

func Ztpttr(ul blas.Uplo, n int, ap []complex128, a []complex128, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ztpttr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda)))
}

func Strcon(norm byte, ul blas.Uplo, d blas.Diag, n int, a []float32, lda int, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_strcon((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&rcond[0])))
}

func Dtrcon(norm byte, ul blas.Uplo, d blas.Diag, n int, a []float64, lda int, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_dtrcon((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&rcond[0])))
}

func Ctrcon(norm byte, ul blas.Uplo, d blas.Diag, n int, a []complex64, lda int, rcond []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ctrcon((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&rcond[0])))
}

func Ztrcon(norm byte, ul blas.Uplo, d blas.Diag, n int, a []complex128, lda int, rcond []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ztrcon((C.int)(rowMajor), (C.char)(norm), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&rcond[0])))
}

func Strexc(compq lapack.CompSV, n int, t []float32, ldt int, q []float32, ldq int, ifst []int32, ilst []int32) bool {
	return isZero(C.LAPACKE_strexc((C.int)(rowMajor), (C.char)(compq), (C.lapack_int)(n), (*C.float)(&t[0]), (C.lapack_int)(ldt), (*C.float)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_int)(&ifst[0]), (*C.lapack_int)(&ilst[0])))
}

func Dtrexc(compq lapack.CompSV, n int, t []float64, ldt int, q []float64, ldq int, ifst []int32, ilst []int32) bool {
	return isZero(C.LAPACKE_dtrexc((C.int)(rowMajor), (C.char)(compq), (C.lapack_int)(n), (*C.double)(&t[0]), (C.lapack_int)(ldt), (*C.double)(&q[0]), (C.lapack_int)(ldq), (*C.lapack_int)(&ifst[0]), (*C.lapack_int)(&ilst[0])))
}

func Ctrexc(compq lapack.CompSV, n int, t []complex64, ldt int, q []complex64, ldq int, ifst int, ilst int) bool {
	return isZero(C.LAPACKE_ctrexc((C.int)(rowMajor), (C.char)(compq), (C.lapack_int)(n), (*C.lapack_complex_float)(&t[0]), (C.lapack_int)(ldt), (*C.lapack_complex_float)(&q[0]), (C.lapack_int)(ldq), (C.lapack_int)(ifst), (C.lapack_int)(ilst)))
}

func Ztrexc(compq lapack.CompSV, n int, t []complex128, ldt int, q []complex128, ldq int, ifst int, ilst int) bool {
	return isZero(C.LAPACKE_ztrexc((C.int)(rowMajor), (C.char)(compq), (C.lapack_int)(n), (*C.lapack_complex_double)(&t[0]), (C.lapack_int)(ldt), (*C.lapack_complex_double)(&q[0]), (C.lapack_int)(ldq), (C.lapack_int)(ifst), (C.lapack_int)(ilst)))
}

func Strrfs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, a []float32, lda int, b []float32, ldb int, x []float32, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_strrfs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Dtrrfs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, a []float64, lda int, b []float64, ldb int, x []float64, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_dtrrfs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Ctrrfs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, a []complex64, lda int, b []complex64, ldb int, x []complex64, ldx int, ferr []float32, berr []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ctrrfs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.float)(&ferr[0]), (*C.float)(&berr[0])))
}

func Ztrrfs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, a []complex128, lda int, b []complex128, ldb int, x []complex128, ldx int, ferr []float64, berr []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ztrrfs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.double)(&ferr[0]), (*C.double)(&berr[0])))
}

func Strsyl(trana byte, tranb byte, isgn int, m int, n int, a []float32, lda int, b []float32, ldb int, c []float32, ldc int, scale []float32) bool {
	return isZero(C.LAPACKE_strsyl((C.int)(rowMajor), (C.char)(trana), (C.char)(tranb), (C.lapack_int)(isgn), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&c[0]), (C.lapack_int)(ldc), (*C.float)(&scale[0])))
}

func Dtrsyl(trana byte, tranb byte, isgn int, m int, n int, a []float64, lda int, b []float64, ldb int, c []float64, ldc int, scale []float64) bool {
	return isZero(C.LAPACKE_dtrsyl((C.int)(rowMajor), (C.char)(trana), (C.char)(tranb), (C.lapack_int)(isgn), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&c[0]), (C.lapack_int)(ldc), (*C.double)(&scale[0])))
}

func Ctrsyl(trana byte, tranb byte, isgn int, m int, n int, a []complex64, lda int, b []complex64, ldb int, c []complex64, ldc int, scale []float32) bool {
	return isZero(C.LAPACKE_ctrsyl((C.int)(rowMajor), (C.char)(trana), (C.char)(tranb), (C.lapack_int)(isgn), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc), (*C.float)(&scale[0])))
}

func Ztrsyl(trana byte, tranb byte, isgn int, m int, n int, a []complex128, lda int, b []complex128, ldb int, c []complex128, ldc int, scale []float64) bool {
	return isZero(C.LAPACKE_ztrsyl((C.int)(rowMajor), (C.char)(trana), (C.char)(tranb), (C.lapack_int)(isgn), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc), (*C.double)(&scale[0])))
}

func Strtri(ul blas.Uplo, d blas.Diag, n int, a []float32, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_strtri((C.int)(rowMajor), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda)))
}

func Dtrtri(ul blas.Uplo, d blas.Diag, n int, a []float64, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_dtrtri((C.int)(rowMajor), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda)))
}

func Ctrtri(ul blas.Uplo, d blas.Diag, n int, a []complex64, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ctrtri((C.int)(rowMajor), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda)))
}

func Ztrtri(ul blas.Uplo, d blas.Diag, n int, a []complex128, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ztrtri((C.int)(rowMajor), (C.char)(ul), (C.char)(d), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda)))
}

func Strtrs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, a []float32, lda int, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_strtrs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dtrtrs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, a []float64, lda int, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_dtrtrs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Ctrtrs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, a []complex64, lda int, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ctrtrs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Ztrtrs(ul blas.Uplo, trans blas.Transpose, d blas.Diag, n int, nrhs int, a []complex128, lda int, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch d {
	case blas.Unit:
		d = 'U'
	case blas.NonUnit:
		d = 'N'
	default:
		panic("lapack: illegal diagonal")
	}
	return isZero(C.LAPACKE_ztrtrs((C.int)(rowMajor), (C.char)(ul), (C.char)(trans), (C.char)(d), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Strttf(transr blas.Transpose, ul blas.Uplo, n int, a []float32, lda int, arf []float32) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_strttf((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&arf[0])))
}

func Dtrttf(transr blas.Transpose, ul blas.Uplo, n int, a []float64, lda int, arf []float64) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dtrttf((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&arf[0])))
}

func Ctrttf(transr blas.Transpose, ul blas.Uplo, n int, a []complex64, lda int, arf []complex64) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ctrttf((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&arf[0])))
}

func Ztrttf(transr blas.Transpose, ul blas.Uplo, n int, a []complex128, lda int, arf []complex128) bool {
	switch transr {
	case blas.NoTrans:
		transr = 'N'
	case blas.Trans:
		transr = 'T'
	case blas.ConjTrans:
		transr = 'C'
	default:
		panic("lapack: bad trans")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ztrttf((C.int)(rowMajor), (C.char)(transr), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&arf[0])))
}

func Strttp(ul blas.Uplo, n int, a []float32, lda int, ap []float32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_strttp((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&ap[0])))
}

func Dtrttp(ul blas.Uplo, n int, a []float64, lda int, ap []float64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dtrttp((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&ap[0])))
}

func Ctrttp(ul blas.Uplo, n int, a []complex64, lda int, ap []complex64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ctrttp((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&ap[0])))
}

func Ztrttp(ul blas.Uplo, n int, a []complex128, lda int, ap []complex128) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ztrttp((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&ap[0])))
}

func Stzrzf(m int, n int, a []float32, lda int, tau []float32) bool {
	return isZero(C.LAPACKE_stzrzf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&tau[0])))
}

func Dtzrzf(m int, n int, a []float64, lda int, tau []float64) bool {
	return isZero(C.LAPACKE_dtzrzf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&tau[0])))
}

func Ctzrzf(m int, n int, a []complex64, lda int, tau []complex64) bool {
	return isZero(C.LAPACKE_ctzrzf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Ztzrzf(m int, n int, a []complex128, lda int, tau []complex128) bool {
	return isZero(C.LAPACKE_ztzrzf((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Cungbr(vect byte, m int, n int, k int, a []complex64, lda int, tau []complex64) bool {
	return isZero(C.LAPACKE_cungbr((C.int)(rowMajor), (C.char)(vect), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Zungbr(vect byte, m int, n int, k int, a []complex128, lda int, tau []complex128) bool {
	return isZero(C.LAPACKE_zungbr((C.int)(rowMajor), (C.char)(vect), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Cunghr(n int, ilo int, ihi int, a []complex64, lda int, tau []complex64) bool {
	return isZero(C.LAPACKE_cunghr((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Zunghr(n int, ilo int, ihi int, a []complex128, lda int, tau []complex128) bool {
	return isZero(C.LAPACKE_zunghr((C.int)(rowMajor), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Cunglq(m int, n int, k int, a []complex64, lda int, tau []complex64) bool {
	return isZero(C.LAPACKE_cunglq((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Zunglq(m int, n int, k int, a []complex128, lda int, tau []complex128) bool {
	return isZero(C.LAPACKE_zunglq((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Cungql(m int, n int, k int, a []complex64, lda int, tau []complex64) bool {
	return isZero(C.LAPACKE_cungql((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Zungql(m int, n int, k int, a []complex128, lda int, tau []complex128) bool {
	return isZero(C.LAPACKE_zungql((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Cungqr(m int, n int, k int, a []complex64, lda int, tau []complex64) bool {
	return isZero(C.LAPACKE_cungqr((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Zungqr(m int, n int, k int, a []complex128, lda int, tau []complex128) bool {
	return isZero(C.LAPACKE_zungqr((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Cungrq(m int, n int, k int, a []complex64, lda int, tau []complex64) bool {
	return isZero(C.LAPACKE_cungrq((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Zungrq(m int, n int, k int, a []complex128, lda int, tau []complex128) bool {
	return isZero(C.LAPACKE_zungrq((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Cungtr(ul blas.Uplo, n int, a []complex64, lda int, tau []complex64) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cungtr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0])))
}

func Zungtr(ul blas.Uplo, n int, a []complex128, lda int, tau []complex128) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zungtr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0])))
}

func Cunmbr(vect byte, s blas.Side, trans blas.Transpose, m int, n int, k int, a []complex64, lda int, tau []complex64, c []complex64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cunmbr((C.int)(rowMajor), (C.char)(vect), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0]), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc)))
}

func Zunmbr(vect byte, s blas.Side, trans blas.Transpose, m int, n int, k int, a []complex128, lda int, tau []complex128, c []complex128, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zunmbr((C.int)(rowMajor), (C.char)(vect), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0]), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc)))
}

func Cunmhr(s blas.Side, trans blas.Transpose, m int, n int, ilo int, ihi int, a []complex64, lda int, tau []complex64, c []complex64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cunmhr((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0]), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc)))
}

func Zunmhr(s blas.Side, trans blas.Transpose, m int, n int, ilo int, ihi int, a []complex128, lda int, tau []complex128, c []complex128, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zunmhr((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(ilo), (C.lapack_int)(ihi), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0]), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc)))
}

func Cunmlq(s blas.Side, trans blas.Transpose, m int, n int, k int, a []complex64, lda int, tau []complex64, c []complex64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cunmlq((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0]), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc)))
}

func Zunmlq(s blas.Side, trans blas.Transpose, m int, n int, k int, a []complex128, lda int, tau []complex128, c []complex128, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zunmlq((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0]), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc)))
}

func Cunmql(s blas.Side, trans blas.Transpose, m int, n int, k int, a []complex64, lda int, tau []complex64, c []complex64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cunmql((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0]), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc)))
}

func Zunmql(s blas.Side, trans blas.Transpose, m int, n int, k int, a []complex128, lda int, tau []complex128, c []complex128, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zunmql((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0]), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc)))
}

func Cunmqr(s blas.Side, trans blas.Transpose, m int, n int, k int, a []complex64, lda int, tau []complex64, c []complex64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cunmqr((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0]), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc)))
}

func Zunmqr(s blas.Side, trans blas.Transpose, m int, n int, k int, a []complex128, lda int, tau []complex128, c []complex128, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zunmqr((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0]), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc)))
}

func Cunmrq(s blas.Side, trans blas.Transpose, m int, n int, k int, a []complex64, lda int, tau []complex64, c []complex64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cunmrq((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0]), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc)))
}

func Zunmrq(s blas.Side, trans blas.Transpose, m int, n int, k int, a []complex128, lda int, tau []complex128, c []complex128, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zunmrq((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0]), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc)))
}

func Cunmrz(s blas.Side, trans blas.Transpose, m int, n int, k int, l int, a []complex64, lda int, tau []complex64, c []complex64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cunmrz((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0]), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc)))
}

func Zunmrz(s blas.Side, trans blas.Transpose, m int, n int, k int, l int, a []complex128, lda int, tau []complex128, c []complex128, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zunmrz((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0]), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc)))
}

func Cunmtr(s blas.Side, ul blas.Uplo, trans blas.Transpose, m int, n int, a []complex64, lda int, tau []complex64, c []complex64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cunmtr((C.int)(rowMajor), (C.char)(s), (C.char)(ul), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&tau[0]), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc)))
}

func Zunmtr(s blas.Side, ul blas.Uplo, trans blas.Transpose, m int, n int, a []complex128, lda int, tau []complex128, c []complex128, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zunmtr((C.int)(rowMajor), (C.char)(s), (C.char)(ul), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&tau[0]), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc)))
}

func Cupgtr(ul blas.Uplo, n int, ap []complex64, tau []complex64, q []complex64, ldq int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cupgtr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&tau[0]), (*C.lapack_complex_float)(&q[0]), (C.lapack_int)(ldq)))
}

func Zupgtr(ul blas.Uplo, n int, ap []complex128, tau []complex128, q []complex128, ldq int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zupgtr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&tau[0]), (*C.lapack_complex_double)(&q[0]), (C.lapack_int)(ldq)))
}

func Cupmtr(s blas.Side, ul blas.Uplo, trans blas.Transpose, m int, n int, ap []complex64, tau []complex64, c []complex64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cupmtr((C.int)(rowMajor), (C.char)(s), (C.char)(ul), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&ap[0]), (*C.lapack_complex_float)(&tau[0]), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc)))
}

func Zupmtr(s blas.Side, ul blas.Uplo, trans blas.Transpose, m int, n int, ap []complex128, tau []complex128, c []complex128, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zupmtr((C.int)(rowMajor), (C.char)(s), (C.char)(ul), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&ap[0]), (*C.lapack_complex_double)(&tau[0]), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc)))
}

func SlapmrWork(forwrd int32, m int, n int, x []float32, ldx int, k []int32) bool {
	return isZero(C.LAPACKE_slapmr((C.int)(rowMajor), (C.lapack_logical)(forwrd), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&x[0]), (C.lapack_int)(ldx), (*C.lapack_int)(&k[0])))
}

func DlapmrWork(forwrd int32, m int, n int, x []float64, ldx int, k []int32) bool {
	return isZero(C.LAPACKE_dlapmr((C.int)(rowMajor), (C.lapack_logical)(forwrd), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&x[0]), (C.lapack_int)(ldx), (*C.lapack_int)(&k[0])))
}

func ClapmrWork(forwrd int32, m int, n int, x []complex64, ldx int, k []int32) bool {
	return isZero(C.LAPACKE_clapmr((C.int)(rowMajor), (C.lapack_logical)(forwrd), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(ldx), (*C.lapack_int)(&k[0])))
}

func ZlapmrWork(forwrd int32, m int, n int, x []complex128, ldx int, k []int32) bool {
	return isZero(C.LAPACKE_zlapmr((C.int)(rowMajor), (C.lapack_logical)(forwrd), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(ldx), (*C.lapack_int)(&k[0])))
}

func SlartgpWork(f float32, g float32, cs []float32, sn []float32, r []float32) bool {
	return isZero(C.LAPACKE_slartgp((C.float)(f), (C.float)(g), (*C.float)(&cs[0]), (*C.float)(&sn[0]), (*C.float)(&r[0])))
}

func DlartgpWork(f float64, g float64, cs []float64, sn []float64, r []float64) bool {
	return isZero(C.LAPACKE_dlartgp((C.double)(f), (C.double)(g), (*C.double)(&cs[0]), (*C.double)(&sn[0]), (*C.double)(&r[0])))
}

func SlartgsWork(x float32, y float32, sigma float32, cs []float32, sn []float32) bool {
	return isZero(C.LAPACKE_slartgs((C.float)(x), (C.float)(y), (C.float)(sigma), (*C.float)(&cs[0]), (*C.float)(&sn[0])))
}

func DlartgsWork(x float64, y float64, sigma float64, cs []float64, sn []float64) bool {
	return isZero(C.LAPACKE_dlartgs((C.double)(x), (C.double)(y), (C.double)(sigma), (*C.double)(&cs[0]), (*C.double)(&sn[0])))
}

func Slapy2Work(x float32, y float32) float32 {
	return float32(C.LAPACKE_slapy2((C.float)(x), (C.float)(y)))
}

func Dlapy2Work(x float64, y float64) float64 {
	return float64(C.LAPACKE_dlapy2((C.double)(x), (C.double)(y)))
}

func Slapy3Work(x float32, y float32, z float32) float32 {
	return float32(C.LAPACKE_slapy3((C.float)(x), (C.float)(y), (C.float)(z)))
}

func Dlapy3Work(x float64, y float64, z float64) float64 {
	return float64(C.LAPACKE_dlapy3((C.double)(x), (C.double)(y), (C.double)(z)))
}

func Cbbcsd(jobu1 lapack.Job, jobu2 lapack.Job, jobv1t lapack.Job, jobv2t lapack.Job, trans blas.Transpose, m int, p int, q int, theta []float32, phi []float32, u1 []complex64, ldu1 int, u2 []complex64, ldu2 int, v1t []complex64, ldv1t int, v2t []complex64, ldv2t int, b11d []float32, b11e []float32, b12d []float32, b12e []float32, b21d []float32, b21e []float32, b22d []float32, b22e []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cbbcsd((C.int)(rowMajor), (C.char)(jobu1), (C.char)(jobu2), (C.char)(jobv1t), (C.char)(jobv2t), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(q), (*C.float)(&theta[0]), (*C.float)(&phi[0]), (*C.lapack_complex_float)(&u1[0]), (C.lapack_int)(ldu1), (*C.lapack_complex_float)(&u2[0]), (C.lapack_int)(ldu2), (*C.lapack_complex_float)(&v1t[0]), (C.lapack_int)(ldv1t), (*C.lapack_complex_float)(&v2t[0]), (C.lapack_int)(ldv2t), (*C.float)(&b11d[0]), (*C.float)(&b11e[0]), (*C.float)(&b12d[0]), (*C.float)(&b12e[0]), (*C.float)(&b21d[0]), (*C.float)(&b21e[0]), (*C.float)(&b22d[0]), (*C.float)(&b22e[0])))
}

func Cheswapr(ul blas.Uplo, n int, a []complex64, i1 int, i2 int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_cheswapr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(i1), (C.lapack_int)(i2)))
}

func Chetri2(ul blas.Uplo, n int, a []complex64, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chetri2((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Chetri2x(ul blas.Uplo, n int, a []complex64, lda int, ipiv []int32, nb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chetri2x((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (C.lapack_int)(nb)))
}

func Chetrs2(ul blas.Uplo, n int, nrhs int, a []complex64, lda int, ipiv []int32, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_chetrs2((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Csyconv(ul blas.Uplo, way byte, n int, a []complex64, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csyconv((C.int)(rowMajor), (C.char)(ul), (C.char)(way), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Csyswapr(ul blas.Uplo, n int, a []complex64, i1 int, i2 int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csyswapr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(i1), (C.lapack_int)(i2)))
}

func Csytri2(ul blas.Uplo, n int, a []complex64, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csytri2((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Csytri2x(ul blas.Uplo, n int, a []complex64, lda int, ipiv []int32, nb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csytri2x((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (C.lapack_int)(nb)))
}

func Csytrs2(ul blas.Uplo, n int, nrhs int, a []complex64, lda int, ipiv []int32, b []complex64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csytrs2((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Cunbdb(trans blas.Transpose, signs byte, m int, p int, q int, x11 []complex64, ldx11 int, x12 []complex64, ldx12 int, x21 []complex64, ldx21 int, x22 []complex64, ldx22 int, theta []float32, phi []float32, taup1 []complex64, taup2 []complex64, tauq1 []complex64, tauq2 []complex64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cunbdb((C.int)(rowMajor), (C.char)(trans), (C.char)(signs), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(q), (*C.lapack_complex_float)(&x11[0]), (C.lapack_int)(ldx11), (*C.lapack_complex_float)(&x12[0]), (C.lapack_int)(ldx12), (*C.lapack_complex_float)(&x21[0]), (C.lapack_int)(ldx21), (*C.lapack_complex_float)(&x22[0]), (C.lapack_int)(ldx22), (*C.float)(&theta[0]), (*C.float)(&phi[0]), (*C.lapack_complex_float)(&taup1[0]), (*C.lapack_complex_float)(&taup2[0]), (*C.lapack_complex_float)(&tauq1[0]), (*C.lapack_complex_float)(&tauq2[0])))
}

func Cuncsd(jobu1 lapack.Job, jobu2 lapack.Job, jobv1t lapack.Job, jobv2t lapack.Job, trans blas.Transpose, signs byte, m int, p int, q int, x11 []complex64, ldx11 int, x12 []complex64, ldx12 int, x21 []complex64, ldx21 int, x22 []complex64, ldx22 int, theta []float32, u1 []complex64, ldu1 int, u2 []complex64, ldu2 int, v1t []complex64, ldv1t int, v2t []complex64, ldv2t int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cuncsd((C.int)(rowMajor), (C.char)(jobu1), (C.char)(jobu2), (C.char)(jobv1t), (C.char)(jobv2t), (C.char)(trans), (C.char)(signs), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(q), (*C.lapack_complex_float)(&x11[0]), (C.lapack_int)(ldx11), (*C.lapack_complex_float)(&x12[0]), (C.lapack_int)(ldx12), (*C.lapack_complex_float)(&x21[0]), (C.lapack_int)(ldx21), (*C.lapack_complex_float)(&x22[0]), (C.lapack_int)(ldx22), (*C.float)(&theta[0]), (*C.lapack_complex_float)(&u1[0]), (C.lapack_int)(ldu1), (*C.lapack_complex_float)(&u2[0]), (C.lapack_int)(ldu2), (*C.lapack_complex_float)(&v1t[0]), (C.lapack_int)(ldv1t), (*C.lapack_complex_float)(&v2t[0]), (C.lapack_int)(ldv2t)))
}

func Dbbcsd(jobu1 lapack.Job, jobu2 lapack.Job, jobv1t lapack.Job, jobv2t lapack.Job, trans blas.Transpose, m int, p int, q int, theta []float64, phi []float64, u1 []float64, ldu1 int, u2 []float64, ldu2 int, v1t []float64, ldv1t int, v2t []float64, ldv2t int, b11d []float64, b11e []float64, b12d []float64, b12e []float64, b21d []float64, b21e []float64, b22d []float64, b22e []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dbbcsd((C.int)(rowMajor), (C.char)(jobu1), (C.char)(jobu2), (C.char)(jobv1t), (C.char)(jobv2t), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(q), (*C.double)(&theta[0]), (*C.double)(&phi[0]), (*C.double)(&u1[0]), (C.lapack_int)(ldu1), (*C.double)(&u2[0]), (C.lapack_int)(ldu2), (*C.double)(&v1t[0]), (C.lapack_int)(ldv1t), (*C.double)(&v2t[0]), (C.lapack_int)(ldv2t), (*C.double)(&b11d[0]), (*C.double)(&b11e[0]), (*C.double)(&b12d[0]), (*C.double)(&b12e[0]), (*C.double)(&b21d[0]), (*C.double)(&b21e[0]), (*C.double)(&b22d[0]), (*C.double)(&b22e[0])))
}

func Dorbdb(trans blas.Transpose, signs byte, m int, p int, q int, x11 []float64, ldx11 int, x12 []float64, ldx12 int, x21 []float64, ldx21 int, x22 []float64, ldx22 int, theta []float64, phi []float64, taup1 []float64, taup2 []float64, tauq1 []float64, tauq2 []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dorbdb((C.int)(rowMajor), (C.char)(trans), (C.char)(signs), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(q), (*C.double)(&x11[0]), (C.lapack_int)(ldx11), (*C.double)(&x12[0]), (C.lapack_int)(ldx12), (*C.double)(&x21[0]), (C.lapack_int)(ldx21), (*C.double)(&x22[0]), (C.lapack_int)(ldx22), (*C.double)(&theta[0]), (*C.double)(&phi[0]), (*C.double)(&taup1[0]), (*C.double)(&taup2[0]), (*C.double)(&tauq1[0]), (*C.double)(&tauq2[0])))
}

func Dorcsd(jobu1 lapack.Job, jobu2 lapack.Job, jobv1t lapack.Job, jobv2t lapack.Job, trans blas.Transpose, signs byte, m int, p int, q int, x11 []float64, ldx11 int, x12 []float64, ldx12 int, x21 []float64, ldx21 int, x22 []float64, ldx22 int, theta []float64, u1 []float64, ldu1 int, u2 []float64, ldu2 int, v1t []float64, ldv1t int, v2t []float64, ldv2t int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dorcsd((C.int)(rowMajor), (C.char)(jobu1), (C.char)(jobu2), (C.char)(jobv1t), (C.char)(jobv2t), (C.char)(trans), (C.char)(signs), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(q), (*C.double)(&x11[0]), (C.lapack_int)(ldx11), (*C.double)(&x12[0]), (C.lapack_int)(ldx12), (*C.double)(&x21[0]), (C.lapack_int)(ldx21), (*C.double)(&x22[0]), (C.lapack_int)(ldx22), (*C.double)(&theta[0]), (*C.double)(&u1[0]), (C.lapack_int)(ldu1), (*C.double)(&u2[0]), (C.lapack_int)(ldu2), (*C.double)(&v1t[0]), (C.lapack_int)(ldv1t), (*C.double)(&v2t[0]), (C.lapack_int)(ldv2t)))
}

func Dsyconv(ul blas.Uplo, way byte, n int, a []float64, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsyconv((C.int)(rowMajor), (C.char)(ul), (C.char)(way), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Dsyswapr(ul blas.Uplo, n int, a []float64, i1 int, i2 int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsyswapr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(i1), (C.lapack_int)(i2)))
}

func Dsytri2(ul blas.Uplo, n int, a []float64, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsytri2((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Dsytri2x(ul blas.Uplo, n int, a []float64, lda int, ipiv []int32, nb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsytri2x((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (C.lapack_int)(nb)))
}

func Dsytrs2(ul blas.Uplo, n int, nrhs int, a []float64, lda int, ipiv []int32, b []float64, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_dsytrs2((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Sbbcsd(jobu1 lapack.Job, jobu2 lapack.Job, jobv1t lapack.Job, jobv2t lapack.Job, trans blas.Transpose, m int, p int, q int, theta []float32, phi []float32, u1 []float32, ldu1 int, u2 []float32, ldu2 int, v1t []float32, ldv1t int, v2t []float32, ldv2t int, b11d []float32, b11e []float32, b12d []float32, b12e []float32, b21d []float32, b21e []float32, b22d []float32, b22e []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sbbcsd((C.int)(rowMajor), (C.char)(jobu1), (C.char)(jobu2), (C.char)(jobv1t), (C.char)(jobv2t), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(q), (*C.float)(&theta[0]), (*C.float)(&phi[0]), (*C.float)(&u1[0]), (C.lapack_int)(ldu1), (*C.float)(&u2[0]), (C.lapack_int)(ldu2), (*C.float)(&v1t[0]), (C.lapack_int)(ldv1t), (*C.float)(&v2t[0]), (C.lapack_int)(ldv2t), (*C.float)(&b11d[0]), (*C.float)(&b11e[0]), (*C.float)(&b12d[0]), (*C.float)(&b12e[0]), (*C.float)(&b21d[0]), (*C.float)(&b21e[0]), (*C.float)(&b22d[0]), (*C.float)(&b22e[0])))
}

func Sorbdb(trans blas.Transpose, signs byte, m int, p int, q int, x11 []float32, ldx11 int, x12 []float32, ldx12 int, x21 []float32, ldx21 int, x22 []float32, ldx22 int, theta []float32, phi []float32, taup1 []float32, taup2 []float32, tauq1 []float32, tauq2 []float32) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sorbdb((C.int)(rowMajor), (C.char)(trans), (C.char)(signs), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(q), (*C.float)(&x11[0]), (C.lapack_int)(ldx11), (*C.float)(&x12[0]), (C.lapack_int)(ldx12), (*C.float)(&x21[0]), (C.lapack_int)(ldx21), (*C.float)(&x22[0]), (C.lapack_int)(ldx22), (*C.float)(&theta[0]), (*C.float)(&phi[0]), (*C.float)(&taup1[0]), (*C.float)(&taup2[0]), (*C.float)(&tauq1[0]), (*C.float)(&tauq2[0])))
}

func Sorcsd(jobu1 lapack.Job, jobu2 lapack.Job, jobv1t lapack.Job, jobv2t lapack.Job, trans blas.Transpose, signs byte, m int, p int, q int, x11 []float32, ldx11 int, x12 []float32, ldx12 int, x21 []float32, ldx21 int, x22 []float32, ldx22 int, theta []float32, u1 []float32, ldu1 int, u2 []float32, ldu2 int, v1t []float32, ldv1t int, v2t []float32, ldv2t int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sorcsd((C.int)(rowMajor), (C.char)(jobu1), (C.char)(jobu2), (C.char)(jobv1t), (C.char)(jobv2t), (C.char)(trans), (C.char)(signs), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(q), (*C.float)(&x11[0]), (C.lapack_int)(ldx11), (*C.float)(&x12[0]), (C.lapack_int)(ldx12), (*C.float)(&x21[0]), (C.lapack_int)(ldx21), (*C.float)(&x22[0]), (C.lapack_int)(ldx22), (*C.float)(&theta[0]), (*C.float)(&u1[0]), (C.lapack_int)(ldu1), (*C.float)(&u2[0]), (C.lapack_int)(ldu2), (*C.float)(&v1t[0]), (C.lapack_int)(ldv1t), (*C.float)(&v2t[0]), (C.lapack_int)(ldv2t)))
}

func Ssyconv(ul blas.Uplo, way byte, n int, a []float32, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssyconv((C.int)(rowMajor), (C.char)(ul), (C.char)(way), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Ssyswapr(ul blas.Uplo, n int, a []float32, i1 int, i2 int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssyswapr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(i1), (C.lapack_int)(i2)))
}

func Ssytri2(ul blas.Uplo, n int, a []float32, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssytri2((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Ssytri2x(ul blas.Uplo, n int, a []float32, lda int, ipiv []int32, nb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssytri2x((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (C.lapack_int)(nb)))
}

func Ssytrs2(ul blas.Uplo, n int, nrhs int, a []float32, lda int, ipiv []int32, b []float32, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_ssytrs2((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Zbbcsd(jobu1 lapack.Job, jobu2 lapack.Job, jobv1t lapack.Job, jobv2t lapack.Job, trans blas.Transpose, m int, p int, q int, theta []float64, phi []float64, u1 []complex128, ldu1 int, u2 []complex128, ldu2 int, v1t []complex128, ldv1t int, v2t []complex128, ldv2t int, b11d []float64, b11e []float64, b12d []float64, b12e []float64, b21d []float64, b21e []float64, b22d []float64, b22e []float64) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zbbcsd((C.int)(rowMajor), (C.char)(jobu1), (C.char)(jobu2), (C.char)(jobv1t), (C.char)(jobv2t), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(q), (*C.double)(&theta[0]), (*C.double)(&phi[0]), (*C.lapack_complex_double)(&u1[0]), (C.lapack_int)(ldu1), (*C.lapack_complex_double)(&u2[0]), (C.lapack_int)(ldu2), (*C.lapack_complex_double)(&v1t[0]), (C.lapack_int)(ldv1t), (*C.lapack_complex_double)(&v2t[0]), (C.lapack_int)(ldv2t), (*C.double)(&b11d[0]), (*C.double)(&b11e[0]), (*C.double)(&b12d[0]), (*C.double)(&b12e[0]), (*C.double)(&b21d[0]), (*C.double)(&b21e[0]), (*C.double)(&b22d[0]), (*C.double)(&b22e[0])))
}

func Zheswapr(ul blas.Uplo, n int, a []complex128, i1 int, i2 int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zheswapr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(i1), (C.lapack_int)(i2)))
}

func Zhetri2(ul blas.Uplo, n int, a []complex128, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhetri2((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Zhetri2x(ul blas.Uplo, n int, a []complex128, lda int, ipiv []int32, nb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhetri2x((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (C.lapack_int)(nb)))
}

func Zhetrs2(ul blas.Uplo, n int, nrhs int, a []complex128, lda int, ipiv []int32, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zhetrs2((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Zsyconv(ul blas.Uplo, way byte, n int, a []complex128, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsyconv((C.int)(rowMajor), (C.char)(ul), (C.char)(way), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Zsyswapr(ul blas.Uplo, n int, a []complex128, i1 int, i2 int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsyswapr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(i1), (C.lapack_int)(i2)))
}

func Zsytri2(ul blas.Uplo, n int, a []complex128, lda int, ipiv []int32) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsytri2((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0])))
}

func Zsytri2x(ul blas.Uplo, n int, a []complex128, lda int, ipiv []int32, nb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsytri2x((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (C.lapack_int)(nb)))
}

func Zsytrs2(ul blas.Uplo, n int, nrhs int, a []complex128, lda int, ipiv []int32, b []complex128, ldb int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsytrs2((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_int)(nrhs), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_int)(&ipiv[0]), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Zunbdb(trans blas.Transpose, signs byte, m int, p int, q int, x11 []complex128, ldx11 int, x12 []complex128, ldx12 int, x21 []complex128, ldx21 int, x22 []complex128, ldx22 int, theta []float64, phi []float64, taup1 []complex128, taup2 []complex128, tauq1 []complex128, tauq2 []complex128) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zunbdb((C.int)(rowMajor), (C.char)(trans), (C.char)(signs), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(q), (*C.lapack_complex_double)(&x11[0]), (C.lapack_int)(ldx11), (*C.lapack_complex_double)(&x12[0]), (C.lapack_int)(ldx12), (*C.lapack_complex_double)(&x21[0]), (C.lapack_int)(ldx21), (*C.lapack_complex_double)(&x22[0]), (C.lapack_int)(ldx22), (*C.double)(&theta[0]), (*C.double)(&phi[0]), (*C.lapack_complex_double)(&taup1[0]), (*C.lapack_complex_double)(&taup2[0]), (*C.lapack_complex_double)(&tauq1[0]), (*C.lapack_complex_double)(&tauq2[0])))
}

func Zuncsd(jobu1 lapack.Job, jobu2 lapack.Job, jobv1t lapack.Job, jobv2t lapack.Job, trans blas.Transpose, signs byte, m int, p int, q int, x11 []complex128, ldx11 int, x12 []complex128, ldx12 int, x21 []complex128, ldx21 int, x22 []complex128, ldx22 int, theta []float64, u1 []complex128, ldu1 int, u2 []complex128, ldu2 int, v1t []complex128, ldv1t int, v2t []complex128, ldv2t int) bool {
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zuncsd((C.int)(rowMajor), (C.char)(jobu1), (C.char)(jobu2), (C.char)(jobv1t), (C.char)(jobv2t), (C.char)(trans), (C.char)(signs), (C.lapack_int)(m), (C.lapack_int)(p), (C.lapack_int)(q), (*C.lapack_complex_double)(&x11[0]), (C.lapack_int)(ldx11), (*C.lapack_complex_double)(&x12[0]), (C.lapack_int)(ldx12), (*C.lapack_complex_double)(&x21[0]), (C.lapack_int)(ldx21), (*C.lapack_complex_double)(&x22[0]), (C.lapack_int)(ldx22), (*C.double)(&theta[0]), (*C.lapack_complex_double)(&u1[0]), (C.lapack_int)(ldu1), (*C.lapack_complex_double)(&u2[0]), (C.lapack_int)(ldu2), (*C.lapack_complex_double)(&v1t[0]), (C.lapack_int)(ldv1t), (*C.lapack_complex_double)(&v2t[0]), (C.lapack_int)(ldv2t)))
}

func Sgemqrt(s blas.Side, trans blas.Transpose, m int, n int, k int, nb int, v []float32, ldv int, t []float32, ldt int, c []float32, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_sgemqrt((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(nb), (*C.float)(&v[0]), (C.lapack_int)(ldv), (*C.float)(&t[0]), (C.lapack_int)(ldt), (*C.float)(&c[0]), (C.lapack_int)(ldc)))
}

func Dgemqrt(s blas.Side, trans blas.Transpose, m int, n int, k int, nb int, v []float64, ldv int, t []float64, ldt int, c []float64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dgemqrt((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(nb), (*C.double)(&v[0]), (C.lapack_int)(ldv), (*C.double)(&t[0]), (C.lapack_int)(ldt), (*C.double)(&c[0]), (C.lapack_int)(ldc)))
}

func Cgemqrt(s blas.Side, trans blas.Transpose, m int, n int, k int, nb int, v []complex64, ldv int, t []complex64, ldt int, c []complex64, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_cgemqrt((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(nb), (*C.lapack_complex_float)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_float)(&t[0]), (C.lapack_int)(ldt), (*C.lapack_complex_float)(&c[0]), (C.lapack_int)(ldc)))
}

func Zgemqrt(s blas.Side, trans blas.Transpose, m int, n int, k int, nb int, v []complex128, ldv int, t []complex128, ldt int, c []complex128, ldc int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_zgemqrt((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(nb), (*C.lapack_complex_double)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_double)(&t[0]), (C.lapack_int)(ldt), (*C.lapack_complex_double)(&c[0]), (C.lapack_int)(ldc)))
}

func Sgeqrt(m int, n int, nb int, a []float32, lda int, t []float32, ldt int) bool {
	return isZero(C.LAPACKE_sgeqrt((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nb), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&t[0]), (C.lapack_int)(ldt)))
}

func Dgeqrt(m int, n int, nb int, a []float64, lda int, t []float64, ldt int) bool {
	return isZero(C.LAPACKE_dgeqrt((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nb), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&t[0]), (C.lapack_int)(ldt)))
}

func Cgeqrt(m int, n int, nb int, a []complex64, lda int, t []complex64, ldt int) bool {
	return isZero(C.LAPACKE_cgeqrt((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nb), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&t[0]), (C.lapack_int)(ldt)))
}

func Zgeqrt(m int, n int, nb int, a []complex128, lda int, t []complex128, ldt int) bool {
	return isZero(C.LAPACKE_zgeqrt((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(nb), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&t[0]), (C.lapack_int)(ldt)))
}

func Sgeqrt2(m int, n int, a []float32, lda int, t []float32, ldt int) bool {
	return isZero(C.LAPACKE_sgeqrt2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&t[0]), (C.lapack_int)(ldt)))
}

func Dgeqrt2(m int, n int, a []float64, lda int, t []float64, ldt int) bool {
	return isZero(C.LAPACKE_dgeqrt2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&t[0]), (C.lapack_int)(ldt)))
}

func Cgeqrt2(m int, n int, a []complex64, lda int, t []complex64, ldt int) bool {
	return isZero(C.LAPACKE_cgeqrt2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&t[0]), (C.lapack_int)(ldt)))
}

func Zgeqrt2(m int, n int, a []complex128, lda int, t []complex128, ldt int) bool {
	return isZero(C.LAPACKE_zgeqrt2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&t[0]), (C.lapack_int)(ldt)))
}

func Sgeqrt3(m int, n int, a []float32, lda int, t []float32, ldt int) bool {
	return isZero(C.LAPACKE_sgeqrt3((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&t[0]), (C.lapack_int)(ldt)))
}

func Dgeqrt3(m int, n int, a []float64, lda int, t []float64, ldt int) bool {
	return isZero(C.LAPACKE_dgeqrt3((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&t[0]), (C.lapack_int)(ldt)))
}

func Cgeqrt3(m int, n int, a []complex64, lda int, t []complex64, ldt int) bool {
	return isZero(C.LAPACKE_cgeqrt3((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&t[0]), (C.lapack_int)(ldt)))
}

func Zgeqrt3(m int, n int, a []complex128, lda int, t []complex128, ldt int) bool {
	return isZero(C.LAPACKE_zgeqrt3((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&t[0]), (C.lapack_int)(ldt)))
}

func Stpmqrt(s blas.Side, trans blas.Transpose, m int, n int, k int, l int, nb int, v []float32, ldv int, t []float32, ldt int, a []float32, lda int, b []float32, ldb int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_stpmqrt((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (C.lapack_int)(nb), (*C.float)(&v[0]), (C.lapack_int)(ldv), (*C.float)(&t[0]), (C.lapack_int)(ldt), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dtpmqrt(s blas.Side, trans blas.Transpose, m int, n int, k int, l int, nb int, v []float64, ldv int, t []float64, ldt int, a []float64, lda int, b []float64, ldb int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dtpmqrt((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (C.lapack_int)(nb), (*C.double)(&v[0]), (C.lapack_int)(ldv), (*C.double)(&t[0]), (C.lapack_int)(ldt), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Ctpmqrt(s blas.Side, trans blas.Transpose, m int, n int, k int, l int, nb int, v []complex64, ldv int, t []complex64, ldt int, a []complex64, lda int, b []complex64, ldb int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_ctpmqrt((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (C.lapack_int)(nb), (*C.lapack_complex_float)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_float)(&t[0]), (C.lapack_int)(ldt), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Ztpmqrt(s blas.Side, trans blas.Transpose, m int, n int, k int, l int, nb int, v []complex128, ldv int, t []complex128, ldt int, a []complex128, lda int, b []complex128, ldb int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_ztpmqrt((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (C.lapack_int)(nb), (*C.lapack_complex_double)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_double)(&t[0]), (C.lapack_int)(ldt), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Dtpqrt(m int, n int, l int, nb int, a []float64, lda int, b []float64, ldb int, t []float64, ldt int) bool {
	return isZero(C.LAPACKE_dtpqrt((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(l), (C.lapack_int)(nb), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&t[0]), (C.lapack_int)(ldt)))
}

func Ctpqrt(m int, n int, l int, nb int, a []complex64, lda int, b []complex64, ldb int, t []complex64, ldt int) bool {
	return isZero(C.LAPACKE_ctpqrt((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(l), (C.lapack_int)(nb), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&t[0]), (C.lapack_int)(ldt)))
}

func Ztpqrt(m int, n int, l int, nb int, a []complex128, lda int, b []complex128, ldb int, t []complex128, ldt int) bool {
	return isZero(C.LAPACKE_ztpqrt((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(l), (C.lapack_int)(nb), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&t[0]), (C.lapack_int)(ldt)))
}

func Stpqrt2(m int, n int, l int, a []float32, lda int, b []float32, ldb int, t []float32, ldt int) bool {
	return isZero(C.LAPACKE_stpqrt2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(l), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb), (*C.float)(&t[0]), (C.lapack_int)(ldt)))
}

func Dtpqrt2(m int, n int, l int, a []float64, lda int, b []float64, ldb int, t []float64, ldt int) bool {
	return isZero(C.LAPACKE_dtpqrt2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(l), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb), (*C.double)(&t[0]), (C.lapack_int)(ldt)))
}

func Ctpqrt2(m int, n int, l int, a []complex64, lda int, b []complex64, ldb int, t []complex64, ldt int) bool {
	return isZero(C.LAPACKE_ctpqrt2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(l), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_float)(&t[0]), (C.lapack_int)(ldt)))
}

func Ztpqrt2(m int, n int, l int, a []complex128, lda int, b []complex128, ldb int, t []complex128, ldt int) bool {
	return isZero(C.LAPACKE_ztpqrt2((C.int)(rowMajor), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(l), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb), (*C.lapack_complex_double)(&t[0]), (C.lapack_int)(ldt)))
}

func Stprfb(s blas.Side, trans blas.Transpose, direct byte, storev byte, m int, n int, k int, l int, v []float32, ldv int, t []float32, ldt int, a []float32, lda int, b []float32, ldb int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_stprfb((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.char)(direct), (C.char)(storev), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (*C.float)(&v[0]), (C.lapack_int)(ldv), (*C.float)(&t[0]), (C.lapack_int)(ldt), (*C.float)(&a[0]), (C.lapack_int)(lda), (*C.float)(&b[0]), (C.lapack_int)(ldb)))
}

func Dtprfb(s blas.Side, trans blas.Transpose, direct byte, storev byte, m int, n int, k int, l int, v []float64, ldv int, t []float64, ldt int, a []float64, lda int, b []float64, ldb int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_dtprfb((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.char)(direct), (C.char)(storev), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (*C.double)(&v[0]), (C.lapack_int)(ldv), (*C.double)(&t[0]), (C.lapack_int)(ldt), (*C.double)(&a[0]), (C.lapack_int)(lda), (*C.double)(&b[0]), (C.lapack_int)(ldb)))
}

func Ctprfb(s blas.Side, trans blas.Transpose, direct byte, storev byte, m int, n int, k int, l int, v []complex64, ldv int, t []complex64, ldt int, a []complex64, lda int, b []complex64, ldb int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_ctprfb((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.char)(direct), (C.char)(storev), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (*C.lapack_complex_float)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_float)(&t[0]), (C.lapack_int)(ldt), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_float)(&b[0]), (C.lapack_int)(ldb)))
}

func Ztprfb(s blas.Side, trans blas.Transpose, direct byte, storev byte, m int, n int, k int, l int, v []complex128, ldv int, t []complex128, ldt int, a []complex128, lda int, b []complex128, ldb int) bool {
	switch s {
	case blas.Left:
		s = 'L'
	case blas.Right:
		s = 'R'
	default:
		panic("lapack: bad side")
	}
	switch trans {
	case blas.NoTrans:
		trans = 'N'
	case blas.Trans:
		trans = 'T'
	case blas.ConjTrans:
		trans = 'C'
	default:
		panic("lapack: bad trans")
	}
	return isZero(C.LAPACKE_ztprfb((C.int)(rowMajor), (C.char)(s), (C.char)(trans), (C.char)(direct), (C.char)(storev), (C.lapack_int)(m), (C.lapack_int)(n), (C.lapack_int)(k), (C.lapack_int)(l), (*C.lapack_complex_double)(&v[0]), (C.lapack_int)(ldv), (*C.lapack_complex_double)(&t[0]), (C.lapack_int)(ldt), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda), (*C.lapack_complex_double)(&b[0]), (C.lapack_int)(ldb)))
}

func Csyr(ul blas.Uplo, n int, alpha complex64, x []complex64, incx int, a []complex64, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_csyr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_complex_float)(alpha), (*C.lapack_complex_float)(&x[0]), (C.lapack_int)(incx), (*C.lapack_complex_float)(&a[0]), (C.lapack_int)(lda)))
}

func Zsyr(ul blas.Uplo, n int, alpha complex128, x []complex128, incx int, a []complex128, lda int) bool {
	switch ul {
	case blas.Upper:
		ul = 'U'
	case blas.Lower:
		ul = 'L'
	default:
		panic("lapack: illegal triangle")
	}
	return isZero(C.LAPACKE_zsyr((C.int)(rowMajor), (C.char)(ul), (C.lapack_int)(n), (C.lapack_complex_double)(alpha), (*C.lapack_complex_double)(&x[0]), (C.lapack_int)(incx), (*C.lapack_complex_double)(&a[0]), (C.lapack_int)(lda)))
}
