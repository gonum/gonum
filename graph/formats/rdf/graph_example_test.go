// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"gonum.org/v1/gonum/graph/formats/rdf"
)

func ExampleGraph() {
	f, err := os.Open("path/to/graph.nq")
	if err != nil {
		log.Fatal(err)
	}

	dec := rdf.NewDecoder(f)
	var statements []*rdf.Statement
	for {
		s, err := dec.Unmarshal()
		if err != nil {
			if err != io.EOF {
				log.Fatalf("error during decoding: %v", err)
			}
			break
		}

		// Statements can be filtered at this point to exclude unwanted
		// or irrelevant parts of the graph.
		statements = append(statements, s)
	}
	f.Close()

	// Canonicalize blank nodes to reduce memory footprint.
	statements, err = rdf.URDNA2015(statements, statements)
	if err != nil {
		log.Fatal(err)
	}

	g := rdf.NewGraph()
	for _, s := range statements {
		g.AddStatement(s)
	}

	// Do something with the graph.
}

const gods = `
_:alcmene <l:type> "human" .
_:alcmene <p:name> "Alcmene" .
_:cerberus <a:lives> _:cerberushome .
_:cerberus <l:type> "monster" .
_:cerberus <p:name> "Cerberus" .
_:cerberushome <p:location> _:tartarus .
_:cronos <l:type> "titan" .
_:cronos <p:name> "Cronos" .
_:hades <a:lives> _:hadeshome .
_:hades <h:brother> _:poseidon .
_:hades <h:brother> _:zeus .
_:hades <h:pet> _:cerberus .
_:hades <l:type> "god" .
_:hades <p:name> "Hades" .
_:hadeshome <p:location> _:tartarus .
_:hadeshome <p:reason> "it is peaceful" .
_:heracles <a:battled> _:cerberus .
_:heracles <a:battled> _:hydra .
_:heracles <a:battled> _:nemean .
_:heracles <h:father> _:zeus .
_:heracles <h:mother> _:alcmene .
_:heracles <l:type> "demigod" .
_:heracles <p:name> "Heracles" .
_:hydra <l:type> "monster" .
_:hydra <p:name> "Lernean Hydra" .
_:nemean <l:type> "monster" .
_:nemean <p:name> "Nemean Lion" .
_:olympus <l:type> "location" .
_:olympus <p:name> "Olympus" .
_:poseidon <a:lives> _:poseidonhome .
_:poseidon <h:brother> _:hades .
_:poseidon <h:brother> _:zeus .
_:poseidon <l:type> "god" .
_:poseidon <p:name> "Poseidon" .
_:poseidonhome <p:location> _:sea .
_:poseidonhome <p:reason> "it was given to him" .
_:sea <l:type> "location" .
_:sea <p:name> "Sea" .
_:tartarus <l:type> "location" .
_:tartarus <p:name> "Tartarus" .
_:theseus <a:battled> _:cerberus .
_:theseus <h:father> _:poseidon .
_:theseus <l:type> "human" .
_:theseus <p:name> "Theseus" .
_:zeus <a:lives> _:zeushome .
_:zeus <h:brother> _:hades .
_:zeus <h:brother> _:poseidon .
_:zeus <h:father> _:cronos .
_:zeus <l:type> "god" .
_:zeus <p:name> "Zeus" .
_:zeushome <p:location> _:olympus .
_:zeushome <p:reason> "he can see everything" .
`

