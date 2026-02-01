// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topo

import (
	"errors"

	"gonum.org/v1/gonum/graph"
)

type GraphReducer interface {
	graph.Directed
	graph.EdgeRemover
}

// TransitiveReduce removes redundant edges from g while preserving reachability.
// g must be a DAG; otherwise TransitiveReduce returns an error.
func TransitiveReduce(g GraphReducer) error {
	if _, err := Sort(g); err != nil {
		return errors.New("topo: transitive reduction requires a DAG")
	}
	// Map node IDs to dense indices.
	ids, id2idx := indexNodes(g)
	n := len(ids)
	if n == 0 {
		return nil
	}

	// Generation counters avoid clearing DFS state.
	seen := make([]uint32, n)
	visited := make([]uint32, n)
	var seenGen, visitedGen uint32

	// Reusable buffers.
	dfsStack := make([]int64, 0, 64)
	reached := make([]int64, 0, 64)

	uit := g.Nodes()
	for uit.Next() {
		uid := uit.Node().ID()
		_, ok := id2idx[uid]
		if !ok {
			continue
		}
		// Snapshot successors of uid (as IDs).
		successors := successorIDs(g.From(uid))

		seenGen++

		for _, vid := range successors {
			// If vid already covered via another successor, uid->vid is redundant.
			if vIdx, ok := id2idx[vid]; ok && seen[vIdx] == seenGen {
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
			vIdx, ok := id2idx[vid]
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
					nIdx, ok := id2idx[nid]
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
				if xIdx, ok := id2idx[xid]; ok {
					seen[xIdx] = seenGen
				}
			}

			// Mark vid itself as covered (reachable directly from uid).
			seen[vIdx] = seenGen
		}
	}
	return nil
}

func successorIDs(it graph.Nodes) []int64 {
	var ids []int64
	for it.Next() {
		ids = append(ids, it.Node().ID())
	}
	return ids
}

func indexNodes(g graph.Graph) (ids []int64, id2idx map[int64]int) {
	it := g.Nodes()
	ids = make([]int64, 0, it.Len())
	id2idx = make(map[int64]int, it.Len())
	for it.Next() {
		id := it.Node().ID()
		id2idx[id] = len(ids)
		ids = append(ids, id)
	}
	return ids, id2idx
}
