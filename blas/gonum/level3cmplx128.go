// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math/cmplx"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/internal/asm/c128"
)

var _ blas.Complex128Level3 = Implementation{}

// Zgemm performs one of the matrix-matrix operations
//  C = alpha * op(A) * op(B) + beta * C
// where op(X) is one of
//  op(X) = X  or  op(X) = X^T  or  op(X) = X^H,
// alpha and beta are scalars, and A, B and C are matrices, with op(A) an m×k matrix,
// op(B) a k×n matrix and C an m×n matrix.
func (Implementation) Zgemm(tA, tB blas.Transpose, m, n, k int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta complex128, c []complex128, ldc int) {
	switch tA {
	default:
		panic(badTranspose)
	case blas.NoTrans, blas.Trans, blas.ConjTrans:
	}
	switch tB {
	default:
		panic(badTranspose)
	case blas.NoTrans, blas.Trans, blas.ConjTrans:
	}
	switch {
	case m < 0:
		panic(mLT0)
	case n < 0:
		panic(nLT0)
	case k < 0:
		panic(kLT0)
	}
	rowA, colA := m, k
	if tA != blas.NoTrans {
		rowA, colA = k, m
	}
	if lda < max(1, colA) {
		panic(badLdA)
	}
	rowB, colB := k, n
	if tB != blas.NoTrans {
		rowB, colB = n, k
	}
	if ldb < max(1, colB) {
		panic(badLdB)
	}
	if ldc < max(1, n) {
		panic(badLdC)
	}

	// Quick return if possible.
	if m == 0 || n == 0 {
		return
	}

	// For zero matrix size the following slice length checks are trivially satisfied.
	if len(a) < (rowA-1)*lda+colA {
		panic(shortA)
	}
	if len(b) < (rowB-1)*ldb+colB {
		panic(shortB)
	}
	if len(c) < (m-1)*ldc+n {
		panic(shortC)
	}

	// Quick return if possible.
	if (alpha == 0 || k == 0) && beta == 1 {
		return
	}

	if alpha == 0 {
		if beta == 0 {
			for i := 0; i < m; i++ {
				for j := 0; j < n; j++ {
					c[i*ldc+j] = 0
				}
			}
		} else {
			for i := 0; i < m; i++ {
				for j := 0; j < n; j++ {
					c[i*ldc+j] *= beta
				}
			}
		}
		return
	}

	switch tA {
	case blas.NoTrans:
		switch tB {
		case blas.NoTrans:
			// Form  C = alpha * A * B + beta * C.
			for i := 0; i < m; i++ {
				if beta == 0 {
					for j := 0; j < n; j++ {
						c[i*ldc+j] = 0
					}
				} else if beta != 1 {
					for j := 0; j < n; j++ {
						c[i*ldc+j] *= beta
					}
				}
				for l := 0; l < k; l++ {
					tmp := alpha * a[i*lda+l]
					for j := 0; j < n; j++ {
						c[i*ldc+j] += tmp * b[l*ldb+j]
					}
				}
			}
		case blas.Trans:
			// Form  C = alpha * A * B^T + beta * C.
			for i := 0; i < m; i++ {
				if beta == 0 {
					for j := 0; j < n; j++ {
						c[i*ldc+j] = 0
					}
				} else if beta != 1 {
					for j := 0; j < n; j++ {
						c[i*ldc+j] *= beta
					}
				}
				for l := 0; l < k; l++ {
					tmp := alpha * a[i*lda+l]
					for j := 0; j < n; j++ {
						c[i*ldc+j] += tmp * b[j*ldb+l]
					}
				}
			}
		case blas.ConjTrans:
			// Form  C = alpha * A * B^H + beta * C.
			for i := 0; i < m; i++ {
				if beta == 0 {
					for j := 0; j < n; j++ {
						c[i*ldc+j] = 0
					}
				} else if beta != 1 {
					for j := 0; j < n; j++ {
						c[i*ldc+j] *= beta
					}
				}
				for l := 0; l < k; l++ {
					tmp := alpha * a[i*lda+l]
					for j := 0; j < n; j++ {
						c[i*ldc+j] += tmp * cmplx.Conj(b[j*ldb+l])
					}
				}
			}
		}
	case blas.Trans:
		switch tB {
		case blas.NoTrans:
			// Form  C = alpha * A^T * B + beta * C.
			for i := 0; i < m; i++ {
				for j := 0; j < n; j++ {
					var tmp complex128
					for l := 0; l < k; l++ {
						tmp += a[l*lda+i] * b[l*ldb+j]
					}
					if beta == 0 {
						c[i*ldc+j] = alpha * tmp
					} else {
						c[i*ldc+j] = alpha*tmp + beta*c[i*ldc+j]
					}
				}
			}
		case blas.Trans:
			// Form  C = alpha * A^T * B^T + beta * C.
			for i := 0; i < m; i++ {
				for j := 0; j < n; j++ {
					var tmp complex128
					for l := 0; l < k; l++ {
						tmp += a[l*lda+i] * b[j*ldb+l]
					}
					if beta == 0 {
						c[i*ldc+j] = alpha * tmp
					} else {
						c[i*ldc+j] = alpha*tmp + beta*c[i*ldc+j]
					}
				}
			}
		case blas.ConjTrans:
			// Form  C = alpha * A^T * B^H + beta * C.
			for i := 0; i < m; i++ {
				for j := 0; j < n; j++ {
					var tmp complex128
					for l := 0; l < k; l++ {
						tmp += a[l*lda+i] * cmplx.Conj(b[j*ldb+l])
					}
					if beta == 0 {
						c[i*ldc+j] = alpha * tmp
					} else {
						c[i*ldc+j] = alpha*tmp + beta*c[i*ldc+j]
					}
				}
			}
		}
	case blas.ConjTrans:
		switch tB {
		case blas.NoTrans:
			// Form  C = alpha * A^H * B + beta * C.
			for i := 0; i < m; i++ {
				for j := 0; j < n; j++ {
					var tmp complex128
					for l := 0; l < k; l++ {
						tmp += cmplx.Conj(a[l*lda+i]) * b[l*ldb+j]
					}
					if beta == 0 {
						c[i*ldc+j] = alpha * tmp
					} else {
						c[i*ldc+j] = alpha*tmp + beta*c[i*ldc+j]
					}
				}
			}
		case blas.Trans:
			// Form  C = alpha * A^H * B^T + beta * C.
			for i := 0; i < m; i++ {
				for j := 0; j < n; j++ {
					var tmp complex128
					for l := 0; l < k; l++ {
						tmp += cmplx.Conj(a[l*lda+i]) * b[j*ldb+l]
					}
					if beta == 0 {
						c[i*ldc+j] = alpha * tmp
					} else {
						c[i*ldc+j] = alpha*tmp + beta*c[i*ldc+j]
					}
				}
			}
		case blas.ConjTrans:
			// Form  C = alpha * A^H * B^H + beta * C.
			for i := 0; i < m; i++ {
				for j := 0; j < n; j++ {
					var tmp complex128
					for l := 0; l < k; l++ {
						tmp += cmplx.Conj(a[l*lda+i]) * cmplx.Conj(b[j*ldb+l])
					}
					if beta == 0 {
						c[i*ldc+j] = alpha * tmp
					} else {
						c[i*ldc+j] = alpha*tmp + beta*c[i*ldc+j]
					}
				}
			}
		}
	}
}

