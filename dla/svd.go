package dla

import (
	"github.com/dane-unltd/lapack"
	"github.com/gonum/blas"
	"github.com/gonum/blas/dbw"
)

func SVD(A dbw.General) (U dbw.General, s []float64, Vt dbw.General) {
	m := A.Rows
	n := A.Cols
	U.Stride = 1
	Vt.Stride = 1
	if m >= n {
		Vt = dbw.NewGeneral(n, n, nil)
		s = make([]float64, n)
		U = A
	} else {
		U = dbw.NewGeneral(m, m, nil)
		s = make([]float64, m)
		Vt = A
	}

	impl.Dgesdd(lapack.Overwrite, A.Rows, A.Cols, A.Data, A.Stride, s, U.Data, U.Stride, Vt.Data, Vt.Stride)

	return
}

func SVDbd(uplo blas.Uplo, d, e []float64) (U dbw.General, s []float64, Vt dbw.General) {
	n := len(d)
	if len(e) != n {
		panic("dimensionality missmatch")
	}

	U = dbw.NewGeneral(n, n, nil)
	Vt = dbw.NewGeneral(n, n, nil)

	impl.Dbdsdc(uplo, lapack.Explicit, n, d, e, U.Data, U.Stride, Vt.Data, Vt.Stride, nil, nil)
	s = d
	return
}
