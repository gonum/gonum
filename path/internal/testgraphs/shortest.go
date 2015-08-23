// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testgraphs

import (
	"fmt"
	"math"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

func init() {
	for _, test := range ShortestPathTests {
		if len(test.WantPaths) != 1 && test.HasUniquePath {
			panic(fmt.Sprintf("%q: bad shortest path test: non-unique paths marked unique", test.Name))
		}
	}
}

// ShortestPathTests are graphs used to test the static shortest path routines in path: BellmanFord,
// DijkstraAllPaths, DijkstraFrom, FloydWarshall and Johnson, and the static degenerate case for the
// dynamic shortest path routine in path/dynamic: DStarLite.
var ShortestPathTests = []struct {
	Name              string
	Graph             func() graph.Mutable
	Edges             []concrete.Edge
	HasNegativeWeight bool
	HasNegativeCycle  bool

	Query         concrete.Edge
	Weight        float64
	WantPaths     [][]int
	HasUniquePath bool

	NoPathFor concrete.Edge
}{
	// Positive weighted graphs.
	{
		Name:  "empty directed",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
		Weight: math.Inf(1),

		NoPathFor: concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
	},
	{
		Name:  "empty undirected",
		Graph: func() graph.Mutable { return concrete.NewGraph(0, math.Inf(1)) },

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
		Weight: math.Inf(1),

		NoPathFor: concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
	},
	{
		Name:  "one edge directed",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
		},

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
		Weight: 1,
		WantPaths: [][]int{
			{0, 1},
		},
		HasUniquePath: true,

		NoPathFor: concrete.Edge{F: concrete.Node(2), T: concrete.Node(3)},
	},
	{
		Name:  "one edge self directed",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
		},

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(0)},
		Weight: 0,
		WantPaths: [][]int{
			{0},
		},
		HasUniquePath: true,

		NoPathFor: concrete.Edge{F: concrete.Node(2), T: concrete.Node(3)},
	},
	{
		Name:  "one edge undirected",
		Graph: func() graph.Mutable { return concrete.NewGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
		},

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
		Weight: 1,
		WantPaths: [][]int{
			{0, 1},
		},
		HasUniquePath: true,

		NoPathFor: concrete.Edge{F: concrete.Node(2), T: concrete.Node(3)},
	},
	{
		Name:  "two paths directed",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
			{F: concrete.Node(0), T: concrete.Node(2), W: 2},
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
			{F: concrete.Node(1), T: concrete.Node(2), W: 1},
		},

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(2)},
		Weight: 2,
		WantPaths: [][]int{
			{0, 1, 2},
			{0, 2},
		},
		HasUniquePath: false,

		NoPathFor: concrete.Edge{F: concrete.Node(2), T: concrete.Node(1)},
	},
	{
		Name:  "two paths undirected",
		Graph: func() graph.Mutable { return concrete.NewGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
			{F: concrete.Node(0), T: concrete.Node(2), W: 2},
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
			{F: concrete.Node(1), T: concrete.Node(2), W: 1},
		},

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(2)},
		Weight: 2,
		WantPaths: [][]int{
			{0, 1, 2},
			{0, 2},
		},
		HasUniquePath: false,

		NoPathFor: concrete.Edge{F: concrete.Node(2), T: concrete.Node(4)},
	},
	{
		Name:  "confounding paths directed",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
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

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(5)},
		Weight: 4,
		WantPaths: [][]int{
			{0, 1, 2, 3, 5},
			{0, 2, 3, 5},
			{0, 5},
		},
		HasUniquePath: false,

		NoPathFor: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		Name:  "confounding paths undirected",
		Graph: func() graph.Mutable { return concrete.NewGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
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

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(5)},
		Weight: 4,
		WantPaths: [][]int{
			{0, 1, 2, 3, 5},
			{0, 2, 3, 5},
			{0, 5},
		},
		HasUniquePath: false,

		NoPathFor: concrete.Edge{F: concrete.Node(5), T: concrete.Node(6)},
	},
	{
		Name:  "confounding paths directed 2-step",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
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

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(5)},
		Weight: 4,
		WantPaths: [][]int{
			{0, 1, 2, 3, 5},
			{0, 2, 3, 5},
			{0, 6, 5},
		},
		HasUniquePath: false,

		NoPathFor: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		Name:  "confounding paths undirected 2-step",
		Graph: func() graph.Mutable { return concrete.NewGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
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

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(5)},
		Weight: 4,
		WantPaths: [][]int{
			{0, 1, 2, 3, 5},
			{0, 2, 3, 5},
			{0, 6, 5},
		},
		HasUniquePath: false,

		NoPathFor: concrete.Edge{F: concrete.Node(5), T: concrete.Node(7)},
	},
	{
		Name:  "zero-weight cycle directed",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
			// Add a path from 0->4 of weight 4
			{F: concrete.Node(0), T: concrete.Node(1), W: 1},
			{F: concrete.Node(1), T: concrete.Node(2), W: 1},
			{F: concrete.Node(2), T: concrete.Node(3), W: 1},
			{F: concrete.Node(3), T: concrete.Node(4), W: 1},

			// Add a zero-weight cycle.
			{F: concrete.Node(1), T: concrete.Node(5), W: 0},
			{F: concrete.Node(5), T: concrete.Node(1), W: 0},
		},

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		Weight: 4,
		WantPaths: [][]int{
			{0, 1, 2, 3, 4},
		},
		HasUniquePath: false,

		NoPathFor: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		Name:  "zero-weight cycle^2 directed",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
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

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		Weight: 4,
		WantPaths: [][]int{
			{0, 1, 2, 3, 4},
		},
		HasUniquePath: false,

		NoPathFor: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		Name:  "zero-weight cycle^2 confounding directed",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
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

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		Weight: 4,
		WantPaths: [][]int{
			{0, 1, 2, 3, 4},
			{0, 1, 5, 4},
		},
		HasUniquePath: false,

		NoPathFor: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		Name:  "zero-weight cycle^3 directed",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
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

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		Weight: 4,
		WantPaths: [][]int{
			{0, 1, 2, 3, 4},
		},
		HasUniquePath: false,

		NoPathFor: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		Name:  "zero-weight 3·cycle^2 confounding directed",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
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

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		Weight: 4,
		WantPaths: [][]int{
			{0, 1, 2, 3, 4},
			{0, 1, 5, 4},
			{0, 1, 5, 6, 4},
			{0, 1, 5, 7, 4},
		},
		HasUniquePath: false,

		NoPathFor: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		Name:  "zero-weight reversed 3·cycle^2 confounding directed",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
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

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		Weight: 4,
		WantPaths: [][]int{
			{0, 1, 2, 3, 4},
			{0, 5, 3, 4},
			{0, 6, 5, 3, 4},
			{0, 7, 5, 3, 4},
		},
		HasUniquePath: false,

		NoPathFor: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		Name:  "zero-weight |V|·cycle^(n/|V|) directed",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: func() []concrete.Edge {
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

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		Weight: 4,
		WantPaths: [][]int{
			{0, 1, 2, 3, 4},
		},
		HasUniquePath: false,

		NoPathFor: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		Name:  "zero-weight n·cycle directed",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: func() []concrete.Edge {
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

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		Weight: 4,
		WantPaths: [][]int{
			{0, 1, 2, 3, 4},
		},
		HasUniquePath: false,

		NoPathFor: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},
	{
		Name:  "zero-weight bi-directional tree with single exit directed",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: func() []concrete.Edge {
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

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(4)},
		Weight: 4,
		WantPaths: [][]int{
			{0, 1, 2, 3, 4},
			{0, 1, 2, 6, 10, 14, 20, 4},
		},
		HasUniquePath: false,

		NoPathFor: concrete.Edge{F: concrete.Node(4), T: concrete.Node(5)},
	},

	// Negative weighted graphs.
	{
		Name:  "one edge directed negative",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
			{F: concrete.Node(0), T: concrete.Node(1), W: -1},
		},
		HasNegativeWeight: true,

		Query:  concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
		Weight: -1,
		WantPaths: [][]int{
			{0, 1},
		},
		HasUniquePath: true,

		NoPathFor: concrete.Edge{F: concrete.Node(2), T: concrete.Node(3)},
	},
	{
		Name:  "one edge undirected negative",
		Graph: func() graph.Mutable { return concrete.NewGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
			{F: concrete.Node(0), T: concrete.Node(1), W: -1},
		},
		HasNegativeWeight: true,
		HasNegativeCycle:  true,

		Query: concrete.Edge{F: concrete.Node(0), T: concrete.Node(1)},
	},
	{
		Name:  "wp graph negative", // http://en.wikipedia.org/w/index.php?title=Johnson%27s_algorithm&oldid=564595231
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
			{F: concrete.Node('w'), T: concrete.Node('z'), W: 2},
			{F: concrete.Node('x'), T: concrete.Node('w'), W: 6},
			{F: concrete.Node('x'), T: concrete.Node('y'), W: 3},
			{F: concrete.Node('y'), T: concrete.Node('w'), W: 4},
			{F: concrete.Node('y'), T: concrete.Node('z'), W: 5},
			{F: concrete.Node('z'), T: concrete.Node('x'), W: -7},
			{F: concrete.Node('z'), T: concrete.Node('y'), W: -3},
		},
		HasNegativeWeight: true,

		Query:  concrete.Edge{F: concrete.Node('z'), T: concrete.Node('y')},
		Weight: -4,
		WantPaths: [][]int{
			{'z', 'x', 'y'},
		},
		HasUniquePath: true,

		NoPathFor: concrete.Edge{F: concrete.Node(2), T: concrete.Node(3)},
	},
	{
		Name:  "roughgarden negative",
		Graph: func() graph.Mutable { return concrete.NewDirectedGraph(0, math.Inf(1)) },
		Edges: []concrete.Edge{
			{F: concrete.Node('a'), T: concrete.Node('b'), W: -2},
			{F: concrete.Node('b'), T: concrete.Node('c'), W: -1},
			{F: concrete.Node('c'), T: concrete.Node('a'), W: 4},
			{F: concrete.Node('c'), T: concrete.Node('x'), W: 2},
			{F: concrete.Node('c'), T: concrete.Node('y'), W: -3},
			{F: concrete.Node('z'), T: concrete.Node('x'), W: 1},
			{F: concrete.Node('z'), T: concrete.Node('y'), W: -4},
		},
		HasNegativeWeight: true,

		Query:  concrete.Edge{F: concrete.Node('a'), T: concrete.Node('y')},
		Weight: -6,
		WantPaths: [][]int{
			{'a', 'b', 'c', 'y'},
		},
		HasUniquePath: true,

		NoPathFor: concrete.Edge{F: concrete.Node(2), T: concrete.Node(3)},
	},
}
