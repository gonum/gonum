// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"container/heap"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/traverse"
)

// DijkstraFrom returns a shortest-path tree for a shortest path from u to all nodes in
// the graph g. If the graph does not implement Weighted, UniformCost is used.
// DijkstraFrom will panic if g has a u-reachable negative edge weight.
//
// If g is a graph.Graph, all nodes of the graph will be stored in the shortest-path
// tree, otherwise only nodes reachable from u will be stored.
//
// The time complexity of DijkstraFrom is O(|E|.log|V|).
func DijkstraFrom(u graph.Node, g traverse.Graph) Shortest {
	return dijkstraFrom(u, nil, g)
}

// DijkstraFromTo returns a shortest path from u to t in the graph g. The
// result is equivalent to DijkstraFrom(u, g).To(t.ID()), but DijkstraFromTo
// can be more efficient, as it can terminate early if t is reached. If the
// graph does not implement Weighted, UniformCost is used. DijkstraFromTo will
// panic if g has a u-reachable negative edge weight that is discovered before
// reaching t.
//
// The time complexity of DijkstraFromTo is O(|E|.log|V|).
func DijkstraFromTo(u, t graph.Node, g traverse.Graph) (path []graph.Node, weight float64) {
	// An implementation of Bidirectional Dijkstra's algorithm could be even
	// more efficient, but it requires a transposed (or undirected) graph.
	if t == nil {
		panic("dijkstra: nil target node")
	}
	return dijkstraFrom(u, t, g).To(t.ID())
}

func dijkstraFrom(u, t graph.Node, g traverse.Graph) Shortest {
	var path Shortest
	// Use the incremental version when a target is provided.
	if h, ok := g.(graph.Graph); t == nil && ok {
		if h.Node(u.ID()) == nil {
			return Shortest{from: u}
		}
		path = newShortestFrom(u, graph.NodesOf(h.Nodes()))
	} else {
		if g.From(u.ID()) == graph.Empty {
			return Shortest{from: u}
		}
		path = newShortestFrom(u, []graph.Node{u})
	}

	var weight Weighting
	if wg, ok := g.(Weighted); ok {
		weight = wg.Weight
	} else {
		weight = UniformCost(g)
	}

	// Dijkstra's algorithm here is implemented essentially as
	// described in Function B.2 in figure 6 of UTCS Technical
	// Report TR-07-54.
	//
	// This implementation deviates from the report as follows:
	// - the value of path.dist for the start vertex u is initialized to 0;
	// - outdated elements from the priority queue (i.e. with respect to the dist value)
	//   are skipped.
	//
	// http://www.cs.utexas.edu/ftp/techreports/tr07-54.pdf
	Q := priorityQueue{{node: u, dist: 0}}
	for Q.Len() != 0 {
		mid := heap.Pop(&Q).(distanceNode)
		k := path.indexOf[mid.node.ID()]
		if mid.dist > path.dist[k] {
			continue
		}
		mnid := mid.node.ID()
		if t != nil && mnid == t.ID() {
			break
		}
		to := g.From(mnid)
		for to.Next() {
			v := to.Node()
			vid := v.ID()
			j, ok := path.indexOf[vid]
			if !ok {
				j = path.add(v)
			}
			w, ok := weight(mnid, vid)
			if !ok {
				panic("dijkstra: unexpected invalid weight")
			}
			if w < 0 {
				panic("dijkstra: negative edge weight")
			}
			joint := path.dist[k] + w
			if joint < path.dist[j] {
				heap.Push(&Q, distanceNode{node: v, dist: joint})
				path.set(j, joint, k)
			}
		}
	}

	return path
}

