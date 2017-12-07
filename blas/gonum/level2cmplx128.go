// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math/cmplx"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/internal/asm/c128"
)

// Zgemv performs one of the matrix-vector operations
//  y = alpha * A * x + beta * y    if trans = blas.NoTrans
//  y = alpha * A^T * x + beta * y  if trans = blas.Trans
//  y = alpha * A^H * x + beta * y  if trans = blas.ConjTrans
// where alpha and beta are scalars, x and y are vectors, and A is an m×n dense matrix.
func (Implementation) Zgemv(trans blas.Transpose, m, n int, alpha complex128, a []complex128, lda int, x []complex128, incX int, beta complex128, y []complex128, incY int) {
	checkZMatrix('A', m, n, a, lda)
	switch trans {
	default:
		panic(badTranspose)
	case blas.NoTrans:
		checkZVector('x', n, x, incX)
		checkZVector('y', m, y, incY)
	case blas.Trans, blas.ConjTrans:
		checkZVector('x', m, x, incX)
		checkZVector('y', n, y, incY)
	}

	if m == 0 || n == 0 || (alpha == 0 && beta == 1) {
		return
	}

	var lenX, lenY int
	if trans == blas.NoTrans {
		lenX = n
		lenY = m
	} else {
		lenX = m
		lenY = n
	}
	var kx int
	if incX < 0 {
		kx = (1 - lenX) * incX
	}
	var ky int
	if incY < 0 {
		ky = (1 - lenY) * incY
	}

	// Form y := beta*y.
	if beta != 1 {
		if incY == 1 {
			if beta == 0 {
				for i := range y {
					y[i] = 0
				}
			} else {
				c128.ScalUnitary(beta, y[:lenY])
			}
		} else {
			iy := ky
			if beta == 0 {
				for i := 0; i < lenY; i++ {
					y[iy] = 0
					iy += incY
				}
			} else {
				if incY > 0 {
					c128.ScalInc(beta, y, uintptr(lenY), uintptr(incY))
				} else {
					c128.ScalInc(beta, y, uintptr(lenY), uintptr(-incY))
				}
			}
		}
	}

	if alpha == 0 {
		return
	}

	switch trans {
	default:
		// Form y := alpha*A*x + y.
		iy := ky
		if incX == 1 {
			for i := 0; i < m; i++ {
				y[iy] += alpha * c128.DotuUnitary(a[i*lda:i*lda+n], x[:n])
				iy += incY
			}
			return
		}
		for i := 0; i < m; i++ {
			y[iy] += alpha * c128.DotuInc(a[i*lda:i*lda+n], x, uintptr(n), 1, uintptr(incX), 0, uintptr(kx))
			iy += incY
		}
		return

	case blas.Trans:
		// Form y := alpha*A^T*x + y.
		ix := kx
		if incY == 1 {
			for i := 0; i < m; i++ {
				c128.AxpyUnitary(alpha*x[ix], a[i*lda:i*lda+n], y[:n])
				ix += incX
			}
			return
		}
		for i := 0; i < m; i++ {
			c128.AxpyInc(alpha*x[ix], a[i*lda:i*lda+n], y, uintptr(n), 1, uintptr(incY), 0, uintptr(ky))
			ix += incX
		}
		return

	case blas.ConjTrans:
		// Form y := alpha*A^H*x + y.
		ix := kx
		if incY == 1 {
			for i := 0; i < m; i++ {
				tmp := alpha * x[ix]
				for j := 0; j < n; j++ {
					y[j] += tmp * cmplx.Conj(a[i*lda+j])
				}
				ix += incX
			}
			return
		}
		for i := 0; i < m; i++ {
			tmp := alpha * x[ix]
			jy := ky
			for j := 0; j < n; j++ {
				y[jy] += tmp * cmplx.Conj(a[i*lda+j])
				jy += incY
			}
			ix += incX
		}
		return
	}
}

