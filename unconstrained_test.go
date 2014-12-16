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

// The Biggs' EXP2 function
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

// The Biggs' EXP3 function
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

// The Biggs' EXP4 function
// M.C. Biggs, Minimization algorithms making use of non-quadratic properties
// of the objective function. J. Inst. Maths Applics 8 (1971), 315-327.
// Dim = 4
// X0 = [1, 2, 1, 1]
// OptX = [1, 10, 1, 5]
// OptF = 0
type BiggsEXP4 struct{}

func (BiggsEXP4) F(x []float64) (sum float64) {
	for i := 1; i <= 10; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z)
		f := x[2]*math.Exp(-x[0]*z) - x[3]*math.Exp(-x[1]*z) - y
		sum += math.Pow(f, 2)
	}
	return sum
}

func (BiggsEXP4) Df(x, g []float64) {
	for i := 0; i < len(g); i++ {
		g[i] = 0
	}
	for i := 1; i <= 10; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z)
		f := x[2]*math.Exp(-x[0]*z) - x[3]*math.Exp(-x[1]*z) - y

		dfdx0 := -z * x[2] * math.Exp(-x[0]*z)
		dfdx1 := z * x[3] * math.Exp(-x[1]*z)
		dfdx2 := math.Exp(-x[0] * z)
		dfdx3 := -math.Exp(-x[1] * z)

		g[0] += 2 * f * dfdx0
		g[1] += 2 * f * dfdx1
		g[2] += 2 * f * dfdx2
		g[3] += 2 * f * dfdx3
	}
}

// The Biggs' EXP5 function
// M.C. Biggs, Minimization algorithms making use of non-quadratic properties
// of the objective function. J. Inst. Maths Applics 8 (1971), 315-327.
// Dim = 5
// X0 = [1, 2, 1, 1, 1]
// OptX = [1, 10, 1, 5, 4]
// OptF = 0
type BiggsEXP5 struct{}

func (BiggsEXP5) F(x []float64) (sum float64) {
	for i := 1; i <= 11; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z) + 3*math.Exp(-4*z)
		f := x[2]*math.Exp(-x[0]*z) - x[3]*math.Exp(-x[1]*z) + 3*math.Exp(-x[4]*z) - y
		sum += math.Pow(f, 2)
	}
	return sum
}

func (BiggsEXP5) Df(x, g []float64) {
	for i := 0; i < len(g); i++ {
		g[i] = 0
	}
	for i := 1; i <= 11; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z) + 3*math.Exp(-4*z)
		f := x[2]*math.Exp(-x[0]*z) - x[3]*math.Exp(-x[1]*z) + 3*math.Exp(-x[4]*z) - y

		dfdx0 := -z * x[2] * math.Exp(-x[0]*z)
		dfdx1 := z * x[3] * math.Exp(-x[1]*z)
		dfdx2 := math.Exp(-x[0] * z)
		dfdx3 := -math.Exp(-x[1] * z)
		dfdx4 := -3 * z * math.Exp(-x[4]*z)

		g[0] += 2 * f * dfdx0
		g[1] += 2 * f * dfdx1
		g[2] += 2 * f * dfdx2
		g[3] += 2 * f * dfdx3
		g[4] += 2 * f * dfdx4
	}
}

// The Biggs' EXP6 function
// M.C. Biggs, Minimization algorithms making use of non-quadratic properties
// of the objective function. J. Inst. Maths Applics 8 (1971), 315-327.
// Dim = 6
// X0 = [1, 2, 1, 1, 1, 1]
// OptX = [1, 10, 1, 5, 4, 3]
// OptF = 0
// OptF = 0.005655649925...
type BiggsEXP6 struct{}

func (BiggsEXP6) F(x []float64) (sum float64) {
	for i := 1; i <= 13; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z) + 3*math.Exp(-4*z)
		f := x[2]*math.Exp(-x[0]*z) - x[3]*math.Exp(-x[1]*z) + x[5]*math.Exp(-x[4]*z) - y
		sum += math.Pow(f, 2)
	}
	return sum
}

