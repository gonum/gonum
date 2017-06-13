// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"fmt"
	"math"
	"time"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

// newLocation allocates a new locatian structure of the appropriate size. It
// allocates memory based on the dimension and the values in Needs. The initial
// function value is set to math.Inf(1).
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
		loc.Hessian = mat.NewSymDense(dim, nil)
	}
	return loc
}

func copyLocation(dst, src *Location) {
	dst.X = resize(dst.X, len(src.X))
	copy(dst.X, src.X)

	dst.F = src.F

	dst.Gradient = resize(dst.Gradient, len(src.Gradient))
	copy(dst.Gradient, src.Gradient)

	if src.Hessian != nil {
		if dst.Hessian == nil || dst.Hessian.Symmetric() != len(src.X) {
			dst.Hessian = mat.NewSymDense(len(src.X), nil)
		}
		dst.Hessian.CopySym(src.Hessian)
	}
}

func checkOptimization(p Problem, dim int, method Needser, recorder Recorder) error {
	if p.Func == nil {
		panic(badProblem)
	}
	if dim <= 0 {
		panic("optimize: impossible problem dimension")
	}
	if err := p.satisfies(method); err != nil {
		return err
	}
	if p.Status != nil {
		_, err := p.Status()
		if err != nil {
			return err
		}
	}
	if recorder != nil {
		err := recorder.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

// evaluate evaluates the routines specified by the Operation at loc.X, and stores
// the answer into loc. loc.X is copied into x before
// evaluating in order to prevent the routines from modifying it.
func evaluate(p *Problem, loc *Location, op Operation, x []float64) (Status, error) {
	if !op.isEvaluation() {
		panic(fmt.Sprintf("optimize: invalid evaluation %v", op))
	}
	if p.Status != nil {
		status, err := p.Status()
		if err != nil || status != NotTerminated {
			return status, err
		}
	}
	copy(x, loc.X)
	if op&FuncEvaluation != 0 {
		loc.F = p.Func(x)
	}
	if op&GradEvaluation != 0 {
		p.Grad(loc.Gradient, x)
	}
	if op&HessEvaluation != 0 {
		p.Hess(loc.Hessian, x)
	}
	return NotTerminated, nil
}

// checkConvergence returns NotTerminated if the Location does not satisfy the
// convergence criteria given by settings. Otherwise a corresponding status is
// returned.
// Unlike checkLimits, checkConvergence is called only at MajorIterations.
//
// If local is true, gradient convergence is also checked.
func checkConvergence(loc *Location, settings *Settings, local bool) Status {
	if local && loc.Gradient != nil {
		norm := floats.Norm(loc.Gradient, math.Inf(1))
		if norm < settings.GradientThreshold {
			return GradientThreshold
		}
	}
	if loc.F < settings.FunctionThreshold {
		return FunctionThreshold
	}
	if settings.FunctionConverge != nil {
		return settings.FunctionConverge.FunctionConverged(loc.F)
	}
	return NotTerminated
}

// updateStats updates the statistics based on the operation.
func updateStats(stats *Stats, op Operation) {
	if op&FuncEvaluation != 0 {
		stats.FuncEvaluations++
	}
	if op&GradEvaluation != 0 {
		stats.GradEvaluations++
	}
	if op&HessEvaluation != 0 {
		stats.HessEvaluations++
	}
}

// checkLimits returns NotTerminated status if the various limits given by
// settings have not been reached. Otherwise it returns a corresponding status.
// Unlike checkConvergence, checkLimits is called by Local and Global at _every_
// iteration.
func checkLimits(loc *Location, stats *Stats, settings *Settings) Status {
	// Check the objective function value for negative infinity because it
	// could break the linesearches and -inf is the best we can do anyway.
	if math.IsInf(loc.F, -1) {
		return FunctionNegativeInfinity
	}

	if settings.MajorIterations > 0 && stats.MajorIterations >= settings.MajorIterations {
		return IterationLimit
	}

	if settings.FuncEvaluations > 0 && stats.FuncEvaluations >= settings.FuncEvaluations {
		return FunctionEvaluationLimit
	}

	if settings.GradEvaluations > 0 && stats.GradEvaluations >= settings.GradEvaluations {
		return GradientEvaluationLimit
	}

	if settings.HessEvaluations > 0 && stats.HessEvaluations >= settings.HessEvaluations {
		return HessianEvaluationLimit
	}

	// TODO(vladimir-ch): It would be nice to update Runtime here.
	if settings.Runtime > 0 && stats.Runtime >= settings.Runtime {
		return RuntimeLimit
	}

	return NotTerminated
}

// TODO(btracey): better name
func iterCleanup(status Status, err error, stats *Stats, settings *Settings, statuser Statuser, startTime time.Time, loc *Location, op Operation) (Status, error) {
	if status != NotTerminated || err != nil {
		return status, err
	}

	if settings.Recorder != nil {
		stats.Runtime = time.Since(startTime)
		err = settings.Recorder.Record(loc, op, stats)
		if err != nil {
			if status == NotTerminated {
				status = Failure
			}
			return status, err
		}
	}

	stats.Runtime = time.Since(startTime)
	status = checkLimits(loc, stats, settings)
	if status != NotTerminated {
		return status, nil
	}

	if statuser != nil {
		status, err = statuser.Status()
		if err != nil || status != NotTerminated {
			return status, err
		}
	}
	return status, nil
}
