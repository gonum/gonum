// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3_test

import (
	"fmt"

	"gonum.org/v1/gonum/spatial/r3"
)

func ExampleTriangle_icosphere() {
	// This example generates a 3D icosphere from
	// a starting icosahedron by subdividing surfaces.
	// See https://schneide.blog/2016/07/15/generating-an-icosphere-in-c/.
	const subdivisions = 5
	// vertices is a slice of r3.Vec
	// triangles is a slice of [3]int indices
	// referencing the vertices.
	vertices, triangles := icosahedron()
	for i := 0; i < subdivisions; i++ {
		vertices, triangles = subdivide(vertices, triangles)
	}
	var faces []r3.Triangle
	for _, t := range triangles {
		var face r3.Triangle
		for i := 0; i < 3; i++ {
			face[i] = vertices[t[i]]
		}
		faces = append(faces, face)
	}
	fmt.Println(faces)
	// The 3D rendering of the icosphere is left as an exercise to the reader.
}

// edgeIdx represents an edge of the icosahedron
type edgeIdx [2]int

func subdivide(vertices []r3.Vec, triangles [][3]int) ([]r3.Vec, [][3]int) {
	// We generate a lookup table of all newly generated vertices so as to not
	// duplicate new vertices. edgeIdx has lower index first.
	lookup := make(map[edgeIdx]int)
	var result [][3]int
	for _, triangle := range triangles {
		var mid [3]int
		for edge := 0; edge < 3; edge++ {
			lookup, mid[edge], vertices = subdivideEdge(lookup, vertices, triangle[edge], triangle[(edge+1)%3])
		}
		newTriangles := [][3]int{
			{triangle[0], mid[0], mid[2]},
			{triangle[1], mid[1], mid[0]},
			{triangle[2], mid[2], mid[1]},
			{mid[0], mid[1], mid[2]},
		}
		result = append(result, newTriangles...)
	}
	return vertices, result
}

// subdivideEdge takes the vertices list and indices first and second which
// refer to the edge that will be subdivided.
// lookup is a table of all newly generated vertices from
// previous calls to subdivideEdge so as to not duplicate vertices.
func subdivideEdge(lookup map[edgeIdx]int, vertices []r3.Vec, first, second int) (map[edgeIdx]int, int, []r3.Vec) {
	key := edgeIdx{first, second}
	if first > second {
		// Swap to ensure edgeIdx always has lower index first.
		key[0], key[1] = key[1], key[0]
	}
	vertIdx, vertExists := lookup[key]
	if !vertExists {
		// If edge not already subdivided add
		// new dividing vertex to lookup table.
		edge0 := vertices[first]
		edge1 := vertices[second]
		point := r3.Unit(r3.Add(edge0, edge1)) // vertex at a normalized position.
		vertices = append(vertices, point)
		vertIdx = len(vertices) - 1
		lookup[key] = vertIdx
	}
	return lookup, vertIdx, vertices
}

// icosahedron returns an icosahedron mesh.
func icosahedron() (vertices []r3.Vec, triangles [][3]int) {
	const (
		radiusSqrt = 1.0 // Example designed for unit sphere generation.
		X          = radiusSqrt * .525731112119133606
		Z          = radiusSqrt * .850650808352039932
		N          = 0.0
	)
	return []r3.Vec{
			{-X, N, Z}, {X, N, Z}, {-X, N, -Z}, {X, N, -Z},
			{N, Z, X}, {N, Z, -X}, {N, -Z, X}, {N, -Z, -X},
			{Z, X, N}, {-Z, X, N}, {Z, -X, N}, {-Z, -X, N},
		}, [][3]int{
			{0, 1, 4}, {0, 4, 9}, {9, 4, 5}, {4, 8, 5},
			{4, 1, 8}, {8, 1, 10}, {8, 10, 3}, {5, 8, 3},
			{5, 3, 2}, {2, 3, 7}, {7, 3, 10}, {7, 10, 6},
			{7, 6, 11}, {11, 6, 0}, {0, 6, 1}, {6, 10, 1},
			{9, 11, 0}, {9, 2, 11}, {9, 5, 2}, {7, 11, 2},
		}
}
