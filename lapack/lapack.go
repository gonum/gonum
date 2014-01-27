package latypes

import "gihub.com/dane-unltd/goblas"

type Lapack interface {
	Dgeqrf(order blas.Order, A goblas.General, tau []float64)
}
