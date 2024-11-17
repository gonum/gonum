// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flow

import (
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

var intervalTests = []struct {
	name  string
	edges []struct{ from, to simple.Node }
	root  int64
	want  []map[int64]intset
}{
	{
		// Graph from C. Cifuentes, "Reverse Compilation Techniques", 1994 (figure 6-23)
		// Available via https://eprints.qut.edu.au/36820/
		name: "cifuentes",
		edges: []struct{ from, to simple.Node }{
			{1, 2},
			{1, 5},
			{2, 3},
			{2, 4},
			{3, 5},
			{4, 5},
			{5, 6},
			{6, 7},
			{6, 12},
			{7, 8},
			{7, 9},
			{8, 9},
			{8, 10},
			{9, 10},
			{10, 11},
			{12, 13},
			{14, 13},
			{13, 14},
			{14, 15},
			{15, 6},
		},
		root: 1,
		want: []map[int64]intset{
			{
				1: linksTo(2, 5),
				2: linksTo(3, 4),
				3: linksTo(5),
				4: linksTo(5),
				5: nil,
			},
			{
				6:  linksTo(7, 12),
				7:  linksTo(8, 9),
				8:  linksTo(9, 10),
				9:  linksTo(10),
				10: linksTo(11),
				11: nil,
				12: nil,
			},
			{
				13: linksTo(14),
				14: linksTo(13, 15),
				15: nil,
			},
		},
	},
}

func TestInterval(t *testing.T) {
	for _, test := range intervalTests {
		t.Run(test.name, func(t *testing.T) {
			g := simple.NewDirectedGraph()
			for _, e := range test.edges {
				g.SetEdge(simple.Edge{F: e.from, T: e.to})
			}

			ig := Intervals(g, test.root)
			if len(ig.Intervals) != len(test.want) {
				t.Fatalf("unexpected interval count: got:%d want:%d", len(ig.Intervals), len(test.want))
			}

			for i, iv := range ig.Intervals {
				var got graph.Graph = iv
				want := gFromIntsets(test.want[i])
				if !topo.Equal(got, want) {
					gotDot, err := dot.Marshal(iv, "", "", "\t")
					if err != nil {
						t.Fatalf("unexpected error marshalling got DOT: %v", err)
					}
					wantDot, err := dot.Marshal(want, "", "", "\t")
					if err != nil {
						t.Fatalf("unexpected error marshalling want DOT: %v", err)
					}
					t.Errorf("unexpected topology of interval %d:\ngot:\n%s\nwant:\n%s", i, gotDot, wantDot)
				}
			}
		})
	}
}

func gFromIntsets(s map[int64]intset) graph.Directed {
	g := simple.NewDirectedGraph()
	for u, e := range s {
		// Add nodes that are not defined by an edge.
		if g.Node(int64(u)) == nil {
			g.AddNode(simple.Node(u))
		}

		for v := range e {
			g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
		}
	}
	return g
}

// intset is an integer set.
type intset map[int64]struct{}

func linksTo(i ...int64) intset {
	if len(i) == 0 {
		return nil
	}
	s := make(intset)
	for _, v := range i {
		s[v] = struct{}{}
	}
	return s
}