// Zgerc performs the rank-one operation
//  A += alpha * x * y^H
// where A is an m×n dense matrix, alpha is a scalar, x is an m element vector,
// and y is an n element vector.
func (Implementation) Zgerc(m, n int, alpha complex128, x []complex128, incX int, y []complex128, incY int, a []complex128, lda int) {
	checkZMatrix('A', m, n, a, lda)
	checkZVector('x', m, x, incX)
	checkZVector('y', n, y, incY)

	if m == 0 || n == 0 || alpha == 0 {
		return
	}

	var kx, jy int
	if incX < 0 {
		kx = (1 - m) * incX
	}
	if incY < 0 {
		jy = (1 - n) * incY
	}
	for j := 0; j < n; j++ {
		if y[jy] != 0 {
			tmp := alpha * cmplx.Conj(y[jy])
			c128.AxpyInc(tmp, x, a[j:], uintptr(m), uintptr(incX), uintptr(lda), uintptr(kx), 0)
		}
		jy += incY
	}
}

// Zgeru performs the rank-one operation
//  A += alpha * x * y^T
// where A is an m×n dense matrix, alpha is a scalar, x is an m element vector,
// and y is an n element vector.
func (Implementation) Zgeru(m, n int, alpha complex128, x []complex128, incX int, y []complex128, incY int, a []complex128, lda int) {
	checkZMatrix('A', m, n, a, lda)
	checkZVector('x', m, x, incX)
	checkZVector('y', n, y, incY)

	if m == 0 || n == 0 || alpha == 0 {
		return
	}

	var kx int
	if incX < 0 {
		kx = (1 - m) * incX
	}
	if incY == 1 {
		for i := 0; i < m; i++ {
			if x[kx] != 0 {
				tmp := alpha * x[kx]
				c128.AxpyUnitary(tmp, y[:n], a[i*lda:i*lda+n])
			}
			kx += incX
		}
		return
	}
	var jy int
	if incY < 0 {
		jy = (1 - n) * incY
	}
	for i := 0; i < m; i++ {
		if x[kx] != 0 {
			tmp := alpha * x[kx]
			c128.AxpyInc(tmp, y, a[i*lda:i*lda+n], uintptr(n), uintptr(incY), 1, uintptr(jy), 0)
		}
		kx += incX
	}
}

// Zhemv performs the matrix-vector operation
//  y = alpha * A * x + beta * y
// where alpha and beta are scalars, x and y are vectors, and A is an n×n
// Hermitian matrix. The imaginary parts of the diagonal elements of A are
// ignored and assumed to be zero.
func (Implementation) Zhemv(uplo blas.Uplo, n int, alpha complex128, a []complex128, lda int, x []complex128, incX int, beta complex128, y []complex128, incY int) {
	if uplo != blas.Upper && uplo != blas.Lower {
		panic(badUplo)
	}
	checkZMatrix('A', n, n, a, lda)
	checkZVector('x', n, x, incX)
	checkZVector('y', n, y, incY)

	if n == 0 || (alpha == 0 && beta == 1) {
		return
	}

	// Set up the start indices in X and Y.
	var kx int
	if incX < 0 {
		kx = (1 - n) * incX
	}
	var ky int
	if incY < 0 {
		ky = (1 - n) * incY
	}

	// Form y := beta*y.
	if beta != 1 {
		if incY == 1 {
			if beta == 0 {
				for i := range y {
					y[i] = 0
				}
			} else {
				for i, v := range y {
					y[i] = beta * v
				}
			}
		} else {
			iy := ky
			if beta == 0 {
				for i := 0; i < n; i++ {
					y[iy] = 0
					iy += incY
				}
			} else {
				for i := 0; i < n; i++ {
					y[iy] = beta * y[iy]
					iy += incY
				}
			}
		}
	}

	if alpha == 0 {
		return
	}

	// The elements of A are accessed sequentially with one pass through
	// the triangular part of A.

	if uplo == blas.Upper {
		// Form y when A is stored in upper triangle.
		if incX == 1 && incY == 1 {
			for i := 0; i < n; i++ {
				tmp1 := alpha * x[i]
				var tmp2 complex128
				for j := i + 1; j < n; j++ {
					y[j] += tmp1 * cmplx.Conj(a[i*lda+j])
					tmp2 += a[i*lda+j] * x[j]
				}
				aii := complex(real(a[i*lda+i]), 0)
				y[i] += tmp1*aii + alpha*tmp2
			}
		} else {
			ix := kx
			iy := ky
			for i := 0; i < n; i++ {
				tmp1 := alpha * x[ix]
				var tmp2 complex128
				jx := ix
				jy := iy
				for j := i + 1; j < n; j++ {
					jx += incX
					jy += incY
					y[jy] += tmp1 * cmplx.Conj(a[i*lda+j])
					tmp2 += a[i*lda+j] * x[jx]
				}
				aii := complex(real(a[i*lda+i]), 0)
				y[iy] += tmp1*aii + alpha*tmp2
				ix += incX
				iy += incY
			}
		}
		return
	}

	// Form y when A is stored in lower triangle.
	if incX == 1 && incY == 1 {
		for i := 0; i < n; i++ {
			tmp1 := alpha * x[i]
			var tmp2 complex128
			for j := 0; j < i; j++ {
				y[j] += tmp1 * cmplx.Conj(a[i*lda+j])
				tmp2 += a[i*lda+j] * x[j]
			}
			aii := complex(real(a[i*lda+i]), 0)
			y[i] += tmp1*aii + alpha*tmp2
		}
	} else {
		ix := kx
		iy := ky
		for i := 0; i < n; i++ {
			tmp1 := alpha * x[ix]
			var tmp2 complex128
			jx := kx
			jy := ky
			for j := 0; j < i; j++ {
				y[jy] += tmp1 * cmplx.Conj(a[i*lda+j])
				tmp2 += a[i*lda+j] * x[jx]
				jx += incX
				jy += incY
			}
			aii := complex(real(a[i*lda+i]), 0)
			y[iy] += tmp1*aii + alpha*tmp2
			ix += incX
			iy += incY
		}
	}
}

