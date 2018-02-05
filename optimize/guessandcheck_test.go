// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"testing"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize/functions"
	"gonum.org/v1/gonum/stat/distmv"
)

func TestGuessAndCheck(t *testing.T) {
	dim := 3000
	problem := Problem{
		Func: functions.ExtendedRosenbrock{}.Func,
	}
	mu := make([]float64, dim)
	sigma := mat.NewSymDense(dim, nil)
	for i := 0; i < dim; i++ {
		sigma.SetSym(i, i, 1)
	}
	d, ok := distmv.NewNormal(mu, sigma, nil)
	if !ok {
		panic("bad test")
	}
	Global(problem, dim, nil, &GuessAndCheck{Rander: d})

	settings := DefaultSettingsGlobal()
	settings.Concurrent = 5
	settings.MajorIterations = 15
	Global(problem, dim, settings, &GuessAndCheck{Rander: d})
}
