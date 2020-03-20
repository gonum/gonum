// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linsolve_test

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/linsolve"
	"gonum.org/v1/gonum/mat"
)

// System represents a linear system with a symmetric band matrix
//  A*x = b
type System struct {
	A *mat.SymBandDense
	B *mat.VecDense
}

// L2Projection returns a linear system whose solution is the L2 projection of f
// into the space of piecewise linear functions defined on the given grid.
//
// References:
//  - M. Larson, F. Bengzon, The Finite Element Method: Theory,
//    Implementations, and Applications. Springer (2013), Section 1.3, also
//    available at:
//    http://www.springer.com/cda/content/document/cda_downloaddocument/9783642332869-c1.pdf
func L2Projection(grid []float64, f func(float64) float64) System {
	n := len(grid)

	// Assemble the symmetric banded mass matrix by iterating over all elements.
	A := mat.NewSymBandDense(n, 1, nil)
	for i := 0; i < n-1; i++ {
		// h is the length of the i-th element.
		h := grid[i+1] - grid[i]
		// Add contribution from the i-th element.
		A.SetSymBand(i, i, A.At(i, i)+h/3)
		A.SetSymBand(i, i+1, h/6)
		A.SetSymBand(i+1, i+1, A.At(i+1, i+1)+h/3)
	}

	// Assemble the load vector by iterating over all elements.
	b := mat.NewVecDense(n, nil)
	for i := 0; i < n-1; i++ {
		h := grid[i+1] - grid[i]
		b.SetVec(i, b.AtVec(i)+f(grid[i])*h/2)
		b.SetVec(i+1, b.AtVec(i+1)+f(grid[i+1])*h/2)
	}

	return System{A, b}
}

func ExampleIterative() {
	const (
		n  = 10
		x0 = 0.0
		x1 = 1.0
	)
	// Make a uniform grid.
	grid := make([]float64, n+1)
	floats.Span(grid, x0, x1)
	sys := L2Projection(grid, func(x float64) float64 {
		return x * math.Sin(x)
	})

	result, err := linsolve.Iterative(sys.A, sys.B, &linsolve.CG{}, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("# iterations: %v\n", result.Stats.Iterations)
	fmt.Printf("Final solution: %.6f\n", mat.Formatted(result.X.T()))

	// Output:
	// # iterations: 11
	// Final solution: [-0.003339   0.006677   0.036530   0.085606   0.152981   0.237072   0.337006   0.447616   0.578244   0.682719   0.920847]
}
