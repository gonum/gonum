// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coloring_test

import (
	"fmt"
	"log"

	"gonum.org/v1/gonum/graph/coloring"
	"gonum.org/v1/gonum/graph/graphs/gen"
	"gonum.org/v1/gonum/graph/simple"
)

// A hard sudoku problem graded at a level of difficulty, "not fun".
// https://dingo.sbs.arizona.edu/~sandiway/sudoku/examples.html
var grid = [9][9]int{
	{0, 2, 0 /**/, 0, 0, 0 /**/, 0, 0, 0},
	{0, 0, 0 /**/, 6, 0, 0 /**/, 0, 0, 3},
	{0, 7, 4 /**/, 0, 8, 0 /**/, 0, 0, 0},

	{0, 0, 0 /**/, 0, 0, 3 /**/, 0, 0, 2},
	{0, 8, 0 /**/, 0, 4, 0 /**/, 0, 1, 0},
	{6, 0, 0 /**/, 5, 0, 0 /**/, 0, 0, 0},

	{0, 0, 0 /**/, 0, 1, 0 /**/, 7, 8, 0},
	{5, 0, 0 /**/, 0, 0, 9 /**/, 0, 0, 0},
	{0, 0, 0 /**/, 0, 0, 0 /**/, 0, 4, 0},
}

func Example_sudoku() {
	g := simple.NewUndirectedGraph()

	// Build the sudoku board constraints.
	for i := 0; i < 9; i++ {
		gen.Complete(g, row(i))
		gen.Complete(g, col(i))
	}
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			gen.Complete(g, block{r, c})
		}
	}

	// Add constraints for the digits.
	gen.Complete(g, gen.IDRange{First: -9, Last: -1})

	// Mark constraints onto the graph.
	for r, row := range &grid {
		for c, val := range &row {
			if val == 0 {
				continue
			}
			for i := 1; i <= 9; i++ {
				if i != val {
					g.SetEdge(simple.Edge{F: simple.Node(-i), T: simple.Node(id(r, c))})
				}
			}
		}
	}

	k, colors, err := coloring.DsaturExact(nil, g)
	if err != nil {
		log.Fatal(err)
	}
	if k != 9 {
		log.Fatalln("could not solve problem", k)
	}
	sets := coloring.Sets(colors)
	for r := 0; r < 9; r++ {
		if r != 0 && r%3 == 0 {
			fmt.Println()
		}
		for c := 0; c < 9; c++ {
			if c != 0 {
				fmt.Print(" ")
				if c%3 == 0 {
					fmt.Print(" ")
				}
			}
			got := -int(sets[colors[id(r, c)]][0])
			if want := grid[r][c]; want != 0 && got != want {
				log.Fatalf("mismatch at row=%d col=%d: %d != %d", r, c, got, want)
			}
			fmt.Print(got)
		}
		fmt.Println()
	}

	// Output:
	//
	// 1 2 6  4 3 7  9 5 8
	// 8 9 5  6 2 1  4 7 3
	// 3 7 4  9 8 5  1 2 6
	//
	// 4 5 7  1 9 3  8 6 2
	// 9 8 3  2 4 6  5 1 7
	// 6 1 2  5 7 8  3 9 4
	//
	// 2 6 9  3 1 4  7 8 5
	// 5 4 8  7 6 9  2 3 1
	// 7 3 1  8 5 2  6 4 9
}

// row is a gen.IDer that enumerates the IDs of graph
// nodes representing a row of cells of a sudoku board.
type row int

func (r row) Len() int       { return 9 }
func (r row) ID(i int) int64 { return id(int(r), i) }

// col is a gen.IDer that enumerates the IDs of graph
// nodes representing a column of cells of a sudoku board.
type col int

func (c col) Len() int       { return 9 }
func (c col) ID(i int) int64 { return id(i, int(c)) }

// block is a gen.IDer that enumerates the IDs of graph
// nodes representing a 3×3 block of cells of a sudoku board.
type block struct {
	r, c int
}

func (b block) Len() int { return 9 }
func (b block) ID(i int) int64 {
	return id(b.r*3, b.c*3) + int64(i%3) + int64(i/3)*9
}

// id returns the graph node ID of a cell in a sudoku board.
func id(row, col int) int64 {
	return int64(row*9 + col)
}
