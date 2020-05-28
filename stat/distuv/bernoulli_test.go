// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
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
		{P: 0.0, Src: src},
		{P: 1.0, Src: src},
	} {
		testBernoulli(t, dist, i)
		testBernoulliCDF(t, dist)
		testBernoulliSurvival(t, dist)
		testBernoulliQuantile(t, dist)
	}
}

func testBernoulli(t *testing.T, b Bernoulli, i int) {
	const (
		tol  = 1e-2
		n    = 3e6
		bins = 50
	)
	x := make([]float64, n)
	generateSamples(x, b)
	sort.Float64s(x)

	checkMean(t, i, x, b, tol)
	checkVarAndStd(t, i, x, b, tol)
	checkEntropy(t, i, x, b, tol)
	checkProbDiscrete(t, i, x, b, tol)
	if b.P != 0 && b.P != 1 {
		// Sample kurtosis and skewness are going to be NaN for P = 0 or 1.
		checkExKurtosis(t, i, x, b, tol)
		checkSkewness(t, i, x, b, tol)
	} else {
		if !math.IsInf(b.ExKurtosis(), 1) {
			t.Errorf("Excess kurtosis for P == 0 or 1 is not +Inf")
		}
		skewness := b.Skewness()
		if b.P == 0 {
			if !math.IsInf(skewness, 1) {
				t.Errorf("Skewness for P == 0 is not +Inf")
			}
		} else {
			if !math.IsInf(skewness, -1) {
				t.Errorf("Skewness for P == 1 is not -Inf")
			}
		}
	}
	if b.P != 0.5 {
		checkMedian(t, i, x, b, tol)
	} else if b.Median() != 0.5 {
		t.Errorf("Median for P == 0.5 is not 0.5")
	}
}

func testBernoulliCDF(t *testing.T, b Bernoulli) {
	if b.CDF(-0.000001) != 0 {
		t.Errorf("Bernoulli CDF below zero is not zero")
	}
	if b.CDF(0) != 1-b.P {
		t.Errorf("Bernoulli CDF at zero is not 1 - P(1)")
	}
	if b.CDF(0.0001) != 1-b.P {
		t.Errorf("Bernoulli CDF between zero and one is not 1 - P(1)")
	}
	if b.CDF(1) != 1 {
		t.Errorf("Bernoulli CDF at one is not one")
	}
	if b.CDF(1.00001) != 1 {
		t.Errorf("Bernoulli CDF above one is not one")
	}
}

func testBernoulliSurvival(t *testing.T, b Bernoulli) {
	if b.Survival(-0.000001) != 1 {
		t.Errorf("Bernoulli Survival below zero is not one")
	}
	if b.Survival(0) != b.P {
		t.Errorf("Bernoulli Survival at zero is not P(1)")
	}
	if b.Survival(0.0001) != b.P {
		t.Errorf("Bernoulli Survival between zero and one is not P(1)")
	}
	if b.Survival(1) != 0 {
		t.Errorf("Bernoulli Survival at one is not zero")
	}
	if b.Survival(1.00001) != 0 {
		t.Errorf("Bernoulli Survival above one is not zero")
	}
}

func testBernoulliQuantile(t *testing.T, b Bernoulli) {
	for _, x := range []float64{0., 1.} {
		if b.Quantile(b.CDF(x)) != x {
			t.Errorf("Quantile(CDF(x)) not equal to x for x = %g for P = %g", x, b.P)
		}
	}
	if b.Quantile(1) != 1 {
		t.Errorf("Quantile at 1 not equal to 1")
	}
}

func TestBernoulliEntropySpecial(t *testing.T) {
	src := rand.New(rand.NewSource(1))
	for _, p := range []float64{0, 1} {
		b := Bernoulli{p, src}
		entropy := b.Entropy()
		if entropy != 0 {
			t.Errorf("Entropy of a Bernoulli distribution with P = %g is not zero, got: %g", p, entropy)
		}
	}
}
