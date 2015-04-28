// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traverse_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/traverse"
)

var breadthFirstTests = []struct {
	g     []set
	from  graph.Node
	edge  func(graph.Edge) bool
	until func(graph.Node, int) bool
	final map[graph.Node]bool
	want  [][]int
}{
	{
		g: []set{
			0: linksTo(1, 4),
			1: linksTo(2, 4),
			2: linksTo(3),
			3: linksTo(4, 5),
			4: nil,
			5: nil,
		},
		from:  concrete.Node(1),
		final: map[graph.Node]bool{nil: true},
		want: [][]int{
			{1},
			{0, 2, 4},
			{3},
			{5},
		},
	},
	{
		g: []set{
			0: linksTo(1, 4),
			1: linksTo(2, 4),
			2: linksTo(3),
			3: linksTo(4, 5),
			4: nil,
			5: nil,
		},
		edge: func(e graph.Edge) bool {
			// Do not traverse an edge between 3 and 5.
			return (e.Head().ID() != 3 || e.Tail().ID() != 5) && (e.Head().ID() != 5 || e.Tail().ID() != 3)
		},
		from:  concrete.Node(1),
		final: map[graph.Node]bool{nil: true},
		want: [][]int{
			{1},
			{0, 2, 4},
			{3},
		},
	},
	{
		g: []set{
			0: linksTo(1, 4),
			1: linksTo(2, 4),
			2: linksTo(3),
			3: linksTo(4, 5),
			4: nil,
			5: nil,
		},
		from:  concrete.Node(1),
		until: func(n graph.Node, _ int) bool { return n == concrete.Node(3) },
		final: map[graph.Node]bool{concrete.Node(3): true},
		want: [][]int{
			{1},
			{0, 2, 4},
		},
	},
	{
		g: []set{
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
		},
		from:  concrete.Node(13),
		final: map[graph.Node]bool{nil: true},
		want: [][]int{
			{13},
			{14, 15},
			{6, 7, 8, 16, 17},
			{11, 12, 18, 19, 20},
			{9, 10},
		},
	},
	{
		g: []set{
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
		},
		from:  concrete.Node(13),
		until: func(_ graph.Node, d int) bool { return d > 2 },
		final: map[graph.Node]bool{
			concrete.Node(11): true,
			concrete.Node(12): true,
			concrete.Node(18): true,
			concrete.Node(19): true,
			concrete.Node(20): true,
		},
		want: [][]int{
			{13},
			{14, 15},
			{6, 7, 8, 16, 17},
		},
	},
}

func TestBreadthFirst(t *testing.T) {
	for i, test := range breadthFirstTests {
		g := concrete.NewGraph()
		for u, e := range test.g {
			if !g.NodeExists(concrete.Node(u)) {
				g.AddNode(concrete.Node(u))
			}
			for v := range e {
				if !g.NodeExists(concrete.Node(v)) {
					g.AddNode(concrete.Node(v))
				}
				g.AddUndirectedEdge(concrete.Edge{H: concrete.Node(u), T: concrete.Node(v)}, 0)
			}
		}
		w := traverse.BreadthFirst{
			EdgeFilter: test.edge,
		}
		var got [][]int
		final := w.Walk(g, test.from, func(n graph.Node, d int) bool {
			if test.until != nil && test.until(n, d) {
				return true
			}
			if d >= len(got) {
				got = append(got, []int(nil))
			}
			got[d] = append(got[d], n.ID())
			return false
		})
		if !test.final[final] {
			t.Errorf("unexepected final node for test %d:\ngot:  %v\nwant: %v", i, final, test.final)
		}
		for _, l := range got {
			sort.Ints(l)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("unexepected BFS level structure for test %d:\ngot:  %v\nwant: %v", i, got, test.want)
		}
	}
}

var depthFirstTests = []struct {
	g     []set
	from  graph.Node
	edge  func(graph.Edge) bool
	until func(graph.Node) bool
	final map[graph.Node]bool
	want  []int
}{
	{
		g: []set{
			0: linksTo(1, 4),
			1: linksTo(2, 4),
			2: linksTo(3),
			3: linksTo(4, 5),
			4: nil,
			5: nil,
		},
		from:  concrete.Node(1),
		final: map[graph.Node]bool{nil: true},
		want:  []int{0, 1, 2, 3, 4, 5},
	},
	{
		g: []set{
			0: linksTo(1, 4),
			1: linksTo(2, 4),
			2: linksTo(3),
			3: linksTo(4, 5),
			4: nil,
			5: nil,
		},
		edge: func(e graph.Edge) bool {
			// Do not traverse an edge between 3 and 5.
			return (e.Head().ID() != 3 || e.Tail().ID() != 5) && (e.Head().ID() != 5 || e.Tail().ID() != 3)
		},
		from:  concrete.Node(1),
		final: map[graph.Node]bool{nil: true},
		want:  []int{0, 1, 2, 3, 4},
	},
	{
		g: []set{
			0: linksTo(1, 4),
			1: linksTo(2, 4),
			2: linksTo(3),
			3: linksTo(4, 5),
			4: nil,
			5: nil,
		},
		from:  concrete.Node(1),
		until: func(n graph.Node) bool { return n == concrete.Node(3) },
		final: map[graph.Node]bool{concrete.Node(3): true},
	},
	{
		g: []set{
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
		},
		from:  concrete.Node(0),
		final: map[graph.Node]bool{nil: true},
		want:  []int{0},
	},
	{
		g: []set{
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
		},
		from:  concrete.Node(3),
		final: map[graph.Node]bool{nil: true},
		want:  []int{1, 2, 3, 4, 5},
	},
	{
		g: []set{
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
		},
		from:  concrete.Node(13),
		final: map[graph.Node]bool{nil: true},
		want:  []int{6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
	},
}

func TestDepthFirst(t *testing.T) {
	for i, test := range depthFirstTests {
		g := concrete.NewGraph()
		for u, e := range test.g {
			if !g.NodeExists(concrete.Node(u)) {
				g.AddNode(concrete.Node(u))
			}
			for v := range e {
				if !g.NodeExists(concrete.Node(v)) {
					g.AddNode(concrete.Node(v))
				}
				g.AddUndirectedEdge(concrete.Edge{H: concrete.Node(u), T: concrete.Node(v)}, 0)
			}
		}
		w := traverse.DepthFirst{
			EdgeFilter: test.edge,
		}
		var got []int
		final := w.Walk(g, test.from, func(n graph.Node) bool {
			if test.until != nil && test.until(n) {
				return true
			}
			got = append(got, n.ID())
			return false
		})
		if !test.final[final] {
			t.Errorf("unexepected final node for test %d:\ngot:  %v\nwant: %v", i, final, test.final)
		}
		sort.Ints(got)
		if test.want != nil && !reflect.DeepEqual(got, test.want) {
			t.Errorf("unexepected DFS traversed nodes for test %d:\ngot:  %v\nwant: %v", i, got, test.want)
		}
	}
}

// set is an integer set.
type set map[int]struct{}

func linksTo(i ...int) set {
	if len(i) == 0 {
		return nil
	}
	s := make(set)
	for _, v := range i {
		s[v] = struct{}{}
	}
	return s
}
