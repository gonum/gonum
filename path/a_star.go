// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"container/heap"

	"github.com/gonum/graph"
	"github.com/gonum/graph/internal"
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
func AStar(start, goal graph.Node, g graph.Graph, h Heuristic) (path []graph.Node, cost float64, expanded int) {
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

	p := newShortestFrom(start, g.Nodes())
	gid := goal.ID()

	visited := make(internal.IntSet)
	open := &aStarQueue{indexOf: make(map[int]int)}
	heap.Push(open, aStarNode{node: start, gscore: 0, fscore: h(start, goal)})

	for open.Len() != 0 {
		u := heap.Pop(open).(aStarNode)
		uid := u.node.ID()
		i := p.indexOf[uid]
		expanded++

		if uid == gid {
			break
		}

		visited.Add(uid)
		for _, v := range g.From(u.node) {
			vid := v.ID()
			if visited.Has(vid) {
				continue
			}
			j := p.indexOf[vid]

			g := u.gscore + weight(g.Edge(u.node, v))
			if n, ok := open.node(vid); !ok {
				p.set(j, g, i)
				heap.Push(open, aStarNode{node: v, gscore: g, fscore: g + h(v, goal)})
			} else if g < n.gscore {
				p.set(j, g, i)
				open.update(vid, g, g+h(v, goal))
			}
		}
	}

	path, cost = p.To(goal)
	return path, cost, expanded
}

// NullHeuristic is an admissible, consistent heuristic that will not speed up computation.
func NullHeuristic(_, _ graph.Node) float64 {
	return 0
}

type aStarNode struct {
	node   graph.Node
	gscore float64
	fscore float64
}

type aStarQueue struct {
	indexOf map[int]int
	nodes   []aStarNode
}

func (q *aStarQueue) Less(i, j int) bool {
	return q.nodes[i].fscore < q.nodes[j].fscore
}

func (q *aStarQueue) Swap(i, j int) {
	q.indexOf[q.nodes[i].node.ID()] = j
	q.indexOf[q.nodes[j].node.ID()] = i
	q.nodes[i], q.nodes[j] = q.nodes[j], q.nodes[i]
}

func (q *aStarQueue) Len() int {
	return len(q.nodes)
}

func (q *aStarQueue) Push(x interface{}) {
	n := x.(aStarNode)
	q.nodes = append(q.nodes, n)
	q.indexOf[n.node.ID()] = len(q.nodes) - 1
}

func (q *aStarQueue) Pop() interface{} {
	n := q.nodes[len(q.nodes)-1]
	q.nodes = q.nodes[:len(q.nodes)-1]
	delete(q.indexOf, n.node.ID())
	return n
}

func (q *aStarQueue) update(id int, g, f float64) {
	i, ok := q.indexOf[id]
	if !ok {
		return
	}
	q.nodes[i].gscore = g
	q.nodes[i].fscore = f
	heap.Fix(q, i)
}

func (q *aStarQueue) node(id int) (aStarNode, bool) {
	loc, ok := q.indexOf[id]
	if ok {
		return q.nodes[loc], true
	}
	return aStarNode{}, false
}
