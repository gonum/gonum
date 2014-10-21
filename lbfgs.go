package optimize

import (
	"math"

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
// A LinesearchMethod for LBFGS must satisfy the strong Wolfe conditions at every
// iteration. If LinesearchMethod == nil, an appropriate default is chosen.
type LBFGS struct {
	LinesearchMethod LinesearchMethod
	Store            int // how many past iterations to store

	linesearch *Linesearch

	dim    int
	oldest int // element of the history slices that is the oldest

	x    []float64 // location at the last major iteration
	grad []float64 // gradient at the last major iteration

	y []float64 // holds g_{k+1} - g_k
	s []float64 // holds x_{k+1} - x_k
	a []float64 // holds cache of hessian updates

	// History
	yHist   [][]float64 // last Store iterations of y
	sHist   [][]float64 // last Store iterations of s
	rhoHist []float64   // last Store iterations of rho
}

func (l *LBFGS) Init(loc Location, f *FunctionInfo, xNext []float64) (EvaluationType, IterationType, error) {
	if l.LinesearchMethod == nil {
		l.LinesearchMethod = &Bisection{}
	}
	if l.linesearch == nil {
		l.linesearch = &Linesearch{}
	}
	l.linesearch.Method = l.LinesearchMethod
	l.linesearch.NextDirectioner = l
	return l.linesearch.Init(loc, f, xNext)
}

func (l *LBFGS) Iterate(loc Location, xNext []float64) (EvaluationType, IterationType, error) {
	return l.linesearch.Iterate(loc, xNext)
}

func (l *LBFGS) InitDirection(loc Location, direction []float64) (stepSize float64) {
	dim := len(loc.X)
	l.dim = dim

	if l.Store == 0 {
		l.Store = 15
	}

	l.oldest = l.Store - 1 // the first vector will be put in at 0

	l.x = resize(l.x, dim)
	l.grad = resize(l.grad, dim)
	copy(l.x, loc.X)
	copy(l.grad, loc.Gradient)

	l.y = resize(l.y, dim)
	l.s = resize(l.s, dim)
	l.a = resize(l.a, l.Store)
	l.rhoHist = resize(l.rhoHist, l.Store)

	if cap(l.yHist) < l.Store {
		n := make([][]float64, l.Store-cap(l.yHist))
		l.yHist = append(l.yHist, n...)
	}
	if cap(l.sHist) < l.Store {
		n := make([][]float64, l.Store-cap(l.sHist))
		l.sHist = append(l.sHist, n...)
	}
	l.yHist = l.yHist[:l.Store]
	l.sHist = l.sHist[:l.Store]
	for i := range l.sHist {
		l.sHist[i] = resize(l.sHist[i], dim)
		for j := range l.sHist[i] {
			l.sHist[i][j] = 0
		}
	}
	for i := range l.yHist {
		l.yHist[i] = resize(l.yHist[i], dim)
		for j := range l.yHist[i] {
			l.yHist[i][j] = 0
		}
	}

	copy(direction, loc.Gradient)
	floats.Scale(-1, direction)

	// see BFGS for comment
	floats.Scale(1/math.Sqrt(floats.Norm(direction, 2)), direction)
	return 1
}

func (l *LBFGS) NextDirection(loc Location, direction []float64) (stepSize float64) {
	if len(loc.X) != l.dim {
		panic("lbfgs: unexpected size mismatch")
	}
	if len(loc.Gradient) != l.dim {
		panic("lbfgs: unexpected size mismatch")
	}
	if len(direction) != l.dim {
		panic("lbfgs: unexpected size mismatch")
	}

	// Update direction. Uses two-loop correction as described in
	// Numerical Optimization. Nocedal and Wright, Ch 9, Page 225.
	floats.SubTo(l.y, loc.Gradient, l.grad)
	floats.SubTo(l.s, loc.X, l.x)
	rho := 1 / floats.Dot(l.y, l.s)
	copy(direction, loc.Gradient)

	// Compute gamma for the initial Hessian.
	gamma := floats.Dot(l.s, l.y) / floats.Dot(l.y, l.y)

	// two loop update. First loop starts with the most recent element
	// and goes backward, second starts with the oldest element and goes
	// forward. At the end have computed H^-1 * g, so flip the direction for
	// minimization.
	for i := 0; i < l.Store; i++ {
		idx := l.oldest - i - 1
		if idx < 0 {
			idx += l.Store
		}
		l.a[idx] = l.rhoHist[idx] * floats.Dot(l.sHist[idx], direction)
		floats.AddScaled(direction, -l.a[idx], l.yHist[idx])
	}
	floats.Scale(gamma, direction)
	for i := 0; i < l.Store; i++ {
		idx := i + l.oldest
		if idx >= l.Store {
			idx -= l.Store
		}
		beta := l.rhoHist[idx] * floats.Dot(l.yHist[idx], direction)
		floats.AddScaled(direction, l.a[idx]-beta, l.sHist[idx])
	}
	floats.Scale(-1, direction)

	// store the information from this iteration to prepare for the next one
	copy(l.x, loc.X)
	copy(l.grad, loc.Gradient)
	l.oldest++
	l.oldest = l.oldest % l.Store
	copy(l.sHist[l.oldest], l.s)
	copy(l.yHist[l.oldest], l.y)
	l.rhoHist[l.oldest] = rho
	return 1
}
