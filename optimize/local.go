// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"

	"gonum.org/v1/gonum/floats"
)

// localOptimizer is a helper type for running an optimization using a LocalMethod.
type localOptimizer struct{}

// run controls the optimization run for a localMethod. The calling method
// must close the operation channel at the conclusion of the optimization. This
// provides a happens before relationship between the return of status and the
// closure of operation, and thus a call to method.Status (if necessary).
func (l localOptimizer) run(method localMethod, gradThresh float64, operation chan<- Task, result <-chan Task, tasks []Task) (Status, error) {
	// Local methods start with a fully-specified initial location.
	task := tasks[0]
	task = l.initialLocation(operation, result, task, method)
	if task.Op == PostIteration {
		l.finish(operation, result)
		return NotTerminated, nil
	}
	status, err := l.checkStartingLocation(task, gradThresh)
	if err != nil {
		l.finishMethodDone(operation, result, task)
		return status, err
	}

	// Send a major iteration with the starting location.
	task.Op = MajorIteration
	operation <- task
	task = <-result
	if task.Op == PostIteration {
		l.finish(operation, result)
		return NotTerminated, nil
	}
	op, err := method.initLocal(task.Location)
	if err != nil {
		l.finishMethodDone(operation, result, task)
		return Failure, err
	}
	task.Op = op
	operation <- task
Loop:
	for {
		r := <-result
		switch r.Op {
		case PostIteration:
			break Loop
		case MajorIteration:
			// The last operation was a MajorIteration. Check if the gradient
			// is below the threshold.
			if status := l.checkGradientConvergence(r.Gradient, gradThresh); status != NotTerminated {
				l.finishMethodDone(operation, result, task)
				return GradientThreshold, nil
			}
			fallthrough
		default:
			op, err := method.iterateLocal(r.Location)
			if err != nil {
				l.finishMethodDone(operation, result, r)
				return Failure, err
			}
			r.Op = op
			operation <- r
		}
	}
	l.finish(operation, result)
	return NotTerminated, nil
}

// initialOperation returns the Operation needed to fill the initial location
// based on the needs of the method and the values already supplied.
func (localOptimizer) initialOperation(task Task, n needser) Operation {
	var newOp Operation
	op := task.Op
	if op&FuncEvaluation == 0 {
		newOp |= FuncEvaluation
	}
	needs := n.needs()
	if needs.Gradient && op&GradEvaluation == 0 {
		newOp |= GradEvaluation
	}
	if needs.Hessian && op&HessEvaluation == 0 {
		newOp |= HessEvaluation
	}
	return newOp
}

// initialLocation fills the initial location based on the needs of the method.
// The task passed to initialLocation should be the first task sent in RunGlobal.
func (l localOptimizer) initialLocation(operation chan<- Task, result <-chan Task, task Task, needs needser) Task {
	task.Op = l.initialOperation(task, needs)
	operation <- task
	return <-result
}

func (l localOptimizer) checkStartingLocation(task Task, gradThresh float64) (Status, error) {
	if math.IsInf(task.F, 1) || math.IsNaN(task.F) {
		return Failure, ErrFunc(task.F)
	}
	for i, v := range task.Gradient {
		if math.IsInf(v, 0) || math.IsNaN(v) {
			return Failure, ErrGrad{Grad: v, Index: i}
		}
	}
	status := l.checkGradientConvergence(task.Gradient, gradThresh)
	return status, nil
}

func (localOptimizer) checkGradientConvergence(gradient []float64, gradThresh float64) Status {
	if gradient == nil || math.IsNaN(gradThresh) {
		return NotTerminated
	}
	if gradThresh == 0 {
		gradThresh = defaultGradientAbsTol
	}
	if norm := floats.Norm(gradient, math.Inf(1)); norm < gradThresh {
		return GradientThreshold
	}
	return NotTerminated
}

// finish completes the channel operations to finish an optimization.
func (localOptimizer) finish(operation chan<- Task, result <-chan Task) {
	// Guarantee that result is closed before operation is closed.
	for range result {
	}
}

// finishMethodDone sends a MethodDone signal on operation, reads the result,
// and completes the channel operations to finish an optimization.
func (l localOptimizer) finishMethodDone(operation chan<- Task, result <-chan Task, task Task) {
	task.Op = MethodDone
	operation <- task
	task = <-result
	if task.Op != PostIteration {
		panic("optimize: task should have returned post iteration")
	}
	l.finish(operation, result)
}
