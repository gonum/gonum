// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf

import (
	"crypto/md5"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"text/tabwriter"
	"time"

	"golang.org/x/exp/rand"
)

var (
	origSeed = flag.Int64("seed", 1, "specify random seed to use for each test (negative for Unix time)")
	tests    = flag.String("test", "*-in.n[qt]", "specify test case in testdata")
)

func TestIsoCanonicalHashes(t *testing.T) {
	seed := uint64(*origSeed)
	if *origSeed < 0 {
		seed = uint64(time.Now().UnixNano())
	}
	defer func() {
		if t.Failed() && *origSeed < 0 {
			t.Logf("time based seed: %d", seed)
		}
	}()

	// Number of times to run IsoCanonicalHashes to check consistency.
	const retries = 5

	// Share a global hash function to ensure that we
	// are resetting the function internally on each use.
	hash := md5.New()

	glob, err := filepath.Glob(filepath.Join("testdata", *tests))
	if err != nil {
		t.Fatalf("Failed to open test suite: %v", err)
	}
	for _, path := range glob {
		name := filepath.Base(path)
		t.Run(name, func(t *testing.T) {
			src := rand.NewSource(seed)

			f, err := os.Open(path)
			if err != nil {
				t.Fatalf("Failed to open test suite in %q: %v", path, err)
			}
			var statements []*Statement
			dec := NewDecoder(f)
			for {
				s, err := dec.Unmarshal()
				if err != nil {
					if err == io.EOF {
						break
					}
					t.Fatalf("Unexpected error reading from %q: %v", path, err)
				}
				statements = append(statements, s)
			}
			f.Close()

			for _, decomp := range []bool{false, true} {
				t.Run(fmt.Sprintf("decomp=%t", decomp), func(t *testing.T) {
					var last map[string][]byte
					for i := 0; i < retries; i++ {
						curr, terms := IsoCanonicalHashes(statements, decomp, true, hash, make([]byte, 16))
						if !hashesDisjoint(terms) {
							t.Errorf("IsoCanonicalHashes did not uniquely identify nodes %q with decomp=%t",
								name, decomp)
						}
						if last != nil {
							last := relabelStatements(statements, termsFor(last, hash))
							sort.Sort(simpleLexicalStatements(last))

							curr := relabelStatements(statements, termsFor(curr, hash))
							sort.Sort(simpleLexicalStatements(curr))

							if !reflect.DeepEqual(last, curr) {
								t.Errorf("IsoCanonicalHashes was not stable between runs on %q with decomp=%t",
									name, decomp)

								t.Log("Current run:")
								for _, s := range curr {
									t.Logf("\t%s", s)
								}

								t.Log("Previous run:")
								for _, s := range last {
									t.Logf("\t%s", s)
								}

								break
							}
						}
						last = curr
					}
					hashes := last
					ok := allUnique(hashes)
					if !ok {
						t.Errorf("Failed to get unique hashes for %q disjoint with decomp=%t", name, decomp)
						t.Logf("skipping %q decomp=%t", path, decomp)
						return
					}

					// Test that a graph is not isomorphic with one generated
					// by deleting the last statement.
					t.Run("isomorphic G != G-s", func(t *testing.T) {
						if len(statements) == 0 {
							return
						}
						if Isomorphic(statements, statements[:len(statements)-1], decomp, hash) {
							t.Error("Isomorphic(G, G-s)=true")
						}
					})

					// Test that a graph is not isomorphic with one generated
					// by hashing the first grounded statement.
					t.Run("isomorphic G != Gμ(g)", func(t *testing.T) {
						mangled, mangTerms := mangleFirstIL(statements, hash)
						if mangTerms == nil {
							// All terms were blanks.
							return
						}
						if Isomorphic(statements, mangled, decomp, hash) {
							t.Error("Isomorphic(G, Gμ(g))=true")
						}
					})

					// Test that a graph is not isomorphic with one generated
					// by merging the first two lexically sorted blank nodes
					// into one.
					t.Run("isomorphic G != G(b1∪b2)", func(t *testing.T) {
						mangled, mangTerms := mergeFirst2B(statements)
						if mangTerms == nil {
							// All terms were blanks.
							return
						}
						if Isomorphic(statements, mangled, decomp, hash) {
							t.Error("Isomorphic(G, G(b1∪b2))=true")
						}
					})

					// Relabel a copy of the statements and then sort.
					orig := relabelStatements(statements, termsFor(hashes, hash))
					sort.Sort(simpleLexicalStatements(orig))

					for _, perm := range []struct {
						name string
						data func() ([]*Statement, map[string]string)
					}{
						{
							name: "reverse statements",
							data: func() ([]*Statement, map[string]string) { return reverseStatements(statements) },
						},
						{
							name: "permute statements",
							data: func() ([]*Statement, map[string]string) { return permuteStatements(statements, src) },
						},
						{
							name: "permute blank labels",
							data: func() ([]*Statement, map[string]string) { return permuteBlanks(statements, src) },
						},
						{
							name: "hash blank labels",
							data: func() ([]*Statement, map[string]string) { return hashBlanks(statements, md5.New()) },
						},
						{
							name: "reverse statements and hash blank labels",
							data: func() ([]*Statement, map[string]string) {
								// Reordering must come first since it does not return
								// a non-nil terms map, but hashBlanks does.
								s, _ := reverseStatements(statements)
								return hashBlanks(s, md5.New())
							},
						},
						{
							name: "permute statements and hash blank labels",
							data: func() ([]*Statement, map[string]string) {
								// Reordering must come first since it does not return
								// a non-nil terms map, but hashBlanks does.
								s, _ := permuteStatements(statements, src)
								return hashBlanks(s, md5.New())
							},
						},
					} {
						t.Run(perm.name, func(t *testing.T) {
							if debug {
								fmt.Fprintf(os.Stderr, "\n%q %q decomp=%t:\n", path, perm.name, decomp)
							}

							altStatements, terms := perm.data()
							altHashes, altTerms := IsoCanonicalHashes(altStatements, decomp, true, hash, make([]byte, 16))
							ok := allUnique(altHashes) && hashesDisjoint(altTerms)
							if !ok {
								t.Errorf("Failed to get unique hashes for %q alternative disjoint %q with decomp=%t",
									path, perm.name, decomp)
							}

							if debug {
								fmt.Fprintln(os.Stderr, "Name mappings from original dataset:")
								keys := make([]string, len(hashes))
								var i int
								for k := range hashes {
									keys[i] = k
									i++
								}
								sort.Strings(keys)
								w := tabwriter.NewWriter(os.Stderr, 0, 4, 8, ' ', 0)
								for _, k := range keys {
									fmt.Fprintf(w, "\t%s\t%s\n", k, translate(k, terms))
								}
								w.Flush()
								fmt.Fprintln(os.Stderr)
							}

							// Relabel a copy of the alternative statements and then sort.
							alt := relabelStatements(altStatements, termsFor(altHashes, hash))
							sort.Sort(simpleLexicalStatements(alt))

							for i := range statements {
								if *orig[i] != *alt[i] { // Otherwise we have pointer inequality.
									t.Errorf("Unexpected statement in %q %q decomp=%t:\ngot: %#v\nwant:%#v",
										path, perm.name, decomp, orig[i], alt[i])

									break
								}
							}

							if !Isomorphic(statements, altStatements, decomp, hash) {
								t.Errorf("Isomorphic(G, perm(G))=false in %q %q decomp=%t",
									path, perm.name, decomp)
							}
						})
					}
				})
			}
		})
	}
}

