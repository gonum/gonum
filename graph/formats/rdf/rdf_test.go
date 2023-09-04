// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRDFWorkingGroupSuite(t *testing.T) {
	for _, file := range []string{
		"ntriple_tests.tar.gz",
		"nquad_tests.tar.gz",
	} {
		suite, err := os.Open(filepath.Join("testdata", file))
		if err != nil {
			t.Fatalf("Failed to open test suite in %q: %v", file, err)
		}
		defer suite.Close()

		r, err := gzip.NewReader(suite)
		if err != nil {
			t.Fatalf("Failed to uncompress test suite in %q: %v", file, err)
		}

		tr := tar.NewReader(r)
		for {
			h, err := tr.Next()
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Fatalf("Unexpected error while reading suite archive: %v", err)
			}

			h.Name = filepath.Base(h.Name)
			if filepath.Ext(h.Name) != ".nt" && filepath.Ext(h.Name) != ".nq" {
				continue
			}
			if _, ok := testSuite[h.Name]; !ok {
				t.Errorf("Missing test suite item %q", h.Name)
				continue
			}

			isBad := strings.Contains(h.Name, "bad")

			var got []statement
			dec := NewDecoder(tr)
			for i := 0; ; i++ {
				s, err := dec.Unmarshal()
				if err == io.EOF {
					break
				}
				gotBad := err != nil
				if gotBad != isBad {
					t.Errorf("Unexpected error return for test suite item %q, got: %v", h.Name, err)
				}

				var subj, pred, obj, lab term
				if s != nil {
					subj.text, subj.qual, subj.kind, _ = s.Subject.Parts()
					pred.text, pred.qual, pred.kind, _ = s.Predicate.Parts()
					obj.text, obj.qual, obj.kind, _ = s.Object.Parts()
					lab.text, lab.qual, lab.kind, _ = s.Label.Parts()
					if lab.text == "" {
						lab = term{}
					}
					got = append(got, statement{testSuite[h.Name][i].input, subj, pred, obj, lab})
				}

				if !gotBad {
					_, err = ParseNQuad(s.String())
					if err != nil {
						t.Errorf("Unexpected error return for valid statement in test suite item %q (%#v) got: %v rendered as\n%[2]s", h.Name, s, err)
					}

					st, err := termFor(subj.text, subj.qual, subj.kind)
					if err != nil {
						t.Errorf("Unexpected error return for valid subject in test suite item %q (%#v) got: %v rendered as\n%[2]s", h.Name, s, err)
					}
					pt, err := termFor(pred.text, pred.qual, pred.kind)
					if err != nil {
						t.Errorf("Unexpected error return for valid predicate in test suite item %q (%#v) got: %v rendered as\n%[2]s", h.Name, s, err)
					}
					ot, err := termFor(obj.text, obj.qual, obj.kind)
					if err != nil {
						t.Errorf("Unexpected error return for valid object in test suite item %q (%#v) got: %v rendered as\n%[2]s", h.Name, s, err)
					}
					lt, err := termFor(lab.text, lab.qual, lab.kind)
					if err != nil {
						t.Errorf("Unexpected error return for valid label in test suite item %q (%#v) got: %v rendered as\n%[2]s", h.Name, s, err)
					}

					// We can't check that we recreate the original from the test suite
					// due to escaping, but we can check for a second pass through the
					// round-trip.
					c := &Statement{Subject: st, Predicate: pt, Object: ot, Label: lt}
					pc, err := ParseNQuad(c.String())
					if err != nil {
						t.Errorf("Unexpected error return for reconstructed statement in test suite item %q (%#v) got: %v rendered as\n%[2]s", h.Name, s, err)
					}
					if !reflect.DeepEqual(c, pc) {
						t.Errorf("Unexpected reconstruction:\norig:  %#v\ncons:  %#v\nparsed:%#v", s, c, pc)
					}
				}
			}

			if !reflect.DeepEqual(testSuite[h.Name], got) {
				t.Errorf("Unexpected result for test suite item %q", h.Name)
			}
		}
	}
}

func termFor(text, qual string, kind Kind) (Term, error) {
	switch kind {
	case Invalid:
		return Term{}, nil
	case Blank:
		return NewBlankTerm(text)
	case IRI:
		return NewIRITerm(text)
	case Literal:
		return NewLiteralTerm(text, qual)
	default:
		panic(fmt.Sprintf("bad test kind=%d", kind))
	}
}

var escapeSequenceTests = []struct {
	escaped      string
	unEscaped    string
	canRoundTrip bool
}{
	{escaped: `plain text!`, unEscaped: "plain text!", canRoundTrip: true},
	{escaped: `\t`, unEscaped: "\t", canRoundTrip: false},
	{escaped: `\b`, unEscaped: "\b", canRoundTrip: false},
	{escaped: `\n`, unEscaped: "\n", canRoundTrip: true},
	{escaped: `\r`, unEscaped: "\r", canRoundTrip: true},
	{escaped: `\f`, unEscaped: "\f", canRoundTrip: false},
	{escaped: `\\`, unEscaped: "\\", canRoundTrip: true},
	{escaped: `\u0080`, unEscaped: "\u0080", canRoundTrip: true},
	{escaped: `\U00000080`, unEscaped: "\u0080", canRoundTrip: false},
	{escaped: `\t\b\n\r\f\"'\\`, unEscaped: "\t\b\n\r\f\"'\\", canRoundTrip: false},

	{escaped: `\t\u0080`, unEscaped: "\t\u0080", canRoundTrip: false},
	{escaped: `\b\U00000080`, unEscaped: "\b\u0080", canRoundTrip: false},
	{escaped: `\u0080\n`, unEscaped: "\u0080\n", canRoundTrip: true},
	{escaped: `\U00000080\r`, unEscaped: "\u0080\r", canRoundTrip: false},
	{escaped: `\u00b7\f\U000000b7`, unEscaped: "·\f·", canRoundTrip: false},
	{escaped: `\U000000b7\\\u00b7`, unEscaped: "·\\·", canRoundTrip: false},
	{escaped: `\U00010105\\\U00010106`, unEscaped: "\U00010105\\\U00010106", canRoundTrip: true},
}

func TestUnescape(t *testing.T) {
	for _, test := range escapeSequenceTests {
		got := unEscape([]rune(test.escaped))
		if got != test.unEscaped {
			t.Errorf("Failed to properly unescape %q, got:%q want:%q", test.escaped, got, test.unEscaped)
		}

		if test.canRoundTrip {
			got = escape("", test.unEscaped, "")
			if got != test.escaped {
				t.Errorf("Failed to properly escape %q, got:%q want:%q", test.unEscaped, got, test.escaped)
			}
			got = escape(`"`, test.unEscaped, `"`)
			if got != `"`+test.escaped+`"` {
				t.Errorf("Failed to properly escape %q quoted, got:%q want:%q", test.unEscaped, got, `"`+test.escaped+`"`)
			}
		}
	}
}
