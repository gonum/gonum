// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"math"

	"github.com/gonum/floats"
)

// Linesearch is a linesearch-based optimization method.
// It consists of a NextDirectioner, which specifies the next linesearch method,
// and a LinesearchMethod which performs the linesearch in the direction specified
// by the NextDirectioner.
type Linesearch struct {
	NextDirectioner NextDirectioner
	Method          LinesearchMethod

	initLoc   []float64
	direction []float64

	funInfo *FunctionInfo

	lastEvalType EvaluationType
	finished     bool
	finishedF    float64
}

func (l *Linesearch) Init(loc Location, f *FunctionInfo, xNext []float64) (EvaluationType, IterationType, error) {
	l.initLoc = resize(l.initLoc, len(loc.X))
	copy(l.initLoc, loc.X)

	l.direction = resize(l.direction, len(loc.X))
	stepSize := l.NextDirectioner.InitDirection(loc, l.direction)

	projGrad := math.NaN()
	if loc.Gradient != nil {
		projGrad = floats.Dot(loc.Gradient, l.direction)
		if projGrad > 0 {
			return NoEvaluation, NoIteration, ErrNonNegativeStepDirection
		}
	}
	linesearchLocation := LinesearchLocation{
		F:          loc.F,
		Derivative: projGrad,
	}
	evalType := l.Method.Init(linesearchLocation, stepSize, f)
	floats.AddScaledTo(xNext, l.initLoc, stepSize, l.direction)
	l.funInfo = f
	l.lastEvalType = evalType
	l.finished = false
	return evalType, MajorIteration, nil
}

func (l *Linesearch) Iterate(loc Location, xNext []float64) (EvaluationType, IterationType, error) {
	if l.finished {
		// Means that we needed to evaluate the gradient, so now we have it and can initialize
		l.finished = false
		loc.F = l.finishedF
		return l.initializeNextLinesearch(loc, xNext)
	}
	projGrad := math.NaN()
	if loc.Gradient != nil {
		projGrad = floats.Dot(loc.Gradient, l.direction)
	}
	linesearchLocation := LinesearchLocation{
		F:          loc.F,
		Derivative: projGrad,
	}
	if l.Method.Finished(linesearchLocation) {
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

	// Line search not done, just iterate
	stepSize, evalType, err := l.Method.Iterate(linesearchLocation)
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
	projGrad := math.NaN()
	if loc.Gradient != nil {
		projGrad = floats.Dot(loc.Gradient, l.direction)
	}
	if projGrad > 0 {
		return NoEvaluation, NoIteration, ErrNonNegativeStepDirection
	}
	initLinesearchLocation := LinesearchLocation{
		F:          loc.F,
		Derivative: projGrad,
	}
	evalType := l.Method.Init(initLinesearchLocation, stepsize, l.funInfo)
	floats.AddScaledTo(xNext, l.initLoc, stepsize, l.direction)
	l.lastEvalType = evalType
	return evalType, MajorIteration, nil
}

// ArmijoConditionMet returns true if the Armijo condition (aka sufficient decrease)
// has been met. Under normal conditions, the following should be true, though this is not enforced:
// 	- initGrad < 0
//  - step > 0
//  - 0 < funConst < 1
func ArmijoConditionMet(currObj, initObj, initGrad, step, funConst float64) bool {
	return currObj <= initObj+funConst*step*initGrad
}

// StrongWolfeConditionsMet returns true if the strong Wolfe conditions have been met.
// The strong wolfe conditions ensure sufficient decrease in the function value,
// and sufficient decrease in the magnitude of the projected gradient. Under normal
// conditions, the following should be true, though this is not enforced:
//  - initGrad < 0
//  - step > 0
//  - 0 <= funConst < gradConst < 1
func StrongWolfeConditionsMet(currObj, currGrad, initObj, initGrad, step, funConst, gradConst float64) bool {
	if currObj > initObj+funConst*step*initGrad {
		return false
	}
	return math.Abs(currGrad) < gradConst*math.Abs(initGrad)
}

// WeakWolfeConditionsMet returns true if the weak Wolfe conditions have been met.
// The weak wolfe conditions ensure sufficient decrease in the function value,
// and sufficient decrease in the value of the projected gradient. Under normal
// conditions, the following should be true, though this is not enforced:
//  - initGrad < 0
//  - step > 0
//  - 0 <= funConst < gradConst < 1
func WeakWolfeConditionsMet(currObj, currGrad, initObj, initGrad, step, funConst, gradConst float64) bool {
	if currObj > initObj+funConst*step*initGrad {
		return false
	}
	return currGrad >= gradConst*initGrad
}
