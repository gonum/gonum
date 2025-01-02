// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"math/rand/v2"
	"sort"
	"testing"

	"gonum.org/v1/gonum/floats/scalar"
)

func TestChiProb(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		x, k, want float64
	}{
		{10, 3, 1.538919725341288e-20},
		{2.3, 3, 0.2997000593061405},
		{0.8, 0.2, 0.1702707693447167},
	} {
		pdf := Chi{test.k, nil}.Prob(test.x)
		if !scalar.EqualWithinAbsOrRel(pdf, test.want, 1e-10, 1e-10) {
			t.Errorf("Pdf mismatch, x = %v, K = %v. Got %v, want %v", test.x, test.k, pdf, test.want)
		}
	}
}

func TestChiCDF(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		x, k, want float64
	}{
		// Values calculated with scipy.stats.chi.cdf
		{0, 1, 0},
		{0.01, 5, 5.319040436531812e-12},
		{0.05, 3, 3.3220267268523235e-05},
		{0.5, 2, 0.1175030974154046},
		{0.95, 3, 0.17517554009157732},
		{0.99, 5, 0.035845177452671864},
		{1, 1, 0.6826894921370859},
		{1.5, 4, 0.3101135068635068},
		{10, 10, 1},
		{25, 15, 1},
	} {
		cdf := Chi{test.k, nil}.CDF(test.x)
		if !scalar.EqualWithinAbsOrRel(cdf, test.want, 1e-10, 1e-10) {
			t.Errorf("CDF mismatch, x = %v, K = %v. Got %v, want %v", test.x, test.k, cdf, test.want)
		}
	}
}

func TestChi(t *testing.T) {
	t.Parallel()
	src := rand.New(rand.NewPCG(1, 1))
	for i, b := range []Chi{
		{3, src},
		{1.5, src},
		{0.9, src},
	} {
		testChi(t, b, i)
	}
}

func testChi(t *testing.T, c Chi, i int) {
	const (
		tol  = 1e-2
		n    = 1e6
		bins = 50
	)
	x := make([]float64, n)
	generateSamples(x, c)
	sort.Float64s(x)

	testRandLogProbContinuous(t, i, 0, x, c, tol, bins)
	checkMean(t, i, x, c, tol)
	checkMedian(t, i, x, c, tol)
	checkVarAndStd(t, i, x, c, tol)
	checkEntropy(t, i, x, c, tol)
	checkExKurtosis(t, i, x, c, 7e-2)
	checkProbContinuous(t, i, x, 0, math.Inf(1), c, 1e-5)
	checkQuantileCDFSurvival(t, i, x, c, 1e-2)

	expectedMode := math.Sqrt(c.K - 1)
	if !scalar.Same(c.Mode(), expectedMode) {
		t.Errorf("Mode is not equal to sqrt(k - 1). Got %v, want %v", c.Mode(), expectedMode)
	}
	if c.NumParameters() != 1 {
		t.Errorf("NumParameters is not 1. Got %v", c.NumParameters())
	}
	survival := c.Survival(-0.00001)
	if survival != 1 {
		t.Errorf("Survival is not 1 for negative argument. Got %v", survival)
	}
}
