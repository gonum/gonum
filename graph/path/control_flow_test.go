// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"reflect"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

var dominatorsTests = []struct {
	n     graph.Node
	edges []simple.Edge

	want map[int64]graph.Node
}{
	{ // Example from Lengauer and Tarjan, fig 1.
		n: char('R'),
		edges: []simple.Edge{
			{F: char('A'), T: char('D')},
			{F: char('B'), T: char('A')}, // Direction inferred from fig 3.
			{F: char('B'), T: char('D')},
			{F: char('B'), T: char('E')},
			{F: char('C'), T: char('F')},
			{F: char('C'), T: char('G')},
			{F: char('D'), T: char('L')},
			{F: char('E'), T: char('H')},
			{F: char('F'), T: char('I')},
			{F: char('G'), T: char('I')},
			{F: char('G'), T: char('J')},
			{F: char('H'), T: char('E')},
			{F: char('H'), T: char('K')},
			{F: char('I'), T: char('K')},
			{F: char('J'), T: char('I')},
			{F: char('K'), T: char('I')},
			{F: char('K'), T: char('R')},
			{F: char('L'), T: char('H')},
			{F: char('R'), T: char('A')},
			{F: char('R'), T: char('B')}, // Direction inferred from fig 2.
			{F: char('R'), T: char('C')},
		},

		want: map[int64]graph.Node{
			'A': char('R'),
			'B': char('R'),
			'C': char('R'),
			'D': char('R'),
			'E': char('R'),
			'F': char('C'),
			'G': char('C'),
			'H': char('R'),
			'I': char('R'),
			'J': char('G'),
			'K': char('R'),
			'L': char('D'),
		},
	},
	{ // WP example: https://en.wikipedia.org/w/index.php?title=Dominator_(graph_theory)&oldid=758099236.
		n: simple.Node(1),
		edges: []simple.Edge{
			{F: simple.Node(1), T: simple.Node(2)},
			{F: simple.Node(2), T: simple.Node(3)},
			{F: simple.Node(2), T: simple.Node(4)},
			{F: simple.Node(2), T: simple.Node(6)},
			{F: simple.Node(3), T: simple.Node(5)},
			{F: simple.Node(4), T: simple.Node(5)},
			{F: simple.Node(5), T: simple.Node(2)},
		},

		want: map[int64]graph.Node{
			2: simple.Node(1),
			3: simple.Node(2),
			4: simple.Node(2),
			5: simple.Node(2),
			6: simple.Node(2),
		},
	},
	{ // WP example with node IDs decremented by 1.
		n: simple.Node(0),
		edges: []simple.Edge{
			{F: simple.Node(0), T: simple.Node(1)},
			{F: simple.Node(1), T: simple.Node(2)},
			{F: simple.Node(1), T: simple.Node(3)},
			{F: simple.Node(1), T: simple.Node(5)},
			{F: simple.Node(2), T: simple.Node(4)},
			{F: simple.Node(3), T: simple.Node(4)},
			{F: simple.Node(4), T: simple.Node(1)},
		},

		want: map[int64]graph.Node{
			1: simple.Node(0),
			2: simple.Node(1),
			3: simple.Node(1),
			4: simple.Node(1),
			5: simple.Node(1),
		},
	},
}

type char int64

func (n char) ID() int64      { return int64(n) }
func (n char) String() string { return string(n) }

func TestDominators(t *testing.T) {
	for _, test := range dominatorsTests {
		g := simple.NewDirectedGraph()
		for _, e := range test.edges {
			g.SetEdge(e)
		}

		got := Dominators(test.n, g)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("unexpected dominator tree: got:%v want:%v", got, test.want)
		}
	}
}
