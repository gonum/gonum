// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"container/heap"

	"github.com/gonum/graph"
)

// Heuristic returns an estimate of the cost of travelling between two nodes.
type Heuristic func(x, y graph.Node) float64

// HeuristicCoster wraps the HeuristicCost method. A graph that implementing
// the interface provides a heuristic between any two given nodes.
type HeuristicCoster interface {
	HeuristicCost(x, y graph.Node) float64
}

// Returns an ordered list consisting of the nodes between start and goal. The path will be the
// shortest path assuming the function heuristicCost is admissible. The second return value is the
// cost, and the third is the number of nodes expanded while searching (useful info for tuning
// heuristics). Negative Costs will cause bad things to happen, as well as negative heuristic
// estimates.
//
// A heuristic is admissible if, for any node in the graph, the heuristic estimate of the cost
// between the node and the goal is less than or set to the true cost.
//
// Performance may be improved by providing a consistent heuristic (though one is not needed to
// find the optimal path), a heuristic is consistent if its value for a given node is less than
// (or equal to) the actual cost of reaching its neighbors + the heuristic estimate for the
// neighbor itself. You can force consistency by making your HeuristicCost function return
// max(NonConsistentHeuristicCost(neighbor,goal), NonConsistentHeuristicCost(self,goal) -
// Cost(self,neighbor)). If there are multiple neighbors, take the max of all of them.
//
// If heuristice is nil, A* will use the graph's HeuristicCost method if present, otherwise
// falling back to NullHeuristic. To run Uniform Cost Search, run A* with the NullHeuristic.
//
// To run Breadth First Search, run A* with both the NullHeuristic and UniformCost (or any cost
// function that returns a uniform positive value.)
func AStar(start, goal graph.Node, g graph.Graph, h Heuristic) (path []graph.Node, pathCost float64, nodesExpanded int) {
	var weight graph.WeightFunc
	if g, ok := g.(graph.Weighter); ok {
		weight = g.Weight
	} else {
		weight = graph.UniformCost
	}
	if h == nil {
		if g, ok := g.(HeuristicCoster); ok {
			h = g.HeuristicCost
		} else {
			h = NullHeuristic
		}
	}

	closedSet := make(map[int]aStarNode)
	openSet := &aStarQueue{nodes: make([]aStarNode, 0), indexList: make(map[int]int)}
	heap.Init(openSet)
	node := aStarNode{start, 0, h(start, goal)}
	heap.Push(openSet, node)
	predecessor := make(map[int]graph.Node)

	for openSet.Len() != 0 {
		curr := heap.Pop(openSet).(aStarNode)

		nodesExpanded += 1

		if curr.ID() == goal.ID() {
			return rebuildPath(predecessor, goal), curr.gscore, nodesExpanded
		}

		closedSet[curr.ID()] = curr

		for _, neighbor := range g.From(curr.Node) {
			if _, ok := closedSet[neighbor.ID()]; ok {
				continue
			}

			g := curr.gscore + weight(g.Edge(curr.Node, neighbor))

			if existing, exists := openSet.Find(neighbor.ID()); !exists {
				predecessor[neighbor.ID()] = curr
				node = aStarNode{neighbor, g, g + h(neighbor, goal)}
				heap.Push(openSet, node)
			} else if g < existing.gscore {
				predecessor[neighbor.ID()] = curr
				openSet.Fix(neighbor.ID(), g, g+h(neighbor, goal))
			}
		}
	}

	return nil, 0, nodesExpanded
}

// NullHeuristic is an admissible, consistent heuristic that will not speed up computation.
func NullHeuristic(_, _ graph.Node) float64 {
	return 0
}

type aStarNode struct {
	graph.Node
	gscore, fscore float64
}

type aStarQueue struct {
	indexList map[int]int
	nodes     []aStarNode
}

func (pq *aStarQueue) Less(i, j int) bool {
	// As the heap documentation says, a priority queue is listed if the actual values
	// are treated as if they were negative
	return pq.nodes[i].fscore < pq.nodes[j].fscore
}

func (pq *aStarQueue) Swap(i, j int) {
	pq.indexList[pq.nodes[i].ID()] = j
	pq.indexList[pq.nodes[j].ID()] = i

	pq.nodes[i], pq.nodes[j] = pq.nodes[j], pq.nodes[i]
}

func (pq *aStarQueue) Len() int {
	return len(pq.nodes)
}

func (pq *aStarQueue) Push(x interface{}) {
	node := x.(aStarNode)
	pq.nodes = append(pq.nodes, node)
	pq.indexList[node.ID()] = len(pq.nodes) - 1
}

func (pq *aStarQueue) Pop() interface{} {
	x := pq.nodes[len(pq.nodes)-1]
	pq.nodes = pq.nodes[:len(pq.nodes)-1]
	delete(pq.indexList, x.ID())

	return x
}

func (pq *aStarQueue) Fix(id int, newGScore, newFScore float64) {
	if i, ok := pq.indexList[id]; ok {
		pq.nodes[i].gscore = newGScore
		pq.nodes[i].fscore = newFScore
		heap.Fix(pq, i)
	}
}

func (pq *aStarQueue) Find(id int) (aStarNode, bool) {
	loc, ok := pq.indexList[id]
	if ok {
		return pq.nodes[loc], true
	} else {
		return aStarNode{}, false
	}

}

func (pq *aStarQueue) Exists(id int) bool {
	_, ok := pq.indexList[id]
	return ok
}

// Rebuilds a path backwards from the goal.
func rebuildPath(predecessors map[int]graph.Node, goal graph.Node) []graph.Node {
	if n, ok := goal.(aStarNode); ok {
		goal = n.Node
	}
	path := []graph.Node{goal}
	curr := goal
	for prev, ok := predecessors[curr.ID()]; ok; prev, ok = predecessors[curr.ID()] {
		if n, ok := prev.(aStarNode); ok {
			prev = n.Node
		}
		path = append(path, prev)
		curr = prev
	}

	reverse(path)
	return path
}
