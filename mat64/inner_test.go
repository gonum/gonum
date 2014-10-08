// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"math"
	"testing"
)

func TestInner(t *testing.T) {
	for i, test := range []struct {
		x   []float64
		y   []float64
		m   [][]float64
		ans float64
	}{
		{
			x:   []float64{5},
			y:   []float64{10},
			m:   [][]float64{{2}},
			ans: 100,
		},
		{
			x:   []float64{5, 6, 1},
			y:   []float64{10},
			m:   [][]float64{{2}, {-3}, {5}},
			ans: -30,
		},
		{
			x:   []float64{5},
			y:   []float64{10, 15},
			m:   [][]float64{{2, -3}},
			ans: -125,
		},
		{
			x: []float64{1, 5},
			y: []float64{10, 15},
			m: [][]float64{
				{2, -3},
				{4, -1},
			},
			ans: 100,
		},
		{
			x: []float64{2, 3, 9},
			y: []float64{8, 9},
			m: [][]float64{
				{2, 3},
				{4, 5},
				{6, 7},
			},
			ans: 1316,
		},
		{
			x: []float64{2, 3},
			y: []float64{8, 9, 9},
			m: [][]float64{
				{2, 3, 6},
				{4, 5, 7},
			},
			ans: 614,
		},
	} {
		m := NewDense(flatten(test.m))
		ans := Inner(test.x, m, test.y)
		if math.Abs(ans-test.ans) > 1e-14 {
			t.Errorf("Inner product mismatch case %v. Want: %v, Got: %v", i, test.ans, ans)
		}
	}
}
