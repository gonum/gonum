// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3_test

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/num/quat"
	"gonum.org/v1/gonum/spatial/r3"
)

// slerp returns the spherical interpolation between q0 and q1
// for t in [0,1]; 0 corresponds to q0 and 1 corresponds to q1.
func slerp(r0, r1 r3.Rotation, t float64) r3.Rotation {
	q0 := quat.Number(r0)
	q1 := quat.Number(r1)
	// Based on Simo Särkkä "Notes on Quaternions" Eq. 35
	//  p(t) = (q1 ∗ q0^−1) ^ t ∗ q0
	// https://users.aalto.fi/~ssarkka/pub/quat.pdf
	q1 = quat.Mul(q1, quat.Inv(q0))
	q1 = quat.PowReal(q1, t)
	return r3.Rotation(quat.Mul(q1, q0))
}

// Spherically interpolate between two quaternions to obtain a rotation.
func Example_slerp() {
	const steps = 10
	// An initial rotation of pi/4 around the x-axis (45 degrees).
	initialRot := r3.NewRotation(math.Pi/4, r3.Vec{X: 1})
	// Final rotation is pi around the x-axis (180 degrees).
	finalRot := r3.NewRotation(math.Pi, r3.Vec{X: 1})
	// The vector we are rotating is (1, 1, 1).
	// The result should then be (1, -1, -1) when t=1 (finalRot) since we invert the y and z axes.
	v := r3.Vec{X: 1, Y: 1, Z: 1}
	for i := 0.0; i <= steps; i++ {
		t := i / steps
		rotated := slerp(initialRot, finalRot, t).Rotate(v)
		fmt.Printf("%.2f %.4g\n", t, rotated)
	}

	// Output:
	//
	// 0.00 {1 -1.11e-16 1.414}
	// 0.10 {1 -0.3301 1.375}
	// 0.20 {1 -0.642 1.26}
	// 0.30 {1 -0.9185 1.075}
	// 0.40 {1 -1.144 0.8313}
	// 0.50 {1 -1.307 0.5412}
	// 0.60 {1 -1.397 0.2212}
	// 0.70 {1 -1.41 -0.111}
	// 0.80 {1 -1.345 -0.437}
	// 0.90 {1 -1.206 -0.7389}
	// 1.00 {1 -1 -1}
}
