package goblas

import "github.com/gonum/blas"

// See http://www.netlib.org/lapack/explore-html/d4/de1/_l_i_c_e_n_s_e_source.html
// for more license information

var _ blas.Float64Level2 = Blasser

// TODO: Need to think about loops when doing row-major. Change after tests?

const (
	mLT0         string = "referenceblas: m < 0"
	nLT0         string = "referenceblas: n < 0"
	kLT0         string = "referenceblas: k < 0"
	badUplo      string = "referenceblas: illegal triangularization"
	badTranspose string = "referenceblas: illegal transpose"
	badDiag      string = "referenceblas: illegal diag"
	badSide      string = "referenceblas: illegal side"
	badLdaRow    string = "lda must be greater than max(1,n) for row major"
	badLdaCol    string = "lda must be greater than max(1,m) for col major"
	badLda       string = "lda must be greater than max(1,n)"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// Dgemv computes y = alpha*a*x + beta*y if tA = blas.NoTrans
// or alpha*A^T*x + beta*y if tA = blas.Trans or blas.ConjTrans
func (b Blas) Dgemv(tA blas.Transpose, m, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int) {
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic(badTranspose)
	}
	if m < 0 {
		panic(mLT0)
	}
	if n < 0 {
		panic(nLT0)
	}
	if lda < max(1, n) {
		panic(badLdaRow)
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

	// Set up indexes
	lenX := m
	lenY := n
	if tA == blas.NoTrans {
		lenX = n
		lenY = m
	}
	var kx, ky int
	if incX > 0 {
		kx = 0
	} else {
		kx = -(lenX - 1) * incX
	}
	if incY > 0 {
		ky = 0
	} else {
		ky = -(lenY - 1) * incY
	}

	// First form y := beta * y
	if incY > 0 {
		b.Dscal(lenY, beta, y, incY)
	} else {
		b.Dscal(lenY, beta, y, -incY)
	}

	if alpha == 0 {
		return
	}

	// Form y := alpha * A * x + y
	switch {

	default:
		panic("shouldn't be here")

	case tA == blas.NoTrans:
		iy := ky
		for i := 0; i < m; i++ {
			jx := kx
			var temp float64
			for j := 0; j < n; j++ {
				temp += a[lda*i+j] * x[jx]
				jx += incX
			}
			y[iy] += alpha * temp
			iy += incY
		}
	case tA == blas.Trans || tA == blas.ConjTrans:
		ix := kx
		for i := 0; i < m; i++ {
			jy := ky
			tmp := alpha * x[ix]
			for j := 0; j < n; j++ {
				y[jy] += a[lda*i+j] * tmp
				jy += incY
			}
			ix += incX
		}
	}
}

// Dger   performs the rank 1 operation
//    A := alpha*x*y**T + A,
// where alpha is a scalar, x is an m element vector, y is an n element
// vector and A is an m by n matrix.
func (Blas) Dger(m, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int) {
	// Check inputs
	if m < 0 {
		panic("m < 0")
	}
	if n < 0 {
		panic(negativeN)
	}
	if incX == 0 {
		panic(zeroInc)
	}
	if incY == 0 {
		panic(zeroInc)
	}
	if lda < max(1, n) {
		panic(badLdaRow)
	}

	// Quick return if possible
	if m == 0 || n == 0 || alpha == 0 {
		return
	}

	var ky, kx int
	if incY > 0 {
		ky = 0
	} else {
		ky = -(n - 1) * incY
	}

	if incY > 0 {
		kx = 0
	} else {
		kx = -(m - 1) * incX
	}

	ix := kx
	for i := 0; i < m; i++ {
		if x[ix] == 0 {
			ix += incX
			continue
		}
		tmp := alpha * x[ix]
		jy := ky
		for j := 0; j < n; j++ {
			a[i*lda+j] += y[jy] * tmp
			jy += incY
		}
		ix += incX
	}

}

func (b Blas) Dgbmv(tA blas.Transpose, m, n, kL, kU int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int) {
	// Transform for row major
	m, n = n, m
	kU, kL = kL, kU
	if tA == blas.NoTrans {
		tA = blas.Trans
	} else {
		tA = blas.NoTrans
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
	if lda < max(1, n) {
		panic(badLdaRow)
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

	// Set up indexes
	lenX := m
	lenY := n
	if tA == blas.NoTrans {
		lenX = n
		lenY = m
	}
	var kx, ky int
	if incX > 0 {
		kx = 0
	} else {
		kx = -(lenX - 1) * incX
	}
	if incY > 0 {
		ky = 0
	} else {
		ky = -(lenY - 1) * incY
	}

	// First form y := beta * y
	if incY > 0 {
		b.Dscal(lenY, beta, y, incY)
	} else {
		b.Dscal(lenY, beta, y, -incY)
	}

	if alpha == 0 {
		return
	}

	if tA == blas.NoTrans {
		jx := kx
		if incY == 1 {
			for j := 0; j < n; j++ {
				if x[jx] != 0 {
					temp := alpha * x[jx]
					k := kU - j
					for i := max(0, j-kU); i < min(m, j+kL+1); i++ {
						y[i] += temp * a[k+i+j*lda]
					}
				}
				jx += incX
			}
		} else {
			for j := 0; j < n; j++ {
				if x[jx] != 0 {
					temp := alpha * x[jx]
					iy := ky
					k := kU - j
					for i := max(0, j-kU); i < min(m, j+kL+1); i++ {
						y[iy] += temp * a[k+i+j*lda]
						iy += incY
					}
				}
				jx += incX
				if j >= kU {
					ky += incY
				}
			}
		}
	} else {
		jy := ky
		if incX == 1 {
			for j := 0; j < n; j++ {
				temp := 0.0
				k := kU - j
				for i := max(0, j-kU); i < min(m, j+kL+1); i++ {
					temp += a[k+i+j*lda] * x[i]
				}
				y[jy] += alpha * temp
				jy += incY
			}
		} else {
			for j := 0; j < n; j++ {
				temp := 0.0
				ix := kx
				k := kU - j
				for i := max(0, j-kU); i < min(m, j+kL+1); i++ {
					temp += a[k+i+j*lda] * x[ix]
					ix += incX
				}
				y[jy] += alpha * temp
				jy += incY
				if j > kU {
					kx += incX
				}
			}
		}
	}
}

// DTRMV  performs one of the matrix-vector operations
// 		x := A*x,   or   x := A**T*x,
// where x is an n element vector and  A is an n by n unit, or non-unit,
// upper or lower triangular matrix.
func (Blas) Dtrmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []float64, lda int, x []float64, incX int) {
	// Verify inputs
	if tA == blas.NoTrans {
		tA = blas.Trans
	} else {
		tA = blas.NoTrans
	}
	if ul == blas.Upper {
		ul = blas.Lower
	} else {
		ul = blas.Upper
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
		panic(badLda)
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
	default:
		panic("unreachable")
	case tA == blas.NoTrans && ul == blas.Upper:
		jx := kx
		for j := 0; j < n; j++ {
			ja := j * lda
			if x[jx] != 0 {
				temp := x[jx]
				ix := kx
				for i := 0; i < j; i++ {
					x[ix] += temp * a[i+ja]
					ix += incX
				}
				if d == blas.NonUnit {
					x[jx] *= a[j+ja]
				}
			}
			jx += incX
		}
	case tA == blas.NoTrans && ul == blas.Lower:
		kx += (n - 1) * incX
		jx := kx
		for j := n - 1; j >= 0; j-- {
			ja := j * lda
			if x[jx] != 0 {
				tmp := x[jx]
				ix := kx
				for i := n - 1; i > j; i-- {
					x[ix] += tmp * a[i+ja]
					ix -= incX
				}
				if d == blas.NonUnit {
					x[jx] *= a[j+ja]
				}
			}
			jx -= incX
		}
	case (tA == blas.Trans || tA == blas.ConjTrans) && ul == blas.Upper:
		jx := kx + (n-1)*incX
		for j := n - 1; j >= 0; j-- {
			ja := j * lda
			tmp := x[jx]
			ix := jx
			if d == blas.NonUnit {
				tmp *= a[j+ja]
			}
			for i := j - 1; i >= 0; i-- {
				ix -= incX
				tmp += a[i+ja] * x[ix]
			}
			x[jx] = tmp
			jx -= incX
		}
	case (tA == blas.Trans || tA == blas.ConjTrans) && ul == blas.Lower:
		jx := kx
		for j := 0; j < n; j++ {
			tmp := x[jx]
			ix := jx
			ja := j * lda
			if d == blas.NonUnit {
				tmp *= a[j+ja]
			}
			for i := j + 1; i < n; i++ {
				ix += incX
				tmp += a[i+ja] * x[ix]
			}
			x[jx] = tmp
			jx += incX
		}
	}
}

// Dtrsv  solves one of the systems of equations
//    A*x = b,   or   A**T*x = b,
// where b and x are n element vectors and A is an n by n unit, or
// non-unit, upper or lower triangular matrix.
//
// No test for singularity or near-singularity is included in this
// routine. Such tests must be performed before calling this routine.
func (Blas) Dtrsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []float64, lda int, x []float64, incX int) {
	// Test the input parameters
	// Verify inputs
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
	default:
		panic("goblas: unreachable")
	case tA == blas.NoTrans && ul == blas.Upper:
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
	case tA == blas.NoTrans && ul == blas.Lower:
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
	case (tA == blas.Trans || tA == blas.ConjTrans) && ul == blas.Upper:
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
	case (tA == blas.Trans || tA == blas.ConjTrans) && ul == blas.Lower:
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

// Dsymv  performs the matrix-vector  operation
//    y := alpha*A*x + beta*y,
// where alpha and beta are scalars, x and y are n element vectors and
// A is an n by n symmetric matrix.
func (b Blas) Dsymv(ul blas.Uplo, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int) {
	// Check inputs
	if ul != blas.Lower && ul != blas.Upper {
		panic(badUplo)
	}
	if n < 0 {
		panic(negativeN)
	}
	if lda > 1 && lda > n {
		panic(badLda)
	}
	if incX == 0 {
		panic(zeroInc)
	}
	if incY == 0 {
		panic(zeroInc)
	}
	// Quick return if possible
	if n == 0 || (alpha == 0 && beta == 1) {
		return
	}

	// Set up start points
	var kx, ky int
	if incX > 0 {
		kx = 1
	} else {
		kx = -(n - 1) * incX
	}
	if incY > 0 {
		ky = 1
	} else {
		ky = -(n - 1) * incY
	}

	// Form y = beta * y
	if beta != 1 {
		b.Dscal(n, beta, y, incY)
	}

	if alpha == 0 {
		return
	}

	// TODO: Need to think about changing the major and minor
	// looping when row major (help with cache misses)

	// Form y = Ax + y
	switch {
	default:
		panic("goblas: unreachable")
	case ul == blas.Upper:
		jx := kx
		jy := ky
		for j := 0; j < n; j++ {
			tmp1 := alpha * x[jx]
			var tmp2 float64
			ix := kx
			iy := ky
			for i := 0; i < j-2; i++ {
				y[iy] += tmp1 * a[i*lda+j]
				tmp2 += a[i*lda+j] * x[ix]
				ix += incX
				iy += incY
			}
			y[jy] += tmp1*a[j*lda+j] + alpha*tmp2
			jx += incX
			jy += incY
		}
	case ul == blas.Lower:
		jx := kx
		jy := ky
		for j := 0; j < n; j++ {
			tmp1 := alpha * x[jx]
			var tmp2 float64
			y[jy] += tmp1 * a[j*lda+j]
			ix := jx
			iy := jy
			for i := j; i < n; i++ {
				ix += incX
				iy += incY
				y[iy] += tmp1 * a[i*lda+j]
				tmp2 += a[i*lda+j] * x[ix]
			}
			y[jy] += alpha * tmp2
			jx += incX
			jy += incY
		}
	}
}

// Dtbmv  performs one of the matrix-vector operations
// 		x := A*x,   or   x := A**T*x,
// where x is an n element vector and  A is an n by n unit, or non-unit,
// upper or lower triangular band matrix.
func (Blas) Dtbmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n, k int, a []float64, lda int, x []float64, incX int) {
	// Verify inputs
	// Transform for row major
	if tA == blas.NoTrans {
		tA = blas.Trans
	} else {
		tA = blas.NoTrans
	}
	if ul == blas.Upper {
		ul = blas.Lower
	} else {
		ul = blas.Upper
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
	if k < 0 {
		panic(kLT0)
	}
	if lda < k+1 {
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
	} else if incX != 1 {
		kx = 0
	}

	if tA == blas.NoTrans {
		if ul == blas.Upper {
			if incX == 1 {
				for j := 0; j < n; j++ {
					if x[j] != 0 {
						temp := x[j]
						l := k - j
						for i := max(0, j-k); i < j; i++ {
							x[i] += temp * a[l+i+j*lda]
						}
						if d == blas.NonUnit {
							x[j] *= a[k+j*lda]
						}
					}
				}
			} else {
				jx := kx
				for j := 0; j < n; j++ {
					if x[jx] != 0 {
						temp := x[jx]
						ix := kx
						l := k - j
						for i := max(0, j-k); i < j; i++ {
							x[ix] += temp * a[l+i+j*lda]
							ix += incX
						}
						if d == blas.NonUnit {
							x[jx] *= a[k+j*lda]
						}
					}
					jx += incX
					if j >= k {
						kx += incX
					}
				}
			}
		} else {

			if incX == 1 {
				for j := n - 1; j >= 0; j-- {
					if x[j] != 0 {
						temp := x[j]
						l := -j
						for i := min(n-1, j+k); i >= j+1; i-- {
							x[i] += temp * a[l+i+j*lda]
						}
						if d == blas.NonUnit {
							x[j] *= a[0+j*lda]
						}
					}
				}
			} else {
				kx += (n - 1) * incX
				jx := kx
				for j := n - 1; j >= 0; j-- {
					if x[jx] != 0 {
						temp := x[jx]
						ix := kx
						l := -j
						for i := min(n-1, j+k); i >= j+1; i-- {
							x[ix] += temp * a[l+i+j*lda]
							ix -= incX
						}
						if d == blas.NonUnit {
							x[jx] *= a[0+j*lda]
						}
					}
					jx -= incX
					if n-j > k {
						kx -= incX
					}
				}
			}
		}
	} else {

		if ul == blas.Upper {
			if incX == 1 {
				for j := n - 1; j >= 0; j-- {
					temp := x[j]
					l := k - j
					if d == blas.NonUnit {
						temp *= a[k+j*lda]
					}
					for i := j - 1; i >= max(0, j-k); i-- {
						temp += a[l+i+j*lda] * x[i]
					}
					x[j] = temp
				}
			} else {
				kx += (n - 1) * incX
				jx := kx
				for j := n - 1; j >= 0; j-- {
					temp := x[jx]
					kx -= incX
					ix := kx
					l := k - j
					if d == blas.NonUnit {
						temp *= a[k+j*lda]
					}
					for i := j - 1; i >= max(0, j-k); i-- {
						temp += a[l+i+j*lda] * x[ix]
						ix -= incX
					}
					x[jx] = temp
					jx -= incX
				}
			}
		} else {

			if incX == 1 {
				for j := 0; j < n; j++ {
					temp := x[j]
					l := -j
					if d == blas.NonUnit {
						temp *= a[0+j*lda]
					}
					for i := j + 1; i < min(n, j+k+1); i++ {
						temp += a[l+i+j*lda] * x[i]
					}
					x[j] = temp
				}
			} else {
				jx := kx
				for j := 0; j < n; j++ {
					temp := x[jx]
					kx += incX
					ix := kx
					l := -j
					if d == blas.NonUnit {
						temp *= a[0+j*lda]
					}
					for i := j + 1; i < min(n, j+k+1); i++ {
						temp += a[l+i+j*lda] * x[ix]
						ix += incX
					}
					x[jx] = temp
					jx += incX
				}
			}
		}
	}
}

func (bl Blas) Dtpmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, ap []float64, x []float64, incX int) {
	// Verify inputs
	// Transform for row major
	if tA == blas.NoTrans {
		tA = blas.Trans
	} else {
		tA = blas.NoTrans
	}
	if ul == blas.Upper {
		ul = blas.Lower
	} else {
		ul = blas.Upper
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
	if len(ap) < (n*(n+1))/2 {
		panic("blas: not enough data in ap")
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
	} else if incX != 1 {
		kx = 0
	}

	if tA == blas.NoTrans {

		//x := A*x
		if ul == blas.Upper {
			kk := 0
			jx := kx
			for j := 0; j < n; j++ {
				if x[jx] != 0 {
					if j > 0 {
						offset := max(0, -(n-j)*incX)
						bl.Daxpy(j, x[jx], ap[kk:], 1, x[offset:], incX)
					}
					if d == blas.NonUnit {
						x[jx] *= ap[kk+j]
					}
				}
				kk += j + 1
				jx += incX
			}
		} else {
			kk := (n*(n+1))/2 - 1
			jx := kx + (n-1)*incX
			for j := n - 1; j >= 0; j-- {
				if x[jx] != 0 {
					if j+1 < n {
						offset := max((j+1)*incX, 0)
						bl.Daxpy(n-j-1, x[jx], ap[kk-n+j+2:], 1, x[offset:], incX)
					}
					if d == blas.NonUnit {
						x[jx] *= ap[kk-n+j+1]
					}
				}
				jx -= incX
				kk -= n - j
			}
		}

	} else {

		// x := A**T*x
		if ul == blas.Upper {
			kk := (n*(n+1))/2 - 1
			jx := kx + (n-1)*incX
			for j := n - 1; j >= 0; j-- {
				temp := x[jx]
				if d == blas.NonUnit {
					temp *= ap[kk]
				}
				if j > 0 {
					offset := max(0, -(n-j)*incX)
					temp += bl.Ddot(j, ap[kk-j:], 1, x[offset:], incX)
				}
				x[jx] = temp
				jx -= incX
				kk -= j + 1
			}
		} else {
			kk := 0
			jx := kx
			for j := 0; j < n; j++ {
				temp := x[jx]
				if d == blas.NonUnit {
					temp *= ap[kk]
				}
				if j+1 < n {
					offset := max((j+1)*incX, 0)
					temp += bl.Ddot(n-j-1, ap[kk+1:], 1, x[offset:], incX)
				}
				x[jx] = temp
				jx += incX
				kk += n - j
			}
		}
	}
}

//TODO: Not yet implemented Level 2 routines.
func (Blas) Dtbsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n, k int, a []float64, lda int, x []float64, incX int) {
	panic("referenceblas: function not implemented")
}
func (Blas) Dtpsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, ap []float64, x []float64, incX int) {
	panic("referenceblas: function not implemented")
}
func (Blas) Dsbmv(ul blas.Uplo, n, k int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int) {
	panic("referenceblas: function not implemented")
}
func (Blas) Dspmv(ul blas.Uplo, n int, alpha float64, ap []float64, x []float64, incX int, beta float64, y []float64, incY int) {
	panic("referenceblas: function not implemented")
}
func (Blas) Dspr(ul blas.Uplo, n int, alpha float64, x []float64, incX int, ap []float64) {
	panic("referenceblas: function not implemented")
}
func (Blas) Dspr2(ul blas.Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64) {
	panic("referenceblas: function not implemented")
}
func (Blas) Dsyr(ul blas.Uplo, n int, alpha float64, x []float64, incX int, a []float64, lda int) {
	panic("referenceblas: function not implemented")
}
func (Blas) Dsyr2(ul blas.Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int) {
	panic("referenceblas: function not implemented")
}
