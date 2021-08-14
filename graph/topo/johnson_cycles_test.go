// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topo

import (
	"fmt"
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
			2: {{2, 3, 4, 2},
				{2, 6, 3, 4, 2}},
			3: {{3, 4, 2, 3},
				{3, 4, 2, 6, 3}},
			4: {{4, 2, 3, 4},
				{4, 2, 6, 3, 4}},
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
			fmt.Println(got, test.want[k])
			if !reflect.DeepEqual(got, test.want[k]) {
				t.Errorf("unexpected johnson result for %d:\n\tgot:%#v\n\twant:%#v", i, got, test.want)
			}
		}
	}
}
