// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"github.com/gonum/matrix/mat64"

	check "launchpad.net/gocheck"
	"math"
)

func (s *S) TestEigen(c *check.C) {
	for _, t := range []struct {
		a *mat64.Dense

		epsilon float64

		e, d []float64
		v    *mat64.Dense
	}{
		{
			a: mustDense(mat64.NewDense(3, 3, []float64{
				1, 2, 1,
				6, -1, 0,
				-1, -2, -1,
			})),

			epsilon: math.Pow(2, -52.0),

			d: []float64{3.0000000000000044, -4.000000000000003, -1.0980273383714707e-16},
			e: []float64{0, 0, 0},
			v: mustDense(mat64.NewDense(3, 3, []float64{
				-0.48507125007266627, 0.4164965639175221, 0.11785113019775806,
				-0.7276068751089995, -0.8329931278350428, 0.7071067811865481,
				0.48507125007266627, -0.41649656391752166, -1.5320646925708528,
			})),
		},
		{
			a: mustDense(mat64.NewDense(3, 3, []float64{
				1, 6, -1,
				6, -1, -2,
				-1, -2, -1,
			})),

			epsilon: math.Pow(2, -52.0),

			d: []float64{-6.240753470718579, -1.3995889142010132, 6.640342384919599},
			e: []float64{0, 0, 0},
			v: mustDense(mat64.NewDense(3, 3, []float64{
				-0.6134279348516111, -0.31411097261113, -0.7245967607083111,
				0.7697297716508223, -0.03251534945303795, -0.6375412384185983,
				0.17669818159240022, -0.9488293044247931, 0.2617263908869383,
			})),
		},
	} {
		a := &mat64.Dense{}
		a.Clone(t.a)

		d, e, v := Eigen(t.a, t.epsilon)
		if t.d != nil {
			c.Check(d, check.DeepEquals, t.d)
		}
		if t.e != nil {
			c.Check(e, check.DeepEquals, t.e)
		}
		// dm := BuildD(d, e)

		if t.v != nil {
			c.Check(v.Equals(t.v), check.Equals, true)
		}

		// vt := &mat64.Dense{}
		// vt.TCopy(v)
		// v.Mul(v, dm)
		// v.Mul(v, vt)
		// c.Check(v.EqualsApprox(a, 1e-6), check.Equals, true)
		// fmt.Println(v)
	}
}
