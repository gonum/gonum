// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

// Status represents the status of the optimization. Statuses that are greater than
// one represent a sufficiently good value was found. Statuses less than one
// signify the optimization was terminated before finding a minimum.
type Status int

const NotTerminated Status = 0

const (
	Success Status = iota + 1
	FunctionAbsoluteConvergence
	GradientAbsoluteConvergence
	StepConvergence
)

const (
	Failure Status = -(iota + 1)
	IterationLimit
	RuntimeLimit
	FunctionEvaluationLimit
	GradientEvaluationLimit
	RecorderError
	UserFunctionError
	MethodError
)

func (s Status) String() string {
	return statusMap[s]
}

var statusMap = map[Status]string{
	NotTerminated:               "NotTerminated",
	Success:                     "Success",
	FunctionAbsoluteConvergence: "FunctionAbsoluteConvergence",
	GradientAbsoluteConvergence: "GradientAbsoluteConvergence",
	StepConvergence:             "StepConvergence",

	Failure:                 "Failure",
	IterationLimit:          "IterationLimit",
	RuntimeLimit:            "RuntimeLimit",
	FunctionEvaluationLimit: "FunctionEvaluationLimit",
	GradientEvaluationLimit: "GradientEvaluationLimit",
	RecorderError:           "RecorderError",
	UserFunctionError:       "UserFunctionError",
	MethodError:             "MethodError",
}

var minStatus = -100
var maxStatus = 100

// NewStatus returns a unique Status variable to represent a custom status.
// NewStatus is intended to be called only during package initialization , and
// calls to NewStatus are not thread safe.
//
// NewStatus takes in two arguments, the string that should be output from
// Status.String(), and a boolean if the status should have a positive or
// negative value.
func NewStatus(printString string, good bool) Status {
	var s Status
	if good {
		s = Status(maxStatus)
		maxStatus++
	} else {
		s = Status(minStatus)
		minStatus--
	}
	statusMap[s] = printString
	return s
}
