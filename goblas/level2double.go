package goblas

import "github.com/gonum/blas"

// See http://www.netlib.org/lapack/explore-html/d4/de1/_l_i_c_e_n_s_e_source.html
// for more license information

var _ blas.Float64Level2 = Blasser

/*
	Accessing guidelines:
	Dense matrices are laid out in row-major order. [a11 a12 ... a1n a21 a22 ... amn]
	dense(i, j) = dense[i*ld + j]

	Banded matrices are laid out in a compact format as described in
	http://www.crest.iu.edu/research/mtl/reference/html/banded.html
	(the row-major version)
	In short, all of the rows are scrunched removing the zeros and aligning
	the diagonals
	So, for the matrix
	[
	  1  2   3  0  0  0
	  4  5   6  7  0  0
	  0  8   9 10 11  0
	  0  0  12 13 14 15
	  0  0   0 16 17 18
	  0  0   0  0 19 20
	]

	The layout is
	[
	   *  1  2  3
	   4  5  6  7
	   8  9 10 11
	  12 13 14 15
	  16 17 18  *
	  19 20  *  *
	]
	where entries marked * are never accessed
*/

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

	badLdaTriBand string = "goblas: lda must be greater than k + 1 for general banded"
	badLdaGenBand string = "goblas: lda must be greater than 1 + kL + kU for general banded"
	kLLT0         string = "goblas: kL < 0"
	kULT0         string = "goblas: kU < 0"
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
	if tA == blas.NoTrans {
		if incX == 1 {
			for i := 0; i < m; i++ {
				var tmp float64
				atmp := a[lda*i : lda*i+n]
				for j, v := range atmp {
					tmp += v * x[j]
				}
				y[i] += alpha * tmp
			}
			return
		}
		iy := ky
		for i := 0; i < m; i++ {
			jx := kx
			var tmp float64
			atmp := a[lda*i : lda*i+n]
			for _, v := range atmp {
				tmp += v * x[jx]
				jx += incX
			}
			y[iy] += alpha * tmp
			iy += incY
		}
		return
	}
	// Cases where a is not transposed.
	if incX == 1 {
		for i := 0; i < m; i++ {
			tmp := alpha * x[i]
			if tmp != 0 {
				atmp := a[lda*i : lda*i+n]
				for j, v := range atmp {
					y[j] += v * tmp
				}
			}
		}
		return
	}
	ix := kx
	for i := 0; i < m; i++ {
		tmp := alpha * x[ix]
		if tmp != 0 {
			jy := ky
			atmp := a[lda*i : lda*i+n]
			for _, v := range atmp {
				y[jy] += v * tmp
				jy += incY
			}
		}
		ix += incX
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

	if incX == 1 && incY == 1 {
		x = x[:m]
		y = y[:n]
		for i, xv := range x {
			tmp := alpha * xv
			if tmp != 0 {
				atmp := a[i*lda:]
				for j, yv := range y {
					atmp[j] += yv * tmp
				}
			}
		}
		return
	}

	ix := kx
	for i := 0; i < m; i++ {
		if x[ix] != 0 {
			tmp := alpha * x[ix]
			jy := ky
			atmp := a[i*lda:]
			for j := 0; j < n; j++ {
				atmp[j] += y[jy] * tmp
				jy += incY
			}
		}
		ix += incX
	}
}

