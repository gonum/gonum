package zla

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/zbw"
)

type QRFact struct {
	a   zbw.General
	tau []complex128
}

func QR(A zbw.General, tau []complex128) QRFact {
	impl.Zgeqrf(A.Order, A.Rows, A.Cols, A.Data, A.Stride, tau)
	return QRFact{A, tau}
}

func (f QRFact) R() zbw.Triangular {
	return zbw.Ge2Tr(f.a, blas.NonUnit, blas.Upper)
}

func (f QRFact) Solve(B zbw.General) zbw.General {
	if B.Order != f.a.Order {
		panic("Order missmatch")
	}
	if f.a.Cols != B.Cols {
		panic("dimension missmatch")
	}
	impl.Zunmqr(B.Order, blas.Left, blas.Trans, f.a.Rows, B.Cols, f.a.Cols, f.a.Data, f.a.Stride, f.tau, B.Data, B.Stride)
	B.Rows = f.a.Cols
	zbw.Trsm(blas.Left, blas.NoTrans, 1, f.R(), B)
	return B
}
