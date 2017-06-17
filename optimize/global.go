// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"
	"sync"
	"time"
)

// GlobalMethod is a global optimizer. Typically will require more function
// evaluations and no sense of local convergence
type GlobalMethod interface {
	// Global tells method the max number of tasks, method returns how many it wants.
	// This is needed to sync the Global goroutines and inside goroutines.
	InitGlobal(dim, tasks int) int
	// Global method may assume that the same task id always has the same pointer with it.
	IterateGlobal(task int, loc *Location) (Operation, error)
	Needser
	// Done communicates to the optimization method that the optimization has
	// concluded to allow for shutdown.
	Done()
}

// Global uses a global optimizer to search for the global minimum of a
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
// returned Status is not NotTerminated or the error is not nil, the
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
// function evaluations. If you would like to put limits on this, for example
// maximum runtime or maximum function evaluations, modify the Settings
// input struct.
//
// Something about Global cannot guarantee strict bounds on function evaluations,
// iterations, etc. in the precense of concurrency.
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
	nTasks = method.InitGlobal(dim, nTasks)

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

type globalStatus struct {
	mux       *sync.RWMutex
	stats     *Stats
	status    Status
	p         *Problem
	startTime time.Time
	optLoc    *Location
	settings  *Settings
	statuser  Statuser
	err       error
}

func globalWorker(task int, m GlobalMethod, g *globalStatus, loc *Location, x []float64) {
	for {
		// Find Evaluation location
		op, err := m.IterateGlobal(task, loc)
		if err != nil {
			// TODO(btracey): Figure out how to handle errors properly. Shut
			// everything down? Pass to globalStatus so it can shut everything down?
			g.mux.Lock()
			g.err = err
			g.status = Failure
			g.mux.Unlock()
			break
		}

		// Evaluate location and/or update stats.
		status := g.globalOperation(op, loc, x)
		if status != NotTerminated {
			break
		}
	}
}

// globalOperation updates handles the status received by an individual worker.
// It uses a mutex to protect updates where necessary.
func (g *globalStatus) globalOperation(op Operation, loc *Location, x []float64) Status {
	// Do a quick check to see if one of the other workers converged in the meantime.
	var status Status
	var err error
	g.mux.RLock()
	status = g.status
	g.mux.RUnlock()
	if status != NotTerminated {
		return status
	}
	switch op {
	case NoOperation:
	case InitIteration:
		panic("optimize: Method returned InitIteration")
	case PostIteration:
		panic("optimize: Method returned PostIteration")
	case MajorIteration:
		g.mux.Lock()
		g.stats.MajorIterations++
		copyLocation(g.optLoc, loc)
		g.mux.Unlock()

		g.mux.RLock()
		status = checkConvergence(g.optLoc, g.settings, false)
		g.mux.RUnlock()
	default: // Any of the Evaluation operations.
		status, err = evaluate(g.p, loc, op, x)
		g.mux.Lock()
		updateStats(g.stats, op)
		g.mux.Unlock()
	}

	g.mux.Lock()
	status, err = iterCleanup(status, err, g.stats, g.settings, g.statuser, g.startTime, loc, op)
	// Update the termination status if it hasn't already terminated.
	if g.status == NotTerminated {
		g.status = status
		g.err = err
	}
	g.mux.Unlock()

	return status
}

func DefaultSettingsGlobal() *Settings {
	return &Settings{
		FunctionThreshold: math.Inf(-1),
		FunctionConverge: &FunctionConverge{
			Absolute:   1e-10,
			Iterations: 100,
		},
	}
}
