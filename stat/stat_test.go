// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
)

// It'd be better if we could use distuv.Normal but that would create a cyclic import
// Instead I created a random sample with the following python script
// ```
// import numpy as np
// np.random.seed(28041990)
// s = np.random.normal(0, 1, 100)
// for v in s:
//     print(str(v) + ",")
// ```
var norm = []float64{
	-0.42559355556264555,
	1.710539107218753,
	-0.33074948677055843,
	-0.5161417685463089,
	-0.6193456407278632,
	1.1351354008003518,
	1.3997307933412553,
	-0.029954251382664998,
	0.8434220441471606,
	2.1128080605675374,
	1.2484609152871478,
	1.8344987008482856,
	0.09931267796096076,
	1.428919892713464,
	1.0189271326046607,
	-0.2414299164875163,
	0.9066103078332822,
	0.3401187294324855,
	0.37733889116557606,
	-0.7323545583397159,
	-0.31958107231420857,
	-0.2512518199227104,
	-1.1383397596501,
	1.0035254463416154,
	1.1418421937745262,
	-1.2358788434742574,
	1.0100704641117955,
	-0.08036883653794552,
	-0.58513457961876,
	0.05084968136313145,
	0.5664086258955331,
	-0.43308548122088375,
	-3.1178336091684975,
	0.7603796596063905,
	0.2853512117407265,
	0.2176337503977105,
	2.152757471948443,
	-1.3087387278036695,
	0.9978177912374294,
	1.786766843343617,
	0.2306781426352178,
	1.4431459612084379,
	-0.6566580121689873,
	0.4684545103705484,
	-0.8951779888297282,
	0.822016499896265,
	-0.16799507924776486,
	-0.0020099175487792724,
	0.11204979393586693,
	1.5055567069562388,
	0.10391457337177848,
	-1.3699213591924695,
	-1.0072435369731148,
	0.9237191947896234,
	-0.019525907747287095,
	-1.0887672634881844,
	0.2751885920079707,
	-0.3434944093974466,
	0.4246733152123921,
	-0.6853964678846619,
	0.46427348833740134,
	-1.1769384002155117,
	0.1711116430457653,
	-1.7144590897345933,
	3.522255994195385,
	-0.37646202717640115,
	-2.433505931729811,
	1.256781201303011,
	-0.15655339743939028,
	0.4552009397642383,
	-0.18924946568753523,
	-1.3725653125945136,
	1.6436894604846044,
	0.597721827583972,
	-0.04796417108532695,
	0.5909740089894597,
	-0.32067675553881103,
	-0.9249846629433014,
	-0.5678094294504085,
	-0.6602157563037903,
	1.9230857720458412,
	-0.22916228851590298,
	-0.4570292925300327,
	-0.5904063162660567,
	0.9813766542000734,
	-0.15441746156260575,
	1.3686357504466162,
	-1.775095655630572,
	-1.0537085236430228,
	0.21705168855671522,
	0.10583883865411757,
	0.680143572120755,
	0.043548775401668044,
	-0.8204090347120757,
	-0.9641761107363781,
	-1.2735654936668594,
	0.6970301935177855,
	1.0141887564307395,
	0.49425150550631447,
	0.8237595013007176,
}

func ExampleCircularMean() {
	x := []float64{0, 0.25 * math.Pi, 0.75 * math.Pi}
	weights := []float64{1, 2, 2.5}
	cmean := CircularMean(x, weights)

	fmt.Printf("The circular mean is %.5f.\n", cmean)
	// Output:
	// The circular mean is 1.37037.
}

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

