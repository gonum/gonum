// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// a set of benchmarks to evaluate the performance of the various
// moment statistics: Mean, Variance, StdDev, MeanVariance, MeanStdDev,
// Covariance, Correlation, Skew, ExKurtosis, Moment, MomentAbout, ...
//
// It tests both weighted and unweighted versions by using a slice of
// all ones.

package stat

import (
	"math/rand"
	"testing"
)

const (
	SMALL  = 10
	MEDIUM = 1000
	LARGE  = 100000
	HUGE   = 10000000
)

// tests for unweighted versions

func RandomSlice(l int) []float64 {
	s := make([]float64, l)
	for i := range s {
		s[i] = rand.Float64()
	}
	return s
}

func benchmarkMean(b *testing.B, s, wts []float64) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Mean(s, wts)
	}
}

func BenchmarkMeanSmall(b *testing.B) {
	s := RandomSlice(SMALL)
	benchmarkMean(b, s, nil)
}

func BenchmarkMeanMedium(b *testing.B) {
	s := RandomSlice(MEDIUM)
	benchmarkMean(b, s, nil)
}

func BenchmarkMeanLarge(b *testing.B) {
	s := RandomSlice(LARGE)
	benchmarkMean(b, s, nil)
}

func BenchmarkMeanHuge(b *testing.B) {
	s := RandomSlice(HUGE)
	benchmarkMean(b, s, nil)
}

func BenchmarkMeanSmallWeighted(b *testing.B) {
	s := RandomSlice(SMALL)
	wts := RandomSlice(SMALL)
	benchmarkMean(b, s, wts)
}

func BenchmarkMeanMediumWeighted(b *testing.B) {
	s := RandomSlice(MEDIUM)
	wts := RandomSlice(MEDIUM)
	benchmarkMean(b, s, wts)
}

func BenchmarkMeanLargeWeighted(b *testing.B) {
	s := RandomSlice(LARGE)
	wts := RandomSlice(LARGE)
	benchmarkMean(b, s, wts)
}

func BenchmarkMeanHugeWeighted(b *testing.B) {
	s := RandomSlice(HUGE)
	wts := RandomSlice(HUGE)
	benchmarkMean(b, s, wts)
}

func benchmarkVariance(b *testing.B, s, wts []float64) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Variance(s, wts)
	}
}

func BenchmarkVarianceSmall(b *testing.B) {
	s := RandomSlice(SMALL)
	benchmarkVariance(b, s, nil)
}

func BenchmarkVarianceMedium(b *testing.B) {
	s := RandomSlice(MEDIUM)
	benchmarkVariance(b, s, nil)
}

func BenchmarkVarianceLarge(b *testing.B) {
	s := RandomSlice(LARGE)
	benchmarkVariance(b, s, nil)
}

func BenchmarkVarianceHuge(b *testing.B) {
	s := RandomSlice(HUGE)
	benchmarkVariance(b, s, nil)
}

func BenchmarkVarianceSmallWeighted(b *testing.B) {
	s := RandomSlice(SMALL)
	wts := RandomSlice(SMALL)
	benchmarkVariance(b, s, wts)
}

func BenchmarkVarianceMediumWeighted(b *testing.B) {
	s := RandomSlice(MEDIUM)
	wts := RandomSlice(MEDIUM)
	benchmarkVariance(b, s, wts)
}

func BenchmarkVarianceLargeWeighted(b *testing.B) {
	s := RandomSlice(LARGE)
	wts := RandomSlice(LARGE)
	benchmarkVariance(b, s, wts)
}

func BenchmarkVarianceHugeWeighted(b *testing.B) {
	s := RandomSlice(HUGE)
	wts := RandomSlice(HUGE)
	benchmarkVariance(b, s, wts)
}

func benchmarkStdDev(b *testing.B, s, wts []float64) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StdDev(s, wts)
	}
}

