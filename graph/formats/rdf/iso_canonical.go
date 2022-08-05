// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf

import (
	"bytes"
	"errors"
	"fmt"
	"hash"
	"sort"
)

// See "Canonical Forms for Isomorphic and Equivalent RDF Graphs: Algorithms
// for Leaning and Labelling Blank Nodes" by Aiden Hogan for description of
// the algorithm, https://doi.org/10.1145/3068333 and available free from
// the author's web page http://aidanhogan.com/docs/rdf-canonicalisation.pdf.
//
// Aspects of implementation from discussion in v1.0 of the readme of the PoC
// at https://doi.org/10.5281/zenodo.3154322

// Isomorphic returns whether the RDF graph datasets a and b are isomorphic,
// where there is a bijective mapping between blank nodes in a and b using
// the given hash function. If decomp is true, the graphs are decomposed
// before canonicalization.
func Isomorphic(a, b []*Statement, decomp bool, h hash.Hash) bool {
	if len(a) != len(b) {
		return false
	}

	zero := make([]byte, h.Size())
	ah, _ := IsoCanonicalHashes(a, decomp, true, h, zero)
	bh, _ := IsoCanonicalHashes(b, decomp, true, h, zero)
	if len(ah) != len(bh) {
		return false
	}

	work := make([][]byte, 2*len(ah))
	lexicalHashes(work[:len(ah)], ah)
	lexicalHashes(work[len(ah):], bh)
	for i := range work[:len(ah)] {
		if !bytes.Equal(work[i], work[i+len(ah)]) {
			return false
		}
	}
	return true
}

func lexicalHashes(dst [][]byte, hashes map[string][]byte) {
	i := 0
	for _, s := range hashes {
		dst[i] = s
		i++
	}
	sort.Sort(lexical(dst))
}

// IsoCanonicalHashes returns a mapping between the nodes of the RDF graph
// dataset described by the given statements using the provided hash
// function. If decomp is true, the graphs are decomposed before hashing.
// If dist is true the input graph is decomposed into identical splits, the
// entire graph will be hashed to distinguish nodes. If decomp is false,
// dist has no effect.
// Blank node hashes are initially set to the value of zero. Hash values
// are provided for literal and IRI nodes as well as for blank node. The
// hash input for literal nodes includes the quotes and the input for IRI
// nodes first removes the angle quotes around the IRI, although these are
// included in the map keys.
//
// Note that hashes returned by IsoCanonicalHashes with decomp=true are not
// comparable with hashes returned by IsoCanonicalHashes with decomp=false.
//
// See http://aidanhogan.com/docs/rdf-canonicalisation.pdf for details of
// the hashing algorithm.
func IsoCanonicalHashes(statements []*Statement, decomp, dist bool, h hash.Hash, zero []byte) (hashes map[string][]byte, terms map[string]map[string]bool) {
	if len(statements) == 0 {
		return nil, nil
	}

	if debug {
		debug.log(0, "Statements:")
		for _, s := range statements {
			debug.log(0, s)
		}
		debug.log(0)
	}

	hash, parts, ok := hashBNodesPerSplit(statements, decomp, h, zero)

	if debug {
		debug.log(0, "Blanks:")
		if len(hash.blanks) != 0 {
			for _, b := range hash.blanks {
				debug.log(0, b)
			}
		} else {
			debug.log(0, "none")
		}
		debug.log(0)

		debug.log(0, "Parts:")
		debug.logParts(0, parts)

		debug.logf(0, "Hashes from hashBNodesPerSplit (splitting=%t):\n", decomp)
		debug.logHashes(0, hash.hashOf, h.Size())
	}

	if ok {
		return hash.hashOf, hash.termsFor
	}

	// TODO: remove the triviality exception in distinguish and return
	// the original hashes if this result is nil. Make the triviality
	// exception optional.
	hashes = distinguish(statements, dist, h, zero, hash, parts, nil, 0)

	if hashes == nil {
		// distinguish was given trivial parts and
		// we did not ask it to try to merge them.
		return hash.hashOf, hash.termsFor
	}

	if debug {
		debug.log(0, "Final resolved Hashes:")
		debug.logHashes(0, hashes, h.Size())
	}

	terms = make(map[string]map[string]bool, len(hashes))
	for k, h := range hashes {
		terms[string(h)] = map[string]bool{k: true}
	}

	return hashes, terms
}

