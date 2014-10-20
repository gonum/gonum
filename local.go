// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"
	"time"

	"github.com/gonum/floats"
)

// Local finds a local minimum of a function using a sequential algorithm.
// In order to maximize a function, multiply the output by -1.
//
// The first argument is of Function type representing the function to be minimized.
// Type switching is used to see if the function implements Gradient, FunctionGradient
// and Statuser.
//
// The second argument is the initial location at which to start the minimization.
// The initial location must be supplied, and must have a length equal to the
// problem dimension.
//
// The third argument contains the settings for the minimization. It is here that
// gradient tolerance, etc. are specified. The DefaultSettings() function
// can be called for a Settings struct with the default values initialized.
// If settings == nil, the default settings are used. Please see the documentation
// for the Settings structure for more information. The optimization Method used
// may also contain settings, see documentation for the appropriate optimizer.
//
// The final argument is the optimization method to use. If method == nil, then
// an appropriate default is chosen based on the properties of the other arguments
// (dimension, gradient-free or gradient-based, etc.). The optimization
// methods in this package are designed such that reasonable defaults occur
// if options are not specified explicitly. For example, the code
//  method := &Bfgs{}
// creates a pointer to a new Bfgs struct. When minimize is called, the settings
// in the method will be populated with default values. The methods are also
// designed such that they can be reused in future calls to method.
//
// Local returns a Result struct and any error that occurred. Please see the
// documentation of Result for more information.
//
// Please be aware that the default behavior of Local is to find the minimum.
// For certain functions and optimization methods, this process can take many
// function evaluations. If you would like to put limits on this, for example
// maximum runtime or maximum function evaluations, please modify the Settings
// input struct.
func Local(f Function, initX []float64, settings *Settings, method Method) (*Result, error) {
	if len(initX) == 0 {
		panic("local: initial X has zero length")
	}

	startTime := time.Now()

	funcs, funcStat := findFunctionStats(f)

	if method == nil {
		method = getDefaultMethod(funcStat)
	}

	if settings == nil {
		settings = DefaultSettings()
	}

	location, err := setStartingLocation(f, funcs, funcStat, initX, settings)
	if err != nil {
		return nil, err
	}

	stats := &Stats{}
	optLoc := &Location{}
	// update stats (grad norm, function value, etc.) so that things are
	// initialized for the first convergence check
	update(location, optLoc, stats, funcStat, NoEvaluation, NoIteration, startTime)

	if settings.Recorder != nil {
		err = settings.Recorder.Init(funcStat)
		if err != nil {
			return &Result{Status: Failure}, err
		}
	}

	// actually perform the minimization
	status, err := minimize(settings, location, method, funcStat, stats, funcs, optLoc, startTime)

	// cleanup at exit
	if settings.Recorder != nil && err == nil {
		err = settings.Recorder.Record(*optLoc, NoEvaluation, PostIteration, stats)
	}
	stats.Runtime = time.Since(startTime)
	return &Result{
		Stats:    *stats,
		Location: *optLoc,
		Status:   status,
	}, err
}

func minimize(settings *Settings, location Location, method Method, funcStat *FunctionInfo, stats *Stats, funcs functions, optLoc *Location, startTime time.Time) (status Status, err error) {
	methodStatus, methodIsStatuser := method.(Statuser)
	xNext := make([]float64, len(location.X))

	evalType, iterType, err := method.Init(location, funcStat, xNext)
	if err != nil {
		return Failure, err
	}
	copyLocation(optLoc, location)

	for {
		if settings.Recorder != nil {
			err = settings.Recorder.Record(location, evalType, iterType, stats)
			if err != nil {
				status = Failure
				return
			}
		}

		status = checkConvergence(location, iterType, stats, settings)
		if status != NotTerminated {
			return
		}

		if funcStat.IsStatuser {
			status, err = funcs.status.Status()
			if err != nil || status != NotTerminated {
				return
			}
		}

		if methodIsStatuser {
			status, err = methodStatus.Status()
			if err != nil || status != NotTerminated {
				return
			}
		}

		// Compute the new function and update the statistics
		err = evaluate(funcs, funcStat, evalType, xNext, &location)
		if err != nil {
			status = Failure
			return
		}
		update(location, optLoc, stats, funcStat, evalType, iterType, startTime)

		// Find the next location
		evalType, iterType, err = method.Iterate(location, xNext)
		if err != nil {
			status = Failure
			return
		}
	}
	panic("unreachable")
}

func copyLocation(dst *Location, src Location) {
	dst.X = resize(dst.X, len(src.X))
	copy(dst.X, src.X)

	dst.F = src.F

	dst.Gradient = resize(dst.Gradient, len(src.Gradient))
	copy(dst.Gradient, src.Gradient)
}

func findFunctionStats(f Function) (functions, *FunctionInfo) {
	// Not sure how/if we want to compute timing to be used with functions
	gradient, isGradient := f.(Gradient)
	gradFunc, isFunGrad := f.(FunctionGradient)

	status, isStatuser := f.(Statuser)

	stats := &FunctionInfo{
		IsGradient:         isGradient,
		IsFunctionGradient: isFunGrad,
		IsStatuser:         isStatuser,
	}
	funcs := functions{
		function: f,
		gradient: gradient,
		gradFunc: gradFunc,
		status:   status,
	}

	return funcs, stats
}

