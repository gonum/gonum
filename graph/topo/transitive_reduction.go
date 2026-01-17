package topo

import (
	"errors"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/traverse"
)

type GraphReducer interface {
	graph.Directed
	graph.EdgeRemover
}

// TransitiveReduce removes redundant edges from g while preserving reachability.
// The input graph must be a DAG. Behavior is undefined if g contains cycles.
func TransitiveReduce(g GraphReducer) error {
	if _, err := Sort(g); err != nil {
		return errors.New("topo: transitive reduction requires a DAG")
	}
	uit := g.Nodes()
	for uit.Next() {
		uid := uit.Node().ID()

		// Snapshot successors of u to avoid iterator invalidation
		// while removing edges out of uid.
		succ := successorIDs(g.From(uid))

		for _, vid := range succ {
			v := g.Node(vid)
			if v == nil {
				continue
			}
			// Walk from v; for every reachable x, remove direct edge u->x.
			df := traverse.DepthFirst{
				Traverse: func(e graph.Edge) bool {
					xid := e.To().ID()
					if xid != vid && g.HasEdgeFromTo(uid, xid) {
						g.RemoveEdge(uid, xid)
					}
					return true
				},
			}
			df.Walk(g, v, nil)
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

// TransitiveReduce removes redundant edges from g while preserving reachability.
// g must be a DAG. Behavior is undefined if g contains cycles.
func TransitiveReduce2(g GraphReducer) error {
	if _, err := Sort(g); err != nil {
		return errors.New("topo: transitive reduction requires a DAG")
	}
	// From node IDs build mapping to dense indices
	ids, id2idx := indexNodes(g)
	n := len(ids)
	if n == 0 {
		return nil
	}

	// Use generation/epoch trick for allocation-free, copy-free resetting of DFS
	seen := make([]uint32, n)
	visited := make([]uint32, n)
	var seenGen uint32
	var visitedGen uint32

	// Reusable data-structures
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
			// - report reached nodes even if already in seen
			// - but do not descend into nodes already in seen
			visitedGen++
			dfsStack = dfsStack[:0]
			reached = reached[:0]

			// Start node to check the successors-of-successors of u via DFS
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

					// If already covered by prior successors of uid, prune descendants --> prune DFS
					if seen[nIdx] == seenGen {
						continue
					}
					dfsStack = append(dfsStack, nid)
				}
			}

			// Remove redundant uid->x edges (if existing) that have been already reached via successor v.
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
