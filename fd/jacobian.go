// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fd

import (
	"runtime"
	"sync"

	"github.com/gonum/floats"
	"github.com/gonum/matrix/mat64"
)

type JacobianSettings struct {
	Formula     Formula
	OriginValue []float64
	Step        float64
	Concurrent  bool
}

// Jacobian approximates the Jacobian matrix of a vector-valued function f at
// the location x.
//
// The Jacobian matrix J is the matrix of all first-order partial derivatives of f.
// If f maps an n-dimensional vector x to an m-dimensional vector y = f(x), J is
// an m×n matrix whose elements are given as
//  J_{i,j} = ∂f_i/∂x_j,
// or expanded out
//      [ ∂f_1/∂x_1 ... ∂f_1/∂x_n ]
//      [     .  .          .     ]
//  J = [     .      .      .     ]
//      [     .          .  .     ]
//      [ ∂f_m/∂x_1 ... ∂f_m/∂x_n ]
//
// If dst is not nil, the result will be stored in-place into dst and returned,
// otherwise a new matrix will be allocated first. Finite difference formula and
// other options are specified by settings. If settings is nil, the Jacobian
// will be estimated using the Forward formula and a default step size.
//
// Jacobian panics if dst is not nil and its size is not m × len(x), or if the
// derivative order of the formula is not 1.
func Jacobian(dst *mat64.Dense, f func(y, x []float64), m int, x []float64, settings *JacobianSettings) *mat64.Dense {
	n := len(x)
	if dst == nil {
		dst = mat64.NewDense(m, n, nil)
	}
	r, c := dst.Dims()
	if r != m || c != n {
		panic("jacobian: mismatched matrix size")
	}

	if settings == nil {
		settings = &JacobianSettings{}
	}
	if settings.OriginValue != nil && len(settings.OriginValue) != m {
		panic("jacobian: mismatched OriginValue slice length")
	}

	formula := settings.Formula
	if formula.isZero() {
		formula = Forward
	}
	if formula.Derivative == 0 || formula.Stencil == nil || formula.Step == 0 {
		panic("jacobian: bad formula")
	}
	if formula.Derivative != 1 {
		panic("jacobian: invalid derivative order")
	}

	step := settings.Step
	if step == 0 {
		step = formula.Step
	}

	var hasOrigin bool
	for _, pt := range formula.Stencil {
		if pt.Loc == 0 {
			hasOrigin = true
			break
		}
	}

	xcopy := make([]float64, n)
	origin := settings.OriginValue
	if hasOrigin && origin == nil {
		origin = make([]float64, m)
		copy(xcopy, x)
		f(origin, xcopy)
	}

	evals := n * len(formula.Stencil)
	if hasOrigin {
		evals -= n
	}
	nWorkers := 1
	if settings.Concurrent {
		nWorkers = runtime.GOMAXPROCS(0)
		if nWorkers > evals {
			nWorkers = evals
		}
	}

	if nWorkers == 1 {
		jacobianSerial(dst, f, x, xcopy, origin, formula, step)
	} else {
		jacobianConcurrent(dst, f, x, origin, formula, step, nWorkers)
	}
	return dst
}

func jacobianSerial(dst *mat64.Dense, f func([]float64, []float64), x, xcopy, origin []float64, formula Formula, step float64) {
	r, c := dst.Dims()
	y := make([]float64, r)
	col := make([]float64, r)
	for j := 0; j < c; j++ {
		for i := range col {
			col[i] = 0
		}
		for _, pt := range formula.Stencil {
			if pt.Loc == 0 {
				floats.AddScaled(col, pt.Coeff, origin)
			} else {
				copy(xcopy, x)
				xcopy[j] += pt.Loc * step
				f(y, xcopy)
				floats.AddScaled(col, pt.Coeff, y)
			}
		}
		dst.SetCol(j, col)
	}
	dst.Scale(1/step, dst)
}

func jacobianConcurrent(dst *mat64.Dense, f func([]float64, []float64), x, origin []float64, formula Formula, step float64, nWorkers int) {
	r, c := dst.Dims()
	originVec := mat64.NewVector(r, origin)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			dst.Set(i, j, 0)
		}
	}

	var wg sync.WaitGroup
	jobs := make(chan jacJob, nWorkers)
	mus := make([]sync.Mutex, c) // Guard access to individual columns.
	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			xcopy := make([]float64, len(x))
			y := make([]float64, r)
			yVec := mat64.NewVector(r, y)
			for job := range jobs {
				if job.pt.Loc != 0 {
					// Not an origin, we need to evaluate.
					copy(xcopy, x)
					xcopy[job.j] += job.pt.Loc * step
					f(y, xcopy)
				}

				col := dst.ColView(job.j)
				mus[job.j].Lock()
				if job.pt.Loc == 0 {
					col.AddScaledVec(col, job.pt.Coeff, originVec)
				} else {
					col.AddScaledVec(col, job.pt.Coeff, yVec)
				}
				mus[job.j].Unlock()
			}
		}()
	}
	for _, pt := range formula.Stencil {
		for j := 0; j < c; j++ {
			jobs <- jacJob{j, pt}
		}
	}
	close(jobs)
	wg.Wait()

	dst.Scale(1/step, dst)
}

type jacJob struct {
	j  int
	pt Point
}
