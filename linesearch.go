// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"

	"github.com/gonum/floats"
)

// LinesearchMethod represents an abstract optimization method in which
// a function is optimized through successive line search optimizations.
// It consists of a NextDirectioner, which specifies the search direction
// of each linesearch, and a Linesearcher which performs a linesearch along
// the search direction.
type LinesearchMethod struct {
	NextDirectioner NextDirectioner
	Linesearcher    Linesearcher

	x     []float64 // Starting point for the current iteration.
	dir   []float64 // Search direction for the current iteration.
	first bool      // Indicator of the first iteration.

	evalType EvaluationType
	iterType IterationType
}

func (ls *LinesearchMethod) Init(loc *Location, xNext []float64) (EvaluationType, IterationType, error) {
	if loc.Gradient == nil {
		panic("linesearch: gradient is nil")
	}

	dim := len(loc.X)
	ls.x = resize(ls.x, dim)
	ls.dir = resize(ls.dir, dim)
	ls.first = true

	return ls.initNextLinesearch(loc, xNext)
}

func (ls *LinesearchMethod) Iterate(loc *Location, xNext []float64) (EvaluationType, IterationType, error) {
	if ls.iterType == SubIteration {
		// We needed to evaluate invalid fields of Location. Now we have them
		// and can announce MajorIteration.
		copy(xNext, loc.X)
		ls.evalType = NoEvaluation
		ls.iterType = MajorIteration
		return ls.evalType, ls.iterType, nil
	}

	if ls.iterType == MajorIteration {
		// The linesearch previously signaled MajorIteration. Since we're here,
		// it means that the previous location is not good enough to converge,
		// so start the next linesearch.
		return ls.initNextLinesearch(loc, xNext)
	}

	projGrad := floats.Dot(loc.Gradient, ls.dir)
	if ls.Linesearcher.Finished(loc.F, projGrad) {
		copy(xNext, loc.X)
		// Check if the last evaluation evaluated all fields of Location.
		ls.evalType = complementEval(loc, ls.evalType)
		if ls.evalType == NoEvaluation {
			// Location is complete and MajorIteration can be announced directly.
			ls.iterType = MajorIteration
		} else {
			// Location is not complete, evaluate its invalid fields in SubIteration.
			ls.iterType = SubIteration
		}
		return ls.evalType, ls.iterType, nil
	}

	// Line search not done, just iterate.
	stepSize, evalType, err := ls.Linesearcher.Iterate(loc.F, projGrad)
	if err != nil {
		ls.evalType = NoEvaluation
		ls.iterType = NoIteration
		return ls.evalType, ls.iterType, err
	}

	floats.AddScaledTo(xNext, ls.x, stepSize, ls.dir)
	// Compare the starting point for the current iteration with the next
	// evaluation point to make sure that rounding errors do not prevent progress.
	if floats.Equal(ls.x, xNext) {
		ls.evalType = NoEvaluation
		ls.iterType = NoIteration
		return ls.evalType, ls.iterType, ErrNoProgress
	}

	ls.evalType = evalType
	ls.iterType = MinorIteration
	return ls.evalType, ls.iterType, nil
}

func (ls *LinesearchMethod) initNextLinesearch(loc *Location, xNext []float64) (EvaluationType, IterationType, error) {
	copy(ls.x, loc.X)

	var stepSize float64
	if ls.first {
		stepSize = ls.NextDirectioner.InitDirection(loc, ls.dir)
		ls.first = false
	} else {
		stepSize = ls.NextDirectioner.NextDirection(loc, ls.dir)
	}

	projGrad := floats.Dot(loc.Gradient, ls.dir)
	if projGrad >= 0 {
		ls.evalType = NoEvaluation
		ls.iterType = NoIteration
		return ls.evalType, ls.iterType, ErrNonNegativeStepDirection
	}

	ls.evalType = ls.Linesearcher.Init(loc.F, projGrad, stepSize)

	floats.AddScaledTo(xNext, ls.x, stepSize, ls.dir)
	// Compare the starting point for the current iteration with the next
	// evaluation point to make sure that rounding errors do not prevent progress.
	if floats.Equal(ls.x, xNext) {
		ls.evalType = NoEvaluation
		ls.iterType = NoIteration
		return ls.evalType, ls.iterType, ErrNoProgress
	}

	ls.iterType = MinorIteration
	return ls.evalType, ls.iterType, nil
}

// ArmijoConditionMet returns true if the Armijo condition (aka sufficient
// decrease) has been met. Under normal conditions, the following should be
// true, though this is not enforced:
//  - initGrad < 0
//  - step > 0
//  - 0 < funcConst < 1
func ArmijoConditionMet(currObj, initObj, initGrad, step, funcConst float64) bool {
	return currObj <= initObj+funcConst*step*initGrad
}

// StrongWolfeConditionsMet returns true if the strong Wolfe conditions have been met.
// The strong Wolfe conditions ensure sufficient decrease in the function
// value, and sufficient decrease in the magnitude of the projected gradient.
// Under normal conditions, the following should be true, though this is not
// enforced:
//  - initGrad < 0
//  - step > 0
//  - 0 <= funcConst < gradConst < 1
func StrongWolfeConditionsMet(currObj, currGrad, initObj, initGrad, step, funcConst, gradConst float64) bool {
	if currObj > initObj+funcConst*step*initGrad {
		return false
	}
	return math.Abs(currGrad) < gradConst*math.Abs(initGrad)
}

// WeakWolfeConditionsMet returns true if the weak Wolfe conditions have been met.
// The weak Wolfe conditions ensure sufficient decrease in the function value,
// and sufficient decrease in the value of the projected gradient. Under normal
// conditions, the following should be true, though this is not enforced:
//  - initGrad < 0
//  - step > 0
//  - 0 <= funcConst < gradConst < 1
func WeakWolfeConditionsMet(currObj, currGrad, initObj, initGrad, step, funcConst, gradConst float64) bool {
	if currObj > initObj+funcConst*step*initGrad {
		return false
	}
	return currGrad >= gradConst*initGrad
}
