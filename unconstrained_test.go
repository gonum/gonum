// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

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
	//	fmt.Println("sum ", sum, "norm", floats.Norm(deriv, 2)/math.Sqrt(float64(len(deriv))))

	return sum
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
				FunctionAbsoluteTolerance: math.Inf(-1),
				GradientAbsoluteTolerance: 1e-13,
			},
		},
		{
			F:      Rosenbrock{2},
			X:      []float64{15, 10},
			OptVal: 0,
			OptLoc: []float64{1, 1},
			Tol:    1e-4,

			Settings: &Settings{
				FunctionAbsoluteTolerance: math.Inf(-1),
				GradientAbsoluteTolerance: 1e-13,
			},
		},
		/*
			// TODO: Fix optimizers so that they don't fail on this
			{
				F:      Linear{8},
				X:      []float64{9, 8, 7, 6, 5, 4, 3, 2},
				OptVal: negInf,
				OptLoc: []float64{negInf, negInf, negInf, negInf, negInf, negInf, negInf, negInf},

				Settings: &Settings{
					FunctionAbsoluteTolerance: math.Inf(-1),
				},
			},
		*/
	} {
		test.Settings.Recorder = nil
		result, err := Minimize(test.F, test.X, test.Settings, method)
		if err != nil {
			t.Errorf("error finding minimum: %v", err.Error())
			continue
		}
		// TODO: Better tests
		if math.Abs(result.F-test.OptVal) > test.Tol {
			t.Errorf("Minimum not found, exited with status: %v. Want: %v, Got: %v", result.Status, test.OptVal, result.F)
			continue
		}
		if result == nil {
			t.Errorf("Case %v: nil result without error", i)
			continue
		}

		// rerun it again to ensure it gets the same answer with the same starting
		// condition
		result2, err2 := Minimize(test.F, test.X, test.Settings, method)
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
