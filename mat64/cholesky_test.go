// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import "gopkg.in/check.v1"

func (s *S) TestCholesky(c *check.C) {
	for _, t := range []struct {
		a     *SymDense
		upper bool
		f     *TriDense

		want *TriDense
		pd   bool
	}{
		{
			a: NewSymDense(3, []float64{
				4, 1, 1,
				0, 2, 3,
				0, 0, 6,
			}),
			upper: false,
			f:     &TriDense{},

			want: NewTriDense(3, false, []float64{
				2, 0, 0,
				0.5, 1.3228756555322954, 0,
				0.5, 2.0788046015507495, 1.195228609334394,
			}),
			pd: true,
		},
		{
			a: NewSymDense(3, []float64{
				4, 1, 1,
				0, 2, 3,
				0, 0, 6,
			}),
			upper: true,
			f:     &TriDense{},

			want: NewTriDense(3, true, []float64{
				2, 0.5, 0.5,
				0, 1.3228756555322954, 2.0788046015507495,
				0, 0, 1.195228609334394,
			}),
			pd: true,
		},
		{
			a: NewSymDense(3, []float64{
				4, 1, 1,
				0, 2, 3,
				0, 0, 6,
			}),
			upper: false,
			f:     NewTriDense(3, false, nil),

			want: NewTriDense(3, false, []float64{
				2, 0, 0,
				0.5, 1.3228756555322954, 0,
				0.5, 2.0788046015507495, 1.195228609334394,
			}),
			pd: true,
		},
		{
			a: NewSymDense(3, []float64{
				4, 1, 1,
				0, 2, 3,
				0, 0, 6,
			}),
			upper: true,
			f:     NewTriDense(3, false, nil),

			want: NewTriDense(3, true, []float64{
				2, 0.5, 0.5,
				0, 1.3228756555322954, 2.0788046015507495,
				0, 0, 1.195228609334394,
			}),
			pd: true,
		},
	} {
		ok := t.f.Cholesky(t.a, t.upper)
		c.Check(ok, check.Equals, t.pd)
		fc := DenseCopyOf(t.f)
		c.Check(fc.Equals(t.want), check.Equals, true)

		ft := &Dense{}
		ft.TCopy(t.f)

		if t.upper {
			fc.Mul(ft, fc)
		} else {
			fc.Mul(fc, ft)
		}
		c.Check(fc.EqualsApprox(t.a, 1e-12), check.Equals, true)

		var x Dense
		x.SolveCholesky(t.f, eye())

		var res Dense
		res.Mul(t.a, &x)
		c.Check(res.EqualsApprox(eye(), 1e-12), check.Equals, true)

		x = Dense{}
		x.SolveTri(t.f, t.upper, eye())
		x.SolveTri(t.f, !t.upper, &x)

		res.Mul(t.a, &x)
		c.Check(res.EqualsApprox(eye(), 1e-12), check.Equals, true)
	}
}

func (s *S) TestCholeskySolve(c *check.C) {
	for _, t := range []struct {
		a   *SymDense
		b   *Dense
		ans *Dense
	}{
		{
			a: NewSymDense(2, []float64{
				1, 0,
				0, 1,
			}),
			b:   NewDense(2, 1, []float64{5, 6}),
			ans: NewDense(2, 1, []float64{5, 6}),
		},
	} {
		var f TriDense
		ok := f.Cholesky(t.a, false)
		c.Assert(ok, check.Equals, true)

		var x Dense
		x.SolveCholesky(&f, t.b)
		c.Check(x.EqualsApprox(t.ans, 1e-12), check.Equals, true)

		x = Dense{}
		x.SolveTri(&f, false, t.b)
		x.SolveTri(&f, true, &x)
		c.Check(x.EqualsApprox(t.ans, 1e-12), check.Equals, true)
	}
}
