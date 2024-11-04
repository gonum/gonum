// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flow

import (
	"testing"

	"gonum.org/v1/gonum/graph/simple"
)

func TestInterval(t *testing.T) {
	// Graph from C. Cifuentes, "Reverse Compilation Techniques", 1994 (figure 6-23)
	g := simple.NewDirectedGraph()
	edges := []struct{ from, to int }{
		{1, 2},
		{1, 5},
		{2, 3},
		{2, 4},
		{3, 5},
		{4, 5},
		{5, 6},
		{6, 7},
		{6, 12},
		{7, 8},
		{7, 9},
		{8, 9},
		{8, 10},
		{9, 10},
		{10, 11},
		{12, 13},
		{14, 13},
		{13, 14},
		{14, 15},
		{15, 6},
	}

	nodes := [][]int64{
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10, 11, 12},
		{13, 14, 15},
	}

	for _, e := range edges {
		g.SetEdge(g.NewEdge(simple.Node(e.from), simple.Node(e.to)))
	}

	// test number of intervals
	ig := Intervals(g, 1)
	if len(ig.Intervals) != 3 {
		t.Fatalf("Expected 3 intervals, got %d", len(ig.Intervals))
	}

	// test number of nodes
	interval := ig.Intervals[0]
	if len(interval.nodes) != 5 {
		t.Errorf("Expected 5 nodes in interval 1, got %d", len(interval.nodes))
	}

	for _, node := range nodes[0] {
		if interval.nodes[node] == nil {
			t.Errorf("Unexpected node %d found in interval 1", node)
		}
	}

	interval2 := ig.Intervals[1]
	if len(interval2.nodes) != 7 {
		t.Errorf("Expected 7 nodes in interval 2, got %d", len(interval2.nodes))
	}

	for _, node := range nodes[1] {
		if interval2.nodes[node] == nil {
			t.Errorf("Unexpected node %d found in interval 2", node)
		}
	}

	interval3 := ig.Intervals[2]
	if len(interval3.nodes) != 3 {
		t.Errorf("Expected 3 nodes in interval 3, got %d", len(interval3.nodes))
	}

	for _, node := range nodes[2] {
		if interval3.nodes[node] == nil {
			t.Errorf("Unexpected node %d found in interval 3", node)
		}
	}

	// test interval edges
	toEdges := ig.To(6)
	if toEdges.Len() != 2 {
		t.Errorf("Expected 2 edges to node 6 in interval graph, got %d", len(ig.to))
	}

	toEdges2 := ig.To(13)
	if toEdges2.Len() != 1 {
		t.Errorf("Expected 1 edge to node 13 in interval graph, got %d", len(ig.to))
	}
}
