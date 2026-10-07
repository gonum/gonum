package examples

import (
	"fmt"

	"gonum.org/v1/gonum/stat"
)

// LinearRegressionWithWeights demonstrates fitting a line through measured
// data points using stat.LinearRegression with inverse-variance weights.
//
// Assumptions:
//   - Each measurement y_i has independent Gaussian error with known variance σ_i².
//   - The model is y = a + b·x.
//   - Weights w_i = 1/σ_i² give the best linear unbiased estimator under these assumptions.
func LinearRegressionWithWeights() {
	// Synthetic measurements of a linear relationship with known noise levels.
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{2.1, 3.9, 6.2, 7.8, 10.1}
	// Known standard deviations of each measurement.
	sigma := []float64{0.1, 0.2, 0.15, 0.25, 0.1}

	// Compute inverse-variance weights.
	w := make([]float64, len(x))
	for i := range x {
		w[i] = 1.0 / (sigma[i] * sigma[i])
	}

	a, b := stat.LinearRegression(x, y, w, nil)

	fmt.Printf("intercept a = %.4f\n", a)
	fmt.Printf("slope     b = %.4f\n", b)

	// Compare with unweighted fit.
	aUnw, bUnw := stat.LinearRegression(x, y, nil, nil)
	fmt.Printf("unweighted intercept = %.4f\n", aUnw)
	fmt.Printf("unweighted slope     = %.4f\n", bUnw)

	// Predicted values using weighted fit.
	for i := range x {
		pred := a + b*x[i]
		resid := y[i] - pred
		fmt.Printf("  x=%.0f  y=%.2f  pred=%.2f  resid=%.2f\n",
			x[i], y[i], pred, resid)
	}
}
