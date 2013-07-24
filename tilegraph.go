package discrete

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
	if loc > len(graph.tiles) || row < 0 || col < 0 {
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
				outString += " " // White square
			}
		}

		outString += "\n"
	}

	return outString[:len(outString)-1] // Kill final newline
}
