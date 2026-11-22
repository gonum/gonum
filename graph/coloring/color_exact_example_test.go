// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coloring_test

import (
	"context"
	"errors"
	"fmt"
	"log"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/coloring"
	"gonum.org/v1/gonum/graph/encoding/graph6"
)

func Example_color_exact() {
	// Queen 6-6 graph.
	g := graph6.Graph("c~~}FDrMw~`~goSwtMYhvIF{SN{dEAQfCehrcTMyPO~ca`~acgoPQSwcCtMWcahvaQIF|CcSN{KSdEAAIQfC__ehrCCcTMwSQPO~ogca`~")
	if !graph6.IsValid(g) {
		log.Fatal("invalid g6 graph")
	}

	k, _, err := nDsaturExact(context.Background(), g, 2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(k)

	// Output:
	// 7
}

// nDsaturExact starts n DsaturExact routines and returns the results
// from the first one that completes, terminating the others. This can
// be used to improve performance for difficult graphs at the cost of
// increased memory use and discarded computation.
func nDsaturExact(ctx context.Context, g graph.Undirected, n int) (k int, colors map[int64]int, err error) {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(context.Canceled)

	type result struct {
		k      int
		colors map[int64]int
		err    error
	}

	r := make(chan result, n)
	done := errors.New("done")
	for range n {
		go func() {
			k, colors, err := coloring.DsaturExact(ctx, g)
			r <- result{k, colors, err}
			if err == nil {
				cancel(done)
			}
		}()
	}

	for range n {
		coloring := <-r
		if coloring.err == nil {
			return coloring.k, coloring.colors, nil
		}

		if context.Cause(ctx) != done {
			return coloring.k, coloring.colors, coloring.err
		}
	}

	panic("unreachable")
}
