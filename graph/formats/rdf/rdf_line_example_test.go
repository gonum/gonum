// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf_test

import (
	"fmt"
	"log"
	"strings"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/formats/rdf"
	"gonum.org/v1/gonum/graph/multi"
)

// foodNode implements graph.Node, dot.Node and encoding.Attributer
// to allow the RDF term value to be given to the DOT encoder.
type foodNode struct {
	rdf.Term
}

func (n foodNode) DOTID() string {
	text, _, kind, err := n.Term.Parts()
	if err != nil {
		return fmt.Sprintf("error:%s", n.Term.Value)
	}
	switch kind {
	case rdf.Blank:
		return n.Term.Value
	case rdf.IRI:
		return text
	case rdf.Literal:
		return fmt.Sprintf("%q", text)
	default:
		return fmt.Sprintf("invalid:%s", n.Term.Value)
	}
}

func (n foodNode) Attributes() []encoding.Attribute {
	_, qual, _, err := n.Term.Parts()
	if err != nil {
		return []encoding.Attribute{{Key: "error", Value: err.Error()}}
	}
	if qual == "" {
		return nil
	}
	parts := strings.Split(qual, ":")
	return []encoding.Attribute{{Key: parts[0], Value: parts[1]}}
}

// foodLine implements graph.Line and encoding.Attributer to
// allow the line's RDF term value to be given to the DOT
// encoder and for the nodes to be shimmed to the foodNode
// type.
//
// It also implements line reversal for the semantics of
// a food web with some taxonomic information.
type foodLine struct {
	*rdf.Statement
}

func (l foodLine) From() graph.Node { return foodNode{l.Subject} }
func (l foodLine) To() graph.Node   { return foodNode{l.Object} }
func (l foodLine) ReversedLine() graph.Line {
	if l.Predicate.Value == "<tax:is>" {
		// This should remain unreversed, so return as is.
		return l
	}
	s := *l.Statement
	// Reverse the line end points.
	s.Subject, s.Object = s.Object, s.Subject
	// Invert the semantics of the predicate.
	switch s.Predicate.Value {
	case "<eco:eats>":
		s.Predicate.Value = "<eco:eaten-by>"
	case "<eco:eaten-by>":
		s.Predicate.Value = "<eco:eats>"
	case "<tax:is-a>":
		s.Predicate.Value = "<tax:includes>"
	case "<tax:includes>":
		s.Predicate.Value = "<tax:is-a>"
	default:
		panic("invalid predicate")
	}
	// All IDs returned by the RDF parser are positive, so
	// sign reverse the edge ID to avoid any collisions.
	s.Predicate.UID *= -1
	return foodLine{&s}
}

func (l foodLine) Attributes() []encoding.Attribute {
	text, _, _, err := l.Predicate.Parts()
	if err != nil {
		return []encoding.Attribute{{Key: "error", Value: err.Error()}}
	}
	parts := strings.Split(text, ":")
	return []encoding.Attribute{{Key: parts[0], Value: parts[1]}}
}

// expand copies src into dst, adding the reversal of each line if it is
// distinct.
func expand(dst, src *multi.DirectedGraph) {
	it := src.Edges()
	for it.Next() {
		lit := it.Edge().(multi.Edge)
		for lit.Next() {
			l := lit.Line()
			r := l.ReversedLine()
			dst.SetLine(l)
			if l == r {
				continue
			}
			dst.SetLine(r)
		}
	}
}

func ExampleStatement_ReversedLine() {
	const statements = `
_:wolf <tax:is-a> _:animal .
_:wolf <tax:is> "Wolf"^^<tax:common> .
_:wolf <tax:is> "Canis lupus"^^<tax:binomial> .
_:wolf <eco:eats> _:sheep .
_:sheep <tax:is-a> _:animal .
_:sheep <tax:is> "Sheep"^^<tax:common> .
_:sheep <tax:is> "Ovis aries"^^<tax:binomial> .
_:sheep <eco:eats> _:grass .
_:grass <tax:is-a> _:plant .
_:grass <tax:is> "Grass"^^<tax:common> .
_:grass <tax:is> "Lolium perenne"^^<tax:binomial> .
_:grass <tax:is> "Festuca rubra"^^<tax:binomial> .
_:grass <tax:is> "Poa pratensis"^^<tax:binomial> .
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
		g.SetLine(foodLine{l})
	}

	h := multi.NewDirectedGraph()
	expand(h, g)

	// Marshal the graph into DOT.
	b, err := dot.MarshalMulti(h, "food web", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n\n", b)

	// Output:
	//
	// digraph "food web" {
	// 	// Node definitions.
	// 	"_:wolf";
	// 	"_:animal";
	// 	"Wolf" [tax=common];
	// 	"Canis lupus" [tax=binomial];
	// 	"_:sheep";
	// 	"Sheep" [tax=common];
	// 	"Ovis aries" [tax=binomial];
	// 	"_:grass";
	// 	"_:plant";
	// 	"Grass" [tax=common];
	// 	"Lolium perenne" [tax=binomial];
	// 	"Festuca rubra" [tax=binomial];
	// 	"Poa pratensis" [tax=binomial];
	//
	// 	// Edge definitions.
	// 	"_:wolf" -> "_:animal" [tax="is-a"];
	// 	"_:wolf" -> "Wolf" [tax=is];
	// 	"_:wolf" -> "Canis lupus" [tax=is];
	// 	"_:wolf" -> "_:sheep" [eco=eats];
	// 	"_:animal" -> "_:wolf" [tax=includes];
	// 	"_:animal" -> "_:sheep" [tax=includes];
	// 	"_:sheep" -> "_:wolf" [eco="eaten-by"];
	// 	"_:sheep" -> "_:animal" [tax="is-a"];
	// 	"_:sheep" -> "Sheep" [tax=is];
	// 	"_:sheep" -> "Ovis aries" [tax=is];
	// 	"_:sheep" -> "_:grass" [eco=eats];
	// 	"_:grass" -> "_:sheep" [eco="eaten-by"];
	// 	"_:grass" -> "_:plant" [tax="is-a"];
	// 	"_:grass" -> "Grass" [tax=is];
	// 	"_:grass" -> "Lolium perenne" [tax=is];
	// 	"_:grass" -> "Festuca rubra" [tax=is];
	// 	"_:grass" -> "Poa pratensis" [tax=is];
	// 	"_:plant" -> "_:grass" [tax=includes];
	// }
}
