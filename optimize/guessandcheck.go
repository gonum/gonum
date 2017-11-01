// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"
	"sync"

	"gonum.org/v1/gonum/stat/distmv"
)

// GuessAndCheck is a global optimizer that evaluates the function at random
// locations. Not a good optimizer, but useful for comparison and debugging.
type GuessAndCheck struct {
	Rander distmv.Rander

	eval []bool

	mux   *sync.Mutex
	bestF float64
	bestX []float64
}

func (g *GuessAndCheck) Needs() struct{ Gradient, Hessian bool } {
	return struct{ Gradient, Hessian bool }{false, false}
}

func (g *GuessAndCheck) Done() {
	// No cleanup needed
}

func (g *GuessAndCheck) InitGlobal(dim, tasks int) int {
	g.eval = make([]bool, tasks)
	g.bestF = math.Inf(1)
	g.bestX = resize(g.bestX, dim)
	g.mux = &sync.Mutex{}
	return tasks
}

func (g *GuessAndCheck) IterateGlobal(task int, loc *Location) (Operation, error) {
	// Task is true if it contains a new function evaluation.
	if g.eval[task] {
		g.eval[task] = false
		g.mux.Lock()
		if loc.F < g.bestF {
			g.bestF = loc.F
			copy(g.bestX, loc.X)
		} else {
			loc.F = g.bestF
			copy(loc.X, g.bestX)
		}
		g.mux.Unlock()
		return MajorIteration, nil
	}
	g.eval[task] = true
	g.Rander.Rand(loc.X)
	return FuncEvaluation, nil
}
