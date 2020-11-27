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
func (p Vec) Add(q Vec) Vec {
	p.X += q.X
	p.Y += q.Y
	return p
}

// Sub returns the vector sum of p and -q.
func (p Vec) Sub(q Vec) Vec {
	p.X -= q.X
	p.Y -= q.Y
	return p
}

// Scale returns the vector p scaled by f.
func (p Vec) Scale(f float64) Vec {
	p.X *= f
	p.Y *= f
	return p
}

// Dot returns the dot product p·q.
func (p Vec) Dot(q Vec) float64 {
	return p.X*q.X + p.Y*q.Y
}

// Cross returns the cross product p×q.
func (p Vec) Cross(q Vec) float64 {
	return p.X*q.Y - p.Y*q.X
}

// Rotate returns a new vector, rotated by alpha around the provided vector.
func (p Vec) Rotate(alpha float64, q Vec) Vec {
	return NewRotation(alpha, q).Rotate(p)
}

// Norm returns the Euclidean norm of p
//  |p| = sqrt(p_x^2 + p_y^2).
func Norm(p Vec) float64 {
	return math.Hypot(p.X, p.Y)
}

// Norm returns the Euclidean squared norm of p
//  |p|^2 = p_x^2 + p_y^2.
func Norm2(p Vec) float64 {
	return p.X*p.X + p.Y*p.Y
}

// Unit returns the unit vector colinear to p.
// Unit returns {NaN,NaN} for the zero vector.
func Unit(p Vec) Vec {
	if p.X == 0 && p.Y == 0 {
		return Vec{X: math.NaN(), Y: math.NaN()}
	}
	return p.Scale(1 / Norm(p))
}

// Cos returns the cosine of the opening angle between p and q.
func Cos(p, q Vec) float64 {
	return p.Dot(q) / (Norm(p) * Norm(q))
}

// Box is a 2D bounding box.
type Box struct {
	Min, Max Vec
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

// Rotate returns the rotated vector according to the definition of rot.
func (r Rotation) Rotate(p Vec) Vec {
	if r.isIdentity() {
		return p
	}
	o := p.Sub(r.p)
	return Vec{
		X: (o.X*r.cos - o.Y*r.sin),
		Y: (o.X*r.sin + o.Y*r.cos),
	}.Add(r.p)
}

func (r Rotation) isIdentity() bool {
	return r.sin == 0 && r.cos == 1
}
