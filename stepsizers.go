// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"

	"github.com/gonum/floats"
)

const (
	initialStepFactor = 1

	quadraticMinimumStepSize = 1e-3
	quadraticMaximumStepSize = 1.01
	quadraticThreshold       = 1e-12
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

// QuadraticStepSize estimates the initial line search step size as the minimum
// of a quadratic that interpolates f(x_{k-1}), f(x_k) and ∇f_k⋅p_k.
// This is useful for line search methods that do not produce well-scaled
// descent directions, such as gradient descent or conjugate gradient methods.
// The step size is bounded away from zero.
type QuadraticStepSize struct {
	// If the relative change in the objective function is larger than
	// Threshold, the step size is estimated by quadratic interpolation,
	// otherwise it is set to 2*previous step size.
	// The default value is 1e-12.
	Threshold float64
	// The step size at the first iteration is estimated as InitialStepFactor / |g|_∞.
	// If InitialStepFactor is zero, it will be set to one.
	InitialStepFactor float64
	// The estimated step size is always greater than or equal to MinStepSize.
	// MinStepSize times GradientAbsTol should always be greater than machine epsilon.
	// If MinStepSize is zero, it will be set to 1e-3.
	MinStepSize float64
	// The estimated step size is always lower than or equal to MaxStepSize.
	// If MaxStepSize is zero, it will be set to 1.01.
	MaxStepSize float64

	fPrev        float64
	dirPrevNorm  float64
	projGradPrev float64
	xPrev        []float64
}

func (q *QuadraticStepSize) Init(l Location, dir []float64) (stepSize float64) {
	q.xPrev = resize(q.xPrev, len(l.X))
	if q.Threshold == 0 {
		q.Threshold = quadraticThreshold
	}
	if q.InitialStepFactor == 0 {
		q.InitialStepFactor = initialStepFactor
	}
	if q.MinStepSize == 0 {
		q.MinStepSize = quadraticMinimumStepSize
	}
	if q.MaxStepSize == 0 {
		q.MaxStepSize = quadraticMaximumStepSize
	}

	gNorm := floats.Norm(l.Gradient, math.Inf(1))
	stepSize = math.Max(q.MinStepSize, math.Min(q.InitialStepFactor/gNorm, q.MaxStepSize))

	q.fPrev = l.F
	q.dirPrevNorm = floats.Norm(dir, 2)
	q.projGradPrev = floats.Dot(l.Gradient, dir)
	copy(q.xPrev, l.X)
	return stepSize
}

func (q *QuadraticStepSize) StepSize(l Location, dir []float64) (stepSize float64) {
	stepSizePrev := floats.Distance(l.X, q.xPrev, 2) / q.dirPrevNorm
	projGrad := floats.Dot(l.Gradient, dir)

	stepSize = 2 * stepSizePrev
	if !floats.EqualWithinRel(q.fPrev, l.F, q.Threshold) {
		// Two consecutive function values are not relatively equal, so
		// computing the minimum of a quadratic interpolant might make sense

		df := (l.F - q.fPrev) / stepSizePrev
		quadTest := df - q.projGradPrev
		if quadTest > 0 {
			// There is a chance of approximating the function well by a
			// quadratic only if the finite difference (f_k-f_{k-1})/stepSizePrev
			// is larger than ∇f_{k-1}⋅p_{k-1}

			// Set the step size to the minimizer of the quadratic function that
			// interpolates f_{k-1}, ∇f_{k-1}⋅p_{k-1} and f_k
			stepSize = -q.projGradPrev * stepSizePrev / quadTest / 2
		}
	}
	// Bound the step size to lie in [MinStepSize, MaxStepSize]
	stepSize = math.Max(q.MinStepSize, math.Min(stepSize, q.MaxStepSize))

	q.fPrev = l.F
	q.dirPrevNorm = floats.Norm(dir, 2)
	q.projGradPrev = projGrad
	copy(q.xPrev, l.X)
	return stepSize
}
