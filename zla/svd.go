package zla

import (
	"github.com/dane-unltd/lapack"
	"github.com/gonum/blas"
	"github.com/gonum/blas/zbw"
)

func SVD(A zbw.General) (U zbw.General, s []float64, Vt zbw.General) {
	m := A.Rows
	n := A.Cols
	U.Stride = 1
	Vt.Stride = 1
	if m >= n {
		Vt = zbw.NewGeneral(A.Order, n, n, nil)
		s = make([]float64, n)
		U = A
	} else {
		U = zbw.NewGeneral(A.Order, m, m, nil)
		s = make([]float64, m)
		Vt = A
	}

	impl.Zgesdd(A.Order, lapack.Overwrite, A.Rows, A.Cols, A.Data, A.Stride, s, U.Data, U.Stride, Vt.Data, Vt.Stride)

	return
}

//Lanczos bidiagonalization with full reorthogonalization
func LanczosBi(L zbw.General, u []complex128, numIter int) (U zbw.General, V zbw.General, a []float64, b []float64) {

	m := L.Rows
	n := L.Cols

	uv := zbw.NewVector(u)
	zbw.Scal(complex(1/zbw.Nrm2(uv), 0), uv)

	U = zbw.NewGeneral(blas.ColMajor, m, numIter, nil)
	V = zbw.NewGeneral(blas.ColMajor, n, numIter, nil)

	a = make([]float64, numIter)
	b = make([]float64, numIter)

	zbw.Copy(uv, U.Col(0))

	tr := zbw.NewVector(zbw.Allocate(n))
	zbw.Gemv(blas.ConjTrans, 1, L, uv, 0, tr)
	a[0] = zbw.Nrm2(tr)
	zbw.Copy(tr, V.Col(0))
	zbw.Scal(complex(1/a[0], 0), V.Col(0))

	tl := zbw.NewVector(zbw.Allocate(m))
	for k := 0; k < numIter-1; k++ {
		zbw.Copy(U.Col(k), tl)
		zbw.Scal(complex(-a[k], 0), tl)
		zbw.Gemv(blas.NoTrans, 1, L, V.Col(k), 1, tl)

		for i := 0; i <= k; i++ {
			zbw.Axpy(-zbw.Dotc(U.Col(i), tl), U.Col(i), tl)
		}

		b[k] = zbw.Nrm2(tl)
		zbw.Copy(tl, U.Col(k+1))
		zbw.Scal(complex(1/b[k], 0), U.Col(k+1))

		zbw.Copy(V.Col(k), tr)
		zbw.Scal(complex(-b[k], 0), tr)
		zbw.Gemv(blas.ConjTrans, 1, L, U.Col(k+1), 1, tr)

		for i := 0; i <= k; i++ {
			zbw.Axpy(-zbw.Dotc(V.Col(i), tr), V.Col(i), tr)
		}

		a[k+1] = zbw.Nrm2(tr)
		zbw.Copy(tr, V.Col(k+1))
		zbw.Scal(complex(1/a[k+1], 0), V.Col(k+1))
	}
	return
}
