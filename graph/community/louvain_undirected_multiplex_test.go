// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package community

import (
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/gonum/floats"
	"github.com/gonum/graph"
	"github.com/gonum/graph/internal/ordered"
	"github.com/gonum/graph/simple"
)

var communityUndirectedMultiplexQTests = []struct {
	name       string
	layers     []layer
	structures []structure

	wantLevels []level
}{
	{
		name:   "unconnected",
		layers: []layer{{g: unconnected, weight: 1}},
		structures: []structure{
			{
				resolution: 1,
				memberships: []set{
					0: linksTo(0),
					1: linksTo(1),
					2: linksTo(2),
					3: linksTo(3),
					4: linksTo(4),
					5: linksTo(5),
				},
				want: math.NaN(),
			},
		},
		wantLevels: []level{
			{
				q: math.Inf(-1), // Here math.Inf(-1) is used as a place holder for NaN to allow use of reflect.DeepEqual.
				communities: [][]graph.Node{
					{simple.Node(0)},
					{simple.Node(1)},
					{simple.Node(2)},
					{simple.Node(3)},
					{simple.Node(4)},
					{simple.Node(5)},
				},
			},
		},
	},
	{
		name: "small_dumbell",
		layers: []layer{
			{g: smallDumbell, edgeWeight: 1, weight: 1},
			{g: dumbellRepulsion, edgeWeight: -1, weight: -1},
		},
		structures: []structure{
			{
				resolution: 1,
				memberships: []set{
					0: linksTo(0, 1, 2),
					1: linksTo(3, 4, 5),
				},
				want: 7.0, tol: 1e-10,
			},
			{
				resolution: 1,
				memberships: []set{
					0: linksTo(0, 1, 2, 3, 4, 5),
				},
				want: 0, tol: 1e-14,
			},
		},
		wantLevels: []level{
			{
				q: 7.0,
				communities: [][]graph.Node{
					{simple.Node(0), simple.Node(1), simple.Node(2)},
					{simple.Node(3), simple.Node(4), simple.Node(5)},
				},
			},
			{
				q: -1.4285714285714284,
				communities: [][]graph.Node{
					{simple.Node(0)},
					{simple.Node(1)},
					{simple.Node(2)},
					{simple.Node(3)},
					{simple.Node(4)},
					{simple.Node(5)},
				},
			},
		},
	},
	{
		name: "small_dumbell_twice",
		layers: []layer{
			{g: smallDumbell, weight: 0.5},
			{g: smallDumbell, weight: 0.5},
		},
		structures: []structure{
			{
				resolution: 1,
				memberships: []set{
					0: linksTo(0, 1, 2),
					1: linksTo(3, 4, 5),
				},
				want: 5, tol: 1e-10,
			},
			{
				resolution: 1,
				memberships: []set{
					0: linksTo(0, 1, 2, 3, 4, 5),
				},
				want: 0, tol: 1e-14,
			},
		},
		wantLevels: []level{
			{
				q: 0.35714285714285715 * 14,
				communities: [][]graph.Node{
					{simple.Node(0), simple.Node(1), simple.Node(2)},
					{simple.Node(3), simple.Node(4), simple.Node(5)},
				},
			},
			{
				q: -0.17346938775510204 * 14,
				communities: [][]graph.Node{
					{simple.Node(0)},
					{simple.Node(1)},
					{simple.Node(2)},
					{simple.Node(3)},
					{simple.Node(4)},
					{simple.Node(5)},
				},
			},
		},
	},
	{
		name:   "repulsion",
		layers: []layer{{g: repulsion, edgeWeight: -1, weight: -1}},
		structures: []structure{
			{
				resolution: 1,
				memberships: []set{
					0: linksTo(0, 1, 2),
					1: linksTo(3, 4, 5),
				},
				want: 9.0, tol: 1e-10,
			},
			{
				resolution: 1,
				memberships: []set{
					0: linksTo(0),
					1: linksTo(1),
					2: linksTo(2),
					3: linksTo(3),
					4: linksTo(4),
					5: linksTo(5),
				},
				want: 3, tol: 1e-14,
			},
		},
		wantLevels: []level{
			{
				q: 9.0,
				communities: [][]graph.Node{
					{simple.Node(0), simple.Node(1), simple.Node(2)},
					{simple.Node(3), simple.Node(4), simple.Node(5)},
				},
			},
			{
				q: 3.0,
				communities: [][]graph.Node{
					{simple.Node(0)},
					{simple.Node(1)},
					{simple.Node(2)},
					{simple.Node(3)},
					{simple.Node(4)},
					{simple.Node(5)},
				},
			},
		},
	},
	{
		name: "middle_east",
		layers: []layer{
			{g: middleEast.friends, edgeWeight: 1, weight: 1},
			{g: middleEast.enemies, edgeWeight: -1, weight: -1},
		},
		structures: []structure{
			{
				resolution: 1,
				memberships: []set{
					0: linksTo(0, 6),
					1: linksTo(1, 7, 9, 12),
					2: linksTo(2, 8, 11),
					3: linksTo(3, 4, 5, 10),
				},
				want: 33.8180574555, tol: 1e-9,
			},
			{
				resolution: 1,
				memberships: []set{
					0: linksTo(0, 2, 3, 4, 5, 10),
					1: linksTo(1, 7, 9, 12),
					2: linksTo(6),
					3: linksTo(8, 11),
				},
				want: 30.92749658, tol: 1e-7,
			},
			{
				resolution: 1,
				memberships: []set{
					0: linksTo(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12),
				},
				want: 0, tol: 1e-14,
			},
		},
		wantLevels: []level{
			{
				q: 33.818057455540355,
				communities: [][]graph.Node{
					{simple.Node(0), simple.Node(6)},
					{simple.Node(1), simple.Node(7), simple.Node(9), simple.Node(12)},
					{simple.Node(2), simple.Node(8), simple.Node(11)},
					{simple.Node(3), simple.Node(4), simple.Node(5), simple.Node(10)},
				},
			},
			{
				q: 3.8071135430916545,
				communities: [][]graph.Node{
					{simple.Node(0)},
					{simple.Node(1)},
					{simple.Node(2)},
					{simple.Node(3)},
					{simple.Node(4)},
					{simple.Node(5)},
					{simple.Node(6)},
					{simple.Node(7)},
					{simple.Node(8)},
					{simple.Node(9)},
					{simple.Node(10)},
					{simple.Node(11)},
					{simple.Node(12)},
				},
			},
		},
	},
}