func BenchmarkStdDevSmall(b *testing.B) {
	s := RandomSlice(SMALL)
	benchmarkStdDev(b, s, nil)
}

func BenchmarkStdDevMedium(b *testing.B) {
	s := RandomSlice(MEDIUM)
	benchmarkStdDev(b, s, nil)
}

func BenchmarkStdDevLarge(b *testing.B) {
	s := RandomSlice(LARGE)
	benchmarkStdDev(b, s, nil)
}

func BenchmarkStdDevHuge(b *testing.B) {
	s := RandomSlice(HUGE)
	benchmarkStdDev(b, s, nil)
}

func BenchmarkStdDevSmallWeighted(b *testing.B) {
	s := RandomSlice(SMALL)
	wts := RandomSlice(SMALL)
	benchmarkStdDev(b, s, wts)
}

func BenchmarkStdDevMediumWeighted(b *testing.B) {
	s := RandomSlice(MEDIUM)
	wts := RandomSlice(MEDIUM)
	benchmarkStdDev(b, s, wts)
}

func BenchmarkStdDevLargeWeighted(b *testing.B) {
	s := RandomSlice(LARGE)
	wts := RandomSlice(LARGE)
	benchmarkStdDev(b, s, wts)
}

func BenchmarkStdDevHugeWeighted(b *testing.B) {
	s := RandomSlice(HUGE)
	wts := RandomSlice(HUGE)
	benchmarkStdDev(b, s, wts)
}

func benchmarkMeanVariance(b *testing.B, s, wts []float64) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MeanVariance(s, wts)
	}
}

func BenchmarkMeanVarianceSmall(b *testing.B) {
	s := RandomSlice(SMALL)
	benchmarkMeanVariance(b, s, nil)
}

func BenchmarkMeanVarianceMedium(b *testing.B) {
	s := RandomSlice(MEDIUM)
	benchmarkMeanVariance(b, s, nil)
}

func BenchmarkMeanVarianceLarge(b *testing.B) {
	s := RandomSlice(LARGE)
	benchmarkMeanVariance(b, s, nil)
}

func BenchmarkMeanVarianceHuge(b *testing.B) {
	s := RandomSlice(HUGE)
	benchmarkMeanVariance(b, s, nil)
}

func BenchmarkMeanVarianceSmallWeighted(b *testing.B) {
	s := RandomSlice(SMALL)
	wts := RandomSlice(SMALL)
	benchmarkMeanVariance(b, s, wts)
}

func BenchmarkMeanVarianceMediumWeighted(b *testing.B) {
	s := RandomSlice(MEDIUM)
	wts := RandomSlice(MEDIUM)
	benchmarkMeanVariance(b, s, wts)
}

func BenchmarkMeanVarianceLargeWeighted(b *testing.B) {
	s := RandomSlice(LARGE)
	wts := RandomSlice(LARGE)
	benchmarkMeanVariance(b, s, wts)
}

func BenchmarkMeanVarianceHugeWeighted(b *testing.B) {
	s := RandomSlice(HUGE)
	wts := RandomSlice(HUGE)
	benchmarkMeanVariance(b, s, wts)
}

func benchmarkMeanStdDev(b *testing.B, s, wts []float64) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MeanStdDev(s, wts)
	}
}

func BenchmarkMeanStdDevSmall(b *testing.B) {
	s := RandomSlice(SMALL)
	benchmarkMeanStdDev(b, s, nil)
}

func BenchmarkMeanStdDevMedium(b *testing.B) {
	s := RandomSlice(MEDIUM)
	benchmarkMeanStdDev(b, s, nil)
}

func BenchmarkMeanStdDevLarge(b *testing.B) {
	s := RandomSlice(LARGE)
	benchmarkMeanStdDev(b, s, nil)
}

func BenchmarkMeanStdDevHuge(b *testing.B) {
	s := RandomSlice(HUGE)
	benchmarkMeanStdDev(b, s, nil)
}

