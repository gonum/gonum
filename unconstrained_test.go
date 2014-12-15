// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"
	"testing"

	"github.com/gonum/blas/goblas"
	"github.com/gonum/floats"
	"github.com/gonum/matrix/mat64"
)

var negInf = math.Inf(-1)

func init() {
	mat64.Register(goblas.Blas{})
}

type Rosenbrock struct {
	nDim int
}

func (r Rosenbrock) F(x []float64) (sum float64) {
	deriv := make([]float64, len(x))
	return r.FDf(x, deriv)
}

func (r Rosenbrock) FDf(x []float64, deriv []float64) (sum float64) {
	sum = 0

	for i := range deriv {
		deriv[i] = 0
	}

	for i := 0; i < len(x)-1; i++ {
		sum += math.Pow(1-x[i], 2) + 100*math.Pow(x[i+1]-math.Pow(x[i], 2), 2)
	}
	for i := 0; i < len(x)-1; i++ {
		deriv[i] += -1 * 2 * (1 - x[i])
		deriv[i] += 2 * 100 * (x[i+1] - math.Pow(x[i], 2)) * (-2 * x[i])
	}
	for i := 1; i < len(x); i++ {
		deriv[i] += 2 * 100 * (x[i] - math.Pow(x[i-1], 2))
	}

	return sum
}

// The Fletcher-Powell helical valley function
// Dim = 3
// X0 = [-1, 0, 0]
// OptX = [1, 0, 0]
// OptF = 0
type HelicalValley struct{}

func (HelicalValley) F(x []float64) float64 {
	θ := 0.5 * math.Atan2(x[1], x[0]) / math.Pi
	r := math.Sqrt(math.Pow(x[0], 2) + math.Pow(x[1], 2))

	f1 := 10 * (x[2] - 10*θ)
	f2 := 10 * (r - 1)
	f3 := x[2]

	return math.Pow(f1, 2) + math.Pow(f2, 2) + math.Pow(f3, 2)
}

func (HelicalValley) Df(x, g []float64) {
	θ := 0.5 * math.Atan2(x[1], x[0]) / math.Pi
	r := math.Sqrt(math.Pow(x[0], 2) + math.Pow(x[1], 2))
	s := x[2] - 10*θ
	t := 5 * s / math.Pow(r, 2) / math.Pi

	g[0] = 200 * (x[0] - x[0]/r + x[1]*t)
	g[1] = 200 * (x[1] - x[1]/r - x[0]*t)
	g[2] = 2 * (x[2] + 100*s)
}

// Biggs' EXP2 function
// M.C. Biggs, Minimization algorithms making use of non-quadratic properties
// of the objective function. J. Inst. Maths Applics 8 (1971), 315-327.
// Dim = 2
// X0 = [1, 2]
// OptX = [1, 10]
// OptF = 0
type BiggsEXP2 struct{}

func (BiggsEXP2) F(x []float64) (sum float64) {
	for i := 1; i <= 10; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z)
		f := math.Exp(-x[0]*z) - 5*math.Exp(-x[1]*z) - y
		sum += math.Pow(f, 2)
	}
	return sum
}

func (BiggsEXP2) Df(x, g []float64) {
	for i := 0; i < len(g); i++ {
		g[i] = 0
	}
	for i := 1; i <= 10; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z)
		f := math.Exp(-x[0]*z) - 5*math.Exp(-x[1]*z) - y

		dfdx0 := -z * math.Exp(-x[0]*z)
		dfdx1 := 5 * z * math.Exp(-x[1]*z)

		g[0] += 2 * f * dfdx0
		g[1] += 2 * f * dfdx1
	}
}

// Biggs' EXP3 function
// M.C. Biggs, Minimization algorithms making use of non-quadratic properties
// of the objective function. J. Inst. Maths Applics 8 (1971), 315-327.
// Dim = 3
// X0 = [1, 2, 1]
// OptX = [1, 10, 5]
// OptF = 0
type BiggsEXP3 struct{}

func (BiggsEXP3) F(x []float64) (sum float64) {
	for i := 1; i <= 10; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z)
		f := math.Exp(-x[0]*z) - x[2]*math.Exp(-x[1]*z) - y
		sum += math.Pow(f, 2)
	}
	return sum
}

func (BiggsEXP3) Df(x, g []float64) {
	for i := 0; i < len(g); i++ {
		g[i] = 0
	}
	for i := 1; i <= 10; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z)
		f := math.Exp(-x[0]*z) - x[2]*math.Exp(-x[1]*z) - y

		dfdx0 := -z * math.Exp(-x[0]*z)
		dfdx1 := x[2] * z * math.Exp(-x[1]*z)
		dfdx2 := -math.Exp(-x[1] * z)

		g[0] += 2 * f * dfdx0
		g[1] += 2 * f * dfdx1
		g[2] += 2 * f * dfdx2
	}
}

