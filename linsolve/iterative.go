// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linsolve

import (
	"errors"

	"gonum.org/v1/gonum/mat"
)

const defaultTolerance = 1e-8

// ErrIterationLimit is returned when a maximum number of iterations were done
// without converging to a solution.
var ErrIterationLimit = errors.New("linsolve: iteration limit reached")

// MulVecToer represents a square matrix A by means of a matrix-vector
// multiplication.
type MulVecToer interface {
	// MulVecTo computes A*x or Aᵀ*x and stores the result into dst.
	MulVecTo(dst *mat.VecDense, trans bool, x mat.Vector)
}

// Settings holds settings for solving a linear system.
type Settings struct {
	// InitX holds the initial guess. If it is nil or empty, the zero vector
	// will be used, otherwise its length must be equal to the dimension of
	// the system.
	InitX *mat.VecDense

	// Dst, if not nil, will be used for storing the approximate solution,
	// otherwise a new vector will be allocated. In both cases the vector will
	// also be returned in Result. If Dst is not empty, its length must be equal
	// to the dimension of the system.
	Dst *mat.VecDense

	// Tolerance specifies error tolerance for the final (approximate)
	// solution produced by the iterative method. The iteration will be
	// stopped when
	//  |r_i| < Tolerance * |b|
	// where r_i is the residual at i-th iteration.
	//
	// If Tolerance is zero, a default value of 1e-8 will be used, otherwise
	// it must be positive and less than 1.
	Tolerance float64

	// MaxIterations is the limit on the number of iterations. If it is
	// zero, a default value of twice the dimension of the system will be
	// used.
	MaxIterations int

	// PreconSolve describes a preconditioner solve that stores into dst the
	// solution of the system
	//  M  * dst = rhs, or
	//  Mᵀ * dst = rhs,
	// where M is the preconditioning matrix. If PreconSolve is nil, no
	// preconditioning will be used (M is the identity).
	PreconSolve func(dst *mat.VecDense, trans bool, rhs mat.Vector) error

	// Work context can be provided to reduce memory allocation when solving
	// multiple linear systems. If Work is not nil, its fields must be
	// either empty or their length must be equal to the dimension of the
	// system.
	Work *Context
}

// defaultSettings fills zero fields of s with default values.
func defaultSettings(s *Settings, dim int) {
	if s.InitX != nil && s.InitX.Len() == 0 {
		s.InitX.ReuseAsVec(dim)
	}
	if s.Dst == nil {
		s.Dst = mat.NewVecDense(dim, nil)
	} else if s.Dst.Len() == 0 {
		s.Dst.ReuseAsVec(dim)
	}
	if s.Tolerance == 0 {
		s.Tolerance = defaultTolerance
	}
	if s.MaxIterations == 0 {
		s.MaxIterations = 4 * dim
	}
	if s.PreconSolve == nil {
		s.PreconSolve = NoPreconditioner
	}
	if s.Work == nil {
		s.Work = NewContext(dim)
	} else {
		if s.Work.X.Len() == 0 {
			s.Work.X.ReuseAsVec(dim)
		}
		if s.Work.Src.Len() == 0 {
			s.Work.Src.ReuseAsVec(dim)
		}
		if s.Work.Dst.Len() == 0 {
			s.Work.Dst.ReuseAsVec(dim)
		}
	}
}

func checkSettings(s *Settings, dim int) {
	if s.InitX != nil && s.InitX.Len() != dim {
		panic("linsolve: mismatched length of initial guess")
	}
	if s.Dst.Len() != dim {
		panic("linsolve: mismatched destination length")
	}
	if s.Tolerance <= 0 || 1 <= s.Tolerance {
		panic("linsolve: invalid tolerance")
	}
	if s.MaxIterations <= 0 {
		panic("linsolve: negative iteration limit")
	}
	if s.Work.X.Len() != dim || s.Work.Src.Len() != dim || s.Work.Dst.Len() != dim {
		panic("linsolve: mismatched work context length")
	}
}

// Result holds the result of an iterative solve.
type Result struct {
	// X is the approximate solution.
	X *mat.VecDense

	// ResidualNorm is an approximation to the norm of the final residual.
	ResidualNorm float64

	// Stats holds statistics about the iterative solve.
	Stats Stats
}