func permuteStatements(s []*Statement, src rand.Source) ([]*Statement, map[string]string) {
	rnd := rand.New(src)
	m := make([]*Statement, len(s))
	for x, y := range rnd.Perm(len(s)) {
		m[x] = s[y]
	}
	return m, nil
}

func reverseStatements(s []*Statement) ([]*Statement, map[string]string) {
	m := make([]*Statement, len(s))
	for i, j := 0, len(s)-1; i < len(s); i, j = i+1, j-1 {
		m[j] = s[i]
	}
	return m, nil
}

func permuteBlanks(s []*Statement, src rand.Source) ([]*Statement, map[string]string) {
	rnd := rand.New(src)
	terms := make(map[string]string)
	for _, e := range s {
		for _, t := range []string{
			e.Subject.Value,
			e.Predicate.Value,
			e.Object.Value,
			e.Label.Value,
		} {
			if t == "" {
				continue
			}
			terms[t] = t
		}
	}

	var blanks []string
	for t := range terms {
		if isBlank(t) {
			blanks = append(blanks, t)
		}
	}
	sort.Strings(blanks)
	for x, y := range rnd.Perm(len(blanks)) {
		terms[blanks[x]] = blanks[y]
	}

	m := relabelStatements(s, terms)
	return m, terms
}

