// Copyright ©2015 The Gonum Authors. All rights reserved.
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

type Dgetc2er interface {
	Dgetc2(n int, a []float64, lda int, ipiv, jpiv []int) (k int)
}

func Dgetc2Test(t *testing.T, impl Dgetc2er) {
	const tol = 1e-12
	rnd := rand.New(rand.NewSource(1))
	for _, n := range []int{2} { // []int{0, 1, 2, 3, 4, 5, 10, 20}
		for _, lda := range []int{n} {
			dgetc2Test(t, impl, rnd, n, lda, tol)
		}
	}
}

func dgetc2Test(t *testing.T, impl Dgetc2er, rnd *rand.Rand, n, lda int, tol float64) {
	name := fmt.Sprintf("n=%v,lda=%v", n, lda)
	if lda == 0 {
		lda = 1
	}
	// Generate a random general matrix A.
	a := randomGeneral(n, n, lda, rnd)

	// ipib and jpiv are outputs.
	ipiv := make([]int, n)
	jpiv := make([]int, n)
	for i := 0; i < n; i++ {
		ipiv[i], jpiv[i] = -1, -1 // set them to non-indices
	}
	// Copy to store output (LU decomposition)
	lu := make([]float64, len(a.Data))
	copy(lu, a.Data)
	k := impl.Dgetc2(n, lu, lda, ipiv, jpiv)
	if k >= 0 {
		t.Fatalf("%v: matrix was perturbed at %d", name, k)
	}

	// Verify all indices are set.
	for i := 0; i < n; i++ {
		if ipiv[i] < 0 {
			t.Errorf("%v: ipiv[%d] is negative", name, i)
		}
		if jpiv[i] < 0 {
			t.Errorf("%v: jpiv[%d] is negative", name, i)
		}
	}
	bi := blas64.Implementation()

	// Reconstruct P and Q permutation matrices from ipiv and jpiv, respectively.
	P := make([]float64, n*n)
	Q := make([]float64, n*n)
	for i := 0; i < n; i++ {
		ipv, jpv := ipiv[i], jpiv[i]
		// ipiv/jpiv indicates column set to one (one per row, see Permutation Matrix).
		P[i*n+ipv] = 1.0
		Q[i*n+jpv] = 1.0
	}
	// Construct L and U triangular matrices from Dgetc2 output.
	L := make([]float64, n*n) //
	U := make([]float64, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			idx := i*n + j
			if j >= i { // on upper triangle and setting of L's unit diagonal elements
				U[idx] = lu[idx]
				if j == i {
					L[idx] = 1.0
				}
			} else if i > j { // on diagonal or lower triangle
				L[idx] = lu[idx]
			}
		}
	}
	// results for multiplication matrix for P * L * U * Q
	result1 := make([]float64, n*n)
	result2 := make([]float64, n*n)
	// Dgemm does  C = alpha * A * B + beta * C
	bi.Dgemm(blas.NoTrans, blas.NoTrans, n, n, n, 1, P, lda, L, lda, 0, result1, lda)
	bi.Dgemm(blas.NoTrans, blas.NoTrans, n, n, n, 1, result1, lda, U, lda, 0, result2, lda)
	bi.Dgemm(blas.NoTrans, blas.NoTrans, n, n, n, 1, result2, lda, Q, lda, 0, result1, lda)
	// result1 should now be equal to A
	for i := range lu {
		if math.Abs(result1[i]-a.Data[i]) > tol {
			t.Errorf("%v: matrix %d idx not equal after reconstruction. got %g, expected %g", name, i, result1[i], a.Data[i])
		}
	}
}
