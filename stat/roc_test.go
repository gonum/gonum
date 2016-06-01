// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"testing"

	"gonum.org/v1/gonum/floats"
)

// Test cases were calculated manually.
func TestROC(t *testing.T) {
	cases := []struct {
		y       []float64
		c       []bool
		w       []float64
		cutoffs []float64
		wantTPR []float64
		wantFPR []float64
	}{
		{ // 0
			y:       []float64{0, 3, 5, 6, 7.5, 8},
			c:       []bool{true, false, true, false, false, false},
			wantTPR: []float64{0, 0.5, 0.5, 1, 1, 1, 1},
			wantFPR: []float64{0, 0, 0.25, 0.25, 0.5, 0.75, 1},
		},
		{ // 1
			y:       []float64{0, 3, 5, 6, 7.5, 8},
			c:       []bool{true, false, true, false, false, false},
			w:       []float64{4, 1, 6, 3, 2, 2},
			wantTPR: []float64{0, 0.4, 0.4, 1, 1, 1, 1},
			wantFPR: []float64{0, 0, 0.125, 0.125, 0.5, 0.75, 1},
		},
		{ // 2
			y:       []float64{0, 3, 5, 6, 7.5, 8},
			c:       []bool{true, false, true, false, false, false},
			cutoffs: []float64{-1, 2, 4, 6, 8},
			wantTPR: []float64{0, 0.5, 0.5, 1, 1},
			wantFPR: []float64{0, 0, 0.25, 0.5, 1},
		},
		{ // 3
			y:       []float64{0, 3, 5, 6, 7.5, 8},
			c:       []bool{true, false, true, false, false, false},
			cutoffs: []float64{-1, 1, 2, 3, 4, 5, 6, 7, 8},
			wantTPR: []float64{0, 0.5, 0.5, 0.5, 0.5, 1, 1, 1, 1},
			wantFPR: []float64{0, 0, 0, 0.25, 0.25, 0.25, 0.5, 0.5, 1},
		},
		{ // 4
			y:       []float64{0, 3, 5, 6, 7.5, 8},
			c:       []bool{true, false, true, false, false, false},
			w:       []float64{4, 1, 6, 3, 2, 2},
			cutoffs: []float64{-1, 2, 4, 6, 8},
			wantTPR: []float64{0, 0.4, 0.4, 1, 1},
			wantFPR: []float64{0, 0, 0.125, 0.5, 1},
		},
		{ // 5
			y:       []float64{0, 3, 5, 6, 7.5, 8},
			c:       []bool{true, false, true, false, false, false},
			w:       []float64{4, 1, 6, 3, 2, 2},
			cutoffs: []float64{-1, 1, 2, 3, 4, 5, 6, 7, 8},
			wantTPR: []float64{0, 0.4, 0.4, 0.4, 0.4, 1, 1, 1, 1},
			wantFPR: []float64{0, 0, 0, 0.125, 0.125, 0.125, 0.5, 0.5, 1},
		},
		{ // 6
			y:       []float64{0, 3, 6, 6, 6, 8},
			c:       []bool{true, false, true, false, false, false},
			wantTPR: []float64{0, 0.5, 0.5, 1, 1},
			wantFPR: []float64{0, 0, 0.25, 0.75, 1},
		},
		{ // 7
			y:       []float64{0, 3, 6, 6, 6, 8},
			c:       []bool{true, false, true, false, false, false},
			w:       []float64{4, 1, 6, 3, 2, 2},
			wantTPR: []float64{0, 0.4, 0.4, 1, 1},
			wantFPR: []float64{0, 0, 0.125, 0.75, 1},
		},
		{ // 8
			y:       []float64{0, 3, 6, 6, 6, 8},
			c:       []bool{true, false, true, false, false, false},
			cutoffs: []float64{-1, 2, 4, 6, 8},
			wantTPR: []float64{0, 0.5, 0.5, 1, 1},
			wantFPR: []float64{0, 0, 0.25, 0.75, 1},
		},
		{ // 9
			y:       []float64{0, 3, 6, 6, 6, 8},
			c:       []bool{true, false, true, false, false, false},
			cutoffs: []float64{-1, 1, 2, 3, 4, 5, 6, 7, 8},
			wantTPR: []float64{0, 0.5, 0.5, 0.5, 0.5, 0.5, 1, 1, 1},
			wantFPR: []float64{0, 0, 0, 0.25, 0.25, 0.25, 0.75, 0.75, 1},
		},
		{ // 10
			y:       []float64{0, 3, 6, 6, 6, 8},
			c:       []bool{true, false, true, false, false, false},
			w:       []float64{4, 1, 6, 3, 2, 2},
			cutoffs: []float64{-1, 2, 4, 6, 8},
			wantTPR: []float64{0, 0.4, 0.4, 1, 1},
			wantFPR: []float64{0, 0, 0.125, 0.75, 1},
		},
		{ // 11
			y:       []float64{0, 3, 6, 6, 6, 8},
			c:       []bool{true, false, true, false, false, false},
			w:       []float64{4, 1, 6, 3, 2, 2},
			cutoffs: []float64{-1, 1, 2, 3, 4, 5, 6, 7, 8},
			wantTPR: []float64{0, 0.4, 0.4, 0.4, 0.4, 0.4, 1, 1, 1},
			wantFPR: []float64{0, 0, 0, 0.125, 0.125, 0.125, 0.75, 0.75, 1},
		},
		{ // 12
			y:       []float64{1, 2},
			c:       []bool{true, true},
			wantTPR: []float64{0, 0.5, 1},
			wantFPR: []float64{0, 0, 0},
		},
		{ // 13
			y:       []float64{1, 2},
			c:       []bool{true, true},
			cutoffs: []float64{-1, 2},
			wantTPR: []float64{0, 1},
			wantFPR: []float64{0, 0},
		},
		{ // 14
			y:       []float64{1, 2},
			c:       []bool{true, true},
			cutoffs: []float64{0, 1.2, 1.4, 1.6, 1.8, 2},
			wantTPR: []float64{0, 0.5, 0.5, 0.5, 0.5, 1},
			wantFPR: []float64{0, 0, 0, 0, 0, 0},
		},
		{ // 15
			y:       []float64{1},
			c:       []bool{true},
			wantTPR: []float64{0, 1},
			wantFPR: []float64{0, 0},
		},
		{ // 16
			y:       []float64{1},
			c:       []bool{true},
			cutoffs: []float64{-1, 1},
			wantTPR: []float64{0, 1},
			wantFPR: []float64{0, 0},
		},
		{ // 17
			y:       []float64{1},
			c:       []bool{false},
			wantTPR: []float64{0, 0},
			wantFPR: []float64{0, 1},
		},
		{ // 18
			y:       []float64{0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 10},
			c:       []bool{true, false, true, true, false, false, true},
			cutoffs: []float64{-1, 2.5, 5, 7.5, 10},
			wantTPR: []float64{0, 0.75, 0.75, 0.75, 1},
			wantFPR: []float64{0, 1, 1, 1, 1},
		},
		{ // 19
			y:       []float64{},
			c:       []bool{},
			wantTPR: nil,
			wantFPR: nil,
		},
		{ // 20
			y:       []float64{},
			c:       []bool{},
			cutoffs: []float64{-1, 2.5, 5, 7.5, 10},
			wantTPR: nil,
			wantFPR: nil,
		},
	}
	for i, test := range cases {
		gotTPR, gotFPR := ROC(test.cutoffs, test.y, test.c, test.w)
		if !floats.Same(gotTPR, test.wantTPR) {
			t.Errorf("%d: unexpected TPR got:%v want:%v", i, gotTPR, test.wantTPR)
		}
		if !floats.Same(gotFPR, test.wantFPR) {
			t.Errorf("%d: unexpected FPR got:%v want:%v", i, gotFPR, test.wantFPR)
		}
	}
}
