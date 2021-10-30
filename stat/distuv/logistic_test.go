package distuv

import (
	"math"
	"testing"
)

func TestLogisticParameters(t *testing.T) {
	t.Parallel()

	l := Logistic{Mu: 0, S: 0}

	if l.NumParameters() != 2 {
		t.Fail()
	}

	if l.ExKurtosis() != 6/5 {
		t.Fail()
	}

	if l.Skewness() != 0 {
		t.Fail()
	}

	if l.Mean() != l.Mu {
		t.Fail()
	}

	if l.Mean() != l.Median() {
		t.Fail()
	}

	if l.Mean() != l.Mode() {
		t.Fail()
	}
}

func TestLogisticStdDev(t *testing.T) {
	t.Parallel()

	const sq3 = 1.732050807568877293527446341505872366942805253810380 // sqrt(3)

	l := Logistic{Mu: 0, S: sq3 / math.Pi}
	if l.StdDev() != 1 {
		t.Fail()
	}

	if l.Variance() != 1 {
		t.Fail()
	}
}

func TestLogisticCDF(t *testing.T) {
	t.Parallel()

	// edge case of zero in denominator
	l := Logistic{Mu: 0, S: 0}

	if !math.IsNaN(l.CDF(0)) {
		t.Fail()
	}

	if l.CDF(1) != 1 {
		t.Fail()
	}

	l = Logistic{Mu: 0, S: 1}

	if l.CDF(0) != 0.5 {
		t.Fail()
	}
}

func TestLogisticSurvival(t *testing.T) {
	t.Parallel()

	l := Logistic{Mu: 0, S: 1}

	if l.Survival(0) != 0.5 {
		t.Fail()
	}
}

func TestLogisticProb(t *testing.T) {
	t.Parallel()

	// edge case of zero in denominator
	l := Logistic{Mu: 0, S: 0}

	if !math.IsNaN(l.Prob(0)) {
		t.Fail()
	}

	if !math.IsNaN(l.Prob(1)) {
		t.Fail()
	}

	l = Logistic{Mu: 0, S: 1}

	if l.Prob(0) != 0.25 {
		t.Fail()
	}

	if l.LogProb(0) != -math.Log(4) {
		t.Fail()
	}
}

func TestQuantile(t *testing.T) {
	t.Parallel()

	l := Logistic{Mu: 0, S: 0}

	if !math.IsNaN(l.Quantile(0)) {
		t.Fail()
	}

	l = Logistic{Mu: 0, S: 1}

	if !math.IsInf(l.Quantile(0), -1) {
		t.Fail()
	}

	if !math.IsInf(l.Quantile(1), 1) {
		t.Fail()
	}
}
