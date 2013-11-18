package goblas

import "github.com/ziutek/blas"

func (Blas) Dcopy(N int, X []float64, incX int, Y []float64, incY int) {
	blas.Dcopy(N, X, incX, Y, incY)
}
