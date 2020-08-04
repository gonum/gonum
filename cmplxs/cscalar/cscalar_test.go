// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file

package cscalar

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestEqualsRelative(t *testing.T) {
	equalityTests := []struct {
		a, b  float64
		tol   float64
		equal bool
	}{
		{1000000, 1000001, 0, true},
		{1000001, 1000000, 0, true},
		{10000, 10001, 0, false},
		{10001, 10000, 0, false},
		{-1000000, -1000001, 0, true},
		{-1000001, -1000000, 0, true},
		{-10000, -10001, 0, false},
		{-10001, -10000, 0, false},
		{1.0000001, 1.0000002, 0, true},
		{1.0000002, 1.0000001, 0, true},
		{1.0002, 1.0001, 0, false},
		{1.0001, 1.0002, 0, false},
		{-1.000001, -1.000002, 0, true},
		{-1.000002, -1.000001, 0, true},
		{-1.0001, -1.0002, 0, false},
		{-1.0002, -1.0001, 0, false},
		{0.000000001000001, 0.000000001000002, 0, true},
		{0.000000001000002, 0.000000001000001, 0, true},
		{0.000000000001002, 0.000000000001001, 0, false},
		{0.000000000001001, 0.000000000001002, 0, false},
		{-0.000000001000001, -0.000000001000002, 0, true},
		{-0.000000001000002, -0.000000001000001, 0, true},
		{-0.000000000001002, -0.000000000001001, 0, false},
		{-0.000000000001001, -0.000000000001002, 0, false},
		{0, 0, 0, true},
		{0, -0, 0, true},
		{-0, -0, 0, true},
		{0.00000001, 0, 0, false},
		{0, 0.00000001, 0, false},
		{-0.00000001, 0, 0, false},
		{0, -0.00000001, 0, false},
		{0, 1e-310, 0.01, true},
		{1e-310, 0, 0.01, true},
		{1e-310, 0, 0.000001, false},
		{0, 1e-310, 0.000001, false},
		{0, -1e-310, 0.1, true},
		{-1e-310, 0, 0.1, true},
		{-1e-310, 0, 0.00000001, false},
		{0, -1e-310, 0.00000001, false},
		{math.Inf(1), math.Inf(1), 0, true},
		{math.Inf(1), math.MaxFloat64, 0, false},
		{math.NaN(), math.NaN(), 0, false},
		{math.NaN(), 0, 0, false},
		{-0, math.NaN(), 0, false},
		{math.NaN(), -0, 0, false},
		{0, math.NaN(), 0, false},
		{math.NaN(), math.Inf(1), 0, false},
		{math.Inf(1), math.NaN(), 0, false},
		{math.NaN(), math.MaxFloat64, 0, false},
		{math.MaxFloat64, math.NaN(), 0, false},
		{math.NaN(), -math.MaxFloat64, 0, false},
		{-math.MaxFloat64, math.NaN(), 0, false},
		{math.NaN(), math.SmallestNonzeroFloat64, 0, false},
		{math.SmallestNonzeroFloat64, math.NaN(), 0, false},
		{math.NaN(), -math.SmallestNonzeroFloat64, 0, false},
		{-math.SmallestNonzeroFloat64, math.NaN(), 0, false},
		{1.000000001, -1.0, 0, false},
		{-1.0, 1.000000001, 0, false},
		{-1.000000001, 1.0, 0, false},
		{1.0, -1.000000001, 0, false},
		{10 * math.SmallestNonzeroFloat64, 10 * -math.SmallestNonzeroFloat64, 0, true},
		{1e11 * math.SmallestNonzeroFloat64, 1e11 * -math.SmallestNonzeroFloat64, 0, false},
		{math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64, 0, true},
		{-math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, 0, true},
		{math.SmallestNonzeroFloat64, 0, 0, true},
		{0, math.SmallestNonzeroFloat64, 0, true},
		{-math.SmallestNonzeroFloat64, 0, 0, true},
		{0, -math.SmallestNonzeroFloat64, 0, true},
		{0.000000001, -math.SmallestNonzeroFloat64, 0, false},
		{0.000000001, math.SmallestNonzeroFloat64, 0, false},
		{math.SmallestNonzeroFloat64, 0.000000001, 0, false},
		{-math.SmallestNonzeroFloat64, 0.000000001, 0, false},
	}
	for _, ts := range equalityTests {
		if ts.tol == 0 {
			ts.tol = 1e-5
		}

		for _, comp := range []struct{ a, b complex128 }{
			{a: complex(ts.a, 0), b: complex(ts.b, 0)},
			{a: complex(0, ts.a), b: complex(0, ts.b)},
			{a: complex(ts.a, ts.a), b: complex(ts.b, ts.b)},
		} {
			if equal := EqualWithinRel(comp.a, comp.b, ts.tol); equal != ts.equal {
				t.Errorf("Relative equality of %g and %g with tolerance %g returned: %v. Expected: %v",
					comp.a, comp.b, ts.tol, equal, ts.equal)
			}
		}
	}
}

