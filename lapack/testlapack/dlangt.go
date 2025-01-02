// Copyright ©2020 The Gonum Authors. All rights reserved.
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

type Dlangter interface {
	Dlangt(norm lapack.MatrixNorm, n int, dl, d, du []float64) float64
}

func DlangtTest(t *testing.T, impl Dlangter) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, norm := range []lapack.MatrixNorm{lapack.MaxAbs, lapack.MaxRowSum, lapack.MaxColumnSum, lapack.Frobenius} {
		t.Run(normToString(norm), func(t *testing.T) {
			for _, n := range []int{0, 1, 2, 3, 4, 5, 10} {
				for iter := 0; iter < 10; iter++ {
					dlangtTest(t, impl, rnd, norm, n)
				}
			}
		})
	}
}

func dlangtTest(t *testing.T, impl Dlangter, rnd *rand.Rand, norm lapack.MatrixNorm, n int) {
	const (
		tol   = 1e-14
		extra = 10
	)

	name := fmt.Sprintf("n=%v", n)

	// Generate three random diagonals.
	dl := randomSlice(n+extra, rnd)
	dlCopy := make([]float64, len(dl))
	copy(dlCopy, dl)

	d := randomSlice(n+1+extra, rnd)
	// Sometimes put a NaN into the matrix.
	if n > 0 && rnd.Float64() < 0.5 {
		d[rnd.IntN(n)] = math.NaN()
	}
	dCopy := make([]float64, len(d))
	copy(dCopy, d)

	du := randomSlice(n+extra, rnd)
	duCopy := make([]float64, len(du))
	copy(duCopy, du)

	// Deal with zero-sized matrices early.
	if n == 0 {
		got := impl.Dlangt(norm, n, nil, nil, nil)
		if got != 0 {
			t.Errorf("%v: unexpected result for zero-sized matrix with nil input", name)
		}
		got = impl.Dlangt(norm, n, dl, d, du)
		if !floats.Same(dl, dlCopy) {
			t.Errorf("%v: unexpected modification in dl", name)
		}
		if !floats.Same(d, dCopy) {
			t.Errorf("%v: unexpected modification in d", name)
		}
		if !floats.Same(du, duCopy) {
			t.Errorf("%v: unexpected modification in du", name)
		}
		if got != 0 {
			t.Errorf("%v: unexpected result for zero-sized matrix with non-nil input", name)
		}
		return
	}

	// Generate a dense representation of the matrix and compute the wanted result.
	a := zeros(n, n, n)
	for i := 0; i < n-1; i++ {
		a.Data[i*a.Stride+i] = d[i]
		a.Data[i*a.Stride+i+1] = du[i]
		a.Data[(i+1)*a.Stride+i] = dl[i]
	}
	a.Data[(n-1)*a.Stride+n-1] = d[n-1]

	got := impl.Dlangt(norm, n, dl, d, du)

	if !floats.Same(dl, dlCopy) {
		t.Errorf("%v: unexpected modification in dl", name)
	}
	if !floats.Same(d, dCopy) {
		t.Errorf("%v: unexpected modification in d", name)
	}
	if !floats.Same(du, duCopy) {
		t.Errorf("%v: unexpected modification in du", name)
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
