// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"sort"
	"testing"

	"golang.org/x/exp/rand"
)

func TestLognormal(t *testing.T) {
	t.Parallel()
	src := rand.New(rand.NewSource(1))
	for i, dist := range []LogNormal{
		{Mu: 0.1, Sigma: 0.3, Src: src},
		{Mu: 0.01, Sigma: 0.01, Src: src},
		{Mu: 2, Sigma: 0.01, Src: src},
	} {
		const (
			tol = 1e-2
			n   = 1e5
		)
		x := make([]float64, n)
		generateSamples(x, dist)
		sort.Float64s(x)

		checkMean(t, i, x, dist, tol)
		checkVarAndStd(t, i, x, dist, tol)
		checkEntropy(t, i, x, dist, tol)
		checkExKurtosis(t, i, x, dist, 2e-1)
		checkSkewness(t, i, x, dist, 5e-2)
		checkMedian(t, i, x, dist, tol)
		checkQuantileCDFSurvival(t, i, x, dist, tol)
		checkProbContinuous(t, i, x, 0, math.Inf(1), dist, 1e-10)
		checkProbQuantContinuous(t, i, x, dist, tol)
		checkMode(t, i, x, dist, 1e-2, 1e-2)

		logProb := dist.LogProb(-0.0001)
		if !math.IsInf(logProb, -1) {
			t.Errorf("Expected LogProb == -Inf for x < 0, got %v", logProb)
		}
		if dist.NumParameters() != 2 {
			t.Errorf("Mismatch in NumParameters: got %v, want 2", dist.NumParameters())
		}
	}
}

// See https://github.com/gonum/gonum/issues/577 for details.
func TestLognormalIssue577(t *testing.T) {
	t.Parallel()
	x := 1.0e-16
	max := 1.0e-295
	cdf := LogNormal{Mu: 0, Sigma: 1}.CDF(x)
	if cdf <= 0 {
		t.Errorf("LogNormal{0,1}.CDF(%e) should be positive. got: %e", x, cdf)
	}
	if cdf > max {
		t.Errorf("LogNormal{0,1}.CDF(%e) is greater than %e. got: %e", x, max, cdf)
	}
}
