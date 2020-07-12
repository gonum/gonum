// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file

package cmplxs_test

import (
	"fmt"

	"gonum.org/v1/gonum/cmplxs"
)

// Set of examples for all the functions

func ExampleAdd_simple() {
	// Adding three slices together. Note that
	// the result is stored in the first slice
	s1 := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i}
	s2 := []complex128{5 + 7i, 6 + 7i, 7 + 7i, 8 + 8i}
	s3 := []complex128{1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i}
	cmplxs.Add(s1, s2)
	cmplxs.Add(s1, s3)

	fmt.Println("s1 =", s1)
	fmt.Println("s2 =", s2)
	fmt.Println("s3 =", s3)

	// Output:
	// s1 = [(7+10i) (9+11i) (11+12i) (13+14i)]
	// s2 = [(5+7i) (6+7i) (7+7i) (8+8i)]
	// s3 = [(1+2i) (1+2i) (1+2i) (1+2i)]
}

func ExampleAdd_newslice() {
	// If one wants to store the result in a
	// new container, just make a new slice
	s1 := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i}
	s2 := []complex128{5 + 7i, 6 + 7i, 7 + 7i, 8 + 8i}
	s3 := []complex128{1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i}
	dst := make([]complex128, len(s1))

	cmplxs.AddTo(dst, s1, s2)
	cmplxs.Add(dst, s3)

	fmt.Println("dst =", dst)
	fmt.Println("s1 =", s1)
	fmt.Println("s2 =", s2)
	fmt.Println("s3 =", s3)

	// Output:
	// dst = [(7+10i) (9+11i) (11+12i) (13+14i)]
	// s1 = [(1+1i) (2+2i) (3+3i) (4+4i)]
	// s2 = [(5+7i) (6+7i) (7+7i) (8+8i)]
	// s3 = [(1+2i) (1+2i) (1+2i) (1+2i)]
}

func ExampleAdd_unequallengths() {
	// If the lengths of the slices are unknown,
	// use EqualLengths to check
	s1 := []complex128{1 + 1i, 2 + 2i, 3 + 3i}
	s2 := []complex128{5 + 5i, 6 + 6i, 7 + 7i, 8 + 8i}

	eq := cmplxs.EqualLengths(s1, s2)
	if eq {
		cmplxs.Add(s1, s2)
	} else {
		fmt.Println("Unequal lengths")
	}

	// Output:
	// Unequal lengths
}

func ExampleAddConst() {
	s := []complex128{1 - 1i, -2 - 1i, 3 - 1i, -4 - 1i}
	c := 5 + 1i

	cmplxs.AddConst(c, s)

	fmt.Println("s =", s)

	// Output:
	// s = [(6+0i) (3+0i) (8+0i) (1+0i)]
}

func ExampleCumProd() {
	s := []complex128{1 + 1i, -2 - 2i, 3 + 3i, -4 - 4i}
	dst := make([]complex128, len(s))

	cmplxs.CumProd(dst, s)

	fmt.Println("dst =", dst)
	fmt.Println("s =", s)

	// Output:
	// dst = [(1+1i) (0-4i) (12-12i) (-96+0i)]
	// s = [(1+1i) (-2-2i) (3+3i) (-4-4i)]
}

func ExampleCumSum() {
	s := []complex128{1 + 1i, -2 - 2i, 3 + 3i, -4 - 4i}
	dst := make([]complex128, len(s))

	cmplxs.CumSum(dst, s)

	fmt.Println("dst =", dst)
	fmt.Println("s =", s)

	// Output:
	// dst = [(1+1i) (-1-1i) (2+2i) (-2-2i)]
	// s = [(1+1i) (-2-2i) (3+3i) (-4-4i)]
}
