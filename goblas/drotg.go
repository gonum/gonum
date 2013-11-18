package goblas

import "github.com/ziutek/blas"

func (Blas) Drotg(a, b float64) (c, s, r, z float64) {
	return blas.Drotg(a, b)
}
