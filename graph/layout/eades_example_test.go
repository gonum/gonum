// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package layout

import (
	"fmt"
	"path/filepath"

	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

func ExampleEadesR2() {

	g := makeCompleteGraph(6)

	// Graph layout.
	// Explained in "A heuristic for graph drawing",
	// Congressus numerantium 42:149-160.
	// Experiment with these values for more/less nodes.
	eades := EadesR2{Repulsion: 1, Rate: 0.05, Updates: 30, Theta: 0.2}

	// Contains graph, layout and updater function.
	optimizer := NewOptimizerR2(g, eades.Update)

	// Reposition nodes until eades.Updates == 0
	// by calling layout updater in Update method.
	for optimizer.Update() {

	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	// Add to plot.
	p.Add(render{optimizer})
	p.HideAxes()

	path := filepath.Join("testdata", "k6_eades.png")

	// Render graph on save.
	if err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, path); err != nil {
		panic(err)
	}

	fmt.Println("Saved plot to testdata/k6_eades.png")

	// Output: Saved plot to testdata/k6_eades.png
}

// Each node is connected to all other nodes
func makeCompleteGraph(n int) *simple.UndirectedGraph {
	g := simple.NewUndirectedGraph()

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			g.SetEdge(g.NewEdge(simple.Node(i), simple.Node(j)))
		}
	}

	return g
}
