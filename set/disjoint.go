// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package set

// A disjoint set is a collection of non-overlapping sets. That is, for any two sets in the
// disjoint set, their intersection is the empty set.
//
// A disjoint set has three principle operations: Make Set, Find, and Union.
//
// Make set creates a new set for an element (presuming it does not already exist in any set in
// the disjoint set), Find finds the set containing that element (if any), and Union merges two
// sets in the disjoint set. In general, algorithms operating on disjoint sets are "union-find"
// algorithms, where two sets are found with Find, and then joined with Union.
//
// A concrete example of a union-find algorithm can be found as discrete.Kruskal -- which unions
// two sets when an edge is created between two vertices, and refuses to make an edge between two
// vertices if they're part of the same set.
type DisjointSet struct {
	master map[interface{}]*DisjointSetNode
}

type DisjointSetNode struct {
	parent *DisjointSetNode
	rank   int
}

func NewDisjointSet() *DisjointSet {
	return &DisjointSet{master: make(map[interface{}]*DisjointSetNode)}
}

// If the element isn't already somewhere in there, adds it to the master set and its own tiny set.
func (ds *DisjointSet) MakeSet(el interface{}) {
	if _, ok := ds.master[el]; ok {
		return
	}
	dsNode := &DisjointSetNode{rank: 0}
	dsNode.parent = dsNode
	ds.master[el] = dsNode
}

// Returns the set the element belongs to, or nil if none.
func (ds *DisjointSet) Find(el interface{}) *DisjointSetNode {
	dsNode, ok := ds.master[el]
	if !ok {
		return nil
	}

	return find(dsNode)
}

func find(dsNode *DisjointSetNode) *DisjointSetNode {
	if dsNode.parent != dsNode {
		dsNode.parent = find(dsNode.parent)
	}

	return dsNode.parent
}

// Unions two subsets within the DisjointSet.
//
// If x or y are not in this disjoint set, the behavior is undefined. If either pointer is nil,
// this function will panic.
func (ds *DisjointSet) Union(x, y *DisjointSetNode) {
	if x == nil || y == nil {
		panic("Disjoint Set union on nil sets")
	}
	xRoot := find(x)
	yRoot := find(y)
	if xRoot == nil || yRoot == nil {
		return
	}

	if xRoot == yRoot {
		return
	}

	if xRoot.rank < yRoot.rank {
		xRoot.parent = yRoot
	} else if yRoot.rank < xRoot.rank {
		yRoot.parent = xRoot
	} else {
		yRoot.parent = xRoot
		xRoot.rank += 1
	}
}
