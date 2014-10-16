// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fd

import (
	"math"
	"runtime"
	"sync"

	"github.com/gonum/floats"
)

// A Point is a stencil location in a difference method.
type Point struct {
	Loc   float64
	Coeff float64
}

// Method is a specific finite difference method. Method specifies the stencil,
// that is, the function locations (relative to x) which will be used to estimate
// the derivative. It also specifies the order of derivative it estimates. Order = 1
// represents the derivative, Order = 2 represents the curvature, etc.
type Method struct {
	Stencil []Point
	Order   int     // The order of the difference method (first derivative, second derivative, etc.)
	Step    float64 // Default step size for the method.
}

// Settings is the settings structure for computing finite differences.
type Settings struct {
	OriginKnown bool    // Flag that the value at the origin x is known
	OriginValue float64 // Value at the origin (only used if OriginKnown is true)
	Concurrent  bool    // Should the function calls be executed concurrently
	Workers     int     // Maximum number of concurrent executions when evaluating concurrently
	Method      Method  // Finite difference method to use
}

// DefaultSettings is a basic set of settings for computing finite differences.
// Computes a central difference approximation for the first derivative
// of the function.
func DefaultSettings() *Settings {
	return &Settings{
		Method:  Central,
		Workers: runtime.GOMAXPROCS(0),
	}
}

// Derivative estimates the derivative of the function f at the given location.
// The order of derivative, sample locations, and other options are specified
// by settings.
func Derivative(f func(float64) float64, x float64, settings *Settings) float64 {
	if settings == nil {
		settings = DefaultSettings()
	}
	step := settings.Method.Step
	var deriv float64
	method := settings.Method
	if !settings.Concurrent {
		for _, pt := range method.Stencil {
			if settings.OriginKnown && pt.Loc == 0 {
				deriv += pt.Coeff * settings.OriginValue
				continue
			}
			deriv += pt.Coeff * f(x+step*pt.Loc)
		}
		return deriv / math.Pow(step, float64(method.Order))
	}

	wg := &sync.WaitGroup{}
	mux := &sync.Mutex{}
	for _, pt := range method.Stencil {
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
	return deriv / math.Pow(step, float64(method.Order))
}

// Gradient estimates the derivative of a multivariate function f at the location
// x. The resulting estimate is stored in-place into gradient. The order of derivative,
// sample locations, and other options are specified by settings.
// If the step size is zero, then the step size of the method will
// be used.
// Gradient panics if len(deriv) != len(x).
func Gradient(f func([]float64) float64, x []float64, settings *Settings, gradient []float64) {
	if len(gradient) != len(x) {
		panic("fd: location and gradient length mismatch")
	}
	if settings == nil {
		settings = DefaultSettings()
	}
	step := settings.Method.Step
	if !settings.Concurrent {
		xcopy := make([]float64, len(x)) // So that x is not modified during the call
		copy(xcopy, x)
		for i := range xcopy {
			var deriv float64
			for _, pt := range settings.Method.Stencil {
				if settings.OriginKnown && pt.Loc == 0 {
					deriv += pt.Coeff * settings.OriginValue
					continue
				}
				xcopy[i] += pt.Loc * step
				deriv += pt.Coeff * f(xcopy)
				xcopy[i] = x[i]
			}
			gradient[i] = deriv / math.Pow(step, float64(settings.Method.Order))
		}
		return
	}

	nWorkers := settings.Workers
	expect := len(settings.Method.Stencil) * len(x)
	if nWorkers > expect {
		nWorkers = expect
	}

	quit := make(chan struct{})
	defer close(quit)
	sendChan := make(chan fdrun, expect)
	ansChan := make(chan fdrun, expect)

	// Launch workers. Workers receive an index and a step, and compute the answer
	for i := 0; i < settings.Workers; i++ {
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
			for _, pt := range settings.Method.Stencil {
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

	for i := range gradient {
		gradient[i] = 0
	}
	// Read in all of the results
	for i := 0; i < expect; i++ {
		run := <-ansChan
		gradient[run.idx] += run.pt.Coeff * run.result
	}
	floats.Scale(1/math.Pow(step, float64(settings.Method.Order)), gradient)
}

type fdrun struct {
	idx    int
	pt     Point
	result float64
}

// Forward represents a first-order forward difference.
var Forward = Method{
	Stencil: []Point{{Loc: 0, Coeff: -1}, {Loc: 1, Coeff: 1}},
	Order:   1,
	Step:    1e-6,
}

// Backward represents a first-order backward difference.
var Backward = Method{
	Stencil: []Point{{Loc: -1, Coeff: -1}, {Loc: 0, Coeff: 1}},
	Order:   1,
	Step:    1e-6,
}

// Central represents a first-order central difference.
var Central = Method{
	Stencil: []Point{{Loc: -1, Coeff: -0.5}, {Loc: 1, Coeff: 0.5}},
	Order:   1,
	Step:    1e-6,
}

// Central2nd represents a secord-order central difference.
var Central2nd = Method{
	Stencil: []Point{{Loc: -1, Coeff: 1}, {Loc: 0, Coeff: -2}, {Loc: 1, Coeff: 1}},
	Order:   2,
	Step:    1e-3,
}
