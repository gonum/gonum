// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scalar

import (
	"math"
)

// Status represents the status of the optimization.
type Status int

const (
	// NotTerminated indicates that the optimization has not yet finished.
	NotTerminated Status = iota
	// Converged indicates that the optimization converged successfully.
	Converged
	// IterationLimit indicates that the maximum number of iterations was reached.
	IterationLimit
	// Failure indicates that the optimization failed.
	Failure
)

func (s Status) String() string {
	switch s {
	case NotTerminated:
		return "NotTerminated"
	case Converged:
		return "Converged"
	case IterationLimit:
		return "IterationLimit"
	case Failure:
		return "Failure"
	default:
		return "Unknown"
	}
}

// Settings allows fine-grained control over the optimization process.
type Settings struct {
	// Tol is the tolerance for the stopping criterion.
	// The algorithm stops when the interval width falls below:
	//   tol_act = Tol*|x| + epsilon
	// If Tol <= 0, a default value based on square root of machine epsilon is used.
	Tol float64

	// MaxIterations is the maximum number of iterations allowed.
	// If MaxIterations <= 0, a default safety limit of 100 is used.
	MaxIterations int
}

// Result holds the result of the minimization.
type Result struct {
	X          float64 // The input value that minimizes the function
	F          float64 // The minimum function value found
	Iterations int     // Number of iterations performed
	Status     Status  // The status of the result
}

// BrentMin finds the local minimum of the scalar function f in the interval [min, max].
// It uses Brent's method, which combines golden section search and successive parabolic interpolation.
// This implementation is a port of Netlib's fmin.f.
//
// See: http://www.netlib.org/opt/fmin.f
// Reference: Brent, Richard P. Algorithms for minimization without derivatives. Courier Corporation, 2013.
func BrentMin(f func(float64) float64, min, max float64, settings *Settings) (Result, error) {
	const (
		// Machine epsilon for float64: 2**-52
		epsilon = 0x1p-52
		// Sqrt of machine epsilon: 2**-26
		sqrtEpsilon = 0x1p-26

		// c is the squared inverse of the golden ratio: (3 - sqrt(5))/2
		c = 0.3819660112501051517954131656343618822796908201942371
	)

	// Default settings
	tol := sqrtEpsilon
	maxIter := 100

	if settings != nil {
		if settings.Tol > 0 {
			tol = settings.Tol
		}
		if settings.MaxIterations > 0 {
			maxIter = settings.MaxIterations
		}
	}

	a := min
	b := max
	x := a + c*(b-a)
	w := x
	v := x
	var (
		e float64 // step size of the step before last
		d float64 // step size of the last step
	)

	fx := f(x)
	fw := fx
	fv := fx

	for i:= 1; ; i++ {

		if iterations > maxIter {
			return Result{
				X:          x,
				F:          fx,
				Iterations: iterations,
				Status:     IterationLimit,
			}, nil
		}

		xm := (a + b) / 2
		tol1 := epsilon*math.Abs(x) + tol/3
		tol2 := 2 * tol1

		// Check stopping criterion
		if math.Abs(x-xm) <= (tol2 - (b-a)/2) {
			return Result{
				X:          x,
				F:          fx,
				Iterations: iterations,
				Status:     Converged,
			}, nil
		}

		// Is golden-section necessary?
		if math.Abs(e) > tol1 {
			r := (x - w) * (fx - fv)
			q := (x - v) * (fx - fw)
			p := (x-v)*q - (x-w)*r
			q = 2.0 * (q - r)
			if q > 0 {
				p = -p
			}
			q = math.Abs(q)
			etemp := e
			e = d

			// Check acceptability of the parabolic fit
			if math.Abs(p) >= math.Abs(0.5*q*etemp) || p <= q*(a-x) || p >= q*(b-x) {
				// Parabolic interpolation step rejected, use golden-section
				if x >= xm {
					e = a - x
				} else {
					e = b - x
				}
				d = c * e
			} else {
				// Parabolic interpolation step accepted
				d = p / q
				u := x + d
				// Evaluate f must not be close to a or b
				if u-a < tol2 || b-u < tol2 {
					if xm-x >= 0 {
						d = tol1
					} else {
						d = -tol1
					}
				}
			}
		} else {
			// Golden-section step
			if x >= xm {
				e = a - x
			} else {
				e = b - x
			}
			d = c * e
		}

		// f must not be evaluated too close to x
		u := 0.0
		if math.Abs(d) >= tol1 {
			u = x + d
		} else {
			if d > 0 {
				u = x + tol1
			} else {
				u = x - tol1
			}
		}

		fu := f(u)

		// Update a, b, v, w, and x
		if fu <= fx {
			if u >= x {
				a = x
			} else {
				b = x
			}
			v = w
			fv = fw
			w = x
			fw = fx
			x = u
			fx = fu
		} else {
			if u < x {
				a = u
			} else {
				b = u
			}
			if fu <= fw || w == x {
				v = w
				fv = fw
				w = u
				fw = fu
			} else if fu <= fv || v == x || v == w {
				v = u
				fv = fu
			}
		}
	}
}
