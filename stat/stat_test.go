// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"math"
	"reflect"
	"strconv"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
)

func TestCircularMean(t *testing.T) {
	for i, test := range []struct {
		x   []float64
		wts []float64
		ans float64
	}{
		// Values compared against scipy.
		{
			x:   []float64{0, 2 * math.Pi},
			ans: 0,
		},
		{
			x:   []float64{0, 0.5 * math.Pi},
			ans: 0.78539816339744,
		},
		{
			x:   []float64{-1.5 * math.Pi, 0.5 * math.Pi, 2.5 * math.Pi},
			wts: []float64{1, 2, 3},
			ans: 0.5 * math.Pi,
		},
		{
			x:   []float64{0, 0.5 * math.Pi},
			wts: []float64{1, 2},
			ans: 1.10714871779409,
		},
	} {
		c := CircularMean(test.x, test.wts)
		if math.Abs(c-test.ans) > 1e-14 {
			t.Errorf("Circular mean mismatch case %d: Expected %v, Found %v", i, test.ans, c)
		}
	}
	if !panics(func() { CircularMean(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("CircularMean did not panic with x, wts length mismatch")
	}
}

func TestCorrelation(t *testing.T) {
	for i, test := range []struct {
		x   []float64
		y   []float64
		w   []float64
		ans float64
	}{
		{
			x:   []float64{8, -3, 7, 8, -4},
			y:   []float64{8, -3, 7, 8, -4},
			w:   nil,
			ans: 1,
		},
		{
			x:   []float64{8, -3, 7, 8, -4},
			y:   []float64{8, -3, 7, 8, -4},
			w:   []float64{1, 1, 1, 1, 1},
			ans: 1,
		},
		{
			x:   []float64{8, -3, 7, 8, -4},
			y:   []float64{8, -3, 7, 8, -4},
			w:   []float64{1, 6, 7, 0.8, 2.1},
			ans: 1,
		},
		{
			x:   []float64{8, -3, 7, 8, -4},
			y:   []float64{10, 15, 4, 5, -1},
			w:   nil,
			ans: 0.0093334660769059,
		},
		{
			x:   []float64{8, -3, 7, 8, -4},
			y:   []float64{10, 15, 4, 5, -1},
			w:   nil,
			ans: 0.0093334660769059,
		},
		{
			x:   []float64{8, -3, 7, 8, -4},
			y:   []float64{10, 15, 4, 5, -1},
			w:   []float64{1, 3, 1, 2, 2},
			ans: -0.13966633352689,
		},
	} {
		c := Correlation(test.x, test.y, test.w)
		if math.Abs(test.ans-c) > 1e-14 {
			t.Errorf("Correlation mismatch case %d. Expected %v, Found %v", i, test.ans, c)
		}
	}
	if !panics(func() { Correlation(make([]float64, 2), make([]float64, 3), make([]float64, 3)) }) {
		t.Errorf("Correlation did not panic with length mismatch")
	}
	if !panics(func() { Correlation(make([]float64, 2), make([]float64, 3), nil) }) {
		t.Errorf("Correlation did not panic with length mismatch")
	}
	if !panics(func() { Correlation(make([]float64, 3), make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("Correlation did not panic with weights length mismatch")
	}
}

func TestKendall(t *testing.T) {
	for i, test := range []struct {
		x       []float64
		y       []float64
		weights []float64
		ans     float64
	}{
		{
			x:       []float64{0, 1, 2, 3},
			y:       []float64{0, 1, 2, 3},
			weights: nil,
			ans:     1,
		},
		{
			x:       []float64{0, 1},
			y:       []float64{1, 0},
			weights: nil,
			ans:     -1,
		},
		{
			x:       []float64{8, -3, 7, 8, -4},
			y:       []float64{10, 15, 4, 5, -1},
			weights: nil,
			ans:     0.2,
		},
		{
			x:       []float64{8, -3, 7, 8, -4},
			y:       []float64{10, 5, 6, 3, -1},
			weights: nil,
			ans:     0.4,
		},
		{
			x:       []float64{1, 2, 3, 4, 5},
			y:       []float64{2, 3, 4, 5, 6},
			weights: []float64{1, 1, 1, 1, 1},
			ans:     1,
		},
		{
			x:       []float64{1, 2, 3, 2, 1},
			y:       []float64{2, 3, 2, 1, 0},
			weights: []float64{1, 1, 0, 0, 0},
			ans:     1,
		},
	} {
		c := Kendall(test.x, test.y, test.weights)
		if math.Abs(test.ans-c) > 1e-14 {
			t.Errorf("Correlation mismatch case %d. Expected %v, Found %v", i, test.ans, c)
		}
	}
	if !panics(func() { Kendall(make([]float64, 2), make([]float64, 3), make([]float64, 3)) }) {
		t.Errorf("Kendall did not panic with length mismatch")
	}
	if !panics(func() { Kendall(make([]float64, 2), make([]float64, 3), nil) }) {
		t.Errorf("Kendall did not panic with length mismatch")
	}
	if !panics(func() { Kendall(make([]float64, 3), make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("Kendall did not panic with weights length mismatch")
	}
}

func TestCovariance(t *testing.T) {
	for i, test := range []struct {
		p       []float64
		q       []float64
		weights []float64
		ans     float64
	}{
		{
			p:   []float64{0.75, 0.1, 0.05},
			q:   []float64{0.5, 0.25, 0.25},
			ans: 0.05625,
		},
		{
			p:   []float64{1, 2, 3},
			q:   []float64{2, 4, 6},
			ans: 2,
		},
		{
			p:   []float64{1, 2, 3},
			q:   []float64{1, 4, 9},
			ans: 4,
		},
		{
			p:       []float64{1, 2, 3},
			q:       []float64{1, 4, 9},
			weights: []float64{1, 1.5, 1},
			ans:     3.2,
		},
		{
			p:       []float64{1, 4, 9},
			q:       []float64{1, 4, 9},
			weights: []float64{1, 1.5, 1},
			ans:     13.142857142857146,
		},
	} {
		c := Covariance(test.p, test.q, test.weights)
		if math.Abs(c-test.ans) > 1e-14 {
			t.Errorf("Covariance mismatch case %d: Expected %v, Found %v", i, test.ans, c)
		}
	}

	// test the panic states
	if !panics(func() { Covariance(make([]float64, 2), make([]float64, 3), nil) }) {
		t.Errorf("Covariance did not panic with x, y length mismatch")
	}
	if !panics(func() { Covariance(make([]float64, 3), make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("Covariance did not panic with x, weights length mismatch")
	}

}

func TestCrossEntropy(t *testing.T) {
	for i, test := range []struct {
		p   []float64
		q   []float64
		ans float64
	}{
		{
			p:   []float64{0.75, 0.1, 0.05},
			q:   []float64{0.5, 0.25, 0.25},
			ans: 0.7278045395879426,
		},
		{
			p:   []float64{0.75, 0.1, 0.05, 0, 0, 0},
			q:   []float64{0.5, 0.25, 0.25, 0, 0, 0},
			ans: 0.7278045395879426,
		},
		{
			p:   []float64{0.75, 0.1, 0.05, 0, 0, 0.1},
			q:   []float64{0.5, 0.25, 0.25, 0, 0, 0},
			ans: math.Inf(1),
		},
		{
			p:   nil,
			q:   nil,
			ans: 0,
		},
	} {
		c := CrossEntropy(test.p, test.q)
		if math.Abs(c-test.ans) > 1e-14 {
			t.Errorf("Cross entropy mismatch case %d: Expected %v, Found %v", i, test.ans, c)
		}
	}
	if !panics(func() { CrossEntropy(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("CrossEntropy did not panic with p, q length mismatch")
	}
}

func TestExKurtosis(t *testing.T) {
	// the example does a good job, this just has to cover the panic
	if !panics(func() { ExKurtosis(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("ExKurtosis did not panic with x, weights length mismatch")
	}
}

func TestGeometricMean(t *testing.T) {
	for i, test := range []struct {
		x   []float64
		wts []float64
		ans float64
	}{
		{
			x:   []float64{2, 8},
			ans: 4,
		},
		{
			x:   []float64{3, 81},
			wts: []float64{2, 1},
			ans: 9,
		},
	} {
		c := GeometricMean(test.x, test.wts)
		if math.Abs(c-test.ans) > 1e-14 {
			t.Errorf("Geometric mean mismatch case %d: Expected %v, Found %v", i, test.ans, c)
		}
	}
	if !panics(func() { GeometricMean(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("GeometricMean did not panic with x, wts length mismatch")
	}
}

func TestHarmonicMean(t *testing.T) {
	for i, test := range []struct {
		x   []float64
		wts []float64
		ans float64
	}{
		{
			x:   []float64{.5, .125},
			ans: .2,
		},
		{
			x:   []float64{.5, .125},
			wts: []float64{2, 1},
			ans: .25,
		},
	} {
		c := HarmonicMean(test.x, test.wts)
		if math.Abs(c-test.ans) > 1e-14 {
			t.Errorf("Harmonic mean mismatch case %d: Expected %v, Found %v", i, test.ans, c)
		}
	}
	if !panics(func() { HarmonicMean(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("HarmonicMean did not panic with x, wts length mismatch")
	}
}

func TestHistogram(t *testing.T) {
	for i, test := range []struct {
		x        []float64
		weights  []float64
		dividers []float64
		ans      []float64
	}{
		{
			x:        []float64{1, 3, 5, 6, 7, 8},
			dividers: []float64{0, 2, 4, 6, 7, 9},
			ans:      []float64{1, 1, 1, 1, 2},
		},
		{
			x:        []float64{1, 3, 5, 6, 7, 8},
			dividers: []float64{1, 2, 4, 6, 7, 9},
			weights:  []float64{1, 2, 1, 1, 1, 2},
			ans:      []float64{1, 2, 1, 1, 3},
		},
		{
			x:        []float64{1, 8},
			dividers: []float64{0, 2, 4, 6, 7, 9},
			weights:  []float64{1, 2},
			ans:      []float64{1, 0, 0, 0, 2},
		},
		{
			x:        []float64{1, 8},
			dividers: []float64{0, 2, 4, 6, 7, 9},
			ans:      []float64{1, 0, 0, 0, 1},
		},
		{
			x:        []float64{},
			dividers: []float64{1, 3},
			ans:      []float64{0},
		},
	} {
		hist := Histogram(nil, test.dividers, test.x, test.weights)
		if !floats.Equal(hist, test.ans) {
			t.Errorf("Hist mismatch case %d. Expected %v, Found %v", i, test.ans, hist)
		}
		// Test with non-zero values
		Histogram(hist, test.dividers, test.x, test.weights)
		if !floats.Equal(hist, test.ans) {
			t.Errorf("Hist mismatch case %d. Expected %v, Found %v", i, test.ans, hist)
		}
	}
	// panic cases
	for _, test := range []struct {
		name     string
		x        []float64
		weights  []float64
		dividers []float64
		count    []float64
	}{
		{
			name:    "len(x) != len(weights)",
			x:       []float64{1, 3, 5, 6, 7, 8},
			weights: []float64{1, 1, 1, 1},
		},
		{
			name:     "len(count) != len(dividers) - 1",
			x:        []float64{1, 3, 5, 6, 7, 8},
			dividers: []float64{1, 4, 9},
			count:    make([]float64, 6),
		},
		{
			name:     "dividers not sorted",
			x:        []float64{1, 3, 5, 6, 7, 8},
			dividers: []float64{0, -1, 0},
		},
		{
			name:     "x not sorted",
			x:        []float64{1, 5, 2, 9, 7, 8},
			dividers: []float64{1, 4, 9},
		},
		{
			name:     "fewer than 2 dividers",
			x:        []float64{1, 2, 3},
			dividers: []float64{5},
		},
		{
			name:     "x too large",
			x:        []float64{1, 2, 3},
			dividers: []float64{1, 3},
		},
		{
			name:     "x too small",
			x:        []float64{1, 2, 3},
			dividers: []float64{2, 3},
		},
	} {
		if !panics(func() { Histogram(test.count, test.dividers, test.x, test.weights) }) {
			t.Errorf("Histogram did not panic when %s", test.name)
		}
	}
}

func TestJensenShannon(t *testing.T) {
	for i, test := range []struct {
		p []float64
		q []float64
	}{
		{
			p: []float64{0.5, 0.1, 0.3, 0.1},
			q: []float64{0.1, 0.4, 0.25, 0.25},
		},
		{
			p: []float64{0.4, 0.6, 0.0},
			q: []float64{0.2, 0.2, 0.6},
		},
		{
			p: []float64{0.1, 0.1, 0.0, 0.8},
			q: []float64{0.6, 0.3, 0.0, 0.1},
		},
		{
			p: []float64{0.5, 0.1, 0.3, 0.1},
			q: []float64{0.5, 0, 0.25, 0.25},
		},
		{
			p: []float64{0.5, 0.1, 0, 0.4},
			q: []float64{0.1, 0.4, 0.25, 0.25},
		},
	} {

		m := make([]float64, len(test.p))
		p := test.p
		q := test.q
		floats.Add(m, p)
		floats.Add(m, q)
		floats.Scale(0.5, m)

		js1 := 0.5*KullbackLeibler(p, m) + 0.5*KullbackLeibler(q, m)
		js2 := JensenShannon(p, q)

		if math.IsNaN(js2) {
			t.Errorf("In case %v, JS distance is NaN", i)
		}

		if math.Abs(js1-js2) > 1e-14 {
			t.Errorf("JS mismatch case %v. Expected %v, found %v.", i, js1, js2)
		}
	}
	if !panics(func() { JensenShannon(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("JensenShannon did not panic with p, q length mismatch")
	}
}

func TestKolmogorovSmirnov(t *testing.T) {
	for i, test := range []struct {
		x        []float64
		xWeights []float64
		y        []float64
		yWeights []float64
		dist     float64
	}{

		{
			dist: 0,
		},
		{
			x:    []float64{1},
			dist: 1,
		},
		{
			y:    []float64{1},
			dist: 1,
		},
		{
			x:        []float64{1},
			xWeights: []float64{8},
			dist:     1,
		},
		{
			y:        []float64{1},
			yWeights: []float64{8},
			dist:     1,
		},
		{
			x:        []float64{1},
			xWeights: []float64{8},
			y:        []float64{1},
			yWeights: []float64{8},
			dist:     0,
		},
		{
			x:        []float64{1, 1, 1},
			xWeights: []float64{2, 3, 7},
			y:        []float64{1},
			yWeights: []float64{8},
			dist:     0,
		},
		{
			x:        []float64{1, 1, 1, 1, 1},
			y:        []float64{1, 1, 1},
			yWeights: []float64{2, 5, 2},
			dist:     0,
		},

		{
			x:    []float64{1, 2, 3},
			y:    []float64{1, 2, 3},
			dist: 0,
		},
		{
			x:        []float64{1, 2, 3},
			y:        []float64{1, 2, 3},
			yWeights: []float64{1, 1, 1},
			dist:     0,
		},

		{
			x:        []float64{1, 2, 3},
			xWeights: []float64{1, 1, 1},
			y:        []float64{1, 2, 3},
			yWeights: []float64{1, 1, 1},
			dist:     0,
		},
		{
			x:        []float64{1, 2},
			xWeights: []float64{2, 5},
			y:        []float64{1, 1, 2, 2, 2, 2, 2},
			dist:     0,
		},
		{
			x:        []float64{1, 1, 2, 2, 2, 2, 2},
			y:        []float64{1, 2},
			yWeights: []float64{2, 5},
			dist:     0,
		},
		{
			x:        []float64{1, 1, 2, 2, 2},
			xWeights: []float64{0.5, 1.5, 1, 2, 2},
			y:        []float64{1, 2},
			yWeights: []float64{2, 5},
			dist:     0,
		},
		{
			x:    []float64{1, 2, 3, 4},
			y:    []float64{5, 6},
			dist: 1,
		},
		{
			x:    []float64{5, 6},
			y:    []float64{1, 2, 3, 4},
			dist: 1,
		},
		{
			x:        []float64{5, 6},
			xWeights: []float64{8, 7},
			y:        []float64{1, 2, 3, 4},
			dist:     1,
		},
		{
			x:        []float64{5, 6},
			xWeights: []float64{8, 7},
			y:        []float64{1, 2, 3, 4},
			yWeights: []float64{9, 2, 1, 6},
			dist:     1,
		},
		{
			x:        []float64{-4, 5, 6},
			xWeights: []float64{0, 8, 7},
			y:        []float64{1, 2, 3, 4},
			yWeights: []float64{9, 2, 1, 6},
			dist:     1,
		},
		{
			x:        []float64{-4, -2, -2, 5, 6},
			xWeights: []float64{0, 0, 0, 8, 7},
			y:        []float64{1, 2, 3, 4},
			yWeights: []float64{9, 2, 1, 6},
			dist:     1,
		},
		{
			x:    []float64{1, 2, 3},
			y:    []float64{1, 1, 3},
			dist: 1.0 / 3.0,
		},
		{
			x:        []float64{1, 2, 3},
			y:        []float64{1, 3},
			yWeights: []float64{2, 1},
			dist:     1.0 / 3.0,
		},
		{
			x:        []float64{1, 2, 3},
			xWeights: []float64{2, 2, 2},
			y:        []float64{1, 3},
			yWeights: []float64{2, 1},
			dist:     1.0 / 3.0,
		},
		{
			x:    []float64{2, 3, 4},
			y:    []float64{1, 5},
			dist: 1.0 / 2.0,
		},
		{
			x:    []float64{1, 2, math.NaN()},
			y:    []float64{1, 1, 3},
			dist: math.NaN(),
		},
		{
			x:    []float64{1, 2, 3},
			y:    []float64{1, 1, math.NaN()},
			dist: math.NaN(),
		},
	} {
		dist := KolmogorovSmirnov(test.x, test.xWeights, test.y, test.yWeights)
		if math.Abs(dist-test.dist) > 1e-14 && !(math.IsNaN(test.dist) && math.IsNaN(dist)) {
			t.Errorf("Distance mismatch case %v: Expected: %v, Found: %v", i, test.dist, dist)
		}
	}
	// panic cases
	for _, test := range []struct {
		name     string
		x        []float64
		xWeights []float64
		y        []float64
		yWeights []float64
	}{
		{
			name:     "len(x) != len(xWeights)",
			x:        []float64{1, 3, 5, 6, 7, 8},
			xWeights: []float64{1, 1, 1, 1},
		},
		{
			name:     "len(y) != len(yWeights)",
			x:        []float64{1, 3, 5, 6, 7, 8},
			y:        []float64{1, 3, 5, 6, 7, 8},
			yWeights: []float64{1, 1, 1, 1},
		},
		{
			name: "x not sorted",
			x:    []float64{10, 3, 5, 6, 7, 8},
			y:    []float64{1, 3, 5, 6, 7, 8},
		},
		{
			name: "y not sorted",
			x:    []float64{1, 3, 5, 6, 7, 8},
			y:    []float64{10, 3, 5, 6, 7, 8},
		},
	} {
		if !panics(func() { KolmogorovSmirnov(test.x, test.xWeights, test.y, test.yWeights) }) {
			t.Errorf("KolmogorovSmirnov did not panic when %s", test.name)
		}
	}
}

func TestKullbackLeibler(t *testing.T) {
	if !panics(func() { KullbackLeibler(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("KullbackLeibler did not panic with p, q length mismatch")
	}
}

var linearRegressionTests = []struct {
	name string

	x, y    []float64
	weights []float64
	origin  bool

	alpha float64
	beta  float64
	r     float64

	tol float64
}{
	{
		name: "faithful",

		x: faithful.waiting,
		y: faithful.eruptions,

		// Values calculated by R using lm(eruptions ~ waiting, data=faithful).
		alpha: -1.87402,
		beta:  0.07563,
		r:     0.8114608,

		tol: 1e-5,
	},
	{
		name: "faithful through origin",

		x:      faithful.waiting,
		y:      faithful.eruptions,
		origin: true,

		// Values calculated by R using lm(eruptions ~ waiting - 1, data=faithful).
		alpha: 0,
		beta:  0.05013,
		r:     0.9726036,

		tol: 1e-5,
	},
	{
		name: "faithful explicit weights",

		x: faithful.waiting,
		y: faithful.eruptions,
		weights: func() []float64 {
			w := make([]float64, len(faithful.eruptions))
			for i := range w {
				w[i] = 1
			}
			return w
		}(),

		// Values calculated by R using lm(eruptions ~ waiting, data=faithful).
		alpha: -1.87402,
		beta:  0.07563,
		r:     0.8114608,

		tol: 1e-5,
	},
	{
		name: "faithful non-uniform weights",

		x:       faithful.waiting,
		y:       faithful.eruptions,
		weights: faithful.waiting, // Just an arbitrary set of non-uniform weights.

		// Values calculated by R using lm(eruptions ~ waiting, data=faithful, weights=faithful$waiting).
		alpha: -1.79268,
		beta:  0.07452,
		r:     0.7840372,

		tol: 1e-5,
	},
}

func TestLinearRegression(t *testing.T) {
	for _, test := range linearRegressionTests {
		alpha, beta := LinearRegression(test.x, test.y, test.weights, test.origin)
		var r float64
		if test.origin {
			r = RNoughtSquared(test.x, test.y, test.weights, beta)
		} else {
			r = RSquared(test.x, test.y, test.weights, alpha, beta)
			ests := make([]float64, len(test.y))
			for i, x := range test.x {
				ests[i] = alpha + beta*x
			}
			rvals := RSquaredFrom(ests, test.y, test.weights)
			if r != rvals {
				t.Errorf("%s: RSquared and RSquaredFrom mismatch: %v != %v", test.name, r, rvals)
			}
		}
		if !floats.EqualWithinAbsOrRel(alpha, test.alpha, test.tol, test.tol) {
			t.Errorf("%s: unexpected alpha estimate: want:%v got:%v", test.name, test.alpha, alpha)
		}
		if !floats.EqualWithinAbsOrRel(beta, test.beta, test.tol, test.tol) {
			t.Errorf("%s: unexpected beta estimate: want:%v got:%v", test.name, test.beta, beta)
		}
		if !floats.EqualWithinAbsOrRel(r, test.r, test.tol, test.tol) {
			t.Errorf("%s: unexpected r estimate: want:%v got:%v", test.name, test.r, r)
		}
	}
}

func BenchmarkLinearRegression(b *testing.B) {
	rnd := rand.New(rand.NewSource(1))
	slope, offset := 2.0, 3.0

	maxn := 10000
	xs := make([]float64, maxn)
	ys := make([]float64, maxn)
	weights := make([]float64, maxn)
	for i := range xs {
		x := rnd.Float64()
		xs[i] = x
		ys[i] = slope*x + offset
		weights[i] = rnd.Float64()
	}

	for _, n := range []int{10, 100, 1000, maxn} {
		for _, weighted := range []bool{true, false} {
			for _, origin := range []bool{true, false} {
				name := "n" + strconv.Itoa(n)
				if weighted {
					name += "wt"
				} else {
					name += "wf"
				}
				if origin {
					name += "ot"
				} else {
					name += "of"
				}
				b.Run(name, func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						var ws []float64
						if weighted {
							ws = weights[:n]
						}
						LinearRegression(xs[:n], ys[:n], ws, origin)
					}
				})
			}
		}
	}
}

func TestChiSquare(t *testing.T) {
	for i, test := range []struct {
		p   []float64
		q   []float64
		res float64
	}{
		{
			p:   []float64{16, 18, 16, 14, 12, 12},
			q:   []float64{16, 16, 16, 16, 16, 8},
			res: 3.5,
		},
		{
			p:   []float64{16, 18, 16, 14, 12, 12},
			q:   []float64{8, 20, 20, 16, 12, 12},
			res: 9.25,
		},
		{
			p:   []float64{40, 60, 30, 45},
			q:   []float64{50, 50, 50, 50},
			res: 12.5,
		},
		{
			p:   []float64{40, 60, 30, 45, 0, 0},
			q:   []float64{50, 50, 50, 50, 0, 0},
			res: 12.5,
		},
	} {
		resultpq := ChiSquare(test.p, test.q)

		if math.Abs(resultpq-test.res) > 1e-10 {
			t.Errorf("ChiSquare distance mismatch in case %d. Expected %v, Found %v", i, test.res, resultpq)
		}
	}
	if !panics(func() { ChiSquare(make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("ChiSquare did not panic with length mismatch")
	}
}

// panics returns true if the called function panics during evaluation.
func panics(fun func()) (b bool) {
	defer func() {
		err := recover()
		if err != nil {
			b = true
		}
	}()
	fun()
	return
}

func TestBhattacharyya(t *testing.T) {
	for i, test := range []struct {
		p   []float64
		q   []float64
		res float64
	}{
		{
			p:   []float64{0.5, 0.1, 0.3, 0.1},
			q:   []float64{0.1, 0.4, 0.25, 0.25},
			res: 0.15597338718671386,
		},
		{
			p:   []float64{0.4, 0.6, 0.0},
			q:   []float64{0.2, 0.2, 0.6},
			res: 0.46322207765351153,
		},
		{
			p:   []float64{0.1, 0.1, 0.0, 0.8},
			q:   []float64{0.6, 0.3, 0.0, 0.1},
			res: 0.3552520032137785,
		},
	} {
		resultpq := Bhattacharyya(test.p, test.q)
		resultqp := Bhattacharyya(test.q, test.p)

		if math.Abs(resultpq-test.res) > 1e-10 {
			t.Errorf("Bhattacharyya distance mismatch in case %d. Expected %v, Found %v", i, test.res, resultpq)
		}
		if math.Abs(resultpq-resultqp) > 1e-10 {
			t.Errorf("Bhattacharyya distance is assymmetric in case %d.", i)
		}
	}
	// Bhattacharyya should panic if the inputs have different length
	if !panics(func() { Bhattacharyya(make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("Bhattacharyya did not panic with length mismatch")
	}
}

func TestHellinger(t *testing.T) {
	for i, test := range []struct {
		p   []float64
		q   []float64
		res float64
	}{
		{
			p:   []float64{0.5, 0.1, 0.3, 0.1},
			q:   []float64{0.1, 0.4, 0.25, 0.25},
			res: 0.3800237367441919,
		},
		{
			p:   []float64{0.4, 0.6, 0.0},
			q:   []float64{0.2, 0.2, 0.6},
			res: 0.6088900771170487,
		},
		{
			p:   []float64{0.1, 0.1, 0.0, 0.8},
			q:   []float64{0.6, 0.3, 0.0, 0.1},
			res: 0.5468118803484205,
		},
	} {
		resultpq := Hellinger(test.p, test.q)
		resultqp := Hellinger(test.q, test.p)

		if math.Abs(resultpq-test.res) > 1e-10 {
			t.Errorf("Hellinger distance mismatch in case %d. Expected %v, Found %v", i, test.res, resultpq)
		}
		if math.Abs(resultpq-resultqp) > 1e-10 {
			t.Errorf("Hellinger distance is assymmetric in case %d.", i)
		}
	}
	if !panics(func() { Hellinger(make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("Hellinger did not panic with length mismatch")
	}
}

func TestMean(t *testing.T) {
	if !panics(func() { Mean(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("Mean did not panic with x, weights length mismatch")
	}
}

func TestMode(t *testing.T) {
	for i, test := range []struct {
		x       []float64
		weights []float64
		ans     float64
		count   float64
	}{
		{},
		{
			x:     []float64{1, 6, 1, 9, -2},
			ans:   1,
			count: 2,
		},
		{
			x:       []float64{1, 6, 1, 9, -2},
			weights: []float64{1, 7, 3, 5, 0},
			ans:     6,
			count:   7,
		},
	} {
		m, count := Mode(test.x, test.weights)
		if test.ans != m {
			t.Errorf("Mode mismatch case %d. Expected %v, found %v", i, test.ans, m)
		}
		if test.count != count {
			t.Errorf("Mode count mismatch case %d. Expected %v, found %v", i, test.count, count)
		}
	}
	if !panics(func() { Mode(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("Mode did not panic with x, weights length mismatch")
	}
}

func TestMixedMoment(t *testing.T) {
	for i, test := range []struct {
		x, y, weights []float64
		r, s          float64
		ans           float64
	}{
		{
			x:   []float64{10, 2, 1, 8, 5},
			y:   []float64{8, 15, 1, 6, 3},
			r:   1,
			s:   1,
			ans: 0.48,
		},
		{
			x:       []float64{10, 2, 1, 8, 5},
			y:       []float64{8, 15, 1, 6, 3},
			weights: []float64{1, 1, 1, 1, 1},
			r:       1,
			s:       1,
			ans:     0.48,
		},
		{
			x:       []float64{10, 2, 1, 8, 5},
			y:       []float64{8, 15, 1, 6, 3},
			weights: []float64{2, 3, 0.2, 8, 4},
			r:       1,
			s:       1,
			ans:     -4.786371011357490,
		},
		{
			x:       []float64{10, 2, 1, 8, 5},
			y:       []float64{8, 15, 1, 6, 3},
			weights: []float64{2, 3, 0.2, 8, 4},
			r:       2,
			s:       3,
			ans:     1.598600579313326e+03,
		},
	} {
		m := BivariateMoment(test.r, test.s, test.x, test.y, test.weights)
		if math.Abs(test.ans-m) > 1e-14 {
			t.Errorf("Moment mismatch case %d. Expected %v, found %v", i, test.ans, m)
		}
	}
	if !panics(func() { BivariateMoment(1, 1, make([]float64, 3), make([]float64, 2), nil) }) {
		t.Errorf("Moment did not panic with x, y length mismatch")
	}
	if !panics(func() { BivariateMoment(1, 1, make([]float64, 2), make([]float64, 3), nil) }) {
		t.Errorf("Moment did not panic with x, y length mismatch")
	}
	if !panics(func() { BivariateMoment(1, 1, make([]float64, 2), make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("Moment did not panic with x, weights length mismatch")
	}
	if !panics(func() { BivariateMoment(1, 1, make([]float64, 2), make([]float64, 2), make([]float64, 1)) }) {
		t.Errorf("Moment did not panic with x, weights length mismatch")
	}
}

func TestMoment(t *testing.T) {
	for i, test := range []struct {
		x       []float64
		weights []float64
		moment  float64
		ans     float64
	}{
		{
			x:      []float64{6, 2, 4, 8, 10},
			moment: 5,
			ans:    0,
		},
		{
			x:       []float64{6, 2, 4, 8, 10},
			weights: []float64{1, 2, 2, 2, 1},
			moment:  5,
			ans:     121.875,
		},
	} {
		m := Moment(test.moment, test.x, test.weights)
		if math.Abs(test.ans-m) > 1e-14 {
			t.Errorf("Moment mismatch case %d. Expected %v, found %v", i, test.ans, m)
		}
	}
	if !panics(func() { Moment(1, make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("Moment did not panic with x, weights length mismatch")
	}
	if !panics(func() { Moment(1, make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("Moment did not panic with x, weights length mismatch")
	}
}

func TestMomentAbout(t *testing.T) {
	for i, test := range []struct {
		x       []float64
		weights []float64
		moment  float64
		mean    float64
		ans     float64
	}{
		{
			x:      []float64{6, 2, 4, 8, 9},
			mean:   3,
			moment: 5,
			ans:    2.2288e3,
		},
		{
			x:       []float64{6, 2, 4, 8, 9},
			weights: []float64{1, 2, 2, 2, 1},
			mean:    3,
			moment:  5,
			ans:     1.783625e3,
		},
	} {
		m := MomentAbout(test.moment, test.x, test.mean, test.weights)
		if math.Abs(test.ans-m) > 1e-14 {
			t.Errorf("MomentAbout mismatch case %d. Expected %v, found %v", i, test.ans, m)
		}
	}
	if !panics(func() { MomentAbout(1, make([]float64, 3), 0, make([]float64, 2)) }) {
		t.Errorf("MomentAbout did not panic with x, weights length mismatch")
	}
}

func TestCDF(t *testing.T) {
	cumulantKinds := []CumulantKind{Empirical}
	for i, test := range []struct {
		q       []float64
		x       []float64
		weights []float64
		ans     [][]float64
	}{
		{},
		{
			q:   []float64{0, 0.9, 1, 1.1, 2.9, 3, 3.1, 4.9, 5, 5.1},
			x:   []float64{1, 2, 3, 4, 5},
			ans: [][]float64{{0, 0, 0.2, 0.2, 0.4, 0.6, 0.6, 0.8, 1, 1}},
		},
		{
			q:       []float64{0, 0.9, 1, 1.1, 2.9, 3, 3.1, 4.9, 5, 5.1},
			x:       []float64{1, 2, 3, 4, 5},
			weights: []float64{1, 1, 1, 1, 1},
			ans:     [][]float64{{0, 0, 0.2, 0.2, 0.4, 0.6, 0.6, 0.8, 1, 1}},
		},
		{
			q:   []float64{0, 0.9, 1},
			x:   []float64{math.NaN()},
			ans: [][]float64{{math.NaN(), math.NaN(), math.NaN()}},
		},
	} {
		copyX := make([]float64, len(test.x))
		copy(copyX, test.x)
		var copyW []float64
		if test.weights != nil {
			copyW = make([]float64, len(test.weights))
			copy(copyW, test.weights)
		}
		for j, q := range test.q {
			for k, kind := range cumulantKinds {
				v := CDF(q, kind, test.x, test.weights)
				if !floats.Equal(copyX, test.x) && !math.IsNaN(v) {
					t.Errorf("x changed for case %d kind %d percentile %v", i, k, q)
				}
				if !floats.Equal(copyW, test.weights) {
					t.Errorf("x changed for case %d kind %d percentile %v", i, k, q)
				}
				if v != test.ans[k][j] && !(math.IsNaN(v) && math.IsNaN(test.ans[k][j])) {
					t.Errorf("mismatch case %d kind %d percentile %v. Expected: %v, found: %v", i, k, q, test.ans[k][j], v)
				}
			}
		}
	}

	// these test cases should all result in a panic
	for i, test := range []struct {
		name    string
		q       float64
		kind    CumulantKind
		x       []float64
		weights []float64
	}{
		{
			name:    "len(x) != len(weights)",
			q:       1.5,
			kind:    Empirical,
			x:       []float64{1, 2, 3, 4, 5},
			weights: []float64{1, 2, 3},
		},
		{
			name: "unsorted x",
			q:    1.5,
			kind: Empirical,
			x:    []float64{3, 2, 1},
		},
		{
			name: "unknown CumulantKind",
			q:    1.5,
			kind: CumulantKind(1000), // bogus
			x:    []float64{1, 2, 3},
		},
	} {
		if !panics(func() { CDF(test.q, test.kind, test.x, test.weights) }) {
			t.Errorf("did not panic as expected with %s for case %d kind %d percentile %v x %v weights %v", test.name, i, test.kind, test.q, test.x, test.weights)
		}
	}

}

func TestQuantile(t *testing.T) {
	cumulantKinds := []CumulantKind{
		Empirical,
		LinInterp,
	}
	for i, test := range []struct {
		p   []float64
		x   []float64
		w   []float64
		ans [][]float64
	}{
		{
			p: []float64{0, 0.05, 0.1, 0.15, 0.45, 0.5, 0.55, 0.85, 0.9, 0.95, 1},
			x: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			w: nil,
			ans: [][]float64{
				{1, 1, 1, 2, 5, 5, 6, 9, 9, 10, 10},
				{1, 1, 1, 1.5, 4.5, 5, 5.5, 8.5, 9, 9.5, 10},
			},
		},
		{
			p: []float64{0, 0.05, 0.1, 0.15, 0.45, 0.5, 0.55, 0.85, 0.9, 0.95, 1},
			x: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			w: []float64{3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
			ans: [][]float64{
				{1, 1, 1, 2, 5, 5, 6, 9, 9, 10, 10},
				{1, 1, 1, 1.5, 4.5, 5, 5.5, 8.5, 9, 9.5, 10},
			},
		},
		{
			p: []float64{0, 0.05, 0.1, 0.15, 0.45, 0.5, 0.55, 0.85, 0.9, 0.95, 1},
			x: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			w: []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0},
			ans: [][]float64{
				{1, 2, 3, 4, 7, 7, 8, 10, 10, 10, 10},
				{1, 1.875, 2.833333333333333, 3.5625, 6.535714285714286, 6.928571428571429, 7.281250000000001, 9.175, 9.45, 9.725, 10},
			},
		},
		{
			p: []float64{0.5},
			x: []float64{1, 2, 3, 4, 5, 6, 7, 8, math.NaN(), 10},
			ans: [][]float64{
				{math.NaN()},
				{math.NaN()},
			},
		},
	} {
		copyX := make([]float64, len(test.x))
		copy(copyX, test.x)
		var copyW []float64
		if test.w != nil {
			copyW = make([]float64, len(test.w))
			copy(copyW, test.w)
		}
		for j, p := range test.p {
			for k, kind := range cumulantKinds {
				v := Quantile(p, kind, test.x, test.w)
				if !floats.Same(copyX, test.x) {
					t.Errorf("x changed for case %d kind %d percentile %v", i, k, p)
				}
				if !floats.Same(copyW, test.w) {
					t.Errorf("x changed for case %d kind %d percentile %v", i, k, p)
				}
				if v != test.ans[k][j] && !(math.IsNaN(v) && math.IsNaN(test.ans[k][j])) {
					t.Errorf("mismatch case %d kind %d percentile %v. Expected: %v, found: %v", i, k, p, test.ans[k][j], v)
				}
			}
		}
	}
}

func TestQuantileInvalidInput(t *testing.T) {
	cumulantKinds := []CumulantKind{
		Empirical,
		LinInterp,
	}
	for _, test := range []struct {
		name string
		p    float64
		x    []float64
		w    []float64
	}{
		{
			name: "p < 0",
			p:    -1,
		},
		{
			name: "p > 1",
			p:    2,
		},
		{
			name: "p is NaN",
			p:    math.NaN(),
		},
		{
			name: "len(x) != len(weights)",
			p:    .5,
			x:    make([]float64, 4),
			w:    make([]float64, 2),
		},
		{
			name: "x not sorted",
			p:    .5,
			x:    []float64{3, 2, 1},
		},
	} {
		for _, kind := range cumulantKinds {
			if !panics(func() { Quantile(test.p, kind, test.x, test.w) }) {
				t.Errorf("Quantile did not panic when %s", test.name)
			}
		}
	}
}

func TestQuantileInvalidCumulantKind(t *testing.T) {
	if !panics(func() { Quantile(0.5, CumulantKind(1000), []float64{1, 2, 3}, nil) }) {
		t.Errorf("Quantile did not panic when CumulantKind is unknown")
	}
}

func TestSkew(t *testing.T) {
	for i, test := range []struct {
		x       []float64
		weights []float64
		ans     float64
	}{
		{
			x:       []float64{8, 3, 7, 8, 4},
			weights: nil,
			ans:     -0.581456499151665,
		},
		{
			x:       []float64{8, 3, 7, 8, 4},
			weights: []float64{1, 1, 1, 1, 1},
			ans:     -0.581456499151665,
		},
		{
			x:       []float64{8, 3, 7, 8, 4},
			weights: []float64{2, 1, 2, 1, 1},
			ans:     -1.12066646837198,
		},
	} {
		skew := Skew(test.x, test.weights)
		if math.Abs(skew-test.ans) > 1e-14 {
			t.Errorf("Skew mismatch case %d. Expected %v, Found %v", i, test.ans, skew)
		}
	}
	if !panics(func() { Skew(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("Skew did not panic with x, weights length mismatch")
	}
}

func TestSortWeighted(t *testing.T) {
	for i, test := range []struct {
		x    []float64
		w    []float64
		ansx []float64
		answ []float64
	}{
		{
			x:    []float64{8, 3, 7, 8, 4},
			ansx: []float64{3, 4, 7, 8, 8},
		},
		{
			x:    []float64{8, 3, 7, 8, 4},
			w:    []float64{.5, 1, 1, .5, 1},
			ansx: []float64{3, 4, 7, 8, 8},
			answ: []float64{1, 1, 1, .5, .5},
		},
	} {
		SortWeighted(test.x, test.w)
		if !floats.Same(test.x, test.ansx) {
			t.Errorf("SortWeighted mismatch case %d. Expected x %v, Found x %v", i, test.ansx, test.x)
		}
		if !(test.w == nil) && !floats.Same(test.w, test.answ) {
			t.Errorf("SortWeighted mismatch case %d. Expected w %v, Found w %v", i, test.answ, test.w)
		}
	}
	if !panics(func() { SortWeighted(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("SortWeighted did not panic with x, weights length mismatch")
	}
}

func TestSortWeightedLabeled(t *testing.T) {
	for i, test := range []struct {
		x    []float64
		l    []bool
		w    []float64
		ansx []float64
		ansl []bool
		answ []float64
	}{
		{
			x:    []float64{8, 3, 7, 8, 4},
			ansx: []float64{3, 4, 7, 8, 8},
		},
		{
			x:    []float64{8, 3, 7, 8, 4},
			w:    []float64{.5, 1, 1, .5, 1},
			ansx: []float64{3, 4, 7, 8, 8},
			answ: []float64{1, 1, 1, .5, .5},
		},
		{
			x:    []float64{8, 3, 7, 8, 4},
			l:    []bool{false, false, true, false, true},
			ansx: []float64{3, 4, 7, 8, 8},
			ansl: []bool{false, true, true, false, false},
		},
		{
			x:    []float64{8, 3, 7, 8, 4},
			l:    []bool{false, false, true, false, true},
			w:    []float64{.5, 1, 1, .5, 1},
			ansx: []float64{3, 4, 7, 8, 8},
			ansl: []bool{false, true, true, false, false},
			answ: []float64{1, 1, 1, .5, .5},
		},
	} {
		SortWeightedLabeled(test.x, test.l, test.w)
		if !floats.Same(test.x, test.ansx) {
			t.Errorf("SortWeightedLabelled mismatch case %d. Expected x %v, Found x %v", i, test.ansx, test.x)
		}
		if (test.l != nil) && !reflect.DeepEqual(test.l, test.ansl) {
			t.Errorf("SortWeightedLabelled mismatch case %d. Expected l %v, Found l %v", i, test.ansl, test.l)
		}
		if (test.w != nil) && !floats.Same(test.w, test.answ) {
			t.Errorf("SortWeightedLabelled mismatch case %d. Expected w %v, Found w %v", i, test.answ, test.w)
		}
	}
	if !panics(func() { SortWeightedLabeled(make([]float64, 3), make([]bool, 2), make([]float64, 3)) }) {
		t.Errorf("SortWeighted did not panic with x, labels length mismatch")
	}
	if !panics(func() { SortWeightedLabeled(make([]float64, 3), make([]bool, 2), nil) }) {
		t.Errorf("SortWeighted did not panic with x, labels length mismatch")
	}
	if !panics(func() { SortWeightedLabeled(make([]float64, 3), make([]bool, 3), make([]float64, 2)) }) {
		t.Errorf("SortWeighted did not panic with x, weights length mismatch")
	}
	if !panics(func() { SortWeightedLabeled(make([]float64, 3), nil, make([]float64, 2)) }) {
		t.Errorf("SortWeighted did not panic with x, weights length mismatch")
	}
}

func TestVariance(t *testing.T) {
	for i, test := range []struct {
		x       []float64
		weights []float64
		ans     float64
	}{
		{
			x:       []float64{8, -3, 7, 8, -4},
			weights: nil,
			ans:     37.7,
		},
		{
			x:       []float64{8, -3, 7, 8, -4},
			weights: []float64{1, 1, 1, 1, 1},
			ans:     37.7,
		},
		{
			x:       []float64{8, 3, 7, 8, 4},
			weights: []float64{2, 1, 2, 1, 1},
			ans:     4.2857142857142865,
		},
		{
			x:       []float64{1, 4, 9},
			weights: []float64{1, 1.5, 1},
			ans:     13.142857142857146,
		},
		{
			x:       []float64{1, 2, 3},
			weights: []float64{1, 1.5, 1},
			ans:     .8,
		},
	} {
		variance := Variance(test.x, test.weights)
		if math.Abs(variance-test.ans) > 1e-14 {
			t.Errorf("Variance mismatch case %d. Expected %v, Found %v", i, test.ans, variance)
		}
	}
	if !panics(func() { Variance(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("Variance did not panic with x, weights length mismatch")
	}

}

func TestStdScore(t *testing.T) {
	for i, test := range []struct {
		x float64
		u float64
		s float64
		z float64
	}{
		{
			x: 4,
			u: -6,
			s: 5,
			z: 2,
		},
		{
			x: 1,
			u: 0,
			s: 1,
			z: 1,
		},
	} {
		z := StdScore(test.x, test.u, test.s)
		if math.Abs(z-test.z) > 1e-14 {
			t.Errorf("StdScore mismatch case %d. Expected %v, Found %v", i, test.z, z)
		}
	}
}
