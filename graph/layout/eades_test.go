// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package layout_test

import (
	"path/filepath"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/spatial/r2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"

	. "gonum.org/v1/gonum/graph/layout"
)

func TestEadesR2(t *testing.T) {
	eadesR2Tests := []struct {
		name      string
		g         graph.Graph
		param     EadesR2
		wantIters int
	}{
		{
			name: "line",
			g: func() graph.Graph {
				edges := []simple.Edge{
					{F: simple.Node(0), T: simple.Node(1)},
				}
				g := simple.NewUndirectedGraph()
				for _, e := range edges {
					g.SetEdge(e)
				}
				return orderedGraph{g}
			}(),
			param:     EadesR2{Repulsion: 1, Rate: 0.1, Updates: 100, Theta: 0.1, Src: rand.NewSource(1)},
			wantIters: 100,
		},
		{
			name: "square",
			g: func() graph.Graph {
				edges := []simple.Edge{
					{F: simple.Node(0), T: simple.Node(1)},
					{F: simple.Node(0), T: simple.Node(2)},
					{F: simple.Node(1), T: simple.Node(3)},
					{F: simple.Node(2), T: simple.Node(3)},
				}
				g := simple.NewUndirectedGraph()
				for _, e := range edges {
					g.SetEdge(e)
				}
				return orderedGraph{g}
			}(),
			param:     EadesR2{Repulsion: 1, Rate: 0.1, Updates: 100, Theta: 0.1, Src: rand.NewSource(1)},
			wantIters: 100,
		},
		{
			name: "tetrahedron",
			g: func() graph.Graph {
				edges := []simple.Edge{
					{F: simple.Node(0), T: simple.Node(1)},
					{F: simple.Node(0), T: simple.Node(2)},
					{F: simple.Node(0), T: simple.Node(3)},
					{F: simple.Node(1), T: simple.Node(2)},
					{F: simple.Node(1), T: simple.Node(3)},
					{F: simple.Node(2), T: simple.Node(3)},
				}
				g := simple.NewUndirectedGraph()
				for _, e := range edges {
					g.SetEdge(e)
				}
				return orderedGraph{g}
			}(),
			param:     EadesR2{Repulsion: 1, Rate: 0.1, Updates: 100, Theta: 0.1, Src: rand.NewSource(1)},
			wantIters: 100,
		},
		{
			name: "sheet",
			g: func() graph.Graph {
				edges := []simple.Edge{
					{F: simple.Node(0), T: simple.Node(1)},
					{F: simple.Node(0), T: simple.Node(3)},
					{F: simple.Node(1), T: simple.Node(2)},
					{F: simple.Node(1), T: simple.Node(4)},
					{F: simple.Node(2), T: simple.Node(5)},
					{F: simple.Node(3), T: simple.Node(4)},
					{F: simple.Node(3), T: simple.Node(6)},
					{F: simple.Node(4), T: simple.Node(5)},
					{F: simple.Node(4), T: simple.Node(7)},
					{F: simple.Node(5), T: simple.Node(8)},
					{F: simple.Node(6), T: simple.Node(7)},
					{F: simple.Node(7), T: simple.Node(8)},
				}
				g := simple.NewUndirectedGraph()
				for _, e := range edges {
					g.SetEdge(e)
				}
				return orderedGraph{g}
			}(),
			param:     EadesR2{Repulsion: 1, Rate: 0.1, Updates: 100, Theta: 0.1, Src: rand.NewSource(1)},
			wantIters: 100,
		},
		{
			name: "tube",
			g: func() graph.Graph {
				edges := []simple.Edge{
					{F: simple.Node(0), T: simple.Node(1)},
					{F: simple.Node(0), T: simple.Node(2)},
					{F: simple.Node(0), T: simple.Node(3)},
					{F: simple.Node(1), T: simple.Node(2)},
					{F: simple.Node(1), T: simple.Node(4)},
					{F: simple.Node(2), T: simple.Node(5)},
					{F: simple.Node(3), T: simple.Node(4)},
					{F: simple.Node(3), T: simple.Node(5)},
					{F: simple.Node(3), T: simple.Node(6)},
					{F: simple.Node(4), T: simple.Node(5)},
					{F: simple.Node(4), T: simple.Node(7)},
					{F: simple.Node(5), T: simple.Node(8)},
					{F: simple.Node(6), T: simple.Node(7)},
					{F: simple.Node(6), T: simple.Node(8)},
					{F: simple.Node(7), T: simple.Node(8)},
				}
				g := simple.NewUndirectedGraph()
				for _, e := range edges {
					g.SetEdge(e)
				}
				return orderedGraph{g}
			}(),
			param:     EadesR2{Repulsion: 1, Rate: 0.1, Updates: 100, Theta: 0.1, Src: rand.NewSource(1)},
			wantIters: 100,
		},
		{
			// This test does not produce a good layout, but is here to
			// ensure that Update does not panic with steep decent rates.
			name: "tube-steep",
			g: func() graph.Graph {
				edges := []simple.Edge{
					{F: simple.Node(0), T: simple.Node(1)},
					{F: simple.Node(0), T: simple.Node(2)},
					{F: simple.Node(0), T: simple.Node(3)},
					{F: simple.Node(1), T: simple.Node(2)},
					{F: simple.Node(1), T: simple.Node(4)},
					{F: simple.Node(2), T: simple.Node(5)},
					{F: simple.Node(3), T: simple.Node(4)},
					{F: simple.Node(3), T: simple.Node(5)},
					{F: simple.Node(3), T: simple.Node(6)},
					{F: simple.Node(4), T: simple.Node(5)},
					{F: simple.Node(4), T: simple.Node(7)},
					{F: simple.Node(5), T: simple.Node(8)},
					{F: simple.Node(6), T: simple.Node(7)},
					{F: simple.Node(6), T: simple.Node(8)},
					{F: simple.Node(7), T: simple.Node(8)},
				}
				g := simple.NewUndirectedGraph()
				for _, e := range edges {
					g.SetEdge(e)
				}
				return orderedGraph{g}
			}(),
			param:     EadesR2{Repulsion: 1, Rate: 1, Updates: 100, Theta: 0.1, Src: rand.NewSource(1)},
			wantIters: 99,
		},

		{
			name: "wp_page", // https://en.wikipedia.org/wiki/PageRank#/media/File:PageRanks-Example.jpg
			g: func() graph.Graph {
				edges := []simple.Edge{
					{F: simple.Node(0), T: simple.Node(3)},
					{F: simple.Node(1), T: simple.Node(2)},
					{F: simple.Node(1), T: simple.Node(3)},
					{F: simple.Node(1), T: simple.Node(4)},
					{F: simple.Node(1), T: simple.Node(5)},
					{F: simple.Node(1), T: simple.Node(6)},
					{F: simple.Node(1), T: simple.Node(7)},
					{F: simple.Node(1), T: simple.Node(8)},
					{F: simple.Node(3), T: simple.Node(4)},
					{F: simple.Node(4), T: simple.Node(5)},
					{F: simple.Node(4), T: simple.Node(6)},
					{F: simple.Node(4), T: simple.Node(7)},
					{F: simple.Node(4), T: simple.Node(8)},
					{F: simple.Node(4), T: simple.Node(9)},
					{F: simple.Node(4), T: simple.Node(10)},
				}
				g := simple.NewUndirectedGraph()
				for _, e := range edges {
					g.SetEdge(e)
				}
				return orderedGraph{g}
			}(),
			param:     EadesR2{Repulsion: 1, Rate: 0.1, Updates: 100, Theta: 0.1, Src: rand.NewSource(1)},
			wantIters: 100,
		},
	}

	for _, test := range eadesR2Tests {
		eades := test.param
		o := NewOptimizerR2(test.g, eades.Update)
		var n int
		for o.Update() {
			n++
		}
		if n != test.wantIters {
			t.Errorf("unexpected number of iterations for %q: got:%d want:%d", test.name, n, test.wantIters)
		}

		p := plot.New()
		p.Add(render{o})
		p.HideAxes()
		path := filepath.Join("testdata", test.name+".png")
		err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, path)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}
		ok := checkRenderedLayout(t, path)
		if !ok {
			got := make(map[int64]r2.Vec)
			nodes := test.g.Nodes()
			for nodes.Next() {
				id := nodes.Node().ID()
				got[id] = o.Coord2(id)
			}
			t.Logf("got node positions: %#v", got)
		}
	}
}
