// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"sort"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats/scalar"
)

func TestFProb(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		x, d1, d2, want float64
	}{
		// Values calculated with scipy.stats.f
		{0.0001, 4, 6, 0.00053315559110558126},
		{0.1, 1, 1, 0.91507658371794609},
		{0.5, 11, 7, 0.66644660411410883},
		{0.9, 20, 15, 0.88293424959522437},
		{1, 1, 1, 0.15915494309189535},
		{2, 15, 12, 0.16611971273429088},
		{5, 4, 8, 0.013599775603702537},
		{10, 12, 9, 0.00032922887567957289},
		{100, 7, 7, 6.08037637806889e-08},
		{1000, 2, 1, 1.1171959870312232e-05},
	} {
		pdf := F{test.d1, test.d2, nil}.Prob(test.x)
		if !scalar.EqualWithinAbsOrRel(pdf, test.want, 1e-10, 1e-10) {
			t.Errorf("Prob mismatch, x = %v, d1 = %v, d2 = %v. Got %v, want %v", test.x, test.d1, test.d2, pdf, test.want)
		}
	}
}

func TestFCDF(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		x, d1, d2, want float64
	}{
		// Values calculated with scipy.stats.f
		{0.0001, 4, 6, 2.6660741629519019e-08},
		{0.1, 1, 1, 0.19498222904213672},
		{0.5, 11, 7, 0.14625028471336987},
		{0.9, 20, 15, 0.40567939897287852},
		{1, 1, 1, 0.50000000000000011},
		{2, 15, 12, 0.8839384428956264},
		{5, 4, 8, 0.97429642410900219},
		{10, 12, 9, 0.99915733385467187},
		{100, 7, 7, 0.99999823560259171},
		{1000, 2, 1, 0.97764490829950534},
	} {
		cdf := F{test.d1, test.d2, nil}.CDF(test.x)
		if !scalar.EqualWithinAbsOrRel(cdf, test.want, 1e-10, 1e-10) {
			t.Errorf("CDF mismatch, x = %v, d1 = %v, d2 = %v. Got %v, want %v", test.x, test.d1, test.d2, cdf, test.want)
		}
	}
}

func TestF(t *testing.T) {
	t.Parallel()
	src := rand.New(rand.NewSource(1))
	for i, f := range []F{
		{13, 16, src},
		{42, 31, src},
		{77, 92, src},
	} {
		testF(t, f, i)
	}
}

func testF(t *testing.T, f F, i int) {
	const (
		tol  = 1e-2
		n    = 1e6
		bins = 50
	)
	x := make([]float64, n)
	generateSamples(x, f)
	sort.Float64s(x)

	testRandLogProbContinuous(t, i, 0, x, f, tol, bins)
	checkProbContinuous(t, i, x, 0, math.Inf(1), f, 1e-4)
	checkMean(t, i, x, f, tol)
	checkVarAndStd(t, i, x, f, tol)
	checkExKurtosis(t, i, x, f, 1e-1)
	checkSkewness(t, i, x, f, 5e-2)
	checkQuantileCDFSurvival(t, i, x, f, 5e-3)
	checkMode(t, i, x, f, 2e-2, 3e-2)

	if f.NumParameters() != 2 {
		t.Errorf("Wrong number of parameters. Got %v, want 2", f.NumParameters())
	}
}

func TestFUndefined(t *testing.T) {
	t.Parallel()
	for _, d1 := range []float64{1, 100} {
		for _, d2 := range []float64{4, 8} {
			f := F{d1, d2, nil}
			exKurt := f.ExKurtosis()
			if !math.IsNaN(exKurt) {
				t.Errorf("Expected NaN excess kurtosis for D1 = %g and D2 = %g, got %v", d1, d2, exKurt)
			}
		}
	}
	for _, d1 := range []float64{1, 100} {
		for _, d2 := range []float64{1, 2} {
			f := F{d1, d2, nil}
			mean := f.Mean()
			if !math.IsNaN(mean) {
				t.Errorf("Expected NaN mean for D1 = %g and D2 = %g, got %v", d1, d2, mean)
			}
		}
	}
	for _, d1 := range []float64{1, 2} {
		for _, d2 := range []float64{1, 100} {
			f := F{d1, d2, nil}
			mode := f.Mode()
			if !math.IsNaN(mode) {
				t.Errorf("Expected NaN mode for D1 = %g and D2 = %g, got %v", d1, d2, mode)
			}
		}
	}
	for _, d1 := range []float64{1, 100} {
		for _, d2 := range []float64{3, 6} {
			f := F{d1, d2, nil}
			skewness := f.Skewness()
			if !math.IsNaN(skewness) {
				t.Errorf("Expected NaN skewness for D1 = %g and D2 = %g, got %v", d1, d2, skewness)
			}
		}
	}
	for _, d1 := range []float64{1, 100} {
		for _, d2 := range []float64{2, 4} {
			f := F{d1, d2, nil}
			variance := f.Variance()
			if !math.IsNaN(variance) {
				t.Errorf("Expected NaN variance for D1 = %g and D2 = %g, got %v", d1, d2, variance)
			}
			stdDev := f.StdDev()
			if !math.IsNaN(stdDev) {
				t.Errorf("Expected NaN standard deviation for D1 = %g and D2 = %g, got %v", d1, d2, variance)
			}
		}
	}
}
