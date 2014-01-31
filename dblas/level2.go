package dblas

import "github.com/gonum/blas"

func Dgemv(tA blas.Transpose, alpha float64, A General, x Vector, beta float64, y Vector) {
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
	impl.Dgemv(A.Order, tA, A.Rows, A.Cols, alpha, A.Data, A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Dgbmv(tA blas.Transpose, alpha float64, A GeneralBand, x Vector, beta float64, y Vector) {
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
	impl.Dgbmv(A.Order, tA, A.Rows, A.Cols, A.KL, A.KU, alpha, A.Data,
		A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Dtrmv(tA blas.Transpose, A Triangular, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dtrmv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, A.Stride, x.Data, x.Inc)
}

func Dtbmv(tA blas.Transpose, A TriangularBand, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dtbmv(A.Order, A.Uplo, tA, A.Diag, A.N, A.K, A.Data, A.Stride, x.Data, x.Inc)
}

func Dtpmv(tA blas.Transpose, A TriangularPacked, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dtpmv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, x.Data, x.Inc)
}

func Dtrsv(tA blas.Transpose, A Triangular, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dtrsv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, A.Stride, x.Data, x.Inc)
}

func Dtbsv(tA blas.Transpose, A TriangularBand, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dtbsv(A.Order, A.Uplo, tA, A.Diag, A.N, A.K, A.Data, A.Stride, x.Data, x.Inc)
}
func Dtpsv(tA blas.Transpose, A TriangularPacked, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dtpsv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, x.Data, x.Inc)
}

func Dsymv(alpha float64, A Symmetric, x Vector, beta float64, y Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dsymv(A.Order, A.Uplo, A.N, alpha, A.Data, A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Dsbmv(alpha float64, A SymmetricBand, x Vector, beta float64, y Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dsbmv(A.Order, A.Uplo, A.N, A.K, alpha, A.Data, A.Stride, x.Data,
		x.Inc, beta, y.Data, y.Inc)
}

func Dspmv(alpha float64, A SymmetricPacked, x Vector, beta float64, y Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dspmv(A.Order, A.Uplo, A.N, alpha, A.Data, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Dger(alpha float64, x Vector, y Vector, A General) {
	if x.N != A.Rows {
		panic("blas: dimension mismatch")
	}
	if y.N != A.Cols {
		panic("blas: dimension mismatch")
	}
	impl.Dger(A.Order, A.Rows, A.Cols, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data, A.Stride)
}

func Dsyr(alpha float64, x Vector, A Symmetric) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dsyr(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, A.Data, A.Stride)
}

func Dspr(alpha float64, x Vector, A SymmetricPacked) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dspr(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, A.Data)
}

func Dsyr2(alpha float64, x Vector, y Vector, A Symmetric) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	if y.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dsyr2(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data, A.Stride)
}

func Dspr2(alpha float64, x Vector, y Vector, A SymmetricPacked) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	if y.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dspr2(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data)
}
