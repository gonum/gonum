// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/gonum/floats"
)

func TestChiSquaredProb(t *testing.T) {
	for _, test := range []struct {
		x, k, want float64
	}{
		{10, 3, 0.0085003666025203432},
		{2.3, 3, 0.19157345407042367},
		{0.8, 0.2, 0.080363259903912673},
	} {
		pdf := ChiSquared{test.k, nil}.Prob(test.x)
		if !floats.EqualWithinAbsOrRel(pdf, test.want, 1e-10, 1e-10) {
			t.Errorf("Pdf mismatch, x = %v, K = %v. Got %v, want %v", test.x, test.k, pdf, test.want)
		}
	}
}

func TestChiSquared(t *testing.T) {
	src := rand.New(rand.NewSource(1))
	for i, b := range []ChiSquared{
		{3, src},
		{1.5, src},
		{0.9, src},
	} {
		testChiSquared(t, b, i)
	}
}

func testChiSquared(t *testing.T, c ChiSquared, i int) {
	tol := 1e-2
	const n = 2e6
	const bins = 50
	x := make([]float64, n)
	generateSamples(x, c)
	sort.Float64s(x)

	testRandLogProbContinuous(t, i, 0, x, c, tol, bins)
	checkMean(t, i, x, c, tol)
	checkVarAndStd(t, i, x, c, tol)
	checkExKurtosis(t, i, x, c, 5e-2)
	checkProbContinuous(t, i, x, c, 1e-3)
}
