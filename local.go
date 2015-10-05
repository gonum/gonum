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

// Local finds a local minimum of a minimization problem using a sequential
// algorithm. A maximization problem can be transformed into a minimization
// problem by multiplying the function by -1.
//
// The first argument represents the problem to be minimized. Its fields are
// routines that evaluate the objective function, gradient, and other
// quantities related to the problem. The objective function, p.Func, must not
// be nil. The optimization method used may require other fields to be non-nil
// as specified by method.Needs. Local will panic if these are not met. The
// method can be determined automatically from the supplied problem which is
// described below.
//
// If p.Status is not nil, it is called before every evaluation. If the
// returned Status is not NotTerminated or the error is not nil, the
// optimization run is terminated.
//
// The second argument is the initial location at which to start the minimization.
// The initial location must be supplied, and must have a length equal to the
// problem dimension.
//
// The third argument contains the settings for the minimization. It is here that
// gradient tolerance, etc. are specified. The DefaultSettings function
// can be called for a Settings struct with the default values initialized.
// If settings == nil, the default settings are used. See the documentation
// for the Settings structure for more information. The optimization Method used
// may also contain settings, see documentation for the appropriate optimizer.
//
// The final argument is the optimization method to use. If method == nil, then
// an appropriate default is chosen based on the properties of the other arguments
// (dimension, gradient-free or gradient-based, etc.). The optimization
// methods in this package are designed such that reasonable defaults occur
// if options are not specified explicitly. For example, the code
//  method := &optimize.BFGS{}
// creates a pointer to a new BFGS struct. When Local is called, the settings
// in the method will be populated with default values. The methods are also
// designed such that they can be reused in future calls to Local.
//
// If method implements Statuser, method.Status is called before every call
// to method.Iterate. If the returned Status is not NotTerminated or the
// error is non-nil, the optimization run is terminated.
//
// Local returns a Result struct and any error that occurred. See the
// documentation of Result for more information.
//
// Be aware that the default behavior of Local is to find the minimum.
// For certain functions and optimization methods, this process can take many
// function evaluations. If you would like to put limits on this, for example
// maximum runtime or maximum function evaluations, modify the Settings
// input struct.
func Local(p Problem, initX []float64, settings *Settings, method Method) (*Result, error) {
	if p.Func == nil {
		panic("optimize: objective function is undefined")
	}
	if len(initX) == 0 {
		panic("optimize: initial X has zero length")
	}

	startTime := time.Now()

	if method == nil {
		method = getDefaultMethod(&p)
	}
	if err := p.satisfies(method); err != nil {
		return nil, err
	}

	if p.Status != nil {
		_, err := p.Status()
		if err != nil {
			return nil, err
		}
	}

	if settings == nil {
		settings = DefaultSettings()
	}

	if settings.Recorder != nil {
		// Initialize Recorder first. If it fails, we avoid the (possibly
		// time-consuming) evaluation at the starting location.
		err := settings.Recorder.Init()
		if err != nil {
			return nil, err
		}
	}

	stats := &Stats{}
	optLoc, err := getStartingLocation(&p, method, initX, stats, settings)
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
		err = settings.Recorder.Record(optLoc, InitIteration, stats)
	}

	// Check if the starting location satisfies the convergence criteria.
	status := checkConvergence(optLoc, settings)
	if status == NotTerminated && err == nil {
		// The starting location is not good enough, we need to perform a
		// minimization. The optimal location will be stored in-place in
		// optLoc.
		status, err = minimize(&p, method, settings, stats, optLoc, startTime)
	}

	if settings.Recorder != nil && err == nil {
		// Send the optimal location to Recorder.
		err = settings.Recorder.Record(optLoc, PostIteration, stats)
	}
	stats.Runtime = time.Since(startTime)
	return &Result{
		Location: *optLoc,
		Stats:    *stats,
		Status:   status,
	}, err
}

