package zbw

func Dotu(x, y Vector) complex128 {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	return impl.Zdotu(x.N, x.Data, x.Inc, y.Data, y.Inc)
}

func Dotc(x, y Vector) complex128 {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	return impl.Zdotc(x.N, x.Data, x.Inc, y.Data, y.Inc)
}

func Nrm2(x Vector) float64 {
	return impl.Dznrm2(x.N, x.Data, x.Inc)
}

func Asum(x Vector) float64 {
	return impl.Dzasum(x.N, x.Data, x.Inc)
}

func Iamax(x Vector) int {
	return impl.Izamax(x.N, x.Data, x.Inc)
}

func Swap(x, y Vector) {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Zswap(x.N, x.Data, x.Inc, y.Data, y.Inc)
}

func Copy(x, y Vector) {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Zcopy(x.N, x.Data, x.Inc, y.Data, y.Inc)
}

func Axpy(alpha complex128, x, y Vector) {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Zaxpy(x.N, alpha, x.Data, x.Inc, y.Data, y.Inc)
}

func Scal(alpha complex128, x Vector) {
	impl.Zscal(x.N, alpha, x.Data, x.Inc)
}

func Dscal(alpha float64, x Vector) {
	impl.Zdscal(x.N, alpha, x.Data, x.Inc)
}
