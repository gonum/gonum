package clapack

/*
//#cgo linux LDFLAGS: -Wl,--no-undefined -lmkl_core -lmkl_rt
#cgo linux LDFLAGS: -L/home/users/dane/builds/OpenBLAS -lopenblas
#cgo darwin LDFLAGS: -DYA_BLAS -DYA_LAPACK -DYA_BLASMULT -framework vecLib
#include <stdlib.h>
#include "lapacke.h"
*/
import "C"
import (
	"github.com/dane-unltd/lapack"
	"github.com/gonum/blas"
)

type La struct{}

func init() {
	_ = lapack.Complex128(La{})
	_ = lapack.Float64(La{})
}

func (La) Dgeqrf(o blas.Order, m, n int, a []float64, lda int, tau []float64) {
	C.LAPACKE_dgeqrf(C.int(o), C.int(m), C.int(n), (*C.double)(&a[0]), C.int(lda), (*C.double)(&tau[0]))
}

func (La) Dormqr(o blas.Order, s blas.Side, t blas.Transpose, m, n, k int, a []float64, lda int, tau []float64, c []float64, ldc int) {
	var cs, ct C.char
	if s == blas.Left {
		cs = 'l'
	} else {
		cs = 'r'
	}
	if t == blas.NoTrans {
		ct = 'n'
	} else {
		ct = 't'
	}

	C.LAPACKE_dormqr(C.int(o), cs, ct, C.int(m),
		C.int(n), C.int(k), (*C.double)(&a[0]),
		C.int(lda), (*C.double)(&tau[0]), (*C.double)(&c[0]), C.int(ldc))
}

func (La) Dgesdd(o blas.Order, job lapack.Job, m, n int, a []float64, lda int, s []float64, u []float64, ldu int, vt []float64, ldvt int) {
	pU := (*float64)(nil)
	if len(u) > 0 {
		pU = &u[0]
	}
	pVt := (*float64)(nil)
	if len(vt) > 0 {
		pVt = &vt[0]
	}
	C.LAPACKE_dgesdd(
		C.int(o), C.char(job),
		C.int(m), C.int(n), (*C.double)(&a[0]), C.int(lda),
		(*C.double)(&s[0]),
		(*C.double)(pU), C.int(ldu),
		(*C.double)(pVt), C.int(ldvt))
}

func (La) Dgebrd(o blas.Order, m, n int, a []float64, lda int, d, e, tauq, taup []float64) {
	C.LAPACKE_dgebrd(
		C.int(o), C.int(m), C.int(n), (*C.double)(&a[0]), C.int(lda),
		(*C.double)(&d[0]),
		(*C.double)(&e[0]),
		(*C.double)(&tauq[0]),
		(*C.double)(&taup[0]))
}

func (La) Dbdsdc(o blas.Order, uplo blas.Uplo, compq lapack.CompSV, n int,
	d, e []float64, u []float64, ldu int, vt []float64, ldvt int, q []float64, iq []int32) {
	pU := (*float64)(nil)
	if len(u) > 0 {
		pU = &u[0]
	}
	pVt := (*float64)(nil)
	if len(vt) > 0 {
		pVt = &vt[0]
	}
	pq := (*float64)(nil)
	if len(q) > 0 {
		pU = &q[0]
	}
	piq := (*int32)(nil)
	if len(iq) > 0 {
		piq = &iq[0]
	}

	cuplo := C.char('u')
	if uplo == blas.Lower {
		cuplo = 'l'
	}

	C.LAPACKE_dbdsdc(C.int(o), cuplo, C.char(compq),
		(C.int)(n),
		(*C.double)(&d[0]),
		(*C.double)(&e[0]),
		(*C.double)(pU),
		(C.int)(ldu),
		(*C.double)(pVt),
		(C.int)(ldvt),
		(*C.double)(pq),
		(*C.int)(piq))
}

func (La) Zgeqrf(o blas.Order, m, n int, a []complex128, lda int, tau []complex128) {
	C.LAPACKE_zgeqrf(C.int(o), C.int(m), C.int(n), (*C.complex)(&a[0]), C.int(lda), (*C.complex)(&tau[0]))
}

func (La) Zunmqr(o blas.Order, s blas.Side, t blas.Transpose, m, n, k int, a []complex128, lda int, tau []complex128, c []complex128, ldc int) {
	var cs, ct C.char
	if s == blas.Left {
		cs = 'l'
	} else {
		cs = 'r'
	}
	if t == blas.NoTrans {
		ct = 'n'
	} else {
		ct = 'c'
	}

	C.LAPACKE_zunmqr(C.int(o), cs, ct, C.int(m),
		C.int(n), C.int(k), (*C.complex)(&a[0]),
		C.int(lda), (*C.complex)(&tau[0]), (*C.complex)(&c[0]), C.int(ldc))
}

func (La) Zgesdd(o blas.Order, job lapack.Job, m, n int, a []complex128, lda int, s []float64, u []complex128, ldu int, vt []complex128, ldvt int) {
	pU := (*complex128)(nil)
	if len(u) > 0 {
		pU = &u[0]
	}
	pVt := (*complex128)(nil)
	if len(vt) > 0 {
		pVt = &vt[0]
	}
	C.LAPACKE_zgesdd(
		C.int(o), C.char(job),
		C.int(m), C.int(n), (*C.complex)(&a[0]), C.int(lda),
		(*C.double)(&s[0]),
		(*C.complex)(pU), C.int(ldu),
		(*C.complex)(pVt), C.int(ldvt))
}

func (La) Zgebrd(o blas.Order, m, n int, a []complex128, lda int, d, e []float64, tauq, taup []complex128) {
	C.LAPACKE_zgebrd(
		C.int(o),
		C.int(m), C.int(n), (*C.complex)(&a[0]), C.int(lda),
		(*C.double)(&d[0]),
		(*C.double)(&e[0]),
		(*C.complex)(&tauq[0]),
		(*C.complex)(&taup[0]))
}
