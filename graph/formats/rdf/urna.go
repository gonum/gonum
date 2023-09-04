// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"sort"

	"gonum.org/v1/gonum/stat/combin"
)

// Deduplicate removes duplicate statements in s, working in place, and returns
// the deduplicated slice with statements sorted in lexical order. Term UID
// fields are not considered and their values may be lost during deduplication.
func Deduplicate(s []*Statement) []*Statement {
	if len(s) < 2 {
		return s
	}
	sort.Sort(c14nStatements(s))
	curr := 0
	for i, e := range s {
		if isSameStatement(e, s[curr]) {
			continue
		}
		curr++
		if curr < i {
			s[curr], s[i] = s[i], nil
		}
	}
	return s[:curr+1]
}

func isSameStatement(a, b *Statement) bool {
	if a == b {
		return true
	}
	return a.Subject.Value == b.Subject.Value &&
		a.Predicate.Value == b.Predicate.Value &&
		a.Object.Value == b.Object.Value &&
		a.Label.Value == b.Label.Value
}

// Note on implementation details: The comment numbering in the code relates the
// implementation to the steps of the algorithm described in the specification.

// URGNA2012 applies the Universal RDF Graph Normalization Algorithm 2012
// to the statements in src, placing the result in dst and returning it.
// If dst is nil a slice of statements will be allocated. If dst is not
// nil and not the same length as src, URGNA2012 will return an error.
//
// See https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html for details.
func URGNA2012(dst, src []*Statement) ([]*Statement, error) {
	if dst == nil {
		dst = make([]*Statement, len(src))
	} else if len(dst) != len(src) {
		return dst, errors.New("rdf: slice length mismatch")
	}
	// 1. https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html#algorithm
	u := &urna{
		canon:         newIssuer("_:c14n"),
		hashes:        make(map[string]string),
		statementsFor: make(map[string][]*Statement),
		hash:          sha1.New(),
		label:         "_:g",
	}
	u.hashToRelated = u.hashToRelatedURGNA2012
	return u.relabel(dst, src)
}

// URDNA2015 applies the Universal RDF Dataset Normalization Algorithm 2015
// to the statements in src, placing the result in dst and returning it.
// If dst is nil a slice of statements will be allocated. If dst is not
// nil and not the same length as src, URDNA2015 will return an error.
//
// See https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html for details.
func URDNA2015(dst, src []*Statement) ([]*Statement, error) {
	if dst == nil {
		dst = make([]*Statement, len(src))
	} else if len(dst) != len(src) {
		return dst, errors.New("rdf: slice length mismatch")
	}
	// 1. https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html#algorithm
	u := &urna{
		canon:         newIssuer("_:c14n"),
		hashes:        make(map[string]string),
		statementsFor: make(map[string][]*Statement),
		hash:          sha256.New(),
	}
	u.hashToRelated = u.hashToRelatedURDNA2015
	return u.relabel(dst, src)
}

// urna is the canonicalization state for the URGNA2012 and URDNA2015
// algorithms. The urna type implements both algorithms through the state
// of the label and hashToRelated fields.
//
// See https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html#canonicalization-state
// for details.
type urna struct {
	// canon is the canonical issuer.
	canon *issuer

	// hashes holds already calculated hashes
	// for hashing first degree quads.
	hashes map[string]string

	// statementsFor is the blank node to quads map.
	statementsFor map[string][]*Statement

	// hash is the hash function used by the
	// canonicalization function.
	hash hash.Hash
	// hashToRelated holds URGNA2012 and URDNA2015-
	// specific hashing routines.
	hashToRelated relatedHashCreator
	// label holds "_:g" when running URGNA2012.
	// Otherwise it is empty.
	label string
}

