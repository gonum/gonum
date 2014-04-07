// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package concrete_test

import (
	"testing"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

var _ graph.Graph = &concrete.MutableGraph{}
var _ graph.MutableGraph = &concrete.MutableGraph{}

func TestAssertMutableNotDirected(t *testing.T) {
	graph.VetMutableGraph(concrete.NewMutableGraph())
}

// var _ gr.EdgeListGraph = &concrete.Graph{}
