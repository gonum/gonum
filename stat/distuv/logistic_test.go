// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats/scalar"
)

func TestLogisticParameters(t *testing.T) {
	t.Parallel()

	var want float64

	l := Logistic{Mu: 1, S: 0}

	want = 2
	if result := l.NumParameters(); result != int(want) {
		t.Errorf("Wrong number of parameters: %d != %.0f", result, want)
	}

	want = 6.0 / 5.0
	if result := l.ExKurtosis(); result != want {
		t.Errorf("Wrong excess kurtosis: %f != %f", result, want)
	}

	want = 0.0
	if result := l.Skewness(); result != want {
		t.Errorf("Wrong skewness: %f != %f", result, want)
	}

	want = l.Mu
	if result := l.Mean(); result != want {
		t.Errorf("Wrong mean value: %f != %f", result, want)
	}

	want = l.Mu
	if result := l.Median(); result != want {
		t.Errorf("Wrong median value: %f != %f", result, want)
	}

	want = l.Mu
	if result := l.Mode(); result != want {
		t.Errorf("Wrong mode value: %f != %f", result, want)
	}
}

func TestLogisticStdDev(t *testing.T) {
	t.Parallel()

	l := Logistic{Mu: 0, S: sqrt3 / math.Pi}

	want := 1.0
	if result := l.StdDev(); !scalar.EqualWithinAbs(result, want, 1e-10) {
		t.Errorf("Wrong StdDev with Mu=%f, S=%f: %f != %f", l.Mu, l.S, result, want)
	}

	want = 1.0
	if result := l.Variance(); !scalar.EqualWithinAbs(result, want, 1e-10) {
		t.Errorf("Wrong Variance with Mu=%f, S=%f: %f != %f", l.Mu, l.S, result, want)
	}
}

func TestLogisticCDF(t *testing.T) {
	t.Parallel()

	// Values for "want" are taken from WolframAlpha: CDF[LogisticDistribution[mu,s], input] to 10 digits.
	for _, v := range []struct {
		mu, s, input, want float64
	}{
		{0.0, 0.0, 1.0, 1.0},
		{0.0, 1.0, 0.0, 0.5},
		{-0.5, 1.0, 0.0, 0.6224593312},
		{69.0, 420.0, 42.0, 0.4839341039},
	} {
		l := Logistic{Mu: v.mu, S: v.s}
		if result := l.CDF(v.input); !scalar.EqualWithinAbs(result, v.want, 1e-10) {
			t.Errorf("Wrong CDF(%f) with Mu=%f, S=%f: %f != %f", v.input, l.Mu, l.S, result, v.want)
		}
	}

	// Edge case of zero in denominator.
	l := Logistic{Mu: 0, S: 0}

	input := 0.0
	if result := l.CDF(input); !math.IsNaN(result) {
		t.Errorf("Wrong CDF(%f) with Mu=%f, S=%f: %f is not NaN", input, l.Mu, l.S, result)
	}
}

// TestLogisticSurvival doesn't need excessive testing since it's just 1-CDF.
func TestLogisticSurvival(t *testing.T) {
	t.Parallel()

	l := Logistic{Mu: 0, S: 1}

	input, want := 0.0, 0.5
	if result := l.Survival(input); result != want {
		t.Errorf("Wrong Survival(%f) with Mu=%f, S=%f: %f != %f", input, l.Mu, l.S, result, want)
	}
}

func TestLogisticProb(t *testing.T) {
	t.Parallel()

	// Values for "want" are taken from WolframAlpha: PDF[LogisticDistribution[mu,s], input] to 10 digits.
	for _, v := range []struct {
		mu, s, input, want float64
	}{
		{0.0, 1.0, 0.0, 0.25},
		{-0.5, 1.0, 0.0, 0.2350037122},
		{69.0, 420.0, 42.0, 0.0005946235404},
	} {
		l := Logistic{Mu: v.mu, S: v.s}
		if result := l.Prob(v.input); !scalar.EqualWithinAbs(result, v.want, 1e-10) {
			t.Errorf("Wrong Prob(%f) with Mu=%f, S=%f: %.09f != %.09f", v.input, l.Mu, l.S, result, v.want)
		}
	}

	// Edge case of zero in denominator.
	l := Logistic{Mu: 0, S: 0}

	input := 0.0
	if result := l.Prob(input); !math.IsNaN(result) {
		t.Errorf("Wrong Prob(%f) with Mu=%f, S=%f: %f is not NaN", input, l.Mu, l.S, result)
	}

	input = 1.0
	if result := l.Prob(input); !math.IsNaN(result) {
		t.Errorf("Wrong Prob(%f) with Mu=%f, S=%f: %f is not NaN", input, l.Mu, l.S, result)
	}
}

