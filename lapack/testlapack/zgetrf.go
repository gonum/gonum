// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/cmplx"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/cblas128"
)

type Zgetrfer interface {
	Zgetrf(m, n int, a []complex128, lda int, ipiv []int) bool
}

func ZgetrfTest(t *testing.T, impl Zgetrfer) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, test := range []struct {
		m, n, lda int
	}{
		{1, 1, 0},
		{2, 2, 0},
		{3, 3, 0},
		{5, 5, 0},
		{16, 16, 0},
		{64, 64, 0},
		{10, 5, 0},
		{5, 10, 0},
		{10, 10, 0},
		{3, 500, 0},
		{500, 3, 0},
		{100, 100, 0},
		{100, 200, 0},
		{200, 100, 0},
		{10, 10, 20},
		{5, 10, 20},
		{10, 5, 20},
		{100, 50, 130},
	} {
		m := test.m
		n := test.n
		lda := test.lda
		if lda == 0 {
			lda = max(1, n)
		}
		a := make([]complex128, m*lda)
		for i := range a {
			a[i] = complex(rnd.NormFloat64(), rnd.NormFloat64())
		}
		mn := min(m, n)
		ipiv := make([]int, mn)
		for i := range ipiv {
			ipiv[i] = rnd.Int()
		}

		aCopy := make([]complex128, len(a))
		copy(aCopy, a)
		ok := impl.Zgetrf(m, n, a, lda, ipiv)
		name := fmt.Sprintf("m=%v,n=%v,lda=%v", m, n, lda)
		checkZPLU(t, name, ok, m, n, lda, ipiv, a, aCopy, 1e-10)
	}

	// Test singular matrix: zero column forces zero pivot.
	t.Run("singular", func(t *testing.T) {
		n := 5
		a := make([]complex128, n*n)
		for i := range a {
			a[i] = complex(rnd.NormFloat64(), rnd.NormFloat64())
		}
		// Zero out column 2 to make A singular.
		for i := 0; i < n; i++ {
			a[i*n+2] = 0
		}
		ipiv := make([]int, n)
		if ok := impl.Zgetrf(n, n, a, n, ipiv); ok {
			t.Errorf("Zgetrf returned ok=true on singular matrix")
		}
	})
}

// checkZPLU verifies that the factorized matrix (L in lower triangle, U in
// upper triangle, unit diagonal of L implicit) together with ipiv reproduces
// the original matrix when multiplied as P*L*U.
func checkZPLU(t *testing.T, name string, ok bool, m, n, lda int, ipiv []int, factorized, original []complex128, tol float64) {
	t.Helper()

	var hasZeroDiag bool
	mn := min(m, n)
	for i := 0; i < mn; i++ {
		if factorized[i*lda+i] == 0 {
			hasZeroDiag = true
			break
		}
	}
	if hasZeroDiag && ok {
		t.Errorf("%s: zero diagonal but returned ok", name)
	}
	if !hasZeroDiag && !ok {
		t.Errorf("%s: nonzero diagonal but returned !ok", name)
	}

	// Build L (m×mn) and U (mn×n) from factorized data.
	l := make([]complex128, m*mn)
	ldl := mn
	u := make([]complex128, mn*n)
	ldu := n
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			v := factorized[i*lda+j]
			switch {
			case i == j:
				l[i*ldl+i] = 1
				u[i*ldu+i] = v
			case i > j:
				if j < mn {
					l[i*ldl+j] = v
				}
			case i < j:
				if i < mn {
					u[i*ldu+j] = v
				}
			}
		}
	}

	L := cblas128.General{
		Rows:   m,
		Cols:   mn,
		Stride: ldl,
		Data:   l,
	}
	U := cblas128.General{
		Rows:   mn,
		Cols:   n,
		Stride: ldu,
		Data:   u,
	}
	LU := cblas128.General{
		Rows:   m,
		Cols:   n,
		Stride: n,
		Data:   make([]complex128, m*n),
	}
	cblas128.Gemm(blas.NoTrans, blas.NoTrans, 1, L, U, 0, LU)

	// Build permutation matrix P from ipiv.
	p := make([]complex128, m*m)
	ldp := m
	for i := 0; i < m; i++ {
		p[i*ldp+i] = 1
	}
	for i := len(ipiv) - 1; i >= 0; i-- {
		v := ipiv[i]
		cblas128.Swap(cblas128.Vector{N: m, Inc: 1, Data: p[i*ldp:]},
			cblas128.Vector{N: m, Inc: 1, Data: p[v*ldp:]})
	}
	P := cblas128.General{
		Rows:   m,
		Cols:   m,
		Stride: ldp,
		Data:   p,
	}

	// Compute P*L*U.
	PLU := cblas128.General{
		Rows:   m,
		Cols:   n,
		Stride: n,
		Data:   make([]complex128, m*n),
	}
	cblas128.Gemm(blas.NoTrans, blas.NoTrans, 1, P, LU, 0, PLU)

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			diff := PLU.Data[i*n+j] - original[i*lda+j]
			if cmplx.Abs(diff) > tol {
				t.Errorf("%s: P*L*U mismatch at (%d,%d): got %v want %v (diff=%v)", name, i, j, PLU.Data[i*n+j], original[i*lda+j], cmplx.Abs(diff))
				return
			}
		}
	}
}