func (BiggsEXP6) Df(x, g []float64) {
	for i := 0; i < len(g); i++ {
		g[i] = 0
	}
	for i := 1; i <= 13; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z) + 3*math.Exp(-4*z)
		f := x[2]*math.Exp(-x[0]*z) - x[3]*math.Exp(-x[1]*z) + x[5]*math.Exp(-x[4]*z) - y

		dfdx0 := -z * x[2] * math.Exp(-x[0]*z)
		dfdx1 := z * x[3] * math.Exp(-x[1]*z)
		dfdx2 := math.Exp(-x[0] * z)
		dfdx3 := -math.Exp(-x[1] * z)
		dfdx4 := -z * x[5] * math.Exp(-x[4]*z)
		dfdx5 := math.Exp(-x[4] * z)

		g[0] += 2 * f * dfdx0
		g[1] += 2 * f * dfdx1
		g[2] += 2 * f * dfdx2
		g[3] += 2 * f * dfdx3
		g[4] += 2 * f * dfdx4
		g[5] += 2 * f * dfdx5
	}
}

// Gaussian function
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = 3
// X0 = [0.4, 1, 0]
// OptX = [0.3989561..., 1.0000191..., 0]
// OptF = 1.12793...e-8
type Gaussian struct{}

func (Gaussian) y(i int) (yi float64) {
	switch i {
	case 1, 15:
		yi = 0.0009
	case 2, 14:
		yi = 0.0044
	case 3, 13:
		yi = 0.0175
	case 4, 12:
		yi = 0.0540
	case 5, 11:
		yi = 0.1295
	case 6, 10:
		yi = 0.2420
	case 7, 9:
		yi = 0.3521
	case 8:
		yi = 0.3989
	}
	return yi
}

func (g Gaussian) F(x []float64) (sum float64) {
	for i := 1; i <= 15; i++ {
		c := float64(8-i) / 2
		d := math.Pow(c-x[2], 2)
		e := math.Exp(-x[1] * d / 2)
		f := x[0]*e - g.y(i)
		sum += f * f
	}
	return sum
}

func (g Gaussian) Df(x, grad []float64) {
	grad[0] = 0
	grad[1] = 0
	grad[2] = 0
	for i := 1; i <= 15; i++ {
		c := float64(8-i) / 2
		d := math.Pow(c-x[2], 2)
		e := math.Exp(-x[1] * d / 2)
		f := x[0]*e - g.y(i)

		grad[0] += 2 * f * e
		grad[1] -= f * e * d * x[0]
		grad[2] += 2 * f * e * x[0] * x[1] * (c - x[2])
	}
}

// The Powell's badly scaled function
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = 2
// X0 = [0, 1]
// OptX = [1.09815933...e-5, 9.10614674...]
// OptF = 0
type Powell struct{}

func (Powell) F(x []float64) float64 {
	f1 := 1e4*x[0]*x[1] - 1
	f2 := math.Exp(-x[0]) + math.Exp(-x[1]) - 1.0001
	return math.Pow(f1, 2) + math.Pow(f2, 2)
}

func (Powell) Df(x, grad []float64) {
	f1 := 1e4*x[0]*x[1] - 1
	f2 := math.Exp(-x[0]) + math.Exp(-x[1]) - 1.0001

	grad[0] = 2 * (1e4*f1*x[1] - f2*math.Exp(-x[0]))
	grad[1] = 2 * (1e4*f1*x[0] - f2*math.Exp(-x[1]))
}

// The Box' three-dimensional function
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = 3
// X0 = [0, 10, 20]
// OptX = [1, 10, 1], [10, 1, -1], [a, a, 0]
// OptF = 0
type Box struct{}

func (Box) F(x []float64) (sum float64) {
	for i := 1; i <= 10; i++ {
		c := float64(i) / 10
		y := math.Exp(-c) - math.Exp(10*c)
		f := math.Exp(-c*x[0]) - math.Exp(-c*x[1]) - x[2]*y
		sum += math.Pow(f, 2)
	}
	return sum
}

