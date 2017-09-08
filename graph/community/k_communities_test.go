// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package community

import (
	"reflect"
	"sort"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/ordered"
	"gonum.org/v1/gonum/graph/simple"
)

// batageljZaversnikGraph is the example graph from
// figure 1 of http://arxiv.org/abs/cs/0310049v1
var batageljZaversnikGraph = []intset{
	0: nil,

	1: linksTo(2, 3),
	2: linksTo(4),
	3: linksTo(4),
	4: linksTo(5),
	5: nil,

	6:  linksTo(7, 8, 14),
	7:  linksTo(8, 11, 12, 14),
	8:  linksTo(14),
	9:  linksTo(11),
	10: linksTo(11),
	11: linksTo(12),
	12: linksTo(18),
	13: linksTo(14, 15),
	14: linksTo(15, 17),
	15: linksTo(16, 17),
	16: nil,
	17: linksTo(18, 19, 20),
	18: linksTo(19, 20),
	19: linksTo(20),
	20: nil,
}

var kCliqueCommunitiesTests = []struct {
	g    []intset
	k    int
	want [][]graph.Node
}{
	{
		g: []intset{
			0: linksTo(1, 2, 4, 6),
			1: linksTo(2, 4, 6),
			2: linksTo(3, 6),
			3: linksTo(4, 5),
			4: linksTo(6),
			5: nil,
			6: nil,
		},
		k: 3,
		want: [][]graph.Node{
			{simple.Node(0), simple.Node(1), simple.Node(2), simple.Node(4), simple.Node(6)},
			{simple.Node(3)},
			{simple.Node(5)},
		},
	},
	{
		g: batageljZaversnikGraph,
		k: 3,
		want: [][]graph.Node{
			{simple.Node(0)},
			{simple.Node(1)},
			{simple.Node(2)},
			{simple.Node(3)},
			{simple.Node(4)},
			{simple.Node(5)},
			{simple.Node(6), simple.Node(7), simple.Node(8), simple.Node(14)},
			{simple.Node(7), simple.Node(11), simple.Node(12)},
			{simple.Node(9)},
			{simple.Node(10)},
			{simple.Node(13), simple.Node(14), simple.Node(15), simple.Node(17)},
			{simple.Node(16)},
			{simple.Node(17), simple.Node(18), simple.Node(19), simple.Node(20)},
		},
	},
	{
		g: batageljZaversnikGraph,
		k: 4,
		want: [][]graph.Node{
			{simple.Node(0)},
			{simple.Node(1)},
			{simple.Node(2)},
			{simple.Node(3)},
			{simple.Node(4)},
			{simple.Node(5)},
			{simple.Node(6), simple.Node(7), simple.Node(8), simple.Node(14)},
			{simple.Node(9)},
			{simple.Node(10)},
			{simple.Node(11)},
			{simple.Node(12)},
			{simple.Node(13)},
			{simple.Node(15)},
			{simple.Node(16)},
			{simple.Node(17), simple.Node(18), simple.Node(19), simple.Node(20)},
		},
	},
}

func TestKCliqueCommunities(t *testing.T) {
	for _, test := range kCliqueCommunitiesTests {
		g := simple.NewUndirectedGraph()
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if !g.Has(simple.Node(u)) {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}
		got := KCliqueCommunities(test.k, g)

		for _, c := range got {
			sort.Sort(ordered.ByID(c))
		}
		sort.Sort(ordered.BySliceIDs(got))

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("unexpected k-connected components:\ngot: %v\nwant:%v", got, test.want)
		}
	}
}
