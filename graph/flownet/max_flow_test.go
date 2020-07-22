package flownet

import (
	"testing"

	"gonum.org/v1/gonum/graph/simple"
)

func makeGraph() *SimpleFlowNetwork {
	// Create graph:
	//
	//    0
	//  / | \
	// 3  |  1
	// \ | /
	//  2
	//
	g := simple.NewWeightedDirectedGraph(0, -1)
	g.AddNode(simple.Node(0))
	g.AddNode(simple.Node(1))
	g.AddNode(simple.Node(2))
	g.AddNode(simple.Node(3))
	g.SetWeightedEdge(simple.WeightedEdge{
		F: simple.Node(0), T: simple.Node(1),
		W: 5.0,
	})
	g.SetWeightedEdge(simple.WeightedEdge{
		F: simple.Node(0), T: simple.Node(2),
		W: 15.0,
	})
	g.SetWeightedEdge(simple.WeightedEdge{
		F: simple.Node(2), T: simple.Node(1),
		W: 10.0,
	})
	g.SetWeightedEdge(simple.WeightedEdge{
		F: simple.Node(3), T: simple.Node(2),
		W: 5.0,
	})
	g.SetWeightedEdge(simple.WeightedEdge{
		F: simple.Node(3), T: simple.Node(0),
		W: 10.0,
	})
	return &SimpleFlowNetwork{*g}
}

func TestEdgesPresenceSimpleGraph(t *testing.T) {
	var sg FordFulkersonGraph = makeGraph()
	for _, edge := range [][]int64{{0, 1}, {0, 2}, {2, 1}, {3, 2}, {3, 0}} {
		e := sg.Edge(edge[0], edge[1])
		if e == nil {
			t.Errorf("Edge %d %d is missing.", edge[0], edge[1])
		}
	}
}

func TestEdgesPresenceResidualGraph(t *testing.T) {
	var g ResidualGraph = makeGraph().ResidualGraph()
	for _, edge := range [][]int64{{0, 1, 5}, {0, 2, 15}, {2, 1, 10}, {3, 2, 5}, {3, 0, 10}} {
		e := g.ResidualEdge(edge[0], edge[1])
		capacity := int32(edge[2])
		if e == nil {
			t.Errorf("Edge %d %d is missing.", edge[0], edge[1])
		}
		if e.MaxCapacity() != capacity {
			t.Errorf("Edge %d %d capacity is %d. Must be %d", edge[0], edge[1], e.MaxCapacity(), capacity)
		}

		if e.CurrentFlow() != 0 {
			t.Errorf("Edge %d %d flow is %d. Must be 0", edge[0], edge[1], e.CurrentFlow())
		}
	}
}

func TestEdgesFromResidualGraph(t *testing.T) {
	var g ResidualGraph = makeGraph().ResidualGraph()

	for frm, to := range map[int64]map[int64]int32{
		0: {1: 5, 2: 15},
		2: {1: 10},
		3: {2: 5, 0: 10},
	} {
		fromNodes := g.From(frm)
		for fromNodes.Next() {
			toNode := fromNodes.Node().ID()

			edge := g.ResidualEdge(frm, toNode)
			expectedCapacity := to[toNode]
			if edge.MaxCapacity() != expectedCapacity {
				t.Errorf("Edge %d %d capacity is %d. Must be %d", frm, toNode, edge.MaxCapacity(), expectedCapacity)
			}
		}
	}

}

func TestSetFlowResidualGraph(t *testing.T) {
	var g ResidualGraph = makeGraph().ResidualGraph()

	g.SetFlow(3, 0, 7)

	if g.ResidualEdge(3, 0).CurrentFlow() != 7 {
		t.Errorf("Edge %d %d flow is %d. Must be %d", 3, 0, g.ResidualEdge(3, 0).CurrentFlow(), 7)
	}

	if g.ResidualEdge(0, 3).CurrentFlow() != 3 {
		t.Errorf("Edge %d %d flow is %d. Must be %d", 0, 3, g.ResidualEdge(0, 3).CurrentFlow(), 3)
	}

}

func TestAugmentingPath(t *testing.T) {
	var g ResidualGraph = makeGraph().ResidualGraph()
	source := simple.Node(3)
	sink := simple.Node(1)
	path := AugmentingPath(g, source, sink)
	flow := AvailableCapacity(path)
	if flow <= 0 {
		t.Errorf("Flow is expected to be > 0")
	}
	for _, edge := range path {
		oldFlow := g.ResidualEdge(edge.From().ID(), edge.To().ID()).CurrentFlow()
		g.SetFlow(edge.From().ID(), edge.To().ID(), flow)
		if oldFlow == g.ResidualEdge(edge.From().ID(), edge.To().ID()).CurrentFlow() {
			t.Errorf("Flow in the edge %d %d hasn't changed", edge.From().ID(), edge.From().ID())
		}
	}

	path2 := AugmentingPath(g, source, sink)
	for _, edge := range path2 {
		if edge.CurrentFlow() == edge.MaxCapacity() {
			t.Errorf("Flow in the edge %d %d is equal to its max capacity.", edge.From().ID(), edge.From().ID())
		}
	}
}

func TestMaxFlow(t *testing.T) {
	ek := EdmondsKarp{
		makeGraph(),
	}
	flow := ek.MaxFlow(
		simple.Node(3),
		simple.Node(1),
	)
	if flow != 15 {
		t.Errorf("Max flow = %d. Must be 15", flow)
	}

}