// Dgbmv performs y = alpha*A*x + beta*y where a is an mxn band matrix with
// kl subdiagonals and ku super-diagonals. m and n refer to the size of the full
// dense matrix, not the actual number of rows and columns in the contracted banded
// matrix.
func (b Blas) Dgbmv(tA blas.Transpose, m, n, kL, kU int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int) {
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic(badTranspose)
	}
	if m < 0 {
		panic(mLT0)
	}
	if n < 0 {
		panic(nLT0)
	}
	if kL < 0 {
		panic(kLLT0)
	}
	if kL < 0 {
		panic(kULT0)
	}
	if lda < kL+kU+1 {
		panic(badLdaGenBand)
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

	// i and j are indices of the compacted banded matrix.
	// off is the offset into the dense matrix (off + j = densej)
	ld := min(m, n)
	nCol := kU + 1 + kL
	if tA == blas.NoTrans {
		iy := ky
		if incX == 1 {
			for i := 0; i < m; i++ {
				l := max(0, kL-i)
				u := min(nCol, ld+kL-i)
				off := max(0, i-kL)
				atmp := a[i*lda+l : i*lda+u]
				xtmp := x[off : off+u-l]
				var sum float64
				for j, v := range atmp {
					sum += xtmp[j] * v
				}
				y[iy] += sum * alpha
				iy += incY
			}
			return
		}
		for i := 0; i < m; i++ {
			l := max(0, kL-i)
			u := min(nCol, ld+kL-i)
			off := max(0, i-kL)
			atmp := a[i*lda+l : i*lda+u]
			jx := kx
			var sum float64
			for _, v := range atmp {
				sum += x[off*incX+jx] * v
				jx += incX
			}
			y[iy] += sum * alpha
			iy += incY
		}
		return
	}
	if incX == 1 {
		for i := 0; i < m; i++ {
			l := max(0, kL-i)
			u := min(nCol, ld+kL-i)
			off := max(0, i-kL)
			atmp := a[i*lda+l : i*lda+u]
			tmp := alpha * x[i]
			jy := ky
			for _, v := range atmp {
				y[jy+off*incY] += tmp * v
				jy += incY
			}
		}
		return
	}
	ix := kx
	for i := 0; i < m; i++ {
		l := max(0, kL-i)
		u := min(nCol, ld+kL-i)
		off := max(0, i-kL)
		atmp := a[i*lda+l : i*lda+u]
		tmp := alpha * x[ix]
		jy := ky
		for _, v := range atmp {
			y[jy+off*incY] += tmp * v
			jy += incY
		}
		ix += incX
	}
}