// C14n performs a relabeling of the statements in src based on the terms
// obtained from IsoCanonicalHashes, placing the results in dst and returning
// them. The relabeling scheme is the same as for the Universal RDF Dataset
// Normalization Algorithm, blank terms are ordered lexically by their hash
// value and then given a blank label with the prefix "_:c14n" and an
// identifier counter corresponding to the label's sort rank.
//
// If dst is nil, it is allocated, otherwise the length of dst must match the
// length of src.
func C14n(dst, src []*Statement, terms map[string]map[string]bool) ([]*Statement, error) {
	if dst == nil {
		dst = make([]*Statement, len(src))
	}

	if len(dst) != len(src) {
		return dst, errors.New("rdf: slice length mismatch")
	}

	need := make(map[string]bool)
	for _, s := range src {
		for _, t := range []string{
			s.Subject.Value,
			s.Object.Value,
			s.Label.Value,
		} {
			if !isBlank(t) {
				continue
			}
			need[t] = true
		}
	}

	blanks := make([]string, len(need))
	i := 0
	for h, m := range terms {
		var ok bool
		for t := range m {
			if isBlank(t) {
				ok = true
				break
			}
		}
		if !ok {
			continue
		}
		if i == len(blanks) {
			return dst, errors.New("rdf: too many blanks in terms")
		}
		blanks[i] = h
		i++
	}
	sort.Strings(blanks)

	c14n := make(map[string]string)
	for i, b := range blanks {
		if len(terms[b]) == 0 {
			return nil, fmt.Errorf("rdf: no term for blank with hash %x", b)
		}
		for t := range terms[b] {
			if !isBlank(t) {
				continue
			}
			if _, exists := c14n[t]; exists {
				continue
			}
			delete(need, t)
			c14n[t] = fmt.Sprintf("_:c14n%d", i)
		}
	}

	if len(need) != 0 {
		return dst, fmt.Errorf("rdf: missing term hashes for %d terms", len(need))
	}

	for i, s := range src {
		if dst[i] == nil {
			dst[i] = &Statement{}
		}
		n := dst[i]
		n.Subject = Term{Value: translate(s.Subject.Value, c14n)}
		n.Predicate = s.Predicate
		n.Object = Term{Value: translate(s.Object.Value, c14n)}
		n.Label = Term{Value: translate(s.Label.Value, c14n)}
	}
	sort.Sort(c14nStatements(dst))

	return dst, nil
}

func translate(term string, mapping map[string]string) string {
	if term, ok := mapping[term]; ok {
		return term
	}
	return term
}

type c14nStatements []*Statement

