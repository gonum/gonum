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

// NodeIDGraphBuild is a graph that can create new nodes with
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
func Complete(dst NodeIDGraphBuilder, ids IDer) error {
	if ids.Len() == 1 {
		u, new := dst.NodeWithID(ids.ID(0))
		if new {
			dst.AddNode(u)
		}
		return nil
	}
	for i := 0; i < ids.Len(); i++ {
		uid := ids.ID(i)
		u, _ := dst.NodeWithID(uid)
		for j := i + 1; j < ids.Len(); j++ {
			vid := ids.ID(j)
			if uid == vid {
				return fmt.Errorf("gen: node ID collision i=%d j=%d: %d", i, j, uid)
			}
			v, _ := dst.NodeWithID(vid)
			dst.SetEdge(dst.NewEdge(u, v))
		}
	}
	return nil
}

// Cycle constructs a cycle in dst using the node IDs in cycle.
// If dst is a directed graph, edges are directed from earlier nodes to later
// nodes in cycle.
func Cycle(dst NodeIDGraphBuilder, cycle IDer) error {
	if cycle.Len() == 1 {
		u, new := dst.NodeWithID(cycle.ID(0))
		if new {
			dst.AddNode(u)
		}
		return nil
	}
	for i := 0; i < cycle.Len(); i++ {
		uid := cycle.ID(i)
		vid := cycle.ID((i + 1) % cycle.Len())
		if uid == vid {
			return fmt.Errorf("gen: adjacent node IDs equal at %d: %d", i, uid)
		}
		u, _ := dst.NodeWithID(uid)
		v, _ := dst.NodeWithID(vid)
		dst.SetEdge(dst.NewEdge(u, v))
	}
	return nil
}

// Path constructs a path graph in dst with
// If dst is a directed graph, edges are directed from earlier nodes to later
// nodes in path.
func Path(dst NodeIDGraphBuilder, path IDer) error {
	if path.Len() == 1 {
		u, new := dst.NodeWithID(path.ID(0))
		if new {
			dst.AddNode(u)
		}
		return nil
	}
	for i := 0; i < path.Len()-1; i++ {
		uid := path.ID(i)
		vid := path.ID(i + 1)
		if uid == vid {
			return fmt.Errorf("gen: adjacent node IDs equal at %d: %d", i, uid)
		}
		u, _ := dst.NodeWithID(uid)
		v, _ := dst.NodeWithID(vid)
		dst.SetEdge(dst.NewEdge(u, v))
	}
	return nil
}

// Star constructs a star graph in dst with edges between the center node ID to
// node with IDs specified in leaves.
// If dst is a directed graph, edges are directed from the center node to the
// leaves.
func Star(dst NodeIDGraphBuilder, center int64, leaves IDer) error {
	c, new := dst.NodeWithID(center)
	if new {
		dst.AddNode(c)
	}
	for i := 0; i < leaves.Len(); i++ {
		id := leaves.ID(i)
		if id == center {
			return fmt.Errorf("gen: leaf %d ID matches central node ID: %d", i, center)
		}
		n, _ := dst.NodeWithID(id)
		dst.SetEdge(dst.NewEdge(c, n))
	}
	return nil
}

// Wheel constructs a wheel graph in dst with edges from the center
// node ID to node with IDs specified in cycle and between nodes with IDs
// adjacent in the cycle.
// If dst is a directed graph, edges are directed from the center node to the
// cycle and from earlier nodes to later nodes in cycle.
func Wheel(dst NodeIDGraphBuilder, center int64, cycle IDer) error {
	err := Star(dst, center, cycle)
	if err != nil {
		return err
	}
	return Cycle(dst, cycle)
}
