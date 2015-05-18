// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"container/heap"
	"errors"
	"sort"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/internal"
	"github.com/gonum/graph/traverse"
)

// Returns an ordered list consisting of the nodes between start and goal. The path will be the
// shortest path assuming the function heuristicCost is admissible. The second return value is the
// cost, and the third is the number of nodes expanded while searching (useful info for tuning
// heuristics). Negative Costs will cause bad things to happen, as well as negative heuristic
// estimates.
//
// A heuristic is admissible if, for any node in the graph, the heuristic estimate of the cost
// between the node and the goal is less than or set to the true cost.
//
// Performance may be improved by providing a consistent heuristic (though one is not needed to
// find the optimal path), a heuristic is consistent if its value for a given node is less than
// (or equal to) the actual cost of reaching its neighbors + the heuristic estimate for the
// neighbor itself. You can force consistency by making your HeuristicCost function return
// max(NonConsistentHeuristicCost(neighbor,goal), NonConsistentHeuristicCost(self,goal) -
// Cost(self,neighbor)). If there are multiple neighbors, take the max of all of them.
//
// Cost and HeuristicCost take precedence for evaluating cost/heuristic distance. If one is not
// present (i.e. nil) the function will check the graph's interface for the respective interface:
// Coster for Cost and HeuristicCoster for HeuristicCost. If the correct one is present, it will
// use the graph's function for evaluation.
//
// Finally, if neither the argument nor the interface is present, the function will assume
// UniformCost for Cost and NullHeuristic for HeuristicCost.
//
// To run Uniform Cost Search, run A* with the NullHeuristic.
//
// To run Breadth First Search, run A* with both the NullHeuristic and UniformCost (or any cost
// function that returns a uniform positive value.)
func AStar(start, goal graph.Node, g graph.Graph, cost graph.CostFunc, heuristicCost graph.HeuristicCostFunc) (path []graph.Node, pathCost float64, nodesExpanded int) {
	sf := setupFuncs(g, cost, heuristicCost)
	successors, cost, heuristicCost, edgeTo := sf.successors, sf.cost, sf.heuristicCost, sf.edgeTo

	closedSet := make(map[int]internalNode)
	openSet := &aStarPriorityQueue{nodes: make([]internalNode, 0), indexList: make(map[int]int)}
	heap.Init(openSet)
	node := internalNode{start, 0, heuristicCost(start, goal)}
	heap.Push(openSet, node)
	predecessor := make(map[int]graph.Node)

	for openSet.Len() != 0 {
		curr := heap.Pop(openSet).(internalNode)

		nodesExpanded += 1

		if curr.ID() == goal.ID() {
			return rebuildPath(predecessor, goal), curr.gscore, nodesExpanded
		}

		closedSet[curr.ID()] = curr

		for _, neighbor := range successors(curr.Node) {
			if _, ok := closedSet[neighbor.ID()]; ok {
				continue
			}

			g := curr.gscore + cost(edgeTo(curr.Node, neighbor))

			if existing, exists := openSet.Find(neighbor.ID()); !exists {
				predecessor[neighbor.ID()] = curr
				node = internalNode{neighbor, g, g + heuristicCost(neighbor, goal)}
				heap.Push(openSet, node)
			} else if g < existing.gscore {
				predecessor[neighbor.ID()] = curr
				openSet.Fix(neighbor.ID(), g, g+heuristicCost(neighbor, goal))
			}
		}
	}

	return nil, 0, nodesExpanded
}

// BreadthFirstSearch finds a path with a minimal number of edges from from start to goal.
//
// BreadthFirstSearch returns the path found and the number of nodes visited in the search.
// The returned path is nil if no path exists.
func BreadthFirstSearch(start, goal graph.Node, g graph.Graph) ([]graph.Node, int) {
	path, _, visited := AStar(start, goal, g, UniformCost, NullHeuristic)
	return path, visited
}