func ExampleCorrelation() {
	x := []float64{8, -3, 7, 8, -4}
	y := []float64{10, 5, 6, 3, -1}
	w := []float64{2, 1.5, 3, 3, 2}

	fmt.Println("Correlation computes the degree to which two datasets move together")
	fmt.Println("about their mean. For example, x and y above move similarly.")

	c := Correlation(x, y, w)
	fmt.Printf("Correlation is %.5f\n", c)

	// Output:
	// Correlation computes the degree to which two datasets move together
	// about their mean. For example, x and y above move similarly.
	// Correlation is 0.59915
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

func ExampleKendall() {
	x := []float64{8, -3, 7, 8, -4}
	y := []float64{10, 5, 6, 3, -1}
	w := []float64{2, 1.5, 3, 3, 2}

	fmt.Println("Kendall correlation computes the number of ordered pairs")
	fmt.Println("between two datasets.")

	c := Kendall(x, y, w)
	fmt.Printf("Kendall correlation is %.5f\n", c)

	// Output:
	// Kendall correlation computes the number of ordered pairs
	// between two datasets.
	// Kendall correlation is 0.25000
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

func ExampleCovariance() {
	fmt.Println("Covariance computes the degree to which datasets move together")
	fmt.Println("about their mean.")
	x := []float64{8, -3, 7, 8, -4}
	y := []float64{10, 2, 2, 4, 1}
	cov := Covariance(x, y, nil)
	fmt.Printf("Cov = %.4f\n", cov)
	fmt.Println("If datasets move perfectly together, the variance equals the covariance")
	y2 := []float64{12, 1, 11, 12, 0}
	cov2 := Covariance(x, y2, nil)
	varX := Variance(x, nil)
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
	kurt := ExKurtosis(x, nil)
	fmt.Printf("ExKurtosis = %.5f\n", kurt)
	weights := []float64{1, 2, 3, 5}
	wKurt := ExKurtosis(x, weights)
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
	if !panics(func() { ExKurtosis(make([]float64, 3), make([]float64, 2)) }) {
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

func ExampleHistogram() {
	x := make([]float64, 101)
	for i := range x {
		x[i] = 1.1 * float64(i) // x data ranges from 0 to 110
	}
	dividers := []float64{0, 7, 20, 100, 1000}
	fmt.Println(`Histogram counts the amount of data in the bins specified by
the dividers. In this data set, there are 7 data points less than 7 (between dividers[0]
and dividers[1]), 12 data points between 7 and 20 (dividers[1] and dividers[2]),
and 0 data points above 1000. Since dividers has length 5, there will be 4 bins.`)
	hist := Histogram(nil, dividers, x, nil)
	fmt.Printf("Hist = %v\n", hist)

	fmt.Println()
	fmt.Println("For ease, the floats Span function can be used to set the dividers")
	nBins := 10
	dividers = make([]float64, nBins+1)
	min := floats.Min(x)
	max := floats.Max(x)
	// Increase the maximum divider so that the maximum value of x is contained
	// within the last bucket.
	max++
	floats.Span(dividers, min, max)
	// Span includes the min and the max. Trim the dividers to create 10 buckets
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
	// the dividers. In this data set, there are 7 data points less than 7 (between dividers[0]
	// and dividers[1]), 12 data points between 7 and 20 (dividers[1] and dividers[2]),
	// and 0 data points above 1000. Since dividers has length 5, there will be 4 bins.
	// Hist = [7 12 72 10]
	//
	// For ease, the floats Span function can be used to set the dividers
	// Hist = [11 10 10 10 10 10 10 10 10 10]
	//
	// Histogram also works with weighted data, and allows reusing of
	// the count field in order to avoid extra garbage
	// Weighted Hist = [66 165 265 365 465 565 665 765 865 965]
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

func ExampleStdDev() {
	x := []float64{8, 2, -9, 15, 4}
	stdev := StdDev(x, nil)
	fmt.Printf("The standard deviation of the samples is %.4f\n", stdev)

	weights := []float64{2, 2, 6, 7, 1}
	weightedStdev := StdDev(x, weights)
	fmt.Printf("The weighted standard deviation of the samples is %.4f\n", weightedStdev)
	// Output:
	// The standard deviation of the samples is 8.8034
	// The weighted standard deviation of the samples is 10.5733
}

func ExampleStdErr() {
	x := []float64{8, 2, -9, 15, 4}
	weights := []float64{2, 2, 6, 7, 1}
	mean := Mean(x, weights)
	stdev := StdDev(x, weights)
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

func ExampleVariance() {
	x := []float64{8, 2, -9, 15, 4}
	variance := Variance(x, nil)
	fmt.Printf("The variance of the samples is %.4f\n", variance)

	weights := []float64{2, 2, 6, 7, 1}
	weightedVariance := Variance(x, weights)
	fmt.Printf("The weighted variance of the samples is %.4f\n", weightedVariance)
	// Output:
	// The variance of the samples is 77.5000
	// The weighted variance of the samples is 111.7941
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

func TestNormalSkewTest(t *testing.T) {
	// even though the skew test and kurtosis test are a copy from the scipy's the results
	// are not exactly the same because the skew/kurtosis corrections are not exactly the same
	// but the statistic value is close enough
	for i, test := range []struct {
		values   []float64
		expected float64
	}{
		{
			values:   []float64{0, 1, 2, 3, 4, 5, 6, 7, 8},
			expected: 1.018464355396213,
		},
		{
			values:   []float64{2, 8, 0, 4, 1, 9, 9, 0},
			expected: 0.5562585562766172,
		},
		{
			values:   []float64{1, 2, 3, 4, 5, 6, 7, 8000},
			expected: 4.319816401673864,
		},
		{
			values:   []float64{100, 100, 100, 100, 100, 100, 100, 101},
			expected: 4.319820025201098,
		},
	} {
		z := NormalSkewTest(test.values)
		if math.Abs(z-test.expected) > 1e-7 {
			t.Errorf("NormalSkewTest mismatch case %d. Expected %v, Found %v", i, test.expected, z)
		}
	}
}

func TestNormalKurtosisTest(t *testing.T) {
	for i, test := range []struct {
		values   []float64
		expected float64
	}{
		{
			values:   []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			expected: -1.6867033202073243,
		},
		{
			values:   norm,
			expected: 1.638462148403484,
		},
	} {
		z := NormalKurtosisTest(test.values)
		if math.Abs(z-test.expected) > 1e-7 {
			t.Errorf("NormalKurtosisTest mismatch case %d. Expected %v, Found %v", i, test.expected, z)
		}
	}
}

func TestNormalTest(t *testing.T) {
	// normal test with scipy yield similar results:
	// ```
	// from scipy.stats import normaltest
	// import numpy as np
	// np.random.seed(28041990)
	// s = np.random.normal(0, 1, 100)
	// print(normaltest(list(range(9))).statistic)
	// print(normaltest([2, 8, 0, 4, 1, 9, 9, 0]).statistic)
	// print(normaltest([1, 2, 3, 4, 5, 6, 7, 8000]).statistic)
	// print(normaltest([100, 100, 100, 100, 100, 100, 100, 101]).statistic)
	// print(normaltest(list(range(20))).statistic)
	// print(normaltest(s).statistic)
	// ```

	for i, test := range []struct {
		values   []float64
		expected float64
	}{
		{
			values:   []float64{0, 1, 2, 3, 4, 5, 6, 7, 8},
			expected: 1.7509245653153176,
		},
		{
			values:   []float64{2, 8, 0, 4, 1, 9, 9, 0},
			expected: 11.454757293481551,
		},
		{
			values:   []float64{1, 2, 3, 4, 5, 6, 7, 8000},
			expected: 40.53534243515444,
		},
		{
			values:   []float64{100, 100, 100, 100, 100, 100, 100, 101},
			expected: 40.53539760601764,
		},
		{
			values:   []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			expected: 3.9272951079276743,
		},
		{
			values:   norm,
			expected: 2.695532408378151,
		},
	} {
		z := NormalTest(test.values)
		if math.Abs(z-test.expected) > 1e-7 {
			t.Errorf("NormalTest mismatch case %d. Expected %v, Found %v", i, test.expected, z)
		}
	}
}
