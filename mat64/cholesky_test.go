// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	check "launchpad.net/gocheck"
)

func (s *S) TestCholesky(c *check.C) {
	for _, t := range []struct {
		a   *Dense
		spd bool
	}{
		{
			a: mustDense(NewDense(3, 3, []float64{
				4, 1, 1,
				1, 2, 3,
				1, 3, 6,
			})),

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
