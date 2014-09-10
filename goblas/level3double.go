package goblas

import "github.com/gonum/blas"

var _ blas.Float64Level3 = Blasser

func (bl Blas) Dgemm(tA, tB blas.Transpose, m, n, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {

	nota := tA == blas.NoTrans
	notb := tB == blas.NoTrans

	if !nota && tA != blas.ConjTrans && tA != blas.Trans {
		panic(badTranspose)
	}
	if !notb && tB != blas.ConjTrans && tB != blas.Trans {
		panic(badTranspose)
	}
	if m < 0 {
		panic(mLT0)
	}
	if n < 0 {
		panic(nLT0)
	}

	nrowa := k
	ncola := m
	if nota {
		nrowa, ncola = ncola, nrowa
	}

	nrowb := n
	ncolb := k
	if notb {
		nrowb, ncolb = ncolb, nrowb
	}

	if lda < max(1, ncola) {
		panic("lda must be at least the column dimension")
	}
	if ldb < max(1, ncolb) {
		panic("lda must be at least the column dimension")
	}
	if ldc < max(1, n) {
		panic("ldc must be at least the column dimension")
	}

	if m == 0 || n == 0 || ((alpha == 0 || k == 0) && beta == 1) {
		return
	}

	if beta != 1 {
		// Scale c
		if beta == 0 {
			for i := range c {
				c[i] = 0
			}
		} else {
			for i := range c {
				c[i] *= beta
			}
		}
	}
	if alpha == 0 {
		return
	}

	// This code is more complicated than it strictly needs to be, but
	// the loops are constructed so that the elements of the slices are
	// indexed in order. This minimizes cache misses and means that a future
	// go compiler can eliminate many (all) of the bounds checks
	if notb {
		if nota {
			// C += alpha * A*B
			for i := 0; i < m; i++ {
				for l := 0; l < k; l++ {
					tmp := a[i*lda+l] * alpha
					if tmp != 0 {
						for j := 0; j < n; j++ {
							c[i*ldc+j] += tmp * b[l*ldb+j]
						}
					}
				}
			}
			return

		}
		// C += A^T * B
		for i := 0; i < m; i++ {
			for l := 0; l < k; l++ {
				tmp := a[l*lda+i]
				if tmp != 0 {
					tmp *= alpha
					for j := 0; j < n; j++ {
						c[i*ldc+j] += tmp * b[l*ldb+j]
					}
				}
			}
		}
		return
	}
	if nota {
		// C += A * B^T
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				var tmp float64
				for l := 0; l < k; l++ {
					tmp += a[i*lda+l] * b[j*ldb+l]
				}
				c[i*ldc+j] += alpha * tmp
			}
		}
		return
	}
	// C += A^T * B^T
	for i := 0; i < m; i++ {
		for l := 0; l < k; l++ {
			aval := a[l*lda+i]
			if aval != 0 {
				tmp := alpha * aval
				for j := 0; j < n; j++ {
					c[i*ldc+j] += tmp * b[j*ldb+l]
				}
			}
		}
	}
	return
}

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

func (Blas) Dsymm(s blas.Side, ul blas.Uplo, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {
	panic("blas: function not implemented")
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