// Dijkstra's Algorithm is essentially a goalless Uniform Cost Search. That is, its results are
// roughly equivalent to running A* with the Null Heuristic from a single node to every other node
// in the graph -- though it's a fair bit faster because running A* in that way will recompute
// things it's already computed every call. Note that you won't necessarily get the same path
// you would get for A*, but the cost is guaranteed to be the same (that is, if multiple shortest
// paths exist, you may get a different shortest path).
//
// Like A*, Dijkstra's Algorithm likely won't run correctly with negative edge weights -- use
// Bellman-Ford for that instead.
//
// Dijkstra's algorithm usually only returns a cost map, however, since the data is available
// this version will also reconstruct the path to every node.
func Dijkstra(source graph.Node, g graph.Graph, cost graph.CostFunc) (paths map[int][]graph.Node, costs map[int]float64) {

	sf := setupFuncs(g, cost, nil)
	successors, cost, edgeTo := sf.successors, sf.cost, sf.edgeTo

	nodes := g.NodeList()
	openSet := &aStarPriorityQueue{nodes: make([]internalNode, 0), indexList: make(map[int]int)}
	costs = make(map[int]float64, len(nodes)) // May overallocate, will change if it becomes a problem
	predecessor := make(map[int]graph.Node, len(nodes))
	nodeIDMap := make(map[int]graph.Node, len(nodes))
	heap.Init(openSet)

	costs[source.ID()] = 0
	heap.Push(openSet, internalNode{source, 0, 0})

	for openSet.Len() != 0 {
		node := heap.Pop(openSet).(internalNode)

		nodeIDMap[node.ID()] = node

		for _, neighbor := range successors(node) {
			tmpCost := costs[node.ID()] + cost(edgeTo(node, neighbor))
			if cost, ok := costs[neighbor.ID()]; !ok {
				costs[neighbor.ID()] = tmpCost
				predecessor[neighbor.ID()] = node
				heap.Push(openSet, internalNode{neighbor, tmpCost, tmpCost})
			} else if tmpCost < cost {
				costs[neighbor.ID()] = tmpCost
				predecessor[neighbor.ID()] = node
				openSet.Fix(neighbor.ID(), tmpCost, tmpCost)
			}
		}
	}

	paths = make(map[int][]graph.Node, len(costs))
	for node := range costs { // Only reconstruct the path if one exists
		paths[node] = rebuildPath(predecessor, nodeIDMap[node])
	}
	return paths, costs
}

// The Bellman-Ford Algorithm is the same as Dijkstra's Algorithm with a key difference. They both
// take a single source and find the shortest path to every other (reachable) node in the graph.
// Bellman-Ford, however, will detect negative edge loops and abort if one is present. A negative
// edge loop occurs when there is a cycle in the graph such that it can take an edge with a
// negative cost over and over. A -(-2)> B -(2)> C isn't a loop because A->B can only be taken once,
// but A<-(-2)->B-(2)>C is one because A and B have a bi-directional edge, and algorithms like
// Dijkstra's will infinitely flail between them getting progressively lower costs.
//
// That said, if you do not have a negative edge weight, use Dijkstra's Algorithm instead, because
// it's faster.
//
// Like Dijkstra's, along with the costs this implementation will also construct all the paths for
// you. In addition, it has a third return value which will be true if the algorithm was aborted
// due to the presence of a negative edge weight cycle.
func BellmanFord(source graph.Node, g graph.Graph, cost graph.CostFunc) (paths map[int][]graph.Node, costs map[int]float64, err error) {
	sf := setupFuncs(g, cost, nil)
	successors, cost, edgeTo := sf.successors, sf.cost, sf.edgeTo

	predecessor := make(map[int]graph.Node)
	costs = make(map[int]float64)
	nodeIDMap := make(map[int]graph.Node)
	nodeIDMap[source.ID()] = source
	costs[source.ID()] = 0
	nodes := g.NodeList()

	for i := 1; i < len(nodes)-1; i++ {
		for _, node := range nodes {
			nodeIDMap[node.ID()] = node
			succs := successors(node)
			for _, succ := range succs {
				weight := cost(edgeTo(node, succ))
				nodeIDMap[succ.ID()] = succ

				if dist := costs[node.ID()] + weight; dist < costs[succ.ID()] {
					costs[succ.ID()] = dist
					predecessor[succ.ID()] = node
				}
			}

		}
	}

	for _, node := range nodes {
		for _, succ := range successors(node) {
			weight := cost(edgeTo(node, succ))
			if costs[node.ID()]+weight < costs[succ.ID()] {
				return nil, nil, errors.New("Negative edge cycle detected")
			}
		}
	}

	paths = make(map[int][]graph.Node, len(costs))
	for node := range costs {
		paths[node] = rebuildPath(predecessor, nodeIDMap[node])
	}
	return paths, costs, nil
}

