package search_test

import (
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/search"
	"testing"
)

func TestFloydWarshall(t *testing.T) {
	dg := concrete.NewDenseGraph(5, false)
	// Adds two paths from 0->2 of equal length
	dg.SetEdgeCost(concrete.GonumNode(0), concrete.GonumNode(2), 2.0, true)
	dg.SetEdgeCost(concrete.GonumNode(0), concrete.GonumNode(1), 1.0, true)
	dg.SetEdgeCost(concrete.GonumNode(1), concrete.GonumNode(2), 1.0, true)

	search.FloydWarshall(dg, nil)
}
