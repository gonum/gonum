package concrete_test

import (
	"testing"

	gr "github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

// exercise methods for 100% coverage.
func TestGonumGraphCoverage(t *testing.T) {
	var n1 gr.Node = concrete.GonumNode(1)
	n1.ID()
	var e gr.Edge = concrete.GonumEdge{n1, n1}
	e.Head()
	e.Tail()
	concrete.NewPreAllocatedGonumGraph(true, 1)
	var dm gr.MutableGraph = concrete.NewGonumGraph(true)  // directed mutable
	var um gr.MutableGraph = concrete.NewGonumGraph(false) // undirected
	dm.AddNode(n1, nil)                                    // node 1
	n0 := dm.NewNode(nil)                                  // fill in hole
	dm.NewNode(nil)                                        // node 2
	dm.AddNode(n1, nil)                                    // no op
	um.AddNode(n1, nil)                                    // undirected
	// both node and successor are new
	dm.AddNode(concrete.GonumNode(3), []gr.Node{concrete.GonumNode(4)})
	um.AddNode(n0, []gr.Node{n1}) // new undirected edge
	dm.AddEdge(e)
	n5 := concrete.GonumNode(5)
	dm.AddEdge(concrete.GonumEdge{n5, n1}) // head not in graph
	dm.AddEdge(concrete.GonumEdge{n1, n5}) // tail not in graph
	um.AddEdge(concrete.GonumEdge{n1, n0}) // undirected
	dm.SetEdgeCost(e, 0)
	n6 := concrete.GonumNode(6)
	dm.SetEdgeCost(concrete.GonumEdge{n6, n1}, 0) // n6 not in graph
	dm.SetEdgeCost(concrete.GonumEdge{n1, n6}, 0)
	um.SetEdgeCost(concrete.GonumEdge{n0, n1}, 0) // undirected
	dm.AddEdge(concrete.GonumEdge{n5, n0})
	dm.RemoveNode(n5)
	dm.RemoveNode(n5)
	// 29.6% so far
}
