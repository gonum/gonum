package goblas

import (
	"github.com/gonum/blas"
	level1 "github.com/ziutek/blas"
)

func (Blas) Dtrsm(o blas.Order, s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int) {

	if o == blas.RowMajor && s == blas.Left && ul == blas.Lower && tA == blas.NoTrans {
		for i := 0; i < m; i++ {
			level1.Dscal(n, 1/a[i*lda+i], b[i*ldb:], 1)
			for j := i + 1; j < m; j++ {
				level1.Daxpy(n, -a[j*lda+i], b[i*ldb:], 1, b[j*ldb:], 1)
			}
		}
		return
	}
	if o == blas.RowMajor && s == blas.Left && ul == blas.Lower && tA == blas.Trans {
		for i := m - 1; i >= 0; i-- {
			level1.Dscal(n, 1/a[i*lda+i], b[i*ldb:], 1)
			for j := i - 1; j >= 0; j-- {
				level1.Daxpy(n, -a[(m-j)*lda+i], b[i*ldb:], 1, b[j*ldb:], 1)
			}
		}
		return
	}
	if o == blas.RowMajor && s == blas.Left && ul == blas.Upper && tA == blas.NoTrans {
		for i := m - 1; i >= 0; i-- {
			level1.Dscal(n, 1/a[i*lda+i], b[i*ldb:], 1)
			for j := i - 1; j >= 0; j-- {
				level1.Daxpy(n, -a[j*lda+i], b[i*ldb:], 1, b[j*ldb:], 1)
			}
		}
		return
	}
	if o == blas.RowMajor && s == blas.Left && ul == blas.Upper && tA == blas.Trans {
		for i := 0; i < m; i++ {
			level1.Dscal(n, 1/a[i*lda+i], b[i*ldb:], 1)
			for j := i + 1; j < m; j++ {
				level1.Daxpy(n, -a[(m-j)*lda+i], b[i*ldb:], 1, b[j*ldb:], 1)
			}
		}
		return
	}
	//ColMajor,Left,Low,NoTrans
	//RowMajor,Right,Up,NoTrans
	if o == blas.RowMajor && s == blas.Right && ul == blas.Upper && tA == blas.NoTrans {
		for i := 0; i < n; i++ {
			level1.Dscal(m, 1/a[i*lda+i], b[i:], ldb)
			for j := i + 1; j < m; j++ {
				level1.Daxpy(m, -a[i*lda+j], b[i:], ldb, b[j:], ldb)
			}
		}
		return
	}
	//ColMajor,Left,Up,NoTrans
	//RowMajor,Right,Low,NoTrans
	if o == blas.RowMajor && s == blas.Right && ul == blas.Lower && tA == blas.NoTrans {
		for i := n - 1; i >= 0; i-- {
			level1.Dscal(m, 1/a[i*lda+i], b[i:], ldb)
			for j := i - 1; j >= 0; j-- {
				level1.Daxpy(m, -a[i*lda+j], b[i:], ldb, b[j:], ldb)
			}
			return
		}
	}
	//ColMajor,Left,Low,Trans
	//RowMajor,Right,Up,Trans
	if o == blas.RowMajor && s == blas.Right && ul == blas.Upper && tA == blas.NoTrans {
		for i := n - 1; i >= 0; i-- {
			level1.Dscal(m, 1/a[i*lda+i], b[i:], ldb)
			for j := i - 1; j >= 0; j-- {
				level1.Daxpy(m, -a[i*lda+(m-j)], b[i:], ldb, b[j:], ldb)
			}
			return
		}
		return
	}
	//ColMajor,Left,Up,Trans
	//RowMajor,Right,Low,Trans
	if o == blas.RowMajor && s == blas.Right && ul == blas.Upper && tA == blas.NoTrans {
		for i := 0; i < n; i++ {
			level1.Dscal(m, 1/a[i*lda+i], b[i:], ldb)
			for j := i + 1; j < m; j++ {
				level1.Daxpy(m, -a[i*lda+(m-j)], b[i:], ldb, b[j:], ldb)
			}
		}
		return
	}
	panic("unreachable")
}
