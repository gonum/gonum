// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	check "launchpad.net/gocheck"
	"math"
)

func (s *S) TestSVD(c *check.C) {
	for _, t := range []struct {
		a *Dense

		epsilon float64
		small   float64

		wantu bool
		u     *Dense

		sigma []float64

		wantv bool
		v     *Dense
	}{
		{
			a: mustDense(NewDense(4, 2, []float64{2, 4, 1, 3, 0, 0, 0, 0})),

			epsilon: math.Pow(2, -52.0),
			small:   math.Pow(2, -966.0),

			wantu: true,
			u: mustDense(NewDense(4, 2, []float64{
				0.8174155604703632, -0.5760484367663209,
				0.5760484367663209, 0.8174155604703633,
				0, 0,
				0, 0,
			})),

			sigma: []float64{5.464985704219041, 0.365966190626258},

			wantv: true,
			v: mustDense(NewDense(2, 2, []float64{
				0.4045535848337571, -0.9145142956773044,
				0.9145142956773044, 0.4045535848337571,
			})),
		},
		{
			a: mustDense(NewDense(4, 2, []float64{2, 4, 1, 3, 0, 0, 0, 0})),

			epsilon: math.Pow(2, -52.0),
			small:   math.Pow(2, -966.0),

			wantu: true,
			u: mustDense(NewDense(4, 2, []float64{
				0.8174155604703632, -0.5760484367663209,
				0.5760484367663209, 0.8174155604703633,
				0, 0,
				0, 0,
			})),

			sigma: []float64{5.464985704219041, 0.365966190626258},

			wantv: false,
		},
		{
			a: mustDense(NewDense(4, 2, []float64{2, 4, 1, 3, 0, 0, 0, 0})),

			epsilon: math.Pow(2, -52.0),
			small:   math.Pow(2, -966.0),

			wantu: false,

			sigma: []float64{5.464985704219041, 0.365966190626258},

			wantv: true,
			v: mustDense(NewDense(2, 2, []float64{
				0.4045535848337571, -0.9145142956773044,
				0.9145142956773044, 0.4045535848337571,
			})),
		},
		{
			a: mustDense(NewDense(4, 2, []float64{2, 4, 1, 3, 0, 0, 0, 0})),

			epsilon: math.Pow(2, -52.0),
			small:   math.Pow(2, -966.0),

			sigma: []float64{5.464985704219041, 0.365966190626258},
		},
		{
			// FIXME(kortschak)
			// This test will fail if t.sigma is set to the real expected values
			// or if u and v are requested, due to a bug in the original Jama code
			// forcing a to be a tall or square matrix.
			//
			// This is a failing case to use to fix that bug.
			a: mustDense(NewDense(3, 11, []float64{
				1, 1, 0, 1, 0, 0, 0, 0, 0, 11, 1,
				1, 0, 0, 0, 0, 0, 1, 0, 0, 12, 2,
				1, 1, 0, 0, 0, 0, 0, 0, 1, 13, 3,
			})),

			epsilon: math.Pow(2, -52.0),
			small:   math.Pow(2, -966.0),

			// FIXME(kortschak) sigma is one element longer than it should be.
			sigma: []float64{21.25950088109745, 1.5415021616856577, 1.2873979074613637, 0},
		},
	} {
		svd := SVD(DenseCopyOf(t.a), t.epsilon, t.small, t.wantu, t.wantv)
		if t.sigma != nil {
			c.Check(svd.Sigma, check.DeepEquals, t.sigma)
		}
		s := svd.S()

		if svd.U != nil {
			c.Check(svd.U.Equals(t.u), check.Equals, true)
		} else {
			c.Check(t.wantu, check.Equals, false)
			c.Check(t.u, check.IsNil)
		}
		if svd.V != nil {
			c.Check(svd.V.Equals(t.v), check.Equals, true)
		} else {
			c.Check(t.wantv, check.Equals, false)
			c.Check(t.v, check.IsNil)
		}

		if t.wantu && t.wantv {
			c.Assert(svd.U, check.NotNil)
			c.Assert(svd.V, check.NotNil)
			vt := &Dense{}
			vt.TCopy(svd.V)
			svd.U.Mul(svd.U, s)
			svd.U.Mul(svd.U, vt)
			c.Check(svd.U.EqualsApprox(t.a, 1e-12), check.Equals, true)
		}
	}
}
