// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"testing"

	"github.com/gonum/gonum/matrix/mat64"
	"github.com/gonum/gonum/optimize/functions"
	"github.com/gonum/gonum/stat/distmv"
)

func TestGuessAndCheck(t *testing.T) {
	dim := 3000
	problem := Problem{
		Func: functions.ExtendedRosenbrock{}.Func,
	}
	mu := make([]float64, dim)
	sigma := mat64.NewSymDense(dim, nil)
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
