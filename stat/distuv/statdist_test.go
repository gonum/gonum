// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
)

func TestKullbackLeiblerBeta(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for cas, test := range []struct {
		a, b    Beta
		samples int
		tol     float64
	}{
		{
			a:       Beta{Alpha: 1, Beta: 2, Src: rnd},
			b:       Beta{Alpha: 1, Beta: 4, Src: rnd},
			samples: 100000,
			tol:     1e-2,
		},
		{
			a:       Beta{Alpha: 0.5, Beta: 0.4, Src: rnd},
			b:       Beta{Alpha: 0.7, Beta: 0.2, Src: rnd},
			samples: 100000,
			tol:     1e-2,
		},
		{
			a:       Beta{Alpha: 3, Beta: 5, Src: rnd},
			b:       Beta{Alpha: 5, Beta: 3, Src: rnd},
			samples: 100000,
			tol:     1e-2,
		},
	} {
		a, b := test.a, test.b
		want := klSample(test.samples, a, b)
		got := KullbackLeibler{}.DistBeta(a, b)
		if !floats.EqualWithinAbsOrRel(want, got, test.tol, test.tol) {
			t.Errorf("Kullback-Leibler mismatch, case %d: got %v, want %v", cas, got, want)
		}
	}
}

func TestKullbackLeiblerNormal(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for cas, test := range []struct {
		a, b    Normal
		samples int
		tol     float64
	}{
		{
			a:       Normal{Mu: 1, Sigma: 2, Src: rnd},
			b:       Normal{Mu: 1, Sigma: 4, Src: rnd},
			samples: 100000,
			tol:     1e-2,
		},
		{
			a:       Normal{Mu: 0, Sigma: 2, Src: rnd},
			b:       Normal{Mu: 2, Sigma: 2, Src: rnd},
			samples: 100000,
			tol:     1e-2,
		},
		{
			a:       Normal{Mu: 0, Sigma: 5, Src: rnd},
			b:       Normal{Mu: 2, Sigma: 0.1, Src: rnd},
			samples: 100000,
			tol:     1e-2,
		},
	} {
		a, b := test.a, test.b
		want := klSample(test.samples, a, b)
		got := KullbackLeibler{}.DistNormal(a, b)
		if !floats.EqualWithinAbsOrRel(want, got, test.tol, test.tol) {
			t.Errorf("Kullback-Leibler mismatch, case %d: got %v, want %v", cas, got, want)
		}
	}
}

// klSample finds an estimate of the Kullback-Leibler divergence through sampling.
func klSample(samples int, l RandLogProber, r LogProber) float64 {
	var klmc float64
	for i := 0; i < samples; i++ {
		x := l.Rand()
		pa := l.LogProb(x)
		pb := r.LogProb(x)
		klmc += pa - pb
	}
	return klmc / float64(samples)
}
