package goblas

import "github.com/gonum/blas"

var _ blas.Float64Level3 = Blasser

const (
	// TODO (btracey): Fix the ld panic messages to be consistent across the package
	badLd string = "goblas: ld must be greater than the number of columns"
)

// Dtrsm solves
//  A X = alpha B
// or
//  X A = alpha B
// where X and B are m x n matrices, and A is a unit or non unit upper or lower
// triangular matrix. The result is stored in place into B. No check is made
// that A is invertible.
func (bl Blas) Dtrsm(s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int) {
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
	if lda < n || ldb < n {
		panic(badLd)
	}

	if m == 0 || n == 0 {
		return
	}

	if alpha == 0 {
		for i := 0; i < m; i++ {
			btmp := b[i*ldb : i*ldb+n]
			for j := range btmp {
				btmp[j] = 0
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
	if m < 0 {
		panic(mLT0)
	}
	if n < 0 {
		panic(nLT0)
	}
	if (lda < m && s == blas.Left) || (lda < n && s == blas.Right) {
		panic(badLd)
	}
	if ldb < n || ldc < n {
		panic(badLd)
	}
	if m == 0 || n == 0 {
		return
	}
	if alpha == 0 && beta == 1 {
		return
	}
	if alpha == 0 {
		if beta == 0 {
			for i := 0; i < m; i++ {
				ctmp := c[i*ldc : i*ldc+n]
				for j := range ctmp {
					ctmp[j] = 0
				}
			}
			return
		}
		for i := 0; i < m; i++ {
			ctmp := c[i*ldc : i*ldc+n]
			for j := 0; j < n; j++ {
				ctmp[j] *= beta
			}
		}
		return
	}

	isUpper := ul == blas.Upper
	if s == blas.Left {
		for i := 0; i < m; i++ {
			atmp := alpha * a[i*lda+i]
			btmp := b[i*ldb : i*ldb+n]
			ctmp := c[i*ldc : i*ldc+n]
			for j, v := range btmp {
				ctmp[j] *= beta
				ctmp[j] += atmp * v
			}
			for k := 0; k < i; k++ {
				var atmp float64
				if isUpper {
					atmp = a[k*lda+i]
				} else {
					atmp = a[i*lda+k]
				}
				atmp *= alpha
				btmp := b[k*ldb : k*ldb+n]
				ctmp := c[i*ldc : i*ldc+n]
				for j, v := range btmp {
					ctmp[j] += atmp * v
				}
			}
			for k := i + 1; k < m; k++ {
				var atmp float64
				if isUpper {
					atmp = a[i*lda+k]
				} else {
					atmp = a[k*lda+i]
				}
				atmp *= alpha
				btmp := b[k*ldb : k*ldb+n]
				ctmp := c[i*ldc : i*ldc+n]
				for j, v := range btmp {
					ctmp[j] += atmp * v
				}
			}
		}
		return
	}
	if isUpper {
		for i := 0; i < m; i++ {
			for j := n - 1; j >= 0; j-- {
				tmp := alpha * b[i*ldb+j]
				var tmp2 float64
				atmp := a[j*lda+j+1 : j*lda+n]
				btmp := b[i*ldb+j+1 : i*ldb+n]
				ctmp := c[i*ldc+j+1 : i*ldc+n]
				for k, v := range atmp {
					ctmp[k] += tmp * v
					tmp2 += btmp[k] * v
				}
				c[i*ldc+j] *= beta
				c[i*ldc+j] += tmp*a[j*lda+j] + alpha*tmp2
			}
		}
		return
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			tmp := alpha * b[i*ldb+j]
			var tmp2 float64
			atmp := a[j*lda : j*lda+j]
			btmp := b[i*ldb : i*ldb+j]
			ctmp := c[i*ldc : i*ldc+j]
			for k, v := range atmp {
				ctmp[k] += tmp * v
				tmp2 += btmp[k] * v
			}
			c[i*ldc+j] *= beta
			c[i*ldc+j] += tmp*a[j*lda+j] + alpha*tmp2
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
