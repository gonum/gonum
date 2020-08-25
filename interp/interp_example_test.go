// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interp_test

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/interp"
)

func ExampleFit() {
	xs := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	ys := []float64{0, 0.001, 0.002, 0.1, 1, 2, 2.5, -10, -10.01, 2.49, 2.53, 2.55}

	var pc interp.PiecewiseConstant
	var pl interp.PiecewiseLinear
	var as interp.AkimaSpline
	var fb interp.FritschButland

	pc.Fit(xs, ys)
	pl.Fit(xs, ys)
	as.Fit(xs, ys)
	fb.Fit(xs, ys)

	predictors := []interp.Predictor{&pc, &pl, &as, &fb}

	n := len(xs)
	dx := 0.25
	nPts := int(math.Round(float64(n-1)/dx)) + 1

	fmt.Printf("x\t\tPC\t\tPL\t\tAS\t\tFB\n")
	for i := 0; i < nPts; i++ {
		x := xs[0] + float64(i)*dx
		fmt.Printf("%.2f", x)
		for _, predictor := range predictors {
			y := predictor.Predict(x)
			fmt.Printf("\t%.3f", y)
		}
		fmt.Printf("\n")
	}
	// Output:
	// x		PC		PL		AS		FB
	// 0.00	0.000	0.000	0.000	0.000
	// 0.25	0.001	0.000	0.000	0.000
	// 0.50	0.001	0.001	0.001	0.001
	// 0.75	0.001	0.001	0.001	0.001
	// 1.00	0.001	0.001	0.001	0.001
	// 1.25	0.002	0.001	0.001	0.001
	// 1.50	0.002	0.002	0.002	0.001
	// 1.75	0.002	0.002	0.002	0.002
	// 2.00	0.002	0.002	0.002	0.002
	// 2.25	0.100	0.027	-0.006	0.009
	// 2.50	0.100	0.051	-0.010	0.029
	// 2.75	0.100	0.076	0.015	0.060
	// 3.00	0.100	0.100	0.100	0.100
	// 3.25	1.000	0.325	0.265	0.221
	// 3.50	1.000	0.550	0.491	0.454
	// 3.75	1.000	0.775	0.747	0.734
	// 4.00	1.000	1.000	1.000	1.000
	// 4.25	2.000	1.250	1.245	1.258
	// 4.50	2.000	1.500	1.496	1.535
	// 4.75	2.000	1.750	1.749	1.794
	// 5.00	2.000	2.000	2.000	2.000
	// 5.25	2.500	2.125	2.218	2.172
	// 5.50	2.500	2.250	2.375	2.333
	// 5.75	2.500	2.375	2.469	2.453
	// 6.00	2.500	2.500	2.500	2.500
	// 6.25	-10.000	-0.625	0.834	0.548
	// 6.50	-10.000	-3.750	-2.983	-3.748
	// 6.75	-10.000	-6.875	-7.184	-8.044
	// 7.00	-10.000	-10.000	-10.000	-10.000
	// 7.25	-10.010	-10.002	-11.157	-10.004
	// 7.50	-10.010	-10.005	-11.553	-10.007
	// 7.75	-10.010	-10.008	-11.175	-10.009
	// 8.00	-10.010	-10.010	-10.010	-10.010
	// 8.25	2.490	-6.885	-7.180	-8.061
	// 8.50	2.490	-3.760	-2.986	-3.770
	// 8.75	2.490	-0.635	0.822	0.526
	// 9.00	2.490	2.490	2.490	2.490
	// 9.25	2.530	2.500	2.504	2.506
	// 9.50	2.530	2.510	2.515	2.517
	// 9.75	2.530	2.520	2.524	2.524
	// 10.00	2.530	2.530	2.530	2.530
	// 10.25	2.550	2.535	2.535	2.536
	// 10.50	2.550	2.540	2.541	2.542
	// 10.75	2.550	2.545	2.546	2.547
	// 11.00	2.550	2.550	2.550	2.550
}
