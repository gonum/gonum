// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iterator_test

import (
	"reflect"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/internal/order"
)

type line struct{ f, t, id int64 }

func (l line) From() graph.Node         { return simple.Node(l.f) }
func (l line) To() graph.Node           { return simple.Node(l.t) }
func (l line) ReversedLine() graph.Line { l.f, l.t = l.t, l.f; return l }
func (l line) ID() int64                { return l.id }

var linesTests = []struct {
	lines map[int64]graph.Line
}{
	{lines: nil},
	{lines: map[int64]graph.Line{1: line{f: 1, t: 2, id: 1}}},
	{lines: map[int64]graph.Line{1: line{f: 1, t: 2, id: 1}, 2: line{f: 2, t: 3, id: 2}, 3: line{f: 3, t: 4, id: 3}, 4: line{f: 4, t: 5, id: 4}}},
	{lines: map[int64]graph.Line{4: line{f: 5, t: 4, id: 4}, 3: line{f: 4, t: 3, id: 3}, 2: line{f: 3, t: 2, id: 2}, 1: line{f: 2, t: 1, id: 1}}},
}

func TestLinesIterate(t *testing.T) {
	for _, test := range linesTests {
		it := iterator.NewLines(test.lines)
		for i := 0; i < 2; i++ {
			if it.Len() != len(test.lines) {
				t.Errorf("unexpected iterator length for round %d: got:%d want:%d", i, it.Len(), len(test.lines))
			}
			var got map[int64]graph.Line
			if it.Len() != 0 {
				got = make(map[int64]graph.Line)
			}
			for it.Next() {
				got[it.Line().ID()] = it.Line()
			}
			want := test.lines
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, want)
			}
			it.Reset()
		}
	}
}

func TestLinesSlice(t *testing.T) {
	for _, test := range linesTests {
		it := iterator.NewLines(test.lines)
		for i := 0; i < 2; i++ {
			got := it.LineSlice()
			var want []graph.Line
			for _, l := range test.lines {
				want = append(want, l)
			}
			order.ByID(got)
			order.ByID(want)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, want)
			}
			it.Reset()
		}
	}
}

var orderedLinesTests = []struct {
	lines []graph.Line
}{
	{lines: nil},
	{lines: []graph.Line{line{f: 1, t: 2}}},
	{lines: []graph.Line{line{f: 1, t: 2}, line{f: 2, t: 3}, line{f: 3, t: 4}, line{f: 4, t: 5}}},
	{lines: []graph.Line{line{f: 5, t: 4}, line{f: 4, t: 3}, line{f: 3, t: 2}, line{f: 2, t: 1}}},
}

func TestOrderedLinesIterate(t *testing.T) {
	for _, test := range orderedLinesTests {
		it := iterator.NewOrderedLines(test.lines)
		for i := 0; i < 2; i++ {
			if it.Len() != len(test.lines) {
				t.Errorf("unexpected iterator length for round %d: got:%d want:%d", i, it.Len(), len(test.lines))
			}
			var got []graph.Line
			for it.Next() {
				got = append(got, it.Line())
			}
			want := test.lines
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, want)
			}
			it.Reset()
		}
	}
}

func TestOrderedLinesSlice(t *testing.T) {
	for _, test := range orderedLinesTests {
		it := iterator.NewOrderedLines(test.lines)
		for i := 0; i < 2; i++ {
			got := it.LineSlice()
			want := test.lines
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, want)
			}
			it.Reset()
		}
	}
}

type weightedLine struct {
	f, t, id int64
	w        float64
}

func (l weightedLine) From() graph.Node         { return simple.Node(l.f) }
func (l weightedLine) To() graph.Node           { return simple.Node(l.t) }
func (l weightedLine) ReversedLine() graph.Line { l.f, l.t = l.t, l.f; return l }
func (l weightedLine) Weight() float64          { return l.w }
func (l weightedLine) ID() int64                { return l.id }

var weightedLinesTests = []struct {
	lines map[int64]graph.WeightedLine
}{
	{lines: nil},
	{lines: map[int64]graph.WeightedLine{2: weightedLine{f: 1, t: 2, w: 1, id: 2}}},
	{lines: map[int64]graph.WeightedLine{2: weightedLine{f: 1, t: 2, w: 1, id: 2}, 4: weightedLine{f: 2, t: 3, w: 2, id: 4}, 6: weightedLine{f: 3, t: 4, w: 3, id: 6}, 8: weightedLine{f: 4, t: 5, w: 4, id: 8}}},
	{lines: map[int64]graph.WeightedLine{8: weightedLine{f: 5, t: 4, w: 4, id: 8}, 6: weightedLine{f: 4, t: 3, w: 3, id: 6}, 4: weightedLine{f: 3, t: 2, w: 2, id: 4}, 2: weightedLine{f: 2, t: 1, w: 1, id: 2}}},
}

func TestWeightedLinesIterate(t *testing.T) {
	for _, test := range weightedLinesTests {
		it := iterator.NewWeightedLines(test.lines)
		for i := 0; i < 2; i++ {
			if it.Len() != len(test.lines) {
				t.Errorf("unexpected iterator length for round %d: got:%d want:%d", i, it.Len(), len(test.lines))
			}
			var got map[int64]graph.WeightedLine
			if it.Len() != 0 {
				got = make(map[int64]graph.WeightedLine)
			}
			for it.Next() {
				got[it.WeightedLine().ID()] = it.WeightedLine()
			}
			want := test.lines
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, want)
			}
			it.Reset()
		}
	}
}

func TestWeightedLinesSlice(t *testing.T) {
	for _, test := range weightedLinesTests {
		it := iterator.NewWeightedLines(test.lines)
		for i := 0; i < 2; i++ {
			got := it.WeightedLineSlice()
			var want []graph.WeightedLine
			for _, l := range test.lines {
				want = append(want, l)
			}
			order.ByID(got)
			order.ByID(want)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, want)
			}
			it.Reset()
		}
	}
}

var orderedWeightedLinesTests = []struct {
	lines []graph.WeightedLine
}{
	{lines: nil},
	{lines: []graph.WeightedLine{weightedLine{f: 1, t: 2, w: 1}}},
	{lines: []graph.WeightedLine{weightedLine{f: 1, t: 2, w: 1}, weightedLine{f: 2, t: 3, w: 2}, weightedLine{f: 3, t: 4, w: 3}, weightedLine{f: 4, t: 5, w: 4}}},
	{lines: []graph.WeightedLine{weightedLine{f: 5, t: 4, w: 4}, weightedLine{f: 4, t: 3, w: 3}, weightedLine{f: 3, t: 2, w: 2}, weightedLine{f: 2, t: 1, w: 1}}},
}

func TestOrderedWeightedLinesIterate(t *testing.T) {
	for _, test := range orderedWeightedLinesTests {
		it := iterator.NewOrderedWeightedLines(test.lines)
		for i := 0; i < 2; i++ {
			if it.Len() != len(test.lines) {
				t.Errorf("unexpected iterator length for round %d: got:%d want:%d", i, it.Len(), len(test.lines))
			}
			var got []graph.WeightedLine
			for it.Next() {
				got = append(got, it.WeightedLine())
			}
			want := test.lines
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, want)
			}
			it.Reset()
		}
	}
}

func TestOrderedWeightedLinesSlice(t *testing.T) {
	for _, test := range orderedWeightedLinesTests {
		it := iterator.NewOrderedWeightedLines(test.lines)
		for i := 0; i < 2; i++ {
			got := it.WeightedLineSlice()
			want := test.lines
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, want)
			}
			it.Reset()
		}
	}
}
