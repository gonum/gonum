package it

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// Entropy calculates the entropy of a probability distribution.
//   H(X) = -\sum_x p(x) log(p(x))
func Entropy(p mat.Vector) float64 {
	var r float64
	for i := 0; i < p.Len(); i++ {
		v := p.AtVec(i)
		if v > 0.0 {
			r -= v * math.Log(v)
		}
	}
	return r
}

// EntropyMLBC is maximum likelihood estimator with bias correction
// It takes discretised data as input.
// Implemented from
// A. Chao and T.-J. Shen. Nonparametric estimation of Shannon’s
// index of diversity when there are unseen species in sample.
// Environmental and Ecological Statistics, 10(4):429–443, 2003.
func EntropyMLBC(data []int) float64 {
	p := Emperical1D(data)
	n := float64(len(data))
	S := float64(p.Len())

	r := 0.0

	for i := 0; i < p.Len(); i++ {
		v := p.AtVec(i)
		if v > 0.0 {
			r -= v * math.Log(v)
		}
	}

	return r + (S-1.0)/(2.0*n)
}

// EntropyHorvitzThompson is the Horvitz-Thompson entropy estimator.
// It takes discretised data as input.
// Implemented from
// A. Chao and T.-J. Shen. Nonparametric estimation of shannon’s
// index of diversity when there are unseen species in sample.
// Environmental and Ecological Statistics, 10(4):429–443, 2003.
func EntropyHorvitzThompson(data []int) float64 {
	p := Emperical1D(data)
	n := float64(len(data))
	r := 0.0

	for i := 0; i < p.Len(); i++ {
		v := p.AtVec(i)
		if v > 0.0 {
			N := v * math.Log(v)
			D := 1.0 - math.Pow(1.0-v, n)
			r -= N / D
		}
	}

	return r
}

// EntropyChaoShen is the Chao-Shen entropy estimator. It take discretised data
// and the log-function as input
// Implemented from
// A. Chao and T.-J. Shen. Nonparametric estimation of shannon’s
// index of diversity when there are unseen species in sample.
// Environmental and Ecological Statistics, 10(4):429–443, 2003.
func EntropyChaoShen(data []int) float64 {
	n := float64(len(data))
	nrOfSingletons := 0.0
	histogram := map[int]float64{}
	for _, v := range data {
		histogram[v] += 1.0
	}

	p := make([]float64, len(histogram), len(histogram))

	var keys []int
	for k, v := range histogram {
		keys = append(keys, k)
		if v == 1.0 {
			nrOfSingletons += 1.0
		}
	}

	if nrOfSingletons == n {
		nrOfSingletons -= 1.0
	}

	for i := range histogram {
		p[i] = histogram[keys[i]] / n
	}

	C := 1.0 - nrOfSingletons/n

	for i := range p {
		p[i] = p[i] * C
	}

	var z float64
	var r float64

	for i := range p {
		if p[i] > 0.0 {
			z = math.Pow((1.0 - p[i]), n)
			z = (1.0 - z)
			r -= p[i] * math.Log(p[i]) / z
		}
	}

	return r
}
