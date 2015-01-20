// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"fmt"
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
	funcs, funcInfo := getFunctionInfo(f)

	if funcInfo.IsStatuser {
		_, err := funcs.status.Status()
		if err != nil {
			return nil, err
		}
	}

	if settings.Recorder != nil {
		// Initialize Recorder first. If it fails, we avoid the (possibly
		// time-consuming) evaluation of F and DF at the starting location.
		err := settings.Recorder.Init(funcInfo)
		if err != nil {
			return nil, err
		}
	}

	if settings == nil {
		settings = DefaultSettings()
	}
	stats := &Stats{}
	optLoc, evalType, err := getStartingLocation(funcs, funcInfo, initX, stats, settings)
	if err != nil {
		return nil, err
	}

	// Runtime is the only Stats field that needs to be updated here.
	stats.Runtime = time.Since(startTime)
	// Send optLoc to Recorder before checking it for convergence.
	if settings.Recorder != nil {
		// TODO(vladimir-ch): Replace NoIteration with InitialIteration when it
		// is added.
		err = settings.Recorder.Record(optLoc, evalType, NoIteration, stats)
	}

	// Check if the starting location satisfies the convergence criteria.
	// TODO(vladimir-ch): Replace MajorIteration with InitialIteration when it
	// is added and checkConvergence() updated to handle it.
	status := checkConvergence(optLoc, MajorIteration, stats, settings)
	if status == NotTerminated && err == nil {
		if method == nil {
			method = getDefaultMethod(funcInfo)
		}
		// The starting location is not good enough, we need to perform a
		// minimization. The optimal location will be stored in-place in
		// optLoc.
		status, err = minimize(settings, method, funcInfo, stats, funcs, optLoc, startTime)
	}

	if settings.Recorder != nil && err == nil {
		// Send the optimal location to Recorder.
		err = settings.Recorder.Record(optLoc, NoEvaluation, PostIteration, stats)
	}
	stats.Runtime = time.Since(startTime)
	return &Result{
		Location: *optLoc,
		Stats:    *stats,
		Status:   status,
	}, err
}

func minimize(settings *Settings, method Method, funcInfo *FunctionInfo, stats *Stats, funcs functions, optLoc *Location, startTime time.Time) (status Status, err error) {
	location := &Location{}
	copyLocation(location, optLoc)
	xNext := make([]float64, len(location.X))

	methodStatus, methodIsStatuser := method.(Statuser)

	evalType, iterType, err := method.Init(location, funcInfo, xNext)
	if err != nil {
		return Failure, err
	}

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

		if funcInfo.IsStatuser {
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
		evaluate(funcs, funcInfo, evalType, xNext, location, stats)
		update(location, optLoc, stats, iterType, startTime)

		// Find the next location
		evalType, iterType, err = method.Iterate(location, xNext)
		if err != nil {
			status = Failure
			return
		}
	}
	panic("unreachable")
}

func copyLocation(dst, src *Location) {
	dst.X = resize(dst.X, len(src.X))
	copy(dst.X, src.X)

	dst.F = src.F

	dst.Gradient = resize(dst.Gradient, len(src.Gradient))
	copy(dst.Gradient, src.Gradient)
}

