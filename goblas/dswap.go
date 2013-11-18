package goblas

import "github.com/ziutek/blas"

func (Blas) Dswap(N int, X []float64, incX int, Y []float64, incY int) {
	blas.Dswap(N, X, incX, Y, incY)
}
