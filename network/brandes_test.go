// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network

import (
	"math"
	"testing"

	"github.com/gonum/floats"
	"github.com/gonum/graph/concrete"
)

var betweennessTests = []struct {
	g []set

	wantTol float64
	want    map[int]float64
}{
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

		wantTol: 1e-1,
		want: map[int]float64{
			A: 0,
			B: 32,
			C: 0,
			D: 18,
			E: 48,
			F: 0,
			G: 0,
			H: 0,
			I: 0,
			J: 0,
			K: 0,
		},
	},
	{
		// Example graph from http://en.wikipedia.org/w/index.php?title=PageRank&oldid=659286279#Power_Method
		g: []set{
			A: linksTo(B, C),
			B: linksTo(D),
			C: linksTo(D, E),
			D: linksTo(E),
			E: linksTo(A),
		},

		wantTol: 1e-3,
		want: map[int]float64{
			A: 2,
			B: 0.6667,
			C: 0.6667,
			D: 2,
			E: 0.6667,
		},
	},
	{
		g: []set{
			A: linksTo(B),
			B: linksTo(C),
			C: nil,
		},

		wantTol: 1e-3,
		want: map[int]float64{
			A: 0,
			B: 2,
			C: 0,
		},
	},
	{
		g: []set{
			A: linksTo(B),
			B: linksTo(C),
			C: linksTo(D),
			D: linksTo(E),
			E: nil,
		},

		wantTol: 1e-3,
		want: map[int]float64{
			A: 0,
			B: 6,
			C: 8,
			D: 6,
			E: 0,
		},
	},
	{
		g: []set{
			A: linksTo(C),
			B: linksTo(C),
			C: nil,
			D: linksTo(C),
			E: linksTo(C),
		},

		wantTol: 1e-3,
		want: map[int]float64{
			A: 0,
			B: 0,
			C: 12,
			D: 0,
			E: 0,
		},
	},
	{
		g: []set{
			A: linksTo(B, C, D, E),
			B: linksTo(C, D, E),
			C: linksTo(D, E),
			D: linksTo(E),
			E: nil,
		},

		wantTol: 1e-3,
		want: map[int]float64{
			A: 0,
			B: 0,
			C: 0,
			D: 0,
			E: 0,
		},
	},
}

func TestBetweenness(t *testing.T) {
	for i, test := range betweennessTests {
		g := concrete.NewGraph()
		for u, e := range test.g {
			if !g.NodeExists(concrete.Node(u)) {
				g.AddNode(concrete.Node(u))
			}
			for v := range e {
				if !g.NodeExists(concrete.Node(v)) {
					g.AddNode(concrete.Node(v))
				}
				g.AddUndirectedEdge(concrete.Edge{H: concrete.Node(u), T: concrete.Node(v)}, 0)
			}
		}
		got := Betweenness(g)
		prec := 1 - int(math.Log10(test.wantTol))
		for n := range test.g {
			if !floats.EqualWithinAbsOrRel(got[n], test.want[n], test.wantTol, test.wantTol) {
				t.Errorf("unexpected PageRank result for test %d:\ngot: %v\nwant:%v",
					i, orderedFloats(got, prec), orderedFloats(test.want, prec))
				break
			}
		}
	}
}
