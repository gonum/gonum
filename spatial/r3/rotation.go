// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import (
	"math"

	"gonum.org/v1/gonum/num/quat"
)

// TODO: possibly useful additions to the current rotation API:
//  - create rotations from Euler angles (NewRotationFromEuler?)
//  - create rotations from rotation matrices (NewRotationFromMatrix?)
//  - return the equivalent Euler angles from a Rotation
//
// Euler angles have issues (see [1] for a discussion).
// We should think carefully before adding them in.
// [1]: http://www.euclideanspace.com/maths/geometry/rotations/conversions/quaternionToEuler/

// Rotation describes a rotation in space.
type Rotation quat.Number

// NewRotation creates a rotation by alpha, around axis.
func NewRotation(alpha float64, axis Vec) Rotation {
	if alpha == 0 {
		return Rotation{Real: 1}
	}
	q := raise(axis)
	sin, cos := math.Sincos(0.5 * alpha)
	q = quat.Scale(sin/quat.Abs(q), q)
	q.Real += cos
	if len := quat.Abs(q); len != 1 {
		q = quat.Scale(1/len, q)
	}

	return Rotation(q)
}

// Rotate returns p rotated according to the parameters used to construct
// the receiver.
func (r Rotation) Rotate(p Vec) Vec {
	if r.isIdentity() {
		return p
	}
	qq := quat.Number(r)
	pp := quat.Mul(quat.Mul(qq, raise(p)), quat.Conj(qq))
	return Vec{X: pp.Imag, Y: pp.Jmag, Z: pp.Kmag}
}

func (r Rotation) isIdentity() bool {
	return r == Rotation{Real: 1}
}

func raise(p Vec) quat.Number {
	return quat.Number{Imag: p.X, Jmag: p.Y, Kmag: p.Z}
}
