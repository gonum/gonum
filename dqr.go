package lapack

import "github.com/gonum/blas"

type QRFact struct {
	A   blas.General
	tau []float64
}

func QR(A blas.General, tau []float64) QRFact {
	impl.Dgeqrf(A, tau)
	return QRFact{A, tau}
}

func (f QRFact) Solve(B blas.General) blas.General {
	impl.Dormqr(blas.Left, blas.Trans, f.A, f.tau, B)
	blas.Dtrsm(blas.Left, blas.NoTrans, 1,
		blas.Ge2Tr(f.A, blas.NonUnit, blas.Upper), B)
	B.Rows = f.A.Cols
	return B
}
