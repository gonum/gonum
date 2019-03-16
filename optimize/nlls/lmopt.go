// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package nlls implements optimization routines for non-linear least squares problems
using the Levenberg-Marquardt method. 

Given function f:Rn -> Rm, where m is the number of non-linear functions and n parameters,
the Levenberg-Marquardt method is used to seek a point X that minimizes F(x) = 0.5 * f.T * f. 

The user supplies a non-linear function. The jacobian may also be supplied by the user or
approximated by finite differences.
*/
package nlls

import (
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/optimize"
)

type Settings struct {
	// Iterations represents the maximum number of iterations. Defaults to 100.
	Iterations int

	// ObjectiveTol represents the value for the obejective function after which
	// the algorithm can stop. Defaults to 1e-16.
	ObjectiveTol float64
}

func defaultSettings(set *Settings) {
	set.Iterations = 100
	set.ObjectiveTol = 1e-16
}

type Result struct {
	X      []float64
	Status optimize.Status
}

// NumJac is used if the user doesn't wish to provide a fucnction that evaluates 
// the jacobian matrix. NumJac provides a method Jac that computes the jacobian matrix
// by finite differences.
type NumJac struct {
	Func func(dst, param []float64)
}

func (nj *NumJac) Jac(dst *mat.Dense, param []float64) {
	fd.Jacobian(dst, nj.Func, param, &fd.JacobianSettings{
		Formula:    fd.Central,
		Concurrent: true,
	})
}

func maxDiagElem(m *mat.Dense) float64 {
	r, c := m.Dims()
	if r != c {
		panic("lm: matrix is not square")
	}
	maxElem := m.At(0, 0)
	for i := 1; i < r; i++ {
		if m.At(i, i) > maxElem {
			maxElem = m.At(i, i)
		}
	}
	return maxElem
}

func addToDiag(m *mat.Dense, v float64) {
	r, c := m.Dims()
	if r != c {
		panic("lm: matrix is not square")
	}
	for i := 0; i < r; i++ {
		m.Set(i, i, m.At(i, i)+v)
	}
}

func updateParams(dst []float64, params []float64, h *mat.VecDense) {
	if len(params) != h.Len() {
		panic("lm: lenghts don't match")
	}
	for i := 0; i < len(params); i++ {
		dst[i] = params[i] - h.At(i, 0)
	}
}

func calcRho(fParams []float64, fParamsNew []float64, h *mat.VecDense, grad *mat.VecDense, mu float64) float64 {
	rho := floats.Dot(fParams, fParams) - floats.Dot(fParamsNew, fParamsNew)
	tmpVec := mat.NewVecDense(h.Len(), nil)
	tmpVec.AddScaledVec(grad, mu, h)
	lDiff := mat.Dot(h, tmpVec)
	rho /= lDiff
	return rho
}

// LM is a function that solves non-linear least squares problems using the Levenberg-Marquardt
// Method. 
//
// References:
//  - Madsen, Kaj, Hans Bruun Nielsen, and Ole Tingleff. "Methods for non-linear least squares
//    problems.", 2nd edition, 2004.
//  - Lourakis, Manolis. "A Brief Description of the Levenberg-Marquardt Algorithm Implemened 
//    by levmar", 2005.
func LM(problem LMProblem, settings *Settings) (*Result, error) {
	var set Settings
	if settings != nil {
		set = *settings
	} else {
		defaultSettings(&set)
	}
	dim := problem.Dim
	if problem.Dim == 0 {
		panic("lm: problem dimension is 0")
	}
	size := problem.Size
	if problem.Size == 0 {
		panic("lm: problem size is 0")
	}
	status := optimize.NotTerminated

	dstFunc := make([]float64, size)
	dstFuncNew := make([]float64, size)
	dstJac := mat.NewDense(size, dim, nil)
	dstA := mat.NewDense(dim, dim, nil)
	dstGrad := mat.NewVecDense(dim, nil)
	dstH := mat.NewVecDense(dim, nil)
	nu := 2.0
	var mu float64
	found := false

	// The inital guess is the zero vector by default.
	parameters := make([]float64, dim)
	parametersNew := make([]float64, dim)
	if problem.InitParams != nil {
		copy(parameters, problem.InitParams)
	}

	// Initial evaluation of A = J.T * J and g = J.T * f.
	problem.Func(dstFunc, parameters)
	problem.Jac(dstJac, parameters)
	dstA.Mul(dstJac.T(), dstJac)
	dstGrad.MulVec(dstJac.T(), mat.NewVecDense(size, dstFunc))

	found = (mat.Norm(dstGrad, math.Inf(1)) <= problem.Eps1)
	mu = problem.Tau * maxDiagElem(dstA)

	for iter := 0; ; iter++ {
		if iter == set.Iterations {
			status = optimize.IterationLimit
			break
		}
		if found {
			status = optimize.StepConvergence
			break
		}

		// Solve (A + mu * I) * h_lm = g.
		addToDiag(dstA, mu)
		err := dstH.SolveVec(dstA, dstGrad)
		if err != nil {
			panic("singular")
		}

		// Return A to its original state for the next steps. This is done in order not to copy A.
		addToDiag(dstA, -mu)

		if mat.Norm(dstH, 2) <= (floats.Norm(parameters, 2)+problem.Eps2)*problem.Eps2 {
			found = true
		} else {
			updateParams(parametersNew, parameters, dstH)

			// Calculate rho = (F(x) - F(x_new)) / (L(0) - L(h_lm)), where
			// F = 0.5 * f.T * f, L = 0.5 * h_lm.T * (mu * h_lm - g).
			problem.Func(dstFuncNew, parametersNew)
			rho := calcRho(dstFunc, dstFuncNew, dstH, dstGrad, mu)

			if rho > 0 { // step is acceptable
				copy(parameters, parametersNew)
				problem.Func(dstFunc, parameters)
				problem.Jac(dstJac, parameters)
				dstA.Mul(dstJac.T(), dstJac)
				dstGrad.MulVec(dstJac.T(), mat.NewVecDense(size, dstFunc))
				found = (mat.Norm(dstGrad, math.Inf(1)) <= problem.Eps1) ||
					(0.5*floats.Dot(dstFunc, dstFunc) <= set.ObjectiveTol)
				mu = mu * math.Max(1.0/3.0, 1-math.Pow(2*rho-1, 3))
				nu = 2.0
			} else {
				mu *= nu
				nu *= 2.0
			}
		}
	}
	return &Result{
		X:      parameters,
		Status: status,
	}, nil
}

// LMProblem is used for running LM optimization. The objective function is
// F = 0.5 * f.T * f, where f:Rn -> Rm and m >= n.
type LMProblem struct {
	// Dim is the dimension of the parameters of the problem (n).
	Dim int
	// Size specifies the number of nonlinear functions (m).
	Size int
	// Func computes the function value at params.
	Func func(dst, param []float64)
	// Jac computes the jacobian matrix of Func.
	Jac func(dst *mat.Dense, param []float64)
	// InitParams stores the users inital guess. Defaults to the zero vector when nil.
	InitParams []float64
	// Tau scales the initial damping parameter.
	Tau float64
	// Eps1 is a stopping criterion for the gradient of F.
	Eps1 float64
	// Eps2 is a stopping criterion for the step size.
	Eps2 float64
}
