package referenceblas

import "github.com/gonum/blas"

func (Blas) Dgemm(o blas.Order, tA, tB blas.Transpose, m, n, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {

	nota := tA == blas.NoTrans
	notb := tB == blas.NoTrans

	if o != blas.RowMajor && o != blas.ColMajor {
		panic(badOrder)
	}

	if nota && tA != blas.ConjTrans && tA != blas.Trans {
		panic(badTranspose)
	}
	if nota && tB != blas.ConjTrans && tB != blas.Trans {
		panic(badTranspose)
	}
	if m < 0 {
		panic(mLT0)
	}
	if n < 0 {
		panic(nLT0)
	}
	if m == 0 || n == 0 || ((alpha == 0 || k == 0) && beta == 1) {
		return
	}

	if alpha == 0 {
		if beta == 0 {
			if o == blas.RowMajor {
				for i := 0; i < m; i++ {
					for j := 0; j < n; j++ {
						c[i*ldac+j] = 0
					}
				}
				return
			}
			for j := 0; j < n; j++ {
				for i := 0; i < m; i++ {
					c[j*ldac+i] = 0
				}
			}
			return
		}
		if o == blas.RowMajor {
			for i := 0; i < m; i++ {
				for j := 0; j < n; j++ {
					c[i*ldac+j] *= beta
				}
			}
			return
		}
		for j := 0; j < n; j++ {
			for i := 0; i < m; i++ {
				c[j*ldac+i] *= beta
			}
		}
		return
	}

	switch {
	case o == blas.RowMajor && nota && notb:

	}
}

// Dsymm(o Order, s Side, ul Uplo, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int)
// Dsyrk(o Order, ul Uplo, t Transpose, n, k int, alpha float64, a []float64, lda int, beta float64, c []float64, ldc int)
// Dsyr2k(o Order, ul Uplo, t Transpose, n, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int)
// Dtrmm(o Order, s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int)
// Dtrsm(o Order, s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int)
