package blas

func Ddot(x, y Vector) float64 {
	must(x.Check())
	must(y.Check())
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	return impl.Ddot(x.N, x.Data, x.Inc, y.Data, y.Inc)
}

func Dnrm2(x Vector) float64 {
	must(x.Check())
	return impl.Dnrm2(x.N, x.Data, x.Inc)
}

func Dasum(x Vector) float64 {
	must(x.Check())
	return impl.Dasum(x.N, x.Data, x.Inc)
}

func Idamax(x Vector) int {
	must(x.Check())
	return impl.Idamax(x.N, x.Data, x.Inc)
}

func Dswap(x, y Vector) {
	must(x.Check())
	must(y.Check())
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Dswap(x.N, x.Data, x.Inc, y.Data, y.Inc)
}

func Dcopy(x, y Vector) {
	must(x.Check())
	must(y.Check())
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Dcopy(x.N, x.Data, x.Inc, y.Data, y.Inc)
}

func Daxpy(alpha float64, x, y Vector) {
	must(x.Check())
	must(y.Check())
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Daxpy(x.N, alpha, x.Data, x.Inc, y.Data, y.Inc)
}

func Drotg(a, b float64) (c, s, r, z float64) {
	return impl.Drotg(a, b)
}

func Drotmg(d1, d2, b1, b2 float64) (p DrotmParams, rd1, rd2, rb1 float64) {
	return impl.Drotmg(d1, d2, b1, b2)
}

func Drot(x, y Vector, c, s float64) {
	must(x.Check())
	must(y.Check())
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Drot(x.N, x.Data, x.Inc, y.Data, y.Inc, c, s)
}

func Drotm(x, y Vector, p DrotmParams) {
	must(x.Check())
	must(y.Check())
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Drotm(x.N, x.Data, x.Inc, y.Data, y.Inc, p)
}

func Dscal(alpha float64, x Vector) {
	must(x.Check())
	impl.Dscal(x.N, alpha, x.Data, x.Inc)
}
