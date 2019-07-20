// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.12

// TODO(kortschak): Merge this into nodes_test.go when go1.11 is no longer supported.

package iterator_test

import (
	"reflect"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/simple"
)

var nodesTests = []struct {
	nodes map[int64]graph.Node
}{
	{nodes: nil},
	{nodes: make(map[int64]graph.Node)},
	{nodes: map[int64]graph.Node{1: simple.Node(1)}},
	{nodes: map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 5: simple.Node(5)}},
	{nodes: map[int64]graph.Node{5: simple.Node(5), 3: simple.Node(3), 2: simple.Node(2), 1: simple.Node(1)}},
}

func TestNodesIterate(t *testing.T) {
	for _, test := range nodesTests {
		it := iterator.NewNodes(test.nodes)
		for i := 0; i < 2; i++ {
			if it.Len() != len(test.nodes) {
				t.Errorf("unexpected iterator length for round %d: got:%d want:%d", i, it.Len(), len(test.nodes))
			}
			var got map[int64]graph.Node
			if test.nodes != nil {
				got = make(map[int64]graph.Node)
			}
			for it.Next() {
				n := it.Node()
				got[n.ID()] = n
				if len(got)+it.Len() != len(test.nodes) {
					t.Errorf("unexpected iterator length during iteration for round %d: got:%d want:%d", i, it.Len(), len(test.nodes))
				}
			}
			want := test.nodes
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, want)
			}
			func() {
				defer func() {
					r := recover()
					if r != nil {
						t.Errorf("unexpected panic: %v", r)
					}
				}()
				it.Next()
			}()
			it.Reset()
		}
	}
}
