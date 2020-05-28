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
	for i, b := range []Bernoulli{
		{P: 0.5, Src: src},
		{P: 0.9, Src: src},
		{P: 0.2, Src: src},
		{P: 0.0, Src: src},
		{P: 1.0, Src: src},
	} {
		testBernoulli(t, b, i)
		testBernoulliCDF(t, b)
		testBernoulliSurvival(t, b)
		testBernoulliQuantile(t, b)
		if b.P == 0 || b.P == 1 {
			entropy := b.Entropy()
			if entropy != 0 {
				t.Errorf("Entropy of a Bernoulli distribution with P = %g is not zero, got: %g", b.P, entropy)
			}
		}
		if b.NumParameters() != 1 {
			t.Errorf("Wrong number of parameters")
		}
		for _, x := range []float64{-0.2, 0.5, 1.1} {
			logP := b.LogProb(x)
			p := b.Prob(x)
			if !math.IsInf(logP, -1) {
				t.Errorf("Log-probability for x = %g is not -Inf, got: %g", x, logP)
			}
			if p != 0 {
				t.Errorf("Probability for x = %g is not 0, got: %g", x, p)
			}
		}
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
	if b.CDF(0.9999) != 1-b.P {
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
	if !panics(func() { b.Quantile(-0.0001) }) {
		t.Errorf("Expected panic with negative argument")
	}
	if !panics(func() { b.Quantile(1.0001) }) {
		t.Errorf("Expected panic with argument above 1")
	}
	for _, x := range []float64{0., 1.} {
		want := x
		if b.P == 0 {
			want = 0
		}
		if b.Quantile(b.CDF(x)) != want {
			t.Errorf("Quantile(CDF(x)) not equal to %g for x = %g for P = %g", want, x, b.P)
		}
	}
	expectedQuantile1 := 1.
	if b.P == 0 {
		expectedQuantile1 = 0.
	}
	if b.Quantile(1) != expectedQuantile1 {
		t.Errorf("Quantile at 1 not equal to 1 for P = %g", b.P)
	}
	eps := 1e-12
	if b.P > eps && b.P < 1-eps {
		if b.Quantile(1-b.P-eps) != 0 {
			t.Errorf("Quantile slightly below 0 < 1-P < 1 is not zero")
		}
		if b.Quantile(1-b.P+eps) != 1 {
			t.Errorf("Quantile slightly above 0 < 1-P < 1 is not one")
		}
		if b.Quantile(1-b.P) != 0 {
			t.Errorf("Quantile at 0 < 1-P < 1 is not zero")
		}
	}
}
