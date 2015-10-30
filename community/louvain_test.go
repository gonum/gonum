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
	"github.com/gonum/graph/graphs/gen"
	"github.com/gonum/graph/internal/ordered"
	"github.com/gonum/graph/simple"
)

// set is an integer set.
type set map[int]struct{}

func linksTo(i ...int) set {
	if len(i) == 0 {
		return nil
	}
	s := make(set)
	for _, v := range i {
		s[v] = struct{}{}
	}
	return s
}

type structure struct {
	resolution  float64
	memberships []set
	want, tol   float64
}

type level struct {
	q           float64
	communities [][]graph.Node
}

var communityQTests = []struct {
	name       string
	g          []set
	structures []structure

	wantLevels []level
}{
	// The java reference implementation is available from http://www.ludowaltman.nl/slm/.
	{
		name: "unconnected",
		g: []set{
			/* Nodes 0-4 are implicit .*/ 5: nil,
		},
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
		g: []set{
			0: linksTo(1, 2),
			1: linksTo(2),
			2: linksTo(3),
			3: linksTo(4, 5),
			4: linksTo(5),
			5: nil,
		},
		structures: []structure{
			{
				resolution: 1,
				// community structure and modularity calculated by java reference implementation.
				memberships: []set{
					0: linksTo(0, 1, 2),
					1: linksTo(3, 4, 5),
				},
				want: 0.357, tol: 1e-3,
			},
			{
				resolution: 1,
				memberships: []set{
					0: linksTo(0, 1, 2, 3, 4, 5),
				},
				// theoretical expectation.
				want: 0, tol: 1e-14,
			},
		},
		wantLevels: []level{
			{
				q: 0.35714285714285715,
				communities: [][]graph.Node{
					{simple.Node(0), simple.Node(1), simple.Node(2)},
					{simple.Node(3), simple.Node(4), simple.Node(5)},
				},
			},
			{
				q: -0.17346938775510204,
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
		// W. W. Zachary, An information flow model for conflict and fission in small groups,
		// Journal of Anthropological Research 33, 452-473 (1977).
		name: "zachary",
		g: []set{
			0:  linksTo(1, 2, 3, 4, 5, 6, 7, 8, 10, 11, 12, 13, 17, 19, 21, 31),
			1:  linksTo(2, 3, 7, 13, 17, 19, 21, 30),
			2:  linksTo(3, 7, 8, 9, 13, 27, 28, 32),
			3:  linksTo(7, 12, 13),
			4:  linksTo(6, 10),
			5:  linksTo(6, 10, 16),
			6:  linksTo(16),
			8:  linksTo(30, 32, 33),
			9:  linksTo(33),
			13: linksTo(33),
			14: linksTo(32, 33),
			15: linksTo(32, 33),
			18: linksTo(32, 33),
			19: linksTo(33),
			20: linksTo(32, 33),
			22: linksTo(32, 33),
			23: linksTo(25, 27, 29, 32, 33),
			24: linksTo(25, 27, 31),
			25: linksTo(31),
			26: linksTo(29, 33),
			27: linksTo(33),
			28: linksTo(31, 33),
			29: linksTo(32, 33),
			30: linksTo(32, 33),
			31: linksTo(32, 33),
			32: linksTo(33),
			33: nil,
		},
		structures: []structure{
			{
				resolution: 1,
				// community structure and modularity from doi: 10.1140/epjb/e2013-40829-0
				memberships: []set{
					0: linksTo(0, 1, 2, 3, 7, 11, 12, 13, 17, 19, 21),
					1: linksTo(4, 5, 6, 10, 16),
					2: linksTo(8, 9, 14, 15, 18, 20, 22, 26, 29, 30, 32, 33),
					3: linksTo(23, 24, 25, 27, 28, 31),
				},
				// Noted to be the optimal modularisation in the paper above.
				want: 0.4198, tol: 1e-4,
			},
			{
				resolution: 0.5,
				// community structure and modularity calculated by java reference implementation.
				memberships: []set{
					0: linksTo(0, 1, 2, 3, 4, 5, 6, 7, 9, 10, 11, 12, 13, 16, 17, 19, 21),
					1: linksTo(8, 14, 15, 18, 20, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33),
				},
				want: 0.6218, tol: 1e-3,
			},
			{
				resolution: 2,
				// community structure and modularity calculated by java reference implementation.
				memberships: []set{
					0: linksTo(14, 18, 20, 22, 32, 33, 15),
					1: linksTo(0, 1, 11, 17, 19, 21),
					2: linksTo(2, 3, 7, 9, 12, 13),
					3: linksTo(4, 5, 6, 10, 16),
					4: linksTo(24, 25, 28, 31),
					5: linksTo(23, 26, 27, 29),
					6: linksTo(8, 30),
				},
				want: 0.1645, tol: 1e-3,
			},
		},
		wantLevels: []level{
			{
				q: 0.4197896120973044,
				communities: [][]graph.Node{
					{simple.Node(0), simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(7), simple.Node(11), simple.Node(12), simple.Node(13), simple.Node(17), simple.Node(19), simple.Node(21)},
					{simple.Node(4), simple.Node(5), simple.Node(6), simple.Node(10), simple.Node(16)},
					{simple.Node(8), simple.Node(9), simple.Node(14), simple.Node(15), simple.Node(18), simple.Node(20), simple.Node(22), simple.Node(26), simple.Node(29), simple.Node(30), simple.Node(32), simple.Node(33)},
					{simple.Node(23), simple.Node(24), simple.Node(25), simple.Node(27), simple.Node(28), simple.Node(31)},
				},
			},
			{
				q: 0.39907955292570674,
				communities: [][]graph.Node{
					{simple.Node(0), simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(7), simple.Node(11), simple.Node(12), simple.Node(13), simple.Node(17), simple.Node(19), simple.Node(21)},
					{simple.Node(4), simple.Node(10)},
					{simple.Node(5), simple.Node(6), simple.Node(16)},
					{simple.Node(8), simple.Node(9), simple.Node(14), simple.Node(15), simple.Node(18), simple.Node(20), simple.Node(22), simple.Node(26), simple.Node(29), simple.Node(30), simple.Node(32), simple.Node(33)},
					{simple.Node(23), simple.Node(24), simple.Node(25), simple.Node(27), simple.Node(28), simple.Node(31)},
				},
			},
			{
				q: -0.04980276134122286,
				communities: [][]graph.Node{
					[]graph.Node{simple.Node(0)},
					[]graph.Node{simple.Node(1)},
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
					{simple.Node(13)},
					{simple.Node(14)},
					{simple.Node(15)},
					{simple.Node(16)},
					{simple.Node(17)},
					{simple.Node(18)},
					{simple.Node(19)},
					{simple.Node(20)},
					{simple.Node(21)},
					{simple.Node(22)},
					{simple.Node(23)},
					{simple.Node(24)},
					{simple.Node(25)},
					{simple.Node(26)},
					{simple.Node(27)},
					{simple.Node(28)},
					{simple.Node(29)},
					{simple.Node(30)},
					{simple.Node(31)},
					{simple.Node(32)},
					{simple.Node(33)},
				},
			},
		},
	},
	{
		// doi:10.1088/1742-5468/2008/10/P10008 figure 1
		name: "blondel",
		g: []set{
			0:  linksTo(2, 3, 4, 5),
			1:  linksTo(2, 4, 7),
			2:  linksTo(4, 5, 6),
			3:  linksTo(7),
			4:  linksTo(10),
			5:  linksTo(7, 11),
			6:  linksTo(7, 11),
			8:  linksTo(9, 10, 11, 14, 15),
			9:  linksTo(12, 14),
			10: linksTo(11, 12, 13, 14),
			11: linksTo(13),
			15: nil,
		},
		structures: []structure{
			{
				resolution: 1,
				// community structure and modularity calculated by java reference implementation.
				memberships: []set{
					0: linksTo(0, 1, 2, 3, 4, 5, 6, 7),
					1: linksTo(8, 9, 10, 11, 12, 13, 14, 15),
				},
				want: 0.3922, tol: 1e-4,
			},
		},
		wantLevels: []level{
			{
				q: 0.39221938775510207,
				communities: [][]graph.Node{
					[]graph.Node{simple.Node(0), simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(4), simple.Node(5), simple.Node(6), simple.Node(7)},
					[]graph.Node{simple.Node(8), simple.Node(9), simple.Node(10), simple.Node(11), simple.Node(12), simple.Node(13), simple.Node(14), simple.Node(15)},
				},
			},
			{
				q: 0.34630102040816324,
				communities: [][]graph.Node{
					{simple.Node(0), simple.Node(1), simple.Node(2), simple.Node(4), simple.Node(5)},
					{simple.Node(3), simple.Node(6), simple.Node(7)},
					{simple.Node(8), simple.Node(9), simple.Node(10), simple.Node(12), simple.Node(14), simple.Node(15)},
					{simple.Node(11), simple.Node(13)},
				},
			},
			{
				q: -0.07142857142857144,
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
					{simple.Node(13)},
					{simple.Node(14)},
					{simple.Node(15)},
				},
			},
		},
	},
}

func TestCommunityQ(t *testing.T) {
	for _, test := range communityQTests {
		g := simple.NewUndirectedGraph(0, 0)
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if !g.Has(simple.Node(u)) {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v), W: 1})
			}
		}
		for _, structure := range test.structures {
			communities := make([][]graph.Node, len(structure.memberships))
			for i, c := range structure.memberships {
				for n := range c {
					communities[i] = append(communities[i], simple.Node(n))
				}
			}
			got := Q(g, communities, structure.resolution)
			if !floats.EqualWithinAbsOrRel(got, structure.want, structure.tol, structure.tol) && math.IsNaN(got) != math.IsNaN(structure.want) {
				for _, c := range communities {
					sort.Sort(ordered.ByID(c))
				}
				t.Errorf("unexpected Q value for %q %v: got: %v want: %v",
					test.name, communities, got, structure.want)
			}
		}
	}
}