// Zher performs the Hermitian rank-one operation
//  A += alpha * x * x^H
// where A is an n×n Hermitian matrix, alpha is a real scalar, and x is an n
// element vector. On entry, the imaginary parts of the diagonal elements of A
// are ignored and assumed to be zero, on return they will be set to zero.
func (Implementation) Zher(uplo blas.Uplo, n int, alpha float64, x []complex128, incX int, a []complex128, lda int) {
	if uplo != blas.Upper && uplo != blas.Lower {
		panic(badUplo)
	}
	checkZMatrix('A', n, n, a, lda)
	checkZVector('x', n, x, incX)

	if n == 0 || alpha == 0 {
		return
	}

	var kx int
	if incX < 0 {
		kx = (1 - n) * incX
	}
	if uplo == blas.Upper {
		if incX == 1 {
			for i := 0; i < n; i++ {
				if x[i] != 0 {
					tmp := complex(alpha*real(x[i]), alpha*imag(x[i]))
					aii := real(a[i*lda+i])
					xtmp := real(tmp * cmplx.Conj(x[i]))
					a[i*lda+i] = complex(aii+xtmp, 0)
					for j := i + 1; j < n; j++ {
						a[i*lda+j] += tmp * cmplx.Conj(x[j])
					}
				} else {
					aii := real(a[i*lda+i])
					a[i*lda+i] = complex(aii, 0)
				}
			}
			return
		}

		ix := kx
		for i := 0; i < n; i++ {
			if x[ix] != 0 {
				tmp := complex(alpha*real(x[ix]), alpha*imag(x[ix]))
				aii := real(a[i*lda+i])
				xtmp := real(tmp * cmplx.Conj(x[ix]))
				a[i*lda+i] = complex(aii+xtmp, 0)
				jx := ix + incX
				for j := i + 1; j < n; j++ {
					a[i*lda+j] += tmp * cmplx.Conj(x[jx])
					jx += incX
				}
			} else {
				aii := real(a[i*lda+i])
				a[i*lda+i] = complex(aii, 0)
			}
			ix += incX
		}
		return
	}

	if incX == 1 {
		for i := 0; i < n; i++ {
			if x[i] != 0 {
				tmp := complex(alpha*real(x[i]), alpha*imag(x[i]))
				for j := 0; j < i; j++ {
					a[i*lda+j] += tmp * cmplx.Conj(x[j])
				}
				aii := real(a[i*lda+i])
				xtmp := real(tmp * cmplx.Conj(x[i]))
				a[i*lda+i] = complex(aii+xtmp, 0)
			} else {
				aii := real(a[i*lda+i])
				a[i*lda+i] = complex(aii, 0)
			}
		}
		return
	}

	ix := kx
	for i := 0; i < n; i++ {
		if x[ix] != 0 {
			tmp := complex(alpha*real(x[ix]), alpha*imag(x[ix]))
			jx := kx
			for j := 0; j < i; j++ {
				a[i*lda+j] += tmp * cmplx.Conj(x[jx])
				jx += incX
			}
			aii := real(a[i*lda+i])
			xtmp := real(tmp * cmplx.Conj(x[ix]))
			a[i*lda+i] = complex(aii+xtmp, 0)

		} else {
			aii := real(a[i*lda+i])
			a[i*lda+i] = complex(aii, 0)
		}
		ix += incX
	}
}

