// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quat

import "testing"

var expTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{Real: 1}},
}

func TestExp(t *testing.T) {
	for _, test := range expTests {
		got := Exp(test.q)
		if got != test.want {
			t.Errorf("unexpected result for Exp(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var logTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{Real: -inf}},
}

func TestLog(t *testing.T) {
	for _, test := range logTests {
		got := Log(test.q)
		if got != test.want {
			t.Errorf("unexpected result for Log(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var powTests = []struct {
	q, r Quat
	want Quat
}{
	{q: Quat{}, r: Quat{}, want: Quat{Real: 1}},
}

func TestPow(t *testing.T) {
	for _, test := range powTests {
		got := Pow(test.q, test.r)
		if got != test.want {
			t.Errorf("unexpected result for Pow(%v, %v): got:%v want:%v", test.q, test.r, got, test.want)
		}
	}
}

var sqrtTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
}

func TestSqrt(t *testing.T) {
	for _, test := range sqrtTests {
		got := Sqrt(test.q)
		if got != test.want {
			t.Errorf("unexpected result for Sqrt(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}