func getFunctionInfo(f Function) (functions, *FunctionInfo) {
	// Not sure how/if we want to compute timing to be used with functions
	gradient, isGradient := f.(Gradient)
	gradFunc, isFunGrad := f.(FunctionGradient)

	status, isStatuser := f.(Statuser)

	funcInfo := &FunctionInfo{
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

	return funcs, funcInfo
}

func getDefaultMethod(funcInfo *FunctionInfo) Method {
	if funcInfo.IsGradient || funcInfo.IsFunctionGradient {
		return &BFGS{}
	}
	// TODO: Implement a gradient-free method
	panic("optimize: gradient-free methods not yet coded")
}

// getStartingLocation allocates and initializes the starting location for the minimization.
func getStartingLocation(funcs functions, funcInfo *FunctionInfo, initX []float64, stats *Stats, settings *Settings) (*Location, EvaluationType, error) {
	dim := len(initX)
	l := &Location{
		X: make([]float64, dim),
	}
	if funcInfo.IsGradient || funcInfo.IsFunctionGradient {
		l.Gradient = make([]float64, dim)
	}
	copy(l.X, initX)

	evalType := NoEvaluation
	if settings.UseInitialData {
		l.F = settings.InitialFunctionValue
		if l.Gradient != nil {
			initG := settings.InitialGradient
			if len(initG) != dim {
				panic("local: initial location size mismatch")
			}
			copy(l.Gradient, initG)
		}
	} else {
		evalType = FunctionEval
		if l.Gradient != nil {
			evalType = FunctionAndGradientEval
		}
		evaluate(funcs, funcInfo, evalType, l.X, l, stats)
	}

	if math.IsNaN(l.F) {
		return l, evalType, ErrNaN
	}
	if math.IsInf(l.F, 1) {
		return l, evalType, ErrInf
	}
	return l, evalType, nil
}

func checkConvergence(loc *Location, itertype IterationType, stats *Stats, settings *Settings) Status {
	if itertype == MajorIteration && loc.Gradient != nil {
		norm := floats.Norm(loc.Gradient, math.Inf(1))
		if norm < settings.GradientAbsTol {
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
		// TODO(vladimir-ch): It would be nice to update Runtime here.
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

// evaluate evaluates the function given by funcs at xNext, stores the answer
// into location and updates stats. If location.X is not equal to xNext, then
// unused fields of location are set to NaN.
// evaluate panics if the function does not support the requested evalType.
func evaluate(funcs functions, funcInfo *FunctionInfo, evalType EvaluationType, xNext []float64, location *Location, stats *Stats) {
	different := !floats.Equal(location.X, xNext)
	if different {
		copy(location.X, xNext)
	}
	switch evalType {
	case FunctionEval:
		if different {
			// Invalidate the gradient because it will not be evaluated.
			for i := range location.Gradient {
				location.Gradient[i] = math.NaN()
			}
		}
		location.F = funcs.function.F(location.X)
		stats.FunctionEvals++
		return
	case GradientEval:
		if funcInfo.IsGradient {
			if different {
				// Invalidate the function value because it will not be
				// evaluated.
				location.F = math.NaN()
			}
			funcs.gradient.Df(location.X, location.Gradient)
			stats.GradientEvals++
			return
		}
		if funcInfo.IsFunctionGradient {
			location.F = funcs.gradFunc.FDf(location.X, location.Gradient)
			stats.FunctionGradientEvals++
			return
		}
	case FunctionAndGradientEval:
		if funcInfo.IsFunctionGradient {
			location.F = funcs.gradFunc.FDf(location.X, location.Gradient)
			stats.FunctionGradientEvals++
			return
		}
		if funcInfo.IsGradient {
			location.F = funcs.function.F(location.X)
			stats.FunctionEvals++
			funcs.gradient.Df(location.X, location.Gradient)
			stats.GradientEvals++
			return
		}
	case NoEvaluation:
		if different {
			// The evaluation point xNext is not equal to location.X, so we
			// fill location with NaNs to indicate that it does not contain a
			// valid evaluation result.
			location.F = math.NaN()
			for i := range location.X {
				location.X[i] = math.NaN()
			}
			for i := range location.Gradient {
				location.Gradient[i] = math.NaN()
			}
		}
	default:
		panic(fmt.Sprintf("unknown evaluation type %v", evalType))
	}
	panic(fmt.Sprintf("objective function does not support %v", evalType))
}

// update updates the stats given the new evaluation
func update(location *Location, optLoc *Location, stats *Stats, iterType IterationType, startTime time.Time) {
	if iterType == MajorIteration {
		stats.MajorIterations++
	}
	if location.F <= optLoc.F {
		copyLocation(optLoc, location)
	}
	stats.Runtime = time.Since(startTime)
}
