// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package barneshut_test

import (
	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/spatial/barneshut"
)

type mass struct {
	d barneshut.Vector2
	v barneshut.Vector2
	m float64
}

func (m *mass) Coord2() barneshut.Vector2 { return m.d }
func (m *mass) Mass() float64             { return m.m }
func (m *mass) move(f barneshut.Vector2) {
	m.v = m.v.Add(f.Scale(1 / m.m))
	m.d = m.d.Add(m.v)
}

func Example_galaxy() {
	rnd := rand.New(rand.NewSource(1))

	// Make 1000 stars in random locations.
	stars := make([]*mass, 1000)
	p := make([]barneshut.Particle2, len(stars))
	for i := range stars {
		s := &mass{
			d: barneshut.Vector2{
				X: 100 * rnd.Float64(),
				Y: 100 * rnd.Float64(),
			},
			v: barneshut.Vector2{
				X: rnd.NormFloat64(),
				Y: rnd.NormFloat64(),
			},
			m: 10 * rnd.Float64(),
		}
		stars[i] = s
		p[i] = s
	}
	vectors := make([]barneshut.Vector2, len(stars))

	// Make a plane to calculate approximate forces
	plane := barneshut.Plane{Particles: p}

	// Run a simulation for 100 updates.
	for i := 0; i < 1000; i++ {
		// Build the data structure. For small systems
		// this step may be omitted and ForceOn will
		// perform the naive quadratic calculation
		// without building the data structure.
		plane.Reset()

		// Calculate the force vectors using the theta
		// parameter...
		const theta = 0.5
		// and an imaginary gravitational constant.
		const G = 10
		for j, s := range stars {
			vectors[j] = plane.ForceOn(s, theta, barneshut.Gravity2).Scale(G)
		}

		// Update positions.
		for j, s := range stars {
			s.move(vectors[j])
		}

		// Rendering stars is left as an exercise for
		// the reader.
	}
}
