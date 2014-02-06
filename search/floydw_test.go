package search_test

import (
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/search"
	"math"
	"testing"
)

func TestFWOneEdge(t *testing.T) {
	dg := concrete.NewDenseGraph(2, true)
	aPaths, sPath := search.FloydWarshall(dg, nil)

	path, cost, err := sPath(concrete.GonumNode(0), concrete.GonumNode(1))
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(cost-1.0) > .000001 {
		t.Errorf("FW got wrong cost %f", cost)
	}

	if len(path) != 2 || path[0].ID() != 0 && path[1].ID() != 1 {
		t.Errorf("Wrong path in FW %v", path)
	}

	paths, cost, err := aPaths(concrete.GonumNode(0), concrete.GonumNode(1))
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(cost-1.0) > .000001 {
		t.Errorf("FW got wrong cost %f", cost)
	}

	if len(paths) != 1 {
		t.Errorf("Didn't get right paths in FW %v", paths)
	}

	path = paths[0]
	if len(path) != 2 || path[0].ID() != 0 && path[1].ID() != 1 {
		t.Errorf("Wrong path in FW allpaths %v", path)
	}
}

func TestFWTwoPaths(t *testing.T) {
	dg := concrete.NewDenseGraph(5, false)
	// Adds two paths from 0->2 of equal length
	dg.SetEdgeCost(concrete.GonumNode(0), concrete.GonumNode(2), 2.0, true)
	dg.SetEdgeCost(concrete.GonumNode(0), concrete.GonumNode(1), 1.0, true)
	dg.SetEdgeCost(concrete.GonumNode(1), concrete.GonumNode(2), 1.0, true)

	aPaths, sPath := search.FloydWarshall(dg, nil)
	path, cost, err := sPath(concrete.GonumNode(0), concrete.GonumNode(2))
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(cost-2.0) > .00001 {
		t.Errorf("Path has incorrect cost, %f", cost)
	}

	if len(path) == 2 && path[0].ID() == 0 && path[1].ID() == 2 {
		t.Logf("Got correct path: %v", path)
	} else if len(path) == 3 && path[0].ID() == 0 && path[1].ID() == 1 && path[2].ID() == 2 {
		t.Logf("Got correct path %v", path)
	} else {
		t.Errorf("Got wrong path %v", path)
	}

	paths, cost, err := aPaths(concrete.GonumNode(0), concrete.GonumNode(2))

	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(cost-2.0) > .00001 {
		t.Errorf("All paths function gets incorrect cost, %f", cost)
	}

	if len(paths) != 2 {
		t.Fatalf("Didn't get all shortest paths %v", paths)
	}

	for _, path := range paths {
		if len(path) == 2 && path[0].ID() == 0 && path[1].ID() == 2 {
			t.Logf("Got correct path for all paths: %v", path)
		} else if len(path) == 3 && path[0].ID() == 0 && path[1].ID() == 1 && path[2].ID() == 2 {
			t.Logf("Got correct path for all paths %v", path)
		} else {
			t.Errorf("Got wrong path for all paths %v", path)
		}
	}
}

// Tests with multiple right paths, but also one dead-end path
// and one path that reaches the goal, but not optimally
func TestFWConfoundingPath(t *testing.T) {
	dg := concrete.NewDenseGraph(6, false)

	// Add a path from 0->5 of cost 4
	dg.SetEdgeCost(concrete.GonumNode(0), concrete.GonumNode(1), 1.0, true)
	dg.SetEdgeCost(concrete.GonumNode(1), concrete.GonumNode(2), 1.0, true)
	dg.SetEdgeCost(concrete.GonumNode(2), concrete.GonumNode(3), 1.0, true)
	dg.SetEdgeCost(concrete.GonumNode(3), concrete.GonumNode(5), 1.0, true)

	// Add direct edge to goal of cost 4
	dg.SetEdgeCost(concrete.GonumNode(0), concrete.GonumNode(5), 4.0, true)

	// Add edge to 3 that's overpriced
	dg.SetEdgeCost(concrete.GonumNode(0), concrete.GonumNode(3), 4.0, true)

	// Add very cheap edge to 4 which is a dead end
	dg.SetEdgeCost(concrete.GonumNode(0), concrete.GonumNode(4), 0.25, true)

	aPaths, sPath := search.FloydWarshall(dg, nil)

	path, cost, err := sPath(concrete.GonumNode(0), concrete.GonumNode(5))
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(cost-4.0) > .000001 {
		t.Error("Incorrect cost %f", cost)
	}

	if len(path) == 5 && path[0].ID() == 0 && path[1].ID() == 1 && path[2].ID() == 2 && path[3].ID() == 3 && path[4].ID() == 5 {
		t.Log("Correct path found for single path %v", path)
	} else if len(path) == 2 && path[0].ID() == 0 && path[1].ID() == 5 {
		t.Log("Correct path found for single path %v", path)
	} else {
		t.Error("Wrong path found for single path %v", path)
	}

	paths, cost, err := aPaths(concrete.GonumNode(0), concrete.GonumNode(5))
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(cost-4.0) > .000001 {
		t.Error("Incorrect cost %f", cost)
	}

	if len(paths) != 2 {
		t.Error("Wrong paths gooten for all paths %v", paths)
	}

	for _, path := range paths {
		if len(path) == 5 && path[0].ID() == 0 && path[1].ID() == 1 && path[2].ID() == 2 && path[3].ID() == 3 && path[4].ID() == 5 {
			t.Log("Correct path found for all paths %v", path)
		} else if len(path) == 2 && path[0].ID() == 0 && path[1].ID() == 5 {
			t.Log("Correct path found for all paths %v", path)
		} else {
			t.Error("Wrong path found for all paths %v", path)
		}
	}

	path, _, err = sPath(concrete.GonumNode(4), concrete.GonumNode(5))
	if err != nil {
		t.Log("Success!", err)
	} else {
		t.Error("Path was found by FW single path where one shouldn't be %v", path)
	}

	paths, _, err = aPaths(concrete.GonumNode(4), concrete.GonumNode(5))
	if err != nil {
		t.Log("Success!", err)
	} else {
		t.Error("Path was found by FW multi-path where one shouldn't be %v", paths)
	}
}
