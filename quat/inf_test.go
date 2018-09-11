// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quat

import (
	"math"
	"testing"
)

var inf = math.Inf(1)

var infTests = []struct {
	q    Quat
	want bool
}{
	{q: Inf(), want: true},
	{q: Quat{Real: inf, Imag: inf, Jmag: inf, Kmag: inf}, want: true},
	{q: Quat{Real: -inf, Imag: -inf, Jmag: -inf, Kmag: -inf}, want: true},
	{q: Quat{Real: inf, Imag: nan, Jmag: nan, Kmag: nan}, want: true},
	{q: Quat{Real: nan, Imag: inf, Jmag: nan, Kmag: nan}, want: true},
	{q: Quat{Real: nan, Imag: nan, Jmag: inf, Kmag: nan}, want: true},
	{q: Quat{Real: nan, Imag: nan, Jmag: nan, Kmag: inf}, want: true},
	{q: Quat{Real: -inf, Imag: nan, Jmag: nan, Kmag: nan}, want: true},
	{q: Quat{Real: nan, Imag: -inf, Jmag: nan, Kmag: nan}, want: true},
	{q: Quat{Real: nan, Imag: nan, Jmag: -inf, Kmag: nan}, want: true},
	{q: Quat{Real: nan, Imag: nan, Jmag: nan, Kmag: -inf}, want: true},
	{q: Quat{Real: inf}, want: true},
	{q: Quat{Imag: inf}, want: true},
	{q: Quat{Jmag: inf}, want: true},
	{q: Quat{Kmag: inf}, want: true},
	{q: Quat{Real: -inf}, want: true},
	{q: Quat{Imag: -inf}, want: true},
	{q: Quat{Jmag: -inf}, want: true},
	{q: Quat{Kmag: -inf}, want: true},
	{q: Quat{}, want: false},
}

func TestIsInf(t *testing.T) {
	for _, test := range infTests {
		got := IsInf(test.q)
		if got != test.want {
			t.Errorf("unexpected result for IsInf(%v): got:%t want:%t", test.q, got, test.want)
		}
	}
}
