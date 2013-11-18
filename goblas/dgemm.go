package goblas

import "github.com/gonum/blas"
import level1 "github.com/ziutek/blas"

func (Blas) Dgemm(o blas.Order, tA, tB blas.Transpose, m int, n int, k int,
	alpha float64, a []float64, lda int, b []float64, ldb int,
	beta float64, c []float64, ldc int) {

	var inner, outer, veclen int
	if o == blas.ColMajor {
		outer = n
		veclen = m
	} else {
		veclen = n
		outer = m
		a, b = b, a
		ldb, lda = lda, ldb
		tA, tB = tB, tA
	}
	inner = k

	for j := 0; j < outer; j++ {
		cj := c[j*ldc:]
		if beta != 1 {
			level1.Dscal(veclen, beta, cj, 1)
		}

		for l := 0; l < inner; l++ {
			al := a[l*lda:]
			if tA == blas.Trans {
				al = a[l:]
			}
			blj := b[j*ldb+l]
			if tB == blas.Trans {
				blj = b[l*ldb+j]
			}
			if tA == blas.NoTrans {
				level1.Daxpy(veclen, blj*alpha, al, 1, cj, 1)
			} else {
				level1.Daxpy(veclen, blj*alpha, al, lda, cj, 1)
			}
		}
	}
}
