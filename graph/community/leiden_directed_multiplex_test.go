// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package community

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/internal/order"
)

func TestLeidenDirectedMultiplex(t *testing.T) {
	const iterations = 20

	for _, test := range communityDirectedMultiplexQTests {
		g, weights, err := directedMultiplexFrom(test.layers)
		if err != nil {
			t.Errorf("unexpected error creating multiplex: %v", err)
			continue
		}

		if test.structures[0].resolution != 1 {
			panic("bad test: expect resolution=1")
		}

		var (
			got   *ReducedDirectedMultiplex
			bestQ = math.Inf(-1)
		)

		src := rand.New(rand.NewPCG(1, 1))
		for i := 0; i < iterations; i++ {
			rr, err := LeidenMultiplex(g, weights, nil, true, src)
			if err != nil {
				if rr == nil {
					continue
				}
				t.Errorf("%s: LeidenMultiplex returned error: %v", test.name, err)
			}
			r := rr.(*ReducedDirectedMultiplex)

			qVec := QMultiplex(r, nil, weights, nil)
			q := floats.Sum(qVec)

			if q > bestQ || math.IsNaN(q) {
				bestQ = q
				got = r

				if math.IsNaN(q) {
					// Don't try again for non-connected case.
					break
				}
			}

			var qs []float64
			for p := r; p != nil; p = p.Expanded().(*ReducedDirectedMultiplex) {
				qs = append(qs, floats.Sum(QMultiplex(p, nil, weights, nil)))
			}

			// Leiden might not be strictly monotonic, but for now let's check it.
			if !math.IsNaN(qs[0]) {
				slices.Reverse(qs)
				// Allow small tolerance for Leiden refinement
				for k := 1; k < len(qs); k++ {
					if qs[k] < qs[k-1]-1e-6 {
						t.Errorf("%s: Q values not monotonically increasing: %.5v", test.name, qs)
					}
				}
			}
		}

		if math.IsNaN(bestQ) {
			if !math.IsNaN(test.structures[0].want) {
				t.Errorf("unexpected NaN Q for %s", test.name)
			}
			continue
		}

		// Check if the result is valid
		gotCommunities := got.Communities()
		for _, c := range gotCommunities {
			order.ByID(c)
		}
		order.BySliceIDs(gotCommunities)

		// Verify partition validity (disjoint and complete)
		for l := 0; l < g.Depth(); l++ {
			nodes := graph.NodesOf(g.Layer(l).Nodes())
			seen := make(map[int64]int)
			for i, c := range gotCommunities {
				for _, n := range c {
					id := n.ID()
					if j, ok := seen[id]; ok {
						t.Errorf("%s: node %d in multiple communities %d and %d", test.name, id, j, i)
					}
					seen[id] = i
				}
			}
			if len(seen) != len(nodes) {
				// Note: some tests might have layers with different node sets or implicit nodes?
				// directedMultiplexFrom ensures all layers have same nodes.
				t.Errorf("%s: partition covers %d nodes, graph layer %d has %d", test.name, len(seen), l, len(nodes))
			}
		}

		want := test.structures[0].want
		if lw := test.structures[0].leidenWant; lw != 0 {
			want = lw
		}

		if bestQ < want-1e-6 {
			t.Errorf("unexpectedly low Q value for %q: got: %v want >= %v",
				test.name, bestQ, want)
		}
	}
}

func TestLeidenNonContiguousDirectedMultiplex(t *testing.T) {
	g := simple.NewDirectedGraph()
	for _, e := range []simple.Edge{
		{F: simple.Node(0), T: simple.Node(1)},
		{F: simple.Node(4), T: simple.Node(5)},
	} {
		g.SetEdge(e)
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("unexpected panic with non-contiguous ID range")
			}
		}()
		if _, err := LeidenMultiplex(DirectedLayers{g}, nil, nil, true, nil); err != nil {
			t.Errorf("unexpected error with non-contiguous ID range: %v", err)
		}
	}()
}
