// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dist

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/floats"
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
		if !absEq(dist.CDF(pt.loc), pt.cumProb) {
			t.Errorf("Cumulative Probability doesn't match for " + name)
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

type PriorFittable interface {
	FitPrior(samples, weights, priorValue, priorWeight []float64) ([]float64, []float64)
	MarshalSlice([]float64)
	Rand() float64
	NumParameters() int
}

func testFitPrior(t *testing.T, dist PriorFittable, newFittable func() PriorFittable) {
	samps := make([]float64, 10)
	for i := range samps {
		samps[i] = dist.Rand()
	}

	n2 := newFittable()
	p2, w2 := n2.FitPrior(samps, nil, nil, nil)
	n3 := newFittable()
	p3, w3 := n3.FitPrior(samps[:7], nil, nil, nil)
	n4 := newFittable()
	p4, w4 := n4.FitPrior(samps[7:], nil, p3, w3)

	params2 := make([]float64, n2.NumParameters())
	n2.MarshalSlice(params2)
	if len(params2) == 0 {
		panic("len 0 params")
	}
	params4 := make([]float64, n4.NumParameters())
	n4.MarshalSlice(params4)

	if !floats.EqualApprox(params2, params4, 1e-14) {
		t.Errorf("parameters don't match: First is %v, second is %v", params2, params4)
	}
	if !floats.EqualApprox(p2, p4, 1e-14) {
		t.Errorf("prior doesn't match after two step update. First is %v, second is %v", p2, p4)
	}
	if !floats.EqualApprox(w2, w4, 1e-14) {
		t.Errorf("prior weight doesn't match after two step update. First is %v, second is %v", w2, w4)
	}

	// Try with weights = 1
	n5 := newFittable()
	ones := make([]float64, len(samps))
	for i := range ones {
		ones[i] = 1
	}
	p5, w5 := n5.FitPrior(samps, ones, nil, nil)
	params5 := make([]float64, n5.NumParameters())
	n5.MarshalSlice(params5)
	if !floats.EqualApprox(params2, params5, 1e-14) {
		t.Errorf("parameters don't match: First is %v, second is %v", params2, params5)
	}
	if !floats.EqualApprox(p2, p5, 1e-14) {
		t.Errorf("prior doesn't match after unitary weights. First is %v, second is %v", p2, p5)
	}
	if !floats.EqualApprox(w2, w5, 1e-14) {
		t.Errorf("prior weight doesn't match unitary weights. First is %v, second is %v", w2, w5)
	}

	// Lastly, make sure it's okay with a bunch of random weights
	weights := make([]float64, len(samps))
	for i := range weights {
		weights[i] = rand.Float64()
	}
	n6 := newFittable()
	p6, w6 := n6.FitPrior(samps, weights, nil, nil)
	n7 := newFittable()
	p7, w7 := n7.FitPrior(samps[:7], weights[:7], nil, nil)
	n8 := newFittable()
	p8, w8 := n8.FitPrior(samps[7:], weights[7:], p7, w7)

	params6 := make([]float64, n6.NumParameters())
	n6.MarshalSlice(params6)
	params8 := make([]float64, n8.NumParameters())
	n8.MarshalSlice(params8)

	if !floats.EqualApprox(params6, params8, 1e-14) {
		t.Errorf("parameters don't match: First is %v, second is %v", params6, params8)
	}
	if !floats.EqualApprox(p6, p8, 1e-14) {
		t.Errorf("prior doesn't match after two step update. First is %v, second is %v", p6, p8)
	}
	if !floats.EqualApprox(w6, w8, 1e-14) {
		t.Errorf("prior weight doesn't match after two step update. First is %v, second is %v", w6, w8)
	}
}
