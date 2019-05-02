// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package barneshut

import (
	"fmt"
	"math"
)

const (
	ne = iota
	se
	sw
	nw
)

// Vector2 is a 2D vector.
type Vector2 struct {
	X, Y float64
}

// Add returns the vector sum of p and q.
func (p Vector2) Add(q Vector2) Vector2 {
	p.X += q.X
	p.Y += q.Y
	return p
}

// Sub returns the vector sum of p and -q.
func (p Vector2) Sub(q Vector2) Vector2 {
	p.X -= q.X
	p.Y -= q.Y
	return p
}

// Scale returns the vector p scaled by f.
func (p Vector2) Scale(f float64) Vector2 {
	p.X *= f
	p.Y *= f
	return p
}

// Box2 is a 2D bounding box.
type Box2 struct {
	Min, Max Vector2
}

// quadrant returns which quadrant of b that p should be placed in.
func (b Box2) quadrant(p Particle2) int {
	center := Vector2{
		X: (b.Min.X + b.Max.X) / 2,
		Y: (b.Min.Y + b.Max.Y) / 2,
	}
	c := p.Coord2()
	if checkBounds && (c.X < b.Min.X || b.Max.X < c.X || c.Y < b.Min.Y || b.Max.Y < c.Y) {
		panic(fmt.Sprintf("p out of range %+v: %#v", b, p))
	}
	if c.X < center.X {
		if c.Y < center.Y {
			return nw
		} else {
			return sw
		}
	} else {
		if c.Y < center.Y {
			return ne
		} else {
			return se
		}
	}
}

// split returns a quadrant subdivision of b in the given direction.
func (b Box2) split(dir int) Box2 {
	halfX := (b.Max.X - b.Min.X) / 2
	halfY := (b.Max.Y - b.Min.Y) / 2
	switch dir {
	case ne:
		b.Min.X += halfX
		b.Max.Y -= halfY
	case se:
		b.Min.X += halfX
		b.Min.Y += halfY
	case sw:
		b.Max.X -= halfX
		b.Min.Y += halfY
	case nw:
		b.Max.X -= halfX
		b.Max.Y -= halfY
	}
	return b
}

// Particle2 is a particle in a plane.
type Particle2 interface {
	Coord2() Vector2
	Mass() float64
}

// Force2 is a force modeling function for interactions between p1 and p2,
// m1 is the mass of p1 and m2 of p2. The vector v is the vector from p1 to
// p2. The returned value is the force vector acting on p1.
//
// In models where the identity of particles must be known, p1 and p2 may be
// compared. Force2 may be passed nil for p2 when the Barnes-Hut approximation
// is being used. A nil p2 indicates that the second mass center is an
// aggregate.
type Force2 func(p1, p2 Particle2, m1, m2 float64, v Vector2) Vector2

// Gravity2 returns a vector force on m1 by m2, equal to (m1⋅m2)/‖v‖²
// in the directions of v. Gravity2 ignores the identity of the interacting
// particles and returns a zero vector when the two particles are
// coincident, but performs no other sanity checks.
func Gravity2(_, _ Particle2, m1, m2 float64, v Vector2) Vector2 {
	d2 := v.X*v.X + v.Y*v.Y
	if d2 == 0 {
		return Vector2{}
	}
	return v.Scale((m1 * m2) / (d2 * math.Sqrt(d2)))
}

// Plane implements Barnes-Hut force approximation calculations.
type Plane struct {
	root tile

	Particles []Particle2
}

// NewPlane returns a new Plane.
func NewPlane(p []Particle2) *Plane {
	q := Plane{Particles: p}
	q.Reset()
	return &q
}

