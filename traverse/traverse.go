// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package traverse provides basic graph traversal primitives.
package traverse

import (
	"github.com/gonum/graph"
	"github.com/gonum/graph/internal"
)

// BreadthFirst implements stateful breadth-first graph traversal.
type BreadthFirst struct {
	EdgeFilter func(graph.Edge) bool
	Visit      func(u, v graph.Node)
	queue      nodeQueue
	visited    internal.IntSet
}

// Walk performs a breadth-first traversal of the graph g starting from the given node,
// depending on the the EdgeFilter field and the until parameter if they are non-nil. The
// traversal follows edges for which EdgeFilter(edge) is true and returns the first node
// for which until(node, depth) is true. During the traversal, if the Visit field is
// non-nil, it is called with the nodes joined by each followed edge.
func (b *BreadthFirst) Walk(g graph.Graph, from graph.Node, until func(n graph.Node, d int) bool) graph.Node {
	var neighbors func(graph.Node) []graph.Node
	switch g := g.(type) {
	case graph.DirectedGraph:
		neighbors = g.Successors
	default:
		neighbors = g.Neighbors
	}

	if b.visited == nil {
		b.visited = make(internal.IntSet)
	}
	b.queue.enqueue(from)
	b.visited.Add(from.ID())

	var (
		depth     int
		children  int
		untilNext = 1
	)
	for b.queue.len() > 0 {
		t := b.queue.dequeue()
		if until != nil && until(t, depth) {
			return t
		}
		for _, n := range neighbors(t) {
			if b.EdgeFilter != nil && !b.EdgeFilter(g.EdgeBetween(t, n)) {
				continue
			}
			if b.visited.Has(n.ID()) {
				continue
			}
			if b.Visit != nil {
				b.Visit(t, n)
			}
			b.visited.Add(n.ID())
			children++
			b.queue.enqueue(n)
		}
		if untilNext--; untilNext == 0 {
			depth++
			untilNext = children
			children = 0
		}
	}

	return nil
}

// WalkAll calls Walk for each unvisited node of the graph g using edges independent
// of their direction. The functions before and after are called prior to commencing
// and after completing each walk if they are non-nil respectively. The function
// during is called on each node as it is traversed.
func (b *BreadthFirst) WalkAll(g graph.Graph, before, after func(), during func(graph.Node)) {
	// Ensure that when we pass a directed graph
	// we use neighbors and not successors.
	g = struct{ graph.Graph }{g}

	b.Reset()
	for _, from := range g.NodeList() {
		if b.Visited(from) {
			continue
		}
		if before != nil {
			before()
		}
		b.Walk(g, from, func(n graph.Node, _ int) bool {
			if during != nil {
				during(n)
			}
			return false
		})
		if after != nil {
			after()
		}
	}
}

// Visited returned whether the node n was visited during a traverse.
func (b *BreadthFirst) Visited(n graph.Node) bool {
	_, ok := b.visited[n.ID()]
	return ok
}

// Reset resets the state of the traverser for reuse.
func (b *BreadthFirst) Reset() {
	b.queue.head = 0
	b.queue.data = b.queue.data[:0]
	b.visited = nil
}

// DepthFirst implements stateful depth-first graph traversal.
type DepthFirst struct {
	EdgeFilter func(graph.Edge) bool
	Visit      func(u, v graph.Node)
	stack      nodeStack
	visited    internal.IntSet
}

// Walk performs a depth-first traversal of the graph g starting from the given node,
// depending on the the EdgeFilter field and the until parameter if they are non-nil. The
// traversal follows edges for which EdgeFilter(edge) is true and returns the first node
// for which until(node) is true. During the traversal, if the Visit field is non-nil, it
// is called with the nodes joined by each followed edge.
func (d *DepthFirst) Walk(g graph.Graph, from graph.Node, until func(graph.Node) bool) graph.Node {
	var neighbors func(graph.Node) []graph.Node
	switch g := g.(type) {
	case graph.DirectedGraph:
		neighbors = g.Successors
	default:
		neighbors = g.Neighbors
	}

	if d.visited == nil {
		d.visited = make(internal.IntSet)
	}
	d.stack.push(from)
	d.visited.Add(from.ID())

	for d.stack.len() > 0 {
		t := d.stack.pop()
		if until != nil && until(t) {
			return t
		}
		for _, n := range neighbors(t) {
			if d.EdgeFilter != nil && !d.EdgeFilter(g.EdgeBetween(t, n)) {
				continue
			}
			if d.visited.Has(n.ID()) {
				continue
			}
			if d.Visit != nil {
				d.Visit(t, n)
			}
			d.visited.Add(n.ID())
			d.stack.push(n)
		}
	}

	return nil
}

// WalkAll calls Walk for each unvisited node of the graph g using edges independent
// of their direction. The functions before and after are called prior to commencing
// and after completing each walk if they are non-nil respectively. The function
// during is called on each node as it is traversed.
func (d *DepthFirst) WalkAll(g graph.Graph, before, after func(), during func(graph.Node)) {
	// Ensure that when we pass a directed graph
	// we use neighbors and not successors.
	g = struct{ graph.Graph }{g}

	d.Reset()
	for _, from := range g.NodeList() {
		if d.Visited(from) {
			continue
		}
		if before != nil {
			before()
		}
		d.Walk(g, from, func(n graph.Node) bool {
			if during != nil {
				during(n)
			}
			return false
		})
		if after != nil {
			after()
		}
	}
}

// Visited returned whether the node n was visited during a traverse.
func (d *DepthFirst) Visited(n graph.Node) bool {
	_, ok := d.visited[n.ID()]
	return ok
}

// Reset resets the state of the traverser for reuse.
func (d *DepthFirst) Reset() {
	d.stack = d.stack[:0]
	d.visited = nil
}

// nodeStack implements a LIFO stack.
type nodeStack []graph.Node

func (s *nodeStack) len() int { return len(*s) }
func (s *nodeStack) pop() graph.Node {
	v := *s
	v, n := v[:len(v)-1], v[len(v)-1]
	*s = v
	return n
}
func (s *nodeStack) push(n graph.Node) { *s = append(*s, n) }

// nodeQueue implements a FIFO queue.
type nodeQueue struct {
	head int
	data []graph.Node
}

func (q *nodeQueue) len() int { return len(q.data) - q.head }
func (q *nodeQueue) enqueue(n graph.Node) {
	if len(q.data) == cap(q.data) && q.head > 0 {
		l := q.len()
		copy(q.data, q.data[q.head:])
		q.head = 0
		q.data = append(q.data[:l], n)
	} else {
		q.data = append(q.data, n)
	}
}
func (q *nodeQueue) dequeue() graph.Node {
	if q.len() == 0 {
		panic("queue: empty queue")
	}

	var n graph.Node
	n, q.data[q.head] = q.data[q.head], nil
	q.head++

	if q.len() == 0 {
		q.head = 0
		q.data = q.data[:0]
	}

	return n
}
