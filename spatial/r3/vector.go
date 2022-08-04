// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import "math"

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
//
//	|p| = sqrt(p_x^2 + p_y^2 + p_z^2).
func Norm(p Vec) float64 {
	return math.Hypot(p.X, math.Hypot(p.Y, p.Z))
}

// Norm2 returns the Euclidean squared norm of p
//
//	|p|^2 = p_x^2 + p_y^2 + p_z^2.
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

// Divergence returns the divergence of the vector field at the point p,
// approximated using finite differences with the given step sizes.
func Divergence(p, step Vec, field func(Vec) Vec) float64 {
	sx := Vec{X: step.X}
	divx := (field(Add(p, sx)).X - field(Sub(p, sx)).X) / step.X
	sy := Vec{Y: step.Y}
	divy := (field(Add(p, sy)).Y - field(Sub(p, sy)).Y) / step.Y
	sz := Vec{Z: step.Z}
	divz := (field(Add(p, sz)).Z - field(Sub(p, sz)).Z) / step.Z
	return 0.5 * (divx + divy + divz)
}

// Gradient returns the gradient of the scalar field at the point p,
// approximated using finite differences with the given step sizes.
func Gradient(p, step Vec, field func(Vec) float64) Vec {
	dx := Vec{X: step.X}
	dy := Vec{Y: step.Y}
	dz := Vec{Z: step.Z}
	return Vec{
		X: (field(Add(p, dx)) - field(Sub(p, dx))) / (2 * step.X),
		Y: (field(Add(p, dy)) - field(Sub(p, dy))) / (2 * step.Y),
		Z: (field(Add(p, dz)) - field(Sub(p, dz))) / (2 * step.Z),
	}
}

// minElem return a vector with the minimum components of two vectors.
func minElem(a, b Vec) Vec {
	return Vec{
		X: math.Min(a.X, b.X),
		Y: math.Min(a.Y, b.Y),
		Z: math.Min(a.Z, b.Z),
	}
}

// maxElem return a vector with the maximum components of two vectors.
func maxElem(a, b Vec) Vec {
	return Vec{
		X: math.Max(a.X, b.X),
		Y: math.Max(a.Y, b.Y),
		Z: math.Max(a.Z, b.Z),
	}
}

// absElem returns the vector with components set to their absolute value.
func absElem(a Vec) Vec {
	return Vec{
		X: math.Abs(a.X),
		Y: math.Abs(a.Y),
		Z: math.Abs(a.Z),
	}
}

// mulElem returns the Hadamard product between vectors a and b.
//
//	v = {a.X*b.X, a.Y*b.Y, a.Z*b.Z}
func mulElem(a, b Vec) Vec {
	return Vec{
		X: a.X * b.X,
		Y: a.Y * b.Y,
		Z: a.Z * b.Z,
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
		Z: a.Z / b.Z,
	}
}
