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

type ConjugateUpdater interface {
	MarshalParameters([]Parameter)
	UnmarshalParameters([]Parameter)
	Rand() float64
	NumParameters() int
	SuffStat([]float64, []float64, []float64) float64
	ConjugateUpdate([]float64, float64, []float64)
	NumSuffStat() int
}

func testConjugateUpdate(t *testing.T, dist ConjugateUpdater, newFittable func() ConjugateUpdater) {
	samps := make([]float64, 10)
	for i := range samps {
		samps[i] = dist.Rand()
	}
	nParams := dist.NumParameters()
	nSuffStat := dist.NumSuffStat()

	stats2 := make([]float64, nSuffStat)
	p2 := make([]Parameter, nParams)
	w2 := make([]float64, nParams)
	n2 := newFittable()
	ns2 := n2.SuffStat(samps, nil, stats2)
	n2.ConjugateUpdate(stats2, ns2, w2)
	n2.MarshalParameters(p2)

	n3 := newFittable()
	p3 := make([]Parameter, nParams)
	w3 := make([]float64, nParams)
	stats3 := make([]float64, nSuffStat)
	ns3 := n3.SuffStat(samps[:7], nil, stats3)
	n3.ConjugateUpdate(stats3, ns3, w3)
	n3.MarshalParameters(p3)

	n4 := newFittable()
	n4.UnmarshalParameters(p3)
	p4 := make([]Parameter, nParams)
	w4 := make([]float64, nParams)
	stats4 := make([]float64, nSuffStat)
	copy(w4, w3)
	ns4 := n4.SuffStat(samps[7:], nil, stats4)
	n4.ConjugateUpdate(stats4, ns4, w4)
	n4.MarshalParameters(p4)

	if !parametersEqual(p2, p4, 1e-14) {
		t.Errorf("prior doesn't match after two step update. First is %v, second is %v", p2, p4)
	}
	if !floats.EqualApprox(w2, w4, 1e-14) {
		t.Errorf("prior weight doesn't match after two step update. First is %v, second is %v", w2, w4)
	}

	// Try with weights = 1
	ones := make([]float64, len(samps))
	for i := range ones {
		ones[i] = 1
	}

	n5 := newFittable()

	p5 := make([]Parameter, nParams)
	w5 := make([]float64, nParams)
	stats5 := make([]float64, nSuffStat)
	ns5 := n5.SuffStat(samps, ones, stats5)
	n5.ConjugateUpdate(stats5, ns5, w5)
	n5.MarshalParameters(p5)

	if !parametersEqual(p2, p5, 1e-14) {
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

	p6 := make([]Parameter, nParams)
	w6 := make([]float64, nParams)
	n6 := newFittable()
	stats6 := make([]float64, nSuffStat)
	ns6 := n6.SuffStat(samps, weights, stats6)
	n6.ConjugateUpdate(stats6, ns6, w6)
	n6.MarshalParameters(p6)

	p7 := make([]Parameter, nParams)
	w7 := make([]float64, nParams)
	n7 := newFittable()
	stats7 := make([]float64, nSuffStat)
	ns7 := n7.SuffStat(samps[:7], weights[:7], stats7)
	n7.ConjugateUpdate(stats7, ns7, w7)
	n7.MarshalParameters(p7)

	p8 := make([]Parameter, nParams)
	w8 := make([]float64, nParams)
	n8 := newFittable()
	n8.UnmarshalParameters(p7)
	stats8 := make([]float64, nSuffStat)
	ns8 := n7.SuffStat(samps[7:], weights[7:], stats8)
	copy(w8, w7)
	n8.ConjugateUpdate(stats8, ns8, w8)
	n8.MarshalParameters(p8)

	if !parametersEqual(p6, p8, 1e-14) {
		t.Errorf("prior doesn't match after two step update. First is %v, second is %v", p6, p8)
	}
	if !floats.EqualApprox(w6, w8, 1e-14) {
		t.Errorf("prior weight doesn't match after two step update. First is %v, second is %v", w6, w8)
	}
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
