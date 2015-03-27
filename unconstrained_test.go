// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/gonum/floats"
	"github.com/gonum/optimize/functions"
)

type unconstrainedTest struct {
	// f is the function that is being minimized.
	f Function
	// x is the initial guess.
	x []float64
	// gradTol is the absolute gradient tolerance for the test. If gradTol == 0,
	// the default tolerance 1e-12 will be used.
	gradTol float64
	// long indicates that the test takes long time to finish and will be
	// excluded if testing.Short() is true.
	long bool
}

func (t unconstrainedTest) String() string {
	dim := len(t.x)
	if dim <= 10 {
		// Print the initial X only for small-dimensional problems.
		return fmt.Sprintf("F: %v\nDim: %v\nInitial X: %v\nGradientAbsTol: %v",
			reflect.TypeOf(t.f), dim, t.x, t.gradTol)
	}
	return fmt.Sprintf("F: %v\nDim: %v\nGradientAbsTol: %v",
		reflect.TypeOf(t.f), dim, t.gradTol)
}

var gradientDescentTests = []unconstrainedTest{
	{
		f: functions.Beale{},
		x: []float64{1, 1},
	},
	{
		f: functions.Beale{},
		x: []float64{3.00001, 0.50001},
	},
	{
		f: functions.BiggsEXP2{},
		x: []float64{1, 2},
	},
	{
		f: functions.BiggsEXP2{},
		x: []float64{1.00001, 10.00001},
	},
	{
		f: functions.BiggsEXP3{},
		x: []float64{1, 2, 1},
	},
	{
		f: functions.BiggsEXP3{},
		x: []float64{1.00001, 10.00001, 3.00001},
	},
	{
		f: functions.ExtendedRosenbrock{},
		x: []float64{-1.2, 1},
	},
	{
		f: functions.ExtendedRosenbrock{},
		x: []float64{1.00001, 1.00001},
	},
	{
		f: functions.ExtendedRosenbrock{},
		x: []float64{-1.2, 1, -1.2},
	},
	{
		f:    functions.ExtendedRosenbrock{},
		x:    []float64{-120, 100, 50},
		long: true,
	},
	{
		f: functions.ExtendedRosenbrock{},
		x: []float64{1, 1, 1},
	},
	{
		f:       functions.ExtendedRosenbrock{},
		x:       []float64{1.00001, 1.00001, 1.00001},
		gradTol: 1e-8,
	},
	{
		f:       functions.Gaussian{},
		x:       []float64{0.4, 1, 0},
		gradTol: 1e-9,
	},
	{
		f:       functions.Gaussian{},
		x:       []float64{0.3989561, 1.0000191, 0},
		gradTol: 1e-9,
	},
	{
		f: functions.HelicalValley{},
		x: []float64{-1, 0, 0},
	},
	{
		f: functions.HelicalValley{},
		x: []float64{1.00001, 0.00001, 0.00001},
	},
	{
		f:       functions.Trigonometric{},
		x:       []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1},
		gradTol: 1e-8,
	},
	{
		f: functions.Trigonometric{},
		x: []float64{0.042964, 0.043976, 0.045093, 0.046338, 0.047744,
			0.049354, 0.051237, 0.195209, 0.164977, 0.060148},
		gradTol: 1e-8,
	},
	newVariablyDimensioned(2, 0),
	{
		f: functions.VariablyDimensioned{},
		x: []float64{1.00001, 1.00001},
	},
	newVariablyDimensioned(10, 0),
	{
		f: functions.VariablyDimensioned{},
		x: []float64{1.00001, 1.00001, 1.00001, 1.00001, 1.00001, 1.00001, 1.00001, 1.00001, 1.00001, 1.00001},
	},
}

