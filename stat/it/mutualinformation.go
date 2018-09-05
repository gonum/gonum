package it

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// MutualInformation calculates the mutual information with the given lnFunc function
//   I(X,Y) = \sum_x,y p(x,y) (lnFunc(p(x,y)) - lnFunc(p(x)p(y)))
func MutualInformation(pxy mat.Matrix) float64 {

	xDim, yDim := pxy.Dims()

	px := make([]float64, xDim)
	py := make([]float64, yDim)

	for x := 0; x < xDim; x++ {
		for y := 0; y < yDim; y++ {
			px[x] += pxy.At(x, y)
		}
	}

	for x := 0; x < xDim; x++ {
		for y := 0; y < yDim; y++ {
			py[y] += pxy.At(x, y)
		}
	}

	mi := 0.0

	for x := 0; x < xDim; x++ {
		if px[x] > 0.0 {
			for y := 0; y < yDim; y++ {
				v := pxy.At(x, y)
				if py[y] > 0.0 && pxy.At(x, y) > 0.0 {
					mi += v * (math.Log(v) - math.Log(px[x]*py[y]))
				}
			}
		}
	}
	return mi
}
