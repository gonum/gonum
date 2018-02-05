// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"
	"time"
)

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

// GlobalTask is a type to communicate between the GlobalMethod and the outer
// calling script.
type GlobalTask struct {
	ID int
	Op Operation
	*Location
}

// GlobalMethod is a type which can search for a global optimum for an objective function.
type GlobalMethod interface {
	Needser
	// InitGlobal takes as input the problem dimension and number of available
	// concurrent tasks, and returns the number of concurrent processes to be used.
	// The returned value must be less than or equal to tasks.
	InitGlobal(dim, tasks int) int
	// RunGlobal runs a global optimization. The method sends GlobalTasks on
	// the operation channel (for performing function evaluations, major
	// iterations, etc.). The result of the tasks will be returned on Result.
	// See the documentation for Operation types for the possible tasks.
	//
	// The caller of RunGlobal will signal the termination of the optimization
	// (i.e. convergence from user settings) by sending a task with a PostIteration
	// Op field on result. More tasks may still be sent on operation after this
	// occurs, but only MajorIteration operations will still be conducted
	// appropriately. Thus, it can not be guaranteed that all Evaluations sent
	// on operation will be evaluated, however if an Evaluation is started,
	// the results of that evaluation will be sent on results.
	//
	// The GlobalMethod must read from the result channel until it is closed.
	// During this, the GlobalMethod may want to send new MajorIteration(s) on
	// operation. GlobalMethod then must close operation, and return from RunGlobal.
	//
	// The las parameter to RunGlobal is a slice of tasks with length equal to
	// the return from InitGlobal. GlobalTask has an ID field which may be
	// set and modified by GlobalMethod, and must not be modified by the caller.
	//
	// GlobalMethod may have its own specific convergence criteria, which can
	// be communicated using a MethodDone operation. This will trigger a
	// PostIteration to be sent on result, and the MethodDone task will not be
	// returned on result. The GlobalMethod must implement Statuser, and the
	// call to Status must return a Status other than NotTerminated.
	//
	// The operation and result tasks are guaranteed to have a buffer length
	// equal to the return from InitGlobal.
	RunGlobal(operation chan<- GlobalTask, result <-chan GlobalTask, tasks []GlobalTask)
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
// returned Status is other than NotTerminated or if the error is not nil, the
// optimization run is terminated.
//
// The third argument contains the settings for the minimization. The
// DefaultGlobalSettings function can be called for a Settings struct with the
// default values initialized. If settings == nil, the default settings are used.
// All of the settings will be followed, but many of them may be counterproductive
// to use (such as GradientThreshold). Global cannot guarantee strict adherence
// to the bounds specified when performing concurrent evaluations and updates.
//
// The final argument is the optimization method to use. If method == nil, then
// an appropriate default is chosen based on the properties of the other arguments
// (dimension, gradient-free or gradient-based, etc.).
//
// Global returns a Result struct and any error that occurred. See the
// documentation of Result for more information.
//
// See the documentation for GlobalMethod for the details on implementing a method.
//
// Be aware that the default behavior of Global is to find the minimum.
// For certain functions and optimization methods, this process can take many
// function evaluations. The Settings input struct can be used to limit this,
// for example by modifying the maximum runtime or maximum function evaluations.
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

// minimizeGlobal performs a Global optimization. minimizeGlobal updates the
// settings and optLoc, and returns the final Status and error.
func minimizeGlobal(prob *Problem, method GlobalMethod, settings *Settings, stats *Stats, optLoc *Location, startTime time.Time) (Status, error) {
	dim := len(optLoc.X)
	nTasks := settings.Concurrent
	if nTasks == 0 {
		nTasks = 1
	}
	newNTasks := method.InitGlobal(dim, nTasks)
	if newNTasks > nTasks {
		panic("global: too many tasks returned by GlobalMethod")
	}
	nTasks = newNTasks

	// Launch the method. The method communicates tasks using the operations
	// channel, and results is used to return the evaluated results.
	operations := make(chan GlobalTask, nTasks)
	results := make(chan GlobalTask, nTasks)
	go func() {
		tasks := make([]GlobalTask, nTasks)
		for i := range tasks {
			tasks[i].Location = newLocation(dim, method)
		}
		method.RunGlobal(operations, results, tasks)
	}()

	// Algorithmic Overview:
	// There are three pieces to performing a concurrent global optimization,
	// the distributor, the workers, and the stats combiner. At a high level,
	// the distributor reads in tasks sent by method, sending evaluations to the
	// workers, and forwarding other operations to the statsCombiner. The workers
	// read these forwarded evaluation tasks, evaluate the relevant parts of Problem
	// and forward the results on to the stats combiner. The stats combiner reads
	// in results from the workers, as well as tasks from the distributor, and
	// uses them to update optimization statistics (function evaluations, etc.)
	// and to check optimization convergence.
	//
	// The complicated part is correctly shutting down the optimization. The
	// procedure is as follows. First, the stats combiner closes done and sends
	// a PostIteration to the method. The distributor then reads that done has
	// been closed, and closes the channel with the workers. At this point, no
	// more evaluation operations will be executed. As the workers finish their
	// evaluations, they forward the results onto the stats combiner, and then
	// signal their shutdown to the stats combiner. When all workers have successfully
	// finished, the stats combiner closes the results channel, signaling to the
	// method that all results have been collected. At this point, the method
	// may send MajorIteration(s) to update an optimum location based on these
	// last returned results, and then the method will close the operations channel.
	// Now that no more tasks will be commanded by the method, the distributor
	// closes statsChan, and with no more statistics to update the optimization
	// concludes.

	workerChan := make(chan GlobalTask) // Delegate tasks to the workers.
	statsChan := make(chan GlobalTask)  // Send evaluation updates.
	done := make(chan struct{})         // Communicate the optimization is done.

	// Read tasks from the method and distribute as appropriate.
	distributor := func() {
		for {
			select {
			case task := <-operations:
				switch task.Op {
				case InitIteration:
					panic("optimize: GlobalMethod returned InitIteration")
				case PostIteration:
					panic("optimize: GlobalMethod returned PostIteration")
				case NoOperation, MajorIteration, MethodDone:
					statsChan <- task
				default:
					if !task.Op.isEvaluation() {
						panic("global: expecting evaluation operation")
					}
					workerChan <- task
				}
			case <-done:
				// No more evaluations will be sent, shut down the workers, and
				// read the final tasks.
				close(workerChan)
				for task := range operations {
					if task.Op == MajorIteration {
						statsChan <- task
					}
				}
				close(statsChan)
				return
			}
		}
	}
	go distributor()

	// Evaluate the Problem concurrently.
	worker := func() {
		x := make([]float64, dim)
		for task := range workerChan {
			evaluate(prob, task.Location, task.Op, x)
			statsChan <- task
		}
		// Signal successful worker completion.
		statsChan <- GlobalTask{Op: signalDone}
	}
	for i := 0; i < nTasks; i++ {
		go worker()
	}

	var (
		workersDone int // effective wg for the workers
		status      Status
		err         error
		finalStatus Status
		finalError  error
	)

	// Update optimization statistics and check convergence.
	for task := range statsChan {
		switch task.Op {
		default:
			if !task.Op.isEvaluation() {
				panic("global: evaluation task expected")
			}
			updateEvaluationStats(stats, task.Op)
			status, err = checkEvaluationLimits(prob, stats, settings)
		case signalDone:
			workersDone++
			if workersDone == nTasks {
				close(results)
			}
			continue
		case NoOperation:
			// Just send the task back.
		case MajorIteration:
			status = performMajorIteration(optLoc, task.Location, stats, startTime, settings)
		case MethodDone:
			statuser, ok := method.(Statuser)
			if !ok {
				panic("optimize: global method returned MethodDone but is not a Statuser")
			}
			status, err = statuser.Status()
			if status == NotTerminated {
				panic("optimize: global method returned MethodDone but a NotTerminated status")
			}
		}
		if settings.Recorder != nil && status == NotTerminated && err == nil {
			stats.Runtime = time.Since(startTime)
			// Allow err to be overloaded if the Recorder fails.
			err = settings.Recorder.Record(task.Location, task.Op, stats)
			if err != nil {
				status = Failure
			}
		}
		// If this is the first termination status, trigger the conclusion of
		// the optimization.
		if status != NotTerminated || err != nil {
			select {
			case <-done:
			default:
				finalStatus = status
				finalError = err
				results <- GlobalTask{
					Op: PostIteration,
				}
				close(done)
			}
		}

		// Send the result back to the Problem if there are still active workers.
		if workersDone != nTasks && task.Op != MethodDone {
			results <- task
		}
	}
	return finalStatus, finalError
}
