// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import "math"

// Bisection is a Linesearcher that uses a bisection to find a point that
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

	lastOp Operation
}

func (b *Bisection) Init(f, g float64, step float64) Operation {
	if step <= 0 {
		panic("bisection: bad step size")
	}
	if g >= 0 {
		panic("bisection: initial derivative is non-negative")
	}

	if b.GradConst == 0 {
		b.GradConst = 0.9
	}
	if b.GradConst <= 0 || b.GradConst >= 1 {
		panic("bisection: GradConst not between 0 and 1")
	}

	b.minStep = 0
	b.maxStep = math.Inf(1)
	b.currStep = step

	b.initF = f
	b.minF = f
	b.maxF = math.NaN()

	b.initGrad = g
	b.minGrad = g
	b.maxGrad = math.NaN()

	b.lastOp = FuncEvaluation | GradEvaluation
	return b.lastOp
}

func (b *Bisection) Iterate(f, g float64) (Operation, float64, error) {
	if b.lastOp != FuncEvaluation|GradEvaluation {
		panic("bisection: Init has not been called")
	}

	// Don't finish the linesearch until a minimum is found that is better than
	// the best point found so far. We want to end up in the lowest basin of
	// attraction
	minF := b.initF
	if b.maxF < minF {
		minF = b.maxF
	}
	if b.minF < minF {
		minF = b.minF
	}
	if StrongWolfeConditionsMet(f, g, minF, b.initGrad, b.currStep, 0, b.GradConst) {
		b.lastOp = MajorIteration
		return b.lastOp, b.currStep, nil
	}

	// Deciding on the next step size
	if math.IsInf(b.maxStep, 1) {
		// Have not yet bounded the minimum
		switch {
		case g > 0:
			// Found a change in derivative sign, so this is the new maximum
			b.maxStep = b.currStep
			b.maxF = f
			b.maxGrad = g
			return b.nextStep((b.minStep + b.maxStep) / 2)
		case f <= b.minF:
			// Still haven't found an upper bound, but there is not an increase in
			// function value and the gradient is still negative, so go more in
			// that direction.
			b.minStep = b.currStep
			b.minF = f
			b.minGrad = g
			return b.nextStep(b.currStep * 2)
		default:
			// Increase in function value, but the gradient is still negative.
			// Means we must have skipped over a local minimum, so set this point
			// as the new maximum
			b.maxStep = b.currStep
			b.maxF = f
			b.maxGrad = g
			return b.nextStep((b.minStep + b.maxStep) / 2)
		}
	}

	// Already bounded the minimum, but wolfe conditions not met. Need to step to
	// find minimum.
	if f <= b.minF && f <= b.maxF {
		if g < 0 {
			b.minStep = b.currStep
			b.minF = f
			b.minGrad = g
		} else {
			b.maxStep = b.currStep
			b.maxF = f
			b.maxGrad = g
		}
	} else {
		// We found a higher point. Want to push toward the minimal bound
		if b.minF <= b.maxF {
			b.maxStep = b.currStep
			b.maxF = f
			b.maxGrad = g
		} else {
			b.minStep = b.currStep
			b.minF = f
			b.minGrad = g
		}
	}
	return b.nextStep((b.minStep + b.maxStep) / 2)
}

// nextStep checks if the new step is equal to the old step.
// This can happen if min and max are the same, or if the step size is infinity,
// both of which indicate the minimization must stop. If the steps are different,
// it sets the new step size and returns the evaluation type and the step. If the steps
// are the same, it returns an error.
func (b *Bisection) nextStep(step float64) (Operation, float64, error) {
	if b.currStep == step {
		b.lastOp = NoOperation
		return b.lastOp, b.currStep, ErrLinesearchFailure
	}
	b.currStep = step
	b.lastOp = FuncEvaluation | GradEvaluation
	return b.lastOp, b.currStep, nil
}