func hashBlanks(s []*Statement, h hash.Hash) ([]*Statement, map[string]string) {
	terms := make(map[string]string)
	for _, e := range s {
		for _, t := range []string{
			e.Subject.Value,
			e.Predicate.Value,
			e.Object.Value,
			e.Label.Value,
		} {
			if !isBlank(t) {
				continue
			}
			h.Reset()
			h.Write([]byte(t))
			terms[t] = fmt.Sprintf("_:%0*x", 2*h.Size(), h.Sum(nil))
		}
	}

	m := relabelStatements(s, terms)
	return m, terms
}

func mangleFirstIL(s []*Statement, h hash.Hash) ([]*Statement, map[string]string) {
	terms := make(map[string]string)
	for _, e := range s {
		for _, t := range []string{
			e.Subject.Value,
			e.Predicate.Value,
			e.Object.Value,
			e.Label.Value,
		} {
			if isBlank(t) {
				continue
			}
			h.Reset()
			h.Write([]byte(t))
			terms[t] = fmt.Sprintf(`"%0*x"`, 2*h.Size(), h.Sum(nil))
			return relabelStatements(s, terms), terms
		}
	}

	m := relabelStatements(s, nil)
	return m, nil
}

func mergeFirst2B(s []*Statement) ([]*Statement, map[string]string) {
	terms := make(map[string]string)
	for _, e := range s {
		for _, t := range []string{
			e.Subject.Value,
			e.Predicate.Value,
			e.Object.Value,
			e.Label.Value,
		} {
			if !isBlank(t) {
				continue
			}
			terms[t] = t
		}
	}
	if len(terms) < 2 {
		return relabelStatements(s, nil), nil
	}

	blanks := make([]string, len(terms))
	i := 0
	for _, b := range terms {
		blanks[i] = b
		i++
	}
	sort.Strings(blanks)
	terms[blanks[1]] = terms[blanks[0]]

	m := relabelStatements(s, terms)
	return m, nil
}

func hashesDisjoint(terms map[string]map[string]bool) bool {
	for _, t := range terms {
		if len(t) != 1 {
			return false
		}
	}
	return true
}

func TestLexicalStatements(t *testing.T) {
	if *tests == "" {
		*tests = "*"
	}

	hash := md5.New()

	glob, err := filepath.Glob(filepath.Join("testdata", *tests))
	if err != nil {
		t.Fatalf("Failed to open test suite: %v", err)
	}
	for _, path := range glob {
		f, err := os.Open(path)
		if err != nil {
			t.Fatalf("Failed to open test suite in %q: %v", path, err)
		}
		var statements []*Statement
		dec := NewDecoder(f)
		for {
			s, err := dec.Unmarshal()
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Fatalf("Unexpected error reading from %q: %v", path, err)
			}
			statements = append(statements, s)
		}
		f.Close()

		for _, decomp := range []bool{false, true} {
			hashes, _ := IsoCanonicalHashes(statements, decomp, true, hash, make([]byte, 16))

			terms := termsFor(hashes, hash)

			// Sort a copy of the statements based on hashes and then relabel.
			indirect := make([]*Statement, len(statements))
			copy(indirect, statements)
			sort.Sort(lexicalStatements{indirect, hashes})
			indirect = relabelStatements(indirect, terms)

			// Relabel a copy of the statements and then sort.
			direct := relabelStatements(statements, terms)
			sort.Sort(simpleLexicalStatements(direct))

			for i := range statements {
				if *indirect[i] != *direct[i] { // Otherwise we have pointer inequality.
					t.Errorf("Unexpected ordering of indirect sort in %q:\ngot: %#v\nwant:%#v",
						path, indirect[i], direct[i])
				}
			}
		}
	}
}

