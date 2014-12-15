// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"gopkg.in/check.v1"
)

func (s *S) TestCholesky(c *check.C) {
	for _, t := range []struct {
		a   *Dense
		spd bool
	}{
		{
			a: NewDense(3, 3, []float64{
				4, 1, 1,
				1, 2, 3,
				1, 3, 6,
			}),

			spd: true,
		},
	} {
		cf := Cholesky(t.a)
		c.Check(cf.SPD, check.Equals, t.spd)

		lt := &Dense{}
		lt.TCopy(cf.L)
		lc := DenseCopyOf(cf.L)

		lc.Mul(lc, lt)
		c.Check(lc.EqualsApprox(t.a, 1e-12), check.Equals, true)

		x := cf.Solve(eye())

		t.a.Mul(t.a, x)
		c.Check(t.a.EqualsApprox(eye(), 1e-12), check.Equals, true)
	}
}

func (s *S) TestCholeskySolve(c *check.C) {
	for _, t := range []struct {
		a   *Dense
		b   *Dense
		ans *Dense
	}{
		{
			a: NewDense(2, 2, []float64{
				1, 0,
				0, 1,
			}),
			b:   NewDense(2, 1, []float64{5, 6}),
			ans: NewDense(2, 1, []float64{5, 6}),
		},
	} {
		ans := Cholesky(t.a).Solve(t.b)
		c.Check(ans.EqualsApprox(t.ans, 1e-12), check.Equals, true)
	}
}