func TestCommunityQUndirectedMultiplex(t *testing.T) {
	for _, test := range communityUndirectedMultiplexQTests {
		g, weights, err := undirectedMultiplexFrom(test.layers)
		if err != nil {
			t.Errorf("unexpected error creating multiplex: %v", err)
			continue
		}

		for _, structure := range test.structures {
			communities := make([][]graph.Node, len(structure.memberships))
			for i, c := range structure.memberships {
				for n := range c {
					communities[i] = append(communities[i], simple.Node(n))
				}
			}
			q := QMultiplex(g, communities, weights, []float64{structure.resolution})
			got := floats.Sum(q)
			if !floats.EqualWithinAbsOrRel(got, structure.want, structure.tol, structure.tol) && !math.IsNaN(structure.want) {
				for _, c := range communities {
					sort.Sort(ordered.ByID(c))
				}
				t.Errorf("unexpected Q value for %q %v: got: %v %.3v want: %v",
					test.name, communities, got, q, structure.want)
			}
		}
	}
}

func TestCommunityDeltaQUndirectedMultiplex(t *testing.T) {
tests:
	for _, test := range communityUndirectedMultiplexQTests {
		g, weights, err := undirectedMultiplexFrom(test.layers)
		if err != nil {
			t.Errorf("unexpected error creating multiplex: %v", err)
			continue
		}

		rnd := rand.New(rand.NewSource(1)).Intn
		for _, structure := range test.structures {
			communityOf := make(map[int]int)
			communities := make([][]graph.Node, len(structure.memberships))
			for i, c := range structure.memberships {
				for n := range c {
					communityOf[n] = i
					communities[i] = append(communities[i], simple.Node(n))
				}
				sort.Sort(ordered.ByID(communities[i]))
			}
			resolution := []float64{structure.resolution}

			before := QMultiplex(g, communities, weights, resolution)

			// We test exhaustively.
			const all = true

			l := newUndirectedMultiplexLocalMover(
				reduceUndirectedMultiplex(g, nil, weights),
				communities, weights, resolution, all)
			if l == nil {
				if !math.IsNaN(floats.Sum(before)) {
					t.Errorf("unexpected nil localMover with non-NaN Q graph: Q=%.4v", before)
				}
				continue tests
			}

			// This is done to avoid run-to-run
			// variation due to map iteration order.
			sort.Sort(ordered.ByID(l.nodes))

			l.shuffle(rnd)

			for _, target := range l.nodes {
				got, gotDst, gotSrc := l.deltaQ(target)

				want, wantDst := math.Inf(-1), -1
				migrated := make([][]graph.Node, len(structure.memberships))
				for i, c := range structure.memberships {
					for n := range c {
						if n == target.ID() {
							continue
						}
						migrated[i] = append(migrated[i], simple.Node(n))
					}
					sort.Sort(ordered.ByID(migrated[i]))
				}

				for i, c := range structure.memberships {
					if i == communityOf[target.ID()] {
						continue
					}
					if !(all && hasNegative(weights)) {
						connected := false
					search:
						for l := 0; l < g.Depth(); l++ {
							if weights[l] < 0 {
								connected = true
								break search
							}
							layer := g.Layer(l)
							for n := range c {
								if layer.HasEdgeBetween(simple.Node(n), target) {
									connected = true
									break search
								}
							}
						}
						if !connected {
							continue
						}
					}
					migrated[i] = append(migrated[i], target)
					after := QMultiplex(g, migrated, weights, resolution)
					migrated[i] = migrated[i][:len(migrated[i])-1]
					if delta := floats.Sum(after) - floats.Sum(before); delta > want {
						want = delta
						wantDst = i
					}
				}

				if !floats.EqualWithinAbsOrRel(got, want, structure.tol, structure.tol) || gotDst != wantDst {
					t.Errorf("unexpected result moving n=%d in c=%d of %s/%.4v: got: %.4v,%d want: %.4v,%d"+
						"\n\t%v\n\t%v",
						target.ID(), communityOf[target.ID()], test.name, structure.resolution, got, gotDst, want, wantDst,
						communities, migrated)
				}
				if gotSrc.community != communityOf[target.ID()] {
					t.Errorf("unexpected source community index: got: %d want: %d", gotSrc, communityOf[target.ID()])
				} else if communities[gotSrc.community][gotSrc.node].ID() != target.ID() {
					wantNodeIdx := -1
					for i, n := range communities[gotSrc.community] {
						if n.ID() == target.ID() {
							wantNodeIdx = i
							break
						}
					}
					t.Errorf("unexpected source node index: got: %d want: %d", gotSrc.node, wantNodeIdx)
				}
			}
		}
	}
}

