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

// Norm returns the euclidian norm of p:
//  |p| = sqrt(p_x^2 + p_y^2)
func Norm(p Vec) float64 {
	return math.Sqrt(Norm2(p))
}

// Norm returns the euclidian squared norm of p:
//  |p|^2 = p_x^2 + p_y^2
func Norm2(p Vec) float64 {
	return p.X*p.X + p.Y*p.Y
}

// Box is a 2D bounding box.
type Box struct {
	Min, Max Vec
}
