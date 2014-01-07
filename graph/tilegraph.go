package graph

import (
	"errors"
	"strings"
)

type TileGraph struct {
	tiles            []bool
	numRows, numCols int
}

func NewTileGraph(dimX, dimY int, isPassable bool) *TileGraph {
	tiles := make([]bool, dimX*dimY)
	if isPassable {
		for i, _ := range tiles {
			tiles[i] = true
		}
	}

	return &TileGraph{
		tiles:   tiles,
		numRows: dimX,
		numCols: dimY,
	}
}

func GenerateTileGraph(template string) (*TileGraph, error) {
	rows := strings.Split(strings.TrimSpace(template), "\n")

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
			case '\u2580':
				tiles = append(tiles, false)
			case ' ':
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

func (graph *TileGraph) SetPassability(row, col int, passability bool) {
	loc := row*graph.numCols + col
	if loc >= len(graph.tiles) || row < 0 || col < 0 {
		return
	}

	graph.tiles[loc] = passability
}

func (graph *TileGraph) String() string {
	var outString string
	for r := 0; r < graph.numRows; r++ {
		for c := 0; c < graph.numCols; c++ {
			if graph.tiles[r*graph.numCols+c] == false {
				outString += "\u2580" // Black square
			} else {
				outString += " " // Space
			}
		}

		outString += "\n"
	}

	return outString[:len(outString)-1] // Kill final newline
}

func (graph *TileGraph) PathString(path []Node) string {
	if path == nil || len(path) == 0 {
		return graph.String()
	}

	var outString string
	for r := 0; r < graph.numRows; r++ {
		for c := 0; c < graph.numCols; c++ {
			if id := r*graph.numCols + c; graph.tiles[id] == false {
				outString += "\u2580" // Black square
			} else if id == path[0].ID() {
				outString += "s"
			} else if id == path[len(path)-1].ID() {
				outString += "g"
			} else {
				toAppend := " "
				for _, num := range path[1 : len(path)-1] {
					if id == num.ID() {
						toAppend = "♥"
					}
				}
				outString += toAppend
			}
		}

		outString += "\n"
	}

	return outString[:len(outString)-1]
}

func (graph *TileGraph) Dimensions() (rows, cols int) {
	return graph.numRows, graph.numCols
}

func (graph *TileGraph) IDToCoords(id int) (row, col int) {
	col = (id % graph.numCols)
	row = (id - col) / graph.numCols

	return row, col
}

func (graph *TileGraph) CoordsToID(row, col int) (id int) {
	if row < 0 || row >= graph.numRows || col < 0 || col >= graph.numCols {
		return -1
	}
	id = row*graph.numCols + col

	return id
}

func (graph *TileGraph) CoordsToNode(row, col int) (node Node) {
	id := graph.CoordsToID(row, col)
	if id == -1 {
		return nil
	} else {
		return GonumNode(id)
	}
}

func (graph *TileGraph) Successors(node Node) []Node {
	id := node.ID()
	if id < 0 || id >= len(graph.tiles) || graph.tiles[id] == false {
		return nil
	}

	row, col := graph.IDToCoords(id)

	neighbors := []Node{graph.CoordsToNode(row-1, col), graph.CoordsToNode(row+1, col), graph.CoordsToNode(row, col-1), graph.CoordsToNode(row, col+1)}
	realNeighbors := make([]Node, 0, 4) // Will overallocate sometimes, but not by much. Not a big deal
	for _, neighbor := range neighbors {
		if neighbor != nil && graph.tiles[neighbor.ID()] == true {
			realNeighbors = append(realNeighbors, neighbor)
		}
	}

	return realNeighbors
}

func (graph *TileGraph) IsSuccessor(node, successor Node) bool {
	id, succ := node.ID(), successor.ID()
	return (id >= 0 && id < len(graph.tiles) && graph.tiles[id] == true) && (succ >= 0 && succ < len(graph.tiles) && graph.tiles[succ] == true)
}

func (graph *TileGraph) Predecessors(node Node) []Node {
	return graph.Successors(node)
}

func (graph *TileGraph) IsPredecessor(node, pred Node) bool {
	return graph.IsSuccessor(node, pred)
}

func (graph *TileGraph) IsAdjacent(id, neighbor Node) bool {
	return graph.IsSuccessor(id, neighbor)
}

func (graph *TileGraph) NodeExists(node Node) bool {
	id := node.ID()
	return id >= 0 && id < len(graph.tiles) && graph.tiles[id] == true
}

func (graph *TileGraph) Degree(node Node) int {
	return len(graph.Successors(node)) * 2
}

func (graph *TileGraph) EdgeList() []Edge {
	edges := make([]Edge, 0)
	for id, passable := range graph.tiles {
		if !passable {
			continue
		}

		for _, succ := range graph.Successors(GonumNode(id)) {
			edges = append(edges, GonumEdge{GonumNode(id), succ})
		}
	}

	return edges
}

func (graph *TileGraph) NodeList() []Node {
	nodes := make([]Node, 0)
	for id, passable := range graph.tiles {
		if !passable {
			continue
		}

		nodes = append(nodes, GonumNode(id))
	}

	return nodes
}

func (graph *TileGraph) IsDirected() bool {
	return false
}
