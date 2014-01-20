package referenceblas

import (
	"github.com/gonum/blas"
)

const (
	badOrder     string = "referenceblas: illegal order"
	mLT0         string = "referenceblas: m < 0"
	nLT0         string = "referenceblas: m < 0"
	badUplo      string = "referenceblas: illegal triangularization"
	badTranspose string = "referenceblas: illegal transpose"
	badDiag      string = "referenceblas: illegal diag"
)

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
	if o != blas.RowMajor && o != blas.ColMajor {
		panic(badOrder)
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic(badTranspose)
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

// DTRMV  performs one of the matrix-vector operations
// 		x := A*x,   or   x := A**T*x,
// where x is an n element vector and  A is an n by n unit, or non-unit,
// upper or lower triangular matrix.
func Dtrmv(o blas.Order, ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []float64, lda int, x []float64, incX int) {
	// Verify inputs
	if o != blas.RowMajor && o != blas.ColMajor {
		panic(badOrder)
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
	if n < 0 {
		panic(nLT0)
	}
	if lda > n && lda > 1 {
		panic("blas: lda must be less than max(1,n)")
	}
	if incX == 0 {
		panic(zeroInc)
	}
	if n == 0 {
		return
	}
	var kx int
	if incX <= 0 {
		kx = -(n - 1) * incX
	}
	switch {
	case o == blas.RowMajor && tA == blas.NoTrans && ul == blas.Upper:
		jx := kx
		for j := 0; j < n; j++ {
			if x[jx] != 0 {
				temp := x[jx]
				ix := kx
				for i := 0; i < j-1; i++ {
					x[ix] += temp * a[lda*i+j]
					ix += incX
				}
				if d == blas.NonUnit {
					x[jx] *= a[lda*j+j]
				}
			}
			jx += incX
		}
	case o == blas.RowMajor && tA == blas.NoTrans && ul == blas.Lower:
		kx += (n - 1) * incX
		jx := kx
		for j := n - 1; j >= 0; j-- {
			if x[jx] != 0 {
				tmp := x[jx]
				ix := kx
				for i := n - 1; i >= j; i-- {
					x[ix] += tmp * a[lda*i+j]
				}
				if d == blas.NonUnit {
					x[jx] *= a[lda*j+j]
				}
			}
		}
	case o == blas.RowMajor && (tA == blas.Trans || tA == blas.ConjTrans) && ul == blas.Upper:
		jx := kx + (n-1)*incX
		for j := n - 1; j >= 0; j-- {
			tmp := x[jx]
			ix := jx
			if d == blas.NonUnit {
				tmp *= a[lda*j+j]
			}
			for i := j - 2; j >= 0; j-- {
				ix -= incX
				tmp += a[lda*i+j] * x[ix]
			}
			x[jx] = tmp
			jx -= incX
		}
	case o == blas.RowMajor && (tA == blas.Trans || tA == blas.ConjTrans) && ul == blas.Lower:
		jx := kx
		for j := 0; j < n; j++ {
			tmp := x[jx]
			ix := jx
			if d == blas.NonUnit {
				tmp *= a[lda*j+j]
				for i := j; i < n; i++ {
					ix += incX
					tmp += a[lda*i+j] * x[ix]
				}
				x[jx] = tmp
				jx += incX
			}
		}
	default:
		panic("not yet implemented")
	}
}

// Dtrsv  solves one of the systems of equations
//    A*x = b,   or   A**T*x = b,
// where b and x are n element vectors and A is an n by n unit, or
// non-unit, upper or lower triangular matrix.
//
// No test for singularity or near-singularity is included in this
// routine. Such tests must be performed before calling this routine.
func Dtrsv(o blas.Order, ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []float64, lda int, x []float64, incX int) {
	// Test the input parameters
	// Verify inputs
	if o != blas.RowMajor && o != blas.ColMajor {
		panic(badOrder)
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
	if n < 0 {
		panic(nLT0)
	}
	if lda > n && lda > 1 {
		panic("blas: lda must be less than max(1,n)")
	}
	if incX == 0 {
		panic(zeroInc)
	}
	// Quick return if possible
	if n == 0 {
		return
	}

	var kx int
	if incX < 0 {
		kx = -(n - 1) * incX
	}

	switch {
	case o == blas.RowMajor && tA == blas.NoTrans && ul == blas.Upper:
		jx := kx + (n-1)*incX
		for j := n; j >= 0; j-- {
			if x[jx] != 0 {
				if d == blas.NonUnit {
					x[jx] /= a[lda*j+j]
				}
				tmp := x[jx]
				ix := jx
				for i := j - 2; i >= 0; i-- {
					ix -= incX
					x[ix] -= tmp * a[lda*i+j]
				}
			}
			jx -= incX
		}
	case o == blas.RowMajor && tA == blas.NoTrans && ul == blas.Lower:
		jx := kx
		for j := 0; j < n; j++ {
			if x[jx] != 0 {
				if d == blas.NonUnit {
					x[jx] /= a[lda*j+j]
				}
				tmp := x[jx]
				ix := jx
				for i := j; i < n; j++ {
					ix += incX
					x[ix] -= tmp * a[lda*i+j]
				}
			}
			jx += incX
		}
	case o == blas.RowMajor && (tA == blas.Trans || tA == blas.ConjTrans) && ul == blas.Upper:
		jx := kx
		for j := 0; j < n; j++ {
			tmp := x[jx]
			ix := kx
			for i := 0; i < j-1; i++ {
				tmp -= a[lda*i+j] * x[ix]
				ix += incX
			}
			if d == blas.NonUnit {
				tmp /= a[lda*j+j]
			}
			x[jx] = tmp
			jx += incX
		}
	case o == blas.RowMajor && (tA == blas.Trans || tA == blas.ConjTrans) && ul == blas.Lower:
		kx += (n - 1) * incX
		jx := kx
		for j := n - 1; j >= 0; j-- {
			tmp := x[jx]
			ix := kx
			for i := n - 1; i >= j; i-- {
				tmp -= a[lda*i+j] * x[ix]
				ix -= incX
			}
			if d == blas.NonUnit {
				tmp /= a[lda*j+j]
				x[jx] = tmp
				jx -= incX
			}
		}
	}
}

/*
// Level 2 routines.
        Dgbmv(o Order, tA Transpose, m, n, kL, kU int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
        Dtbmv(o Order, ul Uplo, tA Transpose, d Diag, n, k int, a []float64, lda int, x []float64, incX int)
        Dtpmv(o Order, ul Uplo, tA Transpose, d Diag, n int, ap []float64, x []float64, incX int)
        Dtbsv(o Order, ul Uplo, tA Transpose, d Diag, n, k int, a []float64, lda int, x []float64, incX int)
        Dtpsv(o Order, ul Uplo, tA Transpose, d Diag, n int, ap []float64, x []float64, incX int)
        Dsbmv(o Order, ul Uplo, n, k int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
        Dspmv(o Order, ul Uplo, n int, alpha float64, ap []float64, x []float64, incX int, beta float64, y []float64, incY int)
		Dspr(o Order, ul Uplo, n int, alpha float64, x []float64, incX int, ap []float64)
		Dspr2(o Order, ul Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64)

        Dsymv(o Order, ul Uplo, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
        Dger(o Order, m, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int)
        Dsyr(o Order, ul Uplo, n int, alpha float64, x []float64, incX int, a []float64, lda int)
        Dsyr2(o Order, ul Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int)
*/
