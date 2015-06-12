// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/internal"
	"github.com/gonum/graph/path"
	"github.com/gonum/graph/topo"
)

func TestSimpleAStar(t *testing.T) {
	tg, err := internal.NewTileGraphFrom("" +
		"▀  ▀\n" +
		"▀▀ ▀\n" +
		"▀▀ ▀\n" +
		"▀▀ ▀",
	)
	if err != nil {
		t.Fatalf("Couldn't generate tilegraph: %v", err)
	}

	p, cost, _ := path.AStar(concrete.Node(1), concrete.Node(14), tg, nil)
	if math.Abs(cost-4) > 1e-5 {
		t.Errorf("A* reports incorrect cost for simple tilegraph search")
	}

	if p == nil {
		t.Fatalf("A* fails to find path for for simple tilegraph search")
	} else {
		correctPath := []int{1, 2, 6, 10, 14}
		if len(p) != len(correctPath) {
			t.Fatalf("Astar returns wrong length path for simple tilegraph search")
		}
		for i, node := range p {
			if node.ID() != correctPath[i] {
				t.Errorf("Astar returns wrong path at step", i, "got:", node, "actual:", correctPath[i])
			}
		}
	}
}

func TestBiggerAStar(t *testing.T) {
	tg := internal.NewTileGraph(3, 3, true)

	p, cost, _ := path.AStar(concrete.Node(0), concrete.Node(8), tg, nil)

	if math.Abs(cost-4) > 1e-5 || !topo.IsPathIn(tg, p) {
		t.Error("Non-optimal or impossible path found for 3x3 grid")
	}

	tg = internal.NewTileGraph(1000, 1000, true)
	p, cost, _ = path.AStar(concrete.Node(0), concrete.Node(999*1000+999), tg, nil)
	if !topo.IsPathIn(tg, p) || cost != 1998 {
		t.Error("Non-optimal or impossible path found for 100x100 grid; cost:", cost, "path:\n"+tg.PathString(p))
	}
}

func TestObstructedAStar(t *testing.T) {
	tg := internal.NewTileGraph(10, 10, true)

	// Creates a partial "wall" down the middle row with a gap down the left side
	tg.SetPassability(4, 1, false)
	tg.SetPassability(4, 2, false)
	tg.SetPassability(4, 3, false)
	tg.SetPassability(4, 4, false)
	tg.SetPassability(4, 5, false)
	tg.SetPassability(4, 6, false)
	tg.SetPassability(4, 7, false)
	tg.SetPassability(4, 8, false)
	tg.SetPassability(4, 9, false)

	rows, cols := tg.Dimensions()
	p, cost1, expanded := path.AStar(concrete.Node(5), tg.CoordsToNode(rows-1, cols-1), tg, nil)

	if !topo.IsPathIn(tg, p) {
		t.Error("Path doesn't exist in obstructed graph")
	}

	ManhattanHeuristic := func(n1, n2 graph.Node) float64 {
		id1, id2 := n1.ID(), n2.ID()
		r1, c1 := tg.IDToCoords(id1)
		r2, c2 := tg.IDToCoords(id2)

		return math.Abs(float64(r1)-float64(r2)) + math.Abs(float64(c1)-float64(c2))
	}

	p, cost2, expanded2 := path.AStar(concrete.Node(5), tg.CoordsToNode(rows-1, cols-1), tg, ManhattanHeuristic)
	if !topo.IsPathIn(tg, p) {
		t.Error("Path doesn't exist when using heuristic on obstructed graph")
	}

	if math.Abs(cost1-cost2) > 1e-5 {
		t.Error("Cost when using admissible heuristic isn't approximately equal to cost without it")
	}

	if expanded2 > expanded {
		t.Error("Using admissible, consistent heuristic expanded more nodes than null heuristic (possible, but unlikely -- suggests an error somewhere)")
	}

}

func TestNoPathAStar(t *testing.T) {
	tg := internal.NewTileGraph(5, 5, true)

	// Creates a "wall" down the middle row
	tg.SetPassability(2, 0, false)
	tg.SetPassability(2, 1, false)
	tg.SetPassability(2, 2, false)
	tg.SetPassability(2, 3, false)
	tg.SetPassability(2, 4, false)

	rows, _ := tg.Dimensions()
	path, _, _ := path.AStar(tg.CoordsToNode(0, 2), tg.CoordsToNode(rows-1, 2), tg, nil)

	if len(path) > 0 { // Note that a nil slice will return len of 0, this won't panic
		t.Error("A* finds path where none exists")
	}
}

func TestSmallAStar(t *testing.T) {
	g := newSmallGonumGraph()
	heur := newSmallHeuristic()
	if ok, edge, goal := monotonic(g, heur); !ok {
		t.Fatalf("non-monotonic heuristic.  edge: %v goal: %v", edge, goal)
	}

	ps := path.DijkstraAllPaths(g)
	for _, start := range g.Nodes() {
		for _, goal := range g.Nodes() {
			gotPath, gotWeight, _ := path.AStar(start, goal, g, heur)
			wantPath, wantWeight, _ := ps.Between(start, goal)
			if gotWeight != wantWeight {
				t.Errorf("unexpected A* path weight from %v to %v result: got:%s want:%s",
					start, goal, gotWeight, wantWeight)
			}
			if !reflect.DeepEqual(gotPath, wantPath) {
				t.Errorf("unexpected A* path from %v to %v result:\ngot: %v\nwant:%v",
					start, goal, gotPath, wantPath)
			}
		}
	}
}

func newSmallGonumGraph() *concrete.Graph {
	eds := []struct{ n1, n2, edgeCost int }{
		{1, 2, 7},
		{1, 3, 9},
		{1, 6, 14},
		{2, 3, 10},
		{2, 4, 15},
		{3, 4, 11},
		{3, 6, 2},
		{4, 5, 7},
		{5, 6, 9},
	}
	g := concrete.NewGraph()
	for n := concrete.Node(1); n <= 6; n++ {
		g.AddNode(n)
	}
	for _, ed := range eds {
		e := concrete.Edge{
			concrete.Node(ed.n1),
			concrete.Node(ed.n2),
		}
		g.SetEdge(e, float64(ed.edgeCost))
	}
	return g
}

func newSmallHeuristic() func(n1, n2 graph.Node) float64 {
	nds := []struct{ id, x, y int }{
		{1, 0, 6},
		{2, 1, 0},
		{3, 8, 7},
		{4, 16, 0},
		{5, 17, 6},
		{6, 9, 8},
	}
	return func(n1, n2 graph.Node) float64 {
		i1 := n1.ID() - 1
		i2 := n2.ID() - 1
		dx := nds[i2].x - nds[i1].x
		dy := nds[i2].y - nds[i1].y
		return math.Hypot(float64(dx), float64(dy))
	}
}

type costEdgeListGraph interface {
	graph.Weighter
	path.EdgeListGraph
}

func monotonic(g costEdgeListGraph, heur func(n1, n2 graph.Node) float64) (bool, graph.Edge, graph.Node) {
	for _, goal := range g.Nodes() {
		for _, edge := range g.Edges() {
			from := edge.From()
			to := edge.To()
			if heur(from, goal) > g.Weight(edge)+heur(to, goal) {
				return false, edge, goal
			}
		}
	}
	return true, nil, nil
}
