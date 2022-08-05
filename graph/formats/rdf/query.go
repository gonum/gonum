// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf

import (
	"sort"

	"gonum.org/v1/gonum/graph"
)

// Query represents a step in an RDF graph query. The methods on Query
// provide a simple graph query language.
type Query struct {
	g graph.Directed

	terms []Term
}

// NewQuery returns a query of g starting from the given nodes.
// Queries may not be mixed between distinct graphs. The type of
// g must be comparable. Query operations only consider edges that
// are represented by a *Statement or is an edge with lines held
// in a graph.Lines with at least one *Statement.
func NewQuery(g graph.Directed, from ...Term) Query {
	return Query{g: g, terms: from}
}

// Query returns a query of the receiver starting from the given nodes.
// Queries may not be mixed between distinct graphs.
func (g *Graph) Query(from ...Term) Query {
	return Query{g: g, terms: from}
}

// Out returns a query holding nodes reachable out from the receiver's
// starting nodes via statements that satisfy fn.
func (q Query) Out(fn func(s *Statement) bool) Query {
	r := Query{g: q.g}
	for _, s := range q.terms {
		it := q.g.From(s.ID())
		for it.Next() {
			if ConnectedByAny(q.g.Edge(s.ID(), it.Node().ID()), fn) {
				r.terms = append(r.terms, it.Node().(Term))
			}
		}
	}
	return r
}

// In returns a query holding nodes reachable in from the receiver's
// starting nodes via statements that satisfy fn.
func (q Query) In(fn func(s *Statement) bool) Query {
	r := Query{g: q.g}
	for _, s := range q.terms {
		it := q.g.To(s.ID())
		for it.Next() {
			if ConnectedByAny(q.g.Edge(it.Node().ID(), s.ID()), fn) {
				r.terms = append(r.terms, it.Node().(Term))
			}
		}
	}
	return r
}

// HasAllOut returns a query holding nodes from the receiver's
// initial set where all outgoing statements satisfy fn. The
// query short circuits, so fn is not called after the first
// failure to match.
func (q Query) HasAllOut(fn func(s *Statement) bool) Query {
	r := Query{g: q.g}
	notFn := not(fn)
loop:
	for _, s := range q.terms {
		it := q.g.From(s.ID())
		for it.Next() {
			if ConnectedByAny(q.g.Edge(s.ID(), it.Node().ID()), notFn) {
				continue loop
			}
		}
		r.terms = append(r.terms, s)
	}
	return r
}

// HasAllIn returns a query holding nodes from the receiver's
// initial set where all incoming statements satisfy fn. The
// query short circuits, so fn is not called after the first
// failure to match.
func (q Query) HasAllIn(fn func(s *Statement) bool) Query {
	r := Query{g: q.g}
	notFn := not(fn)
loop:
	for _, s := range q.terms {
		it := q.g.To(s.ID())
		for it.Next() {
			if ConnectedByAny(q.g.Edge(it.Node().ID(), s.ID()), notFn) {
				continue loop
			}
		}
		r.terms = append(r.terms, s)
	}
	return r
}

// HasAnyOut returns a query holding nodes from the receiver's
// initial set where any outgoing statements satisfies fn. The
// query short circuits, so fn is not called after the first match.
func (q Query) HasAnyOut(fn func(s *Statement) bool) Query {
	r := Query{g: q.g}
	for _, s := range q.terms {
		it := q.g.From(s.ID())
		for it.Next() {
			if ConnectedByAny(q.g.Edge(s.ID(), it.Node().ID()), fn) {
				r.terms = append(r.terms, s)
				break
			}
		}
	}
	return r
}

// HasAnyIn returns a query holding nodes from the receiver's
// initial set where any incoming statements satisfies fn. The
// query short circuits, so fn is not called after the first match.
func (q Query) HasAnyIn(fn func(s *Statement) bool) Query {
	r := Query{g: q.g}
	for _, s := range q.terms {
		it := q.g.To(s.ID())
		for it.Next() {
			if ConnectedByAny(q.g.Edge(it.Node().ID(), s.ID()), fn) {
				r.terms = append(r.terms, s)
				break
			}
		}
	}
	return r
}

