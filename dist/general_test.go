// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dist

import (
	"fmt"
	"math"
	"testing"
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
	if math.Abs(a-b) > 1e-14 {
		return false
	}
	return true
}

// TODO: Implement a better test for Quantile
func testDistributionProbs(t *testing.T, dist UniProbDist, name string, pts []univariateProbPoint) {
	for _, pt := range pts {
		logProb := dist.LogProb(pt.loc)
		if !absEq(logProb, pt.logProb) {
			t.Errorf("Log probability doesnt match for "+name+". Expected %v. Found %v", pt.logProb, logProb)
		}
		prob := dist.Prob(pt.loc)
		if !absEq(prob, pt.prob) {
			t.Errorf("Probability doesn't match for "+name+". Expected %v. Found %v", pt.prob, prob)
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
			nsInc := incDist.SuffStat(test.samps[j:j+1], incWeights, stats)
			incDist.ConjugateUpdate(stats, nsInc, prior)

			allDist := newFittable()
			nsAll := allDist.SuffStat(test.samps[0:j+1], allWeights, stats)
			allDist.ConjugateUpdate(stats, nsAll, make([]float64, allDist.NumParameters()))
			if !parametersEqual(incDist.parameters(nil), allDist.parameters(nil), 1e-14) {
				t.Errorf("prior doesn't match after incremental update for (%d, %d). Incremental is %v, all at once is %v", i, j, incDist, allDist)
			}

			if test.weights == nil {
				onesDist := newFittable()
				nsOnes := onesDist.SuffStat(test.samps[0:j+1], ones(j+1), stats)
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
}

// rander can generate random samples from a given distribution
type Rander interface {
	Rand() float64
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
