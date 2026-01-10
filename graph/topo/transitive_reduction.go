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

// TransitiveReduce removes redundant edges from g while preserving reachability.
// g must be a DAG. Behavior is undefined if g contains cycles.
func TransitiveReduce2(g GraphReducer) error {
	if _, err := Sort(g); err != nil {
		return errors.New("topo: transitive reduction requires a DAG")
	}
	uit := g.Nodes()
	for uit.Next() {
		uid := uit.Node().ID()

		// Snapshot successors of uid; we will remove edges out of uid.
		succ := succIDs(g.From(uid))

		// Nodes reachable from any already-processed successor of uid.
		seen := make(map[int64]struct{}, len(succ))

		for _, vid := range succ {
			// If vid is already reachable through another successor, edge uid->vid is redundant.
			if _, ok := seen[vid]; ok {
				if g.HasEdgeFromTo(uid, vid) {
					g.RemoveEdge(uid, vid)
				}
			}

			// Walk from vid, pruning subgraphs already covered by prior successors.
			newly := prunedDFS(g, vid, seen)

			for _, xid := range newly {
				if g.HasEdgeFromTo(uid, xid) {
					g.RemoveEdge(uid, xid)
				}
				seen[xid] = struct{}{}
			}

			// Mark vid itself as covered (reachable directly from uid).
			seen[vid] = struct{}{}
		}
	}
	return nil
}

// prunedDFS walks nodes reachable from start. It does not descend into nodes already
// present in seen, but it still reports them as reached.
// It returns all nodes first reached by this walk (excluding start).
func prunedDFS(g graph.Directed, start int64, seen map[int64]struct{}) []int64 {
	var reached []int64

	localVisited := make(map[int64]struct{}, len(seen)+1)
	stack := []int64{start}
	localVisited[start] = struct{}{}

	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		from := g.From(n)
		for from.Next() {
			x := from.Node().ID()

			// If we already visited x in this walk, skip.
			if _, ok := localVisited[x]; ok {
				continue
			}
			localVisited[x] = struct{}{}

			// Report x as reached (even if x is already in seen).
			reached = append(reached, x)

			// If x is already covered by prior successors, prune its descendants.
			if _, ok := seen[x]; ok {
				continue
			}

			// Otherwise, explore further.
			stack = append(stack, x)
		}
	}
	return reached
}
