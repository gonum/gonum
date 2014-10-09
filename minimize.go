// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"math"
	"time"

	"github.com/gonum/floats"
)

// Minimize finds a local minimum of a function using a sequential algorithm.
// In order to maximize a function, just multiply the output by -1.
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
// Minimize returns a Result struct and any error that occured. Please see the
// documentation of Result for more information.
//
// Please be aware that the default behavior of minimize is to find the minimum.
// For certain functions and optimization methods, this process can take many
// function evaluations. If you would like to put limits on this, for example
// maximum runtime or maximum function evaluations, please modify the Settings
// input struct.
func Minimize(f Function, initX []float64,
	settings *Settings, method Method) (*Result, error) {

	startTime := time.Now()

	funcs, funcStat := findFunctionStats(f)

	if method == nil {
		method = getDefaultMethod(funcStat)
	}

	methodStatus, methodIsStatuser := method.(Statuser)

	if settings == nil {
		settings = DefaultSettings()
	}

	location, err := setStartingLocation(f, funcs, funcStat, initX, settings)
	if err != nil {
		return nil, err
	}

	optLoc := &Location{}
	copyLocation(optLoc, location)

	stats := &Stats{}
	// update stats (grad norm, function value, etc.) so that things are
	// initialized for the first convergence check
	update(location, optLoc, stats, funcStat, NoEvaluation, NoIteration, startTime)

	if settings.Recorder != nil {
		err = settings.Recorder.Init(funcStat)
		if err != nil {
			return &Result{Status: RecorderError}, err
		}
	}

	xNext := make([]float64, len(location.X))

	evalType, iterType, err := method.Init(location, funcStat, xNext)
	if err != nil {
		return nil, err
	}

	var status Status
	for {
		if settings.Recorder != nil {
			err = settings.Recorder.Record(location, evalType, iterType, stats)
			if err != nil {
				status = RecorderError
				break
			}
		}

		status = checkConvergence(location, iterType, stats, settings)
		if status != NotTerminated {
			break
		}

		if funcStat.IsStatuser {
			status, err = funcs.status.Status()
			if err != nil || status != NotTerminated {
				break
			}
		}

		if methodIsStatuser {
			status, err = methodStatus.Status()
			if err != nil || status != NotTerminated {
				break
			}
		}

		// Compute the new function and update the statistics
		err = evaluate(funcs, funcStat, evalType, xNext, &location)
		if err != nil {
			break
		}
		update(location, optLoc, stats, funcStat, evalType, iterType, startTime)

		// Find the next location
		evalType, iterType, err = method.Iterate(location, xNext)
		if err != nil {
			break
		}
	}

	stats.Runtime = time.Since(startTime)
	return &Result{
		Stats:    *stats,
		Location: *optLoc,
		Status:   status,
	}, err
}

func copyLocation(dst *Location, src Location) {
	dst.X = resize(dst.X, len(src.X))
	copy(dst.X, src.X)

	dst.F = src.F

	dst.Gradient = resize(dst.Gradient, len(src.Gradient))
	copy(dst.Gradient, src.Gradient)
}

func findFunctionStats(f Function) (functions, *FunctionStats) {
	// Not sure how/if we want to compute timing to be used with functions
	gradient, isGradient := f.(Gradient)
	gradFunc, isFunGrad := f.(FunctionGradient)

	status, isStatuser := f.(Statuser)

	stats := &FunctionStats{
		IsGradient: isGradient,
		IsFunGrad:  isFunGrad,
		IsStatuser: isStatuser,
	}
	funcs := functions{
		function: f,
		gradient: gradient,
		gradFunc: gradFunc,
		status:   status,
	}

	return funcs, stats
}

func getDefaultMethod(f *FunctionStats) Method {
	if f.IsFunGrad {
		return &BFGS{}
	}
	// TODO: Implement a gradient-free method
	panic("gradient-free methods not yet coded")
}

