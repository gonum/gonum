// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file

package complex_test

import (
	"fmt"

	. "gonum.org/v1/gonum/complex"
)

// Set of examples for all the functions

func ExampleAdd_simple() {
	// Adding three slices together. Note that
	// the result is stored in the first slice
	s1 := []complex128{1 + 2i, 2 + 1i, 3 - 4i, 4 - 3i}
	s2 := []complex128{5 + 6i, 6 - 5i, 7 + 8i, 8 - 7i}
	s3 := []complex128{1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i}
	Add(s1, s2)
	Add(s1, s3)

	fmt.Println("s1 =", s1)
	fmt.Println("s2 =", s2)
	fmt.Println("s3 =", s3)
	// Output:
	// s1 = [(7+9i) (9-3i) (11+5i) (13-9i)]
	// s2 = [(5+6i) (6-5i) (7+8i) (8-7i)]
	// s3 = [(1+1i) (1+1i) (1+1i) (1+1i)]
}

func ExampleAdd_newslice() {
	// If one wants to store the result in a
	// new container, just make a new slice
	s1 := []complex128{1 + 2i, 2 + 1i, 3 - 4i, 4 - 3i}
	s2 := []complex128{5 + 6i, 6 - 5i, 7 + 8i, 8 - 7i}
	s3 := []complex128{1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i}
	dst := make([]complex128, len(s1))

	AddTo(dst, s1, s2)
	Add(dst, s3)

	fmt.Println("dst =", dst)
	fmt.Println("s1 =", s1)
	fmt.Println("s2 =", s2)
	fmt.Println("s3 =", s3)
	// Output:
	// dst = [(7+9i) (9-3i) (11+5i) (13-9i)]
	// s1 = [(1+2i) (2+1i) (3-4i) (4-3i)]
	// s2 = [(5+6i) (6-5i) (7+8i) (8-7i)]
	// s3 = [(1+1i) (1+1i) (1+1i) (1+1i)]
}

func ExampleAdd_unequallengths() {
	// If the lengths of the slices are unknown,
	// use Eqlen to check
	s1 := []complex128{1 + 2i, 2 + 1i, 3 - 4i}
	s2 := []complex128{5 + 6i, 6 - 5i, 7 + 8i, 8 - 7i}

	eq := EqualLengths(s1, s2)
	if eq {
		Add(s1, s2)
	} else {
		fmt.Println("Unequal lengths")
	}
	// Output:
	// Unequal lengths
}

func ExampleAddConst() {
	s := []complex128{1 + 2i, -2 - 1i, 3 + 4i, -4 - 3i}
	c := 5 + 3i

	AddConst(c, s)

	fmt.Println("s =", s)
	// Output:
	// s = [(6+5i) (3+2i) (8+7i) (1+0i)]
}

func ExampleCumProd() {
	s := []complex128{1 + 2i, -2 - 1i, 3 + 4i, -4 - 3i}
	dst := make([]complex128, len(s))

	CumProd(dst, s)

	fmt.Println("dst =", dst)
	fmt.Println("s =", s)
	// Output:
	// dst = [(1+2i) (0-5i) (20-15i) (-125+0i)]
	// s = [(1+2i) (-2-1i) (3+4i) (-4-3i)]
}

func ExampleCumSum() {
	s := []complex128{1 + 2i, -2 - 1i, 3 + 4i, -4 - 3i}
	dst := make([]complex128, len(s))

	CumSum(dst, s)

	fmt.Println("dst =", dst)
	fmt.Println("s =", s)
	// Output:
	// dst = [(1+2i) (-1+1i) (2+5i) (-2+2i)]
	// s = [(1+2i) (-2-1i) (3+4i) (-4-3i)]
}
