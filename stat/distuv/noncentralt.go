// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"math/rand/v2"

	"gonum.org/v1/exp/root"
	"gonum.org/v1/gonum/mathext"
)

// NoncentralT is the noncentral t-distribution.
//
// See https://en.wikipedia.org/wiki/Noncentral_t-distribution for more details.
type NoncentralT struct {
	// Nu is the degrees of freedom.
	Nu float64

	// Ncp is the noncentral parameter.
	Ncp float64

	// Src is the random source used to generate samples.
	Src rand.Source
}

// Rand samples from the noncentral t-distribution.
func (dist NoncentralT) Rand() float64 {
	n := dist.Ncp + rand.New(dist.Src).NormFloat64()
	c2 := ChiSquared{K: dist.Nu, Src: dist.Src}.Rand()
	return n / math.Sqrt(c2/dist.Nu)
}

// Mean returns the mean of the noncentral t-distribution.
func (dist NoncentralT) Mean() float64 {
	nu := dist.Nu
	if nu <= 1 {
		return math.NaN()
	}
	return dist.Ncp * math.Sqrt(nu/2) * gammaADivB((nu-1)/2, nu/2)
}

// Variance returns the variance of the noncentral t-distribution.
func (dist NoncentralT) Variance() float64 {
	nu := dist.Nu
	if nu <= 2 {
		return math.NaN()
	}
	mean := dist.Mean()
	return nu*(1+dist.Ncp*dist.Ncp)/(nu-2) - mean*mean
}

// Prob returns the probability density function of the noncentral t-distribution.
func (dist NoncentralT) Prob(x float64) float64 {
	return math.Exp(dist.LogProb(x))
}

// LogProb returns the log-probability density function of the noncentral t-distribution.
func (dist NoncentralT) LogProb(x float64) float64 {
	var epsilon = math.Nextafter(1, 2) - 1
	ax := math.Abs(x)
	if ax > math.Sqrt(dist.Nu*epsilon) {
		a := NoncentralT{Nu: dist.Nu + 2, Ncp: dist.Ncp}.CDF(x * math.Sqrt(1+2/dist.Nu))
		b := dist.CDF(x)
		return math.Log(dist.Nu) - math.Log(ax) + math.Log(math.Abs(a-b))
	} else {
		return lgamma((dist.Nu+1)/2) - lgamma(dist.Nu/2) - 0.5*(math.Log(math.Pi)+math.Log(dist.Nu)+dist.Ncp*dist.Ncp)
	}
}

// CDF is the cumulative distribution function of the noncentral t-distribution.
// This implementation is based on:
// Russell Lenth, Cumulative Distribution Function of the Non-Central T Distribution, Algorithm AS 243.
func (dist NoncentralT) CDF(t float64) float64 {
	df, delta := dist.Nu, dist.Ncp

	const itrmax = 1000
	const errmax = 1e-12

	if df <= 0 {
		return math.NaN()
	}

	var negdel bool
	var del float64
	if t >= 0 {
		negdel, del = false, delta
	} else {
		negdel, del = true, -delta
	}

	// Initialize twin series.
	// Guenther, J. (1978). Statist. Computn. Simuln. vol.6, 199.
	x := t * t / (t*t + df)
	lambda := del * del
	p := 0.5 * math.Exp(-0.5*lambda)
	if p == 0 {
		// We overflowed, so use approximation from equation 26.7.10, Abramowitz & Stegun.
		x := (t*(1-1./(4*df)) - delta) / math.Sqrt(1+t*t/(2*df))
		return normCDF(x, 0, 1, !negdel)
	}
	q := math.Sqrt(2/math.Pi) * p * del
	s := 0.5 - p
	a := 0.5
	b := 0.5 * df
	rxb := math.Pow(1-x, b)
	albeta := math.Log(math.Sqrt(math.Pi)) + lgamma(b) - lgamma(0.5+b)
	xodd := mathext.RegIncBeta(a, b, x)
	godd := 2 * rxb * math.Exp(a*math.Log(x)-albeta)
	xeven := 1 - rxb
	geven := b * x * rxb
	tnc := p*xodd + q*xeven

	// Repeat until convergence.
	for en := 1; en <= itrmax; en++ {
		a += 1
		xodd -= godd
		xeven -= geven
		godd *= x * (a + b - 1) / a
		geven *= x * (a + b - 0.5) / (a + 0.5)
		p *= lambda / (2 * float64(en))
		q *= lambda / (2*float64(en) + 1)
		s -= p
		tnc += p*xodd + q*xeven
		errbd := 2 * s * (xodd - godd)
		if math.Abs(errbd) < errmax {
			break
		}

		if s < 0 { // loss of precision
			break
		}
	}

	tnc += normCDF(-del, 0, 1, true)

	if negdel {
		return 1 - tnc
	}
	return tnc
}

// Quantile is the quantile function.
func (dist NoncentralT) Quantile(p float64) float64 {
	if dist.Nu <= 0 {
		return math.NaN()
	}

	f := func(x float64) float64 { return dist.CDF(x) - p }

	// Find a, b where f(a)f(b) < 0.
	// Start the find by making a rough guess assuming a gaussian.
	var guess float64 = 1
	if dist.Nu > 3 {
		sigma := math.Sqrt(dist.Variance())
		guess = dist.Mean() + sigma*mathext.NormalQuantile(p)
	}
	a, b := findBracketMono(f, guess)

	t, err := root.Brent(f, a, b, 1e-13)
	if err != nil {
		return math.NaN()
	}

	return t
}

// findBracketMono finds a bracket interval [a, b] where f(a)f(b) < 0.
// f must be a monotonically increasing function.
func findBracketMono(f func(float64) float64, guess float64) (float64, float64) {
	// Make sure initial guess has the same sign as the root.
	f0 := f(0)
	if (guess < 0 && f0 < 0) || (guess > 0 && f0 > 0) {
		guess *= -1
	}

	// r is the rate in which we adjust the interval.
	var r float64
	a, fa := guess, f(guess)
	if (a > 0) == (fa < 0) {
		r = 2
	} else {
		r = 1. / 2
	}

	b := a * r
	fb := f(b)
	for range 200 {
		if math.Signbit(fa) != math.Signbit(fb) || fa == 0 || fb == 0 {
			break
		}
		a, fa = b, fb
		b *= r
		fb = f(b)
	}

	return a, b
}

func normCDF(x, mu, sigma float64, lowerTail bool) float64 {
	p := 0.5 * math.Erfc(-(x-mu)/(sigma*math.Sqrt2))
	if lowerTail {
		return p
	}
	return 1 - p
}

func lgamma(x float64) float64 {
	y, _ := math.Lgamma(x)
	return y
}

func gammaADivB(a, b float64) float64 {
	var ly float64
	var sign int = 1

	ga, sa := math.Lgamma(a)
	ly, sign = ly+ga, sign*sa

	gb, sb := math.Lgamma(b)
	ly, sign = ly-gb, sign*sb

	return float64(sign) * math.Exp(ly)
}