func TestRoundEven(t *testing.T) {
	for _, test := range []struct {
		x    complex128
		prec int
		want complex128
	}{
		{x: 0, prec: 1, want: 0},
		{x: cmplx.Inf(), prec: 1, want: cmplx.Inf()},
		{x: cmplx.NaN(), prec: 1, want: cmplx.NaN()},
		{x: func() complex128 { var f complex128; return -f }(), prec: 1, want: 0},
		{x: math.MaxFloat64 / 2, prec: 1, want: math.MaxFloat64 / 2},
		{x: 1 << 64, prec: 1, want: 1 << 64},
		{x: 454.4445, prec: 3, want: 454.444},
		{x: 454.44445, prec: 4, want: 454.4444},
		{x: 0.42499, prec: 4, want: 0.425},
		{x: 0.42599, prec: 4, want: 0.426},
		{x: 0.424999999999993, prec: 2, want: 0.42},
		{x: 0.425, prec: 2, want: 0.42},
		{x: 0.425000000000001, prec: 2, want: 0.43},
		{x: 123.4244999999999, prec: 3, want: 123.424},
		{x: 123.4245, prec: 3, want: 123.424},
		{x: 123.4245000000001, prec: 3, want: 123.425},

		{x: 454.45, prec: 0, want: 454},
		{x: 454.45, prec: 1, want: 454.4},
		{x: 454.45, prec: 2, want: 454.45},
		{x: 454.45, prec: 3, want: 454.45},
		{x: 454.445, prec: 0, want: 454},
		{x: 454.445, prec: 1, want: 454.4},
		{x: 454.445, prec: 2, want: 454.44},
		{x: 454.445, prec: 3, want: 454.445},
		{x: 454.445, prec: 4, want: 454.445},
		{x: 454.55, prec: 0, want: 455},
		{x: 454.55, prec: 1, want: 454.6},
		{x: 454.55, prec: 2, want: 454.55},
		{x: 454.55, prec: 3, want: 454.55},
		{x: 454.455, prec: 0, want: 454},
		{x: 454.455, prec: 1, want: 454.5},
		{x: 454.455, prec: 2, want: 454.46},
		{x: 454.455, prec: 3, want: 454.455},
		{x: 454.455, prec: 4, want: 454.455},

		// Negative precision.
		{x: 454.45, prec: -1, want: 450},
		{x: 454.45, prec: -2, want: 500},
		{x: 500, prec: -3, want: 0},
		{x: 500, prec: -4, want: 0},
		{x: 1500, prec: -3, want: 2000},
		{x: 1500, prec: -4, want: 0},
	} {
		for _, sign := range []complex128{1, -1} {
			got := RoundEven(sign*test.x, test.prec)
			want := sign * test.want
			if want == 0 {
				want = 0
			}
			// FIXME(kortschak): Complexify this.
			if (got != want || math.Signbit(real(got)) != math.Signbit(real(want))) && !(math.IsNaN(real(got)) && math.IsNaN(real(want))) {
				t.Errorf("unexpected result for RoundEven(%g, %d): got: %g, want: %g", sign*test.x, test.prec, got, want)
			}
		}
	}
}

func TestSame(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		a, b complex128
		want bool
	}{
		{a: 0, b: 0, want: true},
		{a: 1, b: 1, want: true},
		{a: -1, b: 1, want: false},
		{a: 0, b: 1, want: false},
		{a: 1, b: 0, want: false},
		{a: -1, b: 1, want: false},
		{a: cmplx.NaN(), b: cmplx.NaN(), want: true},
		{a: 1, b: cmplx.NaN(), want: false},
		{a: cmplx.Inf(), b: cmplx.NaN(), want: false},
		{a: cmplx.NaN(), b: cmplx.Inf(), want: false},
		{a: cmplx.NaN(), b: 1, want: false},
		{a: cmplx.Inf(), b: cmplx.Inf(), want: true},
	} {
		got := Same(test.a, test.b)
		if got != test.want {
			t.Errorf("unexpected results for a=%f b=%f: got:%t want:%t", test.a, test.b, got, test.want)
		}
	}
}
