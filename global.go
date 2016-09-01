// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gonum/matrix/mat64"
)

// GlobalMethod is a global optimizer. Typically will require more function
// evaluations and no sense of local convergence
type GlobalMethod interface {
	// Global tells method the max number of tasks, method returns how many it wants.
	// This is needed to sync the Global goroutines and inside goroutines.
	InitGlobal(tasks int) int
	// Global method may assume that the same task id always has the same pointer with it.
	IterateGlobal(task int, loc *Location) (Operation, error)
	Needser
	// Done communicates to the optimization method that the optimization has
	// concluded to allow for shutdown.
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
	if p.Func == nil {
		panic("optimize: objective function is undefined")
	}
	if dim <= 0 {
		panic("optimize: impossible problem dimension")
	}
	startTime := time.Now()
	if method == nil {
		method = &GuessAndCheck{}
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
		settings = DefaultSettingsGlobal()
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
	optLoc := newLocation(dim, method)
	optLoc.F = math.Inf(1)

	if settings.FunctionConverge != nil {
		settings.FunctionConverge.Init(optLoc.F)
	}

	stats.Runtime = time.Since(startTime)

	// Don't need to check convergence because it can't possibly have converged.
	// (No function evaluations and no starting location).
	var err error
	if settings.Recorder != nil {
		err = settings.Recorder.Record(optLoc, InitIteration, stats)
		// TODO(btracey): Handle this error? Fix when merge with Local.
	}

	var status Status
	status, err = minimizeGlobal(&p, method, settings, stats, optLoc, startTime)

	// Cleanup and collect results
	if settings.Recorder != nil && err == nil {
		// Send the optimal location to Recorder.
		err = settings.Recorder.Record(optLoc, PostIteration, stats)
		// TODO(btracey): Handle this error? Fix when merge with Local.
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
	gs := &globalStatus{
		mux:       &sync.RWMutex{},
		stats:     stats,
		status:    NotTerminated,
		p:         p,
		startTime: startTime,
		optLoc:    optLoc,
		settings:  settings,
	}

	nTasks := settings.Concurrent
	nTasks = method.InitGlobal(nTasks)

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
	method    GlobalMethod
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
	g.mux.RLock()
	s := g.status
	g.mux.RUnlock()
	if s != NotTerminated {
		return s
	}
	switch op {
	case NoOperation:
	case InitIteration:
		panic("optimize: GlobalMethod return InitIteration")
	case PostIteration:
		panic("optimize: Method returned PostIteration")
	case MajorIteration:
		g.mux.Lock()
		g.stats.MajorIterations++
		if loc.F < g.optLoc.F {
			copyLocation(g.optLoc, loc)
		}
		g.mux.Unlock()
		g.mux.RLock()
		status := checkConvergenceGlobal(g.optLoc, g.settings)
		g.mux.RUnlock()
		if status != NotTerminated {
			// Update g.status, preserving the first termination status.
			g.mux.Lock()
			if g.status == NotTerminated {
				g.status = status
			}
			status = g.status
			g.mux.Unlock()
			return status
		}
	default:
		if !op.isEvaluation() {
			panic(fmt.Sprintf("optimize: invalid evaluation %v", op))
		}
		copy(x, loc.X)
		if op&FuncEvaluation != 0 {
			loc.F = g.p.Func(x)
			g.mux.Lock()
			g.stats.FuncEvaluations++
			g.mux.Unlock()
		}
		if op&GradEvaluation != 0 {
			g.p.Grad(loc.Gradient, x)
			g.mux.Lock()
			g.stats.GradEvaluations++
			g.mux.Unlock()
		}
		if op&HessEvaluation != 0 {
			g.p.Hess(loc.Hessian, x)
			g.mux.Lock()
			g.stats.HessEvaluations++
			g.mux.Unlock()
		}
	}

	// TODO(btracey): Need to fix all these things to avoid deadlock.
	// When re-do, need to make sure aren't overwritting a converged status.
	g.mux.Lock()
	g.stats.Runtime = time.Since(g.startTime)
	if g.settings.Recorder != nil {
		err := g.settings.Recorder.Record(loc, op, g.stats)
		if err != nil {
			if g.status == NotTerminated && g.err != nil {
				g.status = Failure
				g.err = err
			}
		}
	}
	s = checkLimits(loc, g.stats, g.settings)
	if g.status == NotTerminated {
		g.status = s
	}
	methodStatus, methodIsStatuser := g.method.(Statuser)
	if methodIsStatuser {
		s, err := methodStatus.Status()
		if err != nil && g.status == NotTerminated {
			g.status = s
			g.err = err
		}
	}
	s = g.status
	g.mux.Unlock()
	return s
}

func newLocation(dim int, method Needser) *Location {
	// TODO(btracey): combine this with Local.
	loc := &Location{
		X: make([]float64, dim),
	}
	loc.F = math.Inf(1)
	if method.Needs().Gradient {
		loc.Gradient = make([]float64, dim)
	}
	if method.Needs().Hessian {
		loc.Hessian = mat64.NewSymDense(dim, nil)
	}
	return loc
}

func checkConvergenceGlobal(loc *Location, settings *Settings) Status {
	if loc.F < settings.FunctionThreshold {
		return FunctionThreshold
	}
	if settings.FunctionConverge != nil {
		status := settings.FunctionConverge.FunctionConverged(loc.F)
		if status != NotTerminated {
			return NotTerminated
		}
	}
	return NotTerminated
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
