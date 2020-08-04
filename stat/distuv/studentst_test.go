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
	"gonum.org/v1/gonum/floats/scalar"
)

func TestStudentsTProb(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		x, mu, sigma, nu, want float64
	}{
		// Values comparison with scipy.
		{0.01, 0, 1, 2.74, 0.364778548181318},
		{-0.01, 0, 1, 2.74, 0.364778548181318},
		{0.4, 0, 1, 1.6, 0.30376391362582678},
		{-0.4, 0, 1, 1.6, 0.30376391362582678},
		{0.2, 15, 5, 10, 0.0024440848858034393},
	} {
		pdf := StudentsT{test.mu, test.sigma, test.nu, nil}.Prob(test.x)
		if !scalar.EqualWithinAbsOrRel(pdf, test.want, 1e-10, 1e-10) {
			t.Errorf("Pdf mismatch, x = %v, Nu = %v. Got %v, want %v", test.x, test.nu, pdf, test.want)
		}
	}
}

func TestStudentsT(t *testing.T) {
	t.Parallel()
	src := rand.New(rand.NewSource(1))
	for i, b := range []StudentsT{
		{0, 1, 3.3, src},
		{0, 1, 7.2, src},
		{0, 1, 12, src},
		{0.9, 0.8, 6, src},
	} {
		testStudentsT(t, b, i)
	}
}

func testStudentsT(t *testing.T, c StudentsT, i int) {
	const (
		tol  = 1e-2
		n    = 1e5
		bins = 50
	)
	x := make([]float64, n)
	generateSamples(x, c)
	sort.Float64s(x)

	testRandLogProbContinuous(t, i, math.Inf(-1), x, c, tol, bins)
	checkMean(t, i, x, c, tol)
	if c.Nu > 2 {
		checkVarAndStd(t, i, x, c, 5e-2)
	}
	checkProbContinuous(t, i, x, math.Inf(-1), math.Inf(1), c, 1e-10)
	checkQuantileCDFSurvival(t, i, x, c, tol)
	checkProbQuantContinuous(t, i, x, c, tol)
	if c.Mu != c.Mode() {
		t.Errorf("Mismatch in mode value: got %v, want %g", c.Mode(), c.Mu)
	}
	if c.NumParameters() != 3 {
		t.Errorf("Mismatch in NumParameters: got %v, want 3", c.NumParameters())
	}
}

func TestStudentsTQuantile(t *testing.T) {
	t.Parallel()
	nSteps := 101
	probs := make([]float64, nSteps)
	floats.Span(probs, 0, 1)
	for i, b := range []StudentsT{
		{0, 1, 3.3, nil},
		{0, 1, 7.2, nil},
		{0, 1, 12, nil},
		{0.9, 0.8, 6, nil},
	} {
		for _, p := range probs {
			x := b.Quantile(p)
			p2 := b.CDF(x)
			if !scalar.EqualWithinAbsOrRel(p, p2, 1e-10, 1e-10) {
				t.Errorf("mismatch between CDF and Quantile. Case %v. Want %v, got %v", i, p, p2)
				break
			}
		}
	}
}

func TestStudentsVarianceSpecial(t *testing.T) {
	t.Parallel()
	dist := StudentsT{0, 1, 1, nil}
	variance := dist.Variance()
	if !math.IsNaN(variance) {
		t.Errorf("Expected NaN variance for Nu <= 1, got %v", variance)
	}
	dist = StudentsT{0, 1, 2, nil}
	variance = dist.Variance()
	if !math.IsInf(variance, 1) {
		t.Errorf("Expected +Inf variance for 1 < Nu <= 2, got %v", variance)
	}
}
