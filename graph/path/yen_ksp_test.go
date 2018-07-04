// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"fmt"
	"math"
	"testing"

	"gonum.org/v1/gonum/graph/simple"
)


func TestYenKSP(t *testing.T) {
	g := simple.NewWeightedDirectedGraph(0, math.Inf(1))

	edges := []simple.WeightedEdge{
		{F: simple.Node(0), T: simple.Node(1), W: 3},
		{F: simple.Node(0), T: simple.Node(2), W: 2},
		{F: simple.Node(2), T: simple.Node(1), W: 1},
		{F: simple.Node(1), T: simple.Node(3), W: 4},
		{F: simple.Node(2), T: simple.Node(3), W: 2},
		{F: simple.Node(2), T: simple.Node(4), W: 3},
		{F: simple.Node(3), T: simple.Node(4), W: 2},
		{F: simple.Node(3), T: simple.Node(5), W: 1},
		{F: simple.Node(4), T: simple.Node(5), W: 2},
	}

	for _, edge := range edges {
		g.SetWeightedEdge(edge)
	}

	shortests := YenKSP(simple.Node(0), simple.Node(5), g, 3)
	
	fmt.Printf("%v", shortests)
}
