// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"bytes"
	"errors"
	"math"
	"strings"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

const (
	blackSquare = '\u2580'
	space       = ' '
	startChar   = 's'
	goalChar    = 'g'
	pathChar    = '\u2665'
)

var inf = math.Inf(1)

type TileGraph struct {
	tiles            []bool
	numRows, numCols int
}

func NewTileGraph(dimX, dimY int, isPassable bool) *TileGraph {
	tiles := make([]bool, dimX*dimY)
	if isPassable {
		for i := range tiles {
			tiles[i] = true
		}
	}

	return &TileGraph{
		tiles:   tiles,
		numRows: dimX,
		numCols: dimY,
	}
}

func NewTileGraphFrom(text string) (*TileGraph, error) {
	rows := strings.Split(text, "\n")

	tiles := make([]bool, 0)

	colCheck := -1
	for _, colString := range rows {
		colCount := 0
		cols := strings.NewReader(colString)
		for cols.Len() != 0 {
			colCount += 1
			ch, _, err := cols.ReadRune()
			if err != nil {
				return nil, errors.New("Error while reading rune from input string")
			}

			switch ch {
			case blackSquare:
				tiles = append(tiles, false)
			case space:
				tiles = append(tiles, true)
			default:
				return nil, errors.New("Unrecognized character while reading input string")
			}
		}

		if colCheck == -1 {
			colCheck = colCount
		} else if colCheck != colCount {
			return nil, errors.New("Jagged rows, cannot generate graph.")
		}
	}

	return &TileGraph{
		tiles:   tiles,
		numRows: len(rows),
		numCols: colCheck,
	}, nil
}

func (g *TileGraph) SetPassability(row, col int, passability bool) {
	loc := row*g.numCols + col
	if loc >= len(g.tiles) || row < 0 || col < 0 {
		return
	}

	g.tiles[loc] = passability
}

func (g *TileGraph) String() string {
	var b bytes.Buffer
	for r := 0; r < g.numRows; r++ {
		if r != 0 {
			b.WriteByte('\n')
		}
		for c := 0; c < g.numCols; c++ {
			if g.tiles[r*g.numCols+c] == false {
				b.WriteRune(blackSquare)
			} else {
				b.WriteByte(space)
			}
		}
	}

	return b.String()
}

func (g *TileGraph) PathString(path []graph.Node) string {
	if len(path) == 0 {
		return g.String()
	}

	var b bytes.Buffer
	for r := 0; r < g.numRows; r++ {
		if r != 0 {
			b.WriteByte('\n')
		}
	row:
		for c := 0; c < g.numCols; c++ {
			if id := r*g.numCols + c; g.tiles[id] == false {
				b.WriteRune(blackSquare)
			} else if id == path[0].ID() {
				b.WriteRune(startChar)
			} else if id == path[len(path)-1].ID() {
				b.WriteRune(goalChar)
			} else {
				for _, num := range path[1 : len(path)-1] {
					if id == num.ID() {
						b.WriteRune(pathChar)
						continue row
					}
				}
				b.WriteByte(space)
			}
		}
	}

	return b.String()
}

func (g *TileGraph) Dimensions() (rows, cols int) {
	return g.numRows, g.numCols
}

func (g *TileGraph) IDToCoords(id int) (row, col int) {
	col = (id % g.numCols)
	row = (id - col) / g.numCols

	return row, col
}

func (g *TileGraph) CoordsToID(row, col int) int {
	if row < 0 || row >= g.numRows || col < 0 || col >= g.numCols {
		return -1
	}

	return row*g.numCols + col
}

func (g *TileGraph) CoordsToNode(row, col int) graph.Node {
	id := g.CoordsToID(row, col)
	if id == -1 {
		return nil
	}
	return concrete.Node(id)
}

func (g *TileGraph) Neighbors(n graph.Node) []graph.Node {
	id := n.ID()
	if !g.Has(n) {
		return nil
	}

	row, col := g.IDToCoords(id)

	neighbors := []graph.Node{g.CoordsToNode(row-1, col), g.CoordsToNode(row+1, col), g.CoordsToNode(row, col-1), g.CoordsToNode(row, col+1)}
	realNeighbors := make([]graph.Node, 0, 4) // Will overallocate sometimes, but not by much. Not a big deal
	for _, neigh := range neighbors {
		if neigh != nil && g.tiles[neigh.ID()] == true {
			realNeighbors = append(realNeighbors, neigh)
		}
	}

	return realNeighbors
}

func (g *TileGraph) EdgeBetween(n, neigh graph.Node) graph.Edge {
	if !g.Has(n) || !g.Has(neigh) {
		return nil
	}

	r1, c1 := g.IDToCoords(n.ID())
	r2, c2 := g.IDToCoords(neigh.ID())
	if (c1 == c2 && (r2 == r1+1 || r2 == r1-1)) || (r1 == r2 && (c2 == c1+1 || c2 == c1-1)) {
		return concrete.Edge{n, neigh}
	}

	return nil
}

func (g *TileGraph) Has(n graph.Node) bool {
	id := n.ID()
	return id >= 0 && id < len(g.tiles) && g.tiles[id] == true
}

func (g *TileGraph) Degree(n graph.Node) int {
	return len(g.Neighbors(n)) * 2
}

func (g *TileGraph) EdgeList() []graph.Edge {
	edges := make([]graph.Edge, 0)
	for id, passable := range g.tiles {
		if !passable {
			continue
		}

		for _, succ := range g.Neighbors(concrete.Node(id)) {
			edges = append(edges, concrete.Edge{concrete.Node(id), succ})
		}
	}

	return edges
}

func (g *TileGraph) Nodes() []graph.Node {
	nodes := make([]graph.Node, 0)
	for id, passable := range g.tiles {
		if !passable {
			continue
		}

		nodes = append(nodes, concrete.Node(id))
	}

	return nodes
}

func (g *TileGraph) Cost(e graph.Edge) float64 {
	if edge := g.EdgeBetween(e.From(), e.To()); edge != nil {
		return 1
	}

	return inf
}