// Zher2 performs the Hermitian rank-two operation
//  A += alpha*x*y^H + conj(alpha)*y*x^H
// where alpha is a scalar, x and y are n element vectors and A is an n×n
// Hermitian matrix. On entry, the imaginary parts of the diagonal elements are
// ignored and assumed to be zero. On return they will be set to zero.
func (Implementation) Zher2(uplo blas.Uplo, n int, alpha complex128, x []complex128, incX int, y []complex128, incY int, a []complex128, lda int) {
	if uplo != blas.Upper && uplo != blas.Lower {
		panic(badUplo)
	}
	checkZMatrix('A', n, n, a, lda)
	checkZVector('x', n, x, incX)
	checkZVector('y', n, y, incY)

	if n == 0 || alpha == 0 {
		return
	}

	var kx, ky int
	var ix, iy int
	if incX != 1 || incY != 1 {
		if incX < 0 {
			kx = (1 - n) * incX
		}
		if incY < 0 {
			ky = (1 - n) * incY
		}
		ix = kx
		iy = ky
	}
	if uplo == blas.Upper {
		if incX == 1 && incY == 1 {
			for i := 0; i < n; i++ {
				if x[i] != 0 || y[i] != 0 {
					tmp1 := alpha * x[i]
					tmp2 := cmplx.Conj(alpha) * y[i]
					aii := real(a[i*lda+i]) + real(tmp1*cmplx.Conj(y[i])) + real(tmp2*cmplx.Conj(x[i]))
					a[i*lda+i] = complex(aii, 0)
					for j := i + 1; j < n; j++ {
						a[i*lda+j] += tmp1*cmplx.Conj(y[j]) + tmp2*cmplx.Conj(x[j])
					}
				} else {
					aii := real(a[i*lda+i])
					a[i*lda+i] = complex(aii, 0)
				}
			}
			return
		}
		for i := 0; i < n; i++ {
			if x[i] != 0 || y[i] != 0 {
				tmp1 := alpha * x[ix]
				tmp2 := cmplx.Conj(alpha) * y[iy]
				aii := real(a[i*lda+i]) + real(tmp1*cmplx.Conj(y[iy])) + real(tmp2*cmplx.Conj(x[ix]))
				a[i*lda+i] = complex(aii, 0)
				jx := ix + incX
				jy := iy + incY
				for j := i + 1; j < n; j++ {
					a[i*lda+j] += tmp1*cmplx.Conj(y[jy]) + tmp2*cmplx.Conj(x[jx])
					jx += incX
					jy += incY
				}
			} else {
				aii := real(a[i*lda+i])
				a[i*lda+i] = complex(aii, 0)
			}
			ix += incX
			iy += incY
		}
		return
	}

	if incX == 1 && incY == 1 {
		for i := 0; i < n; i++ {
			if x[i] != 0 || y[i] != 0 {
				tmp1 := alpha * x[i]
				tmp2 := cmplx.Conj(alpha) * y[i]
				for j := 0; j < i; j++ {
					a[i*lda+j] += tmp1*cmplx.Conj(y[j]) + tmp2*cmplx.Conj(x[j])
				}
				aii := real(a[i*lda+i]) + real(tmp1*cmplx.Conj(y[i])) + real(tmp2*cmplx.Conj(x[i]))
				a[i*lda+i] = complex(aii, 0)
			} else {
				aii := real(a[i*lda+i])
				a[i*lda+i] = complex(aii, 0)
			}
		}
		return
	}
	for i := 0; i < n; i++ {
		if x[i] != 0 || y[i] != 0 {
			tmp1 := alpha * x[ix]
			tmp2 := cmplx.Conj(alpha) * y[iy]
			jx := kx
			jy := ky
			for j := 0; j < i; j++ {
				a[i*lda+j] += tmp1*cmplx.Conj(y[jy]) + tmp2*cmplx.Conj(x[jx])
				jx += incX
				jy += incY
			}
			aii := real(a[i*lda+i]) + real(tmp1*cmplx.Conj(y[iy])) + real(tmp2*cmplx.Conj(x[ix]))
			a[i*lda+i] = complex(aii, 0)
		} else {
			aii := real(a[i*lda+i])
			a[i*lda+i] = complex(aii, 0)
		}
		ix += incX
		iy += incY
	}
}

