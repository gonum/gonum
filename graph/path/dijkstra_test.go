// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"math"
	"reflect"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/ordered"
	"gonum.org/v1/gonum/graph/path/internal/testgraphs"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

func TestDijkstraFrom(t *testing.T) {
	t.Parallel()
	for _, test := range testgraphs.ShortestPathTests {
		g := test.Graph()
		for _, e := range test.Edges {
			g.SetWeightedEdge(e)
		}

		for _, tg := range []struct {
			typ string
			g   traverse.Graph
		}{
			{"complete", g.(graph.Graph)},
			{"incremental", incremental{g.(graph.Weighted)}},
		} {
			var (
				pt Shortest

				panicked bool
			)
			func() {
				defer func() {
					panicked = recover() != nil
				}()
				pt = DijkstraFrom(test.Query.From(), tg.g)
			}()
			if panicked || test.HasNegativeWeight {
				if !test.HasNegativeWeight {
					t.Errorf("%q %s: unexpected panic", test.Name, tg.typ)
				}
				if !panicked {
					t.Errorf("%q %s: expected panic for negative edge weight", test.Name, tg.typ)
				}
				continue
			}

			if pt.From().ID() != test.Query.From().ID() {
				t.Fatalf("%q %s: unexpected from node ID: got:%d want:%d", test.Name, tg.typ, pt.From().ID(), test.Query.From().ID())
			}

			p, weight := pt.To(test.Query.To().ID())
			if weight != test.Weight {
				t.Errorf("%q %s: unexpected weight from To: got:%f want:%f",
					test.Name, tg.typ, weight, test.Weight)
			}
			if weight := pt.WeightTo(test.Query.To().ID()); weight != test.Weight {
				t.Errorf("%q %s: unexpected weight from Weight: got:%f want:%f",
					test.Name, tg.typ, weight, test.Weight)
			}

			var got []int64
			for _, n := range p {
				got = append(got, n.ID())
			}
			ok := len(got) == 0 && len(test.WantPaths) == 0
			for _, sp := range test.WantPaths {
				if reflect.DeepEqual(got, sp) {
					ok = true
					break
				}
			}
			if !ok {
				t.Errorf("%q %s: unexpected shortest path:\ngot: %v\nwant from:%v",
					test.Name, tg.typ, p, test.WantPaths)
			}

			np, weight := pt.To(test.NoPathFor.To().ID())
			if pt.From().ID() == test.NoPathFor.From().ID() && (np != nil || !math.IsInf(weight, 1)) {
				t.Errorf("%q %s: unexpected path:\ngot: path=%v weight=%f\nwant:path=<nil> weight=+Inf",
					test.Name, tg.typ, np, weight)
			}
		}
	}
}

func TestDijkstraAllFrom(t *testing.T) {
	t.Parallel()
	for _, test := range testgraphs.ShortestPathTests {
		g := test.Graph()
		for _, e := range test.Edges {
			g.SetWeightedEdge(e)
		}

		for _, tg := range []struct {
			typ string
			g   traverse.Graph
		}{
			{"complete", g.(graph.Graph)},
			{"incremental", incremental{g.(graph.Weighted)}},
		} {
			var (
				pt ShortestAlts

				panicked bool
			)
			func() {
				defer func() {
					panicked = recover() != nil
				}()
				pt = DijkstraAllFrom(test.Query.From(), tg.g)
			}()
			if panicked || test.HasNegativeWeight {
				if !test.HasNegativeWeight {
					t.Errorf("%q %s: unexpected panic", test.Name, tg.typ)
				}
				if !panicked {
					t.Errorf("%q %s: expected panic for negative edge weight", test.Name, tg.typ)
				}
				continue
			}

			if pt.From().ID() != test.Query.From().ID() {
				t.Fatalf("%q %s: unexpected from node ID: got:%d want:%d", test.Name, tg.typ, pt.From().ID(), test.Query.From().ID())
			}

			// Test single path results.
			p, weight, unique := pt.To(test.Query.To().ID())
			if weight != test.Weight {
				t.Errorf("%q %s: unexpected weight from To: got:%f want:%f",
					test.Name, tg.typ, weight, test.Weight)
			}
			if weight := pt.WeightTo(test.Query.To().ID()); weight != test.Weight {
				t.Errorf("%q %s: unexpected weight from Weight: got:%f want:%f",
					test.Name, tg.typ, weight, test.Weight)
			}

			var gotPath []int64
			for _, n := range p {
				gotPath = append(gotPath, n.ID())
			}
			ok := len(gotPath) == 0 && len(test.WantPaths) == 0
			for _, sp := range test.WantPaths {
				if reflect.DeepEqual(gotPath, sp) {
					ok = true
					break
				}
			}
			if !ok {
				t.Errorf("%q: unexpected shortest path:\ngot: %v\nwant from:%v",
					test.Name, p, test.WantPaths)
			}
			if unique != test.HasUniquePath {
				t.Errorf("%q: unexpected uniqueness from To: got:%t want:%t (%d paths)",
					test.Name, unique, test.HasUniquePath, len(test.WantPaths))
			}

			// Test multiple path results.
			paths, weight := pt.AllTo(test.Query.To().ID())
			if weight != test.Weight {
				t.Errorf("%q: unexpected weight from AllTo: got:%f want:%f",
					test.Name, weight, test.Weight)
			}
			if weight := pt.WeightTo(test.Query.To().ID()); weight != test.Weight {
				t.Errorf("%q: unexpected weight from Weight: got:%f want:%f",
					test.Name, weight, test.Weight)
			}

			var gotPaths [][]int64
			if len(paths) != 0 {
				gotPaths = make([][]int64, len(paths))
			}
			for i, p := range paths {
				for _, v := range p {
					gotPaths[i] = append(gotPaths[i], v.ID())
				}
			}
			ordered.BySliceValues(gotPaths)
			if !reflect.DeepEqual(gotPaths, test.WantPaths) {
				t.Errorf("testing %q: unexpected shortest paths:\ngot: %v\nwant:%v",
					test.Name, gotPaths, test.WantPaths)
			}

			// Test absent paths.
			np, weight, unique := pt.To(test.NoPathFor.To().ID())
			if pt.From().ID() == test.NoPathFor.From().ID() && !(np == nil && math.IsInf(weight, 1) && !unique) {
				t.Errorf("%q: unexpected path:\ngot: path=%v weight=%f unique=%t\nwant:path=<nil> weight=+Inf unique=false",
					test.Name, np, weight, unique)
			}
		}
	}
}

