// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/floats/scalar"
)

func TestBhattacharyyaBeta(t *testing.T) {
	t.Parallel()
	rnd := rand.New(rand.NewPCG(1, 1))
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
		want := bhattacharyyaSample(test.samples, test.a, test.b)
		got := Bhattacharyya{}.DistBeta(test.a, test.b)
		if !scalar.EqualWithinAbsOrRel(want, got, test.tol, test.tol) {
			t.Errorf("Bhattacharyya mismatch, case %d: got %v, want %v", cas, got, want)
		}

		// Bhattacharyya should be symmetric
		got2 := Bhattacharyya{}.DistBeta(test.b, test.a)
		if math.Abs(got-got2) > 1e-14 {
			t.Errorf("Bhattacharyya distance not symmetric")
		}
	}
}

func TestBhattacharyyaNormal(t *testing.T) {
	t.Parallel()
	rnd := rand.New(rand.NewPCG(1, 1))
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
			samples: 200000,
			tol:     1e-2,
		},
	} {
		want := bhattacharyyaSample(test.samples, test.a, test.b)
		got := Bhattacharyya{}.DistNormal(test.a, test.b)
		if !scalar.EqualWithinAbsOrRel(want, got, test.tol, test.tol) {
			t.Errorf("Bhattacharyya mismatch, case %d: got %v, want %v", cas, got, want)
		}

		// Bhattacharyya should be symmetric
		got2 := Bhattacharyya{}.DistNormal(test.b, test.a)
		if math.Abs(got-got2) > 1e-14 {
			t.Errorf("Bhattacharyya distance not symmetric")
		}
	}
}

// bhattacharyyaSample finds an estimate of the Bhattacharyya coefficient through
// sampling.
func bhattacharyyaSample(samples int, l RandLogProber, r LogProber) float64 {
	lBhatt := make([]float64, samples)
	for i := 0; i < samples; i++ {
		// Do importance sampling over a: \int sqrt(a*b)/a * a dx
		x := l.Rand()
		pa := l.LogProb(x)
		pb := r.LogProb(x)
		lBhatt[i] = 0.5*pb - 0.5*pa
	}
	logBc := floats.LogSumExp(lBhatt) - math.Log(float64(samples))
	return -logBc
}

func TestKullbackLeiblerBeta(t *testing.T) {
	t.Parallel()
	rnd := rand.New(rand.NewPCG(1, 1))
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
		if !scalar.EqualWithinAbsOrRel(want, got, test.tol, test.tol) {
			t.Errorf("Kullback-Leibler mismatch, case %d: got %v, want %v", cas, got, want)
		}
	}
	good := Beta{0.5, 0.5, nil}
	bad := Beta{0, 1, nil}
	if !panics(func() { KullbackLeibler{}.DistBeta(bad, good) }) {
		t.Errorf("Expected Kullback-Leibler to panic when called with invalid left Beta distribution")
	}
	if !panics(func() { KullbackLeibler{}.DistBeta(good, bad) }) {
		t.Errorf("Expected Kullback-Leibler to panic when called with invalid right Beta distribution")
	}
	bad = Beta{1, 0, nil}
	if !panics(func() { KullbackLeibler{}.DistBeta(bad, good) }) {
		t.Errorf("Expected Kullback-Leibler to panic when called with invalid left Beta distribution")
	}
	if !panics(func() { KullbackLeibler{}.DistBeta(good, bad) }) {
		t.Errorf("Expected Kullback-Leibler to panic when called with invalid right Beta distribution")
	}
}

func TestKullbackLeiblerNormal(t *testing.T) {
	t.Parallel()
	rnd := rand.New(rand.NewPCG(1, 1))
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
		if !scalar.EqualWithinAbsOrRel(want, got, test.tol, test.tol) {
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

func TestHellingerBeta(t *testing.T) {
	t.Parallel()
	rnd := rand.New(rand.NewPCG(1, 1))
	const tol = 1e-15
	for cas, test := range []struct {
		a, b Beta
	}{
		{
			a: Beta{Alpha: 1, Beta: 2, Src: rnd},
			b: Beta{Alpha: 1, Beta: 4, Src: rnd},
		},
		{
			a: Beta{Alpha: 0.5, Beta: 0.4, Src: rnd},
			b: Beta{Alpha: 0.7, Beta: 0.2, Src: rnd},
		},
		{
			a: Beta{Alpha: 3, Beta: 5, Src: rnd},
			b: Beta{Alpha: 5, Beta: 3, Src: rnd},
		},
	} {
		got := Hellinger{}.DistBeta(test.a, test.b)
		want := math.Sqrt(1 - math.Exp(-Bhattacharyya{}.DistBeta(test.a, test.b)))
		if !scalar.EqualWithinAbsOrRel(got, want, tol, tol) {
			t.Errorf("Hellinger mismatch, case %d: got %v, want %v", cas, got, want)
		}
	}
}

func TestHellingerNormal(t *testing.T) {
	t.Parallel()
	rnd := rand.New(rand.NewPCG(1, 1))
	const tol = 1e-15
	for cas, test := range []struct {
		a, b Normal
	}{
		{
			a: Normal{Mu: 1, Sigma: 2, Src: rnd},
			b: Normal{Mu: 1, Sigma: 4, Src: rnd},
		},
		{
			a: Normal{Mu: 0, Sigma: 2, Src: rnd},
			b: Normal{Mu: 2, Sigma: 2, Src: rnd},
		},
		{
			a: Normal{Mu: 0, Sigma: 5, Src: rnd},
			b: Normal{Mu: 2, Sigma: 0.1, Src: rnd},
		},
	} {
		got := Hellinger{}.DistNormal(test.a, test.b)
		want := math.Sqrt(1 - math.Exp(-Bhattacharyya{}.DistNormal(test.a, test.b)))
		if !scalar.EqualWithinAbsOrRel(got, want, tol, tol) {
			t.Errorf("Hellinger mismatch, case %d: got %v, want %v", cas, got, want)
		}
	}
}
