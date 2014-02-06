package search_test

import (
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/search"
	"math"
	"testing"
)

func TestFloydWarshall(t *testing.T) {
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
