package zbw

import "github.com/gonum/blas"

func Gemv(tA blas.Transpose, alpha complex128, A General, x Vector, beta complex128, y Vector) {
	if tA == blas.NoTrans {
		if x.N != A.Cols || y.N != A.Rows {
			panic("blas: dimension mismatch")
		}
	} else if tA == blas.ConjTrans || tA == blas.ConjTrans {
		if x.N != A.Rows || y.N != A.Cols {
			panic("blas: dimension mismatch")
		}
	} else {
		panic("blas: illegal value for tA")
	}

	impl.Zgemv(tA, A.Rows, A.Cols, alpha, A.Data, A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Gbmv(tA blas.Transpose, alpha complex128, A GeneralBand, x Vector, beta complex128, y Vector) {
	if tA == blas.NoTrans {
		if x.N != A.Cols || y.N != A.Rows {
			panic("blas: dimension mismatch")
		}
	} else if tA == blas.ConjTrans {
		if x.N != A.Rows || y.N != A.Cols {
			panic("blas: dimension mismatch")
		}
	} else {
		panic("blas: illegal value for tA")
	}
	impl.Zgbmv(tA, A.Rows, A.Cols, A.KL, A.KU, alpha, A.Data,
		A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Trmv(tA blas.Transpose, A Triangular, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Ztrmv(A.Uplo, tA, A.Diag, A.N, A.Data, A.Stride, x.Data, x.Inc)
}

func Tbmv(tA blas.Transpose, A TriangularBand, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Ztbmv(A.Uplo, tA, A.Diag, A.N, A.K, A.Data, A.Stride, x.Data, x.Inc)
}

func Tpmv(tA blas.Transpose, A TriangularPacked, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Ztpmv(A.Uplo, tA, A.Diag, A.N, A.Data, x.Data, x.Inc)
}

func Trsv(tA blas.Transpose, A Triangular, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Ztrsv(A.Uplo, tA, A.Diag, A.N, A.Data, A.Stride, x.Data, x.Inc)
}

func Tbsv(tA blas.Transpose, A TriangularBand, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Ztbsv(A.Uplo, tA, A.Diag, A.N, A.K, A.Data, A.Stride, x.Data, x.Inc)
}
func Tpsv(tA blas.Transpose, A TriangularPacked, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Ztpsv(A.Uplo, tA, A.Diag, A.N, A.Data, x.Data, x.Inc)
}

func Hemv(alpha complex128, A Hermitian, x Vector, beta complex128, y Vector) {
	if x.N != A.N || y.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Zhemv(A.Uplo, A.N, alpha, A.Data, A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Hbmv(alpha complex128, A HermitianBand, x Vector, beta complex128, y Vector) {
	if x.N != A.N || y.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Zhbmv(A.Uplo, A.N, A.K, alpha, A.Data, A.Stride, x.Data,
		x.Inc, beta, y.Data, y.Inc)
}

func Hpmv(alpha complex128, A HermitianPacked, x Vector, beta complex128, y Vector) {
	if x.N != A.N || y.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Zhpmv(A.Uplo, A.N, alpha, A.Data, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Zgerc(alpha complex128, x Vector, y Vector, A General) {
	if x.N != A.Rows || y.N != A.Cols {
		panic("blas: dimension mismatch")
	}
	impl.Zgerc(A.Rows, A.Cols, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data, A.Stride)
}

func Zgeru(alpha complex128, x Vector, y Vector, A General) {
	if x.N != A.Rows || y.N != A.Cols {
		panic("blas: dimension mismatch")
	}
	impl.Zgeru(A.Rows, A.Cols, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data, A.Stride)
}

func Zher(alpha float64, x Vector, A Hermitian) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Zher(A.Uplo, A.N, alpha, x.Data, x.Inc, A.Data, A.Stride)
}

func Zhpr(alpha float64, x Vector, A HermitianPacked) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Zhpr(A.Uplo, A.N, alpha, x.Data, x.Inc, A.Data)
}

func Zher2(alpha complex128, x Vector, y Vector, A Hermitian) {
	if x.N != A.N || y.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Zher2(A.Uplo, A.N, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data, A.Stride)
}

func Zhpr2(alpha complex128, x Vector, y Vector, A HermitianPacked) {
	if x.N != A.N || y.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Zhpr2(A.Uplo, A.N, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data)
}
