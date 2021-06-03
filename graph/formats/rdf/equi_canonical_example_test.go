// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf_test

import (
	"fmt"
	"strings"

	"gonum.org/v1/gonum/graph/formats/rdf"
)

func ExampleLean() {
	for i, statements := range []string{
		0: `
_:author1 <ex:contributesTo> _:go .
_:author2 <ex:contributesTo> _:go .
_:author2 <ex:contributesTo> _:gonum .
_:author3 <ex:contributesTo> _:gonum .
`,
		1: `
_:author1 <ex:contributesTo> _:go .
_:author2 <ex:contributesTo> _:go .
_:author2 <ex:contributesTo> _:gonum .
_:author3 <ex:contributesTo> _:gonum .
_:gonum <ex:dependsOn> _:go .
`,
		2: `
_:author1 <ex:contributesTo> _:go .
_:author1 <ex:notContributesTo> _:gonum .
_:author2 <ex:contributesTo> _:go .
_:author2 <ex:contributesTo> _:gonum .
_:author3 <ex:contributesTo> _:gonum .
_:gonum <ex:dependsOn> _:go .
`,
		3: `
_:author1 <ex:is> "Alice" .
_:author1 <ex:contributesTo> _:go .
_:author2 <ex:contributesTo> _:go .
_:author2 <ex:contributesTo> _:gonum .
_:author3 <ex:contributesTo> _:gonum .
_:gonum <ex:dependsOn> _:go .
`,
		4: `
_:author1 <ex:contributesTo> _:go .
_:author1 <ex:notContributesTo> _:gonum .
_:author2 <ex:contributesTo> _:go .
_:author2 <ex:contributesTo> _:gonum .
_:author3 <ex:contributesTo> _:gonum .
_:author3 <ex:notContributesTo> _:go .
_:gonum <ex:dependsOn> _:go .
`,
		5: `
_:author1 <ex:is> "Alice" .
_:author2 <ex:is> "Bob" .
_:author1 <ex:contributesTo> _:go .
_:author2 <ex:contributesTo> _:go .
_:author2 <ex:contributesTo> _:gonum .
_:author3 <ex:contributesTo> _:gonum .
_:gonum <ex:dependsOn> _:go .
`,
		6: `
_:author1 <ex:is> "Alice" .
_:author2 <ex:is> "Bob" .
_:author3 <ex:is> "Charlie" .
_:author1 <ex:contributesTo> _:go .
_:author2 <ex:contributesTo> _:go .
_:author2 <ex:contributesTo> _:gonum .
_:author3 <ex:contributesTo> _:gonum .
_:gonum <ex:dependsOn> _:go .
`,
	} {
		// Decode the statement stream.
		dec := rdf.NewDecoder(strings.NewReader(statements))
		var s []*rdf.Statement
		for {
			l, err := dec.Unmarshal()
			if err != nil {
				break
			}
			s = append(s, l)
		}

		// Lean the graph to remove redundant statements.
		lean, err := rdf.Lean(s)
		if err != nil {
			fmt.Println(err)
		}

		// Canonicalize the blank nodes in-place.
		_, err = rdf.URDNA2015(lean, lean)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("%d:\n", i)
		for _, s := range lean {
			fmt.Println(s)
		}
		fmt.Println()
	}

	// Output:
	//
	// 0:
	// _:c14n0 <ex:contributesTo> _:c14n1 .
	//
	// 1:
	// _:c14n0 <ex:contributesTo> _:c14n1 .
	// _:c14n0 <ex:contributesTo> _:c14n2 .
	// _:c14n2 <ex:dependsOn> _:c14n1 .
	//
	// 2:
	// _:c14n0 <ex:contributesTo> _:c14n1 .
	// _:c14n0 <ex:contributesTo> _:c14n3 .
	// _:c14n2 <ex:contributesTo> _:c14n1 .
	// _:c14n2 <ex:notContributesTo> _:c14n3 .
	// _:c14n3 <ex:dependsOn> _:c14n1 .
	//
	// 3:
	// _:c14n0 <ex:contributesTo> _:c14n1 .
	// _:c14n0 <ex:contributesTo> _:c14n3 .
	// _:c14n2 <ex:contributesTo> _:c14n1 .
	// _:c14n2 <ex:is> "Alice" .
	// _:c14n3 <ex:dependsOn> _:c14n1 .
	//
	// 4:
	// _:c14n0 <ex:contributesTo> _:c14n1 .
	// _:c14n0 <ex:contributesTo> _:c14n2 .
	// _:c14n2 <ex:dependsOn> _:c14n1 .
	// _:c14n3 <ex:contributesTo> _:c14n1 .
	// _:c14n3 <ex:notContributesTo> _:c14n2 .
	// _:c14n4 <ex:contributesTo> _:c14n2 .
	// _:c14n4 <ex:notContributesTo> _:c14n1 .
	//
	// 5:
	// _:c14n1 <ex:contributesTo> _:c14n0 .
	// _:c14n1 <ex:contributesTo> _:c14n3 .
	// _:c14n1 <ex:is> "Bob" .
	// _:c14n2 <ex:contributesTo> _:c14n0 .
	// _:c14n2 <ex:is> "Alice" .
	// _:c14n3 <ex:dependsOn> _:c14n0 .
	//
	// 6:
	// _:c14n0 <ex:dependsOn> _:c14n1 .
	// _:c14n2 <ex:contributesTo> _:c14n0 .
	// _:c14n2 <ex:contributesTo> _:c14n1 .
	// _:c14n2 <ex:is> "Bob" .
	// _:c14n3 <ex:contributesTo> _:c14n1 .
	// _:c14n3 <ex:is> "Alice" .
	// _:c14n4 <ex:contributesTo> _:c14n0 .
	// _:c14n4 <ex:is> "Charlie" .
}