func (Box) Df(x, grad []float64) {
	grad[0] = 0
	grad[1] = 0
	grad[2] = 0

	for i := 1; i <= 10; i++ {
		c := float64(i) / 10
		y := math.Exp(-c) - math.Exp(10*c)
		f := math.Exp(-c*x[0]) - math.Exp(-c*x[1]) - x[2]*y

		grad[0] += -2 * f * c * math.Exp(-c*x[0])
		grad[1] += -2 * f * c * math.Exp(-c*x[1])
		grad[2] += -2 * f * y
	}
}

// Variably dimensioned function
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = n
// X0 = [..., (n-i)/n, ...]
// OptX = [1, ..., 1]
// OptF = 0
type VariablyDimensioned struct{}

func (v VariablyDimensioned) F(x []float64) (sum float64) {
	for i := 0; i < len(x); i++ {
		sum += math.Pow(x[i]-1, 2)
	}
	s := 0.0
	for i := 0; i < len(x); i++ {
		s += float64(i+1) * (x[i] - 1)
	}
	sum += math.Pow(s, 2) + math.Pow(s, 4)
	return sum
}

func (v VariablyDimensioned) Df(x, grad []float64) {
	for i := 0; i < len(grad); i++ {
		grad[i] = 0
	}
	s := 0.0
	for i := 0; i < len(x); i++ {
		s += float64(i+1) * (x[i] - 1)
	}
	for i := 0; i < len(grad); i++ {
		grad[i] = 2 * ((x[i] - 1) + s*float64(i+1)*(1+2*math.Pow(s, 2)))
	}
}

// The Watson function
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// For Dim = 9, the problem of minimizing the Watson function is very ill conditioned.
// Dim = n, 2 <= n <= 31
// X0 = [0, ..., 0]
// OptX = [1, ..., 1]
// For Dim = 6 also:
// OptX = [-0.015725, 1.012435, -0.232992, 1.260430, -1.513729, 0.992996]
// OptF = 2.287687...e-3
// For Dim = 9 also:
// OptX = [-0.000015, 0.999790, 0.014764, 0.146342, 1.000821, -2.617731, 4.104403, -3.143612, 1.052627]
// OptF = 1.39976...e-6
// For Dim = 12 also:
// OptF = 4.72238...e-10
type Watson struct{}

func (Watson) F(x []float64) (sum float64) {
	for i := 1; i <= 29; i++ {
		c := float64(i) / 29

		s1 := 0.0
		for j := 1; j < len(x); j++ {
			s1 += float64(j) * x[j] * math.Pow(c, float64(j-1))
		}

		s2 := 0.0
		for j := 0; j < len(x); j++ {
			s2 += x[j] * math.Pow(c, float64(j))
		}
		s2 = math.Pow(s2, 2)

		sum += math.Pow(s1-s2-1, 2)
	}
	sum += math.Pow(x[0], 2)
	sum += math.Pow(x[1]-math.Pow(x[0], 2)-1, 2)
	return sum
}

func (Watson) Df(x, grad []float64) {
	for i := 0; i < len(grad); i++ {
		grad[i] = 0
	}

	for i := 1; i <= 29; i++ {
		c := float64(i) / 29

		s1 := 0.0
		for j := 1; j < len(x); j++ {
			s1 += float64(j) * x[j] * math.Pow(c, float64(j-1))
		}

		s2 := 0.0
		for j := 0; j < len(x); j++ {
			s2 += x[j] * math.Pow(c, float64(j))
		}

		t := s1 - math.Pow(s2, 2) - 1
		for j := 0; j < len(x); j++ {
			grad[j] += 2 * t * math.Pow(c, float64(j-1)) * (float64(j) - 2*s2*c)
		}
	}
	t := x[1] - math.Pow(x[0], 2) - 1
	grad[0] += 2 * (1 - 2*t) * x[0]
	grad[1] += 2 * t
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