// Reset reconstructs the Barnes-Hut tree. Reset must be called if the
// Particles field or elements of Particles have been altered, unless
// ForceOn is called with theta=0 or no data structures have been
// previously built.
func (q *Plane) Reset() {
	if len(q.Particles) == 0 {
		q.root = tile{}
		return
	}

	q.root = tile{
		particle: q.Particles[0],
		center:   q.Particles[0].Coord2(),
		mass:     q.Particles[0].Mass(),
	}
	q.root.bounds.Min = q.root.center
	q.root.bounds.Max = q.root.center
	for _, e := range q.Particles[1:] {
		c := e.Coord2()
		if c.X < q.root.bounds.Min.X {
			q.root.bounds.Min.X = c.X
		}
		if c.X > q.root.bounds.Max.X {
			q.root.bounds.Max.X = c.X
		}
		if c.Y < q.root.bounds.Min.Y {
			q.root.bounds.Min.Y = c.Y
		}
		if c.Y > q.root.bounds.Max.Y {
			q.root.bounds.Max.Y = c.Y
		}
	}

	// TODO(kortschak): Partially parallelise this by
	// choosing the direction and using one of four
	// goroutines to work on each root quadrant.
	for _, e := range q.Particles[1:] {
		q.root.insert(e)
	}
	q.root.summarize()
}

// ForceOn returns a force vector on p given p's mass and the force function, f,
// using the Barnes-Hut theta approximation parameter.
//
// Calls to f will include p in the p1 position and a non-nil p2 if the force
// interaction is with a non-aggregate mass center, otherwise p2 will be nil.
//
// It is safe to call ForceOn concurrently.
func (q *Plane) ForceOn(p Particle2, theta float64, f Force2) (force Vector2) {
	var empty tile
	if theta > 0 && q.root != empty {
		return q.root.forceOn(p, p.Coord2(), p.Mass(), theta, f)
	}

	// For the degenerate case, just iterate over the
	// slice of particles rather than walking the tree.
	var v Vector2
	m := p.Mass()
	pv := p.Coord2()
	for _, e := range q.Particles {
		v = v.Add(f(p, e, m, e.Mass(), e.Coord2().Sub(pv)))
	}
	return v
}

// tile is a quad tree quadrant with Barnes-Hut extensions.
type tile struct {
	particle Particle2

	bounds Box2

	nodes [4]*tile

	center Vector2
	mass   float64
}

// insert inserts p into the subtree rooted at t.
func (t *tile) insert(p Particle2) {
	if t.particle == nil {
		for _, q := range t.nodes {
			if q != nil {
				t.passDown(p)
				return
			}
		}
		t.particle = p
		t.center = p.Coord2()
		t.mass = p.Mass()
		return
	}
	t.passDown(p)
	t.passDown(t.particle)
	t.particle = nil
	t.center = Vector2{}
	t.mass = 0
}

func (t *tile) passDown(p Particle2) {
	dir := t.bounds.quadrant(p)
	if t.nodes[dir] == nil {
		t.nodes[dir] = &tile{bounds: t.bounds.split(dir)}
	}
	t.nodes[dir].insert(p)
}

// summarize updates node masses and centers of mass.
func (t *tile) summarize() (center Vector2, mass float64) {
	for _, d := range &t.nodes {
		if d == nil {
			continue
		}
		c, m := d.summarize()
		t.center.X += c.X * m
		t.center.Y += c.Y * m
		t.mass += m
	}
	t.center.X /= t.mass
	t.center.Y /= t.mass
	return t.center, t.mass
}

// forceOn returns a force vector on p given p's mass m and the force
// calculation function, using the Barnes-Hut theta approximation parameter.
func (t *tile) forceOn(p Particle2, pt Vector2, m, theta float64, f Force2) (vector Vector2) {
	s := ((t.bounds.Max.X - t.bounds.Min.X) + (t.bounds.Max.Y - t.bounds.Min.Y)) / 2
	d := math.Hypot(pt.X-t.center.X, pt.Y-t.center.Y)
	if s/d < theta || t.particle != nil {
		return f(p, t.particle, m, t.mass, t.center.Sub(pt))
	}

	var v Vector2
	for _, d := range &t.nodes {
		if d == nil {
			continue
		}
		v = v.Add(d.forceOn(p, pt, m, theta, f))
	}
	return v
}
