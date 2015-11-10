// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"github.com/gonum/floats"
)

// LBFGS implements the limited-memory BFGS algorithm. While the normal BFGS algorithm
// makes a full approximation to the inverse hessian, LBFGS instead approximates the
// hessian from the last Store optimization steps. The Store parameter is a tradeoff
// between cost of the method and accuracy of the hessian approximation.
// LBFGS has a cost (both in memory and time) of O(Store * inputDimension).
// Since BFGS has a cost of O(inputDimension^2), LBFGS is more appropriate
// for very large problems. This "forgetful" nature of LBFGS may also make it perform
// better than BFGS for functions with Hessians that vary rapidly spatially.
//
// If Store is 0, Store is defaulted to 15.
// A Linesearcher for LBFGS must satisfy the strong Wolfe conditions at every
// iteration. If Linesearcher == nil, an appropriate default is chosen.
type LBFGS struct {
	Linesearcher Linesearcher
	Store        int // how many past iterations to store

	ls *LinesearchMethod

	dim  int       // Dimension of the problem
	x    []float64 // Location at the last major iteration
	grad []float64 // Gradient at the last major iteration

	// History
	oldest int         // Index of the oldest element of the history
	y      [][]float64 // Last Store values of y
	s      [][]float64 // Last Store values of s
	rho    []float64   // Last Store values of rho
	a      []float64   // Cache of Hessian updates
}

func (l *LBFGS) Init(loc *Location) (Operation, error) {
	if l.Linesearcher == nil {
		l.Linesearcher = &Bisection{}
	}
	if l.Store == 0 {
		l.Store = 15
	}

	if l.ls == nil {
		l.ls = &LinesearchMethod{}
	}
	l.ls.Linesearcher = l.Linesearcher
	l.ls.NextDirectioner = l

	return l.ls.Init(loc)
}

func (l *LBFGS) Iterate(loc *Location) (Operation, error) {
	return l.ls.Iterate(loc)
}

func (l *LBFGS) InitDirection(loc *Location, dir []float64) (stepSize float64) {
	dim := len(loc.X)
	l.dim = dim
	l.oldest = 0

	l.a = resize(l.a, l.Store)
	l.rho = resize(l.rho, l.Store)
	l.y = l.initHistory(l.y)
	l.s = l.initHistory(l.s)

	l.x = resize(l.x, dim)
	copy(l.x, loc.X)

	l.grad = resize(l.grad, dim)
	copy(l.grad, loc.Gradient)

	copy(dir, loc.Gradient)
	floats.Scale(-1, dir)
	return 1 / floats.Norm(dir, 2)
}

func (l *LBFGS) initHistory(hist [][]float64) [][]float64 {
	c := cap(hist)
	if c < l.Store {
		n := make([][]float64, l.Store-c)
		hist = append(hist[:c], n...)
	}
	hist = hist[:l.Store]
	for i := range hist {
		hist[i] = resize(hist[i], l.dim)
		for j := range hist[i] {
			hist[i][j] = 0
		}
	}
	return hist
}

func (l *LBFGS) NextDirection(loc *Location, dir []float64) (stepSize float64) {
	// Uses two-loop correction as described in
	// Nocedal, J., Wright, S.: Numerical Optimization (2nd ed). Springer (2006), chapter 7, page 178.

	if len(loc.X) != l.dim {
		panic("lbfgs: unexpected size mismatch")
	}
	if len(loc.Gradient) != l.dim {
		panic("lbfgs: unexpected size mismatch")
	}
	if len(dir) != l.dim {
		panic("lbfgs: unexpected size mismatch")
	}

	y := l.y[l.oldest]
	floats.SubTo(y, loc.Gradient, l.grad)
	s := l.s[l.oldest]
	floats.SubTo(s, loc.X, l.x)
	sDotY := floats.Dot(s, y)
	l.rho[l.oldest] = 1 / sDotY

	l.oldest = (l.oldest + 1) % l.Store

	copy(l.x, loc.X)
	copy(l.grad, loc.Gradient)
	copy(dir, loc.Gradient)

	// Start with the most recent element and go backward,
	for i := 0; i < l.Store; i++ {
		idx := l.oldest - i - 1
		if idx < 0 {
			idx += l.Store
		}
		l.a[idx] = l.rho[idx] * floats.Dot(l.s[idx], dir)
		floats.AddScaled(dir, -l.a[idx], l.y[idx])
	}

	// Scale the initial Hessian.
	gamma := sDotY / floats.Dot(y, y)
	floats.Scale(gamma, dir)

	// Start with the oldest element and go forward.
	for i := 0; i < l.Store; i++ {
		idx := i + l.oldest
		if idx >= l.Store {
			idx -= l.Store
		}
		beta := l.rho[idx] * floats.Dot(l.y[idx], dir)
		floats.AddScaled(dir, l.a[idx]-beta, l.s[idx])
	}

	// dir contains H^{-1} * g, so flip the direction for minimization.
	floats.Scale(-1, dir)

	return 1
}

func (*LBFGS) Needs() struct {
	Gradient bool
	Hessian  bool
} {
	return struct {
		Gradient bool
		Hessian  bool
	}{true, false}
}
