package concrete

import (
	"errors"
	"math"
	"strings"

	gr "github.com/gonum/graph"
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

func (graph *TileGraph) PathString(path []gr.Node) string {
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

func (graph *TileGraph) CoordsToNode(row, col int) (node gr.Node) {
	id := graph.CoordsToID(row, col)
	if id == -1 {
		return nil
	} else {
		return Node(id)
	}
}

func (graph *TileGraph) Neighbors(node gr.Node) []gr.Node {
	id := node.ID()
	if !graph.NodeExists(node) {
		return nil
	}

	row, col := graph.IDToCoords(id)

	neighbors := []gr.Node{graph.CoordsToNode(row-1, col), graph.CoordsToNode(row+1, col), graph.CoordsToNode(row, col-1), graph.CoordsToNode(row, col+1)}
	realNeighbors := make([]gr.Node, 0, 4) // Will overallocate sometimes, but not by much. Not a big deal
	for _, neighbor := range neighbors {
		if neighbor != nil && graph.tiles[neighbor.ID()] == true {
			realNeighbors = append(realNeighbors, neighbor)
		}
	}

	return realNeighbors
}

func (graph *TileGraph) EdgeBetween(node, neighbor gr.Node) gr.Edge {
	if !graph.NodeExists(node) || !graph.NodeExists(neighbor) {
		return nil
	}

	r1, c1 := graph.IDToCoords(node.ID())
	r2, c2 := graph.IDToCoords(neighbor.ID())
	if (c1 == c2 && (r2 == r1+1 || r2 == r1-1)) || (r1 == r2 && (c2 == c1+1 || c2 == c1-1)) {
		return Edge{node, neighbor}
	}

	return nil
}

func (graph *TileGraph) NodeExists(node gr.Node) bool {
	id := node.ID()
	return id >= 0 && id < len(graph.tiles) && graph.tiles[id] == true
}

func (graph *TileGraph) Degree(node gr.Node) int {
	return len(graph.Neighbors(node)) * 2
}

func (graph *TileGraph) EdgeList() []gr.Edge {
	edges := make([]gr.Edge, 0)
	for id, passable := range graph.tiles {
		if !passable {
			continue
		}

		for _, succ := range graph.Neighbors(Node(id)) {
			edges = append(edges, Edge{Node(id), succ})
		}
	}

	return edges
}

func (graph *TileGraph) NodeList() []gr.Node {
	nodes := make([]gr.Node, 0)
	for id, passable := range graph.tiles {
		if !passable {
			continue
		}

		nodes = append(nodes, Node(id))
	}

	return nodes
}

func (graph *TileGraph) Cost(e gr.Edge) float64 {
	if edge := graph.EdgeBetween(e.Head(), e.Tail()); edge != nil {
		return 1.0
	}

	return math.Inf(1)
}
