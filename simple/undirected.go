// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simple

import (
	"fmt"

	"github.com/gonum/graph"
)

// UndirectedGraph implements a generalized undirected graph.
type UndirectedGraph struct {
	neighbors map[int]map[int]graph.Edge
	nodeMap   map[int]graph.Node

	self, absent float64

	// Node add/remove convenience vars
	maxID   int
	freeMap map[int]struct{}
}

// NewUndiectedGraph returns an UndirectedGraph with the specified self and absent edge costs.
func NewUndirectedGraph(self, absent float64) *UndirectedGraph {
	return &UndirectedGraph{
		neighbors: make(map[int]map[int]graph.Edge),
		nodeMap:   make(map[int]graph.Node),

		self:   self,
		absent: absent,

		maxID:   0,
		freeMap: make(map[int]struct{}),
	}
}

func (g *UndirectedGraph) NewNodeID() int {
	if g.maxID != maxInt {
		g.maxID++
		return g.maxID
	}

	// Implicitly checks if len(g.freeMap) == 0
	for id := range g.freeMap {
		return id
	}

	// I cannot foresee this ever happening, but just in case, we check.
	if len(g.nodeMap) == maxInt {
		panic("cannot allocate node: graph too large")
	}

	for i := 0; i < maxInt; i++ {
		if _, ok := g.nodeMap[i]; !ok {
			return i
		}
	}

	// Should not happen.
	panic("cannot allocate node id: no free id found")
}

func (g *UndirectedGraph) AddNode(n graph.Node) {
	if _, exists := g.nodeMap[n.ID()]; exists {
		panic(fmt.Sprintf("simple: node ID collision: %d", n.ID()))
	}
	g.nodeMap[n.ID()] = n
	g.neighbors[n.ID()] = make(map[int]graph.Edge)

	delete(g.freeMap, n.ID())
	g.maxID = max(g.maxID, n.ID())
}

func (g *UndirectedGraph) SetEdge(e graph.Edge) {
	var (
		from = e.From()
		fid  = from.ID()
		to   = e.To()
		tid  = to.ID()
	)

	if fid == tid {
		panic("simple: adding self edge")
	}

	if !g.Has(from) {
		g.AddNode(from)
	}

	if !g.Has(to) {
		g.AddNode(to)
	}

	g.neighbors[fid][tid] = e
	g.neighbors[tid][fid] = e
}

func (g *UndirectedGraph) RemoveNode(n graph.Node) {
	if _, ok := g.nodeMap[n.ID()]; !ok {
		return
	}
	delete(g.nodeMap, n.ID())

	for neigh := range g.neighbors[n.ID()] {
		delete(g.neighbors[neigh], n.ID())
	}
	delete(g.neighbors, n.ID())

	g.freeMap[n.ID()] = struct{}{}
}

func (g *UndirectedGraph) RemoveEdge(e graph.Edge) {
	from, to := e.From(), e.To()
	if _, ok := g.nodeMap[from.ID()]; !ok {
		return
	} else if _, ok := g.nodeMap[to.ID()]; !ok {
		return
	}

	delete(g.neighbors[from.ID()], to.ID())
	delete(g.neighbors[to.ID()], from.ID())
}

func (g *UndirectedGraph) EmptyGraph() {
	g.neighbors = make(map[int]map[int]graph.Edge)
	g.nodeMap = make(map[int]graph.Node)
}

/* UndirectedGraph implementation */

func (g *UndirectedGraph) From(n graph.Node) []graph.Node {
	if !g.Has(n) {
		return nil
	}

	neighbors := make([]graph.Node, len(g.neighbors[n.ID()]))
	i := 0
	for id := range g.neighbors[n.ID()] {
		neighbors[i] = g.nodeMap[id]
		i++
	}

	return neighbors
}

func (g *UndirectedGraph) HasEdgeBetween(n, neigh graph.Node) bool {
	_, ok := g.neighbors[n.ID()][neigh.ID()]
	return ok
}

func (g *UndirectedGraph) Edge(u, v graph.Node) graph.Edge {
	return g.EdgeBetween(u, v)
}

func (g *UndirectedGraph) EdgeBetween(u, v graph.Node) graph.Edge {
	// We don't need to check if neigh exists because
	// it's implicit in the neighbors access.
	if !g.Has(u) {
		return nil
	}

	return g.neighbors[u.ID()][v.ID()]
}

func (g *UndirectedGraph) Node(id int) graph.Node {
	return g.nodeMap[id]
}

func (g *UndirectedGraph) Has(n graph.Node) bool {
	_, ok := g.nodeMap[n.ID()]

	return ok
}

func (g *UndirectedGraph) Nodes() []graph.Node {
	nodes := make([]graph.Node, len(g.nodeMap))
	i := 0
	for _, n := range g.nodeMap {
		nodes[i] = n
		i++
	}

	return nodes
}

func (g *UndirectedGraph) Weight(x, y graph.Node) (w float64, ok bool) {
	xid := x.ID()
	yid := y.ID()
	if xid == yid {
		return g.self, true
	}
	if n, ok := g.neighbors[xid]; ok {
		if e, ok := n[yid]; ok {
			return e.Weight(), true
		}
	}
	return g.absent, false
}

func (g *UndirectedGraph) Edges() []graph.Edge {
	var edges []graph.Edge

	seen := make(map[[2]int]struct{})
	for _, u := range g.neighbors {
		for _, e := range u {
			uid := e.From().ID()
			vid := e.To().ID()
			if _, ok := seen[[2]int{uid, vid}]; ok {
				continue
			}
			seen[[2]int{uid, vid}] = struct{}{}
			seen[[2]int{vid, uid}] = struct{}{}
			edges = append(edges, e)
		}
	}

	return edges
}

func (g *UndirectedGraph) Degree(n graph.Node) int {
	if _, ok := g.nodeMap[n.ID()]; !ok {
		return 0
	}

	return len(g.neighbors[n.ID()])
}
