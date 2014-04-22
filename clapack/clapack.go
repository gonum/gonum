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
	"github.com/gonum/blas"
	"github.com/gonum/blas/dbw"
	"github.com/gonum/blas/zbw"
)

type La struct{}

func (La) Dgeqrf(A dbw.General, tau []float64) {
	C.LAPACKE_dgeqrf(C.int(A.Order), C.int(A.Rows), C.int(A.Cols),
		(*C.double)(&A.Data[0]), C.int(A.Stride), (*C.double)(&tau[0]))
}

func (La) Dormqr(s blas.Side, t blas.Transpose, A dbw.General, tau []float64, B dbw.General) {
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

	C.LAPACKE_dormqr(C.int(A.Order), cs, ct, C.int(B.Rows),
		C.int(B.Cols), C.int(A.Cols), (*C.double)(&A.Data[0]),
		C.int(A.Stride), (*C.double)(&tau[0]), (*C.double)(&B.Data[0]), C.int(B.Stride))
}

func (La) Zgesvd(jobz byte, A zbw.General, s []float64, U zbw.General, Vt zbw.General) {
	pU := (*complex128)(nil)
	if len(U.Data) > 0 {
		pU = &U.Data[0]
	}
	pVt := (*complex128)(nil)
	if len(Vt.Data) > 0 {
		pVt = &Vt.Data[0]
	}
	C.LAPACKE_zgesdd(
		C.int(A.Order), C.char(jobz),
		C.int(A.Rows), C.int(A.Cols), (*C.complex)(&A.Data[0]), C.int(A.Stride),
		(*C.double)(&s[0]),
		(*C.complex)(pU), C.int(U.Stride),
		(*C.complex)(pVt), C.int(Vt.Stride))
}
