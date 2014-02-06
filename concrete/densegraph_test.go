package concrete_test

import (
	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
	"testing"
)

var _ graph.DirectedGraph = &concrete.DenseGraph{}
var _ graph.CrunchGraph = &concrete.DenseGraph{}

func TestBasicDenseImpassable(t *testing.T) {
	dg := concrete.NewDenseGraph(5, false)
	if dg == nil {
		t.Fatal("Directed graph could not be made")
	}

	for i := 0; i < 5; i++ {
		if !dg.NodeExists(concrete.GonumNode(i)) {
			t.Errorf("Node that should exist doesn't: %d", i)
		}

		if neighbors := dg.Neighbors(concrete.GonumNode(i)); len(neighbors) != 0 {
			t.Errorf("Node in impassable graph has a neighbor. Node: %d Neighbors: %v", i, neighbors)
		}
	}

	for i := 5; i < 10; i++ {
		if dg.NodeExists(concrete.GonumNode(i)) {
			t.Errorf("Node exists that shouldn't: %d")
		}
	}
}

func TestBasicDensePassable(t *testing.T) {
	dg := concrete.NewDenseGraph(5, true)
	if dg == nil {
		t.Fatal("Directed graph could not be made")
	}

	for i := 0; i < 5; i++ {
		if !dg.NodeExists(concrete.GonumNode(i)) {
			t.Errorf("Node that should exist doesn't: %d", i)
		}

		if neighbors := dg.Neighbors(concrete.GonumNode(i)); len(neighbors) != 5 {
			t.Errorf("Node in passable graph doesn't have the right number of neighbors. Node: %d Neighbors: %v", i, neighbors)
		}
	}

	for i := 5; i < 10; i++ {
		if dg.NodeExists(concrete.GonumNode(i)) {
			t.Errorf("Node exists that shouldn't: %d")
		}
	}
}
