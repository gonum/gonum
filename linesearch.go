// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import "github.com/gonum/floats"

const maxLinesearchIterations = 30

// Linesearch is a linesearch-based optimization method.
// It consists of a NextDirectioner, which specifies the next linesearch method,
// and a LinesearchMethod which performs the linesearch in the direction specified
// by the NextDirectioner.
type Linesearch struct {
	NextDirectioner NextDirectioner
	Method          LinesearchMethod

	initLoc   []float64
	direction []float64

	f *FunctionStats

	lastEvalType EvaluationType
	finished     bool
	finishedF    float64
}

func (l *Linesearch) Init(loc Location, f *FunctionStats, xNext []float64) (EvaluationType, IterationType, error) {
	l.initLoc = resize(l.initLoc, len(loc.X))
	copy(l.initLoc, loc.X)

	l.direction = resize(l.direction, len(loc.X))
	stepSize := l.NextDirectioner.InitDirection(loc, l.direction)

	projGrad := floats.Dot(loc.Gradient, l.direction)
	if projGrad > 0 {
		return NoEvaluation, NoIteration, ErrNonNegativeStepDirection
	}

	evalType := l.Method.Init(loc.F, projGrad, stepSize, f)

	floats.AddScaledTo(xNext, l.initLoc, stepSize, l.direction)

	l.f = f
	l.lastEvalType = evalType
	l.finished = false
	return evalType, Major, nil
}

func (l *Linesearch) Iterate(loc Location, xNext []float64) (EvaluationType, IterationType, error) {
	if l.finished {
		// Means that we needed to evaluate the gradient, so now we have it and can initialize
		l.finished = false
		loc.F = l.finishedF
		return l.initializeNextLinesearch(loc, xNext)
	}

	projGrad := floats.Dot(loc.Gradient, l.direction)
	if l.Method.Finished(loc.F, projGrad) {
		if l.lastEvalType == JustFunction {
			l.finished = true
			l.finishedF = loc.F
			// We have the function value at the current location, but we don't
			// have the gradient, so get it before starting the next major iteration.
			copy(xNext, loc.X)
			/*
				if l.f.IsGradFunction {
					return FunctionAndGradient, Sub, nil
				}
			*/
			return JustGradient, Sub, nil
		}
		return l.initializeNextLinesearch(loc, xNext)
	}

	// Line search not done, just iterate
	stepSize, evalType := l.Method.Iterate(loc.F, projGrad)
	floats.AddScaledTo(xNext, l.initLoc, stepSize, l.direction)
	l.lastEvalType = evalType
	return evalType, Minor, nil
}

func (l *Linesearch) initializeNextLinesearch(loc Location, xNext []float64) (EvaluationType, IterationType, error) {
	// Line search is finished, so find the next direction, and
	// start the next line search
	copy(l.initLoc, loc.X)
	stepsize := l.NextDirectioner.NextDirection(loc, l.direction)
	projGrad := floats.Dot(loc.Gradient, l.direction)

	if projGrad > 0 {
		return NoEvaluation, NoIteration, ErrNonNegativeStepDirection
	}
	evalType := l.Method.Init(loc.F, projGrad, stepsize, l.f)

	floats.AddScaledTo(xNext, l.initLoc, stepsize, l.direction)

	l.lastEvalType = evalType

	return evalType, Major, nil
}