// relabel is the algorithm described in section 4.4.2 of the spec at
// https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html#algorithm.
func (u *urna) relabel(dst, src []*Statement) ([]*Statement, error) {
	// termsFor is the hash to blank nodes map.
	// It is not held in the urna struct, but is
	// part of the canonicalization state.
	//
	// https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html#dfn-hash-to-blank-nodes-map
	var termsFor map[string][]string // 1.

	for _, s := range src { // 2.
	terms:
		for _, t := range []string{
			s.Subject.Value,
			s.Object.Value,
			s.Label.Value,
		} {
			if !isBlank(t) {
				continue
			}
			for _, e := range u.statementsFor[t] {
				if e == s {
					continue terms
				}
			}
			u.statementsFor[t] = append(u.statementsFor[t], s)
		}
	}

	// todo is the list of non-normalized blank node identifiers.
	todo := make(map[string]bool) // 3.
	for b := range u.statementsFor {
		todo[b] = true
	}

	simple := true // 4.
	for simple {   // 5.
		simple = false // 5.1

		termsFor = make(map[string][]string) // 5.2

		for b := range todo { // 5.3
			hash := u.hashFirstDegreeQuads(b)          // 5.3.1
			termsFor[hash] = append(termsFor[hash], b) // 5.3.2
		}

		for _, h := range lexicallySortedTermHashes(termsFor) { // 5.4
			terms := termsFor[h]
			if len(terms) > 1 { // 5.4.1
				continue
			}
			u.canon.issueFor(terms[0]) // 5.4.2
			delete(todo, terms[0])     // 5.4.3
			delete(termsFor, h)        // 5.4.4
			simple = true              // 5.4.5
		}
	}

	for _, hash := range lexicallySortedTermHashes(termsFor) { // 6.
		paths := make(map[string][]*issuer) // 6.1
		for _, b := range termsFor[hash] {  // 6.2
			if u.canon.has(b) { // 6.2.1
				continue
			}
			names := newIssuer("_:b") // 6.2.2
			names.issueFor(b)         // 6.2.3

			// 6.2.4
			hash, issuer := u.hashNDegreeQuads(b, names)
			paths[string(hash)] = append(paths[string(hash)], issuer)
		}

		for _, hash := range lexicallySortedPathHashes(paths) { // 6.3
			for _, i := range paths[hash] {
				for _, existing := range i.ordered { // 6.3.1
					u.canon.issueFor(existing)
				}
			}
		}
	}

	// 7.
	for i, s := range src {
		if dst[i] == nil {
			dst[i] = &Statement{}
		}
		n := dst[i]
		n.Subject = Term{Value: translateURNA(s.Subject.Value, u.canon.issued), UID: s.Subject.UID}
		n.Predicate = s.Predicate
		n.Object = Term{Value: translateURNA(s.Object.Value, u.canon.issued), UID: s.Object.UID}
		n.Label = Term{Value: translateURNA(s.Label.Value, u.canon.issued), UID: s.Label.UID}
	}
	sort.Sort(c14nStatements(dst))

	return dst, nil
}

// lexicallySortedPathHashes returns the lexically sorted hashes of paths.
func lexicallySortedPathHashes(paths map[string][]*issuer) []string {
	lexicalHashPaths := make([]string, len(paths))
	i := 0
	for h := range paths {
		lexicalHashPaths[i] = h
		i++
	}
	sort.Strings(lexicalHashPaths)
	return lexicalHashPaths
}

func translateURNA(term string, mapping map[string]string) string {
	term = translate(term, mapping)
	if term == "" {
		return ""
	}
	text, qual, kind, err := extract([]rune(term))
	var t Term
	switch kind {
	case Blank:
		return term
	case IRI:
		t, err = NewIRITerm(text)
	case Literal:
		t, err = NewLiteralTerm(text, qual)
	}
	if err != nil {
		panic(fmt.Errorf("rdf: invalid term %q: %w", term, err))
	}
	return t.Value
}

// hashFirstDegreeQuads is the algorithm described in section 4.6 of the spec
// at https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html#algorithm-1.
func (u *urna) hashFirstDegreeQuads(b string) string {
	if h, ok := u.hashes[b]; ok {
		return h
	}

	var statements []*Statement // 1.

	for _, s := range u.statementsFor[b] { // 2. and 3.
		var n Statement
		n.Subject.Value = replaceBlank(s.Subject.Value, b, "")
		n.Predicate.Value = s.Predicate.Value
		n.Object.Value = replaceBlank(s.Object.Value, b, "")
		n.Label.Value = replaceBlank(s.Label.Value, b, u.label)
		statements = append(statements, &n)
	}

	sort.Sort(c14nStatements(statements)) // 4.

	// 5.
	u.hash.Reset()
	for _, s := range statements {
		fmt.Fprintln(u.hash, s)
	}
	u.hashes[b] = string(hex(u.hash.Sum(nil)))

	return u.hashes[b]
}

