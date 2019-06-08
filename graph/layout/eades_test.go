// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package layout

import (
	"path/filepath"
	"reflect"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/spatial/r2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

var eadesR2Tests = []struct {
	name  string
	g     graph.Graph
	param EadesR2
	want  map[int64]r2.Vec
}{
	{
		name: "line",
		g: func() graph.Graph {
			edges := []simple.Edge{
				{simple.Node(0), simple.Node(1)},
			}
			g := simple.NewUndirectedGraph()
			for _, e := range edges {
				g.SetEdge(e)
			}
			return orderedGraph{g}
		}(),
		param: EadesR2{C1: 2, C2: 1, C3: 1, C4: 0.1, M: 100, Theta: 0.1, Src: rand.NewSource(1)},
		want: map[int64]r2.Vec{
			0: {X: 0.05965621545523485, Y: 1.1222058498562848},
			1: {X: 1.0260524834222655, Y: 0.08311554782816481},
		},
	},
	{
		name: "square",
		g: func() graph.Graph {
			edges := []simple.Edge{
				{simple.Node(0), simple.Node(1)},
				{simple.Node(0), simple.Node(2)},
				{simple.Node(1), simple.Node(3)},
				{simple.Node(2), simple.Node(3)},
			}
			g := simple.NewUndirectedGraph()
			for _, e := range edges {
				g.SetEdge(e)
			}
			return orderedGraph{g}
		}(),
		param: EadesR2{C1: 2, C2: 1, C3: 1, C4: 0.1, M: 100, Theta: 0.1, Src: rand.NewSource(1)},
		want: map[int64]r2.Vec{
			0: {X: 0.010610173427526952, Y: 1.396249314504884},
			1: {X: 1.4348226226119472, Y: 0.9351207866753373},
			2: {X: -0.45010867679614364, Y: -0.028095714144673396},
			3: {X: 0.9741037723882764, Y: -0.4892242419742201},
		},
	},
	{
		name: "tetrahedron",
		g: func() graph.Graph {
			edges := []simple.Edge{
				{simple.Node(0), simple.Node(1)},
				{simple.Node(0), simple.Node(2)},
				{simple.Node(0), simple.Node(3)},
				{simple.Node(1), simple.Node(2)},
				{simple.Node(1), simple.Node(3)},
				{simple.Node(2), simple.Node(3)},
			}
			g := simple.NewUndirectedGraph()
			for _, e := range edges {
				g.SetEdge(e)
			}
			return orderedGraph{g}
		}(),
		param: EadesR2{C1: 2, C2: 1, C3: 1, C4: 0.1, M: 100, Theta: 0.1, Src: rand.NewSource(1)},
		want: map[int64]r2.Vec{
			0: {X: 0.05418735696054162, Y: 1.1968877215637272},
			1: {X: 1.2357321582062974, Y: 0.8916821522126922},
			2: {X: -0.2510182123904937, Y: 0.015342920317971562},
			3: {X: 0.9305265888552621, Y: -0.28986264903306364},
		},
	},
	{
		name: "sheet",
		g: func() graph.Graph {
			edges := []simple.Edge{
				{simple.Node(0), simple.Node(1)},
				{simple.Node(0), simple.Node(3)},
				{simple.Node(1), simple.Node(2)},
				{simple.Node(1), simple.Node(4)},
				{simple.Node(2), simple.Node(5)},
				{simple.Node(3), simple.Node(4)},
				{simple.Node(3), simple.Node(6)},
				{simple.Node(4), simple.Node(5)},
				{simple.Node(4), simple.Node(7)},
				{simple.Node(5), simple.Node(8)},
				{simple.Node(6), simple.Node(7)},
				{simple.Node(7), simple.Node(8)},
			}
			g := simple.NewUndirectedGraph()
			for _, e := range edges {
				g.SetEdge(e)
			}
			return orderedGraph{g}
		}(),
		param: EadesR2{C1: 2, C2: 1, C3: 1, C4: 0.1, M: 100, Theta: 0.1, Src: rand.NewSource(1)},
		want: map[int64]r2.Vec{
			0: {X: -0.9784723315844219, Y: 2.4228759281103347},
			1: {X: 0.6046747521369412, Y: 2.148857452899514},
			2: {X: 2.1716142408195864, Y: 1.7905903669465866},
			3: {X: -1.2641103194603636, Y: 0.829651744350034},
			4: {X: 0.3593028500368119, Y: 0.5174079857268895},
			5: {X: 1.98222256702394, Y: 0.2123640542636961},
			6: {X: -1.2873221144771265, Y: -0.7884352233940609},
			7: {X: 0.2919205689515562, Y: -1.133168729322618},
			8: {X: 1.8947514938889922, Y: -1.3761752981118403},
		},
	},
	{
		name: "tube",
		g: func() graph.Graph {
			edges := []simple.Edge{
				{simple.Node(0), simple.Node(1)},
				{simple.Node(0), simple.Node(2)},
				{simple.Node(0), simple.Node(3)},
				{simple.Node(1), simple.Node(2)},
				{simple.Node(1), simple.Node(4)},
				{simple.Node(2), simple.Node(5)},
				{simple.Node(3), simple.Node(4)},
				{simple.Node(3), simple.Node(5)},
				{simple.Node(3), simple.Node(6)},
				{simple.Node(4), simple.Node(5)},
				{simple.Node(4), simple.Node(7)},
				{simple.Node(5), simple.Node(8)},
				{simple.Node(6), simple.Node(7)},
				{simple.Node(6), simple.Node(8)},
				{simple.Node(7), simple.Node(8)},
			}
			g := simple.NewUndirectedGraph()
			for _, e := range edges {
				g.SetEdge(e)
			}
			return orderedGraph{g}
		}(),
		param: EadesR2{C1: 2, C2: 1, C3: 1, C4: 0.1, M: 100, Theta: 0.1, Src: rand.NewSource(1)},
		want: map[int64]r2.Vec{
			0: {X: -0.14160080740028816, Y: 2.13495397568788},
			1: {X: 0.9627848362888625, Y: 2.825431950959263},
			2: {X: 1.6434481244516415, Y: 1.7154427454691445},
			3: {X: -0.5622664107092017, Y: 0.4508021982264986},
			4: {X: 0.5469163170006696, Y: 1.0668505906376997},
			5: {X: 1.2688450976607686, Y: 0.023173378675218948},
			6: {X: -0.8615297878037711, Y: -1.3069652073547353},
			7: {X: 0.15539891254606483, Y: -0.6029115961279473},
			8: {X: 0.7625404743124667, Y: -1.682810696307514},
		},
	},
}

func TestEadesR2(t *testing.T) {
	for _, test := range eadesR2Tests {
		eades := test.param
		o := NewOptimizerR2(test.g, eades.Update)
		var n int
		for o.Update() {
			n++
		}
		if n != test.param.M {
			t.Errorf("unexpected number of iterations for %q: got:%d want:%d", test.name, n, test.param.M)
		}
		got := make(map[int64]r2.Vec)
		nodes := test.g.Nodes()
		for nodes.Next() {
			id := nodes.Node().ID()
			got[id] = o.Coord2(id)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("unexpected result for %q:\ngot: %#v\nwant:%#v", test.name, got, test.want)
		}

		p, err := plot.New()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}
		p.Add(render{o})
		p.HideAxes()
		path := filepath.Join("testdata", test.name+".png")
		err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, path)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}
		checkRenderedLayout(t, path)
	}
}
