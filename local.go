// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"fmt"
	"math"
	"time"

	"github.com/gonum/floats"
	"github.com/gonum/matrix/mat64"
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
		panic("optimize: initial X has zero length")
	}

	startTime := time.Now()
	funcInfo := newFunctionInfo(f)
	if method == nil {
		method = getDefaultMethod(funcInfo)
	}
	if err := funcInfo.satisfies(method); err != nil {
		return nil, err
	}

	if funcInfo.IsStatuser {
		_, err := funcInfo.statuser.Status()
		if err != nil {
			return nil, err
		}
	}

	if settings == nil {
		settings = DefaultSettings()
	}

	if settings.Recorder != nil {
		// Initialize Recorder first. If it fails, we avoid the (possibly
		// time-consuming) evaluation of F and DF at the starting location.
		err := settings.Recorder.Init(&funcInfo.FunctionInfo)
		if err != nil {
			return nil, err
		}
	}

	stats := &Stats{}
	optLoc, evalType, err := getStartingLocation(funcInfo, method, initX, stats, settings)
	if err != nil {
		return nil, err
	}

	if settings.FunctionConverge != nil {
		settings.FunctionConverge.Init(optLoc.F)
	}

	// Runtime is the only Stats field that needs to be updated here.
	stats.Runtime = time.Since(startTime)
	// Send optLoc to Recorder before checking it for convergence.
	if settings.Recorder != nil {
		err = settings.Recorder.Record(optLoc, evalType, InitIteration, stats)
	}

	// Check if the starting location satisfies the convergence criteria.
	status := checkConvergence(optLoc, InitIteration, stats, settings)
	if status == NotTerminated && err == nil {
		// The starting location is not good enough, we need to perform a
		// minimization. The optimal location will be stored in-place in
		// optLoc.
		status, err = minimize(settings, method, funcInfo, stats, optLoc, startTime)
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

func minimize(settings *Settings, method Method, funcInfo *functionInfo, stats *Stats, optLoc *Location, startTime time.Time) (status Status, err error) {
	loc := &Location{}
	copyLocation(loc, optLoc)
	xNext := make([]float64, len(loc.X))

	methodStatus, methodIsStatuser := method.(Statuser)

	evalType, iterType, err := method.Init(loc, &funcInfo.FunctionInfo, xNext)
	if err != nil {
		return Failure, err
	}

	for {
		if funcInfo.IsStatuser {
			// Check the function status before evaluating.
			status, err = funcInfo.statuser.Status()
			if err != nil || status != NotTerminated {
				return
			}
		}

		// Perform evalType evaluation of the function at xNext and store the
		// result in location.
		evaluate(funcInfo, evalType, xNext, loc, stats)
		// Update the stats and optLoc.
		update(loc, optLoc, stats, iterType, startTime)
		// Get the convergence status before recording the new location.
		status = checkConvergence(optLoc, iterType, stats, settings)

		if settings.Recorder != nil {
			err = settings.Recorder.Record(loc, evalType, iterType, stats)
			if err != nil {
				if status == NotTerminated {
					status = Failure
				}
				return
			}
		}

		if status != NotTerminated {
			return
		}

		if methodIsStatuser {
			status, err = methodStatus.Status()
			if err != nil || status != NotTerminated {
				return
			}
		}

		// Find the next location (stored in-place into xNext).
		evalType, iterType, err = method.Iterate(loc, xNext)
		if err != nil {
			status = Failure
			return
		}
	}
	panic("optimize: unreachable")
}

func copyLocation(dst, src *Location) {
	dst.X = resize(dst.X, len(src.X))
	copy(dst.X, src.X)

	dst.F = src.F

	dst.Gradient = resize(dst.Gradient, len(src.Gradient))
	copy(dst.Gradient, src.Gradient)

	if src.Hessian != nil {
		if dst.Hessian == nil || dst.Hessian.Symmetric() != len(src.X) {
			dst.Hessian = mat64.NewSymDense(len(src.X), nil)
		}
		dst.Hessian.CopySym(src.Hessian)
	}
}

func getDefaultMethod(funcInfo *functionInfo) Method {
	if funcInfo.IsGradient {
		return &BFGS{}
	}
	return &NelderMead{}
}

// getStartingLocation allocates and initializes the starting location for the minimization.
func getStartingLocation(funcInfo *functionInfo, method Method, initX []float64, stats *Stats, settings *Settings) (*Location, EvaluationType, error) {
	dim := len(initX)
	loc := &Location{
		X: make([]float64, dim),
	}
	copy(loc.X, initX)
	if method.Needs().Gradient {
		loc.Gradient = make([]float64, dim)
	}
	if method.Needs().Hessian {
		loc.Hessian = mat64.NewSymDense(dim, nil)
	}

	evalType := NoEvaluation
	if settings.UseInitialData {
		loc.F = settings.InitialValue
		if loc.Gradient != nil {
			initG := settings.InitialGradient
			if initG == nil {
				panic("optimize: initial gradient is nil")
			}
			if len(initG) != dim {
				panic("optimize: initial gradient size mismatch")
			}
			copy(loc.Gradient, initG)
		}
		if loc.Hessian != nil {
			initH := settings.InitialHessian
			if initH == nil {
				panic("optimize: initial Hessian is nil")
			}
			if initH.Symmetric() != dim {
				panic("optimize: initial Hessian size mismatch")
			}
			loc.Hessian.CopySym(initH)
		}
	} else {
		switch {
		case loc.Hessian != nil:
			evalType = FuncGradHessEvaluation
		case loc.Gradient != nil:
			evalType = FuncGradEvaluation
		default:
			evalType = FuncEvaluation
		}
		evaluate(funcInfo, evalType, loc.X, loc, stats)
	}

	if math.IsNaN(loc.F) {
		return loc, evalType, ErrNaN
	}
	if math.IsInf(loc.F, 1) {
		return loc, evalType, ErrInf
	}
	for _, v := range loc.Gradient {
		if math.IsInf(v, 0) {
			return loc, evalType, ErrGradInf
		}
		if math.IsNaN(v) {
			return loc, evalType, ErrGradNaN
		}
	}

	return loc, evalType, nil
}

func checkConvergence(loc *Location, iterType IterationType, stats *Stats, settings *Settings) Status {
	if iterType == MajorIteration || iterType == InitIteration {
		if loc.Gradient != nil {
			norm := floats.Norm(loc.Gradient, math.Inf(1))
			if norm < settings.GradientThreshold {
				return GradientThreshold
			}
		}
		if loc.F < settings.FunctionThreshold {
			return FunctionThreshold
		}
	}

	if iterType == MajorIteration && settings.FunctionConverge != nil {
		status := settings.FunctionConverge.FunctionConverged(loc.F)
		if status != NotTerminated {
			return status
		}
	}

	// Check every step for negative infinity because it could break the
	// linesearches and -inf is the best you can do anyway.
	if math.IsInf(loc.F, -1) {
		return FunctionNegativeInfinity
	}

	if settings.FuncEvaluations > 0 {
		totalFun := stats.FuncEvaluations + stats.FuncGradEvaluations + stats.FuncGradHessEvaluations
		if totalFun >= settings.FuncEvaluations {
			return FunctionEvaluationLimit
		}
	}

	if settings.GradEvaluations > 0 {
		totalGrad := stats.GradEvaluations + stats.FuncGradEvaluations + stats.FuncGradHessEvaluations
		if totalGrad >= settings.GradEvaluations {
			return GradientEvaluationLimit
		}
	}

	if settings.HessEvaluations > 0 {
		totalHess := stats.HessEvaluations + stats.FuncGradHessEvaluations
		if totalHess >= settings.HessEvaluations {
			return HessianEvaluationLimit
		}
	}

	if settings.Runtime > 0 {
		// TODO(vladimir-ch): It would be nice to update Runtime here.
		if stats.Runtime >= settings.Runtime {
			return RuntimeLimit
		}
	}

	if iterType == MajorIteration && settings.MajorIterations > 0 {
		if stats.MajorIterations >= settings.MajorIterations {
			return IterationLimit
		}
	}
	return NotTerminated
}

// invalidate marks all fields of Location with NaNs and it is the
// responsibility of Method to issue such EvaluationTypes that make Location
// valid again. Invalidation is done internally as a help for implementers to
// detect silent bugs in Methods using inconsistent Location, e.g., using
// Gradient without a GradEvaluation request.
func invalidate(loc *Location) {
	loc.F = math.NaN()
	if loc.Gradient != nil {
		loc.Gradient[0] = math.NaN()
	}
	if loc.Hessian != nil {
		loc.Hessian.SetSym(0, 0, math.NaN())
	}
}

// evaluate evaluates the function given by f at xNext, stores the answer into
// loc and updates stats. If loc.X is not equal to xNext, then unused fields of
// loc are set to NaN.
// evaluate panics if the function does not support the requested evalType.
func evaluate(f *functionInfo, evalType EvaluationType, xNext []float64, loc *Location, stats *Stats) {
	different := !floats.Equal(loc.X, xNext)
	if different {
		if evalType == NoEvaluation {
			// Optimizers should not request NoEvaluation at a new location.
			// The intent and therefore an appropriate action are both unclear.
			panic("optimize: no evaluation requested at new location")
		}
		invalidate(loc)
		copy(loc.X, xNext)
	}

	if evalType == NoEvaluation {
		return
	}

	if evalType == FuncGradHessEvaluation && f.IsFunctionGradientHessian {
		loc.F = f.functionGradientHessian.FuncGradHess(loc.X, loc.Gradient, loc.Hessian)
		stats.FuncGradHessEvaluations++
		return
	}
	if evalType == FuncGradEvaluation && f.IsFunctionGradient {
		loc.F = f.functionGradient.FuncGrad(loc.X, loc.Gradient)
		stats.FuncGradEvaluations++
		return
	}

	if evalType == FuncGradHessEvaluation && f.IsFunctionGradient {
		loc.F = f.functionGradient.FuncGrad(loc.X, loc.Gradient)
		stats.FuncGradEvaluations++
		evalType = HessEvaluation
	}
	if evalType&FuncEvaluation != 0 {
		loc.F = f.function.Func(loc.X)
		stats.FuncEvaluations++
		evalType &= ^FuncEvaluation
	}
	if evalType&GradEvaluation != 0 {
		f.gradient.Grad(loc.X, loc.Gradient)
		stats.GradEvaluations++
		evalType &= ^GradEvaluation
	}
	if evalType&HessEvaluation != 0 {
		f.hessian.Hess(loc.X, loc.Hessian)
		stats.HessEvaluations++
		evalType &= ^HessEvaluation
	}

	if evalType != NoEvaluation {
		panic(fmt.Sprintf("optimize: unknown evaluation type %v", evalType))
	}
}

// update updates the stats given the new evaluation
func update(loc *Location, optLoc *Location, stats *Stats, iterType IterationType, startTime time.Time) {
	if iterType == MajorIteration {
		stats.MajorIterations++
	}
	if loc.F <= optLoc.F {
		copyLocation(optLoc, loc)
	}
	stats.Runtime = time.Since(startTime)
}
