package stat

import (
	"math"
	"sort"

	"github.com/gonum/floats"
)

// Correlation returns the weighted correlation between the samples of x and y
// with the given means.
// 		sum_i {w_i (x_i - meanX) * (y_i - meanY)} / ((sum_j {w_j} - 1) * stdX * stdY)
// The lengths of x and y must be equal
// If weights is nil then all of the weights are 1
// If weights is not nil, then len(x) must equal len(weights)
func Correlation(x []float64, meanX, stdX float64, y []float64, meanY, stdY float64, weights []float64) float64 {
	return Covariance(x, meanX, y, meanY, weights) / (stdX * stdY)
}

// Covariance returns the weighted covariance between the samples of x and y
// with the given means.
// 		sum_i {w_i (x_i - meanX) * (y_i - meanY)} / (sum_j {w_j} - 1)
// The lengths of x and y must be equal
// If weights is nil then all of the weights are 1
// If weights is not nil, then len(x) must equal len(weights)
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
	var s float64
	var sumWeights float64
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
//		- sum_i (p_i * log_e(p_i))
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
// If weights is nil then all of the weights are 1
// If weights is not nil, then len(x) must equal len(weights)
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
	var e float64
	var sumWeights float64
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

// GeoMean returns the weighted geometric mean of the dataset
// 		\prod_i {x_i ^ w_i}
// This only applies with positive x and positive weights
// If weights is nil then all of the weights are 1
// If weights is not nil, then len(x) must equal len(weights)
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
	var s float64
	var sumWeights float64
	for i, v := range x {
		s += weights[i] * math.Log(v)
		sumWeights += weights[i]
	}
	s /= sumWeights
	return math.Exp(s)
}

// GeoMean returns the weighted harmonic mean of the dataset
// 		\sum_i {w_i} / ( sum_i {w_i / x_i} )
// This only applies with positive x and positive weights
// If weights is nil then all of the weights are 1
// If weights is not nil, then len(x) must equal len(weights)
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
			W += 1
			continue
		}
		logs[i] = math.Log(weights[i]) - math.Log(x[i])
		W += weights[i]
	}

	// Sum all of the logs
	v := floats.LogSumExp(logs) // this computes log(\sum_i { w_i / x_i})
	return math.Exp(math.Log(W) - v)
}

// Histogram sums up the weighted number of data points in each bin.
// The weight of data point x[i] will be placed into count[j] if
// dividers[j-1] <= x < dividers[j]. The "span" function in the floats package can assist
// with bin creation. The count variable must either be nil or have length of
// one less than dividers.
// If weights is nil then all of the weights are 1
// If weights is not nil, then len(x) must equal len(weights)
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

	sortX, sortWeight := sortXandWeight(x, weights)

	idx := 0
	comp := dividers[idx]
	if sortWeight == nil {
		for _, v := range sortX {
			if v < comp || idx == len(count)-1 {
				// Still in the current bucket
				count[idx] += 1
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
			count[idx] += 1
		}
		return count
	}

	for i, v := range sortX {
		if v < comp || idx == len(count)-1 {
			// Still in the current bucket
			count[idx] += sortWeight[i]
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
		count[idx] += sortWeight[i]
	}
	return count

	return count
}

// KulbeckLeibler computes the Kulbeck-Leibler distance between the
// distributions p and q. The natural logarithm is used.
//		sum_i(p_i * log(p_i / q_i))
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
//     sum_i {w_i * x_i} / sum_i {w_i}
// If weights is nil then all of the weights are 1
// If weights is not nil, then len(x) must equal len(weights)
func Mean(x, weights []float64) float64 {
	if weights == nil {
		return floats.Sum(x) / float64(len(x))
	}
	if len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}
	var sumValues float64
	var sumWeights float64
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
			m[v] += 1
		}
	} else {
		for i, v := range x {
			m[v] += weights[i]
		}
	}
	var maxCount float64
	var max float64
	for val, count := range m {
		if count > maxCount {
			maxCount = count
			max = val
		}
	}
	return max, maxCount
}

// Moment computes the weighted n^th moment of the samples,
// 		E[(x - μ)^N]
// No degrees of freedom correction is done.
// If weights is nil then all of the weights are 1
// If weights is not nil, then len(x) must equal len(weights)
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
	var m float64
	var sumWeights float64
	for i, v := range x {
		m += weights[i] * math.Pow(v-mean, moment)
		sumWeights += weights[i]
	}
	return m / sumWeights
}

