// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"

	"github.com/gonum/floats"
)

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

// QuadraticInterpolateStepSize estimates the initial line search step size as
// the minimum of a quadratic that interpolates f(x_{k-1}), f(x_k) and grad f_k
// \dot p_k. This is useful for line search methods that do not produce
// well-scaled descent directions, such as gradient descent or conjugate
// gradient methods. The step size will lie in the interval (0, 1], no minimum
// step size is imposed as the line search along the descent direction should
// take care of that.
//
// See also Nocedal, Wright (2006), Numerical Optimization (2nd ed.), sec.
// 3.5, page 59.
type QuadraticInterpolateStepSize struct {
	// If the relative change in the objective function is larger than
	// InterpolationCutOff, the step size is estimated by quadratic
	// interpolation, otherwise it is set to one.
	InterpolationCutOff float64
	// If x0 is not at the origin, the step size at the first iteration is
	// estimated as InitialStepFactor * |x0|_∞ / |g|_∞.
	InitialStepFactor float64

	fPrev float64
}

func (q *QuadraticInterpolateStepSize) Init(l Location, dir []float64) (stepsize float64) {
	if q.InterpolationCutOff == 0 {
		q.InterpolationCutOff = 1e-10
	}
	if q.InitialStepFactor == 0 {
		q.InitialStepFactor = 0.01
	}

	// In case we lack any other information, take a unit step
	stepsize = 1
	xNorm := floats.Norm(l.X, math.Inf(1))
	if xNorm != 0 {
		gNorm := floats.Norm(l.Gradient, math.Inf(1))
		// If the initial location x0 is not at the origin, use the |x0|_∞ and
		// |g|_∞ to estimate the initial step size. Divide by |g|_∞ to take
		// shorter step if the magnitude of the gradient is large, multiply by
		// |x0|_∞ to avoid rounding errors in the computation of x0 + stepsize*dir0.
		stepsize = q.InitialStepFactor * xNorm / gNorm
	} else if l.F != 0 {
		gNorm := floats.Norm(l.Gradient, 2)
		// If x0 is at the origin and F(x0) != 0, use |F(x0)| and |g|_2 to
		// estimate the initial step size.
		stepsize = 2 * math.Abs(l.F) / math.Pow(gNorm, 2)
	}

	q.fPrev = l.F
	return stepsize
}

func (q *QuadraticInterpolateStepSize) StepSize(l Location, dir []float64) (stepsize float64) {
	stepsize = 1
	t := 1.0
	if l.F != 0 {
		t = (q.fPrev - l.F) / math.Abs(l.F)
	}
	if t > q.InterpolationCutOff {
		// The relative change between two consecutive function values compared to
		// the function value itself is large enough, so compute the minimum of
		// a quadratic interpolant.
		// Assuming that the received direction is a descent direction,
		// stepsize will be positive.
		stepsize = 2 * (l.F - q.fPrev) / floats.Dot(l.Gradient, dir)
		// Bound the step size from above by 1. We do not impose any lower
		// bound, line search should take care of that.
		stepsize = math.Min(1.01*stepsize, 1)
	}

	q.fPrev = l.F
	return stepsize
}