// Zhpr performs the Hermitian rank-1 operation
//  A += alpha * x * x^H,
// where alpha is a real scalar, x is a vector, and A is an n×n hermitian matrix
// in packed form. On entry, the imaginary parts of the diagonal elements are
// assumed to be zero, and on return they are set to zero.
func (Implementation) Zhpr(uplo blas.Uplo, n int, alpha float64, x []complex128, incX int, ap []complex128) {
	if uplo != blas.Upper && uplo != blas.Lower {
		panic(badUplo)
	}
	if n < 0 {
		panic(nLT0)
	}
	checkZVector('x', n, x, incX)
	if len(ap) < n*(n+1)/2 {
		panic("blas: insufficient A packed matrix slice length")
	}

	if n == 0 || alpha == 0 {
		return
	}

	// Set up start index in X.
	var kx int
	if incX < 0 {
		kx = (1 - n) * incX
	}

	// The elements of A are accessed sequentially with one pass through ap.

	var kk int
	if uplo == blas.Upper {
		// Form A when upper triangle is stored in AP.
		// Here, kk points to the current diagonal element in ap.
		if incX == 1 {
			for i := 0; i < n; i++ {
				xi := x[i]
				if xi != 0 {
					aii := real(ap[kk]) + alpha*real(cmplx.Conj(xi)*xi)
					ap[kk] = complex(aii, 0)

					tmp := complex(alpha, 0) * xi
					a := ap[kk+1 : kk+n-i]
					x := x[i+1 : n]
					for j, v := range x {
						a[j] += tmp * cmplx.Conj(v)
					}
				} else {
					ap[kk] = complex(real(ap[kk]), 0)
				}
				kk += n - i
			}
		} else {
			ix := kx
			for i := 0; i < n; i++ {
				xi := x[ix]
				if xi != 0 {
					aii := real(ap[kk]) + alpha*real(cmplx.Conj(xi)*xi)
					ap[kk] = complex(aii, 0)

					tmp := complex(alpha, 0) * xi
					jx := ix + incX
					a := ap[kk+1 : kk+n-i]
					for k := range a {
						a[k] += tmp * cmplx.Conj(x[jx])
						jx += incX
					}
				} else {
					ap[kk] = complex(real(ap[kk]), 0)
				}
				ix += incX
				kk += n - i
			}
		}
		return
	}

	// Form A when lower triangle is stored in AP.
	// Here, kk points to the beginning of current row in ap.
	if incX == 1 {
		for i := 0; i < n; i++ {
			xi := x[i]
			if xi != 0 {
				tmp := complex(alpha, 0) * xi
				a := ap[kk : kk+i]
				for j, v := range x[:i] {
					a[j] += tmp * cmplx.Conj(v)
				}

				aii := real(ap[kk+i]) + alpha*real(cmplx.Conj(xi)*xi)
				ap[kk+i] = complex(aii, 0)
			} else {
				ap[kk+i] = complex(real(ap[kk+i]), 0)
			}
			kk += i + 1
		}
	} else {
		ix := kx
		for i := 0; i < n; i++ {
			xi := x[ix]
			if xi != 0 {
				tmp := complex(alpha, 0) * xi
				a := ap[kk : kk+i]
				jx := kx
				for k := range a {
					a[k] += tmp * cmplx.Conj(x[jx])
					jx += incX
				}

				aii := real(ap[kk+i]) + alpha*real(cmplx.Conj(xi)*xi)
				ap[kk+i] = complex(aii, 0)
			} else {
				ap[kk+i] = complex(real(ap[kk+i]), 0)
			}
			ix += incX
			kk += i + 1
		}
	}
}

