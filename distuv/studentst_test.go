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
		x, nu, want float64
	}{
		// Values comparison with scipy.
		{0.01, 2.74, 0.364778548181318},
		{-0.01, 2.74, 0.364778548181318},
		{0.4, 1.6, 0.30376391362582678},
		{-0.4, 1.6, 0.30376391362582678},
	} {
		pdf := StudentsT{test.nu, nil}.Prob(test.x)
		if !floats.EqualWithinAbsOrRel(pdf, test.want, 1e-10, 1e-10) {
			t.Errorf("Pdf mismatch, x = %v, Nu = %v. Got %v, want %v", test.x, test.nu, pdf, test.want)
		}
	}
}

func TestStudentsT(t *testing.T) {
	src := rand.New(rand.NewSource(1))
	for i, b := range []StudentsT{
		{3.3, src},
		{7.2, src},
		{12, src},
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
	if c.Nu > 4 {
		checkExKurtosis(t, i, x, c, 5e-2)
	}
	checkProbContinuous(t, i, x, c, 1e-3)
}
