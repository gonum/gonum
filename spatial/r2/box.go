// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r2

import "math"

// Box is a 2D bounding box. Well formed Boxes have
// Min components smaller than Max components.
type Box struct {
	Min, Max Vec
}

// NewBox is shorthand for Box{Min:Vec{x0,y0}, Max:Vec{x1,y1}}.
// The sides are swapped so that the resulting Box is well formed.
func NewBox(x0, y0, x1, y1 float64) Box {
	return Box{
		Min: Vec{X: math.Min(x0, x1), Y: math.Min(y0, y1)},
		Max: Vec{X: math.Max(x0, x1), Y: math.Max(y0, y1)},
	}
}

// Size returns the size of the Box.
func (a Box) Size() Vec {
	return Sub(a.Max, a.Min)
}

// Center returns the center of the Box.
func (a Box) Center() Vec {
	return Scale(0.5, Add(a.Min, a.Max))
}

// IsEmpty returns true if a Box's volume is zero
// or if a Min component is greater than its Max component.
func (a Box) Empty() bool {
	return a.Min.X >= a.Max.X || a.Min.Y >= a.Max.Y
}

// Vertices returns a slice of the 4 vertices
// corresponding to each of the Box's corners.
//
// The order of the vertices is CCW in the XY plane starting at the box minimum.
// If viewing box from +Z position the ordering is as follows:
//  1. Bottom left.
//  2. Bottom right.
//  3. Top right.
//  4. Top left.
func (a Box) Vertices() []Vec {
	return []Vec{
		0: a.Min,
		1: {a.Max.X, a.Min.Y},
		2: a.Max,
		3: {a.Min.X, a.Max.Y},
	}
}

// Union returns a box enclosing both the receiver and argument Boxes.
func (a Box) Union(b Box) Box {
	if a.Empty() {
		return b
	}
	if b.Empty() {
		return a
	}
	return Box{
		Min: minElem(a.Min, b.Min),
		Max: maxElem(a.Max, b.Max),
	}
}

// Add adds v to the bounding box components.
// It is the equivalent of translating the Box by v.
func (a Box) Add(v Vec) Box {
	return Box{Add(a.Min, v), Add(a.Max, v)}
}

// Scale returns a new Box scaled by a size vector around its center.
// The scaling is done element wise, which is to say the Box's X size is
// scaled by v.X. Negative components of v are interpreted as zero.
func (a Box) Scale(v Vec) Box {
	v = maxElem(v, Vec{})
	// TODO(soypat): Probably a better way to do this.
	return centeredBox(a.Center(), mulElem(v, a.Size()))
}

// centeredBox creates a Box with a given center and size. Size's negative
// components are interpreted as zero so that resulting box is well formed.
func centeredBox(center, size Vec) Box {
	size = maxElem(size, Vec{})
	half := Scale(0.5, absElem(size))
	return Box{Min: Sub(center, half), Max: Add(center, half)}
}

// Contains returns true if v is contained within the bounds of the Box.
func (a Box) Contains(v Vec) bool {
	if a.Empty() {
		return v == a.Min && v == a.Max
	}
	return a.Min.X <= v.X && v.X <= a.Max.X &&
		a.Min.Y <= v.Y && v.Y <= a.Max.Y
}

// Canon returns the canonical version of b. The returned Box has minimum
// and maximum coordinates swapped if necessary so that it is well-formed.
func (b Box) Canon() Box {
	return Box{
		Min: minElem(b.Min, b.Max),
		Max: maxElem(b.Min, b.Max),
	}
}
