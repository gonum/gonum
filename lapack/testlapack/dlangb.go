// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/lapack"
)

type Dlangber interface {
	Dlangb(norm lapack.MatrixNorm, n, kl, ku int, ab []float64, ldab int, work []float64) float64
}

func DlangbTest(t *testing.T, impl Dlangber) {
	rnd := rand.New(rand.NewSource(1))
	for _, norm := range []lapack.MatrixNorm{lapack.MaxAbs, lapack.MaxRowSum, lapack.MaxColumnSum, lapack.Frobenius} {
		t.Run(normToString(norm), func(t *testing.T) {
			for _, n := range []int{0, 1, 2, 3, 4, 5, 10} {
				for _, kl := range []int{0, 1, 2, 3, 4, 5, 10} {
					for _, ku := range []int{0, 1, 2, 3, 4, 5, 10} {
						for _, ldab := range []int{kl + ku + 1, kl + ku + 1 + 7} {
							for iter := 0; iter < 10; iter++ {
								dlangbTest(t, impl, rnd, norm, n, kl, ku, ldab)
							}
						}
					}
				}
			}
		})
	}
}

func dlangbTest(t *testing.T, impl Dlangber, rnd *rand.Rand, norm lapack.MatrixNorm, n, kl, ku, ldab int) {
	const tol = 1e-14

	name := fmt.Sprintf("n=%v,kl=%v,ku=%v,ldab=%v", n, kl, ku, ldab)

	// Generate a random band matrix.
	ab := randomSlice(n*ldab, rnd)
	// Sometimes put a NaN into the matrix.
	if n > 0 && rnd.Float64() < 0.5 {
		i := rnd.Intn(n)
		ab[i*ldab+kl] = math.NaN()
	}
	abCopy := make([]float64, len(ab))
	copy(abCopy, ab)

	// Deal with zero-sized matrices early.
	if n == 0 {
		got := impl.Dlangb(norm, n, kl, ku, nil, ldab, nil)
		if got != 0 {
			t.Errorf("%v: unexpected result for zero-sized matrix with nil input", name)
		}
		got = impl.Dlangb(norm, n, kl, ku, ab, ldab, nil)
		if !floats.Same(ab, abCopy) {
			t.Errorf("%v: unexpected modification in dl", name)
		}
		if got != 0 {
			t.Errorf("%v: unexpected result for zero-sized matrix with non-nil input", name)
		}
		return
	}

	// Generate a dense representation of the matrix and compute the wanted result.
	a := zeros(n, n, n)
	for i := 0; i < n; i++ {
		for j := max(0, i-kl); j < min(i+ku+1, n); j++ {
			a.Data[i*a.Stride+j] = ab[i*ldab+j-i+kl]
		}
	}

	var work []float64
	if norm == lapack.MaxColumnSum {
		work = make([]float64, n)
	}
	got := impl.Dlangb(norm, n, kl, ku, ab, ldab, work)

	if !floats.Same(ab, abCopy) {
		t.Errorf("%v: unexpected modification in ab", name)
	}

	want := dlange(norm, n, n, a.Data, a.Stride)

	if math.IsNaN(want) {
		if !math.IsNaN(got) {
			t.Errorf("%v: unexpected result with NaN element; got %v, want %v", name, got, want)
		}
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
