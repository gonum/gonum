// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topo

import (
	"reflect"
	"sort"
	"testing"

	"gonum.org/v1/gonum/graph/internal/ordered"
	"gonum.org/v1/gonum/graph/simple"
)

var cyclesInTests = []struct {
	g    []intset
	want [][]int64
}{
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(2, 7),
			2: linksTo(3, 6),
			3: linksTo(4),
			4: linksTo(2, 5),
			6: linksTo(3, 5),
			7: linksTo(0, 6),
		},
		want: [][]int64{
			{0, 1, 7, 0},
			{2, 3, 4, 2},
			{2, 6, 3, 4, 2},
		},
	},
	{
		g: []intset{
			0: linksTo(1, 2, 3),
			1: linksTo(2),
			2: linksTo(3),
			3: linksTo(1),
		},
		want: [][]int64{
			{1, 2, 3, 1},
		},
	},
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(0, 2),
			2: linksTo(1),
		},
		want: [][]int64{
			{0, 1, 0},
			{1, 2, 1},
		},
	},
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(2, 3),
			2: linksTo(4, 5),
			3: linksTo(4, 5),
			4: linksTo(6),
			5: nil,
			6: nil,
		},
		want: nil,
	},
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(2, 3, 4),
			2: linksTo(0, 3),
			3: linksTo(4),
			4: linksTo(3),
		},
		want: [][]int64{
			{0, 1, 2, 0},
			{3, 4, 3},
		},
	},
}

func TestDirectedCyclesIn(t *testing.T) {
	for i, test := range cyclesInTests {
		g := simple.NewDirectedGraph()
		g.AddNode(simple.Node(-10)) // Make sure we test graphs with sparse IDs.
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if g.Node(int64(u)) == nil {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}
		cycles := DirectedCyclesIn(g)
		var got [][]int64
		if cycles != nil {
			got = make([][]int64, len(cycles))
		}
		// johnson.circuit does range iteration over maps,
		// so sort to ensure consistent ordering.
		for j, c := range cycles {
			ids := make([]int64, len(c))
			for k, n := range c {
				ids[k] = n.ID()
			}
			got[j] = ids
		}
		sort.Sort(ordered.BySliceValues(got))
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("unexpected johnson result for %d:\n\tgot:%#v\n\twant:%#v", i, got, test.want)
		}
	}
}

var cyclesContainingTests = []struct {
	g    []intset
	want [][][]int64
}{
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(2, 7),
			2: linksTo(3, 6),
			3: linksTo(4),
			4: linksTo(2, 5),
			6: linksTo(3, 5),
			7: linksTo(0, 6),
		},
		want: [][][]int64{
			0: {{0, 1, 7, 0}},
			1: {{1, 7, 0, 1}},
			2: {{2, 3, 4, 2}, {2, 6, 3, 4, 2}},
			3: {{3, 4, 2, 3}, {3, 4, 2, 6, 3}},
			4: {{4, 2, 3, 4}, {4, 2, 6, 3, 4}},
			5: nil,
			6: {{6, 3, 4, 2, 6}},
			7: {{7, 0, 1, 7}},
		},
	},
	{
		g: []intset{
			0: linksTo(1, 2, 3),
			1: linksTo(2),
			2: linksTo(3),
			3: linksTo(1),
		},
		want: [][][]int64{
			0: nil,
			1: {{1, 2, 3, 1}},
			2: {{2, 3, 1, 2}},
			3: {{3, 1, 2, 3}},
		},
	},
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(0, 2),
			2: linksTo(1),
		},
		want: [][][]int64{
			0: {{0, 1, 0}},
			1: {{1, 0, 1}, {1, 2, 1}},
			2: {{2, 1, 2}},
		},
	},
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(2, 3),
			2: linksTo(4, 5),
			3: linksTo(4, 5),
			4: linksTo(6),
			5: nil,
			6: nil,
		},
		want: [][][]int64{
			0: nil,
			1: nil,
			2: nil,
			3: nil,
			4: nil,
			5: nil,
			6: nil,
		},
	},
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(2, 3, 4),
			2: linksTo(0, 3),
			3: linksTo(4),
			4: linksTo(3),
		},
		want: [][][]int64{
			0: {{0, 1, 2, 0}},
			1: {{1, 2, 0, 1}},
			2: {{2, 0, 1, 2}},
			3: {{3, 4, 3}},
			4: {{4, 3, 4}},
		},
	},
}

