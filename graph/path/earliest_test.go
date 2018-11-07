// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"gonum.org/v1/gonum/graph/temporal"
	"testing"
	"gonum.org/v1/gonum/graph"
	"reflect"
)

var (
	// sketchStream is the edge stream for the arbitrary
	// graph I sketched in a notebook to demonstrate that
	// my method for earliest arrival path construction
	// worked OK.
	sketchStream = []struct{
		F, T int64
		S, E uint64
	}{
		{1, 2, 1, 1},
		{1, 3, 1, 1},
		{4, 7, 1, 1},
		{1, 4, 2, 2},
		{2, 6, 3, 3},
		{2, 5, 4, 4},
		{5, 8, 5, 5},
		{6, 8, 6, 6},
		{7, 9, 7, 7},
		{6, 9, 8, 8},
	}

	// wuChengHuangKeLuXuStream is the example temporal graph from
	// figure 1(a) of https://www.vldb.org/pvldb/vol7/p721-wu.pdf
	wuChengHuangKeLuXuStream = []struct{
		F, T int64
		S, E uint64
	}{
		{1, 2,  1,  1 },
		{1, 2,  2,  2 },
		{7, 10, 2,  2 },
		{1, 6,  3,  3 },
		{2, 7,  3,  3 },
		{2, 8,  3,  3 },
		{1, 3,  4,  4 },
		{6, 9,  5,  5 },
		{3, 8,  6,  6 },
		{7, 11, 6,  6 },
		{8, 9,  7,  7 },
		{9, 12, 8,  8 },
		{9, 12, 9,  9 },
		{1, 9,  10, 10},
	}
)

func lineAt(fid, tid int64, s, e uint64) temporal.TemporalLine {
	f := temporal.Node(fid)
	t := temporal.Node(tid)
	return temporal.TemporalLine{
		F: &f,
		T: &t,
		S: s,
		E: e,
	}
}

func idsOf(nodes []graph.Node) []int64 {
	ids := make([]int64, len(nodes))
	for i := range nodes {
		ids[i] = nodes[i].ID()
	}
	return ids
}

func TestEarliestArrivalFrom(t *testing.T) {
	tests := []struct {
		name string
		stream []struct{
			F, T int64
			S, E uint64
		}
		nodes map[int64]struct{
			earliest uint64
			path     []int64
		}
		from  int64
		at    uint64
		until uint64
	}{
		{
			"sketchStream",
			sketchStream,
			map[int64]struct{
				earliest uint64
				path     []int64
			}{
				1: {0, []int64{1}},
				2: {1, []int64{1,2}},
				3: {1, []int64{1,3}},
				4: {2, []int64{1,4}},
				5: {4, []int64{1,2,5}},
				6: {3, []int64{1,2,6}},
				8: {5, []int64{1,2,5,8}},
				9: {8, []int64{1,2,6,9}},
			},
			1,
			0,
			^uint64(0),
		},
		{
			"wuChengHuangKeLuXuStream",
			wuChengHuangKeLuXuStream,
			map[int64]struct{
				earliest uint64
				path     []int64
			}{
				1:  {0, []int64{1}},
				2:  {1, []int64{1,2}},
				3:  {4, []int64{1,3}},
				6:  {3, []int64{1,6}},
				7:  {3, []int64{1,2,7}},
				8:  {3, []int64{1,2,8}},
				9:  {5, []int64{1,6,9}},
				11: {6, []int64{1,2,7,11}},
				12: {8, []int64{1,6,9,12}},
			},
			1,
			0,
			^uint64(0),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lines := make([]graph.TemporalLine, len(test.stream))
			for i, line := range test.stream {
				lines[i] = lineAt(line.F, line.T, line.S, line.E)
			}
			it := temporal.NewLineStream(lines)
			from := temporal.Node(test.from)
			out := EarliestArrivalFrom(it, &from, test.at, test.until)
			if out.Len() != len(test.nodes) {
				t.Fatalf("unexpected reachable node: expected %v, got %v", len(test.nodes), out.Len())
			}
			for vid, ev := range test.nodes {
				path, earliest := out.To(vid)
				if path == nil {
					t.Fatalf("expected reachable node: %v", vid)
				}
				if earliest != ev.earliest {
					t.Fatalf("expected earliest arrival time: (%v) expected %v, got %v", vid, ev.earliest, earliest)
				}
				if len(path) != len(ev.path) || !reflect.DeepEqual(idsOf(path), ev.path) {
					t.Fatalf("unexpected earliest arrival path: (%v) expected %v, got %v", vid, ev.path, idsOf(path))
				}
			}
		})
	}
}
