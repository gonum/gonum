package search_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/search"
	//	"github.com/gonum/graph/set"
)

func TestSimpleAStar(t *testing.T) {
	tg, err := concrete.GenerateTileGraph("▀  ▀\n▀▀ ▀\n▀▀ ▀\n▀▀ ▀")
	if err != nil {
		t.Fatal("Couldn't generate tilegraph")
	}

	path, cost, _ := search.AStar(concrete.GonumNode(1), concrete.GonumNode(14), tg, nil, nil)
	if math.Abs(cost-4.0) > .00001 {
		t.Errorf("A* reports incorrect cost for simple tilegraph search")
	}

	if path == nil {
		t.Fatalf("A* fails to find path for for simple tilegraph search")
	} else {
		correctPath := []int{1, 2, 6, 10, 14}
		if len(path) != len(correctPath) {
			t.Fatalf("Astar returns wrong length path for simple tilegraph search")
		}
		for i, node := range path {
			if node.ID() != correctPath[i] {
				t.Errorf("Astar returns wrong path at step", i, "got:", node, "actual:", correctPath[i])
			}
		}
	}
}

func TestBiggerAStar(t *testing.T) {
	tg := concrete.NewTileGraph(3, 3, true)

	path, cost, _ := search.AStar(concrete.GonumNode(0), concrete.GonumNode(8), tg, nil, nil)

	if math.Abs(cost-4.0) > .00001 || !search.IsPath(path, tg) {
		t.Error("Non-optimal or impossible path found for 3x3 grid")
	}

	tg = concrete.NewTileGraph(1000, 1000, true)
	path, cost, _ = search.AStar(concrete.GonumNode(00), concrete.GonumNode(999*1000+999), tg, nil, nil)
	if !search.IsPath(path, tg) || cost != 1998.0 {
		t.Error("Non-optimal or impossible path found for 100x100 grid; cost:", cost, "path:\n"+tg.PathString(path))
	}
}

func TestObstructedAStar(t *testing.T) {
	tg := concrete.NewTileGraph(10, 10, true)

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
	path, cost1, expanded := search.AStar(concrete.GonumNode(5), tg.CoordsToNode(rows-1, cols-1), tg, nil, nil)

	if !search.IsPath(path, tg) {
		t.Error("Path doesn't exist in obstructed graph")
	}

	ManhattanHeuristic := func(n1, n2 graph.Node) float64 {
		id1, id2 := n1.ID(), n2.ID()
		r1, c1 := tg.IDToCoords(id1)
		r2, c2 := tg.IDToCoords(id2)

		return math.Abs(float64(r1)-float64(r2)) + math.Abs(float64(c1)-float64(c2))
	}

	path, cost2, expanded2 := search.AStar(concrete.GonumNode(5), tg.CoordsToNode(rows-1, cols-1), tg, nil, ManhattanHeuristic)
	if !search.IsPath(path, tg) {
		t.Error("Path doesn't exist when using heuristic on obstructed graph")
	}

	if math.Abs(cost1-cost2) > .00001 {
		t.Error("Cost when using admissible heuristic isn't approximately equal to cost without it")
	}

	if expanded2 > expanded {
		t.Error("Using admissible, consistent heuristic expanded more nodes than null heuristic (possible, but unlikely -- suggests an error somewhere)")
	}

}

func TestNoPathAStar(t *testing.T) {
	tg := concrete.NewTileGraph(5, 5, true)

	// Creates a "wall" down the middle row
	tg.SetPassability(2, 0, false)
	tg.SetPassability(2, 1, false)
	tg.SetPassability(2, 2, false)
	tg.SetPassability(2, 3, false)
	tg.SetPassability(2, 4, false)

	rows, _ := tg.Dimensions()
	path, _, _ := search.AStar(tg.CoordsToNode(0, 2), tg.CoordsToNode(rows-1, 2), tg, nil, nil)

	if len(path) > 0 { // Note that a nil slice will return len of 0, this won't panic
		t.Error("A* finds path where none exists")
	}
}

/*func TestSmallAStar(t *testing.T) {
	gg := newSmallGonumGraph()
	heur := newSmallHeuristic()
	if ok, edge, goal := monotonic(gg, heur); !ok {
		t.Fatalf("non-monotonic heuristic.  edge: %v goal: %v", edge, goal)
	}
	for _, start := range gg.NodeList() {
		// get reference paths by Dijkstra
		dPaths, dCosts := search.Dijkstra(start, gg, nil)
		// assert that AStar finds each path
		for goalID, dPath := range dPaths {
			exp := fmt.Sprintln(dPath, dCosts[goalID])
			aPath, aCost, work := search.AStar(start, concrete.GonumNode(goalID), gg, nil, heur)
			fmt.Println()
			got := fmt.Sprintln(aPath, aCost)
			if got != exp {
				t.Error("expected", exp, "got", got)
			}
			t.Log(aPath, work)
		}
	}
}*/

func ExampleBreadthFirstSearch() {
	g := concrete.NewGonumGraph(true)
	var n0, n1, n2, n3 concrete.GonumNode = 0, 1, 2, 3
	g.AddNode(n0, []graph.Node{n1, n2})
	g.AddEdge(concrete.GonumEdge{n2, n3})
	path, v := search.BreadthFirstSearch(n0, n3, g)
	fmt.Println("path:", path)
	fmt.Println("nodes visited:", v)
	// Output:
	// path: [0 2 3]
	// nodes visited: 4
}

func newSmallGonumGraph() *concrete.GonumGraph {
	eds := []struct{ n1, n2, edgeCost int }{
		{1, 2, 7},
		{1, 3, 9},
		{1, 6, 14},
		{2, 3, 10},
		{2, 4, 15},
		{3, 4, 11},
		{3, 6, 2},
		{4, 5, 6},
		{5, 6, 9},
	}
	g := concrete.NewGonumGraph(false)
	for n := concrete.GonumNode(1); n <= 6; n++ {
		g.AddNode(n, nil)
	}
	for _, ed := range eds {
		e := concrete.GonumEdge{
			concrete.GonumNode(ed.n1),
			concrete.GonumNode(ed.n2),
		}
		g.AddEdge(e)
		g.SetEdgeCost(e, float64(ed.edgeCost))
	}
	return g
}

func newSmallHeuristic() func(n1, n2 graph.Node) float64 {
	nds := []struct{ id, x, y int }{
		{1, 0, 6},
		{2, 1, 0},
		{3, 8, 7},
		{4, 16, 0},
		{5, 17, 5},
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

func monotonic(g graph.CostGraph, heur func(n1, n2 graph.Node) float64) (bool, graph.Edge, graph.Node) {
	for _, goal := range g.NodeList() {
		for _, edge := range g.EdgeList() {
			head := edge.Head()
			tail := edge.Tail()
			if heur(head, goal) > g.Cost(head, tail)+heur(tail, goal) {
				return false, edge, goal
			}
		}
	}
	return true, nil, nil
}

// Test for correct result on a small graph easily solvable by hand
func TestDijkstraSmall(t *testing.T) {
	g := newSmallGonumGraph()
	paths, lens := search.Dijkstra(concrete.GonumNode(1), g, nil)
	s := fmt.Sprintln(len(paths), len(lens))
	for i := 1; i <= 6; i++ {
		s += fmt.Sprintln(paths[i], lens[i])
	}
	if s != `6 6
[1] 0
[1 2] 7
[1 3] 9
[1 3 4] 20
[1 3 6 5] 20
[1 3 6] 11
` {
		t.Fatal(s)
	}
}
