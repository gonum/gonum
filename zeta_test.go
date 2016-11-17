// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import (
	"math"
	"testing"
)

func TestZeta(t *testing.T) {
	for i, test := range []struct {
		x, q, want float64
	}{
		// Results computed using scipy.special.zeta
		{1, 1, math.MaxFloat64},
		{10, 0.5, 1024.0174503557578},
		{5, 2.5, 0.013073166646113805},
		{1.5, 2, 1.6123753486854886},
		{1.5, 20, 0.45287361712938717},
		{2.5, 0.5, 6.2471106345688137},
		{10, 7.5, 2.5578265694201971e-9},
		{12, 2.5, 1.7089167198843551e-5},
		{20, 0.75, 315.3368689825316},
		{25, 0.25, 1125899906842624.0},
	} {
		if got := Zeta(test.x, test.q); math.Abs(got-test.want) > 1e-10 {
			t.Errorf("test %d Zeta(%g, %g) failed: got %g want %g", i, test.x, test.q, got, test.want)
		}
	}
}
