package concrete_test

import (
	"testing"

	gr "github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

// exercise methods for 100% coverage.
func TestGonumGraphCoverage(t *testing.T) {
	var n1 gr.Node = concrete.GonumNode(1)
	var e gr.Edge = concrete.GonumEdge{n1, n1}
	e.Head()
	e.Tail()
	concrete.NewPreAllocatedGonumGraph(true, 1)
	dg := concrete.NewGonumGraph(true)                     // directed
	var dm gr.MutableGraph = dg                            // directed mutable
	var um gr.MutableGraph = concrete.NewGonumGraph(false) // undirected
	dm.AddNode(n1, nil)
	n0 := dm.NewNode(nil) // fill in hole, node 0
	dm.NewNode(nil)       // node 2
	dm.AddNode(n1, nil)   // no op
	um.AddNode(n1, nil)   // undirected
	n3 := concrete.GonumNode(3)
	n4 := concrete.GonumNode(4)
	dm.AddNode(n3, []gr.Node{n4}) // both node and successor are new
	um.AddNode(n0, []gr.Node{n1}) // new undirected edge
	dm.AddEdge(e)
	n5 := concrete.GonumNode(5)
	dm.AddEdge(concrete.GonumEdge{n5, n1}) // head not in graph
	dm.AddEdge(concrete.GonumEdge{n1, n5}) // tail not in graph
	um.AddEdge(concrete.GonumEdge{n1, n0}) // undirected
	n6 := concrete.GonumNode(6)
	dm.SetEdgeCost(concrete.GonumEdge{n6, n1}, 0) // n6 not in graph
	dm.SetEdgeCost(concrete.GonumEdge{n1, n6}, 0)
	um.SetEdgeCost(concrete.GonumEdge{n0, n1}, 0) // undirected
	dm.AddEdge(concrete.GonumEdge{n5, n0})
	dm.RemoveNode(n5)
	dm.RemoveNode(n5)
	dm.RemoveEdge(concrete.GonumEdge{n6, n0})
	dm.RemoveEdge(concrete.GonumEdge{n0, n6})
	um.RemoveEdge(concrete.GonumEdge{n0, n1})
	dm.SetDirected(false)
	dm.EmptyGraph()
	dm.EmptyGraph()
	dm.SetDirected(true)
	d := dm.(gr.DirectedGraph)
	d.Successors(n0)
	dm.AddNode(n0, []gr.Node{n1})
	d.Successors(n0)
	d.IsSuccessor(n3, n1)
	d.IsSuccessor(n1, n3)
	d.Predecessors(n3)
	d.Predecessors(n1)
	d.IsPredecessor(n3, n1)
	d.IsPredecessor(n1, n3)
	u := um.(gr.UndirectedGraph)
	u.Neighbors(n3)
	um.AddEdge(concrete.GonumEdge{n0, n1})
	u.Neighbors(n1)
	dg.Neighbors(n1) // directed graph neighbors
	d.IsNeighbor(n3, n1)
	d.IsNeighbor(n1, n3)
	d.NodeExists(n1)
	d.Degree(n3)
	um.AddEdge(concrete.GonumEdge{n0, n0})
	d.Degree(n0)
	u.Degree(n0)
	u.NodeList()
	u.EdgeList()
	dg.IsDirected()
	dm.Cost(n0, n1)
	dm.Cost(n6, n1)
}
