// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"math/rand"

	"gopkg.in/check.v1"
)

func (s *S) TestLUD(c *check.C) {
	for _, n := range []int{1, 5, 10, 11, 50} {
		a := NewDense(n, n, nil)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				a.Set(i, j, rand.NormFloat64())
			}
		}
		var want Dense
		want.Clone(a)

		lu := &LU{}
		lu.Factorize(a)

		var l, u TriDense
		l.LFrom(lu)
		u.UFrom(lu)
		var p Dense
		pivot := lu.Pivot(nil)
		p.Permutation(n, pivot)
		var got Dense
		got.Mul(&p, &l)
		got.Mul(&got, &u)
		if !got.EqualsApprox(&want, 1e-12) {
			c.Errorf("PLU does not equal original matrix.\nWant: %v\n Got: %v", want, got)
		}
	}
}

func (s *S) TestSolveLU(c *check.C) {
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
				a.Set(i, j, rand.NormFloat64())
			}
		}
		b := NewDense(n, bc, nil)
		for i := 0; i < n; i++ {
			for j := 0; j < bc; j++ {
				b.Set(i, j, rand.NormFloat64())
			}
		}
		var lu LU
		lu.Factorize(a)
		var x Dense
		if err := x.SolveLU(&lu, false, b); err != nil {
			continue
		}
		var got Dense
		got.Mul(a, &x)
		if !got.EqualsApprox(b, 1e-12) {
			c.Error("Solve mismatch for non-singular matrix. n = %v, bc = %v.\nWant: %v\nGot: %v", n, bc, b, got)
		}
	}
	// TODO(btracey): Add testOneInput test when such a function exists.
}

func (s *S) TestSolveLUVec(c *check.C) {
	for _, n := range []int{5, 10} {
		a := NewDense(n, n, nil)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				a.Set(i, j, rand.NormFloat64())
			}
		}
		b := NewVector(n, nil)
		for i := 0; i < n; i++ {
			b.SetVec(i, rand.NormFloat64())
		}
		var lu LU
		lu.Factorize(a)
		var x Vector
		if err := x.SolveLUVec(&lu, false, b); err != nil {
			continue
		}
		var got Vector
		got.MulVec(a, &x)
		if !got.EqualsApproxVec(b, 1e-12) {
			c.Error("Solve mismatch n = %v.\nWant: %v\nGot: %v", n, b, got)
		}
	}
	// TODO(btracey): Add testOneInput test when such a function exists.
}
