// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"
	"sync"
	"time"
)

// GlobalMethod is an optimization method which seeks to find the global minimum
// of an objective function.
//
// At the beginning of the optimization, InitGlobal is called to communicate
// the dimension of the input and maximum number of concurrent tasks.
// The actual number of concurrent tasks will be set from the return of InitGlobal,
// which must not be greater than the input tasks.
//
// During the optimization, a reverse-communication interface is used between
// the GlobalMethod and the caller.
// GlobalMethod acts as a client that asks the caller to perform
// needed operations given the return from IterateGlobal.
// This allows and enforces automation of maintaining statistics and checking for
// (various types of) convergence.
//
// The return from IterateGlobal can be an Evaluation, a MajorIteration or NoOperation.
//
// An evaluation is one or more of the Evaluation operations (FuncEvaluation,
// GradEvaluation, etc.) combined with the bitwise or operator. In an evaluation
// operation, the requested fields of Problem will be evaluated at the value
// in Location.X, filling the corresponding fields of Location. These values
// can be retrieved and used upon the next call to IterateGlobal with that task id.
// The GlobalMethod interface requires that entries of Location are not modified
// aside from the commanded evaluations. Thus, the type implementing GlobalMethod
// may use multiple Operations to set the Location fields at a particular x value.
//
// When IterateGlobal declares MajorIteration, the caller updates the optimal
// location to the values in Location, and checks for convergence. The type
// implementing GlobalMethod must make sure that the fields of Location are valid
// and consistent.
//
// IterateGlobal must not return InitIteration and PostIteration operations. These are
// reserved for the clients to be passed to Recorders. A Method must also not
// combine the Evaluation operations with the Iteration operations.
type GlobalMethod interface {
	Needser
	// InitGlobal communicates the input dimension and maximum number of tasks,
	// and returns the number of concurrent processes. The return must be less
	// than or equal to tasks.
	InitGlobal(dim, tasks int) int

	// IterateGlobal retrieves information from the location associated with
	// the given task ID, and returns the next operation to perform with that
	// Location. IterateGlobal may assume that the same pointer is associated
	// with the same task.
	IterateGlobal(task int, loc *Location) (Operation, error)

	// Done communicates that the optimization has concluded to allow for shutdown.
	// After Done is called, no more calls to IterateGlobal will be made.
	Done()
}

// Global uses a global optimizer to search for the gloabl minimum of a
// function. A maximization problem can be transformed into a
// minimization problem by multiplying the function by -1.
//
// The first argument represents the problem to be minimized. Its fields are
// routines that evaluate the objective function, gradient, and other
// quantities related to the problem. The objective function, p.Func, must not
// be nil. The optimization method used may require other fields to be non-nil
// as specified by method.Needs. Global will panic if these are not met. The
// method can be determined automatically from the supplied problem which is
// described below.
//
// If p.Status is not nil, it is called before every evaluation. If the
// returned Status is other than NotTerminated or if the error is not nil, the
// optimization run is terminated.
//
// The third argument contains the settings for the minimization. The
// DefaultGlobalSettings function can be called for a Settings struct with the
// default values initialized. If settings == nil, the default settings are used.
// Global optimization methods typically do not make assumptions about the number
// and location of local minima. Thus, the only convergence metric used is the
// function values found at major iterations of the optimization. Bounds on the
// length of optimization are obeyed, such as the number of allowed function
// evaluations.
//
// The final argument is the optimization method to use. If method == nil, then
// an appropriate default is chosen based on the properties of the other arguments
// (dimension, gradient-free or gradient-based, etc.).
//
// If method implements Statuser, method.Status is called before every call
// to method.Iterate. If the returned Status is not NotTerminated or the
// error is non-nil, the optimization run is terminated.
//
// Global returns a Result struct and any error that occurred. See the
// documentation of Result for more information.
//
// Be aware that the default behavior of Global is to find the minimum.
// For certain functions and optimization methods, this process can take many
// function evaluations. The Settings input struct can be used to limit this,
// for example by modifying the maximum runtime or maximum function evaluations.
//
// Global cannot guarantee strict adherence to the bounds specified in Settings
// when performing concurrent evaluations and updates.
func Global(p Problem, dim int, settings *Settings, method GlobalMethod) (*Result, error) {
	startTime := time.Now()
	if method == nil {
		method = &GuessAndCheck{}
	}
	if settings == nil {
		settings = DefaultSettingsGlobal()
	}
	stats := &Stats{}
	err := checkOptimization(p, dim, method, settings.Recorder)
	if err != nil {
		return nil, err
	}

	optLoc := newLocation(dim, method)
	optLoc.F = math.Inf(1)

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

	// Run optimization
	var status Status
	status, err = minimizeGlobal(&p, method, settings, stats, optLoc, startTime)

	// Cleanup and collect results
	if settings.Recorder != nil && err == nil {
		err = settings.Recorder.Record(optLoc, PostIteration, stats)
	}
	stats.Runtime = time.Since(startTime)
	return &Result{
		Location: *optLoc,
		Stats:    *stats,
		Status:   status,
	}, err
}