type Linear struct {
	nDim int
}

func (l Linear) F(x []float64) float64 {
	return floats.Sum(x)
}

func (l Linear) FDf(x []float64, deriv []float64) float64 {
	for i := range deriv {
		deriv[i] = 1
	}
	return floats.Sum(x)
}

func TestMinimize(t *testing.T) {
	testMinimize(t, nil)
}

func TestGradientDescent(t *testing.T) {
	testMinimize(t, &GradientDescent{})
}

func TestGradientDescentBacktracking(t *testing.T) {
	testMinimize(t, &GradientDescent{
		LinesearchMethod: &Backtracking{
			FunConst: 0.1,
		},
	})
}

func TestGradientDescentBisection(t *testing.T) {
	testMinimize(t, &GradientDescent{
		LinesearchMethod: &Bisection{},
	})
}

func TestBFGS(t *testing.T) {
	testMinimize(t, &BFGS{})
}

func TestLBFGS(t *testing.T) {
	testMinimize(t, &LBFGS{})
}

func testMinimize(t *testing.T, method Method) {
	// This should be replaced with a more general testing framework with
	// a plugable method

	for i, test := range []struct {
		F Function
		X []float64

		OptVal float64
		OptLoc []float64

		Tol      float64
		Settings *Settings
	}{
		{
			F:      Rosenbrock{2},
			X:      []float64{15, 10},
			OptVal: 0,
			OptLoc: []float64{1, 1},
			Tol:    1e-4,

			Settings: DefaultSettings(),
		},
		{
			F:      Rosenbrock{4},
			X:      []float64{-150, 100, 5, -6},
			OptVal: 0,
			OptLoc: []float64{1, 1, 1, 1},
			Tol:    1e-4,

			Settings: &Settings{
				FunctionAbsTol: math.Inf(-1),
				GradientAbsTol: 1e-12,
			},
		},
		{
			F:      Rosenbrock{2},
			X:      []float64{15, 10},
			OptVal: 0,
			OptLoc: []float64{1, 1},
			Tol:    1e-4,

			Settings: &Settings{
				FunctionAbsTol: math.Inf(-1),
				GradientAbsTol: 1e-12,
			},
		},
		{
			F:      Rosenbrock{2},
			X:      []float64{-1.2, 1},
			OptVal: 0,
			OptLoc: []float64{1, 1},
			Tol:    1e-4,

			Settings: &Settings{
				FunctionAbsTol: math.Inf(-1),
				GradientAbsTol: 1e-3,
			},
		},
		/*
			// TODO: Turn this on when we have an adaptive linsearch method.
			// Gradient descent with backtracking will basically never finish
			{
				F:      Linear{8},
				X:      []float64{9, 8, 7, 6, 5, 4, 3, 2},
				OptVal: negInf,
				OptLoc: []float64{negInf, negInf, negInf, negInf, negInf, negInf, negInf, negInf},

				Settings: &Settings{
					FunctionAbsTol: math.Inf(-1),
				},
			},
		*/
	} {
		test.Settings.Recorder = nil
		result, err := Local(test.F, test.X, test.Settings, method)
		if err != nil {
			t.Errorf("error finding minimum: %v", err.Error())
			continue
		}
		// fmt.Println("%#v\n", result) // for debugging
		// TODO: Better tests
		if math.Abs(result.F-test.OptVal) > test.Tol {
			t.Errorf("Case %v: Minimum not found, exited with status: %v. Want: %v, Got: %v", i, result.Status, test.OptVal, result.F)
			continue
		}
		if result == nil {
			t.Errorf("Case %v: nil result without error", i)
			continue
		}

		// rerun it again to ensure it gets the same answer with the same starting
		// condition
		result2, err2 := Local(test.F, test.X, test.Settings, method)
		if err2 != nil {
			t.Errorf("error finding minimum second time: %v", err2.Error())
			continue
		}
		if result2 == nil {
			t.Errorf("Case %v: nil result without error", i)
			continue
		}
		/*
			// For debugging purposes, can't use DeepEqual naively becaus of NaNs
			// kill the runtime before the check, because those don't need to be equal
			result.Runtime = 0
			result2.Runtime = 0
			if !reflect.DeepEqual(result, result2) {
				t.Error(eqString)
				continue
			}
		*/
	}
}
