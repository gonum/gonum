// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf_test

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

	"gonum.org/v1/gonum/graph/formats/rdf"
)

func ExampleIsoCanonicalHashes_isomorphisms() {
	for _, statements := range []string{
		`
<https://example.com/1> <https://example.com/2> <https://example.com/3> .
<https://example.com/3> <https://example.com/4> <https://example.com/5> .
`,
		`
_:a <ex:q> <ex:p> .
_:b <ex:q> <ex:p> .
_:c <ex:p> _:a .
_:d <ex:p> _:b .
_:c <ex:r> _:d .
`,
		`
_:a1 <ex:q> <ex:p> .
_:b1 <ex:q> <ex:p> .
_:c1 <ex:p> _:a1 .
_:d1 <ex:p> _:b1 .
_:c1 <ex:r> _:d1 .
`,
		`
# G
<ex:p> <ex:q> _:a .
<ex:p> <ex:q> _:b .
<ex:s> <ex:p> _:a .
<ex:s> <ex:r> _:c .
_:c <ex:p> _:b .

# H
<ex:p> <ex:q> _:d .
<ex:p> <ex:q> _:e .
_:f <ex:p> _:d .
_:f <ex:r> _:g .
_:g <ex:p> _:e .
`,
		`
_:greet <l:is> "hola"@es .
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

		// Get the node label to hash look-up table.
		hashes, _ := rdf.IsoCanonicalHashes(s, false, true, md5.New(), make([]byte, 16))

		// Get all the blank nodes.
		var blanks []string
		for k := range hashes {
			if strings.HasPrefix(k, "_:") {
				blanks = append(blanks, k)
			}
		}
		slices.Sort(blanks)

		if len(blanks) == 0 {
			fmt.Println("No blank nodes.")
		} else {
			w := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', 0)
			fmt.Fprintln(w, "Node\tHash")
			for _, k := range blanks {
				fmt.Fprintf(w, "%s\t%032x\n", k, hashes[k])
			}
			w.Flush()
		}
		fmt.Println()
	}

	// Output:
	//
	// No blank nodes.
	//
	// Node Hash
	// _:a  d4db6df055d5611e9d8aa6ea621561d1
	// _:b  ad70e47f2b026064c7f0922060512b9a
	// _:c  dafd81e6fa603d3e11c898d631e8673f
	// _:d  7e318557b09444e88791721becc2a8e7
	//
	// Node Hash
	// _:a1 d4db6df055d5611e9d8aa6ea621561d1
	// _:b1 ad70e47f2b026064c7f0922060512b9a
	// _:c1 dafd81e6fa603d3e11c898d631e8673f
	// _:d1 7e318557b09444e88791721becc2a8e7
	//
	// Node Hash
	// _:a  44ad49b6df3aea91ddbcef932c93e3b4
	// _:b  ba3ffd8b271a8545b1a3a9042e75ce4b
	// _:c  34e1bd90b6758b4a766e000128caa6a6
	// _:d  eb2a47c1032f623647d0497a2ff74052
	// _:e  1d9ed02f28d87e555feb904688bc2449
	// _:f  d3b00d36ea503dcc8d234e4405feab81
	// _:g  55127e4624c0a4fe5990933a48840af8
	//
	// Node    Hash
	// _:greet 0d9ba18a037a3fa67e46fce821fe51b4
}

func ExampleIsoCanonicalHashes_isomorphicParts() {
	for _, statements := range []string{
		`
# Part 1
_:a <ex:q> <ex:p> .
_:b <ex:q> <ex:p> .
_:c <ex:p> _:a .
_:d <ex:p> _:b .
_:c <ex:r> _:d .

# Part 2
_:a1 <ex:q> <ex:p> .
_:b1 <ex:q> <ex:p> .
_:c1 <ex:p> _:a1 .
_:d1 <ex:p> _:b1 .
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

		// Get the node label to hash look-up table. This time
		// we will decompose the dataset into splits and not
		// distinguish nodes. This will then group nodes from
		// the two isomorphic parts. Otherwise each node in
		// the complete dataset would get a unique hash.
		hashes, _ := rdf.IsoCanonicalHashes(s, true, false, md5.New(), make([]byte, 16))

		// Get all the blank nodes.
		var blanks []string
		for k := range hashes {
			if strings.HasPrefix(k, "_:") {
				blanks = append(blanks, k)
			}
		}
		slices.Sort(blanks)

		if len(blanks) == 0 {
			fmt.Println("No blank nodes.")
		} else {
			w := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', 0)
			fmt.Fprintln(w, "Node\tHash")
			for _, k := range blanks {
				fmt.Fprintf(w, "%s\t%032x\n", k, hashes[k])
			}
			w.Flush()
		}
		fmt.Println()
	}

	// Output:
	//
	// Node Hash
	// _:a  d4db6df055d5611e9d8aa6ea621561d1
	// _:a1 d4db6df055d5611e9d8aa6ea621561d1
	// _:b  ad70e47f2b026064c7f0922060512b9a
	// _:b1 ad70e47f2b026064c7f0922060512b9a
	// _:c  dafd81e6fa603d3e11c898d631e8673f
	// _:c1 dafd81e6fa603d3e11c898d631e8673f
	// _:d  7e318557b09444e88791721becc2a8e7
	// _:d1 7e318557b09444e88791721becc2a8e7
}

func ExampleC14n() {
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

		// Get the hash to term label look-up table.
		_, terms := rdf.IsoCanonicalHashes(s, false, true, md5.New(), make([]byte, 16))

		relabeled, err := rdf.C14n(nil, s, terms)
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
	// _:c14n0 <ex:p> _:c14n1 .
	// _:c14n1 <ex:q> <ex:p> .
	// _:c14n2 <ex:q> <ex:p> .
	// _:c14n3 <ex:p> _:c14n2 .
	// _:c14n3 <ex:r> _:c14n0 .
	//
	// _:c14n0 <ex:p> _:c14n1 .
	// _:c14n1 <ex:q> <ex:p> .
	// _:c14n2 <ex:q> <ex:p> .
	// _:c14n3 <ex:p> _:c14n2 .
	// _:c14n3 <ex:r> _:c14n0 .
}
