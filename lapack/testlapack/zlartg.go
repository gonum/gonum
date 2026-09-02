// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"
)

type Zlartger interface {
	Zlartg(f, g complex128) (cs float64, sn, r complex128)
}

func ZlartgTest(t *testing.T, impl Zlartger) {
	const tol = 40 * ulp

	reals := []float64{
		-1 / ulp,
		-1,
		-1.0 / 3,
		-ulp,
		0,
		ulp,
		1.0 / 3,
		1,
		1 / ulp,
	}
	// Build a small collection of complex test values.
	var values []complex128
	for _, re := range reals {
		for _, im := range reals {
			values = append(values, complex(re, im))
		}
	}
	// A couple of safmin/safmax magnitudes.
	values = append(values,
		complex(safmin, 0),
		complex(0, safmin),
		complex(safmax, 0),
		complex(0, safmax),
	)

	for _, f := range values {
		for _, g := range values {
			name := fmt.Sprintf("Case f=%v,g=%v", f, g)

			cs, sn, r := impl.Zlartg(f, g)

			// cs must be real non-negative and in [0,1].
			if math.IsNaN(cs) || cs < 0 || cs > 1 {
				t.Errorf("%v: bad cs=%v", name, cs)
				continue
			}
			// cs^2 + |sn|^2 should equal 1.
			mag := cs*cs + real(sn)*real(sn) + imag(sn)*imag(sn)
			if math.Abs(mag-1) > tol {
				t.Errorf("%v: cs^2+|sn|^2=%v not 1", name, mag)
			}

			// Skip residual check when either input is zero — covered by the
			// magnitude check above and the special properties below.
			if f == 0 {
				if cs != 1 && !(g == 0 && cs == 1) {
					// For f=0, g!=0: cs should be 0.
					if g != 0 && cs != 0 {
						t.Errorf("%v: f=0,g!=0 but cs=%v != 0", name, cs)
					}
				}
			}
			if g == 0 {
				if cs != 1 {
					t.Errorf("%v: g=0 but cs=%v != 1", name, cs)
				}
				if sn != 0 {
					t.Errorf("%v: g=0 but sn=%v != 0", name, sn)
				}
				if r != f {
					t.Errorf("%v: g=0 but r=%v != f=%v", name, r, f)
				}
				continue
			}

			// Scale inputs to a safe range for residual comparison.
			d := math.Max(zabs1Test(f), zabs1Test(g))
			d = math.Min(math.Max(safmin, d), safmax)
			fs := complex(real(f)/d, imag(f)/d)
			gs := complex(real(g)/d, imag(g)/d)
			rs := complex(real(r)/d, imag(r)/d)

			// Check that cs*f + sn*g = r.
			lhs := complex(cs, 0)*fs + sn*gs
			rnorm := cmplx.Abs(rs)
			if rnorm == 0 {
				rnorm = math.Max(cmplx.Abs(fs), cmplx.Abs(gs))
				if rnorm == 0 {
					rnorm = 1
				}
			}
			resid := cmplx.Abs(rs-lhs) / rnorm
			if resid > tol {
				t.Errorf("%v: cs*f+sn*g!=r; resid=%v", name, resid)
			}

			// Check that -conj(sn)*f + cs*g = 0.
			zero := -complex(real(sn), -imag(sn))*fs + complex(cs, 0)*gs
			resid = cmplx.Abs(zero)
			if resid > tol*math.Max(cmplx.Abs(fs), cmplx.Abs(gs)) && resid > tol {
				t.Errorf("%v: -conj(sn)*f+cs*g != 0; resid=%v", name, resid)
			}
		}
	}
}

func zabs1Test(z complex128) float64 {
	return math.Abs(real(z)) + math.Abs(imag(z))
}
