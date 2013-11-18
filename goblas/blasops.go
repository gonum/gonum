package goblas

import "github.com/gonum/blas"

type Blas struct{}

var (
	_ blas.Float64 = Blas{}
)
