// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat_test

import (
	"fmt"
	"math"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

func ExampleCircularMean() {
	x := []float64{0, 0.25 * math.Pi, 0.75 * math.Pi}
	weights := []float64{1, 2, 2.5}
	cmean := stat.CircularMean(x, weights)

	fmt.Printf("The circular mean is %.5f.\n", cmean)
	// Output:
	// The circular mean is 1.37037.
}

func ExampleCorrelation() {
	x := []float64{8, -3, 7, 8, -4}
	y := []float64{10, 5, 6, 3, -1}
	w := []float64{2, 1.5, 3, 3, 2}

	fmt.Println("Correlation computes the degree to which two datasets move together")
	fmt.Println("about their mean. For example, x and y above move similarly.")

	c := stat.Correlation(x, y, w)
	fmt.Printf("Correlation is %.5f\n", c)

	// Output:
	// Correlation computes the degree to which two datasets move together
	// about their mean. For example, x and y above move similarly.
	// Correlation is 0.59915
}

func ExampleCovariance() {
	fmt.Println("Covariance computes the degree to which datasets move together")
	fmt.Println("about their mean.")
	x := []float64{8, -3, 7, 8, -4}
	y := []float64{10, 2, 2, 4, 1}
	cov := stat.Covariance(x, y, nil)
	fmt.Printf("Cov = %.4f\n", cov)
	fmt.Println("If datasets move perfectly together, the variance equals the covariance")
	y2 := []float64{12, 1, 11, 12, 0}
	cov2 := stat.Covariance(x, y2, nil)
	varX := stat.Variance(x, nil)
	fmt.Printf("Cov2 is %.4f, VarX is %.4f", cov2, varX)

	// Output:
	// Covariance computes the degree to which datasets move together
	// about their mean.
	// Cov = 13.8000
	// If datasets move perfectly together, the variance equals the covariance
	// Cov2 is 37.7000, VarX is 37.7000
}

func ExampleEntropy() {

	p := []float64{0.05, 0.1, 0.9, 0.05}
	entP := stat.Entropy(p)

	q := []float64{0.2, 0.4, 0.25, 0.15}
	entQ := stat.Entropy(q)

	r := []float64{0.2, 0, 0, 0.5, 0, 0.2, 0.1, 0, 0, 0}
	entR := stat.Entropy(r)

	s := []float64{0, 0, 1, 0}
	entS := stat.Entropy(s)

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
	kurt := stat.ExKurtosis(x, nil)
	fmt.Printf("ExKurtosis = %.5f\n", kurt)
	weights := []float64{1, 2, 3, 5}
	wKurt := stat.ExKurtosis(x, weights)
	fmt.Printf("Weighted ExKurtosis is %.4f", wKurt)

	// Output:
	// Kurtosis is a measure of the 'peakedness' of a distribution, and the
	// excess kurtosis is the kurtosis above or below that of the standard normal
	// distribution
	// ExKurtosis = -5.41200
	// Weighted ExKurtosis is -0.6779
}

func ExampleGeometricMean() {
	x := []float64{8, 2, 9, 15, 4}
	weights := []float64{2, 2, 6, 7, 1}
	mean := stat.Mean(x, weights)
	gmean := stat.GeometricMean(x, weights)

	logx := make([]float64, len(x))
	for i, v := range x {
		logx[i] = math.Log(v)
	}
	expMeanLog := math.Exp(stat.Mean(logx, weights))
	fmt.Printf("The arithmetic mean is %.4f, but the geometric mean is %.4f.\n", mean, gmean)
	fmt.Printf("The exponential of the mean of the logs is %.4f\n", expMeanLog)

	// Output:
	// The arithmetic mean is 10.1667, but the geometric mean is 8.7637.
	// The exponential of the mean of the logs is 8.7637
}

func ExampleHarmonicMean() {
	x := []float64{8, 2, 9, 15, 4}
	weights := []float64{2, 2, 6, 7, 1}
	mean := stat.Mean(x, weights)
	hmean := stat.HarmonicMean(x, weights)

	fmt.Printf("The arithmetic mean is %.5f, but the harmonic mean is %.4f.\n", mean, hmean)
	// Output:
	// The arithmetic mean is 10.16667, but the harmonic mean is 6.8354.
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
	hist := stat.Histogram(nil, dividers, x, nil)
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
	hist = stat.Histogram(nil, dividers, x, nil)
	fmt.Printf("Hist = %v\n", hist)
	fmt.Println()
	fmt.Println(`Histogram also works with weighted data, and allows reusing of
the count field in order to avoid extra garbage`)
	weights := make([]float64, len(x))
	for i := range weights {
		weights[i] = float64(i + 1)
	}
	stat.Histogram(hist, dividers, x, weights)
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

func ExampleKendall() {
	x := []float64{8, -3, 7, 8, -4}
	y := []float64{10, 5, 6, 3, -1}
	w := []float64{2, 1.5, 3, 3, 2}

	fmt.Println("Kendall correlation computes the number of ordered pairs")
	fmt.Println("between two datasets.")

	c := stat.Kendall(x, y, w)
	fmt.Printf("Kendall correlation is %.5f\n", c)

	// Output:
	// Kendall correlation computes the number of ordered pairs
	// between two datasets.
	// Kendall correlation is 0.25000
}

func ExampleKullbackLeibler() {

	p := []float64{0.05, 0.1, 0.9, 0.05}
	q := []float64{0.2, 0.4, 0.25, 0.15}
	s := []float64{0, 0, 1, 0}

	klPQ := stat.KullbackLeibler(p, q)
	klPS := stat.KullbackLeibler(p, s)
	klPP := stat.KullbackLeibler(p, p)

	fmt.Println("Kullback-Leibler is one measure of the difference between two distributions")
	fmt.Printf("The K-L distance between p and q is %.4f\n", klPQ)
	fmt.Println("It is impossible for s and p to be the same distribution, because")
	fmt.Println("the first bucket has zero probability in s and non-zero in p. Thus,")
	fmt.Printf("the K-L distance between them is %.4f\n", klPS)
	fmt.Printf("The K-L distance between identical distributions is %.4f\n", klPP)

	// Output:
	// Kullback-Leibler is one measure of the difference between two distributions
	// The K-L distance between p and q is 0.8900
	// It is impossible for s and p to be the same distribution, because
	// the first bucket has zero probability in s and non-zero in p. Thus,
	// the K-L distance between them is +Inf
	// The K-L distance between identical distributions is 0.0000
}

func ExampleLinearRegression() {
	var (
		xs      = make([]float64, 100)
		ys      = make([]float64, 100)
		weights []float64
	)

	line := func(x float64) float64 {
		return 1 + 3*x
	}

	for i := range xs {
		xs[i] = float64(i)
		ys[i] = line(xs[i]) + 0.1*rand.NormFloat64()
	}

	// Do not force the regression line to pass through the origin.
	origin := false

	alpha, beta := stat.LinearRegression(xs, ys, weights, origin)
	r2 := stat.RSquared(xs, ys, weights, alpha, beta)

	fmt.Printf("Estimated slope is:  %.6f\n", alpha)
	fmt.Printf("Estimated offset is: %.6f\n", beta)
	fmt.Printf("R^2: %.6f\n", r2)

	// Output:
	// Estimated slope is:  0.988572
	// Estimated offset is: 3.000154
	// R^2: 0.999999
}

func ExampleMean() {
	x := []float64{8.2, -6, 5, 7}
	mean := stat.Mean(x, nil)
	fmt.Printf("The mean of the samples is %.4f\n", mean)
	w := []float64{2, 6, 3, 5}
	weightedMean := stat.Mean(x, w)
	fmt.Printf("The weighted mean of the samples is %.4f\n", weightedMean)
	x2 := []float64{8.2, 8.2, -6, -6, -6, -6, -6, -6, 5, 5, 5, 7, 7, 7, 7, 7}
	mean2 := stat.Mean(x2, nil)
	fmt.Printf("The mean of x2 is %.4f\n", mean2)
	fmt.Println("The weights act as if there were more samples of that number")

	// Output:
	// The mean of the samples is 3.5500
	// The weighted mean of the samples is 1.9000
	// The mean of x2 is 1.9000
	// The weights act as if there were more samples of that number
}

func ExampleStdDev() {
	x := []float64{8, 2, -9, 15, 4}
	stdev := stat.StdDev(x, nil)
	fmt.Printf("The standard deviation of the samples is %.4f\n", stdev)

	weights := []float64{2, 2, 6, 7, 1}
	weightedStdev := stat.StdDev(x, weights)
	fmt.Printf("The weighted standard deviation of the samples is %.4f\n", weightedStdev)

	// Output:
	// The standard deviation of the samples is 8.8034
	// The weighted standard deviation of the samples is 10.5733
}

func ExampleStdErr() {
	x := []float64{8, 2, -9, 15, 4}
	weights := []float64{2, 2, 6, 7, 1}
	mean := stat.Mean(x, weights)
	stdev := stat.StdDev(x, weights)
	nSamples := floats.Sum(weights)
	stdErr := stat.StdErr(stdev, nSamples)
	fmt.Printf("The standard deviation is %.4f and there are %g samples, so the mean\nis likely %.4f ± %.4f.", stdev, nSamples, mean, stdErr)

	// Output:
	// The standard deviation is 10.5733 and there are 18 samples, so the mean
	// is likely 4.1667 ± 2.4921.
}

func ExampleVariance() {
	x := []float64{8, 2, -9, 15, 4}
	variance := stat.Variance(x, nil)
	fmt.Printf("The variance of the samples is %.4f\n", variance)

	weights := []float64{2, 2, 6, 7, 1}
	weightedVariance := stat.Variance(x, weights)
	fmt.Printf("The weighted variance of the samples is %.4f\n", weightedVariance)

	// Output:
	// The variance of the samples is 77.5000
	// The weighted variance of the samples is 111.7941
}
