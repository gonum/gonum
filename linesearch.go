// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

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
	iterType     IterationType
}

func (l *Linesearch) Init(loc *Location, f *FunctionInfo, xNext []float64) (EvaluationType, IterationType, error) {
	l.initLoc = resize(l.initLoc, len(loc.X))
	copy(l.initLoc, loc.X)

	l.direction = resize(l.direction, len(loc.X))
	stepSize := l.NextDirectioner.InitDirection(loc, l.direction)

	projGrad := math.NaN()
	if loc.Gradient != nil {
		projGrad = floats.Dot(loc.Gradient, l.direction)
		if projGrad >= 0 {
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
	l.iterType = MinorIteration
	return evalType, l.iterType, nil
}

func (l *Linesearch) Iterate(loc *Location, xNext []float64) (EvaluationType, IterationType, error) {
	if l.iterType == SubIteration {
		// We needed to evaluate the gradient, so now we have it and can
		// announce MajorIteration.
		l.iterType = MajorIteration
		copy(xNext, loc.X)
		return NoEvaluation, l.iterType, nil
	}
	if l.iterType == MajorIteration {
		// The linesearch previously signaled MajorIteration. Since we're here,
		// it means that the previous location is not good enough to converge,
		// so start the next linesearch.
		return l.initNextLinesearch(loc, xNext)
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
		if l.lastEvalType == FunctionEval && loc.Gradient != nil {
			// We have the function value at the current location, but we don't
			// have the gradient, so get it before announcing MajorIteration.
			l.iterType = SubIteration
			copy(xNext, loc.X)
			return GradientEval, l.iterType, nil
		}
		// The linesearch is finished. Announce so with an update to
		// MajorIteration. The function value and gradient is already known, so
		// no function evaluations are necessary.
		l.iterType = MajorIteration
		copy(xNext, loc.X)
		return NoEvaluation, l.iterType, nil
	}

	// Line search not done, just iterate.
	stepSize, evalType, err := l.Method.Iterate(linesearchLocation)
	if err != nil {
		return NoEvaluation, NoIteration, err
	}
	floats.AddScaledTo(xNext, l.initLoc, stepSize, l.direction)
	l.lastEvalType = evalType
	l.iterType = MinorIteration
	return evalType, l.iterType, nil
}

func (l *Linesearch) initNextLinesearch(loc *Location, xNext []float64) (EvaluationType, IterationType, error) {
	// Find the next direction, and start the next line search.
	copy(l.initLoc, loc.X)
	stepsize := l.NextDirectioner.NextDirection(loc, l.direction)
	projGrad := math.NaN()
	if loc.Gradient != nil {
		projGrad = floats.Dot(loc.Gradient, l.direction)
	}
	if projGrad >= 0 {
		return NoEvaluation, NoIteration, ErrNonNegativeStepDirection
	}
	initLinesearchLocation := LinesearchLocation{
		F:          loc.F,
		Derivative: projGrad,
	}
	evalType := l.Method.Init(initLinesearchLocation, stepsize, l.funInfo)
	floats.AddScaledTo(xNext, l.initLoc, stepsize, l.direction)
	l.lastEvalType = evalType
	l.iterType = MinorIteration
	return evalType, l.iterType, nil
}

// ArmijoConditionMet returns true if the Armijo condition (aka sufficient decrease)
// has been met. Under normal conditions, the following should be true, though this is not enforced:
//  - initGrad < 0
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
