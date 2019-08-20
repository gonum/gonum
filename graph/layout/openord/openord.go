// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openord

import (
	"log"
	"math"
	"time"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/spatial/r2"
)

const maxProcs = 256

type Layout struct {
	workers []*worker
}

func NewLayout(g graph.Graph, n int, src rand.Source) Layout {
	d := newDescription(g)
	w := make([]*worker, n)
	seed := rand.New(src).Uint64()
	for i := range w {
		w[i] = newWorker(i, n, d, rand.NewSource(seed))
	}
	return Layout{w}
}

type worker struct {
	id, workers int

	rnd *rand.Rand

	*description
	neighbors map[int64]map[int64]float64

	grid *densityGrid

	stage int
	LayoutSchedule

	EdgeCuts

	firstAdd, fineFirstAdd, fineDensity bool

	liquid    LayoutSchedule
	expansion LayoutSchedule
	cooldown  LayoutSchedule
	crunch    LayoutSchedule
	simmer    LayoutSchedule

	start, stop time.Time

	totalIters int

	fixedUntil int // real_iterations
	fixed      bool
}

type description struct {
	g             graph.Graph
	indexOf       map[int64]int
	highestWeight float64
	positions     []*node
}

type EdgeCuts struct {
	MinEdges  float64 // min_edges
	End       float64 // CUR_END
	MaxLength float64 // cut_length_end
	Length    float64 // cut_off_length
	Rate      float64 // cut_rate
}

type LayoutSchedule struct {
	Iters       int
	Temperature float64
	Attraction  float64
	Damping     float64
	elapsed     time.Duration
}

func newWorker(id, workers int, d *description, src rand.Source) *worker {
	return &worker{
		id: id, workers: workers,

		rnd: rand.New(src),

		description: d,
		neighbors:   neighborsFrom(id, workers, d),

		LayoutSchedule: LayoutSchedule{
			Temperature: 2000,
			Attraction:  10,
			Damping:     1,
		},

		EdgeCuts: EdgeCuts{
			MinEdges: 20,
		},

		liquid: LayoutSchedule{
			Iters:       200,
			Temperature: 2000,
			Attraction:  2,
			Damping:     1,
		},
		expansion: LayoutSchedule{
			Iters:       200,
			Temperature: 2000,
			Attraction:  10,
			Damping:     1,
		},
		cooldown: LayoutSchedule{
			Iters:       200,
			Temperature: 2000,
			Attraction:  1,
			Damping:     0.1,
		},
		crunch: LayoutSchedule{
			Iters:       50,
			Temperature: 250,
			Attraction:  1,
			Damping:     0.25,
		},
		simmer: LayoutSchedule{
			Iters:       100,
			Temperature: 250,
			Attraction:  0.5,
			Damping:     0.0,
		},

		firstAdd: true, fineFirstAdd: true,
		grid: newDensityGrid(),
	}
}

func newDescription(g graph.Graph) *description {
	nodes := g.Nodes()
	if nodes.Len() == 0 {
		return nil
	}

	indexOf := make(map[int64]int, nodes.Len())
	positions := make([]*node, nodes.Len())
	i := 0
	for nodes.Next() {
		n := nodes.Node()
		indexOf[n.ID()] = i
		positions[i] = &node{node: n}
		i++
	}

	return &description{
		g:             g,
		indexOf:       indexOf,
		highestWeight: highestWeight(g, positions),
		positions:     positions,
	}
}

func highestWeight(g graph.Graph, positions []*node) float64 {
	wg, ok := g.(graph.Weighted)
	if !ok {
		return 1
	}

	highestWeight := -1.0
	switch eg := g.(type) {
	case interface{ WeightedEdges() graph.WeightedEdges }:
		edges := eg.WeightedEdges()
		for edges.Next() {
			w := edges.WeightedEdge().Weight()
			if w < 0 {
				panic("openord: negative edge weight")
			}
			highestWeight = math.Max(highestWeight, w)
		}
	case interface{ Edges() graph.Edges }:
		edges := eg.Edges()
		for edges.Next() {
			e := edges.Edge()
			w, ok := wg.Weight(e.From().ID(), e.To().ID())
			if !ok {
				panic("openord: missing weight for existing edge")
			}
			if w < 0 {
				panic("openord: negative edge weight")
			}
			highestWeight = math.Max(highestWeight, w)
		}
	default:
		for _, u := range positions {
			to := g.From(u.node.ID())
			for to.Next() {
				v := to.Node()
				w, ok := wg.Weight(u.node.ID(), v.ID())
				if !ok {
					panic("openord: missing weight for existing edge")
				}
				if w < 0 {
					panic("openord: negative edge weight")
				}
				highestWeight = math.Max(highestWeight, w)
			}
		}
	}

	return highestWeight
}

