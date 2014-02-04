package clapack

/*
#cgo linux LDFLAGS: -llapacke -lblas
#cgo darwin LDFLAGS: -DYA_BLAS -DYA_LAPACK -DYA_BLASMULT -framework vecLib
#include <stdlib.h>
#include "lapacke.h"
*/
import "C"
import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/dbw"
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