func (s c14nStatements) Len() int { return len(s) }
func (s c14nStatements) Less(i, j int) bool {
	si := s[i]
	sj := s[j]
	switch {
	case si.Subject.Value < sj.Subject.Value:
		return true
	case si.Subject.Value > sj.Subject.Value:
		return false
	}
	switch { // Always IRI.
	case si.Predicate.Value < sj.Predicate.Value:
		return true
	case si.Predicate.Value > sj.Predicate.Value:
		return false
	}
	switch {
	case si.Object.Value < sj.Object.Value:
		return true
	case si.Object.Value > sj.Object.Value:
		return false
	}
	return si.Label.Value < sj.Label.Value
}
func (s c14nStatements) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// hashBNodes returns the hashed blank nodes of the graph described by statements
// using the provided hash function. Hashes are initialised with zero.
//
// This is algorithm 1 in doi:10.1145/3068333.
func hashBNodes(statements []*Statement, h hash.Hash, zero []byte, hash0 map[string][]byte) (hash *table, disjoint bool) {
	curr := newTable()
	for _, s := range statements {
		for i, t := range []string{
			s.Subject.Value,
			s.Predicate.Value,
			s.Object.Value,
			s.Label.Value,
		} {
			switch {
			case i == 3 && t == "":
				continue
			case isBlank(t):
				if hash0 == nil {
					curr.set(t, zero)
				} else {
					curr.set(t, hash0[t])
				}
			case isIRI(t):
				h.Reset()
				h.Write([]byte(t[1 : len(t)-1]))
				curr.set(t, h.Sum(nil))
			default:
				h.Reset()
				h.Write([]byte(t))
				curr.set(t, h.Sum(nil))
			}
		}
	}

	bag := newHashBag(h, curr)
	last := curr.clone()
	for {
		curr, last = last, curr
		for _, s := range statements {
			if isBlank(s.Subject.Value) {
				var lab []byte
				if s.Label.Value != "" {
					lab = last.hashOf[s.Label.Value]
				}
				c := hashTuple(h, last.hashOf[s.Object.Value], last.hashOf[s.Predicate.Value], lab, []byte{'+'})
				bag.add(s.Subject.Value, c)
			}

			if isBlank(s.Object.Value) {
				var lab []byte
				if s.Label.Value != "" {
					lab = last.hashOf[s.Label.Value]
				}
				c := hashTuple(h, last.hashOf[s.Subject.Value], last.hashOf[s.Predicate.Value], lab, []byte{'-'})
				bag.add(s.Object.Value, c)
			}

			// This and the lab value above implement the label hashing
			// required for RDF dataset hashing as described in
			// https://doi.org/10.5281/zenodo.3154322 v1.0
			// Readme.md#adaptation-of-the-algorithms-to-handle-datasets.
			if isBlank(s.Label.Value) {
				c := hashTuple(h, last.hashOf[s.Subject.Value], last.hashOf[s.Predicate.Value], last.hashOf[s.Object.Value], []byte{'.'})
				bag.add(s.Label.Value, c)
			}
		}

		for t := range bag.hashesFor {
			curr.set(t, bag.sum(t))
		}

		disjoint = curr.allUnique()
		if disjoint || !curr.changedFrom(last) {
			return curr, disjoint
		}
	}
}

// table is a collision aware hash collection for RDF terms.
type table struct {
	// hashOf holds the hash for each term.
	hashOf map[string][]byte
	// termsFor holds the set of nodes in
	// the second key for terms that share
	// the hash in the first key.
	termsFor map[string]map[string]bool

	// isBlank and blanks are the set of blank
	// nodes.
	// isBlank is nil for cloned tables.
	isBlank map[string]bool
	// blanks is nil for tables created
	// with newTable.
	blanks []string
}

// newTable returns a new hash table.
func newTable() *table {
	return &table{
		hashOf:   make(map[string][]byte),
		termsFor: make(map[string]map[string]bool),
		isBlank:  make(map[string]bool),
	}
}

// wasCloned returns whether t is a parent or child of a cloning operation.
func (t *table) wasCloned() bool { return t.isBlank == nil }

// isNew returns whether t is a new table.
func (t *table) isNew() bool { return t.blanks == nil }

// clone returns a clone of the receiver.
func (t *table) clone() *table {
	new := &table{
		hashOf:   make(map[string][]byte),
		termsFor: make(map[string]map[string]bool),
	}
	for term, hash := range t.hashOf {
		new.hashOf[term] = hash
	}
	for hash, coll := range t.termsFor {
		if len(coll) == 0 {
			continue
		}
		terms := make(map[string]bool)
		for term := range coll {
			terms[term] = true
		}
		new.termsFor[hash] = terms
	}
	if t.isNew() {
		t.blanks = make([]string, len(t.isBlank))
		i := 0
		for n := range t.isBlank {
			t.blanks[i] = n
			i++
		}
		t.isBlank = nil
	}
	new.blanks = t.blanks
	return new
}

// TODO(kortschak): Make hash table in table.hashOf reuse the []byte on update.
// This is not trivial since we need to check for changes, so we can't just get
// the current hash buffer and write into it. So if this is done we probably
// a pair of buffers, a current and a waiting.

// set sets the hash of the term, removing any previously set hash.
func (t *table) set(term string, hash []byte) {
	prev := t.hashOf[term]
	if bytes.Equal(prev, hash) {
		return
	}
	t.hashOf[term] = hash

	// Delete any existing hashes for this term.
	switch terms := t.termsFor[string(prev)]; {
	case len(terms) == 1:
		delete(t.termsFor, string(prev))
	case len(terms) > 1:
		delete(terms, term)
	}

	terms, ok := t.termsFor[string(hash)]
	if ok {
		terms[term] = true
	} else {
		t.termsFor[string(hash)] = map[string]bool{term: true}
	}

	if !t.wasCloned() && isBlank(term) {
		// We are in the original table, so note
		// any blank node label that we see.
		t.isBlank[term] = true
	}
}

