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

	// SetCoord2 sets the coordinates of the node with the given
	// id to coords.
	SetCoord2(id int64, coords r2.Vec)

	// Coord2 returns the coordinated of the node with the given
	// id in the graph layout.
	Coord2(id int64) r2.Vec
}

// NewOptimizerR2 returns a new layout optimizer. If g implements LayoutR2 the layout
// will be updated into g, otherwise the OptimizerR2 will hold the graph layout.
func NewOptimizerR2(g graph.Graph, update func(graph.Graph, LayoutR2) bool) OptimizerR2 {
	l, ok := g.(LayoutR2)
	if !ok {
		l = make(coordinates)
	}
	return OptimizerR2{
		g:       g,
		layout:  l,
		Updater: update,
	}
}

// OptimizerR2 is a helper type that holds a graph and layout
// optimization state.
type OptimizerR2 struct {
	g      graph.Graph
	layout LayoutR2

	// Updater is the function called for each call to Update.
	// It updates the OptimizerR2's spatial distribution of the
	// nodes in the backing graph.
	Updater func(graph.Graph, LayoutR2) bool
}

// Coord2 returns the location of the node with the given
// ID. The returned value is only valid if the node exists
// in the graph.
func (g OptimizerR2) Coord2(id int64) r2.Vec {
	return g.layout.Coord2(id)
}

// Update updates the locations of the nodes in the graph
// according to the provided update function. It returns whether
// the update function is able to further refine the graph's
// node locations.
func (g OptimizerR2) Update() bool {
	if g.Updater == nil {
		return false
	}
	return g.Updater(g.g, g.layout)
}

// coordinates is the default layout store.
type coordinates map[int64]r2.Vec

func (c coordinates) IsInitialized() bool            { return len(c) != 0 }
func (c coordinates) SetCoord2(id int64, pos r2.Vec) { c[id] = pos }
func (c coordinates) Coord2(id int64) r2.Vec         { return c[id] }
