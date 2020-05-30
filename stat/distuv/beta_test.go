// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"sort"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
)

func TestBetaProb(t *testing.T) {
	t.Parallel()
	// Values a comparison with scipy
	for _, test := range []struct {
		x, alpha, beta, want float64
	}{
		{0.1, 2, 0.5, 0.079056941504209499},
		{0.5, 1, 5.1, 0.29740426605235754},
		{0.1, 0.5, 0.5, 1.0610329539459691},
		{1, 0.5, 0.5, math.Inf(1)},
		{-1, 0.5, 0.5, 0},
	} {
		pdf := Beta{Alpha: test.alpha, Beta: test.beta}.Prob(test.x)
		if !floats.EqualWithinAbsOrRel(pdf, test.want, 1e-10, 1e-10) {
			t.Errorf("Pdf mismatch. Got %v, want %v", pdf, test.want)
		}
	}
}

func TestBetaRand(t *testing.T) {
	t.Parallel()
	src := rand.New(rand.NewSource(1))
	for i, b := range []Beta{
		{Alpha: 0.5, Beta: 0.5, Src: src},
		{Alpha: 5, Beta: 1, Src: src},
		{Alpha: 2, Beta: 2, Src: src},
		{Alpha: 2, Beta: 5, Src: src},
	} {
		testBeta(t, b, i)
	}
}

func testBeta(t *testing.T, b Beta, i int) {
	const (
		tol  = 1e-2
		n    = 5e4
		bins = 10
	)
	x := make([]float64, n)
	generateSamples(x, b)
	sort.Float64s(x)

	testRandLogProbContinuous(t, i, 0, x, b, tol, bins)
	checkMean(t, i, x, b, tol)
	checkVarAndStd(t, i, x, b, tol)
	checkExKurtosis(t, i, x, b, 5e-2)
	checkEntropy(t, i, x, b, 5e-3)
	checkProbContinuous(t, i, x, b, 1e-3)
	checkQuantileCDFSurvival(t, i, x, b, tol)
	checkProbQuantContinuous(t, i, x, b, tol)

	if b.NumParameters() != 2 {
		t.Errorf("Wrong number of parameters")
	}

	if b.CDF(-0.01) != 0 {
		t.Errorf("CDF below 0 is not 0")
	}
	if b.CDF(0) != 0 {
		t.Errorf("CDF at 0 is not 0")
	}
	if b.CDF(1) != 1 {
		t.Errorf("CDF at 1 is not 1")
	}
	if b.CDF(1.01) != 1 {
		t.Errorf("CDF above 1 is not 1")
	}

	if b.Survival(-0.01) != 1 {
		t.Errorf("Survival below 0 is not 1")
	}
	if b.Survival(0) != 1 {
		t.Errorf("Survival at 0 is not 1")
	}
	if b.Survival(1) != 0 {
		t.Errorf("Survival at 1 is not 0")
	}
	if b.Survival(1.01) != 0 {
		t.Errorf("Survival above 1 is not 0")
	}

	if !panics(func() { b.Quantile(-0.01) }) {
		t.Errorf("Quantile did not panic for negative argument")
	}
	if !panics(func() { b.Quantile(1.01) }) {
		t.Errorf("Quantile did not panic for argument above 1")
	}
}

func TestBetaBadParams(t *testing.T) {
	t.Parallel()
	src := rand.New(rand.NewSource(1))
	for _, alpha := range []float64{0, -0.1} {
		testBetaBadParams(t, alpha, 1, src)
		testBetaBadParams(t, 1, alpha, src)
		for _, beta := range []float64{0, -0.1} {
			testBetaBadParams(t, alpha, beta, src)
		}
	}
}

func testBetaBadParams(t *testing.T, alpha float64, beta float64, src rand.Source) {
	b := Beta{alpha, beta, src}
	if !panics(func() { b.Entropy() }) {
		t.Errorf("Entropy did not panic for Beta(%g, %g)", alpha, beta)
	}
	if !panics(func() { b.LogProb(0.5) }) {
		t.Errorf("LogProb did not panic for Beta(%g, %g)", alpha, beta)
	}
}

func TestBetaMode(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		alpha, beta, want float64
	}{
		{1, 2, 0},
		{0.5, 2, 0},
		{2, 1, 1},
		{2, 0.5, 1},
		{4, 5, 3. / 7.},
	} {
		mode := Beta{Alpha: test.alpha, Beta: test.beta}.Mode()
		if !floats.EqualWithinAbsOrRel(mode, test.want, 1e-10, 1e-10) {
			t.Errorf("Mode mismatch for Beta(%g, %g). Got %v, want %g", test.alpha, test.beta, mode, test.want)
		}
	}
	for _, test := range []struct {
		alpha, beta float64
	}{
		{1, 1},
		{0.5, 0.5},
		{1, 0.5},
		{0.5, 1},
	} {
		mode := Beta{Alpha: test.alpha, Beta: test.beta}.Mode()
		if !math.IsNaN(mode) {
			t.Errorf("Mode is not NaN for Beta(%g, %g). Got: %v", test.alpha, test.beta, mode)
		}
	}
}
