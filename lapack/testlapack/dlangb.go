// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/lapack"
)

type Dlangber interface {
	Dlangb(norm lapack.MatrixNorm, m, n, kl, ku int, ab []float64, ldab int) float64
}

func DlangbTest(t *testing.T, impl Dlangber) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, norm := range []lapack.MatrixNorm{lapack.MaxAbs, lapack.MaxRowSum, lapack.MaxColumnSum, lapack.Frobenius} {
		t.Run(normToString(norm), func(t *testing.T) {
			for _, n := range []int{0, 1, 2, 3, 4, 5, 10} {
				for _, m := range []int{0, 1, 2, 3, 4, 5, 10} {
					for _, kl := range []int{0, 1, 2, 3, 4, 5, 10} {
						for _, ku := range []int{0, 1, 2, 3, 4, 5, 10} {
							for _, ldab := range []int{kl + ku + 1, kl + ku + 1 + 7} {
								dlangbTest(t, impl, rnd, norm, m, n, kl, ku, ldab)
							}
						}
					}
				}
			}
		})
	}
}

func dlangbTest(t *testing.T, impl Dlangber, rnd *rand.Rand, norm lapack.MatrixNorm, m, n, kl, ku, ldab int) {
	const tol = 1e-14

	name := fmt.Sprintf("m=%v,n=%v,kl=%v,ku=%v,ldab=%v", m, n, kl, ku, ldab)

	// Generate a random band matrix.
	ab := randomSlice(m*ldab, rnd)
	// Sometimes put a NaN into the matrix.
	if m > 0 && n > 0 && rnd.Float64() < 0.5 {
		i := rnd.IntN(m)
		ab[i*ldab+kl] = math.NaN()
	}
	abCopy := make([]float64, len(ab))
	copy(abCopy, ab)

	// Deal with zero-sized matrices early.
	if m == 0 || n == 0 {
		got := impl.Dlangb(norm, m, n, kl, ku, nil, ldab)
		if got != 0 {
			t.Errorf("%v: unexpected result for zero-sized matrix with nil input", name)
		}
		got = impl.Dlangb(norm, m, n, kl, ku, ab, ldab)
		if !floats.Same(ab, abCopy) {
			t.Errorf("%v: unexpected modification in dl", name)
		}
		if got != 0 {
			t.Errorf("%v: unexpected result for zero-sized matrix with non-nil input", name)
		}
		return
	}

	got := impl.Dlangb(norm, m, n, kl, ku, ab, ldab)

	if !floats.Same(ab, abCopy) {
		t.Errorf("%v: unexpected modification in ab", name)
	}

	// Generate a dense representation of the matrix and compute the wanted result.
	a := zeros(m, n, n)
	for i := 0; i < m; i++ {
		for j := max(0, i-kl); j < min(i+ku+1, n); j++ {
			a.Data[i*a.Stride+j] = ab[i*ldab+j-i+kl]
		}
	}
	want := dlange(norm, a.Rows, a.Cols, a.Data, a.Stride)

	if math.IsNaN(want) {
		if !math.IsNaN(got) {
			t.Errorf("%v: unexpected result with NaN element; got %v, want %v", name, got, want)
		}
		return
	}

	if math.IsNaN(got) {
		t.Errorf("%v: unexpected NaN; want %v", name, want)
		return
	}

	if norm == lapack.MaxAbs {
		if got != want {
			t.Errorf("%v: unexpected result; got %v, want %v", name, got, want)
		}
		return
	}
	diff := math.Abs(got - want)
	if diff > tol {
		t.Errorf("%v: unexpected result; got %v, want %v, diff=%v", name, got, want, diff)
	}
}