type weightedTraverseGraph interface {
	traverse.Graph
	Weighted
}

type incremental struct {
	weightedTraverseGraph
}

func TestDijkstraAllPaths(t *testing.T) {
	t.Parallel()
	for _, test := range testgraphs.ShortestPathTests {
		g := test.Graph()
		for _, e := range test.Edges {
			g.SetWeightedEdge(e)
		}

		var (
			pt AllShortest

			panicked bool
		)
		func() {
			defer func() {
				panicked = recover() != nil
			}()
			pt = DijkstraAllPaths(g.(graph.Graph))
		}()
		if panicked || test.HasNegativeWeight {
			if !test.HasNegativeWeight {
				t.Errorf("%q: unexpected panic", test.Name)
			}
			if !panicked {
				t.Errorf("%q: expected panic for negative edge weight", test.Name)
			}
			continue
		}

		// Check all random paths returned are OK.
		for i := 0; i < 10; i++ {
			p, weight, unique := pt.Between(test.Query.From().ID(), test.Query.To().ID())
			if weight != test.Weight {
				t.Errorf("%q: unexpected weight from Between: got:%f want:%f",
					test.Name, weight, test.Weight)
			}
			if weight := pt.Weight(test.Query.From().ID(), test.Query.To().ID()); weight != test.Weight {
				t.Errorf("%q: unexpected weight from Weight: got:%f want:%f",
					test.Name, weight, test.Weight)
			}
			if unique != test.HasUniquePath {
				t.Errorf("%q: unexpected number of paths: got: unique=%t want: unique=%t",
					test.Name, unique, test.HasUniquePath)
			}

			var got []int64
			for _, n := range p {
				got = append(got, n.ID())
			}
			ok := len(got) == 0 && len(test.WantPaths) == 0
			for _, sp := range test.WantPaths {
				if reflect.DeepEqual(got, sp) {
					ok = true
					break
				}
			}
			if !ok {
				t.Errorf("%q: unexpected shortest path:\ngot: %v\nwant from:%v",
					test.Name, p, test.WantPaths)
			}
		}

		np, weight, unique := pt.Between(test.NoPathFor.From().ID(), test.NoPathFor.To().ID())
		if np != nil || !math.IsInf(weight, 1) || unique {
			t.Errorf("%q: unexpected path:\ngot: path=%v weight=%f unique=%t\nwant:path=<nil> weight=+Inf unique=false",
				test.Name, np, weight, unique)
		}

		paths, weight := pt.AllBetween(test.Query.From().ID(), test.Query.To().ID())
		if weight != test.Weight {
			t.Errorf("%q: unexpected weight from Between: got:%f want:%f",
				test.Name, weight, test.Weight)
		}

		var got [][]int64
		if len(paths) != 0 {
			got = make([][]int64, len(paths))
		}
		for i, p := range paths {
			for _, v := range p {
				got[i] = append(got[i], v.ID())
			}
		}
		ordered.BySliceValues(got)
		if !reflect.DeepEqual(got, test.WantPaths) {
			t.Errorf("testing %q: unexpected shortest paths:\ngot: %v\nwant:%v",
				test.Name, got, test.WantPaths)
		}

		nps, weight := pt.AllBetween(test.NoPathFor.From().ID(), test.NoPathFor.To().ID())
		if nps != nil || !math.IsInf(weight, 1) {
			t.Errorf("%q: unexpected path:\ngot: paths=%v weight=%f\nwant:path=<nil> weight=+Inf",
				test.Name, nps, weight)
		}
	}
}

func TestAllShortestAbsentNode(t *testing.T) {
	t.Parallel()
	g := simple.NewUndirectedGraph()
	g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(2)})
	paths := DijkstraAllPaths(g)
	// Confirm we have a good paths tree.
	if _, cost := paths.AllBetween(1, 2); cost != 1 {
		t.Errorf("unexpected cost between existing nodes: got:%v want:1", cost)
	}

	gotPath, cost, unique := paths.Between(0, 0)
	if cost != 0 {
		t.Errorf("unexpected cost from absent node to itself: got:%v want:0", cost)
	}
	if !unique {
		t.Error("unexpected non-unique path from absent node to itself")
	}
	wantPath := []graph.Node{node(0)}
	if !reflect.DeepEqual(gotPath, wantPath) {
		t.Errorf("unexpected path from absent node to itself: got:%#v want:%#v", gotPath, wantPath)
	}

	gotPaths, cost := paths.AllBetween(0, 0)
	if cost != 0 {
		t.Errorf("unexpected cost from absent node to itself: got:%v want:0", cost)
	}
	wantPaths := [][]graph.Node{{node(0)}}
	if !reflect.DeepEqual(gotPaths, wantPaths) {
		t.Errorf("unexpected paths from absent node to itself: got:%#v want:%#v", gotPaths, wantPaths)
	}
}
