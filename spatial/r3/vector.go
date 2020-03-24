// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import "math"

// Vec is a 3D vector.
type Vec [3]float64

func (p Vec) X() float64 { return p[0] }
func (p Vec) Y() float64 { return p[1] }
func (p Vec) Z() float64 { return p[2] }

// Add returns the vector sum of p and q.
func (p Vec) Add(q Vec) Vec {
	p[0] += q[0]
	p[1] += q[1]
	p[2] += q[2]
	return p
}

// Sub returns the vector sum of p and -q.
func (p Vec) Sub(q Vec) Vec {
	p[0] -= q[0]
	p[1] -= q[1]
	p[2] -= q[2]
	return p
}

// Scale returns the vector p scaled by f.
func (p Vec) Scale(f float64) Vec {
	p[0] *= f
	p[1] *= f
	p[2] *= f
	return p
}

// Norm returns the vector's norm.
func (v Vec) Norm() float64 {
	return math.Sqrt(Dot(v, v))
}

// Norm2 returns the vector's squared norm.
func (v Vec) Norm2() float64 {
	return Dot(v, v)
}

// Dot returns the u.v dot product of u and v.
func Dot(u, v Vec) float64 {
	return u[0]*v[0] + u[1]*v[1] + u[2]*v[2]
}

// Cross returns the u x v cross product of u and v.
func Cross(u, v Vec) Vec {
	return Vec{
		u[1]*v[2] - u[2]*v[1],
		u[2]*v[0] - u[0]*v[2],
		u[0]*v[1] - u[1]*v[0],
	}
}

// Box is a 3D bounding box.
type Box struct {
	Min, Max Vec
}
