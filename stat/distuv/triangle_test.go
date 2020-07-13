// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"sort"
	"testing"

	"golang.org/x/exp/rand"
)

func TestTriangleConstraint(t *testing.T) {
	t.Parallel()
	for _, test := range []struct{ a, b, c float64 }{
		{a: 1, b: 1, c: 1},
		{a: 1, b: 1, c: 0},
		{a: 1, b: 2, c: 3},
		{a: 1, b: 2, c: 0},
	} {
		if !panics(func() { NewTriangle(test.a, test.b, test.c, nil) }) {
			t.Errorf("expected panic for NewTriangle(%f, %f, %f, nil)", test.a, test.b, test.c)
		}
	}
}

func TestTriangle(t *testing.T) {
	t.Parallel()
	src := rand.New(rand.NewSource(1))
	for i, test := range []struct {
		a, b, c float64
	}{
		{
			a: 0.0,
			b: 1.0,
			c: 0.5,
		},
		{
			a: 0.1,
			b: 0.3,
			c: 0.2,
		},
		{
			a: 1.0,
			b: 2.0,
			c: 1.5,
		},
		{
			a: 0.0,
			b: 1.0,
			c: 0.0,
		},
		{
			a: 0.0,
			b: 1.2,
			c: 1.2,
		},
	} {
		f := NewTriangle(test.a, test.b, test.c, src)
		const (
			tol = 1e-2
			n   = 1e6
		)
		x := make([]float64, n)
		generateSamples(x, f)
		sort.Float64s(x)

		checkMean(t, i, x, f, tol)
		checkVarAndStd(t, i, x, f, tol)
		checkEntropy(t, i, x, f, tol)
		checkExKurtosis(t, i, x, f, tol)
		checkSkewness(t, i, x, f, 5e-2)
		checkMedian(t, i, x, f, tol)
		checkQuantileCDFSurvival(t, i, x, f, tol)
		checkProbContinuous(t, i, x, f.a, f.b, f, 1e-10)
		checkProbQuantContinuous(t, i, x, f, tol)

		if f.c != f.Mode() {
			t.Errorf("Mismatch in mode value: got %v, want %g", f.Mode(), f.c)
		}
	}
}

func TestTriangleProb(t *testing.T) {
	t.Parallel()
	pts := []univariateProbPoint{
		{
			loc:     0.5,
			prob:    0,
			cumProb: 0,
			logProb: math.Inf(-1),
		},
		{
			loc:     1,
			prob:    0,
			cumProb: 0,
			logProb: math.Inf(-1),
		},
		{
			loc:     2,
			prob:    1.0,
			cumProb: 0.5,
			logProb: 0,
		},
		{
			loc:     3,
			prob:    0,
			cumProb: 1,
			logProb: math.Inf(-1),
		},
		{
			loc:     3.5,
			prob:    0,
			cumProb: 1,
			logProb: math.Inf(-1),
		},
	}
	testDistributionProbs(t, NewTriangle(1, 3, 2, nil), "Standard 1,2,3 Triangle", pts)
}

func TestTriangleScore(t *testing.T) {
	const (
		h   = 1e-6
		tol = 1e-6
	)
	t.Parallel()

	f := Triangle{a: -0.5, b: 0.7, c: 0.1}
	testDerivParam(t, &f)

	f = Triangle{a: 0, b: 1, c: 0}
	x := 0.5
	score := f.Score(nil, x)
	if !math.IsNaN(score[0]) {
		t.Errorf("Expected score over A to be NaN for A == C, got %v", score[0])
	}
	if !math.IsNaN(score[2]) {
		t.Errorf("Expected score over C to be NaN for A == C, got %v", score[2])
	}
	expectedScore := logProbDerivative(f, x, 1, h)
	if math.Abs(expectedScore-score[1]) > tol {
		t.Errorf("Mismatch in score over B for A == C: want %g, got %v", expectedScore, score[1])
	}

	f = Triangle{a: 0, b: 1, c: 1}
	score = f.Score(nil, x)
	if !math.IsNaN(score[1]) {
		t.Errorf("Expected score over B to be NaN for B == C, got %v", score[1])
	}
	if !math.IsNaN(score[2]) {
		t.Errorf("Expected score over C to be NaN for B == C, got %v", score[2])
	}
	expectedScore = logProbDerivative(f, x, 0, h)
	if math.Abs(expectedScore-score[0]) > tol {
		t.Errorf("Mismatch in score over A for B == C: want %g, got %v", expectedScore, score[0])
	}

	f = Triangle{a: 0, b: 1, c: 0.5}
	score = f.Score(nil, f.a-0.01)
	if !math.IsNaN(score[0]) {
		t.Errorf("Expected score over B to be NaN for x < A, got %v", score[0])
	}
	if !math.IsNaN(score[1]) {
		t.Errorf("Expected score over B to be NaN for x < A, got %v", score[1])
	}
	if !math.IsNaN(score[2]) {
		t.Errorf("Expected score over C to be NaN for x < A, got %v", score[2])
	}

	score = f.Score(nil, f.b+0.01)
	if !math.IsNaN(score[0]) {
		t.Errorf("Expected score over B to be NaN for x > B, got %v", score[0])
	}
	if !math.IsNaN(score[1]) {
		t.Errorf("Expected score over B to be NaN for x > B, got %v", score[1])
	}
	if !math.IsNaN(score[2]) {
		t.Errorf("Expected score over C to be NaN for x > B, got %v", score[2])
	}

	score = f.Score(nil, f.a)
	if !math.IsNaN(score[0]) {
		t.Errorf("Expected score over C to be NaN for x == A, got %v", score[0])
	}
	score = f.Score(nil, f.b)
	if !math.IsNaN(score[1]) {
		t.Errorf("Expected score over C to be NaN for x == B, got %v", score[1])
	}
	score = f.Score(nil, f.c)
	if !math.IsNaN(score[2]) {
		t.Errorf("Expected score over C to be NaN for x == C, got %v", score[2])
	}
}

func logProbDerivative(t Triangle, x float64, i int, h float64) float64 {
	origParams := t.parameters(nil)
	params := make([]Parameter, len(origParams))
	copy(params, origParams)
	params[i].Value = origParams[i].Value + h
	t.setParameters(params)
	lpUp := t.LogProb(x)
	params[i].Value = origParams[i].Value - h
	t.setParameters(params)
	lpDown := t.LogProb(x)
	t.setParameters(origParams)
	return (lpUp - lpDown) / (2 * h)
}

func TestTriangleScoreInput(t *testing.T) {
	t.Parallel()
	f := Triangle{a: -0.5, b: 0.7, c: 0.1}
	xs := []float64{f.a, f.b, f.c, f.a - 0.0001, f.b + 0.0001}
	for _, x := range xs {
		scoreInput := f.ScoreInput(x)
		if !math.IsNaN(scoreInput) {
			t.Errorf("Expected NaN input score for x == %g, got %v", x, scoreInput)
		}
	}
}
