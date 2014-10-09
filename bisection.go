// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import "math"

// Bisection is a LinesearchMethod that uses a bisection to find a point that
// satisfies the strong Wolfe conditions with the given gradient constant and
// function constant of zero. If GradConst is zero, it will be set to a reasonable
// value. Bisection will panic if GradConst is not between zero and one.
type Bisection struct {
	GradConst float64

	minStep  float64
	maxStep  float64
	currStep float64

	initF float64
	minF  float64
	maxF  float64

	initGrad float64
	minGrad  float64
	maxGrad  float64
}

func (b *Bisection) Init(initF, initG, initStepSize float64, f *FunctionStats) EvaluationType {
	if initG > 0 {
		panic("init G greater than 0")
	}

	if b.GradConst == 0 {
		b.GradConst = 0.9
	}
	if b.GradConst <= 0 || b.GradConst > 1 {
		panic("bisection: GradConst not between 0 and 1")
	}

	b.minStep = 0
	b.maxStep = math.Inf(1)
	b.currStep = initStepSize

	b.initF = initF
	b.minF = initF
	b.maxF = math.NaN()

	b.initGrad = initG
	b.minGrad = initG
	b.maxGrad = math.NaN()

	return FunctionAndGradient
}

func (b *Bisection) Finished(f, g float64) bool {
	isFinished := StrongWolfeConditionsMet(f, g, b.initF, b.initGrad, b.currStep, 0, b.GradConst)
	return isFinished
}

func (b *Bisection) Iterate(f, g float64) (float64, EvaluationType) {
	// Deciding on the next step size
	if math.IsInf(b.maxStep, 1) {
		// Have not yet bounded the minimum
		switch {
		case g > 0:
			// Found a change in derivative sign, so this is the new maximum
			b.maxStep = b.currStep
			b.maxF = f
			b.maxGrad = g
			b.currStep = (b.minStep + b.maxStep) / 2
			return b.currStep, FunctionAndGradient
		case f <= b.minF:
			// Still haven't found an upper bound, but there is not an increase in
			// function value and the gradient is still negative, so go more in
			// that direction.
			b.minStep = b.currStep
			b.minF = f
			b.minGrad = g
			b.currStep *= 2

			return b.currStep, FunctionAndGradient
		default:
			// Increase in function value, but the gradient is still negative.
			// Means we must have skipped over a local minimum, so set this point
			// as the new maximum
			b.maxStep = b.currStep
			b.maxF = f
			b.maxGrad = g
			b.currStep = (b.minStep + b.maxStep) / 2

			return b.currStep, FunctionAndGradient
		}
	}
	// We have already bounded the minimum, so we're just working to find one
	// close enough to the minimum to meet the strong wolfe conditions
	if g < 0 {
		if f <= b.minF {
			b.minStep = b.currStep
			b.minF = f
			b.minGrad = g
		} else {
			// Negative gradient, but increase in function value, so must have
			// skipped over a local minimum. Set this as the new maximum location
			b.maxStep = b.currStep
			b.maxF = f
			b.maxGrad = g
		}
	} else {
		// Gradient is positive, so minimum must be between the max point and
		// the minimum point
		b.maxStep = b.currStep
		b.maxF = f
		b.maxGrad = g
	}
	b.currStep = (b.minStep + b.maxStep) / 2

	return b.currStep, FunctionAndGradient
}
