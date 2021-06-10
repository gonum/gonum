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

type Dgetc2er interface {
	Dgetc2(n int, a []float64, lda int, ipiv, jpiv []int) (k int)
}

func Dgetc2Test(t *testing.T, impl Dgetc2er) {
	const tol = 1e-12
	rnd := rand.New(rand.NewSource(1))
	for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 20} {
		for _, lda := range []int{n, n + 5} {
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
	// ipiv and jpiv are outputs.
	ipiv := make([]int, n)
	jpiv := make([]int, n)
	for i := 0; i < n; i++ {
		ipiv[i], jpiv[i] = -1, -1 // Set to non-indices.
	}
	// Copy to store output (LU decomposition).
	lu := cloneGeneral(a)
	k := impl.Dgetc2(n, lu.Data, lu.Stride, ipiv, jpiv)
	if k >= 0 {
		t.Logf("%v: matrix was perturbed at %d", name, k)
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
	// Construct L and U triangular matrices from Dgetc2 output.
	L := zeros(n, n, lda)
	U := zeros(n, n, lda)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			idx := i*lda + j
			if j >= i { // On upper triangle and setting of L's unit diagonal elements.
				U.Data[idx] = lu.Data[idx]
				if j == i {
					L.Data[idx] = 1.0
				}
			} else if i > j { // On diagonal or lower triangle.
				L.Data[idx] = lu.Data[idx]
			}
		}
	}
	work := zeros(n, n, lda)
	bi.Dgemm(blas.NoTrans, blas.NoTrans, n, n, n, 1, L.Data, L.Stride, U.Data, U.Stride, 0, work.Data, work.Stride)

	// Apply Permutations P and Q to L*U.
	for i := n - 1; i >= 0; i-- {
		ipv, jpv := ipiv[i], jpiv[i]
		if ipv != i {
			bi.Dswap(n, work.Data[i*lda:], 1, work.Data[ipv*lda:], 1)
		}
		if jpv != i {
			bi.Dswap(n, work.Data[i:], work.Stride, work.Data[jpv:], work.Stride)
		}
	}

	// A should be reconstructed by now.
	for i := range work.Data {
		if math.Abs(work.Data[i]-a.Data[i]) > tol {
			t.Errorf("%v: matrix %d idx not equal after reconstruction. got %g, expected %g", name, i, work.Data[i], a.Data[i])
		}
	}
}
