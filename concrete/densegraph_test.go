package concrete_test

import (
	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

var _ graph.DirectedGraph = &concrete.DenseGraph{}
var _ graph.CrunchGraph = &concrete.DenseGraph{}
