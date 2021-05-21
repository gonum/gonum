// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ivp_test

import (
	"fmt"
	"log"
	"math"
	"testing"

	"gonum.org/v1/gonum/ivp"
	"gonum.org/v1/gonum/mat"
)

type TestModel struct {
	*ivp.Model
	solution *ivp.Model
	err      func(h, i float64) float64
}

func TestQuadratic(t *testing.T) {
	Quadratic := quadTestModel(t)
	var solver ivp.Integrator = new(ivp.RK4)

	err := solver.Set(0, Quadratic.Model)
	if err != nil {
		log.Fatal(err)
	}

	nx, _ := Quadratic.Model.Dims()
	steps := 10
	dt := 0.1

	results := make([]float64, nx)

	solmodel := Quadratic.solution
	soleq, _ := solmodel.Equations()
	solDims, _ := solmodel.Dims()
	solution := make([]float64, solDims)
	for i := 1.; i < float64(steps+1); i++ {
		dom := i * dt
		solver.Step(results, dt)
		soleq(solution, dom, results, nil)
		for j := range results {
			got := math.Abs(solution[j] - results[j])
			expected := Quadratic.err(dt, i)
			if got > expected {
				t.Errorf("error %e greater than threshold %e", got, expected)
			}

		}
	}
}

func Example_fallingBall() {
	const (
		g = -10. // gravity field [m.s^-2]
	)
	// we declare our physical model in the following function
	ballModel, err := ivp.NewModel(mat.NewVecDense(2, []float64{100., 0.}),
		nil, func(yvec []float64, _ float64, xvec, _ []float64) {
			// this anonymous function defines the physics.
			// The first variable xvec[0] corresponds to position
			// second variable xvec[1] is velocity.
			Dx := xvec[1]
			// yvec represents change in xvec, or derivative with respect to domain
			// Change in position will be equal to velocity, which is the second variable:
			// thus yvec[0] = xvec[1], which is the same as saying "change in xvec[0]" is equal to xvec[1]
			yvec[0] = Dx
			// change in velocity is acceleration. We suppose our ball is on earth accelerating at `g`
			yvec[1] = g
		}, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Here we choose our algorithm. Runge-Kutta 4th order is used
	var solver ivp.Integrator = new(ivp.RK4)
	// Before integrating the IVP is passed to the integrator (a.k.a solver). Domain (time) starts at 0
	dom := 0.
	err = solver.Set(dom, ballModel)
	if err != nil {
		log.Fatal(err)
	}

	// we define the domain over which we integrate: 10 steps with step length of 0.1
	steps := 10 // number of steps
	dt := 0.1   // step length
	// we will store position in xvec and domain (time) in dvec
	xvec := mat.NewVecDense(steps, nil)
	dvec := mat.NewVecDense(steps, nil)
	// results is an auxiliary vector that stores integration results
	nx, _ := ballModel.Dims()
	results := make([]float64, nx)

	for i := 0; i < steps; i++ {
		// Step integrates the IVP. Each step advances the solution
		step, _ := solver.Step(results, dt) // for non-adaptive algorithms step == dt
		dom += step                         // calculate domain at current step
		xvec.SetVec(i, results[0])          // set x value
		dvec.SetVec(i, dom)
	}
	// print results
	fmt.Println(mat.Formatted(xvec), "\n\n", mat.Formatted(dvec))
}

// Quadratic model may be used for future algorithms
func quadTestModel(t *testing.T) *TestModel {
	Quadratic := new(TestModel)
	quad, err := ivp.NewModel(mat.NewVecDense(2, []float64{0, 0}),
		nil, func(y []float64, dom float64, x, u []float64) {
			y[0], y[1] = x[1], 1.
		}, nil)
	if err != nil {
		t.Fatal(err)
	}
	Quadratic.Model = quad
	quadsol, err := ivp.NewModel(mat.NewVecDense(2, []float64{0, 0}),
		nil, func(y []float64, dom float64, x, u []float64) {
			y[0], y[1] = dom*dom/2., dom
		}, nil)
	if err != nil {
		t.Fatal(err)
	}
	Quadratic.solution = quadsol
	Quadratic.err = func(h, i float64) float64 { return math.Pow(h*i, 4) }
	return Quadratic
}
