// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"math"
	"sort"

	"github.com/gonum/floats"
)

// CumulantKind specifies the behavior for calculating the empirical CDF or Quantile
type CumulantKind int

const (
	// Constant values should match the R nomenclature. See
	// https://en.wikipedia.org/wiki/Quantile#Estimating_the_quantiles_of_a_population

	// Empirical treats the distribution as the actual empirical distribution.
	Empirical CumulantKind = 1
)

// bhattacharyyaCoeff computes the Bhattacharyya Coefficient for probability distributions given by:
//  \sum_i \sqrt{p_i q_i}
//
// It is assumed that p and q have equal length.
func bhattacharyyaCoeff(p, q []float64) float64 {
	var bc float64
	for i, a := range p {
		b := q[i]
		if a == 0 && b == 0 {
			continue
		}
		bc += math.Sqrt(a * b)
	}
	return bc
}

// Bhattacharyya computes the distance between the probability distributions p and q given by:
//  -\ln ( \sum_i \sqrt{p_i q_i} )
//
// The lengths of p and q must be equal. It is assumed that p and q sum to 1.
func Bhattacharyya(p, q []float64) float64 {
	if len(p) != len(q) {
		panic("stat: slice length mismatch")
	}
	bc := bhattacharyyaCoeff(p, q)
	return -math.Log(bc)
}

// CDF returns the empirical cumulative distribution function value of x, that is
// the fraction of the samples less than or equal to q. The
// exact behavior is determined by the CumulantKind. CDF is theoretically
// the inverse of the Quantile function, though it may not be the actual inverse
// for all values q and CumululantKinds.
//
// The x data must be sorted in increasing order. If weights is nil then all
// of the weights are 1. If weights is not nil, then len(x) must equal len(weights).
//
// CumulantKind behaviors:
//  - Empirical: Returns the lowest fraction for which q is greater than or equal
//  to that fraction of samples
func CDF(q float64, c CumulantKind, x, weights []float64) float64 {
	if weights != nil && len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}
	if !sort.Float64sAreSorted(x) {
		panic("x data are not sorted")
	}
	if floats.HasNaN(x) {
		return math.NaN()
	}

	if q < x[0] {
		return 0
	}
	if q >= x[len(x)-1] {
		return 1
	}

	var sumWeights float64
	if weights == nil {
		sumWeights = float64(len(x))
	} else {
		sumWeights = floats.Sum(weights)
	}

	// Calculate the index
	switch c {
	case Empirical:
		// Find the smallest value that is greater than that percent of the samples
		var w float64
		for i, v := range x {
			if v > q {
				return w / sumWeights
			}
			if weights == nil {
				w++
			} else {
				w += weights[i]
			}
		}
		panic("impossible")
	default:
		panic("stat: bad cumulant kind")
	}
}

// ChiSquare computes the chi-square distance between the observed frequences 'obs' and
// expected frequences 'exp' given by:
//  \sum_i (obs_i-exp_i)^2 / exp_i
//
// The lengths of obs and exp must be equal.
func ChiSquare(obs, exp []float64) float64 {
	if len(obs) != len(exp) {
		panic("stat: slice length mismatch")
	}
	var result float64
	for i, a := range obs {
		b := exp[i]
		if a == 0 && b == 0 {
			continue
		}
		result += (a - b) * (a - b) / b
	}
	return result
}

// Correlation returns the weighted correlation between the samples of x and y
// with the given means.
//  sum_i {w_i (x_i - meanX) * (y_i - meanY)} / ((sum_j {w_j} - 1) * stdX * stdY)
// The lengths of x and y must be equal. If weights is nil then all of the
// weights are 1. If weights is not nil, then len(x) must equal len(weights).
func Correlation(x []float64, meanX, stdX float64, y []float64, meanY, stdY float64, weights []float64) float64 {
	return Covariance(x, meanX, y, meanY, weights) / (stdX * stdY)
}

// Covariance returns the weighted covariance between the samples of x and y
// with the given means.
//  sum_i {w_i (x_i - meanX) * (y_i - meanY)} / (sum_j {w_j} - 1)
// The lengths of x and y must be equal. If weights is nil then all of the
// weights are 1. If weights is not nil, then len(x) must equal len(weights).
func Covariance(x []float64, meanX float64, y []float64, meanY float64, weights []float64) float64 {
	if len(x) != len(y) {
		panic("stat: slice length mismatch")
	}
	if weights == nil {
		var s float64
		for i, v := range x {
			s += (v - meanX) * (y[i] - meanY)
		}
		s /= float64(len(x) - 1)
		return s
	}
	if weights != nil && len(weights) != len(x) {
		panic("stat: slice length mismatch")
	}
	var (
		s          float64
		sumWeights float64
	)
	for i, v := range x {
		s += weights[i] * (v - meanX) * (y[i] - meanY)
		sumWeights += weights[i]
	}
	return s / (sumWeights - 1)
}