// Percentile returns the lowest sample of x such that x is greater than or
// equal to the fraction p of samples. p should be a number between 0 and 1
// If no such sample exists, the lowest value is returned
// If weights is nil then all of the weights are 1
// If weights is not nil, then len(x) must equal len(weights)
func Percentile(p float64, x, weights []float64) float64 {
	if p < 0 || p > 1 {
		panic("stat: percentile out of bounds")
	}

	if weights != nil && len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}

	sortX, sortWeight := sortXandWeight(x, weights)
	if weights == nil {
		loc := p * float64(len(x))
		idx := int(math.Floor(loc))
		if (loc == float64(idx) && idx != 0) || idx == len(x) {
			idx--
		}
		return sortX[idx]
	}

	idx := p * floats.Sum(weights)
	var cumsum float64
	for i, w := range sortWeight {
		cumsum += w
		if cumsum >= idx {
			return sortX[i]
		}
	}
	panic("shouldn't be here")
}

// Quantile returns the lowest number p such that q is >= the fraction p of samples
// It is the inverse of the Percentile function.
// If weights is nil then all of the weights are 1
// If weights is not nil, then len(x) must equal len(weights)
func Quantile(q float64, x, weights []float64) float64 {
	if weights != nil && len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}
	sortX, sortWeight := sortXandWeight(x, weights)

	// Find the first x that is greater than the supplied x
	if q < sortX[0] {
		return 0
	}
	if q >= sortX[len(sortX)-1] {
		return 1
	}

	if weights == nil {
		for i, v := range sortX {
			if v > q {
				return float64(i) / float64(len(x))
			}
		}
	}
	sumWeights := floats.Sum(weights)
	var w float64
	for i, v := range sortX {
		if v > q {
			return w / sumWeights
		}
		w += sortWeight[i]
	}
	panic("Impossible. Maybe x contains NaNs.")
}

// Skew computes the skewness of the sample data
// If weights is nil then all of the weights are 1
// If weights is not nil, then len(x) must equal len(weights)
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
	var s float64
	var sumWeights float64
	for i, v := range x {
		z := (v - mean) / stdev
		s += weights[i] * z * z * z
		sumWeights += weights[i]
	}
	return s * skewCorrection(sumWeights)
}

func skewCorrection(n float64) float64 {
	// http://www.amstat.org/publications/jse/v19n2/doane.pdf page 7
	return (n / (n - 1)) * (1 / (n - 2))
}

// StdDev returns the population standard deviation with the provided mean
func StDev(x []float64, mean float64, weights []float64) float64 {
	return math.Sqrt(Variance(x, mean, weights))
}

// StandardError returns the standard error in the mean with the given values
func StdErr(stdev, sampleSize float64) float64 {
	return stdev / math.Sqrt(sampleSize)
}

// StdScore returns the standard score (a.k.a. z-score, z-value) for the value x
// with the givem mean and variance, i.e.
//		(x - mean) / variance
func StdScore(x, mean, variance float64) float64 {
	return (x - mean) / variance
}

// Variance computes the weighted sample variance with the provided mean.
//    \sum_i w_i (x_i - mean)^2 / (sum_i w_i - 1)
// If weights is nil, then all of the weights are 1.
// If weights in not nil, then len(x) must equal len(weights).
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
	var ss float64
	var sumWeights float64
	for i, v := range x {
		ss += weights[i] * (v - mean) * (v - mean)
		sumWeights += weights[i]
	}
	return ss / (sumWeights - 1)
}

// Quartile returns
//func Quartile(x []float64, weights []float64) float64 {}

func sortXandWeight(x, weights []float64) (sortX, sortWeight []float64) {

	sorted := sort.Float64sAreSorted(x)
	if !sorted {
		sortX = make([]float64, len(x))
		copy(sortX, x)
		inds := make([]int, len(x))
		floats.Argsort(sortX, inds)
		if weights != nil {
			sortWeight = make([]float64, len(x))
			for i, v := range inds {
				sortWeight[i] = weights[v]
			}
		}
	} else {
		sortX = x
		sortWeight = weights
	}
	return
}