// Zherk performs one of the hermitian rank-k operations
//  C = alpha*A*A^H + beta*C  if trans == blas.NoTrans
//  C = alpha*A^H*A + beta*C  if trans == blas.ConjTrans
// where alpha and beta are real scalars, C is an n×n hermitian matrix and A is
// an n×k matrix in the first case and a k×n matrix in the second case.
//
// The imaginary parts of the diagonal elements of C are assumed to be zero, and
// on return they will be set to zero.
func (Implementation) Zherk(uplo blas.Uplo, trans blas.Transpose, n, k int, alpha float64, a []complex128, lda int, beta float64, c []complex128, ldc int) {
	var rowA, colA int
	switch trans {
	default:
		panic(badTranspose)
	case blas.NoTrans:
		rowA, colA = n, k
	case blas.ConjTrans:
		rowA, colA = k, n
	}
	switch {
	case uplo != blas.Lower && uplo != blas.Upper:
		panic(badUplo)
	case n < 0:
		panic(nLT0)
	case k < 0:
		panic(kLT0)
	case lda < max(1, colA):
		panic(badLdA)
	case ldc < max(1, n):
		panic(badLdC)
	}

	// Quick return if possible.
	if n == 0 {
		return
	}

	// For zero matrix size the following slice length checks are trivially satisfied.
	if len(a) < (rowA-1)*lda+colA {
		panic(shortA)
	}
	if len(c) < (n-1)*ldc+n {
		panic(shortC)
	}

	// Quick return if possible.
	if (alpha == 0 || k == 0) && beta == 1 {
		return
	}

	if alpha == 0 {
		if uplo == blas.Upper {
			if beta == 0 {
				for i := 0; i < n; i++ {
					ci := c[i*ldc+i : i*ldc+n]
					for j := range ci {
						ci[j] = 0
					}
				}
			} else {
				for i := 0; i < n; i++ {
					ci := c[i*ldc+i : i*ldc+n]
					ci[0] = complex(beta*real(ci[0]), 0)
					if i != n-1 {
						c128.DscalUnitary(beta, ci[1:])
					}
				}
			}
		} else {
			if beta == 0 {
				for i := 0; i < n; i++ {
					ci := c[i*ldc : i*ldc+i+1]
					for j := range ci {
						ci[j] = 0
					}
				}
			} else {
				for i := 0; i < n; i++ {
					ci := c[i*ldc : i*ldc+i+1]
					if i != 0 {
						c128.DscalUnitary(beta, ci[:i])
					}
					ci[i] = complex(beta*real(ci[i]), 0)
				}
			}
		}
		return
	}

	calpha := complex(alpha, 0)
	if trans == blas.NoTrans {
		// Form  C = alpha*A*A^H + beta*C.
		cbeta := complex(beta, 0)
		if uplo == blas.Upper {
			for i := 0; i < n; i++ {
				ci := c[i*ldc+i : i*ldc+n]
				ai := a[i*lda : i*lda+k]
				// Handle the i-th diagonal element of C.
				cii := calpha*c128.DotcUnitary(ai, ai) + cbeta*ci[0]
				ci[0] = complex(real(cii), 0)
				// Handle the remaining elements on the i-th row of C.
				for jc, cij := range ci[1:] {
					j := i + 1 + jc
					ci[jc+1] = calpha*c128.DotcUnitary(a[j*lda:j*lda+k], ai) + cbeta*cij
				}
			}
		} else {
			for i := 0; i < n; i++ {
				ci := c[i*ldc : i*ldc+i+1]
				ai := a[i*lda : i*lda+k]
				// Handle the first i-1 elements on the i-th row of C.
				for j, cij := range ci[:i] {
					ci[j] = calpha*c128.DotcUnitary(a[j*lda:j*lda+k], ai) + cbeta*cij
				}
				// Handle the i-th diagonal element of C.
				cii := calpha*c128.DotcUnitary(ai, ai) + cbeta*ci[i]
				ci[i] = complex(real(cii), 0)
			}
		}
	} else {
		// Form  C = alpha*A^H*A + beta*C.
		if uplo == blas.Upper {
			for i := 0; i < n; i++ {
				ci := c[i*ldc+i : i*ldc+n]
				if beta == 0 {
					for jc := range ci {
						ci[jc] = 0
					}
				} else if beta != 1 {
					c128.DscalUnitary(beta, ci)
					ci[0] = complex(real(ci[0]), 0)
				} else {
					ci[0] = complex(real(ci[0]), 0)
				}
				for j := 0; j < k; j++ {
					aji := cmplx.Conj(a[j*lda+i])
					if aji != 0 {
						c128.AxpyUnitary(calpha*aji, a[j*lda+i:j*lda+n], ci)
					}
				}
				c[i*ldc+i] = complex(real(c[i*ldc+i]), 0)

			}
		} else {
			for i := 0; i < n; i++ {
				ci := c[i*ldc : i*ldc+i+1]
				if beta == 0 {
					for j := range ci {
						ci[j] = 0
					}
				} else if beta != 1 {
					c128.DscalUnitary(beta, ci)
					ci[i] = complex(real(ci[i]), 0)
				} else {
					ci[i] = complex(real(ci[i]), 0)
				}
				for j := 0; j < k; j++ {
					aji := cmplx.Conj(a[j*lda+i])
					if aji != 0 {
						c128.AxpyUnitary(calpha*aji, a[j*lda:j*lda+i+1], ci)
					}
				}
				c[i*ldc+i] = complex(real(c[i*ldc+i]), 0)
			}
		}
	}
}

