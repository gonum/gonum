package concrete_test

import (
	_ "testing"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

var _ graph.Graph = &concrete.Graph{}
var _ graph.DirectedGraph = &concrete.Graph{}
var _ graph.MutableGraph = &concrete.Graph{}

// var _ gr.EdgeListGraph = &concrete.Graph{}
