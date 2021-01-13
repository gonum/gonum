// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/formats/rdf"
	"gonum.org/v1/gonum/graph/multi"
)

// dotNode implements graph.Node and dot.Node to allow the
// RDF term value to be given to the DOT encoder.
type dotNode struct {
	rdf.Term
}

func (n dotNode) DOTID() string { return n.Term.Value }

// dotLine implements graph.Line and encoding.Attributer to
// allow the line's RDF term value to be given to the DOT
// encoder and for the nodes to be shimmed to the dotNode
// type.
//
// Because the graph here is directed and we are not performing
// any line reversals, it is safe not to implement the
// ReversedLine method on dotLine; it will never be called.
type dotLine struct {
	*rdf.Statement
}

func (l dotLine) From() graph.Node { return dotNode{l.Subject} }
func (l dotLine) To() graph.Node   { return dotNode{l.Object} }

func (l dotLine) Attributes() []encoding.Attribute {
	return []encoding.Attribute{{Key: "label", Value: l.Predicate.Value}}
}

func Example_graph() {
	const statements = `
_:alice <http://xmlns.com/foaf/0.1/knows> _:bob .
_:alice <http://xmlns.com/foaf/0.1/givenName> "Alice" .
_:alice <http://xmlns.com/foaf/0.1/familyName> "Smith" .
_:bob <http://xmlns.com/foaf/0.1/knows> _:alice .
_:bob <http://xmlns.com/foaf/0.1/givenName> "Bob" .
_:bob <http://xmlns.com/foaf/0.1/familyName> "Smith" .
`

	// Decode the statement stream and insert the lines into a multigraph.
	g := multi.NewDirectedGraph()
	dec := rdf.NewDecoder(strings.NewReader(statements))
	for {
		l, err := dec.Unmarshal()
		if err != nil {
			break
		}

		// Wrap the line with a shim type to allow the RDF values
		// to be passed to the DOT marshaling routine.
		g.SetLine(dotLine{l})
	}

	// Marshal the graph into DOT.
	b, err := dot.MarshalMulti(g, "smiths", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n\n", b)

	// Get the ID look-up table.
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', 0)
	fmt.Fprintln(w, "Term\tID")
	for t, id := range dec.Terms() {
		fmt.Fprintf(w, "%s\t%d\n", t, id)
	}
	w.Flush()

	// Unordered output:
	//
	// digraph smiths {
	// 	// Node definitions.
	// 	"_:alice";
	// 	"_:bob";
	// 	"Alice";
	// 	"Smith";
	// 	"Bob";
	//
	// 	// Edge definitions.
	// 	"_:alice" -> "_:bob" [label=<http://xmlns.com/foaf/0.1/knows>];
	// 	"_:alice" -> "Alice" [label=<http://xmlns.com/foaf/0.1/givenName>];
	// 	"_:alice" -> "Smith" [label=<http://xmlns.com/foaf/0.1/familyName>];
	// 	"_:bob" -> "_:alice" [label=<http://xmlns.com/foaf/0.1/knows>];
	// 	"_:bob" -> "Smith" [label=<http://xmlns.com/foaf/0.1/familyName>];
	// 	"_:bob" -> "Bob" [label=<http://xmlns.com/foaf/0.1/givenName>];
	// }
	//
	// Term                                   ID
	// _:alice                                1
	// _:bob                                  2
	// <http://xmlns.com/foaf/0.1/knows>      3
	// "Alice"                                4
	// <http://xmlns.com/foaf/0.1/givenName>  5
	// "Smith"                                6
	// <http://xmlns.com/foaf/0.1/familyName> 7
	// "Bob"                                  8
}
