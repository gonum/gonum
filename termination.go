// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import "errors"

// Status represents the status of the optimization. Programs
// should not rely on the underlying numeric value of the Status being constant.
type Status int

const (
	NotTerminated Status = iota
	Success
	FunctionAbsoluteConvergence
	GradientAbsoluteConvergence
	StepConvergence
	FunctionNegativeInfinity
	Failure
	IterationLimit
	RuntimeLimit
	FunctionEvaluationLimit
	GradientEvaluationLimit
	RecorderError
	UserFunctionError
	MethodError
	lastProvidedStatus // this exists to make sure that len(statuses) is correct. New Stauses should be added above this
)

func (s Status) String() string {
	return statuses[s].name
}

// Early returns true if the status indicates the optimization ended before a
// minimum was found. As an example, if the maximum iterations was reached, a
// minimum was not found, but if the gradient norm was reached then a minimum
// was found.
func (s Status) Early() bool {
	return statuses[s].early
}

// Err returns the error associated with an early ending to the minimization. If
// Early returns false, Err will return nil.
func (s Status) Err() error {
	return statuses[s].err
}

var statuses = []struct {
	name  string
	early bool
	err   error
}{
	{
		name: "NotTerminated",
	},
	{
		name: "Success",
	},
	{
		name: "FunctionAbsoluteConvergence",
	},
	{
		name: "GradientAbsoluteConvergence",
	},
	{
		name: "StepConvergence",
	},
	{
		name: "FunctionNegativeInfinity",
	},
	{
		name:  "Failure",
		early: true,
		err:   errors.New("opt: termination ended in failure"),
	},
	{
		name:  "IterationLimit",
		early: true,
		err:   errors.New("opt: maximum number of major iterations reached."),
	},
	{
		name:  "RuntimeLimit",
		early: true,
		err:   errors.New("opt: maximum runtime reached."),
	},
	{
		name:  "FunctionEvaluationLimit",
		early: true,
		err:   errors.New("opt: maximum number of function evaluations reached."),
	},
	{
		name:  "GradientEvaluationLimit",
		early: true,
		err:   errors.New("opt: maximum number of gradient evaluations reached."),
	},
	{
		name:  "RecorderError",
		early: true,
		err:   errors.New("opt: minimizaton stopped due to error in the recorder."),
	},
	{
		name:  "UserFunctionError",
		early: true,
		err:   errors.New("opt: minimizaton stopped due to error in the user function."),
	},
	{
		name:  "MethodError",
		early: true,
		err:   errors.New("opt: minimizaton stopped due to error in the optimizer."),
	},
	{
		name: "lastProvidedStatus",
	},
}

var lastExistingStatus Status

func init() {
	lastExistingStatus = lastProvidedStatus
	if int(lastExistingStatus) != len(statuses)-1 {
		panic("length of statuses is not correct")
	}
}

// NewStatus returns a unique Status variable to represent a custom status.
// NewStatus is intended to be called only during package initialization, and
// calls to NewStatus are not thread safe.
//
// NewStatus takes in three arguments, the string that should be output from
// Status.String(), a boolean if the status indicates early optimization conclusion,
// and the error to return from Err (if any).
func NewStatus(name string, early bool, err error) Status {
	lastExistingStatus++
	statuses = append(statuses, struct {
		name  string
		early bool
		err   error
	}{name, early, err})
	return lastExistingStatus
}
