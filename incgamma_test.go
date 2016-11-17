// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import (
	"math"
	"testing"
)

func TestIncGamma(t *testing.T) {
	for i, test := range []struct {
		a, x, want float64
	}{
		// Results computed using scipy.special.gamminc
		{0, 0, 0},
		{1, 1, 0.63212055882855778},
		{0.5, 2, 0.95449973610364147},
		{1.5, 0.75, 0.31772966966378746},
		{0.1, 10, 0.99999944520142825},
		{10, 5, 0.031828057306204811},
		{3, 7, 0.97036383611947818},
		{5, 50, 1},
		{2.5, 1, 0.15085496391539038},
		{0.01, 10, 0.99999995718295021},
	} {
		if got := IncGamma(test.a, test.x); math.Abs(got-test.want) > 1e-10 {
			t.Errorf("test %d IncGamma(%g, %g) failed: got %g want %g", i, test.a, test.x, got, test.want)
		}
	}
}
