package concrete

import (
	"errors"
	"math"
	"strings"

	"github.com/gonum/graph"
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

func (gr *TileGraph) SetPassability(row, col int, passability bool) {
	loc := row*gr.numCols + col
	if loc >= len(gr.tiles) || row < 0 || col < 0 {
		return
	}

	gr.tiles[loc] = passability
}

func (gr *TileGraph) String() string {
	var outString string
	for r := 0; r < gr.numRows; r++ {
		for c := 0; c < gr.numCols; c++ {
			if gr.tiles[r*gr.numCols+c] == false {
				outString += "\u2580" // Black square
			} else {
				outString += " " // Space
			}
		}

		outString += "\n"
	}

	return outString[:len(outString)-1] // Kill final newline
}

func (gr *TileGraph) PathString(path []graph.Node) string {
	if path == nil || len(path) == 0 {
		return gr.String()
	}

	var outString string
	for r := 0; r < gr.numRows; r++ {
		for c := 0; c < gr.numCols; c++ {
			if id := r*gr.numCols + c; gr.tiles[id] == false {
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

func (gr *TileGraph) Dimensions() (rows, cols int) {
	return gr.numRows, gr.numCols
}

func (gr *TileGraph) IDToCoords(id int) (row, col int) {
	col = (id % gr.numCols)
	row = (id - col) / gr.numCols

	return row, col
}

func (gr *TileGraph) CoordsToID(row, col int) (id int) {
	if row < 0 || row >= gr.numRows || col < 0 || col >= gr.numCols {
		return -1
	}
	id = row*gr.numCols + col

	return id
}

func (gr *TileGraph) CoordsToNode(row, col int) (node graph.Node) {
	id := gr.CoordsToID(row, col)
	if id == -1 {
		return nil
	} else {
		return Node(id)
	}
}

func (gr *TileGraph) Neighbors(node graph.Node) []graph.Node {
	id := node.ID()
	if !gr.NodeExists(node) {
		return nil
	}

	row, col := gr.IDToCoords(id)

	neighbors := []graph.Node{gr.CoordsToNode(row-1, col), gr.CoordsToNode(row+1, col), gr.CoordsToNode(row, col-1), gr.CoordsToNode(row, col+1)}
	realNeighbors := make([]graph.Node, 0, 4) // Will overallocate sometimes, but not by much. Not a big deal
	for _, neighbor := range neighbors {
		if neighbor != nil && gr.tiles[neighbor.ID()] == true {
			realNeighbors = append(realNeighbors, neighbor)
		}
	}

	return realNeighbors
}

func (gr *TileGraph) EdgeBetween(node, neighbor graph.Node) graph.Edge {
	if !gr.NodeExists(node) || !gr.NodeExists(neighbor) {
		return nil
	}

	r1, c1 := gr.IDToCoords(node.ID())
	r2, c2 := gr.IDToCoords(neighbor.ID())
	if (c1 == c2 && (r2 == r1+1 || r2 == r1-1)) || (r1 == r2 && (c2 == c1+1 || c2 == c1-1)) {
		return Edge{node, neighbor}
	}

	return nil
}

func (gr *TileGraph) NodeExists(node graph.Node) bool {
	id := node.ID()
	return id >= 0 && id < len(gr.tiles) && gr.tiles[id] == true
}

func (gr *TileGraph) Degree(node graph.Node) int {
	return len(gr.Neighbors(node)) * 2
}

func (gr *TileGraph) EdgeList() []graph.Edge {
	edges := make([]graph.Edge, 0)
	for id, passable := range gr.tiles {
		if !passable {
			continue
		}

		for _, succ := range gr.Neighbors(Node(id)) {
			edges = append(edges, Edge{Node(id), succ})
		}
	}

	return edges
}

func (gr *TileGraph) NodeList() []graph.Node {
	nodes := make([]graph.Node, 0)
	for id, passable := range gr.tiles {
		if !passable {
			continue
		}

		nodes = append(nodes, Node(id))
	}

	return nodes
}

func (gr *TileGraph) Cost(e graph.Edge) float64 {
	if edge := gr.EdgeBetween(e.Head(), e.Tail()); edge != nil {
		return 1.0
	}

	return math.Inf(1)
}
