// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testgraph

import (
	"reflect"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

var randomNodesTests = []struct {
	n    int
	seed uint64
	new  func(int64) graph.Node
	want []graph.Node
}{
	{
		n:    0,
		want: nil,
	},
	{
		n:    1,
		seed: 1,
		new:  newSimpleNode,
		want: []graph.Node{simple.Node(-1890700816702069259)},
	},
	{
		n:    1,
		seed: 2,
		new:  newSimpleNode,
		want: []graph.Node{simple.Node(-9080340245984136982)},
	},
	{
		n:    4,
		seed: 1,
		new:  newSimpleNode,
		want: []graph.Node{
			simple.Node(-1890700816702069259),
			simple.Node(7618499319381327068),
			simple.Node(-8006975346196781910),
			simple.Node(1952627761515405933),
		},
	},
	{
		n:    4,
		seed: 2,
		new:  newSimpleNode,
		want: []graph.Node{
			simple.Node(-9080340245984136982),
			simple.Node(8080303881793168789),
			simple.Node(6742726861847166348),
			simple.Node(-1570298006490451715),
		},
	},
}

func newSimpleNode(id int64) graph.Node { return simple.Node(id) }

func TestRandomNodesIterate(t *testing.T) {
	for _, test := range randomNodesTests {
		for i := 0; i < 2; i++ {
			it := NewRandomNodes(test.n, test.seed, test.new)
			if it.Len() != len(test.want) {
				t.Errorf("unexpected iterator length for round %d: got:%d want:%d", i, it.Len(), len(test.want))
			}
			var got []graph.Node
			for it.Next() {
				got = append(got, it.Node())
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, test.want)
			}
			it.Reset()
		}
	}
}