// CrossEntropy computes the cross-entropy between the two distributions specified
// in p and q.
func CrossEntropy(p, q []float64) float64 {
	if len(p) != len(q) {
		panic("stat: slice length mismatch")
	}
	var ce float64
	for i, v := range p {
		w := q[i]
		if v == 0 && w == 0 {
			continue
		}
		ce -= v * math.Log(w)
	}
	return ce
}

// Entropy computes the Shannon entropy of a distribution or the distance between
// two distributions. The natural logarithm is used.
//  - sum_i (p_i * log_e(p_i))
func Entropy(p []float64) float64 {
	var e float64
	for _, v := range p {
		if v != 0 { // Entropy needs 0 * log(0) == 0
			e -= v * math.Log(v)
		}
	}
	return e
}

// ExKurtosis returns the population excess kurtosis of the sample.
// The kurtosis is defined by the 4th moment of the mean divided by the squared
// variance. The excess kurtosis subtracts 3.0 so that the excess kurtosis of
// the normal distribution is zero.
// If weights is nil then all of the weights are 1. If weights is not nil, then
// len(x) must equal len(weights).
func ExKurtosis(x []float64, mean, stdev float64, weights []float64) float64 {
	if weights == nil {
		var e float64
		for _, v := range x {
			z := (v - mean) / stdev
			e += z * z * z * z
		}
		mul, offset := kurtosisCorrection(float64(len(x)))
		return e*mul - offset
	}
	if len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}
	var (
		e          float64
		sumWeights float64
	)
	for i, v := range x {
		z := (v - mean) / stdev
		e += weights[i] * z * z * z * z
		sumWeights += weights[i]
	}
	mul, offset := kurtosisCorrection(sumWeights)
	return e*mul - offset
}

// n is the number of samples
// see https://en.wikipedia.org/wiki/Kurtosis
func kurtosisCorrection(n float64) (mul, offset float64) {
	return ((n + 1) / (n - 1)) * (n / (n - 2)) * (1 / (n - 3)), 3 * ((n - 1) / (n - 2)) * ((n - 1) / (n - 3))
}

// GeometricMean returns the weighted geometric mean of the dataset
//  \prod_i {x_i ^ w_i}
// This only applies with positive x and positive weights. If weights is nil
// then all of the weights are 1. If weights is not nil, then len(x) must equal
// len(weights).
func GeometricMean(x, weights []float64) float64 {
	if weights == nil {
		var s float64
		for _, v := range x {
			s += math.Log(v)
		}
		s /= float64(len(x))
		return math.Exp(s)
	}
	if len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}
	var (
		s          float64
		sumWeights float64
	)
	for i, v := range x {
		s += weights[i] * math.Log(v)
		sumWeights += weights[i]
	}
	s /= sumWeights
	return math.Exp(s)
}

// HarmonicMean returns the weighted harmonic mean of the dataset
//  \sum_i {w_i} / ( sum_i {w_i / x_i} )
// This only applies with positive x and positive weights.
// If weights is nil then all of the weights are 1. If weights is not nil, then
// len(x) must equal len(weights).
func HarmonicMean(x, weights []float64) float64 {
	if weights != nil && len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}
	// TODO: Fix this to make it more efficient and avoid allocation

	// This can be numerically unstable (for exapmle if x is very small)
	// W = \sum_i {w_i}
	// hm = exp(log(W) - log(\sum_i w_i / x_i))

	logs := make([]float64, len(x))
	var W float64
	for i := range x {
		if weights == nil {
			logs[i] = -math.Log(x[i])
			W++
			continue
		}
		logs[i] = math.Log(weights[i]) - math.Log(x[i])
		W += weights[i]
	}

	// Sum all of the logs
	v := floats.LogSumExp(logs) // this computes log(\sum_i { w_i / x_i})
	return math.Exp(math.Log(W) - v)
}

// Hellinger computes the distance between the probability distributions p and q given by:
//  \sqrt{ 1 - \sum_i \sqrt{p_i q_i} }
//
// The lengths of p and q must be equal. It is assumed that p and q sum to 1.
func Hellinger(p, q []float64) float64 {
	if len(p) != len(q) {
		panic("stat: slice length mismatch")
	}
	bc := bhattacharyyaCoeff(p, q)
	return math.Sqrt(1 - bc)
}

