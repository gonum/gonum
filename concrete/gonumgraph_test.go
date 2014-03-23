package concrete_test

import (
	_ "testing"

	gr "github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

var _ gr.Graph = &concrete.Graph{}
var _ gr.DirectedGraph = &concrete.Graph{}
var _ gr.MutableGraph = &concrete.Graph{}

// var _ gr.EdgeListGraph = &concrete.Graph{}
