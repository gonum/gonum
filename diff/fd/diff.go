// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fd provides functions to approximate derivatives using finite differences.
package fd // import "gonum.org/v1/gonum/diff/fd"

import (
	"math"
	"runtime"
	"sync"

	"gonum.org/v1/gonum/floats"
)

// A Point is a stencil location in a finite difference formula.
type Point struct {
	Loc   float64
	Coeff float64
}

// Formula represents a finite difference formula on a regularly spaced grid
// that approximates the derivative of order k of a function f at x as
//  d^k f(x) ≈ (1 / Step^k) * \sum_i Coeff_i * f(x + Step * Loc_i).
type Formula struct {
	// Stencil is the set of sampling Points which are used to estimate the
	// derivative. The locations will be scaled by Step and are relative to x.
	Stencil    []Point
	Derivative int     // The order of the approximated derivative.
	Step       float64 // Default step size for the formula.
}

func (f Formula) isZero() bool {
	return f.Stencil == nil && f.Derivative == 0 && f.Step == 0
}

// Settings is the settings structure for computing finite differences.
type Settings struct {
	// Formula is the finite difference formula used
	// for approximating the derivative.
	// Zero value indicates a default formula.
	Formula Formula
	// Step is the distance between points of the stencil.
	// If equal to 0, formula's default step will be used.
	Step float64

	OriginKnown bool    // Flag that the value at the origin x is known.
	OriginValue float64 // Value at the origin (only used if OriginKnown is true).

	Concurrent bool // Should the function calls be executed concurrently.
}

// Derivative estimates the derivative of the function f at the given location.
// The finite difference formula, the step size, and other options are
// specified by settings. If settings is nil, the first derivative will be
// estimated using the Forward formula and a default step size.
func Derivative(f func(float64) float64, x float64, settings *Settings) float64 {
	if settings == nil {
		settings = &Settings{}
	}
	formula := settings.Formula
	if formula.isZero() {
		formula = Forward
	}
	if formula.Derivative == 0 || formula.Stencil == nil || formula.Step == 0 {
		panic("fd: bad formula")
	}
	step := settings.Step
	if step == 0 {
		step = formula.Step
	}

	var deriv float64
	if !settings.Concurrent || runtime.GOMAXPROCS(0) == 1 {
		for _, pt := range formula.Stencil {
			if settings.OriginKnown && pt.Loc == 0 {
				deriv += pt.Coeff * settings.OriginValue
				continue
			}
			deriv += pt.Coeff * f(x+step*pt.Loc)
		}
		return deriv / math.Pow(step, float64(formula.Derivative))
	}

	wg := &sync.WaitGroup{}
	mux := &sync.Mutex{}
	for _, pt := range formula.Stencil {
		if settings.OriginKnown && pt.Loc == 0 {
			mux.Lock()
			deriv += pt.Coeff * settings.OriginValue
			mux.Unlock()
			continue
		}
		wg.Add(1)
		go func(pt Point) {
			defer wg.Done()
			fofx := f(x + step*pt.Loc)
			mux.Lock()
			defer mux.Unlock()
			deriv += pt.Coeff * fofx
		}(pt)
	}
	wg.Wait()
	return deriv / math.Pow(step, float64(formula.Derivative))
}

