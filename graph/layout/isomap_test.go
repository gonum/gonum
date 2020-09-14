// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package layout

import (
	"path/filepath"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/spatial/r2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

var (
	// tag is modified in isomap_noasm_test.go to "_noasm" when any
	// build tag prevents use of the assembly numerical kernels.
	tag string

	// arch is modified in isomap_arm64_test.go to "_arm64" on arm64
	// and "_386" on 386 to allow differences in numerical precision
	// to be allowed for.
	arch string
)

var isomapR2Tests = []struct {
	name string
	g    graph.Graph
}{
	{
		name: "line_isomap",
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
	},
	{
		name: "square_isomap",
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
	},
	{
		name: "tetrahedron_isomap",
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
	},
	{
		name: "sheet_isomap",
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
	},
	{
		name: "tube_isomap",
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
	},
	{
		name: "wp_page_isomap", // https://en.wikipedia.org/wiki/PageRank#/media/File:PageRanks-Example.jpg
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
	},
}

func TestIsomapR2(t *testing.T) {
	for _, test := range isomapR2Tests {
		o := NewOptimizerR2(test.g, IsomapR2{}.Update)
		var n int
		for o.Update() {
			n++
		}
		p, err := plot.New()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}
		p.Add(render{o})
		p.HideAxes()
		path := filepath.Join("testdata", test.name+tag+arch+".png")
		err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, path)
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
