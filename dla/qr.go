package dla

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

type QRFact struct {
	a   blas64.General
	tau []float64
}

func QR(a blas64.General, tau []float64) QRFact {
	impl.Dgeqrf(a.Rows, a.Cols, a.Data, a.Stride, tau)
	return QRFact{a: a, tau: tau}
}

func (f QRFact) R() blas64.Triangular {
	n := f.a.Rows
	if f.a.Cols < n {
		n = f.a.Cols
	}
	return blas64.Triangular{
		Data:   f.a.Data,
		N:      n,
		Stride: f.a.Stride,
		Uplo:   blas.Upper,
		Diag:   blas.NonUnit,
	}
}

func (f QRFact) Solve(b blas64.General) blas64.General {
	if f.a.Cols != b.Cols {
		panic("dimension missmatch")
	}
	impl.Dormqr(blas.Left, blas.Trans, b.Rows, b.Cols, f.a.Cols, f.a.Data, f.a.Stride, f.tau, b.Data, b.Stride)
	b.Rows = f.a.Cols
	blas64.Trsm(blas.Left, blas.NoTrans, 1, f.R(), b)
	return b
}