// Combine location and stats because maybe in the future we'll add evaluation times
// to functionStats?
func setStartingLocation(f Function, funcs functions, stats *FunctionStats, initX []float64, settings *Settings) (Location, error) {

	var l Location

	if len(initX) == 0 {
		// maybe panic? This is clearly a mistake
		return l, ErrZeroDimensional
	}

	l.X = make([]float64, len(initX))
	copy(l.X, initX)

	if settings.UseInitData {
		initF := settings.InitialFunctionValue
		// Do we allow Inf initial function value?
		if math.IsNaN(initF) {
			return l, ErrNaN
		}
		if math.IsInf(initF, 1) {
			return l, ErrInf
		}
		l.F = initF

		initG := settings.IntialGradient
		if stats.IsGradient {
			if len(initX) != len(initG) {
				panic("minimize: initial location size mismatch")
			}

			l.Gradient = make([]float64, len(initG))
			copy(l.Gradient, initG)
		}
		return l, nil
	}
	// Compute any missing information in the inital solution

	if stats.IsFunGrad {
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

	if itertype == Major && loc.Gradient != nil {
		if stats.GradNorm <= settings.GradientAbsoluteTolerance {
			return GradientAbsoluteConvergence
		}
	}

	if itertype == Major && loc.F < settings.FunctionAbsoluteTolerance {
		return FunctionAbsoluteConvergence
	}

	if settings.MaximumFunctionEvaluations > 0 {
		totalFun := stats.NumFunEvals + stats.NumFunGradEvals
		if totalFun > settings.MaximumFunctionEvaluations {
			return FunctionEvaluationLimit
		}
	}

	if settings.MaximumGradientEvaluations > 0 {
		totalGrad := stats.NumGradEvals + stats.NumFunGradEvals
		if totalGrad > settings.MaximumGradientEvaluations {
			return GradientEvaluationLimit
		}
	}

	if settings.MaximumRuntime > 0 {
		if stats.Runtime > settings.MaximumRuntime {
			return RuntimeLimit
		}
	}

	if itertype == Major && settings.MaximumMajorIterations > 0 {
		if stats.NumMajorIterations >= settings.MaximumMajorIterations {
			return IterationLimit
		}
	}
	return NotTerminated
}

// evaluate evaluates the function and stores the answer in place
func evaluate(funcs functions, funcStat *FunctionStats, evalType EvaluationType, xNext []float64, location *Location) error {

	copy(location.X, xNext)
	switch evalType {
	case JustFunction:

		location.F = funcs.function.F(xNext)
		for i := range location.Gradient {
			location.Gradient[i] = math.NaN()
		}
		return nil
	case JustGradient:
		location.F = math.NaN()
		if funcStat.IsGradient {
			funcs.gradient.Df(location.X, location.Gradient)
			return nil
		}
		if funcStat.IsFunGrad {
			location.F = funcs.gradFunc.FDf(location.X, location.Gradient)
			return nil
		}
		return ErrMismatch{Type: evalType}
	case FunctionAndGradient:
		if funcStat.IsFunGrad {
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
func update(location Location, optLoc *Location, stats *Stats, funcStat *FunctionStats, evalType EvaluationType, iterType IterationType, startTime time.Time) {
	switch evalType {
	case JustFunction:
		stats.NumFunEvals++
	case JustGradient:
		if funcStat.IsGradient {
			stats.NumGradEvals++
		}
		if funcStat.IsFunGrad {
			stats.NumFunGradEvals++
		}
	case FunctionAndGradient:
		if funcStat.IsFunGrad {
			stats.NumFunGradEvals++
		}
		if funcStat.IsGradient {
			stats.NumFunEvals++
			stats.NumFunGradEvals++
		}
	}
	if iterType == Major {
		stats.NumMajorIterations++
	}
	if location.F < optLoc.F {
		copyLocation(optLoc, location)
	}
	stats.Runtime = time.Since(startTime)
	stats.GradNorm = floats.Norm(location.Gradient, 2) / math.Sqrt(float64(len(location.Gradient)))
}
