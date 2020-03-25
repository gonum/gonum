// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

// Vec is a 3D vector.
type Vec struct {
	X, Y, Z float64
}

// Add returns the vector sum of p and q.
func (p Vec) Add(q Vec) Vec {
	p.X += q.X
	p.Y += q.Y
	p.Z += q.Z
	return p
}

// Sub returns the vector sum of p and -q.
func (p Vec) Sub(q Vec) Vec {
	p.X -= q.X
	p.Y -= q.Y
	p.Z -= q.Z
	return p
}

// Scale returns the vector p scaled by f.
func (p Vec) Scale(f float64) Vec {
	p.X *= f
	p.Y *= f
	p.Z *= f
	return p
}

// Dot returns the dot product p·q.
func (p Vec) Dot(q Vec) float64 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z
}

// Cross returns the cross product p×q.
func (p Vec) Cross(q Vec) Vec {
	return Vec{
		p.Y*q.Z - p.Z*q.Y,
		p.Z*q.X - p.X*q.Z,
		p.X*q.Y - p.Y*q.X,
	}
}

// Box is a 3D bounding box.
type Box struct {
	Min, Max Vec
}
