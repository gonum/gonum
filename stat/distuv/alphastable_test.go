// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"sort"
	"testing"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat"
)

func TestAlphaStable(t *testing.T) {
	t.Parallel()
	src := rand.New(rand.NewSource(1))
	for i, dist := range []AlphaStable{
		{Alpha: 0.5, Beta: 0, Scale: 1, Mu: 0, Src: src},
		{Alpha: 1, Beta: 0, Scale: 1, Mu: 0, Src: src},
		{Alpha: 2, Beta: 0, Scale: 1, Mu: 0, Src: src},
		{Alpha: 0.5, Beta: 1, Scale: 1, Mu: 0, Src: src},
		{Alpha: 1, Beta: 1, Scale: 1, Mu: 0, Src: src},
		{Alpha: 2, Beta: 1, Scale: 1, Mu: 0, Src: src},
		{Alpha: 0.5, Beta: 0, Scale: 1, Mu: 1, Src: src},
		{Alpha: 1, Beta: 0, Scale: 1, Mu: 1, Src: src},
		{Alpha: 2, Beta: 0, Scale: 1, Mu: 1, Src: src},
		{Alpha: 0.5, Beta: 0.5, Scale: 1, Mu: 1, Src: src},
		{Alpha: 1, Beta: 0.5, Scale: 1, Mu: 1, Src: src},
		{Alpha: 1.1, Beta: 0.5, Scale: 1, Mu: 1, Src: src},
		{Alpha: 2, Beta: 0.5, Scale: 1, Mu: 1, Src: src},
	} {
		testAlphaStableAnalytic(t, i, dist)
	}
}

func TestAlphaStability(t *testing.T) {
	t.Parallel()
	const (
		n     = 10000
		ksTol = 2e-2
	)
	for i, test := range []struct {
		alpha, beta1, beta2, scale1, scale2, mu1, mu2 float64
	}{
		{2, 0, 0, 1, 2, 0.5, 0.25},
		{2, 0.9, -0.4, 1, 2, 0.5, 0.25},
		{1.9, 0, 0, 1, 2, 0.5, 0.25},
		{1, 0, 0, 1, 2, 0.5, 0.25},
		{1, -0.5, 0.5, 1, 2, 0.5, 0.25},
		{0.5, 0, 0, 1, 2, 0.5, 0.25},
		{0.5, -0.5, 0.5, 1, 2, 0.5, 0.25},
	} {
		testStability(t, i, n, test.alpha, test.beta1, test.beta2, test.scale1, test.scale2, test.mu1, test.mu2, ksTol)
	}
}

func TestAlphaStableGaussian(t *testing.T) {
	t.Parallel()
	src := rand.New(rand.NewSource(1))
	d := AlphaStable{Alpha: 2, Beta: 0, Scale: 1.5, Mu: -0.4, Src: src}
	n := 100000
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = d.Rand()
	}
	checkSkewness(t, 0, x, d, 1e-2)
	checkExKurtosis(t, 0, x, d, 1e-2)
	checkMean(t, 0, x, d, 1e-2)
	checkVarAndStd(t, 0, x, d, 1e-2)
	checkMode(t, 0, x, d, 5e-2, 1e-1)
}

func TestAlphaStableMean(t *testing.T) {
	t.Parallel()
	src := rand.New(rand.NewSource(1))
	d := AlphaStable{Alpha: 1.75, Beta: 0.2, Scale: 1.2, Mu: 0.3, Src: src}
	n := 100000
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = d.Rand()
	}
	checkMean(t, 0, x, d, 1e-2)
}

func TestAlphaStableCauchy(t *testing.T) {
	t.Parallel()
	src := rand.New(rand.NewSource(1))
	d := AlphaStable{Alpha: 1, Beta: 0, Scale: 1, Mu: 0, Src: src}
	n := 1000000
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = d.Rand()
	}
	checkMode(t, 0, x, d, 1e-2, 5e-2)
}

func testAlphaStableAnalytic(t *testing.T, i int, dist AlphaStable) {
	if dist.NumParameters() != 4 {
		t.Errorf("%d: expected NumParameters == 4, got %v", i, dist.NumParameters())
	}
	if dist.Beta == 0 {
		if dist.Mode() != dist.Mu {
			t.Errorf("%d: expected Mode == Mu for Beta == 0, got %v", i, dist.Mode())
		}
		if dist.Median() != dist.Mu {
			t.Errorf("%d: expected Median == Mu for Beta == 0, got %v", i, dist.Median())
		}
	} else {
		if !panics(func() { dist.Mode() }) {
			t.Errorf("%d: expected Mode to panic for Beta != 0", i)
		}
		if !panics(func() { dist.Median() }) {
			t.Errorf("%d: expected Median to panic for Beta != 0", i)
		}
	}
	if dist.Alpha > 1 {
		if dist.Mean() != dist.Mu {
			t.Errorf("%d: expected Mean == Mu for Alpha > 1, got %v", i, dist.Mean())
		}
	} else {
		if !math.IsNaN(dist.Mean()) {
			t.Errorf("%d: expected NaN Mean for Alpha <= 1, got %v", i, dist.Mean())
		}
	}
	if dist.Alpha == 2 {
		got := dist.Variance()
		want := 2 * dist.Scale * dist.Scale
		if got != want {
			t.Errorf("%d: mismatch in Variance for Alpha == 2: got %v, want %g", i, got, want)
		}
		got = dist.StdDev()
		want = math.Sqrt(2) * dist.Scale
		if got != want {
			t.Errorf("%d: mismatch in StdDev for Alpha == 2: got %v, want %g", i, got, want)
		}
	}
}

func testStability(t *testing.T, i, n int, alpha, beta1, beta2, scale1, scale2, mu1, mu2, ksTol float64) {
	src := rand.New(rand.NewSource(1))
	d1 := AlphaStable{alpha, beta1, scale1, mu1, src}
	d2 := AlphaStable{alpha, beta2, scale2, mu2, src}
	scale := math.Pow(math.Pow(scale1, alpha)+math.Pow(scale2, alpha), 1/alpha)
	beta := (beta1*math.Pow(scale1, alpha) + beta2*math.Pow(scale2, alpha)) / math.Pow(scale, alpha)
	// Sum of d1 and d2.
	d := AlphaStable{alpha, beta, scale, mu1 + mu2, src}
	sample1 := make([]float64, n)
	sample2 := make([]float64, n)
	for i := 0; i < n; i++ {
		sample1[i] = d1.Rand() + d2.Rand()
		sample2[i] = d.Rand()
	}
	sort.Float64s(sample1)
	sort.Float64s(sample2)
	ks := stat.KolmogorovSmirnov(sample1, nil, sample2, nil)
	if ks > ksTol {
		t.Errorf("%d: Kolmogorov-Smirnov distance %g exceeding tolerance %g", i, ks, ksTol)
	}
}
