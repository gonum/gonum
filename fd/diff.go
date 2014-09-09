package diff

import (
	"math"
	"sync"
)

// A Point is a stencil location in a difference method
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
	Order   int // The order of the difference method (first derivative, second derivative, etc.)
}

// Settings is the settings structure for computing finite differences.
type Settings struct {
	OriginKnown bool    // Flag that the value at the origin x is known
	OriginValue float64 // Value at the origin (only used if OriginKnown is true)
	Step        float64 // step size
	Concurrent  bool    // Should the function calls be executed concurrently
	Method      Method  // Finite difference method to use
}

// DefaultSettings is a basic set of settings for computing finite differences.
// Computes a central difference approximation for the first derivative
// of the function.
func DefaultSettings() *Settings {
	return &Settings{
		Step:   1e-6,
		Method: Central,
	}
}

// Derivative estimates the derivative of the function f at the given location.
// The order of derivative, sample locations, and other options are specified
// by settings.
func Derivative(f func(float64) float64, x float64, settings *Settings) float64 {
	var deriv float64
	method := settings.Method
	if !settings.Concurrent {
		for _, pt := range method.Stencil {
			if settings.OriginKnown && pt.Loc == 0 {
				deriv += pt.Coeff * settings.OriginValue
				continue
			}
			deriv += pt.Coeff * f(x+settings.Step*pt.Loc)
		}
		return deriv / math.Pow(settings.Step, float64(method.Order))
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
			fofx := f(x + settings.Step*pt.Loc)
			mux.Lock()
			defer mux.Unlock()
			deriv += pt.Coeff * fofx

		}(pt)
	}
	wg.Wait()
	return deriv / math.Pow(settings.Step, float64(method.Order))
}

// Forward represents a first-order forward difference.
var Forward = Method{
	Stencil: []Point{{Loc: 0, Coeff: -1}, {Loc: 1, Coeff: 1}},
	Order:   1,
}

// Backward represents a first-order backward difference
var Backward = Method{
	Stencil: []Point{{Loc: -1, Coeff: -1}, {Loc: 0, Coeff: 1}},
	Order:   1,
}

// Central represents a first-order central difference.
var Central = Method{
	Stencil: []Point{{Loc: -1, Coeff: -0.5}, {Loc: 1, Coeff: 0.5}},
	Order:   1,
}

// Central2nd represents a secord-order central difference.
var Central2nd = Method{
	Stencil: []Point{{Loc: -1, Coeff: 1}, {Loc: 0, Coeff: -2}, {Loc: 1, Coeff: 1}},
	Order:   2,
}
