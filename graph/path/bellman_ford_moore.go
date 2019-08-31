// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"errors"
	"gonum.org/v1/gonum/graph"
)

// BellmanFordFrom returns a shortest-path tree for a shortest path from u to all nodes in
// the graph g, or false indicating that a negative cycle exists in the graph. If the graph
// does not implement Weighted, UniformCost is used.
//
// The time complexity of BellmanFordFrom is O(|V|.|E|).
func BellmanFordFrom(u graph.Node, g graph.Graph) (path Shortest, ok bool) {
	if g.Node(u.ID()) == nil {
		return Shortest{from: u}, true
	}
	var weight Weighting
	if wg, ok := g.(Weighted); ok {
		weight = wg.Weight
	} else {
		weight = UniformCost(g)
	}

	nodes := graph.NodesOf(g.Nodes())

	path = newShortestFrom(u, nodes)
	path.dist[path.indexOf[u.ID()]] = 0

	// queue to keep track which nodes need to be relaxed
	// only nodes whose vertex distance changed in the previous iterations need to be relaxed again
	queue := newBellmanFordQueue()
	queue.enqueue(u)

	// bool array to keep track whether a node is on the queue or not
	onQueue := make([]bool, len(nodes))
	onQueue[path.indexOf[u.ID()]] = true

	n := len(nodes)
	// the maximum of edges in a graph is |V| * (|V| -1)
	// which is also the worst case complexity
	// íf the queue-loop has more iterations than the amount of maximum edges
	// it indicates that we have a negative cycle
	maxEdges := n * (n - 1)
	negativeCycle := false

	loops := 0
	for queue.len() != 0 {
		u, err := queue.dequeue()
		if err != nil {
			panic(err)
		}
		uid := u.ID()
		onQueue[path.indexOf[uid]] = false

		for _, v := range graph.NodesOf(g.From(uid)) {
			vid := v.ID()
			k := path.indexOf[vid]
			w, ok := weight(uid, vid)
			if !ok {
				panic("bellman-ford: unexpected invalid weight")
			}

			j := path.indexOf[uid]
			joint := path.dist[j] + w
			if joint < path.dist[k] {
				path.set(k, joint, j)
				index := path.indexOf[vid]

				// check if node is already in the queue
				// we do not want any duplicates in the queue
				if !onQueue[index] {
					onQueue[index] = true
					queue.enqueue(v)
				}
			}
		}

		loops++
		if loops > maxEdges {
			negativeCycle = true
			break
		}
	}

	if negativeCycle {
		for j, u := range nodes {
			uid := u.ID()
			for _, v := range graph.NodesOf(g.From(uid)) {
				vid := v.ID()
				k := path.indexOf[vid]
				w, ok := weight(uid, vid)
				if !ok {
					panic("bellman-ford: unexpected invalid weight")
				}
				if path.dist[j]+w < path.dist[k] {
					path.hasNegativeCycle = true
					return path, false
				}
			}
		}
	}

	return path, true
}

// bellmanFordQueue is a queue for the Queue-based Bellman ford algorithm
type bellmanFordQueue struct {
	nodes []graph.Node
}

// enqueue adds a node to the bellmanFordQueue
func (b *bellmanFordQueue) enqueue(n graph.Node) {
	b.nodes = append(b.nodes, n)
}

// dequeue returns the first value of the bellmanFordQueue
func (b *bellmanFordQueue) dequeue() (graph.Node, error) {
	if len(b.nodes) == 0 {
		return nil, errors.New("queue is empty!")
	}

	u := b.nodes[0]
	b.nodes = b.nodes[1:]
	return u, nil
}

// len returns the amount of nodes in the bellmanFordQueue
func (b *bellmanFordQueue) len() int {
	return len(b.nodes)
}

// newBellmanFordQueue creates a new bellmanFordQueue
func newBellmanFordQueue() bellmanFordQueue {
	return bellmanFordQueue{
		nodes: make([]graph.Node, 0),
	}
}
