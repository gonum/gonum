// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package product_test

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/product"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

// person is a graph.Node representing a person.
type person struct {
	name string // name is the name of the person.
	id   int64
}

// ID satisfies the graph.Node interface.
func (n person) ID() int64 { return n.id }

func ExampleModularExt_subgraphIsomorphism() {
	// Extended attributes of the graph can be used to refine
	// subgraph isomorphism identification. By filtering edge
	// agreement by weight we can identify social network
	// motifs within a larger graph.
	//
	// This example extracts sources of conflict from the
	// relationships of Julius Caesar, Mark Antony and
	// Cleopatra.

	// Make a graph describing people's relationships.
	//
	// Edge weight indicates love/animosity.
	people := simple.NewDirectedGraph()
	for _, relationship := range []simple.WeightedEdge{
		{F: person{name: "Julius Caesar", id: 0}, T: person{name: "Cleopatra", id: 1}, W: 1},
		{F: person{name: "Cleopatra", id: 1}, T: person{name: "Julius Caesar", id: 0}, W: 1},
		{F: person{name: "Julius Caesar", id: 0}, T: person{name: "Cornelia", id: 3}, W: 1},
		{F: person{name: "Cornelia", id: 3}, T: person{name: "Julius Caesar", id: 0}, W: 1},
		{F: person{name: "Mark Antony", id: 2}, T: person{name: "Cleopatra", id: 1}, W: 1},
		{F: person{name: "Cleopatra", id: 1}, T: person{name: "Mark Antony", id: 2}, W: 1},
		{F: person{name: "Fulvia", id: 4}, T: person{name: "Mark Antony", id: 2}, W: 1},
		{F: person{name: "Fulvia", id: 4}, T: person{name: "Cleopatra", id: 1}, W: -1},
		{F: person{name: "Octavia", id: 5}, T: person{name: "Mark Antony", id: 2}, W: 1},
		{F: person{name: "Octavia", id: 5}, T: person{name: "Cleopatra", id: 1}, W: -1},
	} {
		people.SetEdge(relationship)
	}

	// Make a graph for the query pattern: a love triangle.
	pattern := simple.NewDirectedGraph()
	for _, relationsip := range []simple.WeightedEdge{
		{F: person{name: "A", id: -1}, T: person{name: "B", id: -2}, W: 1},
		{F: person{name: "B", id: -2}, T: person{name: "A", id: -1}, W: 1},
		{F: person{name: "C", id: -3}, T: person{name: "A", id: -1}, W: -1},
		{F: person{name: "C", id: -3}, T: person{name: "B", id: -2}, W: 1},
	} {
		pattern.SetEdge(relationsip)
	}

	// Produce the modular product of the two graphs.
	p := simple.NewDirectedGraph()
	product.ModularExt(p, people, pattern, func(a, b graph.Edge) bool {
		return a.(simple.WeightedEdge).Weight() == b.(simple.WeightedEdge).Weight()
	})

	// Find the maximal cliques in the undirected induction
	// of the modular product.
	mc := topo.BronKerbosch(undirected{p})

	// Report the cliques that are identical in order to the pattern.
	fmt.Println("Person — Relationship position:")
	for _, c := range mc {
		if len(c) != pattern.Nodes().Len() {
			continue
		}
		for _, p := range c {
			// Extract the mapping between the
			// inputs from the product.
			p := p.(product.Node)
			people := p.A.(person)
			pattern := p.B.(person)
			fmt.Printf(" %s — %s\n", people.name, pattern.name)
		}
		fmt.Println()
	}

	// Unordered output:
	// Person — Relationship position:
	//  Cleopatra — A
	//  Mark Antony — B
	//  Octavia — C
	//
	//  Cleopatra — A
	//  Mark Antony — B
	//  Fulvia — C
}

// undirected converts a directed graph to an undirected graph
// with edges between nodes only where directed edges exist in
// both directions in the original graph.
type undirected struct {
	graph.Directed
}

func (g undirected) From(uid int64) graph.Nodes {
	nodes := graph.NodesOf(g.Directed.From(uid))
	for i := 0; i < len(nodes); {
		if g.Directed.Edge(nodes[i].ID(), uid) != nil {
			i++
		} else {
			nodes[i], nodes = nodes[len(nodes)-1], nodes[:len(nodes)-1]
		}
	}
	return iterator.NewOrderedNodes(nodes)
}
func (g undirected) Edge(xid, yid int64) graph.Edge {
	e := g.Directed.Edge(xid, yid)
	if e != nil && g.Directed.Edge(yid, xid) != nil {
		return e
	}
	return nil
}
func (g undirected) EdgeBetween(xid, yid int64) graph.Edge {
	return g.Edge(xid, yid)
}
