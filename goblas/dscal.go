package goblas

import "github.com/ziutek/blas"

func (Blas) Dscal(N int, alpha float64, X []float64, incX int) {
	blas.Dscal(N, alpha, X, incX)
}
