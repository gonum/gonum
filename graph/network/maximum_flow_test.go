// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network

import (
	"fmt"
	"reflect"
	"testing"

	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

const tol = 1e-10

var dinicTests = []struct {
	name      string
	g         graph.WeightedDirected
	testFlows []testFlow
}{
	{
		name: "source_is_target",
		g: func() graph.WeightedDirected {
			g := simple.NewWeightedDirectedGraph(0, 0)
			g.AddNode(simple.Node(0))
			return g
		}(),
		testFlows: []testFlow{
			{s: simple.Node(0), t: simple.Node(0), wantPanic: "no cut between s and t"},
		},
	},
	{
		name: "negative_capacity",
		g: func() graph.WeightedDirected {
			g := simple.NewWeightedDirectedGraph(0, 0)
			for _, e := range []simple.WeightedEdge{
				{F: simple.Node(0), T: simple.Node(1), W: 0.3},
				{F: simple.Node(1), T: simple.Node(2), W: -0.6},
			} {
				g.SetWeightedEdge(e)
			}
			return g
		}(),
		testFlows: []testFlow{
			{s: simple.Node(0), t: simple.Node(1), wantPanic: "negative edge weight"},
		},
	},
	{
		name: "three_disjoint_paths",
		g: func() graph.WeightedDirected {
			g := simple.NewWeightedDirectedGraph(0, 0)
			for _, e := range []struct{ u, v int64 }{
				{0, 1},
				{1, 4},
				{0, 2},
				{2, 4},
				{0, 3},
				{3, 4},
			} {
				g.SetWeightedEdge(simple.WeightedEdge{F: simple.Node(e.u), T: simple.Node(e.v), W: 1})
			}
			return g
		}(),
		testFlows: []testFlow{
			{s: simple.Node(0), t: simple.Node(4), want: 3},
		},
	},
	{
		name: "three_disjoint_paths_parallel",
		g: func() graph.WeightedDirected {
			g := simple.NewWeightedDirectedGraph(0, 0)
			for _, e := range []struct{ u, v int64 }{
				{0, 1},
				{1, 0},
				{1, 4},
				{0, 2},
				{2, 0},
				{2, 4},
				{0, 3},
				{3, 0},
				{3, 4},
			} {
				g.SetWeightedEdge(simple.WeightedEdge{F: simple.Node(e.u), T: simple.Node(e.v), W: 1})
			}
			return g
		}(),
		testFlows: []testFlow{
			{s: simple.Node(0), t: simple.Node(4), want: 3},
		},
	},
	{
		name: "cycle_with_tail",
		g: func() graph.WeightedDirected {
			g := simple.NewWeightedDirectedGraph(0, 0)
			// Cycle: 0->1->2->0 and tail 2->3
			for _, e := range []simple.WeightedEdge{
				{F: simple.Node(0), T: simple.Node(1), W: 0.3},
				{F: simple.Node(1), T: simple.Node(2), W: 0.6},
				{F: simple.Node(2), T: simple.Node(0), W: 0.9},
				{F: simple.Node(2), T: simple.Node(3), W: 0.7},
			} {
				g.SetWeightedEdge(e)
			}
			return g
		}(),
		testFlows: []testFlow{
			{s: simple.Node(0), t: simple.Node(3), want: 0.3},
			{s: simple.Node(1), t: simple.Node(3), want: 0.6},
		},
	},
	{
		name: "cycle_with_tail_parallel",
		g: func() graph.WeightedDirected {
			g := simple.NewWeightedDirectedGraph(0, 0)
			// Layers: 0->{1,2,3}, {1,2}->{4,5}, {3}->{5,6}, {4,5,6}->{7}
			for _, e := range []simple.WeightedEdge{
				{F: simple.Node(0), T: simple.Node(1), W: 0.3},
				{F: simple.Node(1), T: simple.Node(2), W: 0.6},
				{F: simple.Node(2), T: simple.Node(0), W: 0.9},
				{F: simple.Node(2), T: simple.Node(3), W: 0.7},
				{F: simple.Node(1), T: simple.Node(0), W: 1.3},
				{F: simple.Node(2), T: simple.Node(1), W: 1.6},
				{F: simple.Node(0), T: simple.Node(2), W: 1.9},
			} {
				g.SetWeightedEdge(e)
			}
			return g
		}(),
		testFlows: []testFlow{
			{s: simple.Node(0), t: simple.Node(3), want: 0.7},
			{s: simple.Node(1), t: simple.Node(3), want: 0.7},
		},
	},
	{
		name: "four_layer_dag",
		g: func() graph.WeightedDirected {
			g := simple.NewWeightedDirectedGraph(0, 0)
			for _, e := range []struct{ u, v int64 }{
				{0, 1},
				{0, 2},
				{0, 3},
				{1, 4},
				{2, 4},
				{2, 5},
				{3, 5},
				{3, 6},
				{4, 7},
				{5, 7},
				{6, 7},
			} {
				g.SetWeightedEdge(simple.WeightedEdge{F: simple.Node(e.u), T: simple.Node(e.v), W: 1})
			}
			return g
		}(),
		testFlows: []testFlow{
			{s: simple.Node(0), t: simple.Node(7), want: 3},
			{s: simple.Node(3), t: simple.Node(7), want: 2},
			{s: simple.Node(0), t: simple.Node(5), want: 2},
			{s: simple.Node(2), t: simple.Node(4), want: 1},
		},
	},
	{
		name: "diamond_with_cross",
		g: func() graph.WeightedDirected {
			g := simple.NewWeightedDirectedGraph(0, 0)
			for _, e := range []simple.WeightedEdge{
				{F: simple.Node(0), T: simple.Node(1), W: 10},
				{F: simple.Node(0), T: simple.Node(2), W: 10},
				{F: simple.Node(1), T: simple.Node(2), W: 5},
				{F: simple.Node(1), T: simple.Node(3), W: 10},
				{F: simple.Node(2), T: simple.Node(3), W: 10},
			} {
				g.SetWeightedEdge(e)
			}
			return g
		}(),
		testFlows: []testFlow{
			{s: simple.Node(0), t: simple.Node(3), want: 20},
			{s: simple.Node(0), t: simple.Node(2), want: 15},
		},
	},
	{
		name: "disconnected",
		g: func() graph.WeightedDirected {
			g := simple.NewWeightedDirectedGraph(0, 0)
			for _, e := range []simple.WeightedEdge{
				{F: simple.Node(0), T: simple.Node(1), W: 10},
				{F: simple.Node(1), T: simple.Node(2), W: 5},
				{F: simple.Node(2), T: simple.Node(3), W: 7},
				{F: simple.Node(4), T: simple.Node(5), W: 11},
				{F: simple.Node(5), T: simple.Node(6), W: 10},
			} {
				g.SetWeightedEdge(e)
			}
			return g
		}(),
		testFlows: []testFlow{
			{s: simple.Node(0), t: simple.Node(5), want: 0},
		},
	},
}

type testFlow struct {
	s, t      graph.Node
	want      float64
	wantPanic any
}

func TestMaxFlowDinic(t *testing.T) {
	const tol = 1e-10

	for _, test := range dinicTests {
		t.Run(test.name, func(t *testing.T) {
			for _, flow := range test.testFlows {
				t.Run(fmt.Sprintf("%d_to_%d", flow.s, flow.t), func(t *testing.T) {
					defer func() {
						r := recover()
						if !reflect.DeepEqual(r, flow.wantPanic) {
							t.Errorf("unexpected panic: got:%v want:%v", r, flow.wantPanic)
						}
					}()
					got := MaxFlowDinic(test.g, flow.s, flow.t)
					if !scalar.EqualWithinAbs(got, flow.want, tol) {
						t.Errorf("unexpected maximum flow: got = %v, want = %v", got, flow.want)
					}
				})
			}
		})
	}
}
