// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf

import (
	"errors"
	"sort"
)

// Throughout, the comments refer to doi:10.1145/3068333 which should be
// understood as a synonym for http://aidanhogan.com/docs/rdf-canonicalisation.pdf
// although there are differences between the two, see http://aidanhogan.com/#errataH17.
// Where there are differences, the document at http://aidanhogan.com/ is the
// canonical truth. The DOI reference is referred to for persistence.

// Lean returns an RDF core of g that entails g. If g contains any non-zero
// labels, Lean will return a non-nil error and a core of g assuming no graph
// labels exist.
//
// See http://aidanhogan.com/docs/rdf-canonicalisation.pdf for details of
// the algorithm.
func Lean(g []*Statement) ([]*Statement, error) {
	// BUG(kortschak): Graph leaning does not take into account graph label terms
	// since the formal semantics for a multiple graph data model have not been
	// defined. See https://www.w3.org/TR/rdf11-datasets/#declaring.

	var (
		hasBlanks bool
		err       error
	)
	for _, s := range g {
		if isBlank(s.Subject.Value) || isBlank(s.Object.Value) {
			hasBlanks = true
			if err != nil {
				break
			}
		}
		if s.Label.Value != "" && err == nil {
			err = errors.New("rdf: data-set contains graph names")
			if hasBlanks {
				break
			}
		}
	}
	if hasBlanks {
		g = lean(&dfs{}, g)
	}
	return g, err
}

// removeRedundantBnodes removes blank nodes whose edges are a subset of
// another term in the RDF graph.
//
// This is algorithm 4 in doi:10.1145/3068333.
func removeRedundantBnodes(g []*Statement) []*Statement {
	g = append(g[:0:0], g...)
	for {
		edges := make(map[string]map[triple]bool)
		for _, s := range g {
			for i, t := range []string{
				s.Subject.Value,
				s.Object.Value,
			} {
				e, ok := edges[t]
				if !ok {
					e = make(map[triple]bool)
					edges[t] = e
				}
				switch i {
				case 0:
					e[triple{s.Predicate.Value, s.Object.Value, "+"}] = true
				case 1:
					e[triple{s.Predicate.Value, s.Subject.Value, "-"}] = true
				}
			}
		}

		seen := make(map[string]bool)
		bNodes := make(map[string]bool)
		terms := make(map[string]bool)
		for _, s := range g {
			for _, t := range []string{
				s.Subject.Value,
				s.Predicate.Value,
				s.Object.Value,
			} {
				terms[t] = true
				if isBlank(t) {
					bNodes[t] = true
				} else {
					seen[t] = true
				}
			}
		}

		redundant := make(map[string]bool)
		for x := range bNodes {
			for xp := range terms {
				if isProperSubset(edges[x], edges[xp]) || (seen[xp] && isEqualEdges(edges[x], edges[xp])) {
					redundant[x] = true
					break
				}
			}
			seen[x] = true
		}

		n := len(g)
		for i := 0; i < len(g); {
			if !redundant[g[i].Subject.Value] && !redundant[g[i].Object.Value] {
				i++
				continue
			}
			g[i], g = g[len(g)-1], g[:len(g)-1]
		}
		if n == len(g) {
			return g
		}
	}
}

type triple [3]string

func isProperSubset(a, b map[triple]bool) bool {
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return len(a) < len(b)
}

