// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simple

import (
	"github.com/gonum/graph"
	"github.com/gonum/matrix/mat64"
)

// DirectedDenseGraph represents a graph such that all IDs are in a contiguous
// block from 0 to n-1.
type DirectedDenseGraph struct {
	self   float64
	absent float64
	mat    *mat64.Dense
}

// NewDirectedDenseGraph creates a directed dense graph with n nodes.
// If passable is true all pairs of nodes will be connected by an edge
// with unit cost, otherwise every node will start unconnected with
// the cost specified by absent. The self parameter specifies the cost
// of self connection.
func NewDirectedDenseGraph(n int, passable bool, self, absent float64) *DirectedDenseGraph {
	mat := make([]float64, n*n)
	v := 1.
	if !passable {
		v = absent
	}
	for i := range mat {
		mat[i] = v
	}
	return &DirectedDenseGraph{
		mat:    mat64.NewDense(n, n, mat),
		self:   self,
		absent: absent,
	}
}

func (g *DirectedDenseGraph) Has(n graph.Node) bool {
	return g.has(n.ID())
}

func (g *DirectedDenseGraph) has(id int) bool {
	r, _ := g.mat.Dims()
	return 0 <= id && id < r
}

func (g *DirectedDenseGraph) Nodes() []graph.Node {
	r, _ := g.mat.Dims()
	nodes := make([]graph.Node, r)
	for i := 0; i < r; i++ {
		nodes[i] = Node(i)
	}
	return nodes
}

func (g *DirectedDenseGraph) Edges() []graph.Edge {
	var edges []graph.Edge
	r, _ := g.mat.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < r; j++ {
			if i == j {
				continue
			}
			if w := g.mat.At(i, j); !isSame(w, g.absent) {
				edges = append(edges, Edge{F: Node(i), T: Node(j), W: w})
			}
		}
	}
	return edges
}

func (g *DirectedDenseGraph) From(n graph.Node) []graph.Node {
	var neighbors []graph.Node
	id := n.ID()
	_, c := g.mat.Dims()
	for j := 0; j < c; j++ {
		if j == id {
			continue
		}
		if !isSame(g.mat.At(id, j), g.absent) {
			neighbors = append(neighbors, Node(j))
		}
	}
	return neighbors
}

func (g *DirectedDenseGraph) To(n graph.Node) []graph.Node {
	var neighbors []graph.Node
	id := n.ID()
	r, _ := g.mat.Dims()
	for i := 0; i < r; i++ {
		if i == id {
			continue
		}
		if !isSame(g.mat.At(i, id), g.absent) {
			neighbors = append(neighbors, Node(i))
		}
	}
	return neighbors
}

func (g *DirectedDenseGraph) HasEdge(x, y graph.Node) bool {
	xid := x.ID()
	yid := y.ID()
	return xid != yid && (!isSame(g.mat.At(xid, yid), g.absent) || !isSame(g.mat.At(yid, xid), g.absent))
}

func (g *DirectedDenseGraph) Edge(u, v graph.Node) graph.Edge {
	if g.HasEdge(u, v) {
		return Edge{F: u, T: v, W: g.mat.At(u.ID(), v.ID())}
	}
	return nil
}

func (g *DirectedDenseGraph) HasEdgeFromTo(u, v graph.Node) bool {
	uid := u.ID()
	vid := v.ID()
	return uid != vid && !isSame(g.mat.At(uid, vid), g.absent)
}

func (g *DirectedDenseGraph) Weight(x, y graph.Node) (w float64, ok bool) {
	xid := x.ID()
	yid := y.ID()
	if xid == yid {
		return 0, true
	}
	if g.has(xid) && g.has(yid) {
		return g.mat.At(xid, yid), true
	}
	return g.absent, false
}

func (g *DirectedDenseGraph) SetEdgeWeight(e graph.Edge) {
	fid := e.From().ID()
	tid := e.To().ID()
	if fid == tid {
		panic("simple: set edge cost of illegal edge")
	}
	g.mat.Set(fid, tid, e.Weight())
}

func (g *DirectedDenseGraph) RemoveEdge(e graph.Edge) {
	g.mat.Set(e.From().ID(), e.To().ID(), g.absent)
}

func (g *DirectedDenseGraph) Matrix() mat64.Matrix {
	// Prevent alteration of dimensions of the returned matrix.
	m := *g.mat
	return &m
}
