package blas

func Dgemv(tA Transpose, alpha float64, A General, x Vector, beta float64, y Vector) {
	if tA == NoTrans {
		if x.N != A.Cols {
			panic("blas: dimension mismatch")
		}
	} else if tA == Trans {
		if x.N != A.Rows {
			panic("blas: dimension mismatch")
		}
	} else {
		panic("blas: illegal value for tA")
	}

	impl.Dgemv(A.Order, tA, A.Rows, A.Cols, alpha, A.Data, A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Dgbmv(tA Transpose, alpha float64, A GeneralBand, x Vector, beta float64, y Vector) {
	impl.Dgbmv(A.Order, tA, A.Rows, A.Cols, A.KL, A.KU, alpha, A.Data,
		A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Dtrmv(tA Transpose, A Triangular, x Vector) {
	impl.Dtrmv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, A.Stride, x.Data, x.Inc)
}

func Dtbmv(tA Transpose, A TriangularBand, x Vector) {
	impl.Dtbmv(A.Order, A.Uplo, tA, A.Diag, A.N, A.K, A.Data, A.Stride, x.Data, x.Inc)
}

func Dtpmv(tA Transpose, A TriangularPacked, x Vector) {
	impl.Dtpmv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, x.Data, x.Inc)
}

func Dtrsv(tA Transpose, A Triangular, x Vector) {
	impl.Dtrsv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, A.Stride, x.Data, x.Inc)
}

func Dtbsv(tA Transpose, A TriangularBand, x Vector) {
	impl.Dtbsv(A.Order, A.Uplo, tA, A.Diag, A.N, A.K, A.Data, A.Stride, x.Data, x.Inc)
}
func Dtpsv(tA Transpose, A TriangularPacked, x Vector) {
	impl.Dtpsv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, x.Data, x.Inc)
}

func Dsymv(alpha float64, A Symmetric, x Vector, beta float64, y Vector) {
	impl.Dsymv(A.Order, A.Uplo, A.N, alpha, A.Data, A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Dsbmv(alpha float64, A SymmetricBand, x Vector, beta float64, y Vector) {
	impl.Dsbmv(A.Order, A.Uplo, A.N, A.K, alpha, A.Data, A.Stride, x.Data,
		x.Inc, beta, y.Data, y.Inc)
}

func Dspmv(alpha float64, A SymmetricPacked, x Vector, beta float64, y Vector) {
	impl.Dspmv(A.Order, A.Uplo, A.N, alpha, A.Data, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Dger(alpha float64, x Vector, y Vector, A General) {
	impl.Dger(A.Order, A.Rows, A.Cols, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data, A.Stride)
}

func Dsyr(alpha float64, x Vector, A Symmetric) {
	impl.Dsyr(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, A.Data, A.Stride)
}

func Dspr(alpha float64, x Vector, A SymmetricPacked) {
	impl.Dspr(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, A.Data)
}

func Dsyr2(alpha float64, x Vector, y Vector, A Symmetric) {
	impl.Dsyr2(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data, A.Stride)
}

func Dspr2(alpha float64, x Vector, y Vector, A SymmetricPacked) {
	impl.Dspr2(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data)
}