func isEqualEdges(a, b map[triple]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

// findCandidates finds candidates for blank nodes and blank nodes that are fixed.
//
// This is algorithm 5 in doi:10.1145/3068333.
func findCandidates(g []*Statement) ([]*Statement, map[string]bool, map[string]map[string]bool, bool) {
	g = removeRedundantBnodes(g)

	edges := make(map[triple]bool)
	f := make(map[string]bool)
	for _, s := range g {
		sub := s.Subject.Value
		prd := s.Predicate.Value
		obj := s.Object.Value

		edges[triple{sub, prd, obj}] = true
		edges[triple{sub, prd, "*"}] = true
		edges[triple{"*", prd, obj}] = true
		switch {
		case isBlank(sub) && isBlank(obj):
			f[sub] = false
			f[obj] = false
		case isBlank(sub):
			if _, ok := f[sub]; !ok {
				f[sub] = true
			}
		case isBlank(obj):
			if _, ok := f[obj]; !ok {
				f[obj] = true
			}
		}
	}
	for k, v := range f {
		if !v {
			delete(f, k)
		}
	}
	if len(f) == 0 {
		f = nil
	}

	cands := make(map[string]map[string]bool)
	bnodes := make(map[string]bool)
	for _, s := range g {
		for _, b := range []string{
			s.Subject.Value,
			s.Object.Value,
		} {
			if !isBlank(b) {
				continue
			}
			bnodes[b] = true
			if f[b] {
				cands[b] = map[string]bool{b: true}
			} else {
				terms := make(map[string]bool)
				for _, s := range g {
					for _, t := range []string{
						s.Subject.Value,
						s.Predicate.Value,
						s.Object.Value,
					} {
						terms[t] = true
					}
				}
				cands[b] = terms
			}
		}
	}
	if isEqualTerms(f, bnodes) {
		return g, f, cands, true
	}

	for {
		bb := make(map[string]bool)
		for b := range bnodes {
			if !f[b] {
				bb[b] = true
			}
		}
		for b := range bb {
			for x := range cands[b] {
				if x == b {
					continue
				}
				for _, s := range g {
					if s.Subject.Value != b {
						continue
					}
					prd := s.Predicate.Value
					obj := s.Object.Value
					if (inILF(obj, f) && !edges[triple{x, prd, obj}]) || (bb[obj] && !edges[triple{x, prd, "*"}]) {
						delete(cands[b], x)
						break
					}
				}
				if !cands[b][x] {
					continue
				}
				for _, s := range g {
					if s.Object.Value != b {
						continue
					}
					sub := s.Subject.Value
					prd := s.Predicate.Value
					if (inIF(sub, f) && !edges[triple{sub, prd, x}]) || (bb[sub] && !edges[triple{"*", prd, x}]) {
						delete(cands[b], x)
						break
					}
				}
			}
		}

		fp := f
		f = make(map[string]bool)
		for b := range fp {
			f[b] = true
		}
		for b := range bb { // Mark newly fixed blank nodes.
			if len(cands[b]) == 1 && cands[b][b] {
				f[b] = true
			}
		}
		allFixed := isEqualTerms(f, bnodes)
		if isEqualTerms(fp, f) || allFixed {
			if len(f) == 0 {
				f = nil
			}
			return g, f, cands, allFixed
		}
	}
}

// inILF returns whether t is in IL or F.
func inILF(t string, f map[string]bool) bool {
	return isIRI(t) || isLiteral(t) || f[t]
}

// inIF returns whether t is in I or F.
func inIF(t string, f map[string]bool) bool {
	return isIRI(t) || f[t]
}

// dfs is a depth-first search strategy.
type dfs struct{}

// lean returns a core of the RDF graph g using the given strategy.
//
// This is lines 1-9 of algorithm 6 in doi:10.1145/3068333.
func lean(strategy *dfs, g []*Statement) []*Statement {
	foundBnode := false
search:
	for _, s := range g {
		for _, t := range []string{
			s.Subject.Value,
			s.Object.Value,
		} {
			if isBlank(t) {
				foundBnode = true
				break search
			}
		}
	}
	if !foundBnode {
		return g
	}
	g, fixed, cands, allFixed := findCandidates(g)
	if allFixed {
		return g
	}
	for _, s := range g {
		if isBlank(s.Subject.Value) && isBlank(s.Object.Value) {
			mu := make(map[string]string, len(fixed))
			for b := range fixed {
				mu[b] = b
			}
			mu = findCoreEndomorphism(strategy, g, cands, mu)
			return applyMu(g, mu)
		}
	}
	return g
}

// findCoreEndomorphism returns a core solution using the given strategy.
//
// This is lines 10-14 of algorithm 6 in doi:10.1145/3068333.
func findCoreEndomorphism(strategy *dfs, g []*Statement, cands map[string]map[string]bool, mu map[string]string) map[string]string {
	var q []*Statement
	preds := make(map[string]int)
	seen := make(map[triple]bool)
	for _, s := range g {
		preds[s.Predicate.Value]++
		if isBlank(s.Subject.Value) && isBlank(s.Object.Value) {
			if seen[triple{s.Subject.Value, s.Predicate.Value, s.Object.Value}] {
				continue
			}
			seen[triple{s.Subject.Value, s.Predicate.Value, s.Object.Value}] = true
			q = append(q, s)
		}
	}
	sort.Slice(q, func(i, j int) bool {
		return selectivity(q[i], cands, preds) < selectivity(q[j], cands, preds)
	})
	return strategy.evaluate(g, q, cands, mu)
}

// selectivity returns the selectivity heuristic score for s. Lower scores
// are more selective.
func selectivity(s *Statement, cands map[string]map[string]bool, preds map[string]int) int {
	return min(len(cands[s.Subject.Value])*len(cands[s.Object.Value]), preds[s.Predicate.Value])
}

// evaluate returns an endomorphism using a DFS strategy.
//
// This is lines 25-32 of algorithm 6 in doi:10.1145/3068333.
func (st *dfs) evaluate(g, q []*Statement, cands map[string]map[string]bool, mu map[string]string) map[string]string {
	mu = st.search(g, q, cands, mu)
	for len(mu) != len(codom(mu)) {
		mupp := fixedFrom(cands)
		mup := findCoreEndomorphism(st, applyMu(g, mu), cands, mupp)
		if isAutomorphism(mup) {
			return mu
		}
		for b, x := range mu {
			if _, ok := mup[b]; !ok {
				mup[b] = x
			}
		}
		mu = mup
	}
	return mu
}

func fixedFrom(cands map[string]map[string]bool) map[string]string {
	fixed := make(map[string]string)
	for b, m := range cands {
		if len(m) == 1 && m[b] {
			fixed[b] = b
		}
	}
	return fixed
}

// applyMu applies mu to g returning the result.
func applyMu(g []*Statement, mu map[string]string) []*Statement {
	back := make([]Statement, 0, len(g))
	dst := make([]*Statement, 0, len(g))
	seen := make(map[Statement]bool)
	for _, s := range g {
		n := Statement{
			Subject:   Term{Value: translate(s.Subject.Value, mu)},
			Predicate: Term{Value: s.Predicate.Value},
			Object:    Term{Value: translate(s.Object.Value, mu)},
			Label:     Term{Value: s.Label.Value},
		}
		if seen[n] {
			continue
		}
		seen[n] = true
		back = append(back, n)
		dst = append(dst, &back[len(back)-1])
	}
	return dst
}

// search returns a minimum endomorphism using a DFS strategy.
//
// This is lines 33-46 of algorithm 6 in doi:10.1145/3068333.
func (st *dfs) search(g, q []*Statement, cands map[string]map[string]bool, mu map[string]string) map[string]string {
	qMin := q[0]
	m := st.join(qMin, g, cands, mu)
	if len(m) == 0 {
		// Early exit if no mapping found.
		return nil
	}
	sortByCodom(m)
	mMin := m[0]
	qp := q[1:]
	if len(qp) != 0 {
		for len(m) != 0 {
			mMin = m[0]
			mup := st.search(g, qp, cands, mMin)
			if !isAutomorphism(mup) {
				return mup
			}
			m = m[1:]
		}
	}
	return mMin
}

// isAutomorphism returns whether mu is an automorphism, this is equivalent to
// dom(mu) == codom(mu).
func isAutomorphism(mu map[string]string) bool {
	return isEqualTerms(dom(mu), codom(mu))
}

// dom returns the domain of mu.
func dom(mu map[string]string) map[string]bool {
	d := make(map[string]bool, len(mu))
	for v := range mu {
		d[v] = true
	}
	return d
}

// codom returns the codomain of mu.
func codom(mu map[string]string) map[string]bool {
	cd := make(map[string]bool, len(mu))
	for _, v := range mu {
		cd[v] = true
	}
	return cd
}

// isEqualTerms returns whether a and b are identical.
func isEqualTerms(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

// sortByCodom performs a sort of maps ordered by fewest blank nodes in
// codomain, then fewest self mappings.
func sortByCodom(maps []map[string]string) {
	m := orderedByCodom{
		maps:  maps,
		attrs: make([]attrs, len(maps)),
	}
	for i, mu := range maps {
		m.attrs[i].blanks = make(map[string]bool)
		for x, y := range mu {
			if isBlank(y) {
				m.attrs[i].blanks[y] = true
			}
			if x == y {
				m.attrs[i].selfs++
			}
		}
	}
	sort.Sort(m)
}

type orderedByCodom struct {
	maps  []map[string]string
	attrs []attrs
}

type attrs struct {
	blanks map[string]bool
	selfs  int
}

func (m orderedByCodom) Len() int { return len(m.maps) }
func (m orderedByCodom) Less(i, j int) bool {
	attrI := m.attrs[i]
	attrJ := m.attrs[j]
	switch {
	case len(attrI.blanks) < len(attrJ.blanks):
		return true
	case len(attrI.blanks) > len(attrJ.blanks):
		return false
	default:
		return attrI.selfs < attrJ.selfs
	}
}
func (m orderedByCodom) Swap(i, j int) {
	m.maps[i], m.maps[j] = m.maps[j], m.maps[i]
	m.attrs[i], m.attrs[j] = m.attrs[j], m.attrs[i]
}

// join evaluates the given pattern, q, joining with solutions in m.
// This takes only a single mapping and so only works for the DFS strategy.
//
// This is lines 47-51 of algorithm 6 in doi:10.1145/3068333.
func (st *dfs) join(q *Statement, g []*Statement, cands map[string]map[string]bool, m map[string]string) []map[string]string {
	var mp []map[string]string
	isLoop := q.Subject.Value == q.Object.Value
	for _, s := range g {
		// Line 45: M_q ← {µ | µ(q) ∈ G}
		//  | µ(q) ∈ G
		//
		//    µ(q) ∈ G ↔ (µ(q_s),q_p,µ(q_o)) ∈ G
		if q.Predicate.Value != s.Predicate.Value {
			continue
		}
		//    q_s = q_o ↔ µ(q_s) =_µ(q_o)
		if isLoop && s.Subject.Value != s.Object.Value {
			continue
		}

		// Line 46: M_q' ← {µ ∈ M_q | for all b ∈ bnodes({q}), µ(b) ∈ cands[b]}
		//  | for all b ∈ bnodes({q}), µ(b) ∈ cands[b]
		if !cands[q.Subject.Value][s.Subject.Value] || !cands[q.Object.Value][s.Object.Value] {
			continue
		}

		// Line 47: M' ← M_q' ⋈ M
		// M₁ ⋈ M₂ = {μ₁ ∪ μ₂ | μ₁ ∈ M₁, μ₂ ∈ M₂ and μ₁, μ₂ are compatible mappings}
		//  | μ₁ ∈ M₁, μ₂ ∈ M₂ and μ₁, μ₂ are compatible mappings
		if mq, ok := m[q.Subject.Value]; ok && mq != s.Subject.Value {
			continue
		}
		if !isLoop {
			if mq, ok := m[q.Object.Value]; ok && mq != s.Object.Value {
				continue
			}
		}
		// Line 47: μ₁ ∪ μ₂
		var mu map[string]string
		if isLoop {
			mu = map[string]string{
				q.Subject.Value: s.Subject.Value,
			}
		} else {
			mu = map[string]string{
				q.Subject.Value: s.Subject.Value,
				q.Object.Value:  s.Object.Value,
			}
		}
		for b, mb := range m {
			mu[b] = mb
		}
		mp = append(mp, mu)
	}
	return mp
}
