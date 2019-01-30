// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flow

import "gonum.org/v1/gonum/graph"

// Graph is a control flow graph, a directed graph rooted at an entry node.
type Graph interface {
	// Entry returns the entry node of the control flow graph.
	Entry() graph.Node
	graph.Directed
}

// NewGraph returns a new control flow graph based on the given directed graph
// and entry node ID.
func NewGraph(g graph.Directed, entryid int64) Graph {
	return &cfg{
		Directed: g,
		entry:    g.Node(entryid),
	}
}

// cfg is a control flow graph rooted at an entry node.
type cfg struct {
	// Entry node of control flow graph.
	entry graph.Node
	graph.Directed
}

// Entry returns the entry node of the control flow graph.
func (g *cfg) Entry() graph.Node {
	return g.entry
}
