// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"sort"
	"testing"

	"github.com/gonum/floats"
	"github.com/gonum/stat"
)

// fullDist is a distribution that implements all the basic functions.
type fullDist interface {
	CDF(x float64) float64
	Entropy() float64
	ExKurtosis() float64
	LogProb(x float64) float64
	Mean() float64
	Median() float64
	Prob(x float64) float64
	Quantile(p float64) float64
	Rand() float64
	Skewness() float64
	StdDev() float64
	Survival(x float64) float64
	Variance() float64
}

// testFullDist tests all of the functions of a fullDist.
func testFullDist(t *testing.T, f fullDist, i int) {
	tol := 1e-1
	const n = 1e6
	xs := make([]float64, n)
	for i := range xs {
		xs[i] = f.Rand()
	}
	sortedXs := make([]float64, n)
	copy(sortedXs, xs)
	sort.Float64s(sortedXs)
	tmp := make([]float64, n)

	// Mean check.
	mean := stat.Mean(xs, nil)
	if !floats.EqualWithinAbsOrRel(mean, f.Mean(), tol, tol) {
		t.Errorf("Mean mismatch case %v: want: %v, got: %v", i, mean, f.Mean())
	} else {
		mean = f.Mean()
	}

	// Median check.
	median := stat.Quantile(0.5, stat.Empirical, sortedXs, nil)
	if !floats.EqualWithinAbsOrRel(median, f.Median(), tol, tol) {
		t.Errorf("Median mismatch case %v: want: %v, got: %v", i, median, f.Median())
	}

	// Variance check.
	variance := stat.Variance(xs, nil)
	if !floats.EqualWithinAbsOrRel(variance, f.Variance(), tol, tol) {
		t.Errorf("Variance mismatch case %v: want: %v, got: %v", i, mean, f.Variance())
	} else {
		variance = f.Variance()
	}

	std := math.Sqrt(variance)
	if !floats.EqualWithinAbsOrRel(std, f.StdDev(), tol, tol) {
		t.Errorf("StdDev mismatch case %v: want: %v, got: %v", i, mean, f.StdDev())
	} else {
		std = f.StdDev()
	}

	// Entropy check.
	for i, x := range xs {
		tmp[i] = -f.LogProb(x)
	}
	entropy := stat.Mean(tmp, nil)
	if !floats.EqualWithinAbsOrRel(entropy, f.Entropy(), tol, tol) {
		t.Errorf("Entropy mismatch case %v: want: %v, got: %v", i, entropy, f.Entropy())
	}

	// Excess Kurtosis check.
	for i, x := range xs {
		tmp[i] = math.Pow(x-mean, 4)
	}
	mu4 := stat.Mean(tmp, nil)
	kurtosis := mu4/(variance*variance) - 3
	if !floats.EqualWithinAbsOrRel(kurtosis, f.ExKurtosis(), tol, tol) {
		t.Errorf("ExKurtosis mismatch case %v: want: %v, got: %v", i, kurtosis, f.ExKurtosis())
	}

	// Skewness check.
	for i, x := range xs {
		tmp[i] = math.Pow(x-mean, 3)
	}
	mu3 := stat.Mean(tmp, nil)
	skewness := mu3 / math.Pow(std, 3)
	if !floats.EqualWithinAbsOrRel(skewness, f.Skewness(), tol, tol) {
		t.Errorf("ExKurtosis mismatch case %v: want: %v, got: %v", i, skewness, f.Skewness())
	}

	// Quantile, CDF, and survival check.
	for i, p := range []float64{0.1, 0.25, 0.5, 0.75, 0.9} {
		x := f.Quantile(p)
		cdf := f.CDF(x)
		estCDF := stat.CDF(x, stat.Empirical, sortedXs, nil)
		if !floats.EqualWithinAbsOrRel(cdf, estCDF, tol, tol) {
			t.Errorf("CDF mismatch case %v: want: %v, got: %v", i, estCDF, cdf)
		}
		if !floats.EqualWithinAbsOrRel(cdf, p, tol, tol) {
			t.Errorf("Quantile/CDF mismatch case %v: want: %v, got: %v", i, p, cdf)
		}
		if math.Abs(1-cdf-f.Survival(x)) > 1e-14 {
			t.Errorf("Survival/CDF mismatch case %v: want: %v, got: %v", i, 1-cdf, f.Survival(x))
		}
	}

	// Prob and LogProb check.
	m := 1001
	bins := make([]float64, m)
	dividers := make([]float64, m)
	floats.Span(bins, 0, 1)
	for i, v := range bins {
		dividers[i] = f.Quantile(v)
	}
	counts := stat.Histogram(nil, dividers, sortedXs, nil)
	// Test PDf against normalized count
	for i, v := range counts {
		v /= float64(n)
		at := f.Quantile((bins[i] + bins[i+1]) / 2)
		prob := f.Prob(at)
		if !floats.EqualWithinAbsOrRel(skewness, f.Skewness(), tol, tol) {
			t.Errorf("Prob mismatch case %v at %v: want: %v, got: %v", i, at, v, prob)
			break
		}
		if math.Abs(math.Log(prob)-f.LogProb(at)) > 1e-14 {
			t.Errorf("Prob and LogProb mismatch case %v at %v: want %v, got %v", i, at, math.Log(prob), f.LogProb(at))
			break
		}
	}
}
