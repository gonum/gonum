// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"fmt"
	"math"
	"testing"

	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/floats"
)

type univariateProbPoint struct {
	loc     float64
	logProb float64
	cumProb float64
	prob    float64
}

type UniProbDist interface {
	Prob(float64) float64
	CDF(float64) float64
	LogProb(float64) float64
	Quantile(float64) float64
	Survival(float64) float64
}

func absEq(a, b float64) bool {
	return absEqTol(a, b, 1e-14)
}

func absEqTol(a, b, tol float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		// NaN is not equal to anything.
		return false
	}
	// This is expressed as the inverse to catch the
	// case a = Inf and b = Inf of the same sign.
	return !(math.Abs(a-b) > tol)
}

// TODO: Implement a better test for Quantile
func testDistributionProbs(t *testing.T, dist UniProbDist, name string, pts []univariateProbPoint) {
	for _, pt := range pts {
		logProb := dist.LogProb(pt.loc)
		if !absEq(logProb, pt.logProb) {
			t.Errorf("Log probability doesnt match for "+name+" at %v. Expected %v. Found %v", pt.loc, pt.logProb, logProb)
		}
		prob := dist.Prob(pt.loc)
		if !absEq(prob, pt.prob) {
			t.Errorf("Probability doesn't match for "+name+" at %v. Expected %v. Found %v", pt.loc, pt.prob, prob)
		}
		cumProb := dist.CDF(pt.loc)
		if !absEq(cumProb, pt.cumProb) {
			t.Errorf("Cumulative Probability doesn't match for "+name+". Expected %v. Found %v", pt.cumProb, cumProb)
		}
		if !absEq(dist.Survival(pt.loc), 1-pt.cumProb) {
			t.Errorf("Survival doesn't match for %v. Expected %v, Found %v", name, 1-pt.cumProb, dist.Survival(pt.loc))
		}
		if pt.prob != 0 {
			if math.Abs(dist.Quantile(pt.cumProb)-pt.loc) > 1e-4 {
				fmt.Println("true =", pt.loc)
				fmt.Println("calculated=", dist.Quantile(pt.cumProb))
				t.Errorf("Quantile doesn't match for "+name+", loc =  %v", pt.loc)
			}
		}
	}
}

type ConjugateUpdater interface {
	NumParameters() int
	parameters([]Parameter) []Parameter

	NumSuffStat() int
	SuffStat([]float64, []float64, []float64) float64
	ConjugateUpdate([]float64, float64, []float64)

	Rand() float64
}

func testConjugateUpdate(t *testing.T, newFittable func() ConjugateUpdater) {
	for i, test := range []struct {
		samps   []float64
		weights []float64
	}{
		{
			samps:   randn(newFittable(), 10),
			weights: nil,
		},
		{
			samps:   randn(newFittable(), 10),
			weights: ones(10),
		},
		{
			samps:   randn(newFittable(), 10),
			weights: randn(&Exponential{Rate: 1}, 10),
		},
	} {
		// ensure that conjugate produces the same result both incrementally and all at once
		incDist := newFittable()
		stats := make([]float64, incDist.NumSuffStat())
		prior := make([]float64, incDist.NumParameters())
		for j := range test.samps {
			var incWeights, allWeights []float64
			if test.weights != nil {
				incWeights = test.weights[j : j+1]
				allWeights = test.weights[0 : j+1]
			}
			nsInc := incDist.SuffStat(stats, test.samps[j:j+1], incWeights)
			incDist.ConjugateUpdate(stats, nsInc, prior)

			allDist := newFittable()
			nsAll := allDist.SuffStat(stats, test.samps[0:j+1], allWeights)
			allDist.ConjugateUpdate(stats, nsAll, make([]float64, allDist.NumParameters()))
			if !parametersEqual(incDist.parameters(nil), allDist.parameters(nil), 1e-12) {
				t.Errorf("prior doesn't match after incremental update for (%d, %d). Incremental is %v, all at once is %v", i, j, incDist, allDist)
			}

			if test.weights == nil {
				onesDist := newFittable()
				nsOnes := onesDist.SuffStat(stats, test.samps[0:j+1], ones(j+1))
				onesDist.ConjugateUpdate(stats, nsOnes, make([]float64, onesDist.NumParameters()))
				if !parametersEqual(onesDist.parameters(nil), incDist.parameters(nil), 1e-14) {
					t.Errorf("nil and uniform weighted prior doesn't match for incremental update for (%d, %d). Uniform weighted is %v, nil is %v", i, j, onesDist, incDist)
				}
				if !parametersEqual(onesDist.parameters(nil), allDist.parameters(nil), 1e-14) {
					t.Errorf("nil and uniform weighted prior doesn't match for all at once update for (%d, %d). Uniform weighted is %v, nil is %v", i, j, onesDist, incDist)
				}
			}
		}
	}
	testSuffStatPanics(t, newFittable)
	testConjugateUpdatePanics(t, newFittable)
}

