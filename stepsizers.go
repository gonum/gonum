// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"

	"github.com/gonum/floats"
)

const minStepSize = 1e-8

// ConstantStepSize is a StepSizer that returns the same step size for
// every iteration.
type ConstantStepSize struct {
	Size float64
}

func (c ConstantStepSize) Init(l Location, dir []float64) float64 {
	return c.Size
}

func (c ConstantStepSize) StepSize(l Location, dir []float64) float64 {
	return c.Size
}

// QuadraticStepSize estimates the initial line search step size as the minimum
// of a quadratic that interpolates f(x_{k-1}), f(x_k) and grad f_k \dot p_k.
// This is useful for line search methods that do not produce well-scaled
// descent directions, such as gradient descent or conjugate gradient methods.
// The step size will lie in the closed interval [MinStepSize, 1].
//
// See also Nocedal, Wright (2006), Numerical Optimization (2nd ed.), sec.
// 3.5, page 59.
type QuadraticStepSize struct {
	// If the relative change in the objective function is larger than
	// InterpolationCutOff, the step size is estimated by quadratic
	// interpolation, otherwise it is set to one.
	InterpolationCutOff float64
	// The step size at the first iteration is estimated as
	// max(MinStepSize, InitialStepFactor / |g|_∞).
	InitialStepFactor float64
	// Minimum step size that is returned.
	// Default value: 1e-8
	MinStepSize float64

	fPrev float64
}

func (q *QuadraticStepSize) Init(l Location, dir []float64) (stepsize float64) {
	if q.InterpolationCutOff == 0 {
		q.InterpolationCutOff = 1e-10
	}
	if q.InitialStepFactor == 0 {
		q.InitialStepFactor = 1
	}
	if q.MinStepSize == 0 {
		q.MinStepSize = minStepSize
	}

	gNorm := floats.Norm(l.Gradient, math.Inf(1))
	stepsize = math.Max(q.MinStepSize, q.InitialStepFactor/gNorm)

	q.fPrev = l.F
	return stepsize
}

func (q *QuadraticStepSize) StepSize(l Location, dir []float64) (stepsize float64) {
	stepsize = 1
	t := 1.0
	if !floats.EqualWithinAbsOrRel(l.F, 0, 1e-8, 1e-6) {
		t = (q.fPrev - l.F) / math.Abs(l.F)
	}
	if t > q.InterpolationCutOff {
		// The relative change between two consecutive function values compared to
		// the function value itself is large enough, so compute the minimum of
		// a quadratic interpolant.
		// Assuming that the received direction is a descent direction,
		// stepsize will be positive.
		stepsize = 2 * (l.F - q.fPrev) / floats.Dot(l.Gradient, dir)
		// Trim the step size to lie in [MinStepSize, 1]
		stepsize = math.Max(q.MinStepSize, math.Min(1.01*stepsize, 1))
	}

	q.fPrev = l.F
	return stepsize
}
