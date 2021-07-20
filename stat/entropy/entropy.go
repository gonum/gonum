package entropy

import (
	"math"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

// Shannon returns the Shannon-entropy. The function takes a probability
// distribution p(x) as input.
//   H(X) = -\sum_x p(x) log(p(x))
func Shannon(p mat.Vector) float64 {
	var r float64
	for i := 0; i < p.Len(); i++ {
		v := p.AtVec(i)
		if v > 0 {
			r -= v * math.Log(v)
		}
	}
	return r
}

// MLBC returns maximum likelihood estimator with bias correction. The function
// requires discretized data as its input.
//    MLBC(X) = (|X| - 1) / (2*n) - \sum p(x) ln p(x)
// where n is the number of data points.
func MLBC(data []int) float64 {
	// Implemented from A. Chao and T.-J. Shen. Nonparametric estimation of
	// Shannon’s index of diversity when there are unseen species in sample.
	// Environmental and Ecological Statistics, 10(4):429–443, 2003.
	var r float64

	c := stat.Counts(data)
	p := stat.CountsToDist(c)
	n := float64(len(data))
	s := float64(p.Len())

	for i := 0; i < p.Len(); i++ {
		v := p.AtVec(i)
		if v > 0 {
			r -= v * math.Log(v)
		}
	}

	return r + (s-1)/(2*n)
}

// HorvitzThompson returns the entropy based on the Horvitz-Thompson estimator.
// The function requires discretized data as its input.
//   H(X) = - \sum (p(x) log p(x))/(1 - (1 - p(x)))
func HorvitzThompson(data []int) float64 {
	// Implemented from A. Chao and T.-J. Shen. Nonparametric estimation of
	// shannon’s index of diversity when there are unseen species in sample.
	// Environmental and Ecological Statistics, 10(4):429–443, 2003.
	var r float64

	c := stat.Counts(data)
	p := stat.CountsToDist(c)

	n := float64(len(data))

	for i := 0; i < p.Len(); i++ {
		v := p.AtVec(i)
		if v > 0 {
			N := v * math.Log(v)
			D := 1 - math.Pow(1-v, n)
			r -= N / D
		}
	}

	return r
}

// ChaoShen returns the entropy based on the Chao-Shen entropy estimator. The
// function requires discretized data as its input.
//   H(X) = - \sum (C p(x) log p(x))/((1 - (1 - C p(x))))
// where C = f_1/n, f_1 = number random variables that occur only once in the
// data, and n is the total number of samples.
func ChaoShen(data []int) float64 {
	// Implemented from A. Chao and T.-J. Shen. Nonparametric estimation of
	// shannon’s index of diversity when there are unseen species in sample.
	// Environmental and Ecological Statistics, 10(4):429–443, 2003.
	var nrOfSingletons, z, r float64

	n := float64(len(data))
	histogram := map[int]float64{}
	for _, v := range data {
		histogram[v]++
	}

	p := make([]float64, len(histogram), len(histogram))

	var keys []int
	for k, v := range histogram {
		keys = append(keys, k)
		if v == 1 {
			nrOfSingletons++
		}
	}

	if nrOfSingletons == n {
		nrOfSingletons--
	}

	for i := range histogram {
		p[i] = histogram[keys[i]] / n
	}

	c := 1 - nrOfSingletons/n

	for i := range p {
		p[i] *= c
	}

	for i := range p {
		if p[i] > 0 {
			z = (1 - math.Pow((1-p[i]), n))
			r -= p[i] * math.Log(p[i]) / z
		}
	}

	return r
}
