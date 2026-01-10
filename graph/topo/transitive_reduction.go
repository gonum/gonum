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
		succ := succIDs(g.From(uid))

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

func succIDs(it graph.Nodes) []int64 {
	var ids []int64
	for it.Next() {
		ids = append(ids, it.Node().ID())
	}
	return ids
}
