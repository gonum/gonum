// +build !amd64

package goblas

func ddotUnitary(x []float64, y []float64) (sum float64) {
	for i, v := range x {
		sum += y[i] * v
	}
	return
}

func ddotInc(x, y []float64, n, incX, incY, ix, iy uintptr) (sum float64) {
	for i := 0; i < int(n); i++ {
		sum += y[iy] * x[ix]
		ix += incX
		iy += incY
	}
	return
}
