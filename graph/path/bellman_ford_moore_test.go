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
	"gonum.org/v1/gonum/graph/traverse"
)

func TestBellmanFordFrom(t *testing.T) {
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
			pt, ok := BellmanFordFrom(test.Query.From(), tg.g)
			if test.HasNegativeCycle {
				if ok {
					t.Errorf("%q %s: expected negative cycle", test.Name, tg.typ)
				}
			} else if !ok {
				t.Fatalf("%q %s: unexpected negative cycle", test.Name, tg.typ)
			}

			if pt.From().ID() != test.Query.From().ID() {
				t.Fatalf("%q %s: unexpected from node ID: got:%d want:%d", test.Name, tg.typ, pt.From().ID(), test.Query.From().ID())
			}

			p, weight := pt.To(test.Query.To().ID())
			if weight != test.Weight {
				t.Errorf("%q %s: unexpected weight from To: got:%f want:%f",
					test.Name, tg.typ, weight, test.Weight)
			}
			if weight := pt.WeightTo(test.Query.To().ID()); !math.IsInf(test.Weight, -1) && weight != test.Weight {
				t.Errorf("%q %s: unexpected weight from Weight: got:%f want:%f",
					test.Name, tg.typ, weight, test.Weight)
			}

			var got []int64
			for _, n := range p {
				got = append(got, n.ID())
			}
			ok = len(got) == 0 && len(test.WantPaths) == 0
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

func TestBellmanFordAllFrom(t *testing.T) {
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
			pt, ok := BellmanFordAllFrom(test.Query.From(), tg.g)
			if test.HasNegativeCycle {
				if ok {
					t.Errorf("%q %s: expected negative cycle", test.Name, tg.typ)
				}
			} else if !ok {
				t.Fatalf("%q %s: unexpected negative cycle", test.Name, tg.typ)
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
			if weight := pt.WeightTo(test.Query.To().ID()); !math.IsInf(test.Weight, -1) && weight != test.Weight {
				t.Errorf("%q %s: unexpected weight from Weight: got:%f want:%f",
					test.Name, tg.typ, weight, test.Weight)
			}

			var gotPath []int64
			for _, n := range p {
				gotPath = append(gotPath, n.ID())
			}
			ok = len(gotPath) == 0 && len(test.WantPaths) == 0
			for _, sp := range test.WantPaths {
				if reflect.DeepEqual(gotPath, sp) {
					ok = true
					break
				}
			}
			if !ok {
				t.Errorf("%q %s: unexpected shortest path:\ngot: %v\nwant from:%v",
					test.Name, tg.typ, p, test.WantPaths)
			}
			if unique != test.HasUniquePath {
				t.Errorf("%q %s: unexpected uniqueness from To: got:%t want:%t (%d paths)",
					test.Name, tg.typ, unique, test.HasUniquePath, len(test.WantPaths))
			}

			// Test multiple path results.
			paths, weight := pt.AllTo(test.Query.To().ID())
			if weight != test.Weight {
				t.Errorf("%q %s: unexpected weight from AllTo: got:%f want:%f",
					test.Name, tg.typ, weight, test.Weight)
			}
			if weight := pt.WeightTo(test.Query.To().ID()); !math.IsInf(test.Weight, -1) && weight != test.Weight {
				t.Errorf("%q %s: unexpected weight from Weight: got:%f want:%f",
					test.Name, tg.typ, weight, test.Weight)
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
			if test.HasNegativeCycleInPath {
				if gotPaths != nil {
					t.Errorf("testing %q %s: unexpected shortest paths:\ngot: %v\nwant: []",
						test.Name, tg.typ, gotPaths)
				}
			} else {
				ordered.BySliceValues(gotPaths)
				if !reflect.DeepEqual(gotPaths, test.WantPaths) {
					t.Errorf("testing %q %s: unexpected shortest paths:\ngot: %v\nwant:%v",
						test.Name, tg.typ, gotPaths, test.WantPaths)
				}
			}

			// Test absent paths.
			np, weight, unique := pt.To(test.NoPathFor.To().ID())
			if pt.From().ID() == test.NoPathFor.From().ID() && !(np == nil && math.IsInf(weight, 1) && !unique) {
				t.Errorf("%q %s: unexpected path:\ngot: path=%v weight=%f unique=%t\nwant:path=<nil> weight=+Inf unique=false",
					test.Name, tg.typ, np, weight, unique)
			}
		}
	}
}
