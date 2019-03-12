// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize_test

import (
	"fmt"
	"log"
	"os"
	"time"

	"gonum.org/v1/gonum/optimize"
)

func ExampleBrent() {
	p := optimize.Problem{
		Func: func(xs []float64) float64 { return (xs[0] + 4) * (xs[0] + 4) },
	}

	printer := optimize.NewPrinter()
	printer.Writer = os.Stderr

	settings := &optimize.Settings{
		Recorder: printer,
		Runtime:  10 * time.Second,
	}

	// extracted from scipy.optimize.brent example:
	// >>> from scipy import optimize
	// >>> optimize.brent(lambda x: (x+4)**2, brack=(1,2))
	// -3.9999999999999982

	x := []float64{1}
	result, err := optimize.Minimize(p, x, settings, &optimize.Brent{
		Min: -10, // FIXME: should be +1, if we want exact scipy.optimize.brent example...
		Max: +20,
	})
	if err != nil {
		log.Fatal(err)
	}
	if err = result.Status.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("result.Status: %v\n", result.Status)
	fmt.Printf("result.X: %0.4g\n", result.X)
	fmt.Printf("result.F: %0.4g\n", result.F)
	fmt.Printf("result.Stats.FuncEvaluations: %d\n", result.Stats.FuncEvaluations)

	// Output:
	// result.Status: FunctionConvergence
	// result.X: [-4]
	// result.F: 0
	// result.Stats.FuncEvaluations: ???
}