// minimizeGlobal is the high-level function for a Global optimization. It launches
// concurrent workers to perform the mimization, and shuts them down properly
// at the conclusion.
func minimizeGlobal(p *Problem, method GlobalMethod, settings *Settings, stats *Stats, optLoc *Location, startTime time.Time) (status Status, err error) {
	dim := len(optLoc.X)
	statuser, _ := method.(Statuser)
	gs := &globalStatus{
		mux:       &sync.RWMutex{},
		stats:     stats,
		status:    NotTerminated,
		p:         p,
		startTime: startTime,
		optLoc:    optLoc,
		settings:  settings,
		statuser:  statuser,
	}

	nTasks := settings.Concurrent
	newNTasks := method.InitGlobal(dim, nTasks)
	if newNTasks > nTasks {
		panic("global: too many tasks returned by GlobalMethod")
	}
	nTasks = newNTasks

	// Launch optimization workers
	var wg sync.WaitGroup
	for task := 0; task < nTasks; task++ {
		wg.Add(1)
		go func(task int) {
			defer wg.Done()
			loc := newLocation(dim, method)
			x := make([]float64, dim)
			globalWorker(task, method, gs, loc, x)
		}(task)
	}
	wg.Wait()
	method.Done()
	return gs.status, gs.err
}

// globalWorker runs the optimization steps for a single (concurrently-executing)
// optimization task.
func globalWorker(task int, m GlobalMethod, g *globalStatus, loc *Location, x []float64) {
	for {
		// Find Evaluation location
		op, err := m.IterateGlobal(task, loc)
		if err != nil {
			g.updateStatus(Failure, err)
			break
		}

		// Evaluate location and/or update stats.
		status := g.globalOperation(op, loc, x)
		if status != NotTerminated {
			break
		}
	}
}

// globalStatus coordinates access to information shared between concurrently
// executing optimization tasks.
type globalStatus struct {
	mux       *sync.RWMutex
	stats     *Stats
	status    Status
	p         *Problem
	startTime time.Time
	optLoc    *Location
	settings  *Settings
	method    GlobalMethod
	statuser  Statuser
	err       error
}

// getStatus returns the current status of the optimization.
func (g *globalStatus) getStatus() Status {
	var status Status
	g.mux.RLock()
	defer g.mux.RUnlock()
	status = g.status
	return status
}

func (g *globalStatus) incrementMajorIteration() {
	g.mux.Lock()
	defer g.mux.Unlock()
	g.stats.MajorIterations++
}

func (g *globalStatus) updateOptLoc(loc *Location) {
	g.mux.Lock()
	defer g.mux.Unlock()
	copyLocation(g.optLoc, loc)
}

// checkConvergence checks the convergence of the global optimization and returns
// the status
func (g *globalStatus) checkConvergence() Status {
	g.mux.RLock()
	defer g.mux.RUnlock()
	return checkConvergence(g.optLoc, g.settings, false)
}

// updateStats updates the evaluation statistics for the given operation.
func (g *globalStatus) updateStats(op Operation) {
	g.mux.Lock()
	defer g.mux.Unlock()
	updateEvaluationStats(g.stats, op)
}

// updateStatus updates the status and error fields of g. This update only happens
// if status == NotTerminated, so that the first different status is the one
// maintained.
func (g *globalStatus) updateStatus(s Status, err error) {
	g.mux.Lock()
	defer g.mux.Unlock()
	if g.status != NotTerminated {
		g.status = s
		g.err = err
	}
}

func (g *globalStatus) finishIteration(status Status, err error, loc *Location, op Operation) (Status, error) {
	g.mux.Lock()
	defer g.mux.Unlock()
	return finishIteration(status, err, g.stats, g.settings, g.statuser, g.startTime, loc, op)
}

// globalOperation executes the requested operation at the given location.
// When modifying this function, keep in mind that it can be called concurrently.
// Uses of the internal fields should be through the methods of globalStatus and
// protected by a mutex where appropriate.
func (g *globalStatus) globalOperation(op Operation, loc *Location, x []float64) Status {
	// Do a quick check to see if one of the other workers converged in the meantime.
	status := g.getStatus()
	if status != NotTerminated {
		return status
	}
	var err error
	switch op {
	case NoOperation:
	case InitIteration:
		panic("optimize: GlobalMethod returned InitIteration")
	case PostIteration:
		panic("optimize: GlobalMethod returned PostIteration")
	case MajorIteration:
		g.incrementMajorIteration()
		g.updateOptLoc(loc)
		status = g.checkConvergence()
	default: // Any of the Evaluation operations.
		status, err = evaluate(g.p, loc, op, x)
		g.updateStats(op)
	}

	status, err = g.finishIteration(status, err, loc, op)
	if status != NotTerminated || err != nil {
		g.updateStatus(status, err)
	}
	return status
}

// DefaultSettingsGlobal returns the default settings for Global optimization.
func DefaultSettingsGlobal() *Settings {
	return &Settings{
		FunctionThreshold: math.Inf(-1),
		FunctionConverge: &FunctionConverge{
			Absolute:   1e-10,
			Iterations: 100,
		},
	}
}
