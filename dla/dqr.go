package dla

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/d"
)

type QRFact struct {
	a   d.General
	tau []float64
}

func QR(A d.General, tau []float64) QRFact {
	impl.Dgeqrf(A, tau)
	return QRFact{A, tau}
}

func (f QRFact) R() d.Triangular {
	return d.Ge2Tr(f.a, blas.NonUnit, blas.Upper)
}

func (f QRFact) Solve(B d.General) d.General {
	impl.Dormqr('L', 'T', f.a, f.tau, B)
	B.Rows = f.a.Cols
	f.R().SolveM(blas.Left, blas.NoTrans, 1, B)
	return B
}