func neighborsFrom(id, workers int, d *description) map[int64]map[int64]float64 {
	ug := make(ugraph)
	switch eg := d.g.(type) {
	case interface{ WeightedEdges() graph.WeightedEdges }:
		edges := eg.WeightedEdges()
		for edges.Next() {
			e := edges.WeightedEdge()
			uid := e.From().ID()
			vid := e.To().ID()
			w := e.Weight() / d.highestWeight
			w *= w
			if uid%int64(workers) == int64(id) {
				ug.setEdge(uid, vid, w)
			}
			if vid%int64(workers) == int64(id) {
				ug.setEdge(vid, uid, w)
			}
		}
	case interface{ Edges() graph.Edges }:
		weightFunc := weightFuncFrom(d.g)
		edges := eg.Edges()
		for edges.Next() {
			e := edges.Edge()
			uid := e.From().ID()
			vid := e.To().ID()
			w, _ := weightFunc(uid, vid) // Existence is already checked.
			w /= d.highestWeight
			w *= w
			if uid%int64(workers) == int64(id) {
				ug.setEdge(uid, vid, w)
			}
			if vid%int64(workers) == int64(id) {
				ug.setEdge(vid, uid, w)
			}
		}
	default:
		weightFunc := weightFuncFrom(d.g)
		for _, u := range d.positions {
			uid := u.node.ID()
			to := d.g.From(uid)
			for to.Next() {
				vid := to.Node().ID()
				w, _ := weightFunc(uid, vid) // Existence is already checked.
				w /= d.highestWeight
				w *= w
				if uid%int64(workers) == int64(id) {
					ug.setEdge(uid, vid, w)
				}
				if vid%int64(workers) == int64(id) {
					ug.setEdge(vid, uid, w)
				}
			}
		}
	}
	return ug
}

func weightFuncFrom(g graph.Graph) func(uid, vid int64) (w float64, ok bool) {
	if wg, ok := g.(graph.Weighted); ok {
		return wg.Weight
	}
	return func(_, _ int64) (w float64, ok bool) {
		return 1, true
	}
}

type ugraph map[int64]map[int64]float64

func (g ugraph) setEdge(uid, vid int64, weight float64) {
	u, ok := g[uid]
	if !ok {
		u = make(map[int64]float64)
		g[uid] = u
	}
	u[vid] = weight
}

