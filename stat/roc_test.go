// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"gonum.org/v1/gonum/floats"
)

func TestROC(t *testing.T) {
	const tol = 1e-14

	cases := []struct {
		y          []float64
		c          []bool
		w          []float64
		cutoffs    []float64
		wantTPR    []float64
		wantFPR    []float64
		wantThresh []float64
	}{
		// Test cases were informed by using sklearn metrics.roc_curve when
		// cutoffs is nil, but all test cases (including when cutoffs is not
		// nil) were calculated manually.
		// Some differences exist between unweighted ROCs from our function
		// and metrics.roc_curve which appears to use integer cutoffs in that
		// case. sklearn also appears to do some magic that trims leading zeros
		// sometimes.
		{ // 0
			y:          []float64{0, 3, 5, 6, 7.5, 8},
			c:          []bool{false, true, false, true, true, true},
			wantTPR:    []float64{0, 0.25, 0.5, 0.75, 0.75, 1, 1},
			wantFPR:    []float64{0, 0, 0, 0, 0.5, 0.5, 1},
			wantThresh: []float64{math.Inf(1), 8, 7.5, 6, 5, 3, 0},
		},
		{ // 1
			y:          []float64{0, 3, 5, 6, 7.5, 8},
			c:          []bool{false, true, false, true, true, true},
			w:          []float64{4, 1, 6, 3, 2, 2},
			wantTPR:    []float64{0, 0.25, 0.5, 0.875, 0.875, 1, 1},
			wantFPR:    []float64{0, 0, 0, 0, 0.6, 0.6, 1},
			wantThresh: []float64{math.Inf(1), 8, 7.5, 6, 5, 3, 0},
		},
		{ // 2
			y:          []float64{0, 3, 5, 6, 7.5, 8},
			c:          []bool{false, true, false, true, true, true},
			cutoffs:    []float64{-1, 2, 4, 6, 8},
			wantTPR:    []float64{0.25, 0.75, 0.75, 1, 1},
			wantFPR:    []float64{0, 0, 0.5, 0.5, 1},
			wantThresh: []float64{8, 6, 4, 2, -1},
		},
		{ // 3
			y:          []float64{0, 3, 5, 6, 7.5, 8},
			c:          []bool{false, true, false, true, true, true},
			cutoffs:    []float64{-1, 1, 2, 3, 4, 5, 6, 7, 8},
			wantTPR:    []float64{0.25, 0.5, 0.75, 0.75, 0.75, 1, 1, 1, 1},
			wantFPR:    []float64{0, 0, 0, 0.5, 0.5, 0.5, 0.5, 0.5, 1},
			wantThresh: []float64{8, 7, 6, 5, 4, 3, 2, 1, -1},
		},
		{ // 4
			y:          []float64{0, 3, 5, 6, 7.5, 8},
			c:          []bool{false, true, false, true, true, true},
			w:          []float64{4, 1, 6, 3, 2, 2},
			cutoffs:    []float64{-1, 2, 4, 6, 8},
			wantTPR:    []float64{0.25, 0.875, 0.875, 1, 1},
			wantFPR:    []float64{0, 0, 0.6, 0.6, 1},
			wantThresh: []float64{8, 6, 4, 2, -1},
		},
		{ // 5
			y:          []float64{0, 3, 5, 6, 7.5, 8},
			c:          []bool{false, true, false, true, true, true},
			w:          []float64{4, 1, 6, 3, 2, 2},
			cutoffs:    []float64{-1, 1, 2, 3, 4, 5, 6, 7, 8},
			wantTPR:    []float64{0.25, 0.5, 0.875, 0.875, 0.875, 1, 1, 1, 1},
			wantFPR:    []float64{0, 0, 0, 0.6, 0.6, 0.6, 0.6, 0.6, 1},
			wantThresh: []float64{8, 7, 6, 5, 4, 3, 2, 1, -1},
		},
		{ // 6
			y:          []float64{0, 3, 6, 6, 6, 8},
			c:          []bool{false, true, false, true, true, true},
			wantTPR:    []float64{0, 0.25, 0.75, 1, 1},
			wantFPR:    []float64{0, 0, 0.5, 0.5, 1},
			wantThresh: []float64{math.Inf(1), 8, 6, 3, 0},
		},
		{ // 7
			y:          []float64{0, 3, 6, 6, 6, 8},
			c:          []bool{false, true, false, true, true, true},
			w:          []float64{4, 1, 6, 3, 2, 2},
			wantTPR:    []float64{0, 0.25, 0.875, 1, 1},
			wantFPR:    []float64{0, 0, 0.6, 0.6, 1},
			wantThresh: []float64{math.Inf(1), 8, 6, 3, 0},
		},
		{ // 8
			y:          []float64{0, 3, 6, 6, 6, 8},
			c:          []bool{false, true, false, true, true, true},
			cutoffs:    []float64{-1, 2, 4, 6, 8},
			wantTPR:    []float64{0.25, 0.75, 0.75, 1, 1},
			wantFPR:    []float64{0, 0.5, 0.5, 0.5, 1},
			wantThresh: []float64{8, 6, 4, 2, -1},
		},
		{ // 9
			y:          []float64{0, 3, 6, 6, 6, 8},
			c:          []bool{false, true, false, true, true, true},
			cutoffs:    []float64{-1, 1, 2, 3, 4, 5, 6, 7, 8},
			wantTPR:    []float64{0.25, 0.25, 0.75, 0.75, 0.75, 1, 1, 1, 1},
			wantFPR:    []float64{0, 0, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 1},
			wantThresh: []float64{8, 7, 6, 5, 4, 3, 2, 1, -1},
		},
		{ // 10
			y:          []float64{0, 3, 6, 6, 6, 8},
			c:          []bool{false, true, false, true, true, true},
			w:          []float64{4, 1, 6, 3, 2, 2},
			cutoffs:    []float64{-1, 2, 4, 6, 8},
			wantTPR:    []float64{0.25, 0.875, 0.875, 1, 1},
			wantFPR:    []float64{0, 0.6, 0.6, 0.6, 1},
			wantThresh: []float64{8, 6, 4, 2, -1},
		},
		{ // 11
			y:          []float64{0, 3, 6, 6, 6, 8},
			c:          []bool{false, true, false, true, true, true},
			w:          []float64{4, 1, 6, 3, 2, 2},
			cutoffs:    []float64{-1, 1, 2, 3, 4, 5, 6, 7, 8},
			wantTPR:    []float64{0.25, 0.25, 0.875, 0.875, 0.875, 1, 1, 1, 1},
			wantFPR:    []float64{0, 0, 0.6, 0.6, 0.6, 0.6, 0.6, 0.6, 1},
			wantThresh: []float64{8, 7, 6, 5, 4, 3, 2, 1, -1},
		},
		{ // 12
			y:          []float64{0.1, 0.35, 0.4, 0.8},
			c:          []bool{true, false, true, false},
			wantTPR:    []float64{0, 0, 0.5, 0.5, 1},
			wantFPR:    []float64{0, 0.5, 0.5, 1, 1},
			wantThresh: []float64{math.Inf(1), 0.8, 0.4, 0.35, 0.1},
		},
		{ // 13
			y:          []float64{0.1, 0.35, 0.4, 0.8},
			c:          []bool{false, false, true, true},
			wantTPR:    []float64{0, 0.5, 1, 1, 1},
			wantFPR:    []float64{0, 0, 0, 0.5, 1},
			wantThresh: []float64{math.Inf(1), 0.8, 0.4, 0.35, 0.1},
		},
		{ // 14
			y:          []float64{0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 10},
			c:          []bool{false, true, false, false, true, true, false},
			cutoffs:    []float64{-1, 2.5, 5, 7.5, 10},
			wantTPR:    []float64{0, 0, 0, 0, 1},
			wantFPR:    []float64{0.25, 0.25, 0.25, 0.25, 1},
			wantThresh: []float64{10, 7.5, 5, 2.5, -1},
		},
		{ // 15
			y:          []float64{1, 2},
			c:          []bool{false, false},
			wantTPR:    []float64{math.NaN(), math.NaN(), math.NaN()},
			wantFPR:    []float64{0, 0.5, 1},
			wantThresh: []float64{math.Inf(1), 2, 1},
		},
		{ // 16
			y:          []float64{1, 2},
			c:          []bool{false, false},
			cutoffs:    []float64{-1, 2},
			wantTPR:    []float64{math.NaN(), math.NaN()},
			wantFPR:    []float64{0.5, 1},
			wantThresh: []float64{2, -1},
		},
		{ // 17
			y:          []float64{1, 2},
			c:          []bool{false, false},
			cutoffs:    []float64{0, 1.2, 1.4, 1.6, 1.8, 2},
			wantTPR:    []float64{math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.NaN()},
			wantFPR:    []float64{0.5, 0.5, 0.5, 0.5, 0.5, 1},
			wantThresh: []float64{2, 1.8, 1.6, 1.4, 1.2, 0},
		},
		{ // 18
			y:          []float64{1},
			c:          []bool{false},
			wantTPR:    []float64{math.NaN(), math.NaN()},
			wantFPR:    []float64{0, 1},
			wantThresh: []float64{math.Inf(1), 1},
		},
		{ // 19
			y:          []float64{1},
			c:          []bool{false},
			cutoffs:    []float64{-1, 1},
			wantTPR:    []float64{math.NaN(), math.NaN()},
			wantFPR:    []float64{1, 1},
			wantThresh: []float64{1, -1},
		},
		{ // 20
			y:          []float64{1},
			c:          []bool{true},
			wantTPR:    []float64{0, 1},
			wantFPR:    []float64{math.NaN(), math.NaN()},
			wantThresh: []float64{math.Inf(1), 1},
		},
		{ // 21
			y:          []float64{},
			c:          []bool{},
			wantTPR:    nil,
			wantFPR:    nil,
			wantThresh: nil,
		},
		{ // 22
			y:          []float64{},
			c:          []bool{},
			cutoffs:    []float64{-1, 2.5, 5, 7.5, 10},
			wantTPR:    nil,
			wantFPR:    nil,
			wantThresh: nil,
		},
		{ // 23
			y:          []float64{0.1, 0.35, 0.4, 0.8},
			c:          []bool{true, false, true, false},
			cutoffs:    []float64{-1, 0.1, 0.35, 0.4, 0.8, 0.9, 1},
			wantTPR:    []float64{0, 0, 0, 0.5, 0.5, 1, 1},
			wantFPR:    []float64{0, 0, 0.5, 0.5, 1, 1, 1},
			wantThresh: []float64{1, 0.9, 0.8, 0.4, 0.35, 0.1, -1},
		},
		{ // 24
			y:          []float64{0.1, 0.35, 0.4, 0.8},
			c:          []bool{true, false, true, false},
			cutoffs:    []float64{math.Inf(-1), 0.1, 0.36, 0.8},
			wantTPR:    []float64{0, 0.5, 1, 1},
			wantFPR:    []float64{0.5, 0.5, 1, 1},
			wantThresh: []float64{0.8, 0.36, 0.1, math.Inf(-1)},
		},
		{ // 25
			y:          []float64{0, 3, 5, 6, 7.5, 8},
			c:          []bool{false, true, false, true, true, true},
			cutoffs:    make([]float64, 0, 10),
			wantTPR:    []float64{0, 0.25, 0.5, 0.75, 0.75, 1, 1},
			wantFPR:    []float64{0, 0, 0, 0, 0.5, 0.5, 1},
			wantThresh: []float64{math.Inf(1), 8, 7.5, 6, 5, 3, 0},
		},
		{ // 26
			y:          []float64{0.1, 0.35, 0.4, 0.8},
			c:          []bool{true, false, true, false},
			cutoffs:    []float64{-1, 0.1, 0.35, 0.4, 0.8, 0.9, 1, 1.1, 1.2},
			wantTPR:    []float64{0, 0, 0, 0, 0, 0.5, 0.5, 1, 1},
			wantFPR:    []float64{0, 0, 0, 0, 0.5, 0.5, 1, 1, 1},
			wantThresh: []float64{1.2, 1.1, 1, 0.9, 0.8, 0.4, 0.35, 0.1, -1},
		},
	}
	for i, test := range cases {
		gotTPR, gotFPR, gotThresh := ROC(test.cutoffs, test.y, test.c, test.w)
		if !floats.Same(gotTPR, test.wantTPR) && !floats.EqualApprox(gotTPR, test.wantTPR, tol) {
			t.Errorf("%d: unexpected TPR got:%v want:%v", i, gotTPR, test.wantTPR)
		}
		if !floats.Same(gotFPR, test.wantFPR) && !floats.EqualApprox(gotFPR, test.wantFPR, tol) {
			t.Errorf("%d: unexpected FPR got:%v want:%v", i, gotFPR, test.wantFPR)
		}
		if !floats.Same(gotThresh, test.wantThresh) {
			t.Errorf("%d: unexpected thresholds got:%#v want:%v", i, gotThresh, test.wantThresh)
		}
	}
}

