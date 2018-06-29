// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import "math"

// LocalController controls the optimization process for a LocalMethod.
// TODO(btracey): Rename?
type LocalController struct {
	Method Method

	status Status
	err    error
}

func (l *LocalController) Status() (Status, error) {
	return l.status, l.err
}

func (l *LocalController) InitGlobal(dim, tasks int) int {
	l.status = NotTerminated
	l.err = nil
	return 1 // Local optimizations always run in serial.
}

func (l *LocalController) RunGlobal(operation chan<- GlobalTask, result <-chan GlobalTask, tasks []GlobalTask) {
	// Local methods start with a fully-specified initial location.
	task := tasks[0]
	task = l.initialLocation(operation, result, task, l.Method)
	if task.Op == PostIteration {
		l.finish(operation, result)
		return
	}
	status, err := l.checkStartingLocation(task)
	if err != nil {
		l.status = status
		l.err = err
		l.finishMethodDone(operation, result, task)
		return
	}

	// Send a major iteration with the starting location.
	task.Op = MajorIteration
	operation <- task
	task = <-result
	if task.Op == PostIteration {
		l.finish(operation, result)
		return
	}

	op, err := l.Method.Init(task.Location)
	if err != nil {
		l.status = Failure
		l.err = err
		l.finishMethodDone(operation, result, task)
		return
	}
	task.Op = op
	operation <- task
Loop:
	for {
		r := <-result
		switch r.Op {
		case PostIteration:
			break Loop
		default:
			op, err := l.Method.Iterate(r.Location)
			if err != nil {
				l.status = Failure
				l.err = err
				l.finishMethodDone(operation, result, r)
				return
			}
			r.Op = op
			operation <- r
		}
	}
	l.finish(operation, result)
}

// initialOperation returns the Operation needed to fill the initial location based on
// the needs of the method and the values already supplied.
func (LocalController) initialOperation(task GlobalTask, needser Needser) Operation {
	needs := needser.Needs()
	var newOp Operation
	op := task.Op
	if op&FuncEvaluation == 0 {
		newOp |= FuncEvaluation
	}
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
func (l *LocalController) initialLocation(operation chan<- GlobalTask, result <-chan GlobalTask, task GlobalTask, needser Needser) GlobalTask {
	op := l.initialOperation(task, needser)
	task.Op = op
	operation <- task
	task = <-result
	return task
}

func (*LocalController) checkStartingLocation(task GlobalTask) (Status, error) {
	if math.IsInf(task.F, 1) || math.IsNaN(task.F) {
		return Failure, ErrFunc(task.F)
	}
	for i, v := range task.Gradient {
		if math.IsInf(v, 0) || math.IsNaN(v) {
			return Failure, ErrGrad{Grad: v, Index: i}
		}
	}
	return NotTerminated, nil
}

// cleaup performs the channel operations to finish an optimization run.
// The Method should return after this is called.
func (*LocalController) finish(operation chan<- GlobalTask, result <-chan GlobalTask) {
	// Guarantee that result is closed before operation is closed.
	for range result {
	}
	close(operation)
}

// cleaup performs the channel operations to finish an optimization run when a
// MethodDone signal must first be sent. The Method should return after this is called.
func (l *LocalController) finishMethodDone(operation chan<- GlobalTask, result <-chan GlobalTask, task GlobalTask) {
	task.Op = MethodDone
	operation <- task
	task = <-result
	if task.Op != PostIteration {
		panic("task should have returned post iteration")
	}
	l.finish(operation, result)
}
