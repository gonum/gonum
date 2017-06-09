// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dot

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

func TestRoundTrip(t *testing.T) {
	golden := []struct {
		want     string
		directed bool
	}{
		{
			want:     directed,
			directed: true,
		},
		{
			want:     undirected,
			directed: false,
		},
	}
	for i, g := range golden {
		var dst Builder
		if g.directed {
			dst = newDotDirectedGraph()
		} else {
			dst = newDotUndirectedGraph()
		}
		data := []byte(g.want)
		if err := Unmarshal(data, dst); err != nil {
			t.Errorf("i=%d: unable to unmarshal DOT graph; %v", i, err)
			continue
		}
		buf, err := Marshal(dst, "", "", "\t", false)
		if err != nil {
			t.Errorf("i=%d: unable to marshal graph; %v", i, dst)
			continue
		}
		got := string(buf)
		if got != g.want {
			t.Errorf("i=%d: graph content mismatch; want:\n%s\n\ngot:\n%s", i, g.want, got)
			continue
		}
	}
}

const directed = `digraph {
	graph [
		outputorder=edgesfirst
	];
	node [
		shape=circle
		style=filled
	];
	edge [
		penwidth=5
		color=gray
	];

	// Node definitions.
	A [label="foo 2"];
	B [label="bar 2"];

	// Edge definitions.
	A -> B [label="baz 2"];
}`

const undirected = `graph {
	graph [
		outputorder=edgesfirst
	];
	node [
		shape=circle
		style=filled
	];
	edge [
		penwidth=5
		color=gray
	];

	// Node definitions.
	A [label="foo 2"];
	B [label="bar 2"];

	// Edge definitions.
	A -- B [label="baz 2"];
}`

// Below follows a minimal implementation of a graph capable of validating the
// round-trip encoding and decoding of DOT graphs with nodes and edges
// containing DOT attributes.

// dotDirectedGraph extends simple.DirectedGraph to add NewNode and NewEdge
// methods for creating user-defined nodes and edges.
//
// dotDirectedGraph implements the dot.Builder interface.
type dotDirectedGraph struct {
	*simple.DirectedGraph
	graph, node, edge attributes
}

// newDotDirectedGraph returns a new directed capable of creating user-defined
// nodes and edges.
func newDotDirectedGraph() *dotDirectedGraph {
	return &dotDirectedGraph{DirectedGraph: simple.NewDirectedGraph(0, 0)}
}

// NewNode adds a new node with a unique node ID to the graph.
func (g *dotDirectedGraph) NewNode() graph.Node {
	n := &dotNode{Node: simple.Node(g.NewNodeID())}
	g.AddNode(n)
	return n
}

// NewEdge adds a new edge from the source to the destination node to the graph,
// or returns the existing edge if already present.
func (g *dotDirectedGraph) NewEdge(from, to graph.Node) graph.Edge {
	if e := g.Edge(from, to); e != nil {
		return e
	}
	e := &dotEdge{Edge: simple.Edge{F: from, T: to}}
	g.SetEdge(e)
	return e
}

// DOTAttributers implements the dot.Attributers interface.
func (g *dotDirectedGraph) DOTAttributers() (graph, node, edge Attributer) {
	return g.graph, g.node, g.edge
}

// DOTUnmarshalerAttrs implements the dot.UnmarshalerAttrs interface.
func (g *dotDirectedGraph) DOTUnmarshalerAttrs() (graph, node, edge UnmarshalerAttr) {
	return &g.graph, &g.node, &g.edge
}

// dotUndirectedGraph extends simple.UndirectedGraph to add NewNode and NewEdge
// methods for creating user-defined nodes and edges.
//
// dotUndirectedGraph implements the dot.Builder interface.
type dotUndirectedGraph struct {
	*simple.UndirectedGraph
	graph, node, edge attributes
}

// newDotUndirectedGraph returns a new undirected capable of creating user-
// defined nodes and edges.
func newDotUndirectedGraph() *dotUndirectedGraph {
	return &dotUndirectedGraph{UndirectedGraph: simple.NewUndirectedGraph(0, 0)}
}

// NewNode adds a new node with a unique node ID to the graph.
func (g *dotUndirectedGraph) NewNode() graph.Node {
	n := &dotNode{Node: simple.Node(g.NewNodeID())}
	g.AddNode(n)
	return n
}

// NewEdge adds a new edge from the source to the destination node to the graph,
// or returns the existing edge if already present.
func (g *dotUndirectedGraph) NewEdge(from, to graph.Node) graph.Edge {
	if e := g.Edge(from, to); e != nil {
		return e
	}
	e := &dotEdge{Edge: simple.Edge{F: from, T: to}}
	g.SetEdge(e)
	return e
}

// DOTAttributers implements the dot.Attributers interface.
func (g *dotUndirectedGraph) DOTAttributers() (graph, node, edge Attributer) {
	return g.graph, g.node, g.edge
}

// DOTUnmarshalerAttrs implements the dot.UnmarshalerAttrs interface.
func (g *dotUndirectedGraph) DOTUnmarshalerAttrs() (graph, node, edge UnmarshalerAttr) {
	return &g.graph, &g.node, &g.edge
}

// dotNode extends simple.Node with a label field to test round-trip encoding
// and decoding of node DOT label attributes.
type dotNode struct {
	simple.Node
	dotID string
	// Node label.
	Label string
}

// DOTID returns the node's DOT ID.
func (n *dotNode) DOTID() string {
	return n.dotID
}

// UnmarshalDOTID decodes a DOT ID.
func (n *dotNode) UnmarshalDOTID(id string) {
	n.dotID = id
}

// UnmarshalDOTAttr decodes a single DOT attribute.
func (n *dotNode) UnmarshalDOTAttr(attr Attribute) error {
	if attr.Key != "label" {
		return fmt.Errorf("unable to unmarshal node DOT attribute with key %q", attr.Key)
	}
	n.Label = attr.Value
	return nil
}

// DOTAttributes returns the DOT attributes of the node.
func (n *dotNode) DOTAttributes() []Attribute {
	if len(n.Label) == 0 {
		return nil
	}
	attr := Attribute{
		Key:   "label",
		Value: n.Label,
	}
	return []Attribute{attr}
}

// dotEdge extends simple.Edge with a label field to test round-trip encoding and
// decoding of edge DOT label attributes.
type dotEdge struct {
	simple.Edge
	// Edge label.
	Label string
}

// UnmarshalDOTAttr decodes a single DOT attribute.
func (e *dotEdge) UnmarshalDOTAttr(attr Attribute) error {
	if attr.Key != "label" {
		return fmt.Errorf("unable to unmarshal node DOT attribute with key %q", attr.Key)
	}
	e.Label = attr.Value
	return nil
}

// DOTAttributes returns the DOT attributes of the edge.
func (e *dotEdge) DOTAttributes() []Attribute {
	if len(e.Label) == 0 {
		return nil
	}
	attr := Attribute{
		Key:   "label",
		Value: e.Label,
	}
	return []Attribute{attr}
}

// attributes is a helper for global attributes.
type attributes []Attribute

func (a attributes) DOTAttributes() []Attribute {
	return []Attribute(a)
}
func (a *attributes) UnmarshalDOTAttr(attr Attribute) error {
	*a = append(*a, attr)
	return nil
}
