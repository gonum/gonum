// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import (
	"fmt"
	"math"
)

// Affine represents a 3D affine spatial transformation.
// The zero value of Affine is the identity transform.
type Affine struct {
	// in order to make the zero value of Transform represent the identity
	// transform we store it with the identity matrix subtracted.
	// These diagonal elements are subtracted such that
	//  d00 = x00-1, d11 = x11-1, d22 = x22-1, d33 = x33-1
	// where x00, x11, x22, x33 are the matrix diagonal elements.
	// We can then check for identity in if blocks like so:
	//  if T == (Transform{})
	d00, x01, x02, x03 float64
	x10, d11, x12, x13 float64
	x20, x21, d22, x23 float64
	x30, x31, x32, d33 float64
}

// Transform applies the transform to the argument vector and returns the result.
func (a Affine) Transform(v Vec) Vec {
	if a.isIdentity() {
		return v
	}
	// See https://github.com/mrdoob/three.js/blob/019fa1ad671a1ffcf9be5828efd518fb06575c2b/src/math/Vector3.js#L263.
	w := 1 / (a.x30*v.X + a.x31*v.Y + a.x32*v.Z + a.d33 + 1)
	if math.IsInf(w, 0) {
		return Vec{}
	}
	return Vec{
		X: ((a.d00+1)*v.X + a.x01*v.Y + a.x02*v.Z + a.x03) * w,
		Y: (a.x10*v.X + (a.d11+1)*v.Y + a.x12*v.Z + a.x13) * w,
		Z: (a.x20*v.X + a.x21*v.Y + (a.d22+1)*v.Z + a.x23) * w,
	}
}

// zeroAffine is the affine transform that returns zeroAffine when multiplied by any transform.
var zeroAffine = Affine{d00: -1, d11: -1, d22: -1, d33: -1}

// NewAffine returns a new Affine type and populates its elements
// with values passed in row-major form. If val is nil then NewAffine
// returns a Affine filled with zeros.
func NewAffine(a []float64) Affine {
	if a == nil {
		return zeroAffine
	}
	if len(a) != 16 {
		panic(fmt.Sprintf("r3: NewTransform invalid slice length: %d != 16", len(a)))
	}
	return Affine{
		d00: a[0] - 1, x01: a[1], x02: a[2], x03: a[3],
		x10: a[4], d11: a[5] - 1, x12: a[6], x13: a[7],
		x20: a[8], x21: a[9], d22: a[10] - 1, x23: a[11],
		x30: a[12], x31: a[13], x32: a[14], d33: a[15] - 1,
	}
}

// AddTranslation adds v to the translational transform of a and returns the result.
func (a Affine) AddTranslation(v Vec) Affine {
	a.x03 += v.X
	a.x13 += v.Y
	a.x23 += v.Z
	return a
}

// WithTranslation returns a with the translatational components set to v.
func (a Affine) WithTranslation(v Vec) Affine {
	a.x03 = v.X
	a.x13 = v.Y
	a.x23 = v.Z
	return a
}

// Scale returns the transform with factor scaling added around origin.
func (a Affine) Scale(origin, factor Vec) Affine {
	if origin == (Vec{}) {
		return a.scale(factor)
	}
	a = a.AddTranslation(Scale(-1, origin))
	a = a.scale(factor)
	return a.AddTranslation(origin)
}

func (a Affine) scale(factor Vec) Affine {
	// scale is a diagonal matrix.
	a.d00 = (a.d00+1)*factor.X - 1
	a.x10 *= factor.X
	a.x20 *= factor.X
	a.x30 *= factor.X

	a.x01 *= factor.Y
	a.d11 = (a.d11+1)*factor.Y - 1
	a.x21 *= factor.Y
	a.x31 *= factor.Y

	a.x02 *= factor.Z
	a.x12 *= factor.Z
	a.d22 = (a.d22+1)*factor.Z - 1
	a.x32 *= factor.Z
	return a
}

