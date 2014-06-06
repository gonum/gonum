package zla

import "github.com/gonum/lapack"

var impl lapack.Complex128

func Register(i lapack.Complex128) {
	impl = i
}
