// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import "math"

// Dlartg generates a plane rotation so that
//
//	[ cs sn] * [f] = [r]
//	[-sn cs]   [g] = [0]
//
// where cs*cs + sn*sn = 1.
//
// This is a more accurate version of BLAS Drotg, with the other differences
// that
//   - if g = 0, then cs = 1 and sn = 0
//   - if f = 0 and g != 0, then cs = 0 and sn = 1
//   - r takes the sign of f and so cs is always non-negative
//
// Dlartg is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dlartg(f, g float64) (cs, sn, r float64) {
	// Implementation based on Supplemental Material to:
	// Edward Anderson. 2017. Algorithm 978: Safe Scaling in the Level 1 BLAS.
	// ACM Trans. Math. Softw. 44, 1, Article 12 (July 2017), 28 pages.
	// DOI: https://doi.org/10.1145/3061665
	const safmin = dlamchS
	const safmax = 1 / safmin
	f1 := math.Abs(f)
	g1 := math.Abs(g)
	switch {
	case g == 0:
		cs = 1
		sn = 0
		r = f
	case f == 0:
		cs = 0
		sn = math.Copysign(1, g)
		r = g1
	case drtmin < f1 && f1 < drtmax && drtmin < g1 && g1 < drtmax:
		d := math.Sqrt(f*f + g*g)
		p := 1 / d
		cs = f1 * p
		sn = g * math.Copysign(p, f)
		r = math.Copysign(d, f)
	default:
		maxfg := math.Max(f1, g1)
		u := math.Min(math.Max(safmin, maxfg), safmax)
		uu := 1 / u
		fs := f * uu
		gs := g * uu
		d := math.Sqrt(fs*fs + gs*gs)
		p := 1 / d
		cs = math.Abs(fs) * p
		sn = gs * math.Copysign(p, f)
		r = math.Copysign(d, f) * u
	}
	return cs, sn, r
}
