package goblas

import "github.com/ziutek/blas"

func (Blas) Idamax(N int, X []float64, incX int) int {
	return blas.Idamax(N, X, incX)
}
