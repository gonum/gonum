// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package concrete

import (
	"testing"

	"github.com/gonum/graph"
)

var _ graph.Graph = (*Graph)(nil)
var _ graph.Graph = (*Graph)(nil)

func TestAssertMutableNotDirected(t *testing.T) {
	var g graph.MutableUndirected = NewGraph()
	if _, ok := g.(graph.Directed); ok {
		t.Fatal("Graph is directed, but a MutableGraph cannot safely be directed!")
	}
}

func TestMaxID(t *testing.T) {
	g := NewGraph()
	nodes := make(map[graph.Node]struct{})
	for i := Node(0); i < 3; i++ {
		g.AddNode(i)
		nodes[i] = struct{}{}
	}
	g.RemoveNode(Node(0))
	delete(nodes, Node(0))
	g.RemoveNode(Node(2))
	delete(nodes, Node(2))
	n := Node(g.NewNodeID())
	g.AddNode(n)
	if !g.Has(n) {
		t.Error("added node does not exist in graph")
	}
	if _, exists := nodes[n]; exists {
		t.Errorf("Created already existing node id: %v", n.ID())
	}
}
