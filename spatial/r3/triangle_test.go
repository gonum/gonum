// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"
)

func TestTriangleDegenerate(t *testing.T) {
	const (
		// tol is how much closer the problematic
		// vertex is placed to avoid floating point error
		// for degeneracy calculation.
		tol = 1e-12
		// This is the argument to Degenerate and represents
		// the minimum permissible distance between the triangle
		// longest edge and the opposite vertex.
		spatialTol = 1e-2
	)
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 200; i++ {
		// Generate a random line for the longest triangle side.
		ln := line{randomVec(rnd), randomVec(rnd)}
		lineDir := Sub(ln[1], ln[0])

		perpendicular := Unit(Cross(lineDir, randomVec(rnd)))
		// generate 3 permutations of needle triangles for
		// each vertex. A needle triangle has two vertices
		// very close to eachother an its third vertex far away.
		var needle Triangle
		for j := 0; j < 3; j++ {
			needle[j] = ln[0]
			needle[(j+1)%3] = ln[1]
			needle[(j+2)%3] = Add(ln[1], Scale((1-tol)*spatialTol, perpendicular))
			if !needle.IsDegenerate(spatialTol) {
				t.Error("needle triangle not degenerate")
			}
		}

		midpoint := ln.vecOnLine(0.5)
		// cap triangles are characterized by having two sides
		// of similar lengths and whose sum is approximately equal
		// to the remaining longest side.
		var cap Triangle
		for j := 0; j < 3; j++ {
			cap[j] = ln[0]
			cap[(j+1)%3] = ln[1]
			cap[(j+2)%3] = Add(midpoint, Scale((1-tol)*spatialTol, perpendicular))
			if !cap.IsDegenerate(spatialTol) {
				t.Error("cap triangle not degenerate")
			}
		}

		var degenerate Triangle
		for j := 0; j < 3; j++ {
			degenerate[j] = ln[0]
			degenerate[(j+1)%3] = ln[1]
			// vertex perpendicular to some random point on longest side.
			degenerate[(j+2)%3] = Add(ln.vecOnLine(rnd.Float64()), Scale((1-tol)*spatialTol, perpendicular))
			if !degenerate.IsDegenerate(spatialTol) {
				t.Error("random degenerate triangle not degenerate")
			}
			// vertex about longest side 0 vertex
			degenerate[(j+2)%3] = Add(ln[0], Scale((1-tol)*spatialTol, Unit(randomVec(rnd))))
			if !degenerate.IsDegenerate(spatialTol) {
				t.Error("needle-like degenerate triangle not degenerate")
			}
			// vertex about longest side 1 vertex
			degenerate[(j+2)%3] = Add(ln[1], Scale((1-tol)*spatialTol, Unit(randomVec(rnd))))
			if !degenerate.IsDegenerate(spatialTol) {
				t.Error("needle-like degenerate triangle not degenerate")
			}
		}
	}
}

func TestTriangleCentroid(t *testing.T) {
	const tol = 1e-12
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		tri := randomTriangle(rnd)
		got := tri.Centroid()
		want := Vec{
			X: (tri[0].X + tri[1].X + tri[2].X) / 3,
			Y: (tri[0].Y + tri[1].Y + tri[2].Y) / 3,
			Z: (tri[0].Z + tri[1].Z + tri[2].Z) / 3,
		}
		if !vecApproxEqual(got, want, tol) {
			t.Fatalf("got %.6g, want %.6g", got, want)
		}
	}
}

func TestTriangleNormal(t *testing.T) {
	const tol = 1e-12
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		tri := randomTriangle(rnd)
		got := tri.Normal()
		expect := goldenNormal(tri)
		if !vecApproxEqual(got, expect, tol) {
			t.Fatalf("got %.6g, want %.6g", got, expect)
		}
	}
}

func TestTriangleArea(t *testing.T) {
	const tol = 1e-16
	for _, test := range []struct {
		T      Triangle
		Expect float64
	}{
		{
			T: Triangle{
				{0, 0, 0},
				{1, 0, 0},
				{0, 1, 0},
			},
			Expect: 0.5,
		},
		{
			T: Triangle{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
			Expect: 0.5,
		},
		{
			T: Triangle{
				{20, 0, 0},
				{0, 0, 20},
				{0, 0, 0},
			},
			Expect: 20 * 20 / 2,
		},
	} {
		got := test.T.Area()
		if math.Abs(got-test.Expect) > tol {
			t.Errorf("got area %g, expected %g", got, test.Expect)
		}
		if test.T.IsDegenerate(tol) {
			t.Error("well-formed triangle is degenerate")
		}
	}
	const tol2 = 1e-12
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		tri := randomTriangle(rnd)
		got := tri.Area()
		want := Norm(Cross(Sub(tri[1], tri[0]), Sub(tri[2], tri[0]))) / 2
		if math.Abs(got-want) > tol2 {
			t.Errorf("got area %g not match half norm of cross product %g", got, want)
		}
	}
}

func TestTriangleOrderedLengths(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 200; i++ {
		tri := randomTriangle(rnd)
		s1, s2, s3 := tri.sides()
		l1 := Norm(s1)
		l2 := Norm(s2)
		l3 := Norm(s3)
		a, b, c := tri.orderedLengths()
		if a != l1 && a != l2 && a != l3 {
			t.Error("shortest ordered length not a side of the triangle")
		}
		if b != l1 && b != l2 && b != l3 {
			t.Error("middle ordered length not a side of the triangle")
		}
		if c != l1 && c != l2 && c != l3 {
			t.Error("longest ordered length not a side of the triangle")
		}
		if a > b || a > c {
			t.Error("ordered short side not shortest side")
		}
		if c < b {
			t.Error("ordered long side not longest side")
		}
	}
}

// taken from soypat/sdf library where it has been thoroughly tested empirically.
func goldenNormal(t Triangle) Vec {
	e1 := Sub(t[1], t[0])
	e2 := Sub(t[2], t[0])
	return Cross(e1, e2)
}

func randomTriangle(rnd *rand.Rand) Triangle {
	return Triangle{
		randomVec(rnd),
		randomVec(rnd),
		randomVec(rnd),
	}
}
