package latypes

import (
	"github.com/dane-unltd/goblas"
	"github.com/gonum/blas"
)

type Lapack interface {
	Dgeqrf(order blas.Order, A goblas.General, tau []float64)
}