func termsFor(hashes map[string][]byte, hash hash.Hash) map[string]string {
	terms := make(map[string]string)
	for t, h := range hashes {
		if isBlank(t) {
			terms[t] = fmt.Sprintf("_:%0*x", 2*hash.Size(), h)
		}
	}
	return terms
}

// simpleLexicalStatements implements lexical statement sorting on the
// literal values without interpolation.
type simpleLexicalStatements []*Statement

func (s simpleLexicalStatements) Len() int { return len(s) }
func (s simpleLexicalStatements) Less(i, j int) bool {
	si := s[i]
	sj := s[j]
	switch {
	case unquoteIRI(si.Subject.Value) < unquoteIRI(sj.Subject.Value):
		return true
	case unquoteIRI(si.Subject.Value) > unquoteIRI(sj.Subject.Value):
		return false
	}
	switch { // Always IRI.
	case si.Predicate.Value < sj.Predicate.Value:
		return true
	case si.Predicate.Value > sj.Predicate.Value:
		return false
	}
	switch {
	case unquoteIRI(si.Object.Value) < unquoteIRI(sj.Object.Value):
		return true
	case unquoteIRI(si.Object.Value) > unquoteIRI(sj.Object.Value):
		return false
	}
	return unquoteIRI(si.Label.Value) < unquoteIRI(sj.Label.Value)
}
func (s simpleLexicalStatements) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func relabelStatements(s []*Statement, terms map[string]string) []*Statement {
	m := make([]*Statement, len(s))
	for i, e := range s {
		n := *e
		n.Subject = Term{Value: translate(n.Subject.Value, terms)}
		n.Predicate = Term{Value: translate(n.Predicate.Value, terms)}
		n.Object = Term{Value: translate(n.Object.Value, terms)}
		n.Label = Term{Value: translate(n.Label.Value, terms)}
		m[i] = &n
	}
	return m
}

func BenchmarkIsoCanonicalHashes(b *testing.B) {
	hash := md5.New()

	benchmarks := []string{
		"test019-in.nq",
		"test044-in.nq",
	}

	for _, name := range benchmarks {
		path := filepath.Join("testdata", name)
		b.Run(name, func(b *testing.B) {
			f, err := os.Open(path)
			if err != nil {
				b.Fatalf("Failed to open test suite in %q: %v", path, err)
			}
			var statements []*Statement
			dec := NewDecoder(f)
			for {
				s, err := dec.Unmarshal()
				if err != nil {
					if err == io.EOF {
						break
					}
					b.Fatalf("Unexpected error reading from %q: %v", path, err)
				}
				statements = append(statements, s)
			}
			f.Close()

			nodes := make(map[string]bool)
			for _, s := range statements {
				for _, t := range []string{
					s.Subject.Value,
					s.Predicate.Value,
					s.Object.Value,
					s.Label.Value,
				} {
					if t != "" {
						nodes[t] = true
					}
				}
			}
			n := len(nodes)

			for _, decomp := range []bool{false, true} {
				b.Run(fmt.Sprintf("decomp=%t", decomp), func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						hashes, _ := IsoCanonicalHashes(statements, decomp, true, hash, make([]byte, 16))
						if len(hashes) != n {
							b.Fatalf("unexpected number of hashes: %d != %d", len(hashes), len(statements))
						}
					}
				})
			}
		})
	}
}
