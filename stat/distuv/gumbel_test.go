// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"sort"
	"testing"

	"gonum.org/v1/gonum/floats"

	"golang.org/x/exp/rand"
)

func TestGumbelProbCDF(t *testing.T) {
	for _, test := range []struct {
		x, mu, beta, wantProb, wantCDF float64
	}{
		// Values calculated with scipy.stats.gumbel_r .
		{-2, 0, 1, 0.0045662814201279153, 0.00061797898933109343},
		{0.01, 0, 1, 0.36786110881643569, 0.37155817442380817},
		{6, 0, 1, 0.0024726155730149077, 0.99752431739275249},
	} {
		g := Gumbel{Mu: test.mu, Beta: test.beta}
		pdf := g.Prob(test.x)
		if !floats.EqualWithinAbsOrRel(pdf, test.wantProb, 1e-12, 1e-12) {
			t.Errorf("Prob mismatch, x = %v, mu = %v, beta = %v. Got %v, want %v", test.x, test.mu, test.beta, pdf, test.wantProb)
		}
		cdf := g.CDF(test.x)
		if !floats.EqualWithinAbsOrRel(cdf, test.wantCDF, 1e-12, 1e-12) {
			t.Errorf("CDF mismatch, x = %v, mu = %v, beta = %v. Got %v, want %v", test.x, test.mu, test.beta, cdf, test.wantCDF)
		}
	}
}

func TestGumbel(t *testing.T) {
	src := rand.New(rand.NewSource(1))
	for i, b := range []Gumbel{
		{0, 1, src},
		{-5, 6, src},
		{3, 0.1, src},
	} {
		testGumbel(t, b, i)
	}
}

func testGumbel(t *testing.T, g Gumbel, i int) {
	const (
		tol  = 1e-2
		n    = 5e5
		bins = 50
	)
	x := make([]float64, n)
	generateSamples(x, g)
	sort.Float64s(x)

	min := math.Inf(-1)
	testRandLogProbContinuous(t, i, min, x, g, tol, bins)
	checkProbContinuous(t, i, x, g, 1e-3)
	checkMean(t, i, x, g, tol)
	checkVarAndStd(t, i, x, g, tol)
	checkExKurtosis(t, i, x, g, 1e-1)
	checkSkewness(t, i, x, g, 5e-2)
	checkQuantileCDFSurvival(t, i, x, g, 5e-3)
}
