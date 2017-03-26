// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topo

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/gonum/graph/internal/ordered"
	"github.com/gonum/graph/simple"
)

var cyclesInTests = []struct {
	g    []intset
	sccs [][]int
	want [][]int
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
		want: [][]int{
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
		want: [][]int{
			{1, 2, 3, 1},
		},
	},
	{
		g: []intset{
			0: linksTo(1),
			1: linksTo(0, 2),
			2: linksTo(1),
		},
		want: [][]int{
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
		want: [][]int{
			{0, 1, 2, 0},
			{3, 4, 3},
		},
	},
}

func TestCyclesIn(t *testing.T) {
	for i, test := range cyclesInTests {
		g := simple.NewDirectedGraph(0, math.Inf(1))
		g.AddNode(simple.Node(-10)) // Make sure we test graphs with sparse IDs.
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if !g.Has(simple.Node(u)) {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}
		cycles := CyclesIn(g)
		var got [][]int
		if cycles != nil {
			got = make([][]int, len(cycles))
		}
		// johnson.circuit does range iteration over maps,
		// so sort to ensure consistent ordering.
		for j, c := range cycles {
			ids := make([]int, len(c))
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
