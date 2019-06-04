// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package layout

import (
	"math"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/spatial/barneshut"
	"gonum.org/v1/gonum/spatial/r2"
)

// Eades implements the graph layout algorithm described in "A
// heuristic for graph drawing", Congressus numerantium 42:149-160.
// The implementation here uses the Barnes-Hut approximation for
// global repulsion calculation
type Eades struct {
	// M is the number of updates to perform.
	M int

	// C1, C2, C3 and C4 are the constants
	// described in the paper.
	C1, C2, C3, C4 float64

	// Theta is the Barnes-Hut theta constant.
	Theta float64

	// Src is the source of randomness used
	// to initialize the nodes' locations. If
	// Src is nil, the global random number
	// generator is used.
	Src rand.Source

	nodes   graph.Nodes
	indexOf map[int64]int

	particles []barneshut.Particle2
	forces    []r2.Vec
}

// Update is the Eades spatial graph update function.
func (u *Eades) Update(g graph.Graph, loc map[int64]r2.Vec) bool {
	if u.M <= 0 {
		return false
	}
	u.M--

	if len(loc) == 0 {
		var rnd func() float64
		if u.Src == nil {
			rnd = rand.Float64
		} else {
			rnd = rand.New(u.Src).Float64
		}
		u.nodes = g.Nodes()
		u.indexOf = make(map[int64]int, u.nodes.Len())
		u.particles = make([]barneshut.Particle2, 0, u.nodes.Len())
		u.forces = make([]r2.Vec, u.nodes.Len())
		for u.nodes.Next() {
			id := u.nodes.Node().ID()
			u.indexOf[id] = len(u.particles)
			u.particles = append(u.particles, eadesNode{id: id, pos: r2.Vec{X: rnd(), Y: rnd()}})
		}
	} else {
		u.nodes.Reset()
	}

	plane := barneshut.NewPlane(u.particles)
	for i, p := range u.particles {
		u.forces[i] = plane.ForceOn(p, u.Theta, barneshut.Gravity2).Scale(-u.C3)
	}

	for u.nodes.Next() {
		xid := u.nodes.Node().ID()
		to := g.From(xid)
		for to.Next() {
			yid := to.Node().ID()

			// Treat all edges as undirected for the purposes of force,
			// so only do this work once for each edge that exists.
			if xid < yid && g.HasEdgeBetween(xid, yid) {
				idx := u.indexOf[xid]

				// Undo repulsion of adjacent node.
				v := u.particles[u.indexOf[yid]].Coord2().Sub(u.particles[idx].Coord2())
				f := u.forces[idx].Add(barneshut.Gravity2(nil, nil, 1, 1, v).Scale(u.C3))

				// Apply adjacent node attraction.
				u.forces[idx] = f.Add(v.Scale(u.C1 * math.Log(math.Hypot(v.X, v.Y))))
			}
		}
	}

	for i, f := range u.forces {
		n := u.particles[i].(eadesNode)
		n.pos = n.pos.Add(f.Scale(u.C4))
		u.particles[i] = n
		loc[n.id] = n.pos
	}

	return true
}

type eadesNode struct {
	id  int64
	pos r2.Vec
}

func (p eadesNode) Coord2() r2.Vec { return p.pos }
func (p eadesNode) Mass() float64  { return 1 }
