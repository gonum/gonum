// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r2

import "math"

// Vec is a 2D vector.
type Vec struct {
	X, Y float64
}

// Add returns the vector sum of p and q.
func Add(p, q Vec) Vec {
	return Vec{
		X: p.X + q.X,
		Y: p.Y + q.Y,
	}
}

// Sub returns the vector sum of p and -q.
func Sub(p, q Vec) Vec {
	return Vec{
		X: p.X - q.X,
		Y: p.Y - q.Y,
	}
}

// Scale returns the vector p scaled by f.
func Scale(f float64, p Vec) Vec {
	return Vec{
		X: f * p.X,
		Y: f * p.Y,
	}
}

// Dot returns the dot product p·q.
func Dot(p, q Vec) float64 {
	return p.X*q.X + p.Y*q.Y
}

// Cross returns the cross product p×q.
func Cross(p, q Vec) float64 {
	return p.X*q.Y - p.Y*q.X
}

// Rotate returns a new vector, rotated by alpha around the provided point, q.
func Rotate(p Vec, alpha float64, q Vec) Vec {
	return NewRotation(alpha, q).Rotate(p)
}

// Norm returns the Euclidean norm of p
//
//	|p| = sqrt(p_x^2 + p_y^2).
func Norm(p Vec) float64 {
	return math.Hypot(p.X, p.Y)
}

// Norm2 returns the Euclidean squared norm of p
//
//	|p|^2 = p_x^2 + p_y^2.
func Norm2(p Vec) float64 {
	return p.X*p.X + p.Y*p.Y
}

// Unit returns the unit vector colinear to p.
// Unit returns {NaN,NaN} for the zero vector.
func Unit(p Vec) Vec {
	if p.X == 0 && p.Y == 0 {
		return Vec{X: math.NaN(), Y: math.NaN()}
	}
	return Scale(1/Norm(p), p)
}

// Cos returns the cosine of the opening angle between p and q.
func Cos(p, q Vec) float64 {
	return Dot(p, q) / (Norm(p) * Norm(q))
}

// Rotation describes a rotation in 2D.
type Rotation struct {
	sin, cos float64
	p        Vec
}

// NewRotation creates a rotation by alpha, around p.
func NewRotation(alpha float64, p Vec) Rotation {
	if alpha == 0 {
		return Rotation{sin: 0, cos: 1, p: p}
	}
	sin, cos := math.Sincos(alpha)
	return Rotation{sin: sin, cos: cos, p: p}
}

// Rotate returns p rotated according to the parameters used to construct
// the receiver.
func (r Rotation) Rotate(p Vec) Vec {
	if r.isIdentity() {
		return p
	}
	o := Sub(p, r.p)
	return Add(Vec{
		X: (o.X*r.cos - o.Y*r.sin),
		Y: (o.X*r.sin + o.Y*r.cos),
	}, r.p)
}

func (r Rotation) isIdentity() bool {
	return r.sin == 0 && r.cos == 1
}

// minElem returns a vector with the element-wise
// minimum components of vectors a and b.
func minElem(a, b Vec) Vec {
	return Vec{
		X: math.Min(a.X, b.X),
		Y: math.Min(a.Y, b.Y),
	}
}

// maxElem returns a vector with the element-wise
// maximum components of vectors a and b.
func maxElem(a, b Vec) Vec {
	return Vec{
		X: math.Max(a.X, b.X),
		Y: math.Max(a.Y, b.Y),
	}
}

// absElem returns the vector with components set to their absolute value.
func absElem(a Vec) Vec {
	return Vec{
		X: math.Abs(a.X),
		Y: math.Abs(a.Y),
	}
}

// mulElem returns the Hadamard product between vectors a and b.
//
//	v = {a.X*b.X, a.Y*b.Y, a.Z*b.Z}
func mulElem(a, b Vec) Vec {
	return Vec{
		X: a.X * b.X,
		Y: a.Y * b.Y,
	}
}

// divElem returns the Hadamard product between vector a
// and the inverse components of vector b.
//
//	v = {a.X/b.X, a.Y/b.Y, a.Z/b.Z}
func divElem(a, b Vec) Vec {
	return Vec{
		X: a.X / b.X,
		Y: a.Y / b.Y,
	}
}
