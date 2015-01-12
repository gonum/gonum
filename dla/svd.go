package dla

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/lapack"
)

func SVD(A blas64.General) (U blas64.General, s []float64, Vt blas64.General) {
	m := A.Rows
	n := A.Cols
	U.Stride = 1
	Vt.Stride = 1
	if m >= n {
		Vt = blas64.General{
			Rows:   n,
			Cols:   n,
			Stride: n,
			Data:   make([]float64, n*n),
		}
		s = make([]float64, n)
		U = A
	} else {
		U = blas64.General{
			Rows:   m,
			Cols:   m,
			Stride: m,
			Data:   make([]float64, n*n),
		}
		s = make([]float64, m)
		Vt = A
	}

	impl.Dgesdd(lapack.Overwrite, A.Rows, A.Cols, A.Data, A.Stride, s, U.Data, U.Stride, Vt.Data, Vt.Stride)

	return
}

func SVDbd(uplo blas.Uplo, d, e []float64) (U blas64.General, s []float64, Vt blas64.General) {
	n := len(d)
	if len(e) != n {
		panic("dimensionality missmatch")
	}

	U = blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, n*n),
	}
	Vt = blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, n*n),
	}

	impl.Dbdsdc(uplo, lapack.Explicit, n, d, e, U.Data, U.Stride, Vt.Data, Vt.Stride, nil, nil)
	s = d
	return
}
