// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package traverse provides basic graph traversal primitives.
package traverse

import "github.com/gonum/graph"

// BreadthFirst implements stateful breadth-first graph traversal.
type BreadthFirst struct {
	EdgeFilter func(graph.Edge) bool
	Visit      func(u, v graph.Node)
	queue      nodeQueue
	visited    intSet
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
		b.visited = make(intSet)
	}
	b.queue.enqueue(from)
	b.visited.add(from.ID())

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
			if b.visited.has(n.ID()) {
				continue
			}
			if b.Visit != nil {
				b.Visit(t, n)
			}
			b.visited.add(n.ID())
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
	visited    intSet
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
		d.visited = make(intSet)
	}
	d.stack.push(from)
	d.visited.add(from.ID())

	for d.stack.len() > 0 {
		t := d.stack.pop()
		if until != nil && until(t) {
			return t
		}
		for _, n := range neighbors(t) {
			if d.EdgeFilter != nil && !d.EdgeFilter(g.EdgeBetween(t, n)) {
				continue
			}
			if d.visited.has(n.ID()) {
				continue
			}
			if d.Visit != nil {
				d.Visit(t, n)
			}
			d.visited.add(n.ID())
			d.stack.push(n)
		}
	}

	return nil
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

// A set is a set of integer identifiers.
type intSet map[int]struct{}

// The simple accessor methods for Set are provided to allow ease of
// implementation change should the need arise.

// add inserts an element into the set.
func (s intSet) add(e int) {
	s[e] = struct{}{}
}

// has reports the existence of the element in the set.
func (s intSet) has(e int) bool {
	_, ok := s[e]
	return ok
}

// remove deletes the specified element from the set.
func (s intSet) remove(e int) {
	delete(s, e)
}

// count reports the number of elements stored in the set.
func (s intSet) count() int {
	return len(s)
}
