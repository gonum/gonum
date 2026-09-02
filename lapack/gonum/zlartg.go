// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import "math"

// Zlartg generates a plane rotation so that
//
//	[  cs        sn] * [f] = [r]
//	[-conj(sn)   cs]   [g] = [0]
//
// where cs is real non-negative, cs*cs + conj(sn)*sn = 1, and |r| = |(f,g)|.
//
// This is a more accurate complex version of BLAS Zrotg that uses scaling to
// avoid overflow or underflow. Properties:
//   - cs >= 0
//   - if g = 0, then cs = 1 and sn = 0
//   - if f = 0 and g != 0, then cs = 0 and |sn| = 1
//
// Zlartg is an internal routine. It is exported for testing purposes.
func (impl Implementation) Zlartg(f, g complex128) (cs float64, sn, r complex128) {
	// Implementation based on Supplemental Material to:
	//
	// Edward Anderson
	// Algorithm 978: Safe Scaling in the Level 1 BLAS
	// ACM Trans. Math. Softw. 44, 1, Article 12 (2017)
	// DOI: https://doi.org/10.1145/3061665
	const safmin = dlamchS
	const safmax = 1 / safmin
	rtmin := math.Sqrt(safmin / dlamchE)
	rtmax := math.Sqrt(safmax * dlamchE)

	if g == 0 {
		return 1, 0, f
	}

	// Use Hypot for overflow-safe magnitudes.
	g1 := math.Hypot(real(g), imag(g))

	if f == 0 {
		cs = 0
		// sn = conj(g)/|g|, r = |g|.
		sn = complex(real(g)/g1, -imag(g)/g1)
		r = complex(g1, 0)
		return cs, sn, r
	}

	f1 := math.Hypot(real(f), imag(f))

	// If all quantities are well-scaled (squaring is safe on f and g), use
	// the direct formula.
	if rtmin < f1 && f1 < rtmax && rtmin < g1 && g1 < rtmax {
		f2 := zabssq(f)
		g2 := zabssq(g)
		h2 := f2 + g2
		var d float64
		if f2 > rtmin && h2 < rtmax {
			d = math.Sqrt(f2 * h2)
		} else {
			d = math.Sqrt(f2) * math.Sqrt(h2)
		}
		p := 1 / d
		cs = f2 * p
		// sn = conj(g) * f / d.
		sn = complex(real(g), -imag(g)) * complex(real(f)*p, imag(f)*p)
		// r = f * h2 / d.
		r = complex(real(f)*(h2*p), imag(f)*(h2*p))
		return cs, sn, r
	}

	// Otherwise scale by u = clamp(max(f1,g1), safmin, safmax).
	u := math.Min(math.Max(safmin, math.Max(f1, g1)), safmax)
	uu := 1 / u
	fs := complex(real(f)*uu, imag(f)*uu)
	gs := complex(real(g)*uu, imag(g)*uu)

	f2 := zabssq(fs)
	g2 := zabssq(gs)
	h2 := f2 + g2

	// If f fully underflows after scaling, fall back to the f == 0 branch.
	if f2 == 0 {
		cs = 0
		sn = complex(real(g)/g1, -imag(g)/g1)
		r = complex(g1, 0)
		return cs, sn, r
	}
	// If g fully underflows after scaling, fall back to the g == 0 branch
	// using the pre-scale f.
	if g2 == 0 {
		return 1, 0, f
	}

	var d float64
	if f2 > rtmin && h2 < rtmax {
		d = math.Sqrt(f2 * h2)
	} else {
		d = math.Sqrt(f2) * math.Sqrt(h2)
	}
	p := 1 / d
	cs = f2 * p
	sn = complex(real(gs), -imag(gs)) * complex(real(fs)*p, imag(fs)*p)
	r = complex(real(fs)*(h2*p), imag(fs)*(h2*p))
	r = complex(real(r)*u, imag(r)*u)
	return cs, sn, r
}

// zabssq returns real(z)^2 + imag(z)^2.
func zabssq(z complex128) float64 {
	re, im := real(z), imag(z)
	return re*re + im*im
}