func getDefaultMethod(f *FunctionInfo) Method {
	if f.IsFunctionGradient {
		return &BFGS{}
	}
	// TODO: Implement a gradient-free method
	panic("optimize: gradient-free methods not yet coded")
}

// Combine location and stats because maybe in the future we'll add evaluation times
// to functionStats?
func setStartingLocation(f Function, funcs functions, stats *FunctionInfo, initX []float64, settings *Settings) (Location, error) {
	var l Location

	l.X = make([]float64, len(initX))
	copy(l.X, initX)

	if settings.UseInitialData {
		initF := settings.InitialFunctionValue
		// Do we allow Inf initial function value?
		if math.IsNaN(initF) {
			return l, ErrNaN
		}
		if math.IsInf(initF, 1) {
			return l, ErrInf
		}
		l.F = initF

		initG := settings.InitialGradient
		if stats.IsGradient {
			if len(initX) != len(initG) {
				panic("local: initial location size mismatch")
			}

			l.Gradient = make([]float64, len(initG))
			copy(l.Gradient, initG)
		}
		return l, nil
	}

	// Compute missing information in the initial state.
	if stats.IsFunctionGradient {
		l.Gradient = make([]float64, len(initX))
		l.F = funcs.gradFunc.FDf(initX, l.Gradient)
		return l, nil
	}
	l.F = funcs.function.F(l.X)
	if math.IsNaN(l.F) {
		return l, ErrNaN
	}
	if math.IsInf(l.F, 1) {
		return l, ErrInf
	}
	return l, nil
}

func checkConvergence(loc Location, itertype IterationType, stats *Stats, settings *Settings) Status {
	if itertype == MajorIteration && loc.Gradient != nil {
		if stats.GradientNorm <= settings.GradientAbsTol {
			return GradientAbsoluteConvergence
		}
	}

	if itertype == MajorIteration && loc.F < settings.FunctionAbsTol {
		return FunctionAbsoluteConvergence
	}

	// Check every step for negative infinity because it could break the
	// linesearches and -inf is the best you can do anyway.
	if math.IsInf(loc.F, -1) {
		return FunctionNegativeInfinity
	}

	if settings.FunctionEvals > 0 {
		totalFun := stats.FunctionEvals + stats.FunctionGradientEvals
		if totalFun >= settings.FunctionEvals {
			return FunctionEvaluationLimit
		}
	}

	if settings.GradientEvals > 0 {
		totalGrad := stats.GradientEvals + stats.FunctionGradientEvals
		if totalGrad >= settings.GradientEvals {
			return GradientEvaluationLimit
		}
	}

	if settings.Runtime > 0 {
		if stats.Runtime >= settings.Runtime {
			return RuntimeLimit
		}
	}

	if itertype == MajorIteration && settings.MajorIterations > 0 {
		if stats.MajorIterations >= settings.MajorIterations {
			return IterationLimit
		}
	}
	return NotTerminated
}

// evaluate evaluates the function and stores the answer in place
func evaluate(funcs functions, funcStat *FunctionInfo, evalType EvaluationType, xNext []float64, location *Location) error {
	copy(location.X, xNext)
	switch evalType {
	case FunctionEval:
		location.F = funcs.function.F(xNext)
		for i := range location.Gradient {
			location.Gradient[i] = math.NaN()
		}
		return nil
	case GradientEval:
		location.F = math.NaN()
		if funcStat.IsGradient {
			funcs.gradient.Df(location.X, location.Gradient)
			return nil
		}
		if funcStat.IsFunctionGradient {
			location.F = funcs.gradFunc.FDf(location.X, location.Gradient)
			return nil
		}
		return ErrMismatch{Type: evalType}
	case FunctionAndGradientEval:
		if funcStat.IsFunctionGradient {
			location.F = funcs.gradFunc.FDf(xNext, location.Gradient)
			return nil
		}
		if funcStat.IsGradient {
			location.F = funcs.function.F(xNext)
			funcs.gradient.Df(xNext, location.Gradient)
			return nil
		}
		return ErrMismatch{Type: evalType}
	default:
		panic("unreachable")
	}
}

// update updates the stats given the new evaluation
func update(location Location, optLoc *Location, stats *Stats, funcStat *FunctionInfo, evalType EvaluationType, iterType IterationType, startTime time.Time) {
	switch evalType {
	case FunctionEval:
		stats.FunctionEvals++
	case GradientEval:
		if funcStat.IsGradient {
			stats.GradientEvals++
		} else if funcStat.IsFunctionGradient {
			stats.FunctionGradientEvals++
		}
	case FunctionAndGradientEval:
		if funcStat.IsFunctionGradient {
			stats.FunctionGradientEvals++
		} else if funcStat.IsGradient {
			stats.FunctionEvals++
			stats.GradientEvals++
		}
	}
	if iterType == MajorIteration {
		stats.MajorIterations++
	}
	if location.F < optLoc.F {
		copyLocation(optLoc, location)
	}
	stats.Runtime = time.Since(startTime)
	if location.Gradient != nil {
		stats.GradientNorm = floats.Norm(location.Gradient, 2) / math.Sqrt(float64(len(location.Gradient)))
	}
}
