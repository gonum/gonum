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

	x   []float64 // Starting point for the current iteration.
	dir []float64 // Search direction for the current iteration.

	first     bool // Indicator of the first iteration.
	nextMajor bool // Indicates that MajorIteration must be requested at the next call to Iterate().

	loc  Location    // Storage for intermediate locations.
	eval RequestType // Indicator of valid fields in loc.

	lastStep    float64     // Step taken from x in the previous call to Iterate().
	lastRequest RequestType // Request returned from the previous call to Iterate().
}

func (ls *LinesearchMethod) Init(loc *Location) (RequestType, error) {
	if loc.Gradient == nil {
		panic("linesearch: gradient is nil")
	}

	dim := len(loc.X)
	ls.x = resize(ls.x, dim)
	ls.dir = resize(ls.dir, dim)

	ls.first = true
	ls.nextMajor = false

	copyLocation(&ls.loc, loc)
	// Indicate that all fields of ls.loc are valid.
	ls.eval = FuncEvaluation | GradEvaluation
	if ls.loc.Hessian != nil {
		ls.eval |= HessEvaluation
	}

	ls.lastStep = math.NaN()
	ls.lastRequest = NoRequest

	return ls.initNextLinesearch(loc.X)
}

func (ls *LinesearchMethod) Iterate(loc *Location) (RequestType, error) {
	switch ls.lastRequest {
	case NoRequest:
		// TODO(vladimir-ch): We have previously returned with an error and
		// Init() was not called. What to do? What about ls's internal state?

	case MajorIteration:
		// We previously requested MajorIteration but since we're here, the
		// previous location was not good enough to converge. So start the next
		// linesearch and store the next evaluation point in loc.X.
		return ls.initNextLinesearch(loc.X)

	default:
		if ls.lastRequest&EvaluationRequest == 0 {
			panic("linesearch: unexpected request")
		}

		// Store the result of the previously requested evaluation into ls.loc.
		if ls.lastRequest&FuncEvaluation != 0 {
			ls.loc.F = loc.F
		}
		if ls.lastRequest&GradEvaluation != 0 {
			copy(ls.loc.Gradient, loc.Gradient)
		}
		if ls.lastRequest&HessEvaluation != 0 {
			ls.loc.Hessian.CopySym(loc.Hessian)
		}
		// Update the indicator of valid fields of ls.loc.
		ls.eval |= ls.lastRequest

		if ls.nextMajor {
			ls.nextMajor = false

			// Linesearcher previously indicated that it had finished, but we
			// needed to evaluate invalid fields of ls.loc. Now we have them and
			// can announce MajorIteration.

			copyLocation(loc, &ls.loc)
			ls.lastRequest = MajorIteration
			return ls.lastRequest, nil
		}
	}

	projGrad := floats.Dot(ls.loc.Gradient, ls.dir)
	if ls.Linesearcher.Finished(ls.loc.F, projGrad) {
		// Form a request that evaluates invalid fields of ls.loc.
		ls.lastRequest = complementEval(&ls.loc, ls.eval)
		if ls.lastRequest == NoRequest {
			// ls.loc is complete and MajorIteration can be announced directly.
			copyLocation(loc, &ls.loc)
			ls.lastRequest = MajorIteration
		} else {
			ls.nextMajor = true
		}
		return ls.lastRequest, nil
	}

	step, request, err := ls.Linesearcher.Iterate(ls.loc.F, projGrad)
	if err != nil {
		return ls.error(err)
	}

	if step != ls.lastStep {
		// We are moving to a new location.

		// Compute the next evaluation point and store it in loc.X.
		floats.AddScaledTo(loc.X, ls.x, step, ls.dir)
		if floats.Equal(ls.x, loc.X) {
			// Step size has become so small that the next evaluation point is
			// indistinguishable from the starting point for the current
			// iteration due to rounding errors.
			return ls.error(ErrNoProgress)
		}

		ls.lastStep = step
		copy(ls.loc.X, loc.X) // Move ls.loc to the next evaluation point
		ls.eval = NoRequest   // and invalidate all its fields.
	} else {
		// Linesearcher is requesting another evaluation at the same point
		// which is stored in ls.loc.X.
		copy(loc.X, ls.loc.X)
	}

	ls.lastRequest = request
	return ls.lastRequest, nil
}

func (ls *LinesearchMethod) error(err error) (RequestType, error) {
	ls.lastRequest = NoRequest
	return ls.lastRequest, err
}

// initNextLinesearch initializes the next linesearch using the previous
// complete location stored in ls.loc. It fills xNext and returns an
// evaluation request to be performed at xNext.
func (ls *LinesearchMethod) initNextLinesearch(xNext []float64) (RequestType, error) {
	copy(ls.x, ls.loc.X)

	var step float64
	if ls.first {
		ls.first = false
		step = ls.NextDirectioner.InitDirection(&ls.loc, ls.dir)
	} else {
		step = ls.NextDirectioner.NextDirection(&ls.loc, ls.dir)
	}

	projGrad := floats.Dot(ls.loc.Gradient, ls.dir)
	if projGrad >= 0 {
		return ls.error(ErrNonNegativeStepDirection)
	}

	ls.lastRequest = ls.Linesearcher.Init(ls.loc.F, projGrad, step)

	floats.AddScaledTo(xNext, ls.x, step, ls.dir)
	if floats.Equal(ls.x, xNext) {
		// Step size is so small that the next evaluation point is
		// indistinguishable from the starting point for the current iteration
		// due to rounding errors.
		return ls.error(ErrNoProgress)
	}

	ls.lastStep = step
	copy(ls.loc.X, xNext) // Move ls.loc to the next evaluation point
	ls.eval = NoRequest   // and invalidate all its fields.

	return ls.lastRequest, nil
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
