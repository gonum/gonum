// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quat

import (
	"math"
	"testing"
)

var sinTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
}

func TestSin(t *testing.T) {
	for _, test := range sinTests {
		got := Sin(test.q)
		if got != test.want {
			t.Errorf("unexpected result for Sin(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var sinhTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
}

func TestSinh(t *testing.T) {
	for _, test := range sinhTests {
		got := Sinh(test.q)
		if got != test.want {
			t.Errorf("unexpected result for Sinh(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var cosTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{Real: 1}},
}

func TestCos(t *testing.T) {
	for _, test := range cosTests {
		got := Cos(test.q)
		if got != test.want {
			t.Errorf("unexpected result for Cos(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var coshTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{Real: 1}},
}

func TestCosh(t *testing.T) {
	for _, test := range coshTests {
		got := Cosh(test.q)
		if got != test.want {
			t.Errorf("unexpected result for Cosh(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var tanTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
}

func TestTan(t *testing.T) {
	for _, test := range tanTests {
		got := Tan(test.q)
		if got != test.want {
			t.Errorf("unexpected result for Tan(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var tanhTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
}

func TestTanh(t *testing.T) {
	for _, test := range tanhTests {
		got := Tanh(test.q)
		if got != test.want {
			t.Errorf("unexpected result for Tanh(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var asinTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
}

func TestAsin(t *testing.T) {
	for _, test := range asinTests {
		got := Asin(test.q)
		if got != test.want {
			t.Errorf("unexpected result for Asin(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var asinhTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
}

func TestAsinh(t *testing.T) {
	for _, test := range asinhTests {
		got := Asinh(test.q)
		if got != test.want {
			t.Errorf("unexpected result for Asinh(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var acosTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{Real: math.Pi / 2}},
}

func TestAcos(t *testing.T) {
	for _, test := range acosTests {
		got := Acos(test.q)
		if got != test.want {
			t.Errorf("unexpected result for Acos(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var acoshTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{Real: math.Pi / 2}},
}

func TestAcosh(t *testing.T) {
	for _, test := range acoshTests {
		got := Acosh(test.q)
		if got != test.want {
			t.Errorf("unexpected result for Acosh(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var atanTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
}

func TestAtan(t *testing.T) {
	for _, test := range atanTests {
		got := Atan(test.q)
		if got != test.want {
			t.Errorf("unexpected result for Atan(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var atanhTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
}

func TestAtanh(t *testing.T) {
	for _, test := range atanhTests {
		got := Atanh(test.q)
		if got != test.want {
			t.Errorf("unexpected result for Atanh(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}
