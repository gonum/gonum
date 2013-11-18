package goblas

import "github.com/ziutek/blas"

func (Blas) Drot(N int, X []float64, incX int, Y []float64, incY int, c, s float64) {
	blas.Drot(N, X, incX, Y, incY, c, s)
}
