package optimize_test

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/optimize"
)

func ExampleNelderMead() {
	p := optimize.Problem{
		Func: func(x []float64) float64 {
			return math.Pow(x[0]-(4.0/5.0), 2)
		},
	}

	withoutBounds := &optimize.NelderMead{}
	result1, _ := optimize.Minimize(p, []float64{2}, &optimize.Settings{}, withoutBounds)
	fmt.Printf("without bounds: %.1f\n", result1.X[0])

	withBounds := &optimize.NelderMead{
		Bounds: []optimize.Bound{
			{Min: 1, Max: 3},
		},
	}
	result2, _ := optimize.Minimize(p, []float64{2}, &optimize.Settings{}, withBounds)
	fmt.Printf("with bounds: %.1f\n", result2.X[0])

	// Output:
	// without bounds: 0.8
	// with bounds: 1.0
}
