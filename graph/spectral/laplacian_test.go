// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spectral

import (
	"testing"

	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/ordered"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/mat"
)

var randomWalkLaplacianTests = []struct {
	g    []set
	damp float64

	want *mat.Dense
}{
	{
		g: []set{
			A: linksTo(B, C),
			B: linksTo(C),
			C: nil,
		},

		want: mat.NewDense(3, 3, []float64{
			1, 0, 0,
			-0.5, 1, 0,
			-0.5, -1, 0,
		}),
	},
	{
		g: []set{
			A: linksTo(B, C),
			B: linksTo(C),
			C: nil,
		},
		damp: 0.85,

		want: mat.NewDense(3, 3, []float64{
			0.15, 0, 0,
			-0.075, 0.15, 0,
			-0.075, -0.15, 0,
		}),
	},
	{
		g: []set{
			A: linksTo(B),
			B: linksTo(C),
			C: linksTo(A),
		},
		damp: 0.85,

		want: mat.NewDense(3, 3, []float64{
			0.15, 0, -0.15,
			-0.15, 0.15, 0,
			0, -0.15, 0.15,
		}),
	},
	{
		// Example graph from http://en.wikipedia.org/wiki/File:PageRanks-Example.svg 16:17, 8 July 2009
		g: []set{
			A: nil,
			B: linksTo(C),
			C: linksTo(B),
			D: linksTo(A, B),
			E: linksTo(D, B, F),
			F: linksTo(B, E),
			G: linksTo(B, E),
			H: linksTo(B, E),
			I: linksTo(B, E),
			J: linksTo(E),
			K: linksTo(E),
		},

		want: mat.NewDense(11, 11, []float64{
			0, 0, 0, -0.5, 0, 0, 0, 0, 0, 0, 0,
			0, 1, -1, -0.5, -1. / 3., -0.5, -0.5, -0.5, -0.5, 0, 0,
			0, -1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, -1. / 3., 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 1, -0.5, -0.5, -0.5, -0.5, -1, -1,
			0, 0, 0, 0, -1. / 3., 1, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		}),
	},
}

func TestRandomWalkLaplacian(t *testing.T) {
	const tol = 1e-14
	for i, test := range randomWalkLaplacianTests {
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
		l := NewRandomWalkLaplacian(g, test.damp)
		_, c := l.Dims()
		for j := 0; j < c; j++ {
			if got := mat.Sum(l.Matrix.(*mat.Dense).ColView(j)); !scalar.EqualWithinAbsOrRel(got, 0, tol, tol) {
				t.Errorf("unexpected column sum for test %d, column %d: got:%v want:0", i, j, got)
			}
		}
		l = NewRandomWalkLaplacian(sortedNodeGraph{g}, test.damp)
		if !mat.EqualApprox(l, test.want, tol) {
			t.Errorf("unexpected result for test %d:\ngot:\n% .2v\nwant:\n% .2v",
				i, mat.Formatted(l), mat.Formatted(test.want))
		}
	}
}

type sortedNodeGraph struct {
	graph.Graph
}

func (g sortedNodeGraph) Nodes() graph.Nodes {
	n := graph.NodesOf(g.Graph.Nodes())
	ordered.ByID(n)
	return iterator.NewOrderedNodes(n)
}

const (
	A = iota
	B
	C
	D
	E
	F
	G
	H
	I
	J
	K
	L
	M
	N
	O
	P
	Q
	R
	S
	T
	U
	V
	W
	X
	Y
	Z
)

// set is an integer set.
type set map[int64]struct{}

func linksTo(i ...int64) set {
	if len(i) == 0 {
		return nil
	}
	s := make(set)
	for _, v := range i {
		s[v] = struct{}{}
	}
	return s
}
