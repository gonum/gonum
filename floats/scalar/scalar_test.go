// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scalar

import (
	"math"
	"testing"
)

func TestEqualsRelative(t *testing.T) {
	t.Parallel()
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
		{math.Inf(-1), math.Inf(-1), 0, true},
		{math.Inf(-1), math.Inf(1), 0, false},
		{math.Inf(1), math.MaxFloat64, 0, false},
		{math.Inf(-1), -math.MaxFloat64, 0, false},
		{math.NaN(), math.NaN(), 0, false},
		{math.NaN(), 0, 0, false},
		{-0, math.NaN(), 0, false},
		{math.NaN(), -0, 0, false},
		{0, math.NaN(), 0, false},
		{math.NaN(), math.Inf(1), 0, false},
		{math.Inf(1), math.NaN(), 0, false},
		{math.NaN(), math.Inf(-1), 0, false},
		{math.Inf(-1), math.NaN(), 0, false},
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
		if equal := EqualWithinRel(ts.a, ts.b, ts.tol); equal != ts.equal {
			t.Errorf("Relative equality of %g and %g with tolerance %g returned: %v. Expected: %v",
				ts.a, ts.b, ts.tol, equal, ts.equal)
		}
	}
}

func TestEqualsULP(t *testing.T) {
	t.Parallel()
	if f := 67329.242; !EqualWithinULP(f, nextAfterN(f, math.Inf(1), 10), 10) {
		t.Errorf("Equal values returned as unequal")
	}
	if f := 67329.242; EqualWithinULP(f, nextAfterN(f, math.Inf(1), 5), 1) {
		t.Errorf("Unequal values returned as equal")
	}
	if f := 67329.242; EqualWithinULP(nextAfterN(f, math.Inf(1), 5), f, 1) {
		t.Errorf("Unequal values returned as equal")
	}
	if f := nextAfterN(0, math.Inf(1), 2); !EqualWithinULP(f, nextAfterN(f, math.Inf(-1), 5), 10) {
		t.Errorf("Equal values returned as unequal")
	}
	if !EqualWithinULP(67329.242, 67329.242, 10) {
		t.Errorf("Equal float64s not returned as equal")
	}
	if EqualWithinULP(1, math.NaN(), 10) {
		t.Errorf("NaN returned as equal")
	}
}

func nextAfterN(x, y float64, n int) float64 {
	for i := 0; i < n; i++ {
		x = math.Nextafter(x, y)
	}
	return x
}

func TestNaNWith(t *testing.T) {
	t.Parallel()
	tests := []struct {
		payload uint64
		bits    uint64
	}{
		{0, math.Float64bits(0 / func() float64 { return 0 }())}, // Hide the division by zero from the compiler.
		{1, math.Float64bits(math.NaN())},
		{1954, 0x7ff80000000007a2}, // R NA.
	}

	for _, test := range tests {
		nan := NaNWith(test.payload)
		if !math.IsNaN(nan) {
			t.Errorf("expected NaN value, got:%f", nan)
		}

		bits := math.Float64bits(nan)

		// Strip sign bit.
		const sign = 1 << 63
		bits &^= sign
		test.bits &^= sign

		if bits != test.bits {
			t.Errorf("expected NaN bit representation: got:%x want:%x", bits, test.bits)
		}
	}
}

func TestNaNPayload(t *testing.T) {
	t.Parallel()
	tests := []struct {
		f       float64
		payload uint64
		ok      bool
	}{
		{0 / func() float64 { return 0 }(), 0, true}, // Hide the division by zero from the compiler.

		// The following two line are written explicitly to defend against potential changes to math.Copysign.
		{math.Float64frombits(math.Float64bits(math.NaN()) | (1 << 63)), 1, true},  // math.Copysign(math.NaN(), -1)
		{math.Float64frombits(math.Float64bits(math.NaN()) &^ (1 << 63)), 1, true}, // math.Copysign(math.NaN(), 1)

		{NaNWith(1954), 1954, true}, // R NA.

		{math.Copysign(0, -1), 0, false},
		{0, 0, false},
		{math.Inf(-1), 0, false},
		{math.Inf(1), 0, false},

		{math.Float64frombits(0x7ff0000000000001), 0, false}, // Signalling NaN.
	}

	for _, test := range tests {
		payload, ok := NaNPayload(test.f)
		if payload != test.payload {
			t.Errorf("expected NaN payload: got:%x want:%x", payload, test.payload)
		}
		if ok != test.ok {
			t.Errorf("expected NaN status: got:%t want:%t", ok, test.ok)
		}
	}
}

