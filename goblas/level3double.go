package goblas

import "github.com/gonum/blas"

var _ blas.Float64Level3 = Blasser

const (
	// TODO (btracey): Fix the ld panic messages to be consistent across the package
	badLd string = "goblas: ld must be greater than the number of columns"
)

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

// Dsyrk performs the symmetric rank-k operation
//  C = alpha * A * A^T + beta*C
// where alpha and beta are scalars, C is an nxn symmetric matrix, and A
// is n x k if NonTrans, and k x n if Trans.
func (Blas) Dsyrk(ul blas.Uplo, tA blas.Transpose, n, k int, alpha float64, a []float64, lda int, beta float64, c []float64, ldc int) {
	if ul != blas.Lower && ul != blas.Upper {
		panic(badUplo)
	}
	if tA != blas.Trans && tA != blas.NoTrans && tA != blas.ConjTrans {
		panic(badTranspose)
	}
	if n < 0 {
		panic(nLT0)
	}
	if k < 0 {
		panic(kLT0)
	}
	if ldc < n {
		panic(badLd)
	}
	if tA == blas.Trans {
		if lda < n {
			panic(badLd)
		}
	} else {
		if lda < k {
			panic(badLd)
		}
	}
	if alpha == 0 {
		if ul == blas.Upper {
			for i := 0; i < n; i++ {
				ctmp := c[i*ldc+i : i*ldc+n]
				for j := range ctmp {
					ctmp[j] *= beta
				}
			}
			return
		}
		for i := 0; i < n; i++ {
			ctmp := c[i*ldc : i*ldc+i+1]
			for j := range ctmp {
				ctmp[j] *= beta
			}
		}
		return
	}
	if tA == blas.NoTrans {
		if ul == blas.Upper {
			for i := 0; i < n; i++ {
				ctmp := c[i*ldc+i : i*ldc+n]
				atmp := a[i*lda : i*lda+k]
				for jc, vc := range ctmp {
					j := jc + i
					var tmp float64
					for l, av := range a[j*lda : j*lda+k] {
						tmp += atmp[l] * av
					}
					tmp *= alpha
					tmp += vc * beta
					ctmp[jc] = tmp
				}
			}
			return
		}
		for i := 0; i < n; i++ {
			atmp := a[i*lda : i*lda+k]
			for j, vc := range c[i*ldc : i*ldc+i+1] {
				var tmp float64
				for l, va := range a[j*lda : j*lda+k] {
					tmp += atmp[l] * va
				}
				tmp *= alpha
				tmp += vc * beta
				c[i*ldc+j] = tmp
			}
		}
		return
	}
	// Cases where a is transposed.
	if ul == blas.Upper {
		for i := 0; i < n; i++ {
			ctmp := c[i*ldc+i : i*ldc+n]
			if beta != 1 {
				for j := range ctmp {
					ctmp[j] *= beta
				}
			}
			for l := 0; l < k; l++ {
				tmp := alpha * a[l*lda+i]
				if tmp != 0 {
					for j, v := range a[l*lda+i : l*lda+n] {
						ctmp[j] += tmp * v
					}
				}
			}
		}
		return
	}
	for i := 0; i < n; i++ {
		ctmp := c[i*ldc : i*ldc+i+1]
		if beta != 0 {
			for j := range ctmp {
				ctmp[j] *= beta
			}
		}
		for l := 0; l < k; l++ {
			tmp := alpha * a[l*lda+i]
			if tmp != 0 {
				for j, v := range a[l*lda : l*lda+i+1] {
					ctmp[j] += tmp * v
				}
			}
		}
	}
}

// Dsyr2k performs a symmetric rank 2k operation
//  C = alpha * A * B^T + alpha * B * A^T + beta *C
// where C is an n x n symmetric matrix and A and B are n x k matrices if
// tA == NoTrans and k x n otherwise.
func (Blas) Dsyr2k(ul blas.Uplo, tA blas.Transpose, n, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {
	if ul != blas.Lower && ul != blas.Upper {
		panic(badUplo)
	}
	if tA != blas.Trans && tA != blas.NoTrans && tA != blas.ConjTrans {
		panic(badTranspose)
	}
	if n < 0 {
		panic(nLT0)
	}
	if k < 0 {
		panic(kLT0)
	}
	if ldc < n {
		panic(badLd)
	}
	if tA == blas.Trans {
		if lda < n || ldb < n {
			panic(badLd)
		}
	} else {
		if lda < k || ldb < k {
			panic(badLd)
		}
	}
	if alpha == 0 {
		if ul == blas.Upper {
			for i := 0; i < n; i++ {
				ctmp := c[i*ldc+i : i*ldc+n]
				for j := range ctmp {
					ctmp[j] *= beta
				}
			}
			return
		}
		for i := 0; i < n; i++ {
			ctmp := c[i*ldc : i*ldc+i+1]
			for j := range ctmp {
				ctmp[j] *= beta
			}
		}
		return
	}
	if tA == blas.NoTrans {
		if ul == blas.Upper {
			for i := 0; i < n; i++ {
				atmp := a[i*lda : i*lda+k]
				btmp := b[i*lda : i*lda+k]
				ctmp := c[i*ldc+i : i*ldc+n]
				for jc := range ctmp {
					j := i + jc
					var tmp1, tmp2 float64
					binner := b[j*ldb : j*ldb+k]
					for l, v := range a[j*lda : j*lda+k] {
						tmp1 += v * btmp[l]
						tmp2 += atmp[l] * binner[l]
					}
					ctmp[jc] *= beta
					ctmp[jc] += alpha * (tmp1 + tmp2)
				}
			}
			return
		}
		for i := 0; i < n; i++ {
			atmp := a[i*lda : i*lda+k]
			btmp := b[i*lda : i*lda+k]
			ctmp := c[i*ldc : i*ldc+i+1]
			for j := 0; j <= i; j++ {
				var tmp1, tmp2 float64
				binner := b[j*ldb : j*ldb+k]
				for l, v := range a[j*lda : j*lda+k] {
					tmp1 += v * btmp[l]
					tmp2 += atmp[l] * binner[l]
				}
				ctmp[j] *= beta
				ctmp[j] += alpha * (tmp1 + tmp2)
			}
		}
		return
	}
	if ul == blas.Upper {
		for i := 0; i < n; i++ {
			ctmp := c[i*ldc+i : i*ldc+n]
			if beta != 1 {
				for j := range ctmp {
					ctmp[j] *= beta
				}
			}
			for l := 0; l < k; l++ {
				tmp1 := alpha * b[l*lda+i]
				tmp2 := alpha * a[l*lda+i]
				btmp := b[l*ldb+i : l*ldb+n]
				if tmp1 != 0 || tmp2 != 0 {
					for j, v := range a[l*lda+i : l*lda+n] {
						ctmp[j] += v*tmp1 + btmp[j]*tmp2
					}
				}
			}
		}
		return
	}
	for i := 0; i < n; i++ {
		ctmp := c[i*ldc : i*ldc+i+1]
		if beta != 1 {
			for j := range ctmp {
				ctmp[j] *= beta
			}
		}
		for l := 0; l < k; l++ {
			tmp1 := alpha * b[l*lda+i]
			tmp2 := alpha * a[l*lda+i]
			btmp := b[l*ldb : l*ldb+i+1]
			if tmp1 != 0 || tmp2 != 0 {
				for j, v := range a[l*lda : l*lda+i+1] {
					ctmp[j] += v*tmp1 + btmp[j]*tmp2
				}
			}
		}
	}
}

func (Blas) Dtrmm(s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int) {
	panic("blas: function not implemented")
}
