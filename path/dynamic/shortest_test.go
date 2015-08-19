// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dynamic

import (
	"fmt"
	"math"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

func init() {
	for _, test := range shortestPathTests {
		if len(test.want) != 1 && test.unique {
			panic(fmt.Sprintf("%q: bad shortest path test: non-unique paths marked unique", test.name))
		}
	}
}

// shortestPathTests are graphs used to test DStarLite.
var shortestPathTests = []struct {
	name             string
	g                func() graph.Mutable
	edges            []concrete.Edge
	negative         bool
	hasNegativeCycle bool

	query  concrete.Edge
	weight float64
	want   [][]int
	unique bool

	none concrete.Edge
}{
	// Positive weighted graphs.
	{
		name: "empty directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
		weight: math.Inf(1),

		none: concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
	},
	{
		name: "empty undirected",
		g:    func() graph.Mutable { return concrete.NewGraph(0, math.Inf(1)) },

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
		weight: math.Inf(1),

		none: concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
	},
	{
		name: "one edge directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
		},

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
		weight: 1,
		want: [][]int{
			{0, 1},
		},
		unique: true,

		none: concrete.Edge{F: concrete.Node(2), T: concrete.Node(3)},
	},
	{
		name: "one edge self directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
		},

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(0)},
		weight: 0,
		want: [][]int{
			{0},
		},
		unique: true,

		none: concrete.Edge{F: concrete.Node(2), T: concrete.Node(3)},
	},
	{
		name: "one edge undirected",
		g:    func() graph.Mutable { return concrete.NewGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
		},

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
		weight: 1,
		want: [][]int{
			{0, 1},
		},
		unique: true,

		none: concrete.Edge{F: concrete.Node(2), T: concrete.Node(3)},
	},
	{
		name: "two paths directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			{F: concrete.Node(0), T: concrete.Node(2), W: 2},
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
			{F: concrete.Node(1), T: concrete.Node(2), W: 1},
		},

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(2)},
		weight: 2,
		want: [][]int{
			{0, 1, 2},
			{0, 2},
		},
		unique: false,

		none: concrete.Edge{F: concrete.Node(2), T: concrete.Node(1)},
	},
	{
		name: "two paths undirected",
		g:    func() graph.Mutable { return concrete.NewGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			{F: concrete.Node(0), T: concrete.Node(2), W: 2},
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
			{F: concrete.Node(1), T: concrete.Node(2), W: 1},
		},

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(2)},
		weight: 2,
		want: [][]int{
			{0, 1, 2},
			{0, 2},
		},
		unique: false,

		none: concrete.Edge{F: concrete.Node(2), T: concrete.Node(4)},
	},
	{
		name: "confounding paths directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			// Add a path from 0->5 of weight 4
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
			{F: concrete.Node(1), T: concrete.Node(2), W: 1},
			{F: concrete.Node(2), T: concrete.Node(3), W: 1},
			{F: concrete.Node(3), T: concrete.Node(5), W: 1},

			// Add direct edge to goal of weight 4
			{F: concrete.Node(0), T: concrete.Node(5), W: 4},

			// Add edge to a node that's still optimal
			{F: concrete.Node(0), T: concrete.Node(2), W: 2},

			// Add edge to 3 that's overpriced
			{F: concrete.Node(0), T: concrete.Node(3), W: 4},

			// Add very cheap edge to 4 which is a dead end
			{F: concrete.Node(0), T: concrete.Node(4), W: 0.25},
		},

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(5)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 5},
			{0, 2, 3, 5},
			{0, 5},
		},
		unique: false,

		none: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		name: "confounding paths undirected",
		g:    func() graph.Mutable { return concrete.NewGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			// Add a path from 0->5 of weight 4
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
			{F: concrete.Node(1), T: concrete.Node(2), W: 1},
			{F: concrete.Node(2), T: concrete.Node(3), W: 1},
			{F: concrete.Node(3), T: concrete.Node(5), W: 1},

			// Add direct edge to goal of weight 4
			{F: concrete.Node(0), T: concrete.Node(5), W: 4},

			// Add edge to a node that's still optimal
			{F: concrete.Node(0), T: concrete.Node(2), W: 2},

			// Add edge to 3 that's overpriced
			{F: concrete.Node(0), T: concrete.Node(3), W: 4},

			// Add very cheap edge to 4 which is a dead end
			{F: concrete.Node(0), T: concrete.Node(4), W: 0.25},
		},

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(5)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 5},
			{0, 2, 3, 5},
			{0, 5},
		},
		unique: false,

		none: concrete.Edge{F: concrete.Node(5), T: concrete.Node(6)},
	},
	{
		name: "confounding paths directed 2-step",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			// Add a path from 0->5 of weight 4
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
			{F: concrete.Node(1), T: concrete.Node(2), W: 1},
			{F: concrete.Node(2), T: concrete.Node(3), W: 1},
			{F: concrete.Node(3), T: concrete.Node(5), W: 1},

			// Add two step path to goal of weight 4
			{F: concrete.Node(0), T: concrete.Node(6), W: 2},
			{F: concrete.Node(6), T: concrete.Node(5), W: 2},

			// Add edge to a node that's still optimal
			{F: concrete.Node(0), T: concrete.Node(2), W: 2},

			// Add edge to 3 that's overpriced
			{F: concrete.Node(0), T: concrete.Node(3), W: 4},

			// Add very cheap edge to 4 which is a dead end
			{F: concrete.Node(0), T: concrete.Node(4), W: 0.25},
		},

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(5)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 5},
			{0, 2, 3, 5},
			{0, 6, 5},
		},
		unique: false,

		none: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		name: "confounding paths undirected 2-step",
		g:    func() graph.Mutable { return concrete.NewGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			// Add a path from 0->5 of weight 4
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
			{F: concrete.Node(1), T: concrete.Node(2), W: 1},
			{F: concrete.Node(2), T: concrete.Node(3), W: 1},
			{F: concrete.Node(3), T: concrete.Node(5), W: 1},

			// Add two step path to goal of weight 4
			{F: concrete.Node(0), T: concrete.Node(6), W: 2},
			{F: concrete.Node(6), T: concrete.Node(5), W: 2},

			// Add edge to a node that's still optimal
			{F: concrete.Node(0), T: concrete.Node(2), W: 2},

			// Add edge to 3 that's overpriced
			{F: concrete.Node(0), T: concrete.Node(3), W: 4},

			// Add very cheap edge to 4 which is a dead end
			{F: concrete.Node(0), T: concrete.Node(4), W: 0.25},
		},

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(5)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 5},
			{0, 2, 3, 5},
			{0, 6, 5},
		},
		unique: false,

		none: concrete.Edge{F: concrete.Node(5), T: concrete.Node(7)},
	},
	{
		name: "zero-weight cycle directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			// Add a path from 0->4 of weight 4
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
			{F: concrete.Node(1), T: concrete.Node(2), W: 1},
			{F: concrete.Node(2), T: concrete.Node(3), W: 1},
			{F: concrete.Node(3), T: concrete.Node(4), W: 1},

			// Add a zero-weight cycle.
			{F: concrete.Node(1), T: concrete.Node(5), W: 0},
			{F: concrete.Node(5), T: concrete.Node(1), W: 0},
		},

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
		},
		unique: false,

		none: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		name: "zero-weight cycle^2 directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			// Add a path from 0->4 of weight 4
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
			{F: concrete.Node(1), T: concrete.Node(2), W: 1},
			{F: concrete.Node(2), T: concrete.Node(3), W: 1},
			{F: concrete.Node(3), T: concrete.Node(4), W: 1},

			// Add a zero-weight cycle.
			{F: concrete.Node(1), T: concrete.Node(5), W: 0},
			{F: concrete.Node(5), T: concrete.Node(1), W: 0},
			// With its own zero-weight cycle.
			{F: concrete.Node(5), T: concrete.Node(6), W: 0},
			{F: concrete.Node(6), T: concrete.Node(5), W: 0},
		},

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
		},
		unique: false,

		none: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		name: "zero-weight cycle^2 confounding directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			// Add a path from 0->4 of weight 4
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
			{F: concrete.Node(1), T: concrete.Node(2), W: 1},
			{F: concrete.Node(2), T: concrete.Node(3), W: 1},
			{F: concrete.Node(3), T: concrete.Node(4), W: 1},

			// Add a zero-weight cycle.
			{F: concrete.Node(1), T: concrete.Node(5), W: 0},
			{F: concrete.Node(5), T: concrete.Node(1), W: 0},
			// With its own zero-weight cycle.
			{F: concrete.Node(5), T: concrete.Node(6), W: 0},
			{F: concrete.Node(6), T: concrete.Node(5), W: 0},
			// But leading to the target.
			{F: concrete.Node(5), T: concrete.Node(4), W: 3},
		},

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
			{0, 1, 5, 4},
		},
		unique: false,

		none: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		name: "zero-weight cycle^3 directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			// Add a path from 0->4 of weight 4
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
			{F: concrete.Node(1), T: concrete.Node(2), W: 1},
			{F: concrete.Node(2), T: concrete.Node(3), W: 1},
			{F: concrete.Node(3), T: concrete.Node(4), W: 1},

			// Add a zero-weight cycle.
			{F: concrete.Node(1), T: concrete.Node(5), W: 0},
			{F: concrete.Node(5), T: concrete.Node(1), W: 0},
			// With its own zero-weight cycle.
			{F: concrete.Node(5), T: concrete.Node(6), W: 0},
			{F: concrete.Node(6), T: concrete.Node(5), W: 0},
			// With its own zero-weight cycle.
			{F: concrete.Node(6), T: concrete.Node(7), W: 0},
			{F: concrete.Node(7), T: concrete.Node(6), W: 0},
		},

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
		},
		unique: false,

		none: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		name: "zero-weight 3·cycle^2 confounding directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			// Add a path from 0->4 of weight 4
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
			{F: concrete.Node(1), T: concrete.Node(2), W: 1},
			{F: concrete.Node(2), T: concrete.Node(3), W: 1},
			{F: concrete.Node(3), T: concrete.Node(4), W: 1},

			// Add a zero-weight cycle.
			{F: concrete.Node(1), T: concrete.Node(5), W: 0},
			{F: concrete.Node(5), T: concrete.Node(1), W: 0},
			// With 3 of its own zero-weight cycles.
			{F: concrete.Node(5), T: concrete.Node(6), W: 0},
			{F: concrete.Node(6), T: concrete.Node(5), W: 0},
			{F: concrete.Node(5), T: concrete.Node(7), W: 0},
			{F: concrete.Node(7), T: concrete.Node(5), W: 0},
			// Each leading to the target.
			{F: concrete.Node(5), T: concrete.Node(4), W: 3},
			{F: concrete.Node(6), T: concrete.Node(4), W: 3},
			{F: concrete.Node(7), T: concrete.Node(4), W: 3},
		},

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
			{0, 1, 5, 4},
			{0, 1, 5, 6, 4},
			{0, 1, 5, 7, 4},
		},
		unique: false,

		none: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		name: "zero-weight reversed 3·cycle^2 confounding directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			// Add a path from 0->4 of weight 4
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
			{F: concrete.Node(1), T: concrete.Node(2), W: 1},
			{F: concrete.Node(2), T: concrete.Node(3), W: 1},
			{F: concrete.Node(3), T: concrete.Node(4), W: 1},

			// Add a zero-weight cycle.
			{F: concrete.Node(3), T: concrete.Node(5), W: 0},
			{F: concrete.Node(5), T: concrete.Node(3), W: 0},
			// With 3 of its own zero-weight cycles.
			{F: concrete.Node(5), T: concrete.Node(6), W: 0},
			{F: concrete.Node(6), T: concrete.Node(5), W: 0},
			{F: concrete.Node(5), T: concrete.Node(7), W: 0},
			{F: concrete.Node(7), T: concrete.Node(5), W: 0},
			// Each leading from the source.
			{F: concrete.Node(0), T: concrete.Node(5), W: 3},
			{F: concrete.Node(0), T: concrete.Node(6), W: 3},
			{F: concrete.Node(0), T: concrete.Node(7), W: 3},
		},

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
			{0, 5, 3, 4},
			{0, 6, 5, 3, 4},
			{0, 7, 5, 3, 4},
		},
		unique: false,

		none: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		name: "zero-weight |V|·cycle^(n/|V|) directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: func() []concrete.Edge {
			e := []concrete.Edge{
				// Add a path from 0->4 of weight 4
				{F: concrete.Node(0), T: concrete.Node(1), W: 1},
				{F: concrete.Node(1), T: concrete.Node(2), W: 1},
				{F: concrete.Node(2), T: concrete.Node(3), W: 1},
				{F: concrete.Node(3), T: concrete.Node(4), W: 1},
			}
			next := len(e) + 1

			// Add n zero-weight cycles.
			const n = 100
			for i := 0; i < n; i++ {
				e = append(e,
					concrete.Edge{F: concrete.Node(next + i), T: concrete.Node(i), W: 0},
					concrete.Edge{F: concrete.Node(i), T: concrete.Node(next + i), W: 0},
				)
			}
			return e
		}(),

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
		},
		unique: false,

		none: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		name: "zero-weight n·cycle directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: func() []concrete.Edge {
			e := []concrete.Edge{
				// Add a path from 0->4 of weight 4
				{F: concrete.Node(0), T: concrete.Node(1), W: 1},
				{F: concrete.Node(1), T: concrete.Node(2), W: 1},
				{F: concrete.Node(2), T: concrete.Node(3), W: 1},
				{F: concrete.Node(3), T: concrete.Node(4), W: 1},
			}
			next := len(e) + 1

			// Add n zero-weight cycles.
			const n = 100
			for i := 0; i < n; i++ {
				e = append(e,
					concrete.Edge{F: concrete.Node(next + i), T: concrete.Node(1), W: 0},
					concrete.Edge{F: concrete.Node(1), T: concrete.Node(next + i), W: 0},
				)
			}
			return e
		}(),

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
		},
		unique: false,

		none: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		name: "zero-weight bi-directional tree with single exit directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: func() []concrete.Edge {
			e := []concrete.Edge{
				// Add a path from 0->4 of weight 4
				{F: concrete.Node(0), T: concrete.Node(1), W: 1},
				{F: concrete.Node(1), T: concrete.Node(2), W: 1},
				{F: concrete.Node(2), T: concrete.Node(3), W: 1},
				{F: concrete.Node(3), T: concrete.Node(4), W: 1},
			}

			// Make a bi-directional tree rooted at node 2 with
			// a single exit to node 4 and co-equal cost from
			// 2 to 4.
			const (
				depth     = 4
				branching = 4
			)

			next := len(e) + 1
			src := 2
			var i, last int
			for l := 0; l < depth; l++ {
				for i = 0; i < branching; i++ {
					last = next + i
					e = append(e, concrete.Edge{F: concrete.Node(src), T: concrete.Node(last), W: 0})
					e = append(e, concrete.Edge{F: concrete.Node(last), T: concrete.Node(src), W: 0})
				}
				src = next + 1
				next += branching
			}
			e = append(e, concrete.Edge{F: concrete.Node(last), T: concrete.Node(4), W: 2})
			return e
		}(),

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
			{0, 1, 2, 6, 10, 14, 20, 4},
		},
		unique: false,

		none: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},

	// Negative weighted graphs.
	{
		name: "one edge directed negative",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			{F: concrete.Node(0), T: concrete.Node(1), W: -1},
		},
		negative: true,

		query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
		weight: -1,
		want: [][]int{
			{0, 1},
		},
		unique: true,

		none: concrete.Edge{F: concrete.Node(2), T: concrete.Node(3)},
	},
	{
		name: "one edge undirected negative",
		g:    func() graph.Mutable { return concrete.NewGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			{F: concrete.Node(0), T: concrete.Node(1), W: -1},
		},
		negative:         true,
		hasNegativeCycle: true,

		query: concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
	},
	{
		name: "wp graph negative", // http://en.wikipedia.org/w/index.php?title=Johnson%27s_algorithm&oldid=564595231
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			{F: concrete.Node('w'), T: concrete.Node('z'), W: 2},
			{F: concrete.Node('x'), T: concrete.Node('w'), W: 6},
			{F: concrete.Node('x'), T: concrete.Node('y'), W: 3},
			{F: concrete.Node('y'), T: concrete.Node('w'), W: 4},
			{F: concrete.Node('y'), T: concrete.Node('z'), W: 5},
			{F: concrete.Node('z'), T: concrete.Node('x'), W: -7},
			{F: concrete.Node('z'), T: concrete.Node('y'), W: -3},
		},
		negative: true,

		query:  concrete.Edge{F: concrete.Node('z'), T: concrete.Node('y')},
		weight: -4,
		want: [][]int{
			{'z', 'x', 'y'},
		},
		unique: true,

		none: concrete.Edge{F: concrete.Node(2), T: concrete.Node(3)},
	},
	{
		name: "roughgarden negative",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			{F: concrete.Node('a'), T: concrete.Node('b'), W: -2},
			{F: concrete.Node('b'), T: concrete.Node('c'), W: -1},
			{F: concrete.Node('c'), T: concrete.Node('a'), W: 4},
			{F: concrete.Node('c'), T: concrete.Node('x'), W: 2},
			{F: concrete.Node('c'), T: concrete.Node('y'), W: -3},
			{F: concrete.Node('z'), T: concrete.Node('x'), W: 1},
			{F: concrete.Node('z'), T: concrete.Node('y'), W: -4},
		},
		negative: true,

		query:  concrete.Edge{F: concrete.Node('a'), T: concrete.Node('y')},
		weight: -6,
		want: [][]int{
			{'a', 'b', 'c', 'y'},
		},
		unique: true,

		none: concrete.Edge{F: concrete.Node(2), T: concrete.Node(3)},
	},
}
