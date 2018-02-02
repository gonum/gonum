// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"
	"time"
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
	startTime := time.Now()
	dim := len(initX)
	if method == nil {
		method = getDefaultMethod(&p)
	}
	if settings == nil {
		settings = DefaultSettings()
	}

	stats := &Stats{}

	err := checkOptimization(p, dim, method, settings.Recorder)
	if err != nil {
		return nil, err
	}

	optLoc, err := getStartingLocation(&p, method, initX, stats, settings)
	if err != nil {
		return nil, err
	}

	if settings.FunctionConverge != nil {
		settings.FunctionConverge.Init(optLoc.F)
	}

	stats.Runtime = time.Since(startTime)

	// Send initial location to Recorder
	if settings.Recorder != nil {
		err = settings.Recorder.Record(optLoc, InitIteration, stats)
		if err != nil {
			return nil, err
		}
	}

	// Check if the starting location satisfies the convergence criteria.
	status := checkLocationConvergence(optLoc, settings)

	// Run optimization
	if status == NotTerminated && err == nil {
		// The starting location is not good enough, we need to perform a
		// minimization. The optimal location will be stored in-place in
		// optLoc.
		status, err = minimize(&p, method, settings, stats, optLoc, startTime)
	}

	// Cleanup and collect results
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
			status = performMajorIteration(optLoc, loc, stats, startTime, settings)
		case MethodDone:
			statuser, ok := method.(Statuser)
			if !ok {
				panic("optimize: method returned MethodDone is not a Statuser")
			}
			status, err = statuser.Status()
			if status == NotTerminated {
				panic("optimize: method returned MethodDone but a NotTerminated status")
			}
		default: // Any of the Evaluation operations.
			evaluate(p, loc, op, x)
			updateEvaluationStats(stats, op)
			status, err = checkEvaluationLimits(p, stats, settings)
		}
		if status != NotTerminated || err != nil {
			return
		}
		if settings.Recorder != nil {
			stats.Runtime = time.Since(startTime)
			err = settings.Recorder.Record(loc, op, stats)
			if err != nil {
				if status == NotTerminated {
					status = Failure
				}
				return status, err
			}
		}

		op, err = method.Iterate(loc)
		if err != nil {
			status = Failure
			return
		}
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
	loc := newLocation(dim, method)
	copy(loc.X, initX)

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
		evaluate(p, loc, eval, x)
		updateEvaluationStats(stats, eval)
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
