package dbw

import "github.com/gonum/blas"

func Gemv(tA blas.Transpose, alpha float64, A General, x Vector, beta float64, y Vector) {
	if tA == blas.NoTrans {
		if x.N != A.Cols || y.N != A.Rows {
			panic("blas: dimension mismatch")
		}
	} else if tA == blas.Trans {
		if x.N != A.Rows || y.N != A.Cols {
			panic("blas: dimension mismatch")
		}
	} else {
		panic("blas: illegal value for tA")
	}
	impl.Dgemv(tA, A.Rows, A.Cols, alpha, A.Data, A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Gbmv(tA blas.Transpose, alpha float64, A GeneralBand, x Vector, beta float64, y Vector) {
	if tA == blas.NoTrans {
		if x.N != A.Cols || y.N != A.Rows {
			panic("blas: dimension mismatch")
		}
	} else if tA == blas.Trans {
		if x.N != A.Rows || y.N != A.Cols {
			panic("blas: dimension mismatch")
		}
	} else {
		panic("blas: illegal value for tA")
	}
	impl.Dgbmv(tA, A.Rows, A.Cols, A.KL, A.KU, alpha, A.Data,
		A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Trmv(tA blas.Transpose, A Triangular, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dtrmv(A.Uplo, tA, A.Diag, A.N, A.Data, A.Stride, x.Data, x.Inc)
}

func Tbmv(tA blas.Transpose, A TriangularBand, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dtbmv(A.Uplo, tA, A.Diag, A.N, A.K, A.Data, A.Stride, x.Data, x.Inc)
}

func Tpmv(tA blas.Transpose, A TriangularPacked, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dtpmv(A.Uplo, tA, A.Diag, A.N, A.Data, x.Data, x.Inc)
}

func Trsv(tA blas.Transpose, A Triangular, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dtrsv(A.Uplo, tA, A.Diag, A.N, A.Data, A.Stride, x.Data, x.Inc)
}

func Tbsv(tA blas.Transpose, A TriangularBand, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dtbsv(A.Uplo, tA, A.Diag, A.N, A.K, A.Data, A.Stride, x.Data, x.Inc)
}

func Tpsv(tA blas.Transpose, A TriangularPacked, x Vector) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dtpsv(A.Uplo, tA, A.Diag, A.N, A.Data, x.Data, x.Inc)
}

func Symv(alpha float64, A Symmetric, x Vector, beta float64, y Vector) {
	if x.N != A.N || y.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dsymv(A.Uplo, A.N, alpha, A.Data, A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Sbmv(alpha float64, A SymmetricBand, x Vector, beta float64, y Vector) {
	if x.N != A.N || y.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dsbmv(A.Uplo, A.N, A.K, alpha, A.Data, A.Stride, x.Data,
		x.Inc, beta, y.Data, y.Inc)
}

func Spmv(alpha float64, A SymmetricPacked, x Vector, beta float64, y Vector) {
	if x.N != A.N || y.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dspmv(A.Uplo, A.N, alpha, A.Data, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Ger(alpha float64, x Vector, y Vector, A General) {
	if x.N != A.Rows {
		panic("blas: dimension mismatch")
	}
	if y.N != A.Cols {
		panic("blas: dimension mismatch")
	}
	impl.Dger(A.Rows, A.Cols, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data, A.Stride)
}

func Syr(alpha float64, x Vector, A Symmetric) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dsyr(A.Uplo, A.N, alpha, x.Data, x.Inc, A.Data, A.Stride)
}

func Spr(alpha float64, x Vector, A SymmetricPacked) {
	if x.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dspr(A.Uplo, A.N, alpha, x.Data, x.Inc, A.Data)
}

func Syr2(alpha float64, x Vector, y Vector, A Symmetric) {
	if x.N != A.N || y.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dsyr2(A.Uplo, A.N, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data, A.Stride)
}

func Spr2(alpha float64, x Vector, y Vector, A SymmetricPacked) {
	if x.N != A.N || y.N != A.N {
		panic("blas: dimension mismatch")
	}
	impl.Dspr2(A.Uplo, A.N, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data)
}