// allUnique returns whether every term has an unique hash. allUnique
// can only be called on a table that was returned by clone.
func (t *table) allUnique() bool {
	if t.isNew() {
		panic("checked hash bag from uncloned table")
	}
	for _, term := range t.blanks {
		if len(t.termsFor[string(t.hashOf[term])]) > 1 {
			return false
		}
	}
	return true
}

// changedFrom returns whether the receiver has been updated from last.
// changedFrom can only be called on a table that was returned by clone.
func (t *table) changedFrom(last *table) bool {
	if t.isNew() {
		panic("checked hash bag from uncloned table")
	}
	for i, x := range t.blanks {
		for _, y := range t.blanks[i+1:] {
			if bytes.Equal(t.hashOf[x], t.hashOf[y]) != bytes.Equal(last.hashOf[x], last.hashOf[y]) {
				return true
			}
		}
	}
	return false
}

// hashBag implements a commutative and associative hash.
// See notes in https://doi.org/10.5281/zenodo.3154322 v1.0
// Readme.md#what-is-the-precise-specification-of-hashbag.
type hashBag struct {
	hash      hash.Hash
	hashesFor map[string][][]byte
}

// newHashBag returns a new hashBag using the provided hash function for
// the given hash table. newHashBag can only take a table parameter that
// was returned by newTable.
func newHashBag(h hash.Hash, t *table) hashBag {
	if t.wasCloned() {
		panic("made hash bag from cloned table")
	}
	b := hashBag{hash: h, hashesFor: make(map[string][][]byte, len(t.isBlank))}
	for n := range t.isBlank {
		b.hashesFor[n] = [][]byte{t.hashOf[n]}
	}
	return b
}

// add adds the hash to the hash bag for the term.
func (b hashBag) add(term string, hash []byte) {
	b.hashesFor[term] = append(b.hashesFor[term], hash)
}

// sum calculates the hash sum for the given term, updates the hash bag
// state and returns the hash.
func (b hashBag) sum(term string) []byte {
	p := b.hashesFor[term]
	sort.Sort(lexical(p))
	h := hashTuple(b.hash, p...)
	b.hashesFor[term] = b.hashesFor[term][:1]
	b.hashesFor[term][0] = h
	return h
}

// lexical implements lexical sorting of [][]byte.
type lexical [][]byte

func (b lexical) Len() int           { return len(b) }
func (b lexical) Less(i, j int) bool { return string(b[i]) < string(b[j]) }
func (b lexical) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

// hashTuple returns the h hash of the concatenation of t.
func hashTuple(h hash.Hash, t ...[]byte) []byte {
	h.Reset()
	for _, b := range t {
		h.Write(b)
	}
	return h.Sum(nil)
}

// hashBNodesPerSplit returns the independently hashed blank nodes of the
// graph described by statements using the provided hash function. Hashes
// are initialised with zero.
//
// This is algorithm 2 in doi:10.1145/3068333.
func hashBNodesPerSplit(statements []*Statement, decomp bool, h hash.Hash, zero []byte) (hash *table, parts byLengthHash, disjoint bool) {
	if !decomp {
		hash, ok := hashBNodes(statements, h, zero, nil)
		parts = appendOrdered(byLengthHash{}, hash.termsFor)
		sort.Sort(parts)
		return hash, parts, ok
	}

	splits := split(statements)

	// Avoid recombination work if there is only one split.
	if len(splits) == 1 {
		hash, ok := hashBNodes(statements, h, zero, nil)
		parts = appendOrdered(byLengthHash{}, hash.termsFor)
		sort.Sort(parts)
		return hash, parts, ok
	}

	hash = &table{hashOf: make(map[string][]byte)}
	disjoint = true
	for _, g := range splits {
		part, ok := hashBNodes(g, h, zero, nil)
		// Each split is guaranteed to be disjoint in its
		// set of blank nodes, so we can just append to our
		// collection of blanks.
		hash.blanks = append(hash.blanks, part.blanks...)
		if !ok {
			// Allow a short-circuit of the allUnique check.
			disjoint = false
		}
		for k, v := range part.hashOf {
			hash.hashOf[k] = v
		}
		parts = appendOrdered(parts, part.termsFor)
	}
	sort.Sort(parts)
	return hash, parts, disjoint && allUnique(hash.hashOf)
}