// replaceBlank implements 3.1 of the algorithm described at
// https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html#algorithm-1.
func replaceBlank(b, matching, label string) string {
	if !isBlank(b) { // 3.1
		return b
	}
	if label != "" { // URGNA2012 modification.
		// When running in URGNA2012 mode, label is "_:g" for Label fields.
		//
		// If any blank node was used in the graph name position in the quad,
		// then the value was serialized using the special blank node identifier,
		// "_:g", instead of "_:z".
		// https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html#urgna2012
		return label
	}
	// 3.1.1.1
	if b == matching {
		return "_:a"
	}
	return "_:z"
}

// hashNDegreeQuads is the algorithm described in section 4.8 of the spec
// at https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html#hash-n-degree-quads.
func (u *urna) hashNDegreeQuads(b string, names *issuer) ([]byte, *issuer) {
	// termsFor is the hash to related blank nodes map.
	termsFor := u.hashToRelated(b, names) // 1., 2. and 3.
	var final []byte                      // 4.

	for _, hash := range lexicallySortedTermHashes(termsFor) { // 5.
		terms := termsFor[hash]
		final = append(final, hash...) // 5.1
		var chosenPath []byte          // 5.2
		var chosenIssuer *issuer       // 5.3
		p := newPermutations(terms)    // 5.4
	permutations:
		for p.next() {
			namesCopy := names.clone()          // 5.4.1
			var path []byte                     // 5.4.2
			var work []string                   // 5.4.3
			for _, b := range p.permutation() { // 5.4.4
				if u.canon.has(b) { // 5.4.4.1
					path = append(path, u.canon.issueFor(b)...)
				} else { // 5.4.4.1
					if !namesCopy.has(b) {
						work = append(work, b)
					}

					path = append(path, namesCopy.issueFor(b)...) // 5.4.4.2.2
				}

				// 5.4.4.3
				if len(chosenPath) != 0 && len(path) >= len(chosenPath) && bytes.Compare(path, chosenPath) > 0 {
					continue permutations
				}
			}

			for _, b := range work { // 5.4.5
				hash, issuer := u.hashNDegreeQuads(b, namesCopy) // 5.4.5.1
				path = append(path, namesCopy.issueFor(b)...)    // 5.4.5.2

				// 5.4.5.3
				path = append(path, '<')
				path = append(path, hash...)
				path = append(path, '>')

				namesCopy = issuer // 5.4.5.4

				// 5.4.5.5
				if len(chosenPath) != 0 && len(path) >= len(chosenPath) && bytes.Compare(path, chosenPath) > 0 {
					continue permutations
				}
			}

			if len(chosenPath) == 0 || bytes.Compare(path, chosenPath) < 0 { // 5.4.6
				chosenPath = path
				chosenIssuer = namesCopy
			}

		}
		// 5.5
		final = append(final, chosenPath...)
		u.hash.Reset()
		u.hash.Write(final)

		names = chosenIssuer // 5.6
	}

	return hex(u.hash.Sum(nil)), names
}

// lexicallySortedTermHashes returns the lexically sorted hashes of termsFor.
func lexicallySortedTermHashes(termsFor map[string][]string) []string {
	lexicalHashes := make([]string, len(termsFor))
	i := 0
	for h := range termsFor {
		lexicalHashes[i] = h
		i++
	}
	sort.Strings(lexicalHashes)
	return lexicalHashes
}

type relatedHashCreator func(b string, names *issuer) map[string][]string

// hashToRelatedURDNA2015 is the section 1. 2. and 3. of 4.8.2 of the spec
// at https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html#hash-n-degree-quads.
func (u *urna) hashToRelatedURDNA2015(b string, names *issuer) map[string][]string {
	// termsFor is the hash to related blank nodes map.
	termsFor := make(map[string][]string) // 1.

	for _, s := range u.statementsFor[b] { // 2. and 3.
		for i, term := range []string{ // 3.1
			s.Subject.Value,
			s.Object.Value,
			s.Label.Value,
		} {
			if !isBlank(term) || term == b {
				continue
			}

			// 3.1.1
			const position = "sog"
			hash := u.hashRelatedBlank(term, s, names, position[i])

			// 3.1.2
			termsFor[string(hash)] = append(termsFor[string(hash)], term)
		}
	}

	return termsFor
}

