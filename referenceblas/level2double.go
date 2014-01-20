package referenceblas

import (
	"github.com/gonum/blas"
)

var badOrder string = "referenceblas: illegal order"
var mLT0 string = "referenceblas: m < 0"
var nLT0 string = "referenceblas: m < 0"

func verifyLevel2inputs(o blas.Order, tA blas.Transpose, m, n, incX, incY, lda int) {
	if o != blas.RowMajor && o != blas.ColMajor {
		panic("blas: illegal order")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if m < 0 {
		panic(mLT0)
	}
	if n < 0 {
		panic(nLT0)
	}
	if o == blas.ColMajor {
		if lda > m && lda > 1 {
			panic("blas: lda must be less than max(1,m) for row major")
		}
	} else {
		if lda > n && lda > 0 {
			panic("blas: lda must be less than max(1,n) for col major")
		}
	}
	if incX == 0 {
		panic(zeroInc)
	}
	if incY == 0 {
		panic(zeroInc)
	}
}

func getLevel2Indexes(o blas.Order, tA blas.Transpose, m, n, incX, incY int) (lenx, leny, kx, ky int) {
	// Set up the lengths of the vectors and start up points
	// TODO: Figure out how this works with order
	lenx = m
	leny = n
	if tA == blas.NoTrans {
		lenx = n
		leny = m
	}
	if incX > 0 {
		kx = 0
	} else {
		kx = -(lenx - 1) * incX
	}
	if incY > 0 {
		ky = 0
	} else {
		ky = -(leny - 1) * incY
	}
	return lenx, leny, kx, ky
}

// Dgemv computes y = alpha*a*x + beta*y if tA = blas.NoTrans
// or alpha*A^T*x + beta*y if tA = blas.Trans or blas.ConjTrans
func (b Blas) Dgemv(o blas.Order, tA blas.Transpose, m, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int) {
	verifyLevel2inputs(o, tA, m, n, incX, incY, lda)

	// Quick return if possible
	if m == 0 || n == 0 || (alpha == 0 && beta == 1) {
		return
	}

	_, lenY, kx, ky := getLevel2Indexes(o, tA, m, n, incX, incY)

	// First form y := beta * y
	b.Dscal(lenY, beta, y, incY)

	if alpha == 0 {
		return
	}

	// Form y := alpha * A * x + y
	switch {
	case tA == blas.NoTrans && o == blas.RowMajor:
		jx := kx
		for j := 0; j < n; j++ {
			temp := alpha * x[jx]
			iy := ky
			for i := 0; i < m; i++ {
				y[iy] += temp * a[lda*i+j]
				iy += incY
			}
			jx += incX
		}
	case (tA == blas.Trans || tA == blas.ConjTrans) && o == blas.RowMajor:
		jy := ky
		for j := 0; j < n; j++ {
			var temp float64
			ix := kx
			for i := 0; i < m; i++ {
				temp += a[lda*i+j] * x[ix]
				ix += incX
			}
			y[jy] += alpha * temp
			jy += incY
		}
	default:
		// TODO: Add in other switch cases
		panic("Not yet implemented")
	}
}

/*
// Level 2 routines.
        Dgemv(o Order, tA Transpose, m, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
        Dgbmv(o Order, tA Transpose, m, n, kL, kU int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
        Dtrmv(o Order, ul Uplo, tA Transpose, d Diag, n int, a []float64, lda int, x []float64, incX int)
        Dtbmv(o Order, ul Uplo, tA Transpose, d Diag, n, k int, a []float64, lda int, x []float64, incX int)
        Dtpmv(o Order, ul Uplo, tA Transpose, d Diag, n int, ap []float64, x []float64, incX int)
        Dtrsv(o Order, ul Uplo, tA Transpose, d Diag, n int, a []float64, lda int, x []float64, incX int)
        Dtbsv(o Order, ul Uplo, tA Transpose, d Diag, n, k int, a []float64, lda int, x []float64, incX int)
        Dtpsv(o Order, ul Uplo, tA Transpose, d Diag, n int, ap []float64, x []float64, incX int)
        Dsymv(o Order, ul Uplo, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
        Dsbmv(o Order, ul Uplo, n, k int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
        Dspmv(o Order, ul Uplo, n int, alpha float64, ap []float64, x []float64, incX int, beta float64, y []float64, incY int)
        Dger(o Order, m, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int)
        Dsyr(o Order, ul Uplo, n int, alpha float64, x []float64, incX int, a []float64, lda int)
        Dspr(o Order, ul Uplo, n int, alpha float64, x []float64, incX int, ap []float64)
        Dsyr2(o Order, ul Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int)
        Dspr2(o Order, ul Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64)
*/
