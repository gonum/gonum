// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package community

import (
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

func TestProfileUndirected(t *testing.T) {
	for _, test := range communityUndirectedQTests {
		g := simple.NewUndirectedGraph()
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if g.Node(int64(u)) == nil {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}

		testProfileUndirected(t, test, g)
	}
}

func TestProfileWeightedUndirected(t *testing.T) {
	for _, test := range communityUndirectedQTests {
		g := simple.NewWeightedUndirectedGraph(0, 0)
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if g.Node(int64(u)) == nil {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetWeightedEdge(simple.WeightedEdge{F: simple.Node(u), T: simple.Node(v), W: 1})
			}
		}

		testProfileUndirected(t, test, g)
	}
}

func testProfileUndirected(t *testing.T, test communityUndirectedQTest, g graph.Undirected) {
	fn := ModularScore(g, Weight, 10, nil)
	p, err := Profile(fn, true, 1e-3, 0.1, 10)
	if err != nil {
		t.Errorf("%s: unexpected error: %v", test.name, err)
	}

	const tries = 1000
	for i, d := range p {
		var score float64
		for i := 0; i < tries; i++ {
			score, _ = fn(d.Low)
			if score >= d.Score {
				break
			}
		}
		if score < d.Score {
			t.Errorf("%s: failed to recover low end score: got: %v want: %v", test.name, score, d.Score)
		}
		if i != 0 && d.Score >= p[i-1].Score {
			t.Errorf("%s: not monotonically decreasing: %v -> %v", test.name, p[i-1], d)
		}
	}
}

func TestProfileDirected(t *testing.T) {
	for _, test := range communityDirectedQTests {
		g := simple.NewDirectedGraph()
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if g.Node(int64(u)) == nil {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}

		testProfileDirected(t, test, g)
	}
}

func TestProfileWeightedDirected(t *testing.T) {
	for _, test := range communityDirectedQTests {
		g := simple.NewWeightedDirectedGraph(0, 0)
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if g.Node(int64(u)) == nil {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetWeightedEdge(simple.WeightedEdge{F: simple.Node(u), T: simple.Node(v), W: 1})
			}
		}

		testProfileDirected(t, test, g)
	}
}

func testProfileDirected(t *testing.T, test communityDirectedQTest, g graph.Directed) {
	fn := ModularScore(g, Weight, 10, nil)
	p, err := Profile(fn, true, 1e-3, 0.1, 10)
	if err != nil {
		t.Errorf("%s: unexpected error: %v", test.name, err)
	}

	const tries = 1000
	for i, d := range p {
		var score float64
		for i := 0; i < tries; i++ {
			score, _ = fn(d.Low)
			if score >= d.Score {
				break
			}
		}
		if score < d.Score {
			t.Errorf("%s: failed to recover low end score: got: %v want: %v", test.name, score, d.Score)
		}
		if i != 0 && d.Score >= p[i-1].Score {
			t.Errorf("%s: not monotonically decreasing: %v -> %v", test.name, p[i-1], d)
		}
	}
}

func TestProfileUndirectedMultiplex(t *testing.T) {
	for _, test := range communityUndirectedMultiplexQTests {
		g, weights, err := undirectedMultiplexFrom(test.layers)
		if err != nil {
			t.Errorf("unexpected error creating multiplex: %v", err)
			continue
		}

		const all = true

		fn := ModularMultiplexScore(g, weights, all, WeightMultiplex, 10, nil)
		p, err := Profile(fn, true, 1e-3, 0.1, 10)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", test.name, err)
		}

		const tries = 1000
		for i, d := range p {
			var score float64
			for i := 0; i < tries; i++ {
				score, _ = fn(d.Low)
				if score >= d.Score {
					break
				}
			}
			if score < d.Score {
				t.Errorf("%s: failed to recover low end score: got: %v want: %v", test.name, score, d.Score)
			}
			if i != 0 && d.Score >= p[i-1].Score {
				t.Errorf("%s: not monotonically decreasing: %v -> %v", test.name, p[i-1], d)
			}
		}
	}
}

func TestProfileDirectedMultiplex(t *testing.T) {
	for _, test := range communityDirectedMultiplexQTests {
		g, weights, err := directedMultiplexFrom(test.layers)
		if err != nil {
			t.Errorf("unexpected error creating multiplex: %v", err)
			continue
		}

		const all = true

		fn := ModularMultiplexScore(g, weights, all, WeightMultiplex, 10, nil)
		p, err := Profile(fn, true, 1e-3, 0.1, 10)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", test.name, err)
		}

		const tries = 1000
		for i, d := range p {
			var score float64
			for i := 0; i < tries; i++ {
				score, _ = fn(d.Low)
				if score >= d.Score {
					break
				}
			}
			if score < d.Score {
				t.Errorf("%s: failed to recover low end score: got: %v want: %v", test.name, score, d.Score)
			}
			if i != 0 && d.Score >= p[i-1].Score {
				t.Errorf("%s: not monotonically decreasing: %v -> %v", test.name, p[i-1], d)
			}
		}
	}
}
