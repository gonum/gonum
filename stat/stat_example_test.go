// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat_test

import (
	"fmt"
	"math/rand/v2"

	"gonum.org/v1/gonum/stat"
)

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
