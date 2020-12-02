// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

// djSet implements a disjoint set finder using the union-find algorithm.
type djSet map[int64]*dsNode

// add adds e to the collection of sets held by the disjoint set.
func (s djSet) add(e int64) {
	if _, ok := s[e]; ok {
		return
	}
	s[e] = &dsNode{}
}

// union joins two sets a and b within the collection of sets held by
// the disjoint set.
func (djSet) union(a, b *dsNode) {
	ra := find(a)
	rb := find(b)
	if ra == rb {
		return
	}
	if ra.rank < rb.rank {
		ra.parent = rb
		return
	}
	rb.parent = ra
	if ra.rank == rb.rank {
		ra.rank++
	}
}

// find returns the root of the set containing e.
func (s djSet) find(e int64) *dsNode {
	n, ok := s[e]
	if !ok {
		return nil
	}
	return find(n)
}

// find returns the root of the set containing the set node, n.
func find(n *dsNode) *dsNode {
	for ; n.parent != nil; n = n.parent {
	}
	return n
}

// dsNode is a disjoint set element.
type dsNode struct {
	parent *dsNode
	rank   int
}
