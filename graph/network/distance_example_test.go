// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network_test

import (
	"fmt"

	"gonum.org/v1/gonum/graph/network"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

func ExampleEccentricity_diameter() {
	// Build a simple graph: n1--n2--n3.
	var n1, n2, n3 simple.Node = 1, 2, 3
	g := simple.NewUndirectedGraph()
	g.SetEdge(simple.Edge{F: n1, T: n2})
	g.SetEdge(simple.Edge{F: n2, T: n3})

	// Find shortest paths and calculate eccentricity values.
	paths := path.DijkstraAllPaths(g)
	e := network.Eccentricity(g, paths)

	// Compute the graph's diameter from the eccentricities.
	var diameter float64
	for _, d := range e {
		diameter = max(d, diameter)
	}
	fmt.Printf("diameter = %v", diameter)

	// Output: diameter = 2
}