func (l *worker) recompute() bool {
	const min = 1

	l.updateNodes()

	l.totalIters++
	if l.totalIters >= l.fixedUntil {
		l.fixed = false
	}

	const (
		liquid = iota
		expansion
		cooldown
		crunch
		simmer
		done
	)
	switch l.stage {
	case liquid:
		if l.Iters == 0 {
			if l.id == 0 {
				log.Println("entering liquid stage")
			}

			l.start = time.Now()

			l.LayoutSchedule = l.liquid
			l.Iters = 0
		}

		if l.Iters < l.liquid.Iters {
			l.Iters++
		} else {
			l.stop = time.Now()
			l.liquid.elapsed += l.stop.Sub(l.start)

			if l.id == 0 {
				log.Printf("liquid stage completed in %v energy=%v", l.liquid.elapsed, l.totalEnergy())
			}

			l.stage = expansion
			l.Iters = 0
		}

	case expansion:
		if l.Iters == 0 {
			if l.id == 0 {
				log.Println("entering expansion stage")
			}

			l.start = time.Now()

			l.LayoutSchedule = l.expansion
			l.Iters = 0
		}

		if l.Iters < l.expansion.Iters {
			if l.Attraction > 1 {
				l.Attraction -= 0.05
			}
			if l.MinEdges > 12 {
				l.MinEdges -= 0.05
			}
			l.Length -= l.Rate
			if l.Damping > 0.1 {
				l.Damping -= 0.005
			}
			l.Iters++
		} else {
			l.stop = time.Now()
			l.expansion.elapsed += l.stop.Sub(l.start)

			if l.id == 0 {
				log.Printf("expansion stage completed in %v energy=%v", l.expansion.elapsed, l.totalEnergy())
			}

			l.stage = cooldown
			l.Iters = 0
		}

	case cooldown:
		if l.Iters == 0 {
			if l.id == 0 {
				log.Println("entering cooldown stage")
			}

			l.start = time.Now()

			l.MinEdges = 12
			l.LayoutSchedule = l.cooldown
			l.Iters = 0
		}

		if l.Iters < l.cooldown.Iters {
			if l.Temperature > 50 {
				l.Temperature -= 10
			}
			if l.Length > l.End {
				l.Length -= 2 * l.Rate
			}
			if l.MinEdges > min {
				l.MinEdges -= 0.2
			}
			l.Iters++
		} else {
			l.stop = time.Now()
			l.cooldown.elapsed += l.stop.Sub(l.start)

			if l.id == 0 {
				log.Printf("cooldown stage completed in %v energy=%v", l.cooldown.elapsed, l.totalEnergy())
			}

			l.stage = crunch
			l.Iters = 0
		}

	case crunch:
		if l.Iters == 0 {
			if l.id == 0 {
				log.Println("entering crunch stage")
			}

			l.start = time.Now()

			l.Length = l.End
			l.MinEdges = min
			l.LayoutSchedule = l.crunch
			l.Iters = 0
		}

		if l.Iters < l.crunch.Iters {
			l.Iters++
		} else {
			l.stop = time.Now()
			l.crunch.elapsed += l.stop.Sub(l.start)

			if l.id == 0 {
				log.Printf("crunch stage completed in %v energy=%v", l.crunch.elapsed, l.totalEnergy())
			}

			l.stage = simmer
			l.Iters = 0
		}

	case simmer:
		if l.Iters == 0 {
			if l.id == 0 {
				log.Println("entering simmer stage")
			}

			l.start = time.Now()

			l.MinEdges = math.NaN()
			l.LayoutSchedule = l.simmer
			l.fineDensity = true
			l.Iters = 0
		}

		if l.Iters < l.simmer.Iters {
			if l.Temperature > 50 {
				l.Temperature -= 2
			}
			l.Iters++
		} else {
			l.stop = time.Now()
			l.simmer.elapsed += l.stop.Sub(l.start)

			if l.id == 0 {
				log.Printf("simmer stage completed in %v energy=%v", l.simmer.elapsed, l.totalEnergy())
			}

			l.stage = done
		}

	case done:
		if l.id == 0 {
			log.Printf("layout completed in %v",
				l.liquid.elapsed+l.expansion.elapsed+l.cooldown.elapsed+l.crunch.elapsed+l.simmer.elapsed)
		}
		return false
	}
	return true
}

func (l *worker) updateNodes() {
	old := make([]r2.Vec, maxProcs)
	new := make([]r2.Vec, maxProcs)

	var allFixed bool

	indices := make([]int, l.workers)
	for i := 0; i < l.workers; i++ {
		indices[i] = i
	}

	n := float64(l.workers)
	squareNumNodes := int(n + n*math.Floor(float64(len(l.indexOf)-1)/n))
	for i := l.id; i < squareNumNodes; i += l.workers {
		l.positions(old, indices)
		l.positions(new, indices)

		if i < len(l.indexOf) {
			for j := 0; j < 2*l.id; j++ {
				l.rnd.Float64()
			}
			if !l.description.positions[i].fixed && l.fixed {
				l.updatePosition(i, old, new)
			}
			for j := 2 * l.id; j < 2*(len(indices)-1); j++ {
				l.rnd.Float64()
			}
		} else {
			for j := 0; j < 2*(len(indices)); j++ {
				l.rnd.Float64()
			}
		}

		allFixed = true
		for _, idx := range indices {
			if !l.description.positions[idx].fixed && l.fixed {
				allFixed = false
			}
		}
		if !allFixed {
			l.updateDensity(indices, old, new)
		}

		for j := range indices {
			indices[j] += l.workers
		}
		for indices[len(indices)-1] >= len(l.indexOf) {
			indices = indices[:len(indices)-1]
		}
	}

	l.firstAdd = false
	if l.fineDensity {
		l.fineFirstAdd = false
	}
}