func ExampleQuery() {
	g := rdf.NewGraph()
	dec := rdf.NewDecoder(strings.NewReader(gods))
	for {
		s, err := dec.Unmarshal()
		if err != nil {
			if err != io.EOF {
				log.Fatalf("error during decoding: %v", err)
			}
			break
		}
		g.AddStatement(s)
	}

	it := g.Nodes()
	nodes := make([]rdf.Term, 0, it.Len())
	for it.Next() {
		nodes = append(nodes, it.Node().(rdf.Term))
	}

	// Construct a query start point. This can be reused. If a specific
	// node is already known it can be used to reduce the work required here.
	heracles := g.Query(nodes...).In(func(s *rdf.Statement) bool {
		// Traverse in from the name "Heracles".
		return s.Predicate.Value == "<p:name>" && s.Object.Value == `"Heracles"`
	})

	// father and name filter statements on their predicate values. These
	// are used in the queries that follow.
	father := func(s *rdf.Statement) bool {
		// Traverse across <h:father>.
		return s.Predicate.Value == "<h:father>"
	}
	name := func(s *rdf.Statement) bool {
		// Traverse across <p:name>.
		return s.Predicate.Value == "<p:name>"
	}

	// g.V().has('name', 'heracles').out('father').out('father').values('name')
	for _, r := range heracles.
		Out(father). // Traverse out across <h:father> to get to Zeus.
		Out(father). // and again to get to Cronos.
		Out(name).   // Retrieve the name by traversing the <p:name> edges.
		Result() {
		fmt.Printf("Heracles' grandfather: %s\n", r.Value)
	}

	// g.V().has('name', 'heracles').repeat(out('father')).emit().values('name')
	var i int
	heracles.Repeat(func(q rdf.Query) (rdf.Query, bool) {
		q = q.Out(father)
		for _, r := range q.Out(name).Result() {
			fmt.Printf("Heracles' lineage %d: %s\n", i, r.Value)
		}
		i++
		return q, true
	})

	// parents and typ are helper filters for queries below.
	parents := func(s *rdf.Statement) bool {
		// Traverse across <h:father> or <h:mother>
		return s.Predicate.Value == "<h:father>" || s.Predicate.Value == "<h:mother>"
	}
	typ := func(s *rdf.Statement) bool {
		// Traverse across <l:type>.
		return s.Predicate.Value == "<l:type>"
	}

	// g.V(heracles).out('father', 'mother').label()
	for _, r := range heracles.Out(parents).Out(typ).Result() {
		fmt.Printf("Heracles' parents' types: %s\n", r.Value)
	}

	// battled is a helper filter for queries below.
	battled := func(s *rdf.Statement) bool {
		// Traverse across <a:battled>.
		return s.Predicate.Value == "<a:battled>"
	}

	// g.V(heracles).out('battled').label()
	for _, r := range heracles.Out(battled).Out(typ).Result() {
		fmt.Printf("Heracles' antagonists' types: %s\n", r.Value)
	}

	// g.V(heracles).out('battled').valueMap()
	for _, r := range heracles.Out(battled).Result() {
		m := make(map[string]string)
		g.Query(r).Out(func(s *rdf.Statement) bool {
			// Store any p: namespace in the map.
			if strings.HasPrefix(s.Predicate.Value, "<p:") {
				prop := strings.TrimSuffix(strings.TrimPrefix(s.Predicate.Value, "<p:"), ">")
				m[prop] = s.Object.Value
			}
			// But don't store the result into the query.
			return false
		})
		fmt.Println(m)
	}

	// g.V(heracles).as('h').out('battled').in('battled').where(neq('h')).values('name')
	for _, r := range heracles.Out(battled).In(battled).Not(heracles).Out(name).Result() {
		fmt.Printf("Heracles' allies: %s\n", r.Value)
	}

	// Construct a query start point for Hades, this time using a restricted
	// starting set only including the name. It would also be possible to
	// start directly from a query with the term _:hades, but that depends
	// on the blank node identity, which may be altered, for example by
	// canonicalization.
	h, ok := g.TermFor(`"Hades"`)
	if !ok {
		log.Fatal("could not find term for Hades")
	}
	hades := g.Query(h).In(name)

	// g.V(hades).as('x').out('lives').in('lives').where(neq('x')).values('name')
	//
	// This is more complex with RDF since properties are encoded by
	// attachment to anonymous blank nodes, so we take two steps, the
	// first to the blank node for where Hades lives and then the second
	// to get the actual location.
	lives := func(s *rdf.Statement) bool {
		// Traverse across <a:lives>.
		return s.Predicate.Value == "<a:lives>"
	}
	location := func(s *rdf.Statement) bool {
		// Traverse across <p:location>.
		return s.Predicate.Value == "<p:location>"
	}
	for _, r := range hades.Out(lives).Out(location).In(location).In(lives).Not(hades).Out(name).Result() {
		fmt.Printf("Hades lives with: %s\n", r.Value)
	}

	// g.V(hades).out('brother').as('god').out('lives').as('place').select('god', 'place').by('name')
	brother := func(s *rdf.Statement) bool {
		// Traverse across <h:brother>.
		return s.Predicate.Value == "<h:brother>"
	}
	for _, r := range hades.Out(brother).Result() {
		m := make(map[string]string)
		as := func(key string) func(s *rdf.Statement) bool {
			return func(s *rdf.Statement) bool {
				// Store any <p:name> objects in the map.
				if s.Predicate.Value == "<p:name>" {
					m[key] = s.Object.Value
				}
				// But don't store the result into the query.
				return false
			}
		}
		sub := g.Query(r)
		sub.Out(as("god"))
		sub.Out(lives).Out(location).Out(as("place"))
		fmt.Println(m)
	}

	// The query above but with the reason for their choice.
	for _, r := range hades.Out(brother).Result() {
		m := make(map[string]string)
		// as stores the query result under the provided key
		// for m, and if cont is not nil, allows the chain
		// to continue.
		as := func(query, key string, cont func(s *rdf.Statement) bool) func(s *rdf.Statement) bool {
			return func(s *rdf.Statement) bool {
				// Store any objects matching the query in the map.
				if s.Predicate.Value == query {
					m[key] = s.Object.Value
				}
				// Continue with chain if cont is not nil and
				// the statement satisfies its condition.
				if cont == nil {
					return false
				}
				return cont(s)
			}
		}
		sub := g.Query(r)
		sub.Out(as("<p:name>", "god", nil))
		sub.Out(lives).
			Out(as("<p:reason>", "reason", location)).
			Out(as("<p:name>", "place", nil))
		fmt.Println(m)
	}

	// Unordered output:
	//
	// Heracles' grandfather: "Cronos"
	// Heracles' lineage 0: "Zeus"
	// Heracles' lineage 1: "Cronos"
	// Heracles' parents' types: "god"
	// Heracles' parents' types: "human"
	// Heracles' antagonists' types: "monster"
	// Heracles' antagonists' types: "monster"
	// Heracles' antagonists' types: "monster"
	// map[name:"Cerberus"]
	// map[name:"Lernean Hydra"]
	// map[name:"Nemean Lion"]
	// Heracles' allies: "Theseus"
	// Hades lives with: "Cerberus"
	// map[god:"Zeus" place:"Olympus"]
	// map[god:"Poseidon" place:"Sea"]
	// map[god:"Zeus" place:"Olympus" reason:"he can see everything"]
	// map[god:"Poseidon" place:"Sea" reason:"it was given to him"]
}
