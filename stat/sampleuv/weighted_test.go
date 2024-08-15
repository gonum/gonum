// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file

package sampleuv

import (
	"flag"
	"reflect"
	"testing"
	"time"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/internal/rand"
)

var prob = flag.Bool("prob", false, "enables probabilistic testing of the random weighted sampler")

const sigChi2 = 16.92 // p = 0.05 df = 9

var (
	newExp = func() []float64 {
		return []float64{1 << 0, 1 << 1, 1 << 2, 1 << 3, 1 << 4, 1 << 5, 1 << 6, 1 << 7, 1 << 8, 1 << 9}
	}
	exp = newExp()

	obt = []float64{980, 1945, 3929, 7835, 15473, 31322, 62602, 124937, 250815, 500162}
)

func newTestWeighted() Weighted {
	weights := make([]float64, len(obt))
	for i := range weights {
		weights[i] = float64(int(1) << uint(i))
	}
	return NewWeighted(weights, nil)
}

func TestWeightedUnseeded(t *testing.T) {
	rand.Seed(0)

	want := Weighted{
		weights: []float64{1 << 0, 1 << 1, 1 << 2, 1 << 3, 1 << 4, 1 << 5, 1 << 6, 1 << 7, 1 << 8, 1 << 9},
		heap: []float64{
			exp[0] + exp[1] + exp[3] + exp[4] + exp[7] + exp[8] + exp[9] + exp[2] + exp[5] + exp[6],
			exp[1] + exp[3] + exp[4] + exp[7] + exp[8] + exp[9],
			exp[2] + exp[5] + exp[6],
			exp[3] + exp[7] + exp[8],
			exp[4] + exp[9],
			exp[5],
			exp[6],
			exp[7],
			exp[8],
			exp[9],
		},
	}

	ts := newTestWeighted()
	if !reflect.DeepEqual(ts, want) {
		t.Fatalf("unexpected new Weighted value:\ngot: %#v\nwant:%#v", ts, want)
	}

	f := make([]float64, len(obt))
	for i := 0; i < 1e6; i++ {
		item, ok := newTestWeighted().Take()
		if !ok {
			t.Fatal("Weighted unexpectedly empty")
		}
		f[item]++
	}

	exp := newExp()
	fac := floats.Sum(f) / floats.Sum(exp)
	for i := range f {
		exp[i] *= fac
	}

	if !reflect.DeepEqual(f, obt) {
		t.Fatalf("unexpected selection:\ngot: %#v\nwant:%#v", f, obt)
	}

	// Check that this is within statistical expectations - we know this is true for this set.
	X := chi2(f, exp)
	if X >= sigChi2 {
		t.Errorf("H₀: d(Sample) = d(Expect), H₁: d(S) ≠ d(Expect). df = %d, p = 0.05, X² threshold = %.2f, X² = %f", len(f)-1, sigChi2, X)
	}
}

func TestWeightedTimeSeeded(t *testing.T) {
	if !*prob {
		t.Skip("probabilistic testing not requested")
	}
	t.Log("Note: This test is stochastic and is expected to fail with probability ≈ 0.05.")

	rand.Seed(uint64(time.Now().Unix()))

	f := make([]float64, len(obt))
	for i := 0; i < 1e6; i++ {
		item, ok := newTestWeighted().Take()
		if !ok {
			t.Fatal("Weighted unexpectedly empty")
		}
		f[item]++
	}

	exp := newExp()
	fac := floats.Sum(f) / floats.Sum(exp)
	for i := range f {
		exp[i] *= fac
	}

	// Check that our obtained values are within statistical expectations for p = 0.05.
	// This will not be true approximately 1 in 20 tests.
	X := chi2(f, exp)
	if X >= sigChi2 {
		t.Errorf("H₀: d(Sample) = d(Expect), H₁: d(S) ≠ d(Expect). df = %d, p = 0.05, X² threshold = %.2f, X² = %f", len(f)-1, sigChi2, X)
	}
}

func TestWeightZero(t *testing.T) {
	rand.Seed(0)

	want := Weighted{
		weights: []float64{1 << 0, 1 << 1, 1 << 2, 1 << 3, 1 << 4, 1 << 5, 0, 1 << 7, 1 << 8, 1 << 9},
		heap: []float64{
			exp[0] + exp[1] + exp[3] + exp[4] + exp[7] + exp[8] + exp[9] + exp[2] + exp[5],
			exp[1] + exp[3] + exp[4] + exp[7] + exp[8] + exp[9],
			exp[2] + exp[5],
			exp[3] + exp[7] + exp[8],
			exp[4] + exp[9],
			exp[5],
			0,
			exp[7],
			exp[8],
			exp[9],
		},
	}

	ts := newTestWeighted()
	ts.Reweight(6, 0)
	if !reflect.DeepEqual(ts, want) {
		t.Fatalf("unexpected new Weighted value:\ngot: %#v\nwant:%#v", ts, want)
	}

	f := make([]float64, len(obt))
	for i := 0; i < 1e6; i++ {
		ts := newTestWeighted()
		ts.Reweight(6, 0)
		item, ok := ts.Take()
		if !ok {
			t.Fatal("Weighted unexpectedly empty")
		}
		f[item]++
	}

	exp := newExp()
	fac := floats.Sum(f) / floats.Sum(exp)
	for i := range f {
		exp[i] *= fac
	}

	if f[6] != 0 {
		t.Errorf("unexpected selection rate for zero-weighted item: got: %v want:%v", f[6], 0)
	}
	if reflect.DeepEqual(f[:6], obt[:6]) {
		t.Fatalf("unexpected selection: too few elements chosen in range:\ngot: %v\nwant:%v",
			f[:6], obt[:6])
	}
	if reflect.DeepEqual(f[7:], obt[7:]) {
		t.Fatalf("unexpected selection: too few elements chosen in range:\ngot: %v\nwant:%v",
			f[7:], obt[7:])
	}
}

