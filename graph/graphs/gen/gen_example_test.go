// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gen_test

import (
	"fmt"
	"log"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/graphs/gen"
	"gonum.org/v1/gonum/graph/simple"
)

func ExampleStar_undirectedRange() {
	dst := simple.NewUndirectedGraph()
	gen.Star(dst, 0, gen.IDRange{First: 1, Last: 6})
	b, err := dot.Marshal(dst, "star", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// strict graph star {
	// 	// Node definitions.
	// 	0;
	// 	1;
	// 	2;
	// 	3;
	// 	4;
	// 	5;
	// 	6;
	//
	// 	// Edge definitions.
	// 	0 -- 1;
	// 	0 -- 2;
	// 	0 -- 3;
	// 	0 -- 4;
	// 	0 -- 5;
	// 	0 -- 6;
	// }
}

func ExampleWheel_directedRange() {
	dst := simple.NewDirectedGraph()
	gen.Wheel(dst, 0, gen.IDRange{First: 1, Last: 6})
	b, err := dot.Marshal(dst, "wheel", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// strict digraph wheel {
	// 	// Node definitions.
	// 	0;
	// 	1;
	// 	2;
	// 	3;
	// 	4;
	// 	5;
	// 	6;
	//
	// 	// Edge definitions.
	// 	0 -> 1;
	// 	0 -> 2;
	// 	0 -> 3;
	// 	0 -> 4;
	// 	0 -> 5;
	// 	0 -> 6;
	// 	1 -> 2;
	// 	2 -> 3;
	// 	3 -> 4;
	// 	4 -> 5;
	// 	5 -> 6;
	// 	6 -> 1;
	// }
}

func ExamplePath_directedSet() {
	dst := simple.NewDirectedGraph()
	gen.Path(dst, gen.IDSet{2, 4, 5, 9})
	b, err := dot.Marshal(dst, "path", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// strict digraph path {
	// 	// Node definitions.
	// 	2;
	// 	4;
	// 	5;
	// 	9;
	//
	// 	// Edge definitions.
	// 	2 -> 4;
	// 	4 -> 5;
	// 	5 -> 9;
	// }
}

func ExampleComplete_directedSet() {
	dst := simple.NewDirectedGraph()
	gen.Complete(dst, gen.IDSet{2, 4, 5, 9})
	b, err := dot.Marshal(dst, "complete", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// strict digraph complete {
	// 	// Node definitions.
	// 	2;
	// 	4;
	// 	5;
	// 	9;
	//
	// 	// Edge definitions.
	// 	2 -> 4;
	// 	2 -> 5;
	// 	2 -> 9;
	// 	4 -> 5;
	// 	4 -> 9;
	// 	5 -> 9;
	// }
}

// Bidirected allows bidirectional directed graph construction.
type Bidirected struct {
	*simple.DirectedGraph
}

func (g Bidirected) SetEdge(e graph.Edge) {
	g.DirectedGraph.SetEdge(e)
	g.DirectedGraph.SetEdge(e.ReversedEdge())
}

func ExampleComplete_biDirectedSet() {
	dst := simple.NewDirectedGraph()
	gen.Complete(Bidirected{dst}, gen.IDSet{2, 4, 5, 9})
	b, err := dot.Marshal(dst, "complete", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// strict digraph complete {
	// 	// Node definitions.
	// 	2;
	// 	4;
	// 	5;
	// 	9;
	//
	// 	// Edge definitions.
	// 	2 -> 4;
	// 	2 -> 5;
	// 	2 -> 9;
	// 	4 -> 2;
	// 	4 -> 5;
	// 	4 -> 9;
	// 	5 -> 2;
	// 	5 -> 4;
	// 	5 -> 9;
	// 	9 -> 2;
	// 	9 -> 4;
	// 	9 -> 5;
	// }
}

func ExampleComplete_undirectedSet() {
	dst := simple.NewUndirectedGraph()
	gen.Complete(dst, gen.IDSet{2, 4, 5, 9})
	b, err := dot.Marshal(dst, "complete", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// strict graph complete {
	// 	// Node definitions.
	// 	2;
	// 	4;
	// 	5;
	// 	9;
	//
	// 	// Edge definitions.
	// 	2 -- 4;
	// 	2 -- 5;
	// 	2 -- 9;
	// 	4 -- 5;
	// 	4 -- 9;
	// 	5 -- 9;
	// }
}

func ExampleTree_undirectedRange() {
	dst := simple.NewUndirectedGraph()
	gen.Tree(dst, 2, gen.IDRange{First: 0, Last: 14})
	b, err := dot.Marshal(dst, "full_binary_tree_undirected", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// strict graph full_binary_tree_undirected {
	// 	// Node definitions.
	// 	0;
	// 	1;
	// 	2;
	// 	3;
	// 	4;
	// 	5;
	// 	6;
	// 	7;
	// 	8;
	// 	9;
	// 	10;
	// 	11;
	// 	12;
	// 	13;
	// 	14;
	//
	// 	// Edge definitions.
	// 	0 -- 1;
	// 	0 -- 2;
	// 	1 -- 3;
	// 	1 -- 4;
	// 	2 -- 5;
	// 	2 -- 6;
	// 	3 -- 7;
	// 	3 -- 8;
	// 	4 -- 9;
	// 	4 -- 10;
	// 	5 -- 11;
	// 	5 -- 12;
	// 	6 -- 13;
	// 	6 -- 14;
	// }
}
