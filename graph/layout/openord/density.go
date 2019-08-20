// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openord

import (
	"math"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/spatial/r2"
)

const (
	radius     = 10
	gridSize   = 1000
	view       = 4000
	viewToGrid = float64(gridSize) / float64(view)
)

type densityGrid struct {
	// The approach taken here is the apparently old
	// static allocation approach used by OpenOrd. The
	// current OpenOrd code dynamically allocates the
	// work spaces.
	//
	// TODO(kortschak): Revisit this.
	fallOff [radius*2 + 1][radius*2 + 1]float64
	density [gridSize][gridSize]float64
	bins    [gridSize][gridSize]queue
}

func newDensityGrid() *densityGrid {
	var g densityGrid
	for i := -radius; i <= radius; i++ {
		for j := -radius; j <= radius; j++ {
			g.fallOff[i+radius][j+radius] = ((radius - math.Abs(float64(i))) / radius) * ((radius - math.Abs(float64(j))) / radius)
		}
	}
	return &g
}

func (g *densityGrid) at(pos r2.Vec, fine bool) float64 {
	x := int((pos.X + view/2 + 0.5) * viewToGrid)
	y := int((pos.Y + view/2 + 0.5) * viewToGrid)

	const boundary = 10
	if y < boundary || gridSize-boundary < y {
		return 1e4
	}
	if x < boundary || gridSize-boundary < x {
		return 1e4
	}

	if !fine {
		d := g.density[y][x]
		return d * d
	}

	var d float64
	for i := y - 1; i <= y+1; i++ {
		for j := x - 1; j <= x+1; j++ {
			for _, r := range g.bins[i][j].slice() {
				v := pos.Sub(r.pos)
				d = v.X*v.X + v.Y*v.Y
				d += 1e-4 / (d + 1e-50)
			}
		}
	}
	return d
}

func (g *densityGrid) add(n *node, fine bool) {
	if fine {
		g.fineAdd(n)
	} else {
		g.coarseAdd(n)
	}
}

func (g *densityGrid) fineAdd(n *node) {
	x := int((n.pos.X + view/2 + 0.5) * viewToGrid)
	y := int((n.pos.Y + view/2 + 0.5) * viewToGrid)
	n.subPos = n.pos
	g.bins[y][x].enqueue(n)
}

func (g *densityGrid) coarseAdd(n *node) {
	x := int((n.pos.X+view/2+0.5)*viewToGrid) - radius
	y := int((n.pos.Y+view/2+0.5)*viewToGrid) - radius
	if x < 0 || gridSize <= x {
		panic("openord: node outside grid")
	}
	if y < 0 || gridSize <= y {
		panic("openord: node outside grid")
	}
	n.subPos = n.pos
	for i := 0; i <= radius*2; i++ {
		for j := 0; j <= radius*2; j++ {
			g.density[y+i][x+j] += g.fallOff[i][j]
		}
	}
}

func (g *densityGrid) sub(n *node, firstAdd, fineFirstAdd, fine bool) {
	if fine && !fineFirstAdd {
		g.fineSub(n)
	} else if !firstAdd {
		g.coarseSub(n)
	}
}

func (g *densityGrid) fineSub(n *node) {
	x := int((n.pos.X + view/2 + 0.5) * viewToGrid)
	y := int((n.pos.Y + view/2 + 0.5) * viewToGrid)
	g.bins[y][x].dequeue()
}

func (g *densityGrid) coarseSub(n *node) {
	x := int((n.pos.X+view/2+0.5)*viewToGrid) - radius
	y := int((n.pos.Y+view/2+0.5)*viewToGrid) - radius
	for i := 0; i <= radius*2; i++ {
		for j := 0; j <= radius*2; j++ {
			g.density[y+i][x+j] -= g.fallOff[i][j]
		}
	}
}

type node struct {
	node graph.Node

	fixed bool

	pos, subPos r2.Vec

	energy float64
}

// queue implements a FIFO queue.
type queue struct {
	head int
	data []*node
}

// len returns the number of nodes in the queue.
func (q *queue) len() int { return len(q.data) - q.head }

// enqueue adds the node n to the back of the queue.
func (q *queue) enqueue(n *node) {
	if len(q.data) == cap(q.data) && q.head > 0 {
		l := q.len()
		copy(q.data, q.data[q.head:])
		q.head = 0
		q.data = append(q.data[:l], n)
	} else {
		q.data = append(q.data, n)
	}
}

// dequeue returns the node at the front of the queue and
// removes it from the queue.
func (q *queue) dequeue() *node {
	if q.len() == 0 {
		panic("openord: empty queue")
	}

	var n *node
	n, q.data[q.head] = q.data[q.head], n
	q.head++

	if q.len() == 0 {
		q.reset()
	}

	return n
}

func (q *queue) slice() []*node {
	return q.data[q.head:]
}

// reset clears the queue for reuse.
func (q *queue) reset() {
	q.head = 0
	q.data = q.data[:0]
}
