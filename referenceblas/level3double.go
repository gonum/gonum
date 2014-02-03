package referenceblas

import "github.com/gonum/blas"

var _ blas.Float64Level3 = Blasser

func (bl Blas) Dgemm(o blas.Order, tA, tB blas.Transpose, m, n, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {

	nota := tA == blas.NoTrans
	notb := tB == blas.NoTrans
	isRowMajor := o == blas.RowMajor

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

	nrowa := k
	ncola := m
	if nota {
		nrowa, ncola = ncola, nrowa
	}

	nrowb := n
	ncolb := k
	if nota {
		nrowb, ncolb = ncolb, nrowb
	}

	if isRowMajor {
		if lda < max(1, ncola) {
			panic("lda must be at least the column dimension for row major")
		}
		if ldb < max(1, ncolb) {
			panic("lda must be at least the column dimension for row major")
		}
		if ldc < max(1, n) {
			panic("ldc must be at least the column dimension for row major")
		}
	} else {
		if lda < max(1, nrowa) {
			panic("lda must be at least the column dimension for col major")
		}
		if ldb < max(1, nrowb) {
			panic("lda must be at least the column dimension for col major")
		}
		if ldc < max(1, m) {
			panic("ldc must be at least the column dimension for col major")
		}
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
			// C += A*B
			if isRowMajor {
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
			for j := 0; j < n; j++ {
				for l := 0; l < k; l++ {
					tmp := b[j*ldb+l]
					if tmp != 0 {
						tmp *= alpha
						for i := 0; i < m; i++ {
							c[j*ldc+i] += tmp * a[l*lda+i]
						}
					}
				}
			}
			return
		}
		// C += A^T * B
		if isRowMajor {
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
		for j := 0; j < n; j++ {
			for i := 0; i < m; i++ {
				var tmp float64
				for l := 0; l < k; l++ {
					tmp += a[i*lda+l] * b[j*ldb+l]
				}
				c[j*ldc+i] += alpha * tmp
			}
		}
		return
	}
	if nota {
		// C += A * B^T
		if isRowMajor {
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
		for j := 0; j < n; j++ {
			for l := 0; l < k; l++ {
				if b[l*ldb+j] != 0 {
					tmp := alpha * b[l*ldb+j]
					for i := 0; i < m; i++ {
						c[j*ldc+i] += tmp * a[l*lda+i]
					}
				}
			}
		}
		return
	}
	// C += A^T * B^T
	if isRowMajor {
		for i := 0; i < m; i++ {
			for l := 0; l < k; l++ {
				aval := a[i*lda+l]
				if aval != 0 {
					tmp := alpha * aval
					for j := 0; j < n; j++ {
						c[i*ldc+j] += tmp * b[i*lda+l]
					}
				}
			}
		}
		return
	}
	for j := 0; j < n; j++ {
		for i := 0; i < m; i++ {
			var tmp float64
			for l := 0; l < k; l++ {
				tmp += a[i*lda+l] * b[l*ldb+l]
			}
			c[j*ldc+i] += alpha * tmp
		}
	}
}

func (bl Blas) Dtrsm(o blas.Order, s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int) {
	if o != blas.RowMajor && o != blas.ColMajor {
		panic(badOrder)
	}

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
			bl.Dtrsv(blas.ColMajor, ul, tA, d, m, a, lda, b[jb:], 1)
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
					for k := 0; k < j-1; k++ {
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
					for j := 0; j < k-1; j++ {
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

func (Blas) Dsymm(o blas.Order, s blas.Side, ul blas.Uplo, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {
	panic("blas: function not implemented")
}
func (Blas) Dsyrk(o blas.Order, ul blas.Uplo, t blas.Transpose, n, k int, alpha float64, a []float64, lda int, beta float64, c []float64, ldc int) {
	panic("blas: function not implemented")
}
func (Blas) Dsyr2k(o blas.Order, ul blas.Uplo, t blas.Transpose, n, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {
	panic("blas: function not implemented")
}
func (Blas) Dtrmm(o blas.Order, s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int) {
	panic("blas: function not implemented")
}