// Mul performs matrix multiplication of affine transforms a and b and returns
// the result c. This operation is the equivalent of creating a new
// Affine that first applies b followed by a.
func (a Affine) Mul(b Affine) Affine {
	if a.isIdentity() {
		return b
	}
	if b.isIdentity() {
		return a
	}
	x00 := a.d00 + 1
	x11 := a.d11 + 1
	x22 := a.d22 + 1
	x33 := a.d33 + 1
	y00 := b.d00 + 1
	y11 := b.d11 + 1
	y22 := b.d22 + 1
	y33 := b.d33 + 1
	var m Affine
	m.d00 = x00*y00 + a.x01*b.x10 + a.x02*b.x20 + a.x03*b.x30 - 1
	m.x10 = a.x10*y00 + x11*b.x10 + a.x12*b.x20 + a.x13*b.x30
	m.x20 = a.x20*y00 + a.x21*b.x10 + x22*b.x20 + a.x23*b.x30
	m.x30 = a.x30*y00 + a.x31*b.x10 + a.x32*b.x20 + x33*b.x30
	m.x01 = x00*b.x01 + a.x01*y11 + a.x02*b.x21 + a.x03*b.x31
	m.d11 = a.x10*b.x01 + x11*y11 + a.x12*b.x21 + a.x13*b.x31 - 1
	m.x21 = a.x20*b.x01 + a.x21*y11 + x22*b.x21 + a.x23*b.x31
	m.x31 = a.x30*b.x01 + a.x31*y11 + a.x32*b.x21 + x33*b.x31
	m.x02 = x00*b.x02 + a.x01*b.x12 + a.x02*y22 + a.x03*b.x32
	m.x12 = a.x10*b.x02 + x11*b.x12 + a.x12*y22 + a.x13*b.x32
	m.d22 = a.x20*b.x02 + a.x21*b.x12 + x22*y22 + a.x23*b.x32 - 1
	m.x32 = a.x30*b.x02 + a.x31*b.x12 + a.x32*y22 + x33*b.x32
	m.x03 = x00*b.x03 + a.x01*b.x13 + a.x02*b.x23 + a.x03*y33
	m.x13 = a.x10*b.x03 + x11*b.x13 + a.x12*b.x23 + a.x13*y33
	m.x23 = a.x20*b.x03 + a.x21*b.x13 + x22*b.x23 + a.x23*y33
	m.d33 = a.x30*b.x03 + a.x31*b.x13 + a.x32*b.x23 + x33*y33 - 1
	return m
}

// Det returns the determinant of the affine transform matrix.
func (a Affine) Det() float64 {
	x00 := a.d00 + 1
	x11 := a.d11 + 1
	x22 := a.d22 + 1
	x33 := a.d33 + 1
	return x00*x11*x22*x33 - x00*x11*a.x23*a.x32 +
		x00*a.x12*a.x23*a.x31 - x00*a.x12*a.x21*x33 +
		x00*a.x13*a.x21*a.x32 - x00*a.x13*x22*a.x31 -
		a.x01*a.x12*a.x23*a.x30 + a.x01*a.x12*a.x20*x33 -
		a.x01*a.x13*a.x20*a.x32 + a.x01*a.x13*x22*a.x30 -
		a.x01*a.x10*x22*x33 + a.x01*a.x10*a.x23*a.x32 +
		a.x02*a.x13*a.x20*a.x31 - a.x02*a.x13*a.x21*a.x30 +
		a.x02*a.x10*a.x21*x33 - a.x02*a.x10*a.x23*a.x31 +
		a.x02*x11*a.x23*a.x30 - a.x02*x11*a.x20*x33 -
		a.x03*a.x10*a.x21*a.x32 + a.x03*a.x10*x22*a.x31 -
		a.x03*x11*x22*a.x30 + a.x03*x11*a.x20*a.x32 -
		a.x03*a.x12*a.x20*a.x31 + a.x03*a.x12*a.x21*a.x30
}

