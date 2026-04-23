// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package community

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"

	"gonum.org/v1/gonum/graph/simple"
)

// ExampleLeiden demonstrates running the Leiden algorithm on a small graph
// with a fixed seed for reproducible output.
func ExampleLeiden() {
	g := simple.NewUndirectedGraph()
	g.SetEdge(simple.Edge{F: simple.Node(0), T: simple.Node(1)})
	g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(2)})
	g.SetEdge(simple.Edge{F: simple.Node(2), T: simple.Node(0)})
	g.SetEdge(simple.Edge{F: simple.Node(2), T: simple.Node(3)})
	g.SetEdge(simple.Edge{F: simple.Node(3), T: simple.Node(4)})
	g.SetEdge(simple.Edge{F: simple.Node(4), T: simple.Node(5)})
	g.SetEdge(simple.Edge{F: simple.Node(5), T: simple.Node(3)})
	src := rand.New(rand.NewPCG(1, 1))
	r := Leiden(g, 1, src)
	communities := r.Communities()
	fmt.Println("communities:", len(communities))
	for i, c := range communities {
		var ids []int64
		for _, n := range c {
			ids = append(ids, n.ID())
		}
		slices.Sort(ids)
		fmt.Printf("  %d: %v\n", i, ids)
	}
	fmt.Printf("Q = %.4f\n", Q(r, nil, 1))
	// Output:
	// communities: 2
	//   0: [0 1 2]
	//   1: [3 4 5]
	// Q = 0.3571
}

// TestLeidenExample runs Leiden on a small graph and logs the result.
func TestLeidenExample(t *testing.T) {
	g := simple.NewUndirectedGraph()
	g.SetEdge(simple.Edge{F: simple.Node(0), T: simple.Node(1)})
	g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(2)})
	g.SetEdge(simple.Edge{F: simple.Node(2), T: simple.Node(0)})
	g.SetEdge(simple.Edge{F: simple.Node(2), T: simple.Node(3)})
	g.SetEdge(simple.Edge{F: simple.Node(3), T: simple.Node(4)})
	g.SetEdge(simple.Edge{F: simple.Node(4), T: simple.Node(5)})
	g.SetEdge(simple.Edge{F: simple.Node(5), T: simple.Node(3)})
	r := Leiden(g, 1, nil)
	communities := r.Communities()
	t.Logf("Leiden found %d communities", len(communities))
	for i, c := range communities {
		var ids []int64
		for _, n := range c {
			ids = append(ids, n.ID())
		}
		t.Logf("  community %d: %v", i, ids)
	}
	q := Q(r, nil, 1)
	t.Logf("Modularity Q = %.4f", q)
	if len(communities) == 0 {
		t.Error("expected at least one community")
	}
}
