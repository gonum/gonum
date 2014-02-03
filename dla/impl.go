package dla

import "github.com/dane-unltd/lapack"

var impl lapack.Float64

func Register(i lapack.Float64) {
	impl = i
}