func BenchmarkMeanStdDevSmallWeighted(b *testing.B) {
	s := RandomSlice(SMALL)
	wts := RandomSlice(SMALL)
	benchmarkMeanStdDev(b, s, wts)
}

func BenchmarkMeanStdDevMediumWeighted(b *testing.B) {
	s := RandomSlice(MEDIUM)
	wts := RandomSlice(MEDIUM)
	benchmarkMeanStdDev(b, s, wts)
}

func BenchmarkMeanStdDevLargeWeighted(b *testing.B) {
	s := RandomSlice(LARGE)
	wts := RandomSlice(LARGE)
	benchmarkMeanStdDev(b, s, wts)
}

func BenchmarkMeanStdDevHugeWeighted(b *testing.B) {
	s := RandomSlice(HUGE)
	wts := RandomSlice(HUGE)
	benchmarkMeanStdDev(b, s, wts)
}

func benchmarkCovariance(b *testing.B, s1, s2, wts []float64) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Covariance(s1, s2, wts)
	}
}

func BenchmarkCovarianceSmall(b *testing.B) {
	s1 := RandomSlice(SMALL)
	s2 := RandomSlice(SMALL)
	benchmarkCovariance(b, s1, s2, nil)
}

func BenchmarkCovarianceMedium(b *testing.B) {
	s1 := RandomSlice(MEDIUM)
	s2 := RandomSlice(MEDIUM)
	benchmarkCovariance(b, s1, s2, nil)
}

func BenchmarkCovarianceLarge(b *testing.B) {
	s1 := RandomSlice(LARGE)
	s2 := RandomSlice(LARGE)
	benchmarkCovariance(b, s1, s2, nil)
}

func BenchmarkCovarianceHuge(b *testing.B) {
	s1 := RandomSlice(HUGE)
	s2 := RandomSlice(HUGE)
	benchmarkCovariance(b, s1, s2, nil)
}

func BenchmarkCovarianceSmallWeighted(b *testing.B) {
	s1 := RandomSlice(SMALL)
	s2 := RandomSlice(SMALL)
	wts := RandomSlice(SMALL)
	benchmarkCovariance(b, s1, s2, wts)
}

func BenchmarkCovarianceMediumWeighted(b *testing.B) {
	s1 := RandomSlice(MEDIUM)
	s2 := RandomSlice(MEDIUM)
	wts := RandomSlice(MEDIUM)
	benchmarkCovariance(b, s1, s2, wts)
}

func BenchmarkCovarianceLargeWeighted(b *testing.B) {
	s1 := RandomSlice(LARGE)
	s2 := RandomSlice(LARGE)
	wts := RandomSlice(LARGE)
	benchmarkCovariance(b, s1, s2, wts)
}

func BenchmarkCovarianceHugeWeighted(b *testing.B) {
	s1 := RandomSlice(HUGE)
	s2 := RandomSlice(HUGE)
	wts := RandomSlice(HUGE)
	benchmarkCovariance(b, s1, s2, wts)
}

func benchmarkCorrelation(b *testing.B, s1, s2, wts []float64) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Correlation(s1, s2, wts)
	}
}

func BenchmarkCorrelationSmall(b *testing.B) {
	s1 := RandomSlice(SMALL)
	s2 := RandomSlice(SMALL)
	benchmarkCorrelation(b, s1, s2, nil)
}

func BenchmarkCorrelationMedium(b *testing.B) {
	s1 := RandomSlice(MEDIUM)
	s2 := RandomSlice(MEDIUM)
	benchmarkCorrelation(b, s1, s2, nil)
}

func BenchmarkCorrelationLarge(b *testing.B) {
	s1 := RandomSlice(LARGE)
	s2 := RandomSlice(LARGE)
	benchmarkCorrelation(b, s1, s2, nil)
}

func BenchmarkCorrelationHuge(b *testing.B) {
	s1 := RandomSlice(HUGE)
	s2 := RandomSlice(HUGE)
	benchmarkCorrelation(b, s1, s2, nil)
}

