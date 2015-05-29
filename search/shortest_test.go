// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search_test

import (
	"fmt"
	"math"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

func init() {
	for _, test := range positiveWeightTests {
		if len(test.want) != 1 && test.unique {
			panic(fmt.Sprintf("%q: bad shortest path test: non-unique paths marked unique", test.name))
		}
	}
}

// positiveWeightTests are positively weighted edge graphs
// used to test DijkstraAllPaths, DijkstraFrom and FloydWarshall.
var positiveWeightTests = []struct {
	name  string
	g     func() graph.Mutable
	edges []concrete.WeightedEdge

	query  concrete.Edge
	weight float64
	want   [][]int
	unique bool

	none concrete.Edge
}{
	{
		name: "empty directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph() },

		query:  concrete.Edge{concrete.Node(0), concrete.Node(1)},
		weight: math.Inf(1),

		none: concrete.Edge{concrete.Node(0), concrete.Node(1)},
	},
	{
		name: "empty undirected",
		g:    func() graph.Mutable { return concrete.NewGraph() },

		query:  concrete.Edge{concrete.Node(0), concrete.Node(1)},
		weight: math.Inf(1),

		none: concrete.Edge{concrete.Node(0), concrete.Node(1)},
	},
	{
		name: "one edge directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph() },
		edges: []concrete.WeightedEdge{
			{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
		},

		query:  concrete.Edge{concrete.Node(0), concrete.Node(1)},
		weight: 1,
		want: [][]int{
			{0, 1},
		},
		unique: true,

		none: concrete.Edge{concrete.Node(2), concrete.Node(3)},
	},
	{
		name: "one edge undirected",
		g:    func() graph.Mutable { return concrete.NewGraph() },
		edges: []concrete.WeightedEdge{
			{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
		},

		query:  concrete.Edge{concrete.Node(0), concrete.Node(1)},
		weight: 1,
		want: [][]int{
			{0, 1},
		},
		unique: true,

		none: concrete.Edge{concrete.Node(2), concrete.Node(3)},
	},
	{
		name: "two paths directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph() },
		edges: []concrete.WeightedEdge{
			{concrete.Edge{concrete.Node(0), concrete.Node(2)}, 2},
			{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
			{concrete.Edge{concrete.Node(1), concrete.Node(2)}, 1},
		},

		query:  concrete.Edge{concrete.Node(0), concrete.Node(2)},
		weight: 2,
		want: [][]int{
			{0, 1, 2},
			{0, 2},
		},
		unique: false,

		none: concrete.Edge{concrete.Node(2), concrete.Node(1)},
	},
	{
		name: "two paths undirected",
		g:    func() graph.Mutable { return concrete.NewGraph() },
		edges: []concrete.WeightedEdge{
			{concrete.Edge{concrete.Node(0), concrete.Node(2)}, 2},
			{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
			{concrete.Edge{concrete.Node(1), concrete.Node(2)}, 1},
		},

		query:  concrete.Edge{concrete.Node(0), concrete.Node(2)},
		weight: 2,
		want: [][]int{
			{0, 1, 2},
			{0, 2},
		},
		unique: false,

		none: concrete.Edge{concrete.Node(2), concrete.Node(4)},
	},
	{
		name: "confounding paths directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph() },
		edges: []concrete.WeightedEdge{
			// Add a path from 0->5 of weight 4
			{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
			{concrete.Edge{concrete.Node(1), concrete.Node(2)}, 1},
			{concrete.Edge{concrete.Node(2), concrete.Node(3)}, 1},
			{concrete.Edge{concrete.Node(3), concrete.Node(5)}, 1},

			// Add direct edge to goal of weight 4
			{concrete.Edge{concrete.Node(0), concrete.Node(5)}, 4},

			// Add edge to a node that's still optimal
			{concrete.Edge{concrete.Node(0), concrete.Node(2)}, 2},

			// Add edge to 3 that's overpriced
			{concrete.Edge{concrete.Node(0), concrete.Node(3)}, 4},

			// Add very cheap edge to 4 which is a dead end
			{concrete.Edge{concrete.Node(0), concrete.Node(4)}, 0.25},
		},

		query:  concrete.Edge{concrete.Node(0), concrete.Node(5)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 5},
			{0, 2, 3, 5},
			{0, 5},
		},
		unique: false,

		none: concrete.Edge{concrete.Node(4), concrete.Node(5)},
	},
	{
		name: "confounding paths undirected",
		g:    func() graph.Mutable { return concrete.NewGraph() },
		edges: []concrete.WeightedEdge{
			// Add a path from 0->5 of weight 4
			{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
			{concrete.Edge{concrete.Node(1), concrete.Node(2)}, 1},
			{concrete.Edge{concrete.Node(2), concrete.Node(3)}, 1},
			{concrete.Edge{concrete.Node(3), concrete.Node(5)}, 1},

			// Add direct edge to goal of weight 4
			{concrete.Edge{concrete.Node(0), concrete.Node(5)}, 4},

			// Add edge to a node that's still optimal
			{concrete.Edge{concrete.Node(0), concrete.Node(2)}, 2},

			// Add edge to 3 that's overpriced
			{concrete.Edge{concrete.Node(0), concrete.Node(3)}, 4},

			// Add very cheap edge to 4 which is a dead end
			{concrete.Edge{concrete.Node(0), concrete.Node(4)}, 0.25},
		},

		query:  concrete.Edge{concrete.Node(0), concrete.Node(5)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 5},
			{0, 2, 3, 5},
			{0, 5},
		},
		unique: false,

		none: concrete.Edge{concrete.Node(5), concrete.Node(6)},
	},
	{
		name: "confounding paths directed 2-step",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph() },
		edges: []concrete.WeightedEdge{
			// Add a path from 0->5 of weight 4
			{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
			{concrete.Edge{concrete.Node(1), concrete.Node(2)}, 1},
			{concrete.Edge{concrete.Node(2), concrete.Node(3)}, 1},
			{concrete.Edge{concrete.Node(3), concrete.Node(5)}, 1},

			// Add two step path to goal of weight 4
			{concrete.Edge{concrete.Node(0), concrete.Node(6)}, 2},
			{concrete.Edge{concrete.Node(6), concrete.Node(5)}, 2},

			// Add edge to a node that's still optimal
			{concrete.Edge{concrete.Node(0), concrete.Node(2)}, 2},

			// Add edge to 3 that's overpriced
			{concrete.Edge{concrete.Node(0), concrete.Node(3)}, 4},

			// Add very cheap edge to 4 which is a dead end
			{concrete.Edge{concrete.Node(0), concrete.Node(4)}, 0.25},
		},

		query:  concrete.Edge{concrete.Node(0), concrete.Node(5)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 5},
			{0, 2, 3, 5},
			{0, 6, 5},
		},
		unique: false,

		none: concrete.Edge{concrete.Node(4), concrete.Node(5)},
	},
	{
		name: "confounding paths undirected 2-step",
		g:    func() graph.Mutable { return concrete.NewGraph() },
		edges: []concrete.WeightedEdge{
			// Add a path from 0->5 of weight 4
			{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
			{concrete.Edge{concrete.Node(1), concrete.Node(2)}, 1},
			{concrete.Edge{concrete.Node(2), concrete.Node(3)}, 1},
			{concrete.Edge{concrete.Node(3), concrete.Node(5)}, 1},

			// Add two step path to goal of weight 4
			{concrete.Edge{concrete.Node(0), concrete.Node(6)}, 2},
			{concrete.Edge{concrete.Node(6), concrete.Node(5)}, 2},

			// Add edge to a node that's still optimal
			{concrete.Edge{concrete.Node(0), concrete.Node(2)}, 2},

			// Add edge to 3 that's overpriced
			{concrete.Edge{concrete.Node(0), concrete.Node(3)}, 4},

			// Add very cheap edge to 4 which is a dead end
			{concrete.Edge{concrete.Node(0), concrete.Node(4)}, 0.25},
		},

		query:  concrete.Edge{concrete.Node(0), concrete.Node(5)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 5},
			{0, 2, 3, 5},
			{0, 6, 5},
		},
		unique: false,

		none: concrete.Edge{concrete.Node(5), concrete.Node(7)},
	},
	{
		name: "zero-weight cycle directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph() },
		edges: []concrete.WeightedEdge{
			// Add a path from 0->4 of weight 4
			{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
			{concrete.Edge{concrete.Node(1), concrete.Node(2)}, 1},
			{concrete.Edge{concrete.Node(2), concrete.Node(3)}, 1},
			{concrete.Edge{concrete.Node(3), concrete.Node(4)}, 1},

			// Add a zero-weight cycle.
			{concrete.Edge{concrete.Node(1), concrete.Node(5)}, 0},
			{concrete.Edge{concrete.Node(5), concrete.Node(1)}, 0},
		},

		query:  concrete.Edge{concrete.Node(0), concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
		},
		unique: false,

		none: concrete.Edge{concrete.Node(4), concrete.Node(5)},
	},
	{
		name: "zero-weight cycle^2 directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph() },
		edges: []concrete.WeightedEdge{
			// Add a path from 0->4 of weight 4
			{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
			{concrete.Edge{concrete.Node(1), concrete.Node(2)}, 1},
			{concrete.Edge{concrete.Node(2), concrete.Node(3)}, 1},
			{concrete.Edge{concrete.Node(3), concrete.Node(4)}, 1},

			// Add a zero-weight cycle.
			{concrete.Edge{concrete.Node(1), concrete.Node(5)}, 0},
			{concrete.Edge{concrete.Node(5), concrete.Node(1)}, 0},
			// With its own zero-weight cycle.
			{concrete.Edge{concrete.Node(5), concrete.Node(5)}, 0},
		},

		query:  concrete.Edge{concrete.Node(0), concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
		},
		unique: false,

		none: concrete.Edge{concrete.Node(4), concrete.Node(5)},
	},
	{
		name: "zero-weight cycle^2 confounding directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph() },
		edges: []concrete.WeightedEdge{
			// Add a path from 0->4 of weight 4
			{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
			{concrete.Edge{concrete.Node(1), concrete.Node(2)}, 1},
			{concrete.Edge{concrete.Node(2), concrete.Node(3)}, 1},
			{concrete.Edge{concrete.Node(3), concrete.Node(4)}, 1},

			// Add a zero-weight cycle.
			{concrete.Edge{concrete.Node(1), concrete.Node(5)}, 0},
			{concrete.Edge{concrete.Node(5), concrete.Node(1)}, 0},
			// With its own zero-weight cycle.
			{concrete.Edge{concrete.Node(5), concrete.Node(5)}, 0},
			// But leading to the target.
			{concrete.Edge{concrete.Node(5), concrete.Node(4)}, 3},
		},

		query:  concrete.Edge{concrete.Node(0), concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
			{0, 1, 5, 4},
		},
		unique: false,

		none: concrete.Edge{concrete.Node(4), concrete.Node(5)},
	},
	{
		name: "zero-weight cycle^3 directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph() },
		edges: []concrete.WeightedEdge{
			// Add a path from 0->4 of weight 4
			{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
			{concrete.Edge{concrete.Node(1), concrete.Node(2)}, 1},
			{concrete.Edge{concrete.Node(2), concrete.Node(3)}, 1},
			{concrete.Edge{concrete.Node(3), concrete.Node(4)}, 1},

			// Add a zero-weight cycle.
			{concrete.Edge{concrete.Node(1), concrete.Node(5)}, 0},
			{concrete.Edge{concrete.Node(5), concrete.Node(1)}, 0},
			// With its own zero-weight cycle.
			{concrete.Edge{concrete.Node(5), concrete.Node(6)}, 0},
			{concrete.Edge{concrete.Node(6), concrete.Node(5)}, 0},
			// With its own zero-weight cycle.
			{concrete.Edge{concrete.Node(6), concrete.Node(6)}, 0},
		},

		query:  concrete.Edge{concrete.Node(0), concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
		},
		unique: false,

		none: concrete.Edge{concrete.Node(4), concrete.Node(5)},
	},
	{
		name: "zero-weight 3·cycle^2 confounding directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph() },
		edges: []concrete.WeightedEdge{
			// Add a path from 0->4 of weight 4
			{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
			{concrete.Edge{concrete.Node(1), concrete.Node(2)}, 1},
			{concrete.Edge{concrete.Node(2), concrete.Node(3)}, 1},
			{concrete.Edge{concrete.Node(3), concrete.Node(4)}, 1},

			// Add a zero-weight cycle.
			{concrete.Edge{concrete.Node(1), concrete.Node(5)}, 0},
			{concrete.Edge{concrete.Node(5), concrete.Node(1)}, 0},
			// With 3 of its own zero-weight cycles.
			{concrete.Edge{concrete.Node(5), concrete.Node(6)}, 0},
			{concrete.Edge{concrete.Node(6), concrete.Node(5)}, 0},
			{concrete.Edge{concrete.Node(5), concrete.Node(7)}, 0},
			{concrete.Edge{concrete.Node(7), concrete.Node(5)}, 0},
			// Each leading to the target.
			{concrete.Edge{concrete.Node(5), concrete.Node(4)}, 3},
			{concrete.Edge{concrete.Node(6), concrete.Node(4)}, 3},
			{concrete.Edge{concrete.Node(7), concrete.Node(4)}, 3},
		},

		query:  concrete.Edge{concrete.Node(0), concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
			{0, 1, 5, 4},
			{0, 1, 5, 6, 4},
			{0, 1, 5, 7, 4},
		},
		unique: false,

		none: concrete.Edge{concrete.Node(4), concrete.Node(5)},
	},
	{
		name: "zero-weight reversed 3·cycle^2 confounding directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph() },
		edges: []concrete.WeightedEdge{
			// Add a path from 0->4 of weight 4
			{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
			{concrete.Edge{concrete.Node(1), concrete.Node(2)}, 1},
			{concrete.Edge{concrete.Node(2), concrete.Node(3)}, 1},
			{concrete.Edge{concrete.Node(3), concrete.Node(4)}, 1},

			// Add a zero-weight cycle.
			{concrete.Edge{concrete.Node(3), concrete.Node(5)}, 0},
			{concrete.Edge{concrete.Node(5), concrete.Node(3)}, 0},
			// With 3 of its own zero-weight cycles.
			{concrete.Edge{concrete.Node(5), concrete.Node(6)}, 0},
			{concrete.Edge{concrete.Node(6), concrete.Node(5)}, 0},
			{concrete.Edge{concrete.Node(5), concrete.Node(7)}, 0},
			{concrete.Edge{concrete.Node(7), concrete.Node(5)}, 0},
			// Each leading from the source.
			{concrete.Edge{concrete.Node(0), concrete.Node(5)}, 3},
			{concrete.Edge{concrete.Node(0), concrete.Node(6)}, 3},
			{concrete.Edge{concrete.Node(0), concrete.Node(7)}, 3},
		},

		query:  concrete.Edge{concrete.Node(0), concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
			{0, 5, 3, 4},
			{0, 6, 5, 3, 4},
			{0, 7, 5, 3, 4},
		},
		unique: false,

		none: concrete.Edge{concrete.Node(4), concrete.Node(5)},
	},
	{
		name: "zero-weight |V|·cycle^(n/|V|) directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph() },
		edges: func() []concrete.WeightedEdge {
			e := []concrete.WeightedEdge{
				// Add a path from 0->4 of weight 4
				{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
				{concrete.Edge{concrete.Node(1), concrete.Node(2)}, 1},
				{concrete.Edge{concrete.Node(2), concrete.Node(3)}, 1},
				{concrete.Edge{concrete.Node(3), concrete.Node(4)}, 1},
			}
			next := len(e) + 1

			// Add n zero-weight cycles.
			const n = 100
			for i := 0; i < n; i++ {
				e = append(e,
					concrete.WeightedEdge{concrete.Edge{concrete.Node(next + i), concrete.Node(i)}, 0},
					concrete.WeightedEdge{concrete.Edge{concrete.Node(i), concrete.Node(next + i)}, 0},
				)
			}
			return e
		}(),

		query:  concrete.Edge{concrete.Node(0), concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
		},
		unique: false,

		none: concrete.Edge{concrete.Node(4), concrete.Node(5)},
	},
	{
		name: "zero-weight n·cycle directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph() },
		edges: func() []concrete.WeightedEdge {
			e := []concrete.WeightedEdge{
				// Add a path from 0->4 of weight 4
				{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
				{concrete.Edge{concrete.Node(1), concrete.Node(2)}, 1},
				{concrete.Edge{concrete.Node(2), concrete.Node(3)}, 1},
				{concrete.Edge{concrete.Node(3), concrete.Node(4)}, 1},
			}
			next := len(e) + 1

			// Add n zero-weight cycles.
			const n = 100
			for i := 0; i < n; i++ {
				e = append(e,
					concrete.WeightedEdge{concrete.Edge{concrete.Node(next + i), concrete.Node(1)}, 0},
					concrete.WeightedEdge{concrete.Edge{concrete.Node(1), concrete.Node(next + i)}, 0},
				)
			}
			return e
		}(),

		query:  concrete.Edge{concrete.Node(0), concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
		},
		unique: false,

		none: concrete.Edge{concrete.Node(4), concrete.Node(5)},
	},
	{
		name: "zero-weight bi-directional tree with single exit directed",
		g:    func() graph.Mutable { return concrete.NewDirectedGraph() },
		edges: func() []concrete.WeightedEdge {
			e := []concrete.WeightedEdge{
				// Add a path from 0->4 of weight 4
				{concrete.Edge{concrete.Node(0), concrete.Node(1)}, 1},
				{concrete.Edge{concrete.Node(1), concrete.Node(2)}, 1},
				{concrete.Edge{concrete.Node(2), concrete.Node(3)}, 1},
				{concrete.Edge{concrete.Node(3), concrete.Node(4)}, 1},
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
					e = append(e, concrete.WeightedEdge{concrete.Edge{concrete.Node(src), concrete.Node(last)}, 0})
					e = append(e, concrete.WeightedEdge{concrete.Edge{concrete.Node(last), concrete.Node(src)}, 0})
				}
				src = next + 1
				next += branching
			}
			e = append(e, concrete.WeightedEdge{concrete.Edge{concrete.Node(last), concrete.Node(4)}, 2})
			return e
		}(),

		query:  concrete.Edge{concrete.Node(0), concrete.Node(4)},
		weight: 4,
		want: [][]int{
			{0, 1, 2, 3, 4},
			{0, 1, 2, 6, 10, 14, 20, 4},
		},
		unique: false,

		none: concrete.Edge{concrete.Node(4), concrete.Node(5)},
	},
}
