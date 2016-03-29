// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fd_test

import (
	"fmt"
	"math"

	"github.com/gonum/diff/fd"
)

func ExampleDerivative() {
	f := func(x float64) float64 {
		return math.Sin(x)
	}
	// Compute the first derivative of f at 0 using the default settings.
	fmt.Println("f'(0) ≈", fd.Derivative(f, 0, nil))
	// Compute the first derivative of f at 0 using the forward approximation
	// with a custom step size.
	df := fd.Derivative(f, 0, &fd.Settings{
		Formula: fd.Forward,
		Step:    1e-3,
	})
	fmt.Println("f'(0) ≈", df)

	f = func(x float64) float64 {
		return math.Pow(math.Cos(x), 3)
	}
	// Compute the second derivative of f at 0 using
	// the centered approximation, concurrent evaluation,
	// and a known function value at x.
	df = fd.Derivative(f, 0, &fd.Settings{
		Formula:     fd.Central2nd,
		Concurrent:  true,
		OriginKnown: true,
		OriginValue: f(0),
	})
	fmt.Println("f''(0) ≈", df)

	// Output:
	// f'(0) ≈ 1
	// f'(0) ≈ 0.9999998333333416
	// f''(0) ≈ -2.999999981767587
}
