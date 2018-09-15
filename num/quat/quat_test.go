// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quat

import (
	"fmt"
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
)

var arithTests = []struct {
	x, y Quat
	f    float64

	wantAdd   Quat
	wantSub   Quat
	wantMul   Quat
	wantScale Quat
}{
	{
		x: Quat{1, 1, 1, 1}, y: Quat{1, 1, 1, 1},
		f: 2,

		wantAdd:   Quat{2, 2, 2, 2},
		wantSub:   Quat{0, 0, 0, 0},
		wantMul:   Quat{-2, 2, 2, 2},
		wantScale: Quat{2, 2, 2, 2},
	},
	{
		x: Quat{1, 1, 1, 1}, y: Quat{2, -1, 1, -1},
		f: -2,

		wantAdd:   Quat{3, 0, 2, 0},
		wantSub:   Quat{-1, 2, 0, 2},
		wantMul:   Quat{3, -1, 3, 3},
		wantScale: Quat{-2, -2, -2, -2},
	},
	{
		x: Quat{1, 2, 3, 4}, y: Quat{4, -3, 2, -1},
		f: 2,

		wantAdd:   Quat{5, -1, 5, 3},
		wantSub:   Quat{-3, 5, 1, 5},
		wantMul:   Quat{8, -6, 4, 28},
		wantScale: Quat{2, 4, 6, 8},
	},
	{
		x: Quat{1, 2, 3, 4}, y: Quat{-4, 3, -2, 1},
		f: -2,

		wantAdd:   Quat{-3, 5, 1, 5},
		wantSub:   Quat{5, -1, 5, 3},
		wantMul:   Quat{-8, 6, -4, -28},
		wantScale: Quat{-2, -4, -6, -8},
	},
	{
		x: Quat{-4, 3, -2, 1}, y: Quat{1, 2, 3, 4},
		f: 0.5,

		wantAdd:   Quat{-3, 5, 1, 5},
		wantSub:   Quat{-5, 1, -5, -3},
		wantMul:   Quat{-8, -16, -24, -2},
		wantScale: Quat{-2, 1.5, -1, 0.5},
	},
}

func TestArithmetic(t *testing.T) {
	for _, test := range arithTests {
		gotAdd := Add(test.x, test.y)
		if gotAdd != test.wantAdd {
			t.Errorf("unexpected result for %v+%v: got:%v, want:%v", test.x, test.y, gotAdd, test.wantAdd)
		}
		gotSub := Sub(test.x, test.y)
		if gotSub != test.wantSub {
			t.Errorf("unexpected result for %v-%v: got:%v, want:%v", test.x, test.y, gotSub, test.wantSub)
		}
		gotMul := Mul(test.x, test.y)
		if gotMul != test.wantMul {
			t.Errorf("unexpected result for %v*%v: got:%v, want:%v", test.x, test.y, gotMul, test.wantMul)
		}
		gotScale := Scale(test.f, test.x)
		if gotScale != test.wantScale {
			t.Errorf("unexpected result for %v*%v: got:%v, want:%v", test.f, test.x, gotScale, test.wantScale)
		}
	}
}

var formatTests = []struct {
	q      Quat
	format string
	want   string
}{
	{q: Quat{1.1, 2.1, 3.1, 4.1}, format: "%#v", want: "quat.Quat{1.1, 2.1, 3.1, 4.1}"},         // Bootstrap test.
	{q: Quat{-1.1, -2.1, -3.1, -4.1}, format: "%#v", want: "quat.Quat{-1.1, -2.1, -3.1, -4.1}"}, // Bootstrap test.
	{q: Quat{1, 2, 3, 4}, format: "%v", want: "(1+2i+3j+4k)"},
	{q: Quat{-1, -2, -3, -4}, format: "%v", want: "(-1-2i-3j-4k)"},
	{q: Quat{1, 2, 3, 4}, format: "%g", want: "(1+2i+3j+4k)"},
	{q: Quat{-1, -2, -3, -4}, format: "%g", want: "(-1-2i-3j-4k)"},
	{q: Quat{1, 2, 3, 4}, format: "%e", want: "(1.000000e+00+2.000000e+00i+3.000000e+00j+4.000000e+00k)"},
	{q: Quat{-1, -2, -3, -4}, format: "%e", want: "(-1.000000e+00-2.000000e+00i-3.000000e+00j-4.000000e+00k)"},
	{q: Quat{1, 2, 3, 4}, format: "%E", want: "(1.000000E+00+2.000000E+00i+3.000000E+00j+4.000000E+00k)"},
	{q: Quat{-1, -2, -3, -4}, format: "%E", want: "(-1.000000E+00-2.000000E+00i-3.000000E+00j-4.000000E+00k)"},
	{q: Quat{1, 2, 3, 4}, format: "%f", want: "(1.000000+2.000000i+3.000000j+4.000000k)"},
	{q: Quat{-1, -2, -3, -4}, format: "%f", want: "(-1.000000-2.000000i-3.000000j-4.000000k)"},
}

func TestFormat(t *testing.T) {
	for _, test := range formatTests {
		got := fmt.Sprintf(test.format, test.q)
		if got != test.want {
			t.Errorf("unexpected result for fmt.Sprintf(%q, %#v): got:%q, want:%q", test.format, test.q, got, test.want)
		}
	}
}

