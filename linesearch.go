// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import "github.com/gonum/floats"

// Linesearch is a linesearch-based optimization method.
// It consists of a NextDirectioner, which specifies the next linesearch method,
// and a LinesearchMethod which performs the linesearch in the direction specified
// by the NextDirectioner.
type Linesearch struct {
	NextDirectioner NextDirectioner
	Method          LinesearchMethod

	initLoc   []float64
	direction []float64

	f *FunctionInfo

	lastEvalType EvaluationType
	finished     bool
	finishedF    float64
	iter         int
}

func (l *Linesearch) Init(loc Location, f *FunctionInfo, xNext []float64) (EvaluationType, IterationType, error) {
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
	l.iter = 0
	return evalType, MajorIteration, nil
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
		l.iter = 0
		if l.lastEvalType == FunctionEval {
			l.finished = true
			l.finishedF = loc.F
			// We have the function value at the current location, but we don't
			// have the gradient, so get it before starting the next major iteration.
			copy(xNext, loc.X)
			return GradientEval, SubIteration, nil
		}
		return l.initializeNextLinesearch(loc, xNext)
	}

	l.iter++
	// Line search not done, just iterate
	stepSize, evalType, err := l.Method.Iterate(loc.F, projGrad)
	if err != nil {
		return NoEvaluation, NoIteration, err
	}
	floats.AddScaledTo(xNext, l.initLoc, stepSize, l.direction)
	l.lastEvalType = evalType
	return evalType, MinorIteration, nil
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

	return evalType, MajorIteration, nil
}
