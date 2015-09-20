// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64_test

import (
	"fmt"

	"github.com/gonum/matrix/mat64"
)

func ExampleCholesky() {
	// Construct a new SymDense. Only the upper triangular data
	// elements are used, so the lower triangular elements may
	// be set zero without altering the semantics.
	a := mat64.NewSymDense(4, []float64{
		108, -21, -9, -96,
		-21, 7, 5, 15,
		-9, 5, 61, 25,
		-96, 15, 25, 142,
	})

	fmt.Printf("a = %0.4v\n", mat64.Formatted(a, mat64.Prefix("    ")))

	// Compute the cholesky factorization.
	var chol mat64.Cholesky
	if ok := chol.Factorize(a); !ok {
		fmt.Println("a matrix is not positive semi-definite.")
	}

	// Find the determinant.
	fmt.Printf("\nThe determinant of a is %0.4g\n\n", chol.Det())

	// Use the factorization to solve the system of equations a * x = b.
	b := mat64.NewVector(4, []float64{1, 2, 3, 4})
	var x mat64.Vector
	if err := x.SolveCholeskyVec(&chol, b); err != nil {
		fmt.Println("Matrix is near singular: ", err)
	}
	fmt.Println("Solve a * x = b")
	fmt.Printf("x = %0.4v\n", mat64.Formatted(&x, mat64.Prefix("    ")))

	// Extract the factorization and check that it equals the original matrix.
	var t mat64.TriDense
	t.LFromCholesky(&chol)
	var test mat64.Dense
	test.Mul(&t, t.T())
	fmt.Println()
	fmt.Printf("L * L^T = %0.4v\n", mat64.Formatted(a, mat64.Prefix("          ")))

	// Output:
	// a = ⎡108  -21   -9  -96⎤
	//     ⎢-21    7    5   15⎥
	//     ⎢ -9    5   61   25⎥
	//     ⎣-96   15   25  142⎦
	//
	// The determinant of a is 7.885e+05
	//
	// Solve a * x = b
	// x = ⎡  0.3524⎤
	//     ⎢    1.02⎥
	//     ⎢-0.05115⎥
	//     ⎣  0.1676⎦
	//
	// L * L^T = ⎡108  -21   -9  -96⎤
	//           ⎢-21    7    5   15⎥
	//           ⎢ -9    5   61   25⎥
	//           ⎣-96   15   25  142⎦
}
