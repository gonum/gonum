// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package community

import (
	"math/rand"

	"github.com/gonum/graph"
	"github.com/gonum/graph/graphs/gen"
	"github.com/gonum/graph/simple"
)

// set is an integer set.
type set map[int]struct{}

func linksTo(i ...int) set {
	if len(i) == 0 {
		return nil
	}
	s := make(set)
	for _, v := range i {
		s[v] = struct{}{}
	}
	return s
}

var (
	unconnected = []set{ /* Nodes 0-4 are implicit .*/ 5: nil}

	smallDumbell = []set{
		0: linksTo(1, 2),
		1: linksTo(2),
		2: linksTo(3),
		3: linksTo(4, 5),
		4: linksTo(5),
		5: nil,
	}

	// W. W. Zachary, An information flow model for conflict and fission in small groups,
	// Journal of Anthropological Research 33, 452-473 (1977).
	//
	// The edge list here is constructed such that all link descriptions
	// head from a node with lower Page Rank to a node with higher Page
	// Rank. This has no impact on undirected tests, but allows a sensible
	// view for directed tests.
	zachary = []set{
		0:  nil,                     // rank=0.097
		1:  linksTo(0, 2),           // rank=0.05288
		2:  linksTo(0, 32),          // rank=0.05708
		3:  linksTo(0, 1, 2),        // rank=0.03586
		4:  linksTo(0, 6, 10),       // rank=0.02198
		5:  linksTo(0, 6),           // rank=0.02911
		6:  linksTo(0, 5),           // rank=0.02911
		7:  linksTo(0, 1, 2, 3),     // rank=0.02449
		8:  linksTo(0, 2, 32, 33),   // rank=0.02977
		9:  linksTo(2, 33),          // rank=0.01431
		10: linksTo(0, 5),           // rank=0.02198
		11: linksTo(0),              // rank=0.009565
		12: linksTo(0, 3),           // rank=0.01464
		13: linksTo(0, 1, 2, 3, 33), // rank=0.02954
		14: linksTo(32, 33),         // rank=0.01454
		15: linksTo(32, 33),         // rank=0.01454
		16: linksTo(5, 6),           // rank=0.01678
		17: linksTo(0, 1),           // rank=0.01456
		18: linksTo(32, 33),         // rank=0.01454
		19: linksTo(0, 1, 33),       // rank=0.0196
		20: linksTo(32, 33),         // rank=0.01454
		21: linksTo(0, 1),           // rank=0.01456
		22: linksTo(32, 33),         // rank=0.01454
		23: linksTo(32, 33),         // rank=0.03152
		24: linksTo(27, 31),         // rank=0.02108
		25: linksTo(23, 24, 31),     // rank=0.02101
		26: linksTo(29, 33),         // rank=0.01504
		27: linksTo(2, 23, 33),      // rank=0.02564
		28: linksTo(2, 31, 33),      // rank=0.01957
		29: linksTo(23, 32, 33),     // rank=0.02629
		30: linksTo(1, 8, 32, 33),   // rank=0.02459
		31: linksTo(0, 32, 33),      // rank=0.03716
		32: linksTo(33),             // rank=0.07169
		33: nil,                     // rank=0.1009
	}

	// doi:10.1088/1742-5468/2008/10/P10008 figure 1
	//
	// The edge list here is constructed such that all link descriptions
	// head from a node with lower Page Rank to a node with higher Page
	// Rank. This has no impact on undirected tests, but allows a sensible
	// view for directed tests.
	blondel = []set{
		0:  linksTo(2),           // rank=0.06858
		1:  linksTo(2, 4, 7),     // rank=0.05264
		2:  nil,                  // rank=0.08249
		3:  linksTo(0, 7),        // rank=0.03884
		4:  linksTo(0, 2, 10),    // rank=0.06754
		5:  linksTo(0, 2, 7, 11), // rank=0.06738
		6:  linksTo(2, 7, 11),    // rank=0.0528
		7:  nil,                  // rank=0.07008
		8:  linksTo(10),          // rank=0.09226
		9:  linksTo(8),           // rank=0.05821
		10: nil,                  // rank=0.1035
		11: linksTo(8, 10),       // rank=0.08538
		12: linksTo(9, 10),       // rank=0.04052
		13: linksTo(10, 11),      // rank=0.03855
		14: linksTo(8, 9, 10),    // rank=0.05621
		15: linksTo(8),           // rank=0.02506
	}
)

type structure struct {
	resolution  float64
	memberships []set
	want, tol   float64
}

type level struct {
	q           float64
	communities [][]graph.Node
}

type moveStructures struct {
	memberships []set
	targetNodes []graph.Node

	resolution float64
	tol        float64
}

func reverse(f []float64) {
	for i, j := 0, len(f)-1; i < j; i, j = i+1, j-1 {
		f[i], f[j] = f[j], f[i]
	}
}

var (
	dupGraph         = simple.NewUndirectedGraph(0, 0)
	dupGraphDirected = simple.NewDirectedGraph(0, 0)
)

func init() {
	err := gen.Duplication(dupGraph, 1000, 0.8, 0.1, 0.5, rand.New(rand.NewSource(1)))
	if err != nil {
		panic(err)
	}

	// Construct a directed graph from dupGraph
	// such that every edge dupGraph is replaced
	// with an edge that flows from the low node
	// ID to the high node ID.
	for _, e := range dupGraph.Edges() {
		if e.To().ID() < e.From().ID() {
			se := e.(simple.Edge)
			se.F, se.T = se.T, se.F
			e = se
		}
		dupGraphDirected.SetEdge(e)
	}
}
