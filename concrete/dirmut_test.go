// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package concrete_test

import (
	_ "testing"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

var _ graph.Graph = &concrete.MutableDirectedGraph{}
var _ graph.DirectedGraph = &concrete.MutableDirectedGraph{}
var _ graph.MutableDirectedGraph = &concrete.MutableDirectedGraph{}
