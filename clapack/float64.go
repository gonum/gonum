package clapack

/*
#cgo linux LDFLAGS: -llapacke -lblas
#cgo darwin LDFLAGS: -DYA_BLAS -DYA_LAPACK -DYA_BLASMULT -framework vecLib
#include <stdlib.h>
#include "lapacke.h"
*/
import "C"
import (
	"github.com/gonum/blas/d"
)

type La struct{}

func (La) Dgeqrf(A d.General, tau []float64) {
	C.LAPACKE_dgeqrf(C.int(A.Order), C.int(A.Rows), C.int(A.Cols),
		(*C.double)(&A.Data[0]), C.int(A.Stride), (*C.double)(&tau[0]))
}

func (La) Dormqr(s byte, t byte, A d.General, tau []float64, B d.General) {
	C.LAPACKE_dormqr(C.int(A.Order), C.char(s), C.char(t), C.int(B.Rows),
		C.int(B.Cols), C.int(A.Cols), (*C.double)(&A.Data[0]),
		C.int(A.Stride), (*C.double)(&tau[0]), (*C.double)(&B.Data[0]), C.int(B.Stride))
}
