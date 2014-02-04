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
	impl.Dgeqrf(A, tau)
	return QRFact{A, tau}
}

func (f QRFact) R() dbw.Triangular {
	return dbw.Ge2Tr(f.a, blas.NonUnit, blas.Upper)
}

func (f QRFact) Solve(B dbw.General) dbw.General {
	impl.Dormqr(blas.Left, blas.Trans, f.a, f.tau, B)
	B.Rows = f.a.Cols
	dbw.Trsm(blas.Left, blas.NoTrans, 1, f.R(), B)
	return B
}
