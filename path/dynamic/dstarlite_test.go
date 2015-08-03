// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dynamic

import (
	"bytes"
	"flag"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/path"
	"github.com/gonum/graph/path/internal"
)

var (
	debug   = flag.Bool("debug", false, "write path progress for failing dynamic case tests")
	vdebug  = flag.Bool("vdebug", false, "write path progress for all dynamic case tests (requires test.v)")
	maxWide = flag.Int("maxwidth", 5, "maximum width grid to dump for debugging")
)

func TestDStarLiteNullHeuristic(t *testing.T) {
	for _, test := range shortestPathTests {
		// Skip zero-weight cycles.
		if strings.HasPrefix(test.name, "zero-weight") {
			continue
		}

		g := test.g()
		for _, e := range test.edges {
			g.SetEdge(e, e.Weight())
		}

		var (
			d *DStarLite

			panicked bool
		)
		func() {
			defer func() {
				panicked = recover() != nil
			}()
			d = NewDStarLite(test.query.From(), test.query.To(), g.(graph.Graph), path.NullHeuristic, concrete.NewDirectedGraph())
		}()
		if panicked || test.negative {
			if !test.negative {
				t.Errorf("%q: unexpected panic", test.name)
			}
			if !panicked {
				t.Errorf("%q: expected panic for negative edge weight", test.name)
			}
			continue
		}

		p, weight := d.Path()

		if !math.IsInf(weight, 1) && p[0].ID() != test.query.From().ID() {
			t.Fatalf("%q: unexpected from node ID: got:%d want:%d", p[0].ID(), test.query.From().ID())
		}
		if weight != test.weight {
			t.Errorf("%q: unexpected weight from Between: got:%f want:%f",
				test.name, weight, test.weight)
		}

		var got []int
		for _, n := range p {
			got = append(got, n.ID())
		}
		ok := len(got) == 0 && len(test.want) == 0
		for _, sp := range test.want {
			if reflect.DeepEqual(got, sp) {
				ok = true
				break
			}
		}
		if !ok {
			t.Errorf("%q: unexpected shortest path:\ngot: %v\nwant from:%v",
				test.name, p, test.want)
		}
	}
}

