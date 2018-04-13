package dot_test

import (
	"fmt"

	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/encoding/dot"
)

type edgeWithPorts struct {
	simple.Edge
	fromPort, toPort string
}

func (e *edgeWithPorts) FromPort() (string, string) {
	return e.fromPort, ""
}

func (e *edgeWithPorts) ToPort() (string, string) {
	return e.toPort, ""
}

func ExamplePorter() {
	g := simple.NewUndirectedGraph()
	g.SetEdge(&edgeWithPorts{
		Edge: simple.Edge{ simple.Node(1), simple.Node(0) },
		fromPort: "p1",
		toPort: "p2",
	})

	result, _ := dot.Marshal(g, "", "", "  ", true)
	fmt.Print(string(result))

	// Output:
	// strict graph {
	//   // Node definitions.
	//   0;
	//   1;
	//
	//   // Edge definitions.
	//   0:p2 -- 1:p1;
	// }
}

