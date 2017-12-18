// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"
	"sync"
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
	Index     int
	Operation Operation
	*Location
}

// GlobalMethod is a type which can search for a global optimum for an objective function.
type GlobalMethod interface {
	Needser
	// InitGlobal communicates the input dimension and maximum number of tasks,
	// and returns the number of concurrent processes. The return must be less
	// than or equal to tasks.
	InitGlobal(dim, tasks int) int
	// RunGlobal runs a global optimization using the GlobalMethod. The method
	// communicates function evaulations, major iterations, etc. using GlobalTasks.
	// The result of these tasks will be returned on Result. The optimization
	// will continue until PostIteration is sent on result, at which point no
	// more Evaluations will be performed. The details of the communication
	// scheme are provided below.
	//
	// The GlobalTask contains an Operation to evaluate, which will be performed
	// at the Location in the task. GlobalTask also has an Index field to identify
	// specific tasks. Also provided is a set of tasks with initialized memory
	// storage for the Location. The number provided tasks is equal to the value
	// returned from InitGlobal.
	//
	// Tasks are sent on the operation channel to be evaluated. The results
	// of these Operations are returned on result. The returned task will have
	// the same Index and Operation as the sent task, which can be used to
	// determine the next appropriate action.
	//
	// Once per RunGlobal, a PostIteration task will be sent on result, signaling the
	// conclusion of the optimization (i.e. convergence from user settings).
	// More tasks may still be sent on operation, after this occurs, but
	// no additional Evaluations will be conducted. MajorIteration operations,
	// however, will still be conducted appropriately.
	//
	// Any Evaluations sent before PostIteration will be evaluated and returned
	// on result. In order to successfully conclude the optimization, GlobalMethod
	// must read from the result channel until it is closed, handling appropriately.
	// For example, many GlobalMethods will want to send one or more
	// MajorIteration on operation while draining the result channel. After
	// all operations have been sent, the GlobalMethod must close the operation
	// channel. Method then must return from RunGlobal.
	//
	// The GlobalMethod may have its own specific convergence criteria. If
	// any of those are met, MethodDone should be sent on the operation channel.
	// This will trigger a PostIteration to be sent on result. The MethodDone
	// task will not be returned on result. If MethodDone is sent, the GlobalMethod
	// must be a Statuser, and the call to Status must return a Status other than
	// NotTerminated.
	//
	// The operation and result tasks have a buffer of the number of tasks.
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
// Please see the documentation for GlobalMethod for the details on implementing
// a method.
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

// minimizeGlobal is the high-level function for a Global optimization. It launches
// concurrent workers to perform the mimization, and shuts them down properly
// at the conclusion.
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

	// Algorithm overview:
	// At a high level, we want to enable function evaluations in parallel.
	// There are a couple challenges with this.
	// - There are statistics that need to be tracked (i.e. function evaluations),
	//   and these are inherently serial.
	// - The method can specify MajorIterations, and there needs to be a clear
	//   order to these calls (see #339, #344).
	// - There are several ways the optimization can be stopped, and when any
	//   of those ways happens the optimization needs to shut down. Specifically:
	//   - MajorIteration can trigger the end (IterationLimit, FunctionThreshold)
	//   - An Evaluation can trigger the end (FunctionEvaluationLimit)
	//   - The method can trigger the end (MethodDone)
	// - The shut down of the optimization can happen while functions are still
	//   being evaluated. We want to enable the results of those evaluations to
	//   be useful.
	//
	// The strategy is as follows:
	// - Pass two channels to GlobalMethod via the RunGlobal call. The method
	//   sends Operations to run on the operation channel, and this caller
	//   script will send the results on the results chan.
	// - A distributor receives tasks from the operation channel, and delegates.
	//   If the task is an Evaluation, it is sent to the workers, otherwise the
	//   Operation is sent to the stats updater.
	// - A set of workers are launched to evaluate functions. A worker evaluates
	//   the function, and then sends the Operation to the stats updater.
	// - A stats updater reads in the completed Operations, updates the stats,
	//   and checks convergence. The single point of combination prevents race
	//   conditions in updating the stats and checking convergence. This also
	//   ensures that MajorIteration updates happen in a predictable order.
	// - When a termination condition is met, a Status is sent on statusChan.
	//   only the first such termination signal is sent. This triggers a
	//   PostIteration to be sent to Method, and triggers a series of shutdown steps.

	var (
		finalStatus Status
		finalError  error
	)
	type final struct {
		Status Status
		Err    error
	}

	operations := make(chan GlobalTask, nTasks) // GlobalMethod sends tasks.
	results := make(chan GlobalTask, nTasks)    // return results to GlobalMethod.
	workerChan := make(chan GlobalTask)         // Delegate tasks to the workers.
	evalStatsChan := make(chan GlobalTask)      // Send evaluation updates.
	statusChan := make(chan final)              // Send a termination status.
	statusSent := make(chan struct{})           // Communicate the optimization is done.
	var statWG sync.WaitGroup                   // All stats are updated.
	var distWG sync.WaitGroup                   // Distributor has successfully terminated.

	// Launch the method.
	go func() {
		tasks := make([]GlobalTask, nTasks)
		for i := range tasks {
			tasks[i].Location = newLocation(dim, method)
		}
		method.RunGlobal(operations, results, tasks)
	}()

	// Launch the distributor
	distWG.Add(1)
	go func() {
		defer distWG.Done()
		// Note: This cannot be a range loop over operations because we want to
		// still send to operation even after a termination Status has been sent.
	Outer:
		for {
			select {
			case <-statusSent:
				break Outer
			case task := <-operations: // Delegate the Operation.
				switch task.Operation {
				case InitIteration:
					panic("optimize: GlobalMethod returned InitIteration")
				case PostIteration:
					panic("optimize: GlobalMethod returned PostIteration")
				case NoOperation, MajorIteration, MethodDone:
					evalStatsChan <- task
				default: // Any of the Evaluation operations.
					workerChan <- task
				}
			}
		}
	}()

	// Launch the workers. After workerChan is closed, each worker signals
	// to the stats updater that they have successfully terminated.
	for worker := 0; worker < nTasks; worker++ {
		go func() {
			x := make([]float64, dim)
			for task := range workerChan {
				evaluate(prob, task.Location, task.Operation, x)
				evalStatsChan <- task
			}
			evalStatsChan <- GlobalTask{Operation: signalOperation}
		}()
	}

	// Launch the stats combiner.
	statWG.Add(1)
	go func() {
		defer statWG.Done()
		var workerDone int // effective wg for the workers
		var status Status
		var err error
		for task := range evalStatsChan {
			switch task.Operation {
			default:
				updateEvaluationStats(stats, task.Operation)
				status, err = checkEvaluationLimits(prob, stats, settings)
			case signalOperation:
				workerDone++
				if workerDone == nTasks {
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
					panic("optimize: global method returned MethodDone is not a Statuser")
				}
				status, err = statuser.Status()
				if status == NotTerminated {
					panic("optimize: global method returned MethodDone but a NotTerminated status")
				}
			}
			if settings.Recorder != nil && status == NotTerminated && err == nil {
				stats.Runtime = time.Since(startTime)
				// Allow err to be overloaded if the Recorder fails.
				err = settings.Recorder.Record(task.Location, task.Operation, stats)
				if err != nil {
					status = Failure
				}
			}
			// If the optimization should be over, send the signal.
			if status != NotTerminated || err != nil {
				f := final{
					Status: status,
					Err:    err,
				}
				// Two cases: 1) A NotTerminated status has already been sent on
				// statusChan, and so we can send the first here. 2) A NotTerminated
				// status has already been sent. This means statusSent is closed
				// and can be received from.
				select {
				case <-statusSent:
				case statusChan <- f:
					// Tell the distributor to quit, and send a PostIteration
					close(statusSent)
					results <- GlobalTask{
						Operation: PostIteration,
					}
				}
			}

			// If there are still workers remaining, send the result back.
			// Don't send MethodDone back.
			if task.Operation != MethodDone {
				if workerDone != nTasks {
					results <- task
				}
			}
		}
	}()

	// Wait until a termination is sent.
	f := <-statusChan
	finalStatus = f.Status
	finalError = f.Err

	// Sending the status triggered the distributor to quit. Wait until it does.
	distWG.Wait()
	close(workerChan)
	// Read the final operations, only performing an action if MajorIteration.
	for task := range operations {
		switch task.Operation {
		case MajorIteration:
			evalStatsChan <- task
		}
	}
	// All stats have been sent. Close the channel and wait until all are updated.
	close(evalStatsChan)
	statWG.Wait()
	return finalStatus, finalError
}
