// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"testing"

	"github.com/gonum/floats"
	"github.com/gonum/optimize/functions"
)

type gradFreeTest struct {
	// f is the function that is being minimized.
	f Function
	// x is the initial guess.
	x []float64
	// absTol is the absolute function convergence for the test. If absTol == 0,
	// the default value of 1e-6 will be used
	absTol float64
	// absIter is the number of iterations for function convergence. If iter == 0,
	// the default value of 5 will be used
	absIter int
	// long indicates that the test takes long time to finish and will be
	// excluded if testing.Short() is true.
	long bool
}

// gradFree ensures that the function is gradient free
type gradFree struct {
	f Function
}

func (g gradFree) Func(x []float64) float64 {
	return g.f.Func(x)
}

// makeGradFree ensures that a function contains no gradient method
func makeGradFree(f Function) gradFree {
	return gradFree{f}
}

// TODO(btracey): The gradient is still evaluated and tested if available
// even if a gradient-free method is being used. This should be fixed. When that
// is fixed, this should include Functions that also have gradients.
var gradFreeTests = []gradFreeTest{
	{
		f: makeGradFree(functions.ExtendedRosenbrock{}),
		x: []float64{-10, 10},
	},
	{
		f: makeGradFree(functions.ExtendedRosenbrock{}),
		x: []float64{-5, 4, 16, 3},
	},
}

func TestLocalGradFree(t *testing.T) {
	testLocalGradFree(t, gradFreeTests, nil)
}

func TestNelderMead(t *testing.T) {
	testLocalGradFree(t, gradFreeTests, &NelderMead{})
}

func testLocalGradFree(t *testing.T, tests []gradFreeTest, method Method) {
	for _, test := range tests {
		if test.long && testing.Short() {
			continue
		}
		settings := DefaultSettings()
		settings.Recorder = nil
		if test.absIter == 0 {
			test.absIter = 5
		}
		if test.absTol == 0 {
			test.absTol = 1e-6
		}
		result, err := Local(test.f, test.x, settings, method)
		if err != nil {
			t.Errorf("error finding minimum (%v) for \n%v", err, test)
		}
		if result == nil {
			t.Errorf("nil result without error for:\n%v", test)
			continue
		}
		if result.Status != FunctionConvergence {
			t.Errorf("Status not %v, %v instead", FunctionConvergence, result.Status)
		}

		result2, err := Local(test.f, test.x, settings, method)
		if err != nil {
			t.Errorf("error finding minimum (%v) when reusing Method for \n%v", err, test)
		}
		if result.FuncEvaluations != result2.FuncEvaluations ||
			result.F != result2.F || !floats.Equal(result.X, result2.X) {
			t.Errorf("Different result when reuse method")
		}
	}
}
