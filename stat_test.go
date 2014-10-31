// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"fmt"
	"math"
	"testing"

	"github.com/gonum/floats"
)

func ExampleCorrelation() {
	x := []float64{8, -3, 7, 8, -4}
	y := []float64{10, 5, 6, 3, -1}
	w := []float64{2, 1.5, 3, 3, 2}

	fmt.Println("Correlation computes the degree to which two datasets move together")
	fmt.Println("about their mean. For example, x and y above move similarly.")
	fmt.Println("Package can be used to compute the mean and standard deviation")
	fmt.Println("or they can be supplied if they are known")
	meanX := Mean(x, w)
	meanY := Mean(x, w)
	c := Correlation(x, meanX, 3, y, meanY, 4, w)
	fmt.Printf("Correlation with set standard deviatons is %.5f\n", c)
	stdX := StdDev(x, meanX, w)
	stdY := StdDev(x, meanY, w)
	c2 := Correlation(x, meanX, stdX, y, meanY, stdY, w)
	fmt.Printf("Correlation with computed standard deviatons is %.5f\n", c2)
	// Output:
	// Correlation computes the degree to which two datasets move together
	// about their mean. For example, x and y above move similarly.
	// Package can be used to compute the mean and standard deviation
	// or they can be supplied if they are known
	// Correlation with set standard deviatons is 0.96894
	// Correlation with computed standard deviatons is 0.39644
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
		meanX := Mean(test.x, test.w)
		meanY := Mean(test.y, test.w)
		stdX := StdDev(test.x, meanX, test.w)
		stdY := StdDev(test.y, meanY, test.w)
		c := Correlation(test.x, meanX, stdX, test.y, meanY, stdY, test.w)
		if math.Abs(test.ans-c) > 1e-14 {
			t.Errorf("Correlation mismatch case %d. Expected %v, Found %v", i, test.ans, c)
		}
	}
}

