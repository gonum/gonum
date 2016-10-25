// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/gonum/floats"
)

func TestStudentsTProb(t *testing.T) {
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
		if !floats.EqualWithinAbsOrRel(pdf, test.want, 1e-10, 1e-10) {
			t.Errorf("Pdf mismatch, x = %v, Nu = %v. Got %v, want %v", test.x, test.nu, pdf, test.want)
		}
	}
}

func TestStudentsT(t *testing.T) {
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
	tol := 1e-2
	const n = 1e6
	const bins = 50
	x := make([]float64, n)
	generateSamples(x, c)
	sort.Float64s(x)

	testRandLogProbContinuous(t, i, math.Inf(-1), x, c, tol, bins)
	checkMean(t, i, x, c, tol)
	if c.Nu > 2 {
		checkVarAndStd(t, i, x, c, tol)
	}
	checkProbContinuous(t, i, x, c, 1e-3)
}
