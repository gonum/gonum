// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gen

import (
	"bytes"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

type nodeIDGraphBuilder interface {
	graph.Graph
	NodeIDGraphBuilder
}

func undirected() nodeIDGraphBuilder { return simple.NewUndirectedGraph() }
func directed() nodeIDGraphBuilder   { return simple.NewDirectedGraph() }

type empty struct{}

func (r empty) Len() int       { return 0 }
func (r empty) ID(i int) int64 { panic("called ID on empty IDer") }

func TestComplete(t *testing.T) {
	tests := []struct {
		name string
		ids  IDer
		dst  func() nodeIDGraphBuilder
		want string
	}{
		{
			name: "empty",
			ids:  empty{},
			dst:  undirected,
			want: `strict graph empty {
}`,
		},
		{
			name: "single",
			ids:  IDRange{First: 1, Last: 1},
			dst:  undirected,
			want: `strict graph single {
 // Node definitions.
 1;
}`,
		},
		{
			name: "pair_undirected",
			ids:  IDRange{First: 1, Last: 2},
			dst:  undirected,
			want: `strict graph pair_undirected {
 // Node definitions.
 1;
 2;

 // Edge definitions.
 1 -- 2;
}`,
		},
		{
			name: "pair_directed",
			ids:  IDRange{First: 1, Last: 2},
			dst:  directed,
			want: `strict digraph pair_directed {
 // Node definitions.
 1;
 2;

 // Edge definitions.
 1 -> 2;
}`,
		},
		{
			name: "quad_undirected",
			ids:  IDRange{First: 1, Last: 4},
			dst:  undirected,
			want: `strict graph quad_undirected {
 // Node definitions.
 1;
 2;
 3;
 4;

 // Edge definitions.
 1 -- 2;
 1 -- 3;
 1 -- 4;
 2 -- 3;
 2 -- 4;
 3 -- 4;
}`,
		},
		{
			name: "quad_directed",
			ids:  IDRange{First: 1, Last: 4},
			dst:  directed,
			want: `strict digraph quad_directed {
 // Node definitions.
 1;
 2;
 3;
 4;

 // Edge definitions.
 1 -> 2;
 1 -> 3;
 1 -> 4;
 2 -> 3;
 2 -> 4;
 3 -> 4;
}`,
		},
	}

	for _, test := range tests {
		dst := test.dst()
		err := Complete(dst, test.ids)
		if err != nil {
			t.Errorf("unexpected error constructing graph: %v", err)
		}
		got, err := dot.Marshal(dst, test.name, "", " ")
		if err != nil {
			t.Errorf("unexpected marshaling graph error: %v", err)
		}
		if !bytes.Equal(got, []byte(test.want)) {
			t.Errorf("unexpected result for test %s:\ngot:\n%s\nwant:\n%s", test.name, got, test.want)
		}
	}
}

func TestCycle(t *testing.T) {
	tests := []struct {
		name string
		ids  IDer
		dst  func() nodeIDGraphBuilder
		want string
	}{
		{
			name: "empty",
			ids:  empty{},
			dst:  undirected,
			want: `strict graph empty {
}`,
		},
		{
			name: "single",
			ids:  IDRange{First: 1, Last: 1},
			dst:  undirected,
			want: `strict graph single {
 // Node definitions.
 1;
}`,
		},
		{
			name: "pair_undirected",
			ids:  IDRange{First: 1, Last: 2},
			dst:  undirected,
			want: `strict graph pair_undirected {
 // Node definitions.
 1;
 2;

 // Edge definitions.
 1 -- 2;
}`,
		},
		{
			name: "pair_directed",
			ids:  IDRange{First: 1, Last: 2},
			dst:  directed,
			want: `strict digraph pair_directed {
 // Node definitions.
 1;
 2;

 // Edge definitions.
 1 -> 2;
 2 -> 1;
}`,
		},
		{
			name: "quad_undirected",
			ids:  IDRange{First: 1, Last: 4},
			dst:  undirected,
			want: `strict graph quad_undirected {
 // Node definitions.
 1;
 2;
 3;
 4;

 // Edge definitions.
 1 -- 2;
 1 -- 4;
 2 -- 3;
 3 -- 4;
}`,
		},
		{
			name: "quad_directed",
			ids:  IDRange{First: 1, Last: 4},
			dst:  directed,
			want: `strict digraph quad_directed {
 // Node definitions.
 1;
 2;
 3;
 4;

 // Edge definitions.
 1 -> 2;
 2 -> 3;
 3 -> 4;
 4 -> 1;
}`,
		},
	}

	for _, test := range tests {
		dst := test.dst()
		err := Cycle(dst, test.ids)
		if err != nil {
			t.Errorf("unexpected error constructing graph: %v", err)
		}
		got, err := dot.Marshal(dst, test.name, "", " ")
		if err != nil {
			t.Errorf("unexpected error marshaling graph: %v", err)
		}
		if !bytes.Equal(got, []byte(test.want)) {
			t.Errorf("unexpected result for test %s:\ngot:\n%s\nwant:\n%s", test.name, got, test.want)
		}
	}
}

func TestPath(t *testing.T) {
	tests := []struct {
		name string
		ids  IDer
		dst  func() nodeIDGraphBuilder
		want string
	}{
		{
			name: "empty",
			ids:  empty{},
			dst:  undirected,
			want: `strict graph empty {
}`,
		},
		{
			name: "single",
			ids:  IDRange{First: 1, Last: 1},
			dst:  undirected,
			want: `strict graph single {
 // Node definitions.
 1;
}`,
		},
		{
			name: "pair_undirected",
			ids:  IDRange{First: 1, Last: 2},
			dst:  undirected,
			want: `strict graph pair_undirected {
 // Node definitions.
 1;
 2;

 // Edge definitions.
 1 -- 2;
}`,
		},
		{
			name: "pair_directed",
			ids:  IDRange{First: 1, Last: 2},
			dst:  directed,
			want: `strict digraph pair_directed {
 // Node definitions.
 1;
 2;

 // Edge definitions.
 1 -> 2;
}`,
		},
		{
			name: "quad_undirected",
			ids:  IDRange{First: 1, Last: 4},
			dst:  undirected,
			want: `strict graph quad_undirected {
 // Node definitions.
 1;
 2;
 3;
 4;

 // Edge definitions.
 1 -- 2;
 2 -- 3;
 3 -- 4;
}`,
		},
		{
			name: "quad_directed",
			ids:  IDRange{First: 1, Last: 4},
			dst:  directed,
			want: `strict digraph quad_directed {
 // Node definitions.
 1;
 2;
 3;
 4;

 // Edge definitions.
 1 -> 2;
 2 -> 3;
 3 -> 4;
}`,
		},
	}

	for _, test := range tests {
		dst := test.dst()
		err := Path(dst, test.ids)
		if err != nil {
			t.Errorf("unexpected error constructing graph: %v", err)
		}
		got, err := dot.Marshal(dst, test.name, "", " ")
		if err != nil {
			t.Errorf("unexpected error marshaling graph: %v", err)
		}
		if !bytes.Equal(got, []byte(test.want)) {
			t.Errorf("unexpected result for test %s:\ngot:\n%s\nwant:\n%s", test.name, got, test.want)
		}
	}
}

func TestStar(t *testing.T) {
	tests := []struct {
		name   string
		center int64
		leaves IDer
		dst    func() nodeIDGraphBuilder
		want   string
	}{
		{
			name:   "empty_leaves",
			center: 0,
			leaves: empty{},
			dst:    undirected,
			want: `strict graph empty_leaves {
 // Node definitions.
 0;
}`,
		},
		{
			name:   "single",
			leaves: IDRange{First: 1, Last: 1},
			dst:    undirected,
			want: `strict graph single {
 // Node definitions.
 0;
 1;

 // Edge definitions.
 0 -- 1;
}`,
		},
		{
			name:   "pair_undirected",
			leaves: IDRange{First: 1, Last: 2},
			dst:    undirected,
			want: `strict graph pair_undirected {
 // Node definitions.
 0;
 1;
 2;

 // Edge definitions.
 0 -- 1;
 0 -- 2;
}`,
		},
		{
			name:   "pair_directed",
			leaves: IDRange{First: 1, Last: 2},
			dst:    directed,
			want: `strict digraph pair_directed {
 // Node definitions.
 0;
 1;
 2;

 // Edge definitions.
 0 -> 1;
 0 -> 2;
}`,
		},
		{
			name:   "quad_undirected",
			leaves: IDRange{First: 1, Last: 4},
			dst:    undirected,
			want: `strict graph quad_undirected {
 // Node definitions.
 0;
 1;
 2;
 3;
 4;

 // Edge definitions.
 0 -- 1;
 0 -- 2;
 0 -- 3;
 0 -- 4;
}`,
		},
		{
			name:   "quad_directed",
			leaves: IDRange{First: 1, Last: 4},
			dst:    directed,
			want: `strict digraph quad_directed {
 // Node definitions.
 0;
 1;
 2;
 3;
 4;

 // Edge definitions.
 0 -> 1;
 0 -> 2;
 0 -> 3;
 0 -> 4;
}`,
		},
	}

	for _, test := range tests {
		dst := test.dst()
		err := Star(dst, test.center, test.leaves)
		if err != nil {
			t.Errorf("unexpected error constructing graph: %v", err)
		}
		got, err := dot.Marshal(dst, test.name, "", " ")
		if err != nil {
			t.Errorf("unexpected error marshaling graph: %v", err)
		}
		if !bytes.Equal(got, []byte(test.want)) {
			t.Errorf("unexpected result for test %s:\ngot:\n%s\nwant:\n%s", test.name, got, test.want)
		}
	}
}

func TestWheel(t *testing.T) {
	tests := []struct {
		name   string
		center int64
		cycle  IDer
		dst    func() nodeIDGraphBuilder
		want   string
	}{
		{
			name:   "empty_cycle",
			center: 0,
			cycle:  empty{},
			dst:    undirected,
			want: `strict graph empty_cycle {
 // Node definitions.
 0;
}`,
		},
		{
			name:  "single",
			cycle: IDRange{First: 1, Last: 1},
			dst:   undirected,
			want: `strict graph single {
 // Node definitions.
 0;
 1;

 // Edge definitions.
 0 -- 1;
}`,
		},
		{
			name:  "pair_undirected",
			cycle: IDRange{First: 1, Last: 2},
			dst:   undirected,
			want: `strict graph pair_undirected {
 // Node definitions.
 0;
 1;
 2;

 // Edge definitions.
 0 -- 1;
 0 -- 2;
 1 -- 2;
}`,
		},
		{
			name:  "pair_directed",
			cycle: IDRange{First: 1, Last: 2},
			dst:   directed,
			want: `strict digraph pair_directed {
 // Node definitions.
 0;
 1;
 2;

 // Edge definitions.
 0 -> 1;
 0 -> 2;
 1 -> 2;
 2 -> 1;
}`,
		},
		{
			name:  "quad_undirected",
			cycle: IDRange{First: 1, Last: 4},
			dst:   undirected,
			want: `strict graph quad_undirected {
 // Node definitions.
 0;
 1;
 2;
 3;
 4;

 // Edge definitions.
 0 -- 1;
 0 -- 2;
 0 -- 3;
 0 -- 4;
 1 -- 2;
 1 -- 4;
 2 -- 3;
 3 -- 4;
}`,
		},
		{
			name:  "quad_directed",
			cycle: IDRange{First: 1, Last: 4},
			dst:   directed,
			want: `strict digraph quad_directed {
 // Node definitions.
 0;
 1;
 2;
 3;
 4;

 // Edge definitions.
 0 -> 1;
 0 -> 2;
 0 -> 3;
 0 -> 4;
 1 -> 2;
 2 -> 3;
 3 -> 4;
 4 -> 1;
}`,
		},
	}

	for _, test := range tests {
		dst := test.dst()
		err := Wheel(dst, test.center, test.cycle)
		if err != nil {
			t.Errorf("unexpected error constructing graph: %v", err)
		}
		got, err := dot.Marshal(dst, test.name, "", " ")
		if err != nil {
			t.Errorf("unexpected error marshaling graph: %v", err)
		}
		if !bytes.Equal(got, []byte(test.want)) {
			t.Errorf("unexpected result for test %s:\ngot:\n%s\nwant:\n%s", test.name, got, test.want)
		}
	}
}
