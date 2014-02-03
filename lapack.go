package lapack

import (
	"github.com/gonum/blas/d"
)

type Float64 interface {
	Dgeqrf(A d.General, tau []float64)
	Dormqr(s byte, t byte, A d.General, tau []float64, B d.General)
}