// Stats holds statistics about an iterative solve.
type Stats struct {
	// Iterations is the number of iterations performed by Method.
	Iterations int

	// MulVec is the number of MulVec operations commanded by Method.
	MulVec int

	// PreconSolve is the number of PreconSolve operations commanded by Method.
	PreconSolve int
}

// Iterative finds an approximate solution of the system of n linear equations
//  A*x = b,
// where A is a nonsingular square matrix of order n and b is the right-hand
// side vector, using an iterative method m. If m is nil, default GMRES will be
// used.
//
// settings provide means for adjusting parameters of the iterative process. See
// the Settings documentation for more information. Iterative will not modify
// the fields of Settings. If settings is nil, default settings will be used.
//
// Note that the default choices of Method and Settings were chosen to provide
// accuracy and robustness, rather than speed. There are many algorithms for
// iterative linear solutions that have different tradeoffs, and can exploit
// special structure in the A matrix. Similarly, in many cases the number of
// iterations can be significantly reduced by using an appropriate
// preconditioner or increasing the error tolerance. Combined, these choices can
// significantly reduce computation time. Thus, while Iterative has supplied
// defaults, users are strongly encouraged to adjust these defaults for their
// problem.
func Iterative(a MulVecToer, b *mat.VecDense, m Method, settings *Settings) (*Result, error) {
	n := b.Len()

	var s Settings
	if settings != nil {
		s = *settings
	}
	defaultSettings(&s, n)
	checkSettings(&s, n)

	var stats Stats
	ctx := s.Work
	rInit := mat.NewVecDense(n, nil)
	if s.InitX != nil {
		// Initial x is provided.
		ctx.X.CloneVec(s.InitX)
		computeResidual(rInit, a, b, ctx.X, &stats)
	} else {
		// Initial x is the zero vector.
		ctx.X.Zero()
		// Residual b-A*x is then equal to b.
		rInit.CopyVec(b)
	}

	if m == nil {
		m = &GMRES{}
	}

	var err error
	ctx.ResidualNorm = mat.Norm(rInit, 2)
	if ctx.ResidualNorm >= s.Tolerance {
		err = iterate(a, b, rInit, s, m, &stats)
	} else {
		s.Dst.CopyVec(ctx.X)
	}

	return &Result{
		X:            s.Dst,
		ResidualNorm: ctx.ResidualNorm,
		Stats:        stats,
	}, err
}

func iterate(a MulVecToer, b, initRes *mat.VecDense, settings Settings, method Method, stats *Stats) error {
	bNorm := mat.Norm(b, 2)
	if bNorm == 0 {
		bNorm = 1
	}

	ctx := settings.Work
	settings.Dst.CopyVec(ctx.X)

	method.Init(ctx.X, initRes)
	for {
		op, err := method.Iterate(ctx)
		if err != nil {
			return err
		}
		switch op {
		case NoOperation:
		case MulVec, MulVec | Trans:
			stats.MulVec++
			a.MulVecTo(ctx.Dst, op&Trans == Trans, ctx.Src)
		case PreconSolve, PreconSolve | Trans:
			stats.PreconSolve++
			err = settings.PreconSolve(ctx.Dst, op&Trans == Trans, ctx.Src)
			if err != nil {
				return err
			}
		case CheckResidualNorm:
			ctx.Converged = ctx.ResidualNorm < settings.Tolerance*bNorm
		case ComputeResidual:
			computeResidual(ctx.Dst, a, b, ctx.X, stats)
		case MajorIteration:
			stats.Iterations++
			if ctx.Converged {
				settings.Dst.CopyVec(ctx.X)
				return nil
			}
			if stats.Iterations == settings.MaxIterations {
				settings.Dst.CopyVec(ctx.X)
				return ErrIterationLimit
			}
		default:
			panic("linsolve: invalid operation")
		}
	}
}

// NoPreconditioner implements the identity preconditioner.
func NoPreconditioner(dst *mat.VecDense, trans bool, rhs mat.Vector) error {
	if dst.Len() != rhs.Len() {
		panic("linsolve: mismatched vector length")
	}
	dst.CloneVec(rhs)
	return nil
}

func computeResidual(dst *mat.VecDense, a MulVecToer, b, x *mat.VecDense, stats *Stats) {
	stats.MulVec++
	a.MulVecTo(dst, false, x)
	dst.AddScaledVec(b, -1, dst)
}