func TestTOC(t *testing.T) {
	cases := []struct {
		c       []bool
		w       []float64
		wantMin []float64
		wantMax []float64
		wantTOC []float64
	}{
		{ // 0
			// This is the example given in the paper's supplement.
			// http://www2.clarku.edu/~rpontius/TOCexample2.xlsx
			// It is also shown in the WP article.
			// https://en.wikipedia.org/wiki/Total_operating_characteristic#/media/File:TOC_labeled.png
			c: []bool{
				false, false, false, false, false, false,
				false, false, false, false, false, false,
				false, false, true, true, true, true,
				true, true, true, false, false, true,
				false, true, false, false, true, false,
			},
			wantMin: []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			wantMax: []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
			wantTOC: []float64{0, 0, 1, 1, 1, 2, 2, 3, 3, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
		},
		{ // 1
			c:       []bool{},
			wantMin: nil,
			wantMax: nil,
			wantTOC: nil,
		},
		{ // 2
			c: []bool{
				true, true, true, true, true,
			},
			wantMin: []float64{0, 1, 2, 3, 4, 5},
			wantMax: []float64{0, 1, 2, 3, 4, 5},
			wantTOC: []float64{0, 1, 2, 3, 4, 5},
		},
		{ // 3
			c: []bool{
				false, false, false, false, false,
			},
			wantMin: []float64{0, 0, 0, 0, 0, 0},
			wantMax: []float64{0, 0, 0, 0, 0, 0},
			wantTOC: []float64{0, 0, 0, 0, 0, 0},
		},
		{ // 4
			c:       []bool{false, false, false, true, false, true},
			w:       []float64{2, 2, 3, 6, 1, 4},
			wantMin: []float64{0, 0, 0, 3, 6, 8, 10},
			wantMax: []float64{0, 4, 5, 10, 10, 10, 10},
			wantTOC: []float64{0, 4, 4, 10, 10, 10, 10},
		},
	}
	for i, test := range cases {
		gotMin, gotTOC, gotMax := TOC(test.c, test.w)
		if !floats.Same(gotMin, test.wantMin) {
			t.Errorf("%d: unexpected minimum bound got:%v want:%v", i, gotMin, test.wantMin)
		}
		if !floats.Same(gotMax, test.wantMax) {
			t.Errorf("%d: unexpected maximum bound got:%v want:%v", i, gotMax, test.wantMax)
		}
		if !floats.Same(gotTOC, test.wantTOC) {
			t.Errorf("%d: unexpected TOC got:%v want:%v", i, gotTOC, test.wantTOC)
		}
	}
}

func BenchmarkROC(b *testing.B) {
	sizes := []int{empty, small, medium, large}
	for _, cutoffsSize := range sizes {
		for _, ySize := range sizes {
			classesSize := ySize
			for _, weightsSize := range slices.Compact([]int{empty, ySize}) {
				benchmarkROC(b, cutoffsSize, ySize, classesSize, weightsSize)
			}
		}
	}
}

func benchmarkROC(b *testing.B, cutoffsSize int, ySize int, classesSize int, weightsSize int) bool {
	return b.Run(
		fmt.Sprintf(
			"cutoffs=%d,y=%d,classes=%d,weights=%d",
			cutoffsSize, ySize, classesSize, weightsSize),
		func(b *testing.B) {
			src := rand.NewPCG(1, 1)

			cutoffs := randomFloats(cutoffsSize, src)
			slices.Sort(cutoffs)

			y := randomFloats(ySize, src)
			slices.Sort(y)

			classes := randomBools(classesSize, src)

			var weights []float64
			if weightsSize != empty {
				weights = randomFloats(weightsSize, src)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ROC(cutoffs, y, classes, weights)
			}
		})
}

func randomFloats(l int, src rand.Source) []float64 {
	rnd := rand.New(src)
	s := make([]float64, l)
	for i := range s {
		s[i] = rnd.Float64()
	}
	return s
}

func randomBools(l int, src rand.Source) []bool {
	rnd := rand.New(src)
	s := make([]bool, l)
	for i := range s {
		s[i] = rnd.Int32N(2) == 1
	}
	return s
}
