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

	initX []float64
	dir   []float64

	funcInfo *FunctionInfo

	lastEvalType EvaluationType
	iterType     IterationType
}

func (ls *Linesearch) Init(loc *Location, f *FunctionInfo, xNext []float64) (EvaluationType, IterationType, error) {
	ls.initX = resize(ls.initX, len(loc.X))
	copy(ls.initX, loc.X)

	ls.dir = resize(ls.dir, len(loc.X))
	stepSize := ls.NextDirectioner.InitDirection(loc, ls.dir)

	projGrad := math.NaN()
	if loc.Gradient != nil {
		projGrad = floats.Dot(loc.Gradient, ls.dir)
		if projGrad >= 0 {
			return NoEvaluation, NoIteration, ErrNonNegativeStepDirection
		}
	}
	lsLoc := LinesearchLocation{
		F:          loc.F,
		Derivative: projGrad,
	}
	evalType := ls.Method.Init(lsLoc, stepSize, f)
	floats.AddScaledTo(xNext, ls.initX, stepSize, ls.dir)
	ls.funcInfo = f
	ls.lastEvalType = evalType
	ls.iterType = MinorIteration
	return evalType, ls.iterType, nil
}

func (ls *Linesearch) Iterate(loc *Location, xNext []float64) (EvaluationType, IterationType, error) {
	if ls.iterType == SubIteration {
		// We needed to evaluate the gradient, so now we have it and can
		// announce MajorIteration.
		ls.iterType = MajorIteration
		copy(xNext, loc.X)
		return NoEvaluation, ls.iterType, nil
	}
	if ls.iterType == MajorIteration {
		// The linesearch previously signaled MajorIteration. Since we're here,
		// it means that the previous location is not good enough to converge,
		// so start the next linesearch.
		return ls.initNextLinesearch(loc, xNext)
	}
	projGrad := math.NaN()
	if loc.Gradient != nil {
		projGrad = floats.Dot(loc.Gradient, ls.dir)
	}
	lsLoc := LinesearchLocation{
		F:          loc.F,
		Derivative: projGrad,
	}
	if ls.Method.Finished(lsLoc) {
		if ls.lastEvalType == FunctionEval && loc.Gradient != nil {
			// We have the function value at the current location, but we don't
			// have the gradient, so get it before announcing MajorIteration.
			ls.iterType = SubIteration
			copy(xNext, loc.X)
			return GradientEval, ls.iterType, nil
		}
		// The linesearch is finished. Announce so with an update to
		// MajorIteration. The function value and gradient is already known, so
		// no function evaluations are necessary.
		ls.iterType = MajorIteration
		copy(xNext, loc.X)
		return NoEvaluation, ls.iterType, nil
	}

	// Line search not done, just iterate.
	stepSize, evalType, err := ls.Method.Iterate(lsLoc)
	if err != nil {
		return NoEvaluation, NoIteration, err
	}
	floats.AddScaledTo(xNext, ls.initX, stepSize, ls.dir)
	// Compare the starting point for the current iteration with the next
	// evaluation point to make sure that rounding errors do not prevent progress.
	if floats.Equal(ls.initX, xNext) {
		return NoEvaluation, NoIteration, ErrNoProgress
	}
	ls.lastEvalType = evalType
	ls.iterType = MinorIteration
	return evalType, ls.iterType, nil
}

func (ls *Linesearch) initNextLinesearch(loc *Location, xNext []float64) (EvaluationType, IterationType, error) {
	// Find the next direction, and start the next line search.
	copy(ls.initX, loc.X)
	stepsize := ls.NextDirectioner.NextDirection(loc, ls.dir)
	projGrad := math.NaN()
	if loc.Gradient != nil {
		projGrad = floats.Dot(loc.Gradient, ls.dir)
	}
	if projGrad >= 0 {
		return NoEvaluation, NoIteration, ErrNonNegativeStepDirection
	}
	lsLoc := LinesearchLocation{
		F:          loc.F,
		Derivative: projGrad,
	}
	evalType := ls.Method.Init(lsLoc, stepsize, ls.funcInfo)
	floats.AddScaledTo(xNext, ls.initX, stepsize, ls.dir)
	// Compare the starting point for the current iteration with the next
	// evaluation point to make sure that rounding errors do not prevent progress.
	if floats.Equal(ls.initX, xNext) {
		return NoEvaluation, NoIteration, ErrNoProgress
	}
	ls.lastEvalType = evalType
	ls.iterType = MinorIteration
	return evalType, ls.iterType, nil
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
