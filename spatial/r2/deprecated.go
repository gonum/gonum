// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r2

// TODO(sbinet): remove this file for Gonum-v0.10.0.

// Add returns the vector sum of p and q.
//
// DEPRECATED: use r2.Add.
func (p Vec) Add(q Vec) Vec {
	return Add(p, q)
}

// Sub returns the vector sum of p and -q.
//
// DEPRECATED: use r2.Sub.
func (p Vec) Sub(q Vec) Vec {
	return Sub(p, q)
}

// Scale returns the vector p scaled by f.
//
// DEPRECATED: use r2.Scale.
func (p Vec) Scale(f float64) Vec {
	return Scale(f, p)
}

// Dot returns the dot product p·q.
//
// DEPRECATED: use r2.Dot.
func (p Vec) Dot(q Vec) float64 {
	return Dot(p, q)
}

// Cross returns the cross product p×q.
//
// DEPRECATED: use r2.Cross.
func (p Vec) Cross(q Vec) float64 {
	return Cross(p, q)
}

// Rotate returns a new vector, rotated by alpha around the provided point, q.
//
// DEPRECATED: use r2.Rotate.
func (p Vec) Rotate(alpha float64, q Vec) Vec {
	return Rotate(p, alpha, q)
}
