package blas

func Zgemv(tA Transpose, alpha complex128, A GeneralCmplx, x VectorCmplx, beta complex128, y VectorCmplx) {
	must(A.Check())
	must(x.Check())
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

	implCmplx.Zgemv(A.Order, tA, A.Rows, A.Cols, alpha, A.Data, A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Zgbmv(tA Transpose, alpha complex128, A GeneralCmplxBand, x VectorCmplx, beta complex128, y VectorCmplx) {
	implCmplx.Zgbmv(A.Order, tA, A.Rows, A.Cols, A.KL, A.KU, alpha, A.Data,
		A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Ztrmv(tA Transpose, A TriangularCmplx, x VectorCmplx) {
	implCmplx.Ztrmv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, A.Stride, x.Data, x.Inc)
}

func Ztbmv(tA Transpose, A TriangularCmplxBand, x VectorCmplx) {
	implCmplx.Ztbmv(A.Order, A.Uplo, tA, A.Diag, A.N, A.K, A.Data, A.Stride, x.Data, x.Inc)
}

func Ztpmv(tA Transpose, A TriangularCmplxPacked, x VectorCmplx) {
	implCmplx.Ztpmv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, x.Data, x.Inc)
}

func Ztrsv(tA Transpose, A TriangularCmplx, x VectorCmplx) {
	implCmplx.Ztrsv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, A.Stride, x.Data, x.Inc)
}

func Ztbsv(tA Transpose, A TriangularCmplxBand, x VectorCmplx) {
	implCmplx.Ztbsv(A.Order, A.Uplo, tA, A.Diag, A.N, A.K, A.Data, A.Stride, x.Data, x.Inc)
}
func Ztpsv(tA Transpose, A TriangularCmplxPacked, x VectorCmplx) {
	implCmplx.Ztpsv(A.Order, A.Uplo, tA, A.Diag, A.N, A.Data, x.Data, x.Inc)
}

func Zhemv(alpha complex128, A Hermitian, x VectorCmplx, beta complex128, y VectorCmplx) {
	implCmplx.Zhemv(A.Order, A.Uplo, A.N, alpha, A.Data, A.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Zhbmv(alpha complex128, A HermitianBand, x VectorCmplx, beta complex128, y VectorCmplx) {
	implCmplx.Zhbmv(A.Order, A.Uplo, A.N, A.K, alpha, A.Data, A.Stride, x.Data,
		x.Inc, beta, y.Data, y.Inc)
}

func Zhpmv(alpha complex128, A HermitianPacked, x VectorCmplx, beta complex128, y VectorCmplx) {
	implCmplx.Zhpmv(A.Order, A.Uplo, A.N, alpha, A.Data, x.Data, x.Inc, beta, y.Data, y.Inc)
}

func Zgerc(alpha complex128, x VectorCmplx, y VectorCmplx, A GeneralCmplx) {
	implCmplx.Zgerc(A.Order, A.Rows, A.Cols, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data, A.Stride)
}

func Zgeru(alpha complex128, x VectorCmplx, y VectorCmplx, A GeneralCmplx) {
	implCmplx.Zgeru(A.Order, A.Rows, A.Cols, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data, A.Stride)
}

func Zher(alpha float64, x VectorCmplx, A Hermitian) {
	implCmplx.Zher(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, A.Data, A.Stride)
}

func Zhpr(alpha float64, x VectorCmplx, A HermitianPacked) {
	implCmplx.Zhpr(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, A.Data)
}

func Zher2(alpha complex128, x VectorCmplx, y VectorCmplx, A Hermitian) {
	implCmplx.Zher2(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data, A.Stride)
}

func Zhpr2(alpha complex128, x VectorCmplx, y VectorCmplx, A HermitianPacked) {
	implCmplx.Zhpr2(A.Order, A.Uplo, A.N, alpha, x.Data, x.Inc, y.Data, y.Inc, A.Data)
}