// appendOrdered adds parts (labels stored in the second key) for each
// hash (stored in the first key) to parts.
func appendOrdered(parts byLengthHash, partSets map[string]map[string]bool) byLengthHash {
	for h, s := range partSets {
		var p []string
		for e := range s {
			if isBlank(e) {
				p = append(p, e)
			}
		}
		if p != nil {
			parts.nodes = append(parts.nodes, p)
			parts.hashes = append(parts.hashes, h)
		}
	}
	return parts
}

// byLengthHash implements ascending length sort of a set of blank RDF
// term partitions with ties broken by lexical ordering of the partitions'
// hashes.
type byLengthHash struct {
	// nodes holds the blank nodes of a part.
	nodes [][]string
	// hashes holds the hashes corresponding
	// to the nodes in the nodes field, using
	// the same index.
	hashes []string
}

func (s byLengthHash) Len() int { return len(s.nodes) }
func (s byLengthHash) Less(i, j int) bool {
	switch {
	case len(s.nodes[i]) < len(s.nodes[j]):
		return true
	case len(s.nodes[i]) > len(s.nodes[j]):
		return false
	}
	return s.hashes[i] < s.hashes[j]
}
func (s byLengthHash) Swap(i, j int) {
	s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i]
	s.hashes[i], s.hashes[j] = s.hashes[j], s.hashes[i]
}

// allUnique returns whether the []byte hash values in hashes are all unique.
func allUnique(hashes map[string][]byte) bool {
	set := make(map[string]bool)
	for _, h := range hashes {
		if set[string(h)] {
			return false
		}
		set[string(h)] = true
	}
	return true
}

// split returns the statements forming connected components in the graph
// described by statements.
//
// This is split in algorithm 2 in doi:10.1145/3068333.
func split(statements []*Statement) [][]*Statement {
	ds := make(djSet)
	for _, s := range statements {
		ds.add(s.Subject.Value)
		ds.add(s.Object.Value)
		if isBlank(s.Subject.Value) && isBlank(s.Object.Value) {
			ds.union(ds.find(s.Subject.Value), ds.find(s.Object.Value))
		}
	}

	var (
		splits [][]*Statement
		ground []*Statement
	)
	idxOf := make(map[*dsNode]int)
	for _, s := range statements {
		var t string
		switch {
		case isBlank(s.Subject.Value):
			t = s.Subject.Value
		case isBlank(s.Object.Value):
			t = s.Object.Value
		default:
			ground = append(ground, s)
			continue
		}
		r := ds.find(t)
		if r == nil {
			panic(fmt.Sprintf("term not found: %q", t))
		}
		i, ok := idxOf[r]
		if !ok {
			i = len(splits)
			idxOf[r] = i
			splits = append(splits, []*Statement{s})
		} else {
			splits[i] = append(splits[i], s)
		}
	}
	if ground != nil {
		splits = append(splits, ground)
	}

	if debug {
		debug.log(0, "Splits:")
		for i, s := range splits {
			for j, t := range s {
				if j == 0 {
					debug.logf(0, "%d.\t%s\n", i+1, t)
				} else {
					debug.logf(0, "\t%s\n", t)
				}
			}
			debug.log(0)
		}
	}

	return splits
}

