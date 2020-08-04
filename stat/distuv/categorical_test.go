// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/floats/scalar"
)

const (
	Tiny   = 2
	Small  = 5
	Medium = 10
	Large  = 100
	Huge   = 1000
)

func TestCategoricalProb(t *testing.T) {
	t.Parallel()
	for _, test := range [][]float64{
		{1, 2, 3, 0},
	} {
		dist := NewCategorical(test, nil)
		norm := make([]float64, len(test))
		floats.Scale(1/floats.Sum(norm), norm)
		for i, v := range norm {
			p := dist.Prob(float64(i))
			if math.Abs(p-v) > 1e-14 {
				t.Errorf("Probability mismatch element %d", i)
			}
			logP := dist.LogProb(float64(i))
			if math.Abs(logP-math.Log(v)) > 1e-14 {
				t.Errorf("Log-probability mismatch element %d", i)
			}
			p = dist.Prob(float64(i) + 0.5)
			if p != 0 {
				t.Errorf("Non-zero probability for non-integer x")
			}
			logP = dist.LogProb(float64(i) + 0.5)
			if !math.IsInf(logP, -1) {
				t.Errorf("Log-probability for non-integer x is not -Inf")
			}
		}
		p := dist.Prob(-1)
		if p != 0 {
			t.Errorf("Non-zero probability for -1")
		}
		logP := dist.LogProb(-1)
		if !math.IsInf(logP, -1) {
			t.Errorf("Log-probability for -1 is not -Inf")
		}
		p = dist.Prob(float64(len(test)))
		if p != 0 {
			t.Errorf("Non-zero probability for len(test)")
		}
		logP = dist.LogProb(float64(len(test)))
		if !math.IsInf(logP, -1) {
			t.Errorf("Log-probability for len(test) is not -Inf")
		}
	}
}

func TestCategoricalRand(t *testing.T) {
	t.Parallel()
	for _, test := range [][]float64{
		{1, 2, 3, 0},
	} {
		dist := NewCategorical(test, nil)
		nSamples := 2000000
		counts := sampleCategorical(t, dist, nSamples)

		probs := make([]float64, len(test))
		for i := range probs {
			probs[i] = dist.Prob(float64(i))
		}
		same := samedDistCategorical(dist, counts, probs, 1e-2)
		if !same {
			t.Errorf("Probability mismatch. Want %v, got %v", probs, counts)
		}

		dist.Reweight(len(test)-1, 10)
		counts = sampleCategorical(t, dist, nSamples)
		probs = make([]float64, len(test))
		for i := range probs {
			probs[i] = dist.Prob(float64(i))
		}
		same = samedDistCategorical(dist, counts, probs, 1e-2)
		if !same {
			t.Errorf("Probability mismatch after Reweight. Want %v, got %v", probs, counts)
		}

		w := make([]float64, len(test))
		for i := range w {
			w[i] = rand.Float64()
		}

		dist.ReweightAll(w)
		counts = sampleCategorical(t, dist, nSamples)
		probs = make([]float64, len(test))
		for i := range probs {
			probs[i] = dist.Prob(float64(i))
		}
		same = samedDistCategorical(dist, counts, probs, 1e-2)
		if !same {
			t.Errorf("Probability mismatch after ReweightAll. Want %v, got %v", probs, counts)
		}
	}
}

func TestCategoricalReweight(t *testing.T) {
	t.Parallel()
	dist := NewCategorical([]float64{1, 1}, nil)
	if !panics(func() { dist.Reweight(0, -1) }) {
		t.Errorf("Reweight did not panic for negative weight")
	}
	dist.Reweight(0, 0)
	if !panics(func() { dist.Reweight(1, 0) }) {
		t.Errorf("Reweight did not panic when trying to set the last positive weight to zero")
	}
}

func TestCategoricalReweightAll(t *testing.T) {
	t.Parallel()
	w := []float64{0, 1, 2, 1}
	dist := NewCategorical(w, nil)
	if !panics(func() { dist.ReweightAll([]float64{1, 1}) }) {
		t.Errorf("ReweightAll did not panic for different number of weights")
	}
	w[0] = -1
	if !panics(func() { dist.ReweightAll(w) }) {
		t.Errorf("ReweightAll did not panic for a negative weight")
	}
	w = []float64{0, 0, 0, 0}
	if !panics(func() { dist.ReweightAll(w) }) {
		t.Errorf("ReweightAll did not panic for weights which are all zero")
	}
}

