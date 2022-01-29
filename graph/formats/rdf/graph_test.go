// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf_test

import (
	"io"
	"math"
	"sort"
	"strings"
	"testing"

	"golang.org/x/exp/rand"

	"github.com/google/go-cmp/cmp"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/formats/rdf"
	"gonum.org/v1/gonum/graph/internal/set"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/multi"
	"gonum.org/v1/gonum/graph/testgraph"
)

func graphBuilder(nodes []graph.Node, edges []testgraph.WeightedLine, _, _ float64) (g graph.Graph, n []graph.Node, e []testgraph.Edge, s, a float64, ok bool) {
	seen := set.NewNodes()
	gg := rdf.NewGraph()
	for _, n := range nodes {
		seen.Add(n)
		gg.AddNode(n)
	}
	for _, edge := range edges {
		f := gg.Node(edge.From().ID())
		if f == nil {
			f = edge.From()
		}
		t := gg.Node(edge.To().ID())
		if t == nil {
			t = edge.To()
		}
		cl := multi.Line{F: f, T: t, UID: edge.ID()}
		seen.Add(cl.F)
		seen.Add(cl.T)
		e = append(e, cl)
		gg.SetLine(cl)
	}
	if len(seen) != 0 {
		n = make([]graph.Node, 0, len(seen))
	}
	for _, sn := range seen {
		n = append(n, sn)
	}
	return gg, n, e, math.NaN(), math.NaN(), true
}

const (
	usesEmpty     = true
	reversesEdges = false // Not tested since Graph is directed, but included for documentation.
)

func TestGraph(t *testing.T) {
	t.Run("EdgeExistence", func(t *testing.T) {
		testgraph.EdgeExistence(t, graphBuilder, reversesEdges)
	})
	t.Run("LineExistence", func(t *testing.T) {
		testgraph.LineExistence(t, graphBuilder, usesEmpty, reversesEdges)
	})
	t.Run("NodeExistence", func(t *testing.T) {
		testgraph.NodeExistence(t, graphBuilder)
	})
	t.Run("ReturnAdjacentNodes", func(t *testing.T) {
		testgraph.ReturnAdjacentNodes(t, graphBuilder, usesEmpty, reversesEdges)
	})
	t.Run("ReturnAllLines", func(t *testing.T) {
		testgraph.ReturnAllLines(t, graphBuilder, usesEmpty)
	})
	t.Run("ReturnAllNodes", func(t *testing.T) {
		testgraph.ReturnAllNodes(t, graphBuilder, usesEmpty)
	})
	t.Run("ReturnNodeSlice", func(t *testing.T) {
		testgraph.ReturnNodeSlice(t, graphBuilder, usesEmpty)
	})

	t.Run("RemoveNodes", func(t *testing.T) {
		g := multi.NewDirectedGraph()
		it := testgraph.NewRandomNodes(100, 1, func(id int64) graph.Node { return multi.Node(id) })
		for it.Next() {
			g.AddNode(it.Node())
		}
		it.Reset()
		rnd := rand.New(rand.NewSource(1))
		for it.Next() {
			u := it.Node()
			d := rnd.Intn(5)
			vit := g.Nodes()
			for d >= 0 && vit.Next() {
				v := vit.Node()
				d--
				g.SetLine(g.NewLine(u, v))
			}
		}
		testgraph.RemoveNodes(t, g)
	})
	t.Run("AddLines", func(t *testing.T) {
		testgraph.AddLines(t, 100,
			multi.NewDirectedGraph(),
			func(id int64) graph.Node { return multi.Node(id) },
			true, // Can update nodes.
		)
	})
	t.Run("RemoveLines", func(t *testing.T) {
		g := multi.NewDirectedGraph()
		it := testgraph.NewRandomNodes(100, 1, func(id int64) graph.Node { return multi.Node(id) })
		for it.Next() {
			g.AddNode(it.Node())
		}
		it.Reset()
		var lines []graph.Line
		rnd := rand.New(rand.NewSource(1))
		for it.Next() {
			u := it.Node()
			d := rnd.Intn(5)
			vit := g.Nodes()
			for d >= 0 && vit.Next() {
				v := vit.Node()
				d--
				l := g.NewLine(u, v)
				g.SetLine(l)
				lines = append(lines, l)
			}
		}
		rnd.Shuffle(len(lines), func(i, j int) {
			lines[i], lines[j] = lines[j], lines[i]
		})
		testgraph.RemoveLines(t, g, iterator.NewOrderedLines(lines))
	})
}