// Inv returns the inverse of the affine transform such that a.Inv() * a is the
// identity transform. If a is singular then Inv() returns the zero transform.
func (a Affine) Inv() Affine {
	if a.isIdentity() {
		return a
	}
	det := a.Det()
	if math.Abs(det) < 1e-16 {
		return zeroAffine
	}
	d := 1 / det
	x00 := a.d00 + 1
	x11 := a.d11 + 1
	x22 := a.d22 + 1
	x33 := a.d33 + 1
	var m Affine
	m.d00 = (a.x12*a.x23*a.x31-a.x13*x22*a.x31+a.x13*a.x21*a.x32-x11*a.x23*a.x32-a.x12*a.x21*x33+x11*x22*x33)*d - 1
	m.x01 = (a.x03*x22*a.x31 - a.x02*a.x23*a.x31 - a.x03*a.x21*a.x32 + a.x01*a.x23*a.x32 + a.x02*a.x21*x33 - a.x01*x22*x33) * d
	m.x02 = (a.x02*a.x13*a.x31 - a.x03*a.x12*a.x31 + a.x03*x11*a.x32 - a.x01*a.x13*a.x32 - a.x02*x11*x33 + a.x01*a.x12*x33) * d
	m.x03 = (a.x03*a.x12*a.x21 - a.x02*a.x13*a.x21 - a.x03*x11*x22 + a.x01*a.x13*x22 + a.x02*x11*a.x23 - a.x01*a.x12*a.x23) * d
	m.x10 = (a.x13*x22*a.x30 - a.x12*a.x23*a.x30 - a.x13*a.x20*a.x32 + a.x10*a.x23*a.x32 + a.x12*a.x20*x33 - a.x10*x22*x33) * d
	m.d11 = (a.x02*a.x23*a.x30-a.x03*x22*a.x30+a.x03*a.x20*a.x32-x00*a.x23*a.x32-a.x02*a.x20*x33+x00*x22*x33)*d - 1
	m.x12 = (a.x03*a.x12*a.x30 - a.x02*a.x13*a.x30 - a.x03*a.x10*a.x32 + x00*a.x13*a.x32 + a.x02*a.x10*x33 - x00*a.x12*x33) * d
	m.x13 = (a.x02*a.x13*a.x20 - a.x03*a.x12*a.x20 + a.x03*a.x10*x22 - x00*a.x13*x22 - a.x02*a.x10*a.x23 + x00*a.x12*a.x23) * d
	m.x20 = (x11*a.x23*a.x30 - a.x13*a.x21*a.x30 + a.x13*a.x20*a.x31 - a.x10*a.x23*a.x31 - x11*a.x20*x33 + a.x10*a.x21*x33) * d
	m.x21 = (a.x03*a.x21*a.x30 - a.x01*a.x23*a.x30 - a.x03*a.x20*a.x31 + x00*a.x23*a.x31 + a.x01*a.x20*x33 - x00*a.x21*x33) * d
	m.d22 = (a.x01*a.x13*a.x30-a.x03*x11*a.x30+a.x03*a.x10*a.x31-x00*a.x13*a.x31-a.x01*a.x10*x33+x00*x11*x33)*d - 1
	m.x23 = (a.x03*x11*a.x20 - a.x01*a.x13*a.x20 - a.x03*a.x10*a.x21 + x00*a.x13*a.x21 + a.x01*a.x10*a.x23 - x00*x11*a.x23) * d
	m.x30 = (a.x12*a.x21*a.x30 - x11*x22*a.x30 - a.x12*a.x20*a.x31 + a.x10*x22*a.x31 + x11*a.x20*a.x32 - a.x10*a.x21*a.x32) * d
	m.x31 = (a.x01*x22*a.x30 - a.x02*a.x21*a.x30 + a.x02*a.x20*a.x31 - x00*x22*a.x31 - a.x01*a.x20*a.x32 + x00*a.x21*a.x32) * d
	m.x32 = (a.x02*x11*a.x30 - a.x01*a.x12*a.x30 - a.x02*a.x10*a.x31 + x00*a.x12*a.x31 + a.x01*a.x10*a.x32 - x00*x11*a.x32) * d
	m.d33 = (a.x01*a.x12*a.x20-a.x02*x11*a.x20+a.x02*a.x10*a.x21-x00*a.x12*a.x21-a.x01*a.x10*x22+x00*x11*x22)*d - 1
	return m
}

func makeAffineRotation(r Rotation) Affine {
	if r == (Rotation{}) {
		return Affine{}
	}
	w, i, j, k := r.Real, r.Imag, r.Jmag, r.Kmag
	ii := 2 * i * i
	jj := 2 * j * j
	kk := 2 * k * k
	wi := 2 * w * i
	wj := 2 * w * j
	wk := 2 * w * k
	ij := 2 * i * j
	jk := 2 * j * k
	ki := 2 * k * i
	return Affine{
		d00: -jj - kk, x01: ij - wk, x02: ki + wj,
		x10: ij + wk, d11: -ii - kk, x12: jk - wi,
		x20: ki - wj, x21: jk + wi, d22: -ii - jj,
	}
}

// transpose returns the transposed transform.
func (a Affine) transpose() Affine {
	return Affine{
		d00: a.d00, x01: a.x10, x02: a.x20, x03: a.x30,
		x10: a.x01, d11: a.d11, x12: a.x21, x13: a.x31,
		x20: a.x02, x21: a.x12, d22: a.d22, x23: a.x32,
		x30: a.x03, x31: a.x13, x32: a.x23, d33: a.d33,
	}
}

// sliceCopy returns a copy of the transform's data
// in row major storage format. It returns 16 elements.
func (a Affine) sliceCopy() []float64 {
	return []float64{
		a.d00 + 1, a.x01, a.x02, a.x03,
		a.x10, a.d11 + 1, a.x12, a.x13,
		a.x20, a.x21, a.d22 + 1, a.x23,
		a.x30, a.x31, a.x32, a.d33 + 1,
	}
}

// isIdentity returns true if receiver is the identity transform.
func (a Affine) isIdentity() bool {
	// The zero value of Affine is guaranteed to be the identity.
	return a == Affine{}
}

// IsZero returns true if a is the zero transform.
func (a Affine) IsZero() bool {
	return a == zeroAffine
}
