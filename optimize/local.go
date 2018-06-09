// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"

	"gonum.org/v1/gonum/floats"
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
	if method == nil {
		method = getDefaultMethod(&p)
	}
	if settings == nil {
		settings = DefaultSettings()
	}
	// Check that the initial location matches the one in settings.
	if settings.InitX != nil && !floats.Equal(settings.InitX, initX) {
		panic("local: initX does not match settings x location")
	}
	lg := &localGlobal{
		Method:   method,
		InitX:    initX,
		Settings: settings,
	}
	return Global(p, len(initX), settings, lg)
}

func getDefaultMethod(p *Problem) Method {
	if p.Grad != nil {
		return &BFGS{}
	}
	return &NelderMead{}
}

// localGlobal is a wrapper for Local methods to allow them to be optimized by Global.
type localGlobal struct {
	Method   Method
	InitX    []float64
	Settings *Settings

	dim    int
	status Status
	err    error
}

func (l *localGlobal) InitGlobal(dim, tasks int) int {
	if dim != len(l.InitX) {
		panic("optimize: initial length mismatch")
	}
	l.dim = dim
	l.status = NotTerminated
	l.err = nil
	return 1 // Local optimizations always run in serial.
}

func (l *localGlobal) Status() (Status, error) {
	return l.status, l.err
}

func (l *localGlobal) Needs() struct {
	Gradient bool
	Hessian  bool
} {
	return l.Method.Needs()
}

func (l *localGlobal) RunGlobal(operations chan<- GlobalTask, results <-chan GlobalTask, tasks []GlobalTask) {
	// Local methods start with a fully-specified initial location.
	task := tasks[0]
	op := l.getStartingLocation(operations, results, task)
	if op == PostIteration {
		l.cleanup(operations, results)
		return
	}
	// Check the starting condition.
	if math.IsInf(task.F, 1) || math.IsNaN(task.F) {
		l.status = Failure
		l.err = ErrFunc(task.F)
	}
	for i, v := range task.Gradient {
		if math.IsInf(v, 0) || math.IsNaN(v) {
			l.status = Failure
			l.err = ErrGrad{Grad: v, Index: i}
			break
		}
	}
	if l.status == Failure {
		l.exitFailure(operations, results, tasks[0])
		return
	}

	// Send a major iteration with the starting location.
	task.Op = MajorIteration
	operations <- task
	task = <-results
	if task.Op == PostIteration {
		l.cleanup(operations, results)
		return
	}

	op, err := l.Method.Init(task.Location)
	if err != nil {
		l.status = Failure
		l.err = err
		l.exitFailure(operations, results, tasks[0])
		return
	}
	task.Op = op
	operations <- task
Loop:
	for {
		result := <-results
		switch result.Op {
		case PostIteration:
			break Loop
		default:
			op, err := l.Method.Iterate(result.Location)
			if err != nil {
				l.status = Failure
				l.err = err
				l.exitFailure(operations, results, result)
				return
			}
			result.Op = op
			operations <- result
		}
	}
	l.cleanup(operations, results)
}

// exitFailure cleans up from a failure of the local method.
func (l *localGlobal) exitFailure(operation chan<- GlobalTask, result <-chan GlobalTask, task GlobalTask) {
	task.Op = MethodDone
	operation <- task
	task = <-result
	if task.Op != PostIteration {
		panic("task should have returned post iteration")
	}
	l.cleanup(operation, result)
}

func (l *localGlobal) cleanup(operation chan<- GlobalTask, result <-chan GlobalTask) {
	// Guarantee that result is closed before operation is closed.
	for range result {
	}
	close(operation)
}

func (l *localGlobal) getStartingLocation(operation chan<- GlobalTask, result <-chan GlobalTask, task GlobalTask) Operation {
	copy(task.X, l.InitX)
	// Construct the operation by what is missing.
	needs := l.Method.Needs()
	initOp := task.Op
	op := NoOperation
	if initOp&FuncEvaluation == 0 {
		op |= FuncEvaluation
	}
	if needs.Gradient && initOp&GradEvaluation == 0 {
		op |= GradEvaluation
	}
	if needs.Hessian && initOp&HessEvaluation == 0 {
		op |= HessEvaluation
	}

	if op == NoOperation {
		return NoOperation
	}
	task.Op = op
	operation <- task
	task = <-result
	return task.Op
}
