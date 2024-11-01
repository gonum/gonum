package flow

import (
	"testing"

	"gonum.org/v1/gonum/graph/simple"
)

func TestInterval(t *testing.T) {
	// copy from book
	g := simple.NewDirectedGraph()
	n1 := simple.Node(1)
	n2 := simple.Node(2)
	n3 := simple.Node(3)
	n4 := simple.Node(4)
	n5 := simple.Node(5)
	n6 := simple.Node(6)
	n7 := simple.Node(7)
	n8 := simple.Node(8)
	n9 := simple.Node(9)
	n10 := simple.Node(10)
	n11 := simple.Node(11)
	n12 := simple.Node(12)
	n13 := simple.Node(13)
	n14 := simple.Node(14)
	n15 := simple.Node(15)
	g.AddNode(n1)
	g.AddNode(n2)
	g.AddNode(n3)
	g.AddNode(n4)
	g.AddNode(n5)
	g.AddNode(n6)
	g.AddNode(n7)
	g.AddNode(n8)
	g.AddNode(n9)
	g.AddNode(n10)
	g.AddNode(n11)
	g.AddNode(n12)
	g.AddNode(n13)
	g.AddNode(n14)
	g.AddNode(n15)

	g.SetEdge(g.NewEdge(n1, n2))
	g.SetEdge(g.NewEdge(n1, n5))
	g.SetEdge(g.NewEdge(n2, n3))
	g.SetEdge(g.NewEdge(n2, n4))
	g.SetEdge(g.NewEdge(n3, n5))
	g.SetEdge(g.NewEdge(n4, n5))
	g.SetEdge(g.NewEdge(n5, n6))
	g.SetEdge(g.NewEdge(n6, n7))
	g.SetEdge(g.NewEdge(n6, n12))
	g.SetEdge(g.NewEdge(n7, n8))
	g.SetEdge(g.NewEdge(n7, n9))
	g.SetEdge(g.NewEdge(n8, n9))
	g.SetEdge(g.NewEdge(n8, n10))
	g.SetEdge(g.NewEdge(n9, n10))
	g.SetEdge(g.NewEdge(n10, n11))
	g.SetEdge(g.NewEdge(n12, n13))
	g.SetEdge(g.NewEdge(n14, n13))
	g.SetEdge(g.NewEdge(n13, n14))
	g.SetEdge(g.NewEdge(n14, n15))
	g.SetEdge(g.NewEdge(n15, n6))

	// test number of intervals
	intervals := Intervals(g, 1)
	if len(intervals) != 3 {
		t.Fatalf("Expected 3 intervals, got %d", len(intervals))
	}

	// test number of nodes
	interval := intervals[0]
	if len(interval.nodes) != 5 {
		t.Errorf("Expected 5 nodes in interval 1, got %d", len(interval.nodes))
	}

	interval2 := intervals[1]
	if len(interval2.nodes) != 7 {
		t.Errorf("Expected 7 nodes in interval 2, got %d", len(interval2.nodes))
	}

	interval3 := intervals[2]
	if len(interval3.nodes) != 3 {
		t.Errorf("Expected 3 nodes in interval 3, got %d", len(interval3.nodes))
	}
}
