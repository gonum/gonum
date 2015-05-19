// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search_test

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/internal"
	"github.com/gonum/graph/search"
)

var floydWarshallTests = []struct {
	name  string
	g     func() graph.Mutable
	edges []concrete.WeightedEdge

	query  concrete.Edge
	weight float64
	want   [][]int

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

		none: concrete.Edge{concrete.Node(5), concrete.Node(6)},
	},
}

func TestFloydWarshall(t *testing.T) {
	for _, test := range floydWarshallTests {
		g := test.g()
		for _, e := range test.edges {
			switch g := g.(type) {
			case graph.MutableDirectedGraph:
				g.AddDirectedEdge(e, e.Cost)
			case graph.MutableGraph:
				g.AddUndirectedEdge(e, e.Cost)
			default:
				panic("floyd warshall: bad graph type")
			}
		}

		pt, ok := search.FloydWarshall(g.(graph.Graph), nil)
		if !ok {
			t.Fatalf("%q: unexpected negative cycle", test.name)
		}

		// Check all random paths returned are OK.
		for i := 0; i < 10; i++ {
			p, weight, unique := pt.Between(test.query.Head(), test.query.Tail())
			if weight != test.weight {
				t.Errorf("%q: unexpected weight from Between: got:%f want:%f",
					test.name, weight, test.weight)
			}
			if weight := pt.Weight(test.query.Head(), test.query.Tail()); weight != test.weight {
				t.Errorf("%q: unexpected weight from Weight: got:%f want:%f",
					test.name, weight, test.weight)
			}
			if unique != (len(test.want) == 1) {
				t.Errorf("%q: unexpected number of paths: got: unique=%t want: unique=%d",
					test.name, unique, len(test.want) == 1)
			}

			var got []int
			for _, n := range p {
				got = append(got, n.ID())
			}
			for _, sp := range test.want {
				if reflect.DeepEqual(got, sp) {
					ok = true
					break
				}
			}
			if !ok {
				t.Errorf("%q: unexpected shortest path:\ngot: %v\nwant from:%v",
					test.name, p, test.want)
			}

			np, weight, unique := pt.Between(test.none.Head(), test.none.Tail())
			if np != nil || !math.IsInf(weight, 1) || unique != false {
				t.Errorf("%q: unexpected path:\ngot: path=%v weight=%f unique=%t\nwant:path=<nil> weight=+Inf unique=false",
					test.name, np, weight, unique)
			}
		}

		paths, weight := pt.AllBetween(test.query.Head(), test.query.Tail())
		if weight != test.weight {
			t.Errorf("%q: unexpected weight from Between: got:%f want:%f",
				test.name, weight, test.weight)
		}

		var got [][]int
		if len(paths) != 0 {
			got = make([][]int, len(paths))
		}
		for i, p := range paths {
			for _, v := range p {
				got[i] = append(got[i], v.ID())
			}
		}
		sort.Sort(internal.BySliceValues(got))
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("testing %q: unexpected shortest paths:\ngot: %v\nwant:%v",
				test.name, got, test.want)
		}

		np, weight := pt.AllBetween(test.none.Head(), test.none.Tail())
		if np != nil || !math.IsInf(weight, 1) {
			t.Errorf("%q: unexpected path:\ngot: paths=%v weight=%f\nwant:path=<nil> weight=+Inf",
				test.name, np, weight)
		}
	}
}
