// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import (
	"math"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/num/quat"
)

// Vec is a 3D vector.
type Vec struct {
	X, Y, Z float64
}

// Add returns the vector sum of p and q.
func Add(p, q Vec) Vec {
	return Vec{
		X: p.X + q.X,
		Y: p.Y + q.Y,
		Z: p.Z + q.Z,
	}
}

// Sub returns the vector sum of p and -q.
func Sub(p, q Vec) Vec {
	return Vec{
		X: p.X - q.X,
		Y: p.Y - q.Y,
		Z: p.Z - q.Z,
	}
}

// Scale returns the vector p scaled by f.
func Scale(f float64, p Vec) Vec {
	return Vec{
		X: f * p.X,
		Y: f * p.Y,
		Z: f * p.Z,
	}
}

// Dot returns the dot product p·q.
func Dot(p, q Vec) float64 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z
}

// Cross returns the cross product p×q.
func Cross(p, q Vec) Vec {
	return Vec{
		p.Y*q.Z - p.Z*q.Y,
		p.Z*q.X - p.X*q.Z,
		p.X*q.Y - p.Y*q.X,
	}
}

// Rotate returns a new vector, rotated by alpha around the provided axis.
func Rotate(p Vec, alpha float64, axis Vec) Vec {
	return NewRotation(alpha, axis).Rotate(p)
}

// Norm returns the Euclidean norm of p
//  |p| = sqrt(p_x^2 + p_y^2 + p_z^2).
func Norm(p Vec) float64 {
	return math.Hypot(p.X, math.Hypot(p.Y, p.Z))
}

// Norm2 returns the Euclidean squared norm of p
//  |p|^2 = p_x^2 + p_y^2 + p_z^2.
func Norm2(p Vec) float64 {
	return p.X*p.X + p.Y*p.Y + p.Z*p.Z
}

// Unit returns the unit vector colinear to p.
// Unit returns {NaN,NaN,NaN} for the zero vector.
func Unit(p Vec) Vec {
	if p.X == 0 && p.Y == 0 && p.Z == 0 {
		return Vec{X: math.NaN(), Y: math.NaN(), Z: math.NaN()}
	}
	return Scale(1/Norm(p), p)
}

// Cos returns the cosine of the opening angle between p and q.
func Cos(p, q Vec) float64 {
	return Dot(p, q) / (Norm(p) * Norm(q))
}

// Box is a 3D bounding box.
type Box struct {
	Min, Max Vec
}

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

// Matrix returns a 3×3 rotation matrix corresponding to the receiver. It
// may be used to perform rotations on a 3-vector or to apply the rotation
// to a 3×n matrix of column vectors. If the receiver is not a unit
// quaternion, the returned matrix will not be a pure rotation.
func (r Rotation) Matrix() mat.Matrix {
	re, im, jm, km := r.Real, r.Imag, r.Jmag, r.Kmag
	im2 := im * im
	jm2 := jm * jm
	km2 := km * km
	rim := re * im
	rjm := re * jm
	rkm := re * km
	ijm := im * jm
	jkm := jm * km
	kim := km * im
	return &matrix{
		1 - 2*(jm2+km2), 2 * (ijm - rkm), 2 * (kim + rjm),
		2 * (ijm + rkm), 1 - 2*(im2+km2), 2 * (jkm - rim),
		2 * (kim - rjm), 2 * (jkm + rim), 1 - 2*(im2+jm2),
	}
}

// matrix is a 3×3 rotation matrix.
type matrix [9]float64

var (
	_ mat.Matrix      = (*matrix)(nil)
	_ mat.RawMatrixer = (*matrix)(nil)
)

func (m *matrix) At(i, j int) float64 {
	if uint(i) >= 3 {
		panic(mat.ErrRowAccess)
	}
	if uint(j) >= 3 {
		panic(mat.ErrColAccess)
	}
	return m[i*3+j]
}
func (m *matrix) Dims() (r, c int) { return 3, 3 }
func (m *matrix) T() mat.Matrix    { return mat.Transpose{Matrix: m} }
func (m *matrix) RawMatrix() blas64.General {
	return blas64.General{Rows: 3, Cols: 3, Data: m[:], Stride: 3}
}
