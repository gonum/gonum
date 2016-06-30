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
	zachary = []set{
		0:  linksTo(1, 2, 3, 4, 5, 6, 7, 8, 10, 11, 12, 13, 17, 19, 21, 31),
		1:  linksTo(2, 3, 7, 13, 17, 19, 21, 30),
		2:  linksTo(3, 7, 8, 9, 13, 27, 28, 32),
		3:  linksTo(7, 12, 13),
		4:  linksTo(6, 10),
		5:  linksTo(6, 10, 16),
		6:  linksTo(16),
		8:  linksTo(30, 32, 33),
		9:  linksTo(33),
		13: linksTo(33),
		14: linksTo(32, 33),
		15: linksTo(32, 33),
		18: linksTo(32, 33),
		19: linksTo(33),
		20: linksTo(32, 33),
		22: linksTo(32, 33),
		23: linksTo(25, 27, 29, 32, 33),
		24: linksTo(25, 27, 31),
		25: linksTo(31),
		26: linksTo(29, 33),
		27: linksTo(33),
		28: linksTo(31, 33),
		29: linksTo(32, 33),
		30: linksTo(32, 33),
		31: linksTo(32, 33),
		32: linksTo(33),
		33: nil,
	}

	// doi:10.1088/1742-5468/2008/10/P10008 figure 1
	blondel = []set{
		0:  linksTo(2, 3, 4, 5),
		1:  linksTo(2, 4, 7),
		2:  linksTo(4, 5, 6),
		3:  linksTo(7),
		4:  linksTo(10),
		5:  linksTo(7, 11),
		6:  linksTo(7, 11),
		8:  linksTo(9, 10, 11, 14, 15),
		9:  linksTo(12, 14),
		10: linksTo(11, 12, 13, 14),
		11: linksTo(13),
		15: nil,
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

var dupGraph = simple.NewUndirectedGraph(0, 0)

func init() {
	err := gen.Duplication(dupGraph, 1000, 0.8, 0.1, 0.5, rand.New(rand.NewSource(1)))
	if err != nil {
		panic(err)
	}
}
