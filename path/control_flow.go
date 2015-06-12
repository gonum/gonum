// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import "github.com/gonum/graph"

// PostDominatores returns all dominators for all nodes in g. It does not
// prune for strict post-dominators, immediate dominators etc.
//
// A dominates B if and only if the only path through B travels through A.
func Dominators(start graph.Node, g graph.Graph) map[int]Set {
	allNodes := make(Set)
	nlist := g.Nodes()
	dominators := make(map[int]Set, len(nlist))
	for _, node := range nlist {
		allNodes.add(node)
	}

	var to func(graph.Node) []graph.Node
	switch g := g.(type) {
	case graph.Directed:
		to = g.To
	default:
		to = g.From
	}

	for _, node := range nlist {
		dominators[node.ID()] = make(Set)
		if node.ID() == start.ID() {
			dominators[node.ID()].add(start)
		} else {
			dominators[node.ID()].copy(allNodes)
		}
	}

	for somethingChanged := true; somethingChanged; {
		somethingChanged = false
		for _, node := range nlist {
			if node.ID() == start.ID() {
				continue
			}
			preds := to(node)
			if len(preds) == 0 {
				continue
			}
			tmp := make(Set).copy(dominators[preds[0].ID()])
			for _, pred := range preds[1:] {
				tmp.intersect(tmp, dominators[pred.ID()])
			}

			dom := make(Set)
			dom.add(node)

			dom.union(dom, tmp)
			if !equal(dom, dominators[node.ID()]) {
				dominators[node.ID()] = dom
				somethingChanged = true
			}
		}
	}

	return dominators
}

// PostDominatores returns all post-dominators for all nodes in g. It does not
// prune for strict post-dominators, immediate post-dominators etc.
//
// A post-dominates B if and only if all paths from B travel through A.
func PostDominators(end graph.Node, g graph.Graph) map[int]Set {
	allNodes := make(Set)
	nlist := g.Nodes()
	dominators := make(map[int]Set, len(nlist))
	for _, node := range nlist {
		allNodes.add(node)
	}

	for _, node := range nlist {
		dominators[node.ID()] = make(Set)
		if node.ID() == end.ID() {
			dominators[node.ID()].add(end)
		} else {
			dominators[node.ID()].copy(allNodes)
		}
	}

	for somethingChanged := true; somethingChanged; {
		somethingChanged = false
		for _, node := range nlist {
			if node.ID() == end.ID() {
				continue
			}
			succs := g.From(node)
			if len(succs) == 0 {
				continue
			}
			tmp := make(Set).copy(dominators[succs[0].ID()])
			for _, succ := range succs[1:] {
				tmp.intersect(tmp, dominators[succ.ID()])
			}

			dom := make(Set)
			dom.add(node)

			dom.union(dom, tmp)
			if !equal(dom, dominators[node.ID()]) {
				dominators[node.ID()] = dom
				somethingChanged = true
			}
		}
	}

	return dominators
}