// Dtrmv performs one of the matrix-vector operations
// 		x := A*x,   or   x := A**T*x,
// where x is an n element vector and  A is an n by n unit, or non-unit,
// upper or lower triangular matrix.
func (Blas) Dtrmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []float64, lda int, x []float64, incX int) {
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
	if lda < n {
		panic(badLda)
	}
	if incX == 0 {
		panic(zeroInc)
	}
	if n == 0 {
		return
	}
	nonUnit := d != blas.Unit
	if n == 1 {
		x[0] *= a[0]
		return
	}
	var kx int
	if incX <= 0 {
		kx = -(n - 1) * incX
	}
	if tA == blas.NoTrans {
		if ul == blas.Upper {
			if incX == 1 {
				for i := 0; i < n; i++ {
					var tmp float64
					if nonUnit {
						tmp = a[i*lda+i] * x[i]
					} else {
						tmp = x[i]
					}
					xtmp := x[i+1:]
					for j, v := range a[i*lda+i+1 : i*lda+n] {
						tmp += v * xtmp[j]
					}
					x[i] = tmp
				}
				return
			}
			ix := kx
			for i := 0; i < n; i++ {
				var tmp float64
				if nonUnit {
					tmp = a[i*lda+i] * x[ix]
				} else {
					tmp = x[ix]
				}
				jx := ix + incX
				for _, v := range a[i*lda+i+1 : i*lda+n] {
					tmp += v * x[jx]
					jx += incX
				}
				x[ix] = tmp
				ix += incX
			}
			return
		}
		if incX == 1 {
			for i := n - 1; i >= 0; i-- {
				var tmp float64
				if nonUnit {
					tmp += a[i*lda+i] * x[i]
				} else {
					tmp = x[i]
				}
				for j, v := range a[i*lda : i*lda+i] {
					tmp += v * x[j]
				}
				x[i] = tmp
			}
			return
		}
		ix := kx + (n-1)*incX
		for i := n - 1; i >= 0; i-- {
			var tmp float64
			if nonUnit {
				tmp += a[i*lda+i] * x[ix]
			} else {
				tmp = x[ix]
			}
			jx := kx
			for _, v := range a[i*lda : i*lda+i] {
				tmp += v * x[jx]
				jx += incX
			}
			x[ix] = tmp
			ix -= incX
		}
		return
	}
	// Cases where a is transposed.
	if ul == blas.Upper {
		if incX == 1 {
			for i := n - 1; i >= 0; i-- {
				xi := x[i]
				atmp := a[i*lda+i+1 : i*lda+n]
				xtmp := x[i+1 : n]
				for j, v := range atmp {
					xtmp[j] += xi * v
				}
				if nonUnit {
					x[i] *= a[i*lda+i]
				}
			}
			return
		}
		ix := kx + (n-1)*incX
		for i := n - 1; i >= 0; i-- {
			xi := x[ix]
			jx := kx + (i+1)*incX
			atmp := a[i*lda+i+1 : i*lda+n]
			for _, v := range atmp {
				x[jx] += xi * v
				jx += incX
			}
			if nonUnit {
				x[ix] *= a[i*lda+i]
			}
			ix -= incX
		}
		return
	}
	if incX == 1 {
		for i := 0; i < n; i++ {
			xi := x[i]
			atmp := a[i*lda : i*lda+i]
			for j, v := range atmp {
				x[j] += xi * v
			}
			if nonUnit {
				x[i] *= a[i*lda+i]
			}
		}
		return
	}
	ix := kx
	for i := 0; i < n; i++ {
		xi := x[ix]
		jx := kx
		atmp := a[i*lda : i*lda+i]
		for _, v := range atmp {
			x[jx] += xi * v
			jx += incX
		}
		if nonUnit {
			x[ix] *= a[i*lda+i]
		}
		ix += incX
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
		kx = 0
	} else {
		kx = -(n - 1) * incX
	}
	if incY > 0 {
		ky = 0
	} else {
		ky = -(n - 1) * incY
	}

	// Form y = beta * y
	if beta != 1 {
		if incY > 0 {
			b.Dscal(n, beta, y, incY)
		} else {
			b.Dscal(n, beta, y, -incY)
		}
	}

	if alpha == 0 {
		return
	}

	if n == 1 {
		y[0] += alpha * a[0] * x[0]
		return
	}

	if ul == blas.Upper {
		if incX == 1 {
			iy := ky
			for i := 0; i < n; i++ {
				xv := x[i] * alpha
				sum := x[i] * a[i*lda+i]
				jy := ky + (i+1)*incY
				atmp := a[i*lda+i+1 : i*lda+n]
				for j, v := range atmp {
					jp := j + i + 1
					sum += x[jp] * v
					y[jy] += xv * v
					jy += incY
				}
				y[iy] += alpha * sum
				iy += incY
			}
			return
		}
		ix := kx
		iy := ky
		for i := 0; i < n; i++ {
			xv := x[ix] * alpha
			sum := x[ix] * a[i*lda+i]
			jx := kx + (i+1)*incX
			jy := ky + (i+1)*incY
			atmp := a[i*lda+i+1 : i*lda+n]
			for _, v := range atmp {
				sum += x[jx] * v
				y[jy] += xv * v
				jx += incX
				jy += incY
			}
			y[iy] += alpha * sum
			ix += incX
			iy += incY
		}
		return
	}
	// Cases where a is lower triangular.
	if incX == 1 {
		iy := ky
		for i := 0; i < n; i++ {
			jy := ky
			xv := alpha * x[i]
			atmp := a[i*lda : i*lda+i]
			var sum float64
			for j, v := range atmp {
				sum += x[j] * v
				y[jy] += xv * v
				jy += incY
			}
			sum += x[i] * a[i*lda+i]
			sum *= alpha
			y[iy] += sum
			iy += incY
		}
		return
	}
	ix := kx
	iy := ky
	for i := 0; i < n; i++ {
		jx := kx
		jy := ky
		xv := alpha * x[ix]
		atmp := a[i*lda : i*lda+i]
		var sum float64
		for _, v := range atmp {
			sum += x[jx] * v
			y[jy] += xv * v
			jx += incX
			jy += incY
		}
		sum += x[ix] * a[i*lda+i]
		sum *= alpha
		y[iy] += sum
		ix += incX
		iy += incY
	}
}

