//+build cblas

package zla

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/cblas128"
	"github.com/gonum/lapack"
)

func SVD(A cblas128.General) (U cblas128.General, s []float64, Vt cblas128.General) {
	m := A.Rows
	n := A.Cols
	U.Stride = 1
	Vt.Stride = 1
	if m >= n {
		Vt = cblas128.General{Rows: n, Cols: n, Stride: n, Data: make([]complex128, n*n)}
		s = make([]float64, n)
		U = A
	} else {
		U = cblas128.General{Rows: n, Cols: n, Stride: n, Data: make([]complex128, n*n)}
		s = make([]float64, m)
		Vt = A
	}

	impl.Zgesdd(lapack.Overwrite, A.Rows, A.Cols, A.Data, A.Stride, s, U.Data, U.Stride, Vt.Data, Vt.Stride)

	return
}

func c128col(i int, a cblas128.General) cblas128.Vector {
	return cblas128.Vector{
		Inc:  a.Stride,
		Data: a.Data[i:],
	}
}

//Lanczos bidiagonalization with full reorthogonalization
func LanczosBi(L cblas128.General, u []complex128, numIter int) (U cblas128.General, V cblas128.General, a []float64, b []float64) {

	m := L.Rows
	n := L.Cols

	uv := cblas128.Vector{Inc: 1, Data: u}
	cblas128.Scal(len(u), complex(1/cblas128.Nrm2(len(u), uv), 0), uv)

	U = cblas128.General{Rows: m, Cols: numIter, Stride: numIter, Data: make([]complex128, m*numIter)}
	V = cblas128.General{Rows: n, Cols: numIter, Stride: numIter, Data: make([]complex128, n*numIter)}

	a = make([]float64, numIter)
	b = make([]float64, numIter)

	cblas128.Copy(len(u), uv, c128col(0, U))

	tr := cblas128.Vector{Inc: 1, Data: make([]complex128, n)}
	cblas128.Gemv(blas.ConjTrans, 1, L, uv, 0, tr)
	a[0] = cblas128.Nrm2(n, tr)
	cblas128.Copy(n, tr, c128col(0, V))
	cblas128.Scal(n, complex(1/a[0], 0), c128col(0, V))

	tl := cblas128.Vector{Inc: 1, Data: make([]complex128, m)}
	for k := 0; k < numIter-1; k++ {
		cblas128.Copy(m, c128col(k, U), tl)
		cblas128.Scal(m, complex(-a[k], 0), tl)
		cblas128.Gemv(blas.NoTrans, 1, L, c128col(k, V), 1, tl)

		for i := 0; i <= k; i++ {
			cblas128.Axpy(m, -cblas128.Dotc(m, c128col(i, U), tl), c128col(i, U), tl)
		}

		b[k] = cblas128.Nrm2(m, tl)
		cblas128.Copy(m, tl, c128col(k+1, U))
		cblas128.Scal(m, complex(1/b[k], 0), c128col(k+1, U))

		cblas128.Copy(n, c128col(k, V), tr)
		cblas128.Scal(n, complex(-b[k], 0), tr)
		cblas128.Gemv(blas.ConjTrans, 1, L, c128col(k+1, U), 1, tr)

		for i := 0; i <= k; i++ {
			cblas128.Axpy(n, -cblas128.Dotc(n, c128col(i, V), tr), c128col(i, V), tr)
		}

		a[k+1] = cblas128.Nrm2(n, tr)
		cblas128.Copy(n, tr, c128col(k+1, V))
		cblas128.Scal(n, complex(1/a[k+1], 0), c128col(k+1, V))
	}
	return
}
