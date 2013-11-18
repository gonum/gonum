package goblas

import "github.com/ziutek/blas"

func (Blas) Dnrm2(N int, X []float64, incX int) float64 {
	return blas.Dnrm2(N, X, incX)
}
