package goblas

import "github.com/gonum/blas"

var _ blas.Float64Level3 = Blasser

func (bl Blas) Dtrsm(s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int) {
	// Transform to row major
	if ul == blas.Upper {
		ul = blas.Lower
	} else {
		ul = blas.Upper
	}
	if s == blas.Left {
		s = blas.Right
	} else {
		s = blas.Left
	}
	m, n = n, m

	var nrowa int
	if s == blas.Left {
		nrowa = m
	} else {
		nrowa = n
	}

	if s != blas.Left && s != blas.Right {
		panic(badSide)
	}
	if ul != blas.Lower && ul != blas.Upper {
		panic(badUplo)
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic(badTranspose)
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic(badDiag)
	}
	if m < 0 {
		panic(mLT0)
	}
	if n < 0 {
		panic(nLT0)
	}
	if lda < nrowa {
		panic(badLda)
	}
	if ldb < m {
		panic(badLda)
	}

	if m == 0 || n == 0 {
		return
	}

	if alpha == 0 {
		for j := 0; j < n; j++ {
			jb := ldb * j
			for i := 0; i < m; i++ {
				b[i+jb] = 0
			}
		}
		return
	}

	if s == blas.Left {
		for j := 0; j < n; j++ {
			jb := j * ldb
			if alpha != 1 {
				bl.Dscal(m, alpha, b[jb:], 1)
			}
			bl.Dtrsv(ul, tA, d, m, a, lda, b[jb:], 1)
		}
	} else {
		if tA == blas.NoTrans {

			//B := alpha*B*inv( A )
			if ul == blas.Upper {
				for j := 0; j < n; j++ {
					jb := j * ldb
					ja := j * lda
					if alpha != 1 {
						bl.Dscal(m, alpha, b[jb:], 1)
					}
					for k := 0; k < j; k++ {
						if a[k+ja] != 0 {
							bl.Daxpy(m, -a[k+ja], b[k*ldb:], 1, b[jb:], 1)
						}
					}
					if d == blas.NonUnit {
						bl.Dscal(m, 1/a[j+ja], b[jb:], 1)
					}
				}
			} else {
				for j := n - 1; j >= 0; j-- {
					jb := j * ldb
					ja := j * lda
					if alpha != 1 {
						bl.Dscal(m, alpha, b[jb:], 1)
					}
					for k := j + 1; k < n; k++ {
						if a[k+ja] != 0 {
							bl.Daxpy(m, -a[k+ja], b[k*ldb:], 1, b[jb:], 1)
						}
					}
					if d == blas.NonUnit {
						bl.Dscal(m, 1/a[j+ja], b[jb:], 1)
					}
				}
			}
		} else {

			//B := alpha*B*inv( A**T )
			if ul == blas.Upper {
				for k := n - 1; k >= 0; k-- {
					ka := k * lda
					kb := k * ldb
					if d == blas.NonUnit {
						bl.Dscal(m, 1/a[k+ka], b[kb:], 1)
					}
					for j := 0; j < k; j++ {
						if a[j+ka] != 0 {
							bl.Daxpy(m, a[j+ka], b[kb:], 1, b[j*ldb:], 1)
						}
					}
					if alpha != 1 {
						bl.Dscal(m, alpha, b[kb:], 1)
					}
				}
			} else {
				for k := 0; k < n; k++ {
					ka := k * lda
					kb := k * ldb
					if d == blas.NonUnit {
						bl.Dscal(m, 1/a[k+ka], b[kb:], 1)
					}
					for j := k + 1; j < n; j++ {
						if a[j+ka] != 0 {
							bl.Daxpy(m, a[j+ka], b[kb:], 1, b[j*ldb:], 1)
						}
					}
					if alpha != 1 {
						bl.Dscal(m, alpha, b[kb:], 1)
					}
				}
			}
		}
	}
}

// Dsymm performs one of
//  C = alpha * A * B + beta * C
//  C = alpha * B * A + beta * C
// where A is a symmetric matrix and B and C are m x n matrices.
func (Blas) Dsymm(s blas.Side, ul blas.Uplo, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {
	if s != blas.Right && s != blas.Left {
		panic("goblas: bad side")
	}
	if ul != blas.Lower && ul != blas.Upper {
		panic(badUplo)
	}
	if alpha == 0 && beta == 0 {
		return
	}
	if beta != 0 {
		for i := 0; i < m; i++ {
			ctmp := c[i*lda : i*lda+n]
			for j := range ctmp {
				ctmp[j] *= beta
			}
		}
	}

	if s == blas.Right {
		if ul == blas.Upper {
			for i := 0; i < m; i++ {
				for j := 0; j < n; j++ {
					for k := i + 1; k < m; k++ {
						v := alpha * a[i*lda+k] * b[k*ldb+j]
						c[i*ldc+j] += v
						if i != j {
							c[j*ldc+i] += v
						}
					}
				}
			}
		}
	}
}
func (Blas) Dsyrk(ul blas.Uplo, t blas.Transpose, n, k int, alpha float64, a []float64, lda int, beta float64, c []float64, ldc int) {
	panic("blas: function not implemented")
}
func (Blas) Dsyr2k(ul blas.Uplo, t blas.Transpose, n, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {
	panic("blas: function not implemented")
}
func (Blas) Dtrmm(s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int) {
	panic("blas: function not implemented")
}
