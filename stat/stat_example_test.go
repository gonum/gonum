// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat_test

import (
	"fmt"
	"math"
	"math/rand/v2"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

func ExampleBivariateMoment() {
	// The bivariate central moment
	//
	//	μ_rs = E[(x-μ_x)^r*(y-μ_y)^s]
	//
	// generalises the covariance, which is the (1,1) moment. The (2,2)
	// moment gives the leading term of the sampling variance of a
	// covariance estimate,
	//
	//	Var(cov) ≈ (μ_22 - μ_11^2)/n
	//
	// so BivariateMoment can be used to attach a standard error to a
	// covariance. See https://stats.stackexchange.com/q/48366.
	const n = 1000

	xs := make([]float64, n)
	ys := make([]float64, n)

	// Draw ys with a fixed linear dependence on xs so that the covariance
	// is known to be non-zero.
	rnd := rand.New(rand.NewPCG(1, 1))
	for i := range xs {
		xs[i] = rnd.NormFloat64()
		ys[i] = 0.5*xs[i] + rnd.NormFloat64()
	}

	// BivariateMoment applies no degrees of freedom correction, so the
	// (1,1) moment divides by n where Covariance divides by n-1.
	m11 := stat.BivariateMoment(1, 1, xs, ys, nil)
	m22 := stat.BivariateMoment(2, 2, xs, ys, nil)
	cov := stat.Covariance(xs, ys, nil)

	stdErr := math.Sqrt((m22 - m11*m11) / n)

	fmt.Printf("μ_11     = %.6f\n", m11)
	fmt.Printf("cov      = %.6f\n", cov)
	fmt.Printf("μ_22     = %.6f\n", m22)
	fmt.Printf("std err  = %.6f\n", stdErr)

	// Output:
	// μ_11     = 0.490614
	// cov      = 0.491105
	// μ_22     = 1.685326
	// std err  = 0.038008
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

	rnd := rand.New(rand.NewPCG(1, 1))
	for i := range xs {
		xs[i] = float64(i)
		ys[i] = line(xs[i]) + 0.1*rnd.NormFloat64()
	}

	// Do not force the regression line to pass through the origin.
	origin := false

	alpha, beta := stat.LinearRegression(xs, ys, weights, origin)
	r2 := stat.RSquared(xs, ys, weights, alpha, beta)

	fmt.Printf("Estimated offset is: %.6f\n", alpha)
	fmt.Printf("Estimated slope is:  %.6f\n", beta)
	fmt.Printf("R^2: %.6f\n", r2)

	// Output:
	// Estimated offset is: 0.999675
	// Estimated slope is:  2.999971
	// R^2: 0.999999
}

func Example_confidenceInterval() {
	// This example shows how one can compute a confidence interval to quantify
	// the uncertainty around an estimated parameter, when working with a small
	// sample from a normally distributed variable.
	//
	// For small samples (N ≤ 30), confidence intervals are computed with
	// the t-distribution:
	//
	//	Conf.Interval = $\hat{x} ± t \frac{s}{\sqrt{n}}
	//
	// where:
	//   - x is the sample mean,
	//   - s is the sample standard deviation,
	//   - n is the sample size, and
	//   - t is the critical value from the t-distribution based on the desired
	//     confidence level and degrees of freedom (df=n-1)
	//
	// For more details, see:
	//
	//	https://en.wikipedia.org/wiki/Student's_t-distribution

	var (
		// First 10 sepal widths from iris data set.
		xs = []float64{3.5, 3.0, 3.2, 3.1, 3.6, 3.9, 3.4, 3.4, 2.9, 3.1}
		ws = []float64(nil) // weights
		n  = float64(len(xs))
		df = n - 1

		µ, std = stat.MeanStdDev(xs, ws)

		lvl = 0.95 // 95% confidence level
		t   = distuv.StudentsT{
			Mu:    0,
			Sigma: 1,
			Nu:    df,
		}.Quantile(0.5 * (1 + lvl))

		err = t * std / math.Sqrt(n)

		lo = µ - err
		hi = µ + err
	)

	fmt.Printf("Mean:     %2.2f\n", µ)
	fmt.Printf("CI(@%g%%): [%.2f, %.2f]\n", lvl*100, lo, hi)

	// Output:
	// Mean:     3.31
	// CI(@95%): [3.09, 3.53]
}
