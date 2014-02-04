package dbw

import "github.com/gonum/blas"

var impl blas.Float64

func Register(i blas.Float64) {
	impl = i
}
