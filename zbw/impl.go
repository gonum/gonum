package zbw

import "github.com/gonum/blas"

var impl blas.Complex128

func Register(i blas.Complex128) {
	impl = i
}
