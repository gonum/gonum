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
)

type Dpstf2er interface {
	Dpstf2(uplo blas.Uplo, n int, a []float64, lda int, piv []int, tol float64, work []float64) (rank int, ok bool)
}

func Dpstf2Test(t *testing.T, impl Dpstf2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower} {
		t.Run(uploToString(uplo), func(t *testing.T) {
			for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 20, 50} {
				for _, lda := range []int{max(1, n), n + 5} {
					for _, rank := range []int{int(0.7 * float64(n)), n} {
						dpstf2Test(t, impl, rnd, uplo, n, lda, rank)
					}
				}
			}
		})
	}
}

func dpstf2Test(t *testing.T, impl Dpstf2er, rnd *rand.Rand, uplo blas.Uplo, n, lda, rankWant int) {
	const tol = 1e-14

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

	// Call Dpstf2 to Compute the Cholesky factorization with complete pivoting.
	rank, ok := impl.Dpstf2(uplo, n, aFac, lda, piv, -1, work)

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