func TestDirectedCyclesContaining(t *testing.T) {
	for i, test := range cyclesContainingTests {
		g := simple.NewDirectedGraph()
		g.AddNode(simple.Node(-10)) // Make sure we test graphs with sparse IDs.
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if g.Node(int64(u)) == nil {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}
		for k := range test.want {
			cycles := DirectedCyclesContaining(g, int64(k))
			var got [][]int64
			if cycles != nil {
				got = make([][]int64, len(cycles))
			}
			// johnson.circuit does range iteration over maps,
			// so sort to ensure consistent ordering.
			for j, c := range cycles {
				ids := make([]int64, len(c))
				for k, n := range c {
					ids[k] = n.ID()
				}
				got[j] = ids
			}
			sort.Sort(ordered.BySliceValues(got))
			if !reflect.DeepEqual(got, test.want[k]) {
				t.Errorf("unexpected johnson result for %d:\n\tgot:%#v\n\twant:%#v", i, got, test.want)
			}
		}
	}
}

var cyclesOfMaxLenTests = []struct {
	g       []intset
	wantMap map[int][][]int64
}{
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(2, 7),
			2: linksTo(3, 6),
			3: linksTo(4),
			4: linksTo(2, 5),
			6: linksTo(3, 5),
			7: linksTo(0, 6),
		},
		wantMap: map[int][][]int64{
			3: nil,
			4: {{0, 1, 7, 0}, {2, 3, 4, 2}},
			5: {{0, 1, 7, 0}, {2, 3, 4, 2}, {2, 6, 3, 4, 2}},
			6: {{0, 1, 7, 0}, {2, 3, 4, 2}, {2, 6, 3, 4, 2}},
		},
	},
	{
		g: []intset{
			0: linksTo(1, 2, 3),
			1: linksTo(2),
			2: linksTo(3),
			3: linksTo(1),
		},
		wantMap: map[int][][]int64{
			3: nil,
			4: {{1, 2, 3, 1}},
			5: {{1, 2, 3, 1}},
		},
	},
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(0, 2),
			2: linksTo(1),
		},
		wantMap: map[int][][]int64{
			2: nil,
			3: {{0, 1, 0}, {1, 2, 1}},
			4: {{0, 1, 0}, {1, 2, 1}},
		},
	},
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(2, 3, 4),
			2: linksTo(0, 3),
			3: linksTo(4),
			4: linksTo(3),
		},
		wantMap: map[int][][]int64{
			2: nil,
			3: {{3, 4, 3}},
			4: {{0, 1, 2, 0}, {3, 4, 3}},
			5: {{0, 1, 2, 0}, {3, 4, 3}},
		},
	},
}

func TestDirectedCyclesOfMaxLen(t *testing.T) {
	for i, test := range cyclesOfMaxLenTests {
		g := simple.NewDirectedGraph()
		g.AddNode(simple.Node(-10)) // Make sure we test graphs with sparse IDs.
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if g.Node(int64(u)) == nil {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}
		for maxLen, want := range test.wantMap {
			cycles := DirectedCyclesOfMaxLen(g, maxLen)
			var got [][]int64
			if cycles != nil {
				got = make([][]int64, len(cycles))
			}
			// johnson.circuit does range iteration over maps,
			// so sort to ensure consistent ordering.
			for j, c := range cycles {
				ids := make([]int64, len(c))
				for k, n := range c {
					ids[k] = n.ID()
				}
				got[j] = ids
			}
			sort.Sort(ordered.BySliceValues(got))
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unexpected johnson result for %d:\n\tgot:%#v\n\twant:%#v", i, got, want)
			}
		}
	}
}

type wantInstance struct {
	maxLen int
	vid    int64
	want   [][]int64
}

