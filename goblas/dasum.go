package goblas

import "github.com/ziutek/blas"

func (Blas) Dasum(N int, X []float64, incX int) float64 {
	return blas.Dasum(N, X, incX)
}
