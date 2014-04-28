package dla

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/dbw"
)

type QRFact struct {
	a   dbw.General
	tau []float64
}

func QR(A dbw.General, tau []float64) QRFact {
	impl.Dgeqrf(A.Order, A.Rows, A.Cols, A.Data, A.Stride, tau)
	return QRFact{A, tau}
}

func (f QRFact) R() dbw.Triangular {
	return dbw.Ge2Tr(f.a, blas.NonUnit, blas.Upper)
}

func (f QRFact) Solve(B dbw.General) dbw.General {
	if B.Order != f.a.Order {
		panic("Order missmatch")
	}
	if f.a.Cols != B.Cols {
		panic("dimension missmatch")
	}
	impl.Dormqr(B.Order, blas.Left, blas.Trans, f.a.Rows, B.Cols, f.a.Cols, f.a.Data, f.a.Stride, f.tau, B.Data, B.Stride)
	B.Rows = f.a.Cols
	dbw.Trsm(blas.Left, blas.NoTrans, 1, f.R(), B)
	return B
}