var dynamicDStarLiteTests = []struct {
	g          *internal.Grid
	radius     float64
	all        bool
	diag, unit bool
	remember   []bool
	modify     func(*internal.LimitedVisionGrid)

	heuristic func(dx, dy float64) float64

	s, t graph.Node

	want        []graph.Node
	weight      float64
	wantedPaths map[int][]graph.Node
}{
	{
		// This is the example shown in figures 6 and 7 of doi:10.1109/tro.2004.838026.
		g: internal.NewGridFrom(
			"...",
			".*.",
			".*.",
			".*.",
			"...",
		),
		radius:   1.5,
		all:      true,
		diag:     true,
		unit:     true,
		remember: []bool{false, true},

		heuristic: func(dx, dy float64) float64 {
			return math.Max(math.Abs(dx), math.Abs(dy))
		},

		s: concrete.Node(3),
		t: concrete.Node(14),

		want: []graph.Node{
			concrete.Node(3),
			concrete.Node(6),
			concrete.Node(9),
			concrete.Node(13),
			concrete.Node(14),
		},
		weight: 4,
	},
	{
		// This is a small example that has the property that the first corner
		// may be taken incorrectly at 90° or correctly at 45° because the
		// calculated rhs values of 12 and 17 are tied when moving from node
		// 16, and the grid is small enough to examine by a dump.
		g: internal.NewGridFrom(
			".....",
			"...*.",
			"**.*.",
			"...*.",
		),
		radius:   1.5,
		all:      true,
		diag:     true,
		remember: []bool{false, true},

		heuristic: func(dx, dy float64) float64 {
			return math.Max(math.Abs(dx), math.Abs(dy))
		},

		s: concrete.Node(15),
		t: concrete.Node(14),

		want: []graph.Node{
			concrete.Node(15),
			concrete.Node(16),
			concrete.Node(12),
			concrete.Node(7),
			concrete.Node(3),
			concrete.Node(9),
			concrete.Node(14),
		},
		weight: 7.242640687119285,
		wantedPaths: map[int][]graph.Node{
			12: []graph.Node{concrete.Node(12), concrete.Node(7), concrete.Node(3), concrete.Node(9), concrete.Node(14)},
		},
	},
	{
		// This is the example shown in figure 2 of doi:10.1109/tro.2004.838026
		// with the exception that diagonal edge weights are calculated with the hypot
		// function instead of a step count and only allowing information to be known
		// from exploration.
		g: internal.NewGridFrom(
			"..................",
			"..................",
			"..................",
			"..................",
			"..................",
			"..................",
			"....*.*...........",
			"*****.***.........",
			"......*...........",
			"......***.........",
			"......*...........",
			"......*...........",
			"......*...........",
			"*****.*...........",
			"......*...........",
		),
		radius:   1.5,
		all:      true,
		diag:     true,
		remember: []bool{false, true},

		heuristic: func(dx, dy float64) float64 {
			return math.Max(math.Abs(dx), math.Abs(dy))
		},

		s: concrete.Node(253),
		t: concrete.Node(122),

		want: []graph.Node{
			concrete.Node(253),
			concrete.Node(254),
			concrete.Node(255),
			concrete.Node(256),
			concrete.Node(239),
			concrete.Node(221),
			concrete.Node(203),
			concrete.Node(185),
			concrete.Node(167),
			concrete.Node(149),
			concrete.Node(131),
			concrete.Node(113),
			concrete.Node(96),

			// The following section depends
			// on map iteration order.
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,

			concrete.Node(122),
		},
		weight: 21.242640687119287,
	},
	{
		// This is the example shown in figure 2 of doi:10.1109/tro.2004.838026
		// with the exception that diagonal edge weights are calculated with the hypot
		// function instead of a step count, not closing the exit and only allowing
		// information to be known from exploration.
		g: internal.NewGridFrom(
			"..................",
			"..................",
			"..................",
			"..................",
			"..................",
			"..................",
			"....*.*...........",
			"*****.***.........",
			"..................", // Keep open.
			"......***.........",
			"......*...........",
			"......*...........",
			"......*...........",
			"*****.*...........",
			"......*...........",
		),
		radius:   1.5,
		all:      true,
		diag:     true,
		remember: []bool{false, true},

		heuristic: func(dx, dy float64) float64 {
			return math.Max(math.Abs(dx), math.Abs(dy))
		},

		s: concrete.Node(253),
		t: concrete.Node(122),

		want: []graph.Node{
			concrete.Node(253),
			concrete.Node(254),
			concrete.Node(255),
			concrete.Node(256),
			concrete.Node(239),
			concrete.Node(221),
			concrete.Node(203),
			concrete.Node(185),
			concrete.Node(167),
			concrete.Node(150),
			concrete.Node(151),
			concrete.Node(152),

			// The following section depends
			// on map iteration order.
			nil,
			nil,
			nil,
			nil,
			nil,

			concrete.Node(122),
		},
		weight: 18.656854249492383,
	},
	{
		// This is the example shown in figure 2 of doi:10.1109/tro.2004.838026
		// with the exception that diagonal edge weights are calculated with the hypot
		// function instead of a step count, the exit is closed at a distance and
		// information is allowed to be known from exploration.
		g: internal.NewGridFrom(
			"..................",
			"..................",
			"..................",
			"..................",
			"..................",
			"..................",
			"....*.*...........",
			"*****.***.........",
			"........*.........",
			"......***.........",
			"......*...........",
			"......*...........",
			"......*...........",
			"*****.*...........",
			"......*...........",
		),
		radius:   1.5,
		all:      true,
		diag:     true,
		remember: []bool{false, true},

		heuristic: func(dx, dy float64) float64 {
			return math.Max(math.Abs(dx), math.Abs(dy))
		},

		s: concrete.Node(253),
		t: concrete.Node(122),

		want: []graph.Node{
			concrete.Node(253),
			concrete.Node(254),
			concrete.Node(255),
			concrete.Node(256),
			concrete.Node(239),
			concrete.Node(221),
			concrete.Node(203),
			concrete.Node(185),
			concrete.Node(167),
			concrete.Node(150),
			concrete.Node(151),
			concrete.Node(150),
			concrete.Node(131),
			concrete.Node(113),
			concrete.Node(96),

			// The following section depends
			// on map iteration order.
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,

			concrete.Node(122),
		},
		weight: 24.07106781186548,
	},
	{
		// This is the example shown in figure 2 of doi:10.1109/tro.2004.838026
		// with the exception that diagonal edge weights are calculated with the hypot
		// function instead of a step count.
		g: internal.NewGridFrom(
			"..................",
			"..................",
			"..................",
			"..................",
			"..................",
			"..................",
			"....*.*...........",
			"*****.***.........",
			"......*...........", // Forget this wall.
			"......***.........",
			"......*...........",
			"......*...........",
			"......*...........",
			"*****.*...........",
			"......*...........",
		),
		radius:   1.5,
		all:      true,
		diag:     true,
		remember: []bool{true},

		modify: func(l *internal.LimitedVisionGrid) {
			all := l.Grid.AllVisible
			l.Grid.AllVisible = false
			for _, n := range l.Nodes() {
				l.Known[n.ID()] = !l.Grid.Has(n)
			}
			l.Grid.AllVisible = all

			const (
				wallRow = 8
				wallCol = 6
			)
			l.Known[l.NodeAt(wallRow, wallCol).ID()] = false

			// Check we have a correctly modified representation.
			for _, u := range l.Nodes() {
				for _, v := range l.Nodes() {
					if l.HasEdge(u, v) != l.Grid.HasEdge(u, v) {
						ur, uc := l.RowCol(u.ID())
						vr, vc := l.RowCol(v.ID())
						if (ur == wallRow && uc == wallCol) || (vr == wallRow && vc == wallCol) {
							if !l.HasEdge(u, v) {
								panic(fmt.Sprintf("expected to believe edge between %v (%d,%d) and %v (%d,%d) is passable",
									u, v, ur, uc, vr, vc))
							}
							continue
						}
						panic(fmt.Sprintf("disagreement about edge between %v (%d,%d) and %v (%d,%d): got:%t want:%t",
							u, v, ur, uc, vr, vc, l.HasEdge(u, v), l.Grid.HasEdge(u, v)))
					}
				}
			}
		},

		heuristic: func(dx, dy float64) float64 {
			return math.Max(math.Abs(dx), math.Abs(dy))
		},

		s: concrete.Node(253),
		t: concrete.Node(122),

		want: []graph.Node{
			concrete.Node(253),
			concrete.Node(254),
			concrete.Node(255),
			concrete.Node(256),
			concrete.Node(239),
			concrete.Node(221),
			concrete.Node(203),
			concrete.Node(185),
			concrete.Node(167),
			concrete.Node(149),
			concrete.Node(131),
			concrete.Node(113),
			concrete.Node(96),

			// The following section depends
			// on map iteration order.
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,

			concrete.Node(122),
		},
		weight: 21.242640687119287,
	},
	{
		g: internal.NewGridFrom(
			"*..*",
			"**.*",
			"**.*",
			"**.*",
		),
		radius:   1,
		all:      true,
		diag:     false,
		remember: []bool{false, true},

		heuristic: func(dx, dy float64) float64 {
			return math.Hypot(dx, dy)
		},

		s: concrete.Node(1),
		t: concrete.Node(14),

		want: []graph.Node{
			concrete.Node(1),
			concrete.Node(2),
			concrete.Node(6),
			concrete.Node(10),
			concrete.Node(14),
		},
		weight: 4,
	},
	{
		g: internal.NewGridFrom(
			"*..*",
			"**.*",
			"**.*",
			"**.*",
		),
		radius:   1.5,
		all:      true,
		diag:     true,
		remember: []bool{false, true},

		heuristic: func(dx, dy float64) float64 {
			return math.Hypot(dx, dy)
		},

		s: concrete.Node(1),
		t: concrete.Node(14),

		want: []graph.Node{
			concrete.Node(1),
			concrete.Node(6),
			concrete.Node(10),
			concrete.Node(14),
		},
		weight: math.Sqrt2 + 2,
	},
	{
		g: internal.NewGridFrom(
			"...",
			".*.",
			".*.",
			".*.",
			".*.",
		),
		radius:   1,
		all:      true,
		diag:     false,
		remember: []bool{false, true},

		heuristic: func(dx, dy float64) float64 {
			return math.Hypot(dx, dy)
		},

		s: concrete.Node(6),
		t: concrete.Node(14),

		want: []graph.Node{
			concrete.Node(6),
			concrete.Node(9),
			concrete.Node(12),
			concrete.Node(9),
			concrete.Node(6),
			concrete.Node(3),
			concrete.Node(0),
			concrete.Node(1),
			concrete.Node(2),
			concrete.Node(5),
			concrete.Node(8),
			concrete.Node(11),
			concrete.Node(14),
		},
		weight: 12,
	},
}