func ExampleCovariance() {
	fmt.Println("Covariance computes the degree to which datasets move together")
	fmt.Println("about their mean.")
	x := []float64{8, -3, 7, 8, -4}
	y := []float64{10, 2, 2, 4, 1}
	meanX := Mean(x, nil)
	meanY := Mean(y, nil)
	cov := Covariance(x, meanX, y, meanY, nil)
	fmt.Printf("Cov = %.4f\n", cov)
	fmt.Println("If datasets move perfectly together, the variance equals the covariance")
	y2 := []float64{12, 1, 11, 12, 0}
	meanY2 := Mean(y2, nil)
	cov2 := Covariance(x, meanX, y2, meanY2, nil)
	varX := Variance(x, meanX, nil)
	fmt.Printf("Cov2 is %.4f, VarX is %.4f", cov2, varX)
	// Output:
	// Covariance computes the degree to which datasets move together
	// about their mean.
	// Cov = 13.8000
	// If datasets move perfectly together, the variance equals the covariance
	// Cov2 is 37.7000, VarX is 37.7000
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
	} {
		c := Covariance(test.p, Mean(test.p, test.weights), test.q, Mean(test.q, test.weights), test.weights)
		if math.Abs(c-test.ans) > 1e-14 {
			t.Errorf("Covariance mismatch case %d: Expected %v, Found %v", i, test.ans, c)
		}
	}

	// test the panic states
	if !Panics(func() { Covariance(make([]float64, 2), 0.0, make([]float64, 3), 0.0, nil) }) {
		t.Errorf("Covariance did not panic with x, y length mismatch")
	}
	if !Panics(func() { Covariance(make([]float64, 3), 0.0, make([]float64, 3), 0.0, make([]float64, 2)) }) {
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
	if !Panics(func() { CrossEntropy(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("CrossEntropy did not panic with p, q length mismatch")
	}
}

func ExampleEntropy() {

	p := []float64{0.05, 0.1, 0.9, 0.05}
	entP := Entropy(p)

	q := []float64{0.2, 0.4, 0.25, 0.15}
	entQ := Entropy(q)

	r := []float64{0.2, 0, 0, 0.5, 0, 0.2, 0.1, 0, 0, 0}
	entR := Entropy(r)

	s := []float64{0, 0, 1, 0}
	entS := Entropy(s)

	fmt.Println("Entropy is a measure of the amount of uncertainty in a distribution")
	fmt.Printf("The second bin of p is very likely to occur. It's entropy is %.4f\n", entP)
	fmt.Printf("The distribution of q is more spread out. It's entropy is %.4f\n", entQ)
	fmt.Println("Adding buckets with zero probability does not change the entropy.")
	fmt.Printf("The entropy of r is: %.4f\n", entR)
	fmt.Printf("A distribution with no uncertainty has entropy %.4f\n", entS)
	// Output:
	// Entropy is a measure of the amount of uncertainty in a distribution
	// The second bin of p is very likely to occur. It's entropy is 0.6247
	// The distribution of q is more spread out. It's entropy is 1.3195
	// Adding buckets with zero probability does not change the entropy.
	// The entropy of r is: 1.2206
	// A distribution with no uncertainty has entropy 0.0000
}

func ExampleExKurtosis() {
	fmt.Println(`Kurtosis is a measure of the 'peakedness' of a distribution, and the
excess kurtosis is the kurtosis above or below that of the standard normal
distribution`)
	x := []float64{5, 4, -3, -2}
	mean := Mean(x, nil)
	stdev := StdDev(x, mean, nil)
	kurt := ExKurtosis(x, mean, stdev, nil)
	fmt.Printf("ExKurtosis = %.5f\n", kurt)
	weights := []float64{1, 2, 3, 5}
	wMean := Mean(x, weights)
	wStdev := StdDev(x, wMean, weights)
	wKurt := ExKurtosis(x, wMean, wStdev, weights)
	fmt.Printf("Weighted ExKurtosis is %.4f", wKurt)
	// Output:
	// Kurtosis is a measure of the 'peakedness' of a distribution, and the
	// excess kurtosis is the kurtosis above or below that of the standard normal
	// distribution
	// ExKurtosis = -5.41200
	// Weighted ExKurtosis is -0.6779
}

func TestExKurtosis(t *testing.T) {
	// the example does a good job, this just has to cover the panic
	if !Panics(func() { ExKurtosis(make([]float64, 3), 0.0, 0.0, make([]float64, 2)) }) {
		t.Errorf("ExKurtosis did not panic with x, weights length mismatch")
	}
}

func ExampleGeometricMean() {
	x := []float64{8, 2, 9, 15, 4}
	weights := []float64{2, 2, 6, 7, 1}
	mean := Mean(x, weights)
	gmean := GeometricMean(x, weights)

	logx := make([]float64, len(x))
	for i, v := range x {
		logx[i] = math.Log(v)
	}
	expMeanLog := math.Exp(Mean(logx, weights))
	fmt.Printf("The arithmetic mean is %.4f, but the geometric mean is %.4f.\n", mean, gmean)
	fmt.Printf("The exponential of the mean of the logs is %.4f\n", expMeanLog)
	// Output:
	// The arithmetic mean is 10.1667, but the geometric mean is 8.7637.
	// The exponential of the mean of the logs is 8.7637
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
	} {
		c := GeometricMean(test.x, test.wts)
		if math.Abs(c-test.ans) > 1e-14 {
			t.Errorf("Geometric mean mismatch case %d: Expected %v, Found %v", i, test.ans, c)
		}
	}
	if !Panics(func() { GeometricMean(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("GeometricMean did not panic with x, wts length mismatch")
	}
}

func ExampleHarmonicMean() {
	x := []float64{8, 2, 9, 15, 4}
	weights := []float64{2, 2, 6, 7, 1}
	mean := Mean(x, weights)
	hmean := HarmonicMean(x, weights)

	fmt.Printf("The arithmetic mean is %.5f, but the harmonic mean is %.4f.\n", mean, hmean)
	// Output:
	// The arithmetic mean is 10.16667, but the harmonic mean is 6.8354.
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
	} {
		c := HarmonicMean(test.x, test.wts)
		if math.Abs(c-test.ans) > 1e-14 {
			t.Errorf("Harmonic mean mismatch case %d: Expected %v, Found %v", i, test.ans, c)
		}
	}
	if !Panics(func() { HarmonicMean(make([]float64, 3), make([]float64, 2)) }) {
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
			dividers: []float64{2, 4, 6, 7},
			ans:      []float64{1, 1, 1, 1, 2},
		},
		{
			x:        []float64{1, 3, 5, 6, 7, 8},
			dividers: []float64{2, 4, 6, 7},
			weights:  []float64{1, 2, 1, 1, 1, 2},
			ans:      []float64{1, 2, 1, 1, 3},
		},
		{
			x:        []float64{1, 8},
			dividers: []float64{2, 4, 6, 7},
			weights:  []float64{1, 2},
			ans:      []float64{1, 0, 0, 0, 2},
		},
		{
			x:        []float64{1, 8},
			dividers: []float64{2, 4, 6, 7},
			ans:      []float64{1, 0, 0, 0, 1},
		},
	} {
		hist := Histogram(nil, test.dividers, test.x, test.weights)
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
			name:     "len(dividers) != len(count)",
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
	} {
		if !Panics(func() { Histogram(test.count, test.dividers, test.x, test.weights) }) {
			t.Errorf("Histogram did not panic when %s", test.name)
		}
	}
}

func ExampleHistogram() {
	x := make([]float64, 101)
	for i := range x {
		x[i] = 1.1 * float64(i) // x data ranges from 0 to 110
	}
	dividers := []float64{7, 20, 100, 1000}
	fmt.Println(`Histogram counts the amount of data in the bins specified by
the dividers. In this data set, there are 7 data points less than 7 (dividers[0]),
12 data points between 7 and 20 (dividers[2] and dividers[1]), and 0 data points
above 1000. Since dividers has length 4, there will be 5 bins.`)
	hist := Histogram(nil, dividers, x, nil)
	fmt.Printf("Hist = %v\n", hist)

	fmt.Println()
	fmt.Println("For ease, the floats Span function can be used to set the dividers")
	nBins := 10
	// Create one fewer divider than bins, but add two to work with Span (see
	// note below)
	dividers = make([]float64, nBins+1)
	min, _ := floats.Min(x)
	max, _ := floats.Max(x)
	floats.Span(dividers, min, max)
	// Span includes the min and the max. Trim the dividers to create 10 buckets
	dividers = dividers[1 : len(dividers)-1]
	fmt.Println("len dividers = ", len(dividers))
	hist = Histogram(nil, dividers, x, nil)
	fmt.Printf("Hist = %v\n", hist)
	fmt.Println()
	fmt.Println(`Histogram also works with weighted data, and allows reusing of
the count field in order to avoid extra garbage`)
	weights := make([]float64, len(x))
	for i := range weights {
		weights[i] = float64(i + 1)
	}
	Histogram(hist, dividers, x, weights)
	fmt.Printf("Weighted Hist = %v\n", hist)

	// Output:
	// Histogram counts the amount of data in the bins specified by
	// the dividers. In this data set, there are 7 data points less than 7 (dividers[0]),
	// 12 data points between 7 and 20 (dividers[2] and dividers[1]), and 0 data points
	// above 1000. Since dividers has length 4, there will be 5 bins.
	// Hist = [7 12 72 10 0]
	//
	// For ease, the floats Span function can be used to set the dividers
	// len dividers =  9
	// Hist = [11 10 10 10 9 11 10 10 9 11]
	//
	// Histogram also works with weighted data, and allows reusing of
	// the count field in order to avoid extra garbage
	// Weighted Hist = [77 175 275 375 423 627 675 775 783 1067]
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
	if !Panics(func() { JensenShannon(make([]float64, 3), make([]float64, 2)) }) {
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
		if !Panics(func() { KolmogorovSmirnov(test.x, test.xWeights, test.y, test.yWeights) }) {
			t.Errorf("KolmogorovSmirnov did not panic when %s", test.name)
		}
	}
}

func ExampleKullbackLeibler() {

	p := []float64{0.05, 0.1, 0.9, 0.05}
	q := []float64{0.2, 0.4, 0.25, 0.15}
	s := []float64{0, 0, 1, 0}

	klPQ := KullbackLeibler(p, q)
	klPS := KullbackLeibler(p, s)
	klPP := KullbackLeibler(p, p)

	fmt.Println("Kullback-Leibler is one measure of the difference between two distributions")
	fmt.Printf("The K-L distance between p and q is %.4f\n", klPQ)
	fmt.Println("It is impossible for s and p to be the same distribution, because")
	fmt.Println("the first bucket has zero probability in s and non-zero in p. Thus,")
	fmt.Printf("the K-L distance between them is %.4f\n", klPS)
	fmt.Printf("The K-L distance between identical distributions is %.4f\n", klPP)

	// Kullback-Leibler is one measure of the difference between two distributions
	// The K-L distance between p and q is 0.8900
	// It is impossible for s and p to be the same distribution, because
	// the first bucket has zero probability in s and non-zero in p. Thus,
	// the K-L distance between them is +Inf
	// The K-L distance between identical distributions is 0.0000
}

func TestKullbackLeibler(t *testing.T) {
	if !Panics(func() { KullbackLeibler(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("KullbackLeibler did not panic with p, q length mismatch")
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
	if !Panics(func() { ChiSquare(make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("ChiSquare did not panic with length mismatch")
	}
}

// Panics returns true if the called function panics during evaluation.
func Panics(fun func()) (b bool) {
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
	if !Panics(func() { Bhattacharyya(make([]float64, 2), make([]float64, 3)) }) {
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
	if !Panics(func() { Hellinger(make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("Hellinger did not panic with length mismatch")
	}
}

func ExampleMean() {
	x := []float64{8.2, -6, 5, 7}
	mean := Mean(x, nil)
	fmt.Printf("The mean of the samples is %.4f\n", mean)
	w := []float64{2, 6, 3, 5}
	weightedMean := Mean(x, w)
	fmt.Printf("The weighted mean of the samples is %.4f\n", weightedMean)
	x2 := []float64{8.2, 8.2, -6, -6, -6, -6, -6, -6, 5, 5, 5, 7, 7, 7, 7, 7}
	mean2 := Mean(x2, nil)
	fmt.Printf("The mean of x2 is %.4f\n", mean2)
	fmt.Println("The weights act as if there were more samples of that number")
	// Output:
	// The mean of the samples is 3.5500
	// The weighted mean of the samples is 1.9000
	// The mean of x2 is 1.9000
	// The weights act as if there were more samples of that number
}
func TestMean(t *testing.T) {
	if !Panics(func() { Mean(make([]float64, 3), make([]float64, 2)) }) {
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
	if !Panics(func() { Mode(make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("Mode did not panic with x, weights length mismatch")
	}
}

func TestMoment(t *testing.T) {
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
		m := Moment(test.moment, test.x, test.mean, test.weights)
		if math.Abs(test.ans-m) > 1e-14 {
			t.Errorf("Moment mismatch case %d. Expected %v, found %v", i, test.ans, m)
		}
	}
	if !Panics(func() { Moment(1, make([]float64, 3), 0, make([]float64, 2)) }) {
		t.Errorf("Moment did not panic with x, weights length mismatch")
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
		if !Panics(func() { CDF(test.q, test.kind, test.x, test.weights) }) {
			t.Errorf("did not panic as expected with %s for case %d kind %d percentile %v x %v weights %v", test.name, i, test.kind, test.q, test.x, test.weights)
		}
	}

}

func TestQuantile(t *testing.T) {
	cumulantKinds := []CumulantKind{Empirical}
	for i, test := range []struct {
		p   []float64
		x   []float64
		w   []float64
		ans [][]float64
	}{
		{
			p:   []float64{0, 0.05, 0.1, 0.15, 0.45, 0.5, 0.55, 0.85, 0.9, 0.95, 1},
			x:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			w:   nil,
			ans: [][]float64{{1, 1, 1, 2, 5, 5, 6, 9, 9, 10, 10}},
		},
		{
			p:   []float64{0, 0.05, 0.1, 0.15, 0.45, 0.5, 0.55, 0.85, 0.9, 0.95, 1},
			x:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			w:   []float64{3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
			ans: [][]float64{{1, 1, 1, 2, 5, 5, 6, 9, 9, 10, 10}},
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
				if !floats.Equal(copyX, test.x) {
					t.Errorf("x changed for case %d kind %d percentile %v", i, k, p)
				}
				if !floats.Equal(copyW, test.w) {
					t.Errorf("x changed for case %d kind %d percentile %v", i, k, p)
				}
				if v != test.ans[k][j] {
					t.Errorf("mismatch case %d kind %d percentile %v. Expected: %v, found: %v", i, k, p, test.ans[k][j], v)
				}
			}
		}
	}
}

func ExampleStdDev() {
	x := []float64{8, 2, -9, 15, 4}
	mean := Mean(x, nil)
	stdev := StdDev(x, mean, nil)
	fmt.Printf("The standard deviation of the samples is %.4f\n", stdev)

	weights := []float64{2, 2, 6, 7, 1}
	weightedMean := Mean(x, weights)
	weightedStdev := StdDev(x, weightedMean, weights)
	fmt.Printf("The weighted standard deviation of the samples is %.4f\n", weightedStdev)
	// Output:
	// The standard deviation of the samples is 8.8034
	// The weighted standard deviation of the samples is 10.5733
}

func ExampleStdErr() {
	x := []float64{8, 2, -9, 15, 4}
	weights := []float64{2, 2, 6, 7, 1}
	mean := Mean(x, weights)
	stdev := StdDev(x, mean, weights)
	nSamples := floats.Sum(weights)
	stdErr := StdErr(stdev, nSamples)
	fmt.Printf("The standard deviation is %.4f and there are %g samples, so the mean\nis likely %.4f ± %.4f.", stdev, nSamples, mean, stdErr)
	// Output:
	// The standard deviation is 10.5733 and there are 18 samples, so the mean
	// is likely 4.1667 ± 2.4921.
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
		mean := Mean(test.x, test.weights)
		std := StdDev(test.x, mean, test.weights)
		skew := Skew(test.x, mean, std, test.weights)
		if math.Abs(skew-test.ans) > 1e-14 {
			t.Errorf("Skew mismatch case %d. Expected %v, Found %v", i, test.ans, skew)
		}
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
	} {
		mean := Mean(test.x, test.weights)
		variance := Variance(test.x, mean, test.weights)
		if math.Abs(variance-test.ans) > 1e-14 {
			t.Errorf("Variance mismatch case %d. Expected %v, Found %v", i, test.ans, variance)
		}
	}
}

func ExampleVariance() {
	x := []float64{8, 2, -9, 15, 4}
	mean := Mean(x, nil)
	variance := Variance(x, mean, nil)
	fmt.Printf("The variance of the samples is %.4f\n", variance)

	weights := []float64{2, 2, 6, 7, 1}
	weightedMean := Mean(x, weights)
	weightedVariance := Variance(x, weightedMean, weights)
	fmt.Printf("The weighted variance of the samples is %.4f\n", weightedVariance)
	// Output:
	// The variance of the samples is 77.5000
	// The weighted variance of the samples is 111.7941
}
