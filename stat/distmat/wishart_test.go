// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmat

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/floats"
	"github.com/gonum/matrix/mat64"
)

func TestWishart(t *testing.T) {
	for c, test := range []struct {
		v   *mat64.SymDense
		nu  float64
		xs  []*mat64.SymDense
		lps []float64
	}{
		// Logprob data compared with scipy.
		{
			v:  mat64.NewSymDense(2, []float64{1, 0, 0, 1}),
			nu: 4,
			xs: []*mat64.SymDense{
				mat64.NewSymDense(2, []float64{0.9, 0.1, 0.1, 0.9}),
			},
			lps: []float64{-4.2357432031863409},
		},
		{
			v:  mat64.NewSymDense(2, []float64{0.8, -0.2, -0.2, 0.7}),
			nu: 5,
			xs: []*mat64.SymDense{
				mat64.NewSymDense(2, []float64{0.9, 0.1, 0.1, 0.9}),
				mat64.NewSymDense(2, []float64{0.3, -0.1, -0.1, 0.7}),
			},
			lps: []float64{-4.2476495605333575, -4.9993285370378633},
		},
		{
			v:  mat64.NewSymDense(3, []float64{0.8, 0.3, 0.1, 0.3, 0.7, -0.1, 0.1, -0.1, 7}),
			nu: 5,
			xs: []*mat64.SymDense{
				mat64.NewSymDense(3, []float64{1, 0.2, -0.3, 0.2, 0.6, -0.2, -0.3, -0.2, 6}),
			},
			lps: []float64{-11.010982249229421},
		},
	} {
		w, ok := NewWishart(test.v, test.nu, nil)
		if !ok {
			panic("bad test")
		}
		for i, x := range test.xs {
			lp := w.LogProbSym(x)

			var chol mat64.Cholesky
			ok := chol.Factorize(x)
			if !ok {
				panic("bad test")
			}
			lpc := w.LogProbSymChol(&chol)

			if math.Abs(lp-lpc) > 1e-14 {
				t.Errorf("Case %d, test %d: probability mismatch between chol and not", c, i)
			}
			if !floats.EqualWithinAbsOrRel(lp, test.lps[i], 1e-14, 1e-14) {
				t.Errorf("Case %d, test %d: got %v, want %v", c, i, lp, test.lps[i])
			}
		}

		ch := w.RandChol(nil)
		w.RandChol(ch)

		s := w.RandSym(nil)
		w.RandSym(s)

	}
}

func TestWishartRand(t *testing.T) {
	for c, test := range []struct {
		v       *mat64.SymDense
		nu      float64
		samples int
		tol     float64
	}{
		{
			v:       mat64.NewSymDense(2, []float64{0.8, -0.2, -0.2, 0.7}),
			nu:      5,
			samples: 30000,
			tol:     3e-2,
		},
		{
			v:       mat64.NewSymDense(3, []float64{0.8, 0.3, 0.1, 0.3, 0.7, -0.1, 0.1, -0.1, 7}),
			nu:      5,
			samples: 300000,
			tol:     3e-2,
		},
		{
			v: mat64.NewSymDense(4, []float64{
				0.8, 0.3, 0.1, -0.2,
				0.3, 0.7, -0.1, 0.4,
				0.1, -0.1, 7, 1,
				-0.2, -0.1, 1, 6}),
			nu:      6,
			samples: 300000,
			tol:     3e-2,
		},
	} {
		rnd := rand.New(rand.NewSource(1))
		dim := test.v.Symmetric()
		w, ok := NewWishart(test.v, test.nu, rnd)
		if !ok {
			panic("bad test")
		}
		mean := mat64.NewSymDense(dim, nil)
		x := mat64.NewSymDense(dim, nil)
		for i := 0; i < test.samples; i++ {
			w.RandSym(x)
			x.ScaleSym(1/float64(test.samples), x)
			mean.AddSym(mean, x)
		}
		trueMean := w.MeanSym(nil)
		if !mat64.EqualApprox(trueMean, mean, test.tol) {
			t.Errorf("Case %d: Mismatch between estimated and true mean. Got\n%0.4v\nWant\n%0.4v\n", c, mat64.Formatted(mean), mat64.Formatted(trueMean))
		}
	}
}
