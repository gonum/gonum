// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file

package cscalar_test

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"gonum.org/v1/gonum/cmplxs"
	"gonum.org/v1/gonum/cmplxs/cscalar"
	"gonum.org/v1/gonum/floats"
)

func ExampleParseWithNA() {
	// Calculate the mean of a list of numbers
	// ignoring missing values.
	const data = `6+2i
missing
4-4i
`

	var (
		vals    []complex128
		weights []float64
	)
	sc := bufio.NewScanner(strings.NewReader(data))
	for sc.Scan() {
		v, w, err := cscalar.ParseWithNA(sc.Text(), "missing")
		if err != nil {
			log.Fatal(err)
		}
		vals = append(vals, v)
		weights = append(weights, w)
	}
	err := sc.Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cmplxs.Sum(vals) / complex(floats.Sum(weights), 0))

	// Output:
	// (5-1i)
}