// Gradient estimates the gradient of the multivariate function f at the
// location x. If dst is not nil, the result will be stored in-place into dst
// and returned, otherwise a new slice will be allocated first. Finite
// difference kernel and other options are specified by settings. If settings is
// nil, the gradient will be estimated using the Forward formula and a default
// step size.
//
// Gradient panics if the length of dst and x is not equal, or if the derivative
// order of the formula is not 1.
func Gradient(dst []float64, f func([]float64) float64, x []float64, settings *Settings) []float64 {
	if dst == nil {
		dst = make([]float64, len(x))
	}
	if len(dst) != len(x) {
		panic("fd: slice length mismatch")
	}
	if settings == nil {
		settings = &Settings{}
	}

	formula := settings.Formula
	if formula.isZero() {
		formula = Forward
	}
	if formula.Derivative == 0 || formula.Stencil == nil || formula.Step == 0 {
		panic("fd: bad formula")
	}
	if formula.Derivative != 1 {
		panic("fd: invalid derivative order")
	}

	step := settings.Step
	if step == 0 {
		step = formula.Step
	}

	expect := len(formula.Stencil) * len(x)
	nWorkers := 1
	if settings.Concurrent {
		nWorkers = runtime.GOMAXPROCS(0)
		if nWorkers > expect {
			nWorkers = expect
		}
	}

	var hasOrigin bool
	for _, pt := range formula.Stencil {
		if pt.Loc == 0 {
			hasOrigin = true
			break
		}
	}
	xcopy := make([]float64, len(x)) // So that x is not modified during the call.
	originValue := settings.OriginValue
	if hasOrigin && !settings.OriginKnown {
		copy(xcopy, x)
		originValue = f(xcopy)
	}

	if nWorkers == 1 {
		for i := range xcopy {
			var deriv float64
			for _, pt := range formula.Stencil {
				if pt.Loc == 0 {
					deriv += pt.Coeff * originValue
					continue
				}
				copy(xcopy, x)
				xcopy[i] += pt.Loc * step
				deriv += pt.Coeff * f(xcopy)
			}
			dst[i] = deriv / step
		}
		return dst
	}

	sendChan := make(chan fdrun, expect)
	ansChan := make(chan fdrun, expect)
	quit := make(chan struct{})
	defer close(quit)

	// Launch workers. Workers receive an index and a step, and compute the answer.
	for i := 0; i < nWorkers; i++ {
		go func(sendChan <-chan fdrun, ansChan chan<- fdrun, quit <-chan struct{}) {
			xcopy := make([]float64, len(x))
			for {
				select {
				case <-quit:
					return
				case run := <-sendChan:
					copy(xcopy, x)
					xcopy[run.idx] += run.pt.Loc * step
					run.result = f(xcopy)
					ansChan <- run
				}
			}
		}(sendChan, ansChan, quit)
	}

	// Launch the distributor. Distributor sends the cases to be computed.
	go func(sendChan chan<- fdrun, ansChan chan<- fdrun) {
		for i := range x {
			for _, pt := range formula.Stencil {
				if pt.Loc == 0 {
					// Answer already known. Send the answer on the answer channel.
					ansChan <- fdrun{
						idx:    i,
						pt:     pt,
						result: originValue,
					}
					continue
				}
				// Answer not known, send the answer to be computed.
				sendChan <- fdrun{
					idx: i,
					pt:  pt,
				}
			}
		}
	}(sendChan, ansChan)

	for i := range dst {
		dst[i] = 0
	}
	// Read in all of the results.
	for i := 0; i < expect; i++ {
		run := <-ansChan
		dst[run.idx] += run.pt.Coeff * run.result
	}
	floats.Scale(1/step, dst)
	return dst
}

type fdrun struct {
	idx    int
	pt     Point
	result float64
}

// Forward represents a first-order accurate forward approximation
// to the first derivative.
var Forward = Formula{
	Stencil:    []Point{{Loc: 0, Coeff: -1}, {Loc: 1, Coeff: 1}},
	Derivative: 1,
	Step:       2e-8,
}

// Backward represents a first-order accurate backward approximation
// to the first derivative.
var Backward = Formula{
	Stencil:    []Point{{Loc: -1, Coeff: -1}, {Loc: 0, Coeff: 1}},
	Derivative: 1,
	Step:       2e-8,
}

// Central represents a second-order accurate centered approximation
// to the first derivative.
var Central = Formula{
	Stencil:    []Point{{Loc: -1, Coeff: -0.5}, {Loc: 1, Coeff: 0.5}},
	Derivative: 1,
	Step:       6e-6,
}

// Central2nd represents a secord-order accurate centered approximation
// to the second derivative.
var Central2nd = Formula{
	Stencil:    []Point{{Loc: -1, Coeff: 1}, {Loc: 0, Coeff: -2}, {Loc: 1, Coeff: 1}},
	Derivative: 2,
	Step:       1e-4,
}