// hashToRelatedURGNA2012 is the section 1., 2. and 3. of 4.8.2 of the spec
// at https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html#hash-n-degree-quads
// with changes made for URGNA2012 shown in the appendix for 4.8 at
// https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html#urgna2012.
// The numbering of steps here corresponds to the spec's numbering in the
// appendix.
func (u *urna) hashToRelatedURGNA2012(b string, names *issuer) map[string][]string {
	// termsFor is the hash to related blank nodes map.
	termsFor := make(map[string][]string)

	for _, s := range u.statementsFor[b] { // 1.
		var (
			term string
			pos  byte
		)
		switch {
		case isBlank(s.Subject.Value) && s.Subject.Value != b: // 1.1
			term = s.Subject.Value
			pos = 'p'
		case isBlank(s.Object.Value) && s.Object.Value != b: // 1.2
			term = s.Object.Value
			pos = 'r'
		default:
			continue // 1.3
		}

		// 1.4
		hash := u.hashRelatedBlank(term, s, names, pos)
		termsFor[string(hash)] = append(termsFor[string(hash)], term)
	}

	return termsFor
}

// hashNDegreeQuads is the algorithm described in section 4.7 of the spec
// https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html#hash-related-blank-node.
func (u *urna) hashRelatedBlank(term string, s *Statement, names *issuer, pos byte) []byte {
	// 1.
	var b string
	switch {
	case u.canon.has(term):
		b = u.canon.issueFor(term)
	case names.has(term):
		b = names.issueFor(term)
	default:
		b = u.hashFirstDegreeQuads(term)
	}

	// 2.
	u.hash.Reset()
	u.hash.Write([]byte{pos})

	if pos != 'g' { // 3.
		if u.label == "" {
			// URDNA2015: Term.Value retained the angle quotes
			// so we don't need to add them.
			u.hash.Write([]byte(s.Predicate.Value))
		} else {
			// URGNA2012 does not delimit predicate by < and >.
			// https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html#urgna2012
			// with reference to 4.7.
			u.hash.Write([]byte(unquoteIRI(s.Predicate.Value)))
		}
	}

	// 4. and 5.
	u.hash.Write([]byte(b))
	return hex(u.hash.Sum(nil))
}

// issuer is an identifier issuer.
type issuer struct {
	prefix  string
	issued  map[string]string
	ordered []string
}

// newIssuer returns a new identifier issuer with the given prefix.
func newIssuer(prefix string) *issuer {
	return &issuer{prefix: prefix, issued: make(map[string]string)}
}

// issueFor implements the issue identifier algorithm.
//
// See https://json-ld.github.io/rdf-dataset-canonicalization/spec/index.html#issue-identifier-algorithm
func (i *issuer) issueFor(b string) string {
	c, ok := i.issued[b]
	if ok {
		return c
	}
	c = fmt.Sprintf("%s%d", i.prefix, len(i.issued))
	i.issued[b] = c
	i.ordered = append(i.ordered, b)
	return c
}

func (i *issuer) has(id string) bool {
	_, ok := i.issued[id]
	return ok
}

func (i *issuer) clone() *issuer {
	new := issuer{
		prefix:  i.prefix,
		issued:  make(map[string]string, len(i.issued)),
		ordered: make([]string, len(i.ordered)),
	}
	copy(new.ordered, i.ordered)
	for k, v := range i.issued {
		new.issued[k] = v
	}
	return &new
}

func hex(data []byte) []byte {
	const digit = "0123456789abcdef"
	buf := make([]byte, 0, len(data)*2)
	for _, b := range data {
		buf = append(buf, digit[b>>4], digit[b&0xf])
	}
	return buf
}

// permutations is a string permutation generator.
type permutations struct {
	src  []string
	dst  []string
	idx  []int
	perm *combin.PermutationGenerator
}

// newPermutation returns a new permutations.
func newPermutations(src []string) *permutations {
	return &permutations{
		src:  src,
		dst:  make([]string, len(src)),
		perm: combin.NewPermutationGenerator(len(src), len(src)),
		idx:  make([]int, len(src)),
	}
}

// next returns whether there is another permutation available.
func (p *permutations) next() bool {
	return p.perm.Next()
}

// permutation returns the permutation. The caller may not retain the
// returned slice between iterations.
func (p *permutations) permutation() []string {
	p.perm.Permutation(p.idx)
	for i, j := range p.idx {
		p.dst[j] = p.src[i]
	}
	return p.dst
}
