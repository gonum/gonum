// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package concrete

import (
	"math"

	"github.com/gonum/graph"
)

// A simple int alias.
type Node int

func (n Node) ID() int {
	return int(n)
}

// Just a collection of two nodes
type Edge struct {
	H, T graph.Node
}

func (edge Edge) Head() graph.Node {
	return edge.H
}

func (edge Edge) Tail() graph.Node {
	return edge.T
}

type WeightedEdge struct {
	graph.Edge
	Cost float64
}

// A GonumGraph is a very generalized graph that can handle an arbitrary number of vertices and
// edges -- as well as act as either directed or undirected.
//
// Internally, it uses a map of successors AND predecessors, to speed up some operations (such as
// getting all successors/predecessors). It also speeds up things like adding edges (assuming both
// edges exist).
//
// However, its generality is also its weakness (and partially a flaw in needing to satisfy
// MutableGraph). For most purposes, creating your own graph is probably better. For instance,
// see TileGraph for an example of an immutable 2D grid of tiles that also implements the Graph
// interface, but would be more suitable if all you needed was a simple undirected 2D grid.
type Graph struct {
	neighbors map[int]map[int]WeightedEdge
	nodeMap   map[int]graph.Node
	maxID     int
}

func NewGraph() *Graph {
	return &Graph{
		neighbors: make(map[int]map[int]WeightedEdge),
		nodeMap:   make(map[int]graph.Node),
		maxID:     0,
	}
}

/* Mutable implementation */

func (g *Graph) NewNode() graph.Node {
	if g.maxID != maxInt {
		g.maxID++
		g.AddNode(Node(g.maxID))
		return Node(g.maxID)
	}

	// I cannot foresee this ever happening, but just in case
	if len(g.nodeMap) == maxInt {
		panic("You have a full graph, so an ID can't be created (good job! You have a lot of memory!)")
	}

	for i := 0; i < maxInt; i++ {
		if _, ok := g.nodeMap[i]; !ok {
			g.AddNode(Node(i))
			return Node(i)
		}
	}

	// Will never happen
	return nil
}

func (g *Graph) AddNode(n graph.Node) {
	g.nodeMap[n.ID()] = n
	g.neighbors[n.ID()] = make(map[int]WeightedEdge)

	g.maxID = max(g.maxID, n.ID())
}

func (g *Graph) AddUndirectedEdge(e graph.Edge, cost float64) {
	head, tail := e.Head(), e.Tail()
	if !g.NodeExists(head) {
		g.AddNode(head)
	}

	if !g.NodeExists(tail) {
		g.AddNode(tail)
	}

	g.neighbors[head.ID()][tail.ID()] = WeightedEdge{Edge: e, Cost: cost}
	g.neighbors[tail.ID()][head.ID()] = WeightedEdge{Edge: e, Cost: cost}
}

func (g *Graph) RemoveNode(n graph.Node) {
	if _, ok := g.nodeMap[n.ID()]; !ok {
		return
	}
	delete(g.nodeMap, n.ID())

	for neigh, _ := range g.neighbors[n.ID()] {
		delete(g.neighbors[neigh], n.ID())
	}
	delete(g.neighbors, n.ID())

	g.maxID-- // Fun facts: even if this ID doesn't exist this still works!
}

func (g *Graph) RemoveUndirectedEdge(e graph.Edge) {
	head, tail := e.Head(), e.Tail()
	if _, ok := g.nodeMap[head.ID()]; !ok {
		return
	} else if _, ok := g.nodeMap[tail.ID()]; !ok {
		return
	}

	delete(g.neighbors[head.ID()], tail.ID())
	delete(g.neighbors[tail.ID()], head.ID())
}

func (g *Graph) EmptyGraph() {
	g.neighbors = make(map[int]map[int]WeightedEdge)
	g.nodeMap = make(map[int]graph.Node)
}

/* Graph implementation */

func (g *Graph) Neighbors(n graph.Node) []graph.Node {
	if !g.NodeExists(n) {
		return nil
	}

	neighbors := make([]graph.Node, len(g.neighbors[n.ID()]))
	i := 0
	for id, _ := range g.neighbors[n.ID()] {
		neighbors[i] = g.nodeMap[id]
		i++
	}

	return neighbors
}

func (g *Graph) EdgeBetween(n, neigh graph.Node) graph.Edge {
	// Don't need to check if neigh exists because
	// it's implicit in the neighbors access.
	if !g.NodeExists(n) {
		return nil
	}

	return g.neighbors[n.ID()][neigh.ID()]
}

func (g *Graph) NodeExists(n graph.Node) bool {
	_, ok := g.nodeMap[n.ID()]

	return ok
}

func (g *Graph) NodeList() []graph.Node {
	nodes := make([]graph.Node, len(g.nodeMap))
	i := 0
	for _, n := range g.nodeMap {
		nodes[i] = n
		i++
	}

	return nodes
}

func (g *Graph) Cost(e graph.Edge) float64 {
	if n, ok := g.neighbors[e.Head().ID()]; ok {
		if we, ok := n[e.Tail().ID()]; ok {
			return we.Cost
		}
	}
	return math.Inf(1)
}

func (g *Graph) EdgeList() []graph.Edge {
	m := make(map[WeightedEdge]struct{})
	toReturn := make([]graph.Edge, 0)

	for _, neighs := range g.neighbors {
		for _, we := range neighs {
			if _, ok := m[we]; !ok {
				m[we] = struct{}{}
				toReturn = append(toReturn, we.Edge)
			}
		}
	}

	return toReturn
}
