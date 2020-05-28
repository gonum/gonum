// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"sort"
	"testing"

	"golang.org/x/exp/rand"
)

func TestBernoulli(t *testing.T) {
	t.Parallel()
	src := rand.New(rand.NewSource(1))
	for i, dist := range []Bernoulli{
		{P: 0.5, Src: src},
		{P: 0.9, Src: src},
		{P: 0.2, Src: src},
	} {
		testBernoulli(t, dist, i)
		testBernoulliCDF(t, dist)
	}
}

func testBernoulli(t *testing.T, dist Bernoulli, i int) {
	const (
		tol  = 1e-2
		n    = 3e6
		bins = 50
	)
	x := make([]float64, n)
	generateSamples(x, dist)
	sort.Float64s(x)

	checkMean(t, i, x, dist, tol)
	checkVarAndStd(t, i, x, dist, tol)
	checkEntropy(t, i, x, dist, tol)
	checkExKurtosis(t, i, x, dist, tol)
	checkSkewness(t, i, x, dist, tol)
	checkProbDiscrete(t, i, x, dist, tol)
}

func testBernoulliCDF(t *testing.T, dist Bernoulli) {
	if dist.CDF(-0.000001) != 0 {
		t.Errorf("Bernoulli CDF below zero is not zero")
	}
	if dist.CDF(0) != 1-dist.P {
		t.Errorf("Bernoulli CDF at zero is not 1 - P(1)")
	}
	if dist.CDF(0.0001) != 1-dist.P {
		t.Errorf("Bernoulli CDF between zero and one is not 1 - P(1)")
	}
	if dist.CDF(1) != 1 {
		t.Errorf("Bernoulli CDF at one is not one")
	}
	if dist.CDF(1.00001) != 1 {
		t.Errorf("Bernoulli CDF above one is not one")
	}
}