// Johnson's Algorithm generates the lowest cost path between every pair of nodes in the graph.
//
// It makes use of Bellman-Ford and a dummy graph. It creates a dummy node containing edges with a
// cost of zero to every other node. Then it runs Bellman-Ford with this dummy node as the source.
// It then modifies the all the nodes' edge weights (which gets rid of all negative weights).
//
// Finally, it removes the dummy node and runs Dijkstra's starting at every node.
//
// This algorithm is fairly slow. Its purpose is to remove negative edge weights to allow
// Dijkstra's to function properly. It's probably not worth it to run this algorithm if you have
// all non-negative edge weights. Also note that this implementation copies your whole graph into
// a GonumGraph (so it can add/remove the dummy node and edges and reweight the graph).
//
// Its return values are, in order: a map from the source node, to the destination node, to the
// path between them; a map from the source node, to the destination node, to the cost of the path
// between them; and a bool that is true if Bellman-Ford detected a negative edge weight cycle --
// thus causing it (and this algorithm) to abort (if aborted is true, both maps will be nil).
func Johnson(g graph.Graph, cost graph.CostFunc) (nodePaths map[int]map[int][]graph.Node, nodeCosts map[int]map[int]float64, err error) {
	sf := setupFuncs(g, cost, nil)
	successors, cost, edgeTo := sf.successors, sf.cost, sf.edgeTo

	/* Copy graph into a mutable one since it has to be altered for this algorithm */
	dummyGraph := concrete.NewDirectedGraph()
	for _, node := range g.NodeList() {
		neighbors := successors(node)
		dummyGraph.NodeExists(node)
		dummyGraph.AddNode(node)
		for _, neighbor := range neighbors {
			e := edgeTo(node, neighbor)
			c := cost(e)
			// Make a new edge with head and tail swapped;
			// works due to the fact that we're not returning
			// any edges in this so the contract doesn't need
			// to be fulfilled.
			if e.Head().ID() != node.ID() {
				e = concrete.Edge{e.Tail(), e.Head()}
			}

			dummyGraph.AddDirectedEdge(e, c)
		}
	}

	/* Step 1: Dummy node with 0 cost edge weights to every other node*/
	dummyNode := dummyGraph.NewNode()
	dummyGraph.AddNode(dummyNode)
	for _, node := range g.NodeList() {
		dummyGraph.AddDirectedEdge(concrete.Edge{dummyNode, node}, 0)
	}

	/* Step 2: Run Bellman-Ford starting at the dummy node, abort if it detects a cycle */
	_, costs, err := BellmanFord(dummyNode, dummyGraph, nil)
	if err != nil {
		return nil, nil, err
	}

	/* Step 3: reweight the graph and remove the dummy node */
	for _, node := range g.NodeList() {
		for _, succ := range successors(node) {
			e := edgeTo(node, succ)
			dummyGraph.AddDirectedEdge(e, cost(e)+costs[node.ID()]-costs[succ.ID()])
		}
	}

	dummyGraph.RemoveNode(dummyNode)

	/* Step 4: Run Dijkstra's starting at every node */
	nodePaths = make(map[int]map[int][]graph.Node, len(g.NodeList()))
	nodeCosts = make(map[int]map[int]float64)

	for _, node := range g.NodeList() {
		nodePaths[node.ID()], nodeCosts[node.ID()] = Dijkstra(node, dummyGraph, nil)
	}

	return nodePaths, nodeCosts, nil
}