func BenchmarkCorrelationSmallWeighted(b *testing.B) {
	s1 := RandomSlice(SMALL)
	s2 := RandomSlice(SMALL)
	wts := RandomSlice(SMALL)
	benchmarkCorrelation(b, s1, s2, wts)
}

func BenchmarkCorrelationMediumWeighted(b *testing.B) {
	s1 := RandomSlice(MEDIUM)
	s2 := RandomSlice(MEDIUM)
	wts := RandomSlice(MEDIUM)
	benchmarkCorrelation(b, s1, s2, wts)
}

func BenchmarkCorrelationLargeWeighted(b *testing.B) {
	s1 := RandomSlice(LARGE)
	s2 := RandomSlice(LARGE)
	wts := RandomSlice(LARGE)
	benchmarkCorrelation(b, s1, s2, wts)
}

func BenchmarkCorrelationHugeWeighted(b *testing.B) {
	s1 := RandomSlice(HUGE)
	s2 := RandomSlice(HUGE)
	wts := RandomSlice(HUGE)
	benchmarkCorrelation(b, s1, s2, wts)
}

func benchmarkSkew(b *testing.B, s, wts []float64) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Skew(s, wts)
	}
}

func BenchmarkSkewSmall(b *testing.B) {
	s := RandomSlice(SMALL)
	benchmarkSkew(b, s, nil)
}

func BenchmarkSkewMedium(b *testing.B) {
	s := RandomSlice(MEDIUM)
	benchmarkSkew(b, s, nil)
}

func BenchmarkSkewLarge(b *testing.B) {
	s := RandomSlice(LARGE)
	benchmarkSkew(b, s, nil)
}

func BenchmarkSkewHuge(b *testing.B) {
	s := RandomSlice(HUGE)
	benchmarkSkew(b, s, nil)
}

func BenchmarkSkewSmallWeighted(b *testing.B) {
	s := RandomSlice(SMALL)
	wts := RandomSlice(SMALL)
	benchmarkSkew(b, s, wts)
}

func BenchmarkSkewMediumWeighted(b *testing.B) {
	s := RandomSlice(MEDIUM)
	wts := RandomSlice(MEDIUM)
	benchmarkSkew(b, s, wts)
}

func BenchmarkSkewLargeWeighted(b *testing.B) {
	s := RandomSlice(LARGE)
	wts := RandomSlice(LARGE)
	benchmarkSkew(b, s, wts)
}

func BenchmarkSkewHugeWeighted(b *testing.B) {
	s := RandomSlice(HUGE)
	wts := RandomSlice(HUGE)
	benchmarkSkew(b, s, wts)
}

func benchmarkExKurtosis(b *testing.B, s, wts []float64) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExKurtosis(s, wts)
	}
}

func BenchmarkExKurtosisSmall(b *testing.B) {
	s := RandomSlice(SMALL)
	benchmarkExKurtosis(b, s, nil)
}

func BenchmarkExKurtosisMedium(b *testing.B) {
	s := RandomSlice(MEDIUM)
	benchmarkExKurtosis(b, s, nil)
}

func BenchmarkExKurtosisLarge(b *testing.B) {
	s := RandomSlice(LARGE)
	benchmarkExKurtosis(b, s, nil)
}

func BenchmarkExKurtosisHuge(b *testing.B) {
	s := RandomSlice(HUGE)
	benchmarkExKurtosis(b, s, nil)
}

func BenchmarkExKurtosisSmallWeighted(b *testing.B) {
	s := RandomSlice(SMALL)
	wts := RandomSlice(SMALL)
	benchmarkExKurtosis(b, s, wts)
}

func BenchmarkExKurtosisMediumWeighted(b *testing.B) {
	s := RandomSlice(MEDIUM)
	wts := RandomSlice(MEDIUM)
	benchmarkExKurtosis(b, s, wts)
}