func minimize(p *Problem, method Method, settings *Settings, stats *Stats, optLoc *Location, startTime time.Time) (status Status, err error) {
	loc := &Location{}
	copyLocation(loc, optLoc)
	x := make([]float64, len(loc.X))

	methodStatus, methodIsStatuser := method.(Statuser)

	var op Operation
	op, err = method.Init(loc)
	if err != nil {
		status = Failure
		return
	}

	for {
		// Sequentially call method.Iterate, performing the operations it has
		// commanded, until convergence.

		switch op {
		case NoOperation:

		case InitIteration:
			panic("optimize: Method returned InitIteration")

		case PostIteration:
			panic("optimize: Method returned PostIteration")

		case MajorIteration:
			stats.MajorIterations++
			copyLocation(optLoc, loc)
			status = checkConvergence(optLoc, settings)

		default: // Any of the Evaluation operations.
			if !op.isEvaluation() {
				panic(fmt.Sprintf("optimize: invalid evaluation %v", op))
			}

			if p.Status != nil {
				status, err = p.Status()
				if err != nil || status != NotTerminated {
					return
				}
			}
			evaluate(p, loc, op, stats, x)
		}

		if settings.Recorder != nil {
			stats.Runtime = time.Since(startTime)
			err = settings.Recorder.Record(loc, op, stats)
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

		stats.Runtime = time.Since(startTime)
		status = checkLimits(loc, stats, settings)
		if status != NotTerminated {
			return
		}

		if methodIsStatuser {
			status, err = methodStatus.Status()
			if err != nil || status != NotTerminated {
				return
			}
		}

		op, err = method.Iterate(loc)
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

func getDefaultMethod(p *Problem) Method {
	if p.Grad != nil {
		return &BFGS{}
	}
	return &NelderMead{}
}

// getStartingLocation allocates and initializes the starting location for the minimization.
func getStartingLocation(p *Problem, method Method, initX []float64, stats *Stats, settings *Settings) (*Location, error) {
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
		eval := FuncEvaluation
		if loc.Gradient != nil {
			eval |= GradEvaluation
		}
		if loc.Hessian != nil {
			eval |= HessEvaluation
		}
		x := make([]float64, len(loc.X))
		evaluate(p, loc, eval, stats, x)
	}

	if math.IsInf(loc.F, 1) || math.IsNaN(loc.F) {
		return loc, ErrFunc(loc.F)
	}
	for i, v := range loc.Gradient {
		if math.IsInf(v, 0) || math.IsNaN(v) {
			return loc, ErrGrad{Grad: v, Index: i}
		}
	}

	return loc, nil
}

// checkConvergence returns NotTerminated if the Location does not satisfy the
// convergence criteria given by settings. Otherwise a corresponding status is
// returned.
// Unlike checkLimits, checkConvergence is called by Local only at MajorIterations.
func checkConvergence(loc *Location, settings *Settings) Status {
	if loc.Gradient != nil {
		norm := floats.Norm(loc.Gradient, math.Inf(1))
		if norm < settings.GradientThreshold {
			return GradientThreshold
		}
	}

	if loc.F < settings.FunctionThreshold {
		return FunctionThreshold
	}

	if settings.FunctionConverge != nil {
		return settings.FunctionConverge.FunctionConverged(loc.F)
	}

	return NotTerminated
}

// checkLimits returns NotTerminated status if the various limits given by
// settings has not been reached. Otherwise it returns a corresponding status.
// Unlike checkConvergence, checkLimits is called by Local at _every_ iteration.
func checkLimits(loc *Location, stats *Stats, settings *Settings) Status {
	// Check the objective function value for negative infinity because it
	// could break the linesearches and -inf is the best we can do anyway.
	if math.IsInf(loc.F, -1) {
		return FunctionNegativeInfinity
	}

	if settings.MajorIterations > 0 && stats.MajorIterations >= settings.MajorIterations {
		return IterationLimit
	}

	if settings.FuncEvaluations > 0 && stats.FuncEvaluations >= settings.FuncEvaluations {
		return FunctionEvaluationLimit
	}

	if settings.GradEvaluations > 0 && stats.GradEvaluations >= settings.GradEvaluations {
		return GradientEvaluationLimit
	}

	if settings.HessEvaluations > 0 && stats.HessEvaluations >= settings.HessEvaluations {
		return HessianEvaluationLimit
	}

	// TODO(vladimir-ch): It would be nice to update Runtime here.
	if settings.Runtime > 0 && stats.Runtime >= settings.Runtime {
		return RuntimeLimit
	}

	return NotTerminated
}

// evaluate evaluates the routines specified by the Operation at loc.X, stores
// the answer into loc and updates stats. loc.X is copied into x before
// evaluating in order to prevent the routines from modifying it.
func evaluate(p *Problem, loc *Location, eval Operation, stats *Stats, x []float64) {
	copy(x, loc.X)
	if eval&FuncEvaluation != 0 {
		loc.F = p.Func(x)
		stats.FuncEvaluations++
	}
	if eval&GradEvaluation != 0 {
		p.Grad(loc.Gradient, x)
		stats.GradEvaluations++
	}
	if eval&HessEvaluation != 0 {
		p.Hess(loc.Hessian, x)
		stats.HessEvaluations++
	}
}
