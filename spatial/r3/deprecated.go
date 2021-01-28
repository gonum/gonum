// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

// TODO(sbinet): remove this file for Gonum-v0.10.0.

// Add returns the vector sum of p and q.
//
// DEPRECATED: use r3.Add.
func (p Vec) Add(q Vec) Vec {
	return Add(p, q)
}

// Sub returns the vector sum of p and -q.
//
// DEPRECATED: use r3.Sub.
func (p Vec) Sub(q Vec) Vec {
	return Sub(p, q)
}

// Scale returns the vector p scaled by f.
//
// DEPRECATED: use r3.Scale.
func (p Vec) Scale(f float64) Vec {
	return Scale(f, p)
}

// Dot returns the dot product p·q.
//
// DEPRECATED: use r3.Dot.
func (p Vec) Dot(q Vec) float64 {
	return Dot(p, q)
}

// Cross returns the cross product p×q.
//
// DEPRECATED: use r3.Cross.
func (p Vec) Cross(q Vec) Vec {
	return Cross(p, q)
}

// Rotate returns a new vector, rotated by alpha around the provided axis.
//
// DEPRECATED: use r3.Rotate
func (p Vec) Rotate(alpha float64, axis Vec) Vec {
	return Rotate(p, alpha, axis)
}
