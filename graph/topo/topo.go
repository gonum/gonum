// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package topo provides graph topology analysis functions.
package topo

import (
	"github.com/gonum/graph"
	"github.com/gonum/graph/traverse"
)

// IsPathIn returns whether path is a path in g.
//
// As special cases, IsPathIn returns true for a zero length path or for
// a path of length 1 when the node in path exists in the graph.
func IsPathIn(g graph.Graph, path []graph.Node) bool {
	switch len(path) {
	case 0:
		return true
	case 1:
		return g.Has(path[0])
	default:
		var canReach func(u, v graph.Node) bool
		switch g := g.(type) {
		case graph.Directed:
			canReach = g.HasEdgeFromTo
		default:
			canReach = g.HasEdgeBetween
		}

		for i, u := range path[:len(path)-1] {
			if !canReach(u, path[i+1]) {
				return false
			}
		}
		return true
	}
}

// PathExistsIn returns whether there is a path in g starting at from extending
// to to.
//
// PathExistsIn exists as a helper function. If many tests for path existence
// are being performed, other approaches will be more efficient.
func PathExistsIn(g graph.Graph, from, to graph.Node) bool {
	var t traverse.BreadthFirst
	return t.Walk(g, from, func(n graph.Node, _ int) bool { return n.ID() == to.ID() }) != nil
}

// ConnectedComponents returns the connected components of the undirected graph g.
func ConnectedComponents(g graph.Undirected) [][]graph.Node {
	var (
		w  traverse.DepthFirst
		c  []graph.Node
		cc [][]graph.Node
	)
	during := func(n graph.Node) {
		c = append(c, n)
	}
	after := func() {
		cc = append(cc, []graph.Node(nil))
		cc[len(cc)-1] = append(cc[len(cc)-1], c...)
		c = c[:0]
	}
	w.WalkAll(g, nil, after, during)

	return cc
}
