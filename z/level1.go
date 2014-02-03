package z

func Zdotu(x, y Vector) complex128 {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	return impl.Zdotu(x.N, x.Data, x.Inc, y.Data, y.Inc)
}

func Zdotc(x, y Vector) complex128 {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	return impl.Zdotc(x.N, x.Data, x.Inc, y.Data, y.Inc)
}

func Znrm2(x Vector) float64 {
	return impl.Dznrm2(x.N, x.Data, x.Inc)
}

func Dzasum(x Vector) float64 {
	return impl.Dzasum(x.N, x.Data, x.Inc)
}

func Izmax(x Vector) int {
	return impl.Izamax(x.N, x.Data, x.Inc)
}

func Zswap(x, y Vector) {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Zswap(x.N, x.Data, x.Inc, y.Data, y.Inc)
}

func Zcopy(x, y Vector) {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Zcopy(x.N, x.Data, x.Inc, y.Data, y.Inc)
}

func Zaxpy(alpha complex128, x, y Vector) {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Zaxpy(x.N, alpha, x.Data, x.Inc, y.Data, y.Inc)
}

func Zscal(alpha complex128, x Vector) {
	impl.Zscal(x.N, alpha, x.Data, x.Inc)
}

func Zdscal(alpha float64, x Vector) {
	impl.Zdscal(x.N, alpha, x.Data, x.Inc)
}