func sampleCategorical(t *testing.T, dist Categorical, nSamples int) []float64 {
	counts := make([]float64, dist.Len())
	for i := 0; i < nSamples; i++ {
		v := dist.Rand()
		if float64(int(v)) != v {
			t.Fatalf("Random number is not an integer")
		}
		counts[int(v)]++
	}
	sum := floats.Sum(counts)
	floats.Scale(1/sum, counts)
	return counts
}

func samedDistCategorical(dist Categorical, counts, probs []float64, tol float64) bool {
	same := true
	for i, prob := range probs {
		if prob == 0 && counts[i] != 0 {
			same = false
			break
		}
		if !scalar.EqualWithinAbsOrRel(prob, counts[i], tol, tol) {
			same = false
			break
		}
	}
	return same
}

func TestCategoricalCDF(t *testing.T) {
	t.Parallel()
	for _, test := range [][]float64{
		{1, 2, 3, 0, 4},
	} {
		c := make([]float64, len(test))
		copy(c, test)
		floats.Scale(1/floats.Sum(c), c)
		sum := make([]float64, len(test))
		floats.CumSum(sum, c)

		dist := NewCategorical(test, nil)
		cdf := dist.CDF(-0.5)
		if cdf != 0 {
			t.Errorf("CDF of negative number not zero")
		}
		for i := range c {
			cdf := dist.CDF(float64(i))
			if math.Abs(cdf-sum[i]) > 1e-14 {
				t.Errorf("CDF mismatch %v. Want %v, got %v.", float64(i), sum[i], cdf)
			}
			cdfp := dist.CDF(float64(i) + 0.5)
			if cdfp != cdf {
				t.Errorf("CDF mismatch for non-integer input")
			}
		}
	}
}

func TestCategoricalEntropy(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		weights []float64
		entropy float64
	}{
		{
			weights: []float64{1, 1},
			entropy: math.Ln2,
		},
		{
			weights: []float64{1, 1, 1, 1},
			entropy: math.Log(4),
		},
		{
			weights: []float64{0, 0, 1, 1, 0, 0},
			entropy: math.Ln2,
		},
	} {
		dist := NewCategorical(test.weights, nil)
		entropy := dist.Entropy()
		if math.IsNaN(entropy) || math.Abs(entropy-test.entropy) > 1e-14 {
			t.Errorf("Entropy mismatch. Want %v, got %v.", test.entropy, entropy)
		}
	}
}

func TestCategoricalMean(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		weights []float64
		mean    float64
	}{
		{
			weights: []float64{10, 0, 0, 0},
			mean:    0,
		},
		{
			weights: []float64{0, 10, 0, 0},
			mean:    1,
		},
		{
			weights: []float64{1, 2, 3, 4},
			mean:    2,
		},
	} {
		dist := NewCategorical(test.weights, nil)
		mean := dist.Mean()
		if math.IsNaN(mean) || math.Abs(mean-test.mean) > 1e-14 {
			t.Errorf("Entropy mismatch. Want %v, got %v.", test.mean, mean)
		}
	}
}

func BenchmarkCategoricalRandTiny(b *testing.B)   { benchmarkCategoricalRand(b, Tiny) }
func BenchmarkCategoricalRandSmall(b *testing.B)  { benchmarkCategoricalRand(b, Small) }
func BenchmarkCategoricalRandMedium(b *testing.B) { benchmarkCategoricalRand(b, Medium) }
func BenchmarkCategoricalRandLarge(b *testing.B)  { benchmarkCategoricalRand(b, Large) }
func BenchmarkCategoricalRandHuge(b *testing.B)   { benchmarkCategoricalRand(b, Huge) }

func benchmarkCategoricalRand(b *testing.B, size int) {
	src := rand.NewSource(1)
	rng := rand.New(src)
	weights := make([]float64, size)
	for i := 0; i < size; i++ {
		weights[i] = rng.Float64() + 0.001
	}
	dist := NewCategorical(weights, src)
	for i := 0; i < b.N; i++ {
		dist.Rand()
	}
}
