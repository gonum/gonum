// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package layout

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/spatial/r2"
)

// TODO(kortschak): Consider whether 3D layouts are an important thing and what to do about them.
// TODO(kortschak): Consider whether Optimizer should be a graph.Graph, renamed layout.Graph.
// TODO(kortschak): Consider whether Optimizer.update should be exported to allow re-optimisation
// with a different algorithm.

// NewOptimizer returns a new layout optimizer.
func NewOptimizer(g graph.Graph, update func(graph.Graph, map[int64]r2.Vec) bool) Optimizer {
	return Optimizer{
		g:         g,
		locations: make(map[int64]r2.Vec),
		update:    update,
	}
}

// Optimizer is a helper type that holds a graph and layout
// optimization state.
type Optimizer struct {
	g         graph.Graph
	locations map[int64]r2.Vec
	update    func(graph.Graph, map[int64]r2.Vec) bool
}

// Location returns the location of the node with the given
// ID. The returned value is only valid if the node exists
// in the graph.
func (g Optimizer) Location(id int64) r2.Vec {
	return g.locations[id]
}

// Update updates the locations of the nodes in the graph
// according to the provided update function. It returns whether
// the update function is able to further refine the graph's
// node locations.
func (g Optimizer) Update() bool {
	return g.update(g.g, g.locations)
}
