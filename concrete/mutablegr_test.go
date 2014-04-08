// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package concrete_test

import (
	"testing"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

var _ graph.Graph = &concrete.Graph{}
var _ graph.Graph = &concrete.Graph{}

func TestAssertMutableNotDirected(t *testing.T) {
	var g graph.MutableGraph = concrete.NewGraph()
	if _, ok := g.(graph.DirectedGraph); ok {
		t.Fatal("concrete.Graph is directed, but a MutableGraph cannot safely be directed!")
	}
}

// var _ gr.EdgeListGraph = &concrete.Graph{}