func TestCommunityDeltaQ(t *testing.T) {
tests:
	for _, test := range communityQTests {
		g := simple.NewUndirectedGraph(0, 0)
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if !g.Has(simple.Node(u)) {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v), W: 1})
			}
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

			before := Q(g, communities, structure.resolution)

			l := newLocalMover(reduce(g, nil), communities, structure.resolution)
			if l == nil {
				if !math.IsNaN(before) {
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
					connected := false
					for n := range c {
						if g.HasEdgeBetween(simple.Node(n), target) {
							connected = true
							break
						}
					}
					if !connected {
						continue
					}
					migrated[i] = append(migrated[i], target)
					after := Q(g, migrated, structure.resolution)
					migrated[i] = migrated[i][:len(migrated[i])-1]
					if after-before > want {
						want = after - before
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

func TestReduceQConsistency(t *testing.T) {
tests:
	for _, test := range communityQTests {
		g := simple.NewUndirectedGraph(0, 0)
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if !g.Has(simple.Node(u)) {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v), W: 1})
			}
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

			gQ := Q(g, communities, structure.resolution)
			gQnull := Q(g, nil, 1)

			cg0 := reduce(g, nil)
			cg0Qnull := Q(cg0, cg0.Structure(), 1)
			if !floats.EqualWithinAbsOrRel(gQnull, cg0Qnull, structure.tol, structure.tol) {
				t.Errorf("disgagreement between null Q from method: %v and function: %v", cg0Qnull, gQnull)
			}
			cg0Q := Q(cg0, communities, structure.resolution)
			if !floats.EqualWithinAbsOrRel(gQ, cg0Q, structure.tol, structure.tol) {
				t.Errorf("unexpected Q result after initial conversion: got: %v want :%v", gQ, cg0Q)
			}

			cg1 := reduce(cg0, communities)
			cg1Q := Q(cg1, cg1.Structure(), structure.resolution)
			if !floats.EqualWithinAbsOrRel(gQ, cg1Q, structure.tol, structure.tol) {
				t.Errorf("unexpected Q result after initial condensation: got: %v want :%v", gQ, cg1Q)
			}
		}
	}
}

