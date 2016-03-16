// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fd provides functions to approximate derivatives using finite differences.
package fd

import (
	"math"
	"runtime"
	"sync"

	"github.com/gonum/floats"
)

// A Point is a stencil location in a finite difference formula.
type Point struct {
	Loc   float64
	Coeff float64
}

// Formula represents a finite difference formula that approximates
// the derivative of order k of a function f at x as
//  d^k f(x) ≅ (1 / h^k) * \sum_i Coeff_i * f(x + h * Loc_i),
// where h is a small positive step.
type Formula struct {
	// Stencil is the set of sampling Points which are used to estimate the
	// derivative. The locations will be scaled by Step and are relative to x.
	Stencil []Point
	Order   int     // The order of the approximated derivative.
	Step    float64 // Default step size for the formula.
}

// Settings is the settings structure for computing finite differences.
type Settings struct {
	OriginKnown bool    // Flag that the value at the origin x is known
	OriginValue float64 // Value at the origin (only used if OriginKnown is true)
	Concurrent  bool    // Should the function calls be executed concurrently.
	Formula     Formula // Finite difference formula to use
}

// DefaultSettings is a basic set of settings for computing finite differences.
// Computes a central difference approximation for the first derivative
// of the function.
func DefaultSettings() *Settings {
	return &Settings{
		Formula: Central,
	}
}

// Derivative estimates the derivative of the function f at the given location.
// The order of the derivative, sample locations, and other options are
// specified by settings. If settings is nil, default settings will be used.
func Derivative(f func(float64) float64, x float64, settings *Settings) float64 {
	if settings == nil {
		settings = DefaultSettings()
	}
	step := settings.Formula.Step
	var deriv float64
	formula := settings.Formula
	if !settings.Concurrent {
		for _, pt := range formula.Stencil {
			if settings.OriginKnown && pt.Loc == 0 {
				deriv += pt.Coeff * settings.OriginValue
				continue
			}
			deriv += pt.Coeff * f(x+step*pt.Loc)
		}
		return deriv / math.Pow(step, float64(formula.Order))
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
	return deriv / math.Pow(step, float64(formula.Order))
}

// Gradient estimates the gradient of the multivariate function f at the
// location x. The result is stored in-place into dst if dst is not nil,
// otherwise a new slice will be allocated and returned. Finite difference
// kernel and other options are specified by settings. If settings is nil,
// default settings will be used.
// Gradient panics if the length of dst and x is not equal.
func Gradient(dst []float64, f func([]float64) float64, x []float64, settings *Settings) []float64 {
	if dst == nil {
		dst = make([]float64, len(x))
	}
	if len(dst) != len(x) {
		panic("fd: slice length mismatch")
	}
	if settings == nil {
		settings = DefaultSettings()
	}
	step := settings.Formula.Step
	if !settings.Concurrent {
		xcopy := make([]float64, len(x)) // So that x is not modified during the call
		copy(xcopy, x)
		for i := range xcopy {
			var deriv float64
			for _, pt := range settings.Formula.Stencil {
				if settings.OriginKnown && pt.Loc == 0 {
					deriv += pt.Coeff * settings.OriginValue
					continue
				}
				xcopy[i] += pt.Loc * step
				deriv += pt.Coeff * f(xcopy)
				xcopy[i] = x[i]
			}
			dst[i] = deriv / math.Pow(step, float64(settings.Formula.Order))
		}
		return dst
	}

	quit := make(chan struct{})
	defer close(quit)

	expect := len(settings.Formula.Stencil) * len(x)
	sendChan := make(chan fdrun, expect)
	ansChan := make(chan fdrun, expect)

	// Launch workers. Workers receive an index and a step, and compute the answer
	nWorkers := runtime.NumCPU()
	if nWorkers > expect {
		nWorkers = expect
	}
	for i := 0; i < nWorkers; i++ {
		go func(sendChan <-chan fdrun, ansChan chan<- fdrun, quit <-chan struct{}) {
			xcopy := make([]float64, len(x))
			copy(xcopy, x)
			for {
				select {
				case <-quit:
					return
				case run := <-sendChan:
					xcopy[run.idx] += run.pt.Loc * step
					run.result = f(xcopy)
					xcopy[run.idx] = x[run.idx]
					ansChan <- run
				}
			}
		}(sendChan, ansChan, quit)
	}

	// Launch the distributor. Distributor sends the cases to be computed
	go func(sendChan chan<- fdrun, ansChan chan<- fdrun) {
		for i := range x {
			for _, pt := range settings.Formula.Stencil {
				if settings.OriginKnown && pt.Loc == 0 {
					// Answer already known. Send the answer on the answer channel
					ansChan <- fdrun{
						idx:    i,
						pt:     pt,
						result: settings.OriginValue,
					}
					continue
				}
				// Answer not known, send the answer to be computed
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
	// Read in all of the results
	for i := 0; i < expect; i++ {
		run := <-ansChan
		dst[run.idx] += run.pt.Coeff * run.result
	}
	floats.Scale(1/math.Pow(step, float64(settings.Formula.Order)), dst)
	return dst
}

type fdrun struct {
	idx    int
	pt     Point
	result float64
}

// Forward represents a first-order forward difference.
var Forward = Formula{
	Stencil: []Point{{Loc: 0, Coeff: -1}, {Loc: 1, Coeff: 1}},
	Order:   1,
	Step:    1e-6,
}

// Backward represents a first-order backward difference.
var Backward = Formula{
	Stencil: []Point{{Loc: -1, Coeff: -1}, {Loc: 0, Coeff: 1}},
	Order:   1,
	Step:    1e-6,
}

// Central represents a first-order central difference.
var Central = Formula{
	Stencil: []Point{{Loc: -1, Coeff: -0.5}, {Loc: 1, Coeff: 0.5}},
	Order:   1,
	Step:    1e-6,
}

// Central2nd represents a secord-order central difference.
var Central2nd = Formula{
	Stencil: []Point{{Loc: -1, Coeff: 1}, {Loc: 0, Coeff: -2}, {Loc: 1, Coeff: 1}},
	Order:   2,
	Step:    1e-3,
}
