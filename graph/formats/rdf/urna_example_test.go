// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf_test

import (
	"fmt"
	"log"
	"strings"

	"gonum.org/v1/gonum/graph/formats/rdf"
)

func ExampleURDNA2015() {
	for _, statements := range []string{
		`
_:a <ex:q> <ex:p> .
_:b <ex:q> <ex:p> .
_:c <ex:p> _:a .
_:d <ex:p> _:b .
_:c <ex:r> _:d .
`,
		`
_:c1 <ex:p> _:a1 .
_:b1 <ex:q> <ex:p> .
_:d1 <ex:p> _:b1 .
_:a1 <ex:q> <ex:p> .
_:c1 <ex:r> _:d1 .
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

		relabeled, err := rdf.URDNA2015(nil, s)
		if err != nil {
			log.Fatal(err)
		}
		for _, s := range relabeled {
			fmt.Println(s)
		}
		fmt.Println()
	}

	// Output:
	//
	// _:c14n0 <ex:p> _:c14n2 .
	// _:c14n1 <ex:p> _:c14n3 .
	// _:c14n1 <ex:r> _:c14n0 .
	// _:c14n2 <ex:q> <ex:p> .
	// _:c14n3 <ex:q> <ex:p> .
	//
	// _:c14n0 <ex:p> _:c14n2 .
	// _:c14n1 <ex:p> _:c14n3 .
	// _:c14n1 <ex:r> _:c14n0 .
	// _:c14n2 <ex:q> <ex:p> .
	// _:c14n3 <ex:q> <ex:p> .
}

func ExampleURGNA2012() {
	for _, statements := range []string{
		`
_:a <ex:q> <ex:p> .
_:b <ex:q> <ex:p> .
_:c <ex:p> _:a .
_:d <ex:p> _:b .
_:c <ex:r> _:d .
`,
		`
_:c1 <ex:p> _:a1 .
_:b1 <ex:q> <ex:p> .
_:d1 <ex:p> _:b1 .
_:a1 <ex:q> <ex:p> .
_:c1 <ex:r> _:d1 .
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

		relabeled, err := rdf.URGNA2012(nil, s)
		if err != nil {
			log.Fatal(err)
		}
		for _, s := range relabeled {
			fmt.Println(s)
		}
		fmt.Println()
	}

	// Output:
	//
	// _:c14n0 <ex:p> _:c14n3 .
	// _:c14n0 <ex:r> _:c14n1 .
	// _:c14n1 <ex:p> _:c14n2 .
	// _:c14n2 <ex:q> <ex:p> .
	// _:c14n3 <ex:q> <ex:p> .
	//
	// _:c14n0 <ex:p> _:c14n3 .
	// _:c14n0 <ex:r> _:c14n1 .
	// _:c14n1 <ex:p> _:c14n2 .
	// _:c14n2 <ex:q> <ex:p> .
	// _:c14n3 <ex:q> <ex:p> .
}
