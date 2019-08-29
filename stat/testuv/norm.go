// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testuv

import (
	"math"

	"gonum.org/v1/gonum/stat"
)

// NormalSkew test whether the skew is different from the normal distribution.
// This function tests the null hypothesis that the skewness of population that the
// sample was drawn from is the same as that of a corresponding normal distribution.
// The length of a must be at least 8.
func NormalSkew(a []float64) float64 {
	// Implementation based on scipy's: https://github.com/scipy/scipy/blob/v1.3.0/scipy/stats/stats.py#L1307-L1382
	b2 := stat.Skew(a, nil)
	if len(a) < 8 {
		panic("testuv: not enough samples for normal skew test")
	}
	n := float64(len(a))
	y := b2 * math.Sqrt(((n+1)*(n+3))/(6*(n-2)))
	beta2 := (3 * (math.Pow(n, 2) + 27*n - 70) * (n + 1) * (n + 3) / ((n - 2) * (n + 5) * (n + 7) * (n + 9)))
	w2 := -1 + math.Sqrt(2*(beta2-1))
	delta := 1 / math.Sqrt(0.5*math.Log(w2))
	alpha := math.Sqrt(2 / (w2 - 1))
	if y == 0 {
		y = 1
	}
	return delta * math.Log(y/alpha+math.Sqrt(math.Pow(y/alpha, 2)+1))
}

// NormalKurtosis test whether a dataset has normal kurtosis.
// This function tests the null hypothesis that the kurtosis of the population from
// which the sample was drawn is that of the normal distribution:
// `kurtosis = 3(n-1)/(n+1)`.
// The length of a must be at least 5.
func NormalKurtosis(a []float64) float64 {
	// Implementation based on scipy's: https://github.com/scipy/scipy/blob/v1.3.0/scipy/stats/stats.py#L1388-L1475
	if len(a) < 5 {
		panic("testuv: not enough samples for normal kurtosis test")
	}
	n := float64(len(a))
	b2 := stat.ExKurtosis(a, nil) + 3
	e := 3.0 * (n - 1) / (n + 1)
	varb2 := 24 * n * (n - 2) * (n - 3) / ((n + 1) * (n + 1.) * (n + 3) * (n + 5))
	x := (b2 - e) / math.Sqrt(varb2)
	sqrtbeta1 := 6 * (n*n - 5*n + 2) / ((n + 7) * (n + 9)) * math.Sqrt((6*(n+3)*(n+5))/(n*(n-2)*(n-3)))
	a2 := 6 + 8/sqrtbeta1*(2/sqrtbeta1+math.Sqrt(1+4/math.Pow(sqrtbeta1, 2)))
	term1 := 1 - 2/(9*a2)
	denom := 1 + x*math.Sqrt(2/(a2-4))
	if denom == 0 {
		return math.NaN()
	}
	term2 := math.Copysign(1, denom) * math.Pow((1-2/a2)/math.Abs(denom), 1/3.0)
	return (term1 - term2) / math.Sqrt(2/(9*a2))
}

// Normal test whether a sample differs from a normal distribution.
// This function tests the null hypothesis that a sample comes from a normal distribution.
// It is based on D'Agostino and Pearson's [1]_, [2]_ test that combines skew and kurtosis to
// produce an omnibus test of normality.
// The length of a must be at least 8.
//
// References
// ----------
// .. [1] D'Agostino, R. B. (1971), "An omnibus test of normality for
//        moderate and large sample size", Biometrika, 58, 341-348
// .. [2] D'Agostino, R. and Pearson, E. S. (1973), "s for departure from
//        normality", Biometrika, 60, 613-622
func Normal(a []float64) float64 {
	// Implementation based on scipy's: https://github.com/scipy/scipy/blob/v1.3.0/scipy/stats/stats.py#L1481-L1549
	s := NormalSkew(a)
	k := NormalKurtosis(a)
	return s*s + k*k
}
