// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dpstrfer interface {
	Dpstrf(uplo blas.Uplo, n int, a []float64, lda int, piv []int, tol float64, work []float64) (rank int, ok bool)
}

func DpstrfTest(t *testing.T, impl Dpstrfer) {
	rnd := rand.New(rand.NewSource(1))
	for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower} {
		t.Run(uploToString(uplo), func(t *testing.T) {
			for _, n := range []int{0, 1, 2, 3, 4, 5, 31, 32, 33, 63, 64, 65, 127, 128, 129} {
				for _, lda := range []int{max(1, n), n + 5} {
					for _, rank := range []int{int(0.7 * float64(n)), n} {
						dpstrfTest(t, impl, rnd, uplo, n, lda, rank)
					}
				}
			}
		})
	}
}

func dpstrfTest(t *testing.T, impl Dpstrfer, rnd *rand.Rand, uplo blas.Uplo, n, lda, rankWant int) {
	const tol = 1e-13

	name := fmt.Sprintf("n=%v,lda=%v", n, lda)
	bi := blas64.Implementation()

	// Generate a random, symmetric A with the given rank by applying rankWant
	// rank-1 updates to the zero matrix.
	a := make([]float64, n*lda)
	for i := 0; i < rankWant; i++ {
		x := randomSlice(n, rnd)
		bi.Dsyr(uplo, n, 1, x, 1, a, lda)
	}

	// Make a copy of A for storing the factorization.
	aFac := make([]float64, len(a))
	copy(aFac, a)

	// Allocate a slice for pivots and fill it with invalid index values.
	piv := make([]int, n)
	for i := range piv {
		piv[i] = -1
	}

	// Allocate the work slice.
	work := make([]float64, 2*n)

	// Call Dpstrf to Compute the Cholesky factorization with complete pivoting.
	rank, ok := impl.Dpstrf(uplo, n, aFac, lda, piv, -1, work)

	if ok != (rank == n) {
		t.Errorf("%v: unexpected ok; got %v, want %v", name, ok, rank == n)
	}
	if rank != rankWant {
		t.Errorf("%v: unexpected rank; got %v, want %v", name, rank, rankWant)
	}

	if n == 0 {
		return
	}

	// Check that the residual |P*Uᵀ*U*Pᵀ - A| / n or |P*L*Lᵀ*Pᵀ - A| / n is
	// sufficiently small.
	resid := residualDpstrf(uplo, n, a, aFac, lda, rank, piv)
	if resid > tol || math.IsNaN(resid) {
		t.Errorf("%v: residual too large; got %v, want<=%v", name, resid, tol)
	}
}

func residualDpstrf(uplo blas.Uplo, n int, a, aFac []float64, lda int, rank int, piv []int) float64 {
	bi := blas64.Implementation()
	// Reconstruct the symmetric positive semi-definite matrix A from its L or U
	// factors and the permutation matrix P.
	perm := zeros(n, n, n)
	if uplo == blas.Upper {
		// Change notation.
		u, ldu := aFac, lda
		// Zero out last n-rank rows of the factor U.
		for i := rank; i < n; i++ {
			for j := i; j < n; j++ {
				u[i*ldu+j] = 0
			}
		}
		// Extract U to aRec.
		aRec := zeros(n, n, n)
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				aRec.Data[i*aRec.Stride+j] = u[i*ldu+j]
			}
		}
		// Multiply U by Uᵀ from the left.
		bi.Dtrmm(blas.Left, blas.Upper, blas.Trans, blas.NonUnit, n, n,
			1, u, ldu, aRec.Data, aRec.Stride)
		// Form P * Uᵀ * U * Pᵀ.
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if piv[i] > piv[j] {
					// Don't set the lower triangle.
					continue
				}
				if i <= j {
					perm.Data[piv[i]*perm.Stride+piv[j]] = aRec.Data[i*aRec.Stride+j]
				} else {
					perm.Data[piv[i]*perm.Stride+piv[j]] = aRec.Data[j*aRec.Stride+i]
				}
			}
		}
		// Compute the difference P*Uᵀ*U*Pᵀ - A.
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				perm.Data[i*perm.Stride+j] -= a[i*lda+j]
			}
		}
	} else {
		// Change notation.
		l, ldl := aFac, lda
		// Zero out last n-rank columns of the factor L.
		for i := rank; i < n; i++ {
			for j := rank; j <= i; j++ {
				l[i*ldl+j] = 0
			}
		}
		// Extract L to aRec.
		aRec := zeros(n, n, n)
		for i := 0; i < n; i++ {
			for j := 0; j <= i; j++ {
				aRec.Data[i*aRec.Stride+j] = l[i*ldl+j]
			}
		}
		// Multiply L by Lᵀ from the right.
		bi.Dtrmm(blas.Right, blas.Lower, blas.Trans, blas.NonUnit, n, n,
			1, l, ldl, aRec.Data, aRec.Stride)
		// Form P * L * Lᵀ * Pᵀ.
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if piv[i] < piv[j] {
					// Don't set the upper triangle.
					continue
				}
				if i >= j {
					perm.Data[piv[i]*perm.Stride+piv[j]] = aRec.Data[i*aRec.Stride+j]
				} else {
					perm.Data[piv[i]*perm.Stride+piv[j]] = aRec.Data[j*aRec.Stride+i]
				}
			}
		}
		// Compute the difference P*L*Lᵀ*Pᵀ - A.
		for i := 0; i < n; i++ {
			for j := 0; j <= i; j++ {
				perm.Data[i*perm.Stride+j] -= a[i*lda+j]
			}
		}
	}
	// Compute |P*Uᵀ*U*Pᵀ - A| / n or |P*L*Lᵀ*Pᵀ - A| / n.
	return dlansy(lapack.MaxColumnSum, uplo, n, perm.Data, perm.Stride) / float64(n)
}
