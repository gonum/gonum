package dbw

import "github.com/gonum/blas"

func Dot(x, y Vector) float64 {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	return impl.Ddot(x.N, x.Data, x.Inc, y.Data, y.Inc)
}

func Nrm2(x Vector) float64 {
	return impl.Dnrm2(x.N, x.Data, x.Inc)
}

func Asum(x Vector) float64 {
	return impl.Dasum(x.N, x.Data, x.Inc)
}

func Iamax(x Vector) int {
	return impl.Idamax(x.N, x.Data, x.Inc)
}

func Swap(x, y Vector) {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Dswap(x.N, x.Data, x.Inc, y.Data, y.Inc)
}

func Copy(x, y Vector) {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Dcopy(x.N, x.Data, x.Inc, y.Data, y.Inc)
}

func Axpy(alpha float64, x, y Vector) {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Daxpy(x.N, alpha, x.Data, x.Inc, y.Data, y.Inc)
}

func Rotg(a, b float64) (c, s, r, z float64) {
	return impl.Drotg(a, b)
}

func Rotmg(d1, d2, b1, b2 float64) (p blas.DrotmParams, rd1, rd2, rb1 float64) {
	return impl.Drotmg(d1, d2, b1, b2)
}

func Rot(x, y Vector, c, s float64) {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Drot(x.N, x.Data, x.Inc, y.Data, y.Inc, c, s)
}

func Rotm(x, y Vector, p blas.DrotmParams) {
	if x.N != y.N {
		panic("blas: dimension mismatch")
	}
	impl.Drotm(x.N, x.Data, x.Inc, y.Data, y.Inc, p)
}

func Scal(alpha float64, x Vector) {
	impl.Dscal(x.N, alpha, x.Data, x.Inc)
}
