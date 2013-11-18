package goblas

import "github.com/ziutek/blas"

func (Blas) Daxpy(N int, alpha float64, X []float64, incX int, Y []float64, incY int) {
	blas.Daxpy(N, alpha, X, incX, Y, incY)
}
