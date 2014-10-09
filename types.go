// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"math"
	"time"
)

// EvalutationType is used by the optimizer to specify information needed
// from the objective function.
type EvaluationType int

const (
	NoEvaluation EvaluationType = iota
	FunctionAndGradient
	JustFunction
	JustGradient
)

func (e EvaluationType) String() string {
	return evaluationMap[e]
}

var evaluationMap = map[EvaluationType]string{
	NoEvaluation:        "NoEvaluation",
	FunctionAndGradient: "FunctionAndGradient",
	JustFunction:        "JustFunction",
	JustGradient:        "JustGradient",
}

// IterationType specifies the type of iteration.
type IterationType int

const (
	NoIteration IterationType = iota
	Major
	Minor
	Sub
)

func (i IterationType) String() string {
	return iterationMap[i]
}

var iterationMap = map[IterationType]string{
	NoIteration: "NoIteration",
	Major:       "Major",
	Minor:       "Minor",
	Sub:         "Sub",
}

// Location represents a location in the optimization procedure.
type Location struct {
	X        []float64
	F        float64
	Gradient []float64
}

// Result represents the answer of an optimizitaion run. It contains the optimum
// location as well as the Status at convergence and Statistics taken during the
// run.
type Result struct {
	Location
	Stats
	Status Status
}

// Stats contains the statistics of the run.
type Stats struct {
	NumMajorIterations int           // Total number of major iterations
	NumFunEvals        int           // Number of evaluations of F()
	NumGradEvals       int           // Number of evaluations of Df()
	NumFunGradEvals    int           // Number of evaluations of FDf()
	Runtime            time.Duration // Total runtime of the optimization
	GradNorm           float64       // 2-norm of the gradient normalized by the sqrt of the length
}

// FunctionStats is data to give to the optimizer about the objective function.
type FunctionStats struct {
	IsGradient bool
	IsFunGrad  bool
	IsStatuser bool

	// note: it's always a function
}

// functions contains the actual methods of F that have been successfully type
// switched
type functions struct {
	function Function
	gradient Gradient
	gradFunc FunctionGradient
	status   Statuser
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
	UseInitData          bool      // Use supplied information about the conitions at the initial x.
	InitialFunctionValue float64   // F(x) at the initial x.
	IntialGradient       []float64 // Df(x) at the initial x.

	// Converge if the objective function is less than this value.
	FunctionAbsoluteTolerance float64
	// Loosely, converge if the 'average' value of the gradient is less than this
	// value. Specifically, converge if ||grad||_2 / sqrt(len(grad)) is less than
	// this value.
	// Has no effect if gradient information is not used.
	GradientAbsoluteTolerance float64
	// Converge if the number of major iterations equals or exceeds this value.
	// If it equals zero, this setting has no effect.
	MaximumMajorIterations int
	// Converge if the duration of the run is longer than this value. Runtime
	// is only checked at iterations of the optimizer. If it equals zero,
	// this setting has no effect.
	MaximumRuntime time.Duration
	// Converge if the total number of function evaluations equals or exceeds this
	// number. Calls to F() and FDf() are both counted as function evaluations
	// for this calculation. If it equals zero, this setting has no effect.
	MaximumFunctionEvaluations int
	// Converge if the total number of gradient evaluations equals or exceeds this
	// number. Calls to D() and FDf() are both counted as gradient evaluations
	// for this calculation. If it equals zero, this setting has no effect.
	MaximumGradientEvaluations int

	Recorder Recorder
}

// DefaultSettings returns a new Settings struct containing the default settings.
func DefaultSettings() *Settings {
	return &Settings{
		GradientAbsoluteTolerance: 1e-6,
		FunctionAbsoluteTolerance: math.Inf(-1),
		Recorder:                  NewPrinter(),
	}
}
