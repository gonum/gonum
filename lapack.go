package lapack

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/dbw"
	"github.com/gonum/blas/zbw"
)

type Float64 interface {
	Dgeqrf(A dbw.General, tau []float64)
	Dormqr(s blas.Side, t blas.Transpose, A dbw.General, tau []float64, B dbw.General)
}

type Complex128 interface {
	Zgesvd(jobu byte, jobvt byte, A zbw.General, s []float64, U zbw.General, Vt zbw.General, superb []float64)
}
