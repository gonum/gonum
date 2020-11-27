// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import (
	"math"

	"gonum.org/v1/gonum/num/quat"
)

// Rotation describes a rotation in space.
type Rotation struct {
	q quat.Number
}

// NewRotation creates a rotation by alpha, around axis.
func NewRotation(alpha float64, axis Vec) Rotation {
	var (
		q        = raise(axis)
		sin, cos = math.Sincos(0.5 * alpha)
	)
	q = quat.Scale(sin/quat.Abs(q), q)
	q.Real += cos
	if len := quat.Abs(q); len != 1 {
		q = quat.Scale(1/len, q)
	}

	return Rotation{q: q}
}

// Rotate returns the rotated vector according to the definition of rot.
func (rot Rotation) Rotate(p Vec) Vec {
	pp := quat.Mul(quat.Mul(rot.q, raise(p)), quat.Conj(rot.q))
	return Vec{X: pp.Imag, Y: pp.Jmag, Z: pp.Kmag}
}
