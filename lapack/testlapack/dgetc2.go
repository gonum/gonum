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
	for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 20} {
		for _, lda := range []int{n} {
			dgetc2Test(t, impl, rnd, n, lda, tol)
		}
	}
	// specific matrix cases
	for _, tc := range []struct {
		name       string
		a, expect  blas64.General
		ipiv, jpiv []int
	}{
		{name: "identity", a: eye(3, 3), expect: eye(3, 3), ipiv: []int{2, 2, 2}, jpiv: []int{2, 2, 2}},
		{
			name:   "small",
			a:      blas64.General{Rows: 3, Cols: 3, Stride: 3, Data: []float64{1, 2, 3, 2, 1, 6, 3, 6, 0}},
			expect: blas64.General{Rows: 3, Cols: 3, Stride: 3, Data: []float64{6, 0, 3, 1. / 6., 6, 1.5, 1. / 3., 0.5, -0.75}},
			ipiv:   []int{2, 1, 2},
			jpiv:   []int{1, 2, 2},
		},
	} {
		name := fmt.Sprintf("%s %dx%d", tc.name, tc.a.Rows, tc.a.Cols)
		n := len(tc.jpiv)
		ipiv, jpiv := make([]int, n), make([]int, n)
		impl.Dgetc2(n, tc.a.Data, tc.a.Stride, ipiv, jpiv)
		for i := 0; i < len(tc.a.Data); i++ {
			got := tc.a.Data[i]
			expect := tc.expect.Data[i]
			if math.Abs(got-expect) > tol {
				t.Errorf("%s: expected %.8g in A matrix. got %g", name, expect, got)
			}
		}
		for i := 0; i < n; i++ {
			if ipiv[i] != tc.ipiv[i] {
				t.Errorf("%s: expected %d in ipiv pivots. got %d", name, tc.ipiv[i], ipiv[i])
			}
			if jpiv[i] != tc.jpiv[i] {
				t.Errorf("%s: expected %d in jpiv pivots. got %d", name, tc.jpiv[i], jpiv[i])
			}
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
	// Construct L and U triangular matrices from Dgetc2 output.
	L := make([]float64, n*n) //
	U := make([]float64, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			idx := i*n + j
			if j >= i { // On upper triangle and setting of L's unit diagonal elements.
				U[idx] = lu[idx]
				if j == i {
					L[idx] = 1.0
				}
			} else if i > j { // On diagonal or lower triangle.
				L[idx] = lu[idx]
			}
		}
	}

	work := make([]float64, n*n)
	bi.Dgemm(blas.NoTrans, blas.NoTrans, n, n, n, 1, L, lda, U, lda, 0, work, lda)

	// Apply Permutations P and Q to L*U.
	for i := n - 1; i >= 0; i-- {
		ipv, jpv := ipiv[i], jpiv[i]
		if ipv != i {
			bi.Dswap(n, work[i*lda:], 1, work[ipv*lda:], 1)
		}
		if jpv != i {
			bi.Dswap(n, work[i:], lda, work[jpv:], lda)
		}
	}

	// A should be reconstructed by now.
	for i := range lu {
		if math.Abs(work[i]-a.Data[i]) > tol {
			t.Errorf("%v: matrix %d idx not equal after reconstruction. got %g, expected %g", name, i, work[i], a.Data[i])
		}
	}
}
