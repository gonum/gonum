// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spectral_test

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/spectral"
	"gonum.org/v1/gonum/mat"
)

func Example_cospectral() {
	// Output the spectra of two isospectral enneahedra.
	// https://commons.wikimedia.org/wiki/File:Isospectral_enneahedra.svg

	enneahedron1 := simple.NewUndirectedMatrix(8, 0, 0, 0)
	for _, e := range []simple.Edge{
		{F: simple.Node(0), T: simple.Node(1)},
		{F: simple.Node(0), T: simple.Node(2)},
		{F: simple.Node(0), T: simple.Node(6)},
		{F: simple.Node(1), T: simple.Node(2)},
		{F: simple.Node(1), T: simple.Node(3)},
		{F: simple.Node(1), T: simple.Node(5)},
		{F: simple.Node(1), T: simple.Node(7)},
		{F: simple.Node(2), T: simple.Node(3)},
		{F: simple.Node(2), T: simple.Node(6)},
		{F: simple.Node(3), T: simple.Node(4)},
		{F: simple.Node(4), T: simple.Node(5)},
		{F: simple.Node(4), T: simple.Node(6)},
		{F: simple.Node(5), T: simple.Node(6)},
		{F: simple.Node(5), T: simple.Node(7)},
		{F: simple.Node(6), T: simple.Node(7)},
	} {
		enneahedron1.SetEdge(e)
	}

	enneahedron2 := simple.NewUndirectedMatrix(8, 0, 0, 0)
	for _, e := range []simple.Edge{
		{F: simple.Node(0), T: simple.Node(1)},
		{F: simple.Node(0), T: simple.Node(2)},
		{F: simple.Node(0), T: simple.Node(6)},
		{F: simple.Node(1), T: simple.Node(2)},
		{F: simple.Node(1), T: simple.Node(3)},
		{F: simple.Node(1), T: simple.Node(7)},
		{F: simple.Node(2), T: simple.Node(3)},
		{F: simple.Node(2), T: simple.Node(4)},
		{F: simple.Node(2), T: simple.Node(5)},
		{F: simple.Node(3), T: simple.Node(5)},
		{F: simple.Node(4), T: simple.Node(5)},
		{F: simple.Node(4), T: simple.Node(6)},
		{F: simple.Node(5), T: simple.Node(6)},
		{F: simple.Node(5), T: simple.Node(7)},
		{F: simple.Node(6), T: simple.Node(7)},
	} {
		enneahedron2.SetEdge(e)
	}

	for _, en := range []*simple.UndirectedMatrix{
		enneahedron1,
		enneahedron2,
	} {
		var ed mat.EigenSym
		ed.Factorize(en.Matrix().(mat.Symmetric), false)
		fmt.Printf("%.2f\n", ed.Values(nil))
	}

	// Output:
	// [-2.41 -2.10 -1.25 -0.74 0.48 0.77 1.36 3.88]
	// [-2.41 -2.10 -1.25 -0.74 0.48 0.77 1.36 3.88]
}

func Example_connecting() {
	g := simple.NewUndirectedMatrix(8, 0, 0, 0)
	edges := []graph.Edge{
		nil,

		// Left lobe.
		simple.Edge{F: simple.Node(0), T: simple.Node(1)},
		simple.Edge{F: simple.Node(0), T: simple.Node(2)},
		simple.Edge{F: simple.Node(0), T: simple.Node(3)},
		simple.Edge{F: simple.Node(1), T: simple.Node(2)},
		simple.Edge{F: simple.Node(1), T: simple.Node(3)},
		simple.Edge{F: simple.Node(2), T: simple.Node(3)},

		// Right lobe.
		simple.Edge{F: simple.Node(4), T: simple.Node(5)},
		simple.Edge{F: simple.Node(4), T: simple.Node(6)},
		simple.Edge{F: simple.Node(4), T: simple.Node(7)},
		simple.Edge{F: simple.Node(5), T: simple.Node(6)},
		simple.Edge{F: simple.Node(5), T: simple.Node(7)},
		simple.Edge{F: simple.Node(6), T: simple.Node(7)},

		// Bridge.
		simple.Edge{F: simple.Node(0), T: simple.Node(4)},
	}
	fmt.Println("eigenvalues as edges are added:")
	for i, e := range edges {
		if e != nil {
			g.SetEdge(e)
		}

		l := spectral.NewLaplacian(g)
		var ed mat.EigenSym
		ed.Factorize(l.Matrix.(mat.Symmetric), i == len(edges)-1)
		vals := ed.Values(nil)
		for i, v := range vals {
			// Zero-out near zero values.
			if math.Abs(v) < 1e-15 {
				vals[i] = 0
			}
		}
		fmt.Printf(" %2.1f\n", vals)

		if i == len(edges)-1 {
			var vecs mat.Dense
			ed.VectorsTo(&vecs)
			fiedler := vecs.ColView(1)
			fmt.Println("Fiedler vector after joining lobes:")
			fmt.Printf(" %2.2f\n", mat.Formatted(fiedler.T()))
		}

	}

	// Output:
	// eigenvalues as edges are added:
	//  [0.0 0.0 0.0 0.0 0.0 0.0 0.0 0.0]
	//  [0.0 0.0 0.0 0.0 0.0 0.0 0.0 2.0]
	//  [0.0 0.0 0.0 0.0 0.0 0.0 1.0 3.0]
	//  [0.0 0.0 0.0 0.0 0.0 1.0 1.0 4.0]
	//  [0.0 0.0 0.0 0.0 0.0 1.0 3.0 4.0]
	//  [0.0 0.0 0.0 0.0 0.0 2.0 4.0 4.0]
	//  [0.0 0.0 0.0 0.0 0.0 4.0 4.0 4.0]
	//  [0.0 0.0 0.0 0.0 2.0 4.0 4.0 4.0]
	//  [0.0 0.0 0.0 1.0 3.0 4.0 4.0 4.0]
	//  [0.0 0.0 1.0 1.0 4.0 4.0 4.0 4.0]
	//  [0.0 0.0 1.0 3.0 4.0 4.0 4.0 4.0]
	//  [0.0 0.0 2.0 4.0 4.0 4.0 4.0 4.0]
	//  [0.0 0.0 4.0 4.0 4.0 4.0 4.0 4.0]
	//  [0.0 0.4 4.0 4.0 4.0 4.0 4.0 5.6]
	// Fiedler vector after joining lobes:
	//  [ 0.25   0.38   0.38   0.38  -0.25  -0.38  -0.38  -0.38]
}