// Histogram sums up the weighted number of data points in each bin.
// The weight of data point x[i] will be placed into count[j] if
// dividers[j-1] <= x < dividers[j]. The "span" function in the floats package can assist
// with bin creation.
//
// The following conditions on the inputs apply:
//  - The count variable must either be nil or have length of one less than dividers.
//  - The values in dividers must be sorted (use the sort package).
//  - The x values must be sorted.
//  - If weights is nil then all of the weights are 1.
//  - If weights is not nil, then len(x) must equal len(weights).
func Histogram(count, dividers, x, weights []float64) []float64 {
	if weights != nil && len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}
	if count == nil {
		count = make([]float64, len(dividers)+1)
	}
	if len(count) != len(dividers)+1 {
		panic("histogram: bin count mismatch")
	}
	if !sort.Float64sAreSorted(dividers) {
		panic("dividers are not sorted")
	}
	if !sort.Float64sAreSorted(x) {
		panic("x data are not sorted")
	}

	idx := 0
	comp := dividers[idx]
	if weights == nil {
		for _, v := range x {
			if v < comp || idx == len(count)-1 {
				// Still in the current bucket
				count[idx]++
				continue
			}
			// Need to find the next divider where v is less than the divider
			// or to set the maximum divider if no such exists
			for j := idx + 1; j < len(count); j++ {
				if j == len(dividers) {
					idx = len(dividers)
					break
				}
				if v < dividers[j] {
					idx = j
					comp = dividers[j]
					break
				}
			}
			count[idx]++
		}
		return count
	}

	for i, v := range x {
		if v < comp || idx == len(count)-1 {
			// Still in the current bucket
			count[idx] += weights[i]
			continue
		}
		// Need to find the next divider where v is less than the divider
		// or to set the maximum divider if no such exists
		for j := idx + 1; j < len(count); j++ {
			if j == len(dividers) {
				idx = len(dividers)
				break
			}
			if v < dividers[j] {
				idx = j
				comp = dividers[j]
				break
			}
		}
		count[idx] += weights[i]
	}
	return count
}

// KulbeckLeibler computes the Kulbeck-Leibler distance between the
// distributions p and q. The natural logarithm is used.
//  sum_i(p_i * log(p_i / q_i))
// Note that the Kulbeck-Leibler distance is not symmetric;
// KulbeckLeibler(p,q) != KulbeckLeibler(q,p)
func KulbeckLeibler(p, q []float64) float64 {
	if len(p) != len(q) {
		panic("stat: slice length mismatch")
	}
	var kl float64
	for i, v := range p {
		if v != 0 { // Entropy needs 0 * log(0) == 0
			kl += v * (math.Log(v) - math.Log(q[i]))
		}
	}
	return kl
}

// Mean computes the weighted mean of the data set.
//  sum_i {w_i * x_i} / sum_i {w_i}
// If weights is nil then all of the weights are 1. If weights is not nil, then
// len(x) must equal len(weights).
func Mean(x, weights []float64) float64 {
	if weights == nil {
		return floats.Sum(x) / float64(len(x))
	}
	if len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}
	var (
		sumValues  float64
		sumWeights float64
	)
	for i, w := range weights {
		sumValues += w * x[i]
		sumWeights += w
	}
	return sumValues / sumWeights
}

// Mode returns the most common value in the dataset specified by x and the
// given weights. Strict float64 equality is used when comparing values, so users
// should take caution. If several values are the mode, any of them may be returned.
func Mode(x []float64, weights []float64) (val float64, count float64) {
	if weights != nil && len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}
	if len(x) == 0 {
		return 0, 0
	}
	m := make(map[float64]float64)
	if weights == nil {
		for _, v := range x {
			m[v]++
		}
	} else {
		for i, v := range x {
			m[v] += weights[i]
		}
	}
	var (
		maxCount float64
		max      float64
	)
	for val, count := range m {
		if count > maxCount {
			maxCount = count
			max = val
		}
	}
	return max, maxCount
}

// Moment computes the weighted n^th moment of the samples,
//  E[(x - μ)^N]
// No degrees of freedom correction is done.
// If weights is nil then all of the weights are 1. If weights is not nil, then
// len(x) must equal len(weights).
func Moment(moment float64, x []float64, mean float64, weights []float64) float64 {
	if weights == nil {
		var m float64
		for _, v := range x {
			m += math.Pow(v-mean, moment)
		}
		m /= float64(len(x))
		return m
	}
	if len(weights) != len(x) {
		panic("stat: slice length mismatch")
	}
	var (
		m          float64
		sumWeights float64
	)
	for i, v := range x {
		m += weights[i] * math.Pow(v-mean, moment)
		sumWeights += weights[i]
	}
	return m / sumWeights
}

