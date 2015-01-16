package zla

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/cblas128"
)

type QRFact struct {
	a   cblas128.General
	tau []complex128
}

func QR(A cblas128.General, tau []complex128) QRFact {
	impl.Zgeqrf(A.Rows, A.Cols, A.Data, A.Stride, tau)
	return QRFact{A, tau}
}

func (f QRFact) R() cblas128.Triangular {
	n := f.a.Rows
	if f.a.Cols < n {
		n = f.a.Cols
	}
	return cblas128.Triangular{
		Data:   f.a.Data,
		N:      n,
		Stride: f.a.Stride,
		Uplo:   blas.Upper,
		Diag:   blas.NonUnit,
	}
}

func (f QRFact) Solve(B cblas128.General) cblas128.General {
	if f.a.Cols != B.Cols {
		panic("dimension missmatch")
	}
	impl.Zunmqr(blas.Left, blas.ConjTrans, f.a.Rows, B.Cols, f.a.Cols, f.a.Data, f.a.Stride, f.tau, B.Data, B.Stride)
	B.Rows = f.a.Cols
	cblas128.Trsm(blas.Left, blas.NoTrans, 1, f.R(), B)
	return B
}
