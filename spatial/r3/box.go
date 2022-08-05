// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import "math"

// Box is a 3D bounding box. Well formed Boxes Min components
// are smaller than Max components.
type Box struct {
	Min, Max Vec
}

// NewBox is shorthand for Box{Min:Vec{x0,y0,z0}, Max:Vec{x1,y1,z1}}.
// The sides are swapped so that the resulting Box is well formed.
func NewBox(x0, y0, z0, x1, y1, z1 float64) Box {
	return Box{
		Min: Vec{X: math.Min(x0, x1), Y: math.Min(y0, y1), Z: math.Min(z0, z1)},
		Max: Vec{X: math.Max(x0, x1), Y: math.Max(y0, y1), Z: math.Max(z0, z1)},
	}
}

// IsEmpty returns true if a Box's volume is zero
// or if a Min component is greater than its Max component.
func (a Box) Empty() bool {
	return a.Min.X >= a.Max.X || a.Min.Y >= a.Max.Y || a.Min.Z >= a.Max.Z
}

// Size returns the size of the Box.
func (a Box) Size() Vec {
	return Sub(a.Max, a.Min)
}

// Center returns the center of the Box.
func (a Box) Center() Vec {
	return Scale(0.5, Add(a.Min, a.Max))
}

// Vertices returns a slice of the 8 vertices
// corresponding to each of the Box's corners.
//
// Ordering of vertices 0-3 is CCW in the XY plane starting at box minimum.
// Ordering of vertices 4-7 is CCW in the XY plane starting at box minimum
// for X and Y values and maximum Z value.
//
// Edges for the box can be constructed with the following indices:
//
//	edges := [12][2]int{
//	 {0, 1}, {1, 2}, {2, 3}, {3, 0},
//	 {4, 5}, {5, 6}, {6, 7}, {7, 4},
//	 {0, 4}, {1, 5}, {2, 6}, {3, 7},
//	}
func (a Box) Vertices() []Vec {
	return []Vec{
		0: a.Min,
		1: {X: a.Max.X, Y: a.Min.Y, Z: a.Min.Z},
		2: {X: a.Max.X, Y: a.Max.Y, Z: a.Min.Z},
		3: {X: a.Min.X, Y: a.Max.Y, Z: a.Min.Z},
		4: {X: a.Min.X, Y: a.Min.Y, Z: a.Max.Z},
		5: {X: a.Max.X, Y: a.Min.Y, Z: a.Max.Z},
		6: a.Max,
		7: {X: a.Min.X, Y: a.Max.Y, Z: a.Max.Z},
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
// The scaling is done element wise which is to say the Box's X dimension
// is scaled by scale.X. Negative elements of scale are interpreted as zero.
func (a Box) Scale(scale Vec) Box {
	scale = maxElem(scale, Vec{})
	// TODO(soypat): Probably a better way to do this.
	return centeredBox(a.Center(), mulElem(scale, a.Size()))
}

// centeredBox creates a Box with a given center and size.
// Negative components of size will be interpreted as zero.
func centeredBox(center, size Vec) Box {
	size = maxElem(size, Vec{}) // set negative values to zero.
	half := Scale(0.5, size)
	return Box{Min: Sub(center, half), Max: Add(center, half)}
}

// Contains returns true if v is contained within the bounds of the Box.
func (a Box) Contains(v Vec) bool {
	if a.Empty() {
		return v == a.Min && v == a.Max
	}
	return a.Min.X <= v.X && v.X <= a.Max.X &&
		a.Min.Y <= v.Y && v.Y <= a.Max.Y &&
		a.Min.Z <= v.Z && v.Z <= a.Max.Z
}

// Canon returns the canonical version of a. The returned Box has minimum
// and maximum coordinates swapped if necessary so that it is well-formed.
func (a Box) Canon() Box {
	return Box{
		Min: minElem(a.Min, a.Max),
		Max: maxElem(a.Min, a.Max),
	}
}