type moveStructures struct {
	memberships []set
	targetNodes []graph.Node

	resolution float64
	tol        float64
}

var localMoveTests = []struct {
	name       string
	g          []set
	structures []moveStructures
}{
	{
		// doi:10.1088/1742-5468/2008/10/P10008 figure 1
		name: "blondel",
		g: []set{
			0:  linksTo(2, 3, 4, 5),
			1:  linksTo(2, 4, 7),
			2:  linksTo(4, 5, 6),
			3:  linksTo(7),
			4:  linksTo(10),
			5:  linksTo(7, 11),
			6:  linksTo(7, 11),
			8:  linksTo(9, 10, 11, 14, 15),
			9:  linksTo(12, 14),
			10: linksTo(11, 12, 13, 14),
			11: linksTo(13),
			15: nil,
		},
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

func TestMoveLocal(t *testing.T) {
	for _, test := range localMoveTests {
		g := simple.NewUndirectedGraph(0, 0)
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if !g.Has(simple.Node(u)) {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v), W: 1})
			}
		}

		for _, structure := range test.structures {
			communities := make([][]graph.Node, len(structure.memberships))
			for i, c := range structure.memberships {
				for n := range c {
					communities[i] = append(communities[i], simple.Node(n))
				}
				sort.Sort(ordered.ByID(communities[i]))
			}

			r := reduce(reduce(g, nil), communities)

			l := newLocalMover(r, r.communities, structure.resolution)
			for _, n := range structure.targetNodes {
				dQ, dst, src := l.deltaQ(n)
				if dQ > 0 {
					before := Q(r, l.communities, structure.resolution)
					l.move(dst, src)
					after := Q(r, l.communities, structure.resolution)
					want := after - before
					if !floats.EqualWithinAbsOrRel(dQ, want, structure.tol, structure.tol) {
						t.Errorf("unexpected deltaQ: got: %v want: %v", dQ, want)
					}
				}
			}
		}
	}
}