// distinguish returns G⊥: smallest hash-labelled graph found thus far.
// The graph is returned as a node to hash lookup.
//
// This is part of algorithm 3 in doi:10.1145/3068333.
//
// The correspondence between the parameters for the function in the paper
// with the implementation here is as follows:
//   - G = statements
//   - hash = hash
//   - P = parts (already sorted by hashBNodesPerSplit)
//   - G⊥ = lowest
//   - B = hash.blanks
//
// The additional parameter dist specifies that distinguish should treat
// coequal trivial parts as a coarse of intermediate part and distinguish
// the nodes in that merged part.
func distinguish(statements []*Statement, dist bool, h hash.Hash, zero []byte, hash *table, parts byLengthHash, lowest map[string][]byte, depth int) map[string][]byte {
	if debug {
		debug.log(depth, "Running Distinguish")
	}

	var small []string
	var k int
	for k, small = range parts.nodes {
		if len(small) > 1 {
			break
		}
	}
	if len(small) < 2 {
		if lowest != nil || !dist {
			if debug {
				debug.log(depth, "Return lowest (no non-trivial parts):")
				debug.logHashes(depth, lowest, h.Size())
			}

			return lowest
		}

		// We have been given a set of fine parts,
		// but to reach here they must have been
		// non-uniquely labeled, so treat them
		// as a single coarse part.
		k, small = 0, parts.nodes[0]
	}

	if debug {
		debug.logf(depth, "Part: %v %x\n\n", small, parts.hashes[k])
		debug.log(depth, "Orig hash:")
		debug.logHashes(depth, hash.hashOf, h.Size())
	}

	smallHash := hash.hashOf[small[0]]
	for _, p := range parts.nodes[k:] {
		if !bytes.Equal(smallHash, hash.hashOf[p[0]]) {

			if debug {
				debug.logf(depth, "End of co-equal hashes: %x != %x\n\n", smallHash, hash.hashOf[p[0]])
			}

			break
		}
		for i, b := range p {

			if debug {
				debug.logf(depth, "Iter: %d — B = %q\n\n", i, b)

				if depth == 0 {
					debug.log(depth, "Current lowest:\n")
					debug.logHashes(depth, lowest, h.Size())
				}
			}

			hashP := hash.clone()
			hashP.set(b, hashTuple(h, hashP.hashOf[b], []byte{'@'}))
			hashPP, ok := hashBNodes(statements, h, zero, hashP.hashOf)
			if ok {

				if debug {
					debug.log(depth, "hashPP is trivial")
					debug.log(depth, "comparing hashPP\n")
					debug.logHashes(depth, hashPP.hashOf, h.Size())
					debug.log(depth, "with previous\n")
					debug.logHashes(depth, lowest, h.Size())
				}

				if lowest == nil || graphLess(statements, hashPP.hashOf, lowest) {
					lowest = hashPP.hashOf
					debug.log(depth, "choose hashPP\n")
				}
			} else {
				partsP := appendOrdered(byLengthHash{}, hashPP.termsFor)
				sort.Sort(partsP)

				if debug {
					debug.log(depth, "Parts':")
					debug.logParts(depth, partsP)
					debug.log(depth, "Recursive distinguish")
					debug.log(depth, "Called with current lowest:\n")
					debug.logHashes(depth, lowest, h.Size())
				}

				lowest = distinguish(statements, dist, h, zero, hashPP, partsP, lowest, depth+1)
			}
		}
	}

	if debug {
		debug.log(depth, "Return lowest:")
		debug.logHashes(depth, lowest, h.Size())
	}

	return lowest
}

// terms ordered syntactically, triples ordered lexicographically, and graphs
// ordered such that G < H if and only if G ⊂ H or there exists a triple
// t ∈ G \ H such that no triple t' ∈ H \ G exists where t' < t.
// p9 https://doi.org/10.1145/3068333
func graphLess(statements []*Statement, a, b map[string][]byte) bool {
	g := newLexicalStatements(statements, a)
	sort.Sort(g)
	h := newLexicalStatements(statements, b)
	sort.Sort(h)

	gSubH := sub(g, h, len(g.statements))
	if len(gSubH) == 0 {
		return true
	}

	hSubG := sub(h, g, 1)
	if len(hSubG) == 0 {
		return true
	}
	lowestH := relabeledStatement{hSubG[0], h.hashes}

	for _, s := range gSubH {
		rs := relabeledStatement{s, g.hashes}
		if rs.less(lowestH) {
			return true
		}
	}
	return false
}

// lexicalStatements is a sort implementation for Statements with blank
// node labels replaced with their hash.
type lexicalStatements struct {
	statements []*Statement
	hashes     map[string][]byte
}