var cyclesOfMaxLenContainingTests = []struct {
	g []intset
	w []wantInstance
}{
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(2, 7),
			2: linksTo(3, 6),
			3: linksTo(4),
			4: linksTo(2, 5),
			6: linksTo(3, 5),
			7: linksTo(0, 6),
		},
		w: []wantInstance{
			{
				maxLen: 3,
				vid:    3,
				want:   nil,
			},
			{
				maxLen: 4,
				vid:    3,
				want:   [][]int64{{3, 4, 2, 3}},
			},
			{
				maxLen: 5,
				vid:    3,
				want:   [][]int64{{3, 4, 2, 3}, {3, 4, 2, 6, 3}},
			},
			{
				maxLen: 4,
				vid:    6,
				want:   nil,
			},
			{
				maxLen: 5,
				vid:    6,
				want:   [][]int64{{6, 3, 4, 2, 6}},
			},
		},
	},
	{
		g: []intset{
			0: linksTo(1, 2, 3),
			1: linksTo(2),
			2: linksTo(3),
			3: linksTo(1),
		},
		w: []wantInstance{
			{
				maxLen: 3,
				vid:    1,
				want:   nil,
			},
			{
				maxLen: 4,
				vid:    1,
				want:   [][]int64{{1, 2, 3, 1}},
			},
			{
				maxLen: 5,
				vid:    1,
				want:   [][]int64{{1, 2, 3, 1}},
			},
		},
	},
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(0, 2),
			2: linksTo(1),
		},
		w: []wantInstance{
			{
				maxLen: 2,
				vid:    0,
				want:   nil,
			},
			{
				maxLen: 3,
				vid:    0,
				want:   [][]int64{{0, 1, 0}},
			},
			{
				maxLen: 3,
				vid:    1,
				want:   [][]int64{{1, 0, 1}, {1, 2, 1}},
			},
		},
	},
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(2, 3, 4),
			2: linksTo(0, 3),
			3: linksTo(4),
			4: linksTo(3),
		},
		w: []wantInstance{
			{
				maxLen: 3,
				vid:    2,
				want:   nil,
			},
			{
				maxLen: 4,
				vid:    2,
				want:   [][]int64{{2, 0, 1, 2}},
			},
			{
				maxLen: 3,
				vid:    3,
				want:   [][]int64{{3, 4, 3}},
			},
			{
				maxLen: 6,
				vid:    3,
				want:   [][]int64{{3, 4, 3}},
			},
		},
	},
}

func TestDirectedCyclesOfMaxLenContaining(t *testing.T) {
	for i, test := range cyclesOfMaxLenContainingTests {
		g := simple.NewDirectedGraph()
		g.AddNode(simple.Node(-10)) // Make sure we test graphs with sparse IDs.
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if g.Node(int64(u)) == nil {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}
		for tuIndex, tu := range test.w {
			cycles := DirectedCyclesOfMaxLenContaining(g, tu.maxLen, tu.vid)
			var got [][]int64
			if cycles != nil {
				got = make([][]int64, len(cycles))
			}
			// johnson.circuit does range iteration over maps,
			// so sort to ensure consistent ordering.
			for j, c := range cycles {
				ids := make([]int64, len(c))
				for k, n := range c {
					ids[k] = n.ID()
				}
				got[j] = ids
			}
			sort.Sort(ordered.BySliceValues(got))
			if !reflect.DeepEqual(got, tu.want) {
				t.Errorf("unexpected johnson result for %d,%d:\n\tgot:%#v\n\twant:%#v", i, tuIndex, got, tu.want)
			}
		}
	}
}

type wantInstanceAO struct {
	maxLen int
	vids   []int64
	want   [][]int64
}

var cyclesOfMaxLenContainingAnyOfTests = []struct {
	g []intset
	w []wantInstanceAO
}{
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(2, 3, 4),
			2: linksTo(0, 3),
			3: linksTo(4),
			4: linksTo(3),
		},
		w: []wantInstanceAO{
			{
				maxLen: 4,
				vids:   []int64{2, 3},
				want:   [][]int64{{2, 0, 1, 2}, {3, 4, 3}},
			},
			{
				maxLen: 4,
				vids:   []int64{4, 1},
				want:   [][]int64{{1, 2, 0, 1}, {4, 3, 4}},
			},
			{
				maxLen: 3,
				vids:   []int64{0, 1, 2, 3, 4},
				want:   [][]int64{{3, 4, 3}},
			},
			{
				maxLen: 4,
				vids:   []int64{0, 1, 2},
				want:   [][]int64{{0, 1, 2, 0}},
			},
		},
	},
}

func TestDirectedCyclesOfMaxLenContainingAnyOf(t *testing.T) {
	for i, test := range cyclesOfMaxLenContainingAnyOfTests {
		g := simple.NewDirectedGraph()
		g.AddNode(simple.Node(-10)) // Make sure we test graphs with sparse IDs.
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if g.Node(int64(u)) == nil {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}
		for tuIndex, tu := range test.w {
			cycles := DirectedCyclesOfMaxLenContainingAnyOf(g, tu.maxLen, tu.vids)
			var got [][]int64
			if cycles != nil {
				got = make([][]int64, len(cycles))
			}
			// johnson.circuit does range iteration over maps,
			// so sort to ensure consistent ordering.
			for j, c := range cycles {
				ids := make([]int64, len(c))
				for k, n := range c {
					ids[k] = n.ID()
				}
				got[j] = ids
			}
			sort.Sort(ordered.BySliceValues(got))
			if !reflect.DeepEqual(got, tu.want) {
				t.Errorf("unexpected johnson result for %d,%d:\n\tgot:%#v\n\twant:%#v", i, tuIndex, got, tu.want)
			}
		}
	}
}
