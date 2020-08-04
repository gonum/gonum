// Copyright ©2018 The Gonum Authors. All rights reserved.
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

func TestInverseGamma(t *testing.T) {
	t.Parallel()
	// Values extracted from a comparison with scipy
	for _, test := range []struct {
		x, alpha, want float64
	}{
		{0.9, 4.5, 0.050521067785046482},
		{0.04, 45, 0.10550644842525572},
		{20, 0.4, 0.0064691988681571536},
	} {
		pdf := InverseGamma{Alpha: test.alpha, Beta: 1}.Prob(test.x)
		if !scalar.EqualWithinAbsOrRel(pdf, test.want, 1e-10, 1e-10) {
			t.Errorf("Pdf mismatch. Got %v, want %v", pdf, test.want)
		}
	}
	src := rand.NewSource(1)
	for i, g := range []InverseGamma{
		{Alpha: 5.6, Beta: 0.5, Src: src},
		{Alpha: 30, Beta: 1.7, Src: src},
		{Alpha: 30.2, Beta: 1.7, Src: src},
	} {
		testInverseGamma(t, g, i)
	}
}

func testInverseGamma(t *testing.T, f InverseGamma, i int) {
	const (
		tol  = 1e-2
		n    = 1e6
		bins = 50
	)
	x := make([]float64, n)
	generateSamples(x, f)
	sort.Float64s(x)

	testRandLogProbContinuous(t, i, 0, x, f, tol, bins)
	checkMean(t, i, x, f, tol)
	checkVarAndStd(t, i, x, f, 2e-2)
	checkExKurtosis(t, i, x, f, 2e-1)
	checkProbContinuous(t, i, x, 0, math.Inf(1), f, 1e-10)
	checkQuantileCDFSurvival(t, i, x, f, 5e-2)
	checkMode(t, i, x, f, 1e-2, 1e-2)

	cdf0 := f.CDF(0)
	if cdf0 != 0 {
		t.Errorf("Expected zero CDF at 0, but got: %v", cdf0)
	}
	cdfNeg := f.CDF(-0.0001)
	if cdfNeg != 0 {
		t.Errorf("Expected zero CDF for a negative argument, but got: %v", cdfNeg)
	}
	surv0 := f.Survival(0)
	if surv0 != 1 {
		t.Errorf("Mismatch in Survival at 0. Got %v, want 1", surv0)
	}
	survNeg := f.Survival(-0.0001)
	if survNeg != 1 {
		t.Errorf("Mismatch in Survival for a negative argument. Got %v, want 1", survNeg)
	}
	if f.NumParameters() != 2 {
		t.Errorf("Mismatch in NumParameters. Got %v, want 2", f.NumParameters())
	}
	pdf0 := f.Prob(0)
	if pdf0 != 0 {
		t.Errorf("Expected zero PDF at 0, but got: %v", pdf0)
	}
	pdfNeg := f.Prob(-0.0001)
	if pdfNeg != 0 {
		t.Errorf("Expected zero PDF for a negative argument, but got: %v", pdfNeg)
	}
}

func TestInverseGammaLowAlpha(t *testing.T) {
	t.Parallel()
	f := InverseGamma{Alpha: 1, Beta: 1}
	mean := f.Mean()
	if !math.IsInf(mean, 1) {
		t.Errorf("Expected +Inf mean for alpha <= 1, got %v", mean)
	}
	f = InverseGamma{Alpha: 2, Beta: 1}
	stdDev := f.StdDev()
	if !math.IsInf(stdDev, 1) {
		t.Errorf("Expected +Inf standard deviation for alpha <= 2, got %v", stdDev)
	}
	variance := f.Variance()
	if !math.IsInf(variance, 1) {
		t.Errorf("Expected +Inf variance for alpha <= 2, got %v", variance)
	}
	f = InverseGamma{Alpha: 4, Beta: 1}
	exKurt := f.ExKurtosis()
	if !math.IsInf(exKurt, 1) {
		t.Errorf("Expected +Inf excess kurtosis for alpha <= 4, got %v", exKurt)
	}
}