func TestRound(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		x    float64
		prec int
		want float64
	}{
		{x: 0, prec: 1, want: 0},
		{x: math.Inf(1), prec: 1, want: math.Inf(1)},
		{x: math.Inf(-1), prec: 1, want: math.Inf(-1)},
		{x: math.NaN(), prec: 1, want: math.NaN()},
		{x: func() float64 { var f float64; return -f }(), prec: 1, want: 0},
		{x: math.MaxFloat64 / 2, prec: 1, want: math.MaxFloat64 / 2},
		{x: 1 << 64, prec: 1, want: 1 << 64},
		{x: 454.4445, prec: 3, want: 454.445},
		{x: 454.44445, prec: 4, want: 454.4445},
		{x: 0.42499, prec: 4, want: 0.425},
		{x: 0.42599, prec: 4, want: 0.426},
		{x: 0.424999999999993, prec: 2, want: 0.42},
		{x: 0.425, prec: 2, want: 0.43},
		{x: 0.425000000000001, prec: 2, want: 0.43},
		{x: 123.4244999999999, prec: 3, want: 123.424},
		{x: 123.4245, prec: 3, want: 123.425},
		{x: 123.4245000000001, prec: 3, want: 123.425},

		{x: 454.45, prec: 0, want: 454},
		{x: 454.45, prec: 1, want: 454.5},
		{x: 454.45, prec: 2, want: 454.45},
		{x: 454.45, prec: 3, want: 454.45},
		{x: 454.445, prec: 0, want: 454},
		{x: 454.445, prec: 1, want: 454.4},
		{x: 454.445, prec: 2, want: 454.45},
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
		{x: math.Inf(1), prec: -1, want: math.Inf(1)},
		{x: math.Inf(-1), prec: -1, want: math.Inf(-1)},
		{x: math.NaN(), prec: -1, want: math.NaN()},
		{x: func() float64 { var f float64; return -f }(), prec: -1, want: 0},
		{x: 454.45, prec: -1, want: 450},
		{x: 454.45, prec: -2, want: 500},
		{x: 500, prec: -3, want: 1000},
		{x: 500, prec: -4, want: 0},
		{x: 1500, prec: -3, want: 2000},
		{x: 1500, prec: -4, want: 0},
	} {
		for _, sign := range []float64{1, -1} {
			got := Round(sign*test.x, test.prec)
			want := sign * test.want
			if want == 0 {
				want = 0
			}
			if (got != want || math.Signbit(got) != math.Signbit(want)) && !(math.IsNaN(got) && math.IsNaN(want)) {
				t.Errorf("unexpected result for Round(%g, %d): got: %g, want: %g", sign*test.x, test.prec, got, want)
			}
		}
	}
}

func TestRoundEven(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		x    float64
		prec int
		want float64
	}{
		{x: 0, prec: 1, want: 0},
		{x: math.Inf(1), prec: 1, want: math.Inf(1)},
		{x: math.Inf(-1), prec: 1, want: math.Inf(-1)},
		{x: math.NaN(), prec: 1, want: math.NaN()},
		{x: func() float64 { var f float64; return -f }(), prec: 1, want: 0},
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
		{x: math.Inf(1), prec: -1, want: math.Inf(1)},
		{x: math.Inf(-1), prec: -1, want: math.Inf(-1)},
		{x: math.NaN(), prec: -1, want: math.NaN()},
		{x: func() float64 { var f float64; return -f }(), prec: -1, want: 0},
		{x: 454.45, prec: -1, want: 450},
		{x: 454.45, prec: -2, want: 500},
		{x: 500, prec: -3, want: 0},
		{x: 500, prec: -4, want: 0},
		{x: 1500, prec: -3, want: 2000},
		{x: 1500, prec: -4, want: 0},
	} {
		for _, sign := range []float64{1, -1} {
			got := RoundEven(sign*test.x, test.prec)
			want := sign * test.want
			if want == 0 {
				want = 0
			}
			if (got != want || math.Signbit(got) != math.Signbit(want)) && !(math.IsNaN(got) && math.IsNaN(want)) {
				t.Errorf("unexpected result for RoundEven(%g, %d): got: %g, want: %g", sign*test.x, test.prec, got, want)
			}
		}
	}
}

func TestSame(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		a, b float64
		want bool
	}{
		{a: 0, b: 0, want: true},
		{a: 1, b: 1, want: true},
		{a: -1, b: 1, want: false},
		{a: 0, b: 1, want: false},
		{a: 1, b: 0, want: false},
		{a: -1, b: 1, want: false},
		{a: math.NaN(), b: math.NaN(), want: true},
		{a: 1, b: math.NaN(), want: false},
		{a: math.Inf(1), b: math.NaN(), want: false},
		{a: math.NaN(), b: math.Inf(1), want: false},
		{a: math.NaN(), b: 1, want: false},
		{a: math.Inf(1), b: math.Inf(1), want: true},
		{a: math.Inf(-1), b: math.Inf(1), want: false},
		{a: math.Inf(1), b: math.Inf(-1), want: false},
	} {
		got := Same(test.a, test.b)
		if got != test.want {
			t.Errorf("unexpected results for a=%f b=%f: got:%t want:%t", test.a, test.b, got, test.want)
		}
	}
}
