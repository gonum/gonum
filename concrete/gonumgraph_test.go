package concrete_test

import (
	"testing"

	gr "github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

var _ gr.Graph = &concrete.GonumGraph{}
var _ gr.DirectedGraph = &concrete.GonumGraph{}
var _ gr.MutableGraph = &concrete.GonumGraph{}
var _ gr.EdgeListGraph = &concrete.GonumGraph{}
