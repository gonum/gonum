// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/gonum/graph"
	"github.com/gonum/graph/internal/ordered"
	"github.com/gonum/graph/path/internal/testgraphs"
)

func TestFloydWarshall(t *testing.T) {
	for _, test := range testgraphs.ShortestPathTests {
		g := test.Graph()
		for _, e := range test.Edges {
			g.SetEdge(e)
		}

		pt, ok := FloydWarshall(g.(graph.Graph))
		if test.HasNegativeCycle {
			if ok {
				t.Errorf("%q: expected negative cycle", test.Name)
			}
			continue
		}
		if !ok {
			t.Fatalf("%q: unexpected negative cycle", test.Name)
		}

		// Check all random paths returned are OK.
		for i := 0; i < 10; i++ {
			p, weight, unique := pt.Between(test.Query.From(), test.Query.To())
			if weight != test.Weight {
				t.Errorf("%q: unexpected weight from Between: got:%f want:%f",
					test.Name, weight, test.Weight)
			}
			if weight := pt.Weight(test.Query.From(), test.Query.To()); weight != test.Weight {
				t.Errorf("%q: unexpected weight from Weight: got:%f want:%f",
					test.Name, weight, test.Weight)
			}
			if unique != test.HasUniquePath {
				t.Errorf("%q: unexpected number of paths: got: unique=%t want: unique=%t",
					test.Name, unique, test.HasUniquePath)
			}

			var got []int
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

		np, weight, unique := pt.Between(test.NoPathFor.From(), test.NoPathFor.To())
		if np != nil || !math.IsInf(weight, 1) || unique != false {
			t.Errorf("%q: unexpected path:\ngot: path=%v weight=%f unique=%t\nwant:path=<nil> weight=+Inf unique=false",
				test.Name, np, weight, unique)
		}

		paths, weight := pt.AllBetween(test.Query.From(), test.Query.To())
		if weight != test.Weight {
			t.Errorf("%q: unexpected weight from Between: got:%f want:%f",
				test.Name, weight, test.Weight)
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
		sort.Sort(ordered.BySliceValues(got))
		if !reflect.DeepEqual(got, test.WantPaths) {
			t.Errorf("testing %q: unexpected shortest paths:\ngot: %v\nwant:%v",
				test.Name, got, test.WantPaths)
		}

		nps, weight := pt.AllBetween(test.NoPathFor.From(), test.NoPathFor.To())
		if nps != nil || !math.IsInf(weight, 1) {
			t.Errorf("%q: unexpected path:\ngot: paths=%v weight=%f\nwant:path=<nil> weight=+Inf",
				test.Name, nps, weight)
		}
	}
}