var cgTests = []unconstrainedTest{
	{
		f: functions.BiggsEXP4{},
		x: []float64{1, 2, 1, 1},
	},
	{
		f: functions.BiggsEXP4{},
		x: []float64{1.00001, 10.00001, 1.00001, 5.00001},
	},
	{
		f:       functions.BiggsEXP5{},
		x:       []float64{1, 2, 1, 1, 1},
		gradTol: 1e-7,
	},
	{
		f: functions.BiggsEXP5{},
		x: []float64{1.00001, 10.00001, 1.00001, 5.00001, 4.00001},
	},
	{
		f:       functions.BiggsEXP6{},
		x:       []float64{1, 2, 1, 1, 1, 1},
		gradTol: 1e-7,
	},
	{
		f:       functions.BiggsEXP6{},
		x:       []float64{1.00001, 10.00001, 1.00001, 5.00001, 4.00001, 3.00001},
		gradTol: 1e-8,
	},
	{
		f: functions.Box3D{},
		x: []float64{0, 10, 20},
	},
	{
		f: functions.Box3D{},
		x: []float64{1.00001, 10.00001, 1.00001},
	},
	{
		f: functions.Box3D{},
		x: []float64{100.00001, 100.00001, 0.00001},
	},
	{
		f: functions.ExtendedPowellSingular{},
		x: []float64{3, -1, 0, 3},
	},
	{
		f: functions.ExtendedPowellSingular{},
		x: []float64{0.00001, 0.00001, 0.00001, 0.00001},
	},
	{
		f:       functions.ExtendedPowellSingular{},
		x:       []float64{3, -1, 0, 3, 3, -1, 0, 3},
		gradTol: 1e-8,
	},
	{
		f: functions.ExtendedPowellSingular{},
		x: []float64{0.00001, 0.00001, 0.00001, 0.00001, 0.00001, 0.00001, 0.00001, 0.00001},
	},
	{
		f: functions.ExtendedRosenbrock{},
		x: []float64{-1.2, 1, -1.2, 1},
	},
	{
		f: functions.ExtendedRosenbrock{},
		x: []float64{1e4, 1e4},
	},
	{
		f: functions.ExtendedRosenbrock{},
		x: []float64{1.00001, 1.00001, 1.00001, 1.00001},
	},
	{
		f:       functions.PenaltyI{},
		x:       []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		gradTol: 1e-10,
	},
	{
		f:       functions.PenaltyI{},
		x:       []float64{0.250007, 0.250007, 0.250007, 0.250007},
		gradTol: 1e-10,
	},
	{
		f: functions.PenaltyI{},
		x: []float64{0.1581, 0.1581, 0.1581, 0.1581, 0.1581, 0.1581,
			0.1581, 0.1581, 0.1581, 0.1581},
		gradTol: 1e-10,
	},
	{
		f:       functions.PenaltyII{},
		x:       []float64{0.5, 0.5, 0.5, 0.5},
		gradTol: 1e-8,
	},
	{
		f:       functions.PenaltyII{},
		x:       []float64{0.19999, 0.19131, 0.4801, 0.51884},
		gradTol: 1e-8,
	},
	{
		f: functions.PenaltyII{},
		x: []float64{0.19998, 0.01035, 0.01960, 0.03208, 0.04993, 0.07651,
			0.11862, 0.19214, 0.34732, 0.36916},
		gradTol: 1e-6,
	},
	{
		f:       functions.PowellBadlyScaled{},
		x:       []float64{1.09815e-05, 9.10614},
		gradTol: 1e-8,
	},
	newVariablyDimensioned(100, 1e-10),
	newVariablyDimensioned(1000, 1e-7),
	newVariablyDimensioned(10000, 1e-4),
	{
		f:       functions.Watson{},
		x:       []float64{0, 0, 0, 0, 0, 0},
		gradTol: 1e-7,
	},
	{
		f:       functions.Watson{},
		x:       []float64{-0.01572, 1.01243, -0.23299, 1.26043, -1.51372, 0.99299},
		gradTol: 1e-7,
	},
	{
		f:       functions.Watson{},
		x:       []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		gradTol: 1e-7,
		long:    true,
	},
	{
		f: functions.Watson{},
		x: []float64{-1.53070e-05, 0.99978, 0.01476, 0.14634, 1.00082,
			-2.61773, 4.10440, -3.14361, 1.05262},
		gradTol: 1e-7,
	},
	{
		f:       functions.Wood{},
		x:       []float64{-3, -1, -3, -1},
		gradTol: 1e-6,
	},
}