func newLexicalStatements(statements []*Statement, hash map[string][]byte) lexicalStatements {
	s := lexicalStatements{
		statements: make([]*Statement, len(statements)),
		hashes:     hash,
	}
	copy(s.statements, statements)
	return s
}

// sub returns the difference between a and b up to max elements long.
func sub(a, b lexicalStatements, max int) []*Statement {
	var d []*Statement
	var i, j int
	for i < len(a.statements) && j < len(b.statements) && len(d) < max {
		ra := relabeledStatement{a.statements[i], a.hashes}
		rb := relabeledStatement{b.statements[j], b.hashes}
		switch {
		case ra.less(rb):
			d = append(d, a.statements[i])
			i++
		case rb.less(ra):
			j++
		default:
			i++
		}
	}
	if len(d) < max {
		d = append(d, a.statements[i:min(len(a.statements), i+max-len(d))]...)
	}
	return d
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s lexicalStatements) Len() int { return len(s.statements) }
func (s lexicalStatements) Less(i, j int) bool {
	return relabeledStatement{s.statements[i], s.hashes}.less(relabeledStatement{s.statements[j], s.hashes})
}
func (s lexicalStatements) Swap(i, j int) {
	s.statements[i], s.statements[j] = s.statements[j], s.statements[i]
}

// relabeledStatement is a statement that is orderable by its blank node
// hash relabeling.
type relabeledStatement struct {
	statement *Statement
	labels    map[string][]byte
}

func (a relabeledStatement) less(b relabeledStatement) bool {
	switch {
	case relabeledTerm{a.statement.Subject, a.labels}.less(relabeledTerm{b.statement.Subject, b.labels}):
		return true
	case relabeledTerm{b.statement.Subject, b.labels}.less(relabeledTerm{a.statement.Subject, a.labels}):
		return false
	}
	switch { // Always IRI.
	case a.statement.Predicate.Value < b.statement.Predicate.Value:
		return true
	case a.statement.Predicate.Value > b.statement.Predicate.Value:
		return false
	}
	switch {
	case relabeledTerm{a.statement.Object, a.labels}.less(relabeledTerm{b.statement.Object, b.labels}):
		return true
	case relabeledTerm{b.statement.Object, b.labels}.less(relabeledTerm{a.statement.Object, a.labels}):
		return false
	}
	return relabeledTerm{a.statement.Label, a.labels}.less(relabeledTerm{b.statement.Label, b.labels})
}

func (s relabeledStatement) String() string {
	subj := relabeledTerm{term: s.statement.Subject, labels: s.labels}
	obj := relabeledTerm{term: s.statement.Object, labels: s.labels}
	if s.statement.Label.Value == "" {
		return fmt.Sprintf("%s %s %s .", subj, s.statement.Predicate.Value, obj)
	}
	lab := relabeledTerm{term: s.statement.Label, labels: s.labels}
	return fmt.Sprintf("%s %s %s %s .", subj, s.statement.Predicate.Value, obj, lab)
}

// relabeledTerm is a term that is orderable by its blank node hash relabeling.
type relabeledTerm struct {
	term   Term
	labels map[string][]byte
}

func (a relabeledTerm) less(b relabeledTerm) bool {
	aIsBlank := isBlank(a.term.Value)
	bIsBlank := isBlank(b.term.Value)
	switch {
	case aIsBlank && bIsBlank:
		return bytes.Compare(a.labels[a.term.Value], b.labels[b.term.Value]) < 0
	case aIsBlank:
		return blankPrefix < unquoteIRI(b.term.Value)
	case bIsBlank:
		return unquoteIRI(a.term.Value) < blankPrefix
	default:
		return unquoteIRI(a.term.Value) < unquoteIRI(b.term.Value)
	}
}

func unquoteIRI(s string) string {
	if len(s) > 1 && s[0] == '<' && s[len(s)-1] == '>' {
		s = s[1 : len(s)-1]
	}
	return s
}

func (t relabeledTerm) String() string {
	if !isBlank(t.term.Value) {
		return t.term.Value
	}
	h, ok := t.labels[t.term.Value]
	if !ok {
		return t.term.Value + "_missing_hash"
	}
	return fmt.Sprintf("_:%0x", h)
}
