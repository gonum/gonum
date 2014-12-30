package testblas

import (
	"math/rand"
	"testing"

	"github.com/gonum/blas"
)

func DgemvBenchmark(b *testing.B, blasser Dgemver, tA blas.Transpose, m, n, incX, incY int) {
	var lenX, lenY int
	if tA == blas.NoTrans {
		lenX = n
		lenY = m
	} else {
		lenX = m
		lenY = n
	}
	xr := make([]float64, lenX)
	for i := range xr {
		xr[i] = rand.Float64()
	}
	x := makeIncremented(xr, incX, 0)
	yr := make([]float64, lenY)
	for i := range yr {
		yr[i] = rand.Float64()
	}
	y := makeIncremented(yr, incY, 0)
	a := make([]float64, m*n)
	for i := range a {
		a[i] = rand.Float64()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		blasser.Dgemv(tA, m, n, 2, a, n, x, incX, 3, y, incY)
	}
}
