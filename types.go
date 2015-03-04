// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"fmt"
	"math"
	"time"
)

const defaultGradientAbsTol = 1e-6

// EvaluationType is used by the optimizer to specify information needed
// from the objective function.
type EvaluationType int

const (
	NoEvaluation EvaluationType = iota
	FuncEvaluation
	GradEvaluation
	HessEvaluation
	FuncGradEvaluation
	FuncGradHessEvaluation
)

func (e EvaluationType) String() string {
	if e < 0 || int(e) >= len(evaluationStrings) {
		return fmt.Sprintf("EvaluationType(%d)", e)
	}
	return evaluationStrings[e]
}

var evaluationStrings = [...]string{
	"NoEvaluation",
	"FuncEvaluation",
	"GradEvaluation",
	"HessEvaluation",
	"FuncGradEvaluation",
	"FuncGradHessEvaluation",
}

// IterationType specifies the type of iteration.
type IterationType int

const (
	NoIteration IterationType = iota
	MajorIteration
	MinorIteration
	SubIteration
	InitIteration
	PostIteration // Iteration after the optimization. Sent to Recorder.
)

func (i IterationType) String() string {
	if i < 0 || int(i) >= len(iterationStrings) {
		return fmt.Sprintf("IterationType(%d)", i)
	}
	return iterationStrings[i]
}

var iterationStrings = [...]string{
	"NoIteration",
	"MajorIteration",
	"MinorIteration",
	"SubIteration",
	"InitIteration",
	"PostIteration",
}

// Location represents a location in the optimization procedure.
type Location struct {
	X        []float64
	F        float64
	Gradient []float64
}

// LinesearchLocation is a location for a linesearch subiteration
type LinesearchLocation struct {
	F          float64 // Function value at the step
	Derivative float64 // Projected gradient in the linesearch direction
}

// Result represents the answer of an optimization run. It contains the optimum
// location as well as the Status at convergence and Statistics taken during the
// run.
type Result struct {
	Location
	Stats
	Status Status
}

// Stats contains the statistics of the run.
type Stats struct {
	MajorIterations         int           // Total number of major iterations
	FuncEvaluations         int           // Number of evaluations of Func()
	GradEvaluations         int           // Number of evaluations of Grad()
	HessEvaluations         int           // Number of evaluations of Hess()
	FuncGradEvaluations     int           // Number of evaluations of FuncGrad()
	FuncGradHessEvaluations int           // Number of evaluations of FuncGradHess()
	Runtime                 time.Duration // Total runtime of the optimization
}

// FunctionInfo is data to give to the optimizer about the objective function.
type FunctionInfo struct {
	IsGradient         bool
	IsFunctionGradient bool
	IsStatuser         bool
}

// functionInfo contains information about which interfaces the objective
// function F implements and the actual methods of F that have been
// successfully type switched.
type functionInfo struct {
	FunctionInfo

	function         Function
	gradient         Gradient
	functionGradient FunctionGradient
	statuser         Statuser
}

// Settings represents settings of the optimization run. It contains initial
// settings, convergence information, and Recorder information. In general, users
// should use DefaultSettings() rather than constructing a Settings literal.
//
// If UseInitData is true, InitialFunctionValue and InitialGradient specify
// function information at the initial location.
//
// If Recorder is nil, no information will be recorded.
type Settings struct {
	UseInitialData       bool      // Use supplied information about the conditions at the initial x.
	InitialFunctionValue float64   // Func(x) at the initial x.
	InitialGradient      []float64 // Grad(x) at the initial x.

	// FunctionAbsTol is the threshold for acceptably small values of the
	// objective function. FunctionAbsoluteConvergence status is returned if
	// the objective function is less than this value.
	// The default value is -inf.
	FunctionAbsTol float64

	// GradientAbsTol determines the accuracy to which the minimum is found.
	// GradientAbsoluteConvergence status is returned if the infinity norm of
	// the gradient is less than this value.
	// Has no effect if gradient information is not used.
	// The default value is 1e-6.
	GradientAbsTol float64

	// MajorIterations is the maximum number of iterations allowed.
	// IterationLimit status is returned if the number of major iterations
	// equals or exceeds this value.
	// If it equals zero, this setting has no effect.
	// The default value is 0.
	MajorIterations int

	// Runtime is the maximum runtime allowed. RuntimeLimit status is returned
	// if the duration of the run is longer than this value. Runtime is only
	// checked at iterations of the Method.
	// If it equals zero, this setting has no effect.
	// The default value is 0.
	Runtime time.Duration

	// FuncEvaluations is the maximum allowed number of function evaluations.
	// FunctionEvaluationLimit status is returned if the total number of
	// function evaluations equals or exceeds this number. Calls to Func() and
	// FuncGrad() are both counted as function evaluations for this calculation.
	// If it equals zero, this setting has no effect.
	// The default value is 0.
	FuncEvaluations int

	// GradEvaluations is the maximum allowed number of gradient evaluations.
	// GradientEvaluationLimit status is returned if the total number of
	// gradient evaluations equals or exceeds this number. Calls to Grad() and
	// FuncGrad() are both counted as gradient evaluations for this calculation.
	// If it equals zero, this setting has no effect.
	// The default value is 0.
	GradEvaluations int

	Recorder Recorder
}

// DefaultSettings returns a new Settings struct containing the default settings.
func DefaultSettings() *Settings {
	return &Settings{
		GradientAbsTol: defaultGradientAbsTol,
		FunctionAbsTol: math.Inf(-1),
		Recorder:       NewPrinter(),
	}
}

// resize takes x and returns a slice of length dim. It returns a resliced x
// if cap(x) >= dim, and a new slice otherwise.
func resize(x []float64, dim int) []float64 {
	if dim > cap(x) {
		return make([]float64, dim)
	}
	return x[:dim]
}