func (l *worker) positions(dst []r2.Vec, indices []int) []r2.Vec {
	for i, idx := range indices {
		dst[i] = l.description.positions[idx].pos
	}
	return dst
}

func (l *worker) updatePosition(idx int, old, new []r2.Vec) {
	jump := 0.010 * l.Temperature

	l.grid.sub(l.description.positions[idx], l.firstAdd, l.fineFirstAdd, l.fineDensity)

	var energies [2]float64
	energies[0] = l.computeNodeEnergy(idx)

	var updatedPos [2]r2.Vec
	updatedPos[0] = l.solveAnalytic(idx)
	l.description.positions[idx].pos = updatedPos[0]
	updatedPos[1] = updatedPos[0].Add(r2.Vec{
		X: 0.5 - l.rnd.Float64()*jump,
		Y: 0.5 - l.rnd.Float64()*jump,
	})
	l.description.positions[idx].pos = updatedPos[1]
	energies[1] = l.computeNodeEnergy(idx)

	l.description.positions[idx].pos = old[l.id]
	if (!l.fineDensity && !l.firstAdd) || !l.fineFirstAdd {
		l.grid.add(l.description.positions[idx], l.fineDensity)
	}

	if energies[0] < energies[1] {
		new[l.id] = updatedPos[0]
		l.description.positions[idx].energy = energies[0]
	} else {
		new[l.id] = updatedPos[1]
		l.description.positions[idx].energy = energies[1]
	}
}

func (l *worker) updateDensity(indices []int, old, new []r2.Vec) {
	for i, idx := range indices {
		l.description.positions[idx].pos = old[i]
		l.grid.sub(l.description.positions[idx], l.firstAdd, l.fineFirstAdd, l.fineDensity)
		l.description.positions[idx].pos = new[i]
		l.grid.add(l.description.positions[idx], l.fineDensity)
	}
}

func (l *worker) computeNodeEnergy(idx int) float64 {
	attractionFactor := l.Attraction * l.Attraction * l.Attraction * l.Attraction * 2e-2

	var energy float64
	u := l.description.positions[idx]
	uid := u.node.ID()
	for vid, w := range l.neighbors[uid] {
		d := u.pos.Sub(l.description.positions[l.description.indexOf[vid]].pos)
		energyDistance := d.X*d.X + d.Y*d.Y

		if l.stage < 2 {
			energyDistance *= energyDistance

			// In the liquid phase we want to discourage long link distances
			if l.stage == 0 {
				energyDistance *= energyDistance
			}
		}

		energy += w * attractionFactor * energyDistance
	}

	energy += l.grid.at(u.pos, l.fineDensity)

	return energy
}

func (l *worker) solveAnalytic(idx int) r2.Vec {
	var weight float64
	var pos r2.Vec
	u := l.description.positions[idx]
	uid := u.node.ID()
	for vid, w := range l.neighbors[uid] {
		weight += w
		pos = pos.Add(l.description.positions[l.description.indexOf[vid]].pos.Scale(w))
	}

	var center r2.Vec
	if weight > 0 {
		center = pos.Scale(1 / weight)
		pos = l.description.positions[idx].pos.Scale(1 - l.Damping).Add(center.Scale(l.Damping))
	}

	if math.IsNaN(l.MinEdges) || l.End >= 39500 {
		return pos
	}
	deg := float64(len(l.neighbors[uid]))
	if deg < l.MinEdges {
		return pos
	}

	nConns := math.Sqrt(deg)
	maxLength := math.Inf(-1)
	var maxID int64
	for vid := range l.neighbors[uid] {
		p := center.Sub(l.description.positions[l.description.indexOf[vid]].pos)
		dis := (p.X*p.X + p.Y*p.Y) * nConns
		if dis > maxLength {
			maxLength = dis
			maxID = vid
		}
	}
	if maxLength > l.Length {
		delete(l.neighbors, maxID)
	}

	return pos
}

func (l *worker) totalEnergy() float64 {
	var energy float64
	for _, n := range l.description.positions {
		energy += n.energy
	}
	return energy
}