var removeStatementTests = []struct {
	name    string
	triples string
	remove  int
	want    string
}{
	{
		name: "triangle",
		triples: `
<ex:a> <ex:p> <ex:b> .
<ex:b> <ex:p> <ex:c> .
<ex:c> <ex:p> <ex:a> .
`,
		remove: 0,
		want: `
<ex:b> <ex:p> <ex:c> .
<ex:c> <ex:p> <ex:a> .
`,
	},
	{
		name: "star",
		triples: `
<ex:a> <ex:p> <ex:b> .
<ex:a> <ex:p> <ex:c> .
<ex:a> <ex:p> <ex:d> .
`,
		remove: 0,
		want: `
<ex:a> <ex:p> <ex:c> .
<ex:a> <ex:p> <ex:d> .
`,
	},
	{
		name: "loop",
		triples: `
<ex:a> <ex:p> <ex:b> .
<ex:b> <ex:p> <ex:a> .
`,
		remove: 0,
		want: `
<ex:b> <ex:p> <ex:a> .
`,
	},
	{
		name: "parallel",
		triples: `
<ex:a> <ex:p> <ex:b> .
<ex:a> <ex:q> <ex:b> .
`,
		remove: 0,
		want: `
<ex:a> <ex:q> <ex:b> .
`,
	},
	{
		name: "dumbell",
		triples: `
<ex:a> <ex:p> <ex:b> .
<ex:b> <ex:p> <ex:c> .
<ex:c> <ex:p> <ex:a> .

<ex:a> <ex:p> <ex:d> .

<ex:d> <ex:p> <ex:e> .
<ex:e> <ex:p> <ex:f> .
<ex:f> <ex:p> <ex:d> .
`,
		remove: 3,
		want: `
<ex:a> <ex:p> <ex:b> .
<ex:b> <ex:p> <ex:c> .
<ex:c> <ex:p> <ex:a> .
<ex:d> <ex:p> <ex:e> .
<ex:e> <ex:p> <ex:f> .
<ex:f> <ex:p> <ex:d> .
`,
	},
}

func TestRemoveStatement(t *testing.T) {
	for _, test := range removeStatementTests {
		g, statements, err := graphFromReader(strings.NewReader(test.triples))
		if err != nil {
			t.Errorf("unexpected error for %q: %v", test.name, err)
		}

		g.RemoveStatement(statements[test.remove])

		var gotStatements []string
		it := g.AllStatements()
		for it.Next() {
			gotStatements = append(gotStatements, it.Statement().String())
		}
		sort.Strings(gotStatements)

		got := strings.TrimSpace(strings.Join(gotStatements, "\n"))
		want := strings.TrimSpace(test.want)

		if got != want {
			t.Errorf("unexpected result for %q:\n%s", test.name, cmp.Diff(got, want))
		}
	}
}

var removeTermTests = []struct {
	name    string
	triples string
	remove  string
	want    string
}{
	{
		name: "triangle",
		triples: `
<ex:a> <ex:p> <ex:b> .
<ex:b> <ex:p> <ex:c> .
<ex:c> <ex:p> <ex:a> .
`,
		remove: "<ex:a>",
		want: `
<ex:b> <ex:p> <ex:c> .
`,
	},
	{
		name: "star",
		triples: `
<ex:a> <ex:p> <ex:b> .
<ex:a> <ex:p> <ex:c> .
<ex:a> <ex:p> <ex:d> .
`,
		remove: "<ex:a>",
		want:   "",
	},
	{
		name: "loop",
		triples: `
<ex:a> <ex:p> <ex:b> .
<ex:b> <ex:p> <ex:a> .
`,
		remove: "<ex:a>",
		want:   "",
	},
	{
		name: "parallel",
		triples: `
<ex:a> <ex:p> <ex:b> .
<ex:a> <ex:q> <ex:b> .
`,
		remove: "<ex:a>",
		want:   "",
	},
	{
		name: "dumbell_1",
		triples: `
<ex:a> <ex:p> <ex:b> .
<ex:b> <ex:p> <ex:c> .
<ex:c> <ex:p> <ex:a> .

<ex:a> <ex:q> <ex:d> .

<ex:d> <ex:p> <ex:e> .
<ex:e> <ex:p> <ex:f> .
<ex:f> <ex:p> <ex:d> .
`,
		remove: "<ex:q>",
		want: `
<ex:a> <ex:p> <ex:b> .
<ex:b> <ex:p> <ex:c> .
<ex:c> <ex:p> <ex:a> .
<ex:d> <ex:p> <ex:e> .
<ex:e> <ex:p> <ex:f> .
<ex:f> <ex:p> <ex:d> .
`,
	},
	{
		name: "dumbell_2",
		triples: `
<ex:a> <ex:p> <ex:b> .
<ex:b> <ex:p> <ex:c> .
<ex:c> <ex:p> <ex:a> .

<ex:a> <ex:q> <ex:d> .

<ex:d> <ex:p> <ex:e> .
<ex:e> <ex:p> <ex:f> .
<ex:f> <ex:p> <ex:d> .
`,
		remove: "<ex:p>",
		want: `
<ex:a> <ex:q> <ex:d> .
`,
	},
}

func TestRemoveTerm(t *testing.T) {
	for _, test := range removeTermTests {
		g, _, err := graphFromReader(strings.NewReader(test.triples))
		if err != nil {
			t.Errorf("unexpected error for %q: %v", test.name, err)
		}

		term, ok := g.TermFor(test.remove)
		if !ok {
			t.Errorf("couldn't find expected term for %q", test.name)
			continue
		}
		g.RemoveTerm(term)

		var gotStatements []string
		it := g.AllStatements()
		for it.Next() {
			gotStatements = append(gotStatements, it.Statement().String())
		}
		sort.Strings(gotStatements)

		got := strings.TrimSpace(strings.Join(gotStatements, "\n"))
		want := strings.TrimSpace(test.want)

		if got != want {
			t.Errorf("unexpected result for %q:\n%s", test.name, cmp.Diff(got, want))
		}
	}
}

func graphFromReader(r io.Reader) (*rdf.Graph, []*rdf.Statement, error) {
	g := rdf.NewGraph()
	var statements []*rdf.Statement
	dec := rdf.NewDecoder(r)
	for {
		s, err := dec.Unmarshal()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, err
		}
		statements = append(statements, s)
		g.AddStatement(s)
	}
	return g, statements, nil
}
