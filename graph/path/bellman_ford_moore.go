// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/linear"
	"gonum.org/v1/gonum/graph/traverse"
)

// BellmanFordFrom returns a shortest-path tree for a shortest path from u to all nodes in
// the graph g, or false indicating that a negative cycle exists in the graph. If the graph
// does not implement Weighted, UniformCost is used.
//
// If g is a graph.Graph, all nodes of the graph will be stored in the shortest-path
// tree, otherwise only nodes reachable from u will be stored.
//
// The time complexity of BellmanFordFrom is O(|V|.|E|).
func BellmanFordFrom(u graph.Node, g traverse.Graph) (path Shortest, ok bool) {
	if h, ok := g.(graph.Graph); ok {
		if h.Node(u.ID()) == nil {
			return Shortest{from: u}, true
		}
		path = newShortestFrom(u, graph.NodesOf(h.Nodes()))
	} else {
		if g.From(u.ID()) == graph.Empty {
			return Shortest{from: u}, true
		}
		path = newShortestFrom(u, []graph.Node{u})
	}
	path.dist[path.indexOf[u.ID()]] = 0
	path.negCosts = make(map[negEdge]float64)

	var weight Weighting
	if wg, ok := g.(Weighted); ok {
		weight = wg.Weight
	} else {
		weight = UniformCost(g)
	}

	// Queue to keep track which nodes need to be relaxed.
	// Only nodes whose vertex distance changed in the previous iterations
	// need to be relaxed again.
	queue := newBellmanFordQueue(path.indexOf)
	queue.enqueue(u)

	// TODO(kortschak): Consider adding further optimisations
	// from http://arxiv.org/abs/1111.5414.
	var loops int64
	for queue.len() != 0 {
		u := queue.dequeue()
		uid := u.ID()
		j := path.indexOf[uid]

		to := g.From(uid)
		for to.Next() {
			v := to.Node()
			vid := v.ID()
			k, ok := path.indexOf[vid]
			if !ok {
				k = path.add(v)
			}
			w, ok := weight(uid, vid)
			if !ok {
				panic("bellman-ford: unexpected invalid weight")
			}

			joint := path.dist[j] + w
			if joint < path.dist[k] {
				path.set(k, joint, j)

				if !queue.has(vid) {
					queue.enqueue(v)
				}
			}
		}

		// The maximum number of edges in the relaxed subgraph is |V_r| * (|V_r|-1).
		// If the queue-loop has more iterations than the maximum number of edges
		// it indicates that we have a negative cycle.
		maxEdges := int64(len(path.nodes)) * int64(len(path.nodes)-1)
		if loops > maxEdges {
			path.hasNegativeCycle = true
			return path, false
		}
		loops++
	}

	return path, true
}

// BellmanFordAllFrom returns a shortest-path tree for shortest paths from u to all nodes in
// the graph g, or false indicating that a negative cycle exists in the graph. If the graph
// does not implement Weighted, UniformCost is used.
//
// If g is a graph.Graph, all nodes of the graph will be stored in the shortest-path
// tree, otherwise only nodes reachable from u will be stored.
//
// The time complexity of BellmanFordAllFrom is O(|V|.|E|).
func BellmanFordAllFrom(u graph.Node, g traverse.Graph) (path ShortestAlts, ok bool) {
	if h, ok := g.(graph.Graph); ok {
		if h.Node(u.ID()) == nil {
			return ShortestAlts{from: u}, true
		}
		path = newShortestAltsFrom(u, graph.NodesOf(h.Nodes()))
	} else {
		if g.From(u.ID()) == graph.Empty {
			return ShortestAlts{from: u}, true
		}
		path = newShortestAltsFrom(u, []graph.Node{u})
	}
	path.dist[path.indexOf[u.ID()]] = 0
	path.negCosts = make(map[negEdge]float64)

	var weight Weighting
	if wg, ok := g.(Weighted); ok {
		weight = wg.Weight
	} else {
		weight = UniformCost(g)
	}

	// Queue to keep track which nodes need to be relaxed.
	// Only nodes whose vertex distance changed in the previous iterations
	// need to be relaxed again.
	queue := newBellmanFordQueue(path.indexOf)
	queue.enqueue(u)

	// TODO(kortschak): Consider adding further optimisations
	// from http://arxiv.org/abs/1111.5414.
	var loops int64
	for queue.len() != 0 {
		u := queue.dequeue()
		uid := u.ID()
		j := path.indexOf[uid]

		for _, v := range graph.NodesOf(g.From(uid)) {
			vid := v.ID()
			k, ok := path.indexOf[vid]
			if !ok {
				k = path.add(v)
			}
			w, ok := weight(uid, vid)
			if !ok {
				panic("bellman-ford: unexpected invalid weight")
			}

			joint := path.dist[j] + w
			if joint < path.dist[k] {
				path.set(k, joint, j)

				if !queue.has(vid) {
					queue.enqueue(v)
				}
			} else if joint == path.dist[k] {
				path.addPath(k, j)
			}
		}

		// The maximum number of edges in the relaxed subgraph is |V_r| * (|V_r|-1).
		// If the queue-loop has more iterations than the maximum number of edges
		// it indicates that we have a negative cycle.
		maxEdges := int64(len(path.nodes)) * int64(len(path.nodes)-1)
		if loops > maxEdges {
			path.hasNegativeCycle = true
			return path, false
		}
		loops++
	}

	return path, true
}

// bellmanFordQueue is a queue for the Queue-based Bellman-Ford algorithm.
type bellmanFordQueue struct {
	// queue holds the nodes which need to be relaxed.
	queue linear.NodeQueue

	// onQueue keeps track whether a node is on the queue or not.
	onQueue []bool

	// indexOf contains a mapping holding the id of a node with its index in the onQueue array.
	indexOf map[int64]int
}

// enqueue adds a node to the bellmanFordQueue.
func (q *bellmanFordQueue) enqueue(n graph.Node) {
	i, ok := q.indexOf[n.ID()]
	switch {
	case !ok:
		panic("bellman-ford: unknown node")
	case i < len(q.onQueue):
		if q.onQueue[i] {
			panic("bellman-ford: already queued")
		}
	case i == len(q.onQueue):
		q.onQueue = append(q.onQueue, false)
	case i < cap(q.onQueue):
		q.onQueue = q.onQueue[:i+1]
	default:
		q.onQueue = append(q.onQueue, make([]bool, i-len(q.onQueue)+1)...)
	}
	q.onQueue[i] = true
	q.queue.Enqueue(n)
}

// dequeue returns the first value of the bellmanFordQueue.
func (q *bellmanFordQueue) dequeue() graph.Node {
	n := q.queue.Dequeue()
	q.onQueue[q.indexOf[n.ID()]] = false
	return n
}

// len returns the number of nodes in the bellmanFordQueue.
func (q *bellmanFordQueue) len() int { return q.queue.Len() }

// has returns whether a node with the given id is in the queue.
func (q bellmanFordQueue) has(id int64) bool {
	idx, ok := q.indexOf[id]
	if !ok || idx >= len(q.onQueue) {
		return false
	}
	return q.onQueue[idx]
}

// newBellmanFordQueue creates a new bellmanFordQueue.
func newBellmanFordQueue(indexOf map[int64]int) bellmanFordQueue {
	return bellmanFordQueue{
		onQueue: make([]bool, len(indexOf)),
		indexOf: indexOf,
	}
}
