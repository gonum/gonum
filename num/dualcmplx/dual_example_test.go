// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dualcmplx_test

import (
	"fmt"
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/num/dualcmplx"
)

// point is a 2-dimensional point/vector.
type point struct {
	x, y float64
}

// raiseDual raises the dimensionality of a point to a dual complex number.
func raiseDual(p point) dualcmplx.Number {
	return dualcmplx.Number{
		Real: 1,
		Dual: complex(p.x, p.y),
	}
}

// transform performs the transformation of p by the given set of dual
// complex transforms in order. The rotations are normalized to unit
// vectors.
func transform(p point, by ...dualcmplx.Number) point {
	if len(by) == 0 {
		return p
	}

	// Ensure the modulus of by is correctly scaled.
	for i := range by {
		if len := cmplx.Abs(by[i].Real); len != 1 {
			by[i].Real *= complex(1/len, 0)
		}
	}

	// Perform the transformations.
	z := by[0]
	for _, o := range by[1:] {
		z = dualcmplx.Mul(o, z)
	}
	pp := dualcmplx.Mul(z, raiseDual(p))

	// Extract the point.
	return point{x: real(pp.Dual), y: imag(pp.Dual)}
}

func Example() {
	// Translate a 1×1 square [3, 4] and rotate it 90° around the
	// origin.
	fmt.Println("square:")
	for i, p := range []point{
		{x: 0, y: 0},
		{x: 0, y: 1},
		{x: 1, y: 0},
		{x: 1, y: 1},
	} {
		pp := transform(p,
			// Displace.
			raiseDual(point{3, 4}),

			// Rotate.
			dualcmplx.Number{Real: complex(math.Cos(math.Pi/2), math.Sin(math.Pi/2))},
		)

		// Clean up floating point error for clarity.
		pp.x = floats.Round(pp.x, 2)
		pp.y = floats.Round(pp.y, 2)

		fmt.Printf(" %d %+v -> %+v\n", i, p, pp)
	}

	// Rotate a line segment 90° around its lower end.
	fmt.Println("\nline segment:")
	// Offset to origin from lower end.
	off := raiseDual(point{-2, -2})
	for i, p := range []point{
		{x: 2, y: 2},
		{x: 2, y: 3},
	} {
		pp := transform(p,
			// Shift origin.
			//
			// Complex number multiplication is commutative,
			// so the offset can be constructed as a single
			// dual complex number.
			dualcmplx.Mul(off, dualcmplx.ConjDual(dualcmplx.ConjCmplx(off))),

			// Rotate.
			dualcmplx.Number{Real: complex(math.Cos(math.Pi/2), math.Sin(math.Pi/2))},
		)

		// Clean up floating point error for clarity.
		pp.x = floats.Round(pp.x, 2)
		pp.y = floats.Round(pp.y, 2)

		fmt.Printf(" %d %+v -> %+v\n", i, p, pp)
	}

	// Output:
	//
	// square:
	//  0 {x:0 y:0} -> {x:-4 y:3}
	//  1 {x:0 y:1} -> {x:-5 y:3}
	//  2 {x:1 y:0} -> {x:-4 y:4}
	//  3 {x:1 y:1} -> {x:-5 y:4}
	//
	// line segment:
	//  0 {x:2 y:2} -> {x:2 y:2}
	//  1 {x:2 y:3} -> {x:1 y:2}
}
