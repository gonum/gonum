// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topo

import (
	"slices"

	"gonum.org/v1/gonum/graph"
)

type DirectedGraphRemover interface {
	graph.Directed
	graph.EdgeRemover
}

// TransitiveReduce removes redundant edges from g while preserving reachability.
// g must be a DAG; otherwise if g contains a cycle, TransitiveReduce panics.
func TransitiveReduce(g DirectedGraphRemover) {
	if slices.ContainsFunc(TarjanSCC(g), func(scc []graph.Node) bool { return len(scc) != 1 }) {
		panic("topo: TransitiveReduce: graph is not directed acyclic")
	}

	// Map node IDs to dense indices.
	ids, indexOf := indexNodes(g)
	n := len(ids)
	if n == 0 {
		return
	}

	// The seen/visited slices use "generation counters" to avoid O(n) clearing between
	// DFS runs: entries equal to the current generation are considered set, and we
	// increment the generation to logically reset the slice.
	seen := make([]uint64, n)
	visited := make([]uint64, n)
	var seenGen, visitedGen uint64

	// Reusable buffers.
	dfsStack := make([]int64, 0, 64)
	reached := make([]int64, 0, 64)

	uit := g.Nodes()
	for uit.Next() {
		uid := uit.Node().ID()
		_, ok := indexOf[uid]
		if !ok {
			continue
		}
		// Snapshot successors of uid (as IDs).
		successors := idsFrom(g.From(uid))

		seenGen++

		for _, vid := range successors {
			// If vid already covered via another successor, uid->vid is redundant.
			if vIdx, ok := indexOf[vid]; ok && seen[vIdx] == seenGen {
				if g.HasEdgeFromTo(uid, vid) {
					g.RemoveEdge(uid, vid)
				}
			}
			// Pruned DFS from vid:
			// - record reached nodes (even if already seen)
			// - but don't descend into nodes already seen
			visitedGen++
			dfsStack = dfsStack[:0]
			reached = reached[:0]

			// DFS starting at vid to find nodes reachable via vid.
			vIdx, ok := indexOf[vid]
			if !ok {
				continue
			}
			visited[vIdx] = visitedGen
			dfsStack = append(dfsStack, vid)

			for len(dfsStack) > 0 {
				cid := dfsStack[len(dfsStack)-1]
				dfsStack = dfsStack[:len(dfsStack)-1]

				it := g.From(cid)
				for it.Next() {
					nid := it.Node().ID()
					nIdx, ok := indexOf[nid]
					if !ok {
						continue
					}
					if visited[nIdx] == visitedGen {
						continue
					}
					visited[nIdx] = visitedGen

					reached = append(reached, nid)

					// If already covered by another successor of uid, prune descendants.
					if seen[nIdx] == seenGen {
						continue
					}
					dfsStack = append(dfsStack, nid)
				}
			}

			// Remove uid->x edges where x is reachable via vid.
			for _, xid := range reached {
				if g.HasEdgeFromTo(uid, xid) {
					g.RemoveEdge(uid, xid)
				}
				if xIdx, ok := indexOf[xid]; ok {
					seen[xIdx] = seenGen
				}
			}

			// Mark vid itself as covered (reachable directly from uid).
			seen[vIdx] = seenGen
		}
	}
}

func idsFrom(it graph.Nodes) []int64 {
	var ids []int64
	if n := it.Len(); n > 0 {
		ids = make([]int64, 0, n)
	}
	for it.Next() {
		ids = append(ids, it.Node().ID())
	}
	return ids
}

// indexNodes returns the node IDs of g and a map from node ID to a dense index
// into the returned ids slice.
func indexNodes(g graph.Graph) (ids []int64, indexOf map[int64]int) {
	it := g.Nodes()

	if n := it.Len(); n > 0 {
		ids = make([]int64, 0, n)
		indexOf = make(map[int64]int, n)
	} else {
		indexOf = make(map[int64]int)
	}

	for it.Next() {
		id := it.Node().ID()
		indexOf[id] = len(ids)
		ids = append(ids, id)
	}
	return ids, indexOf
}