// Zsyrk performs one of the symmetric rank-k operations
//  C = alpha*A*A^T + beta*C  if trans == blas.NoTrans
//  C = alpha*A^T*A + beta*C  if trans == blas.Trans
// where alpha and beta are scalars, C is an n×n symmetric matrix and A is
// an n×k matrix in the first case and a k×n matrix in the second case.
func (Implementation) Zsyrk(uplo blas.Uplo, trans blas.Transpose, n, k int, alpha complex128, a []complex128, lda int, beta complex128, c []complex128, ldc int) {
	var rowA, colA int
	switch trans {
	default:
		panic(badTranspose)
	case blas.NoTrans:
		rowA, colA = n, k
	case blas.Trans:
		rowA, colA = k, n
	}
	switch {
	case uplo != blas.Lower && uplo != blas.Upper:
		panic(badUplo)
	case n < 0:
		panic(nLT0)
	case k < 0:
		panic(kLT0)
	case lda < max(1, colA):
		panic(badLdA)
	case ldc < max(1, n):
		panic(badLdC)
	}

	// Quick return if possible.
	if n == 0 {
		return
	}

	// For zero matrix size the following slice length checks are trivially satisfied.
	if len(a) < (rowA-1)*lda+colA {
		panic(shortA)
	}
	if len(c) < (n-1)*ldc+n {
		panic(shortC)
	}

	// Quick return if possible.
	if (alpha == 0 || k == 0) && beta == 1 {
		return
	}

	if alpha == 0 {
		if uplo == blas.Upper {
			if beta == 0 {
				for i := 0; i < n; i++ {
					ci := c[i*ldc+i : i*ldc+n]
					for j := range ci {
						ci[j] = 0
					}
				}
			} else {
				for i := 0; i < n; i++ {
					ci := c[i*ldc+i : i*ldc+n]
					c128.ScalUnitary(beta, ci)
				}
			}
		} else {
			if beta == 0 {
				for i := 0; i < n; i++ {
					ci := c[i*ldc : i*ldc+i+1]
					for j := range ci {
						ci[j] = 0
					}
				}
			} else {
				for i := 0; i < n; i++ {
					ci := c[i*ldc : i*ldc+i+1]
					c128.ScalUnitary(beta, ci)
				}
			}
		}
		return
	}

	if trans == blas.NoTrans {
		// Form  C = alpha*A*A^T + beta*C.
		if uplo == blas.Upper {
			for i := 0; i < n; i++ {
				ci := c[i*ldc+i : i*ldc+n]
				ai := a[i*lda : i*lda+k]
				for jc, cij := range ci {
					j := i + jc
					ci[jc] = beta*cij + alpha*c128.DotuUnitary(ai, a[j*lda:j*lda+k])
				}
			}
		} else {
			for i := 0; i < n; i++ {
				ci := c[i*ldc : i*ldc+i+1]
				ai := a[i*lda : i*lda+k]
				for j, cij := range ci {
					ci[j] = beta*cij + alpha*c128.DotuUnitary(ai, a[j*lda:j*lda+k])
				}
			}
		}
	} else {
		// Form  C = alpha*A^T*A + beta*C.
		if uplo == blas.Upper {
			for i := 0; i < n; i++ {
				ci := c[i*ldc+i : i*ldc+n]
				if beta == 0 {
					for jc := range ci {
						ci[jc] = 0
					}
				} else if beta != 1 {
					for jc := range ci {
						ci[jc] *= beta
					}
				}
				for j := 0; j < k; j++ {
					aji := a[j*lda+i]
					if aji != 0 {
						c128.AxpyUnitary(alpha*aji, a[j*lda+i:j*lda+n], ci)
					}
				}
			}
		} else {
			for i := 0; i < n; i++ {
				ci := c[i*ldc : i*ldc+i+1]
				if beta == 0 {
					for j := range ci {
						ci[j] = 0
					}
				} else if beta != 1 {
					for j := range ci {
						ci[j] *= beta
					}
				}
				for j := 0; j < k; j++ {
					aji := a[j*lda+i]
					if aji != 0 {
						c128.AxpyUnitary(alpha*aji, a[j*lda:j*lda+i+1], ci)
					}
				}
			}
		}
	}
}
