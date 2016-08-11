// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import "github.com/gonum/stat/distmv"

// GuessAndCheck is a global optimizer that evaluates the function at random
// locations. Not a good optimizer, but useful for comparison and debugging.
type GuessAndCheck struct {
	Rander distmv.Rander

	eval []bool
}

func (g *GuessAndCheck) Needs() struct{ Gradient, Hessian bool } {
	return struct{ Gradient, Hessian bool }{false, false}
}

func (g *GuessAndCheck) Done() {
	// No cleanup needed
}

func (g *GuessAndCheck) InitGlobal(tasks int) int {
	g.eval = make([]bool, tasks)
	return tasks
}

func (g *GuessAndCheck) IterateGlobal(task int, loc *Location) (Operation, error) {
	if g.eval[task] {
		g.eval[task] = false
		return MajorIteration, nil
	}
	g.eval[task] = true
	g.Rander.Rand(loc.X)
	return FuncEvaluation, nil
}
