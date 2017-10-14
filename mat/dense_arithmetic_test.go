// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat_test

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func ExampleDense_Add() {
	// Initialize two matrices, a and b.
	a := mat.NewDense(2, 2, []float64{1, 0, 1, 0})
	b := mat.NewDense(2, 2, []float64{0, 1, 0, 1})

	// Add a and b, placing the result into c.
	// ...Notice that the size is automatically adjusted when the receiver has zero size.
	var c mat.Dense
	c.Add(a, b)

	// Print the result using the formatter.
	fc := mat.Formatted(&c, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("Result:\nc = %v\n\n", fc)
	// Output:
	// Result:
	// c = ⎡1  1⎤
	//     ⎣1  1⎦
	//
}

func ExampleDense_Sub() {
	// Initialize two matrices, a and b.
	a := mat.NewDense(2, 2, []float64{1, 1, 1, 1})
	b := mat.NewDense(2, 2, []float64{1, 0, 0, 1})

	// Subtract b from a, placing the result into a.
	a.Sub(a, b)

	// Print the result using the formatter.
	fa := mat.Formatted(a, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("Result:\na = %v\n\n", fa)
	// Output:
	// Result:
	// a = ⎡0  1⎤
	//     ⎣1  0⎦
	//
}

func ExampleDense_MulElem() {
	// Initialize two matrices, a and b.
	a := mat.NewDense(2, 2, []float64{1, 2, 3, 4})
	b := mat.NewDense(2, 2, []float64{1, 2, 3, 4})

	// Multiply the elements of a and b, placing the result into a.
	a.MulElem(a, b)

	// Print the result using the formatter.
	fa := mat.Formatted(a, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("Result:\na = %v\n\n", fa)
	// Output:
	// Result:
	// a = ⎡1   4⎤
	//     ⎣9  16⎦
	//
}

func ExampleDense_DivElem() {
	// Initialize two matrices, a and b.
	a := mat.NewDense(2, 2, []float64{5, 10, 15, 20})
	b := mat.NewDense(2, 2, []float64{5, 5, 5, 5})

	// Divide the elements of a by b, placing the result into a.
	a.DivElem(a, b)

	// Print the result using the formatter.
	fa := mat.Formatted(a, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("Result:\na = %v\n\n", fa)
	// Output:
	// Result:
	// a = ⎡1  2⎤
	//     ⎣3  4⎦
	//
}

func ExampleDense_Inverse() {
	// Initialize two matrices, a and ia.
	a := mat.NewDense(2, 2, []float64{4, 0, 0, 4})
	var ia mat.Dense

	// Take the inverse of a and place the result in ia.
	ia.Inverse(a)

	// Print the result using the formatter.
	fa := mat.Formatted(&ia, mat.Prefix("     "), mat.Squeeze())
	fmt.Printf("Result:\nia = %.2g\n\n", fa)

	// Show that [A] * [A]^(-1) = [I]
	var r mat.Dense
	r.Mul(a, &ia)
	fr := mat.Formatted(&r, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("Result:\nr = %v\n\n", fr)

	// An alternative to using Inverse is to solve [A] * [X] = [I].
	// ...Note that matrix inversion should generally be accomplished in this way, rather than using Inverse.
	i := mat.NewDense(2, 2, []float64{1, 0, 0, 1})
	var x mat.Dense
	x.Solve(a, i)

	// Print the result using the formatter.
	fx := mat.Formatted(&x, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("Result:\nx = %v\n\n", fx)

	// Output:
	// Result:
	// ia = ⎡0.25    -0⎤
	//      ⎣   0  0.25⎦
	//
	// Result:
	// r = ⎡1  0⎤
	//     ⎣0  1⎦
	//
	// Result:
	// x = ⎡0.25     0⎤
	//     ⎣   0  0.25⎦
	//
}