var newtonTests = []unconstrainedTest{
	{
		f: functions.BiggsEXP4{},
		x: []float64{1, 2, 1, 1},
	},
	{
		f: functions.BiggsEXP4{},
		x: []float64{1.00001, 10.00001, 1.00001, 5.00001},
	},
	{
		f: functions.BiggsEXP5{},
		x: []float64{1, 2, 1, 1, 1},
	},
	{
		f: functions.BiggsEXP5{},
		x: []float64{1.00001, 10.00001, 1.00001, 5.00001, 4.00001},
	},
	{
		f:       functions.BiggsEXP6{},
		x:       []float64{1, 2, 1, 1, 1, 1},
		gradTol: 1e-8,
	},
	{
		f:       functions.BiggsEXP6{},
		x:       []float64{1.00001, 10.00001, 1.00001, 5.00001, 4.00001, 3.00001},
		gradTol: 1e-8,
	},
	{
		f: functions.Box3D{},
		x: []float64{0, 10, 20},
	},
	{
		f: functions.Box3D{},
		x: []float64{1.00001, 10.00001, 1.00001},
	},
	{
		f: functions.Box3D{},
		x: []float64{100.00001, 100.00001, 0.00001},
	},
	{
		f: functions.BrownBadlyScaled{},
		x: []float64{1, 1},
	},
	{
		f: functions.BrownBadlyScaled{},
		x: []float64{1.000001e6, 2.01e-6},
	},
	{
		f: functions.ExtendedPowellSingular{},
		x: []float64{3, -1, 0, 3},
	},
	{
		f: functions.ExtendedPowellSingular{},
		x: []float64{0.00001, 0.00001, 0.00001, 0.00001},
	},
	{
		f: functions.ExtendedPowellSingular{},
		x: []float64{3, -1, 0, 3, 3, -1, 0, 3},
	},
	{
		f: functions.ExtendedPowellSingular{},
		x: []float64{0.00001, 0.00001, 0.00001, 0.00001, 0.00001, 0.00001, 0.00001, 0.00001},
	},
	{
		f: functions.ExtendedRosenbrock{},
		x: []float64{-1.2, 1, -1.2, 1},
	},
	{
		f: functions.ExtendedRosenbrock{},
		x: []float64{1.00001, 1.00001, 1.00001, 1.00001},
	},
	{
		f:       functions.Gaussian{},
		x:       []float64{0.4, 1, 0},
		gradTol: 1e-11,
	},
	{
		f: functions.GulfResearchAndDevelopment{},
		x: []float64{5, 2.5, 0.15},
	},
	{
		f: functions.GulfResearchAndDevelopment{},
		x: []float64{50.00001, 25.00001, 1.50001},
	},
	{
		f: functions.GulfResearchAndDevelopment{},
		x: []float64{99.89529, 60.61453, 9.16124},
	},
	{
		f: functions.GulfResearchAndDevelopment{},
		x: []float64{201.66258, 60.61633, 10.22489},
	},
	{
		f: functions.PenaltyI{},
		x: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	},
	{
		f: functions.PenaltyI{},
		x: []float64{0.250007, 0.250007, 0.250007, 0.250007},
	},
	{
		f: functions.PenaltyI{},
		x: []float64{0.1581, 0.1581, 0.1581, 0.1581, 0.1581, 0.1581,
			0.1581, 0.1581, 0.1581, 0.1581},
	},
	{
		f:       functions.PenaltyII{},
		x:       []float64{0.5, 0.5, 0.5, 0.5},
		gradTol: 1e-10,
	},
	{
		f:       functions.PenaltyII{},
		x:       []float64{0.19999, 0.19131, 0.4801, 0.51884},
		gradTol: 1e-10,
	},
	{
		f:       functions.PenaltyII{},
		x:       []float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
		gradTol: 1e-9,
	},
	{
		f: functions.PenaltyII{},
		x: []float64{0.19998, 0.01035, 0.01960, 0.03208, 0.04993, 0.07651,
			0.11862, 0.19214, 0.34732, 0.36916},
		gradTol: 1e-9,
	},
	{
		f: functions.PowellBadlyScaled{},
		x: []float64{0, 1},
	},
	{
		f:       functions.PowellBadlyScaled{},
		x:       []float64{1.09815e-05, 9.10614},
		gradTol: 1e-10,
	},
	newVariablyDimensioned(100, 1e-10),
	{
		f:       functions.Watson{},
		x:       []float64{0, 0, 0, 0, 0, 0},
		gradTol: 1e-7,
	},
	{
		f:       functions.Watson{},
		x:       []float64{-0.01572, 1.01243, -0.23299, 1.26043, -1.51372, 0.99299},
		gradTol: 1e-7,
	},
	{
		f:       functions.Watson{},
		x:       []float64{0, 0, 0, 0, 0, 0, 0, 0, 0},
		gradTol: 1e-8,
	},
	{
		f: functions.Watson{},
		x: []float64{-1.53070e-05, 0.99978, 0.01476, 0.14634, 1.00082,
			-2.61773, 4.10440, -3.14361, 1.05262},
		gradTol: 1e-8,
	},
}