// Ztrmv performs one of the matrix-vector operations
//  x = A * x    if trans = blas.NoTrans
//  x = A^T * x  if trans = blas.Trans
//  x = A^H * x  if trans = blas.ConjTrans
// where x is a vector, and A is an n×n triangular matrix.
func (Implementation) Ztrmv(uplo blas.Uplo, trans blas.Transpose, diag blas.Diag, n int, a []complex128, lda int, x []complex128, incX int) {
	if uplo != blas.Upper && uplo != blas.Lower {
		panic(badUplo)
	}
	if trans != blas.NoTrans && trans != blas.Trans && trans != blas.ConjTrans {
		panic(badTranspose)
	}
	if diag != blas.Unit && diag != blas.NonUnit {
		panic(badDiag)
	}
	checkZMatrix('A', n, n, a, lda)
	checkZVector('x', n, x, incX)

	if n == 0 {
		return
	}

	// Set up start index in X.
	var kx int
	if incX < 0 {
		kx = (1 - n) * incX
	}

	// The elements of A are accessed sequentially with one pass through A.

	if trans == blas.NoTrans {
		// Form x := A*x.
		if uplo == blas.Upper {
			if incX == 1 {
				for i := 0; i < n; i++ {
					if diag == blas.NonUnit {
						x[i] *= a[i*lda+i]
					}
					if n-i-1 > 0 {
						x[i] += c128.DotuUnitary(a[i*lda+i+1:i*lda+n], x[i+1:n])
					}
				}
			} else {
				ix := kx
				for i := 0; i < n; i++ {
					if diag == blas.NonUnit {
						x[ix] *= a[i*lda+i]
					}
					if n-i-1 > 0 {
						x[ix] += c128.DotuInc(a[i*lda+i+1:i*lda+n], x, uintptr(n-i-1), 1, uintptr(incX), 0, uintptr(ix+incX))
					}
					ix += incX
				}
			}
		} else {
			if incX == 1 {
				for i := n - 1; i >= 0; i-- {
					if diag == blas.NonUnit {
						x[i] *= a[i*lda+i]
					}
					if i > 0 {
						x[i] += c128.DotuUnitary(a[i*lda:i*lda+i], x[:i])
					}
				}
			} else {
				ix := kx + (n-1)*incX
				for i := n - 1; i >= 0; i-- {
					if diag == blas.NonUnit {
						x[ix] *= a[i*lda+i]
					}
					if i > 0 {
						x[ix] += c128.DotuInc(a[i*lda:i*lda+i], x, uintptr(i), 1, uintptr(incX), 0, uintptr(kx))
					}
					ix -= incX
				}
			}
		}
		return
	}

	if trans == blas.Trans {
		// Form x := A^T*x.
		if uplo == blas.Upper {
			if incX == 1 {
				for i := n - 1; i >= 0; i-- {
					xi := x[i]
					if diag == blas.NonUnit {
						x[i] *= a[i*lda+i]
					}
					if n-i-1 > 0 {
						c128.AxpyUnitary(xi, a[i*lda+i+1:i*lda+n], x[i+1:n])
					}
				}
			} else {
				ix := kx + (n-1)*incX
				for i := n - 1; i >= 0; i-- {
					xi := x[ix]
					if diag == blas.NonUnit {
						x[ix] *= a[i*lda+i]
					}
					if n-i-1 > 0 {
						c128.AxpyInc(xi, a[i*lda+i+1:i*lda+n], x, uintptr(n-i-1), 1, uintptr(incX), 0, uintptr(ix+incX))
					}
					ix -= incX
				}
			}
		} else {
			if incX == 1 {
				for i := 0; i < n; i++ {
					if i > 0 {
						c128.AxpyUnitary(x[i], a[i*lda:i*lda+i], x[:i])
					}
					if diag == blas.NonUnit {
						x[i] *= a[i*lda+i]
					}
				}
			} else {
				ix := kx
				for i := 0; i < n; i++ {
					if i > 0 {
						c128.AxpyInc(x[ix], a[i*lda:i*lda+i], x, uintptr(i), 1, uintptr(incX), 0, uintptr(kx))
					}
					if diag == blas.NonUnit {
						x[ix] *= a[i*lda+i]
					}
					ix += incX
				}
			}
		}
		return
	}

	// Form x := A^H*x.
	if uplo == blas.Upper {
		if incX == 1 {
			for i := n - 1; i >= 0; i-- {
				xi := x[i]
				if diag == blas.NonUnit {
					x[i] *= cmplx.Conj(a[i*lda+i])
				}
				for j := i + 1; j < n; j++ {
					x[j] += xi * cmplx.Conj(a[i*lda+j])
				}
			}
		} else {
			ix := kx + (n-1)*incX
			for i := n - 1; i >= 0; i-- {
				xi := x[ix]
				if diag == blas.NonUnit {
					x[ix] *= cmplx.Conj(a[i*lda+i])
				}
				jx := ix + incX
				for j := i + 1; j < n; j++ {
					x[jx] += xi * cmplx.Conj(a[i*lda+j])
					jx += incX
				}
				ix -= incX
			}
		}
	} else {
		if incX == 1 {
			for i := 0; i < n; i++ {
				for j := 0; j < i; j++ {
					x[j] += x[i] * cmplx.Conj(a[i*lda+j])
				}
				if diag == blas.NonUnit {
					x[i] *= cmplx.Conj(a[i*lda+i])
				}
			}
		} else {
			ix := kx
			for i := 0; i < n; i++ {
				jx := kx
				for j := 0; j < i; j++ {
					x[jx] += x[ix] * cmplx.Conj(a[i*lda+j])
					jx += incX
				}
				if diag == blas.NonUnit {
					x[ix] *= cmplx.Conj(a[i*lda+i])
				}
				ix += incX
			}
		}
	}
}

