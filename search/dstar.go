package search

import (
	"container/heap"
	"errors"
	"math"

	gr "github.com/gonum/graph"
)

// A DStarInstance is a way for a general Graph implementer to run D*-Lite. D*-Lite has state over multiple calls, and allows the graph to chnge, so after initialization an instance is returned.
//
// The general flow of D*-lite is that it's initialized and an initial shortest path is computed (essentially the same as A*). This is all done in InitDStar.
// This is followed up by a Step() which returns the node to move to. After this action is taken, the state of the algorithm is Update()d if any edge costs have changed. In general, running DStar
// directly will look like:
//
//     myDStar := InitDStar(... info ...)
//
//     * Repeat the following until Step returns error or you reach your goal *
//         move := myDStar.Step()
//         perform the move returned
//         if the graph changed, call Update(... change info ...)
type DStarInstance struct {
	graph             gr.Graph
	start, goal, last gr.Node
	gScores           map[int]float64
	cost              func(gr.Node, gr.Node) float64
	heuristicCost     func(gr.Node, gr.Node) float64
	successors        func(gr.Node) []gr.Node
	predecessors      func(gr.Node) []gr.Node
	u                 *dStarPriorityQueue
	rhs               map[int]float64
	k_m               float64
}

func (ds *DStarInstance) calculateKey(node gr.Node) key {
	rhs := ds.rhs[node.ID()]
	gScore := ds.gScores[node.ID()]
	return key{math.Min(gScore, rhs) + ds.heuristicCost(ds.start, node) + ds.k_m, math.Min(gScore, rhs)}
}

// Initialized an instance of D*-Lite for running on a graph. Note that this does not match directly with Initialize() in the original D*-Lite paper.
// Instead, it is the lines:
//
//     s_last = s_start
//     Initialize()
//     ComputeShortestPath()
//
// In other words, it's all the lines before the main loop in Main() in the original paper. Essentially a full state initialization.
func InitDStar(start, goal gr.Node, graph gr.Graph, Cost, HeuristicCost func(gr.Node, gr.Node) float64) *DStarInstance {
	successors, predecessors, _, _, _, _, Cost, HeuristicCost := setupFuncs(graph, Cost, HeuristicCost)

	u := &dStarPriorityQueue{indexList: make(map[int]int, 0), nodes: make([]dStarNode, 0)}
	heap.Init(u)

	ds := &DStarInstance{
		graph:         graph,
		start:         start,
		goal:          goal,
		last:          start,
		u:             u,
		k_m:           0.0,
		gScores:       make(map[int]float64, 0),
		rhs:           make(map[int]float64, 0),
		cost:          Cost,
		heuristicCost: HeuristicCost,
		successors:    successors,
		predecessors:  predecessors,
	}

	for _, node := range graph.NodeList() {
		ds.rhs[node.ID()] = math.Inf(1)
		ds.gScores[node.ID()] = math.Inf(1)
	}

	ds.rhs[goal.ID()] = 0.0
	heap.Push(ds.u, dStarNode{Node: goal, key: ds.calculateKey(goal)})
	ds.computeShortestPath()
	return ds
}

func (ds *DStarInstance) updateVertex(node gr.Node) {
	if node.ID() != ds.goal.ID() {
		min := math.Inf(1)
		for _, succ := range ds.successors(node) {
			min = math.Min(min, ds.cost(node, succ)+ds.gScores[succ.ID()])
		}
		ds.rhs[node.ID()] = min
	}

	if math.Abs(ds.gScores[node.ID()]-ds.rhs[node.ID()]) > .000001 {
		ds.u.Fix(node, ds.calculateKey(node))
	} else {
		ds.u.Remove(node)
	}
}

func (ds *DStarInstance) computeShortestPath() {
	for ds.u.Len() > 0 && (ds.u.Peek().Less(dStarNode{Node: ds.start, key: ds.calculateKey(ds.start)}) || math.Abs(ds.rhs[ds.start.ID()]-ds.gScores[ds.start.ID()]) > .000001) {

		vert := heap.Pop(ds.u).(dStarNode)
		newKey := ds.calculateKey(vert.Node)
		if vert.Less(dStarNode{Node: vert.Node, key: newKey}) {

			heap.Push(ds.u, dStarNode{Node: vert.Node, key: newKey})

		} else if ds.gScores[vert.ID()] > ds.rhs[vert.ID()] {

			ds.gScores[vert.ID()] = ds.rhs[vert.ID()]
			for _, pred := range ds.predecessors(vert.Node) {
				ds.updateVertex(pred)
			}

		} else {

			ds.gScores[vert.ID()] = math.Inf(1)
			ds.updateVertex(vert.Node)
			for _, pred := range ds.predecessors(vert.Node) {
				ds.updateVertex(pred)
			}

		}
	}
}

