// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flow

import "gonum.org/v1/gonum/graph"

// Graph is a control flow graph, a directed graph rooted at an entry node.
type Graph struct {
	// Underlying directed graph.
	graph.Directed
	// Root node of the control flow graph.
	RootID int64
}

// Root returns the root node of the control flow graph.
func (g Graph) Root() graph.Node {
	return g.Directed.Node(g.RootID)
}