// Quantile returns the sample of x such that x is greater than or
// equal to the fraction p of samples. The exact behavior is determined by the
// CumulantKind, and p should be a number between 0 and 1. Quantile is theoretically
// the inverse of the CDF function, though it may not be the actual inverse
// for all values p and CumulantKinds.
//
// The x data must be sorted in increasing order. If weights is nil then all
// of the weights are 1. If weights is not nil, then len(x) must equal len(weights).
//
// CumulantKind behaviors:
//  - Empirical: Returns the lowest value q for which q is greater than or equal
//  to the fraction p of samples
func Quantile(p float64, c CumulantKind, x, weights []float64) float64 {
	if p < 0 || p > 1 {
		panic("stat: percentile out of bounds")
	}

	if weights != nil && len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}
	if !sort.Float64sAreSorted(x) {
		panic("x data are not sorted")
	}
	if floats.HasNaN(x) {
		return math.NaN() // This is needed because the algorithm breaks otherwise
	}
	var sumWeights float64
	if weights == nil {
		sumWeights = float64(len(x))
	} else {
		sumWeights = floats.Sum(weights)
	}
	switch c {
	case Empirical:
		var cumsum float64
		fidx := p * sumWeights
		for i := range x {
			if weights == nil {
				cumsum++
			} else {
				cumsum += weights[i]
			}
			if cumsum >= fidx {
				return x[i]
			}
		}
		panic("impossible")
	default:
		panic("stat: bad cumulant kind")
	}
}

// Skew computes the skewness of the sample data.
// If weights is nil then all of the weights are 1. If weights is not nil, then
// len(x) must equal len(weights).
func Skew(x []float64, mean, stdev float64, weights []float64) float64 {
	if weights == nil {
		var s float64
		for _, v := range x {
			z := (v - mean) / stdev
			s += z * z * z
		}
		return s * skewCorrection(float64(len(x)))
	}
	if len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}
	var (
		s          float64
		sumWeights float64
	)
	for i, v := range x {
		z := (v - mean) / stdev
		s += weights[i] * z * z * z
		sumWeights += weights[i]
	}
	return s * skewCorrection(sumWeights)
}

// From: http://www.amstat.org/publications/jse/v19n2/doane.pdf page 7
func skewCorrection(n float64) float64 {
	return (n / (n - 1)) * (1 / (n - 2))
}

// SortWeighted rearranges the data in x along with their corresponding
// weights so that the x data are sorted. The data is sorted in place.
// Weights may be nil, but if weights is non-nil then it must have the same
// length as x.
func SortWeighted(x, weights []float64) {
	if weights == nil {
		sort.Float64s(x)
		return
	}
	if len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}
	sort.Sort(weightSorter{
		x: x,
		w: weights,
	})
}

type weightSorter struct {
	x []float64
	w []float64
}

func (w weightSorter) Less(i, j int) bool {
	return w.x[i] < w.x[j]
}

func (w weightSorter) Swap(i, j int) {
	w.x[i], w.x[j] = w.x[j], w.x[i]
	w.w[i], w.w[j] = w.w[j], w.w[i]
}

func (w weightSorter) Len() int {
	return len(w.x)
}

// StdDev returns the population standard deviation with the provided mean.
func StdDev(x []float64, mean float64, weights []float64) float64 {
	return math.Sqrt(Variance(x, mean, weights))
}

// StdErr returns the standard error in the mean with the given values.
func StdErr(stdev, sampleSize float64) float64 {
	return stdev / math.Sqrt(sampleSize)
}

// StdScore returns the standard score (a.k.a. z-score, z-value) for the value x
// with the givem mean and variance, i.e.
//  (x - mean) / variance
func StdScore(x, mean, variance float64) float64 {
	return (x - mean) / variance
}

// Variance computes the weighted sample variance with the provided mean.
//  \sum_i w_i (x_i - mean)^2 / (sum_i w_i - 1)
// If weights is nil then all of the weights are 1. If weights is not nil, then
// len(x) must equal len(weights).
func Variance(x []float64, mean float64, weights []float64) float64 {
	if weights == nil {
		var s float64
		for _, v := range x {
			s += (v - mean) * (v - mean)
		}
		return s / float64(len(x)-1)
	}
	if len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}
	var (
		ss         float64
		sumWeights float64
	)
	for i, v := range x {
		ss += weights[i] * (v - mean) * (v - mean)
		sumWeights += weights[i]
	}
	return ss / (sumWeights - 1)
}