func TestLouvain(t *testing.T) {
	const louvainIterations = 20

	for _, test := range communityQTests {
		g := simple.NewUndirectedGraph(0, 0)
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if !g.Has(simple.Node(u)) {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v), W: 1})
			}
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
			got   *ReducedUndirected
			bestQ = math.Inf(-1)
		)
		// Louvain is randomised so we do this to
		// ensure the level tests are consistent.
		src := rand.New(rand.NewSource(1))
		for i := 0; i < louvainIterations; i++ {
			r := Louvain(g, 1, src)
			if q := Q(r, nil, 1); q > bestQ || math.IsNaN(q) {
				bestQ = q
				got = r

				if math.IsNaN(q) {
					// Don't try again for non-connected case.
					break
				}
			}

			var qs []float64
			for p := r; p != nil; p = p.Expanded() {
				qs = append(qs, Q(p, nil, 1))
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
		for p := got; p != nil; p = p.Expanded() {
			var communities [][]graph.Node
			if p.parent != nil {
				communities = p.parent.Communities()
				for _, c := range communities {
					sort.Sort(ordered.ByID(c))
				}
				sort.Sort(ordered.BySliceIDs(communities))
			} else {
				communities = reduce(g, nil).Communities()
			}
			q := Q(p, nil, 1)
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

func reverse(f []float64) {
	for i, j := 0, len(f)-1; i < j; i, j = i+1, j-1 {
		f[i], f[j] = f[j], f[i]
	}
}

func TestNonContiguous(t *testing.T) {
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
		Louvain(g, 1, nil)
	}()
}

var dupGraph = simple.NewUndirectedGraph(0, 0)

func init() {
	err := gen.Duplication(dupGraph, 1000, 0.8, 0.1, 0.5, rand.New(rand.NewSource(1)))
	if err != nil {
		panic(err)
	}
}

func BenchmarkLouvain(b *testing.B) {
	src := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		Louvain(dupGraph, 1, src)
	}
}