func TestDStarLiteDynamic(t *testing.T) {
	for i, test := range dynamicDStarLiteTests {
		for _, remember := range test.remember {
			l := &internal.LimitedVisionGrid{
				Grid:         test.g,
				VisionRadius: test.radius,
				Location:     test.s,
			}
			if remember {
				l.Known = make(map[int]bool)
			}

			l.Grid.AllVisible = test.all

			l.Grid.AllowDiagonal = test.diag
			l.Grid.UnitEdgeWeight = test.unit

			if test.modify != nil {
				test.modify(l)
			}

			got := []graph.Node{test.s}
			l.MoveTo(test.s)

			heuristic := func(a, b graph.Node) float64 {
				ax, ay := l.XY(a)
				bx, by := l.XY(b)
				return test.heuristic(ax-bx, ay-by)
			}

			world := concrete.NewDirectedGraph()
			d := NewDStarLite(test.s, test.t, l, heuristic, world)
			var (
				dp  *dumper
				buf bytes.Buffer
			)
			_, c := l.Grid.Dims()
			if c <= *maxWide && (*debug || *vdebug) {
				dp = &dumper{
					w: &buf,

					dStarLite: d,
					grid:      l,
				}
			}

			dp.dump(true)
			dp.printEdges("Initial world knowledge: %s\n\n", concreteEdgesOf(l, world.Edges()))
			for d.Step() {
				changes, _ := l.MoveTo(d.Here())
				got = append(got, l.Location)
				d.UpdateWorld(changes)
				dp.dump(true)
				if wantedPath, ok := test.wantedPaths[l.Location.ID()]; ok {
					gotPath, _ := d.Path()
					if !samePath(gotPath, wantedPath) {
						t.Errorf("unexpected intermediate path estimation for test %d %s memory:\ngot: %v\nwant:%v",
							i, memory(remember), gotPath, wantedPath)
					}
				}
				dp.printEdges("Edges changing after last step:\n%s\n\n", concreteEdgesOf(l, changes))
			}

			if weight := weightOf(got, l.Grid); !samePath(got, test.want) || weight != test.weight {
				t.Errorf("unexpected path for test %d %s memory got weight:%v want weight:%v:\ngot: %v\nwant:%v",
					i, memory(remember), weight, test.weight, got, test.want)
				b, err := l.Render(got)
				t.Errorf("path taken (err:%v):\n%s", err, b)
				if c <= *maxWide && (*debug || *vdebug) {
					t.Error(buf.String())
				}
			} else if c <= *maxWide && *vdebug {
				t.Logf("Test %d:\n%s", i, buf.String())
			}
		}
	}
}

type memory bool

func (m memory) String() string {
	if m {
		return "with"
	}
	return "without"
}

// samePath compares two paths for equality ignoring nodes that are nil.
func samePath(a, b []graph.Node) bool {
	if len(a) != len(b) {
		return false
	}
	for i, e := range a {
		if e == nil || b[i] == nil {
			continue
		}
		if e.ID() != b[i].ID() {
			return false
		}
	}
	return true
}

type weightedGraph interface {
	graph.Graph
	graph.Weighter
}

// weightOf return the weight of the path in g.
func weightOf(path []graph.Node, g weightedGraph) float64 {
	var w float64
	if len(path) > 1 {
		for p, n := range path[1:] {
			e := g.Edge(path[p], n)
			if e == nil {
				return math.Inf(1)
			}
			w += g.Weight(e)
		}
	}
	return w
}

// concreteEdgesOf returns the weighted edges in g corresponding to the given edges.
func concreteEdgesOf(g weightedGraph, edges []graph.Edge) []concrete.Edge {
	w := make([]concrete.Edge, len(edges))
	for i, e := range edges {
		w[i].F = e.From()
		w[i].T = e.To()
		w[i].W = g.Weight(e)
	}
	return w
}