var parseTests = []struct {
	s       string
	want    Quat
	wantErr error
}{
	// Simple error states:
	{s: "", wantErr: parseError{state: -1}},
	{s: "()", wantErr: parseError{string: "()", state: -1}},
	{s: "(1", wantErr: parseError{string: "(1", state: -1}},
	{s: "1)", wantErr: parseError{string: "1)", state: -1}},

	// Ambiguous parse error states:
	{s: "1+2i+3i", wantErr: parseError{string: "1+2i+3i", state: -1}},
	{s: "1+2i3j", wantErr: parseError{string: "1+2i3j", state: -1}},
	{s: "1e-4i-4k+10.3e6j+", wantErr: parseError{string: "1e-4i-4k+10.3e6j+", state: -1}},
	{s: "1e-4i-4k+10.3e6j-", wantErr: parseError{string: "1e-4i-4k+10.3e6j-", state: -1}},

	// Valid input:
	{s: "1+4i", want: Quat{Real: 1, Imag: 4}},
	{s: "4i+1", want: Quat{Real: 1, Imag: 4}},
	{s: "+1+4i", want: Quat{Real: 1, Imag: 4}},
	{s: "+4i+1", want: Quat{Real: 1, Imag: 4}},
	{s: "1e-4-4k+10.3e6j+1i", want: Quat{Real: 1e-4, Imag: 1, Jmag: 10.3e6, Kmag: -4}},
	{s: "1e-4-4k+10.3e6j+i", want: Quat{Real: 1e-4, Imag: 1, Jmag: 10.3e6, Kmag: -4}},
	{s: "1e-4-4k+10.3e6j-i", want: Quat{Real: 1e-4, Imag: -1, Jmag: 10.3e6, Kmag: -4}},
	{s: "1e-4i-4k+10.3e6j-1", want: Quat{Real: -1, Imag: 1e-4, Jmag: 10.3e6, Kmag: -4}},
	{s: "1e-4i-4k+10.3e6j+1", want: Quat{Real: 1, Imag: 1e-4, Jmag: 10.3e6, Kmag: -4}},
	{s: "(1+4i)", want: Quat{Real: 1, Imag: 4}},
	{s: "(4i+1)", want: Quat{Real: 1, Imag: 4}},
	{s: "(+1+4i)", want: Quat{Real: 1, Imag: 4}},
	{s: "(+4i+1)", want: Quat{Real: 1, Imag: 4}},
	{s: "(1e-4-4k+10.3e6j+1i)", want: Quat{Real: 1e-4, Imag: 1, Jmag: 10.3e6, Kmag: -4}},
	{s: "(1e-4-4k+10.3e6j+i)", want: Quat{Real: 1e-4, Imag: 1, Jmag: 10.3e6, Kmag: -4}},
	{s: "(1e-4-4k+10.3e6j-i)", want: Quat{Real: 1e-4, Imag: -1, Jmag: 10.3e6, Kmag: -4}},
	{s: "(1e-4i-4k+10.3e6j-1)", want: Quat{Real: -1, Imag: 1e-4, Jmag: 10.3e6, Kmag: -4}},
	{s: "(1e-4i-4k+10.3e6j+1)", want: Quat{Real: 1, Imag: 1e-4, Jmag: 10.3e6, Kmag: -4}},
	{s: "NaN", want: NaN()},
	{s: "nan", want: NaN()},
	{s: "Inf", want: Inf()},
	{s: "inf", want: Inf()},
	{s: "(Inf+Infi)", want: Quat{Real: math.Inf(1), Imag: math.Inf(1)}},
	{s: "(-Inf+Infi)", want: Quat{Real: math.Inf(-1), Imag: math.Inf(1)}},
	{s: "(+Inf-Infi)", want: Quat{Real: math.Inf(1), Imag: math.Inf(-1)}},
	{s: "(inf+infi)", want: Quat{Real: math.Inf(1), Imag: math.Inf(1)}},
	{s: "(-inf+infi)", want: Quat{Real: math.Inf(-1), Imag: math.Inf(1)}},
	{s: "(+inf-infi)", want: Quat{Real: math.Inf(1), Imag: math.Inf(-1)}},
	{s: "(nan+nani)", want: Quat{Real: math.NaN(), Imag: math.NaN()}},
	{s: "(nan-nani)", want: Quat{Real: math.NaN(), Imag: math.NaN()}},
	{s: "(nan+nani+1k)", want: Quat{Real: math.NaN(), Imag: math.NaN(), Kmag: 1}},
	{s: "(nan-nani+1k)", want: Quat{Real: math.NaN(), Imag: math.NaN(), Kmag: 1}},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		got, err := Parse(test.s)
		if err != test.wantErr {
			t.Errorf("unexpected error for Parse(%q): got:%#v, want:%#v", test.s, err, test.wantErr)
		}
		if err != nil {
			continue
		}
		if !sameQuat(got, test.want) {
			t.Errorf("unexpected result for Parse(%q): got:%v, want:%v", test.s, got, test.want)
		}
	}
}

func equalApprox(a, b Quat, tol float64) bool {
	return floats.EqualWithinAbsOrRel(a.Real, b.Real, tol, tol) &&
		floats.EqualWithinAbsOrRel(a.Imag, b.Imag, tol, tol) &&
		floats.EqualWithinAbsOrRel(a.Jmag, b.Jmag, tol, tol) &&
		floats.EqualWithinAbsOrRel(a.Kmag, b.Kmag, tol, tol)
}

func sameQuat(a, b Quat) bool {
	return a == b || (sameFloat(a.Real, b.Real) &&
		sameFloat(a.Imag, b.Imag) &&
		sameFloat(a.Jmag, b.Jmag) &&
		sameFloat(a.Kmag, b.Kmag))
}

func sameFloat(a, b float64) bool {
	return a == b || (math.IsNaN(a) && math.IsNaN(b))
}
