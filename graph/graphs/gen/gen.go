// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gen

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
)

// GraphBuilder is a graph that can have nodes and edges added.
type GraphBuilder interface {
	HasEdgeBetween(xid, yid int64) bool
	graph.Builder
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// NodeIDGraphBuilder is a graph that can create new nodes with
// specified IDs.
type NodeIDGraphBuilder interface {
	graph.Builder
	graph.NodeWithIDer
}

// IDer is a mapping from an index to a node ID.
type IDer interface {
	// Len returns the length of the set of node IDs.
	Len() int

	// ID returns the ID of the indexed node.
	// ID must be a bijective function. No check
	// is made for this property.
	ID(int) int64
}

// IDRange is an IDer that provides a set of IDs in [First, Last].
type IDRange struct{ First, Last int64 }

func (r IDRange) Len() int       { return int(r.Last - r.First + 1) }
func (r IDRange) ID(i int) int64 { return r.First + int64(i) }

// IDSet is an IDer providing an explicit set of IDs.
type IDSet []int64

func (s IDSet) Len() int       { return len(s) }
func (s IDSet) ID(i int) int64 { return s[i] }

// Complete constructs a complete graph in dst using nodes with the given IDs.
// If any ID appears twice in ids, Complete will panic.
func Complete(dst NodeIDGraphBuilder, ids IDer) {
	switch ids.Len() {
	case 0:
		return
	case 1:
		u, new := dst.NodeWithID(ids.ID(0))
		if new {
			dst.AddNode(u)
		}
		return
	}
	for i := 0; i < ids.Len(); i++ {
		uid := ids.ID(i)
		u, _ := dst.NodeWithID(uid)
		for j := i + 1; j < ids.Len(); j++ {
			vid := ids.ID(j)
			if uid == vid {
				panic(fmt.Errorf("gen: node ID collision i=%d j=%d: id=%d", i, j, uid))
			}
			v, _ := dst.NodeWithID(vid)
			dst.SetEdge(dst.NewEdge(u, v))
		}
	}
}

// Cycle constructs a cycle in dst using the node IDs in cycle.
// If dst is a directed graph, edges are directed from earlier nodes to later
// nodes in cycle. If any ID appears twice in cycle, Cycle will panic.
func Cycle(dst NodeIDGraphBuilder, cycle IDer) {
	switch cycle.Len() {
	case 0:
		return
	case 1:
		u, new := dst.NodeWithID(cycle.ID(0))
		if new {
			dst.AddNode(u)
		}
		return
	}
	err := check(cycle)
	if err != nil {
		panic(err)
	}
	cycleNoCheck(dst, cycle)
}

func cycleNoCheck(dst NodeIDGraphBuilder, cycle IDer) {
	for i := 0; i < cycle.Len(); i++ {
		uid := cycle.ID(i)
		vid := cycle.ID((i + 1) % cycle.Len())
		u, _ := dst.NodeWithID(uid)
		v, _ := dst.NodeWithID(vid)
		dst.SetEdge(dst.NewEdge(u, v))
	}
}

// Path constructs a path graph in dst with
// If dst is a directed graph, edges are directed from earlier nodes to later
// nodes in path. If any ID appears twice in path, Path will panic.
func Path(dst NodeIDGraphBuilder, path IDer) {
	switch path.Len() {
	case 0:
		return
	case 1:
		u, new := dst.NodeWithID(path.ID(0))
		if new {
			dst.AddNode(u)
		}
		return
	}
	err := check(path)
	if err != nil {
		panic(err)
	}
	for i := 0; i < path.Len()-1; i++ {
		uid := path.ID(i)
		vid := path.ID(i + 1)
		u, _ := dst.NodeWithID(uid)
		v, _ := dst.NodeWithID(vid)
		dst.SetEdge(dst.NewEdge(u, v))
	}
}

// Star constructs a star graph in dst with edges between the center node ID to
// node with IDs specified in leaves.
// If dst is a directed graph, edges are directed from the center node to the
// leaves. If any ID appears twice in leaves and center, Star will panic.
func Star(dst NodeIDGraphBuilder, center int64, leaves IDer) {
	c, new := dst.NodeWithID(center)
	if new {
		dst.AddNode(c)
	}
	if leaves.Len() == 0 {
		return
	}
	err := check(leaves, center)
	if err != nil {
		panic(err)
	}
	for i := 0; i < leaves.Len(); i++ {
		id := leaves.ID(i)
		n, _ := dst.NodeWithID(id)
		dst.SetEdge(dst.NewEdge(c, n))
	}
}

// Wheel constructs a wheel graph in dst with edges from the center
// node ID to node with IDs specified in cycle and between nodes with IDs
// adjacent in the cycle.
// If dst is a directed graph, edges are directed from the center node to the
// cycle and from earlier nodes to later nodes in cycle. If any ID appears
// twice in cycle and center, Wheel will panic.
func Wheel(dst NodeIDGraphBuilder, center int64, cycle IDer) {
	Star(dst, center, cycle)
	if cycle.Len() <= 1 {
		return
	}
	cycleNoCheck(dst, cycle)
}

// Tree constructs an n-ary tree breadth-first in dst with the given fan-out, n.
// If the number of nodes does not equal \sum_{i=0}^h n^i, where h is an integer
// indicating the height of the tree, a partial tree will be constructed with not
// all nodes having zero or n children, and not all leaves being h from the root.
// If the number of nodes is greater than one, n must be non-zero and
// less than the number of nodes, otherwise Tree will panic.
// If dst is a directed graph, edges are directed from the root to the leaves.
// If any ID appears more than once in nodes, Tree will panic.
func Tree(dst NodeIDGraphBuilder, n int, nodes IDer) {
	switch nodes.Len() {
	case 0:
		return
	case 1:
		if u, new := dst.NodeWithID(nodes.ID(0)); new {
			dst.AddNode(u)
		}
		return
	}

	if n < 1 || nodes.Len() <= n {
		panic("gen: invalid fan-out")
	}

	err := check(nodes)
	if err != nil {
		panic(err)
	}

	j := 0
	for i := 0; j < nodes.Len(); i++ {
		u, _ := dst.NodeWithID(nodes.ID(i))
		for j = n*i + 1; j <= n*i+n && j < nodes.Len(); j++ {
			v, _ := dst.NodeWithID(nodes.ID(j))
			dst.SetEdge(dst.NewEdge(u, v))
		}
	}
}

// check confirms that no node ID exists more than once in ids and extra.
func check(ids IDer, extra ...int64) error {
	seen := make(map[int64]int, ids.Len()+len(extra))
	for j := 0; j < ids.Len(); j++ {
		uid := ids.ID(j)
		if prev, exists := seen[uid]; exists {
			return fmt.Errorf("gen: node ID collision i=%d j=%d: id=%d", prev, j, uid)
		}
		seen[uid] = j
	}
	lenIDs := ids.Len()
	for j, uid := range extra {
		if prev, exists := seen[uid]; exists {
			if prev < lenIDs {
				if len(extra) == 1 {
					return fmt.Errorf("gen: node ID collision i=%d with extra: id=%d", prev, uid)
				}
				return fmt.Errorf("gen: node ID collision i=%d with extra j=%d: id=%d", prev, j, uid)
			}
			prev -= lenIDs
			return fmt.Errorf("gen: extra node ID collision i=%d j=%d: id=%d", prev, j, uid)
		}
		seen[uid] = j + lenIDs
	}
	return nil
}
