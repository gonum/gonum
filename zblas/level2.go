package zblas

import "github.com/gonum/blas"

func Zgemv(tA blas.Transpose, alpha complex128, A General, x Vector, beta complex128, y Vector) {
	if tA == blas.NoTrans {
		if x.N != A.Cols {
			panic("blas: dimension mismatch")
		}
	} else if tA == blas.Trans {
		if x.N != A.Rows {
			panic("blas: dimension mismatch")
		}
	} else {
		panic("blas: illegal value for tA")
	}

	impl.Zgemv(A.Order, tA, A.Rows, A.Cols, alpha, A.Data, A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Zgbmv(tA blas.Transpose, alpha complex128, A GeneralBand, x Vector, beta complex128, y Vector) {
	if tA == blas.NoTrans {
		if x.N != A.Cols {
			panic("blas: dimension mismatch")
		}
	} else if tA == blas.Trans {
		if x.N != A.Rows {
			panic("blas: dimension mismatch")
		}
	} else {
		panic("blas: illegal value for tA")
	}
	impl.Zgbmv(A.Order, tA, A.Rows, A.Cols, A.KL, A.KU, alpha, A.Data,
		A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Ztrmv(tA blas.Transpose, A Triangular, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Ztrmv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, A.Stride, x.Data, x.Inc)
}

func Ztbmv(tA blas.Transpose, A TriangularBand, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Ztbmv(A.Order, A.Uplo, tA, A.Diag, A.N, A.K, A.Data, A.Stride, x.Data, x.Inc)
}

func Ztpmv(tA blas.Transpose, A TriangularPacked, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Ztpmv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, x.Data, x.Inc)
}

func Ztrsv(tA blas.Transpose, A Triangular, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Ztrsv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, A.Stride, x.Data, x.Inc)
}

func Ztbsv(tA blas.Transpose, A TriangularBand, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Ztbsv(A.Order, A.Uplo, tA, A.Diag, A.N, A.K, A.Data, A.Stride, x.Data, x.Inc)
}
func Ztpsv(tA blas.Transpose, A TriangularPacked, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Ztpsv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, x.Data, x.Inc)
}

func Zhemv(alpha complex128, A Hermitian, x Vector, beta complex128, y Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Zhemv(A.Order, A.Uplo, A.N, alpha, A.Data, A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Zhbmv(alpha complex128, A HermitianBand, x Vector, beta complex128, y Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Zhbmv(A.Order, A.Uplo, A.N, A.K, alpha, A.Data, A.Stride, x.Data,
		x.Inc, beta, y.Data, y.Inc)
}

func Zhpmv(alpha complex128, A HermitianPacked, x Vector, beta complex128, y Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Zhpmv(A.Order, A.Uplo, A.N, alpha, A.Data, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Zgerc(alpha complex128, x Vector, y Vector, A General) {
	if x.N != A.Rows {
		panic("blas: dimension mismatch")
	}
	if y.N != A.Cols {
		panic("blas: dimension mismatch")
	}
	impl.Zgerc(A.Order, A.Rows, A.Cols, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data, A.Stride)
}

func Zgeru(alpha complex128, x Vector, y Vector, A General) {
	if x.N != A.Rows {
		panic("blas: dimension mismatch")
	}
	if y.N != A.Cols {
		panic("blas: dimension mismatch")
	}
	impl.Zgeru(A.Order, A.Rows, A.Cols, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data, A.Stride)
}

func Zher(alpha float64, x Vector, A Hermitian) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Zher(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, A.Data, A.Stride)
}

func Zhpr(alpha float64, x Vector, A HermitianPacked) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Zhpr(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, A.Data)
}

func Zher2(alpha complex128, x Vector, y Vector, A Hermitian) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	if y.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Zher2(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data, A.Stride)
}

func Zhpr2(alpha complex128, x Vector, y Vector, A HermitianPacked) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	if y.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Zhpr2(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data)
}