var bfgsTests = []unconstrainedTest{
	{
		f:       functions.BiggsEXP6{},
		x:       []float64{1, 2, 1, 1, 1, 1},
		gradTol: 1e-10,
	},
	{
		f:       functions.BiggsEXP6{},
		x:       []float64{1.00001, 10.00001, 1.00001, 5.00001, 4.00001, 3.00001},
		gradTol: 1e-10,
	},
	{
		f:       functions.BrownAndDennis{},
		x:       []float64{25, 5, -5, -1},
		gradTol: 1e-5,
	},
	{
		f: functions.ExtendedRosenbrock{},
		x: []float64{1e5, 1e5},
	},
	{
		f:       functions.Gaussian{},
		x:       []float64{0.398, 1, 0},
		gradTol: 1e-11,
	},
	{
		f: functions.Wood{},
		x: []float64{-3, -1, -3, -1},
	},
}

var lbfgsTests = []unconstrainedTest{
	{
		f:       functions.BiggsEXP6{},
		x:       []float64{1, 2, 1, 1, 1, 1},
		gradTol: 1e-8,
	},
	{
		f:       functions.BiggsEXP6{},
		x:       []float64{1.00001, 10.00001, 1.00001, 5.00001, 4.00001, 3.00001},
		gradTol: 1e-8,
	},
	{
		f: functions.ExtendedRosenbrock{},
		x: []float64{1e7, 1e6},
	},
	{
		f:       functions.Gaussian{},
		x:       []float64{0.398, 1, 0},
		gradTol: 1e-10,
	},
	newVariablyDimensioned(1000, 1e-8),
	newVariablyDimensioned(10000, 1e-5),
}

func newVariablyDimensioned(dim int, gradTol float64) unconstrainedTest {
	x := make([]float64, dim)
	for i := range x {
		x[i] = float64(dim-i-1) / float64(dim)
	}
	return unconstrainedTest{
		f:       functions.VariablyDimensioned{},
		x:       x,
		gradTol: gradTol,
	}
}

func TestLocal(t *testing.T) {
	// TODO: When method is nil, Local chooses the method automatically. At
	// present, it always chooses BFGS (or panics if the function does not
	// implement Grad() or FuncGrad()). For now, run this test with the
	// simplest set of problems and revisit this later when more methods are
	// added.
	testLocal(t, gradientDescentTests, nil)
}

func TestGradientDescent(t *testing.T) {
	testLocal(t, gradientDescentTests, &GradientDescent{})
}

