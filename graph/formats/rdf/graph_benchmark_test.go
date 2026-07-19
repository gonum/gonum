// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf_test

import (
	"io"
	"strings"
	"testing"

	"gonum.org/v1/gonum/graph/formats/rdf"
)

func BenchmarkQuery(b *testing.B) {
	g := rdf.NewGraph()
	dec := rdf.NewDecoder(strings.NewReader(gods))
	for {
		s, err := dec.Unmarshal()
		if err != nil {
			if err != io.EOF {
				b.Fatalf("error during decoding: %v", err)
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

	b.Run("grandfather", func(b *testing.B) {
		for range b.N {
			got := len(heracles.Out(father).Out(father).Out(name).Result())
			if got != 1 {
				b.Fatalf("unexpected result count: got %d, want %d", got, 1)
			}
		}
	})

	b.Run("lineage", func(b *testing.B) {
		for range b.N {
			var got int
			heracles.Repeat(func(q rdf.Query) (rdf.Query, bool) {
				q = q.Out(father)
				got++
				return q, true
			})
			if got != 3 {
				b.Fatalf("unexpected generation result count: got %d, want %d", got, 3)
			}
		}
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

	b.Run("parents", func(b *testing.B) {
		for range b.N {
			got := len(heracles.Out(parents).Out(typ).Result())
			if got != 2 {
				b.Fatalf("unexpected result count: got %d, want %d", got, 1)
			}
		}
	})

	// battled is a helper filter for queries below.
	battled := func(s *rdf.Statement) bool {
		// Traverse across <a:battled>.
		return s.Predicate.Value == "<a:battled>"
	}

	b.Run("battled", func(b *testing.B) {
		for range b.N {
			got := len(heracles.Out(battled).Out(typ).Result())
			if got != 3 {
				b.Fatalf("unexpected result count: got %d, want %d", got, 1)
			}
		}
	})

	b.Run("allies", func(b *testing.B) {
		for range b.N {
			got := len(heracles.Out(battled).In(battled).Not(heracles).Out(name).Result())
			if got != 1 {
				b.Fatalf("unexpected result count: got %d, want %d", got, 1)
			}
		}
	})

	h, ok := g.TermFor(`"Hades"`)
	if !ok {
		b.Fatal("could not find term for Hades")
	}
	hades := g.Query(h).In(name)
	lives := func(s *rdf.Statement) bool {
		// Traverse across <a:lives>.
		return s.Predicate.Value == "<a:lives>"
	}
	location := func(s *rdf.Statement) bool {
		// Traverse across <p:location>.
		return s.Predicate.Value == "<p:location>"
	}

	b.Run("lives_with", func(b *testing.B) {
		for range b.N {
			got := len(hades.Out(lives).Out(location).In(location).In(lives).Not(hades).Out(name).Result())
			if got != 1 {
				b.Fatalf("unexpected result count: got %d, want %d", got, 1)
			}
		}
	})

	// g.V(hades).out('brother').as('god').out('lives').as('place').select('god', 'place').by('name')
	brother := func(s *rdf.Statement) bool {
		// Traverse across <h:brother>.
		return s.Predicate.Value == "<h:brother>"
	}

	b.Run("brothers", func(b *testing.B) {
		for range b.N {
			var got int
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
				if len(m) != 3 {
					b.Fatalf("unexpected sub result: got %v want three elements", m)
				}
				got++
			}
			if got != 2 {
				b.Fatalf("unexpected result count: got %d, want %d", got, 2)
			}
		}
	})
}