// Expands the first node it sees trying to find the destination. Depth First Search is *not*
// guaranteed to find the shortest path, however, if a path exists DFS is guaranteed to find it
// (provided you don't find a way to implement a Graph with an infinite depth.)
func DepthFirstSearch(start, goal graph.Node, g graph.Graph) []graph.Node {
	sf := setupFuncs(g, nil, nil)
	successors := sf.successors

	closedSet := make(internal.IntSet)
	predecessor := make(map[int]graph.Node)

	openSet := internal.NodeStack{start}
	for openSet.Len() != 0 {
		curr := openSet.Pop()

		if closedSet.Has(curr.ID()) {
			continue
		}

		if curr.ID() == goal.ID() {
			return rebuildPath(predecessor, goal)
		}

		closedSet.Add(curr.ID())

		for _, neighbor := range successors(curr) {
			if closedSet.Has(neighbor.ID()) {
				continue
			}

			predecessor[neighbor.ID()] = curr
			openSet.Push(neighbor)
		}
	}

	return nil
}

// An admissible, consistent heuristic that won't speed up computation time at all.
func NullHeuristic(_, _ graph.Node) float64 {
	return 0
}

// Assumes all edges in the graph have the same weight (including edges that don't exist!)
func UniformCost(e graph.Edge) float64 {
	if e == nil {
		return inf
	}

	return 1
}

/* Simple operations */

// Copies a graph into the destination; maintaining all node IDs. The destination
// need not be empty, though overlapping node IDs and conflicting edges will overwrite
// existing data.
func CopyUndirectedGraph(dst graph.MutableGraph, src graph.Graph) {
	cost := setupFuncs(src, nil, nil).cost

	for _, node := range src.NodeList() {
		succs := src.Neighbors(node)
		dst.AddNode(node)
		for _, succ := range succs {
			edge := src.EdgeBetween(node, succ)
			dst.AddUndirectedEdge(edge, cost(edge))
		}
	}

}

// Copies a graph into the destination; maintaining all node IDs. The destination
// need not be empty, though overlapping node IDs and conflicting edges will overwrite
// existing data.
func CopyDirectedGraph(dst graph.MutableDirectedGraph, src graph.DirectedGraph) {
	cost := setupFuncs(src, nil, nil).cost

	for _, node := range src.NodeList() {
		succs := src.Successors(node)
		dst.AddNode(node)
		for _, succ := range succs {
			edge := src.EdgeTo(node, succ)
			dst.AddDirectedEdge(edge, cost(edge))
		}
	}

}

/* Basic Graph tests */

// TarjanSCC returns the strongly connected components of the graph g using Tarjan's algorithm.
//
// A strongly connected component of a graph is a set of vertices where it's possible to reach any
// vertex in the set from any other (meaning there's a cycle between them.)
//
// Generally speaking, a directed graph where the number of strongly connected components is equal
// to the number of nodes is acyclic, unless you count reflexive edges as a cycle (which requires
// only a little extra testing.)
//
func TarjanSCC(g graph.DirectedGraph) [][]graph.Node {
	nodes := g.NodeList()
	t := tarjan{
		succ: g.Successors,

		indexTable: make(map[int]int, len(nodes)),
		lowLink:    make(map[int]int, len(nodes)),
		onStack:    make(internal.IntSet, len(nodes)),
	}
	for _, v := range nodes {
		if t.indexTable[v.ID()] == 0 {
			t.strongconnect(v)
		}
	}
	return t.sccs
}

// tarjan implements Tarjan's strongly connected component finding
// algorithm. The implementation is from the pseudocode at
//
// http://en.wikipedia.org/wiki/Tarjan%27s_strongly_connected_components_algorithm?oldid=642744644
//
type tarjan struct {
	succ func(graph.Node) []graph.Node

	index      int
	indexTable map[int]int
	lowLink    map[int]int
	onStack    internal.IntSet

	stack []graph.Node

	sccs [][]graph.Node
}

