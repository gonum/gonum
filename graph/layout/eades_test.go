// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package layout

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/spatial/r2"
)

var eadesR2Tests = []struct {
	name  string
	g     graph.Graph
	param EadesR2
	want  map[int64]r2.Vec
}{
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
			0: r2.Vec{X: 2.15059197551694, Y: -10.220981772158567},
			1: r2.Vec{X: 3.2862425589944833, Y: -11.008330556318482},
			2: r2.Vec{X: 1.363536823864583, Y: -11.356835875434365},
			3: r2.Vec{X: 2.499187407342126, Y: -12.14418465959428},
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
			0: r2.Vec{X: 7.7334883284703375, Y: -11.484406140819413},
			1: r2.Vec{X: 8.425306545865576, Y: -11.140480488404148},
			2: r2.Vec{X: 7.976501845950402, Y: -12.514472252654656},
			3: r2.Vec{X: 9.400882882193617, Y: -11.63154357822779},
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

			0: r2.Vec{X: 14.475884753952018, Y: 11.87254908676763},
			1: r2.Vec{X: 16.08088846104896, Y: 11.888141530966085},
			2: r2.Vec{X: 17.38001523058331, Y: 11.878481798669878},
			3: r2.Vec{X: 14.356099897250903, Y: 13.473142461734563},
			4: r2.Vec{X: 16.01516712526384, Y: 13.54758041594208},
			5: r2.Vec{X: 17.373180959720063, Y: 13.490074894554255},
			6: r2.Vec{X: 14.236725269213267, Y: 14.76680859469247},
			7: r2.Vec{X: 15.843133632993109, Y: 14.895887859745487},
			8: r2.Vec{X: 17.14673527193477, Y: 14.779176044572363},
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
			0: r2.Vec{X: 20.81592592580858, Y: 7.093353250039021},
			1: r2.Vec{X: 21.726993818602256, Y: 7.055585705916057},
			2: r2.Vec{X: 21.810995821173716, Y: 5.783565022523372},
			3: r2.Vec{X: 22.0864584730445, Y: 8.300612169610107},
			4: r2.Vec{X: 23.00743865904811, Y: 8.184144796089953},
			5: r2.Vec{X: 23.20316516390141, Y: 6.828096609467102},
			6: r2.Vec{X: 23.007775564315455, Y: 9.319231864784014},
			7: r2.Vec{X: 23.97755074521381, Y: 9.266038395576885},
			8: r2.Vec{X: 24.15211052582012, Y: 7.809512061906321},
		},
	},
}

func TestEadesR2(t *testing.T) {
	// Factor to dilate positions so graphviz separates nodes.
	const scaling = 100

	for _, test := range eadesR2Tests {
		eades := test.param
		o := NewOptimizerR2(test.g, eades.Update)
		var n int
		for o.Update() {
			n++
		}
		if n != test.param.M {
			t.Errorf("unexpected number of iterations: got:%d want:%d", n, test.param.M)
		}
		got := make(map[int64]r2.Vec)
		nodes := test.g.Nodes()
		for nodes.Next() {
			id := nodes.Node().ID()
			got[id] = o.Coord2(id)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("unexpected result:\ngot: %#v\nwant:%#v", got, test.want)
		}

		// TODO(kortschak): Replace this with Go rendering code.
		if _, err := exec.LookPath("neato"); err != nil {
			continue
		}
		layout := simple.NewUndirectedGraph()
		for _, u := range graph.NodesOf(test.g.Nodes()) {
			for _, v := range graph.NodesOf(test.g.From(u.ID())) {
				if v.ID() < u.ID() {
					continue
				}
				layout.SetEdge(simple.Edge{
					F: positionNode{id: u.ID(), pos: got[u.ID()].Scale(scaling)},
					T: positionNode{id: v.ID(), pos: got[v.ID()].Scale(scaling)},
				})
			}
		}
		b, _ := dot.Marshal(layout, test.name, "", "  ")
		path := filepath.Join("testdata", test.name+".png")
		cmd := exec.Command("neato", "-n", "-Tpng", "-o", path)
		cmd.Stdin = bytes.NewReader(b)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("unexpected error: %v\n%s", err, out)
			continue
		}
		checkRenderedLayout(t, path)
	}
}
