// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dgetc2er interface {
	Dgetc2(n int, a []float64, lda int, ipiv, jpiv []int) (k int)
}

func Dgetc2Test(t *testing.T, impl Dgetc2er) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 20} {
		for _, lda := range []int{n, n + 5} {
			dgetc2Test(t, impl, rnd, n, lda, false)
			dgetc2Test(t, impl, rnd, n, lda, true)
		}
	}
}

func dgetc2Test(t *testing.T, impl Dgetc2er, rnd *rand.Rand, n, lda int, perturb bool) {
	const tol = 1e-14

	name := fmt.Sprintf("n=%v,lda=%v,perturb=%v", n, lda, perturb)

	// Generate a random lower-triangular matrix with unit diagonal.
	l := randomGeneral(n, n, max(1, n), rnd)
	for i := 0; i < n; i++ {
		l.Data[i*l.Stride+i] = 1
		for j := i + 1; j < n; j++ {
			l.Data[i*l.Stride+j] = 0
		}
	}

	// Generate a random upper-triangular matrix.
	u := randomGeneral(n, n, max(1, n), rnd)
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			u.Data[i*u.Stride+j] = 0
		}
	}
	if perturb && n > 0 {
		// Make U singular by randomly placing a zero on the diagonal.
		i := rnd.IntN(n)
		u.Data[i*u.Stride+i] = 0
	}

	// Construct A = L*U.
	a := zeros(n, n, max(1, lda))
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, l, u, 0, a)

	// Allocate slices for pivots and pre-fill them with invalid indices.
	ipiv := make([]int, n)
	jpiv := make([]int, n)
	for i := 0; i < n; i++ {
		ipiv[i] = -1
		jpiv[i] = -1
	}

	// Call Dgetc2 to compute the LU decomposition.
	lu := cloneGeneral(a)
	k := impl.Dgetc2(n, lu.Data, lu.Stride, ipiv, jpiv)

	if n == 0 {
		return
	}

	if perturb && k < 0 {
		t.Errorf("%v: expected matrix perturbation", name)
	}

	// Verify all indices have been set.
	for i := 0; i < n; i++ {
		if ipiv[i] < 0 {
			t.Errorf("%v: ipiv[%d] is not set", name, i)
		}
		if jpiv[i] < 0 {
			t.Errorf("%v: jpiv[%d] is not set", name, i)
		}
	}

	// Construct L and U matrices from Dgetc2 output.
	l = zeros(n, n, n)
	u = zeros(n, n, n)
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			l.Data[i*l.Stride+j] = lu.Data[i*lu.Stride+j]
		}
		l.Data[i*l.Stride+i] = 1
		for j := i; j < n; j++ {
			u.Data[i*u.Stride+j] = lu.Data[i*lu.Stride+j]
		}
	}
	diff := zeros(n, n, n)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, l, u, 0, diff)

	// Apply permutation matrices P and Q to L*U.
	for i := n - 1; i >= 0; i-- {
		ipv := ipiv[i]
		if ipv != i {
			row1 := blas64.Vector{N: n, Data: diff.Data[i*diff.Stride:], Inc: 1}
			row2 := blas64.Vector{N: n, Data: diff.Data[ipv*diff.Stride:], Inc: 1}
			blas64.Swap(row1, row2)
		}
		jpv := jpiv[i]
		if jpv != i {
			col1 := blas64.Vector{N: n, Data: diff.Data[i:], Inc: diff.Stride}
			col2 := blas64.Vector{N: n, Data: diff.Data[jpv:], Inc: diff.Stride}
			blas64.Swap(col1, col2)
		}
	}

	// Compute the residual |P*L*U*Q - A| and check that it is small.
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			diff.Data[i*diff.Stride+j] -= a.Data[i*a.Stride+j]
		}
	}
	resid := dlange(lapack.MaxColumnSum, n, n, diff.Data, diff.Stride)
	if resid > tol || math.IsNaN(resid) {
		t.Errorf("%v: unexpected result; resid=%v, want<=%v", name, resid, tol)
	}
}