func TestGradientDescentBacktracking(t *testing.T) {
	testLocal(t, gradientDescentTests, &GradientDescent{
		LinesearchMethod: &Backtracking{
			FunConst: 0.1,
		},
	})
}

func TestGradientDescentBisection(t *testing.T) {
	testLocal(t, gradientDescentTests, &GradientDescent{
		LinesearchMethod: &Bisection{},
	})
}

func TestCG(t *testing.T) {
	var tests []unconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testLocal(t, tests, &CG{})
}

func TestFletcherReevesQuadStep(t *testing.T) {
	var tests []unconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testLocal(t, tests, &CG{
		Variant:     &FletcherReeves{},
		InitialStep: &QuadraticStepSize{},
	})
}

func TestFletcherReevesFirstOrderStep(t *testing.T) {
	var tests []unconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testLocal(t, tests, &CG{
		Variant:     &FletcherReeves{},
		InitialStep: &FirstOrderStepSize{},
	})
}

func TestHestenesStiefelQuadStep(t *testing.T) {
	var tests []unconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testLocal(t, tests, &CG{
		Variant:     &HestenesStiefel{},
		InitialStep: &QuadraticStepSize{},
	})
}

func TestHestenesStiefelFirstOrderStep(t *testing.T) {
	var tests []unconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testLocal(t, tests, &CG{
		Variant:     &HestenesStiefel{},
		InitialStep: &FirstOrderStepSize{},
	})
}

func TestPolakRibiereQuadStep(t *testing.T) {
	var tests []unconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testLocal(t, tests, &CG{
		Variant:     &PolakRibierePolyak{},
		InitialStep: &QuadraticStepSize{},
	})
}

func TestPolakRibiereFirstOrderStep(t *testing.T) {
	var tests []unconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testLocal(t, tests, &CG{
		Variant:     &PolakRibierePolyak{},
		InitialStep: &FirstOrderStepSize{},
	})
}

func TestDaiYuanQuadStep(t *testing.T) {
	var tests []unconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testLocal(t, tests, &CG{
		Variant:     &DaiYuan{},
		InitialStep: &QuadraticStepSize{},
	})
}

func TestDaiYuanFirstOrderStep(t *testing.T) {
	var tests []unconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testLocal(t, tests, &CG{
		Variant:     &DaiYuan{},
		InitialStep: &FirstOrderStepSize{},
	})
}

func TestHagerZhangQuadStep(t *testing.T) {
	var tests []unconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testLocal(t, tests, &CG{
		Variant:     &HagerZhang{},
		InitialStep: &QuadraticStepSize{},
	})
}

func TestHagerZhangFirstOrderStep(t *testing.T) {
	var tests []unconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testLocal(t, tests, &CG{
		Variant:     &HagerZhang{},
		InitialStep: &FirstOrderStepSize{},
	})
}

func TestBFGS(t *testing.T) {
	var tests []unconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, newtonTests...)
	tests = append(tests, bfgsTests...)
	testLocal(t, tests, &BFGS{})
}

func TestLBFGS(t *testing.T) {
	var tests []unconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, newtonTests...)
	tests = append(tests, lbfgsTests...)
	testLocal(t, tests, &LBFGS{})
}

