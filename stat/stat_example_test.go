// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat_test

import (
	"fmt"
	"math"
	"sort"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

func ExampleKolmogorovSmirnov() {
	// Given a set of observations xobs, we want to test whether
	// they come from a given distribution (here, the normal distribution).

	xobs := []float64{
		-0.16, -0.68, -0.32, -0.85, 0.89, -2.28, 0.63, 0.41, 0.15, 0.74,
		1.30, -0.13, 0.80, -0.75, 0.28, -1.00, 0.14, -1.38, -0.04, -0.25,
		-0.17, 1.29, 0.47, -1.23, 0.21, -0.04, 0.07, -0.08, 0.32, -0.17,
		0.13, -1.94, 0.78, 0.19, -0.12, -0.19, 0.76, -1.48, -0.01, 0.20,
		-1.97, -0.37, 3.08, -0.40, 0.80, 0.01, 1.32, -0.47, 2.29, -0.26,
		-1.52, -0.06, -1.02, 1.06, 0.60, 1.15, 1.92, -0.06, -0.19, 0.67,
		0.29, 0.58, 0.02, 2.18, -0.04, -0.13, -0.79, -1.28, -1.41, -0.23,
		0.65, -0.26, -0.17, -1.53, -1.69, -1.60, 0.09, -1.11, 0.30, 0.71,
		-0.88, -0.03, 0.56, -3.68, 2.40, 0.62, 0.52, -1.25, 0.85, -0.09,
		-0.23, -1.16, 0.22, -1.68, 0.50, -0.35, -0.35, -0.33, -0.24, 0.25,
	}

	sort.Float64s(xobs)

	var (
		n    = len(xobs)
		inv  = 1.0 / float64(n)
		exp  = make([]float64, n)
		xedf = make([]float64, n)
		diff = make([]float64, n)
		norm = distuv.Normal{Mu: 0, Sigma: 1}
	)

	for i, v := range xobs {
		xedf[i] = inv * float64(i)
		exp[i] = norm.CDF(v)
		diff[i] = math.Abs(xedf[i] - exp[i])
	}

	dist := stat.KolmogorovSmirnov(xobs, xedf, xobs, exp)
	fmt.Printf("stat:     %d\n", n)
	fmt.Printf("diff:     %2.3f\n", floats.Max(diff))
	fmt.Printf("dist-KS:  %2.3f\n", dist)

	// critical value at the 95% confidence level
	crit := 1.36 / math.Sqrt(float64(n))
	fmt.Printf("dist(3σ): %2.3f\n", crit)

	h0 := "reject"
	if dist < crit {
		h0 = "do not reject"
	}
	fmt.Printf("H0:       %s\n", h0)

	// Output:
	// stat:     100
	// diff:     0.096
	// dist-KS:  0.027
	// dist(3σ): 0.136
	// H0:       do not reject
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