func BenchmarkExKurtosisLargeWeighted(b *testing.B) {
	s := RandomSlice(LARGE)
	wts := RandomSlice(LARGE)
	benchmarkExKurtosis(b, s, wts)
}

func BenchmarkExKurtosisHugeWeighted(b *testing.B) {
	s := RandomSlice(HUGE)
	wts := RandomSlice(HUGE)
	benchmarkExKurtosis(b, s, wts)
}

func benchmarkMoment(b *testing.B, n float64, s, wts []float64) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Moment(n, s, wts)
	}
}

func BenchmarkMomentSmall(b *testing.B) {
	s := RandomSlice(SMALL)
	benchmarkMoment(b, 5, s, nil)
}

func BenchmarkMomentMedium(b *testing.B) {
	s := RandomSlice(MEDIUM)
	benchmarkMoment(b, 5, s, nil)
}

func BenchmarkMomentLarge(b *testing.B) {
	s := RandomSlice(LARGE)
	benchmarkMoment(b, 5, s, nil)
}

func BenchmarkMomentHuge(b *testing.B) {
	s := RandomSlice(HUGE)
	benchmarkMoment(b, 5, s, nil)
}

func BenchmarkMomentSmallWeighted(b *testing.B) {
	s := RandomSlice(SMALL)
	wts := RandomSlice(SMALL)
	benchmarkMoment(b, 5, s, wts)
}

func BenchmarkMomentMediumWeighted(b *testing.B) {
	s := RandomSlice(MEDIUM)
	wts := RandomSlice(MEDIUM)
	benchmarkMoment(b, 5, s, wts)
}

func BenchmarkMomentLargeWeighted(b *testing.B) {
	s := RandomSlice(LARGE)
	wts := RandomSlice(LARGE)
	benchmarkMoment(b, 5, s, wts)
}

func BenchmarkMomentHugeWeighted(b *testing.B) {
	s := RandomSlice(HUGE)
	wts := RandomSlice(HUGE)
	benchmarkMoment(b, 5, s, wts)
}

func benchmarkMomentAbout(b *testing.B, n float64, s []float64, mean float64, wts []float64) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MomentAbout(n, s, mean, wts)
	}
}

func BenchmarkMomentAboutSmall(b *testing.B) {
	s := RandomSlice(SMALL)
	benchmarkMomentAbout(b, 5, s, 0, nil)
}

func BenchmarkMomentAboutMedium(b *testing.B) {
	s := RandomSlice(MEDIUM)
	benchmarkMomentAbout(b, 5, s, 0, nil)
}

func BenchmarkMomentAboutLarge(b *testing.B) {
	s := RandomSlice(LARGE)
	benchmarkMomentAbout(b, 5, s, 0, nil)
}

func BenchmarkMomentAboutHuge(b *testing.B) {
	s := RandomSlice(HUGE)
	benchmarkMomentAbout(b, 5, s, 0, nil)
}

func BenchmarkMomentAboutSmallWeighted(b *testing.B) {
	s := RandomSlice(SMALL)
	wts := RandomSlice(SMALL)
	benchmarkMomentAbout(b, 5, s, 0, wts)
}

func BenchmarkMomentAboutMediumWeighted(b *testing.B) {
	s := RandomSlice(MEDIUM)
	wts := RandomSlice(MEDIUM)
	benchmarkMomentAbout(b, 5, s, 0, wts)
}

func BenchmarkMomentAboutLargeWeighted(b *testing.B) {
	s := RandomSlice(LARGE)
	wts := RandomSlice(LARGE)
	benchmarkMomentAbout(b, 5, s, 0, wts)
}

func BenchmarkMomentAboutHugeWeighted(b *testing.B) {
	s := RandomSlice(HUGE)
	wts := RandomSlice(HUGE)
	benchmarkMomentAbout(b, 5, s, 0, wts)
}
