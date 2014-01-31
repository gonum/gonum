package blas

var impl Float64
var implCmplx Complex128

func Register(i Float64) {
	impl = i
}

func RegisterCmplx(i Complex128) {
	implCmplx = i
}
