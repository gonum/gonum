// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/matrix/mat64"
)

func TestDirichlet(t *testing.T) {
	// Data from Scipy.
	for cas, test := range []struct {
		Dir  *Dirichlet
		x    []float64
		prob float64
	}{
		{
			NewDirichlet([]float64{1, 1, 1}, nil),
			[]float64{0.2, 0.3, 0.5},
			2.0,
		},
		{
			NewDirichlet([]float64{0.6, 10, 8.7}, nil),
			[]float64{0.2, 0.3, 0.5},
			0.24079612737071665,
		},
	} {
		p := test.Dir.Prob(test.x)
		if math.Abs(p-test.prob) > 1e-14 {
			t.Errorf("Probablility mismatch. Case %v. Got %v, want %v", cas, p, test.prob)
		}
	}

	rnd := rand.New(rand.NewSource(1))
	for cas, test := range []struct {
		Dir *Dirichlet
		N   int
	}{
		{
			NewDirichlet([]float64{1, 1, 1}, rnd),
			1e6,
		},
		{
			NewDirichlet([]float64{2, 3}, rnd),
			1e6,
		},
		{
			NewDirichlet([]float64{0.2, 0.3}, rnd),
			1e6,
		},
		{
			NewDirichlet([]float64{0.2, 4}, rnd),
			1e6,
		},
		{
			NewDirichlet([]float64{0.1, 4, 20}, rnd),
			1e6,
		},
	} {
		d := test.Dir
		dim := d.Dim()
		x := mat64.NewDense(test.N, dim, nil)
		generateSamples(x, d)
		checkMean(t, cas, x, d, 1e-3)
		checkCov(t, cas, x, d, 1e-3)
	}
}