// strongconnect is the strongconnect function described in the
// wikipedia article.
func (t *tarjan) strongconnect(v graph.Node) {
	vID := v.ID()

	// Set the depth index for v to the smallest unused index.
	t.index++
	t.indexTable[vID] = t.index
	t.lowLink[vID] = t.index
	t.stack = append(t.stack, v)
	t.onStack.Add(vID)

	// Consider successors of v.
	for _, w := range t.succ(v) {
		wID := w.ID()
		if t.indexTable[wID] == 0 {
			// Successor w has not yet been visited; recur on it.
			t.strongconnect(w)
			t.lowLink[vID] = min(t.lowLink[vID], t.lowLink[wID])
		} else if t.onStack.Has(wID) {
			// Successor w is in stack s and hence in the current SCC.
			t.lowLink[vID] = min(t.lowLink[vID], t.indexTable[wID])
		}
	}

	// If v is a root node, pop the stack and generate an SCC.
	if t.lowLink[vID] == t.indexTable[vID] {
		// Start a new strongly connected component.
		var (
			scc []graph.Node
			w   graph.Node
		)
		for {
			w, t.stack = t.stack[len(t.stack)-1], t.stack[:len(t.stack)-1]
			t.onStack.Remove(w.ID())
			// Add w to current strongly connected component.
			scc = append(scc, w)
			if w.ID() == vID {
				break
			}
		}
		// Output the current strongly connected component.
		t.sccs = append(t.sccs, scc)
	}
}

// IsPath returns true for a connected path within a graph.
//
// IsPath returns true if, starting at path[0] and ending at path[len(path)-1], all nodes between
// are valid neighbors. That is, for each element path[i], path[i+1] is a valid successor.
//
// As special cases, IsPath returns true for a nil or zero length path, and for a path of length 1
// (only one node) but only if the node listed in path exists within the graph.
//
// Graph must be non-nil.
func IsPath(path []graph.Node, g graph.Graph) bool {
	isSuccessor := setupFuncs(g, nil, nil).isSuccessor

	if path == nil || len(path) == 0 {
		return true
	} else if len(path) == 1 {
		return g.NodeExists(path[0])
	}

	for i := 0; i < len(path)-1; i++ {
		if !isSuccessor(path[i], path[i+1]) {
			return false
		}
	}

	return true
}

/* Implements minimum-spanning tree algorithms;
puts the resulting minimum spanning tree in the dst graph */

// Generates a minimum spanning tree with sets.
//
// As with other algorithms that use Cost, the order of precedence is
// Argument > Interface > UniformCost.
//
// The destination must be empty (or at least disjoint with the node IDs of the input)
func Prim(dst graph.MutableGraph, g graph.EdgeListGraph, cost graph.CostFunc) {
	sf := setupFuncs(g, cost, nil)
	cost = sf.cost

	nlist := g.NodeList()

	if nlist == nil || len(nlist) == 0 {
		return
	}

	dst.AddNode(nlist[0])
	remainingNodes := make(internal.IntSet)
	for _, node := range nlist[1:] {
		remainingNodes.Add(node.ID())
	}

	edgeList := g.EdgeList()
	for remainingNodes.Count() != 0 {
		var edges []concrete.WeightedEdge
		for _, edge := range edgeList {
			if (dst.NodeExists(edge.Head()) && remainingNodes.Has(edge.Tail().ID())) ||
				(dst.NodeExists(edge.Tail()) && remainingNodes.Has(edge.Head().ID())) {

				edges = append(edges, concrete.WeightedEdge{Edge: edge, Cost: cost(edge)})
			}
		}

		sort.Sort(byWeight(edges))
		myEdge := edges[0]

		dst.AddUndirectedEdge(myEdge.Edge, myEdge.Cost)
		remainingNodes.Remove(myEdge.Edge.Head().ID())
	}

}