func TestReduceQConsistencyUndirectedMultiplex(t *testing.T) {
tests:
	for _, test := range communityUndirectedMultiplexQTests {
		g, weights, err := undirectedMultiplexFrom(test.layers)
		if err != nil {
			t.Errorf("unexpected error creating multiplex: %v", err)
			continue
		}

		for _, structure := range test.structures {
			if math.IsNaN(structure.want) {
				continue tests
			}

			communities := make([][]graph.Node, len(structure.memberships))
			for i, c := range structure.memberships {
				for n := range c {
					communities[i] = append(communities[i], simple.Node(n))
				}
				sort.Sort(ordered.ByID(communities[i]))
			}

			gQ := QMultiplex(g, communities, weights, []float64{structure.resolution})
			gQnull := QMultiplex(g, nil, weights, nil)

			cg0 := reduceUndirectedMultiplex(g, nil, weights)
			cg0Qnull := QMultiplex(cg0, cg0.Structure(), weights, nil)
			if !floats.EqualWithinAbsOrRel(floats.Sum(gQnull), floats.Sum(cg0Qnull), structure.tol, structure.tol) {
				t.Errorf("disagreement between null Q from method: %v and function: %v", cg0Qnull, gQnull)
			}
			cg0Q := QMultiplex(cg0, communities, weights, []float64{structure.resolution})
			if !floats.EqualWithinAbsOrRel(floats.Sum(gQ), floats.Sum(cg0Q), structure.tol, structure.tol) {
				t.Errorf("unexpected Q result after initial reduction: got: %v want :%v", cg0Q, gQ)
			}

			cg1 := reduceUndirectedMultiplex(cg0, communities, weights)
			cg1Q := QMultiplex(cg1, cg1.Structure(), weights, []float64{structure.resolution})
			if !floats.EqualWithinAbsOrRel(floats.Sum(gQ), floats.Sum(cg1Q), structure.tol, structure.tol) {
				t.Errorf("unexpected Q result after second reduction: got: %v want :%v", cg1Q, gQ)
			}
		}
	}
}