// not returns the negation of fn.
func not(fn func(s *Statement) bool) func(s *Statement) bool {
	return func(s *Statement) bool { return !fn(s) }
}

// And returns a query that holds the conjunction of q and p.
func (q Query) And(p Query) Query {
	if q.g != p.g {
		panic("rdf: binary query operation parameters from distinct graphs")
	}
	sortByID(q.terms)
	sortByID(p.terms)
	r := Query{g: q.g}
	var i, j int
	for i < len(q.terms) && j < len(p.terms) {
		qi := q.terms[i]
		pj := p.terms[j]
		switch {
		case qi.ID() < pj.ID():
			i++
		case pj.ID() < qi.ID():
			j++
		default:
			r.terms = append(r.terms, qi)
			i++
			j++
		}
	}
	return r
}

// Or returns a query that holds the disjunction of q and p.
func (q Query) Or(p Query) Query {
	if q.g != p.g {
		panic("rdf: binary query operation parameters from distinct graphs")
	}
	sortByID(q.terms)
	sortByID(p.terms)
	r := Query{g: q.g}
	var i, j int
	for i < len(q.terms) && j < len(p.terms) {
		qi := q.terms[i]
		pj := p.terms[j]
		switch {
		case qi.ID() < pj.ID():
			if len(r.terms) == 0 || r.terms[len(r.terms)-1].UID != qi.UID {
				r.terms = append(r.terms, qi)
			}
			i++
		case pj.ID() < qi.ID():
			if len(r.terms) == 0 || r.terms[len(r.terms)-1].UID != pj.UID {
				r.terms = append(r.terms, pj)
			}
			j++
		default:
			if len(r.terms) == 0 || r.terms[len(r.terms)-1].UID != qi.UID {
				r.terms = append(r.terms, qi)
			}
			i++
			j++
		}
	}
	r.terms = append(r.terms, q.terms[i:]...)
	r.terms = append(r.terms, p.terms[j:]...)
	return r
}

// Not returns a query that holds q less p.
func (q Query) Not(p Query) Query {
	if q.g != p.g {
		panic("rdf: binary query operation parameters from distinct graphs")
	}
	sortByID(q.terms)
	sortByID(p.terms)
	r := Query{g: q.g}
	var i, j int
	for i < len(q.terms) && j < len(p.terms) {
		qi := q.terms[i]
		pj := p.terms[j]
		switch {
		case qi.ID() < pj.ID():
			r.terms = append(r.terms, qi)
			i++
		case pj.ID() < qi.ID():
			j++
		default:
			i++
		}
	}
	if len(r.terms) < len(q.terms) {
		r.terms = append(r.terms, q.terms[i:len(q.terms)+min(0, i-len(r.terms))]...)
	}
	return r
}

// Repeat repeatedly calls fn on q until the set of results is empty or
// ok is false, and then returns the result. If the last non-empty result
// is wanted, fn should return its input and false when the partial
// traversal returns an empty result.
//
//	result := start.Repeat(func(q rdf.Query) (rdf.Query, bool) {
//		r := q.Out(condition)
//		if r.Len() == 0 {
//			return q, false
//		}
//		return r, true
//	}).Result()
func (q Query) Repeat(fn func(Query) (q Query, ok bool)) Query {
	for {
		var ok bool
		q, ok = fn(q)
		if !ok || len(q.terms) == 0 {
			return q
		}
	}
}

// Unique returns a copy of the receiver that contains only one instance
// of each term.
func (q Query) Unique() Query {
	sortByID(q.terms)
	r := Query{g: q.g}
	for i, t := range q.terms {
		if i == 0 || t.UID != q.terms[i-1].UID {
			r.terms = append(r.terms, t)
		}
	}
	return r
}

// Len returns the number of terms held by the query.
func (q Query) Len() int {
	return len(q.terms)
}

// Result returns the terms held by the query.
func (q Query) Result() []Term {
	return q.terms
}

func sortByID(terms []Term) {
	sort.Slice(terms, func(i, j int) bool { return terms[i].ID() < terms[j].ID() })
}