// DijkstraAllFrom returns a shortest-path tree for shortest paths from u to all nodes in
// the graph g. If the graph does not implement Weighted, UniformCost is used.
// DijkstraAllFrom will panic if g has a u-reachable negative edge weight.
//
// If g is a graph.Graph, all nodes of the graph will be stored in the shortest-path
// tree, otherwise only nodes reachable from u will be stored.
//
// The time complexity of DijkstraAllFrom is O(|E|.log|V|).
func DijkstraAllFrom(u graph.Node, g traverse.Graph) ShortestAlts {
	var path ShortestAlts
	if h, ok := g.(graph.Graph); ok {
		if h.Node(u.ID()) == nil {
			return ShortestAlts{from: u}
		}
		path = newShortestAltsFrom(u, graph.NodesOf(h.Nodes()))
	} else {
		if g.From(u.ID()) == graph.Empty {
			return ShortestAlts{from: u}
		}
		path = newShortestAltsFrom(u, []graph.Node{u})
	}

	var weight Weighting
	if wg, ok := g.(Weighted); ok {
		weight = wg.Weight
	} else {
		weight = UniformCost(g)
	}

	// Dijkstra's algorithm here is implemented essentially as
	// described in Function B.2 in figure 6 of UTCS Technical
	// Report TR-07-54.
	//
	// This implementation deviates from the report as follows:
	// - the value of path.dist for the start vertex u is initialized to 0;
	// - outdated elements from the priority queue (i.e. with respect to the dist value)
	//   are skipped.
	//
	// http://www.cs.utexas.edu/ftp/techreports/tr07-54.pdf
	Q := priorityQueue{{node: u, dist: 0}}
	for Q.Len() != 0 {
		mid := heap.Pop(&Q).(distanceNode)
		k := path.indexOf[mid.node.ID()]
		if mid.dist > path.dist[k] {
			continue
		}
		mnid := mid.node.ID()
		for _, v := range graph.NodesOf(g.From(mnid)) {
			vid := v.ID()
			j, ok := path.indexOf[vid]
			if !ok {
				j = path.add(v)
			}
			w, ok := weight(mnid, vid)
			if !ok {
				panic("dijkstra: unexpected invalid weight")
			}
			if w < 0 {
				panic("dijkstra: negative edge weight")
			}
			joint := path.dist[k] + w
			if joint < path.dist[j] {
				heap.Push(&Q, distanceNode{node: v, dist: joint})
				path.set(j, joint, k)
			} else if joint == path.dist[j] {
				path.addPath(j, k)
			}
		}
	}

	return path
}

// DijkstraAllPaths returns a shortest-path tree for shortest paths in the graph g.
// If the graph does not implement graph.Weighter, UniformCost is used.
// DijkstraAllPaths will panic if g has a negative edge weight.
//
// The time complexity of DijkstraAllPaths is O(|V|.|E|+|V|^2.log|V|).
func DijkstraAllPaths(g graph.Graph) (paths AllShortest) {
	paths = newAllShortest(graph.NodesOf(g.Nodes()), false)
	dijkstraAllPaths(g, paths)
	return paths
}

// dijkstraAllPaths is the all-paths implementation of Dijkstra. It is shared
// between DijkstraAllPaths and JohnsonAllPaths to avoid repeated allocation
// of the nodes slice and the indexOf map. It returns nothing, but stores the
// result of the work in the paths parameter which is a reference type.
func dijkstraAllPaths(g graph.Graph, paths AllShortest) {
	var weight Weighting
	if wg, ok := g.(graph.Weighted); ok {
		weight = wg.Weight
	} else {
		weight = UniformCost(g)
	}

	var Q priorityQueue
	for i, u := range paths.nodes {
		// Dijkstra's algorithm here is implemented essentially as
		// described in Function B.2 in figure 6 of UTCS Technical
		// Report TR-07-54 with the addition of handling multiple
		// co-equal paths.
		//
		// http://www.cs.utexas.edu/ftp/techreports/tr07-54.pdf

		// Q must be empty at this point.
		heap.Push(&Q, distanceNode{node: u, dist: 0})
		for Q.Len() != 0 {
			mid := heap.Pop(&Q).(distanceNode)
			k := paths.indexOf[mid.node.ID()]
			if mid.dist < paths.dist.At(i, k) {
				paths.dist.Set(i, k, mid.dist)
			}
			mnid := mid.node.ID()
			to := g.From(mnid)
			for to.Next() {
				v := to.Node()
				vid := v.ID()
				j := paths.indexOf[vid]
				w, ok := weight(mnid, vid)
				if !ok {
					panic("dijkstra: unexpected invalid weight")
				}
				if w < 0 {
					panic("dijkstra: negative edge weight")
				}
				joint := paths.dist.At(i, k) + w
				if joint < paths.dist.At(i, j) {
					heap.Push(&Q, distanceNode{node: v, dist: joint})
					paths.set(i, j, joint, k)
				} else if joint == paths.dist.At(i, j) {
					paths.add(i, j, k)
				}
			}
		}
	}
}

type distanceNode struct {
	node graph.Node
	dist float64
}

// priorityQueue implements a no-dec priority queue.
type priorityQueue []distanceNode

func (q priorityQueue) Len() int            { return len(q) }
func (q priorityQueue) Less(i, j int) bool  { return q[i].dist < q[j].dist }
func (q priorityQueue) Swap(i, j int)       { q[i], q[j] = q[j], q[i] }
func (q *priorityQueue) Push(n interface{}) { *q = append(*q, n.(distanceNode)) }
func (q *priorityQueue) Pop() interface{} {
	t := *q
	var n interface{}
	n, *q = t[len(t)-1], t[:len(t)-1]
	return n
}