func testLocal(t *testing.T, tests []unconstrainedTest, method Method) {
	for _, test := range tests {
		if test.long && testing.Short() {
			continue
		}

		settings := &Settings{
			FunctionThreshold: math.Inf(-1),
		}
		if test.gradTol == 0 {
			test.gradTol = 1e-12
		}
		settings.GradientThreshold = test.gradTol

		result, err := Local(test.f, test.x, settings, method)
		if err != nil {
			t.Errorf("error finding minimum (%v) for:\n%v", err, test)
			continue
		}

		if result == nil {
			t.Errorf("nil result without error for:\n%v", test)
			continue
		}

		funcInfo := newFunctionInfo(test.f)

		// Evaluate the norm of the gradient at the found optimum location.
		var optF, optNorm float64
		if funcInfo.IsFunctionGradient {
			g := make([]float64, len(test.x))
			optF = funcInfo.functionGradient.FuncGrad(result.X, g)
			optNorm = floats.Norm(g, math.Inf(1))
		} else {
			optF = funcInfo.function.Func(result.X)
			if funcInfo.IsGradient {
				g := make([]float64, len(test.x))
				funcInfo.gradient.Grad(result.X, g)
				optNorm = floats.Norm(g, math.Inf(1))
			}
		}

		// Check that the function value at the found optimum location is
		// equal to result.F
		if optF != result.F {
			t.Errorf("Function value at the optimum location %v not equal to the returned value %v for:\n%v",
				optF, result.F, test)
		}

		// Check that the norm of the gradient at the found optimum location is
		// smaller than the tolerance.
		if optNorm >= settings.GradientThreshold {
			t.Errorf("Norm of the gradient at the optimum location %v not smaller than tolerance %v for:\n%v",
				optNorm, settings.GradientThreshold, test)
		}

		// We are going to restart the solution using a fixed starting gradient
		// and value, so evaluate them.
		settings.UseInitialData = true
		if funcInfo.IsFunctionGradient {
			settings.InitialGradient = resize(settings.InitialGradient, len(test.x))
			settings.InitialFunctionValue = funcInfo.functionGradient.FuncGrad(test.x, settings.InitialGradient)
		} else {
			settings.InitialFunctionValue = funcInfo.function.Func(test.x)
			if funcInfo.IsGradient {
				settings.InitialGradient = resize(settings.InitialGradient, len(test.x))
				funcInfo.gradient.Grad(test.x, settings.InitialGradient)
			}
		}

		// Rerun the test again to make sure that it gets the same answer with
		// the same starting condition. Moreover, we are using the initial data
		// in settings.InitialFunctionValue and settings.InitialGradient.
		result2, err2 := Local(test.f, test.x, settings, method)
		if err2 != nil {
			t.Errorf("error finding minimum second time (%v) for:\n%v", err2, test)
			continue
		}

		if result2 == nil {
			t.Errorf("second time nil result without error for:\n%v", test)
			continue
		}

		// At the moment all the optimizers are deterministic, so check that we
		// get _exactly_ the same answer second time as well.
		if result.F != result2.F {
			t.Errorf("Different minimum second time. First: %v, Second: %v for:\n%v",
				result.F, result2.F, test)
		}

		// Check that providing initial data reduces the number of function
		// and/or gradient calls exactly by one.
		if funcInfo.IsFunctionGradient {
			if result.FuncGradEvaluations != result2.FuncGradEvaluations+1 {
				t.Errorf("Providing initial data does not reduce the number of function/gradient calls for:\n%v", test)
				continue
			}
		} else {
			if result.FuncEvaluations != result2.FuncEvaluations+1 {
				t.Errorf("Providing initial data does not reduce the number of functions calls for:\n%v", test)
				continue
			}
			if funcInfo.IsGradient {
				if result.GradEvaluations != result2.GradEvaluations+1 {
					t.Errorf("Providing initial data does not reduce the number of gradient calls for:\n%v", test)
					continue
				}
			}
		}
	}
}

func TestIssue76(t *testing.T) {
	f := functions.BrownAndDennis{}
	// Location very close to the minimum.
	x := []float64{-11.594439904886773, 13.203630051265385, -0.40343948776868443, 0.2367787746745986}
	s := &Settings{
		FunctionThreshold: math.Inf(-1),
		GradientThreshold: 1e-14,
		MajorIterations:   1000000,
	}
	m := &GradientDescent{
		LinesearchMethod: &Backtracking{},
	}
	// We are not interested in the error, only in the returned status.
	r, _ := Local(f, x, s, m)
	// With the above stringent tolerance, the optimizer will never
	// successfully reach the minimum. Check if it terminated in a finite
	// number of steps.
	if r.Status == IterationLimit {
		t.Error("Issue https://github.com/gonum/optimize/issues/76 not fixed")
	}
}
