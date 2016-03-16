// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fd

import (
	"math"
	"testing"
)

var xSquared = func(x float64) float64 { return x * x }

type testPoint struct {
	f    func(float64) float64
	loc  float64
	fofx float64
	ans  float64
}

var testsFirst = []testPoint{
	{
		f:    xSquared,
		loc:  0,
		fofx: 0,
		ans:  0,
	},
	{
		f:    xSquared,
		loc:  5,
		fofx: 25,
		ans:  10,
	},
	{
		f:    xSquared,
		loc:  2,
		fofx: 4,
		ans:  4,
	},
	{
		f:    xSquared,
		loc:  -5,
		fofx: 25,
		ans:  -10,
	},
}

var testsSecond = []testPoint{
	{
		f:    xSquared,
		loc:  0,
		fofx: 0,
		ans:  2,
	},
	{
		f:    xSquared,
		loc:  5,
		fofx: 25,
		ans:  2,
	},
	{
		f:    xSquared,
		loc:  2,
		fofx: 4,
		ans:  2,
	},
	{
		f:    xSquared,
		loc:  -5,
		fofx: 25,
		ans:  2,
	},
}

func testDerivative(t *testing.T, formula Formula, tol float64, tests []testPoint) {
	for i, test := range tests {

		settings := DefaultSettings()
		settings.Formula = formula
		ans := Derivative(test.f, test.loc, settings)
		if math.Abs(test.ans-ans) > tol {
			t.Errorf("Case %v: ans mismatch serial: expected %v, found %v", i, test.ans, ans)
		}

		settings.OriginKnown = true
		settings.OriginValue = test.fofx
		ans = Derivative(test.f, test.loc, settings)
		if math.Abs(test.ans-ans) > tol {
			t.Errorf("Case %v: ans mismatch serial origin known: expected %v, found %v", i, test.ans, ans)
		}

		settings.OriginKnown = false
		settings.Concurrent = true
		ans = Derivative(test.f, test.loc, settings)
		if math.Abs(test.ans-ans) > tol {
			t.Errorf("Case %v: ans mismatch concurrent: expected %v, found %v", i, test.ans, ans)
		}

		settings.OriginKnown = true
		ans = Derivative(test.f, test.loc, settings)
		if math.Abs(test.ans-ans) > tol {
			t.Errorf("Case %v: ans mismatch concurrent: expected %v, found %v", i, test.ans, ans)
		}
	}
}

func TestForward(t *testing.T) {
	testDerivative(t, Forward, 2e-4, testsFirst)
}

func TestBackward(t *testing.T) {
	testDerivative(t, Backward, 2e-4, testsFirst)
}

func TestCentral(t *testing.T) {
	testDerivative(t, Central, 1e-6, testsFirst)
}

func TestCentralSecond(t *testing.T) {
	testDerivative(t, Central2nd, 1e-3, testsSecond)
}

// TestDerivativeDefault checks that the derivative works when settings is set to nil.
func TestDerivativeDefault(t *testing.T) {
	tol := 1e-6
	for i, test := range testsFirst {
		ans := Derivative(test.f, test.loc, nil)
		if math.Abs(test.ans-ans) > tol {
			t.Errorf("Case %v: ans mismatch default: expected %v, found %v", i, test.ans, ans)
		}
	}
}