// Returns the next action to be taken, or nil and an error if it's determined that no path exists.
// Should be called before Update every loop
func (ds *DStarInstance) Step() (succ gr.Node, err error) {
	if ds.start.ID() == ds.goal.ID() {
		return ds.start, nil
	} else if ds.gScores[ds.start.ID()] == math.Inf(1) {
		return nil, errors.New("No path exists")
	}

	min := math.Inf(1)
	var next gr.Node
	for _, succ := range ds.successors(ds.start) {
		newMin := math.Min(min, ds.cost(ds.start, succ)+ds.gScores[succ.ID()])
		if newMin < min {
			min = newMin
			next = succ
		}
	}

	return next, nil
}

// Updates D*-Lite if new information has been discovered or the graph has changed in any way. Should be called after each call of Step()
// This is a no-op if changedEdgeCosts is nil or its len is 0.
func (ds *DStarInstance) Update(cost func(gr.Node, gr.Node) float64, changedEdgeCosts []gr.Edge) {
	if changedEdgeCosts == nil || len(changedEdgeCosts) == 0 {
		return
	}

	if cost != nil {
		ds.cost = cost
	}
	ds.k_m += ds.heuristicCost(ds.last, ds.start)
	ds.last = ds.start

	for _, edge := range changedEdgeCosts {
		ds.updateVertex(edge.Head())
	}
	ds.computeShortestPath()
}

// Runs D*-Lite in its entirety on an appropriate graph. What is D*-Lite? It's an incremental heuristic lifelong planning search. What this means is that
// unlike A*, it reacts to new information and can be used when the graph representation changes to replan paths. Perhaps the best way to understand D*-Lite
// is to read the paper it came from[1]. However, here is a (very) brief synopsis:
//
// D*-Lite begins by computing the shortest path between the start and goal, and assuming this is the path it will take. This step is essentially just A*.
// After moving a step, it will "take a look around" so to speak. It will scan and see if anything in the graph has changed, if it had changed, it recomputes the shortest path
// and uses the new one. The big trick in D*-lite is that since it has memory, this is significantly cheaper than rerunning A* at every step (though there are edge cases where performance may suffer).
//
// D*-Lite is a modification of LPA* (Lifelong Planning A*) that has the same behavior as a little-used algorithm known as D*. It is guaranteed to run as fast or faster than D*, and is generally
// known to run faster than LPA* and other similar lifelong planning algorithms.
//
// It is used often in real world robotics path planning -- a modification of this algorithm known as Field D* (which allows more degrees of freedom in movement) is used in the Mars rovers Spirit and Opportunity.
// It was also notably used for traffic path planning in the recent reboot of the SimCity franchise.
//
// As with other algorithms with cost function arguments in this package, Cost and HeuristicCost are optional, and if absent will default to the graph's Cost/HeuristicCost functions (if present), and finally
// to UniformCost and NullHeauristic respectively.
//
// [1] http://www.aaai.org/Papers/AAAI/2002/AAAI02-072.pdf
func DStarLite(start, goal gr.Node, graph gr.DStarGraph, Cost, HeuristicCost func(gr.Node, gr.Node) float64) error {
	ds := InitDStar(start, goal, graph, Cost, HeuristicCost) // InitDStar does s_last = s_start and computeShortestPath for us
	for ds.start.ID() != ds.goal.ID() {
		next, err := ds.Step()
		if err != nil {
			return err
		}

		graph.Move(next)
		newCost, edges := graph.ChangedEdges()
		ds.Update(newCost, edges)
	}

	return nil
}

// Starts a D*-Lite service in a seperate goroutine. When signalled on the step channel, it will perform a single step/move/update cycle.
// If the step channel is closed before the algorithm is done running, this goroutine will abort.
//
// If an error is encountered, it will be sent over the done channel. If D*-lite exits successfully the done channel will be closed with no error written to it.
//
// D*-lite is initialized upon call, albeit in the new goroutine. However, the first step/move/update cycle is not performed until a signal is received.
func SynchronizedDStarLite(start, goal gr.Node, graph gr.DStarGraph, Cost, HeuristicCost func(gr.Node, gr.Node) float64, step <-chan struct{}, done chan<- error) {
	go func() {
		ds := InitDStar(start, goal, graph, Cost, HeuristicCost) // InitDStar does s_last = s_start and computeShortestPath for us
		for ds.start.ID() != ds.goal.ID() {
			_, ok := <-step
			if !ok {
				close(done)
				return
			}
			next, err := ds.Step()
			if err != nil {
				done <- err
				close(done)
				return
			}

			graph.Move(next)
			newCost, edges := graph.ChangedEdges()
			ds.Update(newCost, edges)
		}

		close(done)
	}()
}
