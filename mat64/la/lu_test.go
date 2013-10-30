// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"github.com/gonum/matrix/mat64"

	check "launchpad.net/gocheck"
)

func (s *S) TestLUD(c *check.C) {
	for _, t := range []struct {
		a *mat64.Dense

		l *mat64.Dense
		u *mat64.Dense

		pivot []int
		sign  int
	}{
		{ // This is a hard coded equivalent of the approach used in the Jama LU test.
			a: mustDense(mat64.NewDense(3, 3, []float64{
				0, 2, 3,
				4, 5, 6,
				7, 8, 9,
			})),

			l: mustDense(mat64.NewDense(3, 3, []float64{
				1, 0, 0,
				0, 1, 0,
				0.5714285714285714, 0.2142857142857144, 1,
			})),
			u: mustDense(mat64.NewDense(3, 3, []float64{
				7, 8, 9,
				0, 2, 3,
				0, 0, 0.2142857142857144,
			})),
			pivot: []int{
				2, // 0 0 1
				0, // 1 0 0
				1, // 0 1 0
			},
			sign: 1,
		},
	} {
		a := &mat64.Dense{}
		a.Clone(t.a)

		lu, pivot, pivotsign := LUD(a)
		if t.pivot != nil {
			c.Check(pivot, check.DeepEquals, t.pivot)
			c.Check(pivotsign, check.Equals, t.sign)
		}

		l := LUGetL(lu)
		if t.l != nil {
			c.Check(l.Equals(t.l), check.Equals, true)
		}
		u := LUGetU(lu)
		if t.u != nil {
			c.Check(u.Equals(t.u), check.Equals, true)
		}

		l.Mul(l, u)
		pa := &mat64.Dense{}
		pa.Clone(t.a)
		c.Check(l.EqualsApprox(pivotRows(pa, pivot), 1e-12), check.Equals, true)

		x := LUSolve(lu, eye(), pivot)
		t.a.Mul(t.a, x)
		c.Check(t.a.EqualsApprox(eye(), 1e-12), check.Equals, true)
	}
}

func (s *S) TestLUDGaussian(c *check.C) {
	for _, t := range []struct {
		a *mat64.Dense

		l *mat64.Dense
		u *mat64.Dense

		pivot []int
		sign  int
	}{
		{ // This is a hard coded equivalent of the approach used in the Jama LU test.
			a: mustDense(mat64.NewDense(3, 3, []float64{
				0, 2, 3,
				4, 5, 6,
				7, 8, 9,
			})),

			l: mustDense(mat64.NewDense(3, 3, []float64{
				1, 0, 0,
				0, 1, 0,
				0.5714285714285714, 0.2142857142857144, 1,
			})),
			u: mustDense(mat64.NewDense(3, 3, []float64{
				7, 8, 9,
				0, 2, 3,
				0, 0, 0.2142857142857144,
			})),
			pivot: []int{
				2, // 0 0 1
				0, // 1 0 0
				1, // 0 1 0
			},
			sign: 1,
		},
	} {
		a := &mat64.Dense{}
		a.Clone(t.a)

		lu, pivot, pivotsign := LUDGaussian(a)
		if t.pivot != nil {
			c.Check(pivot, check.DeepEquals, t.pivot)
			c.Check(pivotsign, check.Equals, t.sign)
		}

		l := LUGetL(lu)
		if t.l != nil {
			c.Check(l.Equals(t.l), check.Equals, true)
		}
		u := LUGetU(lu)
		if t.u != nil {
			c.Check(u.Equals(t.u), check.Equals, true)
		}

		l.Mul(l, u)
		pa := &mat64.Dense{}
		pa.Clone(t.a)
		c.Check(l.EqualsApprox(pivotRows(pa, pivot), 1e-12), check.Equals, true)

		x := LUSolve(lu, eye(), pivot)
		t.a.Mul(t.a, x)
		c.Check(t.a.EqualsApprox(eye(), 1e-12), check.Equals, true)
	}
}
