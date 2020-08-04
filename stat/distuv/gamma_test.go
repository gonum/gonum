// Copyright ©2016 The Gonum Authors. All rights reserved.
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

func TestGamma(t *testing.T) {
	t.Parallel()
	// Values a comparison with scipy
	for _, test := range []struct {
		x, alpha, want float64
	}{
		{0.9, 0.1, 0.046986817861555757},
		{0.9, 0.01, 0.0045384353289090401},
		{0.45, 0.01, 0.014137035997241795},
	} {
		pdf := Gamma{Alpha: test.alpha, Beta: 1}.Prob(test.x)
		if !scalar.EqualWithinAbsOrRel(pdf, test.want, 1e-10, 1e-10) {
			t.Errorf("Pdf mismatch. Got %v, want %v", pdf, test.want)
		}
	}
	src := rand.New(rand.NewSource(1))
	for i, g := range []Gamma{
		{Alpha: 0.1, Beta: 0.8, Src: src},
		{Alpha: 0.3, Beta: 0.8, Src: src},
		{Alpha: 0.5, Beta: 0.8, Src: src},
		{Alpha: 0.9, Beta: 6, Src: src},
		{Alpha: 0.9, Beta: 500, Src: src},
		{Alpha: 1, Beta: 1, Src: src},
		{Alpha: 1.6, Beta: 0.4, Src: src},
		{Alpha: 2.6, Beta: 1.5, Src: src},
		{Alpha: 5.6, Beta: 0.5, Src: src},
		{Alpha: 30, Beta: 1.7, Src: src},
		{Alpha: 30.2, Beta: 1.7, Src: src},
	} {
		testGamma(t, g, i)
	}
}

func testGamma(t *testing.T, f Gamma, i int) {
	// TODO(btracey): Replace this when Gamma implements FullDist.
	const (
		tol  = 1e-2
		n    = 1e5
		bins = 50
	)
	x := make([]float64, n)
	generateSamples(x, f)
	sort.Float64s(x)

	var quadTol float64

	if f.Alpha < 0.3 {
		// Gamma PDF has a singularity at 0 for alpha < 1,
		// which gets sharper as alpha -> 0.
		quadTol = 0.2
	} else {
		quadTol = tol
	}
	testRandLogProbContinuous(t, i, 0, x, f, quadTol, bins)
	checkMean(t, i, x, f, tol)
	checkVarAndStd(t, i, x, f, 2e-2)
	checkExKurtosis(t, i, x, f, 2e-1)
	switch {
	case f.Alpha < 0.3:
		quadTol = 0.1
	case f.Alpha < 1:
		quadTol = 1e-3
	default:
		quadTol = 1e-10
	}
	checkProbContinuous(t, i, x, 0, math.Inf(1), f, quadTol)
	checkQuantileCDFSurvival(t, i, x, f, 5e-2)
	if f.Alpha < 1 {
		if !math.IsNaN(f.Mode()) {
			t.Errorf("Expected NaN mode for alpha < 1, got %v", f.Mode())
		}
	} else {
		checkMode(t, i, x, f, 2e-1, 1)
	}
	cdfNegX := f.CDF(-0.0001)
	if cdfNegX != 0 {
		t.Errorf("Expected zero CDF for negative argument, got %v", cdfNegX)
	}
	survNegX := f.Survival(-0.0001)
	if survNegX != 1 {
		t.Errorf("Expected survival function of 1 for negative argument, got %v", survNegX)
	}
	if f.NumParameters() != 2 {
		t.Errorf("Mismatch in NumParameters: got %v, want 2", f.NumParameters())
	}
	lPdf := f.LogProb(-0.0001)
	if !math.IsInf(lPdf, -1) {
		t.Errorf("Expected log(CDF) to be -Infinity for negative argument, got %v", lPdf)
	}
}

func TestGammaPanics(t *testing.T) {
	t.Parallel()
	g := Gamma{1, 0, nil}
	if !panics(func() { g.Rand() }) {
		t.Errorf("Expected Rand panic for Beta <= 0")
	}
	g = Gamma{0, 1, nil}
	if !panics(func() { g.Rand() }) {
		t.Errorf("Expected Rand panic for Alpha <= 0")
	}
}
