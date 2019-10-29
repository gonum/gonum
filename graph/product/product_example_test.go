// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package product_test

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/product"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

// atom is a graph.Node representing an atom in a molecule.
type atom struct {
	name string // name is the name of the atom.
	pos  int    // pos is the position number of the atom.
	id   int64
}

// ID satisfies the graph.Node interface.
func (n atom) ID() int64 { return n.id }

func ExampleModular_subgraphIsomorphism() {
	// The modular product can be used to find subgraph isomorphisms.
	// See https://doi.org/10.1016/0020-0190(76)90049-1 and for a
	// theoretical perspective, https://doi.org/10.1145/990524.990529.

	// We can find the common structure between two organic molecules.
	// For example the purines adenine and guanine from nucleic acids.

	// Make a graph for adenine.
	adenine := simple.NewUndirectedGraph()
	for _, bond := range []simple.Edge{
		// Purine nucleus.
		{F: atom{name: "N", pos: 1, id: 0}, T: atom{name: "C", pos: 2, id: 1}},
		{F: atom{name: "N", pos: 1, id: 0}, T: atom{name: "C", pos: 6, id: 5}},
		{F: atom{name: "C", pos: 2, id: 1}, T: atom{name: "N", pos: 3, id: 2}},
		{F: atom{name: "N", pos: 3, id: 2}, T: atom{name: "C", pos: 4, id: 3}},
		{F: atom{name: "C", pos: 4, id: 3}, T: atom{name: "C", pos: 5, id: 4}},
		{F: atom{name: "C", pos: 4, id: 3}, T: atom{name: "N", pos: 9, id: 8}},
		{F: atom{name: "C", pos: 5, id: 4}, T: atom{name: "C", pos: 6, id: 5}},
		{F: atom{name: "C", pos: 5, id: 4}, T: atom{name: "N", pos: 7, id: 6}},
		{F: atom{name: "N", pos: 7, id: 6}, T: atom{name: "C", pos: 8, id: 7}},
		{F: atom{name: "C", pos: 8, id: 7}, T: atom{name: "N", pos: 9, id: 8}},

		// Amino modification in adenine.
		//
		// Note that the position number of the N is non-standard.
		{F: atom{name: "C", pos: 6, id: 5}, T: atom{name: "N", pos: 10, id: 9}},
	} {
		adenine.SetEdge(bond)
	}

	// Make a graph for guanine.
	//
	// Note that node IDs here have no intersection with
	// the adenine graph to show that they are not being
	// used to map between the graphs.
	guanine := simple.NewUndirectedGraph()
	for _, bond := range []simple.Edge{
		// Purine nucleus.
		{F: atom{name: "N", pos: 1, id: 10}, T: atom{name: "C", pos: 2, id: 11}},
		{F: atom{name: "N", pos: 1, id: 10}, T: atom{name: "C", pos: 6, id: 15}},
		{F: atom{name: "C", pos: 2, id: 11}, T: atom{name: "N", pos: 3, id: 12}},
		{F: atom{name: "N", pos: 3, id: 12}, T: atom{name: "C", pos: 4, id: 13}},
		{F: atom{name: "C", pos: 4, id: 13}, T: atom{name: "C", pos: 5, id: 14}},
		{F: atom{name: "C", pos: 4, id: 13}, T: atom{name: "N", pos: 9, id: 18}},
		{F: atom{name: "C", pos: 5, id: 14}, T: atom{name: "C", pos: 6, id: 15}},
		{F: atom{name: "C", pos: 5, id: 14}, T: atom{name: "N", pos: 7, id: 16}},
		{F: atom{name: "N", pos: 7, id: 16}, T: atom{name: "C", pos: 8, id: 17}},
		{F: atom{name: "C", pos: 8, id: 17}, T: atom{name: "N", pos: 9, id: 18}},

		// Amino and keto modifications in guanine.
		//
		// Note that the position number of the N and O is non-standard.
		{F: atom{name: "C", pos: 2, id: 11}, T: atom{name: "N", pos: 11, id: 19}},
		{F: atom{name: "C", pos: 6, id: 15}, T: atom{name: "O", pos: 10, id: 20}},
	} {
		guanine.SetEdge(bond)
	}

	// Produce the modular product of the two graphs.
	p := simple.NewUndirectedGraph()
	product.Modular(p, adenine, guanine)

	// Find the maximal cliques in the modular product.
	mc := topo.BronKerbosch(p)

	// Report the largest.
	sort.Sort(byLength(mc))
	max := len(mc[0])
	w := tabwriter.NewWriter(os.Stdout, 5, 0, 0, ' ', tabwriter.AlignRight)
	fmt.Println("  Adenine   Guanine")
	fmt.Fprintln(w, "Atom\tPos\tAtom\tPos\t")
	for _, c := range mc {
		if len(c) < max {
			break
		}
		for _, p := range c {
			// Extract the mapping between the
			// inputs from the product.
			p := p.(product.Node)
			adenine := p.A.(atom)
			guanine := p.B.(atom)
			fmt.Fprintf(w, "%s\t%d\t%s\t%d\t\n", adenine.name, adenine.pos, guanine.name, guanine.pos)
		}
	}
	w.Flush()

	// Unordered output:
	//   Adenine   Guanine
	//  Atom  Pos Atom  Pos
	//     N    3    N    3
	//     N    7    N    7
	//     N   10    O   10
	//     C    6    C    6
	//     C    2    C    2
	//     C    8    C    8
	//     C    5    C    5
	//     N    9    N    9
	//     N    1    N    1
	//     C    4    C    4
}

// byLength implements the sort.Interface, sorting the slices
// descending by length.
type byLength [][]graph.Node

func (n byLength) Len() int           { return len(n) }
func (n byLength) Less(i, j int) bool { return len(n[i]) > len(n[j]) }
func (n byLength) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
