// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package layout

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/spatial/r2"
)

// LayoutR2 implements graph layout updates and representations.
type LayoutR2 interface {
	// IsInitialized returns whether the Layout is initialized.
	IsInitialized() bool

	// SetLocation sets the location of the node with the given
	// id to pos.
	SetLocation(id int64, pos r2.Vec)

	// Location returns the location of the node with the given
	// id in the graph layout.
	Location(id int64) r2.Vec
}

// NewOptimizer returns a new layout optimizer. If g implements Layout the layout
// will be updated into g, otherwise the Optimizer will hold the graph layout.
func NewOptimizer(g graph.Graph, update func(graph.Graph, Layout) bool) Optimizer {
	l, ok := g.(Layout)
	if !ok {
		l = make(coordinates)
	}
	return Optimizer{
		g:       g,
		layout:  l,
		Updater: update,
	}
}

// Optimizer is a helper type that holds a graph and layout
// optimization state.
type Optimizer struct {
	g      graph.Graph
	layout Layout

	// Updater is the function called for each call to Update.
	// It updates the Optimizer's spatial distribution of the
	// nodes in the backing graph.
	Updater func(graph.Graph, Layout) bool
}

// Location returns the location of the node with the given
// ID. The returned value is only valid if the node exists
// in the graph.
func (g Optimizer) Location(id int64) r2.Vec {
	return g.layout.Location(id)
}

// Update updates the locations of the nodes in the graph
// according to the provided update function. It returns whether
// the update function is able to further refine the graph's
// node locations.
func (g Optimizer) Update() bool {
	if g.Updater == nil {
		return false
	}
	return g.Updater(g.g, g.layout)
}

// coordinates is the default layout store.
type coordinates map[int64]r2.Vec

func (c coordinates) IsInitialized() bool              { return len(c) != 0 }
func (c coordinates) SetLocation(id int64, pos r2.Vec) { c[id] = pos }
func (c coordinates) Location(id int64) r2.Vec         { return c[id] }