func TestLogisticLogProb(t *testing.T) {
	t.Parallel()

	l := Logistic{Mu: 0, S: 1}

	input, want := 0.0, -math.Log(4)
	if result := l.LogProb(input); result != want {
		t.Errorf("Wrong LogProb(%f) with Mu=%f, S=%f: %f != %f", input, l.Mu, l.S, result, want)
	}

	// At the mean the density is 1/(4*S), so the log density is -log(4*S)
	// for any location and scale.
	for _, test := range []struct {
		dist Logistic
		want float64
	}{
		{dist: Logistic{Mu: 0, S: 3}, want: -2.4849066497880004},   // -log(12)
		{dist: Logistic{Mu: 2, S: 0.5}, want: -0.6931471805599453}, // -log(2)
	} {
		if got := test.dist.LogProb(test.dist.Mu); !scalar.EqualWithinAbsOrRel(got, test.want, 1e-14, 1e-14) {
			t.Errorf("Wrong LogProb(%g) with Mu=%g, S=%g: got %g, want %g", test.dist.Mu, test.dist.Mu, test.dist.S, got, test.want)
		}
	}

	// Check the log density against the density away from the far tails,
	// where Prob is accurate.
	for _, dist := range []Logistic{
		{Mu: 0, S: 1},
		{Mu: 2, S: 1},
		{Mu: 0, S: 3},
		{Mu: 3, S: 630},
		{Mu: -5, S: 0.125},
	} {
		for _, x := range []float64{-8, -1, 0, 0.5, 2, 10} {
			want := math.Log(dist.Prob(x))
			if got := dist.LogProb(x); !scalar.EqualWithinAbsOrRel(got, want, 1e-14, 1e-14) {
				t.Errorf("Wrong LogProb(%g) with Mu=%g, S=%g: got %g, want %g", x, dist.Mu, dist.S, got, want)
			}
		}
	}
}

func TestLogisticLogProbTails(t *testing.T) {
	t.Parallel()

	for _, dist := range []Logistic{
		{Mu: 0, S: 1},
		{Mu: 2, S: 3},
		{Mu: -5, S: 0.125},
	} {
		// At 1000 scale units from the mean, the log density is
		// -1000-log(S) to float64 precision: the remaining term,
		// -2*log1p(exp(-1000)), is far below the resolution at this
		// magnitude.
		for _, z := range []float64{-1000, 1000} {
			x := dist.Mu + dist.S*z
			want := -1000 - math.Log(dist.S)
			if got := dist.LogProb(x); !scalar.EqualWithinAbsOrRel(got, want, 1e-14, 1e-14) {
				t.Errorf("Wrong LogProb(%g) with Mu=%g, S=%g: got %g, want %g", x, dist.Mu, dist.S, got, want)
			}
		}
		for _, x := range []float64{math.Inf(-1), math.Inf(1)} {
			if got := dist.LogProb(x); !math.IsInf(got, -1) {
				t.Errorf("Wrong LogProb(%g) with Mu=%g, S=%g: got %g, want -Inf", x, dist.Mu, dist.S, got)
			}
		}
	}
}

func TestQuantile(t *testing.T) {
	t.Parallel()

	for _, v := range []struct {
		mu, s, input, want float64
	}{
		{0.0, 1.0, 0.5, 0.0},
		{0.0, 1.0, 0.0, math.Inf(-1)},
		{0.0, 1.0, 1.0, math.Inf(+1)},
	} {
		l := Logistic{Mu: v.mu, S: v.s}
		if result := l.Quantile(v.input); result != v.want {
			t.Errorf("Wrong Quantile(%f) with Mu=%f, S=%f: %f != %f", v.input, l.Mu, l.S, result, v.want)
		}
	}

	// Edge case with NaN.
	l := Logistic{Mu: 0, S: 0}

	input := 0.0
	if result := l.Quantile(input); !math.IsNaN(result) {
		t.Errorf("Wrong Quantile(%f) with Mu=%f, S=%f: %f is not NaN", input, l.Mu, l.S, result)
	}
}
