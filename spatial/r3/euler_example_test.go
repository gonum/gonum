// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3_test

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/num/quat"
	"gonum.org/v1/gonum/spatial/r3"
)

// euler returns an r3.Rotation that corresponds to the Euler
// angles alpha, beta and gamma which are rotations around the x,
// y and z axes respectively. The order of rotations is x, y, z;
// there are many conventions for this ordering.
func euler(alpha, beta, gamma float64) r3.Rotation {
	// Note that this function can be algebraically simplified
	// to reduce floating point operations, but is left in this
	// form for clarity.
	var rot1, rot2, rot3 quat.Number
	rot1.Imag, rot1.Real = math.Sincos(alpha / 2) // x-axis rotation
	rot2.Jmag, rot2.Real = math.Sincos(beta / 2)  // y-axis rotation
	rot3.Kmag, rot3.Real = math.Sincos(gamma / 2) // z-axis rotation

	return r3.Rotation(quat.Mul(rot3, quat.Mul(rot2, rot1))) // order of rotations
}

func ExampleRotation_eulerAngles() {
	// It is possible to interconvert between the quaternion representation
	// of a rotation and Euler angles, but this leads to problems.
	//
	// The first of these is that there are a variety of conventions for
	// application of the rotations.
	//
	// The more serious consequence of using Euler angles is that it is
	// possible to put the rotation system into a singularity which results
	// in loss of degrees of freedom and so causes gimbal lock. This happens
	// when the second axis to be rotated around is rotated to 𝝿/2.
	//
	// See https://en.wikipedia.org/wiki/Euler_angles for more details.

	pt := r3.Vec{1, 0, 0}

	// For the Euler conversion function in this example, the second rotation
	// is around the y-axis.
	const singularY = math.Pi / 2

	arb := math.Pi / 4

	fmt.Printf("rotate around x-axis: %.2f\n", euler(arb, 0, 0).Rotate(pt))
	fmt.Printf("rotate around y-axis: %.2f\n", euler(0, arb, 0).Rotate(pt))
	fmt.Printf("rotate around z-axis: %.2f\n", euler(0, 0, arb).Rotate(pt))
	fmt.Printf("rotate around x+y-axes: %.2f\n", euler(arb, arb, 0).Rotate(pt))
	fmt.Printf("rotate around x+z-axes: %.2f\n", euler(arb, 0, arb).Rotate(pt))
	fmt.Printf("rotate around y+z-axes: %.2f\n", euler(0, arb, arb).Rotate(pt))

	fmt.Printf("rotate around y-axis to singularity: %.2f\n", euler(0, singularY, 0).Rotate(pt))
	fmt.Printf("rotate around x+y-axes with singularity → gimbal lock: %.2f\n", euler(arb, singularY, 0).Rotate(pt))
	fmt.Printf("rotate around z+y-axes with singularity → gimbal lock: %.2f\n", euler(0, singularY, arb).Rotate(pt))
	fmt.Printf("rotate around all-axes with singularity → gimbal lock: %.2f\n", euler(arb, singularY, arb).Rotate(pt))

	// Output:
	//
	// rotate around x-axis: {1.00 0.00 0.00}
	// rotate around y-axis: {0.71 0.00 -0.71}
	// rotate around z-axis: {0.71 0.71 0.00}
	// rotate around x+y-axes: {0.71 0.00 -0.71}
	// rotate around x+z-axes: {0.71 0.71 0.00}
	// rotate around y+z-axes: {0.50 0.50 -0.71}
	// rotate around y-axis to singularity: {0.00 0.00 -1.00}
	// rotate around x+y-axes with singularity → gimbal lock: {0.00 0.00 -1.00}
	// rotate around z+y-axes with singularity → gimbal lock: {0.00 0.00 -1.00}
	// rotate around all-axes with singularity → gimbal lock: {0.00 0.00 -1.00}
}