// Generates a minimum spanning tree for a graph using discrete.DisjointSet.
//
// As with other algorithms with Cost, the precedence goes Argument > Interface > UniformCost.
//
// The destination must be empty (or at least disjoint with the node IDs of the input)
func Kruskal(dst graph.MutableGraph, g graph.EdgeListGraph, cost graph.CostFunc) {
	cost = setupFuncs(g, cost, nil).cost

	edgeList := g.EdgeList()
	edges := make([]concrete.WeightedEdge, 0, len(edgeList))
	for _, edge := range edgeList {
		edges = append(edges, concrete.WeightedEdge{Edge: edge, Cost: cost(edge)})
	}

	sort.Sort(byWeight(edges))

	ds := newDisjointSet()
	for _, node := range g.NodeList() {
		ds.makeSet(node.ID())
	}

	for _, edge := range edges {
		// The disjoint set doesn't really care for which is head and which is tail so this
		// should work fine without checking both ways
		if s1, s2 := ds.find(edge.Edge.Head().ID()), ds.find(edge.Edge.Tail().ID()); s1 != s2 {
			ds.union(s1, s2)
			dst.AddUndirectedEdge(edge.Edge, edge.Cost)
		}
	}
}

/* Control flow graph stuff */

// A dominates B if and only if the only path through B travels through A.
//
// This returns all possible dominators for all nodes, it does not prune for strict dominators,
// immediate dominators etc.
//
func Dominators(start graph.Node, g graph.Graph) map[int]Set {
	allNodes := make(Set)
	nlist := g.NodeList()
	dominators := make(map[int]Set, len(nlist))
	for _, node := range nlist {
		allNodes.add(node)
	}

	predecessors := setupFuncs(g, nil, nil).predecessors

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
			preds := predecessors(node)
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

// A Postdominates B if and only if all paths from B travel through A.
//
// This returns all possible post-dominators for all nodes, it does not prune for strict
// postdominators, immediate postdominators etc.
func PostDominators(end graph.Node, g graph.Graph) map[int]Set {
	successors := setupFuncs(g, nil, nil).successors

	allNodes := make(Set)
	nlist := g.NodeList()
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
			succs := successors(node)
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

// VertexOrdering returns the vertex ordering and the k-cores of
// the undirected graph g.
func VertexOrdering(g graph.Graph) (order []graph.Node, cores [][]graph.Node) {
	nodes := g.NodeList()

	// The algorithm used here is essentially as described at
	// http://en.wikipedia.org/w/index.php?title=Degeneracy_%28graph_theory%29&oldid=640308710

	// Initialize an output list L.
	var l []graph.Node

	// Compute a number d_v for each vertex v in G,
	// the number of neighbors of v that are not already in L.
	// Initially, these numbers are just the degrees of the vertices.
	dv := make(map[int]int, len(nodes))
	var (
		maxDegree  int
		neighbours = make(map[int][]graph.Node)
	)
	for _, n := range nodes {
		adj := g.Neighbors(n)
		neighbours[n.ID()] = adj
		dv[n.ID()] = len(adj)
		if len(adj) > maxDegree {
			maxDegree = len(adj)
		}
	}

	// Initialize an array D such that D[i] contains a list of the
	// vertices v that are not already in L for which d_v = i.
	d := make([][]graph.Node, maxDegree+1)
	for _, n := range nodes {
		deg := dv[n.ID()]
		d[deg] = append(d[deg], n)
	}

	// Initialize k to 0.
	k := 0
	// Repeat n times:
	s := []int{0}
	for _ = range nodes { // TODO(kortschak): Remove blank assignment when go1.3.3 is no longer supported.
		// Scan the array cells D[0], D[1], ... until
		// finding an i for which D[i] is nonempty.
		var (
			i  int
			di []graph.Node
		)
		for i, di = range d {
			if len(di) != 0 {
				break
			}
		}

		// Set k to max(k,i).
		if i > k {
			k = i
			s = append(s, make([]int, k-len(s)+1)...)
		}

		// Select a vertex v from D[i]. Add v to the
		// beginning of L and remove it from D[i].
		var v graph.Node
		v, d[i] = di[len(di)-1], di[:len(di)-1]
		l = append(l, v)
		s[k]++
		delete(dv, v.ID())

		// For each neighbor w of v not already in L,
		// subtract one from d_w and move w to the
		// cell of D corresponding to the new value of d_w.
		for _, w := range neighbours[v.ID()] {
			dw, ok := dv[w.ID()]
			if !ok {
				continue
			}
			for i, n := range d[dw] {
				if n.ID() == w.ID() {
					d[dw][i], d[dw] = d[dw][len(d[dw])-1], d[dw][:len(d[dw])-1]
					dw--
					d[dw] = append(d[dw], w)
					break
				}
			}
			dv[w.ID()] = dw
		}
	}

	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
	cores = make([][]graph.Node, len(s))
	offset := len(l)
	for i, n := range s {
		cores[i] = l[offset-n : offset]
		offset -= n
	}
	return l, cores
}

// BronKerbosch returns the set of maximal cliques of the undirected graph g.
func BronKerbosch(g graph.Graph) [][]graph.Node {
	nodes := g.NodeList()

	// The algorithm used here is essentially BronKerbosch3 as described at
	// http://en.wikipedia.org/w/index.php?title=Bron%E2%80%93Kerbosch_algorithm&oldid=656805858

	p := make(Set, len(nodes))
	for _, n := range nodes {
		p.add(n)
	}
	x := make(Set)
	var bk bronKerbosch
	order, _ := VertexOrdering(g)
	for _, v := range order {
		neighbours := g.Neighbors(v)
		nv := make(Set, len(neighbours))
		for _, n := range neighbours {
			nv.add(n)
		}
		bk.maximalCliquePivot(g, []graph.Node{v}, make(Set).intersect(p, nv), make(Set).intersect(x, nv))
		p.remove(v)
		x.add(v)
	}
	return bk
}

type bronKerbosch [][]graph.Node

func (bk *bronKerbosch) maximalCliquePivot(g graph.Graph, r []graph.Node, p, x Set) {
	if len(p) == 0 && len(x) == 0 {
		*bk = append(*bk, r)
		return
	}

	neighbours := bk.choosePivotFrom(g, p, x)
	nu := make(Set, len(neighbours))
	for _, n := range neighbours {
		nu.add(n)
	}
	for _, v := range p {
		if nu.has(v) {
			continue
		}
		neighbours := g.Neighbors(v)
		nv := make(Set, len(neighbours))
		for _, n := range neighbours {
			nv.add(n)
		}

		var found bool
		for _, n := range r {
			if n.ID() == v.ID() {
				found = true
				break
			}
		}
		var sr []graph.Node
		if !found {
			sr = append(r[:len(r):len(r)], v)
		}

		bk.maximalCliquePivot(g, sr, make(Set).intersect(p, nv), make(Set).intersect(x, nv))
		p.remove(v)
		x.add(v)
	}
}

func (*bronKerbosch) choosePivotFrom(g graph.Graph, p, x Set) (neighbors []graph.Node) {
	// TODO(kortschak): Investigate the impact of pivot choice that maximises
	// |p ⋂ neighbours(u)| as a function of input size. Until then, leave as
	// compile time option.
	if !tomitaTanakaTakahashi {
		for _, n := range p {
			return g.Neighbors(n)
		}
		for _, n := range x {
			return g.Neighbors(n)
		}
		panic("bronKerbosch: empty set")
	}

	var (
		max   = -1
		pivot graph.Node
	)
	maxNeighbors := func(s Set) {
	outer:
		for _, u := range s {
			nb := g.Neighbors(u)
			c := len(nb)
			if c <= max {
				continue
			}
			for n := range nb {
				if _, ok := p[n]; ok {
					continue
				}
				c--
				if c <= max {
					continue outer
				}
			}
			max = c
			pivot = u
			neighbors = nb
		}
	}
	maxNeighbors(p)
	maxNeighbors(x)
	if pivot == nil {
		panic("bronKerbosch: empty set")
	}
	return neighbors
}

// ConnectedComponents returns the connected components of the graph g. All
// edges are treated as undirected.
func ConnectedComponents(g graph.Graph) [][]graph.Node {
	var (
		w  traverse.DepthFirst
		c  []graph.Node
		cc [][]graph.Node
	)
	during := func(n graph.Node) {
		c = append(c, n)
	}
	after := func() {
		cc = append(cc, []graph.Node(nil))
		cc[len(cc)-1] = append(cc[len(cc)-1], c...)
		c = c[:0]
	}
	w.WalkAll(g, nil, after, during)

	return cc
}
