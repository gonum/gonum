// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r2

import "math"

// Triangle represents a triangle in 2D space and
// is composed by 3 vectors corresponding to the position
// of each of the vertices.
type Triangle [3]Vec

// Centroid returns the intersection of the three medians of the triangle
// as a point in space.
func (t Triangle) Centroid() Vec {
	return Scale(1.0/3.0, Add(Add(t[0], t[1]), t[2]))
}

// Area returns the surface area of the triangle.
func (t Triangle) Area() float64 {
	// Heron's Formula, see https://en.wikipedia.org/wiki/Heron%27s_formula.
	// Also see William M. Kahan (24 March 2000). "Miscalculating Area and Angles of a Needle-like Triangle"
	// for more discussion. https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf.
	a, b, c := t.orderedLengths()
	A := (c + (b + a)) * (a - (c - b))
	A *= (a + (c - b)) * (c + (b - a))
	return math.Sqrt(A) / 4
}

// orderedLengths returns the lengths of the sides of the triangle such that
// a ≤ b ≤ c.
func (t Triangle) orderedLengths() (a, b, c float64) {
	s1, s2, s3 := t.sides()
	l1 := Norm(s1)
	l2 := Norm(s2)
	l3 := Norm(s3)
	// sort-3
	if l2 < l1 {
		l1, l2 = l2, l1
	}
	if l3 < l2 {
		l2, l3 = l3, l2
		if l2 < l1 {
			l1, l2 = l2, l1
		}
	}
	return l1, l2, l3
}

// sides returns vectors for each of the sides of t.
func (t Triangle) sides() (Vec, Vec, Vec) {
	return Sub(t[1], t[0]), Sub(t[2], t[1]), Sub(t[0], t[2])
}

// IsDegenerate returns true if all of triangle's vertices are
// within tol distance of its longest side.
func (t Triangle) IsDegenerate(tol float64) bool {
	sides := [3]Vec{Sub(t[1], t[0]), Sub(t[2], t[1]), Sub(t[0], t[2])}
	len2 := [3]float64{Norm2(sides[0]), Norm2(sides[1]), Norm2(sides[2])}
	longLen := len2[0]
	longIdx := 0
	if len2[1] > longLen {
		longLen = len2[1]
		longIdx = 1
	}
	if len2[2] > longLen {
		longIdx = 2
	}
	// calculate vertex distance from longest side
	ln := line{t[longIdx], t[(longIdx+1)%3]}
	dist := ln.distance(t[(longIdx+2)%3])
	return dist <= tol
}

// line is an infinite 3D line
// defined by two points on the line.
type line [2]Vec

// vecOnLine takes a value between 0 and 1 to linearly
// interpolate a point on the line.
//
//	vecOnLine(0) returns l[0]
//	vecOnLine(1) returns l[1]
func (l line) vecOnLine(t float64) Vec {
	lineDir := Sub(l[1], l[0])
	return Add(l[0], Scale(t, lineDir))
}

// distance returns the minimum euclidean distance of point p to the line.
func (l line) distance(p Vec) float64 {
	// https://mathworld.wolfram.com/Point-LineDistance3-Dimensional.html
	p1 := l[0]
	p2 := l[1]
	num := math.Abs((p2.X-p1.X)*(p1.Y-p.Y) - (p1.X-p.X)*(p2.Y-p1.Y))
	return num / math.Hypot(p2.X-p1.X, p2.Y-p1.Y)
}
