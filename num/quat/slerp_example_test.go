// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quat_test

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/num/quat"
	"gonum.org/v1/gonum/spatial/r3"
)

// slerp returns the spherical interpolation between q0 and q1
// for t in [0,1].
func slerp(r0, r1 r3.Rotation, t float64) r3.Rotation {
	q0 := quat.Number(r0)
	q1 := quat.Number(r1)
	// Simo Särkkä "Notes on Quaternions", June 28, 2007
	q1 = quat.Mul(q1, quat.Inv(q0))
	return r3.Rotation(quat.Mul(quat.Exp(quat.Scale(t, quat.Log(q1))), q0))
}

// Spherically interpolate between two quaternions to obtain a rotation.
func Example_slerp() {
	const nsamples = 10
	noRot := r3.NewRotation(0, r3.Vec{})
	xRot := r3.NewRotation(math.Pi, r3.Vec{X: 1})
	v := r3.Vec{X: 1, Y: 1, Z: 1}
	for i := 0; i <= nsamples; i++ {
		t := float64(i) / nsamples
		interpRot := slerp(noRot, xRot, t)
		rotated := interpRot.Rotate(v)
		fmt.Printf("%.2f %.4g\n", t, rotated)
	}

	// Output:
	//
	// 0.00 {1 1 1}
	// 0.10 {1 0.642 1.26}
	// 0.20 {1 0.2212 1.397}
	// 0.30 {1 -0.2212 1.397}
	// 0.40 {1 -0.642 1.26}
	// 0.50 {1 -1 1}
	// 0.60 {1 -1.26 0.642}
	// 0.70 {1 -1.397 0.2212}
	// 0.80 {1 -1.397 -0.2212}
	// 0.90 {1 -1.26 -0.642}
	// 1.00 {1 -1 -1}
}
