package concrete_test

import (
	"testing"

	gr "github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

// Asserts that concrete graphs implement the correct interfaces
var _ gr.DirectedGraph = (*concrete.GonumGraph)(nil)
var _ gr.MutableGraph = (*concrete.GonumGraph)(nil)
var _ gr.UndirectedGraph = (*concrete.TileGraph)(nil)

func TestTileGraph(t *testing.T) {
	tg := concrete.NewTileGraph(4, 4, false)

	if tg == nil || tg.String() != "▀▀▀▀\n▀▀▀▀\n▀▀▀▀\n▀▀▀▀" {
		t.Fatal("Tile graph not generated correctly")
	}

	tg.SetPassability(0, 1, true)
	if tg == nil || tg.String() != "▀ ▀▀\n▀▀▀▀\n▀▀▀▀\n▀▀▀▀" {
		t.Fatal("Passability set incorrectly")
	}

	tg.SetPassability(0, 1, false)
	if tg == nil || tg.String() != "▀▀▀▀\n▀▀▀▀\n▀▀▀▀\n▀▀▀▀" {
		t.Fatal("Passability set incorrectly")
	}

	tg.SetPassability(0, 1, true)
	if tg == nil || tg.String() != "▀ ▀▀\n▀▀▀▀\n▀▀▀▀\n▀▀▀▀" {
		t.Fatal("Passability set incorrectly")
	}

	tg.SetPassability(0, 2, true)
	if tg == nil || tg.String() != "▀  ▀\n▀▀▀▀\n▀▀▀▀\n▀▀▀▀" {
		t.Fatal("Passability set incorrectly")
	}

	tg.SetPassability(1, 2, true)
	if tg == nil || tg.String() != "▀  ▀\n▀▀ ▀\n▀▀▀▀\n▀▀▀▀" {
		t.Fatal("Passability set incorrectly")
	}

	tg.SetPassability(2, 2, true)
	if tg == nil || tg.String() != "▀  ▀\n▀▀ ▀\n▀▀ ▀\n▀▀▀▀" {
		t.Fatal("Passability set incorrectly")
	}

	tg.SetPassability(3, 2, true)
	if tg == nil || tg.String() != "▀  ▀\n▀▀ ▀\n▀▀ ▀\n▀▀ ▀" {
		t.Fatal("Passability set incorrectly")
	}

	if tg2, err := concrete.GenerateTileGraph("▀  ▀\n▀▀ ▀\n▀▀ ▀\n▀▀ ▀"); err != nil {
		t.Error("Tile graph errored on interpreting valid template string\n▀  ▀\n▀▀ ▀\n▀▀ ▀\n▀▀ ▀")
	} else if tg2.String() != "▀  ▀\n▀▀ ▀\n▀▀ ▀\n▀▀ ▀" {
		t.Error("Tile graph failed to generate properly with input string\n▀  ▀\n▀▀ ▀\n▀▀ ▀\n▀▀ ▀")
	}

	if tg.CoordsToID(0, 0) != 0 {
		t.Error("Coords to ID fails on 0,0")
	} else if tg.CoordsToID(3, 3) != 15 {
		t.Error("Coords to ID fails on 3,3")
	} else if tg.CoordsToID(0, 3) != 3 {
		t.Error("Coords to ID fails on 0,3")
	} else if tg.CoordsToID(3, 0) != 12 {
		t.Error("Coords to ID fails on 3,0")
	}

	if r, c := tg.IDToCoords(0); r != 0 || c != 0 {
		t.Error("ID to Coords fails on 0,0")
	} else if r, c := tg.IDToCoords(15); r != 3 || c != 3 {
		t.Error("ID to Coords fails on 3,3")
	} else if r, c := tg.IDToCoords(3); r != 0 || c != 3 {
		t.Error("ID to Coords fails on 0,3")
	} else if r, c := tg.IDToCoords(12); r != 3 || c != 0 {
		t.Error("ID to Coords fails on 3,0")
	}

	if succ := tg.Neighbors(concrete.GonumNode(0)); succ != nil || len(succ) != 0 {
		t.Error("Successors for impassable tile not 0")
	}

	if succ := tg.Neighbors(concrete.GonumNode(2)); succ == nil || len(succ) != 2 {
		t.Error("Incorrect number of successors for (0,2)")
	} else {
		for _, s := range succ {
			if s.ID() != 1 && s.ID() != 6 {
				t.Error("Successor for (0,2) neither (0,1) nor (1,2)")
			}
		}
	}

	if tg.Degree(concrete.GonumNode(2)) != 4 {
		t.Error("Degree returns incorrect number for (0,2)")
	}
	if tg.Degree(concrete.GonumNode(1)) != 2 {
		t.Error("Degree returns incorrect number for (0,2)")
	}
	if tg.Degree(concrete.GonumNode(0)) != 0 {
		t.Error("Degree returns incorrect number for impassable tile (0,0)")
	}

}

func TestTileGraphCoverage(t *testing.T) {
	concrete.GenerateTileGraph("▀ x")   // x invalid
	concrete.GenerateTileGraph("▀ \n▀") // row lengths not equal
	concrete.GenerateTileGraph("▀ ▀")
	tg := concrete.NewTileGraph(2, 3, true)
	tg.SetPassability(0, -1, true)
	tg.SetPassability(0, 1, false)
	tg.SetPassability(0, 4, false)
	tg.PathString(nil)
	n1 := concrete.GonumNode(1)
	n2 := concrete.GonumNode(2)
	n3 := concrete.GonumNode(3)
	n5 := concrete.GonumNode(5)
	p := []gr.Node{n3, n2, n1, n5}
	tg.PathString(p)
	tg.Dimensions()
	tg.IDToCoords(0)
	tg.CoordsToNode(-1, 0)
	tg.CoordsToNode(0, 0)
	tg.Neighbors(n1)
	tg.Neighbors(n2)
	tg.IsNeighbor(n1, n2)
	tg.NodeExists(n1)
	tg.Degree(n1)
	tg.EdgeList()
	tg.NodeList()
	tg.IsDirected()
	tg.Cost(n1, n2)
	tg.Cost(n2, n5)
}
