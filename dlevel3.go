package blas

func Dgemm(tA, tB Transpose, alpha float64, A, B General, beta float64, C General) {
	var m, n, k int
	if tA == NoTrans {
		m, k = A.Rows, A.Cols
	} else {
		m, k = A.Cols, A.Rows
	}
	if tB == NoTrans {
		n = B.Cols
	} else {
		n = B.Rows
	}
	impl.Dgemm(A.Order, tA, tB, m, n, k, alpha, A.Data, A.Stride,
		B.Data, B.Stride, beta, C.Data, C.Stride)
}

func Dsymm(s Side, alpha float64, A Symmetric, B General, beta float64, C General) {
	var m, n int
	if s == Left {
		m = A.N
		n = B.Cols
	} else {
		m = B.Rows
		n = A.N
	}
	impl.Dsymm(A.Order, s, A.Uplo, m, n, alpha, A.Data, A.Stride,
		B.Data, B.Stride, beta, C.Data, C.Stride)
}

func Dsyrk(t Transpose, alpha float64, A General, beta float64, C Symmetric) {
	var n, k int
	if t == NoTrans {
		n, k = A.Rows, A.Cols
	} else {
		n, k = A.Cols, A.Rows
	}
	impl.Dsyrk(A.Order, C.Uplo, t, n, k, alpha, A.Data, A.Stride, beta,
		C.Data, C.Stride)
}

func Dsyr2k(t Transpose, alpha float64, A, B General, beta float64, C Symmetric) {
	var n, k int
	if t == NoTrans {
		n, k = A.Rows, A.Cols
	} else {
		n, k = A.Cols, A.Rows
	}
	impl.Dsyr2k(A.Order, C.Uplo, t, n, k, alpha, A.Data, A.Stride,
		B.Data, B.Stride, beta, C.Data, C.Stride)
}

func Dtrmm(s Side, tA Transpose, alpha float64, A Triangular, B General) {
	impl.Dtrmm(A.Order, s, A.Uplo, tA, A.Diag, B.Rows, B.Cols, alpha, A.Data, A.Stride,
		B.Data, B.Stride)
}

func Dtrsm(s Side, tA Transpose, alpha float64, A Triangular, B General) {
	impl.Dtrsm(A.Order, s, A.Uplo, tA, A.Diag, B.Rows, B.Cols, alpha, A.Data, A.Stride,
		B.Data, B.Stride)
}
