// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dot_test

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

// dotGraph provides a shim for interaction between the DOT
// unmarshaler and a simple.WeightedUndirectedGraph.
type dotGraph struct {
	*simple.WeightedUndirectedGraph
}

func newDotGraph() *dotGraph {
	return &dotGraph{WeightedUndirectedGraph: simple.NewWeightedUndirectedGraph(0, 0)}
}

// NewEdge returns a DOT-aware edge.
func (g *dotGraph) NewEdge(from, to graph.Node) graph.Edge {
	e := g.WeightedUndirectedGraph.NewWeightedEdge(from, to, math.NaN()).(simple.WeightedEdge)
	return &weightedEdge{WeightedEdge: e}
}

// NewNode returns a DOT-aware node.
func (g *dotGraph) NewNode() graph.Node {
	return &node{Node: g.WeightedUndirectedGraph.NewNode()}
}

// SetEdge is a shim to allow the DOT unmarshaler to
// add weighted edges to a graph.
func (g *dotGraph) SetEdge(e graph.Edge) {
	g.WeightedUndirectedGraph.SetWeightedEdge(e.(*weightedEdge))
}

// weightedEdge is a DOT-aware weighted edge.
type weightedEdge struct {
	simple.WeightedEdge
}

// SetAttribute sets the weight of the receiver.
func (e *weightedEdge) SetAttribute(attr encoding.Attribute) error {
	if attr.Key != "weight" {
		return fmt.Errorf("unable to unmarshal node DOT attribute with key %q", attr.Key)
	}
	var err error
	e.W, err = strconv.ParseFloat(attr.Value, 64)
	return err
}

// node is a DOT-aware node.
type node struct {
	graph.Node
	dotID string
}

// SetDOTID sets the DOT ID of the node.
func (n *node) SetDOTID(id string) { n.dotID = id }

func (n *node) String() string { return n.dotID }

const ug = `
graph {
	a
	b
	c
	a--b ["weight"=0.5]
	a--c ["weight"=1]
}
`

func ExampleUnmarshal_weighted() {
	dst := newDotGraph()
	err := dot.Unmarshal([]byte(ug), dst)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range graph.EdgesOf(dst.Edges()) {
		fmt.Printf("%+v\n", e.(*weightedEdge).WeightedEdge)
	}

	// Unordered output:
	// {F:a T:b W:0.5}
	// {F:a T:c W:1}
}
