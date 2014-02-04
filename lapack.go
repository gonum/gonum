package lapack

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/dbw"
)

type Float64 interface {
	Dgeqrf(A dbw.General, tau []float64)
	Dormqr(s blas.Side, t blas.Transpose, A dbw.General, tau []float64, B dbw.General)
}
