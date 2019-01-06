// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math/cmplx"

	"gonum.org/v1/gonum/blas"
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
