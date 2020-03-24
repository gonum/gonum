// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

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

// Box is a 3D bounding box.
type Box struct {
	Min, Max Vec
}
