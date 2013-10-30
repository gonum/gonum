// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"github.com/gonum/matrix/mat64"

	check "launchpad.net/gocheck"
)

func eye() *mat64.Dense {
	return mustDense(mat64.NewDense(3, 3, []float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}))
}

func (s *S) TestCholesky(c *check.C) {
	for _, t := range []struct {
		a   *mat64.Dense
		spd bool
	}{
		{
			a: mustDense(mat64.NewDense(3, 3, []float64{
				4, 1, 1,
				1, 2, 3,
				1, 3, 6,
			})),

			spd: true,
		},
	} {
		l, spd := CholeskyL(t.a)
		c.Check(spd, check.Equals, t.spd)

		lt := &mat64.Dense{}
		lt.TCopy(l)
		lc := &mat64.Dense{}
		lc.Clone(l)

		lc.Mul(lc, lt)
		c.Check(lc.EqualsApprox(t.a, 1e-12), check.Equals, true)

		x := CholeskySolve(l, eye())

		t.a.Mul(t.a, x)
		c.Check(t.a.EqualsApprox(eye(), 1e-12), check.Equals, true)
	}
}
