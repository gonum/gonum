// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"sort"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats/scalar"
)

func TestUniformProb(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		min, max, x, want float64
	}{
		{0, 1, 1, 1},
		{2, 4, 0, 0},
		{2, 4, 5, 0},
		{2, 4, 3, 0.5},
		{0, 100, 1, 0.01},
		{-1, 1, -1.5, 0},
		{-1, 1, 1.5, 0},
	} {
		u := Uniform{test.min, test.max, nil}
		pdf := u.Prob(test.x)
		if !scalar.EqualWithinAbsOrRel(pdf, test.want, 1e-15, 1e-15) {
			t.Errorf("PDF mismatch, x = %v, min = %v, max = %v. Got %v, want %v", test.x, test.min, test.max, pdf, test.want)
		}
		logWant := math.Log(test.want)
		logPdf := u.LogProb(test.x)
		if !scalar.EqualWithinAbsOrRel(logPdf, logWant, 1e-15, 1e-15) {
			t.Errorf("Log PDF mismatch, x = %v, min = %v, max = %v. Got %v, want %v", test.x, test.min, test.max, logPdf, logWant)
		}
	}
}

func TestUniformCDFSurvival(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		min, max, x, want float64
	}{
		{0, 1, 1, 1},
		{0, 100, 100, 1},
		{0, 100, 0, 0},
		{0, 100, 50, 0.5},
		{0, 50, 10, 0.2},
		{-1, 1, -1.5, 0},
		{-1, 1, 1.5, 1},
	} {
		u := Uniform{test.min, test.max, nil}
		cdf := u.CDF(test.x)
		if !scalar.EqualWithinAbsOrRel(cdf, test.want, 1e-15, 1e-15) {
			t.Errorf("CDF mismatch, x = %v, min = %v, max = %v. Got %v, want %v", test.x, test.min, test.max, cdf, test.want)
		}
		survival := u.Survival(test.x)
		if !scalar.EqualWithinAbsOrRel(survival, 1-test.want, 1e-15, 1e-15) {
			t.Errorf("CDF mismatch, x = %v, min = %v, max = %v. Got %v, want %v", test.x, test.min, test.max, survival, 1-test.want)
		}
	}
}

func TestUniform(t *testing.T) {
	t.Parallel()
	src := rand.New(rand.NewSource(1))
	for i, b := range []Uniform{
		{1, 2, src},
		{0, 100, src},
		{50, 60, src},
	} {
		testUniform(t, b, i)
	}
}

func testUniform(t *testing.T, u Uniform, i int) {
	const (
		tol  = 1e-2
		n    = 1e5
		bins = 50
	)
	x := make([]float64, n)
	generateSamples(x, u)
	sort.Float64s(x)

	testRandLogProbContinuous(t, i, 0, x, u, tol, bins)
	checkMean(t, i, x, u, tol)
	checkVarAndStd(t, i, x, u, tol)
	checkExKurtosis(t, i, x, u, 7e-2)
	checkProbContinuous(t, i, x, u.Min, u.Max, u, 1e-10)
	checkQuantileCDFSurvival(t, i, x, u, 1e-2)
	checkEntropy(t, i, x, u, tol)
	checkSkewness(t, i, x, u, tol)
	checkMedian(t, i, x, u, tol)
	testDerivParam(t, &u)
}

func TestUniformScore(t *testing.T) {
	t.Parallel()
	u := Uniform{0, 1, nil}
	for _, test := range []struct {
		x, wantMin, wantMax float64
	}{
		{-0.001, math.NaN(), math.NaN()},
		{0, math.NaN(), -1},
		{1, 1, math.NaN()},
		{1.001, math.NaN(), math.NaN()},
	} {
		score := u.Score(nil, test.x)
		if !scalar.Same(score[0], test.wantMin) {
			t.Errorf("Score[0] mismatch for at %g: got %v, want %g", test.x, score[0], test.wantMin)
		}
		if !scalar.Same(score[1], test.wantMax) {
			t.Errorf("Score[1] mismatch for at %g: got %v, want %g", test.x, score[1], test.wantMax)
		}
	}
}

func TestUniformScoreInput(t *testing.T) {
	t.Parallel()
	u := Uniform{0, 1, nil}
	scoreInput := u.ScoreInput(0.5)
	if scoreInput != 0 {
		t.Errorf("Mismatch in input score for U(0, 1) at x == 0.5: got %v, want 0", scoreInput)
	}
	xs := []float64{-0.0001, 0, 1, 1.0001}
	for _, x := range xs {
		scoreInput = u.ScoreInput(x)
		if !math.IsNaN(scoreInput) {
			t.Errorf("Expected NaN score input for U(0, 1) at x == %g, got %v", x, scoreInput)
		}
	}
}