var localUndirectedMultiplexMoveTests = []struct {
	name       string
	layers     []layer
	structures []moveStructures
}{
	{
		name:   "blondel",
		layers: []layer{{g: blondel, weight: 1}, {g: blondel, weight: 0.5}},
		structures: []moveStructures{
			{
				memberships: []set{
					0: linksTo(0, 1, 2, 4, 5),
					1: linksTo(3, 6, 7),
					2: linksTo(8, 9, 10, 12, 14, 15),
					3: linksTo(11, 13),
				},
				targetNodes: []graph.Node{simple.Node(0)},
				resolution:  1,
				tol:         1e-14,
			},
			{
				memberships: []set{
					0: linksTo(0, 1, 2, 4, 5),
					1: linksTo(3, 6, 7),
					2: linksTo(8, 9, 10, 12, 14, 15),
					3: linksTo(11, 13),
				},
				targetNodes: []graph.Node{simple.Node(3)},
				resolution:  1,
				tol:         1e-14,
			},
			{
				memberships: []set{
					0: linksTo(0, 1, 2, 4, 5),
					1: linksTo(3, 6, 7),
					2: linksTo(8, 9, 10, 12, 14, 15),
					3: linksTo(11, 13),
				},
				// Case to demonstrate when A_aa != k_a^𝛼.
				targetNodes: []graph.Node{simple.Node(3), simple.Node(2)},
				resolution:  1,
				tol:         1e-14,
			},
		},
	},
}

func TestMoveLocalUndirectedMultiplex(t *testing.T) {
	for _, test := range localUndirectedMultiplexMoveTests {
		g, weights, err := undirectedMultiplexFrom(test.layers)
		if err != nil {
			t.Errorf("unexpected error creating multiplex: %v", err)
			continue
		}

		for _, structure := range test.structures {
			communities := make([][]graph.Node, len(structure.memberships))
			for i, c := range structure.memberships {
				for n := range c {
					communities[i] = append(communities[i], simple.Node(n))
				}
				sort.Sort(ordered.ByID(communities[i]))
			}

			r := reduceUndirectedMultiplex(reduceUndirectedMultiplex(g, nil, weights), communities, weights)

			l := newUndirectedMultiplexLocalMover(r, r.communities, weights, []float64{structure.resolution}, true)
			for _, n := range structure.targetNodes {
				dQ, dst, src := l.deltaQ(n)
				if dQ > 0 {
					before := floats.Sum(QMultiplex(r, l.communities, weights, []float64{structure.resolution}))
					l.move(dst, src)
					after := floats.Sum(QMultiplex(r, l.communities, weights, []float64{structure.resolution}))
					want := after - before
					if !floats.EqualWithinAbsOrRel(dQ, want, structure.tol, structure.tol) {
						t.Errorf("unexpected deltaQ: got: %v want: %v", dQ, want)
					}
				}
			}
		}
	}
}

