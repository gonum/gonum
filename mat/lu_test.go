// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"math/rand/v2"
	"testing"
)

func TestLU(t *testing.T) {
	t.Parallel()
	const tol = 1e-16
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, n := range []int{1, 2, 3, 4, 5, 10, 11, 50} {
		// Construct a random matrix A.
		a := NewDense(n, n, nil)
		a.Apply(func(_, _ int, _ float64) float64 { return rnd.NormFloat64() }, a)

		// Compute the LU factorization of A.
		var lu LU
		lu.Factorize(a)

		// Compare A and LU using At.
		if !EqualApprox(a, &lu, tol*float64(n)) {
			var diff Dense
			diff.Sub(a, &lu)
			t.Errorf("n=%d: A and LU not equal\ndiff=%v", n, Formatted(&diff, Prefix("     ")))
		}

		// Recover A using RowPivots, LTo and UTo.
		var l, u TriDense
		lu.LTo(&l)
		lu.UTo(&u)
		var got Dense
		got.Mul(&l, &u)
		got.PermuteRows(lu.RowPivots(nil), false)
		if !EqualApprox(&got, a, tol*float64(n)) {
			var diff Dense
			diff.Sub(&got, a)
			t.Errorf("n=%d: A and P*L*U not equal\ndiff=%v", n, Formatted(&diff, Prefix("     ")))
		}
	}
}

func TestLURankOne(t *testing.T) {
	t.Parallel()
	const tol = 1e-14
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, n := range []int{1, 2, 3, 4, 5, 10, 50} {
		// Construct a random matrix A.
		a := NewDense(n, n, nil)
		a.Apply(func(_, _ int, _ float64) float64 { return rnd.NormFloat64() }, a)

		// Compute the LU factorization of A.
		var lu LU
		lu.Factorize(a)

		// Apply a rank one update to A. Ensure the update magnitude is larger than
		// the equal tolerance.
		alpha := rnd.Float64() + 1
		x := NewVecDense(n, nil)
		y := NewVecDense(n, nil)
		for i := 0; i < n; i++ {
			x.setVec(i, rnd.Float64()+1)
			y.setVec(i, rnd.Float64()+1)
		}
		a.RankOne(a, alpha, x, y)

		// Apply the same rank one update to the LU factorization of A.
		var luNew LU
		luNew.RankOne(&lu, alpha, x, y)
		lu.RankOne(&lu, alpha, x, y)

		if !EqualApprox(&lu, a, tol*float64(n)) {
			var diff Dense
			diff.Sub(&lu, a)
			t.Errorf("n=%d: rank one mismatch\ndiff=%v", n, Formatted(&diff, Prefix("     ")))
		}

		if !Equal(&lu, &luNew) {
			var diff Dense
			diff.Sub(&lu, &luNew)
			t.Errorf("n=%d: rank one mismatch with new receiver\ndiff=%v", n, Formatted(&diff, Prefix("     ")))
		}
	}
}

func TestLUSolveTo(t *testing.T) {
	t.Parallel()
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, test := range []struct {
		n, bc int
	}{
		{5, 5},
		{5, 10},
		{10, 5},
	} {
		n := test.n
		bc := test.bc
		a := NewDense(n, n, nil)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				a.Set(i, j, rnd.NormFloat64())
			}
		}
		b := NewDense(n, bc, nil)
		for i := 0; i < n; i++ {
			for j := 0; j < bc; j++ {
				b.Set(i, j, rnd.NormFloat64())
			}
		}
		var lu LU
		lu.Factorize(a)
		var x Dense
		if err := lu.SolveTo(&x, false, b); err != nil {
			continue
		}
		var got Dense
		got.Mul(a, &x)
		if !EqualApprox(&got, b, 1e-12) {
			t.Errorf("SolveTo mismatch for non-singular matrix. n = %v, bc = %v.\nWant: %v\nGot: %v", n, bc, b, got)
		}
	}
	// TODO(btracey): Add testOneInput test when such a function exists.
}

func TestLUSolveToCond(t *testing.T) {
	t.Parallel()
	for _, test := range []*Dense{
		NewDense(2, 2, []float64{1, 0, 0, 1e-20}),
	} {
		m, _ := test.Dims()
		var lu LU
		lu.Factorize(test)
		b := NewDense(m, 2, nil)
		var x Dense
		if err := lu.SolveTo(&x, false, b); err == nil {
			t.Error("No error for near-singular matrix in matrix solve.")
		}

		bvec := NewVecDense(m, nil)
		var xvec VecDense
		if err := lu.SolveVecTo(&xvec, false, bvec); err == nil {
			t.Error("No error for near-singular matrix in matrix solve.")
		}
	}
}

func TestLUSolveVecTo(t *testing.T) {
	t.Parallel()
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, n := range []int{5, 10} {
		a := NewDense(n, n, nil)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				a.Set(i, j, rnd.NormFloat64())
			}
		}
		b := NewVecDense(n, nil)
		for i := 0; i < n; i++ {
			b.SetVec(i, rnd.NormFloat64())
		}
		var lu LU
		lu.Factorize(a)
		var x VecDense
		if err := lu.SolveVecTo(&x, false, b); err != nil {
			continue
		}
		var got VecDense
		got.MulVec(a, &x)
		if !EqualApprox(&got, b, 1e-12) {
			t.Errorf("SolveTo mismatch n = %v.\nWant: %v\nGot: %v", n, b, got)
		}
	}
	// TODO(btracey): Add testOneInput test when such a function exists.
}
