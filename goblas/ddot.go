package goblas

import "github.com/ziutek/blas"

func (Blas) Ddot(N int, X []float64, incX int, Y []float64, incY int) float64 {
	return blas.Ddot(N, X, incX, Y, incY)
}
