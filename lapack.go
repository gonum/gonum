package lapack

import (
	"github.com/gonum/blas"
)

type Lapack interface {
	Dgeqrf(A blas.General, tau []float64)
	Dormqr(s blas.Side, t blas.Transpose, A blas.General, tau []float64, B blas.General)
}