// Dtbmv  performs one of the matrix-vector operations
// 		x := A*x,   or   x := A**T*x,
// where x is an n element vector and  A is an n by n unit, or non-unit,
// upper or lower triangular band matrix.
func (Blas) Dtbmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n, k int, a []float64, lda int, x []float64, incX int) {
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
	_ = kx

	nonunit := d != blas.Unit

	if tA == blas.NoTrans {
		if ul == blas.Upper {
			if incX == 1 {
				for i := 0; i < n; i++ {
					u := min(1+k, n-i)
					var sum float64
					atmp := a[i*lda:]
					xtmp := x[i:]
					for j := 1; j < u; j++ {
						sum += xtmp[j] * atmp[j]
					}
					if nonunit {
						sum += xtmp[0] * atmp[0]
					} else {
						sum += xtmp[0]
					}
					x[i] = sum
				}
				return
			}
			ix := kx
			for i := 0; i < n; i++ {
				u := min(1+k, n-i)
				var sum float64
				atmp := a[i*lda:]
				jx := incX
				for j := 1; j < u; j++ {
					sum += x[ix+jx] * atmp[j]
					jx += incX
				}
				if nonunit {
					sum += x[ix] * atmp[0]
				} else {
					sum += x[ix]
				}
				x[ix] = sum
				ix += incX
			}
			return
		}
		if incX == 1 {
			for i := n - 1; i >= 0; i-- {
				l := max(0, k-i)
				atmp := a[i*lda:]
				var sum float64
				for j := l; j < k; j++ {
					sum += x[i-k+j] * atmp[j]
				}
				if nonunit {
					sum += x[i] * atmp[k]
				} else {
					sum += x[i]
				}
				x[i] = sum
			}
			return
		}
		ix := kx + (n-1)*incX
		for i := n - 1; i >= 0; i-- {
			l := max(0, k-i)
			atmp := a[i*lda:]
			var sum float64
			jx := l * incX
			for j := l; j < k; j++ {
				sum += x[ix-k*incX+jx] * atmp[j]
				jx += incX
			}
			if nonunit {
				sum += x[ix] * atmp[k]
			} else {
				sum += x[ix]
			}
			x[ix] = sum
			ix -= incX
		}
		return
	}
	if ul == blas.Upper {
		if incX == 1 {
			for i := n - 1; i >= 0; i-- {
				u := k + 1
				if i < u {
					u = i + 1
				}
				var sum float64
				for j := 1; j < u; j++ {
					sum += x[i-j] * a[(i-j)*lda+j]
				}
				if nonunit {
					sum += x[i] * a[i*lda]
				} else {
					sum += x[i]
				}
				x[i] = sum
			}
			return
		}
		ix := kx + (n-1)*incX
		for i := n - 1; i >= 0; i-- {
			u := k + 1
			if i < u {
				u = i + 1
			}
			var sum float64
			jx := incX
			for j := 1; j < u; j++ {
				sum += x[ix-jx] * a[(i-j)*lda+j]
				jx += incX
			}
			if nonunit {
				sum += x[ix] * a[i*lda]
			} else {
				sum += x[ix]
			}
			x[ix] = sum
			ix -= incX
		}
		return
	}
	if incX == 1 {
		for i := 0; i < n; i++ {
			u := k
			if i+k >= n {
				u = n - i - 1
			}
			var sum float64
			for j := 0; j < u; j++ {
				sum += x[i+j+1] * a[(i+j+1)*lda+k-j-1]
			}
			if nonunit {
				sum += x[i] * a[i*lda+k]
			} else {
				sum += x[i]
			}
			x[i] = sum
		}
		return
	}
	ix := kx
	for i := 0; i < n; i++ {
		u := k
		if i+k >= n {
			u = n - i - 1
		}
		var (
			sum float64
			jx  int
		)
		for j := 0; j < u; j++ {
			sum += x[ix+jx+incX] * a[(i+j+1)*lda+k-j-1]
			jx += incX
		}
		if nonunit {
			sum += x[ix] * a[i*lda+k]
		} else {
			sum += x[ix]
		}
		x[ix] = sum
		ix += incX
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

// Dtbsv solves A * x = b where A is a triangular banded matrix with k diagonals
// above the main diagonal. A has compact banded storage.
// The result in stored in-place into x.
func (Blas) Dtbsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n, k int, a []float64, lda int, x []float64, incX int) {
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
	if lda < k+1 {
		panic(badLdaTriBand)
	}
	if incX == 0 {
		panic(zeroInc)
	}
	if n == 0 {
		return
	}
	var kx int
	if incX < 0 {
		kx = -(n - 1) * incX
	} else {
		kx = 0
	}
	nonUnit := d == blas.NonUnit
	// Form x = A^-1 x.
	// Several cases below use subslices for speed improvement.
	// The incX != 1 cases usually do not because incX may be negative.
	if tA == blas.NoTrans {
		if ul == blas.Upper {
			if incX == 1 {
				for i := n - 1; i >= 0; i-- {
					bands := k
					if i+bands >= n {
						bands = n - i - 1
					}
					atmp := a[i*lda+1:]
					xtmp := x[i+1 : i+bands+1]
					var sum float64
					for j, v := range xtmp {
						sum += v * atmp[j]
					}
					x[i] -= sum
					if nonUnit {
						x[i] /= a[i*lda]
					}
				}
				return
			}
			ix := kx + (n-1)*incX
			for i := n - 1; i >= 0; i-- {
				max := k + 1
				if i+max > n {
					max = n - i
				}
				atmp := a[i*lda:]
				var (
					jx  int
					sum float64
				)
				for j := 1; j < max; j++ {
					jx += incX
					sum += x[ix+jx] * atmp[j]
				}
				x[ix] -= sum
				if nonUnit {
					x[ix] /= atmp[0]
				}
				ix -= incX
			}
			return
		}
		if incX == 1 {
			for i := 0; i < n; i++ {
				bands := k
				if i-k < 0 {
					bands = i
				}
				atmp := a[i*lda+k-bands:]
				xtmp := x[i-bands : i]
				var sum float64
				for j, v := range xtmp {
					sum += v * atmp[j]
				}
				x[i] -= sum
				if nonUnit {
					x[i] /= atmp[bands]
				}
			}
			return
		}
		ix := kx
		for i := 0; i < n; i++ {
			bands := k
			if i-k < 0 {
				bands = i
			}
			atmp := a[i*lda+k-bands:]
			var (
				sum float64
				jx  int
			)
			for j := 0; j < bands; j++ {
				sum += x[ix-bands*incX+jx] * atmp[j]
				jx += incX
			}
			x[ix] -= sum
			if nonUnit {
				x[ix] /= atmp[bands]
			}
			ix += incX
		}
		return
	}
	// Cases where a is transposed.
	if ul == blas.Upper {
		if incX == 1 {
			for i := 0; i < n; i++ {
				bands := k
				if i-k < 0 {
					bands = i
				}
				var sum float64
				for j := 0; j < bands; j++ {
					sum += x[i-bands+j] * a[(i-bands+j)*lda+bands-j]
				}
				x[i] -= sum
				if nonUnit {
					x[i] /= a[i*lda]
				}
			}
			return
		}
		ix := kx
		for i := 0; i < n; i++ {
			bands := k
			if i-k < 0 {
				bands = i
			}
			var (
				sum float64
				jx  int
			)
			for j := 0; j < bands; j++ {
				sum += x[ix-bands*incX+jx] * a[(i-bands+j)*lda+bands-j]
				jx += incX
			}
			x[ix] -= sum
			if nonUnit {
				x[ix] /= a[i*lda]
			}
			ix += incX
		}
		return
	}
	if incX == 1 {
		for i := n - 1; i >= 0; i-- {
			bands := k
			if i+bands >= n {
				bands = n - i - 1
			}
			var sum float64
			xtmp := x[i+1 : i+1+bands]
			for j, v := range xtmp {
				sum += v * a[(i+j+1)*lda+k-j-1]
			}
			x[i] -= sum
			if nonUnit {
				x[i] /= a[i*lda+k]
			}
		}
		return
	}
	ix := kx + (n-1)*incX
	for i := n - 1; i >= 0; i-- {
		bands := k
		if i+bands >= n {
			bands = n - i - 1
		}
		var (
			sum float64
			jx  int
		)
		for j := 0; j < bands; j++ {
			sum += x[ix+jx+incX] * a[(i+j+1)*lda+k-j-1]
			jx += incX
		}
		x[ix] -= sum
		if nonUnit {
			x[ix] /= a[i*lda+k]
		}
		ix -= incX
	}
}

// Dsbmv performs y = alpha*A*x + beta*y where A is a symmetric banded matrix.
func (b Blas) Dsbmv(ul blas.Uplo, n, k int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int) {
	if ul != blas.Lower && ul != blas.Upper {
		panic(badUplo)
	}
	if n < 0 {
		panic(nLT0)
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

	// Set up indexes
	lenX := n
	lenY := n
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

	if ul == blas.Upper {
		if incX == 1 {
			iy := ky
			for i := 0; i < n; i++ {
				atmp := a[i*lda:]
				tmp := alpha * x[i]
				sum := tmp * atmp[0]
				u := min(k, n-i-1)
				jy := incY
				for j := 1; j <= u; j++ {
					v := atmp[j]
					sum += alpha * x[i+j] * v
					y[iy+jy] += tmp * v
					jy += incY
				}
				y[iy] += sum
				iy += incY
			}
			return
		}
		ix := kx
		iy := ky
		for i := 0; i < n; i++ {
			atmp := a[i*lda:]
			tmp := alpha * x[ix]
			sum := tmp * atmp[0]
			u := min(k, n-i-1)
			jx := incX
			jy := incY
			for j := 1; j <= u; j++ {
				v := atmp[j]
				sum += alpha * x[ix+jx] * v
				y[iy+jy] += tmp * v
				jx += incX
				jy += incY
			}
			y[iy] += sum
			ix += incX
			iy += incY
		}
		return
	}

	// Casses where a has bands below the diagonal.
	if incX == 1 {
		iy := ky
		for i := 0; i < n; i++ {
			l := max(0, k-i)
			tmp := alpha * x[i]
			jy := l * incY
			atmp := a[i*lda:]
			for j := l; j < k; j++ {
				v := atmp[j]
				y[iy] += alpha * v * x[i-k+j]
				y[iy-k*incY+jy] += tmp * v
				jy += incY
			}
			y[iy] += tmp * atmp[k]
			iy += incY
		}
		return
	}
	ix := kx
	iy := ky
	for i := 0; i < n; i++ {
		l := max(0, k-i)
		tmp := alpha * x[ix]
		jx := l * incX
		jy := l * incY
		atmp := a[i*lda:]
		for j := l; j < k; j++ {
			v := atmp[j]
			y[iy] += alpha * v * x[ix-k*incX+jx]
			y[iy-k*incY+jy] += tmp * v
			jx += incX
			jy += incY
		}
		y[iy] += tmp * atmp[k]
		ix += incX
		iy += incY
	}
	return
}

// Dsyr computes a = alpha*x*x^T + a where a is an nxn symmetric matrix
func (Blas) Dsyr(ul blas.Uplo, n int, alpha float64, x []float64, incX int, a []float64, lda int) {
	if ul != blas.Lower && ul != blas.Upper {
		panic(badUplo)
	}
	if n < 0 {
		panic(nLT0)
	}
	if incX == 0 {
		panic(negInc)
	}
	if lda < n {
		panic(badLda)
	}
	if alpha == 0 || n == 0 {
		return
	}

	lenX := n
	var kx int
	if incX > 0 {
		kx = 0
	} else {
		kx = -(lenX - 1) * incX
	}
	if ul == blas.Upper {
		if incX == 1 {
			for i := 0; i < n; i++ {
				tmp := x[i] * alpha
				if tmp != 0 {
					atmp := a[i*lda+i : i*lda+n]
					xtmp := x[i:n]
					for j, v := range xtmp {
						atmp[j] += v * tmp
					}
				}
			}
			return
		}
		ix := kx
		for i := 0; i < n; i++ {
			tmp := x[ix] * alpha
			if tmp != 0 {
				jx := ix
				atmp := a[i*lda:]
				for j := i; j < n; j++ {
					atmp[j] += x[jx] * tmp
					jx += incX
				}
			}
			ix += incX
		}
		return
	}
	// Cases where a is lower triangular.
	if incX == 1 {
		for i := 0; i < n; i++ {
			tmp := x[i] * alpha
			if tmp != 0 {
				atmp := a[i*lda:]
				xtmp := x[:i+1]
				for j, v := range xtmp {
					atmp[j] += tmp * v
				}
			}
		}
		return
	}
	ix := kx
	for i := 0; i < n; i++ {
		tmp := x[ix] * alpha
		if tmp != 0 {
			atmp := a[i*lda:]
			jx := kx
			for j := 0; j < i+1; j++ {
				atmp[j] += tmp * x[jx]
				jx += incX
			}
		}
		ix += incX
	}
}

//TODO: Not yet implemented Level 2 routines.
func (Blas) Dtpsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, ap []float64, x []float64, incX int) {
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
func (Blas) Dsyr2(ul blas.Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int) {
	panic("referenceblas: function not implemented")
}
