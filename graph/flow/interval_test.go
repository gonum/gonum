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

	intervalSize := []int64{
		5, 7, 3,
	}

	for _, e := range edges {
		g.SetEdge(g.NewEdge(simple.Node(e.from), simple.Node(e.to)))
	}

	// Test number of intervals
	ig := Intervals(g, 1)
	if len(ig.Intervals) != 3 {
		t.Fatalf("Expected 3 intervals, got %d", len(ig.Intervals))
	}

	for i, interval := range ig.Intervals {
		// Test correct number of nodes are found in interval
		if len(interval.nodes) != int(intervalSize[i]) {
			t.Errorf("Expected %d nodes in interval %d, got %d", intervalSize[i], i, len(interval.nodes))
		}

		// Test all expected in interval are present
		for _, node := range nodes[i] {
			if interval.nodes[node] == nil {
				t.Errorf("Node %d not found in interval %d", node, i)
			}
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
