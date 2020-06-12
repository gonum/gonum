// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph_test

import (
	"fmt"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/graphs/gen"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/simple"
)

var complementTests = []struct {
	g graph.Graph
}{
	{g: gnp(100, 0, rand.NewSource(1))},
	{g: gnp(100, 0.05, rand.NewSource(1))},
	{g: gnp(100, 0.5, rand.NewSource(1))},
	{g: gnp(100, 0.95, rand.NewSource(1))},
	{g: gnp(100, 1, rand.NewSource(1))},
}

func TestComplement(t *testing.T) {
	for _, test := range complementTests {
		n := len(graph.NodesOf(test.g.Nodes()))
		wantM := n * (n - 1) // Double counting edges, but no self-loops.

		var gotM int
		iter := test.g.Nodes()
		for iter.Next() {
			id := iter.Node().ID()
			to := test.g.From(id)
			for to.Next() {
				gotM++
			}
			toC := graph.Complement{test.g}.From(id)
			for toC.Next() {
				gotM++
			}
		}
		if gotM != wantM {
			t.Errorf("unexpected number of edges in sum of input and complement: got:%d want:%d", gotM, wantM)
		}
	}
}

func gnp(n int, p float64, src rand.Source) *simple.UndirectedGraph {
	g := simple.NewUndirectedGraph()
	err := gen.Gnp(g, n, p, src)
	if err != nil {
		panic(fmt.Sprintf("gnp: bad test: %v", err))
	}
	return g
}

var nodeFilterIteratorTests = []struct {
	src, filter graph.Nodes
	root        int64
	len         int
}{
	{src: iterator.NewOrderedNodes([]graph.Node{simple.Node(0)}), filter: graph.Empty, root: 0, len: 0},
	{src: iterator.NewOrderedNodes([]graph.Node{simple.Node(0), simple.Node(1)}), filter: graph.Empty, root: 0, len: 1},
	{src: iterator.NewOrderedNodes([]graph.Node{simple.Node(0), simple.Node(1), simple.Node(2)}), filter: iterator.NewOrderedNodes([]graph.Node{simple.Node(1)}), root: 0, len: 1},
}

func TestNodeFilterIterator(t *testing.T) {
	for _, test := range nodeFilterIteratorTests {
		it := graph.NewNodeFilterIterator(test.src, test.filter, test.root)
		if it.Len() < 0 {
			t.Logf("don't test indeterminate iterators: %T", it)
			continue
		}
		for i := 0; i < 2; i++ {
			n := it.Len()
			if n != test.len {
				t.Errorf("unexpected length of iterator construction/reset: got:%d want:%d", n, test.len)
			}
			for it.Next() {
				n--
			}
			if n != 0 {
				t.Errorf("unexpected remaining nodes after iterator completion: got:%d want:0", n)
			}
			it.Reset()
		}
	}
}