// Ztrsv solves one of the systems of equations
//  A*x = b     if trans == blas.NoTrans,
//  A^T*x = b,  if trans == blas.Trans,
//  A^H*x = b,  if trans == blas.ConjTrans,
// where b and x are n element vectors and A is an n×n triangular matrix.
//
// On entry, x contains the values of b, and the solution is
// stored in-place into x.
//
// No test for singularity or near-singularity is included in this
// routine. Such tests must be performed before calling this routine.
func (Implementation) Ztrsv(uplo blas.Uplo, trans blas.Transpose, diag blas.Diag, n int, a []complex128, lda int, x []complex128, incX int) {
	if uplo != blas.Upper && uplo != blas.Lower {
		panic(badUplo)
	}
	if trans != blas.NoTrans && trans != blas.Trans && trans != blas.ConjTrans {
		panic(badTranspose)
	}
	if diag != blas.Unit && diag != blas.NonUnit {
		panic(badDiag)
	}
	checkZMatrix('A', n, n, a, lda)
	checkZVector('x', n, x, incX)

	if n == 0 {
		return
	}

	// Set up start index in X.
	var kx int
	if incX < 0 {
		kx = (1 - n) * incX
	}

	// The elements of A are accessed sequentially with one pass through A.

	if trans == blas.NoTrans {
		// Form x := inv(A)*x.
		if uplo == blas.Upper {
			if incX == 1 {
				for i := n - 1; i >= 0; i-- {
					aii := a[i*lda+i]
					if n-i-1 > 0 {
						x[i] -= c128.DotuUnitary(x[i+1:n], a[i*lda+i+1:i*lda+n])
					}
					if diag == blas.NonUnit {
						x[i] /= aii
					}
				}
			} else {
				ix := kx + (n-1)*incX
				for i := n - 1; i >= 0; i-- {
					aii := a[i*lda+i]
					if n-i-1 > 0 {
						x[ix] -= c128.DotuInc(x, a[i*lda+i+1:i*lda+n], uintptr(n-i-1), uintptr(incX), 1, uintptr(ix+incX), 0)
					}
					if diag == blas.NonUnit {
						x[ix] /= aii
					}
					ix -= incX
				}
			}
		} else {
			if incX == 1 {
				for i := 0; i < n; i++ {
					if i > 0 {
						x[i] -= c128.DotuUnitary(x[:i], a[i*lda:i*lda+i])
					}
					if diag == blas.NonUnit {
						x[i] /= a[i*lda+i]
					}
				}
			} else {
				ix := kx
				for i := 0; i < n; i++ {
					if i > 0 {
						x[ix] -= c128.DotuInc(x, a[i*lda:i*lda+i], uintptr(i), uintptr(incX), 1, uintptr(kx), 0)
					}
					if diag == blas.NonUnit {
						x[ix] /= a[i*lda+i]
					}
					ix += incX
				}
			}
		}
		return
	}

	if trans == blas.Trans {
		// Form x := inv(A^T)*x.
		if uplo == blas.Upper {
			if incX == 1 {
				for j := 0; j < n; j++ {
					if diag == blas.NonUnit {
						x[j] /= a[j*lda+j]
					}
					if n-j-1 > 0 {
						c128.AxpyUnitary(-x[j], a[j*lda+j+1:j*lda+n], x[j+1:n])
					}
				}
			} else {
				jx := kx
				for j := 0; j < n; j++ {
					if diag == blas.NonUnit {
						x[jx] /= a[j*lda+j]
					}
					if n-j-1 > 0 {
						c128.AxpyInc(-x[jx], a[j*lda+j+1:j*lda+n], x, uintptr(n-j-1), 1, uintptr(incX), 0, uintptr(jx+incX))
					}
					jx += incX
				}
			}
		} else {
			if incX == 1 {
				for j := n - 1; j >= 0; j-- {
					if diag == blas.NonUnit {
						x[j] /= a[j*lda+j]
					}
					xj := x[j]
					if j > 0 {
						c128.AxpyUnitary(-xj, a[j*lda:j*lda+j], x[:j])
					}
				}
			} else {
				jx := kx + (n-1)*incX
				for j := n - 1; j >= 0; j-- {
					if diag == blas.NonUnit {
						x[jx] /= a[j*lda+j]
					}
					if j > 0 {
						c128.AxpyInc(-x[jx], a[j*lda:j*lda+j], x, uintptr(j), 1, uintptr(incX), 0, uintptr(kx))
					}
					jx -= incX
				}
			}
		}
		return
	}

	// Form x := inv(A^H)*x.
	if uplo == blas.Upper {
		if incX == 1 {
			for j := 0; j < n; j++ {
				if diag == blas.NonUnit {
					x[j] /= cmplx.Conj(a[j*lda+j])
				}
				xj := x[j]
				for i := j + 1; i < n; i++ {
					x[i] -= xj * cmplx.Conj(a[j*lda+i])
				}
			}
		} else {
			jx := kx
			for j := 0; j < n; j++ {
				if diag == blas.NonUnit {
					x[jx] /= cmplx.Conj(a[j*lda+j])
				}
				xj := x[jx]
				ix := jx + incX
				for i := j + 1; i < n; i++ {
					x[ix] -= xj * cmplx.Conj(a[j*lda+i])
					ix += incX
				}
				jx += incX
			}
		}
	} else {
		if incX == 1 {
			for j := n - 1; j >= 0; j-- {
				if diag == blas.NonUnit {
					x[j] /= cmplx.Conj(a[j*lda+j])
				}
				xj := x[j]
				for i := 0; i < j; i++ {
					x[i] -= xj * cmplx.Conj(a[j*lda+i])
				}
			}
		} else {
			jx := kx + (n-1)*incX
			for j := n - 1; j >= 0; j-- {
				if diag == blas.NonUnit {
					x[jx] /= cmplx.Conj(a[j*lda+j])
				}
				xj := x[jx]
				ix := kx
				for i := 0; i < j; i++ {
					x[ix] -= xj * cmplx.Conj(a[j*lda+i])
					ix += incX
				}
				jx -= incX
			}
		}
	}
}
