package diff

import (
	"math"
	"testing"
)

var xSquared = func(x float64) float64 { return x * x }

type testPoint struct {
	f   func(float64) float64
	loc float64
	ans float64
}

var testsFirst = []testPoint{
	{
		f:   xSquared,
		loc: 0,
		ans: 0,
	},
	{
		f:   xSquared,
		loc: 5,
		ans: 10,
	},
	{
		f:   xSquared,
		loc: 2,
		ans: 4,
	},
	{
		f:   xSquared,
		loc: -5,
		ans: -10,
	},
}

var testsSecond = []testPoint{
	{
		f:   xSquared,
		loc: 0,
		ans: 2,
	},
	{
		f:   xSquared,
		loc: 5,
		ans: 2,
	},
	{
		f:   xSquared,
		loc: 2,
		ans: 2,
	},
	{
		f:   xSquared,
		loc: -5,
		ans: 2,
	},
}

func testFirstOrder(t *testing.T, method Method, tol float64, tests []testPoint) {
	for _, test := range tests {
		settings := DefaultSettings()
		settings.Method = method
		ans := Derivative(test.f, test.loc, settings)
		if math.Abs(test.ans-ans) > tol {
			t.Errorf("ans mismatch: expected %v, found %v", test.ans, ans)
		}
		settings.Concurrent = true
		ans = Derivative(test.f, test.loc, settings)
		if math.Abs(test.ans-ans) > tol {
			t.Errorf("ans mismatch: expected %v, found %v", test.ans, ans)
		}
	}
}

func TestForward(t *testing.T) {
	testFirstOrder(t, Forward, 1e-4, testsFirst)
}

func TestBackward(t *testing.T) {
	testFirstOrder(t, Backward, 1e-4, testsFirst)
}

func TestCentral(t *testing.T) {
	testFirstOrder(t, Central, 1e-6, testsFirst)
}

func TestCentralSecond(t *testing.T) {
	testFirstOrder(t, Central2nd, 2e-3, testsSecond)
}
