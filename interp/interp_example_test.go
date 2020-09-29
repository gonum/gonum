// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interp_test

import (
	"fmt"
	"math"
	"os"
	"text/tabwriter"

	"gonum.org/v1/gonum/interp"
)

func ExamplePredictor() {
	// An example of fitting different interpolation
	// algorithms to (X, Y) data with widely varying slope.
	//
	// Cubic interpolators have to balance the smoothness
	// of the generated curve with suppressing ugly wiggles
	// (compare the output of AkimaSpline with that of
	// FritschButland).
	xs := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	ys := []float64{0, 0.001, 0.002, 0.1, 1, 2, 2.5, -10, -10.01, 2.49, 2.53, 2.55}

	var pc interp.PiecewiseConstant
	var pl interp.PiecewiseLinear
	var as interp.AkimaSpline
	var fb interp.FritschButland
	var cs interp.CubicSpline

	predictors := []interp.FittablePredictor{&pc, &pl, &as, &fb, &cs}
	for i, p := range predictors {
		err := p.Fit(xs, ys)
		if err != nil {
			panic(fmt.Sprintf("Error fitting %d-th predictor: %v", i, err))
		}
	}

	n := len(xs)
	dx := 0.25
	nPts := int(math.Round(float64(n-1)/dx)) + 1

	w := tabwriter.NewWriter(os.Stdout, 8, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "x\tPC\tPL\tAS\tFB\tCS\t")
	for i := 0; i < nPts; i++ {
		x := xs[0] + float64(i)*dx
		fmt.Fprintf(w, "%.2f", x)
		for _, predictor := range predictors {
			y := predictor.Predict(x)
			fmt.Fprintf(w, "\t%.2f", y)
		}
		fmt.Fprintln(w, "\t")
	}
	fmt.Fprintln(w)
	w.Flush()
	// Output:
	//        x      PC      PL      AS      FB      CS
	//     0.00    0.00    0.00    0.00    0.00    0.00
	//     0.25    0.00    0.00    0.00    0.00   -0.00
	//     0.50    0.00    0.00    0.00    0.00   -0.00
	//     0.75    0.00    0.00    0.00    0.00   -0.00
	//     1.00    0.00    0.00    0.00    0.00    0.00
	//     1.25    0.00    0.00    0.00    0.00    0.01
	//     1.50    0.00    0.00    0.00    0.00    0.02
	//     1.75    0.00    0.00    0.00    0.00    0.02
	//     2.00    0.00    0.00    0.00    0.00    0.00
	//     2.25    0.10    0.03   -0.01    0.01   -0.02
	//     2.50    0.10    0.05   -0.01    0.03   -0.04
	//     2.75    0.10    0.08    0.02    0.06   -0.01
	//     3.00    0.10    0.10    0.10    0.10    0.10
	//     3.25    1.00    0.33    0.26    0.22    0.30
	//     3.50    1.00    0.55    0.49    0.45    0.56
	//     3.75    1.00    0.78    0.75    0.73    0.81
	//     4.00    1.00    1.00    1.00    1.00    1.00
	//     4.25    2.00    1.25    1.24    1.26    1.11
	//     4.50    2.00    1.50    1.50    1.54    1.23
	//     4.75    2.00    1.75    1.75    1.79    1.48
	//     5.00    2.00    2.00    2.00    2.00    2.00
	//     5.25    2.50    2.12    2.22    2.17    2.80
	//     5.50    2.50    2.25    2.37    2.33    3.49
	//     5.75    2.50    2.38    2.47    2.45    3.56
	//     6.00    2.50    2.50    2.50    2.50    2.50
	//     6.25  -10.00   -0.62    0.83    0.55    0.01
	//     6.50  -10.00   -3.75   -2.98   -3.75   -3.38
	//     6.75  -10.00   -6.88   -7.18   -8.04   -6.96
	//     7.00  -10.00  -10.00  -10.00  -10.00  -10.00
	//     7.25  -10.01  -10.00  -11.16  -10.00  -11.89
	//     7.50  -10.01  -10.00  -11.55  -10.01  -12.52
	//     7.75  -10.01  -10.01  -11.18  -10.01  -11.89
	//     8.00  -10.01  -10.01  -10.01  -10.01  -10.01
	//     8.25    2.49   -6.88   -7.18   -8.06   -6.99
	//     8.50    2.49   -3.76   -2.99   -3.77   -3.43
	//     8.75    2.49   -0.63    0.82    0.53   -0.04
	//     9.00    2.49    2.49    2.49    2.49    2.49
	//     9.25    2.53    2.50    2.50    2.51    3.64
	//     9.50    2.53    2.51    2.51    2.52    3.70
	//     9.75    2.53    2.52    2.52    2.52    3.16
	//    10.00    2.53    2.53    2.53    2.53    2.53
	//    10.25    2.55    2.53    2.54    2.54    2.29
	//    10.50    2.55    2.54    2.54    2.54    2.94
	//    10.75    2.55    2.54    2.55    2.55    4.96
	//    11.00    2.55    2.55    2.55    2.55    2.55
}
