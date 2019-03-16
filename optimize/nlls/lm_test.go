// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nlls

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

// LeastSquares is a type for solving linear least squares problems with LM.
type LeastSquares struct {
	X *mat.Dense
	Y []float64
}

func (l LeastSquares) Problem() LMProblem {
	r, c := l.X.Dims()
	return LMProblem{
		Dim:  c,
		Size: r,
		Func: l.Func,
		Jac:  l.Jac,
	}
}

func (l LeastSquares) Optimal() []float64 {
	r, c := l.X.Dims()
	if len(l.Y) != r {
		panic("size mismatch")
	}
	yVec := mat.NewVecDense(len(l.Y), l.Y)

	tmp := make([]float64, c)
	tmpVec := mat.NewVecDense(len(tmp), tmp)

	err := tmpVec.SolveVec(l.X, yVec)
	if err != nil {
		panic("singular")
	}
	return tmp
}

func (l LeastSquares) Func(dst, params []float64) {
	l.funcJac(nil, dst, params)
}

func (l LeastSquares) Jac(dst *mat.Dense, params []float64) {
	l.funcJac(dst, nil, params)
}

func (l LeastSquares) funcJac(jacDst *mat.Dense, funDst, params []float64) {
	if funDst != nil {
		for i := 0; i < len(funDst); i++ {
			x := l.X.RawRowView(i)
			diff := floats.Dot(x, params) - l.Y[i]
			funDst[i] = diff
		}
	}
	if jacDst != nil {
		jacDst.Copy(l.X)
	}
}

func constructLeastSquares(trueParam []float64, noise float64, offset bool, nData int, source rand.Source) *LeastSquares {
	norm := rand.New(source).NormFloat64
	dim := len(trueParam)
	xs := mat.NewDense(nData, len(trueParam), nil)
	ys := make([]float64, nData)
	for i := 0; i < nData; i++ {
		if offset {
			xs.Set(i, 0, 1)
		} else {
			xs.Set(i, 0, norm())
		}
		for j := 1; j < dim; j++ {
			xs.Set(i, j, norm())
		}

		x := xs.RawRowView(i)
		y := floats.Dot(trueParam, x) + distuv.Normal{Mu: 0, Sigma: noise, Src: source}.Rand()
		ys[i] = y
	}
	return &LeastSquares{
		X: xs,
		Y: ys,
	}
}

// Powell, M. J. D. "A Hybrid Method for Nonlinear Equations", in P. Rabinowitz, ed.,
// "Numerical Methods for Nonlinear Algebraic Equations", Gordon and Breach, 1970.
func powellFunc(dst, x []float64) {
	dst[0] = x[0]
	dst[1] = 10*x[0]/(x[0]+0.1) + 2*x[1]*x[1]
}

func powellJac(dst *mat.Dense, x []float64) {
	dst.Set(0, 0, 1)
	dst.Set(0, 1, 0)
	dst.Set(1, 0, math.Pow(x[0]+0.1, -2))
	dst.Set(1, 1, 4*x[1])
}

// The following test functions are taken form:
// - More, J., Garbow, B.S., Hillstrom, K.E.: Testing unconstrained optimization software. 
//   ACM Trans Math Softw 7 (1981), 17-41.
func bealeFunc(dst, x []float64) {
	dst[0] = 1.5 - x[0]*(1-x[1])
	dst[1] = 2.25 - x[0]*(1-math.Pow(x[1], 2))
	dst[2] = 2.625 - x[0]*(1-math.Pow(x[1], 3))
}

func biggsEXP6Func(dst, x []float64) {
	for i := 0; i < 13; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z) + 3*math.Exp(-4*z)
		dst[i] = x[2]*math.Exp(-x[0]*z) - x[3]*math.Exp(-x[1]*z) + x[5]*math.Exp(-x[4]*z) - y
	}
}

func extendedRosenbrockFunc(dst, x []float64) {
	dst[0] = 10*(x[1]-x[0]*x[0])
	dst[1] = 1-x[0]
}

type LMTest struct {
	prob     LMProblem
	expected []float64
	tol      float64
	settings Settings
}

func TestLM(t *testing.T) {
	// Constract a linear least squares probelm.
	trueParam := []float64{0.7, 0.8, 5.6, -40.8}
	nData := 50
	noise := 1e-2
	ls := constructLeastSquares(trueParam, noise, true, nData, rand.NewSource(1))
	optLSParams := ls.Optimal()

	// Numerical Jacobians for problems.
	bealeNumJac := NumJac{Func: bealeFunc}
	biggsNumJac := NumJac{Func: biggsEXP6Func}
	rosenbrockNumJac := NumJac{Func: extendedRosenbrockFunc}

	problems := []LMTest{
		// Simple linear fit problem.
		LMTest{
			prob: LMProblem{
				Dim:        len(trueParam),
				Size:       nData,
				Func:       ls.Func,
				Jac:        ls.Jac,
				InitParams: nil,
				Tau:        1e-3,
				Eps1:       1e-8,
				Eps2:       1e-8,
			},
			expected: optLSParams,
			tol:      1e-6,
			settings: Settings{Iterations: 100, ObjectiveTol: 1e-16},
		},
		// Powell problem.
		LMTest{
			prob: LMProblem{
				Dim:        2,
				Size:       2,
				Func:       powellFunc,
				Jac:        powellJac,
				InitParams: []float64{3, 1},
				Tau:        1,
				Eps1:       1e-15,
				Eps2:       1e-15,
			},
			expected: []float64{0.0, 0.0},
			tol:      1e-2,
			settings: Settings{Iterations: 100, ObjectiveTol: 1e-16},
		},
		// Beale problem.
		LMTest{
			prob: LMProblem{
				Dim:        2,
				Size:       3,
				Func:       bealeFunc,
				Jac:        bealeNumJac.Jac,
				InitParams: []float64{1, 1},
				Tau:        1,
				Eps1:       1e-15,
				Eps2:       1e-15,
			},
			expected: []float64{3.0, 0.5},
			tol:      1e-5,
			settings: Settings{Iterations: 100, ObjectiveTol: 1e-16},
		},
		// Biggs EXP6 problem.
		LMTest{
			prob: LMProblem{
				Dim:        6,
				Size:       13,
				Func:       biggsEXP6Func,
				Jac:        biggsNumJac.Jac,
				InitParams: []float64{1, 2, 1, 1, 1, 1},
				Tau:        1e-6,
				Eps1:       1e-8,
				Eps2:       1e-8,
			},
			expected: []float64{1, 10, 1, 5, 4, 3},
			tol:      1e-3,
			settings: Settings{Iterations: 100, ObjectiveTol: 1e-16},
		},
		// Extended Rosenbrock problem.
		LMTest{
			prob: LMProblem{
				Dim:        2,
				Size:       2,
				Func:       extendedRosenbrockFunc,
				Jac:        rosenbrockNumJac.Jac,
				InitParams: []float64{-20, 150},
				Tau:        1e-6,
				Eps1:       1e-8,
				Eps2:       1e-8,
			},
			expected: []float64{1, 1},
			tol:      1e-6,
			settings: Settings{Iterations: 100, ObjectiveTol: 1e-16},
		},
	}

	for _, testProb := range problems {
		result, err := LM(testProb.prob, &testProb.settings)
		if err != nil {
			t.Errorf("unexepected error: %v", err)
		}
		if !floats.EqualApprox(result.X, testProb.expected, testProb.tol) {
			t.Errorf("Optimal mismatch: got %v, want %v", result.X, testProb.expected)
		}
	}
}
