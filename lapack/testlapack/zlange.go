// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"math/cmplx"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/lapack"
)

type Zlanger interface {
	Zlange(norm lapack.MatrixNorm, m, n int, a []complex128, lda int, work []float64) float64
}

func ZlangeTest(t *testing.T, impl Zlanger) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, test := range []struct {
		m, n, lda int
	}{
		{4, 3, 0},
		{3, 4, 0},
		{4, 3, 20},
		{3, 4, 20},
		{1, 1, 0},
		{5, 5, 0},
	} {
		m := test.m
		n := test.n
		lda := test.lda
		if lda == 0 {
			lda = n
		}
		a := make([]complex128, m*lda)
		for i := range a {
			a[i] = complex(rnd.Float64()-0.5, rnd.Float64()-0.5)
		}
		work := make([]float64, n)

		// MaxAbs.
		got := impl.Zlange(lapack.MaxAbs, m, n, a, lda, work)
		var want float64
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				want = math.Max(want, cmplx.Abs(a[i*lda+j]))
			}
		}
		if got != want {
			t.Errorf("MaxAbs mismatch. m=%d n=%d Want %v, got %v.", m, n, want, got)
		}

		// MaxColumnSum.
		got = impl.Zlange(lapack.MaxColumnSum, m, n, a, lda, work)
		want = 0
		for j := 0; j < n; j++ {
			var sum float64
			for i := 0; i < m; i++ {
				sum += cmplx.Abs(a[i*lda+j])
			}
			want = math.Max(want, sum)
		}
		if math.Abs(got-want) > 1e-13 {
			t.Errorf("MaxColumnSum mismatch. m=%d n=%d Want %v, got %v.", m, n, want, got)
		}

		// MaxRowSum.
		got = impl.Zlange(lapack.MaxRowSum, m, n, a, lda, work)
		want = 0
		for i := 0; i < m; i++ {
			var sum float64
			for j := 0; j < n; j++ {
				sum += cmplx.Abs(a[i*lda+j])
			}
			want = math.Max(want, sum)
		}
		if math.Abs(got-want) > 1e-13 {
			t.Errorf("MaxRowSum mismatch. m=%d n=%d Want %v, got %v.", m, n, want, got)
		}

		// Frobenius.
		got = impl.Zlange(lapack.Frobenius, m, n, a, lda, work)
		want = 0
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				abs := cmplx.Abs(a[i*lda+j])
				want += abs * abs
			}
		}
		want = math.Sqrt(want)
		if math.Abs(got-want) > 1e-13 {
			t.Errorf("Frobenius mismatch. m=%d n=%d Want %v, got %v.", m, n, want, got)
		}

		// Corner: zero matrix.
		zero := make([]complex128, m*lda)
		for _, norm := range []lapack.MatrixNorm{lapack.MaxAbs, lapack.MaxRowSum, lapack.MaxColumnSum, lapack.Frobenius} {
			got := impl.Zlange(norm, m, n, zero, lda, work)
			if got != 0 {
				t.Errorf("zero-matrix norm %c gave %v, want 0", byte(norm), got)
			}
		}
	}

	// Empty matrix.
	for _, norm := range []lapack.MatrixNorm{lapack.MaxAbs, lapack.MaxRowSum, lapack.MaxColumnSum, lapack.Frobenius} {
		got := impl.Zlange(norm, 0, 0, nil, 1, nil)
		if got != 0 {
			t.Errorf("empty-matrix norm %c gave %v, want 0", byte(norm), got)
		}
	}
}