func TestWeightIncrease(t *testing.T) {
	rand.Seed(0)

	want := Weighted{
		weights: []float64{1 << 0, 1 << 1, 1 << 2, 1 << 3, 1 << 4, 1 << 5, 1 << 9 * 2, 1 << 7, 1 << 8, 1 << 9},
		heap: []float64{
			exp[0] + exp[1] + exp[3] + exp[4] + exp[7] + exp[8] + exp[9] + exp[2] + exp[5] + exp[9]*2,
			exp[1] + exp[3] + exp[4] + exp[7] + exp[8] + exp[9],
			exp[2] + exp[5] + exp[9]*2,
			exp[3] + exp[7] + exp[8],
			exp[4] + exp[9],
			exp[5],
			exp[9] * 2,
			exp[7],
			exp[8],
			exp[9],
		},
	}

	ts := newTestWeighted()
	ts.Reweight(6, ts.weights[len(ts.weights)-1]*2)
	if !reflect.DeepEqual(ts, want) {
		t.Fatalf("unexpected new Weighted value:\ngot: %#v\nwant:%#v", ts, want)
	}

	f := make([]float64, len(obt))
	for i := 0; i < 1e6; i++ {
		ts := newTestWeighted()
		ts.Reweight(6, ts.weights[len(ts.weights)-1]*2)
		item, ok := ts.Take()
		if !ok {
			t.Fatal("Weighted unexpectedly empty")
		}
		f[item]++
	}

	exp := newExp()
	fac := floats.Sum(f) / floats.Sum(exp)
	for i := range f {
		exp[i] *= fac
	}

	if f[6] < f[9] {
		t.Errorf("unexpected selection rate for re-weighted item: got: %v want:%v", f[6], f[9])
	}
	if reflect.DeepEqual(f[:6], obt[:6]) {
		t.Fatalf("unexpected selection: too many elements chosen in range:\ngot: %v\nwant:%v",
			f[:6], obt[:6])
	}
	if reflect.DeepEqual(f[7:], obt[7:]) {
		t.Fatalf("unexpected selection: too many elements chosen in range:\ngot: %v\nwant:%v",
			f[7:], obt[7:])
	}
}

func chi2(ob, ex []float64) (sum float64) {
	for i := range ob {
		x := ob[i] - ex[i]
		sum += (x * x) / ex[i]
	}

	return sum
}

func TestWeightedNoResample(t *testing.T) {
	const (
		tries = 10
		n     = 10e4
	)
	ts := NewWeighted(make([]float64, n), nil)
	w := make([]float64, n)
	for i := 0; i < tries; i++ {
		for j := range w {
			w[j] = rand.Float64() * n
		}
		ts.ReweightAll(w)
		taken := make(map[int]struct{})
		var c int
		for {
			item, ok := ts.Take()
			if !ok {
				if c != n {
					t.Errorf("unexpected number of items: got: %d want: %d", c, int(n))
				}
				break
			}
			c++
			if _, exists := taken[item]; exists {
				t.Errorf("unexpected duplicate sample for item: %d", item)
			}
			taken[item] = struct{}{}
		}
	}
}

var issue1866Tests = []struct {
	weights []float64
	seed    uint64
}{
	{weights: []float64{72000, 166.66666666666666}, seed: 202},
	{weights: []float64{0.6666666666666666, 20412, 0}, seed: 0},
}

// See https://github.com/gonum/gonum/issues/1866
func TestIssue1866(t *testing.T) {
	for i, test := range issue1866Tests {
		w := NewWeighted(test.weights, rand.NewSource(test.seed))
		for j := 0; j < len(test.weights)+1; j++ {
			if panicked(func() { w.Take() }) {
				t.Errorf("unexpected panic for test %d iteration %d", i, j)
			}
		}
	}
}

func panicked(fn func()) (yes bool) {
	defer func() {
		if recover() != nil {
			yes = true
		}
	}()
	fn()
	return false
}

func FuzzWeightedTake(f *testing.F) {
	f.Add(0.000001, 0.000000012)

	f.Fuzz(func(t *testing.T, w1 float64, w2 float64) {
		if w1 < 0 || w2 < 0 {
			t.Skip()
		}

		weights := []float64{w1, w2}

		weighted := NewWeighted(weights, rand.NewSource(0))

		want := calculateNonZeroWeights(weights)
		var got int
		for {
			_, ok := weighted.Take()
			if !ok {
				break
			}
			got++
		}

		if got != want {
			t.Errorf("unexpected taken set size: got=%d, want=%d", got, want)
		}
	})
}

func calculateNonZeroWeights(weights []float64) int {
	count := 0
	for _, w := range weights {
		if w != 0 {
			count++
		}
	}

	return count
}