func TestLouvainMultiplex(t *testing.T) {
	const louvainIterations = 20

	for _, test := range communityUndirectedMultiplexQTests {
		g, weights, err := undirectedMultiplexFrom(test.layers)
		if err != nil {
			t.Errorf("unexpected error creating multiplex: %v", err)
			continue
		}

		if test.structures[0].resolution != 1 {
			panic("bad test: expect resolution=1")
		}
		want := make([][]graph.Node, len(test.structures[0].memberships))
		for i, c := range test.structures[0].memberships {
			for n := range c {
				want[i] = append(want[i], simple.Node(n))
			}
			sort.Sort(ordered.ByID(want[i]))
		}
		sort.Sort(ordered.BySliceIDs(want))

		var (
			got   *ReducedUndirectedMultiplex
			bestQ = math.Inf(-1)
		)
		// Modularize is randomised so we do this to
		// ensure the level tests are consistent.
		src := rand.New(rand.NewSource(1))
		for i := 0; i < louvainIterations; i++ {
			r := ModularizeMultiplex(g, weights, nil, true, src).(*ReducedUndirectedMultiplex)
			if q := floats.Sum(QMultiplex(r, nil, weights, nil)); q > bestQ || math.IsNaN(q) {
				bestQ = q
				got = r

				if math.IsNaN(q) {
					// Don't try again for non-connected case.
					break
				}
			}

			var qs []float64
			for p := r; p != nil; p = p.Expanded().(*ReducedUndirectedMultiplex) {
				qs = append(qs, floats.Sum(QMultiplex(p, nil, weights, nil)))
			}

			// Recovery of Q values is reversed.
			if reverse(qs); !sort.Float64sAreSorted(qs) {
				t.Errorf("Q values not monotonically increasing: %.5v", qs)
			}
		}

		gotCommunities := got.Communities()
		for _, c := range gotCommunities {
			sort.Sort(ordered.ByID(c))
		}
		sort.Sort(ordered.BySliceIDs(gotCommunities))
		if !reflect.DeepEqual(gotCommunities, want) {
			t.Errorf("unexpected community membership for %s Q=%.4v:\n\tgot: %v\n\twant:%v",
				test.name, bestQ, gotCommunities, want)
			continue
		}

		var levels []level
		for p := got; p != nil; p = p.Expanded().(*ReducedUndirectedMultiplex) {
			var communities [][]graph.Node
			if p.parent != nil {
				communities = p.parent.Communities()
				for _, c := range communities {
					sort.Sort(ordered.ByID(c))
				}
				sort.Sort(ordered.BySliceIDs(communities))
			} else {
				communities = reduceUndirectedMultiplex(g, nil, weights).Communities()
			}
			q := floats.Sum(QMultiplex(p, nil, weights, nil))
			if math.IsNaN(q) {
				// Use an equalable flag value in place of NaN.
				q = math.Inf(-1)
			}
			levels = append(levels, level{q: q, communities: communities})
		}
		if !reflect.DeepEqual(levels, test.wantLevels) {
			t.Errorf("unexpected level structure:\n\tgot: %v\n\twant:%v", levels, test.wantLevels)
		}
	}
}

func TestNonContiguousUndirectedMultiplex(t *testing.T) {
	g := simple.NewUndirectedGraph(0, 0)
	for _, e := range []simple.Edge{
		{F: simple.Node(0), T: simple.Node(1), W: 1},
		{F: simple.Node(4), T: simple.Node(5), W: 1},
	} {
		g.SetEdge(e)
	}

	func() {
		defer func() {
			r := recover()
			if r != nil {
				t.Error("unexpected panic with non-contiguous ID range")
			}
		}()
		ModularizeMultiplex(UndirectedLayers{g}, nil, nil, true, nil)
	}()
}

func BenchmarkLouvainMultiplex(b *testing.B) {
	src := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		ModularizeMultiplex(UndirectedLayers{dupGraph}, nil, nil, true, src)
	}
}

func undirectedMultiplexFrom(raw []layer) (UndirectedLayers, []float64, error) {
	var layers []graph.Undirected
	var weights []float64
	for _, l := range raw {
		g := simple.NewUndirectedGraph(0, 0)
		for u, e := range l.g {
			// Add nodes that are not defined by an edge.
			if !g.Has(simple.Node(u)) {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				w := 1.0
				if l.edgeWeight != 0 {
					w = l.edgeWeight
				}
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v), W: w})
			}
		}
		layers = append(layers, g)
		weights = append(weights, l.weight)
	}
	g, err := NewUndirectedLayers(layers...)
	if err != nil {
		return nil, nil, err
	}
	return g, weights, nil
}
