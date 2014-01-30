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
		if k != B.Rows {
			panic("blas: dimension mismatch")
		}
	} else {
		n = B.Rows
		if k != B.Cols {
			panic("blas: dimension mismatch")
		}
	}
	if m != C.Rows {
		panic("blas: dimension mismatch")
	}
	if n != C.Cols {
		panic("blas: dimension mismatch")
	}
	impl.Dgemm(A.Order, tA, tB, m, n, k, alpha, A.Data, A.Stride,
		B.Data, B.Stride, beta, C.Data, C.Stride)
}

func Dsymm(s Side, alpha float64, A Symmetric, B General, beta float64, C General) {
	var m, n int
	if s == Left {
		m = A.N
		n = B.Cols
		if m != B.Rows {
			panic("blas: dimension mismatch")
		}
	} else {
		m = B.Rows
		n = A.N
		if n != B.Cols {
			panic("blas: dimension mismatch")
		}
	}
	if m != C.Rows {
		panic("blas: dimension mismatch")
	}
	if n != C.Cols {
		panic("blas: dimension mismatch")
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
	if n != C.N {
		panic("blas: dimension mismatch")
	}
	impl.Dsyrk(A.Order, C.Uplo, t, n, k, alpha, A.Data, A.Stride, beta,
		C.Data, C.Stride)
}

func Dsyr2k(t Transpose, alpha float64, A, B General, beta float64, C Symmetric) {
	var n, k int
	if t == NoTrans {
		n, k = A.Rows, A.Cols
		if n != B.Rows || k != B.Cols {
			panic("blas: dimension mismatch")
		}
	} else {
		n, k = A.Cols, A.Rows
		if k != B.Rows || n != B.Cols {
			panic("blas: dimension mismatch")
		}
	}
	if n != C.N {
		panic("blas: dimension mismatch")
	}
	impl.Dsyr2k(A.Order, C.Uplo, t, n, k, alpha, A.Data, A.Stride,
		B.Data, B.Stride, beta, C.Data, C.Stride)
}

func Dtrmm(s Side, tA Transpose, alpha float64, A Triangular, B General) {
	if s == Left {
		if A.N != B.Rows {
			panic("blas: dimension mismatch")
		}
	} else {
		if A.N != B.Cols {
			panic("blas: dimension mismatch")
		}
	}
	impl.Dtrmm(A.Order, s, A.Uplo, tA, A.Diag, B.Rows, B.Cols, alpha, A.Data, A.Stride,
		B.Data, B.Stride)
}

func Dtrsm(s Side, tA Transpose, alpha float64, A Triangular, B General) {
	if s == Left {
		if A.N != B.Rows {
			panic("blas: dimension mismatch")
		}
	} else {
		if A.N != B.Cols {
			panic("blas: dimension mismatch")
		}
	}
	impl.Dtrsm(A.Order, s, A.Uplo, tA, A.Diag, B.Rows, B.Cols, alpha, A.Data, A.Stride,
		B.Data, B.Stride)
}
