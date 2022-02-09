// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var deduplicateTests = []struct {
	statements string
	want       string
}{
	{},
	{
		statements: `
_:1 <p:a> _:2 .
`,
		want: `_:1 <p:a> _:2 .
`,
	},
	{
		statements: `
_:1 <p:a> _:2 .
_:1 <p:a> _:2 .
`,
		want: `_:1 <p:a> _:2 .
`,
	},
	{
		statements: `
_:1 <p:a> _:2 .
_:2 <p:a> _:1 .
_:1 <p:a> _:2 .
`,
		want: `_:1 <p:a> _:2 .
_:2 <p:a> _:1 .
`,
	},
	{
		statements: `
_:1 <p:a> _:2 .
_:2 <p:a> _:1 .
_:1 <p:a> _:2 .
_:1 <p:a> _:2 .
`,
		want: `_:1 <p:a> _:2 .
_:2 <p:a> _:1 .
`,
	},
}

func TestDeduplicate(t *testing.T) {
tests:
	for i, test := range deduplicateTests {
		var statements []*Statement
		dec := NewDecoder(strings.NewReader(test.statements))
		for {
			s, err := dec.Unmarshal()
			if err != nil {
				if err != io.EOF {
					t.Errorf("error during decoding: %v", err)
					continue tests
				}
				break
			}
			statements = append(statements, s)
		}

		var buf strings.Builder
		for _, s := range Deduplicate(statements) {
			fmt.Fprintln(&buf, s)
		}

		got := buf.String()
		if got != test.want {
			t.Errorf("unexpected result for test %d:\n%s", i, cmp.Diff(got, test.want))
		}
	}
}

func TestURNA(t *testing.T) {
	glob, err := filepath.Glob(filepath.Join("testdata", *tests))
	if err != nil {
		t.Fatalf("Failed to open test suite: %v", err)
	}
	for _, test := range []struct {
		name  string
		fn    func(dst, src []*Statement) ([]*Statement, error)
		truth string
	}{
		{
			name:  "URDNA2015",
			fn:    URDNA2015,
			truth: "-urdna2015.nq",
		},
		{
			name:  "URGNA2012",
			fn:    URGNA2012,
			truth: "-urgna2012.nq",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, path := range glob {
				name := filepath.Base(path)
				golden := strings.TrimSuffix(path, "-in.nq") + test.truth
				want, err := os.ReadFile(golden)
				if err != nil {
					if !os.IsNotExist(err) {
						t.Errorf("Failed to read golden data: %v", err)
					}
					continue
				}

				t.Run(name, func(t *testing.T) {
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

					relabeled, _ := test.fn(nil, statements)

					var buf bytes.Buffer
					for _, s := range relabeled {
						fmt.Fprintln(&buf, s)
					}
					got := buf.Bytes()

					if !bytes.Equal(got, want) {
						t.Errorf("Unexpected result for %s %s:\ngot:\n%s\nwant:\n%s",
							test.name, path, got, want)
					}
				})
			}
		})
	}
}

func BenchmarkURNA(b *testing.B) {
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

			for _, bench := range []struct {
				name string
				fn   func(dst, src []*Statement) ([]*Statement, error)
			}{
				{
					name: "URDNA2015",
					fn:   URDNA2015,
				},
				{
					name: "URGNA2012",
					fn:   URGNA2012,
				},
			} {
				b.Run(bench.name, func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						relabeled, _ := bench.fn(nil, statements)
						if len(relabeled) != len(statements) {
							b.Fatalf("unexpected number of relabeled statements: %d != %d", len(relabeled), len(statements))
						}
					}
				})
			}
		})
	}
}