func testSuffStatPanics(t *testing.T, newFittable func() ConjugateUpdater) {
	dist := newFittable()
	sample := randn(dist, 10)
	if !panics(func() { dist.SuffStat(make([]float64, dist.NumSuffStat()), sample, make([]float64, len(sample)+1)) }) {
		t.Errorf("Expected panic for mismatch between samples and weights lengths")
	}
	if !panics(func() { dist.SuffStat(make([]float64, dist.NumSuffStat()+1), sample, nil) }) {
		t.Errorf("Expected panic for wrong sufficient statistic length")
	}
}

func testConjugateUpdatePanics(t *testing.T, newFittable func() ConjugateUpdater) {
	dist := newFittable()
	if !panics(func() {
		dist.ConjugateUpdate(make([]float64, dist.NumSuffStat()+1), 100, make([]float64, dist.NumParameters()))
	}) {
		t.Errorf("Expected panic for wrong sufficient statistic length")
	}
	if !panics(func() {
		dist.ConjugateUpdate(make([]float64, dist.NumSuffStat()), 100, make([]float64, dist.NumParameters()+1))
	}) {
		t.Errorf("Expected panic for wrong prior strength length")
	}
}

// randn generates a specified number of random samples
func randn(dist Rander, n int) []float64 {
	x := make([]float64, n)
	for i := range x {
		x[i] = dist.Rand()
	}
	return x
}

func ones(n int) []float64 {
	x := make([]float64, n)
	for i := range x {
		x[i] = 1
	}
	return x
}

func parametersEqual(p1, p2 []Parameter, tol float64) bool {
	for i, p := range p1 {
		if p.Name != p2[i].Name {
			return false
		}
		if math.Abs(p.Value-p2[i].Value) > tol {
			return false
		}
	}
	return true
}

type derivParamTester interface {
	LogProb(x float64) float64
	Score(deriv []float64, x float64) []float64
	ScoreInput(x float64) float64
	Quantile(p float64) float64
	NumParameters() int
	parameters([]Parameter) []Parameter
	setParameters([]Parameter)
}

func testDerivParam(t *testing.T, d derivParamTester) {
	// Tests that the derivative matches for a number of different quantiles
	// along the distribution.
	nTest := 10
	quantiles := make([]float64, nTest)
	floats.Span(quantiles, 0.1, 0.9)

	scoreInPlace := make([]float64, d.NumParameters())
	fdDerivParam := make([]float64, d.NumParameters())

	if !panics(func() { d.Score(make([]float64, d.NumParameters()+1), 0) }) {
		t.Errorf("Expected panic for wrong derivative slice length")
	}
	if !panics(func() { d.parameters(make([]Parameter, d.NumParameters()+1)) }) {
		t.Errorf("Expected panic for wrong parameter slice length")
	}

	initParams := d.parameters(nil)
	tooLongParams := make([]Parameter, len(initParams)+1)
	copy(tooLongParams, initParams)
	if !panics(func() { d.setParameters(tooLongParams) }) {
		t.Errorf("Expected panic for wrong parameter slice length")
	}
	badNameParams := make([]Parameter, len(initParams))
	copy(badNameParams, initParams)
	const badName = "__badName__"
	for i := 0; i < len(initParams); i++ {
		badNameParams[i].Name = badName
		if !panics(func() { d.setParameters(badNameParams) }) {
			t.Errorf("Expected panic for wrong %d-th parameter name", i)
		}
		badNameParams[i].Name = initParams[i].Name
	}

	init := make([]float64, d.NumParameters())
	for i, v := range initParams {
		init[i] = v.Value
	}
	for _, v := range quantiles {
		d.setParameters(initParams)
		x := d.Quantile(v)
		score := d.Score(scoreInPlace, x)
		if &score[0] != &scoreInPlace[0] {
			t.Errorf("Returned a different derivative slice than passed in. Got %v, want %v", score, scoreInPlace)
		}
		logProbParams := func(p []float64) float64 {
			params := d.parameters(nil)
			for i, v := range p {
				params[i].Value = v
			}
			d.setParameters(params)
			return d.LogProb(x)
		}
		fd.Gradient(fdDerivParam, logProbParams, init, nil)
		if !floats.EqualApprox(scoreInPlace, fdDerivParam, 1e-6) {
			t.Errorf("Score mismatch at x = %g. Want %v, got %v", x, fdDerivParam, scoreInPlace)
		}
		d.setParameters(initParams)
		score2 := d.Score(nil, x)
		if !floats.EqualApprox(score2, scoreInPlace, 1e-14) {
			t.Errorf("Score mismatch when input nil Want %v, got %v", score2, scoreInPlace)
		}
		logProbInput := func(x2 float64) float64 {
			return d.LogProb(x2)
		}
		scoreInput := d.ScoreInput(x)
		fdDerivInput := fd.Derivative(logProbInput, x, nil)
		if !absEqTol(scoreInput, fdDerivInput, 1e-6) {
			t.Errorf("ScoreInput mismatch at x = %g. Want %v, got %v", x, fdDerivInput, scoreInput)
		}
	}
}
