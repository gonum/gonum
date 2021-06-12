// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"
)

type Dlartger interface {
	Dlartg(f, g float64) (cs, sn, r float64)
}

func DlartgTest(t *testing.T, impl Dlartger) {
	const tol = 20 * ulp

	values := []float64{
		-safmax,
		-1 / ulp,
		-1,
		-1.0 / 3,
		-ulp,
		-safmin,
		0,
		safmin,
		ulp,
		1.0 / 3,
		1,
		1 / ulp,
		safmax,
		math.Inf(-1),
		math.Inf(1),
		math.NaN(),
	}

	for _, f := range values {
		for _, g := range values {
			name := fmt.Sprintf("Case f=%v,g=%v", f, g)

			// Generate a plane rotation so that
			//  [ cs sn] * [f] = [r]
			//  [-sn cs]   [g] = [0]
			// where cs*cs + sn*sn = 1.
			cs, sn, r := impl.Dlartg(f, g)

			switch {
			case math.IsNaN(f) || math.IsNaN(g):
				if !math.IsNaN(r) {
					t.Errorf("%v: unexpected r=%v; want NaN", name, r)
				}
			case math.IsInf(f, 0) || math.IsInf(g, 0):
				if !math.IsNaN(r) && !math.IsInf(r, 0) {
					t.Errorf("%v: unexpected r=%v; want NaN or Inf", name, r)
				}
			default:
				d := math.Max(math.Abs(f), math.Abs(g))
				d = math.Min(math.Max(safmin, d), safmax)
				fs := f / d
				gs := g / d
				rs := r / d

				// Check that cs*f + sn*g = r.
				rnorm := math.Abs(rs)
				if rnorm == 0 {
					rnorm = math.Max(math.Abs(fs), math.Abs(gs))
					if rnorm == 0 {
						rnorm = 1
					}
				}
				resid := math.Abs(rs-(cs*fs+sn*gs)) / rnorm
				if resid > tol {
					t.Errorf("%v: cs*f + sn*g != r; resid=%v", name, resid)
				}

				// Check that -sn*f + cs*g = 0.
				resid = math.Abs(-sn*fs + cs*gs)
				if resid > tol {
					t.Errorf("%v: -sn*f + cs*g != 0; resid=%v", name, resid)
				}

				// Check that cs*cs + sn*sn = 1.
				resid = math.Abs(1 - (cs*cs + sn*sn))
				if resid > tol {
					t.Errorf("%v: cs*cs + sn*sn != 1; resid=%v", name, resid)
				}

				// Check that cs is non-negative.
				if math.Abs(f) > math.Abs(g) && cs < 0 {
					t.Errorf("%v: cs is negative; cs=%v", name, cs)
				}
			}
		}
	}
}
