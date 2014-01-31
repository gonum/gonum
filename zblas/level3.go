package zblas

import "github.com/gonum/blas"

func Zgemm(tA, tB blas.Transpose, alpha complex128, A, B General, beta complex128, C General) {
	var m, n, k int
	if tA == blas.NoTrans {
		m, k = A.Rows, A.Cols
	} else {
		m, k = A.Cols, A.Rows
	}
	if tB == blas.NoTrans {
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
	impl.Zgemm(A.Order, tA, tB, m, n, k, alpha, A.Data, A.Stride,
		B.Data, B.Stride, beta, C.Data, C.Stride)
}

func Zsymm(s blas.Side, alpha complex128, A Symmetric, B General, beta complex128, C General) {
	var m, n int
	if s == blas.Left {
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
	impl.Zsymm(A.Order, s, A.Uplo, m, n, alpha, A.Data, A.Stride,
		B.Data, B.Stride, beta, C.Data, C.Stride)
}

func Zsyrk(t blas.Transpose, alpha complex128, A General, beta complex128, C Symmetric) {
	var n, k int
	if t == blas.NoTrans {
		n, k = A.Rows, A.Cols
	} else {
		n, k = A.Cols, A.Rows
	}
	if n != C.N {
		panic("blas: dimension mismatch")
	}
	impl.Zsyrk(A.Order, C.Uplo, t, n, k, alpha, A.Data, A.Stride, beta,
		C.Data, C.Stride)
}

func Zsyr2k(t blas.Transpose, alpha complex128, A, B General, beta complex128, C Symmetric) {
	var n, k int
	if t == blas.NoTrans {
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
	impl.Zsyr2k(A.Order, C.Uplo, t, n, k, alpha, A.Data, A.Stride,
		B.Data, B.Stride, beta, C.Data, C.Stride)
}

func Ztrmm(s blas.Side, tA blas.Transpose, alpha complex128, A Triangular, B General) {
	if s == blas.Left {
		if A.N != B.Rows {
			panic("blas: dimension mismatch")
		}
	} else {
		if A.N != B.Cols {
			panic("blas: dimension mismatch")
		}
	}
	impl.Ztrmm(A.Order, s, A.Uplo, tA, A.Diag, B.Rows, B.Cols, alpha, A.Data, A.Stride,
		B.Data, B.Stride)
}

func Ztrsm(s blas.Side, tA blas.Transpose, alpha complex128, A Triangular, B General) {
	if s == blas.Left {
		if A.N != B.Rows {
			panic("blas: dimension mismatch")
		}
	} else {
		if A.N != B.Cols {
			panic("blas: dimension mismatch")
		}
	}
	impl.Ztrsm(A.Order, s, A.Uplo, tA, A.Diag, B.Rows, B.Cols, alpha, A.Data, A.Stride,
		B.Data, B.Stride)
}

func Zhemm(s blas.Side, alpha complex128, A Hermitian, B General, beta complex128, C General) {
	var m, n int
	if s == blas.Left {
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
	impl.Zhemm(A.Order, s, A.Uplo, m, n, alpha, A.Data, A.Stride,
		B.Data, B.Stride, beta, C.Data, C.Stride)
}

func Zherk(t blas.Transpose, alpha float64, A General, beta float64, C Hermitian) {
	var n, k int
	if t == blas.NoTrans {
		n, k = A.Rows, A.Cols
	} else {
		n, k = A.Cols, A.Rows
	}
	if n != C.N {
		panic("blas: dimension mismatch")
	}
	impl.Zherk(A.Order, C.Uplo, t, n, k, alpha, A.Data, A.Stride, beta,
		C.Data, C.Stride)
}

func Zher2k(t blas.Transpose, alpha complex128, A, B General, beta float64, C Hermitian) {
	var n, k int
	if t == blas.NoTrans {
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
	impl.Zher2k(A.Order, C.Uplo, t, n, k, alpha, A.Data, A.Stride,
		B.Data, B.Stride, beta, C.Data, C.Stride)
}
