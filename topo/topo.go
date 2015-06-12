// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topo

import (
	"github.com/gonum/graph"
	"github.com/gonum/graph/traverse"
)

// IsPath returns true for a connected path within a graph.
//
// IsPath returns true if, starting at path[0] and ending at path[len(path)-1], all nodes between
// are valid neighbors. That is, for each element path[i], path[i+1] is a valid successor.
//
// As special cases, IsPath returns true for a nil or zero length path, and for a path of length 1
// (only one node) but only if the node listed in path exists within the graph.
//
// Graph must be non-nil.
func IsPathIn(g graph.Graph, path []graph.Node) bool {
	var canReach func(u, v graph.Node) bool
	switch g := g.(type) {
	case graph.Directed:
		canReach = func(u, v graph.Node) bool {
			return g.EdgeFromTo(u, v) != nil
		}
	default:
		canReach = g.HasEdge
	}

	if path == nil || len(path) == 0 {
		return true
	} else if len(path) == 1 {
		return g.Has(path[0])
	}

	for i := 0; i < len(path)-1; i++ {
		if !canReach(path[i], path[i+1]) {
			return false
		}
	}

	return true
}

// ConnectedComponents returns the connected components of the graph g. All
// edges are treated as undirected.
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
