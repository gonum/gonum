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
		if dist.P == 0 || dist.P == 1 {
			entropy := dist.Entropy()
			if entropy != 0 {
				t.Errorf("Entropy of a Bernoulli distribution with P = %g is not zero, got: %g", dist.P, entropy)
			}
		}
		if dist.NumParameters() != 1 {
			t.Errorf("Wrong number of parameters")
		}
		for _, x := range []float64{-0.2, 0.5, 1.1} {
			logP := dist.LogProb(x)
			p := dist.Prob(x)
			if !math.IsInf(logP, -1) {
				t.Errorf("Log-probability for x = %g is not -Inf, got: %g", x, logP)
			}
			if p != 0 {
				t.Errorf("Probability for x = %g is not 0, got: %g", x, p)
			}
		}
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
	checkProbDiscrete(t, i, x, dist, tol)
	if dist.P != 0 && dist.P != 1 {
		// Sample kurtosis and skewness are going to be NaN for P = 0 or 1.
		checkExKurtosis(t, i, x, dist, tol)
		checkSkewness(t, i, x, dist, tol)
	} else {
		if !math.IsInf(dist.ExKurtosis(), 1) {
			t.Errorf("Excess kurtosis for P == 0 or 1 is not +Inf")
		}
		skewness := dist.Skewness()
		if dist.P == 0 {
			if !math.IsInf(skewness, 1) {
				t.Errorf("Skewness for P == 0 is not +Inf")
			}
		} else {
			if !math.IsInf(skewness, -1) {
				t.Errorf("Skewness for P == 1 is not -Inf")
			}
		}
	}
	if dist.P != 0.5 {
		checkMedian(t, i, x, dist, tol)
	} else if dist.Median() != 0.5 {
		t.Errorf("Median for P == 0.5 is not 0.5")
	}
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
	if dist.CDF(0.9999) != 1-dist.P {
		t.Errorf("Bernoulli CDF between zero and one is not 1 - P(1)")
	}
	if dist.CDF(1) != 1 {
		t.Errorf("Bernoulli CDF at one is not one")
	}
	if dist.CDF(1.00001) != 1 {
		t.Errorf("Bernoulli CDF above one is not one")
	}
}

func testBernoulliSurvival(t *testing.T, dist Bernoulli) {
	if dist.Survival(-0.000001) != 1 {
		t.Errorf("Bernoulli Survival below zero is not one")
	}
	if dist.Survival(0) != dist.P {
		t.Errorf("Bernoulli Survival at zero is not P(1)")
	}
	if dist.Survival(0.0001) != dist.P {
		t.Errorf("Bernoulli Survival between zero and one is not P(1)")
	}
	if dist.Survival(1) != 0 {
		t.Errorf("Bernoulli Survival at one is not zero")
	}
	if dist.Survival(1.00001) != 0 {
		t.Errorf("Bernoulli Survival above one is not zero")
	}
}

func testBernoulliQuantile(t *testing.T, dist Bernoulli) {
	if !panics(func() { dist.Quantile(-0.0001) }) {
		t.Errorf("Expected panic with negative argument")
	}
	if !panics(func() { dist.Quantile(1.0001) }) {
		t.Errorf("Expected panic with argument above 1")
	}
	for _, x := range []float64{0., 1.} {
		want := x
		if dist.P == 0 {
			want = 0
		}
		if dist.Quantile(dist.CDF(x)) != want {
			t.Errorf("Quantile(CDF(x)) not equal to %g for x = %g for P = %g", want, x, dist.P)
		}
	}
	expectedQuantile1 := 1.
	if dist.P == 0 {
		expectedQuantile1 = 0.
	}
	if dist.Quantile(1) != expectedQuantile1 {
		t.Errorf("Quantile at 1 not equal to 1 for P = %g", dist.P)
	}
	eps := 1e-12
	if dist.P > eps && dist.P < 1-eps {
		if dist.Quantile(1-dist.P-eps) != 0 {
			t.Errorf("Quantile slightly below 0 < 1-P < 1 is not zero")
		}
		if dist.Quantile(1-dist.P+eps) != 1 {
			t.Errorf("Quantile slightly above 0 < 1-P < 1 is not one")
		}
		if dist.Quantile(1-dist.P) != 0 {
			t.Errorf("Quantile at 0 < 1-P < 1 is not zero")
		}
	}
}
